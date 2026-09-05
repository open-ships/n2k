package transport

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/brutella/can"
	"github.com/open-ships/n2k/internal/framer"
	"github.com/stretchr/testify/require"
)

// expiredTimer models an AfterFunc callback which has already started and is
// waiting for the manager lock. Stop cannot prevent its later execution.
type expiredTimer struct{}

func (expiredTimer) Stop() bool { return false }

type racingClock struct {
	mu        sync.Mutex
	callbacks map[time.Duration][]func()
}

func (c *racingClock) afterFunc(duration time.Duration, callback func()) Timer {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.callbacks == nil {
		c.callbacks = make(map[time.Duration][]func())
	}
	c.callbacks[duration] = append(c.callbacks[duration], callback)
	return expiredTimer{}
}

func (c *racingClock) callback(duration time.Duration, index int) func() {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.callbacks[duration][index]
}

func awaitTransmit(t *testing.T, result <-chan error) error {
	t.Helper()
	select {
	case err := <-result:
		return err
	case <-time.After(time.Second):
		t.Fatal("transport transmission did not finish")
		return nil
	}
}

func awaitFrame(t *testing.T, frames <-chan can.Frame) can.Frame {
	t.Helper()
	select {
	case frame := <-frames:
		return frame
	case <-time.After(time.Second):
		t.Fatal("transport did not emit an expected frame")
		return can.Frame{}
	}
}

func TestObsoleteReceiveTimersPreserveReplacementAndRearmedTransfer(t *testing.T) {
	clock := &racingClock{}
	h := newTestHelper()
	mgr := NewManager(ManagerConfig{AfterFunc: clock.afterFunc, OnComplete: h.onComplete})
	defer mgr.Close()
	cm := buildCMBAMFrame(16, 3, 126998, 10)
	mgr.HandleFrame(cm)
	oldTransferTimeout := clock.callback(DTTimeout, 0)
	mgr.HandleFrame(cm)
	oldPacketTimeout := clock.callback(DTTimeout, 1)
	frames := buildDTFrameRange(make([]byte, 16), 1, 3, 10, BroadcastAddr)
	mgr.HandleFrame(frames[0])
	oldTransferTimeout()
	oldPacketTimeout()
	mgr.HandleFrame(frames[1])
	mgr.HandleFrame(frames[2])
	require.Len(t, h.getCompleted(), 1)
	require.Zero(t, activeSessionCount(mgr))
}

func TestObsoleteTransmitTimerCannotCompleteReusedKey(t *testing.T) {
	clock := &racingClock{}
	writes := make(chan can.Frame, 10)
	mgr := NewManager(ManagerConfig{AfterFunc: clock.afterFunc, WriteFrame: func(frame can.Frame) error {
		writes <- frame
		return nil
	}})
	defer mgr.Close()
	result := make(chan error, 1)
	go func() { result <- mgr.SendRTSCTS(126998, 10, 42, make([]byte, 9)) }()
	awaitFrame(t, writes)
	oldTimeout := clock.callback(CTSTimeout, 0)
	oldDeadline := clock.callback(DefaultTransferTimeout, 0)
	mgr.HandleFrame(buildCTSFrame(2, 1, 126998, 42, 10))
	awaitFrame(t, writes)
	awaitFrame(t, writes)
	mgr.HandleFrame(buildTestEndOfMsgAckFrame(9, 2, 126998, 42, 10))
	require.NoError(t, awaitTransmit(t, result))

	go func() { result <- mgr.SendRTSCTS(126998, 10, 42, make([]byte, 9)) }()
	awaitFrame(t, writes)
	require.NotPanics(t, oldTimeout)
	require.NotPanics(t, oldDeadline)
	require.Equal(t, 1, activeSessionCount(mgr))
	mgr.HandleFrame(buildCTSFrame(2, 1, 126998, 42, 10))
	mgr.HandleFrame(buildTestEndOfMsgAckFrame(9, 2, 126998, 42, 10))
	require.NoError(t, awaitTransmit(t, result))
}

func TestReceiverHoldsCannotExtendAbsoluteTransferDeadline(t *testing.T) {
	clock := &racingClock{}
	writes := make(chan can.Frame, 1)
	mgr := NewManager(ManagerConfig{AfterFunc: clock.afterFunc, WriteFrame: func(frame can.Frame) error {
		writes <- frame
		return nil
	}})
	defer mgr.Close()
	result := make(chan error, 1)
	go func() { result <- mgr.SendRTSCTS(126998, 10, 42, make([]byte, 9)) }()
	awaitFrame(t, writes)
	for i := 0; i < 10; i++ {
		obsolete := clock.callback(CTSTimeout, i)
		mgr.HandleFrame(buildCTSFrame(0, 1, 126998, 42, 10))
		obsolete()
		require.Equal(t, 1, activeSessionCount(mgr))
	}
	clock.callback(DefaultTransferTimeout, 0)()
	require.ErrorIs(t, awaitTransmit(t, result), context.DeadlineExceeded)
	require.Zero(t, activeSessionCount(mgr))
}

func TestBAMCancellationInterruptsPacing(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	writes := make(chan can.Frame, 4)
	mgr := NewManager(ManagerConfig{WriteFrame: func(frame can.Frame) error {
		writes <- frame
		if framer.ParseCANID(frame.ID).PGN == PGNDT {
			cancel()
		}
		return nil
	}})
	defer mgr.Close()
	result := make(chan error, 1)
	go func() { result <- mgr.SendBAMContext(ctx, 126998, 10, make([]byte, 100)) }()
	awaitFrame(t, writes)
	firstDT := awaitFrame(t, writes)
	require.Equal(t, uint8(1), firstDT.Data[0])
	require.ErrorIs(t, awaitTransmit(t, result), context.Canceled)
	require.Empty(t, writes, "cancellation must suppress subsequent DT packets")
	require.Zero(t, activeSessionCount(mgr))
}

