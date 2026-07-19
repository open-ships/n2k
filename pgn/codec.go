package pgn

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"reflect"
	"sort"
	"strconv"
	"sync"
)

// This file implements the metadata-driven PGN codec. Generated PGN structs
// delegate their DecodePayload/EncodePayload methods to decodeFields and
// encodeFields, which interpret the field metadata registered for the struct
// type (structInfoLookup) instead of relying on generated per-field code.
//
// Field semantics follow the upstream source schema:
//   - Numeric fields of 4+ bits treat the maximum representable value as the
//     null sentinel and decode it to nil. Lookup-family fields and numbers
//     narrower than 4 bits always keep their raw value.
//   - Fields carrying a Match value select PGN variants: decoding validates
//     every Match field at its absolute bit offset before any sequential
//     read, and a mismatch returns an error so dispatch can try the next
//     candidate struct.
//   - Repeating field sets decode into Repeating1/Repeating2 slice fields.
//     The set's count field (always a fixed field preceding the set) supplies
//     the iteration count; a null or absent count means read until the
//     payload is exhausted. Encoding derives the count field from the slice
//     length, overriding any user-set value.
//   - Variable-length binary fields resolve their width at runtime. A
//     BitLengthField reference names an earlier field whose decoded value is
//     the width in BITS (the upstream fields read "Number of Bits in Binary
//     Data Field"). A DYNAMIC_FIELD_VALUE without such a reference takes its
//     width in BYTES from the nearest preceding DYNAMIC_FIELD_LENGTH field.
//     Encoding writes the data verbatim and derives the referenced length
//     field from len(data) in the same unit decoding expects.

// fieldKind selects the decode/encode strategy for one field descriptor.
type fieldKind uint8

const (
	// fieldKindReserved skips bits on decode and writes all-ones on encode.
	fieldKindReserved fieldKind = iota
	// fieldKindSpare skips bits on decode and writes zeros on encode.
	fieldKindSpare
	// fieldKindFixedString is a fixed-width padded string (STRING_FIX).
	fieldKindFixedString
	// fieldKindLengthString is a length-prefixed string (STRING_LZ).
	fieldKindLengthString
	// fieldKindControlString is a length+control prefixed string (STRING_LAU).
	fieldKindControlString
	// fieldKindBinary is an opaque byte field (BINARY, VARIABLE,
	// DYNAMIC_FIELD_VALUE), fixed or variable width.
	fieldKindBinary
	// fieldKindFloat is an IEEE 754 single-precision float (FLOAT).
	fieldKindFloat
	// fieldKindMatch is a numeric field with a variant-selecting Match value.
	fieldKindMatch
	// fieldKindLookup is an enumeration field: raw values survive decoding
	// because every bit pattern has a defined meaning.
	fieldKindLookup
	// fieldKindNullableNumber is a plain number of 4+ bits with null
	// sentinel detection.
	fieldKindNullableNumber
	// fieldKindRawNumber is a plain number narrower than 4 bits, too small
	// to reserve a null sentinel.
	fieldKindRawNumber
)

// fieldCondition selects whether a field is present in the wire payload.
// Conditions are compiled from the source metadata so an unknown predicate
// cannot be silently ignored.
type fieldCondition uint8

const (
	conditionAlways fieldCondition = iota
	conditionPGNIsProprietary
)

// planField is one compiled field of a codecPlan.
type planField struct {
	// order is the 1-based metadata field order.
	order int
	// name is the human-readable metadata field name used in validation errors.
	name string
	// kind selects the decode/encode strategy.
	kind fieldKind
	// fieldIndex is the reflect field index in the containing struct
	// (parent struct for fixed fields, element struct for group fields).
	// It is -1 for reserved and spare fields, which have no struct field.
	fieldIndex int
	// bitLength is the field width in bits. For binary fields it is the
	// raw (not byte-aligned) fixed width -- upstream packs binary fields
	// bit-exact, so a narrower-than-byte width must not consume extra
	// padding bits -- or 0 when the width is variable.
	bitLength uint16
	// signed marks two's-complement numeric fields.
	signed bool
	// condition controls whether this field occupies space in the payload.
	condition fieldCondition
	// resolution and offset convert raw numeric ticks to physical values for
	// source-schema range validation.
	resolution float64
	offset     float64
	rangeMin   *float64
	rangeMax   *float64
	// sentinels are explicit special raw values permitted outside the normal
	// physical range (out-of-range, reserved, and unknown states).
	sentinels map[int64]struct{}
	// match is the variant-selecting value for fieldKindMatch fields.
	match *int
	// refOrder is the order of the field holding this binary field's width
	// in bits (metadata BitLengthField), or 0.
	refOrder int
	// dynLenOrder is the order of the preceding DYNAMIC_FIELD_LENGTH field
	// holding this binary field's width in bytes, or 0.
	dynLenOrder int
}

