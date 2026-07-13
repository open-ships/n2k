package n2k

import (
	"errors"
	"log/slog"
	"time"

	"github.com/brutella/can"
)

type config struct {
	sources           []source
	filterExpr        string
	includeUnknown    bool
	logger            *slog.Logger
	sourceAddress     *uint8         // nil = auto mode
	deviceName        *DeviceName    // nil = use default
	claimTimeout      *time.Duration // nil = use default (1500ms)
	heartbeatInterval *time.Duration // nil = default (60s), 0 = disabled
	productInfo       *ProductInfo   // nil = defaults
	configInfo        *ConfigInfo    // nil = defaults
	bus               Bus            // pre-constructed bus
}

func (c *config) validate() error {
	if len(c.sources) == 0 && c.bus == nil {
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

// File adds a source that replays CAN frames from a candump -L / -l log file.
// By default frames are delivered as fast as they can be read; pass
// OriginalTiming() to pace them by the log's timestamps. File sources are
// read-only: they work with Receive and NewScanner but not NewClient.
func File(path string, opts ...FileOption) Option {
	return optionFunc(func(c *config) {
		src := &fileSource{path: path}
		for _, o := range opts {
			o.applyFile(src)
		}
		c.sources = append(c.sources, src)
	})
}

// TCP adds a source that dials a network gateway (e.g. a Yacht Devices
// YDWG-02 in RAW server mode, or an Actisense gateway) at addr ("host:port")
// and reads its stream. TCP sources are read-only: they work with Receive and
// NewScanner but not NewClient.
func TCP(addr string, format StreamFormat) Option {
	return optionFunc(func(c *config) {
		c.sources = append(c.sources, &tcpSource{addr: addr, format: format})
	})
}

// UDP adds a source that listens on listenAddr (e.g. ":1457" or
// "0.0.0.0:1457") for datagrams broadcast by a network gateway. UDP sources
// are read-only: they work with Receive and NewScanner but not NewClient.
func UDP(listenAddr string, format StreamFormat) Option {
	return optionFunc(func(c *config) {
		c.sources = append(c.sources, &udpSource{addr: listenAddr, format: format})
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

// WithClaimTimeout sets how long NewClient blocks waiting for address claiming
// to complete on a real CAN bus. Default is 1500ms. This allows time for the
// initial 250ms claim window plus several rounds of contention renegotiation.
func WithClaimTimeout(d time.Duration) Option {
	return optionFunc(func(c *config) {
		c.claimTimeout = &d
	})
}

// WithProductInfo sets the product identity (PGN 126996) this client reports
// when another device requests it. Without it, a generic software-gateway
// identity is reported. String fields longer than 32 bytes are truncated on
// the wire.
func WithProductInfo(p ProductInfo) Option {
	return optionFunc(func(c *config) {
		c.productInfo = &p
	})
}

// WithConfigInfo sets the installation description (PGN 126998) this client
// reports when another device requests it.
func WithConfigInfo(ci ConfigInfo) Option {
	return optionFunc(func(c *config) {
		c.configInfo = &ci
	})
}

// WithHeartbeatInterval sets the cadence of the client's automatic heartbeat
// (PGN 126993). The NMEA 2000 standard requires every device to heartbeat at
// least every 60 seconds, which is the default. Pass 0 to disable automatic
// heartbeats. Only bus clients heartbeat; replay clients never do.
func WithHeartbeatInterval(d time.Duration) Option {
	return optionFunc(func(c *config) {
		c.heartbeatInterval = &d
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

// WithBus provides a pre-constructed Bus for the client to use — either a
// custom transport not shipped by this library, or a fake for testing. When
// set, the client uses this bus directly instead of constructing one from
// CAN/USB sources.
func WithBus(bus Bus) Option {
	return optionFunc(func(c *config) {
		c.bus = bus
	})
}
