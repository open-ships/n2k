// Package ebl reads bounded Enhanced Binary Log capture streams.
package ebl

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"time"
)

const (
	escape = 0x1B
	start  = 0x01
	end    = 0x0A

	TagVersion         = 0x01
	TagDescription     = 0x02
	TagTimeUTC         = 0x03
	TagElementType     = 0x04
	TagDirectionMarker = 0x05
	TagBSTRawFrame     = 0x07

	CurrentVersion = 1002
	MaxTagPayload  = 1800
	maxRawChunk    = 4096

	filetimeEpochOffsetSeconds = 11644473600
)

type Direction uint8

const (
	DirectionUnknown Direction = iota
	DirectionReceived
	DirectionTransmitted
)

// Event is an owned raw-stream chunk or one unframed BST record.
type Event struct {
	Timestamp    time.Time
	HasTimestamp bool
	Direction    Direction
	RawBST       bool
	Payload      []byte
}

type Warning struct {
	Offset int64
	Err    error
}

func (w Warning) Error() string { return fmt.Sprintf("EBL byte %d: %v", w.Offset, w.Err) }
func (w Warning) Unwrap() error { return w.Err }

var ErrUnsupportedVersion = errors.New("EBL version is newer than supported")

type state uint8

const (
	stateOutside state = iota
	stateOutsideEscape
	stateTag
	stateTagEscape
	statePayload
	statePayloadEscape
	stateDiscardPayload
	stateDiscardEscape
)

// Reader incrementally parses one EBL file. Retained tag and raw-stream state
// is bounded by MaxTagPayload and maxRawChunk.
type Reader struct {
	state        state
	tag          byte
	tagPayload   []byte
	rawPayload   []byte
	offset       int64
	versionSeen  bool
	timestamp    time.Time
	hasTimestamp bool
	direction    Direction
}

func NewReader() *Reader {
	return &Reader{
		tagPayload: make([]byte, 0, 64),
		rawPayload: make([]byte, 0, maxRawChunk),
	}
}

// Read consumes r to EOF. Recoverable corruption is reported and parsing
// continues; a newer declared EBL version is terminal.
func (r *Reader) Read(input io.Reader, emit func(Event), report func(Warning)) error {
	buf := make([]byte, 4096)
	for {
		n, err := input.Read(buf)
		for _, b := range buf[:n] {
			if feedErr := r.feedByte(b, emit, report); feedErr != nil {
				return feedErr
			}
			r.offset++
		}
		if err != nil {
			if errors.Is(err, io.EOF) {
				r.flushRaw(emit, report)
				if r.state != stateOutside && r.state != stateOutsideEscape {
					r.warn(report, errors.New("truncated metatag discarded"))
				}
				return nil
			}
			return err
		}
	}
}

func (r *Reader) feedByte(b byte, emit func(Event), report func(Warning)) error {
	switch r.state {
	case stateOutside:
		if b == escape {
			r.state = stateOutsideEscape
		} else {
			r.appendRaw(b, emit, report)
		}
	case stateOutsideEscape:
		switch b {
		case start:
			r.flushRaw(emit, report)
			r.tag = 0
			r.tagPayload = r.tagPayload[:0]
			r.state = stateTag
		case escape:
			r.appendRaw(escape, emit, report)
			r.state = stateOutside
		default:
			r.warn(report, fmt.Errorf("invalid stream escape 0x%02X", b))
			r.state = stateOutside
		}
	case stateTag:
		if b == escape {
			r.state = stateTagEscape
		} else {
			r.tag = b
			r.state = statePayload
		}
	case stateTagEscape:
		if b != escape {
			r.warn(report, fmt.Errorf("invalid escaped tag byte 0x%02X", b))
			r.state = stateOutside
		} else {
			r.tag = escape
			r.state = statePayload
		}
	case statePayload:
		if b == escape {
			r.state = statePayloadEscape
		} else {
			r.appendTag(b, report)
		}
	case statePayloadEscape:
		switch b {
		case escape:
			r.appendTag(escape, report)
			if r.state == statePayloadEscape {
				r.state = statePayload
			}
		case end:
			r.state = stateOutside
			return r.handleTag(emit, report)
		default:
			r.warn(report, fmt.Errorf("invalid metatag escape 0x%02X", b))
			r.tagPayload = r.tagPayload[:0]
			r.state = stateOutside
		}
	case stateDiscardPayload:
		if b == escape {
			r.state = stateDiscardEscape
		}
	case stateDiscardEscape:
		if b == end {
			r.state = stateOutside
		} else if b != escape {
			r.state = stateDiscardPayload
		} else {
			r.state = stateDiscardPayload
		}
	}
	return nil
}