// planGroup is a compiled repeating field set.
type planGroup struct {
	// sliceIndex is the reflect field index of the RepeatingN slice in the
	// parent struct.
	sliceIndex int
	// elemType is the slice element struct type.
	elemType reflect.Type
	// countOrder is the order of the fixed field that carries the group
	// count, or 0 when the group repeats until the payload is exhausted.
	countOrder int
	// fields are the group's member fields in order, with fieldIndex values
	// resolved against elemType.
	fields []planField
}

// planStep is one decode/encode step: exactly one of field or group is set.
type planStep struct {
	field *planField
	group *planGroup
}

// codecPlan is the compiled decode/encode program for one PGN struct type.
type codecPlan struct {
	info *PgnInfo
	// matches lists every field descriptor carrying a Match value, in field
	// order, for absolute-offset validation before sequential decoding.
	matches []*FieldDescriptor
	steps   []planStep
}

// codecPlanCache maps reflect struct types to compiled plans. Plans are
// immutable once built, so a lock-free sync.Map suffices.
var codecPlanCache sync.Map // reflect.Type -> *codecPlan

// decodeFields decodes payload into m using the field metadata registered for
// m's struct type. It performs Match-field validation (returning an error so
// dispatch can try the next candidate), reads numeric fields with null-sentinel
// detection, expands repeating field sets into slice fields, and resolves
// variable-length binary fields via their length-field references.
func decodeFields(m PGN, payload []uint8) (err error) {
	// Stash the wire payload so a later encodeFields call on this same
	// struct can reproduce it byte-for-bit in RESERVED/STRING_FIX
	// positions instead of always applying the default fill. This must
	// happen before any error return below: SetMessageInfo also carries
	// this decode's info fields (PGN, source, etc.), which dispatch's
	// decodePGNCandidates already set on the candidate before calling in.
	info := m.MessageInfo()
	info.rawPayload = append([]uint8(nil), payload...)
	info.rawCanonical = nil
	m.SetMessageInfo(info)
	defer func() {
		if err != nil {
			return
		}
		canonical, encodeErr := encodeCurrentFields(m)
		if encodeErr != nil {
			return
		}
		decodedInfo := m.MessageInfo()
		decodedInfo.rawCanonical = append([]uint8(nil), canonical...)
		m.SetMessageInfo(decodedInfo)
	}()

	plan, err := codecPlanForMessage(m)
	if err != nil {
		return err
	}
	for _, desc := range plan.matches {
		if desc.BitLength == 0 {
			return fmt.Errorf("match failed for %s", plan.info.Description)
		}
		raw, err := readRawAt(payload, desc.BitOffset, desc.BitLength)
		if err != nil {
			return err
		}
		if int(raw) != *desc.Match {
			return fmt.Errorf("match failed for %s", plan.info.Description)
		}
	}

	target := reflect.ValueOf(m).Elem()
	stream := NewPgnDataStream(payload)
	// raws records the raw value of every non-null unsigned numeric field by
	// order, so later fields can resolve count and length references.
	raws := make(map[int]uint64)
	for _, step := range plan.steps {
		if step.group != nil {
			done, err := decodeGroup(stream, step.group, target, raws)
			if err != nil {
				return err
			}
			if done {
				return nil
			}
			continue
		}
		included, err := conditionSatisfied(step.field.condition, raws)
		if err != nil {
			return fmt.Errorf("field %d (%s): %w", step.field.order, step.field.name, err)
		}
		if !included {
			continue
		}
		if err := decodeField(stream, step.field, target, raws); err != nil {
			return err
		}
		if stream.isEOF() {
			// Payloads may legitimately end early; trailing fields stay nil.
			return nil
		}
	}
	return nil
}

// encodeFields is the inverse of decodeFields: it writes every field in
// metadata order, writing null sentinels for nil numeric fields, Match values
// for nil match fields, deriving repeating-set count fields from slice
// lengths, and deriving variable-length binary length fields from the data
// slice length.
//
// When m was produced by decodeFields, untouched fields return the original
// wire payload exactly, including reserved bits and trailing filler. If any
// decoded field changes, the current field values are encoded instead.
func encodeFields(m PGN) ([]uint8, error) {
	info := m.MessageInfo()
	current, err := encodeCurrentFields(m)
	if err != nil {
		return nil, err
	}
	if len(info.rawPayload) > 0 && len(info.rawCanonical) > 0 && bytes.Equal(current, info.rawCanonical) {
		return append([]uint8(nil), info.rawPayload...), nil
	}
	return current, nil
}

