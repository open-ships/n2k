package transport

import (
	"testing"
	"time"

	"github.com/brutella/can"
	"github.com/open-ships/n2k/internal/framer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRTSCTSReceive_HappyPath(t *testing.T) {
	h := newTestHelper()
	mgr := NewManager(ManagerConfig{
		WriteFrame: h.writeFrame,
		OnComplete: h.onComplete,
	})
	defer mgr.Close()

	source := uint8(42)      // the remote sender
	destination := uint8(10) // us, the receiver
	transportedPGN := uint32(126998)

	// Create a 16-byte payload (3 DT frames).
	payload := make([]byte, 16)
	for i := range payload {
		payload[i] = uint8(i + 0x30)
	}

	// Step 1: Remote sends RTS.
	rtsFrame := buildTestCMFrame(ControlRTS, 16, 3, transportedPGN, source, destination)
	// Set byte 4 to 0xFF (unlimited frames per CTS).
	rtsFrame.Data[4] = 0xFF
	mgr.HandleFrame(rtsFrame)

	// Verify we sent a CTS in response.
	sentFrames := h.getSentFrames()
	require.Len(t, sentFrames, 1, "should have sent a CTS")

	ctsFrame := sentFrames[0]
	cts := framer.ParseCANID(ctsFrame.ID)
	assert.Equal(t, PGNCM, cts.PGN)
	assert.Equal(t, destination, cts.Source, "CTS source should be us")
	assert.Equal(t, source, cts.Destination, "CTS destination should be the remote sender")
	assert.Equal(t, ControlCTS, ctsFrame.Data[0])
	assert.Equal(t, uint8(3), ctsFrame.Data[1], "request all 3 frames")
	assert.Equal(t, uint8(1), ctsFrame.Data[2], "start from frame 1")

	// Step 2: Remote sends DT frames.
	var dt1Data [7]byte
	copy(dt1Data[:], payload[0:7])
	mgr.HandleFrame(buildTestDTFrame(1, dt1Data, source, destination))

	var dt2Data [7]byte
	copy(dt2Data[:], payload[7:14])
	mgr.HandleFrame(buildTestDTFrame(2, dt2Data, source, destination))

	var dt3Data [7]byte
	for j := range dt3Data {
		dt3Data[j] = 0xFF
	}
	copy(dt3Data[:], payload[14:16])
	mgr.HandleFrame(buildTestDTFrame(3, dt3Data, source, destination))

	// Verify EndOfMsgAck was sent.
	sentFrames = h.getSentFrames()
	require.Len(t, sentFrames, 2, "should have sent CTS + EndOfMsgAck")

	eomFrame := sentFrames[1]
	eom := framer.ParseCANID(eomFrame.ID)
	assert.Equal(t, PGNCM, eom.PGN)
	assert.Equal(t, destination, eom.Source, "EOM source should be us")
	assert.Equal(t, source, eom.Destination, "EOM destination should be the remote sender")
	assert.Equal(t, ControlEndOfMsgAck, eomFrame.Data[0])
	assert.Equal(t, uint16(16), uint16(eomFrame.Data[1])|uint16(eomFrame.Data[2])<<8)
	assert.Equal(t, uint8(3), eomFrame.Data[3])

	// Verify the assembled payload was delivered.
	msg := h.waitComplete(t, time.Second)
	assert.Equal(t, transportedPGN, msg.pgn)
	assert.Equal(t, source, msg.source)
	assert.Equal(t, destination, msg.destination)
	assert.Equal(t, payload, msg.data)
}

func TestRTSCTSReceive_RequestsEachCTSWindow(t *testing.T) {
	h := newTestHelper()
	mgr := NewManager(ManagerConfig{WriteFrame: h.writeFrame, OnComplete: h.onComplete})
	defer mgr.Close()

	const (
		source      = uint8(42)
		destination = uint8(10)
		messagePGN  = uint32(126998)
	)
	payload := make([]byte, 29) // five DT packets
	for i := range payload {
		payload[i] = byte(i + 1)
	}
	rts := buildTestCMFrame(ControlRTS, uint16(len(payload)), 5, messagePGN, source, destination)
	rts.Data[4] = 2
	mgr.HandleFrame(rts)

	assertCTS := func(index int, count, first uint8) {
		t.Helper()
		frames := h.getSentFrames()
		require.Greater(t, len(frames), index)
		assert.Equal(t, ControlCTS, frames[index].Data[0])
		assert.Equal(t, count, frames[index].Data[1])
		assert.Equal(t, first, frames[index].Data[2])
	}
	assertCTS(0, 2, 1)

	for seq := uint8(1); seq <= 5; seq++ {
		var data [7]byte
		start := int(seq-1) * MaxDTDataBytes
		if start < len(payload) {
			copy(data[:], payload[start:])
		}
		mgr.HandleFrame(buildTestDTFrame(seq, data, source, destination))
		if seq == 2 {
			assertCTS(1, 2, 3)
		}
		if seq == 4 {
			assertCTS(2, 1, 5)
		}
	}

	msg := h.waitComplete(t, time.Second)
	assert.Equal(t, payload, msg.data)
	frames := h.getSentFrames()
	require.Len(t, frames, 4)
	assert.Equal(t, ControlEndOfMsgAck, frames[3].Data[0])
}

