package pgn

import (
	"errors"
	"math"
	"reflect"
	"strconv"
	"testing"
)

func TestHeadingValueAccessor(t *testing.T) {
	m := &VesselHeading{}
	if _, ok := m.HeadingValue(); ok {
		t.Fatal("HeadingValue() ok = true on nil field, want false")
	}

	raw := uint64(15708)
	m.Heading = &raw
	v, ok := m.HeadingValue()
	if !ok {
		t.Fatal("HeadingValue() ok = false with field set, want true")
	}
	if math.Abs(v-1.5708) > 1e-9 {
		t.Fatalf("HeadingValue() = %v, want 1.5708", v)
	}
}

func TestSetHeadingValueRoundsToTicks(t *testing.T) {
	m := &VesselHeading{}
	if err := m.SetHeadingValue(1.5708); err != nil {
		t.Fatal(err)
	}
	if m.Heading == nil {
		t.Fatal("SetHeadingValue left Heading nil")
	}
	if *m.Heading != 15708 {
		t.Fatalf("SetHeadingValue(1.5708) raw = %d, want 15708", *m.Heading)
	}

	v, ok := m.HeadingValue()
	if !ok || math.Abs(v-1.5708) > 1e-9 {
		t.Fatalf("round trip HeadingValue() = %v, %v, want 1.5708, true", v, ok)
	}
}

func TestOffsetValueAccessor(t *testing.T) {
	// RealPower is stored with an additive offset of -2e9 W.
	m := &UtilityPhaseCAcPower{}
	if err := m.SetRealPowerValue(-5); err != nil {
		t.Fatal(err)
	}
	if m.RealPower == nil {
		t.Fatal("SetRealPowerValue left RealPower nil")
	}
	if *m.RealPower != 1999999995 {
		t.Fatalf("SetRealPowerValue(-5) raw = %d, want 1999999995", *m.RealPower)
	}

	v, ok := m.RealPowerValue()
	if !ok || v != -5 {
		t.Fatalf("RealPowerValue() = %v, %v, want -5, true", v, ok)
	}
}

func TestRepeatingElementValueAccessor(t *testing.T) {
	elem := AirmarCalibrateSpeedRepeating1{}
	if _, ok := elem.OutputSpeedValue(); ok {
		t.Fatal("OutputSpeedValue() ok = true on nil field, want false")
	}
	if err := elem.SetOutputSpeedValue(2.5); err != nil {
		t.Fatal(err)
	}
	v, ok := elem.OutputSpeedValue()
	if !ok || math.Abs(v-2.5) > 1e-6 {
		t.Fatalf("OutputSpeedValue() = %v, %v, want 2.5, true", v, ok)
	}
}

// TestValueAccessorsMatchPhysicalValue cross-checks every generated
// <Field>Value accessor on every registered PGN struct against the
// metadata-driven PhysicalValue lookup, including availability.
func TestValueAccessorsMatchPhysicalValue(t *testing.T) {
	checked := 0
	for name, newMessage := range structTypeRegistry {
		msg := newMessage()
		target := reflect.ValueOf(msg).Elem()
		targetType := target.Type()
		for i := 0; i < targetType.NumField(); i++ {
			structField := targetType.Field(i)
			order, err := strconv.Atoi(structField.Tag.Get("n2k"))
			if err != nil {
				continue // Info field or repeating-group slice
			}
			method := reflect.ValueOf(msg).MethodByName(structField.Name + "Value")
			if !method.IsValid() {
				continue
			}
			// A sibling field pair like Temperature/SetTemperature makes
			// <Field>Value resolve to the sibling's setter; only inspect
			// getter-shaped methods.
			methodType := method.Type()
			if methodType.NumIn() != 0 || methodType.NumOut() != 2 ||
				methodType.Out(0).Kind() != reflect.Float64 || methodType.Out(1).Kind() != reflect.Bool {
				continue
			}

			switch structField.Type.String() {
			case "*uint64":
				raw := uint64(1)
				target.Field(i).Set(reflect.ValueOf(&raw))
			case "*int64":
				raw := int64(1)
				target.Field(i).Set(reflect.ValueOf(&raw))
			default:
				t.Errorf("%s.%sValue exists on non raw-ticks field type %s", name, structField.Name, structField.Type)
				continue
			}

			out := method.Call(nil)
			got, ok := out[0].Float(), out[1].Bool()
			want, _, wantOk, err := PhysicalValue(msg, order)
			if err != nil {
				t.Errorf("%s field %d: PhysicalValue error: %v", name, order, err)
				continue
			}
			if ok != wantOk {
				t.Errorf("%s field %d: accessor availability %t != PhysicalValue availability %t", name, order, ok, wantOk)
				continue
			}
			tolerance := math.Max(math.Abs(want)*1e-12, 1e-12)
			if math.Abs(got-want) > tolerance {
				t.Errorf("%s.%sValue() = %g, PhysicalValue = %g", name, structField.Name, got, want)
			}
			checked++
		}
	}
	if checked == 0 {
		t.Fatal("no generated value accessors found to check")
	}
	t.Logf("cross-checked %d value accessors", checked)
}