func encodeCurrentFields(m PGN) ([]uint8, error) {
	plan, err := codecPlanForMessage(m)
	if err != nil {
		return nil, err
	}
	source := reflect.ValueOf(m).Elem()
	writer := NewPGNDataStreamWriter()

	// Derived fields (group counts, binary length references) override any
	// user-set value; collect them before writing so a length field that
	// precedes its data field encodes the derived value.
	overrides := make(map[int]uint64)
	for _, step := range plan.steps {
		if step.group != nil {
			if step.group.countOrder != 0 {
				overrides[step.group.countOrder] = uint64(source.Field(step.group.sliceIndex).Len())
			}
			continue
		}
		addLengthOverride(overrides, step.field, source)
	}
	raws := collectEncodeRaws(plan, source, overrides)

	for _, step := range plan.steps {
		if step.group != nil {
			encodeGroup(writer, step.group, source, raws)
			continue
		}
		included, conditionErr := conditionSatisfied(step.field.condition, raws)
		if conditionErr != nil {
			writer.setErr(fmt.Errorf("field %d (%s): %w", step.field.order, step.field.name, conditionErr))
			continue
		}
		if !included {
			continue
		}
		encodeField(writer, step.field, source, overrides)
	}
	return writer.Bytes(), writer.Err()
}

// codecPlanForMessage resolves the compiled plan for a PGN message value.
func codecPlanForMessage(m PGN) (*codecPlan, error) {
	t := reflect.TypeOf(m)
	if t == nil || t.Kind() != reflect.Pointer || t.Elem().Kind() != reflect.Struct {
		return nil, fmt.Errorf("PGN message must be a struct pointer, got %T", m)
	}
	return codecPlanFor(t.Elem())
}

// codecPlanFor returns the cached plan for a PGN struct type, compiling it on
// first use.
func codecPlanFor(t reflect.Type) (*codecPlan, error) {
	if cached, ok := codecPlanCache.Load(t); ok {
		return cached.(*codecPlan), nil
	}
	plan, err := compileCodecPlan(t)
	if err != nil {
		return nil, err
	}
	codecPlanCache.Store(t, plan)
	return plan, nil
}

// compileCodecPlan builds the decode/encode program for a PGN struct type
// from its registered metadata and its `n2k` struct tags.
func compileCodecPlan(t reflect.Type) (*codecPlan, error) {
	info := structInfoLookup[t.Name()]
	if info == nil {
		return nil, fmt.Errorf("no PGN metadata registered for struct %s", t.Name())
	}
	indexByOrder, sliceIndexes, err := codecStructIndexes(t)
	if err != nil {
		return nil, fmt.Errorf("PGN struct %s: %w", t.Name(), err)
	}

	groups, groupBounds, err := compileCodecGroups(t, info, sliceIndexes)
	if err != nil {
		return nil, err
	}

	plan := &codecPlan{info: info}
	orders := sortedFieldOrders(info.Fields)
	for _, order := range orders {
		if info.Fields[order].Match != nil {
			plan.matches = append(plan.matches, info.Fields[order])
		}
	}

	lastDynLen := 0
	for i := 0; i < len(orders); {
		order := orders[i]
		grouped := false
		for g, group := range groups {
			if group == nil || order != groupBounds[g][0] {
				continue
			}
			plan.steps = append(plan.steps, planStep{group: group})
			for i < len(orders) && orders[i] < groupBounds[g][1] {
				i++
			}
			grouped = true
			break
		}
		if grouped {
			continue
		}
		field, err := compileCodecField(info, order, indexByOrder, &lastDynLen)
		if err != nil {
			return nil, fmt.Errorf("PGN struct %s: %w", t.Name(), err)
		}
		plan.steps = append(plan.steps, planStep{field: field})
		i++
	}
	return plan, nil
}

