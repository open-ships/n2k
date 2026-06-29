// Package pgn converts NMEA 2000 messages to strongly typed Go data.
// It provides PGN structs, payload encode/decode methods, and field metadata.
package pgn

import (
	"fmt"
)

// PgnInfo describes a known NMEA 2000 message type. Entries are generated from
// upstream metadata and indexed into PgnInfoLookup at init time for fast access by
// PGN number.
//
// Multiple PgnInfo entries can share the same PGN number. This happens with proprietary PGNs
// where different manufacturers define different payloads for the same PGN, and also with
// "KeyValue" style PGNs that have multiple structural variants.
type PgnInfo struct {
	// SourceID is the upstream source Id for this PGN variant.
	SourceID string `json:"sourceId"`
	// Id is a unique string identifier for this PGN variant, needed to distinguish
	// PGNs that share the same numeric PGN but have different field layouts (KeyValue PGNs).
	Id string `json:"id"`
	// PGN is the NMEA 2000 Parameter Group Number that identifies this message type
	// on the CAN bus. Values range from 0 to ~131071 (0x1FFFF).
	PGN uint32 `json:"pgn"`
	// Description is a human-readable name for this PGN (e.g., "Vessel Heading").
	Description string `json:"description"`
	// Explanation is the source schema's longer PGN description, when present.
	Explanation string `json:"explanation,omitempty"`
	// Fast indicates whether this PGN uses the NMEA 2000 fast-packet protocol.
	// Fast-packet PGNs can carry more than 8 bytes by spanning multiple CAN frames;
	// single-frame PGNs are limited to 8 bytes.
	Fast bool `json:"fast"`
	// Type is the source schema transport type string (Single, Fast, ISO, or Mixed).
	Type string `json:"type"`
	// Complete mirrors the source schema's confidence flag for fully described PGNs.
	Complete bool `json:"complete"`
	// Fallback marks range fallback definitions from the source schema.
	Fallback bool `json:"fallback"`
	// Missing lists source metadata categories still missing for this PGN.
	Missing []string `json:"missing,omitempty"`
	// Length is the nominal payload length in bytes when the source schema provides one.
	Length *int `json:"length,omitempty"`
	// MinLength is the minimum payload length in bytes when the source schema provides one.
	MinLength *int `json:"minLength,omitempty"`
	// Priority is the default CAN priority when the source schema provides one.
	Priority *uint8 `json:"priority,omitempty"`
	// TransmissionInterval is the default transmission interval in milliseconds.
	TransmissionInterval *int `json:"transmissionInterval,omitempty"`
	// TransmissionIrregular marks PGNs whose transmission is event/request driven.
	TransmissionIrregular *bool `json:"transmissionIrregular,omitempty"`
	// RepeatingFieldSet metadata mirrors the source schema's repeating field-set annotations.
	RepeatingFieldSet1StartField *int `json:"repeatingFieldSet1StartField,omitempty"`
	RepeatingFieldSet1CountField *int `json:"repeatingFieldSet1CountField,omitempty"`
	RepeatingFieldSet1Size       *int `json:"repeatingFieldSet1Size,omitempty"`
	RepeatingFieldSet2StartField *int `json:"repeatingFieldSet2StartField,omitempty"`
	RepeatingFieldSet2CountField *int `json:"repeatingFieldSet2CountField,omitempty"`
	RepeatingFieldSet2Size       *int `json:"repeatingFieldSet2Size,omitempty"`
	// ManId is the manufacturer code for proprietary PGNs. It is zero for standard
	// (non-proprietary) PGNs. Used to select the correct variant when multiple
	// manufacturers define different payloads for the same proprietary PGN number.
	ManId ManufacturerCodeConst `json:"manId"`
	// Fields maps field index (1-based, matching the source field order) to
	// FieldDescriptor. This is needed at runtime for variable-length and KeyValue
	// fields where the decoder must inspect field metadata dynamically.
	Fields map[int]*FieldDescriptor `json:"fields"`
}

