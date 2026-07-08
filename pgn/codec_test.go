package pgn

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// vesselHeadingPayload is a hand-built PGN 127250 payload:
//
//	byte 0      SID              = 5
//	bytes 1-2   Heading          = 15708 (0x3D5C, little endian)
//	bytes 3-4   Deviation        = -523  (0xFDF5 two's complement)
//	bytes 5-6   Variation        = 349   (0x015D)
//	byte 7      Reference (2 bits) = 1, reserved (6 bits, all ones)
var vesselHeadingPayload = []uint8{0x05, 0x5C, 0x3D, 0xF5, 0xFD, 0x5D, 0x01, 0xFD}

func TestDecodeFieldsVesselHeading(t *testing.T) {
	var msg VesselHeading
	require.NoError(t, decodeFields(&msg, vesselHeadingPayload))

	require.Equal(t, uint64(5), *msg.Sid)
	require.Equal(t, uint64(15708), *msg.Heading)
	require.Equal(t, int64(-523), *msg.Deviation)
	require.Equal(t, int64(349), *msg.Variation)
	require.Equal(t, uint64(1), *msg.Reference)

	encoded, err := encodeFields(&msg)
	require.NoError(t, err)
	require.Equal(t, vesselHeadingPayload, encoded)
}

func TestDecodeFieldsVesselHeadingNullSentinels(t *testing.T) {
	// Every numeric field carries its null sentinel: unsigned all-ones,
	// signed positive maximum. The 2-bit Reference lookup keeps its raw
	// value because small lookups have no null convention.
	payload := []uint8{0xFF, 0xFF, 0xFF, 0xFF, 0x7F, 0xFF, 0x7F, 0xFF}

	var msg VesselHeading
	require.NoError(t, decodeFields(&msg, payload))

	require.Nil(t, msg.Sid)
	require.Nil(t, msg.Heading)
	require.Nil(t, msg.Deviation)
	require.Nil(t, msg.Variation)
	require.Equal(t, uint64(3), *msg.Reference)

	// Byte fixpoint: re-encoding the decoded struct reproduces the payload.
	encoded, err := encodeFields(&msg)
	require.NoError(t, err)
	require.Equal(t, payload, encoded)
}

func TestDecodeFieldsEOFLeavesTrailingFieldsNil(t *testing.T) {
	var msg VesselHeading
	require.NoError(t, decodeFields(&msg, vesselHeadingPayload[:3]))

	require.Equal(t, uint64(5), *msg.Sid)
	require.Equal(t, uint64(15708), *msg.Heading)
	require.Nil(t, msg.Deviation)
	require.Nil(t, msg.Variation)
	require.Nil(t, msg.Reference)
}

// bgKeyValuePayload is a hand-built PGN 130824 (B&G key-value) payload:
//
//	bytes 0-1   Manufacturer Code = 381 (11 bits), reserved (2 bits, ones),
//	            Industry Code = 4 (3 bits) => 0x7D 0x99
//	bytes 2-3   Key = 1 (12 bits), Length = 2 (4 bits) => 0x01 0x20
//	bytes 4-5   Value = 0x10 0x20 (Length counts bytes)
var bgKeyValuePayload = []uint8{0x7D, 0x99, 0x01, 0x20, 0x10, 0x20}

func TestDecodeFieldsMatchVariantSelection(t *testing.T) {
	var msg BGKeyValueData
	require.NoError(t, decodeFields(&msg, bgKeyValuePayload))
	require.Equal(t, uint64(381), *msg.ManufacturerCode)
	require.Equal(t, uint64(4), *msg.IndustryCode)

	// Same layout with Manufacturer Code = 419 must be rejected so dispatch
	// can try the next candidate struct for the shared PGN.
	mismatched := append([]uint8(nil), bgKeyValuePayload...)
	mismatched[0] = 0xA3
	mismatched[1] = 0x99

	var wrong BGKeyValueData
	err := decodeFields(&wrong, mismatched)
	require.Error(t, err)
	require.Contains(t, err.Error(), "match failed for B&G: key-value data")
}

func TestDecodeFieldsKeyValueLengthDrivenBinary(t *testing.T) {
	var msg BGKeyValueData
	require.NoError(t, decodeFields(&msg, bgKeyValuePayload))

	require.Len(t, msg.Repeating1, 1)
	require.Equal(t, uint64(1), *msg.Repeating1[0].Key)
	require.Equal(t, uint64(2), *msg.Repeating1[0].Length)
	require.Equal(t, []uint8{0x10, 0x20}, msg.Repeating1[0].Value)

	encoded, err := encodeFields(&msg)
	require.NoError(t, err)
	require.Equal(t, bgKeyValuePayload, encoded)
}

