package n2k

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync/atomic"
	"time"

	"github.com/open-ships/n2k/pgn"
)

const (
	defaultRequiredProtocolQueue = 32
	defaultAdvisoryProtocolQueue = 64
	protocolRetryAttempts        = 4
)

// ErrProtocolQueueFull reports exhaustion of a bounded automatic-protocol
// transmission lane.
var ErrProtocolQueueFull = errors.New("n2k: protocol write queue full")

type protocolWriteClass uint8

const (
	protocolRequired protocolWriteClass = iota
	protocolAdvisory
)

// protocolTransmitter owns bounded admission and operational counters for
// automatic network traffic. Required traffic has a dedicated priority lane;
// discovery/enumeration traffic uses a separately bounded advisory lane.
type protocolTransmitter struct {
	required chan writeJob
	advisory chan writeJob
	log      *slog.Logger

	accepted  atomic.Uint64
	completed atomic.Uint64
	failed    atomic.Uint64
	rejected  atomic.Uint64
}

func newProtocolTransmitter(log *slog.Logger) *protocolTransmitter {
	return &protocolTransmitter{
		required: make(chan writeJob, defaultRequiredProtocolQueue),
		advisory: make(chan writeJob, defaultAdvisoryProtocolQueue),
		log:      log,
	}
}

func (tx *protocolTransmitter) admit(job writeJob) error {
	var queue chan writeJob
	switch job.protocolClass {
	case protocolRequired:
		queue = tx.required
	case protocolAdvisory:
		queue = tx.advisory
	default:
		return fmt.Errorf("unknown protocol write class %d", job.protocolClass)
	}
	select {
	case queue <- job:
		tx.accepted.Add(1)
		return nil
	default:
		tx.rejected.Add(1)
		return fmt.Errorf("%w: %s", ErrProtocolQueueFull, job.operation)
	}
}

func (tx *protocolTransmitter) takeReady() (writeJob, bool) {
	select {
	case job := <-tx.required:
		return job, true
	default:
	}
	select {
	case job := <-tx.advisory:
		return job, true
	default:
		return writeJob{}, false
	}
}

func (tx *protocolTransmitter) finish(err error) {
	if err != nil {
		tx.failed.Add(1)
		return
	}
	tx.completed.Add(1)
}

func (c *Client) writeProtocol(operation string, class protocolWriteClass, msg pgn.Message) *WriteResult {
	result := newWriteResult()
	job := writeJob{
		msg:           msg,
		result:        result,
		protocol:      true,
		protocolClass: class,
		operation:     operation,
	}

	c.mu.Lock()
	if c.closed {
		c.mu.Unlock()
		result.complete(ErrClientClosed)
		return result
	}
	if c.terminalErr != nil {
		err := c.terminalErr
		c.mu.Unlock()
		result.complete(err)
		return result
	}
	err := c.protocolTx.admit(job)
	c.mu.Unlock()
	if err == nil {
		return result
	}
	result.complete(err)
	if class == protocolRequired {
		c.fail(fmt.Errorf("n2k: required protocol transmission %s rejected: %w", operation, err))
	}
	return result
}

// retryAdvisoryProtocol retries only bounded admission failures. Once a job
// reaches the writer, any encoding or bus failure is terminal and must not be
// duplicated by retrying an uncertain transmission.
func (c *Client) retryAdvisoryProtocol(operation string, msg pgn.Message) {
	backoff := 25 * time.Millisecond
	for attempt := 1; attempt <= protocolRetryAttempts; attempt++ {
		err := c.writeProtocol(operation, protocolAdvisory, msg).WaitContext(c.ctx)
		if err == nil || c.ctx.Err() != nil {
			return
		}
		if !errors.Is(err, ErrProtocolQueueFull) {
			return
		}
		if attempt == protocolRetryAttempts {
			c.log.Warn("advisory protocol transmission exhausted retries",
				"operation", operation, "attempts", protocolRetryAttempts, "error", err)
			return
		}
		timer := time.NewTimer(backoff)
		select {
		case <-timer.C:
			backoff *= 2
		case <-c.ctx.Done():
			timer.Stop()
			return
		}
	}
}

func (c *Client) drainWriteQueues(err error) {
	for {
		drained := false
		if c.protocolTx != nil {
			if job, ok := c.protocolTx.takeReady(); ok {
				job.result.complete(err)
				drained = true
			}
		}
		select {
		case job := <-c.writeCh:
			job.result.complete(err)
			drained = true
		default:
		}
		if !drained {
			return
		}
	}
}

// waitForWriteJob selects the next job while giving already-queued required
// and advisory protocol traffic priority over application traffic.
func (c *Client) waitForWriteJob(ctx context.Context) (writeJob, bool) {
	if c.protocolTx != nil {
		if job, ok := c.protocolTx.takeReady(); ok {
			return job, true
		}
		select {
		case job := <-c.protocolTx.required:
			return job, true
		case job := <-c.protocolTx.advisory:
			return job, true
		case job := <-c.writeCh:
			return job, true
		case <-ctx.Done():
			return writeJob{}, false
		}
	}
	select {
	case job := <-c.writeCh:
		return job, true
	case <-ctx.Done():
		return writeJob{}, false
	}
}
