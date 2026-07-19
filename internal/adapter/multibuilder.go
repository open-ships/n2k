package adapter

import (
	"fmt"
	"time"

	"github.com/open-ships/n2k/internal/decoder"
)

const (
	defaultFastPacketTTL      = 2 * time.Second
	defaultMaxFastPacketState = 2048
)

// MultiBuilder manages the assembly of multi-frame NMEA 2000 fast packets into complete
// single Packets. NMEA 2000 messages with payloads larger than 8 bytes are transmitted
// as a series of CAN frames that must be reassembled.
//
// MultiBuilder keys in-progress sequences by network identity, source address,
// PGN, and sequence ID. Network identity is essential when several physical
// networks are merged into one read pipeline because their dynamic source
// addresses and three-bit sequence IDs can overlap.
//
// Each unique (source, PGN, seqId) tuple maps to a sequence that accumulates frames
// until the message is complete. Once complete, the sequence is deleted to free memory.
//
// This structure is instantiated by CANAdapter and delegates frame-level assembly to
// the sequence type.
type MultiBuilder struct {
	sequences map[fastPacketKey]*sequence
	now       func() time.Time
	ttl       time.Duration
	maxActive int
}

type fastPacketKey struct {
	network string
	source  uint8
	pgn     uint32
	seqID   uint8
}

// NewMultiBuilder creates and returns a new MultiBuilder with an initialized (empty)
// sequences map, ready to receive fast-packet frames.
func NewMultiBuilder() *MultiBuilder {
	mBuilder := MultiBuilder{
		sequences: make(map[fastPacketKey]*sequence),
		now:       time.Now,
		ttl:       defaultFastPacketTTL,
		maxActive: defaultMaxFastPacketState,
	}
	return &mBuilder
}

// Add processes an incoming fast-packet frame by extracting its sequence ID and frame
// number, locating (or creating) the appropriate sequence, and adding the frame data.
// If the sequence is now complete (all expected bytes received), the assembled data is
// copied into the packet, the packet is marked Complete, and the sequence is deleted
// from the map to free memory.
//
// Parameters:
//   - p: A pointer to the Packet containing the raw CAN frame data. On return, if the
//     sequence is complete, p.Data will contain the fully assembled payload and
//     p.Complete will be true.
func (m *MultiBuilder) Add(p *decoder.Packet) {
	if p == nil || len(p.Data) == 0 {
		if p != nil {
			p.ParseErrors = append(p.ParseErrors, fmt.Errorf("fast-packet frame has no header byte"))
		}
		return
	}
	now := m.now()
	m.evictExpired(now)
	// Extract the 3-bit sequence ID and 5-bit frame number from the first data byte.
	p.GetSeqFrame()
	// Find or create the sequence for this source/PGN/seqId combination.
	seq := m.seqFor(p, now)
	// Add this frame's data to the sequence.
	if !seq.add(p, now) {
		m.deleteSequence(keyForPacket(p))
		return
	}
	// Check if the sequence now has all expected data and finalize if so.
	if seq.complete(p) {
		// Sequence is done -- delete it from the map to free memory and prevent stale data.
		m.deleteSequence(keyForPacket(p))
	}
}

// SeqFor returns the existing sequence for the given packet's (source, PGN, seqId) tuple,
// or creates a new one if it doesn't exist yet. It lazily initializes intermediate map
// levels as needed.
//
// Parameters:
//   - p: The Packet whose Info.SourceId, Info.PGN, and SeqId identify the target sequence.
//
// Returns a pointer to the sequence for this packet's source/PGN/seqId combination.
func (m *MultiBuilder) SeqFor(p *decoder.Packet) *sequence {
	return m.seqFor(p, m.now())
}

func (m *MultiBuilder) seqFor(p *decoder.Packet, now time.Time) *sequence {
	key := keyForPacket(p)
	if m.sequences[key] == nil {
		m.evictOldestIfFull()
	}
	seq := m.sequences[key]
	if seq == nil {
		seq = &sequence{}
		seq.updated = now
		m.sequences[key] = seq
	}
	return seq
}

func keyForPacket(packet *decoder.Packet) fastPacketKey {
	network := packet.Info.NetworkID
	if network == "" {
		network = packet.Info.AdapterID
	}
	return fastPacketKey{network: network, source: packet.Info.SourceId, pgn: packet.Info.PGN, seqID: packet.SeqId}
}

func (m *MultiBuilder) deleteSequence(key fastPacketKey) {
	delete(m.sequences, key)
}

func (m *MultiBuilder) evictExpired(now time.Time) {
	for key, seq := range m.sequences {
		if !seq.updated.IsZero() && now.Sub(seq.updated) >= m.ttl {
			m.deleteSequence(key)
		}
	}
}

func (m *MultiBuilder) evictOldestIfFull() {
	var oldest *sequence
	var oldestKey fastPacketKey
	for key, seq := range m.sequences {
		if oldest == nil || seq.updated.Before(oldest.updated) {
			oldest, oldestKey = seq, key
		}
	}
	if len(m.sequences) >= m.maxActive && oldest != nil {
		m.deleteSequence(oldestKey)
	}
}
