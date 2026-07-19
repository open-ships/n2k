package n2k

// ClientStatus is a point-in-time operational snapshot suitable for health
// endpoints and metrics collectors.
type ClientStatus struct {
	Address            uint8
	AddressClaimed     bool
	Closed             bool
	TerminalError      error
	WriteQueueDepth    int
	WriteQueueCapacity int
	ReceiveSubscribers int
}

// Status returns a concurrency-safe snapshot of the client's lifecycle and
// bounded queues. It does not perform I/O.
func (c *Client) Status() ClientStatus {
	if c == nil {
		return ClientStatus{Closed: true, TerminalError: ErrClientClosed}
	}
	c.mu.Lock()
	status := ClientStatus{
		Address:            c.sourceAddr,
		AddressClaimed:     c.claimed,
		Closed:             c.closed,
		TerminalError:      c.terminalErr,
		WriteQueueDepth:    len(c.writeCh),
		WriteQueueCapacity: cap(c.writeCh),
	}
	c.mu.Unlock()
	if c.msgHub != nil {
		status.ReceiveSubscribers = c.msgHub.subscriberCount()
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