// compileCodecGroups compiles the (up to two) repeating field sets described
// by info against the struct's RepeatingN slice fields. The returned bounds
// hold each compiled group's [start, end) order range. A metadata group whose
// members are all reserved/spare needs no slice field and stays uncompiled,
// leaving its members to decode once as fixed fields.
func compileCodecGroups(t reflect.Type, info *PgnInfo, sliceIndexes [2]int) ([2]*planGroup, [2][2]int, error) {
	var groups [2]*planGroup
	var bounds [2][2]int
	specs := [2]struct{ start, count, size *int }{
		{info.RepeatingFieldSet1StartField, info.RepeatingFieldSet1CountField, info.RepeatingFieldSet1Size},
		{info.RepeatingFieldSet2StartField, info.RepeatingFieldSet2CountField, info.RepeatingFieldSet2Size},
	}
	for g, spec := range specs {
		if spec.start == nil || spec.size == nil || *spec.size <= 0 {
			continue
		}
		start, end := *spec.start, *spec.start+*spec.size
		if sliceIndexes[g] < 0 {
			for order := start; order < end; order++ {
				if desc := info.Fields[order]; desc != nil && desc.GolangType != "" {
					return groups, bounds, fmt.Errorf(
						"PGN struct %s has no repeating slice for field set %d", t.Name(), g+1)
				}
			}
			continue
		}
		sliceType := t.Field(sliceIndexes[g]).Type
		if sliceType.Kind() != reflect.Slice || sliceType.Elem().Kind() != reflect.Struct {
			return groups, bounds, fmt.Errorf(
				"PGN struct %s repeating field set %d must be a slice of structs", t.Name(), g+1)
		}
		elemType := sliceType.Elem()
		elemIndexByOrder, _, err := codecStructIndexes(elemType)
		if err != nil {
			return groups, bounds, fmt.Errorf("PGN struct %s: %w", elemType.Name(), err)
		}
		group := &planGroup{
			sliceIndex: sliceIndexes[g],
			elemType:   elemType,
		}
		if spec.count != nil {
			group.countOrder = *spec.count
		}
		lastDynLen := 0
		for order := start; order < end; order++ {
			if info.Fields[order] == nil {
				return groups, bounds, fmt.Errorf(
					"PGN struct %s repeating field set %d references missing order %d", t.Name(), g+1, order)
			}
			field, err := compileCodecField(info, order, elemIndexByOrder, &lastDynLen)
			if err != nil {
				return groups, bounds, fmt.Errorf("PGN struct %s: %w", t.Name(), err)
			}
			group.fields = append(group.fields, *field)
		}
		groups[g] = group
		bounds[g] = [2]int{start, end}
	}
	return groups, bounds, nil
}

// compileCodecField classifies one metadata field into a planField and
// resolves its struct field index. lastDynLen threads the order of the most
// recent DYNAMIC_FIELD_LENGTH field within the current section (fixed fields
// or one repeating group) so a following DYNAMIC_FIELD_VALUE can reference it.
func compileCodecField(info *PgnInfo, order int, indexByOrder map[int]int, lastDynLen *int) (*planField, error) {
	desc := info.Fields[order]
	condition, err := compileFieldCondition(desc.Condition)
	if err != nil {
		return nil, fmt.Errorf("field %d (%s): %w", order, desc.Name, err)
	}
	field := &planField{
		order:      order,
		name:       desc.Name,
		fieldIndex: -1,
		bitLength:  desc.BitLength,
		signed:     desc.Signed,
		condition:  condition,
		resolution: float64(desc.Resolution),
		rangeMin:   desc.RangeMin,
		rangeMax:   desc.RangeMax,
		sentinels:  make(map[int64]struct{}, 3),
	}
	if field.resolution == 0 {
		field.resolution = 1
	}
	if desc.Offset != nil {
		field.offset = *desc.Offset
	}
	for _, sentinel := range []*int64{desc.OutOfRangeValue, desc.ReservedValue, desc.UnknownValue} {
		if sentinel != nil {
			field.sentinels[*sentinel] = struct{}{}
		}
	}
	switch desc.SourceType {
	case "RESERVED":
		field.kind = fieldKindReserved
	case "SPARE":
		field.kind = fieldKindSpare
	case "STRING_FIX":
		field.kind = fieldKindFixedString
	case "STRING_LZ":
		field.kind = fieldKindLengthString
	case "STRING_LAU":
		field.kind = fieldKindControlString
	case "BINARY", "VARIABLE", "DYNAMIC_FIELD_VALUE":
		field.kind = fieldKindBinary
		if desc.BitLengthField != nil {
			field.refOrder = *desc.BitLengthField
		} else if desc.SourceType == "DYNAMIC_FIELD_VALUE" && *lastDynLen != 0 {
			field.dynLenOrder = *lastDynLen
			*lastDynLen = 0
		}
	case "FLOAT":
		field.kind = fieldKindFloat
	default:
		switch {
		case desc.Match != nil:
			field.kind = fieldKindMatch
			field.match = desc.Match
		case desc.LookupEnumeration != "" || desc.LookupBitEnumeration != "" ||
			desc.LookupIndirectEnumeration != "" || desc.LookupFieldTypeEnumeration != "":
			field.kind = fieldKindLookup
		case desc.BitLength >= 4:
			field.kind = fieldKindNullableNumber
		default:
			field.kind = fieldKindRawNumber
		}
	}
	if desc.SourceType == "DYNAMIC_FIELD_LENGTH" {
		// Navico proprietary PGNs 130822 and 130823 (the record-list
		// variants, e.g. NavicoUdbDatabaseObjectDump and
		// NavicoDataTypeSourceDirectory) declare a record length that
		// includes a 3-byte sub-record header (a class/type byte plus a
		// 16-bit data-type id) which the schema does not expose as its own
		// field, so the declared length over-reads the following
		// DYNAMIC_FIELD_VALUE by about 3 bytes relative to the actual
		// value data. This is tolerated rather than special-cased: those
		// records live in an until-EOF repeating group (no count field),
		// so decodeGroup's blanket "any in-group error ends the group"
		// handling discards whatever partial/misaligned tail results and
		// decode still succeeds; encode independently derives the length
		// field from len(data), so it stays internally symmetric even
		// though it does not reproduce upstream's header-inclusive count.
		*lastDynLen = order
	}
	if desc.GolangType != "" {
		index, ok := indexByOrder[order]
		if !ok {
			return nil, fmt.Errorf("no struct field tagged for metadata order %d (%s)", order, desc.Name)
		}
		field.fieldIndex = index
	}
	return field, nil
}

