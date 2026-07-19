package n2k

import (
	"context"
	"fmt"
	"log/slog"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"

	"github.com/open-ships/n2k/pgn"
)

// Broadcast schedules periodic transmission of a PGN. provide is called on
// each tick to build the message; returning nil skips that tick (data not
// available yet). The first transmission happens as soon as the client's
// address claim completes.
//
// The PGN is learned asynchronously from the first non-nil message. Until
// then, group-function lookup cannot find this schedule; BroadcastPGN avoids
// that delay. Scheduling the same PGN again replaces the earlier schedule.
// The returned stop function is safe to call repeatedly; Close also stops all
// broadcasters.
//
// Another device can retime or pause a broadcast at runtime by sending a
// request group function (PGN 126208) naming the broadcast PGN.
func (c *Client) Broadcast(interval time.Duration, provide func() pgn.Message) (stop func()) {
	return c.newBroadcast(0, false, interval, provide)
}

// BroadcastPGN is Broadcast with an explicit PGN identity. Prefer it when a
// group-function peer must be able to retime the schedule immediately, before
// provide returns its first non-nil message. Declaring the PGN also makes
// replacement of an existing schedule deterministic without probing provide.
func (c *Client) BroadcastPGN(pgnNum uint32, interval time.Duration, provide func() pgn.Message) (stop func()) {
	return c.newBroadcast(pgnNum, true, interval, provide)
}

func (c *Client) newBroadcast(pgnNum uint32, known bool, interval time.Duration, provide func() pgn.Message) func() {
	c.mu.Lock()
	closed := c.closed
	c.mu.Unlock()
	if closed {
		return func() {}
	}
	if known && pgnNum > 0x3FFFF {
		c.log.Error("cannot schedule broadcast with invalid PGN", "pgn", pgnNum)
		return func() {}
	}
	if provide == nil {
		c.log.Error("cannot schedule broadcast with nil provider", "pgn", pgnNum)
		return func() {}
	}
	b := &broadcaster{
		write: c.Write,
		protocolWrite: func(msg pgn.Message) *WriteResult {
			return c.writeProtocol("group-function scheduled PGN response", protocolRequired, msg)
		},
		provide:         provide,
		log:             c.log,
		onPGNKnown:      c.registerBroadcaster,
		onStop:          c.unregisterBroadcaster,
		interval:        interval,
		defaultInterval: interval,
		changed:         make(chan struct{}, 1),
		quit:            make(chan struct{}),
		done:            make(chan struct{}),
		pgn:             pgnNum,
		known:           known,
		order:           c.broadcastSeq.Add(1),
	}
	if known {
		c.registerBroadcaster(b)
	}
	go b.run(c.ctx, c.addrReady)

	return func() {
		b.stop()
		c.unregisterBroadcaster(b)
	}
}

// registerBroadcaster installs b as the broadcaster for its PGN, stopping any
// earlier schedule for the same PGN.
func (c *Client) registerBroadcaster(b *broadcaster) {
	pgnNum := b.pgnNumber()
	c.bMu.Lock()
	if c.broadcasters == nil {
		c.broadcasters = make(map[uint32]*broadcaster)
	}
	prev := c.broadcasters[pgnNum]
	if prev != nil && prev.order > b.order {
		c.bMu.Unlock()
		b.stop()
		return
	}
	c.broadcasters[pgnNum] = b
	c.bMu.Unlock()
	if prev != nil && prev != b {
		prev.stop()
	}
}

// unregisterBroadcaster removes b unless it has already been replaced.
func (c *Client) unregisterBroadcaster(b *broadcaster) {
	pgnNum := b.pgnNumber()
	c.bMu.Lock()
	defer c.bMu.Unlock()
	if c.broadcasters[pgnNum] == b {
		delete(c.broadcasters, pgnNum)
	}
}

// broadcasterFor returns the active broadcaster for a PGN, or nil.
func (c *Client) broadcasterFor(pgnNum uint32) *broadcaster {
	c.bMu.Lock()
	defer c.bMu.Unlock()
	return c.broadcasters[pgnNum]
}

// stopBroadcasters halts every scheduled broadcast (used by Close).
func (c *Client) stopBroadcasters() {
	c.bMu.Lock()
	all := make([]*broadcaster, 0, len(c.broadcasters))
	for _, b := range c.broadcasters {
		all = append(all, b)
	}
	c.broadcasters = nil
	c.bMu.Unlock()
	for _, b := range all {
		b.stop()
	}
	for _, b := range all {
		b.wait()
	}
}

