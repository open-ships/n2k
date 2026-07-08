package transport

import (
	"sync"
	"testing"
	"time"

	"github.com/brutella/can"
	"github.com/open-ships/n2k/internal/framer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testHelper provides common utilities for transport tests.
type testHelper struct {
	mu         sync.Mutex
	sentFrames []can.Frame
	completed  []completedMsg
	completeCh chan completedMsg
}

type completedMsg struct {
	pgn         uint32
	source      uint8
	destination uint8
	data        []byte
}

func newTestHelper() *testHelper {
	return &testHelper{
		completeCh: make(chan completedMsg, 10),
	}
}

func (h *testHelper) writeFrame(f can.Frame) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.sentFrames = append(h.sentFrames, f)
	return nil
}

func (h *testHelper) onComplete(pgn uint32, source uint8, destination uint8, data []byte) {
	msg := completedMsg{pgn: pgn, source: source, destination: destination, data: data}
	h.mu.Lock()
	h.completed = append(h.completed, msg)
	h.mu.Unlock()
	h.completeCh <- msg
}

func (h *testHelper) getSentFrames() []can.Frame {
	h.mu.Lock()
	defer h.mu.Unlock()
	out := make([]can.Frame, len(h.sentFrames))
	copy(out, h.sentFrames)
	return out
}

func (h *testHelper) getCompleted() []completedMsg {
	h.mu.Lock()
	defer h.mu.Unlock()
	out := make([]completedMsg, len(h.completed))
	copy(out, h.completed)
	return out
}

func (h *testHelper) waitComplete(t *testing.T, timeout time.Duration) completedMsg {
	t.Helper()
	select {
	case msg := <-h.completeCh:
		return msg
	case <-time.After(timeout):
		t.Fatal("timed out waiting for completed message")
		return completedMsg{}
	}
}

// buildTestCMFrame creates a CM frame for testing purposes.
func buildTestCMFrame(controlByte uint8, totalSize uint16, frameCount uint8, pgn uint32, source uint8, destination uint8) can.Frame {
	var data [8]uint8
	data[0] = controlByte
	data[1] = uint8(totalSize)
	data[2] = uint8(totalSize >> 8)
	data[3] = frameCount
	data[4] = 0xFF
	encodePGN(data[5:8], pgn)

	canID := framer.BuildCANID(PGNCM, TPPriority, source, destination)
	return can.Frame{
		ID:     canID,
		Length: 8,
		Data:   data,
	}
}

// buildTestDTFrame creates a DT frame for testing purposes.
func buildTestDTFrame(seqNum uint8, data [7]byte, source uint8, destination uint8) can.Frame {
	var frameData [8]uint8
	frameData[0] = seqNum
	copy(frameData[1:], data[:])

	canID := framer.BuildCANID(PGNDT, TPPriority, source, destination)
	return can.Frame{
		ID:     canID,
		Length: 8,
		Data:   frameData,
	}
}

func TestBAMReceive_HappyPath(t *testing.T) {
	h := newTestHelper()
	mgr := NewManager(ManagerConfig{
		WriteFrame: h.writeFrame,
		OnComplete: h.onComplete,
	})
	defer mgr.Close()

	source := uint8(42)
	transportedPGN := uint32(126996) // Product Information PGN

	// Create a 16-byte payload: requires 3 DT frames (7+7+2 bytes).
	payload := make([]byte, 16)
	for i := range payload {
		payload[i] = uint8(i + 0x10)
	}

	// Send CM_BAM announcement.
	bamFrame := buildTestCMFrame(ControlBAM, 16, 3, transportedPGN, source, BroadcastAddr)
	mgr.HandleFrame(bamFrame)

	// Send DT frame 1: bytes 0-6.
	var dt1Data [7]byte
	copy(dt1Data[:], payload[0:7])
	mgr.HandleFrame(buildTestDTFrame(1, dt1Data, source, BroadcastAddr))

	// Send DT frame 2: bytes 7-13.
	var dt2Data [7]byte
	copy(dt2Data[:], payload[7:14])
	mgr.HandleFrame(buildTestDTFrame(2, dt2Data, source, BroadcastAddr))

	// Send DT frame 3: bytes 14-15, padded with 0xFF.
	var dt3Data [7]byte
	for j := range dt3Data {
		dt3Data[j] = 0xFF
	}
	copy(dt3Data[:], payload[14:16])
	mgr.HandleFrame(buildTestDTFrame(3, dt3Data, source, BroadcastAddr))

	// Verify the assembled payload.
	msg := h.waitComplete(t, time.Second)
	assert.Equal(t, transportedPGN, msg.pgn)
	assert.Equal(t, source, msg.source)
	assert.Equal(t, BroadcastAddr, msg.destination)
	assert.Equal(t, payload, msg.data)
}

