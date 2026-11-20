package n2k

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/brutella/can"
	"github.com/open-ships/n2k/internal/framer"
	"github.com/open-ships/n2k/pgn"
)

var (
	// ErrNotReady reports that application transmission cannot be admitted
	// until the current connection has completed address claiming.
	ErrNotReady = errors.New("n2k: network session is not ready")
	// ErrEpochChanged reports work invalidated by disconnect or address change.
	ErrEpochChanged = errors.New("n2k: network session epoch changed")
)

func (c *Client) readyErrorLocked() error {
	if c.closed {
		return ErrClientClosed
	}
	if c.terminalErr != nil {
		return c.terminalErr
	}
	if c.ctx != nil && c.ctx.Err() != nil {
		return c.ctx.Err()
	}
	if c.bus != nil {
		if !c.connected || !c.claimed {
			return ErrNotReady
		}
		select {
		case <-c.txReady:
		default:
			return ErrNotReady
		}
	}
	return nil
}

func (c *Client) operationError() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.terminalErr != nil {
		return c.terminalErr
	}
	if c.closed {
		return ErrClientClosed
	}
	if c.ctx != nil && c.ctx.Err() != nil {
		return c.ctx.Err()
	}
	return ErrNotReady
}

// writeContextLocked binds admission to both caller cancellation and the
// current network epoch. c.mu makes registration atomic with invalidation.
func (c *Client) writeContextLocked(parent context.Context, protocol bool) (context.Context, func()) {
	if parent == nil {
		parent = context.Background()
	}
	ctx, cancel := context.WithCancelCause(parent)
	stop := func() bool { return true }
	if c.epochCtx != nil {
		epoch := c.epochCtx
		stop = context.AfterFunc(epoch, func() { cancel(context.Cause(epoch)) })
		if cause := context.Cause(epoch); cause != nil {
			cancel(cause)
		}
	}
	connection, claim := c.connectionEpoch, c.claimEpoch
	if inherited, ok := parent.Value(writeStampKey{}).(writeStamp); ok {
		// Retain the admission identity even before asynchronous epoch
		// cancellation callbacks have run. Rebinding cannot revive stale work.
		connection, claim = inherited.connection, inherited.claim
	}
	ctx = context.WithValue(ctx, writeStampKey{}, writeStamp{
		connection: connection, claim: claim, protocol: protocol, progress: &writeProgress{},
	})
	return ctx, func() { stop(); cancel(context.Canceled) }
}

func (c *Client) stampProtocolContext(ctx context.Context) (context.Context, func()) {
	if _, ok := ctx.Value(writeStampKey{}).(writeStamp); ok {
		return ctx, func() {}
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.writeContextLocked(ctx, true)
}

func (c *Client) checkWriteContext(ctx context.Context) error {
	if cause := context.Cause(ctx); cause != nil {
		return cause
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return ErrClientClosed
	}
	if c.terminalErr != nil {
		return c.terminalErr
	}
	if stamp, ok := ctx.Value(writeStampKey{}).(writeStamp); ok {
		if stamp.connection != c.connectionEpoch || stamp.claim != c.claimEpoch {
			return ErrEpochChanged
		}
	}
	return nil
}

func (c *Client) currentMessageEpoch(info pgn.MessageInfo) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return !c.closed && c.terminalErr == nil && info.ConnectionEpoch == c.connectionEpoch && info.ClaimEpoch == c.claimEpoch
}

const addressContentionWindow = 250 * time.Millisecond

// finishStartupClaim publishes the address installed by the claimer's change
// callback if it has settled by the caller's deadline. Checking the current
// gate and publishing under one lock prevents a superseded claim from winning.
func (c *Client) finishStartupClaim(connection uint64, settledAt, deadline time.Time) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.terminalErr != nil {
		return c.terminalErr
	}
	if c.closed {
		return ErrClientClosed
	}
	if err := c.ctx.Err(); err != nil {
		return err
	}
	if !c.connected || c.connectionEpoch != connection {
		return fmt.Errorf("n2k: connection changed during initial address claim: %w", ErrEpochChanged)
	}
	if c.txReady != c.addrReady {
		settledAt = c.contentionEndsAt
	}
	if settledAt.IsZero() || settledAt.After(deadline) {
		return fmt.Errorf("n2k: address claim did not settle within claim timeout: %w", context.DeadlineExceeded)
	}
	if c.contentionTimer != nil {
		// The deadline and contention timers may become runnable together.
		c.stopContentionTimerLocked()
		close(c.txReady)
	}
	c.claimed = true
	close(c.addrReady)
	return nil
}

