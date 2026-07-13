package n2k

import (
	"context"
	"sync"
	"time"

	"github.com/open-ships/n2k/pgn"
)

// Broadcast schedules periodic transmission of a PGN. provide is called on
// each tick to build the message; returning nil skips that tick (data not
// available yet). The first transmission happens as soon as the client's
// address claim completes.
//
// provide is called once immediately to learn the PGN it produces; if it
// returns nil there, the PGN is learned from its first non-nil message
// instead. Scheduling the same PGN again replaces the earlier schedule. The
// returned stop function halts transmission and is safe to call repeatedly;
// Close also stops all broadcasters.
//
// Another device can retime or pause a broadcast at runtime by sending a
// request group function (PGN 126208) naming the broadcast PGN.
func (c *Client) Broadcast(interval time.Duration, provide func() pgn.Message) (stop func()) {
	b := &broadcaster{
		write:           c.Write,
		provide:         provide,
		onPGNKnown:      c.registerBroadcaster,
		interval:        interval,
		defaultInterval: interval,
		changed:         make(chan struct{}, 1),
		quit:            make(chan struct{}),
		done:            make(chan struct{}),
	}
	if first := provide(); first != nil {
		b.pgn = first.PGNNumber()
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
	prev := c.broadcasters[pgnNum]
	if c.broadcasters == nil {
		c.broadcasters = make(map[uint32]*broadcaster)
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
}

// broadcaster periodically transmits one PGN. Its run loop mirrors the
// heartbeater's: park while the interval is zero, otherwise send then wait
// out the interval, waking early on retiming.
type broadcaster struct {
	write   func(pgn.Message) *WriteResult
	provide func() pgn.Message
	// onPGNKnown registers the broadcaster for group-function lookup once
	// its PGN is known (deferred when provide returned nil at registration).
	onPGNKnown func(*broadcaster)

	mu              sync.Mutex
	pgn             uint32        // 0 until the first non-nil message
	interval        time.Duration // current cadence; <= 0 means paused
	defaultInterval time.Duration // cadence given to Broadcast, for restore

	changed  chan struct{}
	quit     chan struct{}
	done     chan struct{}
	stopOnce sync.Once
}

func (b *broadcaster) run(ctx context.Context, ready <-chan struct{}) {
	defer close(b.done)

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

		b.sendNow()

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
func (b *broadcaster) sendNow() {
	msg := b.provide()
	if msg == nil {
		return
	}
	b.mu.Lock()
	learned := b.pgn == 0
	if learned {
		b.pgn = msg.PGNNumber()
	}
	b.mu.Unlock()
	if learned && b.onPGNKnown != nil {
		b.onPGNKnown(b)
	}
	b.write(msg)
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

// stop terminates the run loop and waits for it to exit. Safe to call more
// than once.
func (b *broadcaster) stop() {
	b.stopOnce.Do(func() { close(b.quit) })
	<-b.done
}
