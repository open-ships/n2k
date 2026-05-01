package claiming_test

import (
	"encoding/binary"
	"fmt"
	"log/slog"
	"sync"
	"testing"

	"github.com/brutella/can"
	"github.com/open-ships/n2k/internal/claiming"
	"github.com/open-ships/n2k/internal/framer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// frameSink captures CAN frames written by the claimer, safe for concurrent use.
type frameSink struct {
	mu     sync.Mutex
	frames []can.Frame
}

func (s *frameSink) write(f can.Frame) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.frames = append(s.frames, f)
	return nil
}

func (s *frameSink) lastFrame() can.Frame {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.frames[len(s.frames)-1]
}

func (s *frameSink) count() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.frames)
}

func (s *frameSink) all() []can.Frame {
	s.mu.Lock()
	defer s.mu.Unlock()
	dst := make([]can.Frame, len(s.frames))
	copy(dst, s.frames)
	return dst
}

// extractClaimAddress returns the source address from an Address Claim CAN frame.
func extractClaimAddress(f can.Frame) uint8 {
	return uint8(f.ID & 0xFF)
}

// extractClaimName returns the 64-bit NAME from an Address Claim CAN frame payload.
func extractClaimName(f can.Frame) uint64 {
	return binary.LittleEndian.Uint64(f.Data[:])
}

// expectedCANID returns the expected CAN ID for an address claim from the given source.
func expectedCANID(source uint8) uint32 {
	return framer.BuildCANID(60928, 6, source, 255)
}

func discardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(noopWriter{}, nil))
}

type noopWriter struct{}

func (noopWriter) Write(p []byte) (int, error) { return len(p), nil }

func TestClaimer_AutoMode_HappyPath(t *testing.T) {
	sink := &frameSink{}
	ourName := uint64(1000)

	c := claiming.New(claiming.Config{
		Mode:       claiming.ModeAuto,
		Address:    253,
		Name:       ourName,
		WriteFrame: sink.write,
		Logger:     discardLogger(),
	})

	err := c.Start()
	require.NoError(t, err)

	// Should have sent exactly one claim frame.
	require.Equal(t, 1, sink.count())

	f := sink.lastFrame()
	assert.Equal(t, expectedCANID(253), f.ID, "CAN ID should target address 253")
	assert.Equal(t, ourName, extractClaimName(f), "payload should contain our NAME")
	assert.Equal(t, uint8(253), c.Address(), "claimed address should be 253")
}

func TestClaimer_AutoMode_Contention(t *testing.T) {
	sink := &frameSink{}
	ourName := uint64(5000)
	competitorName := uint64(1000) // lower NAME wins

	var newAddr uint8
	c := claiming.New(claiming.Config{
		Mode:    claiming.ModeAuto,
		Address: 253,
		Name:    ourName,
		WriteFrame: sink.write,
		OnAddressChange: func(addr uint8) {
			newAddr = addr
		},
		Logger: discardLogger(),
	})

	err := c.Start()
	require.NoError(t, err)

	// Competitor claims address 253 with a lower NAME — we lose.
	c.HandleAddressClaim(253, competitorName)

	// We should have retried on 252.
	assert.Equal(t, uint8(252), c.Address(), "should have moved to 252")
	assert.Equal(t, uint8(252), newAddr, "OnAddressChange should have been called with 252")

	// Two frames total: initial claim at 253, retry claim at 252.
	require.Equal(t, 2, sink.count())
	assert.Equal(t, uint8(252), extractClaimAddress(sink.lastFrame()))
}

func TestClaimer_AutoMode_MultipleRetries(t *testing.T) {
	sink := &frameSink{}
	ourName := uint64(5000)
	competitorName := uint64(1000)

	var addressHistory []uint8
	c := claiming.New(claiming.Config{
		Mode:    claiming.ModeAuto,
		Address: 253,
		Name:    ourName,
		WriteFrame: sink.write,
		OnAddressChange: func(addr uint8) {
			addressHistory = append(addressHistory, addr)
		},
		Logger: discardLogger(),
	})

	err := c.Start()
	require.NoError(t, err)

	// Lose contention three times in a row.
	c.HandleAddressClaim(253, competitorName)
	assert.Equal(t, uint8(252), c.Address())

	c.HandleAddressClaim(252, competitorName)
	assert.Equal(t, uint8(251), c.Address())

	c.HandleAddressClaim(251, competitorName)
	assert.Equal(t, uint8(250), c.Address())

	assert.Equal(t, []uint8{252, 251, 250}, addressHistory)

	// 1 initial + 3 retries = 4 frames.
	assert.Equal(t, 4, sink.count())
}

