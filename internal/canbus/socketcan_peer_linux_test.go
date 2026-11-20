package canbus

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// N2K_TEST_VCAN names an isolated, already-up Linux virtual CAN interface.
func TestSocketCANLinuxPeer(t *testing.T) {
	iface := os.Getenv("N2K_TEST_VCAN")
	if iface == "" {
		t.Skip("set N2K_TEST_VCAN to run with two independent Linux CAN sockets")
	}
	peer, err := openSocketCAN(iface)
	require.NoError(t, err)
	t.Cleanup(func() { _ = peer.Close() })
	require.NoError(t, peer.(*os.File).SetDeadline(time.Now().Add(5*time.Second)))
	verifySocketCANPeer(t, NewSocketCAN(testLogger(), iface), peer)
}
