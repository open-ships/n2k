package n2k

import (
	"sort"
	"sync"
	"time"

	"github.com/open-ships/n2k/pgn"
)

// Device is one node observed on the NMEA 2000 bus. Devices are identified by
// their 64-bit ISO 11783 NAME — bus addresses are dynamic and can change
// under address contention, so RawName is the stable key.
type Device struct {
	// RawName is the packed 64-bit ISO 11783 NAME from the device's address
	// claim.
	RawName uint64
	// Name is the decoded NAME.
	Name DeviceName
	// Address is the device's most recently claimed bus address.
	Address uint8
	// LastSeen is when the device last transmitted anything.
	LastSeen time.Time
	// ProductInfo is the device's PGN 126996 payload, nil until observed.
	ProductInfo *pgn.ProductInformation
	// ConfigInfo is the device's PGN 126998 payload, nil until observed.
	ConfigInfo *pgn.ConfigurationInformation
}

// Devices returns a snapshot of every device observed on the bus, sorted by
// address. The client builds this map passively from address claims, product
// information, and configuration information; it also requests those PGNs
// from newly seen devices and enumerates the bus once at startup. Replay
// clients always return an empty list.
func (c *Client) Devices() []Device {
	if c.registry == nil {
		return nil
	}
	return c.registry.snapshot()
}

// DeviceAt resolves a bus source address to the device currently holding it.
// Use it to correlate a message's SourceId with a stable device identity:
//
//	if dev, ok := client.DeviceAt(msg.MessageInfo().SourceId); ok { ... }
func (c *Client) DeviceAt(source uint8) (Device, bool) {
	if c.registry == nil {
		return Device{}, false
	}
	return c.registry.deviceAt(source)
}

// registry tracks bus devices keyed by NAME with a secondary index by
// current address.
type registry struct {
	mu     sync.Mutex
	byName map[uint64]*deviceEntry
	byAddr map[uint8]*deviceEntry
}

type deviceEntry struct {
	rawName       uint64
	addr          uint8
	lastSeen      time.Time
	product       *pgn.ProductInformation
	config        *pgn.ConfigurationInformation
	infoRequested bool
}

func newRegistry() *registry {
	return &registry{
		byName: make(map[uint64]*deviceEntry),
		byAddr: make(map[uint8]*deviceEntry),
	}
}

// handleClaim records an address claim. It returns true the first time a
// NAME is seen, signaling the caller to request the device's product and
// configuration info.
func (r *registry) handleClaim(source uint8, rawName uint64, now time.Time) (firstSight bool) {
	// 254 is the "cannot claim" announcement; it does not hold an address.
	if source >= 254 {
		return false
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	entry, known := r.byName[rawName]
	if !known {
		entry = &deviceEntry{rawName: rawName}
		r.byName[rawName] = entry
	}

	// Rebind the address index: drop the entry's old binding and evict any
	// stale holder of the new address.
	if known && r.byAddr[entry.addr] == entry {
		delete(r.byAddr, entry.addr)
	}
	if displaced := r.byAddr[source]; displaced != nil && displaced != entry {
		// Two active NAMEs cannot own the same address. The new claim is the
		// authoritative binding; remove the stale entry so snapshots never
		// report an impossible duplicate-address topology.
		delete(r.byName, displaced.rawName)
	}
	r.byAddr[source] = entry
	entry.addr = source
	entry.lastSeen = now

	if entry.infoRequested {
		return false
	}
	entry.infoRequested = true
	return true
}

// observe attaches product and configuration info to the device currently
// holding the message's source address. Registered as a system handler.
func (r *registry) observe(msg pgn.Message) {
	carrier, ok := msg.(infoCarrier)
	if !ok {
		return
	}
	info := carrier.MessageInfo()

	r.mu.Lock()
	defer r.mu.Unlock()
	entry, ok := r.byAddr[info.SourceId]
	if !ok {
		return
	}
	switch m := msg.(type) {
	case *pgn.ProductInformation:
		entry.product = m
	case *pgn.ConfigurationInformation:
		entry.config = m
	}
	entry.lastSeen = info.Timestamp
}

// touch advances LastSeen for the device holding an address.
func (r *registry) touch(source uint8, now time.Time) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if entry, ok := r.byAddr[source]; ok {
		entry.lastSeen = now
	}
}

// snapshot copies every entry, sorted by address.
func (r *registry) snapshot() []Device {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]Device, 0, len(r.byName))
	for _, entry := range r.byName {
		out = append(out, entry.device())
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Address < out[j].Address })
	return out
}

// deviceAt resolves an address to a device copy.
func (r *registry) deviceAt(source uint8) (Device, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	entry, ok := r.byAddr[source]
	if !ok {
		return Device{}, false
	}
	return entry.device(), true
}

// device copies the entry into the public shape. Nested PGN structs are
// copied so callers cannot mutate registry state through a snapshot.
func (e *deviceEntry) device() Device {
	d := Device{
		RawName:  e.rawName,
		Name:     UnpackDeviceName(e.rawName),
		Address:  e.addr,
		LastSeen: e.lastSeen,
	}
	if e.product != nil {
		productCopy := *e.product
		productCopy.Info = cloneMessageInfo(e.product.Info)
		productCopy.Nmea2000Version = cloneUint64(e.product.Nmea2000Version)
		productCopy.ProductCode = cloneUint64(e.product.ProductCode)
		productCopy.CertificationLevel = cloneUint64(e.product.CertificationLevel)
		productCopy.LoadEquivalency = cloneUint64(e.product.LoadEquivalency)
		d.ProductInfo = &productCopy
	}
	if e.config != nil {
		configCopy := *e.config
		configCopy.Info = cloneMessageInfo(e.config.Info)
		d.ConfigInfo = &configCopy
	}
	return d
}

func cloneUint64(v *uint64) *uint64 {
	if v == nil {
		return nil
	}
	copy := *v
	return &copy
}

func cloneMessageInfo(info pgn.MessageInfo) pgn.MessageInfo {
	if info.Priority != nil {
		info.Priority = pgn.Priority(*info.Priority)
	}
	if info.TargetId != nil {
		info.TargetId = pgn.Target(*info.TargetId)
	}
	return info
}
