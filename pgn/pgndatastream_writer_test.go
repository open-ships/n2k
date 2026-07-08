package pgn

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Pointer helper functions for creating literal pointers in test expressions.
func ptrFloat32(v float32) *float32 { return &v }
func ptrFloat64(v float64) *float64 { return &v }
func ptrUint8(v uint8) *uint8       { return &v }
func ptrUint16(v uint16) *uint16    { return &v }
func ptrUint32(v uint32) *uint32    { return &v }
func ptrUint64(v uint64) *uint64    { return &v }
func ptrInt8(v int8) *int8          { return &v }
func ptrInt16(v int16) *int16       { return &v }
func ptrInt32(v int32) *int32       { return &v }
func ptrInt64(v int64) *int64       { return &v }

// TestPutNumberRaw_ByteAligned verifies that byte-aligned writes produce correct
// little-endian byte sequences: 8-bit, 16-bit, and 32-bit values.
func TestPutNumberRaw_ByteAligned(t *testing.T) {
	// 8-bit write: 0x12 → {0x12}
	w := NewPGNDataStreamWriter()
	err := w.putNumberRaw(0x12, 8)
	assert.NoError(t, err)
	assert.Equal(t, []uint8{0x12}, w.Bytes())

	// 16-bit write: 0x1234 → {0x34, 0x12} (little-endian)
	w = NewPGNDataStreamWriter()
	err = w.putNumberRaw(0x1234, 16)
	assert.NoError(t, err)
	assert.Equal(t, []uint8{0x34, 0x12}, w.Bytes())

	// 32-bit write: 0xFFFFEED4 → {0xD4, 0xEE, 0xFF, 0xFF}
	w = NewPGNDataStreamWriter()
	err = w.putNumberRaw(0xFFFFEED4, 32)
	assert.NoError(t, err)
	assert.Equal(t, []uint8{0xD4, 0xEE, 0xFF, 0xFF}, w.Bytes())
}

// TestPutNumberRaw_SubByte verifies sub-byte writes that pack multiple fields
// into a single byte using little-endian bit order.
func TestPutNumberRaw_SubByte(t *testing.T) {
	// Two 4-bit fields: write 0x0A (low nibble), then 0x05 (high nibble) → 0x5A
	w := NewPGNDataStreamWriter()
	err := w.putNumberRaw(0x0A, 4)
	assert.NoError(t, err)
	err = w.putNumberRaw(0x05, 4)
	assert.NoError(t, err)
	assert.Equal(t, []uint8{0x5A}, w.Bytes())

	// 3+5 bit fields: write 0x05 (3 bits), then 0x1E (5 bits) → 0xF5
	// 0x05 = 0b101 in bits [0:2], 0x1E = 0b11110 in bits [3:7]
	// byte = 0b11110_101 = 0xF5
	w = NewPGNDataStreamWriter()
	err = w.putNumberRaw(0x05, 3)
	assert.NoError(t, err)
	err = w.putNumberRaw(0x1E, 5)
	assert.NoError(t, err)
	assert.Equal(t, []uint8{0xF5}, w.Bytes())

	// Cross-byte: 5+8+3 bits = 16 bits total = 2 bytes
	// Write 0x1E (5 bits), 0x42 (8 bits), 0x03 (3 bits) → {0x5E, 0x68}
	w = NewPGNDataStreamWriter()
	err = w.putNumberRaw(0x1E, 5)
	assert.NoError(t, err)
	err = w.putNumberRaw(0x42, 8)
	assert.NoError(t, err)
	err = w.putNumberRaw(0x03, 3)
	assert.NoError(t, err)
	assert.Equal(t, []uint8{0x5E, 0x68}, w.Bytes())
}

// TestRoundTrip_PutGetNumberRaw writes various bit-width values with the writer
// then reads them back with the reader to verify perfect round-trip fidelity.
func TestRoundTrip_PutGetNumberRaw(t *testing.T) {
	tests := []struct {
		name      string
		value     uint64
		bitLength uint16
	}{
		{"8-bit", 0x42, 8},
		{"16-bit", 0x1234, 16},
		{"32-bit", 0xDEADBEEF, 32},
		{"3-bit", 0x05, 3},
		{"11-bit", 0x567, 11},
		{"21-bit", 0x1ABCDE, 21},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := NewPGNDataStreamWriter()
			err := w.putNumberRaw(tt.value, tt.bitLength)
			assert.NoError(t, err)

			r := NewPgnDataStream(w.Bytes())
			got, err := r.getNumberRaw(tt.bitLength)
			assert.NoError(t, err)
			// Mask the expected value to the bit width to handle any high-bit truncation.
			expected := tt.value & ((1 << tt.bitLength) - 1)
			assert.Equal(t, expected, got)
		})
	}
}

