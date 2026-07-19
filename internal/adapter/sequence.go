package adapter

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/open-ships/n2k/internal/decoder"
)

// MaxFrameNum is the maximum frame number in a multipart NMEA 2000 fast-packet message.
// Frame numbers are encoded as 5 bits (0-31), so a single fast-packet sequence can span
// up to 32 CAN frames. With 6 data bytes in frame 0 and 7 in each subsequent frame,
// this allows payloads up to 6 + (31 * 7) = 223 bytes.
const MaxFrameNum = 31

// sequence defines the data and methods needed to assemble a series of CAN frames into a
// single complete NMEA 2000 fast-packet message.
//
// NMEA 2000 fast-packet protocol overview:
//   - Messages with more than 8 bytes of payload are split across multiple CAN frames.
//   - Each frame's first data byte encodes the sequence ID (bits 7-5) and frame number (bits 4-0).
//   - Frame 0 is special: its second byte contains the total payload length (excluding the
//     header bytes), and it carries 6 bytes of actual data (bytes 2-7).
//   - Continuation frames (1-31) carry 7 bytes of data each (bytes 1-7), with byte 0 used
//     for the sequence ID / frame number header.
//   - The final frame may have unused padding bytes that are trimmed based on the expected
//     length declared in frame 0.
//   - Frame 0 must be received first; other frames may arrive in any order.
//   - The 3-bit sequence ID (0-7) allows up to 8 concurrent sequences for the same PGN/source.
//   - The 5-bit frame number (0-31) identifies each frame's position within the sequence.
type sequence struct {
	// zero holds the initial frame (frame number 0) of this sequence. It must be received
	// before any continuation frames can be accepted. If nil, no frame 0 has been received yet.
	zero *decoder.Packet

	// expected is the total number of payload bytes expected for this message, as declared
	// in byte 1 of frame 0. This is used to determine when all data has been received and
	// to trim padding bytes from the final frame.
	expected uint8

	// contents stores the data bytes from each frame, indexed by frame number. This array
	// allows frames to be received out of order (except frame 0 which must come first).
	// A nil entry means that frame has not been received yet. The array size is MaxFrameNum+1
	// (32 slots) to accommodate frame numbers 0-31.
	contents [MaxFrameNum + 1][]uint8 // need arrays since packets can be received out of order
	updated  time.Time
}

// add copies the payload data from a CAN frame into the appropriate slot in the sequence.
//
// For frame 0:
//   - Stores the packet as the sequence's zero reference
//   - Reads byte 1 as the expected total payload length
//   - Copies bytes 2-7 (6 bytes of data) into contents[0]
//
// For continuation frames (frame number > 0):
//   - Requires frame 0 to have been received first (resets if not)
//   - Detects and resets on duplicate frame numbers (assumes a new sequence started)
//   - Copies bytes 1-7 (7 bytes of data) into contents[frameNum]
//
// If a duplicate frame 0 arrives before the previous sequence completes, the old sequence
// is discarded and a new one starts. Similarly, receiving a duplicate continuation frame
// resets the entire sequence under the assumption that a new transmission has started.
//
// Parameters:
//   - p: The Packet containing the raw CAN frame data with sequence/frame header in byte 0.
func (s *sequence) add(p *decoder.Packet, now time.Time) bool {
	s.updated = now
	if p.FrameNum == 0 {
		if s.zero != nil { // we've received frame zero for a new sequence before completing the previous one.
			slog.Debug("Fast sequence duplicate frame zero detected. Resetting")
			s.reset() // so we toss the old one and start anew
		}
		s.zero = p
		// Byte 1 of frame 0 contains the total expected payload length (not counting
		// the 2 header bytes of frame 0 or the 1 header byte of continuation frames).
		if len(p.Data) < 2 {
			p.ParseErrors = append(p.ParseErrors, fmt.Errorf("fast-packet frame zero is shorter than 2 bytes"))
			return false
		}
		s.expected = p.Data[1]
		if s.expected == 0 || int(s.expected) > 223 {
			p.ParseErrors = append(p.ParseErrors, fmt.Errorf("invalid fast-packet payload length %d", s.expected))
			return false
		}
		// Bytes 2-7 of frame 0 contain the first 6 bytes of the actual payload.
		s.contents[p.FrameNum] = append([]uint8(nil), p.Data[2:]...)
	} else {
		if s.zero == nil { // we've received a subsequent frame before getting the first one
			slog.Debug("fast sequence received subsequent frame before zero frame, resetting",
				"source", p.Info.SourceId, "pgn", p.Info.PGN, "seqId", p.SeqId, "frameNum", p.FrameNum)
			s.reset()
			s.updated = now
			return true
		} else if s.contents[p.FrameNum] != nil { // uh-oh, we've already seen this frame
			// Duplicate frame detected -- likely a new sequence has started with the same
			// sequence ID before the old one completed. Reset to avoid mixing data from
			// two different transmissions.
			slog.Debug("fast sequence received duplicate frame, resetting sequence",
				"source", p.Info.SourceId, "pgn", p.Info.PGN, "seqId", p.SeqId, "frameNum", p.FrameNum)
			s.reset()
			s.updated = now
			return true
		} else {
			// Normal continuation frame: copy bytes 1-7 (7 data bytes, skipping the
			// sequence/frame header in byte 0).
			s.contents[p.FrameNum] = append([]uint8(nil), p.Data[1:]...)
		}
	}
	return true
}

