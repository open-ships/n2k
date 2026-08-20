package n2k

import (
	"errors"
	"fmt"
)

// ActisenseParity configures parity on a directly connected Actisense serial
// Adapter without exposing the implementation's serial library.
type ActisenseParity uint8

const (
	ActisenseParityNone ActisenseParity = iota
	ActisenseParityOdd
	ActisenseParityEven
	ActisenseParityMark
	ActisenseParitySpace
)

// ActisenseStopBits configures stop bits on a directly connected Actisense
// serial Adapter.
type ActisenseStopBits uint8

const (
	ActisenseStopBitsOne ActisenseStopBits = iota
	ActisenseStopBitsOnePointFive
	ActisenseStopBitsTwo
)

// ActisenseSerialConfig is the complete host-side serial configuration. Its
// zero value selects the Actisense binary default of 115200 8N1.
type ActisenseSerialConfig struct {
	BaudRate int
	DataBits int
	Parity   ActisenseParity
	StopBits ActisenseStopBits
}

func defaultActisenseSerialConfig() ActisenseSerialConfig {
	return ActisenseSerialConfig{
		BaudRate: 115200,
		DataBits: 8,
		Parity:   ActisenseParityNone,
		StopBits: ActisenseStopBitsOne,
	}
}

func normalizeActisenseSerialConfig(config ActisenseSerialConfig) ActisenseSerialConfig {
	defaults := defaultActisenseSerialConfig()
	if config.BaudRate == 0 {
		config.BaudRate = defaults.BaudRate
	}
	if config.DataBits == 0 {
		config.DataBits = defaults.DataBits
	}
	return config
}

func (c ActisenseSerialConfig) validate() error {
	if c.BaudRate <= 0 || c.BaudRate > 4_000_000 {
		return fmt.Errorf("n2k: Actisense serial baud rate %d is outside 1-4000000", c.BaudRate)
	}
	if c.DataBits < 5 || c.DataBits > 8 {
		return fmt.Errorf("n2k: Actisense serial data bits %d is outside 5-8", c.DataBits)
	}
	if c.Parity > ActisenseParitySpace {
		return fmt.Errorf("n2k: unknown Actisense serial parity %d", c.Parity)
	}
	if c.StopBits > ActisenseStopBitsTwo {
		return fmt.Errorf("n2k: unknown Actisense serial stop bits %d", c.StopBits)
	}
	return nil
}

// ActisenseSerialOption configures one direct serial Adapter.
type ActisenseSerialOption interface {
	applyActisenseSerial(*ActisenseSerialConfig)
}

type actisenseSerialOptionFunc func(*ActisenseSerialConfig)

func (f actisenseSerialOptionFunc) applyActisenseSerial(config *ActisenseSerialConfig) { f(config) }

// WithActisenseSerialConfig replaces the host-side serial configuration. Zero
// BaudRate or DataBits fields retain the 115200 8N1 defaults.
func WithActisenseSerialConfig(config ActisenseSerialConfig) ActisenseSerialOption {
	return actisenseSerialOptionFunc(func(target *ActisenseSerialConfig) {
		*target = normalizeActisenseSerialConfig(config)
	})
}

// WithActisenseBaudRate changes only the host-side baud rate.
func WithActisenseBaudRate(baudRate int) ActisenseSerialOption {
	return actisenseSerialOptionFunc(func(config *ActisenseSerialConfig) {
		config.BaudRate = baudRate
	})
}

func applyActisenseSerialOptions(options []ActisenseSerialOption) (ActisenseSerialConfig, error) {
	config := defaultActisenseSerialConfig()
	for _, option := range options {
		if option == nil {
			return ActisenseSerialConfig{}, errors.New("n2k: nil ActisenseSerialOption")
		}
		option.applyActisenseSerial(&config)
	}
	return normalizeActisenseSerialConfig(config), nil
}
