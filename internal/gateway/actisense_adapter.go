package gateway

import (
	"errors"
	"io"
	"time"

	"github.com/brutella/can"
	"github.com/open-ships/n2k/internal/actisense"
	"github.com/open-ships/n2k/raw"
)

// ObserveActisenseDatagram delivers an owned observation at the layer the
// gateway actually supplied: BST-95 carries a frame, while BST-93/94/D0
// carry complete messages.
func ObserveActisenseDatagram(datagram actisense.Datagram, handler func(raw.Observation)) {
	if handler == nil {
		return
	}
	now := time.Now()
	if decoded, ok, err := actisense.DecodeCANFrame(datagram); ok {
		if err != nil {
			handler(actisenseTransportError(now, datagram.ID, err))
			return
		}
		handler(raw.Observation{
			Kind:                  raw.KindFrame,
			Timestamp:             now,
			ReceivedAt:            now,
			TransportTimestamp:    decoded.Timestamp,
			HasTransportTimestamp: decoded.HasTimestamp,
			AdapterID:             "actisense",
			Direction:             actisenseDirection(decoded.Direction),
			Frame:                 &decoded.Frame,
			Protocol:              "actisense-bst",
			DatagramID:            uint16(datagram.ID),
		})
		return
	}

	if message, ok, err := actisense.DecodeMessage(datagram); ok {
		if err != nil {
			handler(actisenseTransportError(now, datagram.ID, err))
			return
		}
		destination := message.Destination
		handler(raw.Observation{
			Kind:                  raw.KindMessage,
			Timestamp:             now,
			ReceivedAt:            now,
			TransportTimestamp:    message.Timestamp,
			HasTransportTimestamp: message.HasTimestamp,
			AdapterID:             "actisense",
			Direction:             actisenseDirection(message.Direction),
			PGN:                   message.PGN,
			Priority:              message.Priority,
			Source:                message.Source,
			Destination:           &destination,
			Payload:               append([]byte(nil), message.Data...),
			Protocol:              "actisense-bst",
			DatagramID:            uint16(datagram.ID),
		})
		return
	}

	if response, ok, err := actisense.DecodeBEMResponse(datagram); ok {
		if err != nil {
			handler(actisenseTransportError(now, datagram.ID, err))
			return
		}
		observation := raw.Observation{
			Kind:         raw.KindGateway,
			Timestamp:    now,
			ReceivedAt:   now,
			AdapterID:    "actisense",
			Direction:    raw.DirectionReceived,
			Payload:      append([]byte(nil), response.Data...),
			Protocol:     "actisense-bem",
			DatagramID:   uint16(datagram.ID),
			CommandID:    uint16(response.BEMID),
			Sequence:     response.Sequence,
			DeviceModel:  uint32(response.ModelID),
			DeviceSerial: response.SerialNumber,
			DeviceError:  response.ErrorCode,
		}
		if response.ErrorCode != 0 {
			observation.Error = (&actisense.DeviceError{Command: response.BEMID, Code: response.ErrorCode}).Error()
		}
		handler(observation)
		return
	}

	handler(raw.Observation{
		Kind:       raw.KindGateway,
		Timestamp:  now,
		ReceivedAt: now,
		AdapterID:  "actisense",
		Direction:  raw.DirectionReceived,
		Payload:    append([]byte(nil), datagram.Raw...),
		Protocol:   "actisense-bst",
		DatagramID: uint16(datagram.ID),
	})
}

// ObserveActisenseDecodeError delivers a transport error to the observation reader.
func ObserveActisenseDecodeError(err actisense.DecodeError, handler func(raw.Observation)) {
	if handler != nil {
		handler(actisenseTransportError(time.Now(), err.ID, err))
	}
}

func actisenseTransportError(now time.Time, datagramID byte, err error) raw.Observation {
	return raw.Observation{
		Kind:       raw.KindTransportError,
		Timestamp:  now,
		ReceivedAt: now,
		AdapterID:  "actisense",
		Direction:  raw.DirectionReceived,
		Protocol:   "actisense-bdtp",
		DatagramID: uint16(datagramID),
		Error:      err.Error(),
	}
}

func actisenseDirection(direction actisense.Direction) raw.Direction {
	if direction == actisense.DirectionTransmitted {
		return raw.DirectionTransmitted
	}
	return raw.DirectionReceived
}

// ReadActisense reassembles Actisense messages from r and emits CAN frames.
// Message-level D0 records and gateway diagnostics require the Observation
// Interface and are intentionally omitted by this compatibility helper.
func ReadActisense(r io.Reader, handler func(can.Frame)) error {
	return ReadActisenseObservations(r, actisenseFrameHandler(handler))
}

// Only legacy frame-only readers synthesize CAN frames from BST-93.
func actisenseFrameHandler(handler func(can.Frame)) func(raw.Observation) {
	if handler == nil {
		return nil
	}
	reframe := ReframeEmitter(handler)
	var epoch uint64
	return func(observation raw.Observation) {
		if observation.Frame != nil {
			handler(*observation.Frame)
			return
		}
		if observation.Kind != raw.KindMessage || observation.DatagramID != actisense.BSTN2KReceive {
			return
		}
		if epoch != observation.ConnectionEpoch {
			epoch = observation.ConnectionEpoch
			reframe = ReframeEmitter(handler)
		}
		reframe(N2KMessage{
			PGN: observation.PGN, Priority: observation.Priority, Destination: *observation.Destination,
			Source: observation.Source, Data: observation.Payload,
		})
	}
}

// ReadActisenseObservations decodes all supported BDTP/BST datagrams until EOF.
func ReadActisenseObservations(r io.Reader, handler func(raw.Observation)) error {
	parser := actisense.NewParser()
	buf := make([]byte, 4096)
	for {
		n, err := r.Read(buf)
		if n > 0 {
			parser.Feed(buf[:n],
				func(datagram actisense.Datagram) { ObserveActisenseDatagram(datagram, handler) },
				func(decodeErr actisense.DecodeError) { ObserveActisenseDecodeError(decodeErr, handler) },
			)
		}
		if err != nil {
			parser.End(func(decodeErr actisense.DecodeError) { ObserveActisenseDecodeError(decodeErr, handler) })
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}
	}
}
