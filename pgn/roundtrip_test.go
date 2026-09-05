package pgn

import (
	"fmt"
	"math"
	"reflect"
	"sort"
	"strings"
	"testing"
)

// This file walks every PGN struct registered in structTypeRegistry (built by
// pgngen from the same metadata decodeFields/encodeFields compile against)
// and exercises the metadata-driven codec both ways:
//
//   - TestAllPGNsNullRoundTrip proves the byte fixpoint contract: encoding a
//     zero-value struct, decoding the result, and re-encoding must reproduce
//     the exact same bytes (not struct equality -- lookup fields decode the
//     null sentinel to a non-nil raw value by design; see codec.go).
//   - TestAllPGNsValueRoundTrip proves value fidelity: every field gets a
//     representative in-range value, and after encode -> decode the payload
//     the decoded struct's fields must equal what encode/decode's own
//     derivation rules produce.
//
// Both tests derive their per-field behavior from the same compiled
// *codecPlan that decodeFields/encodeFields use (codecPlanForMessage), so the
// harness never has to re-discover field kinds, bit widths, or repeating-set
// shapes independently -- it drives the exact same interpreter program the
// production code path runs.

// valueRoundTripGroupElements is how many elements TestAllPGNsValueRoundTrip
// populates into every repeating field set. Two elements is the minimum that
// distinguishes "decoded a slice" from "decoded a single leftover element"
// and exercises the per-element encode/decode loop at least once.
const valueRoundTripGroupElements = 2

// TestAllPGNsNullRoundTrip proves the byte fixpoint contract for every
// registered PGN struct: encode(zero) -> decode -> encode must reproduce the
// exact same bytes.
func TestAllPGNsNullRoundTrip(t *testing.T) {
	for _, name := range sortedRegistryNames() {
		factory := structTypeRegistry[name]
		t.Run(name, func(t *testing.T) {
			zero := factory()
			first, err := zero.EncodePayload()
			if err != nil {
				t.Fatalf("encode zero-value struct: %v", err)
			}

			fresh := factory()
			if err := fresh.DecodePayload(first); err != nil {
				t.Fatalf("decode: %v", err)
			}

			second, err := fresh.EncodePayload()
			if err != nil {
				t.Fatalf("re-encode: %v", err)
			}

			if !reflect.DeepEqual(first, second) {
				t.Fatalf("byte fixpoint mismatch:\n first encode: % x\nsecond encode: % x", first, second)
			}
		})
	}
}

// TestAllPGNsValueRoundTrip proves value fidelity for every registered PGN
// struct: populate every field with a
// representative in-range value, encode, decode into a fresh instance, and
// compare against the values encode/decode's own derivation rules (repeating
// group counts and length-reference fields are overridden from the data,
// exactly as encodeFields does) predict for the result.
func TestAllPGNsValueRoundTrip(t *testing.T) {
	for _, name := range sortedRegistryNames() {
		factory := structTypeRegistry[name]
		t.Run(name, func(t *testing.T) {
			input := factory()
			populateMessage(t, input, false)
			encoded, err := input.EncodePayload()
			if err != nil {
				t.Fatalf("encode: %v", err)
			}

			decoded := factory()
			if err := decoded.DecodePayload(encoded); err != nil {
				t.Fatalf("decode: %v", err)
			}

			expected := factory()
			populateMessage(t, expected, true)

			compareMessages(t, name, decoded, expected)
		})
	}
}

