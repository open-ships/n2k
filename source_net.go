package n2k

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"

	"github.com/brutella/can"
	"github.com/open-ships/n2k/internal/gateway"
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
	addr   string
	format StreamFormat
}

func (s *tcpSource) run(ctx context.Context, log *slog.Logger, handler func(can.Frame)) error {
	if !s.format.valid() {
		return fmt.Errorf("n2k: unknown stream format %d", s.format)
	}

	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", s.addr)
	if err != nil {
		return fmt.Errorf("n2k: dialing %s: %w", s.addr, err)
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

	err = readStream(conn, s.format, handler)
	if ctx.Err() != nil {
		return ctx.Err()
	}
	return err
}

// readStream consumes a gateway byte stream until EOF or a read error.
func readStream(r io.Reader, format StreamFormat, handler func(can.Frame)) error {
	switch format {
	case FormatYDRaw:
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			if frame, ok := gateway.ParseYDRaw(scanner.Text()); ok {
				handler(frame)
			}
		}
		return scanner.Err()

	case FormatActisense:
		reader := gateway.NewActisenseReader()
		emit := actisenseEmitter(handler)
		buf := make([]byte, 4096)
		for {
			n, err := r.Read(buf)
			if n > 0 {
				reader.Feed(buf[:n], emit)
			}
			if err != nil {
				if errors.Is(err, io.EOF) {
					return nil
				}
				return err
			}
		}
	}
	return fmt.Errorf("n2k: unknown stream format %d", format)
}

// actisenseEmitter re-frames assembled Actisense messages into CAN frames,
// rotating the fast-packet sequence ID per message.
func actisenseEmitter(handler func(can.Frame)) func(gateway.N2KMessage) {
	var seq uint8
	return func(m gateway.N2KMessage) {
		frames, err := gateway.Reframe(m, seq)
		if err != nil {
			return
		}
		seq = (seq + 1) % 8
		for _, f := range frames {
			handler(f)
		}
	}
}

// udpSource reads gateway datagrams from a local UDP listen address.
type udpSource struct {
	addr   string
	format StreamFormat
}

func (s *udpSource) run(ctx context.Context, log *slog.Logger, handler func(can.Frame)) error {
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

	reader := gateway.NewActisenseReader()
	emit := actisenseEmitter(handler)
	buf := make([]byte, 65536)
	for {
		n, _, err := conn.ReadFrom(buf)
		if n > 0 {
			switch s.format {
			case FormatYDRaw:
				scanner := bufio.NewScanner(bytes.NewReader(buf[:n]))
				for scanner.Scan() {
					if frame, ok := gateway.ParseYDRaw(scanner.Text()); ok {
						handler(frame)
					}
				}
			case FormatActisense:
				reader.Feed(buf[:n], emit)
			}
		}
		if err != nil {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			return err
		}
	}
}
