package n2k

import (
	"context"
	"errors"
	"fmt"

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

// invalidateEpochLocked owns pending operation lifetimes. Protocol parsers are
// reset after releasing c.mu, avoiding callback/Client lock inversions.
func (c *Client) invalidateEpochLocked(err error) {
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
