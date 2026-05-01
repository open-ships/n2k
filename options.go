package n2k

import (
	"errors"
	"log/slog"

	"github.com/brutella/can"
)

type config struct {
	sources        []source
	filterExpr     string
	includeUnknown bool
	logger         *slog.Logger
	sourceAddress  *uint8      // nil = auto mode
	deviceName     *DeviceName // nil = use default
}

func (c *config) validate() error {
	if len(c.sources) == 0 {
		return errors.New("n2k: at least one source (CAN, USB, or Replay) is required")
	}
	return nil
}

// Option configures the behavior of Receive and NewScanner.
type Option interface {
	apply(*config)
}

type optionFunc func(*config)

func (f optionFunc) apply(c *config) { f(c) }

// CAN adds a SocketCAN source for the given Linux CAN interface (e.g., "can0").
func CAN(iface string) Option {
	return optionFunc(func(c *config) {
		c.sources = append(c.sources, &socketCANSource{iface: iface})
	})
}

// USB adds a USB-CAN Analyzer source for the given serial port (e.g., "/dev/ttyUSB0").
func USB(port string) Option {
	return optionFunc(func(c *config) {
		c.sources = append(c.sources, &usbCANSource{port: port})
	})
}

// Replay adds a source that replays the given CAN frames. Useful for testing.
func Replay(frames []can.Frame) Option {
	return optionFunc(func(c *config) {
		c.sources = append(c.sources, &replaySource{frames: frames})
	})
}

// Filter sets a CEL expression to filter messages. The expression is automatically
// partitioned into pre-decode (metadata) and post-decode (struct field) stages.
func Filter(expr string) Option {
	return optionFunc(func(c *config) {
		c.filterExpr = expr
	})
}

// IncludeUnknown includes undecodable messages as *pgn.UnknownPGN in the output stream.
// By default, unknown PGNs are dropped and logged at debug level.
func IncludeUnknown() Option {
	return optionFunc(func(c *config) {
		c.includeUnknown = true
	})
}

// WithLogger overrides the default slog.Default() logger.
func WithLogger(l *slog.Logger) Option {
	return optionFunc(func(c *config) {
		c.logger = l
	})
}

// WithSourceAddress sets an explicit NMEA 2000 source address for the client.
// When set, the client uses this address and treats contention as a fatal error.
// When not set (default), the client uses auto mode — starting at address 253
// and working downward if contention occurs.
func WithSourceAddress(addr uint8) Option {
	return optionFunc(func(c *config) {
		c.sourceAddress = &addr
	})
}

// WithName sets the ISO 11783 device NAME used for address claiming.
// The NAME is a 64-bit identifier that uniquely identifies this device on the
// NMEA 2000 network. In address contention, the device with the lower NAME wins.
// When not set, a default NAME is used (see DefaultDeviceName).
func WithName(name DeviceName) Option {
	return optionFunc(func(c *config) {
		c.deviceName = &name
	})
}
