package adapter

import (
	"strconv"
	"strings"

	"github.com/brutella/can"
	"github.com/open-ships/n2k/internal/framer"
)

// CanFrameFromRaw parses a CSV-formatted raw CAN bus log line into a can.Frame suitable
// for testing and replay. The expected input format (matching common N2K replay tools) is:
//
//	timestamp,priority,pgn,source,destination,length,byte0,byte1,...,byteN
//
// For example:
//
//	"2023-01-21T00:04:17Z,3,127501,224,0,8,00,03,c0,ff,ff,ff,ff,ff"
//
// The function reconstructs the 29-bit CAN extended ID from the parsed fields using
// framer.BuildCANID, and always sets the frame Length to 8 (standard CAN data frame size).
// Data bytes beyond the declared length are left at their zero default value.
//
// Note: This function ignores parse errors on individual fields for simplicity in test
// code. In production, a more robust parser would be needed.
//
// Parameters:
//   - in: A comma-separated string in the format described above.
//
// Returns a can.Frame with the reconstructed CAN ID and parsed data bytes.
func CanFrameFromRaw(in string) can.Frame {
	// Split the CSV line into its component fields.
	elems := strings.Split(in, ",")

	// Parse the individual NMEA 2000 fields from the CSV columns.
	// Column 0: timestamp (ignored for frame construction)
	// Column 1: priority (3-bit value, 0-7)
	priority, _ := strconv.ParseUint(elems[1], 10, 8)
	// Column 2: PGN (Parameter Group Number, up to 18 bits)
	pgn, _ := strconv.ParseUint(elems[2], 10, 32)
	// Column 3: source address (8-bit NMEA 2000 bus address)
	source, _ := strconv.ParseUint(elems[3], 10, 8)
	// Column 4: destination address (8-bit, 255 = broadcast)
	destination, _ := strconv.ParseUint(elems[4], 10, 8)
	// Column 5: data length (number of data bytes following)
	length, _ := strconv.ParseUint(elems[5], 10, 8)

	// Reconstruct the 29-bit CAN extended ID by encoding PGN, source, priority, and
	// destination into the appropriate bit positions.
	id := framer.BuildCANID(uint32(pgn), uint8(priority), uint8(source), uint8(destination))
	retval := can.Frame{
		ID:     id,
		Length: 8, // CAN data frames always carry 8 bytes on NMEA 2000
	}

	// Parse each hex-encoded data byte from columns 6 onward.
	for i := 0; i < int(length); i++ {
		b, _ := strconv.ParseUint(elems[i+6], 16, 8)
		retval.Data[i] = uint8(b)
	}

	return retval
}
