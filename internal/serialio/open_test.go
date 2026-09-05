package serialio

import (
	"path/filepath"
	"testing"

	"go.bug.st/serial"
)

func TestOpenRejectsMissingDevice(t *testing.T) {
	connection, err := Open(filepath.Join(t.TempDir(), "missing-device"), &serial.Mode{BaudRate: 115200})
	if connection != nil {
		_ = connection.Close()
		t.Fatal("opening a missing device returned a port")
	}
	if err == nil {
		t.Fatal("opening a missing device returned no error")
	}
}
