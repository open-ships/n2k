package n2k

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"runtime/debug"
	"sync"
	"time"

	"github.com/open-ships/n2k/pgn"
)

const maxBroadcastWorkers = 64

var (
	// ErrBroadcastLimit reports that all bounded schedule slots are occupied.
	ErrBroadcastLimit = errors.New("n2k: broadcast schedule limit reached")
	// ErrBroadcastQueueFull reports that a required response cannot be queued.
	ErrBroadcastQueueFull = errors.New("n2k: broadcast response queue full")
)

// Broadcast schedules a PGN, learning its identity from the first non-nil
// result. Prefer BroadcastPGN when peers need immediate group-function lookup.
// A nil result skips a periodic tick. Scheduling an existing PGN replaces it.
//
// The provider runs synchronously on its schedule's owned worker. It must
// return promptly when ctx ends and must not call its own stop function or
// Client.Close; those operations wait for the provider to finish. Stop is
// idempotent, cancels the provider, and waits for the worker to exit.
func (c *Client) Broadcast(interval time.Duration, provide func(context.Context) pgn.Message) (stop func(), err error) {
	return c.newBroadcast(0, false, interval, provide)
}

// BroadcastPGN declares a schedule's PGN before its first provider call, so
// group functions can retime or request it immediately. At most 64 schedule
// workers may be alive, including providers exiting after replacement.
func (c *Client) BroadcastPGN(pgnNum uint32, interval time.Duration, provide func(context.Context) pgn.Message) (stop func(), err error) {
	return c.newBroadcast(pgnNum, true, interval, provide)
}

func (c *Client) newBroadcast(pgnNum uint32, known bool, interval time.Duration, provide func(context.Context) pgn.Message) (func(), error) {
	if c == nil {
		return nil, ErrClientClosed
	}
	if provide == nil {
		return nil, errors.New("n2k: cannot schedule broadcast with nil provider")
	}
	if known && (pgnNum > 0x1FFFF || pgnNum>>8&0xFF < 240 && pgnNum&0xFF != 0) {
		return nil, fmt.Errorf("n2k: cannot schedule invalid PGN %d", pgnNum)
	}

	// Admission and worker ownership are atomic with Close. Unknown-PGN and
	// replaced schedules remain owned until their workers actually exit.
	c.mu.Lock()
	if c.closed || c.terminalErr != nil || c.ctx.Err() != nil {
		c.mu.Unlock()
		return nil, c.operationError()
	}
	c.bMu.Lock()
	if len(c.broadcastWorkers) >= maxBroadcastWorkers {
		c.bMu.Unlock()
		c.mu.Unlock()
		return nil, ErrBroadcastLimit
	}
	ctx, cancel := context.WithCancel(c.ctx)
	b := &broadcaster{
		ctx: ctx, cancel: cancel, write: c.WriteContext,
		bindRequest: c.stampProtocolContext,
		protocolWrite: func(ctx context.Context, msg pgn.Message) *WriteResult {
			return c.writeProtocolContext(ctx, "group-function scheduled PGN response", protocolRequired, msg)
		},
		provide: provide, log: c.log, onError: c.fail,
		onPGNKnown: c.registerBroadcaster, onStop: c.unregisterBroadcaster,
		interval: interval, defaultInterval: interval,
		changed: make(chan struct{}, 1), requests: make(chan broadcastRequest, 1),
		done: make(chan struct{}), pgn: pgnNum, known: known,
		order: c.broadcastSeq.Add(1),
	}
	if c.broadcastWorkers == nil {
		c.broadcastWorkers = make(map[*broadcaster]struct{})
	}
	c.broadcastWorkers[b] = struct{}{}
	var previous *broadcaster
	if known {
		if c.broadcasters == nil {
			c.broadcasters = make(map[uint32]*broadcaster)
		}
		previous = c.broadcasters[pgnNum]
		c.broadcasters[pgnNum] = b
	}
	go b.run(c.addrReady)
	c.bMu.Unlock()
	c.mu.Unlock()
	if previous != nil {
		previous.cancel()
	}
	return b.stop, nil
}

