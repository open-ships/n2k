//go:build !linux

package canbus

import (
	"errors"

	"github.com/brutella/can"
)

func newSocketCANBus(string) (*can.Bus, error) {
	return nil, errors.New("SocketCAN is only available on Linux")
}