func compileFieldCondition(condition string) (fieldCondition, error) {
	switch condition {
	case "":
		return conditionAlways, nil
	case "PGNIsProprietary":
		return conditionPGNIsProprietary, nil
	default:
		return conditionAlways, fmt.Errorf("unsupported field condition %q", condition)
	}
}

// conditionSatisfied evaluates a compiled predicate using raw field values
// from the same message. PGNIsProprietary is defined by the commanded PGN in
// field order 2 of NMEA group-function messages.
func conditionSatisfied(condition fieldCondition, raws map[int]uint64) (bool, error) {
	switch condition {
	case conditionAlways:
		return true, nil
	case conditionPGNIsProprietary:
		commandedPGN, ok := raws[2]
		if !ok {
			// A zero-value group-function struct has no commanded PGN. Treat the
			// optional proprietary header as absent; callers still receive normal
			// validation if they populate an invalid PGN value.
			return false, nil
		}
		return IsProprietaryPGN(uint32(commandedPGN)), nil
	default:
		return false, fmt.Errorf("unsupported compiled field condition %d", condition)
	}
}

// codecStructIndexes reads the `n2k` struct tags of t, returning the field
// index for each numeric order tag and the indexes of the rep1/rep2 slice
// fields (-1 when absent).
func codecStructIndexes(t reflect.Type) (map[int]int, [2]int, error) {
	indexByOrder := make(map[int]int)
	sliceIndexes := [2]int{-1, -1}
	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag.Get("n2k")
		if tag == "" {
			continue
		}
		switch tag {
		case "rep1":
			sliceIndexes[0] = i
		case "rep2":
			sliceIndexes[1] = i
		default:
			order, err := strconv.Atoi(tag)
			if err != nil {
				return nil, sliceIndexes, fmt.Errorf("field %s has invalid n2k tag %q", t.Field(i).Name, tag)
			}
			indexByOrder[order] = i
		}
	}
	return indexByOrder, sliceIndexes, nil
}

// sortedFieldOrders returns the metadata field orders in ascending order.
func sortedFieldOrders(fields map[int]*FieldDescriptor) []int {
	orders := make([]int, 0, len(fields))
	for order := range fields {
		orders = append(orders, order)
	}
	sort.Ints(orders)
	return orders
}

// decodeField decodes one field from the stream into target's struct field.
func decodeField(stream *PGNDataStream, field *planField, target reflect.Value, raws map[int]uint64) error {
	switch field.kind {
	case fieldKindReserved, fieldKindSpare:
		stream.skipBits(field.bitLength)
		return nil
	case fieldKindFixedString:
		value, err := stream.readFixedString(field.bitLength)
		if err != nil {
			return err
		}
		target.Field(field.fieldIndex).SetString(value)
	case fieldKindLengthString:
		value, err := stream.readStringWithLength()
		if err != nil {
			return err
		}
		target.Field(field.fieldIndex).SetString(value)
	case fieldKindControlString:
		value, err := stream.readStringWithLengthAndControl()
		if err != nil {
			return err
		}
		target.Field(field.fieldIndex).SetString(value)
	case fieldKindFloat:
		value, err := stream.readFloat32()
		if err != nil {
			return err
		}
		setOptional(target.Field(field.fieldIndex), value)
	case fieldKindBinary:
		value, err := stream.readBinaryData(binaryFieldBits(field, raws, stream))
		if err != nil {
			return err
		}
		target.Field(field.fieldIndex).SetBytes(value)
	case fieldKindMatch, fieldKindLookup, fieldKindRawNumber:
		raw, err := stream.getNumberRaw(field.bitLength)
		if err != nil {
			return err
		}
		raws[field.order] = raw
		setRawNumber(target.Field(field.fieldIndex), raw, field)
	case fieldKindNullableNumber:
		if field.signed {
			value, err := stream.getSignedNullableNumber(field.bitLength)
			if err != nil {
				return err
			}
			setOptional(target.Field(field.fieldIndex), value)
			return nil
		}
		value, err := stream.getUnsignedNullableNumber(field.bitLength)
		if err != nil {
			return err
		}
		if value != nil {
			raws[field.order] = *value
		}
		setOptional(target.Field(field.fieldIndex), value)
	}
	return nil
}