func TestPhysicalSettersRejectInvalidInputWithoutMutation(t *testing.T) {
	for _, value := range []float64{math.NaN(), math.Inf(1), math.Inf(-1), -1, 7, math.MaxFloat64} {
		msg := &VesselHeading{}
		if err := msg.SetHeadingValue(1.25); err != nil {
			t.Fatal(err)
		}
		before := *msg.Heading
		if err := msg.SetHeadingValue(value); !errors.Is(err, ErrInvalidPhysicalValue) {
			t.Fatalf("SetHeadingValue(%g): %v", value, err)
		}
		if *msg.Heading != before {
			t.Fatalf("invalid value %g mutated Heading", value)
		}
	}
	var nilHeading *VesselHeading
	if err := nilHeading.SetHeadingValue(1); !errors.Is(err, ErrInvalidPhysicalValue) {
		t.Fatal(err)
	}
	signed := &HeadingTrackControl{}
	if err := signed.SetCommandedRudderAngleValue(-0.25); err != nil {
		t.Fatal(err)
	}
	before := *signed.CommandedRudderAngle
	if err := signed.SetCommandedRudderAngleValue(math.Inf(1)); !errors.Is(err, ErrInvalidPhysicalValue) {
		t.Fatal(err)
	}
	if *signed.CommandedRudderAngle != before {
		t.Fatal("signed setter mutated on error")
	}
	repeating := &AirmarCalibrateSpeedRepeating1{}
	if err := repeating.SetOutputSpeedValue(2.5); err != nil {
		t.Fatal(err)
	}
	beforeUnsigned := *repeating.OutputSpeed
	if err := repeating.SetOutputSpeedValue(math.NaN()); !errors.Is(err, ErrInvalidPhysicalValue) {
		t.Fatal(err)
	}
	if *repeating.OutputSpeed != beforeUnsigned {
		t.Fatal("repeating setter mutated on error")
	}
}

func TestPhysicalRawTicksChecksIntegerBoundaries(t *testing.T) {
	for _, signed := range []bool{false, true} {
		bits := 64
		if signed {
			bits--
		}
		limit := math.Ldexp(1, bits)
		for _, invalid := range []float64{limit, math.Inf(1), math.NaN()} {
			if _, err := physicalRawTicks(invalid, 1, 0, 64, signed); err == nil {
				t.Fatalf("accepted %g (signed %v)", invalid, signed)
			}
		}
		if _, err := physicalRawTicks(math.Nextafter(limit, 0), 1, 0, 64, signed); err != nil {
			t.Fatal(err)
		}
	}
}

func TestAllPhysicalSettersRejectNonFiniteValues(t *testing.T) {
	checked := 0
	var checkStruct func(reflect.Value)
	checkStruct = func(pointer reflect.Value) {
		value := pointer.Elem()
		for i := 0; i < value.NumField(); i++ {
			field, fieldType := value.Field(i), value.Type().Field(i)
			if field.Kind() == reflect.Slice && field.Type().Elem().Kind() == reflect.Struct {
				checkStruct(reflect.New(field.Type().Elem()))
				continue
			}
			setter := pointer.MethodByName("Set" + fieldType.Name + "Value")
			if !setter.IsValid() {
				continue
			}
			for _, invalid := range []float64{math.NaN(), math.Inf(-1), math.Inf(1)} {
				result := setter.Call([]reflect.Value{reflect.ValueOf(invalid)})
				if len(result) != 1 || result[0].IsNil() {
					t.Fatalf("%s.Set%sValue accepted %g", value.Type(), fieldType.Name, invalid)
				}
				err, ok := result[0].Interface().(error)
				if !ok || !errors.Is(err, ErrInvalidPhysicalValue) || !field.IsNil() {
					t.Fatalf("%s.Set%sValue(%g) = %v; field = %v", value.Type(), fieldType.Name, invalid, err, field)
				}
			}
			checked++
		}
	}
	for _, newMessage := range structTypeRegistry {
		checkStruct(reflect.ValueOf(newMessage()))
	}
	if checked == 0 {
		t.Fatal("no generated setters checked")
	}
	t.Logf("checked %d physical setters including repeating fields", checked)
}

func TestOffsetSetterPreservesMeasurementOnFailure(t *testing.T) {
	m := &UtilityPhaseCAcPower{}
	if err := m.SetRealPowerValue(-5); err != nil {
		t.Fatal(err)
	}
	before := m.RealPower
	for _, invalid := range []float64{math.NaN(), math.Inf(1), -3e9, 3e9} {
		if err := m.SetRealPowerValue(invalid); !errors.Is(err, ErrInvalidPhysicalValue) {
			t.Fatalf("SetRealPowerValue(%g) = %v", invalid, err)
		}
		if m.RealPower != before || *m.RealPower != 1999999995 {
			t.Fatal("failed offset setter changed the original raw ticks")
		}
	}
}

func FuzzPhysicalHeadingSetter(f *testing.F) {
	for _, value := range []float64{0, 1.25, 6.2831852, math.NaN(), math.Inf(1), -1} {
		f.Add(value)
	}
	f.Fuzz(func(t *testing.T, value float64) {
		msg := &VesselHeading{}
		if err := msg.SetHeadingValue(1.25); err != nil {
			t.Fatal(err)
		}
		before := *msg.Heading
		err := msg.SetHeadingValue(value)
		if err != nil {
			if *msg.Heading != before {
				t.Fatal("failed setter mutated field")
			}
			return
		}
		if math.IsNaN(value) || math.IsInf(value, 0) {
			t.Fatal("accepted non-finite value")
		}
		if _, ok := msg.HeadingValue(); !ok {
			t.Fatal("successful setter produced unavailable measurement")
		}
		if _, err := msg.EncodePayload(); err != nil {
			t.Fatal(err)
		}
	})
}
