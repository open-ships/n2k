package actisense

import (
	"sync/atomic"
	"time"
)

// SessionMetrics is a point-in-time snapshot of one Actisense protocol
// session. BSTFrames is copied and safe for the caller to mutate.
type SessionMetrics struct {
	TransportReadCalls   uint64
	TransportReadBytes   uint64
	TransportReadErrors  uint64
	TransportWriteCalls  uint64
	TransportWriteBytes  uint64
	TransportWriteErrors uint64

	Datagrams      uint64
	UnframedBytes  uint64
	FramingErrors  uint64
	ChecksumErrors uint64
	LengthErrors   uint64
	OversizeErrors uint64
	BSTFrames      map[byte]uint64

	BEMRequests               uint64
	BEMResponses              uint64
	BEMCompleted              uint64
	BEMCorrelationMisses      uint64
	BEMDuplicateRequests      uint64
	BEMTimeouts               uint64
	BEMDeviceErrors           uint64
	BEMNegativeAcks           uint64
	BEMResponseTrainOverflows uint64
	BEMInFlight               uint64
	BEMMaxInFlight            uint64
	BEMLatencyMinimum         time.Duration
	BEMLatencyAverage         time.Duration
	BEMLatencyMaximum         time.Duration
}

type sessionMetrics struct {
	transportReadCalls   atomic.Uint64
	transportReadBytes   atomic.Uint64
	transportReadErrors  atomic.Uint64
	transportWriteCalls  atomic.Uint64
	transportWriteBytes  atomic.Uint64
	transportWriteErrors atomic.Uint64

	datagrams      atomic.Uint64
	unframedBytes  atomic.Uint64
	framingErrors  atomic.Uint64
	checksumErrors atomic.Uint64
	lengthErrors   atomic.Uint64
	oversizeErrors atomic.Uint64
	bstFrames      [256]atomic.Uint64

	bemRequests               atomic.Uint64
	bemResponses              atomic.Uint64
	bemCompleted              atomic.Uint64
	bemCorrelationMisses      atomic.Uint64
	bemDuplicateRequests      atomic.Uint64
	bemTimeouts               atomic.Uint64
	bemDeviceErrors           atomic.Uint64
	bemNegativeAcks           atomic.Uint64
	bemResponseTrainOverflows atomic.Uint64
	bemInFlight               atomic.Uint64
	bemMaxInFlight            atomic.Uint64
	bemLatencyTotal           atomic.Uint64
	bemLatencyMinimum         atomic.Uint64
	bemLatencyMaximum         atomic.Uint64
}

func (m *sessionMetrics) observeLatency(elapsed time.Duration) {
	nanos := uint64(max(elapsed.Nanoseconds(), 0))
	m.bemLatencyTotal.Add(nanos)
	for current := m.bemLatencyMinimum.Load(); current == 0 || nanos < current; current = m.bemLatencyMinimum.Load() {
		if m.bemLatencyMinimum.CompareAndSwap(current, nanos) {
			break
		}
	}
	for current := m.bemLatencyMaximum.Load(); nanos > current; current = m.bemLatencyMaximum.Load() {
		if m.bemLatencyMaximum.CompareAndSwap(current, nanos) {
			break
		}
	}
}

func (m *sessionMetrics) incrementInFlight() {
	current := m.bemInFlight.Add(1)
	for maximum := m.bemMaxInFlight.Load(); current > maximum; maximum = m.bemMaxInFlight.Load() {
		if m.bemMaxInFlight.CompareAndSwap(maximum, current) {
			break
		}
	}
}

func (m *sessionMetrics) snapshot() SessionMetrics {
	completed := m.bemCompleted.Load()
	totalLatency := m.bemLatencyTotal.Load()
	result := SessionMetrics{
		TransportReadCalls: m.transportReadCalls.Load(), TransportReadBytes: m.transportReadBytes.Load(), TransportReadErrors: m.transportReadErrors.Load(),
		TransportWriteCalls: m.transportWriteCalls.Load(), TransportWriteBytes: m.transportWriteBytes.Load(), TransportWriteErrors: m.transportWriteErrors.Load(),
		Datagrams: m.datagrams.Load(), UnframedBytes: m.unframedBytes.Load(), FramingErrors: m.framingErrors.Load(),
		ChecksumErrors: m.checksumErrors.Load(), LengthErrors: m.lengthErrors.Load(), OversizeErrors: m.oversizeErrors.Load(),
		BEMRequests: m.bemRequests.Load(), BEMResponses: m.bemResponses.Load(), BEMCompleted: completed,
		BEMCorrelationMisses: m.bemCorrelationMisses.Load(), BEMDuplicateRequests: m.bemDuplicateRequests.Load(),
		BEMTimeouts: m.bemTimeouts.Load(), BEMDeviceErrors: m.bemDeviceErrors.Load(), BEMNegativeAcks: m.bemNegativeAcks.Load(),
		BEMResponseTrainOverflows: m.bemResponseTrainOverflows.Load(), BEMInFlight: m.bemInFlight.Load(), BEMMaxInFlight: m.bemMaxInFlight.Load(),
		BEMLatencyMinimum: time.Duration(m.bemLatencyMinimum.Load()), BEMLatencyMaximum: time.Duration(m.bemLatencyMaximum.Load()),
		BSTFrames: make(map[byte]uint64),
	}
	if completed != 0 {
		result.BEMLatencyAverage = time.Duration(totalLatency / completed)
	}
	for id := range m.bstFrames {
		if count := m.bstFrames[id].Load(); count != 0 {
			result.BSTFrames[byte(id)] = count
		}
	}
	return result
}
