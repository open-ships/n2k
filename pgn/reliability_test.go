package pgn

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"reflect"
	"strings"
	"testing"
)

func TestNilMessageOperationsReturnErrors(t *testing.T) {
	var msg *VesselHeading
	if _, err := EncodeMessage(msg); err == nil {
		t.Fatal("typed nil encoded without error")
	}
	if _, err := msg.EncodePayload(); err == nil {
		t.Fatal("typed nil payload encoded without error")
	}
	if err := msg.DecodePayload([]byte{1}); err == nil {
		t.Fatal("typed nil decoded without error")
	}
	if _, err := CloneMessage(msg); err == nil {
		t.Fatal("typed nil cloned without error")
	}
}

func TestLengthStringIndependentBoundaries(t *testing.T) {
	for _, length := range []int{0, 1, 31, 32, 33, 63, 64, 127, 128, 254} {
		t.Run(fmt.Sprint(length), func(t *testing.T) {
			want := strings.Repeat("x", length)
			payload := append([]byte{byte(length)}, []byte(want)...)
			payload = append(payload, 0x42)
			stream := NewPgnDataStream(payload)
			got, err := stream.readStringWithLength()
			if err != nil || got != want {
				t.Fatalf("length %d: got %q, %v; want %q", length, got, err, want)
			}
			next, err := stream.getNumberRaw(8)
			if err != nil || next != 0x42 {
				t.Fatalf("following field = %#x, %v; string consumed wrong extent", next, err)
			}
		})
	}
}

func TestMeasurementSentinelsRemainInspectableButUnavailable(t *testing.T) {
	for _, ticks := range []uint16{65533, 65534, 65535} {
		t.Run(fmt.Sprint(ticks), func(t *testing.T) {
			payload := []byte{7, byte(ticks), byte(ticks >> 8), 0xff, 0x7f, 0xff, 0x7f, 0xfc}
			var msg VesselHeading
			if err := msg.DecodePayload(payload); err != nil {
				t.Fatal(err)
			}
			if ticks != 65535 && (msg.Heading == nil || *msg.Heading != uint64(ticks)) {
				t.Fatal("special raw ticks were lost")
			}
			if value, valid := msg.HeadingValue(); valid {
				t.Fatalf("sentinel exposed as valid heading %g", value)
			}
			if _, _, valid, err := PhysicalValue(&msg, 2); err != nil || valid {
				t.Fatalf("dynamic availability = %t, %v", valid, err)
			}
			encoded, err := msg.EncodePayload()
			if err != nil || !bytes.Equal(encoded, payload) {
				t.Fatalf("sentinel replay = %x, %v", encoded, err)
			}
		})
	}
}

func TestCoordinateIndependentBoundaries(t *testing.T) {
	for _, latitude := range []float64{-90, -89.9999999, 89.9999999, 90} {
		var msg PositionRapidUpdate
		msg.SetLatitudeValue(latitude)
		msg.SetLongitudeValue(180)
		payload, err := msg.EncodePayload()
		if err != nil {
			t.Fatalf("latitude %g: %v", latitude, err)
		}
		if ticks := int32(binary.LittleEndian.Uint32(payload)); ticks != int32(math.Round(latitude*1e7)) {
			t.Fatalf("latitude %g encoded as %d ticks", latitude, ticks)
		}
		if longitude := int32(binary.LittleEndian.Uint32(payload[4:])); longitude != 1800000000 {
			t.Fatalf("longitude = %d ticks", longitude)
		}
	}
}

func TestInvalidPhysicalPayloadReplaysUntilMutation(t *testing.T) {
	payload := make([]byte, 8)
	binary.LittleEndian.PutUint32(payload, 900000001) // 90.0000001 degrees, one tick beyond latitude range.
	var msg PositionRapidUpdate
	if err := msg.DecodePayload(payload); err != nil {
		t.Fatal(err)
	}
	if msg.Latitude == nil || *msg.Latitude != 900000001 {
		t.Fatal("invalid raw latitude lost")
	}
	if value, valid := msg.LatitudeValue(); valid {
		t.Fatalf("invalid latitude exposed as %g", value)
	}
	encoded, err := msg.EncodePayload()
	if err != nil || !bytes.Equal(encoded, payload) {
		t.Fatalf("unchanged invalid payload replay = %x, %v", encoded, err)
	}
	msg.SetLongitudeValue(1)
	if _, err := msg.EncodePayload(); err == nil {
		t.Fatal("mutated message bypassed outbound latitude validation")
	}
	msg.SetLatitudeValue(90)
	if _, err := msg.EncodePayload(); err != nil {
		t.Fatalf("corrected message did not encode: %v", err)
	}
}

