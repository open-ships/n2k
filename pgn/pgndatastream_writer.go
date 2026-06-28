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

	// err holds the first error encountered. Once set, all subsequent writes
	// are no-ops. Callers check Err() after all writes are complete.
	err error
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

// Err returns the first error encountered during writes, or nil if all writes succeeded.
func (w *PGNDataStreamWriter) Err() error {
	return w.err
}

func (w *PGNDataStreamWriter) setErr(err error) {
	if w.err == nil {
		w.err = err
	}
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
	if w.err != nil {
		return w.err
	}
	if bitLength < 64 {
		maxVal := uint64(1<<bitLength) - 1
		if value > maxVal {
			err := fmt.Errorf("value %d exceeds %d-bit unsigned range (max %d)", value, bitLength, maxVal)
			w.setErr(err)
			return err
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

func (w *PGNDataStreamWriter) writeSpareBits(bitLength uint16) {
	w.setErr(w.putNumberRaw(0, bitLength))
}

func (w *PGNDataStreamWriter) writeReservedBits(bitLength uint16) {
	if bitLength == 0 {
		return
	}
	maxVal := uint64(0xFFFFFFFFFFFFFFFF)
	if bitLength < 64 {
		maxVal = (1 << bitLength) - 1
	}
	w.setErr(w.putNumberRaw(maxVal, bitLength))
}

func (w *PGNDataStreamWriter) skipBits(bitLength uint16) {
	w.writeSpareBits(bitLength)
}

func (w *PGNDataStreamWriter) writeLookupField(value uint64, bitLength uint16) {
	if bitLength > 64 {
		w.setErr(fmt.Errorf("requested %d bitLength in writeLookupField", bitLength))
		return
	}
	w.setErr(w.putNumberRaw(value, bitLength))
}

func (w *PGNDataStreamWriter) writeUnsignedResolution(value *float32, bitLength uint16, resolution float32) {
	if bitLength > 64 {
		w.setErr(fmt.Errorf("requested %d bitLength in writeUnsignedResolution", bitLength))
		return
	}
	if value == nil {
		w.setErr(w.putNullUnsigned(bitLength))
		return
	}
	raw := uint64(math.Round(float64(*value) / float64(resolution)))
	w.setErr(w.putNumberRaw(raw, bitLength))
}

func (w *PGNDataStreamWriter) writeSignedResolution(value *float32, bitLength uint16, resolution float32) {
	if bitLength > 64 {
		w.setErr(fmt.Errorf("requested %d bitLength in writeSignedResolution", bitLength))
		return
	}
	if value == nil {
		w.setErr(w.putNullSigned(bitLength))
		return
	}
	raw := int64(math.Round(float64(*value) / float64(resolution)))
	w.setErr(w.putSignedNumber(raw, bitLength))
}

func (w *PGNDataStreamWriter) writeSignedResolution64Override(value *float64, bitLength uint16, resolution float64) {
	if bitLength > 64 {
		w.setErr(fmt.Errorf("requested %d bitLength in writeSignedResolution64Override", bitLength))
		return
	}
	if value == nil {
		w.setErr(w.putNullSigned(bitLength))
		return
	}
	raw := int64(math.Round(*value / resolution))
	w.setErr(w.putSignedNumber(raw, bitLength))
}

func (w *PGNDataStreamWriter) writeUInt8(value *uint8, bitLength uint16) {
	if value == nil {
		w.setErr(w.putNullUnsigned(bitLength))
		return
	}
	w.setErr(w.putNumberRaw(uint64(*value), bitLength))
}

func (w *PGNDataStreamWriter) writeUInt16(value *uint16, bitLength uint16) {
	if value == nil {
		w.setErr(w.putNullUnsigned(bitLength))
		return
	}
	w.setErr(w.putNumberRaw(uint64(*value), bitLength))
}

func (w *PGNDataStreamWriter) writeUInt32(value *uint32, bitLength uint16) {
	if value == nil {
		w.setErr(w.putNullUnsigned(bitLength))
		return
	}
	w.setErr(w.putNumberRaw(uint64(*value), bitLength))
}

func (w *PGNDataStreamWriter) writeUInt64(value *uint64, bitLength uint16) {
	if value == nil {
		w.setErr(w.putNullUnsigned(bitLength))
		return
	}
	w.setErr(w.putNumberRaw(*value, bitLength))
}

func (w *PGNDataStreamWriter) writeInt8(value *int8, bitLength uint16) {
	if value == nil {
		w.setErr(w.putNullSigned(bitLength))
		return
	}
	w.setErr(w.putSignedNumber(int64(*value), bitLength))
}

func (w *PGNDataStreamWriter) writeInt16(value *int16, bitLength uint16) {
	if value == nil {
		w.setErr(w.putNullSigned(bitLength))
		return
	}
	w.setErr(w.putSignedNumber(int64(*value), bitLength))
}

func (w *PGNDataStreamWriter) writeInt32(value *int32, bitLength uint16) {
	if value == nil {
		w.setErr(w.putNullSigned(bitLength))
		return
	}
	w.setErr(w.putSignedNumber(int64(*value), bitLength))
}

func (w *PGNDataStreamWriter) writeInt64(value *int64, bitLength uint16) {
	if value == nil {
		w.setErr(w.putNullSigned(bitLength))
		return
	}
	w.setErr(w.putSignedNumber(*value, bitLength))
}

func (w *PGNDataStreamWriter) writeFloat32(value *float32) {
	if value == nil {
		w.setErr(w.putNullUnsigned(32))
		return
	}
	bits := math.Float32bits(*value)
	w.setErr(w.putNumberRaw(uint64(bits), 32))
}

func (w *PGNDataStreamWriter) writeBinaryData(data []uint8, bitLength uint16) {
	idx := 0
	for i := uint16(0); i < bitLength; i += 64 {
		num := uint16(64)
		if bitLength-i < 64 {
			num = bitLength - i
		}
		var value uint64
		for h := uint16(0); h < num; h += 8 {
			if idx < len(data) {
				bitsInThisByte := num - h
				if bitsInThisByte > 8 {
					bitsInThisByte = 8
				}
				b := data[idx]
				if bitsInThisByte < 8 {
					b &= uint8((1 << bitsInThisByte) - 1)
				}
				value |= uint64(b) << h
				idx++
			}
		}
		w.setErr(w.putNumberRaw(value, num))
	}
}

func (w *PGNDataStreamWriter) writeFixedString(s string, bitLength uint16) {
	numBytes := bitLength / 8
	buf := make([]uint8, numBytes)
	sBytes := []uint8(s)
	copy(buf, sBytes)
	for i := len(sBytes); i < int(numBytes); i++ {
		buf[i] = 0xFF
	}
	w.writeBinaryData(buf, bitLength)
}

func (w *PGNDataStreamWriter) writeStringWithLength(s string) {
	sBytes := []uint8(s)
	length := uint8(len(sBytes))
	w.setErr(w.putNumberRaw(uint64(length), 8))
	w.writeBinaryData(sBytes, uint16(length)*8)
}

func (w *PGNDataStreamWriter) writeVariableData(data []uint8, pgn uint32, manID ManufacturerCodeConst, fieldIndex uint8) {
	field, err := GetFieldDescriptor(pgn, manID, fieldIndex)
	if err != nil {
		w.setErr(err)
		return
	}
	if field.BitLengthVariable && field.CanboatType == "STRING_LAU" {
		w.writeStringWithLengthAndControl(string(data))
		return
	}
	bitLen := (field.BitLength + 7) &^ 0x7
	w.writeBinaryData(data, bitLen)
}

func (w *PGNDataStreamWriter) writeStringWithLengthAndControl(s string) {
	sBytes := []uint8(s)
	totalLength := uint8(len(sBytes) + 3)
	controlByte := uint8(1)
	header := []uint8{totalLength, controlByte}
	w.writeBinaryData(header, 16)
	payload := make([]uint8, len(sBytes)+1)
	copy(payload, sBytes)
	payload[len(sBytes)] = 0
	w.writeBinaryData(payload, uint16(len(payload))*8)
}
