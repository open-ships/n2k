package n2k

import "github.com/brutella/can"

// ReplayFrameCapacity bounds the recent written-frame capture of a replay
// Client. Older frames are evicted on overflow; Status reports every eviction.
const ReplayFrameCapacity = 4096

// captureReplayFrameLocked retains an owned frame value in constant time.
// The caller holds Client.mu; CAN frames contain fixed arrays, not slices.
func (c *Client) captureReplayFrameLocked(frame can.Frame) {
	if c.writtenFrames == nil {
		c.writtenFrames = make([]can.Frame, 0, ReplayFrameCapacity)
	}
	if len(c.writtenFrames) < ReplayFrameCapacity {
		c.writtenFrames = append(c.writtenFrames, frame)
		return
	}
	c.writtenFrames[c.writtenFramesStart] = frame
	c.writtenFramesStart = (c.writtenFramesStart + 1) % ReplayFrameCapacity
	c.writtenFramesDropped++
}
