package n2k

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/open-ships/n2k/pgn"
)

// defaultRequestTimeout bounds Request when the caller's context carries no
// deadline. ISO 11783-3 allows a responder 1250 ms.
const defaultRequestTimeout = 1250 * time.Millisecond

// Request sends an ISO Request (PGN 59904) for T's PGN to target and waits
// for the first matching reply. T names the expected response struct, e.g.
//
//	pi, err := n2k.Request[*pgn.ProductInformation](ctx, client, 0x23)
//
// target 255 broadcasts the request and accepts a reply from any device.
// When ctx has no deadline, a default of 1250 ms (the ISO 11783 response
// time) is applied. Request works on bus clients only.
func Request[T pgn.Message](ctx context.Context, c *Client, target uint8) (T, error) {
	var zero T

	structType := reflect.TypeFor[T]()
	if structType.Kind() != reflect.Pointer || structType.Elem().Kind() != reflect.Struct {
		return zero, fmt.Errorf("n2k: Request type %v must be a pointer to a PGN struct", structType)
	}
	proto, ok := reflect.New(structType.Elem()).Interface().(pgn.Message)
	if !ok {
		return zero, fmt.Errorf("n2k: %v does not implement pgn.Message", structType)
	}
	pgnNum := proto.PGNNumber()

	if c.system == nil {
		return zero, errors.New("n2k: Request requires a bus client (replay clients cannot receive responses)")
	}

	if _, hasDeadline := ctx.Deadline(); !hasDeadline {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, defaultRequestTimeout)
		defer cancel()
	}

	// Open the system pipeline for the response PGN and register a waiter
	// before transmitting, so a fast responder cannot slip past us.
	c.system.addPGN(pgnNum)
	defer c.system.removePGN(pgnNum)
	w := c.correlator.add(pgnNum, target)
	defer c.correlator.remove(w)

	requested := uint64(pgnNum)
	req := &pgn.IsoRequest{
		Info: pgn.MessageInfo{TargetId: pgn.Target(target)},
		Pgn:  &requested,
	}
	if err := c.Write(req).Wait(); err != nil {
		return zero, fmt.Errorf("n2k: sending ISO request for PGN %d: %w", pgnNum, err)
	}

	for {
		select {
		case msg := <-w.ch:
			// The same PGN can decode to several variants; keep waiting
			// until the requested type shows up.
			if typed, ok := msg.(T); ok {
				return typed, nil
			}
		case <-ctx.Done():
			return zero, ctx.Err()
		}
	}
}

// correlator matches decoded system messages to in-flight Request calls.
type correlator struct {
	mu      sync.Mutex
	waiters []*waiter
}

// waiter is one in-flight Request: a PGN, the address the request went to
// (255 = any responder), and a channel the response is delivered on.
type waiter struct {
	pgn    uint32
	target uint8
	ch     chan pgn.Message
}

func newCorrelator() *correlator {
	return &correlator{}
}

// add registers a waiter for a PGN from target (255 = any source).
func (co *correlator) add(pgnNum uint32, target uint8) *waiter {
	w := &waiter{pgn: pgnNum, target: target, ch: make(chan pgn.Message, 4)}
	co.mu.Lock()
	co.waiters = append(co.waiters, w)
	co.mu.Unlock()
	return w
}

// remove deregisters a waiter.
func (co *correlator) remove(w *waiter) {
	co.mu.Lock()
	defer co.mu.Unlock()
	for i, cand := range co.waiters {
		if cand == w {
			co.waiters = append(co.waiters[:i], co.waiters[i+1:]...)
			return
		}
	}
}

// observe is registered as a system router handler; it fans each decoded
// message out to every matching waiter.
func (co *correlator) observe(msg pgn.Message) {
	carrier, ok := msg.(infoCarrier)
	if !ok {
		return
	}
	info := carrier.MessageInfo()

	co.mu.Lock()
	defer co.mu.Unlock()
	for _, w := range co.waiters {
		if w.pgn != info.PGN {
			continue
		}
		if w.target != 255 && w.target != info.SourceId {
			continue
		}
		select {
		case w.ch <- msg:
		default: // waiter's buffer is full; it has what it needs
		}
	}
}
