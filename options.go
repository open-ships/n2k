package n2k

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/brutella/can"
	"github.com/open-ships/n2k/raw"
)

type config struct {
	sources           []source
	filterExpr        string
	includeUnknown    bool
	logger            *slog.Logger
	sourceAddress     *uint8           // nil = auto mode
	preferredAddress  *uint8           // nil = default auto-mode starting address
	deviceName        *DeviceName      // nil = use default
	claimTimeout      *time.Duration   // nil = use default (1500ms)
	readyTimeout      *time.Duration   // nil = use default (5s)
	heartbeatInterval *time.Duration   // nil = default (60s), 0 = disabled
	receiveBuffer     *int             // nil = default (64)
	writeQueue        *int             // nil = default (64)
	writeTimeout      *time.Duration   // maximum duration of one physical write
	productInfo       *ProductInfo     // nil = defaults
	configInfo        *ConfigInfo      // nil = defaults
	bus               Bus              // pre-constructed bus
	reconnect         *ReconnectPolicy // nil = no auto-reconnect
}

func (c *config) validate() error {
	if len(c.sources) == 0 && c.bus == nil {
		return errors.New("n2k: at least one source or WithBus is required")
	}
	if c.sourceAddress != nil && *c.sourceAddress > 251 {
		return fmt.Errorf("n2k: source address %d is outside 0-251", *c.sourceAddress)
	}
	if c.preferredAddress != nil && *c.preferredAddress > 251 {
		return fmt.Errorf("n2k: preferred address %d is outside 0-251", *c.preferredAddress)
	}
	if c.sourceAddress != nil && c.preferredAddress != nil {
		return errors.New("n2k: WithSourceAddress and WithPreferredAddress are mutually exclusive")
	}
	if c.claimTimeout != nil && *c.claimTimeout <= 0 {
		return errors.New("n2k: claim timeout must be positive")
	}
	if c.readyTimeout != nil && *c.readyTimeout <= 0 {
		return errors.New("n2k: ready timeout must be positive")
	}
	if c.heartbeatInterval != nil && *c.heartbeatInterval < 0 {
		return errors.New("n2k: heartbeat interval cannot be negative")
	}
	if c.receiveBuffer != nil && *c.receiveBuffer <= 0 {
		return errors.New("n2k: receive buffer must be positive")
	}
	if c.writeQueue != nil && *c.writeQueue <= 0 {
		return errors.New("n2k: write queue must be positive")
	}
	if c.writeTimeout != nil && *c.writeTimeout <= 0 {
		return errors.New("n2k: write timeout must be positive")
	}
	hasTCP := false
	for _, src := range c.sources {
		switch s := src.(type) {
		case *socketCANSource:
			if s.iface == "" {
				return errors.New("n2k: CAN interface name cannot be empty")
			}
		case *usbCANSource:
			if s.port == "" {
				return errors.New("n2k: USB serial port cannot be empty")
			}
		case *serialSource:
			if s.configErr != nil {
				return s.configErr
			}
			if s.port == "" || (s.format != FormatActisense && s.format != FormatActisenseRaw &&
				s.format != FormatActisenseCANASCII && s.format != FormatActisenseN2KASCII) {
				return fmt.Errorf("n2k: invalid serial source port or Actisense stream format %d", s.format)
			}
			if err := s.serialConfig.validate(); err != nil {
				return err
			}
		case *actisenseSerialSource:
			if s.configErr != nil {
				return s.configErr
			}
			if s.port == "" {
				return errors.New("n2k: Actisense serial port cannot be empty")
			}
			if err := s.serialConfig.validate(); err != nil {
				return err
			}
		case *tcpSource:
			hasTCP = true
			if s.addr == "" || !s.format.valid() {
				return fmt.Errorf("n2k: invalid TCP source address or stream format %d", s.format)
			}
		case *actisenseTCPSource:
			hasTCP = true
			if s.addr == "" {
				return errors.New("n2k: Actisense TCP address cannot be empty")
			}
		case *udpSource:
			if s.addr == "" || !s.format.valid() {
				return fmt.Errorf("n2k: invalid UDP source address or stream format %d", s.format)
			}
		case *fileSource:
			if s.path == "" {
				return errors.New("n2k: file path cannot be empty")
			}
		case *eblSource:
			if s.path == "" {
				return errors.New("n2k: EBL path cannot be empty")
			}
		}
	}
	if c.reconnect != nil && !hasTCP {
		return errors.New("n2k: WithReconnect requires a TCP source")
	}
	return nil
}

