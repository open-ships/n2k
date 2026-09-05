package canbus

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/brutella/can"
	"golang.org/x/sys/unix"
)

// Own the descriptor setup: a blocking AF_CAN descriptor cannot reliably be
// interrupted by Close from another goroutine. Nonblocking mode enrolls the
// socket in Go's poller, giving reads and writes a real cancellation boundary.
func newSocketCANBus(name string) (*can.Bus, error) {
	iface, err := net.InterfaceByName(name)
	if err != nil {
		return nil, err
	}
	fd, err := unix.Socket(unix.AF_CAN, unix.SOCK_RAW|unix.SOCK_NONBLOCK|unix.SOCK_CLOEXEC, unix.CAN_RAW)
	if err != nil {
		return nil, fmt.Errorf("open CAN socket: %w", err)
	}
	if err := unix.Bind(fd, &unix.SockaddrCAN{Ifindex: iface.Index}); err != nil {
		_ = unix.Close(fd)
		return nil, fmt.Errorf("bind CAN socket: %w", err)
	}
	file := os.NewFile(uintptr(fd), "socketcan:"+name)
	if file == nil {
		_ = unix.Close(fd)
		return nil, fmt.Errorf("open CAN socket file")
	}
	if err := file.SetDeadline(time.Time{}); err != nil {
		_ = file.Close()
		return nil, fmt.Errorf("CAN socket does not support cancellation: %w", err)
	}
	return can.NewBus(can.NewReadWriteCloser(file)), nil
}