func TestEncodeFieldsDerivesDynamicLengthFromValue(t *testing.T) {
	// A stale user-set Length must be overridden by len(Value), and nil
	// match fields must encode their Match values.
	msg := BGKeyValueData{
		Repeating1: []BGKeyValueDataRepeating1{
			{Key: ptrUint64(1), Length: ptrUint64(9), Value: []uint8{0x10, 0x20}},
		},
	}

	encoded, err := encodeFields(&msg)
	require.NoError(t, err)
	require.Equal(t, bgKeyValuePayload, encoded)
}

// gnssSatsInViewPayload is a hand-built PGN 129540 payload with two satellites:
//
//	byte 0      SID = 1
//	byte 1      Range Residual Mode = 0 (2 bits), reserved (6 bits, ones)
//	byte 2      Sats in View = 2 (repeating set count field)
//	bytes 3-14  satellite 1: PRN=1 Elevation=1000 Azimuth=2000 SNR=345
//	            RangeResiduals=100 Status=2 reserved(4 bits, ones)
//	bytes 15-26 satellite 2: PRN=2 Elevation=-1000 Azimuth=3000 SNR=500
//	            RangeResiduals=-100 Status=1 reserved(4 bits, ones)
var gnssSatsInViewPayload = []uint8{
	0x01, 0xFC, 0x02,
	0x01, 0xE8, 0x03, 0xD0, 0x07, 0x59, 0x01, 0x64, 0x00, 0x00, 0x00, 0xF2,
	0x02, 0x18, 0xFC, 0xB8, 0x0B, 0xF4, 0x01, 0x9C, 0xFF, 0xFF, 0xFF, 0xF1,
}

func TestDecodeFieldsRepeatingSet(t *testing.T) {
	var msg GnssSatsInView
	require.NoError(t, decodeFields(&msg, gnssSatsInViewPayload))

	require.Equal(t, uint64(1), *msg.Sid)
	require.Equal(t, uint64(0), *msg.RangeResidualMode)
	require.Equal(t, uint64(2), *msg.SatsInView)
	require.Len(t, msg.Repeating1, 2)

	first := msg.Repeating1[0]
	require.Equal(t, uint64(1), *first.Prn)
	require.Equal(t, int64(1000), *first.Elevation)
	require.Equal(t, uint64(2000), *first.Azimuth)
	require.Equal(t, int64(345), *first.Snr)
	require.Equal(t, int64(100), *first.RangeResiduals)
	require.Equal(t, uint64(2), *first.Status)

	second := msg.Repeating1[1]
	require.Equal(t, uint64(2), *second.Prn)
	require.Equal(t, int64(-1000), *second.Elevation)
	require.Equal(t, uint64(3000), *second.Azimuth)
	require.Equal(t, int64(500), *second.Snr)
	require.Equal(t, int64(-100), *second.RangeResiduals)
	require.Equal(t, uint64(1), *second.Status)

	encoded, err := encodeFields(&msg)
	require.NoError(t, err)
	require.Equal(t, gnssSatsInViewPayload, encoded)
}

func TestEncodeFieldsDerivesCountFromSliceLength(t *testing.T) {
	var msg GnssSatsInView
	require.NoError(t, decodeFields(&msg, gnssSatsInViewPayload))

	// A stale count value must be overridden by len(Repeating1) on encode.
	msg.SatsInView = ptrUint64(9)
	encoded, err := encodeFields(&msg)
	require.NoError(t, err)
	require.Equal(t, gnssSatsInViewPayload, encoded)
}

func TestDecodeFieldsDiscardsPartialRepeatingGroup(t *testing.T) {
	// Truncate the payload in the middle of the second satellite: the
	// partial group is discarded and decoding ends cleanly.
	truncated := gnssSatsInViewPayload[:21]

	var msg GnssSatsInView
	require.NoError(t, decodeFields(&msg, truncated))
	require.Len(t, msg.Repeating1, 1)
	require.Equal(t, uint64(1), *msg.Repeating1[0].Prn)
}

