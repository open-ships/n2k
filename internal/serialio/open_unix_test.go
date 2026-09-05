//go:build linux || darwin || freebsd || openbsd

package serialio

import (
	"errors"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.bug.st/serial"
)

func TestPortCloseInterruptsBlockedWrite(t *testing.T) {
	reader, writer, err := os.Pipe()
	require.NoError(t, err)
	t.Cleanup(func() { _ = reader.Close() })
	libraryPort := &serialStub{}
	connection := &port{Port: libraryPort, writer: writer}
	t.Cleanup(func() { _ = connection.Close() })
	done := make(chan error, 1)
	go func() {
		_, err := connection.Write(make([]byte, 1<<20))
		done <- err
	}()
	select {
	case err := <-done:
		t.Fatalf("write unexpectedly finished without a reader: %v", err)
	case <-time.After(20 * time.Millisecond):
	}
	require.NoError(t, connection.Close())
	select {
	case err := <-done:
		require.ErrorIs(t, err, os.ErrClosed)
	case <-time.After(time.Second):
		t.Fatal("Close did not interrupt a blocked serial write")
	}
	require.Equal(t, int64(1), libraryPort.closes.Load())
}

func TestPortConcurrentCloseOwnsBothDescriptorsOnce(t *testing.T) {
	reader, writer, err := os.Pipe()
	require.NoError(t, err)
	t.Cleanup(func() { _ = reader.Close() })
	want := errors.New("serial close diagnostic")
	libraryPort := &serialStub{err: want}
	connection := &port{Port: libraryPort, writer: writer}
	var wg sync.WaitGroup
	results := make(chan error, 8)
	for range 8 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			results <- connection.Close()
		}()
	}
	wg.Wait()
	close(results)
	for err := range results {
		require.ErrorIs(t, err, want)
	}
	require.Equal(t, int64(1), libraryPort.closes.Load())
	_, err = writer.Write([]byte{1})
	require.ErrorIs(t, err, os.ErrClosed)
}

type serialStub struct {
	serial.Port
	closes atomic.Int64
	err    error
}

func (s *serialStub) Close() error {
	s.closes.Add(1)
	return s.err
}
