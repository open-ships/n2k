package n2k

import (
	"context"
	"errors"
	"iter"
	"log/slog"
	"sync"
	"time"

	"github.com/brutella/can"
	"github.com/open-ships/n2k/internal/adapter"
	"github.com/open-ships/n2k/internal/decoder"
	"github.com/open-ships/n2k/pgn"
	"github.com/open-ships/n2k/raw"
)

// Public aliases keep the common observation surface in package n2k while
// allowing transport Adapter implementations to share the cycle-free raw
// package.
type Observation = raw.Observation
type ObservationKind = raw.Kind
type Direction = raw.Direction

const (
	ObservationFrame       = raw.KindFrame
	ObservationMessage     = raw.KindMessage
	ObservationDecodeError = raw.KindDecodeError

	DirectionUnknown     = raw.DirectionUnknown
	DirectionReceived    = raw.DirectionReceived
	DirectionTransmitted = raw.DirectionTransmitted
)

func normalizeObservation(observation raw.Observation) raw.Observation {
	observation = observation.Clone()
	if observation.ReceivedAt.IsZero() {
		observation.ReceivedAt = time.Now()
	}
	if observation.Timestamp.IsZero() {
		observation.Timestamp = observation.ReceivedAt
	}
	if observation.Direction == "" {
		observation.Direction = raw.DirectionUnknown
	}
	if observation.Kind == "" {
		observation.Kind = raw.KindFrame
	}
	if observation.Frame != nil {
		id := adapter.NewPacketInfoAt(observation.Frame, observation.Timestamp)
		observation.PGN = id.PGN
		if id.Priority != nil {
			observation.Priority = *id.Priority
		}
		observation.Source = id.SourceId
		observation.Destination = id.TargetId
	}
	return observation
}

// Observe returns a bounded iterator of owned transport observations from the
// configured read-only sources. Its bounded channel applies backpressure to
// the source, which makes file capture and replay lossless without unbounded
// memory growth. Use Client.Observations for a writable Client; that live
// protocol path fails only the slow subscriber with ErrObservationOverflow.
func Observe(ctx context.Context, opts ...Option) iter.Seq2[Observation, error] {
	return func(yield func(Observation, error) bool) {
		cfg := config{}
		var optionErr error
		for _, option := range opts {
			if option == nil {
				optionErr = errors.New("n2k: nil Option")
				break
			}
			option.apply(&cfg)
		}
		cfg.applyReconnect()
		if cfg.logger == nil {
			cfg.logger = slog.Default()
		}
		if optionErr == nil {
			optionErr = cfg.validate()
		}
		if optionErr == nil && len(cfg.sources) == 0 {
			optionErr = errors.New("n2k: Observe requires at least one source option")
		}
		if optionErr != nil {
			yield(Observation{}, optionErr)
			return
		}

		observeCtx, cancel := context.WithCancel(ctx)
		defer cancel()
		buffer := defaultReceiveBuffer
		if cfg.receiveBuffer != nil {
			buffer = *cfg.receiveBuffer
		}
		observations := make(chan raw.Observation, buffer)
		var (
			errMu sync.Mutex
			err   error
		)
		setError := func(value error) {
			errMu.Lock()
			if err == nil {
				err = value
			}
			errMu.Unlock()
		}
		go func() {
			defer close(observations)
			runErr := runSources(observeCtx, cfg.logger, cfg.sources, func(observation raw.Observation) {
				select {
				case observations <- normalizeObservation(observation):
				case <-observeCtx.Done():
				}
			})
			if runErr != nil && !errors.Is(runErr, context.Canceled) {
				setError(runErr)
			}
		}()

		for observation := range observations {
			if !yield(observation.Clone(), nil) {
				return
			}
		}
		errMu.Lock()
		terminalErr := err
		errMu.Unlock()
		if terminalErr != nil {
			yield(Observation{}, terminalErr)
		}
	}
}

func messageInfoForObservation(observation raw.Observation) pgn.MessageInfo {
	info := adapter.NewPacketInfoAt(observation.Frame, observation.Timestamp)
	info.ReceivedAt = observation.ReceivedAt
	info.TransportTimestamp = observation.TransportTimestamp
	info.HasTransportTimestamp = observation.HasTransportTimestamp
	info.AdapterID = observation.AdapterID
	info.NetworkID = observation.NetworkID
	info.Direction = observation.Direction
	return info
}

func observationFromPacket(kind raw.Kind, packet decoder.Packet, reason string) raw.Observation {
	info := packet.Info
	observation := raw.Observation{
		Kind:                  kind,
		Timestamp:             info.Timestamp,
		ReceivedAt:            info.ReceivedAt,
		TransportTimestamp:    info.TransportTimestamp,
		HasTransportTimestamp: info.HasTransportTimestamp,
		AdapterID:             info.AdapterID,
		NetworkID:             info.NetworkID,
		Direction:             info.Direction,
		PGN:                   info.PGN,
		Source:                info.SourceId,
		Payload:               append([]byte(nil), packet.Data...),
		Error:                 reason,
	}
	if info.Priority != nil {
		observation.Priority = *info.Priority
	}
	if info.TargetId != nil {
		destination := *info.TargetId
		observation.Destination = &destination
	}
	return observation
}

func frameObservation(frame can.Frame, adapterID, networkID string, direction raw.Direction) raw.Observation {
	now := time.Now()
	return normalizeObservation(raw.Observation{
		Kind:       raw.KindFrame,
		Timestamp:  now,
		ReceivedAt: now,
		AdapterID:  adapterID,
		NetworkID:  networkID,
		Direction:  direction,
		Frame:      &frame,
	})
}

func (c *Client) publishObservation(observation raw.Observation) {
	if c == nil || c.observationHub == nil {
		return
	}
	switch observation.Kind {
	case raw.KindFrame:
		if observation.Direction == raw.DirectionTransmitted {
			c.framesTransmitted.Add(1)
		} else {
			c.framesReceived.Add(1)
		}
	case raw.KindMessage:
		c.messagesObserved.Add(1)
	case raw.KindDecodeError:
		c.decodeErrorsObserved.Add(1)
	}
	c.observationHub.publish(observation)
}

// Observations returns a bounded iterator of raw observations for a live
// client. Frame events are published before protocol handling; assembled
// message and decode-error events follow as the user pipeline processes them.
func (c *Client) Observations() iter.Seq2[Observation, error] {
	return func(yield func(Observation, error) bool) {
		if c == nil || c.observationHub == nil {
			yield(Observation{}, ErrClientClosed)
			return
		}
		subscription := c.observationHub.subscribe()
		defer subscription.unsubscribe()
		for observation := range subscription.ch {
			if !yield(observation.Clone(), nil) {
				return
			}
		}
		if err := subscription.terminalError(); err != nil {
			yield(Observation{}, err)
		}
	}
}
