package n2k

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/open-ships/n2k/internal/framer"
	"github.com/open-ships/n2k/pgn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// respondWhenRequested watches the mock bus for an ISO request for pgnNum
// addressed to responder, then injects the given message from that address.
func respondWhenRequested(t *testing.T, mb *mockBus, pgnNum uint32, responder uint8, msg pgn.Message) {
	t.Helper()
	go func() {
		deadline := time.Now().Add(5 * time.Second)
		for time.Now().Before(deadline) {
			for _, f := range mb.getWritten() {
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
				return
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
	// A generous explicit deadline: under the race detector on a loaded
	// machine the default 1.25s ISO deadline can flake, and this test is
	// about concurrent request coalescing, not deadline behavior.
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
