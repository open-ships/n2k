//go:build !linux

package canbus

import (
	"errors"
	"io"
)

func openSocketCAN(string) (io.ReadWriteCloser, error) {
	return nil, errors.New("SocketCAN is only available on Linux")
}
