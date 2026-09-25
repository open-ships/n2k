package pgn

import (
	"bytes"
	"encoding/json"
	"math"
	"testing"
)

func marshalJSONObject(t *testing.T, value any) map[string]json.RawMessage {
	t.Helper()
	data, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	var object map[string]json.RawMessage
	if err := json.Unmarshal(data, &object); err != nil {
		t.Fatal(err)
	}
	return object
}

func checkJSONMeasurement(t *testing.T, object map[string]json.RawMessage, field string, want float64, unit string) {
	t.Helper()
	var physical map[string]json.RawMessage
	if err := json.Unmarshal(object["physical"], &physical); err != nil {
		t.Fatal(err)
	}
	var measurement struct {
		Value float64 `json:"value"`
		Unit  string  `json:"unit"`
	}
	if len(physical[field]) == 0 || string(physical[field]) == "null" {
		t.Fatalf("missing measurement %s: %s", field, object["physical"])
	}
	if err := json.Unmarshal(physical[field], &measurement); err != nil {
		t.Fatal(err)
	}
	if math.Abs(measurement.Value-want) > 1e-10 || measurement.Unit != unit {
		t.Fatalf("%s = %+v, want %g %s", field, measurement, want, unit)
	}
}

func TestJSONPhysicalMeasurements(t *testing.T) {
	heading, depth, speed, temperature := uint64(15708), uint64(270), uint64(500), uint64(29315)
	latitude, power := int64(-123456789), int64(1999999995)
	for _, test := range []struct {
		name, field, raw, unit string
		message                any
		value                  float64
	}{
		{"heading", "heading", "15708", "rad", VesselHeading{Heading: &heading}, 1.5708},
		{"depth", "depth", "270", "m", WaterDepth{Depth: &depth}, 2.7},
		{"speed", "windSpeed", "500", "m/s", WindData{WindSpeed: &speed}, 5},
		{"temperature", "actualTemperature", "29315", "K", Temperature{ActualTemperature: &temperature}, 293.15},
		{"position", "latitude", "-123456789", "deg", PositionRapidUpdate{Latitude: &latitude}, -12.3456789},
		{"offset", "realPower", "1999999995", "W", UtilityPhaseCAcPower{RealPower: &power}, -5},
	} {
		t.Run(test.name, func(t *testing.T) {
			object := marshalJSONObject(t, test.message)
			if string(object[test.field]) != test.raw {
				t.Fatalf("raw %s = %s, want %s", test.field, object[test.field], test.raw)
			}
			checkJSONMeasurement(t, object, test.field, test.value, test.unit)
		})
	}
}

func TestJSONPhysicalAvailability(t *testing.T) {
	for _, raw := range []uint64{0, 65533, 65534, 65535, 62832} {
		object := marshalJSONObject(t, VesselHeading{Heading: &raw, Sid: &raw, Reference: &raw})
		var physical map[string]json.RawMessage
		if err := json.Unmarshal(object["physical"], &physical); err != nil {
			t.Fatal(err)
		}
		if len(physical) != 3 || string(physical["deviation"]) != "null" || string(physical["variation"]) != "null" {
			t.Fatalf("physical measurements must exclude identifiers/enums and retain absent measurements: %s", object["physical"])
		}
		if raw == 0 {
			checkJSONMeasurement(t, object, "heading", 0, "rad")
		} else if string(physical["heading"]) != "null" {
			t.Fatalf("unavailable raw heading %d became %s", raw, physical["heading"])
		}
	}
	object := marshalJSONObject(t, VesselHeading{})
	if _, present := object["heading"]; present {
		t.Fatal("absent raw heading should retain omitempty behavior")
	}
	if string(object["physical"]) != `{"deviation":null,"heading":null,"variation":null}` {
		t.Fatalf("absent measurements: %s", object["physical"])
	}
}

func TestJSONPhysicalValueAndPointerEncoding(t *testing.T) {
	raw := uint64(15708)
	msg := VesselHeading{Heading: &raw}
	want, err := json.Marshal(&msg)
	if err != nil {
		t.Fatal(err)
	}
	for _, value := range []any{msg, &msg} {
		var stream bytes.Buffer
		if err := json.NewEncoder(&stream).Encode(value); err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(bytes.TrimSpace(stream.Bytes()), want) {
			t.Fatalf("stream/value output differs: %s", stream.Bytes())
		}
	}
	object := marshalJSONObject(t, map[string]VesselHeading{"message": msg})
	if !bytes.Equal(object["message"], want) {
		t.Fatalf("map value lost physical output: %s", object["message"])
	}
	data, err := json.Marshal((*VesselHeading)(nil))
	if err != nil || string(data) != "null" {
		t.Fatalf("nil message = %s, %v", data, err)
	}
}

func TestJSONPhysicalRepeatingEntries(t *testing.T) {
	latitude, longitude := int64(123456789), int64(-234567890)
	id := uint64(7)
	msg := NavigationRouteWpInformation{Repeating1: []NavigationRouteWpInformationRepeating1{
		{WpId: &id, WpLatitude: &latitude, WpLongitude: &longitude},
		{WpId: &id},
	}}
	object := marshalJSONObject(t, msg)
	if _, exists := object["physical"]; exists {
		t.Fatal("parent without physical measurements gained a physical object")
	}
	var entries []map[string]json.RawMessage
	if err := json.Unmarshal(object["repeating1"], &entries); err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 || string(entries[0]["wpId"]) != "7" || string(entries[0]["wpLatitude"]) != "123456789" {
		t.Fatalf("raw repeat entries changed: %s", object["repeating1"])
	}
	checkJSONMeasurement(t, entries[0], "wpLatitude", 12.3456789, "deg")
	checkJSONMeasurement(t, entries[0], "wpLongitude", -23.456789, "deg")
	if string(entries[1]["physical"]) != `{"wpLatitude":null,"wpLongitude":null}` {
		t.Fatalf("absent repeat measurements: %s", entries[1]["physical"])
	}
}

func TestJSONPhysicalIsDerivedAndPreservesWirePayload(t *testing.T) {
	payload := []byte{1, 0x5c, 0x3d, 0xff, 0x7f, 0xff, 0x7f, 0xfd}
	msg := &VesselHeading{}
	if err := msg.DecodePayload(payload); err != nil {
		t.Fatal(err)
	}
	checkJSONMeasurement(t, marshalJSONObject(t, msg), "heading", 1.5708, "rad")
	after, err := msg.EncodePayload()
	if err != nil || !bytes.Equal(after, payload) {
		t.Fatalf("JSON changed wire payload: %x, %v", after, err)
	}
	*msg.Heading = 10000
	checkJSONMeasurement(t, marshalJSONObject(t, msg), "heading", 1, "rad")
	if err := json.Unmarshal([]byte(`{"heading":20000,"physical":{"heading":{"value":5,"unit":"rad"}}}`), msg); err != nil {
		t.Fatal(err)
	}
	checkJSONMeasurement(t, marshalJSONObject(t, msg), "heading", 2, "rad")
	if *msg.Heading != 20000 {
		t.Fatal("physical input overrode raw heading")
	}

	// The JSON path must not round raw integers through a float64 intermediary.
	latitude := int64(123456789012345678)
	object := marshalJSONObject(t, GnssPositionData{Latitude: &latitude})
	if string(object["latitude"]) != "123456789012345678" {
		t.Fatalf("raw integer lost precision: %s", object["latitude"])
	}
}