// channelSourceConfigurationPayload is a hand-built PGN 130061 (Channel
// Source Configuration) payload, derived bit-by-bit from the field
// descriptors in upstream_definitions.go (Order 1-11, offsets from that
// table):
//
//	Order 1  DataSourceChannelId              8 bits  @0    = 0x11
//	Order 2  SourceSelectionStatus             2 bits @8    = 0b01
//	Order 3  Reserved                          2 bits @10   = 0b11 (encoder
//	                                            always writes reserved bits
//	                                            as all-ones, so the input
//	                                            payload must already carry
//	                                            ones there for a byte-exact
//	                                            round trip)
//	Order 4  NameSelectionCriteriaMask (BINARY) 12 bits @12  = 0xA5C
//	                                            -- NOT byte-aligned: the
//	                                            fixed BitLength is 12, and
//	                                            the next field starts at
//	                                            offset 24 (12+12=24). Byte-
//	                                            aligning this field's width
//	                                            to 16 bits, as the plan
//	                                            compiler used to do, would
//	                                            overrun into the next field
//	                                            by 4 bits on both decode and
//	                                            encode.
//	Order 5  SourceName (ISO_NAME)            64 bits @24   = 0x1122334455667788
//	                                            -- immediately follows the
//	                                            binary field, so a correct
//	                                            decode of this value proves
//	                                            the read cursor landed
//	                                            exactly at bit 24, not 28.
//	Order 6  Pgn                              24 bits @88   = 130061 (0x01FC0D)
//	Order 7  DataSourceInstanceFieldNumber     8 bits @112  = 3
//	Order 8  DataSourceInstanceValue           8 bits @120  = 7
//	Order 9  SecondaryEnumerationFieldNumber   8 bits @128  = 9
//	Order 10 SecondaryEnumerationFieldValue    8 bits @136  = 42 (0x2A)
//	Order 11 ParameterFieldNumber              8 bits @144  = 5
//
// Packing every field's value into a 152-bit little-endian bit stream
// (value << bitOffset, OR'd together, then split into 19 little-endian
// bytes) yields the bytes below. NameSelectionCriteriaMask's 12 bits are
// the upper nibble of byte 1 (0xC) plus all of byte 2 (0xA5): 0xC |
// (0xA5<<4) = 0xA5C.
var channelSourceConfigurationPayload = []uint8{
	0x11, 0xCD, 0xA5, 0x88, 0x77, 0x66, 0x55, 0x44,
	0x33, 0x22, 0x11, 0x0D, 0xFC, 0x01, 0x03, 0x07,
	0x09, 0x2A, 0x05,
}

// TestDecodeFieldsNonByteAlignedBinaryField pins bit-exact handling of a
// fixed-width BINARY field whose BitLength is not a multiple of 8
// (NameSelectionCriteriaMask, 12 bits). Byte-aligning that width (the old
// behavior) advances the decode cursor 4 bits too far, so SourceName --
// the very next field -- would decode from the wrong bit position, and
// re-encoding would diverge from the original wire bytes for the rest of
// the message.
func TestDecodeFieldsNonByteAlignedBinaryField(t *testing.T) {
	var msg ChannelSourceConfiguration
	require.NoError(t, decodeFields(&msg, channelSourceConfigurationPayload))

	// The binary field itself: readBinaryData(12) returns ceil(12/8)=2
	// bytes, LSB first, with the unused upper nibble of the second byte
	// zero-padded: 0xA5C -> {0x5C, 0x0A}.
	require.Equal(t, []uint8{0x5C, 0x0A}, msg.NameSelectionCriteriaMask)

	// The field immediately after the binary field: correct only if the
	// cursor advanced by exactly 12 bits, not 16.
	require.Equal(t, uint64(0x1122334455667788), *msg.SourceName)

	// Fields further downstream, to confirm the rest of the message stays
	// aligned too.
	require.Equal(t, uint64(130061), *msg.Pgn)
	require.Equal(t, uint64(3), *msg.DataSourceInstanceFieldNumber)
	require.Equal(t, uint64(7), *msg.DataSourceInstanceValue)
	require.Equal(t, uint64(9), *msg.SecondaryEnumerationFieldNumber)
	require.Equal(t, uint64(42), *msg.SecondaryEnumerationFieldValue)
	require.Equal(t, uint64(5), *msg.ParameterFieldNumber)

	// Wire-format fidelity: re-encoding must reproduce the exact input
	// bytes, which only holds if the binary field was written back at its
	// raw 12-bit width instead of a padded 16-bit width.
	encoded, err := encodeFields(&msg)
	require.NoError(t, err)
	require.Equal(t, channelSourceConfigurationPayload, encoded)
}