func TestTransferDeadlineInterruptsBlockedWrite(t *testing.T) {
	mgr := NewManager(ManagerConfig{
		TransferTimeout: 20 * time.Millisecond,
		WriteFrameContext: func(ctx context.Context, _ can.Frame) error {
			<-ctx.Done()
			return context.Cause(ctx)
		},
	})
	defer mgr.Close()
	result := make(chan error, 1)
	go func() { result <- mgr.SendRTSCTS(126998, 10, 42, make([]byte, 9)) }()
	require.ErrorIs(t, awaitTransmit(t, result), context.DeadlineExceeded)
	require.Zero(t, activeSessionCount(mgr))
}

func TestCanceledCallerDoesNotStartTransfer(t *testing.T) {
	h := newTestHelper()
	mgr := NewManager(ManagerConfig{WriteFrame: h.writeFrame})
	defer mgr.Close()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	require.ErrorIs(t, mgr.SendBAMContext(ctx, 126998, 10, make([]byte, 9)), context.Canceled)
	require.ErrorIs(t, mgr.SendRTSCTSContext(ctx, 126998, 10, 42, make([]byte, 9)), context.Canceled)
	require.Empty(t, h.getSentFrames())
}

func TestResetCancelsWritesAndAllowsFreshTransfers(t *testing.T) {
	for _, bam := range []bool{false, true} {
		t.Run(map[bool]string{false: "RTSCTS", true: "BAM"}[bam], func(t *testing.T) {
			entered := make(chan struct{}, 1)
			resetErr := errors.New("connection epoch changed")
			mgr := NewManager(ManagerConfig{WriteFrameContext: func(ctx context.Context, _ can.Frame) error {
				entered <- struct{}{}
				<-ctx.Done()
				return context.Cause(ctx)
			}})
			defer mgr.Close()
			result := make(chan error, 1)
			send := func() {
				if bam {
					result <- mgr.SendBAM(126998, 10, make([]byte, 9))
				} else {
					result <- mgr.SendRTSCTS(126998, 10, 42, make([]byte, 9))
				}
			}
			go send()
			<-entered
			mgr.Reset(resetErr)
			require.ErrorIs(t, awaitTransmit(t, result), resetErr)
			go send()
			<-entered
			mgr.Reset(resetErr)
			require.ErrorIs(t, awaitTransmit(t, result), resetErr)
			require.Zero(t, activeSessionCount(mgr))
		})
	}
}

func TestResetDropsReceiveFragmentsAndCloseRejectsAnnouncements(t *testing.T) {
	h := newTestHelper()
	mgr := NewManager(ManagerConfig{OnComplete: h.onComplete, WriteFrame: h.writeFrame})
	cm := buildCMBAMFrame(9, 2, 126998, 10)
	frames := buildDTFrameRange(make([]byte, 9), 1, 2, 10, BroadcastAddr)
	mgr.HandleFrame(cm)
	mgr.HandleFrame(frames[0])
	mgr.Reset(errors.New("new network epoch"))
	mgr.HandleFrame(frames[1])
	require.Empty(t, h.getCompleted())
	mgr.HandleFrame(cm)
	for _, frame := range frames {
		mgr.HandleFrame(frame)
	}
	require.Len(t, h.getCompleted(), 1)
	mgr.Close()
	mgr.HandleFrame(cm)
	mgr.HandleFrame(buildRTSFrame(9, 2, 126998, 11, 42))
	require.Zero(t, activeSessionCount(mgr))
	require.Empty(t, h.getSentFrames())
}

func TestBAMCollisionAndLoopbackCannotReplaceActiveTransfer(t *testing.T) {
	writes := make(chan can.Frame, 4)
	blocked := make(chan struct{})
	release := make(chan struct{})
	mgr := NewManager(ManagerConfig{WriteFrame: func(frame can.Frame) error {
		writes <- frame
		if framer.ParseCANID(frame.ID).PGN == PGNCM {
			close(blocked)
			<-release
		}
		return nil
	}})
	defer mgr.Close()
	result := make(chan error, 1)
	go func() { result <- mgr.SendBAM(126998, 10, make([]byte, 9)) }()
	<-blocked
	cm := awaitFrame(t, writes)
	mgr.HandleFrame(cm)
	err := mgr.SendBAM(126996, 10, make([]byte, 9))
	require.ErrorContains(t, err, "session already active")
	resetErr := errors.New("address changed")
	mgr.Reset(resetErr)
	close(release)
	require.ErrorIs(t, awaitTransmit(t, result), resetErr)
	require.Empty(t, writes, "a stale announcement write must not continue with DT packets")
}

func TestSessionLimitIncludesReceiveAndTransmit(t *testing.T) {
	mgr := NewManager(ManagerConfig{MaxSessions: 1})
	defer mgr.Close()
	mgr.HandleFrame(buildCMBAMFrame(9, 2, 126998, 10))
	mgr.HandleFrame(buildCMBAMFrame(9, 2, 126998, 11))
	require.Equal(t, 1, activeSessionCount(mgr))
	require.ErrorContains(t, mgr.SendBAM(126998, 12, make([]byte, 9)), "table is full")
	// Replacing an existing receive transfer consumes no additional slot.
	mgr.HandleFrame(buildCMBAMFrame(9, 2, 126996, 10))
	require.Equal(t, 1, activeSessionCount(mgr))
}
