package gateway

import (
	"bufio"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/brutella/can"
	"github.com/open-ships/n2k/raw"
)

const (
	maxActisenseASCIIFrameLine = 128
	maxActisenseN2KPayload     = 1785
	maxActisenseN2KASCIILine   = 2*maxActisenseN2KPayload + 128
)

// ParseActisenseCANASCIIObservation parses the Actisense CAN-frame ASCII
// representation used by operating mode 6. Zero-to-eight-byte frames and
// timestamps with or without fractional seconds are accepted.
func ParseActisenseCANASCIIObservation(line string) (raw.Observation, error) {
	fields := strings.Fields(strings.TrimSpace(line))
	if len(fields) < 3 {
		return raw.Observation{}, errors.New("actisense: CAN ASCII line is shorter than timestamp, direction, and CAN ID")
	}
	timestamp, ok := parseActisenseClock(fields[0], true)
	if !ok {
		return raw.Observation{}, fmt.Errorf("actisense: invalid CAN ASCII timestamp %q", fields[0])
	}
	direction := raw.DirectionReceived
	switch fields[1] {
	case "R":
	case "T":
		direction = raw.DirectionTransmitted
	default:
		return raw.Observation{}, fmt.Errorf("actisense: invalid CAN ASCII direction %q", fields[1])
	}
	if len(fields[2]) != 8 {
		return raw.Observation{}, fmt.Errorf("actisense: CAN ASCII identifier %q is not eight hexadecimal digits", fields[2])
	}
	id, err := strconv.ParseUint(fields[2], 16, 32)
	if err != nil || id > maxCANID {
		return raw.Observation{}, fmt.Errorf("actisense: invalid 29-bit CAN ASCII identifier %q", fields[2])
	}
	dataFields := fields[3:]
	if len(dataFields) > 8 {
		return raw.Observation{}, fmt.Errorf("actisense: CAN ASCII line has %d data bytes; maximum is 8", len(dataFields))
	}
	frame := can.Frame{ID: uint32(id), Length: uint8(len(dataFields))}
	for index, field := range dataFields {
		if len(field) != 2 {
			return raw.Observation{}, fmt.Errorf("actisense: CAN ASCII data field %q is not two hexadecimal digits", field)
		}
		decoded, decodeErr := hex.DecodeString(field)
		if decodeErr != nil {
			return raw.Observation{}, fmt.Errorf("actisense: invalid CAN ASCII data field %q", field)
		}
		frame.Data[index] = decoded[0]
	}
	now := time.Now()
	return raw.Observation{
		Kind: raw.KindFrame, Timestamp: now, ReceivedAt: now,
		TransportTimestamp: timestamp, HasTransportTimestamp: true,
		AdapterID: "actisense-can-ascii", Direction: direction, Frame: &frame,
	}, nil
}

// FormatActisenseCANASCII renders a mode-6 CAN frame. A zero timestamp is
// represented as 00:00:00.000, which is suitable for host-to-gateway writes.
func FormatActisenseCANASCII(frame can.Frame, direction raw.Direction, timestamp time.Duration) ([]byte, error) {
	if frame.ID > maxCANID {
		return nil, fmt.Errorf("actisense: CAN identifier 0x%X exceeds 29 bits", frame.ID)
	}
	if frame.Length > 8 {
		return nil, fmt.Errorf("actisense: CAN frame length %d exceeds 8", frame.Length)
	}
	marker := byte('R')
	if direction == raw.DirectionTransmitted {
		marker = 'T'
	} else if direction != raw.DirectionReceived {
		return nil, fmt.Errorf("actisense: cannot format CAN ASCII direction %q", direction)
	}
	buf := make([]byte, 0, 48)
	buf = appendActisenseClock(buf, timestamp, true)
	buf = append(buf, ' ', marker, ' ')
	buf = fmt.Appendf(buf, "%08X", frame.ID)
	for index := range int(frame.Length) {
		buf = fmt.Appendf(buf, " %02X", frame.Data[index])
	}
	return append(buf, '\r', '\n'), nil
}

// ParseActisenseN2KASCIIObservation parses one gateway-assembled N2K ASCII
// message. It preserves source, destination, priority, PGN, and up to the
// complete 1785-byte NMEA 2000 transport payload.
func ParseActisenseN2KASCIIObservation(line string) (raw.Observation, error) {
	fields := strings.Fields(strings.TrimSpace(line))
	if len(fields) < 3 || len(fields[0]) < 8 || fields[0][0] != 'A' {
		return raw.Observation{}, errors.New("actisense: N2K ASCII line has no A/timestamp header")
	}
	timestamp, ok := parseActisenseClock(fields[0][1:], false)
	if !ok {
		return raw.Observation{}, fmt.Errorf("actisense: invalid N2K ASCII timestamp %q", fields[0][1:])
	}
	if len(fields[1]) != 5 {
		return raw.Observation{}, fmt.Errorf("actisense: N2K ASCII address/priority field %q is not five hexadecimal digits", fields[1])
	}
	sourceValue, sourceErr := strconv.ParseUint(fields[1][:2], 16, 8)
	destinationValue, destinationErr := strconv.ParseUint(fields[1][2:4], 16, 8)
	priorityValue, priorityErr := strconv.ParseUint(fields[1][4:], 16, 4)
	if sourceErr != nil || destinationErr != nil || priorityErr != nil {
		return raw.Observation{}, errors.New("actisense: N2K ASCII address/priority field is not hexadecimal")
	}
	pgnValue, err := strconv.ParseUint(fields[2], 16, 32)
	if err != nil || pgnValue > 0x3FFFF {
		return raw.Observation{}, fmt.Errorf("actisense: invalid N2K ASCII PGN %q", fields[2])
	}
	payloadHex := strings.Join(fields[3:], "")
	if len(payloadHex)%2 != 0 {
		return raw.Observation{}, errors.New("actisense: N2K ASCII payload has an odd number of hexadecimal digits")
	}
	if len(payloadHex)/2 > maxActisenseN2KPayload {
		return raw.Observation{}, fmt.Errorf("actisense: N2K ASCII payload is %d bytes; maximum is %d", len(payloadHex)/2, maxActisenseN2KPayload)
	}
	payload, err := hex.DecodeString(payloadHex)
	if err != nil {
		return raw.Observation{}, errors.New("actisense: N2K ASCII payload is not hexadecimal")
	}
	destination := uint8(destinationValue)
	now := time.Now()
	return raw.Observation{
		Kind: raw.KindMessage, Timestamp: now, ReceivedAt: now,
		TransportTimestamp: timestamp, HasTransportTimestamp: true,
		AdapterID: "actisense-n2k-ascii", Direction: raw.DirectionReceived,
		PGN: uint32(pgnValue), Priority: uint8(priorityValue) & 0x07,
		Source: uint8(sourceValue), Destination: &destination, Payload: payload,
	}, nil
}

