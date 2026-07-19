// Package candump parses candump -L / -l log lines into CAN frames.
package candump

import (
	"encoding/hex"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/brutella/can"
)

// Record preserves the source identity and timestamp carried by a candump
// line in addition to its classical CAN frame.
type Record struct {
	Frame     can.Frame
	Timestamp time.Time
	Network   string
}

// Parse parses one candump -L / -l log line, e.g.
//
//	(1720000000.000000) can0 09F50E7F#00FFFFFFFFFFFFFF
//
// It finds the ID#DATA field so extra or missing columns are tolerated, and
// extracts the leading parenthesized timestamp when present (ts is the zero
// time.Time when the timestamp is absent or unparseable; ok is unaffected).
// Returns ok=false for lines that carry no classic CAN data frame (RTR
// frames, CAN FD "##" frames, comments, malformed lines).
func Parse(line string) (frame can.Frame, ts time.Time, ok bool) {
	record, ok := ParseRecord(line)
	return record.Frame, record.Timestamp, ok
}

// ParseRecord parses a candump line and retains its interface name.
func ParseRecord(line string) (record Record, ok bool) {
	var idData string
	var network string
	fields := strings.Fields(line)
	for i, field := range fields {
		if strings.Contains(field, "#") {
			idData = field
			if i > 0 && !strings.HasPrefix(fields[i-1], "(") {
				network = fields[i-1]
			}
			break
		}
	}
	parts := strings.SplitN(idData, "#", 2)
	if len(parts) != 2 || parts[0] == "" {
		return Record{}, false
	}
	payload := parts[1]
	if strings.HasPrefix(payload, "#") || strings.HasPrefix(payload, "R") {
		return Record{}, false
	}

	id, err := strconv.ParseUint(parts[0], 16, 32)
	if err != nil {
		return Record{}, false
	}
	data, err := hex.DecodeString(payload)
	if err != nil || len(data) > 8 {
		return Record{}, false
	}

	frame := can.Frame{ID: uint32(id), Length: uint8(len(data))}
	copy(frame.Data[:], data)
	return Record{Frame: frame, Timestamp: parseTimestamp(line), Network: network}, true
}

// parseTimestamp extracts the "(seconds.fraction)" prefix candump -L writes.
func parseTimestamp(line string) time.Time {
	trimmed := strings.TrimSpace(line)
	if !strings.HasPrefix(trimmed, "(") {
		return time.Time{}
	}
	end := strings.IndexByte(trimmed, ')')
	if end < 0 {
		return time.Time{}
	}
	seconds, err := strconv.ParseFloat(trimmed[1:end], 64)
	if err != nil {
		return time.Time{}
	}
	sec, frac := math.Modf(seconds)
	return time.Unix(int64(sec), int64(math.Round(frac*1e9)))
}
