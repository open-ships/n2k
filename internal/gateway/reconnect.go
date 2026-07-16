package gateway

import (
	"context"
	"time"
)

// Default backoff bounds used when a ReconnectPolicy leaves a field at zero.
const (
	defaultInitialBackoff = 500 * time.Millisecond
	defaultMaxBackoff     = 30 * time.Second
)

// ReconnectPolicy configures automatic re-dialing of a dropped gateway
// connection. A nil *ReconnectPolicy means reconnection is disabled: a
// dropped connection ends the read loop (the historical behavior).
type ReconnectPolicy struct {
	// InitialBackoff is the delay before the first reconnect attempt after a
	// drop, and the delay restored after any successful connection.
	InitialBackoff time.Duration
	// MaxBackoff caps the exponentially growing delay between attempts while a
	// connection cannot be re-established.
	MaxBackoff time.Duration
}

// Backoff is a re-armable exponential backoff sequencer. It is not safe for
// concurrent use; a single reconnect loop owns one Backoff.
type Backoff struct {
	initial, max, cur time.Duration
}

// NewBackoff returns a Backoff seeded from policy, substituting defaults for
// any non-positive field.
func NewBackoff(policy ReconnectPolicy) *Backoff {
	initial := policy.InitialBackoff
	if initial <= 0 {
		initial = defaultInitialBackoff
	}
	max := policy.MaxBackoff
	if max <= 0 {
		max = defaultMaxBackoff
	}
	if max < initial {
		max = initial
	}
	return &Backoff{initial: initial, max: max}
}

// Reset returns the sequencer to its initial delay, so the next Wait sleeps
// for InitialBackoff again. Call it after a connection is successfully
// established so a brief reconnect does not inherit a long prior backoff.
func (b *Backoff) Reset() { b.cur = 0 }

// Wait sleeps for the next backoff interval, doubling it up to MaxBackoff on
// each successive call, and returns true. It returns false immediately if ctx
// is cancelled before the interval elapses.
func (b *Backoff) Wait(ctx context.Context) bool {
	if b.cur == 0 {
		b.cur = b.initial
	} else {
		b.cur *= 2
		if b.cur > b.max {
			b.cur = b.max
		}
	}
	timer := time.NewTimer(b.cur)
	defer timer.Stop()
	select {
	case <-timer.C:
		return true
	case <-ctx.Done():
		return false
	}
}
