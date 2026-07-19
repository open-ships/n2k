package pgn

import "testing"

// isoRequestPayload is a real-world PGN 59904 (IsoRequest) capture: the
// schema declares only a 3-byte "requested PGN" field (Length: 3), but the
// device padded the single-frame CAN message to 8 bytes with 0x00 filler
// that isn't part of any declared field at all -- not RESERVED, not SPARE,
// nothing a field-by-field encode could reconstruct.
var isoRequestPayload = []uint8{0x00, 0xEE, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

// TestEncodeFieldsPassthroughForUndeclaredTrailingBytes proves that decoding
// and re-encoding reproduces bytes the schema has no field for at all, by
// returning the original wire payload verbatim -- this is the only case
// field-by-field encoding could never close, since there's no field
// descriptor to hang a fallback off of.
func TestEncodeFieldsPassthroughForUndeclaredTrailingBytes(t *testing.T) {
	msg := &IsoRequest{}
	if err := msg.DecodePayload(isoRequestPayload); err != nil {
		t.Fatalf("decode: %v", err)
	}

	encoded, err := msg.EncodePayload()
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	if !bytesEqual(encoded, isoRequestPayload) {
		t.Fatalf("byte fixpoint mismatch:\n wire:    % x\n encoded: % x", isoRequestPayload, encoded)
	}
}

// TestEncodeFieldsHonorsFieldMutationAfterDecode proves decoded messages stay
// byte-exact while untouched but switch to field encoding after an edit.
func TestEncodeFieldsHonorsFieldMutationAfterDecode(t *testing.T) {
	msg := &IsoRequest{}
	if err := msg.DecodePayload(isoRequestPayload); err != nil {
		t.Fatalf("decode: %v", err)
	}

	newPgn := uint64(126996)
	msg.Pgn = &newPgn
	encoded, err := msg.EncodePayload()
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	want := []uint8{0x14, 0xF0, 0x01}
	if !bytesEqual(encoded, want) {
		t.Fatalf("expected mutation to be encoded:\n want:    % x\n encoded: % x", want, encoded)
	}
}

func TestDecodePayloadOwnsWireBuffer(t *testing.T) {
	payload := append([]uint8(nil), isoRequestPayload...)
	msg := &IsoRequest{}
	if err := msg.DecodePayload(payload); err != nil {
		t.Fatalf("decode: %v", err)
	}
	payload[0] = 0
	encoded, err := msg.EncodePayload()
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	if !bytesEqual(encoded, isoRequestPayload) {
		t.Fatalf("decoded message retained caller-owned payload: % x", encoded)
	}
}

// TestEncodeFieldsFreshMessageUsesDefaultFill proves the byte-preserving
// behavior only applies to a message that came from DecodePayload: a struct
// built from scratch (the normal way to construct an outbound message,
// covering the manually-crafted use case) has no stashed payload and always
// goes through the normal per-field encode with the library's default fill.
func TestEncodeFieldsFreshMessageUsesDefaultFill(t *testing.T) {
	msg := &VesselHeading{}
	if err := msg.DecodePayload(reservedBitsPayload); err != nil {
		t.Fatalf("decode: %v", err)
	}
	// A struct built directly (never decoded) with the same field values
	// should encode identically to before this feature existed: reserved
	// bits filled all-ones, regardless of what any other message's wire
	// bytes looked like.
	fresh := &VesselHeading{
		Sid:       msg.Sid,
		Heading:   msg.Heading,
		Deviation: msg.Deviation,
		Variation: msg.Variation,
		Reference: msg.Reference,
	}
	encoded, err := fresh.EncodePayload()
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	if encoded[7] != 0xFD {
		t.Fatalf("expected fresh struct to fill reserved bits as all-ones (0xFD), got %#x", encoded[7])
	}
}

// reservedBitsPayload is vesselHeadingPayload with byte 7's reserved 6 bits
// set to a non-conforming pattern (0x07, mirroring the PGN 127245 Rudder
// finding from real capture logs) instead of the all-ones convention
// EncodePayload's default fill produces.
//
//	byte 7: Reference (2 bits) = 1, reserved (6 bits) = 0b000111
var reservedBitsPayload = []uint8{0x05, 0x5C, 0x3D, 0xF5, 0xFD, 0x5D, 0x01, 0x1D}

// TestEncodeFieldsReservedBitsFromWire proves that decoding a message and
// immediately re-encoding it (via the ordinary DecodePayload/EncodePayload
// pair -- no separate API) reproduces the original wire bytes even when the
// reserved bits don't follow the library's own all-ones fill convention.
func TestEncodeFieldsReservedBitsFromWire(t *testing.T) {
	msg := &VesselHeading{}
	if err := msg.DecodePayload(reservedBitsPayload); err != nil {
		t.Fatalf("decode: %v", err)
	}

	encoded, err := msg.EncodePayload()
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	if !bytesEqual(encoded, reservedBitsPayload) {
		t.Fatalf("byte fixpoint mismatch:\n wire:    % x\n encoded: % x", reservedBitsPayload, encoded)
	}
}

// spareBitsPayload is a hand-built PGN 129797 (AisBinaryBroadcastMessage)
// payload where byte 5 packs a RESERVED bit (order 4, value 1) and a SPARE
// field (order 6, 2 bits, value 0b11) both set to non-conforming patterns
// (the library's defaults are all-ones for RESERVED, zero for SPARE),
// mirroring the PGN 129041/129810 SPARE findings from real capture logs.
//
//	byte 0      messageId (6 bits) = 8, repeatIndicator (2 bits) = 1
//	bytes 1-4   sourceId (MMSI)    = 123456789
//	byte 5      reserved (1 bit) = 1, aisTransceiverInfo (5 bits) = 5, spare (2 bits) = 0b11
//	bytes 6-7   numberOfBitsInBinaryDataField = 0
var spareBitsPayload = []uint8{0x48, 0x15, 0xCD, 0x5B, 0x07, 0xCB, 0x00, 0x00}

func TestEncodeFieldsSpareBitsFromWire(t *testing.T) {
	msg := &AisBinaryBroadcastMessage{}
	if err := msg.DecodePayload(spareBitsPayload); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if *msg.AisTransceiverInformation != 5 {
		t.Fatalf("unexpected decode: aisTransceiverInformation=%v", *msg.AisTransceiverInformation)
	}

	encoded, err := msg.EncodePayload()
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	if !bytesEqual(encoded, spareBitsPayload) {
		t.Fatalf("byte fixpoint mismatch:\n wire:    % x\n encoded: % x", spareBitsPayload, encoded)
	}
}

// controlStringPayload is a hand-built PGN 126998 (ConfigurationInformation)
// payload with three STRING_LAU (length + control byte prefixed) fields,
// exercising every case readStringWithLengthAndControl/
// writeStringWithLengthAndControl need to agree on:
//
//	field 1 (installationDescription1): totalLength=6, control=2 (non-default),
//	    content="AB"+NUL+0xFF filler -- decodes to "AB", trimming 2 filler bytes
//	field 2 (installationDescription2): totalLength=2, control=1, no content
//	    -- decodes to "" (the empty-string edge case)
//	field 3 (manufacturerInformation): totalLength=4, control=1, content="XY"
//	    with no filler at all -- already an exact fixpoint under the default
//	    encoding, included to prove the fix doesn't disturb the case that
//	    already worked
var controlStringPayload = []uint8{
	0x06, 0x02, 0x41, 0x42, 0x00, 0xFF, // "AB" + NUL + 0xFF filler, control=2
	0x02, 0x01, // empty string
	0x04, 0x01, 0x58, 0x59, // "XY", no filler
}

func TestEncodeFieldsControlStringFromWire(t *testing.T) {
	msg := &ConfigurationInformation{}
	if err := msg.DecodePayload(controlStringPayload); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if msg.InstallationDescription1 != "AB" || msg.InstallationDescription2 != "" || msg.ManufacturerInformation != "XY" {
		t.Fatalf("unexpected decode: d1=%q d2=%q mfg=%q",
			msg.InstallationDescription1, msg.InstallationDescription2, msg.ManufacturerInformation)
	}

	encoded, err := msg.EncodePayload()
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	if !bytesEqual(encoded, controlStringPayload) {
		t.Fatalf("byte fixpoint mismatch:\n wire:    % x\n encoded: % x", controlStringPayload, encoded)
	}
}

// datumPayload is a hand-built PGN 129044 (Datum) payload with two STRING_FIX
// fields padded using two different device conventions: '@' (0x40) and NUL
// (0x00), both distinct from EncodePayload's default 0xFF fill.
//
//	bytes 0-3   localDatum     = "W84" + '@' pad
//	bytes 4-15  delta lat/lon/alt = 0
//	bytes 16-19 referenceDatum = "NAD" + NUL pad
var datumPayload = []uint8{
	0x57, 0x38, 0x34, 0x40, // "W84@"
	0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00,
	0x4E, 0x41, 0x44, 0x00, // "NAD\x00"
}

func TestEncodeFieldsStringPaddingFromWire(t *testing.T) {
	msg := &Datum{}
	if err := msg.DecodePayload(datumPayload); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if msg.LocalDatum != "W84" || msg.ReferenceDatum != "NAD" {
		t.Fatalf("unexpected decode: localDatum=%q referenceDatum=%q", msg.LocalDatum, msg.ReferenceDatum)
	}

	encoded, err := msg.EncodePayload()
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	if !bytesEqual(encoded, datumPayload) {
		t.Fatalf("byte fixpoint mismatch:\n wire:    % x\n encoded: % x", datumPayload, encoded)
	}
}

// TestEncodeFieldsShortPayloadPassthroughWhenUnmodified decodes a
// legitimately short single-frame payload (NMEA 2000 allows trailing fields
// to simply be absent rather than explicitly filled) and confirms re-encode
// reproduces that same short payload -- the most faithful reproduction of
// what was actually on the wire -- rather than padding it back out to the
// struct's full field-declared width.
func TestEncodeFieldsShortPayloadPassthroughWhenUnmodified(t *testing.T) {
	short := reservedBitsPayload[:5] // Sid, Heading, Deviation only -- a clean field boundary
	msg := &VesselHeading{}
	if err := msg.DecodePayload(short); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if msg.Variation != nil || msg.Reference != nil {
		t.Fatalf("expected trailing fields nil after short decode, got Variation=%v Reference=%v", msg.Variation, msg.Reference)
	}

	encoded, err := msg.EncodePayload()
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	if !bytesEqual(encoded, short) {
		t.Fatalf("byte fixpoint mismatch:\n wire:    % x\n encoded: % x", short, encoded)
	}
}

// TestEncodeFieldsSetMessageInfoClearsStashedPayload proves that replacing a
// decoded message's MessageInfo wholesale (the normal way to prep a message
// for retransmission with a fresh timestamp/source) drops the stashed wire
// payload, moving the message onto the normal, manually-crafted-style encode
// path -- the only supported way back from "decoded" to "encode reflects
// current field values" for a message that was once decoded.
func TestEncodeFieldsSetMessageInfoClearsStashedPayload(t *testing.T) {
	msg := &Datum{}
	if err := msg.DecodePayload(datumPayload); err != nil {
		t.Fatalf("decode: %v", err)
	}

	msg.SetMessageInfo(MessageInfo{PGN: msg.PGNNumber()})
	msg.LocalDatum = "AB"
	encoded, err := msg.EncodePayload()
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	if encoded[0] != 'A' || encoded[1] != 'B' {
		t.Fatalf("expected new field value 'AB' after stash was cleared, got % x", encoded[0:2])
	}
	if encoded[2] != 0xFF || encoded[3] != 0xFF {
		t.Fatalf("expected default 0xFF padding after SetMessageInfo cleared the stashed payload, got % x", encoded[2:4])
	}
}

func bytesEqual(a, b []uint8) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