func TestRTSCTSTransmit_HappyPath(t *testing.T) {
	h := newTestHelper()
	mgr := NewManager(ManagerConfig{
		WriteFrame: h.writeFrame,
		OnComplete: h.onComplete,
	})
	defer mgr.Close()

	source := uint8(10)      // us, the sender
	destination := uint8(42) // the remote receiver
	transportedPGN := uint32(126998)

	// Create a 16-byte payload (3 DT frames).
	payload := make([]byte, 16)
	for i := range payload {
		payload[i] = uint8(i + 0x40)
	}

	// Start the transmit in a goroutine since it blocks.
	errCh := make(chan error, 1)
	go func() {
		errCh <- mgr.SendRTSCTS(transportedPGN, source, destination, payload)
	}()

	// Wait a bit for the RTS to be sent.
	time.Sleep(50 * time.Millisecond)

	// Verify RTS was sent.
	sentFrames := h.getSentFrames()
	require.NotEmpty(t, sentFrames, "should have sent an RTS")

	rtsFrame := sentFrames[0]
	rts := framer.ParseCANID(rtsFrame.ID)
	assert.Equal(t, PGNCM, rts.PGN)
	assert.Equal(t, source, rts.Source)
	assert.Equal(t, destination, rts.Destination)
	assert.Equal(t, ControlRTS, rtsFrame.Data[0])
	assert.Equal(t, uint16(16), uint16(rtsFrame.Data[1])|uint16(rtsFrame.Data[2])<<8)
	assert.Equal(t, uint8(3), rtsFrame.Data[3])

	// Simulate: remote sends CTS requesting all 3 frames starting at 1.
	ctsFrame := buildTestCTSFrame(3, 1, transportedPGN, destination, source)
	mgr.HandleFrame(ctsFrame)

	// Wait for DT frames to be sent.
	time.Sleep(50 * time.Millisecond)

	sentFrames = h.getSentFrames()
	// Should have: RTS + 3 DT frames = 4 frames.
	require.GreaterOrEqual(t, len(sentFrames), 4, "should have sent RTS + 3 DT frames")

	// Verify DT frames.
	var assembled []byte
	for i := 1; i <= 3; i++ {
		dtFrame := sentFrames[i]
		dt := framer.ParseCANID(dtFrame.ID)
		assert.Equal(t, PGNDT, dt.PGN)
		assert.Equal(t, source, dt.Source)
		assert.Equal(t, destination, dt.Destination)
		assert.Equal(t, uint8(i), dtFrame.Data[0], "DT frame %d sequence number", i)
		assembled = append(assembled, dtFrame.Data[1:]...)
	}
	assembled = assembled[:16]
	assert.Equal(t, payload, assembled, "DT frame data should match payload")

	// Simulate: remote sends EndOfMsgAck.
	eomFrame := buildTestEndOfMsgAckFrame(16, 3, transportedPGN, destination, source)
	mgr.HandleFrame(eomFrame)

	// Wait for SendRTSCTS to return.
	select {
	case err := <-errCh:
		assert.NoError(t, err)
	case <-time.After(2 * time.Second):
		t.Fatal("SendRTSCTS did not return after EndOfMsgAck")
	}
}

func TestRTSCTSTransmit_RejectsInvalidCTSRange(t *testing.T) {
	h := newTestHelper()
	mgr := NewManager(ManagerConfig{WriteFrame: h.writeFrame})
	defer mgr.Close()

	errCh := make(chan error, 1)
	go func() {
		errCh <- mgr.SendRTSCTS(126998, 10, 42, make([]byte, 16))
	}()
	require.Eventually(t, func() bool { return len(h.getSentFrames()) == 1 }, time.Second, time.Millisecond)

	// Three packets exist, so a request for packets 2 through 4 is invalid.
	mgr.HandleFrame(buildTestCTSFrame(3, 2, 126998, 42, 10))
	select {
	case err := <-errCh:
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid CTS range")
	case <-time.After(time.Second):
		t.Fatal("SendRTSCTS did not reject invalid CTS")
	}
	assert.Len(t, h.getSentFrames(), 1, "invalid CTS must not emit DT packets")
}