// FieldDescriptor holds metadata about a single field within a PGN definition.
// It is used at runtime by decode and encode helpers to handle fields whose type
// or length cannot be fully resolved at code-generation time.
type FieldDescriptor struct {
	// SourceID is the upstream source field Id.
	SourceID string `json:"sourceId"`
	// Name is the source field name (e.g., "Heading", "SID", "Manufacturer Code").
	Name string `json:"name"`
	// Description is the source schema's per-field description, when present.
	Description string `json:"description,omitempty"`
	// BitLength is the width of this field in bits. For variable-length fields,
	// this may be a nominal/default length.
	BitLength uint16 `json:"bitLength"`
	// BitLengthField is the field order that defines a variable-length field's width.
	BitLengthField *int `json:"bitLengthField,omitempty"`
	// BitOffset is the absolute bit position of this field from the start of the PGN payload.
	BitOffset uint16 `json:"bitOffset"`
	// BitLengthVariable is true when the field's actual length is determined at runtime
	// (e.g., STRING_LAU fields whose length is encoded in a preceding byte).
	BitLengthVariable bool `json:"bitLengthVariable"`
	// SourceType is the source field type string (e.g., "NUMBER", "LOOKUP", "STRING_LAU",
	// "STRING_LZ", "STRING_FIX"). It drives type-specific decoding logic.
	SourceType string `json:"sourceType"`
	// BitStart is the bit position within the containing byte when the source schema provides it.
	BitStart uint16 `json:"bitStart"`
	// PhysicalQuantity is the source schema physical quantity identifier.
	PhysicalQuantity string `json:"physicalQuantity,omitempty"`
	// LookupEnumeration is the source lookup enumeration name, when present.
	LookupEnumeration string `json:"lookupEnumeration,omitempty"`
	// LookupBitEnumeration is the source bit lookup enumeration name, when present.
	LookupBitEnumeration string `json:"lookupBitEnumeration,omitempty"`
	// LookupIndirectEnumeration is the source indirect lookup enumeration name, when present.
	LookupIndirectEnumeration string `json:"lookupIndirectEnumeration,omitempty"`
	// LookupFieldTypeEnumeration is the source field-type lookup name, when present.
	LookupFieldTypeEnumeration string `json:"lookupFieldTypeEnumeration,omitempty"`
	// LookupIndirectEnumerationFieldOrder identifies the field that selects an indirect lookup.
	LookupIndirectEnumerationFieldOrder *int `json:"lookupIndirectEnumerationFieldOrder,omitempty"`
	// Condition is the source field-level inclusion predicate, when present.
	Condition string `json:"condition,omitempty"`
	// Offset is an additive physical-value offset from the source schema.
	Offset *float64 `json:"offset,omitempty"`
	// RangeMin and RangeMax are the physical-value bounds from the source schema.
	RangeMin *float64 `json:"rangeMin,omitempty"`
	RangeMax *float64 `json:"rangeMax,omitempty"`
	// Source sentinel values for nullable and special numeric states.
	OutOfRangeValue *int64 `json:"outOfRangeValue,omitempty"`
	ReservedValue   *int64 `json:"reservedValue,omitempty"`
	UnknownValue    *int64 `json:"unknownValue,omitempty"`
	// PartOfPrimaryKey marks fields the source schema uses to identify repeated rows.
	PartOfPrimaryKey *bool `json:"partOfPrimaryKey,omitempty"`
	// GolangType is the Go type name used in the generated struct field (e.g., "*uint8", "float32").
	GolangType string `json:"golangType"`
	// Resolution is the scaling factor applied to integer fields to produce physical units
	// (e.g., 0.0001 for a heading field stored in units of 1/10000 radian).
	Resolution float32 `json:"resolution"`
	// Signed indicates whether this numeric field uses two's-complement signed encoding.
	Signed bool `json:"signed"`
	// Unit is the physical unit string from the source schema (e.g., "rad", "m/s", "K").
	Unit string `json:"unit"`
	// BitLookupName is the name of the bit-enumeration lookup table, if this field
	// is a bitfield-type enum. Empty for non-lookup fields.
	BitLookupName string `json:"bitLookupName"`
	// Match is non-nil for fields that must match a specific value to select a PGN variant.
	// For example, proprietary PGN variants use a Match on the manufacturer code field
	// to distinguish which decoder to use.
	Match *int `json:"match"`
}

// PgnInfoLookup maps PGN numbers to their PgnInfo descriptors. It is the primary
// lookup table used by callers to find metadata and, when available, decoders for a
// received message. Multiple entries per PGN are possible (proprietary PGNs with
// different manufacturers).
var PgnInfoLookup map[uint32][]*PgnInfo

// UnseenLookup maps PGN numbers that are defined by the source schema but are
// marked incomplete or missing metadata. These entries are still present in
// PgnInfoLookup when a PGN struct exists for them.
var UnseenLookup map[uint32][]*PgnInfo

