package n2k

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/open-ships/n2k/internal/actisense"
	"github.com/open-ships/n2k/internal/ebl"
	"github.com/open-ships/n2k/internal/gateway"
	"github.com/open-ships/n2k/raw"
)

type eblSource struct {
	path           string
	originalTiming bool
}

type contextReader struct {
	ctx    context.Context
	reader io.Reader
}

func (r contextReader) Read(buf []byte) (int, error) {
	select {
	case <-r.ctx.Done():
		return 0, r.ctx.Err()
	default:
		return r.reader.Read(buf)
	}
}

func (s *eblSource) run(ctx context.Context, _ *slog.Logger, handler func(raw.Observation)) error {
	file, err := os.Open(s.path) // #nosec G304 -- reading the caller-selected capture is this source's purpose.
	if err != nil {
		return fmt.Errorf("n2k: opening EBL capture: %w", err)
	}
	defer func() { _ = file.Close() }()

	bstParser := actisense.NewParser()
	adapter := gateway.NewActisenseObservationAdapter()
	var lastTimestamp time.Time
	emitWithContext := func(event ebl.Event, observation raw.Observation) {
		if event.HasTimestamp {
			observation.Timestamp = event.Timestamp
		}
		observation.ReceivedAt = time.Now()
		if event.Direction != ebl.DirectionUnknown {
			observation.Direction = eblDirection(event.Direction)
		}
		observation.AdapterID = "ebl:" + s.path
		observation.NetworkID = s.path
		handler(observation)
	}
	pace := func(timestamp time.Time) error {
		if !s.originalTiming || timestamp.IsZero() || lastTimestamp.IsZero() {
			if !timestamp.IsZero() {
				lastTimestamp = timestamp
			}
			return nil
		}
		delta := timestamp.Sub(lastTimestamp)
		lastTimestamp = timestamp
		if delta <= 0 {
			return nil
		}
		timer := time.NewTimer(delta)
		defer timer.Stop()
		select {
		case <-timer.C:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	reader := ebl.NewReader()
	err = reader.Read(contextReader{ctx: ctx, reader: file}, func(event ebl.Event) {
		if ctx.Err() != nil {
			return
		}
		if paceErr := pace(event.Timestamp); paceErr != nil {
			return
		}
		handleDatagram := func(datagram actisense.Datagram) {
			adapter.HandleDatagram(datagram, func(observation raw.Observation) {
				emitWithContext(event, observation)
			})
		}
		handleDecodeError := func(decodeErr actisense.DecodeError) {
			adapter.HandleDecodeError(decodeErr, func(observation raw.Observation) {
				emitWithContext(event, observation)
			})
		}
		if event.RawBST {
			datagram, decodeErr := actisense.DecodeRaw(event.Payload)
			if decodeErr != nil {
				handleDecodeError(actisense.DecodeError{Kind: actisense.DecodeLength, Err: decodeErr})
				return
			}
			handleDatagram(datagram)
			return
		}
		bstParser.Feed(event.Payload, handleDatagram, handleDecodeError)
	}, func(warning ebl.Warning) {
		now := time.Now()
		handler(raw.Observation{
			Kind:       raw.KindTransportError,
			Timestamp:  now,
			ReceivedAt: now,
			AdapterID:  "ebl:" + s.path,
			NetworkID:  s.path,
			Direction:  raw.DirectionUnknown,
			Protocol:   "ebl",
			Error:      warning.Error(),
		})
	})
	if ctx.Err() == nil {
		bstParser.End(func(decodeErr actisense.DecodeError) {
			adapter.HandleDecodeError(decodeErr, func(observation raw.Observation) {
				observation.AdapterID = "ebl:" + s.path
				observation.NetworkID = s.path
				handler(observation)
			})
		})
	}
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if err != nil {
		return fmt.Errorf("n2k: reading EBL capture: %w", err)
	}
	return nil
}

func eblDirection(direction ebl.Direction) raw.Direction {
	if direction == ebl.DirectionTransmitted {
		return raw.DirectionTransmitted
	}
	if direction == ebl.DirectionReceived {
		return raw.DirectionReceived
	}
	return raw.DirectionUnknown
}

// EBL adds an Actisense Enhanced Binary Log capture source. Both classic raw
// BDTP streams and BSTRawFrame records are decoded. OriginalTiming paces replay
// using FILETIME markers.
func EBL(path string, opts ...FileOption) Option {
	return optionFunc(func(c *config) {
		capture := captureOptions{}
		for _, option := range opts {
			if option != nil {
				option.applyCapture(&capture)
			}
		}
		c.sources = append(c.sources, &eblSource{path: path, originalTiming: capture.originalTiming})
	})
}