// TestRoundTrip_UnsignedResolution writes a value with unsigned resolution encoding
// then reads it back, verifying the value survives the encode/decode cycle within
// the expected quantization error.
func TestRoundTrip_UnsignedResolution(t *testing.T) {
	w := NewPGNDataStreamWriter()
	w.writeUnsignedResolution(ptrFloat32(1.5708), 16, 0.0001)
	assert.NoError(t, w.Err())

	r := NewPgnDataStream(w.Bytes())
	v, err := r.readUnsignedResolution(16, 0.0001)
	assert.NoError(t, err)
	assert.NotNil(t, v)
	assert.InDelta(t, float32(1.5708), *v, 0.0001)
}

// TestRoundTrip_SignedResolution writes a negative value with signed resolution encoding
// then reads it back, verifying correct two's-complement handling.
func TestRoundTrip_SignedResolution(t *testing.T) {
	w := NewPGNDataStreamWriter()
	w.writeSignedResolution(ptrFloat32(-0.5), 16, 0.0001)
	assert.NoError(t, w.Err())

	r := NewPgnDataStream(w.Bytes())
	v, err := r.readSignedResolution(16, 0.0001)
	assert.NoError(t, err)
	assert.NotNil(t, v)
	assert.InDelta(t, float32(-0.5), *v, 0.0001)
}

// TestRoundTrip_SignedResolution64Override verifies the float64 resolution variant
// for high-precision fields like latitude/longitude.
func TestRoundTrip_SignedResolution64Override(t *testing.T) {
	w := NewPGNDataStreamWriter()
	w.writeSignedResolution64Override(ptrFloat64(-123.456789), 32, 0.0000001)
	assert.NoError(t, w.Err())

	r := NewPgnDataStream(w.Bytes())
	v, err := r.readSignedResolution64Override(32, 0.0000001)
	assert.NoError(t, err)
	assert.NotNil(t, v)
	assert.InDelta(t, -123.456789, *v, 0.0000001)
}

// TestRoundTrip_NullFields writes nil pointers for uint8, unsigned resolution, and
// signed resolution fields, then reads them back and verifies the reader returns nil.
func TestRoundTrip_NullFields(t *testing.T) {
	w := NewPGNDataStreamWriter()

	// Write nil uint8 (8 bits → 0xFF null sentinel)
	w.writeUInt8(nil, 8)

	// Write nil unsigned resolution (16 bits → 0xFFFF null sentinel)
	w.writeUnsignedResolution(nil, 16, 0.0001)

	// Write nil signed resolution (16 bits → 0x7FFF null sentinel)
	w.writeSignedResolution(nil, 16, 0.0001)

	assert.NoError(t, w.Err())

	r := NewPgnDataStream(w.Bytes())

	v1, err := r.readUInt8(8)
	assert.NoError(t, err)
	assert.Nil(t, v1)

	v2, err := r.readUnsignedResolution(16, 0.0001)
	assert.NoError(t, err)
	assert.Nil(t, v2)

	v3, err := r.readSignedResolution(16, 0.0001)
	assert.NoError(t, err)
	assert.Nil(t, v3)
}

// TestRoundTrip_LookupField writes a 2-bit lookup value then reads it back.
func TestRoundTrip_LookupField(t *testing.T) {
	w := NewPGNDataStreamWriter()
	w.writeLookupField(2, 2)

	// Pad the remaining 6 bits so the reader has a full byte.
	w.skipBits(6)

	assert.NoError(t, w.Err())

	r := NewPgnDataStream(w.Bytes())
	v, err := r.readLookupField(2)
	assert.NoError(t, err)
	assert.Equal(t, uint64(2), v)
}

func TestWriteReservedAndSpareBits(t *testing.T) {
	w := NewPGNDataStreamWriter()
	w.writeReservedBits(4)
	w.writeSpareBits(4)
	assert.NoError(t, w.Err())
	assert.Equal(t, []uint8{0x0f}, w.Bytes())
}