// decodeGroup decodes one repeating field set into its slice field. The
// returned done flag reports that the payload is exhausted (or ended inside a
// group, whose partial element is discarded) and decoding should stop.
//
// Every decodeField error encountered while decoding an element -- not just
// the "read past end of payload" case -- is treated as end-of-data: the
// partial element is discarded and decodeGroup reports success. This is
// sanctioned for the genuine EOF case (payloads legitimately end mid-group),
// but it also swallows other decode failures (e.g. a malformed length that
// makes a nested binary field over- or under-read) under the same "stop
// cleanly" path. That conflation is deliberate for now -- there is no
// decode-time signal here that distinguishes "ran out of bytes" from "hit
// something we can't parse" -- pending a finer-grained error distinction.
func decodeGroup(stream *PGNDataStream, group *planGroup, target reflect.Value, raws map[int]uint64) (bool, error) {
	count := -1
	if group.countOrder != 0 {
		if raw, ok := raws[group.countOrder]; ok {
			count = int(raw)
		}
	}
	sliceField := target.Field(group.sliceIndex)
	slice := reflect.MakeSlice(sliceField.Type(), 0, 0)
	for i := 0; count < 0 || i < count; i++ {
		if stream.isEOF() {
			break
		}
		element := reflect.New(group.elemType).Elem()
		complete := true
		for f := range group.fields {
			included, err := conditionSatisfied(group.fields[f].condition, raws)
			if err != nil {
				return false, fmt.Errorf("repeating group field %d (%s): %w", group.fields[f].order, group.fields[f].name, err)
			}
			if !included {
				continue
			}
			if err := decodeField(stream, &group.fields[f], element, raws); err != nil {
				if !errors.Is(err, ErrUnexpectedPayloadEnd) {
					return false, fmt.Errorf("decode repeating group field %d: %w", group.fields[f].order, err)
				}
				// Truncated payloads may legitimately end inside a final
				// repeating element; discard only that partial element.
				complete = false
				break
			}
			if stream.isEOF() && f < len(group.fields)-1 {
				complete = false
				break
			}
		}
		if !complete {
			setGroupSlice(sliceField, slice)
			return true, nil
		}
		slice = reflect.Append(slice, element)
	}
	setGroupSlice(sliceField, slice)
	return stream.isEOF(), nil
}

// binaryFieldBits resolves the width of a binary field at decode time.
// Variable widths are clamped to the remaining payload, mirroring upstream
// decoder behavior for truncated or sentinel-valued length fields.
func binaryFieldBits(field *planField, raws map[int]uint64, stream *PGNDataStream) uint16 {
	if field.refOrder == 0 && field.dynLenOrder == 0 {
		return field.bitLength
	}
	remaining := uint64(0)
	if total := uint32(len(stream.data)) * 8; total > stream.getBitOffset() {
		remaining = uint64(total - stream.getBitOffset())
	}
	bits := remaining
	if field.refOrder != 0 {
		if raw, ok := raws[field.refOrder]; ok {
			// The referenced field carries the width in bits.
			bits = raw
		}
	} else if raw, ok := raws[field.dynLenOrder]; ok {
		// DYNAMIC_FIELD_LENGTH fields carry the width in bytes.
		bits = raw * 8
	}
	if bits > remaining {
		bits = remaining
	}
	return uint16(bits)
}

// setOptional stores a typed pointer (or nil) into a pointer-typed struct field.
func setOptional[T any](field reflect.Value, value *T) {
	if value == nil {
		field.Set(reflect.Zero(field.Type()))
		return
	}
	field.Set(reflect.ValueOf(value))
}

// setRawNumber stores a raw numeric value into a *uint64 or *int64 struct
// field, sign-extending when the metadata marks the field as signed.
func setRawNumber(field reflect.Value, raw uint64, planned *planField) {
	if planned.signed {
		value := signExtend(raw, planned.bitLength)
		field.Set(reflect.ValueOf(&value))
		return
	}
	value := raw
	field.Set(reflect.ValueOf(&value))
}