// sortedRegistryNames returns structTypeRegistry's keys in sorted order for
// deterministic subtest iteration and output.
func sortedRegistryNames() []string {
	names := make([]string, 0, len(structTypeRegistry))
	for name := range structTypeRegistry {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// populateMessage walks m's compiled codec plan and sets every value-bearing
// field to a representative value.
//
// wantDerived selects which of two populations to build:
//   - false ("input"): the value encodeFields is handed. Match fields are
//     left nil (proving the encoder fills in the Match value itself, exactly
//     like TestDecodeFieldsMatchVariantSelection / the existing codec_test.go
//     coverage), and repeating-group counts / length-reference fields get
//     the same naive representative value as any other number (proving
//     encodeFields overrides a stale user-set value, exactly like
//     TestEncodeFieldsDerivesCountFromSliceLength /
//     TestEncodeFieldsDerivesDynamicLengthFromValue).
//   - true ("expected"): the value decodeFields is contractually required to
//     produce after a round trip. Match fields hold their raw Match value
//     (decodeField always sets them, non-nil, regardless of what was
//     encoded -- see codec.go's fieldKindMatch case), and repeating-group
//     counts / length-reference fields hold the value encodeFields actually
//     derives from the populated data (len(slice) for counts, len(data)
//     bits/bytes for length references) rather than the naive value, because
//     that derivation happens unconditionally on encode and overrides
//     whatever the input struct held.
func populateMessage(t *testing.T, m PGN, wantDerived bool) {
	t.Helper()
	plan, err := codecPlanForMessage(m)
	if err != nil {
		t.Fatalf("compile codec plan: %v", err)
	}

	var overrides map[int]uint64
	if wantDerived {
		overrides = computeExpectedOverrides(plan)
	}

	target := reflect.ValueOf(m).Elem()
	for _, step := range plan.steps {
		if step.group != nil {
			populateGroup(step.group, target, wantDerived, overrides)
			continue
		}
		populateField(step.field, target, wantDerived, overrides)
	}
	if m.PGNNumber() == 126208 && target.FieldByName("Pgn").IsValid() {
		// The dynamic value width comes from a real commanded parameter, not
		// an arbitrary numeric placeholder: Heading in Vessel Heading is 16 bits.
		commanded := uint64(127250)
		target.FieldByName("Pgn").Set(reflect.ValueOf(&commanded))
		for _, step := range plan.steps {
			if step.field != nil {
				if step.field.condition != conditionAlways && step.field.fieldIndex >= 0 {
					target.Field(step.field.fieldIndex).SetZero()
				}
				continue
			}
			slice := target.Field(step.group.sliceIndex)
			for i := 0; i < slice.Len(); i++ {
				for _, field := range step.group.fields {
					if !field.commandedParameter {
						continue
					}
					parameter := uint64(2)
					slice.Index(i).Field(field.fieldIndex - 1).Set(reflect.ValueOf(&parameter))
					slice.Index(i).Field(field.fieldIndex).SetBytes([]byte{0x5c, 0x3d})
				}
			}
		}
	}
}

// populateGroup fills a repeating field set's slice field with
// valueRoundTripGroupElements elements, each populated the same way as a
// fixed field.
func populateGroup(group *planGroup, target reflect.Value, wantDerived bool, overrides map[int]uint64) {
	slice := reflect.MakeSlice(target.Field(group.sliceIndex).Type(), valueRoundTripGroupElements, valueRoundTripGroupElements)
	for i := 0; i < valueRoundTripGroupElements; i++ {
		element := slice.Index(i)
		for f := range group.fields {
			populateField(&group.fields[f], element, wantDerived, overrides)
		}
	}
	target.Field(group.sliceIndex).Set(slice)
}

// populateField sets one field's representative value. See populateMessage
// for the wantDerived/overrides contract.
func populateField(field *planField, target reflect.Value, wantDerived bool, overrides map[int]uint64) {
	switch field.kind {
	case fieldKindReserved, fieldKindSpare:
		// No struct field: reserved/spare bits carry no value.
	case fieldKindFixedString:
		target.Field(field.fieldIndex).SetString(representativeFixedString(field.bitLength))
	case fieldKindLengthString, fieldKindControlString:
		// STRING_LZ / STRING_LAU are length-prefixed on the wire, so unlike
		// STRING_FIX there is no padding to worry about: "AB" round-trips
		// exactly.
		target.Field(field.fieldIndex).SetString("AB")
	case fieldKindFloat:
		value := float32(1.5)
		target.Field(field.fieldIndex).Set(reflect.ValueOf(&value))
	case fieldKindBinary:
		target.Field(field.fieldIndex).SetBytes(representativeBinary(field))
	case fieldKindMatch:
		if !wantDerived || field.match == nil {
			// Input population: leave nil so encodeFields must write the
			// Match value itself.
			return
		}
		setNumericPointer(target.Field(field.fieldIndex), field.signed, uint64(*field.match))
	case fieldKindLookup, fieldKindRawNumber, fieldKindNullableNumber:
		if wantDerived {
			if raw, ok := overrides[field.order]; ok {
				setNumericPointer(target.Field(field.fieldIndex), field.signed, raw)
				return
			}
		}
		setNumericPointer(target.Field(field.fieldIndex), field.signed, representativeNumericRaw(field))
	}
}

// representativeNumericRaw returns a representable ordinary value inside
// the source-schema physical range. Prefer raw tick 1 for compact payloads;
// fields such as AIS MMSI identifiers have a much larger minimum and use the
// first raw tick at or above that minimum instead.
func representativeNumericRaw(field *planField) uint64 {
	physical := field.resolution + field.offset
	if (field.rangeMin == nil || physical >= *field.rangeMin) &&
		(field.rangeMax == nil || physical <= *field.rangeMax) {
		return 1
	}

	raw := float64(0)
	if field.rangeMin != nil {
		raw = math.Ceil((*field.rangeMin-field.offset)/field.resolution - 1e-9)
	} else if field.rangeMax != nil {
		raw = math.Floor((*field.rangeMax-field.offset)/field.resolution + 1e-9)
	}
	if field.signed {
		return uint64(int64(raw))
	}
	if raw < 0 {
		return 0
	}
	return uint64(raw)
}

// representativeFixedString returns an exact-width fill for a STRING_FIX
// field: a string of precisely bitLength/8 bytes, containing none of the
// characters writeFixedString/readFixedString treat as padding (NUL, 0xFF,
// '@'). Filling to the exact width sidesteps the padding round-trip
// question entirely -- there's no truncation or pad-stripping to account for
// when the value already fills the field.
func representativeFixedString(bitLength uint16) string {
	numBytes := int(bitLength / 8)
	if numBytes <= 0 {
		return ""
	}
	const pattern = "ABCDEFGHIJKLMNOP"
	var b strings.Builder
	for b.Len() < numBytes {
		b.WriteString(pattern)
	}
	return b.String()[:numBytes]
}

// representativeBinary returns a representative value for a BINARY/VARIABLE
// field.
//
//   - Variable-width fields (a BitLengthField/DYNAMIC_FIELD_LENGTH reference,
//     or no fixed width at all) get the two-byte value {0xDE, 0xAD} named in
//     the task brief.
//   - Fixed-width fields get exactly ceil(bitLength/8) bytes with the unused
//     high bits of a partial final byte zeroed. This mirrors the bit-exact
//     contract writeBinaryData/readBinaryData implement (see
//     TestDecodeFieldsNonByteAlignedBinaryField in codec_test.go): a
//     narrower-than-byte width only owns its low bits, and writeBinaryData
//     masks off the rest on encode, so if this function left stray high bits
//     set, the encoded-then-decoded value would come back with those bits
//     cleared and mismatch the input value it was compared against.
func representativeBinary(field *planField) []uint8 {
	if field.refOrder != 0 || field.dynLenOrder != 0 || field.bitLength == 0 {
		return []uint8{0xDE, 0xAD}
	}
	numBytes := int((field.bitLength + 7) / 8)
	data := make([]uint8, numBytes)
	for i := range data {
		data[i] = 0xA5
	}
	if rem := field.bitLength % 8; rem != 0 {
		data[numBytes-1] &= uint8(0xFF >> (8 - rem))
	}
	return data
}

// setNumericPointer stores raw into a *uint64 or *int64 struct field
// depending on signed.
func setNumericPointer(field reflect.Value, signed bool, raw uint64) {
	if signed {
		value := int64(raw)
		field.Set(reflect.ValueOf(&value))
		return
	}
	value := raw
	field.Set(reflect.ValueOf(&value))
}

// computeExpectedOverrides predicts, for every repeating-group count field
// and every binary length-reference field in plan, the value encodeFields
// will actually write -- mirroring encodeFields' own override computation
// (addLengthOverride, encodeFields' overrides map in codec.go) but computed
// structurally instead of from a populated reflect.Value, since the
// representative binary payload's length is already known ahead of
// population. Every entry keyed by field order here takes precedence over
// the naive representative value when building the "expected" struct.
func computeExpectedOverrides(plan *codecPlan) map[int]uint64 {
	overrides := make(map[int]uint64)
	for _, step := range plan.steps {
		if step.group != nil {
			if step.group.countOrder != 0 {
				overrides[step.group.countOrder] = valueRoundTripGroupElements
			}
			for f := range step.group.fields {
				addExpectedLengthOverride(&step.group.fields[f], overrides)
			}
			continue
		}
		addExpectedLengthOverride(step.field, overrides)
	}
	return overrides
}

// addExpectedLengthOverride is the structural counterpart to codec.go's
// addLengthOverride: bit count for BitLengthField references, byte count for
// DYNAMIC_FIELD_LENGTH references.
func addExpectedLengthOverride(field *planField, overrides map[int]uint64) {
	if field.kind != fieldKindBinary {
		return
	}
	length := uint64(len(representativeBinary(field)))
	if field.refOrder != 0 {
		overrides[field.refOrder] = length * 8
	} else if field.dynLenOrder != 0 {
		overrides[field.dynLenOrder] = length
	}
}

// compareMessages compares decoded against expected field by field (fixed
// fields and every repeating-group element), skipping the Info field and
// reserved/spare fields (which have no struct field to compare). Mismatches
// are reported individually via t.Errorf so a single run surfaces every
// diverging field, not just the first.
func compareMessages(t *testing.T, structName string, decoded, expected PGN) {
	t.Helper()
	plan, err := codecPlanForMessage(decoded)
	if err != nil {
		t.Fatalf("compile codec plan: %v", err)
	}

	d := reflect.ValueOf(decoded).Elem()
	e := reflect.ValueOf(expected).Elem()
	for _, step := range plan.steps {
		if step.group != nil {
			compareGroup(t, structName, step.group, d, e)
			continue
		}
		compareField(t, structName, "", step.field, d, e)
	}
}

// compareField compares one field's value on d against e, reporting a
// mismatch with the struct name, field path, and both values.
func compareField(t *testing.T, structName, path string, field *planField, d, e reflect.Value) {
	if field.fieldIndex < 0 {
		return
	}
	dv := d.Field(field.fieldIndex).Interface()
	ev := e.Field(field.fieldIndex).Interface()
	if reflect.DeepEqual(dv, ev) {
		return
	}
	name := d.Type().Field(field.fieldIndex).Name
	t.Errorf("%s%s.%s: got %s, want %s", structName, path, name, formatFieldValue(dv), formatFieldValue(ev))
}

// compareGroup compares every element of a repeating field set.
func compareGroup(t *testing.T, structName string, group *planGroup, d, e reflect.Value) {
	sliceName := d.Type().Field(group.sliceIndex).Name
	dSlice := d.Field(group.sliceIndex)
	eSlice := e.Field(group.sliceIndex)
	if dSlice.Len() != eSlice.Len() {
		t.Errorf("%s.%s: got %d elements, want %d", structName, sliceName, dSlice.Len(), eSlice.Len())
		return
	}
	for i := 0; i < dSlice.Len(); i++ {
		path := fmt.Sprintf(".%s[%d]", sliceName, i)
		de := dSlice.Index(i)
		ee := eSlice.Index(i)
		for f := range group.fields {
			compareField(t, structName, path, &group.fields[f], de, ee)
		}
	}
}

// formatFieldValue renders a field's value for a mismatch message, printing
// pointer targets (the common case: every nullable numeric field is
// *uint64/*int64/*float32) rather than the pointer address %v would
// otherwise show, and byte slices in hex.
func formatFieldValue(v any) string {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Pointer:
		if rv.IsNil() {
			return "<nil>"
		}
		return fmt.Sprintf("%v", rv.Elem().Interface())
	case reflect.Slice:
		if rv.Type().Elem().Kind() == reflect.Uint8 {
			return fmt.Sprintf("% x", rv.Interface())
		}
	}
	return fmt.Sprintf("%v", v)
}
