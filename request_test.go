package n2k

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/open-ships/n2k/internal/framer"
	"github.com/open-ships/n2k/pgn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// respondWhenRequested watches the mock bus for ISO requests for pgnNum and
// injects the given message from responder in reply to each one. It answers
// every matching request, not just the first: concurrent callers each send
// their own ISO request (Request registers its waiter before transmitting),
// so replying per-request guarantees each waiter gets a response regardless of
// the order in which the goroutines are scheduled. A scan cursor makes sure a
// given request frame is answered exactly once, since getWritten returns the
// whole append-only log on every poll.
func respondWhenRequested(t *testing.T, mb *mockBus, pgnNum uint32, responder uint8, msg pgn.Message) {
	t.Helper()
	go func() {
		scanned := 0
		deadline := time.Now().Add(5 * time.Second)
		for time.Now().Before(deadline) {
			written := mb.getWritten()
			for ; scanned < len(written); scanned++ {
				f := written[scanned]
				id := framer.ParseCANID(f.ID)
				if id.PGN != 59904 || f.Length < 3 {
					continue
				}
				requested := uint32(f.Data[0]) | uint32(f.Data[1])<<8 | uint32(f.Data[2])<<16
				if requested != pgnNum {
					continue
				}
				payload, err := pgn.EncodeMessage(msg)
				if err != nil {
					return
				}
				canID := framer.BuildCANID(pgnNum, 6, responder, 255)
				if len(payload) <= 8 {
					mb.inbound <- framer.FrameSingle(canID, payload)
				} else {
					for _, rf := range framer.FrameFastPacket(canID, payload, 0) {
						mb.inbound <- rf
					}
				}
			}
			time.Sleep(5 * time.Millisecond)
		}
	}()
}

func TestRequest_ProductInformation(t *testing.T) {
	c, mb, _ := newCitizenClient(t)

	respondWhenRequested(t, mb, 126996, 0x23, &pgn.ProductInformation{
		ModelId:     "wind sensor",
		ProductCode: u64(9),
	})

	pi, err := Request[*pgn.ProductInformation](context.Background(), c, 0x23)
	require.NoError(t, err)
	assert.Equal(t, "wind sensor", pi.ModelId)
	assert.Equal(t, uint8(0x23), pi.Info.SourceId)
}

func TestRequest_BroadcastTargetAcceptsAnySource(t *testing.T) {
	c, mb, _ := newCitizenClient(t)

	respondWhenRequested(t, mb, 126996, 0x77, &pgn.ProductInformation{ModelId: "gps"})

	pi, err := Request[*pgn.ProductInformation](context.Background(), c, 255)
	require.NoError(t, err)
	assert.Equal(t, "gps", pi.ModelId)
}

func TestRequest_TimesOutWithDefaultDeadline(t *testing.T) {
	c, _, _ := newCitizenClient(t)

	start := time.Now()
	_, err := Request[*pgn.ProductInformation](context.Background(), c, 0x30)
	elapsed := time.Since(start)

	require.Error(t, err)
	assert.True(t, errors.Is(err, context.DeadlineExceeded), "expected DeadlineExceeded, got %v", err)
	assert.Less(t, elapsed, 3*time.Second, "default deadline should be ~1.25s")
	assert.GreaterOrEqual(t, elapsed, time.Second, "should wait out the default deadline")
}

func TestRequest_ResponseFromWrongSourceIgnored(t *testing.T) {
	c, mb, _ := newCitizenClient(t)

	// Respond from 0x24 while the request targets 0x23.
	respondWhenRequested(t, mb, 126996, 0x24, &pgn.ProductInformation{ModelId: "impostor"})

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	_, err := Request[*pgn.ProductInformation](ctx, c, 0x23)
	require.Error(t, err)
	assert.True(t, errors.Is(err, context.DeadlineExceeded))
}

func TestRequest_ConcurrentSamePGN(t *testing.T) {
	c, mb, _ := newCitizenClient(t)

	respondWhenRequested(t, mb, 126996, 0x23, &pgn.ProductInformation{ModelId: "shared"})

	type result struct {
		pi  *pgn.ProductInformation
		err error
	}
	// A generous explicit deadline so a loaded machine under the race
	// detector has room; this test is about two concurrent requests for the
	// same PGN both completing, not deadline behavior.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	results := make(chan result, 2)
	for range 2 {
		go func() {
			pi, err := Request[*pgn.ProductInformation](ctx, c, 0x23)
			results <- result{pi, err}
		}()
	}
	for range 2 {
		r := <-results
		require.NoError(t, r.err)
		assert.Equal(t, "shared", r.pi.ModelId)
	}
}

func TestRequest_ReplayClientRejected(t *testing.T) {
	c, err := NewClient(context.Background(), Replay(nil))
	require.NoError(t, err)
	defer func() { _ = c.Close() }()

	_, err = Request[*pgn.ProductInformation](context.Background(), c, 0x23)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "bus")
}

func TestRequestSkipsOtherVariantsBeforeBuffering(t *testing.T) {
	co := newCorrelator()
	info := pgn.MessageInfo{PGN: 126208, SourceId: 42, TargetId: pgn.Target(251)}
	waiter, err := co.add(info.PGN, 42, 251, 0, 0, reflect.TypeFor[*pgn.NmeaAcknowledgeGroupFunction]())
	require.NoError(t, err)
	for range 8 {
		co.observe(&pgn.NmeaRequestGroupFunction{Info: info})
	}
	co.observe(&pgn.NmeaAcknowledgeGroupFunction{Info: info})
	require.Len(t, waiter.ch, 1)
	require.IsType(t, &pgn.NmeaAcknowledgeGroupFunction{}, <-waiter.ch)
}

func TestGatewayTransmitEchoStaysOutOfProtocolIngress(t *testing.T) {
	c, mb, local := newCitizenClient(t, WithProductInfo(ProductInfo{ModelID: "local"}))
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	reply := make(chan *pgn.ProductInformation, 1)
	go func() { msg, _ := Request[*pgn.ProductInformation](ctx, c, 255); reply <- msg }()
	require.Eventually(t, func() bool { return infoRequestsTo(mb, 126996, 255) > 0 }, time.Second, time.Millisecond)
	c.handleBusObservation(frameObservation(claimFrame(testDeviceName, 42), "test", "bus", DirectionTransmitted))
	_, claimed := c.DeviceAt(42)
	require.False(t, claimed)
	c.handleBusObservation(frameObservation(isoRequestFrame(126996, local, 255), "test", "bus", DirectionTransmitted))

	payload, err := pgn.EncodeMessage(&pgn.ProductInformation{ModelId: "local"})
	require.NoError(t, err)
	c.handleBusObservation(Observation{Kind: ObservationMessage, PGN: 126996, Source: local, Direction: DirectionTransmitted, Payload: payload})
	sub := c.msgHub.subscribe()
	defer sub.unsubscribe()
	foundEcho := false
	for len(sub.ch) > 0 {
		if product, ok := (<-sub.ch).(*pgn.ProductInformation); ok && product.ModelId == "local" && product.Info.Direction == DirectionTransmitted {
			foundEcho = true
		}
	}
	require.True(t, foundEcho)
	for _, f := range fastPacketFrames(t, &pgn.ProductInformation{ModelId: "remote"}, 6, 42, 255) {
		c.handleBusFrame(f)
	}
	select {
	case got := <-reply:
		require.NotNil(t, got)
		require.Equal(t, "remote", got.ModelId)
	case <-ctx.Done():
		t.Fatal("timed out waiting for the received reply")
	}
	require.Empty(t, framesWithPGN(mb.getWritten(), 126996))
}