// IsProprietaryPGN returns true if the given PGN number falls within one of the four
// NMEA 2000 proprietary PGN ranges. Proprietary PGNs are manufacturer-specific messages
// that require knowing the manufacturer code (embedded in the payload) to decode correctly.
// The four ranges cover all combinations of addressed/broadcast and single-frame/fast-packet.
func IsProprietaryPGN(pgn uint32) bool {
	if pgn >= 0x0EF00 && pgn <= 0x0EFFF {
		// proprietary PDU1 (addressed) single-frame range 0EF00 to 0xEFFF (61184 - 61439) messages.
		// Addressed means that you send it to specific node on the bus. This you can easily use for responding,
		// since you know the sender. For sender it is bit more complicate since your device address may change
		// due to address claiming. There is N2kDeviceList module for handling devices on bus and find them by
		// "NAME" (= 64 bit value set by SetDeviceInformation ).
		return true
	} else if pgn >= 0x0FF00 && pgn <= 0x0FFFF {
		// proprietary PDU2 (non addressed) single-frame range 0xFF00 to 0xFFFF (65280 - 65535).
		// Non addressed means that destination wil be 255 (=broadcast) so any cabable device can handle it.
		return true
	} else if pgn >= 0x1EF00 && pgn <= 0x1EFFF {
		// proprietary PDU1 (addressed) fast-packet PGN range 0x1EF00 to 0x1EFFF (126720 - 126975)
		return true
	} else if pgn >= 0x1FF00 && pgn <= 0x1FFFF {
		// proprietary PDU2 (non addressed) fast packet range 0x1FF00 to 0x1FFFF (130816 - 131071)
		return true
	}

	return false
}

// GetProprietaryInfo extracts the Manufacturer Code and Industry Code from the first
// two bytes of a proprietary PGN's payload. The wire layout for proprietary PGNs is:
//
//	Bits  0-10: Manufacturer Code (11 bits)
//	Bits 11-12: Reserved (2 bits, skipped)
//	Bits 13-15: Industry Code (3 bits)
//
// This function should only be called for PGNs that IsProprietaryPGN reports as true.
// If called on a non-proprietary PGN, the returned values will be meaningless since
// those bytes have a different field layout.
func GetProprietaryInfo(data []uint8) (ManufacturerCodeConst, IndustryCodeConst, error) {
	stream := NewPgnDataStream(data)
	var man ManufacturerCodeConst
	var ind IndustryCodeConst
	var err error
	// Read the 11-bit manufacturer code from bits 0-10.
	if v, err := stream.readLookupField(11); err == nil {
		man = ManufacturerCodeConst(v)
	}
	// Skip the 2 reserved bits (bits 11-12).
	stream.skipBits(2)
	// Read the 3-bit industry code from bits 13-15.
	if v, err := stream.readLookupField(3); err == nil {
		ind = IndustryCodeConst(v)
	}
	return man, ind, err
}

// GetFieldDescriptor looks up the FieldDescriptor for a specific field within a PGN.
// Parameters:
//   - pgn: the PGN number to look up
//   - manID: the manufacturer code, used to disambiguate proprietary PGN variants
//     (pass 0 for standard/non-proprietary PGNs)
//   - fieldIndex: the 1-based field index matching the source field order
//
// For non-proprietary PGNs, the first (and usually only) variant is used.
// For proprietary PGNs, the variant matching manID is selected. If manID is 0 and
// multiple variants exist, an error is returned because the correct variant cannot
// be determined without a manufacturer code.
//
// Returns an error if the PGN is unknown, the field index doesn't exist, or the
// variant cannot be disambiguated.
func GetFieldDescriptor(pgn uint32, manID ManufacturerCodeConst, fieldIndex uint8) (*FieldDescriptor, error) {
	var retval *FieldDescriptor
	var err error

	if pi, piKnown := PgnInfoLookup[pgn]; piKnown {
		if !IsProprietaryPGN(pgn) {
			// Standard PGN: use the first variant. In the future, other match fields
			// (beyond manufacturer code) may need validation here.
			retval = pi[0].Fields[int(fieldIndex)]
		} else {
			if manID != 0 {
				// Proprietary PGN with a known manufacturer: find the matching variant.
				for _, p := range pi {
					if p.ManId == manID {
						retval = p.Fields[int(fieldIndex)]
					}
				}
			} else {
				// Proprietary PGN but no manufacturer code provided.
				if len(pi) == 1 {
					// Only one variant known, so use it optimistically.
					retval = pi[0].Fields[int(fieldIndex)]
				} else {
					// Multiple variants exist and we can't tell which one to use.
					err = fmt.Errorf("error: cannot distinguish between variants for pgn: %d", pgn)
					return nil, err
				}
			}

		}
		if retval == nil {
			err = fmt.Errorf("error: Field Index: %d, not found for pgn: %d with manufacturer code: %d", fieldIndex, pgn, manID)
		}
		return retval, err
	}
	return nil, fmt.Errorf("PGN not found")
}

// SearchUnseenList returns true if the given PGN number has incomplete or missing
// upstream metadata.
func SearchUnseenList(pgn uint32) bool {
	return UnseenLookup[pgn] != nil
}
