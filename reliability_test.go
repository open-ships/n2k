package n2k

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/brutella/can"
	"github.com/open-ships/n2k/internal/framer"
	"github.com/open-ships/n2k/pgn"
	"github.com/open-ships/n2k/raw"
	"github.com/stretchr/testify/require"
)

// A controllable physical write verifies admission ownership and cancellation
// without relying on scheduler timing or a real device.
type reliabilityBus struct {
	*mockBus
	entered    chan struct{}
	release    chan struct{}
	once       sync.Once
	closedOnce sync.Once
}

func (b *reliabilityBus) WriteFrameContext(ctx context.Context, frame can.Frame) error {
	if framer.ParseCANID(frame.ID).PGN == 127250 {
		b.once.Do(func() { close(b.entered) })
		select {
		case <-b.release:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return b.WriteFrame(frame)
}

func (b *reliabilityBus) Close() error {
	b.closedOnce.Do(func() { close(b.release) })
	return b.mockBus.Close()
}

func newReliabilityClient(t *testing.T) (*Client, *reliabilityBus) {
	t.Helper()
	b := &reliabilityBus{mockBus: newMockBus(), entered: make(chan struct{}), release: make(chan struct{})}
	c, err := NewClient(context.Background(), WithBus(b), WithClaimTimeout(20*time.Millisecond), WithHeartbeatInterval(0))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, c.Close()) })
	return c, b
}

func TestReliabilityWriteSnapshotsAtAdmission(t *testing.T) {
	c, bus := newReliabilityClient(t)
	first := c.Write(&pgn.VesselHeading{Heading: u64(100)})
	select {
	case <-bus.entered:
	case <-time.After(time.Second):
		t.Fatal("first write did not reach the bus")
	}
	value := uint64(12345)
	priority := uint8(3)
	message := &pgn.VesselHeading{Info: pgn.MessageInfo{Priority: &priority}, Heading: &value}
	second := c.Write(message)
	value, priority = 54321, 7
	message.Heading = nil
	bus.closedOnce.Do(func() { close(bus.release) })
	require.NoError(t, first.Wait())
	require.NoError(t, second.Wait())
	frames := framesWithPGN(bus.getWritten(), 127250)
	require.Len(t, frames, 2)
	require.Equal(t, uint8(3), framer.ParseCANID(frames[1].ID).Priority)
	decoded, err := pgn.DecodeMessage(pgn.MessageInfo{PGN: 127250}, frames[1].Data[:frames[1].Length])
	require.NoError(t, err)
	require.Equal(t, uint64(12345), *decoded.(*pgn.VesselHeading).Heading)
}

func TestReliabilityQueuedWriteDoesNotCrossEpoch(t *testing.T) {
	c, bus := newReliabilityClient(t)
	first := c.Write(&pgn.VesselHeading{Heading: u64(100)})
	<-bus.entered
	queued := c.Write(&pgn.VesselHeading{Heading: u64(200)})
	c.handleConnectionChange(false, 1)
	require.ErrorIs(t, queued.Wait(), ErrEpochChanged)
	require.ErrorIs(t, first.Wait(), ErrEpochChanged)
	require.ErrorIs(t, c.Write(&pgn.VesselHeading{}).Wait(), ErrNotReady)
	bus.closedOnce.Do(func() { close(bus.release) })
	c.handleConnectionChange(true, 2)
	require.Eventually(t, func() bool { return !c.Status().Rejoining }, time.Second, time.Millisecond)
	require.NoError(t, c.Write(&pgn.VesselHeading{Heading: u64(300)}).Wait())
	require.Len(t, framesWithPGN(bus.getWritten(), 127250), 1)
}

func TestReliabilityRequestInvalidatedPromptly(t *testing.T) {
	for _, transition := range []string{"close", "disconnect", "address", "failure"} {
		t.Run(transition, func(t *testing.T) {
			c, _, addr := newCitizenClient(t)
			ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
			defer cancel()
			done := make(chan error, 1)
			go func() { _, err := Request[*pgn.ProductInformation](ctx, c, 42); done <- err }()
			require.Eventually(t, func() bool {
				c.correlator.mu.Lock()
				defer c.correlator.mu.Unlock()
				return len(c.correlator.waiters) == 1
			}, time.Second, time.Millisecond)
			want := ErrEpochChanged
			switch transition {
			case "close":
				want = ErrClientClosed
				require.NoError(t, c.Close())
			case "disconnect":
				c.handleConnectionChange(false, 1)
			case "address":
				c.handleAddressChange(addr - 1)
			case "failure":
				want = errors.New("injected terminal failure")
				c.fail(want)
			}
			select {
			case err := <-done:
				require.ErrorIs(t, err, want)
			case <-time.After(500 * time.Millisecond):
				t.Fatal("request remained alive after its network session ended")
			}
		})
	}
}