func TestClaimer_AutoMode_Exhaustion(t *testing.T) {
	sink := &frameSink{}
	ourName := uint64(5000)
	competitorName := uint64(1000)

	var lastAddr uint8
	c := claiming.New(claiming.Config{
		Mode:    claiming.ModeAuto,
		Address: 3, // Start low to make exhaustion quick.
		Name:    ourName,
		WriteFrame: sink.write,
		OnAddressChange: func(addr uint8) {
			lastAddr = addr
		},
		Logger: discardLogger(),
	})

	err := c.Start()
	require.NoError(t, err)

	// Contest address 3, should try 2.
	c.HandleAddressClaim(3, competitorName)
	assert.Equal(t, uint8(2), c.Address())

	// Contest address 2, should try 1.
	c.HandleAddressClaim(2, competitorName)
	assert.Equal(t, uint8(1), c.Address())

	// Contest address 1, should wrap to 253.
	c.HandleAddressClaim(1, competitorName)
	assert.Equal(t, uint8(253), c.Address())

	// Now contest every remaining address down from 253 to 4.
	for addr := uint8(253); addr >= 4; addr-- {
		c.HandleAddressClaim(addr, competitorName)
	}

	// All 253 addresses exhausted — should fall back to 254.
	assert.Equal(t, uint8(254), c.Address())
	assert.Equal(t, uint8(254), lastAddr)
}

func TestClaimer_ExplicitMode_HappyPath(t *testing.T) {
	sink := &frameSink{}
	ourName := uint64(1000)

	c := claiming.New(claiming.Config{
		Mode:       claiming.ModeExplicit,
		Address:    42,
		Name:       ourName,
		WriteFrame: sink.write,
		Logger:     discardLogger(),
	})

	err := c.Start()
	require.NoError(t, err)

	require.Equal(t, 1, sink.count())
	assert.Equal(t, expectedCANID(42), sink.lastFrame().ID)
	assert.Equal(t, uint8(42), c.Address())
}

func TestClaimer_ExplicitMode_Contention(t *testing.T) {
	sink := &frameSink{}
	ourName := uint64(5000)
	competitorName := uint64(1000) // lower NAME wins

	var fatalErr error
	c := claiming.New(claiming.Config{
		Mode:    claiming.ModeExplicit,
		Address: 42,
		Name:    ourName,
		WriteFrame: sink.write,
		OnFatalError: func(err error) {
			fatalErr = err
		},
		Logger: discardLogger(),
	})

	err := c.Start()
	require.NoError(t, err)

	c.HandleAddressClaim(42, competitorName)

	require.NotNil(t, fatalErr, "OnFatalError should have been called")
	assert.Contains(t, fatalErr.Error(), "explicit mode")
	assert.Contains(t, fatalErr.Error(), "42")

	// Address should still be 42 — explicit mode does not change the address.
	assert.Equal(t, uint8(42), c.Address())

	// Only the initial claim frame should have been sent (no retry in explicit mode).
	assert.Equal(t, 1, sink.count())
}

func TestClaimer_DefendAddress(t *testing.T) {
	sink := &frameSink{}
	ourName := uint64(1000)
	competitorName := uint64(5000) // higher NAME — we win

	c := claiming.New(claiming.Config{
		Mode:       claiming.ModeAuto,
		Address:    100,
		Name:       ourName,
		WriteFrame: sink.write,
		Logger:     discardLogger(),
	})

	err := c.Start()
	require.NoError(t, err)
	require.Equal(t, 1, sink.count())

	// Competitor claims our address but has a higher NAME — we defend.
	c.HandleAddressClaim(100, competitorName)

	// Should have sent a second claim (defensive re-broadcast).
	require.Equal(t, 2, sink.count())
	assert.Equal(t, expectedCANID(100), sink.lastFrame().ID, "defensive claim should be from our address")
	assert.Equal(t, ourName, extractClaimName(sink.lastFrame()))

	// Address stays the same.
	assert.Equal(t, uint8(100), c.Address())
}

