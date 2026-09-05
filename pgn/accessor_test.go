package pgn

import (
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
	m.SetHeadingValue(1.5708)
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
	m.SetRealPowerValue(-5)
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
	elem.SetOutputSpeedValue(2.5)
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