func TestDecodeReuseClearsAbsentFieldsAndCommitsAtomically(t *testing.T) {
	full := []byte{7, 0x5c, 0x3d, 1, 0, 2, 0, 0xfc}
	var msg VesselHeading
	if err := msg.DecodePayload(full); err != nil {
		t.Fatal(err)
	}
	before := msg.Clone()
	if err := msg.DecodePayload([]byte{1, 2}); err == nil {
		t.Fatal("partial numeric field must fail")
	}
	if !reflect.DeepEqual(&msg, before) {
		t.Fatal("failed decode modified existing message")
	}
	if err := msg.DecodePayload([]byte{1}); err != nil {
		t.Fatal(err)
	}
	if msg.Heading != nil || msg.Deviation != nil || msg.Variation != nil || msg.Reference != nil {
		t.Fatal("short decode retained old measurements")
	}
	encoded, err := msg.EncodePayload()
	if err != nil || !bytes.Equal(encoded, []byte{1}) {
		t.Fatalf("short payload replay = %x, %v", encoded, err)
	}
}

func TestCloneOwnsNestedFieldsAndWireState(t *testing.T) {
	parameter := uint64(2)
	msg := &NmeaCommandGroupFunction{
		Info: MessageInfo{Priority: Priority(3), TargetId: Target(22), ConnectionEpoch: 8, ClaimEpoch: 9,
			rawPayload: []byte{1, 2}, rawCanonical: []byte{3, 4}, DecodeIssues: []string{"partial"}},
		Repeating1: []NmeaCommandGroupFunctionRepeating1{{Parameter: &parameter, Value: []byte{0x5c, 0x3d}}},
	}
	owned, err := CloneMessage(msg)
	if err != nil {
		t.Fatal(err)
	}
	clone := owned.(*NmeaCommandGroupFunction)
	if !reflect.DeepEqual(msg, clone) {
		t.Fatal("clone changed message content")
	}
	*clone.Info.Priority = 7
	*clone.Info.TargetId = 1
	clone.Info.rawPayload[0] = 8
	clone.Info.rawCanonical[0] = 8
	clone.Info.DecodeIssues[0] = "changed"
	*clone.Repeating1[0].Parameter = 9
	clone.Repeating1[0].Value[0] = 0
	if *msg.Info.Priority != 3 || *msg.Info.TargetId != 22 || msg.Info.rawPayload[0] != 1 || msg.Info.rawCanonical[0] != 3 || msg.Info.DecodeIssues[0] != "partial" || *msg.Repeating1[0].Parameter != 2 || msg.Repeating1[0].Value[0] != 0x5c {
		t.Fatal("clone shares mutable state")
	}
}

func TestGroupFunctionResolvesIndependentParameterWidths(t *testing.T) {
	// Command Vessel Heading, two parameters: Heading (16 bits) and Reference
	// (2 bits, occupying one whole byte in a group function).
	payload := []byte{1, 0x12, 0xf1, 1, 0xf3, 2, 2, 0x5c, 0x3d, 5, 1}
	dispatched, err := DecodeMessage(MessageInfo{PGN: 126208}, payload)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := dispatched.(*NmeaCommandGroupFunction); !ok {
		t.Fatalf("generic fallback shadowed recognized group function: %T", dispatched)
	}
	var msg NmeaCommandGroupFunction
	if err := msg.DecodePayload(payload); err != nil {
		t.Fatal(err)
	}
	if len(msg.Repeating1) != 2 || !bytes.Equal(msg.Repeating1[0].Value, []byte{0x5c, 0x3d}) || !bytes.Equal(msg.Repeating1[1].Value, []byte{1}) {
		t.Fatalf("wrong parameter boundaries: %+v", msg.Repeating1)
	}
	encoded, err := msg.EncodePayload()
	if err != nil || !bytes.Equal(encoded, payload) {
		t.Fatalf("group-function replay = %x, %v", encoded, err)
	}
	// Unknown commanded PGN: preserve bytes and explicitly expose partial decode.
	payload[1], payload[2], payload[3] = 0, 0, 0
	if err := msg.DecodePayload(payload); err != nil {
		t.Fatal(err)
	}
	if len(msg.Repeating1) != 0 || !strings.Contains(strings.Join(msg.Info.DecodeIssues, ";"), "unsupported dynamic field") {
		t.Fatalf("unresolved payload lacks explicit partial result: %+v", msg)
	}
	encoded, err = msg.EncodePayload()
	if err != nil || !bytes.Equal(encoded, payload) {
		t.Fatalf("partial payload replay = %x, %v", encoded, err)
	}
}

func TestOpaqueRegisterTailPreservedAsData(t *testing.T) {
	// Manufacturer 358 (Victron), marine industry, register 0x1234, four opaque bytes.
	payload := []byte{0x66, 0x99, 0x34, 0x12, 0xde, 0xad, 0xbe, 0xef}
	var msg VictronVeCanRegister
	if err := msg.DecodePayload(payload); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(msg.Value, payload[4:]) {
		t.Fatalf("register value = %x", msg.Value)
	}
}
