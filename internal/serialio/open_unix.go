//go:build linux || darwin || freebsd || openbsd

// Package serialio opens serial ports with interruptible reads and writes.
package serialio

import (
	"errors"
	"os"
	"sync"
	"syscall"
	"time"

	"go.bug.st/serial"
)

// Open configures a serial port and gives its writes to Go's poller. The serial
// library already provides interruptible reads, but its Unix Write uses a
// blocking syscall that closing another goroutine's descriptor need not stop.
// A separate nonblocking descriptor makes pending writes interruptible by Close.
// It is opened before terminal exclusivity is claimed; both descriptors share
// the same terminal configuration.
func Open(path string, mode *serial.Mode) (serial.Port, error) {
	writer, err := os.OpenFile(path, os.O_RDWR|syscall.O_NONBLOCK|syscall.O_NOCTTY, 0) // #nosec G304 -- path is the explicitly selected serial device; no create or truncate flags, and serial.Open validates the terminal before writes.
	if err != nil {
		return nil, err
	}
	if err := writer.SetWriteDeadline(time.Time{}); err != nil {
		_ = writer.Close()
		return nil, err
	}
	libraryPort, err := serial.Open(path, mode)
	if err != nil {
		_ = writer.Close()
		return nil, err
	}
	return &port{Port: libraryPort, writer: writer}, nil
}

type port struct {
	serial.Port
	writer    *os.File
	closeOnce sync.Once
	closeErr  error
}

var _ serial.Port = (*port)(nil)

func (p *port) Write(data []byte) (int, error) { return p.writer.Write(data) }

func (p *port) Close() error {
	p.closeOnce.Do(func() {
		p.closeErr = errors.Join(p.writer.Close(), p.Port.Close())
	})
	return p.closeErr
}