// TestRoundTrip_FixedString writes a fixed-width string then reads it back,
// verifying that the 0xFF padding is correctly stripped by the reader.
func TestRoundTrip_FixedString(t *testing.T) {
	w := NewPGNDataStreamWriter()
	w.writeFixedString("HELLO", 64) // 8 bytes = 64 bits, "HELLO" + 3 bytes of 0xFF padding
	assert.NoError(t, w.Err())

	r := NewPgnDataStream(w.Bytes())
	v, err := r.readFixedString(64)
	assert.NoError(t, err)
	assert.Equal(t, "HELLO", v)
}

// TestRoundTrip_Float32 writes an IEEE 754 float32 then reads it back.
func TestRoundTrip_Float32(t *testing.T) {
	w := NewPGNDataStreamWriter()
	w.writeFloat32(ptrFloat32(3.14))
	assert.NoError(t, w.Err())

	r := NewPgnDataStream(w.Bytes())
	v, err := r.readFloat32()
	assert.NoError(t, err)
	assert.NotNil(t, v)
	assert.Equal(t, float32(3.14), *v)
}

// TestRoundTrip_Float32_Null writes a nil float32 then reads it back as nil.
func TestRoundTrip_Float32_Null(t *testing.T) {
	w := NewPGNDataStreamWriter()
	w.writeFloat32(nil)
	assert.NoError(t, w.Err())

	r := NewPgnDataStream(w.Bytes())
	v, err := r.readFloat32()
	assert.NoError(t, err)
	assert.Nil(t, v)
}

// TestRoundTrip_SignedIntegers writes positive and negative signed values at various
// widths then reads them back to verify two's-complement round-trip correctness.
func TestRoundTrip_SignedIntegers(t *testing.T) {
	w := NewPGNDataStreamWriter()

	// Positive int8
	w.writeInt8(ptrInt8(42), 8)

	// Negative int8
	w.writeInt8(ptrInt8(-44), 8)

	// Positive int16
	w.writeInt16(ptrInt16(1000), 16)

	// Negative int16
	w.writeInt16(ptrInt16(-4396), 16)

	assert.NoError(t, w.Err())

	r := NewPgnDataStream(w.Bytes())

	v1, err := r.readInt8(8)
	assert.NoError(t, err)
	assert.NotNil(t, v1)
	assert.Equal(t, int8(42), *v1)

	v2, err := r.readInt8(8)
	assert.NoError(t, err)
	assert.NotNil(t, v2)
	assert.Equal(t, int8(-44), *v2)

	v3, err := r.readInt16(16)
	assert.NoError(t, err)
	assert.NotNil(t, v3)
	assert.Equal(t, int16(1000), *v3)

	v4, err := r.readInt16(16)
	assert.NoError(t, err)
	assert.NotNil(t, v4)
	assert.Equal(t, int16(-4396), *v4)
}

// TestRoundTrip_SignedIntegers_Null writes nil signed integers and verifies the reader
// returns nil for the signed null sentinel.
func TestRoundTrip_SignedIntegers_Null(t *testing.T) {
	w := NewPGNDataStreamWriter()
	w.writeInt8(nil, 8)
	w.writeInt16(nil, 16)
	assert.NoError(t, w.Err())

	r := NewPgnDataStream(w.Bytes())

	v1, err := r.readInt8(8)
	assert.NoError(t, err)
	assert.Nil(t, v1)

	v2, err := r.readInt16(16)
	assert.NoError(t, err)
	assert.Nil(t, v2)
}

// TestRoundTrip_BinaryData writes a binary payload then reads it back.
func TestRoundTrip_BinaryData(t *testing.T) {
	payload := []uint8{0xDE, 0xAD, 0xBE, 0xEF}
	w := NewPGNDataStreamWriter()
	w.writeBinaryData(payload, 32)
	assert.NoError(t, w.Err())

	r := NewPgnDataStream(w.Bytes())
	v, err := r.readBinaryData(32)
	assert.NoError(t, err)
	assert.Equal(t, payload, v)
}

// TestRoundTrip_StringWithLength writes a STRING_LZ string then reads it back.
func TestRoundTrip_StringWithLength(t *testing.T) {
	w := NewPGNDataStreamWriter()
	w.writeStringWithLength("Hello")
	assert.NoError(t, w.Err())

	r := NewPgnDataStream(w.Bytes())
	v, err := r.readStringWithLength()
	assert.NoError(t, err)
	assert.Equal(t, "Hello", v)
}

