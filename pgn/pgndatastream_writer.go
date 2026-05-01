package pgn

import (
	"fmt"
	"math"
)

// PGNDataStreamWriter provides a bit-level sequential writer that is the exact inverse
// of PGNDataStream. It builds up an NMEA 2000 message payload by writing typed fields
// in the order defined by a PGN specification. Sub-byte and cross-byte writes are handled
// transparently using the same little-endian bit ordering as the reader.
//
// Generated PGN encoder functions will create a PGNDataStreamWriter, call its typed write
// methods for each field, and then call Bytes() to obtain the finished payload.
type PGNDataStreamWriter struct {
	// data accumulates the raw message payload bytes being encoded.
	data []uint8

	// byteOffset is the whole-byte portion of the current write cursor position.
	byteOffset uint16
	// bitOffset is the sub-byte portion (0-7) of the current write cursor position.
	// Together with byteOffset, the absolute bit position is: byteOffset*8 + bitOffset.
	bitOffset uint8
}

// NewPGNDataStreamWriter creates a PGNDataStreamWriter with an empty payload buffer.
// Callers write fields sequentially, then call Bytes() to retrieve the encoded payload.
func NewPGNDataStreamWriter() *PGNDataStreamWriter {
	return &PGNDataStreamWriter{
		data:       make([]uint8, 0, 8), // pre-allocate for typical 8-byte CAN frame
		byteOffset: 0,
		bitOffset:  0,
	}
}

// Bytes returns the accumulated payload bytes. The returned slice is owned by the writer;
// callers should copy it if they need to retain it after further writes.
func (w *PGNDataStreamWriter) Bytes() []uint8 {
	return w.data
}

// putNumberRaw is the lowest-level write primitive. It writes up to 64 bits into the
// stream at the current cursor position in little-endian bit order. This is the exact
// inverse of PGNDataStream.getNumberRaw.
//
// Algorithm (mirrors getNumberRaw):
//  1. Loop while bitLength > 0.
//  2. Ensure the data slice has room for the current byte (append 0 if needed).
//  3. Calculate bits to write: min(8 - bitOffset, remaining bitLength).
//  4. Extract the low bits from value, mask them, shift into position, OR into current byte.
//  5. Shift value right by the number of bits written, advance cursor.
func (w *PGNDataStreamWriter) putNumberRaw(value uint64, bitLength uint16) error {
	if bitLength < 64 {
		maxVal := uint64(1<<bitLength) - 1
		if value > maxVal {
			return fmt.Errorf("value %d exceeds %d-bit unsigned range (max %d)", value, bitLength, maxVal)
		}
	}

	for bitLength > 0 {
		// Ensure we have a byte to write into.
		for int(w.byteOffset) >= len(w.data) {
			w.data = append(w.data, 0)
		}

		// How many bits can we write in the current byte from the current bit position.
		bitsToWrite := uint8(8) - w.bitOffset
		if bitLength < uint16(bitsToWrite) {
			bitsToWrite = uint8(bitLength)
		}

		// Extract the low bitsToWrite bits from value.
		var mask uint64
		if bitsToWrite < 64 {
			mask = (1 << bitsToWrite) - 1
		} else {
			mask = 0xFFFFFFFFFFFFFFFF
		}
		bits := uint8(value & mask)

		// Shift the bits into position within the current byte and OR them in.
		w.data[w.byteOffset] |= bits << w.bitOffset

		// Shift value right to discard the bits we just wrote.
		if bitsToWrite < 64 {
			value >>= bitsToWrite
		} else {
			value = 0
		}
		bitLength -= uint16(bitsToWrite)

		// Advance the write cursor.
		w.bitOffset += bitsToWrite
		if w.bitOffset >= 8 {
			w.bitOffset -= 8
			w.byteOffset++
		}
	}

	return nil
}

// putSignedNumber writes a signed two's-complement integer of the given bit width.
// Negative values are converted to their two's-complement unsigned representation
// by masking to the specified bit width.
func (w *PGNDataStreamWriter) putSignedNumber(value int64, bitLength uint16) error {
	var raw uint64
	if value >= 0 {
		raw = uint64(value)
	} else {
		// Two's complement: mask to bitLength bits.
		// For negative values, Go's uint64 cast preserves the two's complement representation.
		raw = uint64(value) & ((1 << bitLength) - 1)
	}
	return w.putNumberRaw(raw, bitLength)
}

