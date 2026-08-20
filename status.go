package n2k

import "github.com/open-ships/n2k/internal/gateway"

// ClientStatus is a point-in-time operational snapshot suitable for health
// endpoints and metrics collectors.
type ClientStatus struct {
	Address                       uint8
	AddressClaimed                bool
	Connected                     bool
	ConnectionEpoch               uint64
	Rejoining                     bool
	Closed                        bool
	TerminalError                 error
	WriteQueueDepth               int
	WriteQueueCapacity            int
	ReceiveSubscribers            int
	ObservationSubscribers        int
	ProtocolRequiredQueueDepth    int
	ProtocolRequiredQueueCapacity int
	ProtocolAdvisoryQueueDepth    int
	ProtocolAdvisoryQueueCapacity int
	ApplicationWritesAccepted     uint64
	ApplicationWritesCompleted    uint64
	ApplicationWritesFailed       uint64
	ApplicationWritesRejected     uint64
	ProtocolWritesAccepted        uint64
	ProtocolWritesCompleted       uint64
	ProtocolWritesFailed          uint64
	ProtocolWritesRejected        uint64
	FramesReceived                uint64
	FramesTransmitted             uint64
	MessagesObserved              uint64
	DecodeErrorsObserved          uint64
	GatewayEventsObserved         uint64
	TransportErrorsObserved       uint64
	// Actisense is populated for source-authoritative Actisense binary or CAN
	// ASCII buses. Gateway-owned message sessions expose the same counters on
	// ActisenseGatewaySession.Status instead.
	Actisense *ActisenseSessionMetrics
}

// Status returns a concurrency-safe snapshot of the client's lifecycle and
// bounded queues. It does not perform I/O.
func (c *Client) Status() ClientStatus {
	if c == nil {
		return ClientStatus{Closed: true, TerminalError: ErrClientClosed}
	}
	c.mu.Lock()
	status := ClientStatus{
		Address:                    c.sourceAddr,
		AddressClaimed:             c.claimed,
		Connected:                  c.connected,
		ConnectionEpoch:            c.connectionEpoch,
		Rejoining:                  c.rejoining,
		Closed:                     c.closed,
		TerminalError:              c.terminalErr,
		WriteQueueDepth:            len(c.writeCh),
		WriteQueueCapacity:         cap(c.writeCh),
		ApplicationWritesAccepted:  c.applicationWritesAccepted.Load(),
		ApplicationWritesCompleted: c.applicationWritesCompleted.Load(),
		ApplicationWritesFailed:    c.applicationWritesFailed.Load(),
		ApplicationWritesRejected:  c.applicationWritesRejected.Load(),
		FramesReceived:             c.framesReceived.Load(),
		FramesTransmitted:          c.framesTransmitted.Load(),
		MessagesObserved:           c.messagesObserved.Load(),
		DecodeErrorsObserved:       c.decodeErrorsObserved.Load(),
		GatewayEventsObserved:      c.gatewayEventsObserved.Load(),
		TransportErrorsObserved:    c.transportErrorsObserved.Load(),
	}
	c.mu.Unlock()
	if c.msgHub != nil {
		status.ReceiveSubscribers = c.msgHub.subscriberCount()
	}
	if c.observationHub != nil {
		status.ObservationSubscribers = c.observationHub.subscriberCount()
	}
	if c.protocolTx != nil {
		status.ProtocolRequiredQueueDepth = len(c.protocolTx.required)
		status.ProtocolRequiredQueueCapacity = cap(c.protocolTx.required)
		status.ProtocolAdvisoryQueueDepth = len(c.protocolTx.advisory)
		status.ProtocolAdvisoryQueueCapacity = cap(c.protocolTx.advisory)
		status.ProtocolWritesAccepted = c.protocolTx.accepted.Load()
		status.ProtocolWritesCompleted = c.protocolTx.completed.Load()
		status.ProtocolWritesFailed = c.protocolTx.failed.Load()
		status.ProtocolWritesRejected = c.protocolTx.rejected.Load()
	}
	if provider, ok := c.bus.(interface {
		Metrics() gateway.ActisenseGatewayMetrics
	}); ok {
		metrics := provider.Metrics()
		status.Actisense = &ActisenseSessionMetrics{
			ConnectionEpochs: metrics.ConnectionEpochs,
			Reconnects:       metrics.Reconnects,
			GatewayResets:    metrics.GatewayResets,
			Protocol:         metrics.Protocol,
		}
	}
	return status
}

// Err returns the terminal runtime error, if the bus or address-claiming
// lifecycle failed. A normal Close does not set an error.
func (c *Client) Err() error {
	if c == nil {
		return ErrClientClosed
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.terminalErr
}