// Option configures Receive, Observe, NewScanner, or NewClient.
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

// Serial adds a directly connected Actisense-format gateway at 115200 8N1.
// FormatActisense and FormatActisenseN2KASCII are passive, gateway-owned
// message sources. FormatActisenseRaw and FormatActisenseCANASCII perform an
// acknowledged mode setup and provide source-authoritative CAN frames.
func Serial(port string, format StreamFormat, options ...ActisenseSerialOption) Option {
	return optionFunc(func(c *config) {
		serialConfig, err := applyActisenseSerialOptions(options)
		c.sources = append(c.sources, &serialSource{port: port, format: format, serialConfig: serialConfig, configErr: err})
	})
}

// ActisenseSerial adds a directly connected Actisense gateway at 115200 8N1.
// Receive, Observe, and NewScanner use the compatible gateway-owned message
// session. NewClient requires acknowledged mode 5 and uses source-authoritative
// BST-95 raw CAN; it fails rather than falling back when raw mode is
// unavailable. Use Serial with an explicit Actisense format to override this
// role-based policy.
func ActisenseSerial(port string, options ...ActisenseSerialOption) Option {
	return optionFunc(func(c *config) {
		serialConfig, err := applyActisenseSerialOptions(options)
		c.sources = append(c.sources, &actisenseSerialSource{port: port, serialConfig: serialConfig, configErr: err})
	})
}

// File adds a source that replays CAN frames from a candump -L / -l log file.
// By default frames are delivered as fast as they can be read; pass
// OriginalTiming() to pace them by the log's timestamps. File sources are
// read-only: they work with Receive and NewScanner but not NewClient.
func File(path string, opts ...FileOption) Option {
	return optionFunc(func(c *config) {
		capture := captureOptions{}
		for _, o := range opts {
			if o != nil {
				o.applyCapture(&capture)
			}
		}
		src := &fileSource{path: path, originalTiming: capture.originalTiming}
		c.sources = append(c.sources, src)
	})
}

// TCP adds a source that dials a network gateway (e.g. a Yacht Devices
// YDWG-02 in RAW server mode, or an Actisense gateway) at addr ("host:port").
// TCP works with Receive/NewScanner and can also back NewClient for writes.
// FormatYDRaw, FormatActisenseRaw, and FormatActisenseCANASCII provide
// frame-level access. The two Actisense message formats are read-only here;
// use NewActisenseTCPSession for a writable gateway-owned binary session.
func TCP(addr string, format StreamFormat) Option {
	return optionFunc(func(c *config) {
		c.sources = append(c.sources, &tcpSource{addr: addr, format: format})
	})
}

// YachtDevicesTCP adds a Yacht Devices RAW TCP gateway (for example, a
// YDWG-02 in RAW server mode). Compatible gateways that speak the same
// protocol are also supported. It provides source-authoritative CAN frames
// and can back NewClient.
func YachtDevicesTCP(addr string) Option {
	return TCP(addr, FormatYDRaw)
}

// ActisenseTCP adds an Actisense gateway over TCP. Receive, Observe, and
// NewScanner passively decode all supported BST records without changing the
// gateway's operating mode. NewClient requires acknowledged mode 5 and uses
// source-authoritative BST-95 raw CAN; it fails rather than falling back when
// raw mode is unavailable. Use TCP with an explicit Actisense format to
// override this role-based policy.
func ActisenseTCP(addr string) Option {
	return optionFunc(func(c *config) {
		c.sources = append(c.sources, &actisenseTCPSource{addr: addr})
	})
}

// UDP adds a source that listens on listenAddr (e.g. ":1457" or
// "0.0.0.0:1457") for datagrams broadcast by a network gateway. UDP sources
// are read-only: they work with Receive and NewScanner but not NewClient.
// Raw/ASCII formats require the upstream gateway to already emit the selected
// representation because UDP has no return channel for BEM mode setup.
func UDP(listenAddr string, format StreamFormat) Option {
	return optionFunc(func(c *config) {
		c.sources = append(c.sources, &udpSource{addr: listenAddr, format: format})
	})
}

// YachtDevicesUDP adds a read-only Yacht Devices RAW UDP broadcast source.
// Compatible gateways that speak the same protocol are also supported.
func YachtDevicesUDP(listenAddr string) Option {
	return UDP(listenAddr, FormatYDRaw)
}