// registerBroadcaster completes deferred PGN discovery without allowing an
// older provider to replace a newer schedule for that PGN.
func (c *Client) registerBroadcaster(b *broadcaster) {
	pgnNum := b.pgnNumber()
	c.mu.Lock()
	c.bMu.Lock()
	if c.closed || b.ctx.Err() != nil {
		c.bMu.Unlock()
		c.mu.Unlock()
		b.cancel()
		return
	}
	if c.broadcasters == nil {
		c.broadcasters = make(map[uint32]*broadcaster)
	}
	previous := c.broadcasters[pgnNum]
	if previous != nil && previous.order > b.order {
		c.bMu.Unlock()
		c.mu.Unlock()
		b.cancel()
		return
	}
	c.broadcasters[pgnNum] = b
	c.bMu.Unlock()
	c.mu.Unlock()
	if previous != nil && previous != b {
		previous.cancel()
	}
}

func (c *Client) unregisterBroadcaster(b *broadcaster) {
	pgnNum := b.pgnNumber()
	c.bMu.Lock()
	defer c.bMu.Unlock()
	if c.broadcasters[pgnNum] == b {
		delete(c.broadcasters, pgnNum)
	}
	delete(c.broadcastWorkers, b)
}

func (c *Client) broadcasterFor(pgnNum uint32) *broadcaster {
	c.bMu.Lock()
	defer c.bMu.Unlock()
	return c.broadcasters[pgnNum]
}

func (c *Client) stopBroadcasters() {
	c.bMu.Lock()
	all := make([]*broadcaster, 0, len(c.broadcastWorkers))
	for b := range c.broadcastWorkers {
		all = append(all, b)
	}
	c.bMu.Unlock()
	for _, b := range all {
		b.cancel()
	}
	for _, b := range all {
		b.wait()
	}
}

type broadcaster struct {
	ctx           context.Context
	cancel        context.CancelFunc
	write         func(context.Context, pgn.Message) *WriteResult
	protocolWrite func(context.Context, pgn.Message) *WriteResult
	bindRequest   func(context.Context) (context.Context, func())
	provide       func(context.Context) pgn.Message
	log           *slog.Logger
	onError       func(error)
	onPGNKnown    func(*broadcaster)
	onStop        func(*broadcaster)

	mu              sync.Mutex
	pgn             uint32
	known           bool
	order           uint64
	interval        time.Duration
	defaultInterval time.Duration
	stopped         bool

	changed  chan struct{}
	requests chan broadcastRequest
	done     chan struct{}
}

type broadcastRequest struct {
	ctx  context.Context
	stop func()
}

func (b *broadcaster) run(ready <-chan struct{}) {
	defer close(b.done)
	defer b.cancel()
	defer b.drainRequests()
	defer func() {
		if b.onStop != nil {
			b.onStop(b)
		}
	}()
	select {
	case <-ready:
	case <-b.ctx.Done():
		return
	}

	timer := time.NewTimer(time.Hour)
	defer timer.Stop()
	var ticks <-chan time.Time
	reset := func(immediate bool) {
		timer.Stop()
		ticks = nil
		if interval := b.currentInterval(); interval > 0 {
			if immediate {
				interval = 0
			}
			timer.Reset(interval)
			ticks = timer.C
		}
	}
	reset(true)
	for {
		if b.ctx.Err() != nil {
			return
		}
		// A required response takes the next provider slot before a due tick.
		select {
		case request := <-b.requests:
			if !b.sendRequested(request) {
				return
			}
			continue
		default:
		}
		select {
		case <-b.ctx.Done():
			return
		case request := <-b.requests:
			if !b.sendRequested(request) {
				return
			}
		case <-b.changed:
			reset(true)
		case <-ticks:
			if !b.sendNow(b.ctx, false) {
				return
			}
			reset(false)
		}
	}
}

