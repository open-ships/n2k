package n2k

import (
	"github.com/brutella/can"
	"github.com/open-ships/n2k/internal/framer"
	"github.com/open-ships/n2k/pgn"
)

// handleISORequest reacts to an incoming ISO Request (PGN 59904). Requests
// addressed to another node are ignored. Requests we can serve (address
// claim, product info, configuration info, heartbeat) are answered; anything
// else gets a NAK when the request was addressed to us specifically —
// broadcast requests for unsupported PGNs are ignored per ISO 11783-3.
func (c *Client) handleISORequest(info pgn.MessageInfo, frame can.Frame) {
	if frame.Length < 3 {
		return
	}
	requested := uint32(frame.Data[0]) | uint32(frame.Data[1])<<8 | uint32(frame.Data[2])<<16

	c.mu.Lock()
	ourAddr := c.sourceAddr
	c.mu.Unlock()

	broadcast := info.TargetId == nil || *info.TargetId == framer.BroadcastAddr
	if !broadcast && *info.TargetId != ourAddr {
		return
	}

	switch requested {
	case framer.PGNISOAddressClaim:
		c.claimer.HandleISORequest()
	case 126996:
		c.Write(c.productInfo.message())
	case 126998:
		c.Write(c.configInfo.message())
	case 126993:
		if c.heartbeat != nil {
			c.heartbeat.sendNow()
		}
	default:
		if !broadcast {
			c.Write(nakFor(requested))
		}
	}
}

// nakFor builds the ISO Acknowledgement (PGN 59392) refusing a request for
// an unsupported PGN. Acknowledgements are broadcast per ISO 11783-3; the
// requester is identifiable from the refused PGN.
func nakFor(requested uint32) *pgn.IsoAcknowledgement {
	control := uint64(pgn.NAK)
	groupFunction := uint64(0xFF) // not a group-function response
	refused := uint64(requested)
	return &pgn.IsoAcknowledgement{
		Control:       &control,
		GroupFunction: &groupFunction,
		Pgn:           &refused,
	}
}
