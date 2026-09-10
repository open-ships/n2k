package n2k

import (
	"context"
	"encoding/binary"
	"errors"
	"net"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/open-ships/n2k/internal/actisense"
	"github.com/stretchr/testify/require"
)

func TestCanceledTxConfigurationRestores(t *testing.T) {
	for _, closeSession := range []bool{false, true} {
		name := "caller cancellation"
		if closeSession {
			name = "session Close"
		}
		t.Run(name, func(t *testing.T) { testCanceledTxConfigurationRestores(t, closeSession) })
	}
}

func testCanceledTxConfigurationRestores(t *testing.T, closeSession bool) {
	host, peer := net.Pipe()
	defer func() { _ = peer.Close() }()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var enabled, activated atomic.Uint32
	var received atomic.Uint64
	var firstActivation atomic.Bool
	activation := make(chan struct{})
	go func() {
		defer func() { _ = peer.Close() }()
		parser := actisense.NewParser()
		buf := make([]byte, 4096)
		mode := uint16(1)
		for {
			n, readErr := peer.Read(buf)
			if n > 0 {
				received.Add(uint64(n))
				parser.Feed(buf[:n], func(d actisense.Datagram) {
					if d.ID != actisense.BSTBEMCommand || len(d.Payload) == 0 {
						return
					}
					command, data := d.Payload[0], d.Payload[1:]
					var out []byte
					switch command {
					case actisense.BEMOperatingMode:
						if len(data) == 2 {
							mode = binary.LittleEndian.Uint16(data)
						}
						out = make([]byte, 2)
						binary.LittleEndian.PutUint16(out, mode)
					case actisense.BEMTxPGNEnable:
						if len(data) >= 5 {
							enabled.Store(uint32(data[4]))
						}
						out = make([]byte, 14)
						copy(out, data[:4])
						out[4] = byte(enabled.Load())
					case actisense.BEMActivatePGNLists:
						activated.Store(enabled.Load())
						if !firstActivation.Swap(true) {
							close(activation)
							return
						}
					}
					payload := make([]byte, 12)
					payload[0] = command
					binary.LittleEndian.PutUint16(payload[2:4], uint16(actisense.ModelNGX1))
					wire, _ := actisense.EncodeDatagram(actisense.BSTBEMResponse, append(payload, out...))
					_, _ = peer.Write(wire)
				}, nil)
			}
			if readErr != nil {
				return
			}
		}
	}()
	session, err := NewActisenseGatewaySession(context.Background(), "review-config", func(context.Context) (ActisenseByteStream, error) { return host, nil }, WithActisenseCommandTimeout(time.Second))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = session.Close() })
	result := make(chan error, 1)
	go func() {
		result <- session.ConfigureTransmitPGNs(ctx, []ActisenseTxPGNConfiguration{{PGN: 127250, Flag: ActisensePGNEnabled}})
	}()
	select {
	case <-activation:
	case <-time.After(time.Second):
		t.Fatal("activation never reached")
	}
	require.Eventually(t, func() bool { return session.Status().Metrics.Protocol.TransportWriteBytes >= received.Load() }, time.Second, time.Millisecond)
	if closeSession {
		require.NoError(t, session.Close())
	} else {
		cancel()
	}
	err = <-result
	if !errors.Is(err, context.Canceled) {
		t.Errorf("unexpected result: %v", err)
	}
	if !closeSession && !session.Status().Connected {
		t.Error("test needs same live connection after cancellation")
	}
	require.NoError(t, session.Close())
	if activated.Load() != 0 {
		t.Fatal("canceled transaction left Tx PGN enabled even after Close")
	}
}

type restorationRequester struct {
	calls int
	fail  bool
}

func (r *restorationRequester) Request(_ context.Context, command byte, data []byte) (actisense.BEMResponse, error) {
	r.calls++
	if r.fail {
		return actisense.BEMResponse{}, errors.New("device rejected restoration")
	}
	response := actisense.BEMResponse{BEMID: command}
	if command == actisense.BEMTxPGNEnable {
		response.Data = make([]byte, 14)
		copy(response.Data, data)
	}
	return response, nil
}
func (*restorationRequester) RequestMulti(context.Context, byte, []byte, time.Duration, func([]actisense.BEMResponse) (bool, error)) ([]actisense.BEMResponse, error) {
	return nil, errors.New("unexpected multi request")
}
func TestTransmitPGNRestorationReportsFailures(t *testing.T) {
	requester := &restorationRequester{fail: true}
	commands := actisense.NewCommandSet(requester, actisense.CommandSetConfig{Timeout: time.Second})
	session := &ActisenseGatewaySession{commandTimeout: time.Second}
	originals := map[uint32]ActisenseTxPGNState{127250: {PGN: 127250, Enabled: 0}, 129025: {PGN: 129025, Enabled: 1}}
	err := session.rollbackTransmitPGNs(commands, originals)
	require.ErrorContains(t, err, "restoration is uncertain")
	require.ErrorContains(t, err, "restore Tx PGN 127250")
	require.ErrorContains(t, err, "restore Tx PGN 129025")
	require.ErrorContains(t, err, "activate restored Tx PGNs")
	require.Equal(t, 3, requester.calls)
	requester.fail = false
	require.NoError(t, session.rollbackTransmitPGNs(commands, originals))
}

// txTestDevice models separate staged and active lists and can reject exact
// requests. Only the wire reader mutates it; assertions take the same lock.
type txTestDevice struct {
	mu        sync.Mutex
	staged    map[uint32]ActisenseTxPGNState
	active    map[uint32]ActisenseTxPGNState
	failAt    map[int]bool
	requests  int
	mutations int
}

