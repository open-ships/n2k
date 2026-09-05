package pgn

import (
	"fmt"
	"reflect"
)

// PhysicalValue returns the physical (unit-scaled) value of the numeric field
// with the given source order on a decoded PGN struct, applying the field's
// Resolution and Offset from its metadata: value = raw*Resolution + Offset.
//
// Most callers should prefer the generated typed accessors instead: every
// numeric field with a physical interpretation has <Field>Value() (float64,
// bool) and Set<Field>Value(float64) methods on its struct (for example,
// VesselHeading.HeadingValue returns radians). PhysicalValue remains for
// dynamic, metadata-driven access when the field is only known at runtime.
// The unit string is the metadata Unit label ("rad", "m/s", "K", ...), empty
// for unitless fields.
//
// "Numeric" here means a field whose decoded Go representation is a raw wire
// tick count -- a *uint64 or *int64 struct field, per the codec's
// fieldKindNullableNumber, fieldKindRawNumber, and fieldKindLookup kinds
// (pgn/codec.go). Lookup/enumeration fields qualify and naturally carry
// Resolution 1 and Unit "", so their raw ordinal comes back unchanged.
//
// The following are deliberately NOT numeric for this function and return an
// error:
//   - Match-selector fields (fieldKindMatch): mechanically they decode into
//     the same *uint64/*int64 storage as plain numbers, but semantically they
//     pick a PGN variant rather than carry a measured quantity, so scaling
//     them by Resolution/Offset is not meaningful.
//   - FLOAT fields (fieldKindFloat, *float32 storage): already a physical
//     value on the wire, not a raw tick count, so there is nothing to scale.
//   - Strings, binary data, and reserved/spare padding.
//
// Fields inside a repeating group are not addressable by top-level source
// order alone -- a group can have any number of elements, and a bare order
// number does not select one -- so those orders return an error too. Decode
// the group slice field directly and inspect its elements instead.
//
// The bool result is false (with a nil error) when the field order names a
// numeric field but its decoded value is nil: the wire wrote the field's
// null/out-of-range sentinel, or the payload ended before reaching it.
//
// An error is returned when msg's struct type has no registered metadata,
// when fieldOrder does not name any field of that metadata, or when the
// named field is non-numeric or only reachable inside a repeating group.
func PhysicalValue(msg PGN, fieldOrder int) (float64, string, bool, error) {
	plan, err := codecPlanForMessage(msg)
	if err != nil {
		return 0, "", false, fmt.Errorf("PhysicalValue: %w", err)
	}
	desc, ok := plan.info.Fields[fieldOrder]
	if !ok || desc == nil {
		return 0, "", false, fmt.Errorf("PhysicalValue: %s has no field with order %d", plan.info.Id, fieldOrder)
	}

	field, insideGroup := findPlanField(plan, fieldOrder)
	if field == nil {
		if insideGroup {
			return 0, "", false, fmt.Errorf(
				"PhysicalValue: %s field %d (%s) is inside a repeating group and is not addressable by source order alone",
				plan.info.Id, fieldOrder, desc.Name)
		}
		return 0, "", false, fmt.Errorf("PhysicalValue: %s has no field with order %d", plan.info.Id, fieldOrder)
	}

	switch field.kind {
	case fieldKindLookup, fieldKindNullableNumber, fieldKindRawNumber:
		// Numeric: raw wire ticks in a *uint64/*int64 struct field.
	default:
		return 0, "", false, fmt.Errorf(
			"PhysicalValue: %s field %d (%s) is not a raw-ticks numeric field", plan.info.Id, fieldOrder, desc.Name)
	}

	target := reflect.ValueOf(msg).Elem()
	raw, present := numericFieldRaw(target.Field(field.fieldIndex), field.signed)
	if !present {
		return 0, "", false, nil
	}
	if field.kind != fieldKindLookup && !field.measurementAvailable(target.Field(field.fieldIndex), raw) {
		return 0, desc.Unit, false, nil
	}

	offset := 0.0
	if desc.Offset != nil {
		offset = *desc.Offset
	}
	return raw*float64(desc.Resolution) + offset, desc.Unit, true, nil
}

func (field *planField) measurementAvailable(value reflect.Value, raw float64) bool {
	if field.signed {
		ticks := value.Elem().Int()
		if _, special := field.sentinels[ticks]; special {
			return false
		}
		if field.kind == fieldKindNullableNumber && uint64(ticks) == (uint64(1)<<(field.bitLength-1))-1 {
			return false
		}
	} else {
		ticks := value.Elem().Uint()
		if _, special := field.sentinels[int64(ticks)]; ticks <= uint64(^uint64(0)>>1) && special {
			return false
		}
		if field.kind == fieldKindNullableNumber && ticks == ^uint64(0)>>(64-field.bitLength) {
			return false
		}
	}
	physical := raw*field.resolution + field.offset
	return (field.rangeMin == nil || physical >= *field.rangeMin || approximatelyEqual(physical, *field.rangeMin)) &&
		(field.rangeMax == nil || physical <= *field.rangeMax || approximatelyEqual(physical, *field.rangeMax))
}

// findPlanField locates the top-level planField for a metadata source order.
// It returns a nil field for orders that exist only inside a repeating group
// (insideGroup=true) or not at all (insideGroup=false); the caller has
// already confirmed the order exists in the metadata, so the latter case
// indicates an internal inconsistency between metadata and the compiled plan
// rather than a genuinely unknown order.
func findPlanField(plan *codecPlan, order int) (field *planField, insideGroup bool) {
	for _, step := range plan.steps {
		if step.field != nil {
			if step.field.order == order {
				return step.field, false
			}
			continue
		}
		for i := range step.group.fields {
			if step.group.fields[i].order == order {
				return nil, true
			}
		}
	}
	return nil, false
}

// numericFieldRaw reads a *uint64 or *int64 struct field's value as a
// float64. The second result is false when the pointer is nil.
func numericFieldRaw(fv reflect.Value, signed bool) (float64, bool) {
	if signed {
		value, _ := fv.Interface().(*int64)
		if value == nil {
			return 0, false
		}
		return float64(*value), true
	}
	value, _ := fv.Interface().(*uint64)
	if value == nil {
		return 0, false
	}
	return float64(*value), true
}
