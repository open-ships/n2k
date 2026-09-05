// Package raw defines owned observations from NMEA 2000 transport adapters.
package raw

import (
	"time"

	"github.com/brutella/can"
)

// Kind identifies the layer represented by an Observation.
type Kind string

const (
	KindFrame          Kind = "frame"
	KindMessage        Kind = "message"
	KindDecodeError    Kind = "decode_error"
	KindGateway        Kind = "gateway"
	KindTransportError Kind = "transport_error"
)

// Direction identifies how an observation crossed its Adapter.
type Direction string

const (
	DirectionUnknown     Direction = "unknown"
	DirectionReceived    Direction = "received"
	DirectionTransmitted Direction = "transmitted"
)

// Observation is an owned raw record from a transport Adapter or decode
// pipeline. Frame observations retain the classical CAN frame; message
// observations retain an assembled PGN payload; gateway and transport-error
// observations retain bounded protocol context and stable error text.
type Observation struct {
	Kind Kind `json:"kind"`

	// Timestamp is the best source timestamp available. ReceivedAt is always
	// the host receipt time. When a gateway exposes only a relative clock,
	// TransportTimestamp preserves it without pretending it is wall time.
	Timestamp             time.Time     `json:"timestamp"`
	ReceivedAt            time.Time     `json:"receivedAt"`
	TransportTimestamp    time.Duration `json:"transportTimestamp,omitempty"`
	HasTransportTimestamp bool          `json:"hasTransportTimestamp,omitempty"`

	AdapterID string    `json:"adapterId"`
	NetworkID string    `json:"networkId"`
	Direction Direction `json:"direction"`
	// ConnectionEpoch and ClaimEpoch identify the live Client session that
	// observed this record. Capture sources may leave both zero.
	ConnectionEpoch uint64 `json:"connectionEpoch,omitempty"`
	ClaimEpoch      uint64 `json:"claimEpoch,omitempty"`

	Frame *can.Frame `json:"frame,omitempty"`

	PGN         uint32 `json:"pgn,omitempty"`
	Priority    uint8  `json:"priority,omitempty"`
	Source      uint8  `json:"source,omitempty"`
	Destination *uint8 `json:"destination,omitempty"`
	Payload     []byte `json:"payload,omitempty"`
	Error       string `json:"error,omitempty"`

	// Protocol fields retain gateway control-plane and framing diagnostics
	// without coupling consumers to a concrete Adapter implementation.
	Protocol     string `json:"protocol,omitempty"`
	DatagramID   uint16 `json:"datagramId,omitempty"`
	CommandID    uint16 `json:"commandId,omitempty"`
	Sequence     uint8  `json:"sequence,omitempty"`
	DeviceModel  uint32 `json:"deviceModel,omitempty"`
	DeviceSerial uint32 `json:"deviceSerial,omitempty"`
	DeviceError  int32  `json:"deviceError,omitempty"`
}

// Clone returns a fully owned copy safe for another goroutine or caller to
// retain and mutate.
func (o Observation) Clone() Observation {
	copy := o
	if o.Frame != nil {
		frame := *o.Frame
		copy.Frame = &frame
	}
	if o.Destination != nil {
		destination := *o.Destination
		copy.Destination = &destination
	}
	copy.Payload = append([]byte(nil), o.Payload...)
	return copy
}
