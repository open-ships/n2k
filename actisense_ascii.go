package n2k

import (
	"time"

	"github.com/brutella/can"
	"github.com/open-ships/n2k/internal/gateway"
)

// ParseActisenseCANASCII parses one source-authoritative mode-6 CAN frame.
func ParseActisenseCANASCII(line string) (Observation, error) {
	return gateway.ParseActisenseCANASCIIObservation(line)
}

// EncodeActisenseCANASCII renders one CAN frame in Actisense's documented
// RAW ASCII representation.
func EncodeActisenseCANASCII(frame can.Frame, direction Direction, transportTimestamp time.Duration) ([]byte, error) {
	return gateway.FormatActisenseCANASCII(frame, direction, transportTimestamp)
}

// ParseActisenseN2KASCII parses one gateway-assembled Type-A NMEA 2000
// message.
func ParseActisenseN2KASCII(line string) (Observation, error) {
	return gateway.ParseActisenseN2KASCIIObservation(line)
}

// EncodeActisenseN2KASCII renders one assembled NMEA 2000 message in the
// documented Type-A representation.
func EncodeActisenseN2KASCII(pgn uint32, priority, source, destination uint8, payload []byte, transportTimestamp time.Duration) ([]byte, error) {
	return gateway.FormatActisenseN2KASCII(pgn, priority, source, destination, payload, transportTimestamp)
}
