// Package actisense implements the bounded BDTP, BST, BEM, and EBL wire
// protocols used by Actisense-format gateways. Transport adapters live in
// internal/gateway; this package owns protocol framing and session state.
package actisense

import (
	"encoding/binary"
	"errors"
	"fmt"
)

const (
	DLE = 0x10
	STX = 0x02
	ETX = 0x03

	// MaxDatagramLength bounds an unframed BST datagram excluding its checksum.
	// It covers a 1785-byte NMEA 2000 transport payload plus the 13-byte BST-D0
	// header, while remaining within the EBL tag capacity used by field captures.
	MaxDatagramLength = 1800
)

// DecodeErrorKind identifies a recoverable BDTP/BST stream defect.
type DecodeErrorKind string

const (
	DecodeFraming  DecodeErrorKind = "framing"
	DecodeChecksum DecodeErrorKind = "checksum"
	DecodeLength   DecodeErrorKind = "length"
	DecodeOversize DecodeErrorKind = "oversize"
)

// DecodeError describes one discarded datagram. A Parser remains usable after
// reporting it and continues hunting for the next DLE/STX marker.
type DecodeError struct {
	Kind DecodeErrorKind
	ID   byte
	Err  error
}

func (e DecodeError) Error() string {
	if e.ID != 0 {
		return fmt.Sprintf("actisense: %s error for BST 0x%02X: %v", e.Kind, e.ID, e.Err)
	}
	return fmt.Sprintf("actisense: %s error: %v", e.Kind, e.Err)
}

func (e DecodeError) Unwrap() error { return e.Err }

// Datagram is an owned, checksum-validated BST record. Raw contains the BST
// ID, length field, and payload without BDTP markers or the checksum byte.
type Datagram struct {
	ID      byte
	Payload []byte
	Raw     []byte
}

func (d Datagram) Clone() Datagram {
	return Datagram{
		ID:      d.ID,
		Payload: append([]byte(nil), d.Payload...),
		Raw:     append([]byte(nil), d.Raw...),
	}
}

// Parser incrementally decodes a BDTP byte stream. Its retained state is
// bounded by MaxDatagramLength plus one checksum byte.
type Parser struct {
	inFrame       bool
	escaped       bool
	body          []byte
	dropped       bool
	unframedBytes uint64
}

// NewParser returns a parser ready for one connection epoch.
func NewParser() *Parser {
	return &Parser{body: make([]byte, 0, 258)}
}

// BufferedBytes reports retained partial-frame bytes for diagnostics and
// bounded-state tests.
func (p *Parser) BufferedBytes() int { return len(p.body) }

// UnframedBytes reports how many bytes were discarded while hunting for a
// DLE/STX marker. It includes a pending DLE flushed by End.
func (p *Parser) UnframedBytes() uint64 { return p.unframedBytes }

// End reports and discards an incomplete candidate at an input boundary such
// as connection close or end-of-file.
func (p *Parser) End(report func(DecodeError)) {
	p.EndWithUnframed(report, nil)
}

// EndWithUnframed is End with exact delivery of a pending unframed DLE.
func (p *Parser) EndWithUnframed(report func(DecodeError), unframed func([]byte)) {
	if p.inFrame || len(p.body) > 0 {
		if report != nil {
			report(DecodeError{Kind: DecodeFraming, ID: p.id(), Err: errors.New("truncated datagram")})
		}
	} else if p.escaped {
		p.emitUnframed(DLE, unframed)
	}
	p.reset()
}

// Feed consumes arbitrary stream chunks. Complete datagrams and recoverable
// errors are delivered synchronously in wire order.
func (p *Parser) Feed(buf []byte, emit func(Datagram), report func(DecodeError)) {
	p.FeedWithUnframed(buf, emit, report, nil)
}

// FeedWithUnframed additionally delivers exact bytes found outside BDTP
// frames. This supports boot-banner evidence and mode-6 binary/ASCII demux.
// The callback must not retain its one-byte slice.
func (p *Parser) FeedWithUnframed(buf []byte, emit func(Datagram), report func(DecodeError), unframed func([]byte)) {
	for _, b := range buf {
		if !p.inFrame {
			if p.escaped {
				if b == STX {
					p.inFrame = true
					p.body = p.body[:0]
					p.escaped = false
					p.dropped = false
					continue
				}
				p.emitUnframed(DLE, unframed)
				p.escaped = false
			}
			if b == DLE {
				p.escaped = true
			} else {
				p.emitUnframed(b, unframed)
			}
			continue
		}

		if p.escaped {
			p.escaped = false
			switch b {
			case DLE:
				p.appendBody(DLE, report)
			case ETX:
				p.finish(emit, report)
			case STX:
				// A fresh start marker resynchronizes an interrupted frame.
				if len(p.body) > 0 && report != nil {
					report(DecodeError{Kind: DecodeFraming, ID: p.id(), Err: errors.New("unexpected frame restart")})
				}
				p.body = p.body[:0]
				p.dropped = false
			default:
				if report != nil {
					report(DecodeError{Kind: DecodeFraming, ID: p.id(), Err: fmt.Errorf("invalid DLE escape 0x%02X", b)})
				}
				p.reset()
			}
			continue
		}

		if b == DLE {
			p.escaped = true
			continue
		}
		p.appendBody(b, report)
	}
}