// putNullUnsigned writes the unsigned null sentinel (all bits set) for the given bit width.
// In NMEA 2000, a field whose bits are all ones means "data not available".
func (w *PGNDataStreamWriter) putNullUnsigned(bitLength uint16) error {
	maxVal := uint64(0xFFFFFFFFFFFFFFFF)
	if bitLength < 64 {
		maxVal = (1 << bitLength) - 1
	}
	return w.putNumberRaw(maxVal, bitLength)
}

// putNullSigned writes the signed null sentinel (positive maximum) for the given bit width.
// In NMEA 2000, the positive maximum for a signed field means "data not available".
// For example, 0x7F for 8-bit, 0x7FFF for 16-bit.
func (w *PGNDataStreamWriter) putNullSigned(bitLength uint16) error {
	maxVal := uint64(0xFFFFFFFFFFFFFFFF)
	maxVal >>= (64 - bitLength)
	maxVal >>= 1 // Exclude the sign bit to get positive maximum.
	return w.putNumberRaw(maxVal, bitLength)
}

// skipBits advances the write cursor by bitLength bits, filling with zeros.
// This is used for reserved or unused fields in a PGN definition.
func (w *PGNDataStreamWriter) skipBits(bitLength uint16) error {
	return w.putNumberRaw(0, bitLength)
}

// writeLookupField writes an unsigned integer value at the given bit width.
// No null detection is performed -- every bit pattern is valid for enum/lookup fields.
func (w *PGNDataStreamWriter) writeLookupField(value uint64, bitLength uint16) error {
	if bitLength > 64 {
		return fmt.Errorf("requested %d bitLength in writeLookupField", bitLength)
	}
	return w.putNumberRaw(value, bitLength)
}

// writeUnsignedResolution writes a scaled unsigned value. The float value is divided by the
// resolution factor, rounded to the nearest integer, and written as an unsigned field.
// A nil pointer writes the unsigned null sentinel (all bits set).
func (w *PGNDataStreamWriter) writeUnsignedResolution(value *float32, bitLength uint16, resolution float32) error {
	if bitLength > 64 {
		return fmt.Errorf("requested %d bitLength in writeUnsignedResolution", bitLength)
	}
	if value == nil {
		return w.putNullUnsigned(bitLength)
	}
	raw := uint64(math.Round(float64(*value) / float64(resolution)))
	return w.putNumberRaw(raw, bitLength)
}

// writeSignedResolution writes a scaled signed value. The float value is divided by the
// resolution factor, rounded to the nearest integer, and written as a two's-complement field.
// A nil pointer writes the signed null sentinel (positive maximum).
func (w *PGNDataStreamWriter) writeSignedResolution(value *float32, bitLength uint16, resolution float32) error {
	if bitLength > 64 {
		return fmt.Errorf("requested %d bitLength in writeSignedResolution", bitLength)
	}
	if value == nil {
		return w.putNullSigned(bitLength)
	}
	raw := int64(math.Round(float64(*value) / float64(resolution)))
	return w.putSignedNumber(raw, bitLength)
}

// writeSignedResolution64Override is the float64 variant of writeSignedResolution.
// It exists for fields where float32 precision is insufficient, such as latitude
// and longitude which use very small resolution values applied to large integers.
func (w *PGNDataStreamWriter) writeSignedResolution64Override(value *float64, bitLength uint16, resolution float64) error {
	if bitLength > 64 {
		return fmt.Errorf("requested %d bitLength in writeSignedResolution64Override", bitLength)
	}
	if value == nil {
		return w.putNullSigned(bitLength)
	}
	raw := int64(math.Round(*value / resolution))
	return w.putSignedNumber(raw, bitLength)
}

// writeUInt8 writes an 8-bit unsigned integer. A nil pointer writes the unsigned null sentinel.
func (w *PGNDataStreamWriter) writeUInt8(value *uint8, bitLength uint16) error {
	if bitLength > 8 {
		return fmt.Errorf("requested %d bitLength in writeUInt8", bitLength)
	}
	if value == nil {
		return w.putNullUnsigned(bitLength)
	}
	return w.putNumberRaw(uint64(*value), bitLength)
}