func TestReliabilityRequestTableBoundAndDestinationIsolation(t *testing.T) {
	co := newCorrelator()
	for range maxPendingRequests {
		_, err := co.add(126996, 42, 251, 1, 2)
		require.NoError(t, err)
	}
	_, err := co.add(126996, 42, 251, 1, 2)
	require.ErrorIs(t, err, ErrRequestQueueFull)
	wrong := &pgn.ProductInformation{Info: pgn.MessageInfo{PGN: 126996, SourceId: 42, TargetId: pgn.Target(250), ConnectionEpoch: 1, ClaimEpoch: 2}}
	co.observe(wrong)
	require.Empty(t, co.waiters[0].ch)
	wrong.Info.TargetId = pgn.Target(251)
	wrong.Info.ClaimEpoch = 1
	co.observe(wrong)
	require.Empty(t, co.waiters[0].ch)
	wrong.Info.ClaimEpoch = 2
	co.observe(wrong)
	first := (<-co.waiters[0].ch).(*pgn.ProductInformation)
	second := (<-co.waiters[1].ch).(*pgn.ProductInformation)
	first.Info.DecodeIssues = append(first.Info.DecodeIssues, "caller change")
	require.Empty(t, second.Info.DecodeIssues)
	co.invalidate(ErrEpochChanged)
	require.Empty(t, co.waiters)
}

func TestReliabilityFastPacketCannotSpanEpochs(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var delivered []pgn.Message
	pipeline, err := newReadPipeline(ctx, config{}, func(msg pgn.Message) { delivered = append(delivered, msg) })
	require.NoError(t, err)
	payload, err := pgn.EncodeMessage(&pgn.ProductInformation{ModelId: "epoch fixture"})
	require.NoError(t, err)
	frames := framer.FrameFastPacket(framer.BuildCANID(126996, 6, 42, 255), payload, 0)
	send := func(frame can.Frame, epoch uint64) {
		observation := frameObservation(frame, "test", "test", raw.DirectionReceived)
		observation.ConnectionEpoch = epoch
		observation.ClaimEpoch = epoch
		pipeline.HandleObservation(observation)
	}
	send(frames[0], 1)
	pipeline.resetEpoch(2, 2)
	for _, frame := range frames[1:] {
		send(frame, 2)
	}
	require.Empty(t, delivered)
	for _, frame := range frames {
		send(frame, 1)
	}
	require.Empty(t, delivered)
	for _, frame := range frames {
		send(frame, 2)
	}
	require.Len(t, delivered, 1)
}

func TestReliabilityRequiredResponseDuringMaximumBAM(t *testing.T) {
	c, bus, address := newCitizenClient(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	transfer := c.WriteContext(ctx, &messageSnapshot{number: 65000, payload: make([]byte, 1785)})
	require.Eventually(t, func() bool { return len(framesWithPGN(bus.getWritten(), 60160)) > 0 }, time.Second, time.Millisecond)
	started := time.Now()
	bus.inbound <- isoRequestFrame(126996, 42, address)
	require.Eventually(t, func() bool { return len(framesWithPGN(bus.getWritten(), 126996)) >= 20 }, 500*time.Millisecond, time.Millisecond)
	require.Less(t, time.Since(started), 500*time.Millisecond)
	select {
	case <-transfer.Done():
		t.Fatal("maximum BAM ended before required response")
	default:
	}
	cancel()
	err := transfer.Wait()
	require.ErrorIs(t, err, context.Canceled)
	var partial *WriteError
	require.ErrorAs(t, err, &partial)
	require.True(t, partial.TransmissionUncertain)
	require.Positive(t, partial.CompletedRecords)
	require.NoError(t, c.Err())
}

func TestReliabilityReconnectAfterSupersededContentionGate(t *testing.T) {
	c, _, address := newCitizenClient(t)
	c.handleConnectionChange(false, 1)
	c.handleConnectionChange(true, 2)
	// Replace the reconnect gate with a contention gate, then invalidate it
	// before its timer fires. Only the latest connection may become ready.
	c.handleAddressChange(address - 1)
	c.handleConnectionChange(false, 2)
	c.handleConnectionChange(true, 3)
	require.Eventually(t, func() bool { return !c.Status().Rejoining }, 2*time.Second, time.Millisecond)
	require.NoError(t, c.Write(&pgn.VesselHeading{}).Wait())
	require.Equal(t, uint64(3), c.Status().ConnectionEpoch)
}
