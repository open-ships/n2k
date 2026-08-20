package n2k

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/open-ships/n2k/internal/actisense"
	"github.com/open-ships/n2k/internal/ebl"
)

type eblWriterConfig struct {
	start       time.Time
	description string
}

// EBLWriterOption configures a new Enhanced Binary Log writer.
type EBLWriterOption interface{ applyEBLWriter(*eblWriterConfig) }

type eblWriterOptionFunc func(*eblWriterConfig)

func (f eblWriterOptionFunc) applyEBLWriter(config *eblWriterConfig) { f(config) }

func WithEBLDescription(description string) EBLWriterOption {
	return eblWriterOptionFunc(func(config *eblWriterConfig) { config.description = description })
}

func WithEBLStartTime(timestamp time.Time) EBLWriterOption {
	return eblWriterOptionFunc(func(config *eblWriterConfig) { config.start = timestamp })
}

type EBLWriterMetrics = ebl.WriterMetrics

// EBLWriter writes interoperable Actisense Enhanced Binary Log records. It is
// safe for concurrent use and retains the first output error.
type EBLWriter struct{ writer *ebl.Writer }

func NewEBLWriter(output io.Writer, options ...EBLWriterOption) (*EBLWriter, error) {
	config := eblWriterConfig{}
	for _, option := range options {
		if option == nil {
			return nil, errors.New("n2k: nil EBLWriterOption")
		}
		option.applyEBLWriter(&config)
	}
	writer, err := ebl.NewWriter(output, ebl.WriterConfig{StartTime: config.start, Description: config.description})
	if err != nil {
		return nil, fmt.Errorf("n2k: initialize EBL writer: %w", err)
	}
	return &EBLWriter{writer: writer}, nil
}

func (w *EBLWriter) WriteRawBST(timestamp time.Time, direction Direction, rawBST []byte) error {
	if w == nil || w.writer == nil {
		return errors.New("n2k: nil EBL writer")
	}
	return w.writer.WriteRawBST(timestamp, eblWriterDirection(direction), append([]byte(nil), rawBST...))
}

func (w *EBLWriter) WriteRawBytes(timestamp time.Time, direction Direction, data []byte) error {
	if w == nil || w.writer == nil {
		return errors.New("n2k: nil EBL writer")
	}
	return w.writer.WriteRawBytes(timestamp, eblWriterDirection(direction), append([]byte(nil), data...))
}

func (w *EBLWriter) WriteDescription(description string) error {
	if w == nil || w.writer == nil {
		return errors.New("n2k: nil EBL writer")
	}
	return w.writer.WriteDescription(description)
}

func (w *EBLWriter) Err() error {
	if w == nil || w.writer == nil {
		return errors.New("n2k: nil EBL writer")
	}
	return w.writer.Err()
}

func (w *EBLWriter) Metrics() EBLWriterMetrics {
	if w == nil || w.writer == nil {
		return EBLWriterMetrics{}
	}
	return w.writer.Metrics()
}

func eblWriterDirection(direction Direction) ebl.Direction {
	if direction == DirectionTransmitted {
		return ebl.DirectionTransmitted
	}
	if direction == DirectionReceived {
		return ebl.DirectionReceived
	}
	return ebl.DirectionUnknown
}

// ActisenseEBLTrace converts raw transport reads/writes into the SDK's EBL
// wire-trace representation. Valid received frames become checksum-stripped
// BSTRawFrame records; invalid or unframed bytes remain exact raw evidence.
type ActisenseEBLTrace struct {
	mu        sync.Mutex
	writer    *EBLWriter
	buffer    []byte
	startedAt time.Time
	err       error
}

func NewActisenseEBLTrace(writer *EBLWriter) (*ActisenseEBLTrace, error) {
	if writer == nil || writer.writer == nil {
		return nil, errors.New("n2k: Actisense EBL trace requires an EBL writer")
	}
	return &ActisenseEBLTrace{writer: writer, buffer: make([]byte, 0, 4096)}, nil
}

func (t *ActisenseEBLTrace) trace(direction actisense.WireDirection, timestamp time.Time, data []byte) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.err != nil || len(data) == 0 {
		return
	}
	if direction == actisense.WireTransmitted {
		t.err = t.writer.WriteRawBytes(timestamp, DirectionTransmitted, data)
		return
	}
	if len(t.buffer) == 0 {
		t.startedAt = timestamp
	}
	t.buffer = append(t.buffer, data...)
	t.processReceivedLocked(timestamp, false)
}

func (t *ActisenseEBLTrace) processReceivedLocked(timestamp time.Time, flush bool) {
	marker := []byte{actisense.DLE, actisense.STX}
	for t.err == nil && len(t.buffer) > 0 {
		start := bytes.Index(t.buffer, marker)
		if start < 0 {
			keep := 0
			if !flush && t.buffer[len(t.buffer)-1] == actisense.DLE {
				keep = 1
			}
			if len(t.buffer) > keep {
				t.err = t.writer.WriteRawBytes(t.startedAt, DirectionReceived, t.buffer[:len(t.buffer)-keep])
				t.buffer = append(t.buffer[:0], t.buffer[len(t.buffer)-keep:]...)
				t.startedAt = timestamp
			}
			return
		}
		if start > 0 {
			t.err = t.writer.WriteRawBytes(t.startedAt, DirectionReceived, t.buffer[:start])
			t.buffer = t.buffer[start:]
			t.startedAt = timestamp
			continue
		}
		end := findActisenseFrameEnd(t.buffer)
		if end < 0 {
			if flush || len(t.buffer) > actisense.MaxDatagramLength*2+16 {
				t.err = t.writer.WriteRawBytes(t.startedAt, DirectionReceived, t.buffer)
				t.buffer = t.buffer[:0]
			}
			return
		}
		frame := append([]byte(nil), t.buffer[:end]...)
		t.buffer = t.buffer[end:]
		var datagrams []actisense.Datagram
		var decodeErrors []actisense.DecodeError
		parser := actisense.NewParser()
		parser.Feed(frame, func(datagram actisense.Datagram) { datagrams = append(datagrams, datagram) }, func(decodeErr actisense.DecodeError) { decodeErrors = append(decodeErrors, decodeErr) })
		if len(datagrams) == 1 && len(decodeErrors) == 0 {
			t.err = t.writer.WriteRawBST(t.startedAt, DirectionReceived, datagrams[0].Raw)
		} else {
			t.err = t.writer.WriteRawBytes(t.startedAt, DirectionReceived, frame)
		}
		t.startedAt = timestamp
	}
}

func findActisenseFrameEnd(data []byte) int {
	for index := 2; index < len(data); index++ {
		if data[index] != actisense.DLE {
			continue
		}
		if index+1 >= len(data) {
			return -1
		}
		switch data[index+1] {
		case actisense.DLE:
			index++
		case actisense.ETX:
			return index + 2
		}
	}
	return -1
}

func (t *ActisenseEBLTrace) Flush() error {
	if t == nil {
		return nil
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.err == nil {
		t.processReceivedLocked(time.Now(), true)
	}
	return t.err
}

func (t *ActisenseEBLTrace) Err() error {
	if t == nil {
		return nil
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.err
}