func TestBAMReceive_Timeout(t *testing.T) {
	h := newTestHelper()
	mgr := NewManager(ManagerConfig{
		WriteFrame: h.writeFrame,
		OnComplete: h.onComplete,
	})
	defer mgr.Close()

	source := uint8(42)
	transportedPGN := uint32(126996)

	// Send CM_BAM for a 16-byte message (3 DT frames expected).
	bamFrame := buildTestCMFrame(ControlBAM, 16, 3, transportedPGN, source, BroadcastAddr)
	mgr.HandleFrame(bamFrame)

	// Send only the first DT frame.
	var dt1Data [7]byte
	for i := range dt1Data {
		dt1Data[i] = uint8(i)
	}
	mgr.HandleFrame(buildTestDTFrame(1, dt1Data, source, BroadcastAddr))

	// Wait for the DT timeout to fire.
	time.Sleep(DTTimeout + 100*time.Millisecond)

	// Verify no message was delivered.
	completed := h.getCompleted()
	assert.Empty(t, completed, "should not deliver incomplete message")

	// Verify the session was cleaned up.
	mgr.mu.Lock()
	assert.Empty(t, mgr.sessions, "session should be cleaned up after timeout")
	mgr.mu.Unlock()
}

func TestBAMTransmit_HappyPath(t *testing.T) {
	h := newTestHelper()
	mgr := NewManager(ManagerConfig{
		WriteFrame: h.writeFrame,
		OnComplete: h.onComplete,
	})
	defer mgr.Close()

	source := uint8(10)
	transportedPGN := uint32(126996)

	// Create a 16-byte payload.
	payload := make([]byte, 16)
	for i := range payload {
		payload[i] = uint8(i + 0x20)
	}

	err := mgr.SendBAM(transportedPGN, source, payload)
	require.NoError(t, err)

	frames := h.getSentFrames()
	// Expect: 1 CM_BAM + 3 DT frames.
	require.Len(t, frames, 4)

	// Verify CM_BAM frame.
	cmFrame := frames[0]
	cm := framer.ParseCANID(cmFrame.ID)
	assert.Equal(t, source, cm.Source)
	assert.Equal(t, BroadcastAddr, cm.Destination)
	assert.Equal(t, ControlBAM, cmFrame.Data[0])
	assert.Equal(t, uint16(16), uint16(cmFrame.Data[1])|uint16(cmFrame.Data[2])<<8)
	assert.Equal(t, uint8(3), cmFrame.Data[3])
	assert.Equal(t, transportedPGN, extractPGN(cmFrame.Data))

	// Verify DT frames carry the correct data.
	var assembled []byte
	for i := 1; i < len(frames); i++ {
		dtFrame := frames[i]
		assert.Equal(t, uint8(i), dtFrame.Data[0], "DT sequence number")
		assembled = append(assembled, dtFrame.Data[1:]...)
	}
	assembled = assembled[:16] // trim padding
	assert.Equal(t, payload, assembled)
}

func TestBAMTransmit_SmallPayload(t *testing.T) {
	h := newTestHelper()
	mgr := NewManager(ManagerConfig{
		WriteFrame: h.writeFrame,
		OnComplete: h.onComplete,
	})
	defer mgr.Close()

	// 9-byte payload: needs 2 DT frames (7 + 2 bytes).
	payload := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09}

	err := mgr.SendBAM(126996, 5, payload)
	require.NoError(t, err)

	frames := h.getSentFrames()
	require.Len(t, frames, 3, "1 CM_BAM + 2 DT frames")

	// Verify CM_BAM.
	assert.Equal(t, ControlBAM, frames[0].Data[0])
	assert.Equal(t, uint8(9), frames[0].Data[1]) // total size low byte
	assert.Equal(t, uint8(0), frames[0].Data[2]) // total size high byte
	assert.Equal(t, uint8(2), frames[0].Data[3]) // frame count

	// Verify DT frame 1: seq=1, data=payload[0:7].
	assert.Equal(t, uint8(1), frames[1].Data[0])
	assert.Equal(t, payload[0:7], frames[1].Data[1:8])

	// Verify DT frame 2: seq=2, data=payload[7:9] + padding.
	assert.Equal(t, uint8(2), frames[2].Data[0])
	assert.Equal(t, payload[7:9], frames[2].Data[1:3])
	assert.Equal(t, uint8(0xFF), frames[2].Data[3])
	assert.Equal(t, uint8(0xFF), frames[2].Data[7])
}

func TestBAMTransmit_LargePayload(t *testing.T) {
	h := newTestHelper()
	mgr := NewManager(ManagerConfig{
		WriteFrame: h.writeFrame,
		OnComplete: h.onComplete,
	})
	defer mgr.Close()

	// 100-byte payload: needs ceil(100/7) = 15 DT frames.
	payload := make([]byte, 100)
	for i := range payload {
		payload[i] = uint8(i)
	}

	err := mgr.SendBAM(126996, 7, payload)
	require.NoError(t, err)

	frames := h.getSentFrames()
	expectedDTFrames := 15 // ceil(100/7)
	require.Len(t, frames, 1+expectedDTFrames, "1 CM_BAM + 15 DT frames")

	// Verify CM_BAM.
	assert.Equal(t, ControlBAM, frames[0].Data[0])
	assert.Equal(t, uint16(100), uint16(frames[0].Data[1])|uint16(frames[0].Data[2])<<8)
	assert.Equal(t, uint8(expectedDTFrames), frames[0].Data[3])

	// Reassemble and verify data integrity.
	var assembled []byte
	for i := 1; i < len(frames); i++ {
		assert.Equal(t, uint8(i), frames[i].Data[0], "DT frame %d sequence number", i)
		assembled = append(assembled, frames[i].Data[1:]...)
	}
	assembled = assembled[:100]
	assert.Equal(t, payload, assembled, "reassembled data should match original payload")
}
