package n2k

import (
	"time"

	"github.com/open-ships/n2k/internal/framer"
	"github.com/open-ships/n2k/pgn"
)

// restoreDefaultIntervalTicks is the reserved transmission-interval value
// (0xFFFFFFFE ms) that commands a device to revert to its default cadence.
const restoreDefaultIntervalTicks = 0xFFFFFFFE

// Heartbeat cadence bounds enforced when a group function retimes PGN 126993.
const (
	minHeartbeatInterval = time.Second
	maxHeartbeatInterval = 60 * time.Second
)

// handleGroupFunction reacts to decoded group functions (PGN 126208) from
// the system router. Request group functions for PGNs this client transmits
// are honored (including retiming); everything else is refused with an
// acknowledge group function when addressed to us specifically.
func (c *Client) handleGroupFunction(msg pgn.Message) {
	switch gf := msg.(type) {
	case *pgn.NmeaRequestGroupFunction:
		c.handleRequestGF(gf)
	case *pgn.NmeaCommandGroupFunction:
		if !c.groupFnForUs(gf.Info) {
			return
		}
		c.ackGroupFunction(gf.Pgn)
	case *pgn.NmeaReadFieldsGroupFunction:
		if !c.groupFnForUs(gf.Info) {
			return
		}
		c.ackGroupFunction(gf.Pgn)
	case *pgn.NmeaWriteFieldsGroupFunction:
		if !c.groupFnForUs(gf.Info) {
			return
		}
		c.ackGroupFunction(gf.Pgn)
	}
}

// groupFnForUs reports whether a group function is addressed to this client
// specifically. Broadcast group functions return false: they are honored for
// requests but never acknowledged, per the group-function usage rules.
func (c *Client) groupFnForUs(info pgn.MessageInfo) bool {
	if info.TargetId == nil || *info.TargetId == framer.BroadcastAddr {
		return false
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	return *info.TargetId == c.sourceAddr
}

// handleRequestGF applies a request group function: retime and/or transmit
// the named PGN when we transmit it, refuse otherwise.
func (c *Client) handleRequestGF(gf *pgn.NmeaRequestGroupFunction) {
	broadcast := gf.Info.TargetId == nil || *gf.Info.TargetId == framer.BroadcastAddr
	if !broadcast && !c.groupFnForUs(gf.Info) {
		return
	}
	if gf.Pgn == nil {
		return
	}
	requested := uint32(*gf.Pgn)

	if !c.transmitsPGN(requested) {
		if !broadcast {
			c.ackGroupFunctionError(requested, uint64(pgn.PGNNotSupported))
		}
		return
	}

	// Parameter-filtered requests are not supported.
	if len(gf.Repeating1) > 0 {
		if !broadcast {
			c.ackGroupFunctionError(requested, uint64(pgn.NotSupported_3))
		}
		return
	}

	c.applyRequestedInterval(requested, gf.TransmissionInterval)
}

// transmitsPGN reports whether this client can transmit the given PGN on
// request: its own identity PGNs plus anything with an active broadcaster.
func (c *Client) transmitsPGN(pgnNum uint32) bool {
	switch pgnNum {
	case 126993, 126996, 126998:
		return true
	}
	return c.broadcasterFor(pgnNum) != nil
}

// applyRequestedInterval retimes and/or transmits a PGN we support.
// interval is in wire ticks (milliseconds): nil = no cadence change (just
// transmit once), 0 = stop periodic transmission, 0xFFFFFFFE = restore the
// default cadence, anything else = the new cadence. Cadence only applies to
// PGNs with a schedule (the heartbeat and active broadcasts); product and
// configuration info have none, so a request simply transmits them.
func (c *Client) applyRequestedInterval(pgnNum uint32, intervalTicks *uint64) {
	switch pgnNum {
	case 126996:
		c.Write(c.productInfo.message())
		return
	case 126998:
		c.Write(c.configInfo.message())
		return
	case 126993:
		if intervalTicks == nil {
			c.heartbeat.sendNow()
			return
		}
		switch *intervalTicks {
		case 0:
			c.heartbeat.setInterval(0)
		case restoreDefaultIntervalTicks:
			// setInterval transmits immediately, covering the request.
			c.heartbeat.setInterval(defaultHeartbeatInterval)
		default:
			c.heartbeat.setInterval(clampHeartbeatInterval(ticksToDuration(*intervalTicks)))
		}
		return
	}

	b := c.broadcasterFor(pgnNum)
	if b == nil {
		return // schedule raced away since the transmitsPGN check
	}
	if intervalTicks == nil {
		b.sendNow()
		return
	}
	switch *intervalTicks {
	case 0:
		b.setInterval(0)
	case restoreDefaultIntervalTicks:
		b.restoreDefaultInterval()
	default:
		b.setInterval(ticksToDuration(*intervalTicks))
	}
}

func ticksToDuration(ms uint64) time.Duration {
	return time.Duration(ms) * time.Millisecond
}

func clampHeartbeatInterval(d time.Duration) time.Duration {
	if d < minHeartbeatInterval {
		return minHeartbeatInterval
	}
	if d > maxHeartbeatInterval {
		return maxHeartbeatInterval
	}
	return d
}

// ackGroupFunction refuses a non-request group function for pgnPtr.
func (c *Client) ackGroupFunction(pgnPtr *uint64) {
	if pgnPtr == nil {
		return
	}
	c.ackGroupFunctionError(uint32(*pgnPtr), uint64(pgn.NotSupported_3))
}

// ackGroupFunctionError broadcasts an acknowledge group function refusing the
// named PGN with the given PGN error code.
func (c *Client) ackGroupFunctionError(pgnNum uint32, errorCode uint64) {
	fc := uint64(2) // acknowledge
	refused := uint64(pgnNum)
	intervalErr := uint64(0) // acknowledge (no interval/priority complaint)
	params := uint64(0)
	c.Write(&pgn.NmeaAcknowledgeGroupFunction{
		Info:                                  pgn.MessageInfo{Priority: pgn.Priority(3)},
		FunctionCode:                          &fc,
		Pgn:                                   &refused,
		PgnErrorCode:                          &errorCode,
		TransmissionIntervalPriorityErrorCode: &intervalErr,
		NumberOfParameters:                    &params,
	})
}
