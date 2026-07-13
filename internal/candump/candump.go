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
	var idData string
	for _, field := range strings.Fields(line) {
		if strings.Contains(field, "#") {
			idData = field
			break
		}
	}
	parts := strings.SplitN(idData, "#", 2)
	if len(parts) != 2 || parts[0] == "" {
		return can.Frame{}, time.Time{}, false
	}
	payload := parts[1]
	if strings.HasPrefix(payload, "#") || strings.HasPrefix(payload, "R") {
		return can.Frame{}, time.Time{}, false
	}

	id, err := strconv.ParseUint(parts[0], 16, 32)
	if err != nil {
		return can.Frame{}, time.Time{}, false
	}
	data, err := hex.DecodeString(payload)
	if err != nil || len(data) > 8 {
		return can.Frame{}, time.Time{}, false
	}

	frame = can.Frame{ID: uint32(id), Length: uint8(len(data))}
	copy(frame.Data[:], data)
	return frame, parseTimestamp(line), true
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