// writeUInt16 writes a 16-bit unsigned integer. A nil pointer writes the unsigned null sentinel.
func (w *PGNDataStreamWriter) writeUInt16(value *uint16, bitLength uint16) error {
	if bitLength > 16 {
		return fmt.Errorf("requested %d bitLength in writeUInt16", bitLength)
	}
	if value == nil {
		return w.putNullUnsigned(bitLength)
	}
	return w.putNumberRaw(uint64(*value), bitLength)
}

// writeUInt32 writes a 32-bit unsigned integer. A nil pointer writes the unsigned null sentinel.
func (w *PGNDataStreamWriter) writeUInt32(value *uint32, bitLength uint16) error {
	if bitLength > 32 {
		return fmt.Errorf("requested %d bitLength in writeUInt32", bitLength)
	}
	if value == nil {
		return w.putNullUnsigned(bitLength)
	}
	return w.putNumberRaw(uint64(*value), bitLength)
}

// writeUInt64 writes a 64-bit unsigned integer. A nil pointer writes the unsigned null sentinel.
func (w *PGNDataStreamWriter) writeUInt64(value *uint64, bitLength uint16) error {
	if bitLength > 64 {
		return fmt.Errorf("requested %d bitLength in writeUInt64", bitLength)
	}
	if value == nil {
		return w.putNullUnsigned(bitLength)
	}
	return w.putNumberRaw(*value, bitLength)
}

// writeInt8 writes an 8-bit signed two's-complement integer.
// A nil pointer writes the signed null sentinel.
func (w *PGNDataStreamWriter) writeInt8(value *int8, bitLength uint16) error {
	if bitLength > 8 {
		return fmt.Errorf("requested %d bitLength in writeInt8", bitLength)
	}
	if value == nil {
		return w.putNullSigned(bitLength)
	}
	return w.putSignedNumber(int64(*value), bitLength)
}

// writeInt16 writes a 16-bit signed two's-complement integer.
// A nil pointer writes the signed null sentinel.
func (w *PGNDataStreamWriter) writeInt16(value *int16, bitLength uint16) error {
	if bitLength > 16 {
		return fmt.Errorf("requested %d bitLength in writeInt16", bitLength)
	}
	if value == nil {
		return w.putNullSigned(bitLength)
	}
	return w.putSignedNumber(int64(*value), bitLength)
}

// writeInt32 writes a 32-bit signed two's-complement integer.
// A nil pointer writes the signed null sentinel.
func (w *PGNDataStreamWriter) writeInt32(value *int32, bitLength uint16) error {
	if bitLength > 32 {
		return fmt.Errorf("requested %d bitLength in writeInt32", bitLength)
	}
	if value == nil {
		return w.putNullSigned(bitLength)
	}
	return w.putSignedNumber(int64(*value), bitLength)
}

// writeInt64 writes a 64-bit signed two's-complement integer.
// A nil pointer writes the signed null sentinel.
func (w *PGNDataStreamWriter) writeInt64(value *int64, bitLength uint16) error {
	if bitLength > 64 {
		return fmt.Errorf("requested %d bitLength in writeInt64", bitLength)
	}
	if value == nil {
		return w.putNullSigned(bitLength)
	}
	return w.putSignedNumber(*value, bitLength)
}

// writeFloat32 writes an IEEE 754 single-precision float as 32 bits in little-endian order.
// A nil pointer writes the unsigned null sentinel (0xFFFFFFFF), which is the same sentinel
// that readFloat32 checks for null detection.
func (w *PGNDataStreamWriter) writeFloat32(value *float32) error {
	if value == nil {
		return w.putNullUnsigned(32)
	}
	bits := math.Float32bits(*value)
	return w.putNumberRaw(uint64(bits), 32)
}