func TestClaimer_HandleISORequest(t *testing.T) {
	sink := &frameSink{}
	ourName := uint64(42424242)

	c := claiming.New(claiming.Config{
		Mode:       claiming.ModeAuto,
		Address:    200,
		Name:       ourName,
		WriteFrame: sink.write,
		Logger:     discardLogger(),
	})

	err := c.Start()
	require.NoError(t, err)
	require.Equal(t, 1, sink.count())

	// Receive ISO Request for PGN 60928 — should respond with our claim.
	c.HandleISORequest()

	require.Equal(t, 2, sink.count())
	f := sink.lastFrame()
	assert.Equal(t, expectedCANID(200), f.ID)
	assert.Equal(t, ourName, extractClaimName(f))
}

func TestClaimer_IgnoresDifferentAddress(t *testing.T) {
	sink := &frameSink{}
	ourName := uint64(5000)

	c := claiming.New(claiming.Config{
		Mode:       claiming.ModeAuto,
		Address:    100,
		Name:       ourName,
		WriteFrame: sink.write,
		Logger:     discardLogger(),
	})

	err := c.Start()
	require.NoError(t, err)

	// Another device claims address 50 — not our problem.
	c.HandleAddressClaim(50, uint64(1000))

	// No additional frames sent.
	assert.Equal(t, 1, sink.count())
	assert.Equal(t, uint8(100), c.Address())
}

func TestClaimer_IgnoresOwnName(t *testing.T) {
	sink := &frameSink{}
	ourName := uint64(5000)

	c := claiming.New(claiming.Config{
		Mode:       claiming.ModeAuto,
		Address:    100,
		Name:       ourName,
		WriteFrame: sink.write,
		Logger:     discardLogger(),
	})

	err := c.Start()
	require.NoError(t, err)

	// Receive a claim with our own NAME on our address (echo). Should be ignored.
	c.HandleAddressClaim(100, ourName)

	assert.Equal(t, 1, sink.count())
	assert.Equal(t, uint8(100), c.Address())
}

func TestClaimer_WriteFrameError(t *testing.T) {
	writeErr := fmt.Errorf("bus offline")
	c := claiming.New(claiming.Config{
		Mode:    claiming.ModeAuto,
		Address: 100,
		Name:    uint64(1000),
		WriteFrame: func(f can.Frame) error {
			return writeErr
		},
		Logger: discardLogger(),
	})

	err := c.Start()
	assert.ErrorIs(t, err, writeErr, "Start should propagate WriteFrame errors")
}

func TestClaimer_AutoMode_WrapAround(t *testing.T) {
	sink := &frameSink{}
	ourName := uint64(5000)
	competitorName := uint64(1000)

	c := claiming.New(claiming.Config{
		Mode:       claiming.ModeAuto,
		Address:    1, // start at lowest valid address
		Name:       ourName,
		WriteFrame: sink.write,
		Logger:     discardLogger(),
	})

	err := c.Start()
	require.NoError(t, err)

	// Lose contention at address 1 — should wrap to 253.
	c.HandleAddressClaim(1, competitorName)
	assert.Equal(t, uint8(253), c.Address(), "should wrap from 1 to 253")
}

func TestClaimer_PayloadMatchesLittleEndianName(t *testing.T) {
	sink := &frameSink{}
	ourName := uint64(0xDEADBEEFCAFEBABE)

	c := claiming.New(claiming.Config{
		Mode:       claiming.ModeAuto,
		Address:    100,
		Name:       ourName,
		WriteFrame: sink.write,
		Logger:     discardLogger(),
	})

	err := c.Start()
	require.NoError(t, err)

	f := sink.lastFrame()
	var expected [8]byte
	binary.LittleEndian.PutUint64(expected[:], ourName)
	assert.Equal(t, expected, f.Data, "payload should be NAME in little-endian byte order")
}