// setGroupSlice stores a decoded group slice, leaving the field nil when no
// complete element was decoded.
func setGroupSlice(field reflect.Value, slice reflect.Value) {
	if slice.Len() == 0 {
		field.Set(reflect.Zero(field.Type()))
		return
	}
	field.Set(slice)
}

// encodeField writes one field from source's struct field to the writer.
// overrides carries derived values (group counts, binary lengths) keyed by
// field order that take precedence over the struct field value.
func encodeField(writer *PGNDataStreamWriter, field *planField, source reflect.Value, overrides map[int]uint64) {
	switch field.kind {
	case fieldKindReserved:
		writer.writeReservedBits(field.bitLength)
	case fieldKindSpare:
		writer.writeSpareBits(field.bitLength)
	case fieldKindFixedString:
		writer.writeFixedString(source.Field(field.fieldIndex).String(), field.bitLength)
	case fieldKindLengthString:
		writer.writeStringWithLength(source.Field(field.fieldIndex).String())
	case fieldKindControlString:
		writer.writeStringWithLengthAndControl(source.Field(field.fieldIndex).String())
	case fieldKindFloat:
		value, _ := source.Field(field.fieldIndex).Interface().(*float32)
		writer.writeFloat32(value)
	case fieldKindBinary:
		data := source.Field(field.fieldIndex).Bytes()
		bits := field.bitLength
		if field.refOrder != 0 || field.dynLenOrder != 0 || bits == 0 {
			bits = uint16(len(data) * 8)
		}
		writer.writeBinaryData(data, bits)
	case fieldKindMatch, fieldKindLookup, fieldKindNullableNumber, fieldKindRawNumber:
		encodeNumericField(writer, field, source, overrides)
	}
}

// encodeNumericField writes a numeric field, applying derived overrides,
// null sentinels for nil values, and Match values for nil match fields.
func encodeNumericField(writer *PGNDataStreamWriter, field *planField, source reflect.Value, overrides map[int]uint64) {
	if raw, ok := overrides[field.order]; ok {
		if err := field.validateUnsigned(raw); err != nil {
			writer.setErr(err)
			return
		}
		writer.setErr(writer.putNumberRaw(raw, field.bitLength))
		return
	}
	if field.signed {
		value, _ := source.Field(field.fieldIndex).Interface().(*int64)
		if value != nil {
			if err := field.validateSigned(*value); err != nil {
				writer.setErr(err)
				return
			}
			writer.writeInt64(value, field.bitLength)
			return
		}
		if field.match != nil {
			writer.writeLookupField(uint64(*field.match), field.bitLength)
			return
		}
		writer.setErr(writer.putNullSigned(field.bitLength))
		return
	}
	value, _ := source.Field(field.fieldIndex).Interface().(*uint64)
	if value != nil {
		if err := field.validateUnsigned(*value); err != nil {
			writer.setErr(err)
			return
		}
		writer.writeUInt64(value, field.bitLength)
		return
	}
	if field.match != nil {
		writer.writeLookupField(uint64(*field.match), field.bitLength)
		return
	}
	writer.setErr(writer.putNullUnsigned(field.bitLength))
}

// encodeGroup writes every element of a repeating field set. Length
// references inside a group always point at fields of the same group, so
// derived overrides are computed per element.
func encodeGroup(writer *PGNDataStreamWriter, group *planGroup, source reflect.Value, parentRaws map[int]uint64) {
	slice := source.Field(group.sliceIndex)
	for i := 0; i < slice.Len(); i++ {
		element := slice.Index(i)
		overrides := make(map[int]uint64)
		for f := range group.fields {
			addLengthOverride(overrides, &group.fields[f], element)
		}
		raws := collectFieldRaws(group.fields, element, overrides, parentRaws)
		for f := range group.fields {
			included, err := conditionSatisfied(group.fields[f].condition, raws)
			if err != nil {
				writer.setErr(fmt.Errorf("repeating group field %d (%s): %w", group.fields[f].order, group.fields[f].name, err))
				continue
			}
			if !included {
				continue
			}
			encodeField(writer, &group.fields[f], element, overrides)
		}
	}
}

func collectEncodeRaws(plan *codecPlan, source reflect.Value, overrides map[int]uint64) map[int]uint64 {
	raws := make(map[int]uint64, len(overrides)+len(plan.steps))
	for order, raw := range overrides {
		raws[order] = raw
	}
	for _, step := range plan.steps {
		if step.field != nil {
			collectFieldRaw(raws, step.field, source, overrides)
		}
	}
	return raws
}