// writeBinaryData writes raw binary bytes into the stream at the current bit position.
// The data is written in chunks of up to 64 bits to mirror the reader's readBinaryData
// which reads in 64-bit chunks. The bitLength parameter specifies the total number of
// bits to write; if it is not a multiple of 8, only the low bits of the final byte are used.
func (w *PGNDataStreamWriter) writeBinaryData(data []uint8, bitLength uint16) error {
	idx := 0
	for i := uint16(0); i < bitLength; i += 64 {
		// Determine how many bits to write in this chunk (up to 64).
		num := uint16(64)
		if bitLength-i < 64 {
			num = bitLength - i
		}

		// Assemble the chunk value from individual bytes, placing each byte
		// at the correct position (little-endian order within the chunk).
		var value uint64
		for h := uint16(0); h < num; h += 8 {
			if idx < len(data) {
				bitsInThisByte := num - h
				if bitsInThisByte > 8 {
					bitsInThisByte = 8
				}
				// Mask the byte to only include the bits we need (for sub-byte final chunk).
				b := data[idx]
				if bitsInThisByte < 8 {
					b &= uint8((1 << bitsInThisByte) - 1)
				}
				value |= uint64(b) << h
				idx++
			}
		}

		err := w.putNumberRaw(value, num)
		if err != nil {
			return err
		}
	}
	return nil
}

// writeFixedString writes a fixed-width string field of exactly bitLength bits.
// The string is written as raw bytes, and any remaining space is padded with 0xFF
// (the NMEA 2000 "data not available" fill byte). This matches the convention used
// by readFixedString, which strips 0xFF, 0x00, and '@' padding on read.
func (w *PGNDataStreamWriter) writeFixedString(s string, bitLength uint16) error {
	numBytes := bitLength / 8
	buf := make([]uint8, numBytes)

	// Copy string bytes into the buffer.
	sBytes := []uint8(s)
	copy(buf, sBytes)

	// Pad remaining bytes with 0xFF.
	for i := len(sBytes); i < int(numBytes); i++ {
		buf[i] = 0xFF
	}

	return w.writeBinaryData(buf, bitLength)
}

// writeStringWithLength writes a STRING_LZ encoded string.
// Wire format:
//   - Byte 0: length of the string data in bytes (does NOT include this length byte itself)
//   - Bytes 1..N: the string character data
func (w *PGNDataStreamWriter) writeStringWithLength(s string) error {
	sBytes := []uint8(s)
	length := uint8(len(sBytes))

	// Write the length byte.
	err := w.putNumberRaw(uint64(length), 8)
	if err != nil {
		return err
	}

	// Write the string data.
	return w.writeBinaryData(sBytes, uint16(length)*8)
}

// writeVariableData writes raw binary bytes back for a variable-data field. It mirrors
// readVariableData by looking up the field descriptor and using the appropriate encoding:
// STRING_LAU fields are written with writeStringWithLengthAndControl, all others are
// written as binary data at the field's bit length (rounded up to byte boundary).
func (w *PGNDataStreamWriter) writeVariableData(data []uint8, pgn uint32, manID ManufacturerCodeConst, fieldIndex uint8) error {
	field, err := GetFieldDescriptor(pgn, manID, fieldIndex)
	if err != nil {
		return err
	}
	if field.BitLengthVariable && field.CanboatType == "STRING_LAU" {
		return w.writeStringWithLengthAndControl(string(data))
	}
	bitLen := (field.BitLength + 7) &^ 0x7
	return w.writeBinaryData(data, bitLen)
}

// writeStringWithLengthAndControl writes a STRING_LAU encoded string.
// Wire format:
//   - Byte 0: total length in bytes (includes this byte, the control byte, the string chars,
//     and a terminating zero)
//   - Byte 1: control/encoding byte (1 = ASCII)
//   - Bytes 2..N: the string character data plus a trailing NUL
func (w *PGNDataStreamWriter) writeStringWithLengthAndControl(s string) error {
	sBytes := []uint8(s)
	// Total length = 1 (length byte) + 1 (control byte) + len(string) + 1 (trailing NUL)
	totalLength := uint8(len(sBytes) + 3)
	controlByte := uint8(1) // 1 = ASCII

	// Write the 2-byte header as binary data to match the reader's readBinaryData(16) call.
	header := []uint8{totalLength, controlByte}
	err := w.writeBinaryData(header, 16)
	if err != nil {
		return err
	}

	// Write the string data followed by a trailing NUL.
	payload := make([]uint8, len(sBytes)+1)
	copy(payload, sBytes)
	payload[len(sBytes)] = 0 // trailing NUL

	return w.writeBinaryData(payload, uint16(len(payload))*8)
}