// Replay adds a source that replays the given CAN frames. Useful for testing.
func Replay(frames []can.Frame) Option {
	return optionFunc(func(c *config) {
		observations := make([]raw.Observation, 0, len(frames))
		for _, frame := range frames {
			observations = append(observations, frameObservation(frame, "replay", "replay", raw.DirectionReceived))
		}
		c.sources = append(c.sources, &replaySource{observations: observations})
	})
}

// ReplayObservations adds owned source-aware observations for deterministic
// tests and capture replay. Each observation is copied before it is retained.
func ReplayObservations(observations []Observation) Option {
	return optionFunc(func(c *config) {
		copied := make([]raw.Observation, len(observations))
		for i := range observations {
			copied[i] = observations[i].Clone()
		}
		c.sources = append(c.sources, &replaySource{observations: copied})
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
// When not set (default), the client uses auto mode — starting at address 251
// and working downward if contention occurs.
func WithSourceAddress(addr uint8) Option {
	return optionFunc(func(c *config) {
		c.sourceAddress = &addr
	})
}

// WithPreferredAddress sets the starting address for automatic address
// claiming while retaining arbitrary-address capability. Persist the last
// Client.Status().Address and pass it here on the next start to reclaim the
// device's prior address when available. Valid addresses are 0 through 251.
// Unlike WithSourceAddress, contention moves the client to another address.
func WithPreferredAddress(addr uint8) Option {
	return optionFunc(func(c *config) {
		c.preferredAddress = &addr
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

// WithReadyTimeout sets how long NewClient waits for an asynchronous Bus to
// open and finish transport-specific readiness work before address claiming
// starts. The default is five seconds. This deadline is intentionally separate
// from WithClaimTimeout so a gateway handshake can return its typed failure
// instead of being hidden by the address-claim deadline.
func WithReadyTimeout(d time.Duration) Option {
	return optionFunc(func(c *config) {
		c.readyTimeout = &d
	})
}

// WithProductInfo sets the product identity (PGN 126996) this client reports
// when another device requests it. Without it, a generic software-gateway
// identity is reported. String fields longer than 32 bytes are rejected when
// the client is created.
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

// WithReceiveBuffer sets the number of decoded messages retained per live
// Client subscription. The default is 64. A subscriber that falls behind
// this bound is closed with ErrReceiveOverflow so it cannot stall protocol
// processing or other subscribers.
func WithReceiveBuffer(size int) Option {
	return optionFunc(func(c *config) {
		c.receiveBuffer = &size
	})
}

// WithWriteQueue sets the number of asynchronous writes that can wait behind
// the active write. The default is 64. Once full, Write completes immediately
// with ErrWriteQueueFull rather than blocking an application goroutine.
func WithWriteQueue(size int) Option {
	return optionFunc(func(c *config) {
		c.writeQueue = &size
	})
}

// WithWriteTimeout bounds one physical frame or gateway record write. The
// default is one second. A stalled legacy Bus is closed to interrupt I/O.
func WithWriteTimeout(timeout time.Duration) Option {
	return optionFunc(func(c *config) { c.writeTimeout = &timeout })
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

// ReconnectPolicy configures automatic reconnection for network gateway
// (TCP) sources after a dropped connection. The zero value is valid: both
// fields fall back to their defaults.
type ReconnectPolicy struct {
	// InitialBackoff is the delay before the first reconnect attempt after a
	// drop, and the delay restored after any successful connection. Defaults
	// to 500ms when zero.
	InitialBackoff time.Duration
	// MaxBackoff caps the exponentially growing delay between attempts while
	// the gateway stays unreachable. Defaults to 30s when zero.
	MaxBackoff time.Duration
}

// WithReconnect enables automatic reconnection for TCP gateway sources. After
// a connection drops, the source re-dials with exponential backoff (starting
// at InitialBackoff, capped at MaxBackoff) until it reconnects or the context
// is cancelled. Without this option, a dropped TCP connection ends the read
// loop and surfaces as an error (the historical behavior).
//
// Reconnection covers connections that drop mid-session; the initial
// connection must still succeed (for NewClient, within the claim timeout).
// It applies only to TCP sources — CAN, USB, UDP, file, and replay sources are
// unaffected. A bus client starts a new network epoch after reconnect: it
// reclaims its address, waits through contention, refreshes device discovery,
// and then resumes scheduled transmissions.
func WithReconnect(policy ReconnectPolicy) Option {
	// The backoff bounds are normalized once, at the point of use, by
	// gateway.NewBackoff (defaults for non-positive fields, MaxBackoff clamped
	// up to InitialBackoff); the policy is stored here as given.
	return optionFunc(func(c *config) {
		c.reconnect = &policy
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