func newTxTestSession(t *testing.T, failAt ...int) (*ActisenseGatewaySession, *txTestDevice) {
	t.Helper()
	host, peer := net.Pipe()
	device := &txTestDevice{
		staged: map[uint32]ActisenseTxPGNState{127250: {PGN: 127250, Rate: 100}, 129025: {PGN: 129025, Enabled: 1, Rate: 200}},
		active: make(map[uint32]ActisenseTxPGNState), failAt: make(map[int]bool),
	}
	for number, state := range device.staged {
		device.active[number] = state
	}
	for _, request := range failAt {
		device.failAt[request] = true
	}
	done := make(chan struct{})
	go func() {
		defer close(done)
		defer func() { _ = peer.Close() }()
		parser := actisense.NewParser()
		mode := uint16(1)
		buf := make([]byte, 4096)
		for {
			n, err := peer.Read(buf)
			parser.Feed(buf[:n], func(datagram actisense.Datagram) {
				if datagram.ID != actisense.BSTBEMCommand || len(datagram.Payload) == 0 {
					return
				}
				command, data := datagram.Payload[0], datagram.Payload[1:]
				response := make([]byte, 12)
				response[0] = command
				binary.LittleEndian.PutUint16(response[2:4], uint16(actisense.ModelNGX1))
				if command == actisense.BEMOperatingMode {
					if len(data) == 2 {
						mode = binary.LittleEndian.Uint16(data)
					}
					response = binary.LittleEndian.AppendUint16(response, mode)
				} else {
					device.mu.Lock()
					device.requests++
					if device.failAt[device.requests] {
						binary.LittleEndian.PutUint32(response[8:12], 1)
					} else if command == actisense.BEMTxPGNEnable && len(data) >= 4 {
						number := binary.LittleEndian.Uint32(data)
						state := device.staged[number]
						if len(data) >= 5 {
							device.mutations++
							state.Enabled = data[4]
							if len(data) >= 9 {
								rate := binary.LittleEndian.Uint32(data[5:9])
								if rate < actisense.TxPGNRateNoChange {
									state.Rate = rate
								}
							}
							device.staged[number] = state
						}
						out := make([]byte, 14)
						binary.LittleEndian.PutUint32(out, state.PGN)
						out[4] = state.Enabled
						binary.LittleEndian.PutUint32(out[5:9], state.Rate)
						response = append(response, out...)
					} else if command == actisense.BEMActivatePGNLists {
						device.mutations++
						for number, state := range device.staged {
							device.active[number] = state
						}
					}
					device.mu.Unlock()
				}
				wire, _ := actisense.EncodeDatagram(actisense.BSTBEMResponse, response)
				_, _ = peer.Write(wire)
			}, nil)
			if err != nil {
				return
			}
		}
	}()
	t.Cleanup(func() { _ = peer.Close(); <-done })
	session, err := NewActisenseGatewaySession(t.Context(), "tx-transaction", func(context.Context) (ActisenseByteStream, error) { return host, nil }, WithActisenseCommandTimeout(time.Second))
	require.NoError(t, err)
	t.Cleanup(func() { _ = session.Close() })
	return session, device
}

func TestTransmitPGNTransactionFailureBoundaries(t *testing.T) {
	for _, test := range []struct {
		name string
		fail []int
		want string
	}{
		{"first snapshot", []int{1}, "snapshot"},
		{"second snapshot", []int{2}, "snapshot"},
		{"first stage", []int{3}, "stage"},
		{"second stage", []int{4}, "stage"},
		{"activation", []int{5}, "activate"},
		{"rollback and close retry", []int{4, 5, 6, 7}, "restoration is uncertain"},
	} {
		t.Run(test.name, func(t *testing.T) {
			session, device := newTxTestSession(t, test.fail...)
			err := session.ConfigureTransmitPGNs(t.Context(), []ActisenseTxPGNConfiguration{
				{PGN: 127250, Flag: ActisensePGNEnabled}, {PGN: 129025, Flag: ActisensePGNDisabled},
			})
			require.ErrorContains(t, err, test.want)
			require.NoError(t, session.Close())
			device.mu.Lock()
			defer device.mu.Unlock()
			require.Equal(t, ActisenseTxPGNState{PGN: 127250, Rate: 100}, device.staged[127250])
			require.Equal(t, ActisenseTxPGNState{PGN: 129025, Enabled: 1, Rate: 200}, device.staged[129025])
			require.Equal(t, device.staged, device.active)
			if test.want == "snapshot" {
				require.Zero(t, device.mutations, "snapshot failure must not mutate any entry")
			}
		})
	}
}

func TestTransmitPGNCloseRetainsEarliestOriginal(t *testing.T) {
	session, device := newTxTestSession(t)
	for _, rate := range []uint32{300, 400} {
		require.NoError(t, session.ConfigureTransmitPGNs(t.Context(), []ActisenseTxPGNConfiguration{{PGN: 127250, Flag: ActisensePGNEnabled, Rate: &rate}}))
	}
	require.NoError(t, session.Close())
	device.mu.Lock()
	defer device.mu.Unlock()
	require.Equal(t, ActisenseTxPGNState{PGN: 127250, Rate: 100}, device.active[127250])
}

func TestTransmitPGNCloseReportsRestorationFailure(t *testing.T) {
	session, _ := newTxTestSession(t, 4) // Get, Set, activate, then Close's Set.
	require.NoError(t, session.ConfigureTransmitPGNs(t.Context(), []ActisenseTxPGNConfiguration{{PGN: 127250, Flag: ActisensePGNEnabled}}))
	err := session.Close()
	require.ErrorContains(t, err, "restoration is uncertain")
	require.Equal(t, err, session.Close())
}
