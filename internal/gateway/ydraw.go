// Package gateway parses the wire formats spoken by NMEA 2000 network
// gateways: the Yacht Devices RAW ASCII line protocol and the Actisense
// binary stream protocol.
package gateway

import (
	"encoding/hex"
	"strconv"
	"strings"

	"github.com/brutella/can"
)

// maxCANID is the largest valid 29-bit extended CAN identifier.
const maxCANID = 0x1FFFFFFF

// ParseYDRaw parses one line of the Yacht Devices RAW protocol (YDWG-02 and
// compatible gateways), e.g.
//
//	17:33:21.107 R 09F11201 01 5C 3D FF 7F FF 7F FC
//
// The direction field is R (received from the bus) or T (transmitted to the
// bus); both carry frames. Returns ok=false for service messages, blank
// lines, and anything else that does not carry a CAN data frame.
func ParseYDRaw(line string) (can.Frame, bool) {
	fields := strings.Fields(strings.TrimSpace(line))
	// time, direction, ID, and at least one data byte
	if len(fields) < 4 {
		return can.Frame{}, false
	}
	if fields[1] != "R" && fields[1] != "T" {
		return can.Frame{}, false
	}
	id, err := strconv.ParseUint(fields[2], 16, 32)
	if err != nil || id > maxCANID {
		return can.Frame{}, false
	}

	dataFields := fields[3:]
	if len(dataFields) > 8 {
		return can.Frame{}, false
	}
	frame := can.Frame{ID: uint32(id), Length: uint8(len(dataFields))}
	for i, b := range dataFields {
		v, err := hex.DecodeString(b)
		if err != nil || len(v) != 1 {
			return can.Frame{}, false
		}
		frame.Data[i] = v[0]
	}
	return frame, true
}
