package n2k

import (
	"fmt"
	"runtime/debug"

	"github.com/open-ships/n2k/pgn"
)

// ProductInfo describes this client as an NMEA 2000 product. Other devices
// (chartplotters, network analyzers) request it via PGN 126996 and display it
// in their device lists.
type ProductInfo struct {
	// N2KVersion is the supported NMEA 2000 standard version in thousandths,
	// e.g. 2101 means version 2.101. Zero means the default (2101).
	N2KVersion uint16
	// ProductCode is the manufacturer-assigned product code.
	ProductCode uint16
	// ModelID names the product (at most 32 bytes on the wire).
	ModelID string
	// SoftwareVersion is the software version string (at most 32 bytes).
	SoftwareVersion string
	// ModelVersion is the model version string (at most 32 bytes).
	ModelVersion string
	// SerialNumber is the device serial code (at most 32 bytes).
	SerialNumber string
	// CertificationLevel is the NMEA 2000 certification level lookup value.
	CertificationLevel uint8
	// LoadEquivalency is the device's bus load in units of 50 mA.
	LoadEquivalency uint8
}

// ConfigInfo describes this installation of the device (PGN 126998). All
// three fields are free-form strings.
type ConfigInfo struct {
	InstallationDescription1 string
	InstallationDescription2 string
	ManufacturerInformation  string
}

// defaultProductInfo builds the product identity used when WithProductInfo is
// not given. identity is the device NAME's identity number, used as a stable
// serial.
func defaultProductInfo(identity uint32) ProductInfo {
	return ProductInfo{
		N2KVersion:      2101,
		ProductCode:     2000,
		ModelID:         "open-ships n2k",
		SoftwareVersion: moduleVersion(),
		ModelVersion:    "n2k",
		SerialNumber:    fmt.Sprintf("%d", identity),
		LoadEquivalency: 1,
	}
}

// defaultConfigInfo builds the configuration info used when WithConfigInfo is
// not given.
func defaultConfigInfo() ConfigInfo {
	return ConfigInfo{
		ManufacturerInformation: "github.com/open-ships/n2k",
	}
}

// moduleVersion reports this library's module version when the embedding
// binary carries build info, else "dev".
func moduleVersion() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, dep := range info.Deps {
			if dep.Path == "github.com/open-ships/n2k" {
				return dep.Version
			}
		}
	}
	return "dev"
}

// message converts the ProductInfo into a wire-ready PGN 126996 struct.
func (p ProductInfo) message() *pgn.ProductInformation {
	version := uint64(p.N2KVersion)
	if version == 0 {
		version = 2101
	}
	code := uint64(p.ProductCode)
	cert := uint64(p.CertificationLevel)
	load := uint64(p.LoadEquivalency)
	return &pgn.ProductInformation{
		Nmea2000Version:     &version,
		ProductCode:         &code,
		ModelId:             p.ModelID,
		SoftwareVersionCode: p.SoftwareVersion,
		ModelVersion:        p.ModelVersion,
		ModelSerialCode:     p.SerialNumber,
		CertificationLevel:  &cert,
		LoadEquivalency:     &load,
	}
}

// message converts the ConfigInfo into a wire-ready PGN 126998 struct.
func (ci ConfigInfo) message() *pgn.ConfigurationInformation {
	return &pgn.ConfigurationInformation{
		InstallationDescription1: ci.InstallationDescription1,
		InstallationDescription2: ci.InstallationDescription2,
		ManufacturerInformation:  ci.ManufacturerInformation,
	}
}
