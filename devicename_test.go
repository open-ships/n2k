package n2k

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeviceName_Pack(t *testing.T) {
	// Round-trip: pack then unpack should return the original struct.
	original := DeviceName{
		IdentityNumber:   123456,
		ManufacturerCode: 999,
		DeviceInstance:   42,
		DeviceFunction:   130,
		DeviceClass:      25,
		SystemInstance:   3,
		IndustryGroup:    4,
	}

	packed := original.Pack(true)
	unpacked := UnpackDeviceName(packed)

	assert.Equal(t, original, unpacked)

	// Verify the arbitrary address capable bit is set.
	assert.Equal(t, uint64(1), (packed>>63)&1)

	// Round-trip without arbitrary address capable.
	packed2 := original.Pack(false)
	unpacked2 := UnpackDeviceName(packed2)

	assert.Equal(t, original, unpacked2)
	assert.Equal(t, uint64(0), (packed2>>63)&1)
}

func TestDeviceName_Validation(t *testing.T) {
	tests := []struct {
		name    string
		dn      DeviceName
		wantErr string
	}{
		{
			name: "valid device name",
			dn: DeviceName{
				IdentityNumber:   2097151,
				ManufacturerCode: 2047,
				DeviceInstance:   255,
				DeviceFunction:   255,
				DeviceClass:      127,
				SystemInstance:   15,
				IndustryGroup:    7,
			},
			wantErr: "",
		},
		{
			name:    "identity number too large",
			dn:      DeviceName{IdentityNumber: 2097152},
			wantErr: "IdentityNumber",
		},
		{
			name:    "manufacturer code too large",
			dn:      DeviceName{ManufacturerCode: 2048},
			wantErr: "ManufacturerCode",
		},
		{
			name:    "device class too large",
			dn:      DeviceName{DeviceClass: 128},
			wantErr: "DeviceClass",
		},
		{
			name:    "system instance too large",
			dn:      DeviceName{SystemInstance: 16},
			wantErr: "SystemInstance",
		},
		{
			name:    "industry group too large",
			dn:      DeviceName{IndustryGroup: 8},
			wantErr: "IndustryGroup",
		},
		{
			name:    "zero values are valid",
			dn:      DeviceName{},
			wantErr: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.dn.Validate()
			if tc.wantErr == "" {
				assert.NoError(t, err)
			} else {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErr)
			}
		})
	}
}

func TestDeviceName_Default(t *testing.T) {
	dn := DefaultDeviceName()

	assert.Equal(t, uint8(4), dn.IndustryGroup, "marine industry group")
	assert.Equal(t, uint16(2000), dn.ManufacturerCode, "experimental manufacturer code")
	assert.Equal(t, uint8(25), dn.DeviceClass, "internetwork device class")
	assert.Equal(t, uint8(130), dn.DeviceFunction, "PC gateway function")
	assert.Equal(t, uint8(0), dn.DeviceInstance)
	assert.Equal(t, uint8(0), dn.SystemInstance)

	// Identity number must fit in 21 bits.
	assert.LessOrEqual(t, dn.IdentityNumber, uint32(0x1FFFFF))

	// Default must pass validation.
	assert.NoError(t, dn.Validate())
}

func TestDeviceName_PackKnownValue(t *testing.T) {
	// Hand-calculated 64-bit NAME for a known DeviceName.
	dn := DeviceName{
		IdentityNumber:   1, // bits 0-20:  0x000001
		ManufacturerCode: 1, // bits 21-31: 1 << 21
		DeviceInstance:   0, // bits 32-39: 0
		DeviceFunction:   0, // bits 40-47: 0
		DeviceClass:      0, // bits 49-55: 0
		SystemInstance:   0, // bits 56-59: 0
		IndustryGroup:    0, // bits 60-62: 0
	}

	// Expected: bit 0 set (identity=1) + bit 21 set (mfr=1)
	// = 0x0000_0000_0020_0001
	expected := uint64(1) | (uint64(1) << 21)
	assert.Equal(t, expected, dn.Pack(false))

	// With arbitrary address capable: bit 63 also set.
	expectedAAC := expected | (uint64(1) << 63)
	assert.Equal(t, expectedAAC, dn.Pack(true))

	// More complex known value.
	dn2 := DeviceName{
		IdentityNumber:   0x1FFFFF, // all 21 bits set
		ManufacturerCode: 0x7FF,    // all 11 bits set
		DeviceInstance:   0xFF,     // all 8 bits set
		DeviceFunction:   0xFF,     // all 8 bits set
		DeviceClass:      0x7F,     // all 7 bits set
		SystemInstance:   0x0F,     // all 4 bits set
		IndustryGroup:    0x07,     // all 3 bits set
	}

	packed := dn2.Pack(true)
	// All field bits set + reserved bit 48 should be 0.
	// Bit 48 is the only zero bit in the upper portion.
	assert.Equal(t, uint64(0), (packed>>48)&1, "reserved bit 48 must be 0")

	// Everything else set means all bits except bit 48 are 1.
	allOnes := ^uint64(0)
	expectedAll := allOnes &^ (uint64(1) << 48) // clear bit 48
	assert.Equal(t, expectedAll, packed)
}

func TestDeviceName_NAMEComparison(t *testing.T) {
	// In NMEA 2000 address claiming, the device with the numerically lower
	// 64-bit NAME wins bus address contention.
	low := DeviceName{
		IdentityNumber:   1,
		ManufacturerCode: 100,
		DeviceInstance:   0,
		DeviceFunction:   130,
		DeviceClass:      25,
		SystemInstance:   0,
		IndustryGroup:    4,
	}
	high := DeviceName{
		IdentityNumber:   2,
		ManufacturerCode: 100,
		DeviceInstance:   0,
		DeviceFunction:   130,
		DeviceClass:      25,
		SystemInstance:   0,
		IndustryGroup:    4,
	}

	lowName := low.Pack(true)
	highName := high.Pack(true)

	assert.Less(t, lowName, highName, "lower NAME should win address contention")

	// Higher manufacturer code produces a higher NAME.
	highMfr := low
	highMfr.ManufacturerCode = 200
	assert.Less(t, low.Pack(true), highMfr.Pack(true))
}

func TestDeviceName_DeviceInstanceRoundTrip(t *testing.T) {
	// DeviceInstance splits across two bit ranges (lower 3 + upper 5).
	// Test all edge cases to ensure the split/recombine is correct.
	for _, di := range []uint8{0, 1, 7, 8, 15, 127, 128, 255} {
		dn := DeviceName{DeviceInstance: di}
		packed := dn.Pack(false)
		unpacked := UnpackDeviceName(packed)
		assert.Equal(t, di, unpacked.DeviceInstance, "DeviceInstance %d round-trip failed", di)
	}
}

func TestDeviceName_DefaultRandomness(t *testing.T) {
	// Two calls to DefaultDeviceName should (almost certainly) produce
	// different identity numbers.
	d1 := DefaultDeviceName()
	d2 := DefaultDeviceName()
	assert.NotEqual(t, d1.IdentityNumber, d2.IdentityNumber,
		"two default NAMEs should have different random identity numbers")
}
