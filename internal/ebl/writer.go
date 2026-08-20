package ebl

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"
)

type WriterConfig struct {
	StartTime   time.Time
	Description string
}

type WriterMetrics struct {
	Records uint64
	Bytes   uint64
	Errors  uint64
}

// Writer serializes Enhanced Binary Log records. Each public write is one
// atomic event group so timestamps and direction markers cannot interleave.
type Writer struct {
	mu      sync.Mutex
	output  io.Writer
	err     error
	metrics WriterMetrics
}

func NewWriter(output io.Writer, config WriterConfig) (*Writer, error) {
	if output == nil {
		return nil, errors.New("ebl: nil output writer")
	}
	if config.StartTime.IsZero() {
		config.StartTime = time.Now()
	}
	writer := &Writer{output: output}
	writer.mu.Lock()
	defer writer.mu.Unlock()
	// Match the SDK wire-trace preamble order.
	if err := writer.writeTimeLocked(config.StartTime); err != nil {
		return nil, err
	}
	version := make([]byte, 4)
	binary.LittleEndian.PutUint32(version, CurrentVersion)
	if err := writer.writeTagLocked(TagVersion, version); err != nil {
		return nil, err
	}
	if config.Description != "" {
		if err := writer.writeTagLocked(TagDescription, []byte(config.Description)); err != nil {
			return nil, err
		}
	}
	return writer, nil
}

func (w *Writer) WriteRawBST(timestamp time.Time, direction Direction, raw []byte) error {
	if len(raw) < 2 || len(raw) > MaxTagPayload {
		return fmt.Errorf("ebl: raw BST record is %d bytes; expected 2-%d", len(raw), MaxTagPayload)
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	if err := w.readyLocked(); err != nil {
		return err
	}
	if err := w.writeEventPrefixLocked(timestamp, direction); err != nil {
		return err
	}
	return w.writeTagLocked(TagBSTRawFrame, raw)
}

func (w *Writer) WriteRawBytes(timestamp time.Time, direction Direction, data []byte) error {
	if len(data) == 0 {
		return nil
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	if err := w.readyLocked(); err != nil {
		return err
	}
	if err := w.writeEventPrefixLocked(timestamp, direction); err != nil {
		return err
	}
	escaped := stuff(data)
	return w.writeLocked(escaped)
}

func (w *Writer) WriteDescription(description string) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if err := w.readyLocked(); err != nil {
		return err
	}
	return w.writeTagLocked(TagDescription, []byte(description))
}

func (w *Writer) Err() error {
	if w == nil {
		return errors.New("ebl: nil writer")
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.err
}

func (w *Writer) Metrics() WriterMetrics {
	if w == nil {
		return WriterMetrics{}
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.metrics
}

func (w *Writer) readyLocked() error {
	if w == nil || w.output == nil {
		return errors.New("ebl: writer is closed")
	}
	return w.err
}

func (w *Writer) writeEventPrefixLocked(timestamp time.Time, direction Direction) error {
	if timestamp.IsZero() {
		timestamp = time.Now()
	}
	if err := w.writeTimeLocked(timestamp); err != nil {
		return err
	}
	var marker byte
	switch direction {
	case DirectionReceived:
		marker = 0
	case DirectionTransmitted:
		marker = 1
	default:
		return fmt.Errorf("ebl: direction %d is unknown", direction)
	}
	return w.writeTagLocked(TagDirectionMarker, []byte{marker})
}

func (w *Writer) writeTimeLocked(timestamp time.Time) error {
	ticks, err := timeToFiletime(timestamp)
	if err != nil {
		return err
	}
	payload := make([]byte, 8)
	binary.LittleEndian.PutUint64(payload, ticks)
	return w.writeTagLocked(TagTimeUTC, payload)
}

func (w *Writer) writeTagLocked(tag byte, payload []byte) error {
	if len(payload) > MaxTagPayload {
		return fmt.Errorf("ebl: tag 0x%02X payload is %d bytes; maximum is %d", tag, len(payload), MaxTagPayload)
	}
	record := make([]byte, 0, 5+len(payload)*2)
	record = append(record, escape, start)
	record = append(record, stuff([]byte{tag})...)
	record = append(record, stuff(payload)...)
	record = append(record, escape, end)
	if err := w.writeLocked(record); err != nil {
		return err
	}
	w.metrics.Records++
	return nil
}

func (w *Writer) writeLocked(data []byte) error {
	for len(data) > 0 {
		written, err := w.output.Write(data)
		if written > 0 {
			w.metrics.Bytes += uint64(written)
			data = data[written:]
		}
		if err != nil {
			w.metrics.Errors++
			w.err = err
			return err
		}
		if written == 0 {
			w.metrics.Errors++
			w.err = io.ErrShortWrite
			return io.ErrShortWrite
		}
	}
	return nil
}

func stuff(data []byte) []byte {
	result := make([]byte, 0, len(data)+4)
	for _, value := range data {
		result = append(result, value)
		if value == escape {
			result = append(result, escape)
		}
	}
	return result
}

func timeToFiletime(timestamp time.Time) (uint64, error) {
	seconds := timestamp.Unix()
	if seconds < -filetimeEpochOffsetSeconds {
		return 0, errors.New("ebl: timestamp predates the FILETIME epoch")
	}
	return uint64(seconds+filetimeEpochOffsetSeconds)*10_000_000 + uint64(timestamp.Nanosecond()/100), nil
}
