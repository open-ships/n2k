package gateway

import (
	"fmt"

	"github.com/brutella/can"
	"github.com/open-ships/n2k/internal/actisense"
	"github.com/open-ships/n2k/internal/framer"
	"github.com/open-ships/n2k/pgn"
)

// Actisense stream framing bytes. Messages are wrapped as
// DLE STX cmd len payload... checksum DLE ETX, with any DLE byte between the
// STX and ETX markers escaped by doubling.
const (
	dle = actisense.DLE
	stx = actisense.STX
	etx = actisense.ETX

	// cmdN2KReceived is the Actisense command byte for an NMEA 2000 message
	// received from the bus.
	cmdN2KReceived = actisense.BSTN2KReceive
	// cmdN2KSend is the Actisense command byte for transmitting an NMEA 2000
	// message onto the bus.
	cmdN2KSend = actisense.BSTN2KTransmit
)

// N2KMessage is one NMEA 2000 message carried by an Actisense stream. Data is
// the assembled PGN payload (fast-packet PGNs arrive already assembled).
type N2KMessage = actisense.Message

// ActisenseReader is an incremental decoder for the Actisense binary stream
// protocol. It tolerates arbitrary buffer boundaries: bytes are fed as they
// arrive and complete messages are emitted as they are recognized.
type ActisenseReader struct {
	parser *actisense.Parser
}

// Maximum unescaped body includes the bounded BST record and checksum.
const maxActisenseBody = actisense.MaxDatagramLength + 1

// NewActisenseReader returns a reader ready to consume stream bytes.
func NewActisenseReader() *ActisenseReader {
	return &ActisenseReader{parser: actisense.NewParser()}
}

// Feed consumes raw stream bytes; emit is called once per complete,
// checksum-valid N2K message (command 0x93). Garbage between messages is
// skipped and messages that fail their checksum are dropped silently.
func (r *ActisenseReader) Feed(buf []byte, emit func(N2KMessage)) {
	if r.parser == nil {
		r.parser = actisense.NewParser()
	}
	r.parser.Feed(buf, func(datagram actisense.Datagram) {
		if datagram.ID != actisense.BSTN2KReceive {
			return
		}
		message, ok, err := actisense.DecodeMessage(datagram)
		if err == nil && ok && emit != nil {
			emit(message)
		}
	}, nil)
}

// EncodeSend renders an assembled N2K message as an Actisense transmit
// command (0x94). The payload layout is priority (1), PGN (3, little-endian),
// destination (1), data length (1), data. There is no source-address byte:
// the gateway transmits under its own claimed bus address, so m.Source is
// not carried on the wire. The gateway also performs fast-packet
// fragmentation itself, which is why sends are whole messages rather than
// CAN frames.
func EncodeSend(m N2KMessage) ([]byte, error) {
	return actisense.EncodeMessage94(m)
}

// EncodeStartup renders the historical unacknowledged mode-2 command. Mode 2
// disables the receive PGN list only; the transmit list remains active. New
// session adapters use an acknowledged BEM handshake instead.
func EncodeStartup() []byte {
	buf, _ := actisense.EncodeBEMRequest(actisense.BEMOperatingMode, actisense.OperatingModeSet(actisense.ModeTransferReceiveAll))
	return buf
}

// ReframeEmitter adapts a CAN-frame handler into an assembled-message
// consumer: each message is re-framed into wire CAN frames (rotating the
// fast-packet sequence ID per message) and fed to the handler. Messages that
// cannot be re-framed are dropped.
func ReframeEmitter(handler func(can.Frame)) func(N2KMessage) {
	var seq uint8
	return func(m N2KMessage) {
		frames, err := Reframe(m, seq)
		if err != nil {
			return
		}
		seq = (seq + 1) % 8
		for _, f := range frames {
			handler(f)
		}
	}
}

// Reframe converts an assembled N2K message back into wire CAN frames:
// a single frame for non-fast PGNs up to 8 bytes, fast-packet frames for
// fast PGNs up to 223 bytes. PGNs absent from the metadata tables are framed
// by size (fast iff the payload exceeds 8 bytes). seq is the fast-packet
// sequence ID (0-7).
func Reframe(m N2KMessage, seq uint8) ([]can.Frame, error) {
	if len(m.Data) > 223 {
		return nil, fmt.Errorf("gateway: PGN %d payload is %d bytes; the fast-packet maximum is 223", m.PGN, len(m.Data))
	}

	fast := len(m.Data) > 8
	if infos, ok := pgn.PgnInfoLookup[m.PGN]; ok && len(infos) > 0 {
		fast = infos[0].Fast
	}

	canID := framer.BuildCANID(m.PGN, m.Priority, m.Source, m.Destination)
	if !fast {
		if len(m.Data) > 8 {
			return nil, fmt.Errorf("gateway: PGN %d is single-frame but payload is %d bytes", m.PGN, len(m.Data))
		}
		return []can.Frame{framer.FrameSingle(canID, m.Data)}, nil
	}
	return framer.FrameFastPacket(canID, m.Data, seq%8), nil
}