// writeAddressClaim starts a replacement address's contention window only
// after the wire writer has accepted its claim, within the same lifetime.
func (c *Client) writeAddressClaim(frame can.Frame) error {
	address := framer.ParseCANID(frame.ID).Source
	c.mu.Lock()
	if address != c.sourceAddr || !c.connected {
		c.mu.Unlock()
		return ErrEpochChanged
	}
	ctx, stop := c.writeContextLocked(c.ctx, true)
	gate := c.txReady
	c.mu.Unlock()
	defer stop()
	if err := c.writeFrameContext(ctx, frame); err != nil {
		return err
	}
	acceptedAt := time.Now()
	c.mu.Lock()
	if c.replacementClaim && gate == c.txReady && c.sourceAddr == address {
		if c.sourceAddr == 254 {
			c.mu.Unlock()
			c.fail(errors.New("n2k: address claiming failed: no address available"))
			return nil
		}
		c.startContentionTimerLocked(acceptedAt)
	}
	c.mu.Unlock()
	return nil
}

func (c *Client) startContentionTimerLocked(acceptedAt time.Time) {
	if c.closed || c.terminalErr != nil || !c.connected || c.contentionTimer != nil {
		return
	}
	c.replacementClaim = false
	ready := c.txReady
	c.contentionEndsAt = acceptedAt.Add(addressContentionWindow)
	c.contentionTimer = time.AfterFunc(time.Until(c.contentionEndsAt), func() {
		c.mu.Lock()
		defer c.mu.Unlock()
		if c.txReady == ready && c.contentionTimer != nil {
			close(ready)
			c.contentionTimer = nil
		}
	})
}

func (c *Client) stopContentionTimerLocked() {
	c.replacementClaim = false
	c.contentionEndsAt = time.Time{}
	if c.contentionTimer != nil {
		c.contentionTimer.Stop()
		c.contentionTimer = nil
	}
}

// heartbeatReadiness returns the lifetime that invalidates a scheduled tick
// and a gate to wait on when the claim is not ready. A nil gate means ready.
func (c *Client) heartbeatReadiness() (context.Context, <-chan struct{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.claimed {
		return c.epochCtx, c.addrReady
	}
	if c.readyErrorLocked() != nil {
		return c.epochCtx, c.txReady
	}
	return c.epochCtx, nil
}

// invalidateEpochLocked owns pending operation lifetimes. Protocol parsers are
// reset after releasing c.mu, avoiding callback/Client lock inversions.
func (c *Client) invalidateEpochLocked(err error) {
	c.stopContentionTimerLocked()
	if c.epochCancel != nil {
		c.epochCancel(err)
	}
	c.claimEpoch++
	if c.ctx != nil {
		c.epochCtx, c.epochCancel = context.WithCancelCause(c.ctx)
	}
	if c.correlator != nil {
		c.correlator.invalidate(err)
	}
	if c.actisenseRemote != nil {
		c.actisenseRemote.invalidate(ErrActisenseRemoteEpochChanged)
	}
}

func (c *Client) resetReadEpoch() {
	c.mu.Lock()
	connection, claim := c.connectionEpoch, c.claimEpoch
	c.mu.Unlock()
	if c.tp != nil {
		c.tp.Reset(ErrEpochChanged)
	}
	if c.pipeline != nil {
		c.pipeline.resetEpoch(connection, claim)
	}
	if c.system != nil {
		c.system.pipeline.resetEpoch(connection, claim)
	}
}

type deviceInfoRequest struct {
	address    uint8
	name       uint64
	connection uint64
	claim      uint64
}

func (c *Client) deviceInfoLoop() {
	defer c.backgroundWg.Done()
	for {
		select {
		case <-c.ctx.Done():
			return
		case request := <-c.deviceInfoCh:
			select {
			case <-c.addrReady:
			case <-c.ctx.Done():
				return
			}
			for _, requested := range []uint32{126996, 126998} {
				c.mu.Lock()
				current := request.connection == c.connectionEpoch && request.claim == c.claimEpoch
				ctx := c.epochCtx
				c.mu.Unlock()
				device, exists := c.DeviceAt(request.address)
				if !current || !exists || device.RawName != request.name {
					break
				}
				value := uint64(requested)
				c.retryAdvisoryProtocolContext(ctx, fmt.Sprintf("device %d information PGN %d", request.address, requested), &pgn.IsoRequest{
					Info: pgn.MessageInfo{TargetId: pgn.Target(request.address)}, Pgn: &value,
				})
			}
		}
	}
}

type rejoinRequest struct {
	epoch uint64
	ready chan struct{}
}

func (c *Client) rejoinLoop() {
	defer c.backgroundWg.Done()
	for {
		select {
		case <-c.ctx.Done():
			return
		case request := <-c.rejoinCh:
			c.rejoinNetwork(request.epoch, request.ready)
		}
	}
}

func (c *Client) queueRejoin(request rejoinRequest) {
	select {
	case c.rejoinCh <- request:
	default:
		select {
		case <-c.rejoinCh:
		default:
		}
		select {
		case c.rejoinCh <- request:
		default:
		}
	}
}