// FormatActisenseN2KASCII renders one assembled message in the documented
// Type-A representation.
func FormatActisenseN2KASCII(pgn uint32, priority, source, destination uint8, payload []byte, timestamp time.Duration) ([]byte, error) {
	if pgn > 0x3FFFF {
		return nil, fmt.Errorf("actisense: PGN %d exceeds 18 bits", pgn)
	}
	if priority > 7 {
		return nil, fmt.Errorf("actisense: priority %d exceeds 7", priority)
	}
	if len(payload) > maxActisenseN2KPayload {
		return nil, fmt.Errorf("actisense: N2K ASCII payload is %d bytes; maximum is %d", len(payload), maxActisenseN2KPayload)
	}
	buf := make([]byte, 0, 24+2*len(payload))
	buf = append(buf, 'A')
	buf = appendActisenseClock(buf, timestamp, false)
	buf = fmt.Appendf(buf, " %02X%02X%X %05X ", source, destination, priority, pgn)
	hexStart := len(buf)
	buf = append(buf, make([]byte, hex.EncodedLen(len(payload)))...)
	hex.Encode(buf[hexStart:], payload)
	for index := hexStart; index < len(buf); index++ {
		if buf[index] >= 'a' && buf[index] <= 'f' {
			buf[index] -= 'a' - 'A'
		}
	}
	return append(buf, '\r', '\n'), nil
}

func parseActisenseClock(value string, colons bool) (time.Duration, bool) {
	layouts := []string{"150405.000", "150405"}
	if colons {
		layouts = []string{"15:04:05.000", "15:04:05"}
	}
	for _, layout := range layouts {
		parsed, err := time.Parse(layout, value)
		if err == nil {
			return time.Duration(parsed.Hour())*time.Hour + time.Duration(parsed.Minute())*time.Minute +
				time.Duration(parsed.Second())*time.Second + time.Duration(parsed.Nanosecond()), true
		}
	}
	return 0, false
}

func appendActisenseClock(buf []byte, timestamp time.Duration, colons bool) []byte {
	if timestamp < 0 {
		timestamp = 0
	}
	timestamp %= 24 * time.Hour
	hours := timestamp / time.Hour
	minutes := timestamp / time.Minute % 60
	seconds := timestamp / time.Second % 60
	millis := timestamp / time.Millisecond % 1000
	if colons {
		return fmt.Appendf(buf, "%02d:%02d:%02d.%03d", hours, minutes, seconds, millis)
	}
	return fmt.Appendf(buf, "%02d%02d%02d.%03d", hours, minutes, seconds, millis)
}

func ReadActisenseCANASCIIObservations(reader io.Reader, handler func(raw.Observation)) error {
	return readActisenseASCIILines(reader, maxActisenseASCIIFrameLine, ParseActisenseCANASCIIObservation, handler)
}

func ReadActisenseN2KASCIIObservations(reader io.Reader, handler func(raw.Observation)) error {
	return readActisenseASCIILines(reader, maxActisenseN2KASCIILine, ParseActisenseN2KASCIIObservation, handler)
}

func readActisenseASCIILines(reader io.Reader, maximum int, parse func(string) (raw.Observation, error), handler func(raw.Observation)) error {
	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, 1024), maximum)
	for scanner.Scan() {
		observation, err := parse(scanner.Text())
		if err == nil && handler != nil {
			handler(observation)
		}
	}
	return scanner.Err()
}

type actisenseCANASCIILines struct {
	buf      []byte
	dropping bool
}

func (a *actisenseCANASCIILines) feed(data []byte, handler func(raw.Observation)) {
	for _, value := range data {
		if value == '\n' {
			if !a.dropping {
				if observation, err := ParseActisenseCANASCIIObservation(string(a.buf)); err == nil && handler != nil {
					handler(observation)
				}
			}
			a.buf = a.buf[:0]
			a.dropping = false
			continue
		}
		if a.dropping {
			continue
		}
		if len(a.buf) >= maxActisenseASCIIFrameLine {
			a.buf = a.buf[:0]
			a.dropping = true
			continue
		}
		a.buf = append(a.buf, value)
	}
}