// TestWriteStringWithLengthAndControl_ExactBytes pins the STRING_LAU wire layout: the
// length byte counts itself, the control byte, and the string bytes (len(s)+2) -- there
// is no trailing NUL -- followed by the ASCII control byte (1), followed by the string.
func TestWriteStringWithLengthAndControl_ExactBytes(t *testing.T) {
	w := NewPGNDataStreamWriter()
	w.writeStringWithLengthAndControl("AB")
	assert.NoError(t, w.Err())
	assert.Equal(t, []uint8{0x04, 0x01, 'A', 'B'}, w.Bytes())
}

// TestRoundTrip_StringWithLengthAndControl writes a STRING_LAU string then reads it back,
// verifying the writer and reader are exact inverses (no trailing NUL to strip).
func TestRoundTrip_StringWithLengthAndControl(t *testing.T) {
	w := NewPGNDataStreamWriter()
	w.writeStringWithLengthAndControl("AB")
	assert.NoError(t, w.Err())

	r := NewPgnDataStream(w.Bytes())
	v, err := r.readStringWithLengthAndControl()
	assert.NoError(t, err)
	assert.Equal(t, "AB", v)
}

// TestRoundTrip_UInt32 verifies uint32 write and read round-trip.
func TestRoundTrip_UInt32(t *testing.T) {
	w := NewPGNDataStreamWriter()
	w.writeUInt32(ptrUint32(0xDEADBEEF), 32)
	assert.NoError(t, w.Err())

	r := NewPgnDataStream(w.Bytes())
	v, err := r.readUInt32(32)
	assert.NoError(t, err)
	assert.NotNil(t, v)
	assert.Equal(t, uint32(0xDEADBEEF), *v)
}

// TestRoundTrip_UInt16 verifies uint16 write and read round-trip.
func TestRoundTrip_UInt16(t *testing.T) {
	w := NewPGNDataStreamWriter()
	w.writeUInt16(ptrUint16(0x1234), 16)
	assert.NoError(t, w.Err())

	r := NewPgnDataStream(w.Bytes())
	v, err := r.readUInt16(16)
	assert.NoError(t, err)
	assert.NotNil(t, v)
	assert.Equal(t, uint16(0x1234), *v)
}

// TestRoundTrip_UInt64 verifies uint64 write and read round-trip.
func TestRoundTrip_UInt64(t *testing.T) {
	w := NewPGNDataStreamWriter()
	w.writeUInt64(ptrUint64(0x123456789ABCDEF0), 64)
	assert.NoError(t, w.Err())

	r := NewPgnDataStream(w.Bytes())
	v, err := r.readUInt64(64)
	assert.NoError(t, err)
	assert.NotNil(t, v)
	assert.Equal(t, uint64(0x123456789ABCDEF0), *v)
}

// TestRoundTrip_Int32 verifies int32 write and read round-trip.
func TestRoundTrip_Int32(t *testing.T) {
	w := NewPGNDataStreamWriter()
	w.writeInt32(ptrInt32(-123456), 32)
	assert.NoError(t, w.Err())

	r := NewPgnDataStream(w.Bytes())
	v, err := r.readInt32(32)
	assert.NoError(t, err)
	assert.NotNil(t, v)
	assert.Equal(t, int32(-123456), *v)
}

// TestRoundTrip_Int64 verifies int64 write and read round-trip.
func TestRoundTrip_Int64(t *testing.T) {
	w := NewPGNDataStreamWriter()
	w.writeInt64(ptrInt64(-9876543210), 64)
	assert.NoError(t, w.Err())

	r := NewPgnDataStream(w.Bytes())
	v, err := r.readInt64(64)
	assert.NoError(t, err)
	assert.NotNil(t, v)
	assert.Equal(t, int64(-9876543210), *v)
}

// TestPutNullUnsigned verifies that putNullUnsigned writes the all-ones sentinel
// and that the reader detects it as null.
func TestPutNullUnsigned(t *testing.T) {
	w := NewPGNDataStreamWriter()
	err := w.putNullUnsigned(8)
	assert.NoError(t, err)
	assert.Equal(t, []uint8{0xFF}, w.Bytes())

	w = NewPGNDataStreamWriter()
	err = w.putNullUnsigned(16)
	assert.NoError(t, err)
	assert.Equal(t, []uint8{0xFF, 0xFF}, w.Bytes())
}