// aisDgnssBroadcastBinaryMessagePayload is a hand-built PGN 129792 (AIS DGNSS
// Broadcast Binary Message) payload, derived bit-by-bit from the field
// descriptors in upstream_definitions.go:
//
//	Order 1  MessageId (LOOKUP)                6 bits  @0    = 17
//	Order 2  RepeatIndicator (NUMBER)           2 bits  @6    = 1
//	Order 3  SourceId (MMSI, unsigned number)  32 bits  @8    = 123456789
//	Order 4  Reserved                           1 bit   @40   = 1 (all-ones)
//	Order 5  AisTransceiverInformation (LOOKUP) 5 bits  @41   = 0
//	Order 6  Spare                              2 bits  @46   = 0
//	Order 7  Longitude (signed number)         32 bits  @48   = 1000000
//	Order 8  Latitude (signed number)          32 bits  @80   = -2000000
//	Order 9  Reserved                           3 bits  @112  = 0b111 (all-ones)
//	Order 10 Spare                              5 bits  @115  = 0
//	Order 11 NumberOfBitsInBinaryDataField      16 bits  @120  = 16 -- this is
//	         the BitLengthField referenced by Order 12: its decoded raw value
//	         is used directly as the bit width of the binary field that
//	         follows (not multiplied by 8 -- the field is literally named
//	         "Number of Bits").
//	Order 12 BinaryData (BINARY, BitLengthField=11) @136 = 16 bits (2 bytes)
//	         = {0xDE, 0xAD}
//
// Packing every fixed field's value into a 136-bit little-endian bit stream
// (value << bitOffset, OR'd together, split into 17 little-endian bytes),
// then appending the 2-byte binary data field (byte-aligned, so it needs no
// further bit packing), yields the 19-byte payload below.
var aisDgnssBroadcastBinaryMessagePayload = []uint8{
	0x51, 0x15, 0xCD, 0x5B, 0x07, 0x01, 0x40, 0x42, 0x0F, 0x00,
	0x80, 0x7B, 0xE1, 0xFF, 0x07, 0x10, 0x00, 0xDE, 0xAD,
}

// TestDecodeFieldsBitLengthFieldDrivenBinary pins the BitLengthField contract:
// a variable-length binary field whose width is the raw decoded value (in
// BITS) of an earlier field, as opposed to the DYNAMIC_FIELD_LENGTH contract
// (pgn/codec_test.go's KeyValueLengthDrivenBinary tests) whose referenced
// field counts BYTES. PGN 129792's "Number of Bits in Binary Data Field" is
// the canonical example: its name says exactly what it holds.
func TestDecodeFieldsBitLengthFieldDrivenBinary(t *testing.T) {
	var msg AisDgnssBroadcastBinaryMessage
	require.NoError(t, decodeFields(&msg, aisDgnssBroadcastBinaryMessagePayload))

	require.Equal(t, uint64(17), *msg.MessageId)
	require.Equal(t, uint64(1), *msg.RepeatIndicator)
	require.Equal(t, uint64(123456789), *msg.SourceId)
	require.Equal(t, uint64(0), *msg.AisTransceiverInformation)
	require.Equal(t, int64(1000000), *msg.Longitude)
	require.Equal(t, int64(-2000000), *msg.Latitude)
	require.Equal(t, uint64(16), *msg.NumberOfBitsInBinaryDataField)
	require.Equal(t, []uint8{0xDE, 0xAD}, msg.BinaryData)

	// Byte-identical re-encode: the BitLengthField reference must round-trip
	// both the referenced field's value and the binary data's exact width.
	encoded, err := encodeFields(&msg)
	require.NoError(t, err)
	require.Equal(t, aisDgnssBroadcastBinaryMessagePayload, encoded)
}

// codecTestOrphan implements PGN but has no registered field metadata.
type codecTestOrphan struct{ Info MessageInfo }

func (m *codecTestOrphan) PGNNumber() uint32               { return 0 }
func (m *codecTestOrphan) MessageInfo() MessageInfo        { return m.Info }
func (m *codecTestOrphan) SetMessageInfo(info MessageInfo) { m.Info = info }
func (m *codecTestOrphan) DecodePayload(payload []uint8) error {
	return decodeFields(m, payload)
}
func (m *codecTestOrphan) EncodePayload() ([]uint8, error) { return encodeFields(m) }

func TestDecodeFieldsUnknownStructErrors(t *testing.T) {
	err := decodeFields(&codecTestOrphan{}, []uint8{0x00})
	require.Error(t, err)
	require.Contains(t, err.Error(), "codecTestOrphan")
}

func TestStructTypeRegistryInstantiatesStructs(t *testing.T) {
	factory, ok := structTypeRegistry["VesselHeading"]
	require.True(t, ok)
	_, isVesselHeading := factory().(*VesselHeading)
	require.True(t, isVesselHeading)
}

func TestDebugDumpPGNRepeatingSlice(t *testing.T) {
	var msg GnssSatsInView
	require.NoError(t, decodeFields(&msg, gnssSatsInViewPayload))

	dump := DebugDumpPGN(&msg)
	require.Contains(t, dump, "GnssSatsInView:")
	require.Contains(t, dump, "Repeating1")
	require.Contains(t, dump, "Prn=")
}