func (b *broadcaster) message(ctx context.Context) (msg pgn.Message, err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			err = fmt.Errorf("provider panic: %v", recovered)
			b.log.Error("recovered panic in broadcast provider", "panic", recovered, "stack", string(debug.Stack()))
		}
	}()
	return b.provide(ctx), nil
}

func (b *broadcaster) sendRequested(request broadcastRequest) bool {
	defer request.stop()
	return b.sendNow(request.ctx, true)
}

func (b *broadcaster) sendNow(ctx context.Context, required bool) bool {
	if b.ctx.Err() != nil {
		return false
	}
	if ctx.Err() != nil {
		return true
	}
	msg, err := b.message(ctx)
	if b.ctx.Err() != nil {
		return false
	}
	if ctx.Err() != nil {
		return true
	}
	if err == nil && msg == nil {
		if !required {
			return true
		}
		err = errors.New("broadcast provider returned no message for a required response")
	}
	if err != nil {
		return b.report(err, required)
	}
	b.mu.Lock()
	learned := !b.known
	if learned {
		b.pgn = msg.PGNNumber()
		b.known = true
	}
	b.mu.Unlock()
	if !learned && msg.PGNNumber() != b.pgnNumber() {
		return b.report(fmt.Errorf("broadcast provider returned PGN %d, expected %d", msg.PGNNumber(), b.pgnNumber()), required)
	}
	if learned && b.onPGNKnown != nil {
		b.onPGNKnown(b)
	}
	if b.ctx.Err() != nil {
		return false
	}
	write := b.write
	if required {
		write = b.protocolWrite
	}
	if err := write(ctx, msg).WaitContext(ctx); err != nil && ctx.Err() == nil && b.ctx.Err() == nil {
		return b.report(err, required)
	}
	return b.ctx.Err() == nil
}

func (b *broadcaster) report(err error, required bool) bool {
	if errors.Is(err, ErrEpochChanged) || errors.Is(err, context.Canceled) {
		return b.ctx.Err() == nil
	}
	if required {
		if b.onError != nil {
			b.onError(fmt.Errorf("n2k: required broadcast response failed: %w", err))
		}
		return false
	}
	b.log.Warn("scheduled broadcast failed", "pgn", b.pgnNumber(), "error", err)
	return true
}

// trigger admits one required response without spawning work. A second queued
// response is rejected explicitly, allowing the protocol owner to fail closed.
func (b *broadcaster) trigger(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if b.ctx.Err() != nil {
		return context.Canceled
	}
	requestCtx, cancel := context.WithCancelCause(ctx)
	stopSchedule := context.AfterFunc(b.ctx, func() { cancel(context.Cause(b.ctx)) })
	requestCtx, stopEpoch := b.bindRequest(requestCtx)
	request := broadcastRequest{ctx: requestCtx, stop: func() {
		stopEpoch()
		stopSchedule()
		cancel(context.Canceled)
	}}
	b.mu.Lock()
	if b.stopped || b.ctx.Err() != nil {
		b.mu.Unlock()
		request.stop()
		return context.Canceled
	}
	select {
	case b.requests <- request:
		b.mu.Unlock()
		return nil
	default:
		b.mu.Unlock()
		request.stop()
		return ErrBroadcastQueueFull
	}
}

func (b *broadcaster) drainRequests() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.stopped = true
	for {
		select {
		case request := <-b.requests:
			request.stop()
		default:
			return
		}
	}
}

func (b *broadcaster) setInterval(d time.Duration) {
	b.mu.Lock()
	b.interval = d
	b.mu.Unlock()
	select {
	case b.changed <- struct{}{}:
	default:
	}
}

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

func (b *broadcaster) stop() {
	b.cancel()
	b.wait()
}

func (b *broadcaster) wait() { <-b.done }