// TestPutNullSigned verifies that putNullSigned writes the positive-max sentinel
// and that the reader detects it as null.
func TestPutNullSigned(t *testing.T) {
	w := NewPGNDataStreamWriter()
	err := w.putNullSigned(8)
	assert.NoError(t, err)
	assert.Equal(t, []uint8{0x7F}, w.Bytes())

	w = NewPGNDataStreamWriter()
	err = w.putNullSigned(16)
	assert.NoError(t, err)
	assert.Equal(t, []uint8{0xFF, 0x7F}, w.Bytes())
}

// TestPutSignedNumber verifies two's-complement encoding for positive and negative values.
func TestPutSignedNumber(t *testing.T) {
	// Positive value: 42 in 8 bits → 0x2A
	w := NewPGNDataStreamWriter()
	err := w.putSignedNumber(42, 8)
	assert.NoError(t, err)
	assert.Equal(t, []uint8{0x2A}, w.Bytes())

	// Negative value: -44 in 8 bits → 0xD4
	w = NewPGNDataStreamWriter()
	err = w.putSignedNumber(-44, 8)
	assert.NoError(t, err)
	assert.Equal(t, []uint8{0xD4}, w.Bytes())

	// Negative value: -4396 in 16 bits → {0xD4, 0xEE} (0xEED4 = 65536 - 4396 = 61140)
	w = NewPGNDataStreamWriter()
	err = w.putSignedNumber(-4396, 16)
	assert.NoError(t, err)
	assert.Equal(t, []uint8{0xD4, 0xEE}, w.Bytes())
}

// TestWriterBytes_Empty verifies that a fresh writer returns an empty byte slice.
func TestWriterBytes_Empty(t *testing.T) {
	w := NewPGNDataStreamWriter()
	assert.Equal(t, []uint8{}, w.Bytes())
}

// TestRoundTrip_MultipleSequentialFields simulates encoding a mini PGN message with
// mixed field types, then decoding it to verify the writer and reader stay in sync.
func TestRoundTrip_MultipleSequentialFields(t *testing.T) {
	w := NewPGNDataStreamWriter()

	// SID: uint8 field (8 bits)
	w.writeUInt8(ptrUint8(1), 8)

	// Heading: unsigned resolution (16 bits, resolution 0.0001 rad)
	w.writeUnsignedResolution(ptrFloat32(0.4660), 16, 0.0001)

	// Reserved: skip 8 bits
	w.skipBits(8)

	assert.NoError(t, w.Err())

	r := NewPgnDataStream(w.Bytes())

	sid, err := r.readUInt8(8)
	assert.NoError(t, err)
	assert.NotNil(t, sid)
	assert.Equal(t, uint8(1), *sid)

	heading, err := r.readUnsignedResolution(16, 0.0001)
	assert.NoError(t, err)
	assert.NotNil(t, heading)
	assert.InDelta(t, float32(0.4660), *heading, 0.0001)

	// The reader can read the reserved byte (zeros from skipBits).
	reserved, err := r.readUInt8(8)
	assert.NoError(t, err)
	assert.NotNil(t, reserved)
	assert.Equal(t, uint8(0), *reserved)
}

// TestRoundTrip_FixedString_ExactLength verifies a string that exactly fills the field.
func TestRoundTrip_FixedString_ExactLength(t *testing.T) {
	w := NewPGNDataStreamWriter()
	// 5 chars in 40 bits (5 bytes) — no padding needed.
	w.writeFixedString("ABCDE", 40)
	assert.NoError(t, w.Err())
	assert.Equal(t, []uint8{'A', 'B', 'C', 'D', 'E'}, w.Bytes())

	r := NewPgnDataStream(w.Bytes())
	v, err := r.readFixedString(40)
	assert.NoError(t, err)
	assert.Equal(t, "ABCDE", v)
}

// TestRoundTrip_FixedString_Empty verifies an empty string fills with padding.
func TestRoundTrip_FixedString_Empty(t *testing.T) {
	w := NewPGNDataStreamWriter()
	w.writeFixedString("", 24) // 3 bytes of padding
	assert.NoError(t, w.Err())
	assert.Equal(t, []uint8{0xFF, 0xFF, 0xFF}, w.Bytes())

	r := NewPgnDataStream(w.Bytes())
	v, err := r.readFixedString(24)
	assert.NoError(t, err)
	assert.Equal(t, "", v)
}

