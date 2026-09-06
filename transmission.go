package n2k

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/brutella/can"
)

const defaultWriteTimeout = time.Second

type writeStampKey struct{}

type writeStamp struct {
	connection uint64
	claim      uint64
	protocol   bool
	progress   *writeProgress
}

type writeProgress struct {
	attempted atomic.Uint64
	completed atomic.Uint64
}

// WriteError retains transmission uncertainty. CompletedRecords counts atomic
// transport records: CAN frames for a Bus, whole messages for MessageWriter.
// It does not imply application-level acceptance by a remote device.
// An error after any physical write attempt is conservative:
// callers must decide whether resending is appropriate.
type WriteError struct {
	Err                   error
	CompletedRecords      uint64
	TransmissionUncertain bool
}

func (e *WriteError) Error() string { return e.Err.Error() }
func (e *WriteError) Unwrap() error { return e.Err }

type wireJob struct {
	ctx   context.Context
	write func(context.Context) error
	done  chan error
}

// wireTransmitter serializes atomic wire records. Protocol records are selected
// first between application frames, including while ISO transfers await peers
// or their next pacing deadline.
type wireTransmitter struct {
	mu          sync.Mutex
	closed      bool
	client      *Client
	required    chan wireJob
	application chan wireJob
	done        chan struct{}
}

func newWireTransmitter(c *Client) *wireTransmitter {
	tx := &wireTransmitter{client: c, required: make(chan wireJob, 64), application: make(chan wireJob, 64), done: make(chan struct{})}
	go tx.run()
	return tx
}

func (tx *wireTransmitter) send(ctx context.Context, write func(context.Context) error) error {
	job := wireJob{ctx: ctx, write: write, done: make(chan error, 1)}
	stamp, _ := ctx.Value(writeStampKey{}).(writeStamp)
	queue := tx.required
	if !stamp.protocol {
		queue = tx.application
	}
	if err := context.Cause(ctx); err != nil {
		return err
	}
	tx.mu.Lock()
	if tx.closed {
		tx.mu.Unlock()
		return tx.client.operationError()
	}
	select {
	case queue <- job:
	default:
		tx.mu.Unlock()
		if stamp.protocol {
			return ErrProtocolQueueFull
		}
		return ErrWriteQueueFull
	}
	tx.mu.Unlock()
	// Completion owns the physical call as well as cancellation. In particular,
	// accepted-frame accounting cannot be finalized while I/O is unwinding.
	return <-job.done
}

func (tx *wireTransmitter) run() {
	defer close(tx.done)
	defer func() {
		tx.mu.Lock()
		defer tx.mu.Unlock()
		tx.closed = true
		for {
			select {
			case job := <-tx.required:
				job.done <- tx.client.operationError()
			case job := <-tx.application:
				job.done <- tx.client.operationError()
			default:
				return
			}
		}
	}()
	for {
		if tx.client.ctx.Err() != nil {
			return
		}
		var job wireJob
		select {
		case job = <-tx.required:
		default:
			select {
			case job = <-tx.required:
			case job = <-tx.application:
			case <-tx.client.ctx.Done():
				return
			}
		}
		err := tx.execute(job)
		job.done <- err
		var busErr *runtimeBusError
		if errors.As(err, &busErr) && tx.client.ctx.Err() == nil {
			tx.client.fail(err)
		}
	}
}

func (tx *wireTransmitter) execute(job wireJob) (err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			err = wrapRuntimeBusError(fmt.Errorf("panic writing message: %v", recovered))
		}
	}()
	if err := tx.client.checkWriteContext(job.ctx); err != nil {
		return err
	}
	timeout := defaultWriteTimeout
	if tx.client.cfg.writeTimeout != nil {
		timeout = *tx.client.cfg.writeTimeout
	}
	ctx, cancel := context.WithTimeout(job.ctx, timeout)
	defer cancel()
	stamp, _ := ctx.Value(writeStampKey{}).(writeStamp)
	if stamp.progress != nil {
		stamp.progress.attempted.Add(1)
	}
	err = job.write(ctx)
	if err == nil && stamp.progress != nil {
		stamp.progress.completed.Add(1)
	}
	if cause := context.Cause(job.ctx); cause != nil {
		return cause
	}
	if cause := context.Cause(ctx); cause != nil {
		return wrapRuntimeBusError(cause)
	}
	return wrapRuntimeBusError(err)
}

func (c *Client) writeFrameContext(ctx context.Context, frame can.Frame) error {
	if c.wire == nil {
		return c.writeFrame(frame)
	}
	ctx, stop := c.stampProtocolContext(ctx)
	defer stop()
	return c.wire.send(ctx, func(ctx context.Context) error {
		var err error
		if bus, ok := c.bus.(ContextBus); ok {
			err = bus.WriteFrameContext(ctx, frame)
		} else {
			// The Bus contract requires Close to interrupt blocked I/O. Wait for
			// a started cancellation callback before reusing the transport.
			closed := make(chan struct{})
			interrupt := context.AfterFunc(ctx, func() { _ = c.bus.Close(); close(closed) })
			defer func() {
				if !interrupt() {
					<-closed
				}
			}()
			err = c.bus.WriteFrame(frame)
		}
		if err == nil {
			observation := frameObservation(frame, "client", "bus", DirectionTransmitted)
			stamp, _ := ctx.Value(writeStampKey{}).(writeStamp)
			observation.ConnectionEpoch, observation.ClaimEpoch = stamp.connection, stamp.claim
			c.publishObservation(observation)
		}
		return err
	})
}
