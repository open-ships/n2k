// Package gateway parses the wire formats spoken by NMEA 2000 network
// gateways: the Yacht Devices RAW ASCII line protocol and the Actisense
// binary stream protocol.
package gateway

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/brutella/can"
	"github.com/open-ships/n2k/raw"
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
	observation, ok := ParseYDRawObservation(line)
	if !ok || observation.Frame == nil {
		return can.Frame{}, false
	}
	return *observation.Frame, true
}

// ParseYDRawObservation preserves the gateway's relative timestamp and
// receive/transmit direction in addition to the CAN frame.
func ParseYDRawObservation(line string) (raw.Observation, bool) {
	fields := strings.Fields(strings.TrimSpace(line))
	// time, direction, ID, and at least one data byte
	if len(fields) < 4 {
		return raw.Observation{}, false
	}
	if fields[1] != "R" && fields[1] != "T" {
		return raw.Observation{}, false
	}
	id, err := strconv.ParseUint(fields[2], 16, 32)
	if err != nil || id > maxCANID {
		return raw.Observation{}, false
	}

	dataFields := fields[3:]
	if len(dataFields) > 8 {
		return raw.Observation{}, false
	}
	frame := can.Frame{ID: uint32(id), Length: uint8(len(dataFields))}
	for i, b := range dataFields {
		v, err := hex.DecodeString(b)
		if err != nil || len(v) != 1 {
			return raw.Observation{}, false
		}
		frame.Data[i] = v[0]
	}
	direction := raw.DirectionReceived
	if fields[1] == "T" {
		direction = raw.DirectionTransmitted
	}
	transportTimestamp, timestampOK := parseYDTime(fields[0])
	return raw.Observation{
		Kind:                  raw.KindFrame,
		ReceivedAt:            time.Now(),
		TransportTimestamp:    transportTimestamp,
		HasTransportTimestamp: timestampOK,
		AdapterID:             "ydraw",
		Direction:             direction,
		Frame:                 &frame,
	}, true
}

func parseYDTime(value string) (time.Duration, bool) {
	parsed, err := time.Parse("15:04:05.000", value)
	if err != nil {
		return 0, false
	}
	return time.Duration(parsed.Hour())*time.Hour +
		time.Duration(parsed.Minute())*time.Minute +
		time.Duration(parsed.Second())*time.Second +
		time.Duration(parsed.Nanosecond()), true
}

// FormatYDRawTX renders a CAN frame as a Yacht Devices RAW transmit line.
// Application-to-gateway lines carry no time or direction field — just the
// 29-bit identifier and the data bytes, CRLF-terminated:
//
//	19F51323 01 02<CR><LF>
//
// Once the gateway transmits the frame onto the bus it echoes the line back
// with a timestamp and direction T, which ParseYDRaw accepts like any other
// frame. Filtered or malformed lines get no echo.
func FormatYDRawTX(frame can.Frame) []byte {
	length := int(frame.Length)
	if length > 8 {
		length = 8
	}
	buf := make([]byte, 0, 8+3*length+2)
	buf = fmt.Appendf(buf, "%08X", frame.ID)
	for i := 0; i < length; i++ {
		buf = fmt.Appendf(buf, " %02X", frame.Data[i])
	}
	return append(buf, '\r', '\n')
}