// complete checks whether all expected payload bytes have been received and, if so,
// assembles the final contiguous data buffer from the per-frame contents.
//
// Assembly process:
//  1. Verify that frame 0 has been received (s.zero != nil).
//  2. Check if received bytes >= expected bytes.
//  3. Concatenate frame contents in order (0, 1, 2, ...), stopping when enough bytes
//     are collected. If any intermediate frame is missing (nil), this indicates a sparse
//     sequence and a parse error is recorded.
//  4. Trim the assembled buffer to exactly s.expected bytes (removes padding from the
//     last CAN frame).
//  5. Store the assembled data in p.Data and set p.Complete = true.
//
// Parameters:
//   - p: The Packet to finalize. On success, p.Data is replaced with the assembled payload
//     and p.Complete is set to true. On sparse-data errors, p.ParseErrors is populated.
//
// Returns true if the sequence is complete (either successfully or with errors), false
// if more frames are still needed.
func (s *sequence) complete(p *decoder.Packet) bool {
	if s.zero == nil {
		p.Complete = false
		return false
	}
	remaining := int(s.expected) - 6
	lastFrame := 0
	if remaining > 0 {
		lastFrame = (remaining + 6) / 7
	}
	results := make([]uint8, 0, s.expected)
	for i := 0; i <= lastFrame; i++ {
		if s.contents[i] == nil {
			p.Complete = false
			return false
		}
		results = append(results, s.contents[i]...)
	}
	if len(results) < int(s.expected) {
		p.ParseErrors = append(p.ParseErrors, fmt.Errorf("fast-packet frames contain %d bytes, expected %d", len(results), s.expected))
		return true
	}
	p.Data = append([]uint8(nil), results[:s.expected]...)
	// Message context belongs to the start of the wire transfer, not whichever
	// continuation frame happened to complete an out-of-order assembly.
	p.Info = s.zero.Info
	p.Complete = true
	return true
}

// reset clears all sequence state to allow reuse of the sequence slot for a new
// transmission. This is called when duplicate frames are detected, indicating that the
// sender has started a new sequence with the same sequence ID before the previous one
// completed (e.g., due to bus errors or missed frames).
func (s *sequence) reset() {
	s.zero = nil
	s.expected = 0
	s.updated = time.Time{}
	// Clear all stored frame data to prevent stale data from mixing with new frames.
	for i := range s.contents {
		s.contents[i] = nil
	}
}