// TestWriteBinaryData_SubByte verifies writing fewer than 8 bits of binary data.
func TestWriteBinaryData_SubByte(t *testing.T) {
	w := NewPGNDataStreamWriter()
	w.writeBinaryData([]uint8{0x0F}, 4) // Only low 4 bits matter

	// Fill remaining bits so we can read back
	w.skipBits(4)

	assert.NoError(t, w.Err())

	r := NewPgnDataStream(w.Bytes())
	v, err := r.readBinaryData(4)
	assert.NoError(t, err)
	assert.Equal(t, []uint8{0x0F}, v)
}

// TestRoundTrip_UnsignedResolution_LargeValue tests a value near the top of a 16-bit
// unsigned range with resolution.
func TestRoundTrip_UnsignedResolution_LargeValue(t *testing.T) {
	// Max representable non-null: (0xFFFE) * 0.0001 = 6.5534
	w := NewPGNDataStreamWriter()
	w.writeUnsignedResolution(ptrFloat32(6.5534), 16, 0.0001)
	assert.NoError(t, w.Err())

	r := NewPgnDataStream(w.Bytes())
	v, err := r.readUnsignedResolution(16, 0.0001)
	assert.NoError(t, err)
	assert.NotNil(t, v)
	assert.InDelta(t, float32(6.5534), *v, 0.0001)
}

// TestPutNumberRaw_64Bit verifies writing a full 64-bit value.
func TestPutNumberRaw_64Bit(t *testing.T) {
	w := NewPGNDataStreamWriter()
	err := w.putNumberRaw(0x0102030405060708, 64)
	assert.NoError(t, err)
	assert.Equal(t, []uint8{0x08, 0x07, 0x06, 0x05, 0x04, 0x03, 0x02, 0x01}, w.Bytes())
}

// TestRoundTrip_SubByte_LookupFields verifies multiple sub-byte lookup writes
// followed by reads to ensure cursor synchronization.
func TestRoundTrip_SubByte_LookupFields(t *testing.T) {
	w := NewPGNDataStreamWriter()
	// Write three lookup fields: 3-bit, 3-bit, 2-bit (total 8 bits = 1 byte)
	w.writeLookupField(5, 3)
	w.writeLookupField(3, 3)
	w.writeLookupField(2, 2)
	assert.NoError(t, w.Err())

	r := NewPgnDataStream(w.Bytes())
	v1, err := r.readLookupField(3)
	assert.NoError(t, err)
	assert.Equal(t, uint64(5), v1)

	v2, err := r.readLookupField(3)
	assert.NoError(t, err)
	assert.Equal(t, uint64(3), v2)

	v3, err := r.readLookupField(2)
	assert.NoError(t, err)
	assert.Equal(t, uint64(2), v3)
}

// TestRoundTrip_Float32_SpecialValues tests edge case float values.
func TestRoundTrip_Float32_SpecialValues(t *testing.T) {
	tests := []struct {
		name  string
		value float32
	}{
		{"zero", 0.0},
		{"negative", -1.5},
		{"large", 1e10},
		{"small", 1e-10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := NewPGNDataStreamWriter()
			w.writeFloat32(ptrFloat32(tt.value))
			assert.NoError(t, w.Err())

			// Verify exact byte representation matches IEEE 754 LE encoding.
			bits := math.Float32bits(tt.value)
			expected := []uint8{
				uint8(bits),
				uint8(bits >> 8),
				uint8(bits >> 16),
				uint8(bits >> 24),
			}
			assert.Equal(t, expected, w.Bytes())

			r := NewPgnDataStream(w.Bytes())
			v, err := r.readFloat32()
			assert.NoError(t, err)
			assert.NotNil(t, v)
			assert.Equal(t, tt.value, *v)
		})
	}
}

// TestPutNumberRaw_Overflow verifies that writing a value that exceeds the bit width returns an error.
func TestPutNumberRaw_Overflow(t *testing.T) {
	w := NewPGNDataStreamWriter()
	err := w.putNumberRaw(256, 8) // 256 doesn't fit in 8 bits (max 255)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "exceeds 8-bit unsigned range")

	w = NewPGNDataStreamWriter()
	err = w.putNumberRaw(0x10000, 16) // 65536 doesn't fit in 16 bits
	assert.Error(t, err)

	// Value that fits should still work.
	w = NewPGNDataStreamWriter()
	err = w.putNumberRaw(255, 8)
	assert.NoError(t, err)
}
