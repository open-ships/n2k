package gateway

import (
	"fmt"

	"github.com/brutella/can"
	"github.com/open-ships/n2k/internal/framer"
	"github.com/open-ships/n2k/pgn"
)

// Actisense stream framing bytes. Messages are wrapped as
// DLE STX cmd len payload... checksum DLE ETX, with any DLE byte between the
// STX and ETX markers escaped by doubling.
const (
	dle = 0x10
	stx = 0x02
	etx = 0x03

	// cmdN2KReceived is the Actisense command byte for an NMEA 2000 message
	// received from the bus.
	cmdN2KReceived = 0x93
)

// n2kHeaderLen is the fixed part of a command 0x93 payload: priority (1),
// PGN (3), destination (1), source (1), timestamp (4), data length (1).
const n2kHeaderLen = 11

// N2KMessage is one NMEA 2000 message carried by an Actisense stream. Data is
// the assembled PGN payload (fast-packet PGNs arrive already assembled).
type N2KMessage struct {
	Priority    uint8
	PGN         uint32
	Destination uint8
	Source      uint8
	Data        []byte
}

// ActisenseReader is an incremental decoder for the Actisense binary stream
// protocol. It tolerates arbitrary buffer boundaries: bytes are fed as they
// arrive and complete messages are emitted as they are recognized.
type ActisenseReader struct {
	inMessage bool
	escaped   bool
	body      []byte
}

// NewActisenseReader returns a reader ready to consume stream bytes.
func NewActisenseReader() *ActisenseReader {
	return &ActisenseReader{}
}

// Feed consumes raw stream bytes; emit is called once per complete,
// checksum-valid N2K message (command 0x93). Garbage between messages is
// skipped and messages that fail their checksum are dropped silently.
func (r *ActisenseReader) Feed(buf []byte, emit func(N2KMessage)) {
	for _, b := range buf {
		if !r.inMessage {
			// Hunt for DLE STX. r.escaped doubles as "previous byte was DLE".
			if r.escaped && b == stx {
				r.inMessage = true
				r.body = r.body[:0]
				r.escaped = false
			} else {
				r.escaped = b == dle
			}
			continue
		}

		if r.escaped {
			r.escaped = false
			switch b {
			case dle:
				// Escaped DLE: a literal 0x10 body byte.
				r.body = append(r.body, dle)
			case etx:
				r.finish(emit)
			case stx:
				// Unexpected restart: begin a fresh message.
				r.body = r.body[:0]
			default:
				// Malformed escape: abandon the message.
				r.inMessage = false
			}
			continue
		}

		if b == dle {
			r.escaped = true
			continue
		}
		r.body = append(r.body, b)
	}
}

// finish validates and emits the accumulated message body, then resets.
func (r *ActisenseReader) finish(emit func(N2KMessage)) {
	defer func() {
		r.inMessage = false
		r.escaped = false
		r.body = r.body[:0]
	}()

	// cmd + len + checksum is the minimum body.
	if len(r.body) < 3 {
		return
	}
	var sum byte
	for _, b := range r.body {
		sum += b
	}
	if sum != 0 {
		return
	}

	cmd, declaredLen, payload := r.body[0], int(r.body[1]), r.body[2:len(r.body)-1]
	if cmd != cmdN2KReceived || declaredLen != len(payload) {
		return
	}
	if len(payload) < n2kHeaderLen {
		return
	}

	dataLen := int(payload[10])
	if len(payload) != n2kHeaderLen+dataLen {
		return
	}

	data := make([]byte, dataLen)
	copy(data, payload[n2kHeaderLen:])
	emit(N2KMessage{
		Priority:    payload[0],
		PGN:         uint32(payload[1]) | uint32(payload[2])<<8 | uint32(payload[3])<<16,
		Destination: payload[4],
		Source:      payload[5],
		Data:        data,
	})
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
