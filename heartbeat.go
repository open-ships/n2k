package n2k

import (
	"context"
	"sync"
	"time"

	"github.com/open-ships/n2k/pgn"
)

// defaultHeartbeatInterval is the NMEA 2000 heartbeat (PGN 126993) cadence
// required of every bus device.
const defaultHeartbeatInterval = 60 * time.Second

// heartbeater periodically transmits PGN 126993 so other devices can tell we
// are alive. It starts once the client's address claim completes and can be
// retimed at runtime (group functions command interval changes).
type heartbeater struct {
	write func(pgn.Message) *WriteResult

	mu       sync.Mutex
	interval time.Duration // current cadence; <= 0 means stopped
	seq      uint8         // rolling sequence counter, 0-252

	// changed wakes the run loop after setInterval; buffered so setInterval
	// never blocks.
	changed chan struct{}
	// quit stops the run loop ahead of context cancellation so Close can
	// guarantee no heartbeat write races the write channel shutdown.
	quit     chan struct{}
	done     chan struct{}
	stopOnce sync.Once
}

func newHeartbeater(interval time.Duration, write func(pgn.Message) *WriteResult) *heartbeater {
	return &heartbeater{
		write:    write,
		interval: interval,
		changed:  make(chan struct{}, 1),
		quit:     make(chan struct{}),
		done:     make(chan struct{}),
	}
}

// run waits for the address claim, then transmits on the configured cadence.
// A disabled heartbeater (interval <= 0) parks until setInterval re-enables it.
func (h *heartbeater) run(ctx context.Context, ready <-chan struct{}) {
	defer close(h.done)

	select {
	case <-ready:
	case <-ctx.Done():
		return
	case <-h.quit:
		return
	}

	for {
		iv := h.currentInterval()
		if iv <= 0 {
			select {
			case <-h.changed:
				continue
			case <-ctx.Done():
				return
			case <-h.quit:
				return
			}
		}

		h.write(h.message())

		timer := time.NewTimer(iv)
		select {
		case <-timer.C:
		case <-h.changed:
			timer.Stop()
		case <-ctx.Done():
			timer.Stop()
			return
		case <-h.quit:
			timer.Stop()
			return
		}
	}
}

// setInterval changes the heartbeat cadence. Zero (or negative) stops
// transmission; a positive value (re)starts it, sending immediately.
func (h *heartbeater) setInterval(d time.Duration) {
	h.mu.Lock()
	h.interval = d
	h.mu.Unlock()
	select {
	case h.changed <- struct{}{}:
	default:
	}
}

// sendNow transmits a single heartbeat outside the regular cadence.
func (h *heartbeater) sendNow() {
	h.write(h.message())
}

// stop terminates the run loop and waits for it to exit. Safe to call more
// than once.
func (h *heartbeater) stop() {
	h.stopOnce.Do(func() { close(h.quit) })
	<-h.done
}

func (h *heartbeater) currentInterval() time.Duration {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.interval
}

// message builds the next heartbeat, advancing the sequence counter (0-252,
// wrapping; 253+ are reserved/sentinel values on the wire).
func (h *heartbeater) message() *pgn.Heartbeat {
	h.mu.Lock()
	offset := uint64(h.interval / time.Millisecond)
	seq := uint64(h.seq)
	if h.seq >= 252 {
		h.seq = 0
	} else {
		h.seq++
	}
	h.mu.Unlock()

	c1, c2, eq := uint64(pgn.ControllerStateErrorActive), uint64(pgn.ControllerStateErrorActive), uint64(pgn.EquipmentStatusOperational)
	return &pgn.Heartbeat{
		Info:               pgn.MessageInfo{Priority: pgn.Priority(7)},
		DataTransmitOffset: &offset,
		SequenceCounter:    &seq,
		Controller1State:   &c1,
		Controller2State:   &c2,
		EquipmentStatus:    &eq,
	}
}