func (r *Reader) appendRaw(b byte, emit func(Event), report func(Warning)) {
	if len(r.rawPayload) == maxRawChunk {
		r.flushRaw(emit, report)
	}
	r.rawPayload = append(r.rawPayload, b)
}

func (r *Reader) flushRaw(emit func(Event), report func(Warning)) {
	if len(r.rawPayload) == 0 {
		return
	}
	if !r.versionSeen {
		r.warn(report, errors.New("raw payload before Version tag discarded"))
		r.rawPayload = r.rawPayload[:0]
		return
	}
	if emit != nil {
		emit(Event{
			Timestamp: r.timestamp, HasTimestamp: r.hasTimestamp,
			Direction: r.direction,
			Payload:   append([]byte(nil), r.rawPayload...),
		})
	}
	r.rawPayload = r.rawPayload[:0]
}

func (r *Reader) appendTag(b byte, report func(Warning)) {
	if len(r.tagPayload) >= MaxTagPayload {
		r.warn(report, fmt.Errorf("metatag 0x%02X exceeds %d bytes and was discarded", r.tag, MaxTagPayload))
		r.tagPayload = r.tagPayload[:0]
		r.state = stateDiscardPayload
		return
	}
	r.tagPayload = append(r.tagPayload, b)
}

func (r *Reader) handleTag(emit func(Event), report func(Warning)) error {
	payload := r.tagPayload
	switch r.tag {
	case TagVersion:
		if len(payload) != 4 {
			r.warn(report, fmt.Errorf("version tag has %d bytes; expected 4", len(payload)))
			return nil
		}
		version := binary.LittleEndian.Uint32(payload)
		if version > CurrentVersion {
			return fmt.Errorf("%w: file=%d reader=%d", ErrUnsupportedVersion, version, CurrentVersion)
		}
		r.versionSeen = true
	case TagTimeUTC:
		if len(payload) != 8 {
			r.warn(report, fmt.Errorf("TimeUTC tag has %d bytes; expected 8", len(payload)))
			return nil
		}
		r.timestamp = filetimeToTime(binary.LittleEndian.Uint64(payload))
		r.hasTimestamp = true
	case TagDirectionMarker:
		if len(payload) != 1 {
			r.warn(report, fmt.Errorf("DirectionMarker tag has %d bytes; expected 1", len(payload)))
			return nil
		}
		switch payload[0] {
		case 0:
			r.direction = DirectionReceived
		case 1:
			r.direction = DirectionTransmitted
		default:
			r.warn(report, fmt.Errorf("invalid direction marker 0x%02X", payload[0]))
		}
	case TagBSTRawFrame:
		if !r.versionSeen {
			r.warn(report, errors.New("BSTRawFrame before Version tag discarded"))
			return nil
		}
		if len(payload) < 2 {
			r.warn(report, errors.New("BSTRawFrame shorter than two bytes discarded"))
			return nil
		}
		if emit != nil {
			emit(Event{
				Timestamp: r.timestamp, HasTimestamp: r.hasTimestamp,
				Direction: r.direction, RawBST: true,
				Payload: append([]byte(nil), payload...),
			})
		}
	case TagDescription, TagElementType:
		// Metadata does not alter replay state.
	default:
		// Unknown tags are forward-compatible and skipped.
	}
	return nil
}

func (r *Reader) warn(report func(Warning), err error) {
	if report != nil {
		report(Warning{Offset: r.offset, Err: err})
	}
}

func filetimeToTime(ticks uint64) time.Time {
	secondsSince1601 := ticks / 10_000_000
	nanoseconds := int64(ticks%10_000_000) * 100
	secondsSinceUnix := int64(secondsSince1601) - filetimeEpochOffsetSeconds
	return time.Unix(secondsSinceUnix, nanoseconds).UTC()
}
