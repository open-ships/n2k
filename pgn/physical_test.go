package pgn

import (
	"math"
	"testing"
)

// TestPhysicalValueAppliesResolution is the base case from the task brief:
// VesselHeading field order 2 is Heading, a 16-bit unsigned NUMBER field with
// Resolution 0.0001 and Unit "rad" (see upstream_definitions.go). 15708 raw
// ticks should scale to ~1.5708 rad.
//
// Resolution in FieldDescriptor is float32 (metadata_runtime.go), so
// float64(raw)*float64(float32(0.0001)) accumulates float32-rounding error on
// the order of 4e-8 relative to the exact decimal product -- confirmed by
// direct computation, not just tolerance-widening on faith. The brief's
// original 1e-9 tolerance is tighter than that rounding error and would flake
// this test, so the tolerance here is widened to 1e-6 (documented per the
// task's ambiguity-resolution instructions); FieldDescriptor.Resolution stays
// float32 as-is.
func TestPhysicalValueAppliesResolution(t *testing.T) {
	h := uint64(15708)
	m := &VesselHeading{Heading: &h}
	v, unit, ok, err := PhysicalValue(m, 2) // field order 2 = Heading, resolution 0.0001 rad
	if err != nil || !ok {
		t.Fatalf("PhysicalValue: ok=%v err=%v", ok, err)
	}
	if unit != "rad" {
		t.Fatalf("unit = %q, want rad", unit)
	}
	if math.Abs(v-1.5708) > 1e-6 {
		t.Fatalf("v = %v, want 1.5708", v)
	}
}

func TestPhysicalValueNullField(t *testing.T) {
	m := &VesselHeading{}
	_, _, ok, err := PhysicalValue(m, 2)
	if err != nil || ok {
		t.Fatalf("want ok=false err=nil, got ok=%v err=%v", ok, err)
	}
}

// TestPhysicalValueSignedFieldWithOffset covers a signed field with a
// non-zero additive Offset: UtilityPhaseCAcPower field order 1 (Real Power)
// is a signed 32-bit NUMBER with Resolution 1 and Offset -2e9 (see
// upstream_definitions.go), so the physical value is raw + offset.
func TestPhysicalValueSignedFieldWithOffset(t *testing.T) {
	raw := int64(2_000_000_100)
	m := &UtilityPhaseCAcPower{RealPower: &raw}
	v, unit, ok, err := PhysicalValue(m, 1)
	if err != nil || !ok {
		t.Fatalf("PhysicalValue: ok=%v err=%v", ok, err)
	}
	if unit != "W" {
		t.Fatalf("unit = %q, want W", unit)
	}
	if math.Abs(v-100) > 1e-6 {
		t.Fatalf("v = %v, want 100", v)
	}
}

// TestPhysicalValueLookupField covers a LOOKUP-kind field: VesselHeading
// field order 5 (Reference) has Resolution 1 and Unit "", so its physical
// value is just the raw enumeration ordinal.
func TestPhysicalValueLookupField(t *testing.T) {
	ref := uint64(1)
	m := &VesselHeading{Reference: &ref}
	v, unit, ok, err := PhysicalValue(m, 5)
	if err != nil || !ok {
		t.Fatalf("PhysicalValue: ok=%v err=%v", ok, err)
	}
	if unit != "" {
		t.Fatalf("unit = %q, want empty", unit)
	}
	if v != 1 {
		t.Fatalf("v = %v, want 1", v)
	}
}

// TestPhysicalValueUnknownFieldOrder covers an order with no field
// descriptor at all on the given struct type.
func TestPhysicalValueUnknownFieldOrder(t *testing.T) {
	m := &VesselHeading{}
	_, _, _, err := PhysicalValue(m, 99)
	if err == nil {
		t.Fatal("want error for unknown field order, got nil")
	}
}

// TestPhysicalValueNonNumericField covers a RESERVED field: VesselHeading
// field order 6 is reserved padding, not a decoded struct field, so it is
// not numeric.
func TestPhysicalValueNonNumericField(t *testing.T) {
	m := &VesselHeading{}
	_, _, _, err := PhysicalValue(m, 6)
	if err == nil {
		t.Fatal("want error for non-numeric field, got nil")
	}
}

// TestPhysicalValueMatchFieldNotNumeric covers a Match-selector field:
// NmeaAcknowledgeGroupFunction field order 1 (Function Code) carries a Match
// value used to pick this PGN variant during decode. Match fields are
// excluded from "numeric" for PhysicalValue even though they decode into the
// same *uint64/*int64 storage as plain numbers and lookups -- see the
// PhysicalValue doc comment for the rationale.
func TestPhysicalValueMatchFieldNotNumeric(t *testing.T) {
	m := &NmeaAcknowledgeGroupFunction{}
	_, _, _, err := PhysicalValue(m, 1)
	if err == nil {
		t.Fatal("want error for Match-selector field, got nil")
	}
}

// TestPhysicalValueRepeatingGroupFieldNotAddressable covers a field order
// that only exists inside a repeating group: NmeaAcknowledgeGroupFunction
// field order 6 (Parameter) is the sole member of repeating field set 1
// (RepeatingFieldSet1StartField=6, Size=1). A source order alone cannot pick
// a group element, so PhysicalValue must reject it.
func TestPhysicalValueRepeatingGroupFieldNotAddressable(t *testing.T) {
	m := &NmeaAcknowledgeGroupFunction{}
	_, _, _, err := PhysicalValue(m, 6)
	if err == nil {
		t.Fatal("want error for repeating-group field order, got nil")
	}
}