func (p *Parser) emitUnframed(value byte, emit func([]byte)) {
	p.unframedBytes++
	if emit != nil {
		one := [1]byte{value}
		emit(one[:])
	}
}

func (p *Parser) appendBody(b byte, report func(DecodeError)) {
	if p.dropped {
		return
	}
	if len(p.body) >= MaxDatagramLength+1 {
		if report != nil {
			report(DecodeError{Kind: DecodeOversize, ID: p.id(), Err: fmt.Errorf("datagram exceeds %d bytes", MaxDatagramLength)})
		}
		p.body = p.body[:0]
		p.dropped = true
		return
	}
	p.body = append(p.body, b)
}

func (p *Parser) id() byte {
	if len(p.body) == 0 {
		return 0
	}
	return p.body[0]
}

func (p *Parser) reset() {
	p.inFrame = false
	p.escaped = false
	p.dropped = false
	p.body = p.body[:0]
}

func (p *Parser) finish(emit func(Datagram), report func(DecodeError)) {
	defer p.reset()
	if p.dropped {
		return
	}
	if len(p.body) < 3 {
		if report != nil {
			report(DecodeError{Kind: DecodeLength, ID: p.id(), Err: errors.New("datagram shorter than header and checksum")})
		}
		return
	}
	var sum byte
	for _, b := range p.body {
		sum += b
	}
	if sum != 0 {
		if report != nil {
			report(DecodeError{Kind: DecodeChecksum, ID: p.id(), Err: fmt.Errorf("sum is 0x%02X, want 0", sum)})
		}
		return
	}
	d, err := decodeRaw(p.body[:len(p.body)-1])
	if err != nil {
		if report != nil {
			report(DecodeError{Kind: DecodeLength, ID: p.id(), Err: err})
		}
		return
	}
	if emit != nil {
		emit(d)
	}
}

// DecodeRaw decodes one unframed BST datagram without a checksum. This is the
// representation stored in EBL BSTRawFrame records.
func DecodeRaw(raw []byte) (Datagram, error) {
	return decodeRaw(raw)
}

func decodeRaw(raw []byte) (Datagram, error) {
	if len(raw) < 2 {
		return Datagram{}, errors.New("BST datagram shorter than two bytes")
	}
	if len(raw) > MaxDatagramLength {
		return Datagram{}, fmt.Errorf("BST datagram is %d bytes; maximum is %d", len(raw), MaxDatagramLength)
	}
	id := raw[0]
	headerLen := 2
	if IsType2(id) {
		if len(raw) < 3 {
			return Datagram{}, errors.New("BST Type 2 datagram shorter than three-byte header")
		}
		headerLen = 3
		declared := int(binary.LittleEndian.Uint16(raw[1:3]))
		if declared != len(raw) {
			return Datagram{}, fmt.Errorf("BST Type 2 length is %d, actual length is %d", declared, len(raw))
		}
	} else if int(raw[1]) != len(raw)-2 {
		return Datagram{}, fmt.Errorf("BST Type 1 payload length is %d, actual length is %d", raw[1], len(raw)-2)
	}
	owned := append([]byte(nil), raw...)
	return Datagram{ID: id, Payload: owned[headerLen:], Raw: owned}, nil
}

// IsType2 reports whether id uses the BST Type 2 16-bit total-length field.
func IsType2(id byte) bool { return id >= 0xD0 && id <= 0xDF }

// EncodeDatagram wraps a BST ID and payload in BDTP framing with the applicable
// Type 1 or Type 2 length field and zero-sum checksum.
func EncodeDatagram(id byte, payload []byte) ([]byte, error) {
	var raw []byte
	if IsType2(id) {
		total := 3 + len(payload)
		if total > MaxDatagramLength {
			return nil, fmt.Errorf("actisense: BST Type 2 datagram is %d bytes; maximum is %d", total, MaxDatagramLength)
		}
		raw = make([]byte, 3, total+1)
		raw[0] = id
		binary.LittleEndian.PutUint16(raw[1:3], uint16(total))
	} else {
		if len(payload) > 255 {
			return nil, fmt.Errorf("actisense: BST Type 1 payload is %d bytes; maximum is 255", len(payload))
		}
		raw = []byte{id, byte(len(payload))}
	}
	raw = append(raw, payload...)
	var sum byte
	for _, b := range raw {
		sum += b
	}
	raw = append(raw, -sum)

	out := make([]byte, 0, len(raw)+6)
	out = append(out, DLE, STX)
	for _, b := range raw {
		out = append(out, b)
		if b == DLE {
			out = append(out, DLE)
		}
	}
	return append(out, DLE, ETX), nil
}
