//go:build !linux && !darwin && !freebsd && !openbsd

// Package serialio opens serial ports with interruptible reads and writes.
package serialio

import "go.bug.st/serial"

// Open uses the platform serial implementation outside Unix.
func Open(path string, mode *serial.Mode) (serial.Port, error) {
	return serial.Open(path, mode)
}
