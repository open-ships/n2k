package n2k

import (
	"fmt"
	"math/rand/v2"
)

// DeviceName represents the ISO 11783 / NMEA 2000 64-bit NAME that uniquely
// identifies a device on a CAN bus network. The NAME is used during address
// claiming (PGN 60928) to arbitrate bus addresses.
type DeviceName struct {
	IdentityNumber   uint32 // 21 bits (0-20)
	ManufacturerCode uint16 // 11 bits (21-31)
	DeviceInstance   uint8  // 8 bits: lower 3 (32-34) + upper 5 (35-39)
	DeviceFunction   uint8  // 8 bits (40-47)
	DeviceClass      uint8  // 7 bits (49-55), bit 48 is reserved
	SystemInstance   uint8  // 4 bits (56-59)
	IndustryGroup    uint8  // 3 bits (60-62)
}

// Validate checks that each field fits within its bit width as defined by the
// ISO 11783 NAME specification. DeviceInstance and DeviceFunction are uint8 and
// always fit their 8-bit fields.
func (d DeviceName) Validate() error {
	if d.IdentityNumber > 0x1FFFFF { // 2^21 - 1
		return fmt.Errorf("n2k: IdentityNumber %d exceeds 21-bit maximum (2097151)", d.IdentityNumber)
	}
	if d.ManufacturerCode > 0x7FF { // 2^11 - 1
		return fmt.Errorf("n2k: ManufacturerCode %d exceeds 11-bit maximum (2047)", d.ManufacturerCode)
	}
	if d.DeviceClass > 0x7F { // 2^7 - 1
		return fmt.Errorf("n2k: DeviceClass %d exceeds 7-bit maximum (127)", d.DeviceClass)
	}
	if d.SystemInstance > 0x0F { // 2^4 - 1
		return fmt.Errorf("n2k: SystemInstance %d exceeds 4-bit maximum (15)", d.SystemInstance)
	}
	if d.IndustryGroup > 0x07 { // 2^3 - 1
		return fmt.Errorf("n2k: IndustryGroup %d exceeds 3-bit maximum (7)", d.IndustryGroup)
	}
	return nil
}

// Pack encodes the DeviceName into the 64-bit NAME integer defined by ISO 11783.
// The arbitraryAddressCapable flag sets bit 63, indicating the device can
// participate in dynamic address negotiation.
//
// 64-bit NAME layout (LSB-first):
//
//	Bits  0-20:  Identity Number (21 bits)
//	Bits 21-31:  Manufacturer Code (11 bits)
//	Bits 32-34:  Device Instance Lower (3 bits)
//	Bits 35-39:  Device Instance Upper (5 bits)
//	Bits 40-47:  Device Function (8 bits)
//	Bit  48:     Reserved (0)
//	Bits 49-55:  Device Class (7 bits)
//	Bits 56-59:  System Instance (4 bits)
//	Bits 60-62:  Industry Group (3 bits)
//	Bit  63:     Arbitrary Address Capable (1 bit)
func (d DeviceName) Pack(arbitraryAddressCapable bool) uint64 {
	var name uint64
	name |= uint64(d.IdentityNumber) & 0x1FFFFF              // bits 0-20
	name |= uint64(d.ManufacturerCode) << 21 & (0x7FF << 21)  // bits 21-31
	name |= uint64(d.DeviceInstance&0x07) << 32                // bits 32-34 (lower 3)
	name |= uint64(d.DeviceInstance>>3) << 35                  // bits 35-39 (upper 5)
	name |= uint64(d.DeviceFunction) << 40                     // bits 40-47
	// bit 48 reserved (0)
	name |= uint64(d.DeviceClass&0x7F) << 49                   // bits 49-55
	name |= uint64(d.SystemInstance&0x0F) << 56                 // bits 56-59
	name |= uint64(d.IndustryGroup&0x07) << 60                  // bits 60-62
	if arbitraryAddressCapable {
		name |= 1 << 63 // bit 63
	}
	return name
}

// UnpackDeviceName extracts a DeviceName from the 64-bit NAME integer defined
// by ISO 11783. The arbitraryAddressCapable flag (bit 63) is not stored in the
// struct; callers can check it directly via (name >> 63) & 1.
func UnpackDeviceName(name uint64) DeviceName {
	lower := uint8((name >> 32) & 0x07) // bits 32-34
	upper := uint8((name >> 35) & 0x1F) // bits 35-39
	return DeviceName{
		IdentityNumber:   uint32(name & 0x1FFFFF),         // bits 0-20
		ManufacturerCode: uint16((name >> 21) & 0x7FF),     // bits 21-31
		DeviceInstance:   (upper << 3) | lower,             // recombine 8 bits
		DeviceFunction:   uint8((name >> 40) & 0xFF),       // bits 40-47
		DeviceClass:      uint8((name >> 49) & 0x7F),       // bits 49-55
		SystemInstance:   uint8((name >> 56) & 0x0F),       // bits 56-59
		IndustryGroup:    uint8((name >> 60) & 0x07),       // bits 60-62
	}
}

// DefaultDeviceName returns a DeviceName pre-configured for a PC-based software
// gateway operating on an NMEA 2000 marine network. The identity number is
// randomized so that multiple processes using the same binary can coexist on
// the same bus without NAME collisions during address claiming.
func DefaultDeviceName() DeviceName {
	return DeviceName{
		// Marine industry group per NMEA 2000 specification.
		IndustryGroup: 4,

		// Manufacturer code 2000 is in the unassigned/experimental range (1864-2045).
		// Codes 2046-2047 are reserved sentinels in the NMEA 2000 spec. Using an
		// unassigned code avoids conflicts with real manufacturers while being clearly
		// identifiable as a software gateway rather than a hardware device.
		ManufacturerCode: 2000,

		// Device Class 25 = "Internetwork Device" per NMEA 2000 device class table.
		// This is the standard classification for software gateways and protocol bridges
		// that operate on NMEA 2000 networks.
		DeviceClass: 25,

		// Device Function 130 = "PC Gateway" within Device Class 25.
		// This is the standard function code for PC-based software that communicates
		// on NMEA 2000 networks via a CAN bus interface.
		DeviceFunction: 130,

		// Random identity number ensures uniqueness across multiple client instances
		// on the same NMEA 2000 network. Without randomization, two clients started
		// from the same binary would have identical NAMEs and one would fail address
		// claiming. The 21-bit field supports up to 2,097,151 unique values.
		IdentityNumber: randomIdentityNumber(),

		DeviceInstance: 0,
		SystemInstance: 0,
	}
}

// randomIdentityNumber returns a random 21-bit value suitable for use as an
// ISO 11783 NAME identity number.
func randomIdentityNumber() uint32 {
	return rand.Uint32() & 0x1FFFFF
}
