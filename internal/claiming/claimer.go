// Package claiming implements the NMEA 2000 / ISO 11783 address claiming
// protocol (PGN 60928). It supports two modes: auto mode, where the claimer
// negotiates downward from a starting address on contention, and explicit mode,
// where contention is treated as a fatal error.
package claiming

import (
	"encoding/binary"
	"fmt"
	"log/slog"
	"sync"

	"github.com/brutella/can"
	"github.com/open-ships/n2k/internal/framer"
)

// pgnAddressClaim is the PGN number for ISO Address Claim (60928 = 0xEE00).
const pgnAddressClaim = 60928

// addressClaimPriority is the CAN priority used for address claim messages.
const addressClaimPriority = 6

// globalDestination is the broadcast destination address (255).
const globalDestination = 255

// cannotClaimAddress is the special address indicating the device cannot claim
// any address on the network.
const cannotClaimAddress = 254

// minAddress is the lowest valid claimable address on an NMEA 2000 network.
const minAddress = 1

// maxAddress is the highest valid claimable address on an NMEA 2000 network.
const maxAddress = 253

// Mode controls how the claimer responds to address contention.
type Mode int

const (
	// ModeAuto starts at the configured address (default 253) and works
	// downward on contention until all addresses are exhausted, at which
	// point it falls back to address 254 (cannot claim).
	ModeAuto Mode = iota

	// ModeExplicit uses a fixed address and surfaces a fatal error if
	// another device with a lower NAME claims the same address.
	ModeExplicit
)

// Config holds the parameters needed to construct a Claimer.
type Config struct {
	// Mode selects auto or explicit address claiming behaviour.
	Mode Mode

	// Address is the starting address in auto mode or the fixed address in
	// explicit mode. Valid range is 1-253.
	Address uint8

	// Name is our device's 64-bit ISO 11783 NAME used for arbitration.
	// Lower NAME values win contention.
	Name uint64

	// WriteFrame sends a CAN frame onto the bus.
	WriteFrame func(can.Frame) error

	// OnAddressChange is called when the claimed address changes (auto mode).
	// It may be nil if the caller does not need notifications.
	OnAddressChange func(newAddr uint8)

	// OnFatalError is called when an unrecoverable error occurs, such as
	// address contention in explicit mode. It may be nil.
	OnFatalError func(err error)

	// Logger is used for structured logging of claiming activity.
	// If nil, a no-op logger is used.
	Logger *slog.Logger
}

// Claimer implements the NMEA 2000 address claiming state machine.
type Claimer struct {
	cfg     Config
	address uint8
	tried   map[uint8]bool
	mu      sync.Mutex
}

// New creates a new Claimer from the given configuration.
func New(cfg Config) *Claimer {
	if cfg.Logger == nil {
		cfg.Logger = slog.Default()
	}
	return &Claimer{
		cfg:     cfg,
		address: cfg.Address,
		tried:   make(map[uint8]bool),
	}
}

// Start sends the initial address claim broadcast. It returns an error if the
// claim frame cannot be sent.
func (c *Claimer) Start() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cfg.Logger.Info("joining network", "address", c.address)
	c.tried[c.address] = true
	return c.sendClaim()
}

// HandleAddressClaim processes an incoming Address Claim (PGN 60928) from
// another device. source is the CAN source address of the competing device and
// name is its 64-bit NAME.
//
// If the competing device claims our address:
//   - If our NAME is lower (we win), we re-broadcast our claim to defend.
//   - If their NAME is lower (we lose), auto mode tries the next available
//     address downward; explicit mode surfaces a fatal error.
//
// If the competing device claims a different address, the message is ignored.
func (c *Claimer) HandleAddressClaim(source uint8, name uint64) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Only relevant if someone else is claiming our address.
	if source != c.address {
		return
	}

	// Ignore our own claims (same NAME).
	if name == c.cfg.Name {
		return
	}

	if c.cfg.Name < name {
		// We win — re-broadcast our claim to defend the address.
		c.cfg.Logger.Info("defended address against contender",
			"address", c.address,
			"contenderName", fmt.Sprintf("0x%016X", name),
		)
		if err := c.sendClaim(); err != nil {
			c.cfg.Logger.Error("failed to send defensive claim", "error", err)
		}
		return
	}

	// We lose — the other device has a lower (winning) NAME.
	switch c.cfg.Mode {
	case ModeAuto:
		c.cfg.Logger.Info("address contested, retrying",
			"address", c.address,
			"contenderName", fmt.Sprintf("0x%016X", name),
		)
		c.tryNextAddress()

	case ModeExplicit:
		err := fmt.Errorf("address %d contested in explicit mode by device 0x%016X", c.address, name)
		c.cfg.Logger.Error("FATAL: address contested in explicit mode",
			"address", c.address,
			"contenderName", fmt.Sprintf("0x%016X", name),
		)
		if c.cfg.OnFatalError != nil {
			c.cfg.OnFatalError(err)
		}
	}
}

// HandleISORequest processes an incoming ISO Request (PGN 59904) asking for
// PGN 60928. It responds with our current address claim.
func (c *Claimer) HandleISORequest() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := c.sendClaim(); err != nil {
		c.cfg.Logger.Error("failed to respond to ISO request", "error", err)
	}
}

// Address returns the currently claimed address. This value may change over
// time in auto mode as contention is resolved.
func (c *Claimer) Address() uint8 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.address
}

// sendClaim broadcasts an Address Claim frame for our current address.
// The caller must hold c.mu.
func (c *Claimer) sendClaim() error {
	payload := nameToPayload(c.cfg.Name)
	canID := framer.BuildCANID(pgnAddressClaim, addressClaimPriority, c.address, globalDestination)
	frame := framer.FrameSingle(canID, payload)
	return c.cfg.WriteFrame(frame)
}

// tryNextAddress searches downward from the current address for an untried
// address. If all addresses (1-253) have been tried, it falls back to 254.
// The caller must hold c.mu.
func (c *Claimer) tryNextAddress() {
	oldAddr := c.address

	// Scan downward from current address-1, wrapping around from minAddress
	// back to maxAddress, covering all 253 valid addresses.
	for candidate := oldAddr - 1; ; candidate-- {
		if candidate < minAddress {
			candidate = maxAddress
		}
		if candidate == oldAddr {
			// We have wrapped all the way around — no addresses available.
			break
		}
		if !c.tried[candidate] {
			c.address = candidate
			c.tried[candidate] = true
			c.cfg.Logger.Info("address contested, retrying with new address",
				"oldAddress", oldAddr,
				"newAddress", c.address,
			)
			if c.cfg.OnAddressChange != nil {
				c.cfg.OnAddressChange(c.address)
			}
			if err := c.sendClaim(); err != nil {
				c.cfg.Logger.Error("failed to send claim for new address", "error", err)
			}
			return
		}
	}

	// All addresses exhausted.
	c.address = cannotClaimAddress
	c.cfg.Logger.Warn("all addresses exhausted, falling back to 254")
	if c.cfg.OnAddressChange != nil {
		c.cfg.OnAddressChange(c.address)
	}
	// Send a "cannot claim" announcement from address 254.
	if err := c.sendClaim(); err != nil {
		c.cfg.Logger.Error("failed to send cannot-claim announcement", "error", err)
	}
}

// nameToPayload encodes a 64-bit NAME as an 8-byte little-endian payload
// suitable for an Address Claim CAN frame. The bit layout matches the
// ISO 11783 NAME specification and the wire format produced by
// pgn.EncodeIsoAddressClaim.
func nameToPayload(name uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, name)
	return b
}