func TestRTSCTSTransmit_RejectsMismatchedEndAck(t *testing.T) {
	h := newTestHelper()
	mgr := NewManager(ManagerConfig{WriteFrame: h.writeFrame})
	defer mgr.Close()

	errCh := make(chan error, 1)
	go func() {
		errCh <- mgr.SendRTSCTS(126998, 10, 42, make([]byte, 16))
	}()
	require.Eventually(t, func() bool { return len(h.getSentFrames()) == 1 }, time.Second, time.Millisecond)
	mgr.HandleFrame(buildTestCTSFrame(3, 1, 126998, 42, 10))
	require.Eventually(t, func() bool { return len(h.getSentFrames()) == 4 }, time.Second, time.Millisecond)

	mgr.HandleFrame(buildTestEndOfMsgAckFrame(15, 3, 126998, 42, 10))
	select {
	case err := <-errCh:
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid EndOfMsgAck")
	case <-time.After(time.Second):
		t.Fatal("SendRTSCTS did not reject mismatched EndOfMsgAck")
	}
}

func TestRTSCTSTransmit_EndAckRequiresEveryPacket(t *testing.T) {
	h := newTestHelper()
	mgr := NewManager(ManagerConfig{WriteFrame: h.writeFrame})
	defer mgr.Close()

	errCh := make(chan error, 1)
	go func() {
		errCh <- mgr.SendRTSCTS(126998, 10, 42, make([]byte, 16))
	}()
	require.Eventually(t, func() bool { return len(h.getSentFrames()) == 1 }, time.Second, time.Millisecond)

	// A hostile peer requests only the last of three packets, then claims the
	// transfer is complete. Seeing sequence 3 must not imply 1 and 2 were sent.
	mgr.HandleFrame(buildTestCTSFrame(1, 3, 126998, 42, 10))
	require.Eventually(t, func() bool { return len(h.getSentFrames()) == 2 }, time.Second, time.Millisecond)
	mgr.HandleFrame(buildTestEndOfMsgAckFrame(16, 3, 126998, 42, 10))

	select {
	case err := <-errCh:
		require.Error(t, err)
		assert.Contains(t, err.Error(), "packets sent=1")
	case <-time.After(time.Second):
		t.Fatal("SendRTSCTS accepted an acknowledgement with missing packets")
	}
}

func TestRTSCTSTransmit_RejectsAmbiguousConcurrentSession(t *testing.T) {
	h := newTestHelper()
	mgr := NewManager(ManagerConfig{WriteFrame: h.writeFrame})

	first := make(chan error, 1)
	go func() { first <- mgr.SendRTSCTS(126998, 10, 42, make([]byte, 16)) }()
	require.Eventually(t, func() bool { return len(h.getSentFrames()) == 1 }, time.Second, time.Millisecond)

	err := mgr.SendRTSCTS(126996, 10, 42, make([]byte, 16))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "session already active")
	assert.Len(t, h.getSentFrames(), 1, "a second RTS would make incoming CTS ownership ambiguous")

	mgr.Close()
	require.Error(t, <-first)
}

// buildTestCTSFrame creates a CTS frame for testing.
func buildTestCTSFrame(numFrames uint8, nextSeqNum uint8, pgn uint32, source uint8, destination uint8) can.Frame {
	var data [8]uint8
	data[0] = ControlCTS
	data[1] = numFrames
	data[2] = nextSeqNum
	data[3] = 0xFF
	data[4] = 0xFF
	encodePGN(data[5:8], pgn)

	canID := framer.BuildCANID(PGNCM, TPPriority, source, destination)
	return can.Frame{
		ID:     canID,
		Length: 8,
		Data:   data,
	}
}

// buildTestEndOfMsgAckFrame creates an EndOfMsgAck frame for testing.
func buildTestEndOfMsgAckFrame(totalSize uint16, numFrames uint8, pgn uint32, source uint8, destination uint8) can.Frame {
	var data [8]uint8
	data[0] = ControlEndOfMsgAck
	data[1] = uint8(totalSize)
	data[2] = uint8(totalSize >> 8)
	data[3] = numFrames
	data[4] = 0xFF
	encodePGN(data[5:8], pgn)

	canID := framer.BuildCANID(PGNCM, TPPriority, source, destination)
	return can.Frame{
		ID:     canID,
		Length: 8,
		Data:   data,
	}
}