func collectFieldRaws(fields []planField, source reflect.Value, overrides, parent map[int]uint64) map[int]uint64 {
	raws := make(map[int]uint64, len(parent)+len(overrides)+len(fields))
	for order, raw := range parent {
		raws[order] = raw
	}
	for order, raw := range overrides {
		raws[order] = raw
	}
	for i := range fields {
		collectFieldRaw(raws, &fields[i], source, overrides)
	}
	return raws
}

func collectFieldRaw(raws map[int]uint64, field *planField, source reflect.Value, overrides map[int]uint64) {
	if raw, ok := overrides[field.order]; ok {
		raws[field.order] = raw
		return
	}
	if field.fieldIndex < 0 {
		return
	}
	value := source.Field(field.fieldIndex)
	if value.Kind() != reflect.Pointer || value.IsNil() {
		if field.match != nil {
			raws[field.order] = uint64(*field.match)
		}
		return
	}
	switch raw := value.Elem().Interface().(type) {
	case uint64:
		raws[field.order] = raw
	case int64:
		raws[field.order] = uint64(raw)
	}
}

func (field *planField) validateSigned(value int64) error {
	if field.bitLength == 0 || field.bitLength > 64 {
		return fmt.Errorf("field %d (%s): invalid signed width %d", field.order, field.name, field.bitLength)
	}
	if field.bitLength < 64 {
		minValue := -(int64(1) << (field.bitLength - 1))
		maxValue := (int64(1) << (field.bitLength - 1)) - 1
		if value < minValue || value > maxValue {
			return fmt.Errorf("field %d (%s): value %d exceeds %d-bit signed range [%d,%d]", field.order, field.name, value, field.bitLength, minValue, maxValue)
		}
		_, declaredSentinel := field.sentinels[value]
		if field.kind == fieldKindNullableNumber && value == maxValue && !declaredSentinel {
			return fmt.Errorf("field %d (%s): value %d is the null sentinel; use nil", field.order, field.name, value)
		}
	}
	return field.validatePhysical(float64(value)*field.resolution+field.offset, value)
}

func (field *planField) validateUnsigned(value uint64) error {
	if field.bitLength == 0 || field.bitLength > 64 {
		return fmt.Errorf("field %d (%s): invalid unsigned width %d", field.order, field.name, field.bitLength)
	}
	if field.bitLength < 64 {
		maxValue := uint64(1)<<field.bitLength - 1
		if value > maxValue {
			return fmt.Errorf("field %d (%s): value %d exceeds %d-bit unsigned range [0,%d]", field.order, field.name, value, field.bitLength, maxValue)
		}
		_, declaredSentinel := field.sentinels[int64(value)]
		if field.kind == fieldKindNullableNumber && value == maxValue && !declaredSentinel {
			return fmt.Errorf("field %d (%s): value %d is the null sentinel; use nil", field.order, field.name, value)
		}
	}
	if value > math.MaxInt64 {
		return field.validatePhysical(float64(value)*field.resolution+field.offset, 0)
	}
	return field.validatePhysical(float64(value)*field.resolution+field.offset, int64(value))
}

func (field *planField) validatePhysical(value float64, raw int64) error {
	if field.kind != fieldKindNullableNumber && field.kind != fieldKindRawNumber {
		return nil
	}
	if _, ok := field.sentinels[raw]; ok {
		return nil
	}
	if field.rangeMin != nil && value < *field.rangeMin && !approximatelyEqual(value, *field.rangeMin) {
		return fmt.Errorf("field %d (%s): physical value %g is below minimum %g", field.order, field.name, value, *field.rangeMin)
	}
	if field.rangeMax != nil && value > *field.rangeMax && !approximatelyEqual(value, *field.rangeMax) {
		return fmt.Errorf("field %d (%s): physical value %g exceeds maximum %g", field.order, field.name, value, *field.rangeMax)
	}
	return nil
}

func approximatelyEqual(left, right float64) bool {
	scale := math.Max(1, math.Max(math.Abs(left), math.Abs(right)))
	return math.Abs(left-right) <= scale*1e-12
}

// addLengthOverride records the derived value of a binary field's length
// reference: bit count for BitLengthField references, byte count for
// DYNAMIC_FIELD_LENGTH references.
func addLengthOverride(overrides map[int]uint64, field *planField, source reflect.Value) {
	if field.kind != fieldKindBinary {
		return
	}
	length := uint64(source.Field(field.fieldIndex).Len())
	if field.refOrder != 0 {
		overrides[field.refOrder] = length * 8
	} else if field.dynLenOrder != 0 {
		overrides[field.dynLenOrder] = length
	}
}
