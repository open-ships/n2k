package n2k

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"

	"github.com/open-ships/n2k/internal/gateway"
	"github.com/open-ships/n2k/raw"
)

// StreamFormat identifies the wire format spoken by a network gateway.
type StreamFormat int

const (
	// FormatYDRaw is the Yacht Devices RAW ASCII line protocol (YDWG-02 and
	// compatible gateways).
	FormatYDRaw StreamFormat = iota
	// FormatActisense is the Actisense binary stream protocol (NGT-1 and
	// compatible gateways). Messages arrive fully assembled and are re-framed
	// internally so they flow through the same decode pipeline as CAN frames.
	FormatActisense
)

func (f StreamFormat) valid() bool {
	return f == FormatYDRaw || f == FormatActisense
}

// tcpSource reads gateway traffic from a TCP connection.
type tcpSource struct {
	addr      string
	format    StreamFormat
	reconnect *gateway.ReconnectPolicy // nil = no auto-reconnect
}

func (s *tcpSource) run(ctx context.Context, log *slog.Logger, handler func(raw.Observation)) error {
	if !s.format.valid() {
		return fmt.Errorf("n2k: unknown stream format %d", s.format)
	}

	var backoff *gateway.Backoff
	if s.reconnect != nil {
		backoff = gateway.NewBackoff(*s.reconnect)
	}

	for {
		connected, err := s.dialAndRead(ctx, handler)
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if backoff == nil {
			return err
		}
		// Reset the backoff after a connection that came up, so a brief drop
		// reconnects promptly rather than inheriting a long prior backoff.
		if connected {
			backoff.Reset()
		}
		if err != nil {
			log.Warn("n2k: gateway connection lost, reconnecting", "addr", s.addr, "error", err)
		} else {
			log.Info("n2k: gateway connection closed, reconnecting", "addr", s.addr)
		}
		if !backoff.Wait(ctx) {
			return ctx.Err()
		}
	}
}

// dialAndRead opens one connection and reads it to completion. The returned
// bool reports whether the connection was established (used to reset backoff);
// the error is the dial or read error (nil on a clean EOF).
func (s *tcpSource) dialAndRead(ctx context.Context, handler func(raw.Observation)) (bool, error) {
	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", s.addr)
	if err != nil {
		return false, fmt.Errorf("n2k: dialing %s: %w", s.addr, err)
	}

	// Closing the connection on cancellation unblocks any pending Read.
	watchDone := make(chan struct{})
	defer close(watchDone)
	go func() {
		select {
		case <-ctx.Done():
			_ = conn.Close()
		case <-watchDone:
			_ = conn.Close()
		}
	}()

	return true, readStreamObservations(conn, s.format, func(observation raw.Observation) {
		observation.AdapterID = "tcp:" + s.addr
		observation.NetworkID = s.addr
		handler(observation)
	})
}

// newBus opens the gateway connection as a read/write Bus for NewClient.
// YD RAW is frame-level in both directions, so the client behaves exactly
// as on CAN hardware (the gateway echoes transmitted frames back with
// direction T). Actisense is message-level: the bus implements
// MessageWriter, and the gateway fragments fast-packet payloads and stamps
// its own claimed source address on transmissions.
func (s *tcpSource) newBus(log *slog.Logger) Bus {
	switch s.format {
	case FormatYDRaw:
		return gateway.NewYDRawTCPBus(log, s.addr, s.reconnect)
	case FormatActisense:
		return gateway.NewActisenseTCPBus(log, s.addr, s.reconnect)
	default:
		return nil
	}
}

// applyReconnect propagates the configured reconnect policy to every TCP
// source — the only source type that maintains a persistent connection. It is
// called after all options are applied, since WithReconnect and TCP may be
// supplied in any order.
func (c *config) applyReconnect() {
	if c.reconnect == nil {
		return
	}
	gp := &gateway.ReconnectPolicy{
		InitialBackoff: c.reconnect.InitialBackoff,
		MaxBackoff:     c.reconnect.MaxBackoff,
	}
	for _, src := range c.sources {
		if ts, ok := src.(*tcpSource); ok {
			ts.reconnect = gp
		}
	}
}

// readStream consumes a gateway byte stream until EOF or a read error,
// delegating to the shared per-protocol readers so the TCP bus and the
// read-only network sources decode identically.
func readStreamObservations(r io.Reader, format StreamFormat, handler func(raw.Observation)) error {
	switch format {
	case FormatYDRaw:
		return gateway.ReadYDRawObservations(r, handler)
	case FormatActisense:
		return gateway.ReadActisenseObservations(r, handler)
	}
	return fmt.Errorf("n2k: unknown stream format %d", format)
}

// udpSource reads gateway datagrams from a local UDP listen address.
type udpSource struct {
	addr   string
	format StreamFormat
}

func (s *udpSource) run(ctx context.Context, log *slog.Logger, handler func(raw.Observation)) error {
	if !s.format.valid() {
		return fmt.Errorf("n2k: unknown stream format %d", s.format)
	}

	var lc net.ListenConfig
	conn, err := lc.ListenPacket(ctx, "udp", s.addr)
	if err != nil {
		return fmt.Errorf("n2k: listening on %s: %w", s.addr, err)
	}

	watchDone := make(chan struct{})
	defer close(watchDone)
	go func() {
		select {
		case <-ctx.Done():
			_ = conn.Close()
		case <-watchDone:
			_ = conn.Close()
		}
	}()

	buf := make([]byte, 65536)
	for {
		n, remote, err := conn.ReadFrom(buf)
		if n > 0 {
			networkID := remote.String()
			emit := func(observation raw.Observation) {
				observation.AdapterID = "udp:" + s.addr
				observation.NetworkID = networkID
				handler(observation)
			}
			_ = readStreamObservations(bytes.NewReader(buf[:n]), s.format, emit)
		}
		if err != nil {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			return err
		}
	}
}