// broadcaster periodically transmits one PGN. Its run loop mirrors the
// heartbeater's: park while the interval is zero, otherwise send then wait
// out the interval, waking early on retiming.
type broadcaster struct {
	write         func(pgn.Message) *WriteResult
	protocolWrite func(pgn.Message) *WriteResult
	provide       func() pgn.Message
	log           *slog.Logger
	// onPGNKnown registers the broadcaster for group-function lookup once
	// its PGN is known (deferred when provide returned nil at registration).
	onPGNKnown func(*broadcaster)
	onStop     func(*broadcaster)

	mu              sync.Mutex
	pgn             uint32 // valid once known is true
	known           bool
	order           uint64        // creation order makes deferred registration deterministic
	interval        time.Duration // current cadence; <= 0 means paused
	defaultInterval time.Duration // cadence given to Broadcast, for restore

	changed  chan struct{}
	quit     chan struct{}
	done     chan struct{}
	stopOnce sync.Once
	sending  atomic.Bool
}

func (b *broadcaster) run(ctx context.Context, ready <-chan struct{}) {
	defer close(b.done)
	defer func() {
		if b.onStop != nil {
			b.onStop(b)
		}
	}()

	select {
	case <-ready:
	case <-ctx.Done():
		return
	case <-b.quit:
		return
	}

	for {
		iv := b.currentInterval()
		if iv <= 0 {
			select {
			case <-b.changed:
				continue
			case <-ctx.Done():
				return
			case <-b.quit:
				return
			}
		}

		if !b.sendNow(ctx, b.write) {
			return
		}

		timer := time.NewTimer(iv)
		select {
		case <-timer.C:
		case <-b.changed:
			timer.Stop()
		case <-ctx.Done():
			timer.Stop()
			return
		case <-b.quit:
			timer.Stop()
			return
		}
	}
}

// sendNow transmits one message immediately (nil from provide skips),
// completing deferred PGN registration on the first non-nil message.
func (b *broadcaster) sendNow(ctx context.Context, write func(pgn.Message) *WriteResult) bool {
	if !b.sending.CompareAndSwap(false, true) {
		return true
	}
	defer b.sending.Store(false)

	type provided struct {
		msg pgn.Message
		err error
	}
	result := make(chan provided, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				result <- provided{err: fmt.Errorf("provider panic: %v", r)}
				b.log.Error("recovered panic in broadcast provider", "panic", r, "stack", string(debug.Stack()))
			}
		}()
		result <- provided{msg: b.provide()}
	}()
	var msg pgn.Message
	select {
	case got := <-result:
		if got.err != nil {
			return true
		}
		msg = got.msg
	case <-ctx.Done():
		return false
	case <-b.quit:
		return false
	}
	if msg == nil {
		return true
	}
	b.mu.Lock()
	learned := !b.known
	if learned {
		b.pgn = msg.PGNNumber()
		b.known = true
	}
	b.mu.Unlock()
	if !learned && msg.PGNNumber() != b.pgnNumber() {
		b.log.Warn("broadcast provider returned a different PGN", "want", b.pgnNumber(), "got", msg.PGNNumber())
		return true
	}
	if learned && b.onPGNKnown != nil {
		b.onPGNKnown(b)
		select {
		case <-b.quit:
			return false
		default:
		}
	}
	if err := write(msg).WaitContext(ctx); err != nil && ctx.Err() == nil {
		b.log.Warn("scheduled broadcast failed", "pgn", msg.PGNNumber(), "error", err)
	}
	return true
}

func (b *broadcaster) trigger(ctx context.Context) {
	go func() { _ = b.sendNow(ctx, b.protocolWrite) }()
}

// setInterval retimes the broadcast. Zero (or negative) pauses it; a positive
// value (re)starts it, sending immediately.
func (b *broadcaster) setInterval(d time.Duration) {
	b.mu.Lock()
	b.interval = d
	b.mu.Unlock()
	select {
	case b.changed <- struct{}{}:
	default:
	}
}

// restoreDefaultInterval reverts to the cadence given to Broadcast.
func (b *broadcaster) restoreDefaultInterval() {
	b.mu.Lock()
	d := b.defaultInterval
	b.mu.Unlock()
	b.setInterval(d)
}

func (b *broadcaster) currentInterval() time.Duration {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.interval
}

func (b *broadcaster) pgnNumber() uint32 {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.pgn
}

// stop requests termination without waiting, so it is safe to invoke from a
// provider callback. Client.Close waits separately for broadcaster exit.
func (b *broadcaster) stop() {
	b.stopOnce.Do(func() { close(b.quit) })
}

func (b *broadcaster) wait() {
	<-b.done
}
