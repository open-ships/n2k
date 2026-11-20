package canbus

import (
	"context"
	"io"
	"log/slog"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.bug.st/serial"
)

func TestSocketCANRunStopsOnContextCancellation(t *testing.T) {
	local, peer := net.Pipe()
	defer func() { _ = peer.Close() }()
	ch := NewSocketCAN(testLogger(), "test")
	ch.open = func(string) (io.ReadWriteCloser, error) { return local, nil }

	assertRunStopsOnCancel(t, func(ctx context.Context) error {
		return ch.Run(ctx, nil)
	})
}

func TestUSBCANRunStopsOnContextCancellation(t *testing.T) {
	port := newBlockingSerialPort()
	ch := newUSBCANChannel(testLogger(), usbCANChannelOptions{
		SerialPortName: "test",
		SerialBaudRate: 2000000,
		openPort: func(string, *serial.Mode) (serial.Port, error) {
			return port, nil
		},
	})

	assertRunStopsOnCancel(t, func(ctx context.Context) error {
		return ch.Run(ctx, nil)
	})
}

func assertRunStopsOnCancel(t *testing.T, run func(context.Context) error) {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- run(ctx) }()

	time.Sleep(20 * time.Millisecond)
	cancel()

	select {
	case err := <-done:
		require.ErrorIs(t, err, context.Canceled)
	case <-time.After(time.Second):
		t.Fatal("Run did not stop after context cancellation")
	}
}

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

type blockingSerialPort struct {
	closed chan struct{}
	once   sync.Once
	mu     sync.Mutex
	writes [][]byte
}

func newBlockingSerialPort() *blockingSerialPort {
	return &blockingSerialPort{closed: make(chan struct{})}
}

func (p *blockingSerialPort) Read([]byte) (int, error) {
	<-p.closed
	return 0, io.EOF
}

func (p *blockingSerialPort) Write(buf []byte) (int, error) {
	p.mu.Lock()
	p.writes = append(p.writes, append([]byte(nil), buf...))
	p.mu.Unlock()
	return len(buf), nil
}

func (p *blockingSerialPort) writtenFrames() [][]byte {
	p.mu.Lock()
	defer p.mu.Unlock()
	frames := make([][]byte, len(p.writes))
	for i := range p.writes {
		frames[i] = append([]byte(nil), p.writes[i]...)
	}
	return frames
}
func (p *blockingSerialPort) Drain() error               { return nil }
func (p *blockingSerialPort) ResetInputBuffer() error    { return nil }
func (p *blockingSerialPort) ResetOutputBuffer() error   { return nil }
func (p *blockingSerialPort) SetDTR(bool) error          { return nil }
func (p *blockingSerialPort) SetRTS(bool) error          { return nil }
func (p *blockingSerialPort) SetMode(*serial.Mode) error { return nil }
func (p *blockingSerialPort) SetReadTimeout(time.Duration) error {
	return nil
}
func (p *blockingSerialPort) GetModemStatusBits() (*serial.ModemStatusBits, error) {
	return &serial.ModemStatusBits{}, nil
}
func (p *blockingSerialPort) Break(time.Duration) error { return nil }
func (p *blockingSerialPort) Close() error {
	p.once.Do(func() { close(p.closed) })
	return nil
}
