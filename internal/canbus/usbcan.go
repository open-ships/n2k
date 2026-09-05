package canbus

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"sync"

	"github.com/brutella/can"
	"github.com/open-ships/n2k/internal/serialio"
	"go.bug.st/serial"
)

// USB-CAN Analyzer Protocol Documentation:
//
// This file implements the USB-CAN Analyzer protocol, a cheap USB-to-CAN dongle commonly
// sold under names like "USB-CAN Analyzer" or "Seeed USB CAN Analyzer". The device
// communicates over a virtual serial port using a proprietary binary framing protocol.
//
// Protocol references:
// * https://github.com/SeeedDocument/USB-CAN-Analyzer/blob/master/res/Document/USB%20(Serial%20port)%20to%20CAN%20protocol%20defines.pdf
// * https://github.com/SeeedDocument/USB-CAN-Analyzer?tab=readme-ov-file
// * https://github.com/kobolt/usb-can/blob/master/canusb.c
//
// Frame format overview:
//   - Every frame starts with 0xAA (start-of-frame marker)
//   - Every data/command frame ends with 0x55 (end-of-frame marker)
//   - There are two major frame categories:
//     1. Command frames: 0xAA 0x55 followed by 18 bytes of settings/response data (20 bytes total)
//     2. Data frames: 0xAA <info_byte> <id_bytes> <data_bytes> 0x55
//
// Data frame info byte layout (byte[1]):
//   - Bits [7:6] = 0b11 indicates a data frame
//   - Bit  [5]   = 1 for extended (29-bit) CAN ID, 0 for standard (11-bit) CAN ID
//   - Bit  [4]   = 1 for remote frame, 0 for data frame
//   - Bits [3:0] = data length (0-8 bytes of CAN payload)

// usbCANChannelOptions contains the configuration required to open and operate a USB-CAN channel.
type usbCANChannelOptions struct {
	// SerialPortName is the OS path to the serial port device, e.g. "/dev/ttyUSB0" on Linux
	// or "/dev/cu.usbserial-1234" on macOS.
	SerialPortName string `json:"serialPortName"`

	// SerialBaudRate is the baud rate for the serial (UART) connection between the host and
	// the USB-CAN dongle. This is NOT the CAN bus bitrate -- it controls only the USB serial link.
	// The USB-CAN Analyzer typically uses 2000000 (2 Mbaud).
	SerialBaudRate int `json:"serialBaudRate"`
	openPort       func(string, *serial.Mode) (serial.Port, error)
}

// usbCANChannel represents a single USB-CAN-based canbus channel for sending/receiving CAN frames.
// It manages the serial port connection, sends the initial configuration frame, and continuously
// reads and parses the proprietary USB-CAN binary protocol from the serial stream.
type usbCANChannel struct {
	// options holds the user-provided configuration (serial port, bitrate, etc.);
	// the message handler is passed separately to Run.
	options usbCANChannelOptions

	// port is the open serial port connection to the USB-CAN dongle.
	// It is nil until Run() is called.
	port serial.Port

	// mu guards the port reference. Reads and writes must not hold it because
	// Close needs to close the port concurrently to unblock them.
	mu        sync.Mutex
	writeMu   sync.Mutex
	ready     chan struct{}
	readyOnce sync.Once

	// log is the structured logger for debug/info messages about frame parsing and errors.
	log *slog.Logger

	// handler is the callback invoked for each successfully parsed CAN frame.
	// It is set by Run() for the duration of the read loop.
	handler func(can.Frame)
}

// newUSBCANChannel creates and returns a new usbCANChannel configured with the given options.
// The channel is not opened until Run() is called.
//
// Parameters:
//   - log: structured logger for diagnostic output
//   - options: required configuration including serial port path, baud rate, and CAN bitrate
//
// Returns a *usbCANChannel (catalog type, not Interface) because USB-CAN-specific
// callers may need access to USB-CAN-specific functionality.
func newUSBCANChannel(log *slog.Logger, options usbCANChannelOptions) *usbCANChannel {
	c := usbCANChannel{
		options: options,
		log:     log,
		ready:   make(chan struct{}),
	}

	return &c
}

// Run opens the serial port, sends the initial configuration/settings frame to set the CAN bitrate,
// and then enters a blocking read loop that continuously reads bytes from the serial port,
// accumulates them in a buffer, and parses complete USB-CAN frames as they arrive.
//
// The read loop uses a 32-byte working buffer per read call. Since USB-CAN frames can span
// multiple reads (or multiple frames can arrive in a single read), the pending buffer accumulates
// partial data across iterations.
//
// This method blocks indefinitely and only returns on a read error or serial port failure.
func (c *usbCANChannel) Run(ctx context.Context, handler func(can.Frame)) error {
	// Configure and open the serial port at the specified baud rate.
	// The USB-CAN Analyzer typically runs at 2 Mbaud on the serial link.
	mode := &serial.Mode{
		BaudRate: c.options.SerialBaudRate,
	}
	openPort := c.options.openPort
	if openPort == nil {
		openPort = serialio.Open
	}
	port, err := openPort(c.options.SerialPortName, mode)
	if err != nil {
		return err
	}

	c.mu.Lock()
	c.port = port
	c.mu.Unlock()
	defer func() {
		c.mu.Lock()
		if c.port == port {
			c.port = nil
			c.mu.Unlock()
			_ = port.Close()
			return
		}
		c.mu.Unlock()
	}()

	watchDone := make(chan struct{})
	defer close(watchDone)
	go func() {
		select {
		case <-ctx.Done():
			_ = c.Close()
		case <-watchDone:
		}
	}()

	// Put the adapter in normal, extended-frame mode at the NMEA 2000 CAN
	// bitrate before making it available to writers. The serial baud rate above
	// configures only the host-to-adapter link.
	settings := usbCANSettingsFrame()
	c.writeMu.Lock()
	written, writeErr := port.Write(settings[:])
	c.writeMu.Unlock()
	if writeErr != nil {
		return fmt.Errorf("configure USB-CAN adapter: %w", writeErr)
	}
	if written != len(settings) {
		return fmt.Errorf("configure USB-CAN adapter: wrote %d of %d bytes", written, len(settings))
	}
	c.readyOnce.Do(func() { close(c.ready) })

	// Wrap the caller's handler so a nil handler becomes a no-op.
	c.handler = func(frame can.Frame) {
		if handler != nil {
			handler(frame)
		}
	}

	c.log.Info("Opened USBCAN and listening", "portName", c.options.SerialPortName)

	// pending accumulates bytes read from the serial port. Frames are parsed and consumed from the
	// front of this buffer. Any incomplete frame data remains in pending for the next iteration.
	pending := []byte{}
	for {
		// Read up to 32 bytes at a time from the serial port.
		// The actual number of bytes returned depends on what the OS has buffered.
		working := make([]byte, 32)
		readBytes, err := port.Read(working)
		if err != nil {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			return err
		}
		// Append only the bytes actually read (not the full 32-byte buffer) to pending.
		pending = append(pending, working[0:readBytes]...)

		// Attempt to parse as many complete frames as possible from the accumulated buffer.
		// parseFrames will update the pending slice to remove consumed bytes.
		if err := c.parseFrames(&pending); err != nil {
			return err
		}
	}
}

// parseFrames processes the accumulated receive buffer, extracting and dispatching complete USB-CAN
// frames. It modifies the buffer in-place via the pointer, removing consumed bytes from the front.
//
// The parser handles three frame types:
//  1. Command frames: start with 0xAA 0x55, are always 20 bytes long, and are logged but not dispatched.
//     These are typically responses to settings frames or device status messages.
//  2. Data frames: start with 0xAA followed by an info byte with bits[7:6]=0b11. The info byte encodes
//     whether this is a standard/extended frame, remote/data frame, and the data length.
//  3. Error recovery: if the buffer doesn't start with 0xAA, the parser scans forward to find the next
//     0xAA byte and discards the garbage bytes before it.
//
// If a frame is incomplete (not enough bytes in the buffer yet), the function returns nil and leaves
// the partial data in the buffer for the next call.
//
// Parameters:
//   - bufAddr: pointer to the byte slice buffer, modified in-place to remove consumed bytes
func (c *usbCANChannel) parseFrames(bufAddr *[]byte) error {
	for {
		buf := *bufAddr

		// If the buffer is empty, there's nothing to parse.
		if len(buf) == 0 {
			return nil
		}

		// Every valid USB-CAN frame must start with 0xAA (start-of-frame marker).
		// If the first byte is not 0xAA, we have garbage/corruption in the stream.
		if buf[0] != 0xaa {
			// Scan forward to find the next 0xAA byte to resynchronize.
			idx := slices.Index(buf, 0xaa)
			if idx == -1 {
				// No 0xAA found anywhere in the buffer -- discard everything.
				c.log.Debug("Error frame", "data", fmt.Sprintf("%+v", buf))
				*bufAddr = []byte{}
				return nil
			} else {
				// Found 0xAA at position idx -- skip the garbage bytes before it.
				c.log.Debug("Error frame, skipping bytes", "skip", idx, "data", fmt.Sprintf("%+v", buf))
			}

			*bufAddr = buf[idx:]
			continue
		}

		// We need at least 2 bytes to determine the frame type (0xAA + type byte).
		if len(buf) < 2 {
			return nil
		}

		// ----- Command Frame -----
		// Format: 0xAA 0x55 <18 bytes of command/response data> (20 bytes total)
		// These are sent by the device in response to settings frames or as status updates.
		if buf[1] == 0x55 {
			// Command frames are always exactly 20 bytes long.
			// If we don't have all 20 bytes yet, wait for more data.
			if len(buf) < 20 {
				return nil
			}

			c.log.Debug("Command frame", "data", fmt.Sprintf("%+v", buf[0:20]))

			// Consume the 20-byte command frame and continue parsing any remaining data.
			*bufAddr = buf[20:]
			continue
		}

		// ----- Data Frame -----
		// Check if bits [7:6] of the info byte are 0b11 (value 3), indicating a data frame.
		// Info byte layout: [7:6]=frame_type, [5]=extended, [4]=remote, [3:0]=data_length
		if (buf[1] >> 6) == 3 {
			// Start with 2 bytes consumed so far (0xAA + info byte)
			frameLen := byte(2)

			// Bit 5 of the info byte: 1 = extended frame (29-bit CAN ID, 4 ID bytes),
			//                          0 = standard frame (11-bit CAN ID, 2 ID bytes)
			extendedFrame := (buf[1] & 0x20) > 0

			// Bit 4 of the info byte: 1 = remote frame (RTR), 0 = data frame.
			// Remote frames are logged but handled the same way as data frames here.
			remoteFrame := (buf[1] & 0x10) > 0

			// Add the ID field length: 4 bytes for extended, 2 bytes for standard.
			if extendedFrame {
				frameLen += 4
			} else {
				frameLen += 2
			}

			// Bits [3:0] of the info byte encode the number of data bytes (0-8).
			dataLen := buf[1] & 0xf
			if dataLen > 8 {
				c.log.Warn("discarding USB-CAN frame with invalid payload length", "length", dataLen)
				*bufAddr = buf[1:]
				continue
			}

			// Total frame length = header (2) + ID bytes (2 or 4) + data bytes + end byte (0x55)
			frameLen += dataLen + 1
			if len(buf) < int(frameLen) {
				// Not enough bytes yet for the complete frame -- wait for more data.
				return nil
			}

			// Extract the CAN ID from the frame. The ID is stored in little-endian byte order.
			var frameID uint32
			if extendedFrame {
				// Extended frame: 4-byte CAN ID in little-endian (bytes 2-5)
				// This produces a 29-bit extended CAN ID used by protocols like NMEA 2000.
				frameID = (uint32(buf[2])) | (uint32(buf[3]) << 8) | (uint32(buf[4]) << 16) | (uint32(buf[5]) << 24)
			} else {
				// Standard frame: 2-byte CAN ID in little-endian (bytes 2-3)
				// This produces an 11-bit standard CAN ID.
				frameID = (uint32(buf[2])) | (uint32(buf[3]) << 8)
			}

			// The last byte of the frame must be 0x55 (end-of-frame marker).
			endByte := buf[frameLen-1]

			// Extract the data payload bytes (located between the ID and the end byte).
			dataBytes := buf[frameLen-1-dataLen : frameLen-1]
			if endByte != 0x55 {
				// If the end byte is wrong, discard only this candidate start byte.
				// The outer resynchronizer can then recover a valid frame already
				// buffered behind the corrupt candidate.
				c.log.Debug("Data frame with bad end byte",
					"extended", extendedFrame, "remote", remoteFrame,
					"frameID", fmt.Sprintf("%X", frameID),
					"data", fmt.Sprintf("%+v", dataBytes),
					"endByte", fmt.Sprintf("%X", endByte))
				*bufAddr = buf[1:]
				continue
			}
			// Build a can.Frame struct with the parsed data. The can.Frame.Data field is a
			// fixed [8]byte array, so we copy only the actual data bytes into it.
			fData := [8]byte{}
			for i := 0; i < int(dataLen); i++ {
				fData[i] = dataBytes[i]
			}
			fd := can.Frame{
				ID:     frameID,
				Length: dataLen,
				Data:   fData,
			}

			// Dispatch the parsed frame to the user-provided handler callback.
			c.handler(fd)

			// Consume this frame from the buffer and continue parsing.
			*bufAddr = buf[frameLen:]
			continue
		}

		// ----- Unknown Frame Type -----
		// If we get here, the byte after 0xAA is neither 0x55 (command) nor has bits[7:6]=0b11 (data).
		// Discard only this false start so a later valid frame can still be found.
		c.log.Debug("Unknown frame", "data", fmt.Sprintf("%+v", buf))
		*bufAddr = buf[1:]
		continue
	}
}

// usbCANSettingsFrame selects 250 kbit/s CAN, 29-bit extended frames, and
// normal (non-loopback, non-silent) operation. The final byte is the protocol
// checksum: the low byte of the sum of bytes 2 through 18.
func usbCANSettingsFrame() [20]byte {
	settings := [20]byte{
		0xAA, 0x55, 0x12,
		0x05,                   // 250 kbit/s CAN bitrate
		0x02,                   // extended CAN frames
		0x00, 0x00, 0x00, 0x00, // acceptance code
		0x00, 0x00, 0x00, 0x00, // acceptance mask
		0x00, // normal mode
		0x01, // apply settings
		0x00, 0x00, 0x00, 0x00,
	}
	for _, b := range settings[2:19] {
		settings[19] += b
	}
	return settings
}

// Close shuts down the USB-CAN channel by closing the underlying serial port.
// It is safe to call Close() even if the port was never opened (port is nil).
func (c *usbCANChannel) Close() error {
	c.mu.Lock()
	port := c.port
	c.port = nil
	c.mu.Unlock()
	if port != nil {
		return port.Close()
	}
	return nil
}

// WriteFrame sends a CAN frame out through the USB-CAN dongle by encoding it into the
// USB-CAN binary protocol format and writing it to the serial port.
//
// Outgoing data frame format:
//   - Byte 0: 0xAA (start-of-frame marker)
//   - Byte 1: 0xE0 | dataLen (data frame + 29-bit extended-ID flag + length)
//   - Bytes 2-5: CAN ID in little-endian byte order
//   - Next N bytes: data payload (N = frame.Length)
//   - Final byte: 0x55 (end-of-frame marker)
func (c *usbCANChannel) WriteFrame(frame can.Frame) error {
	if frame.Length > 8 {
		return fmt.Errorf("WriteFrame: invalid classical CAN payload length %d", frame.Length)
	}
	if frame.ID > 0x1FFFFFFF {
		return fmt.Errorf("WriteFrame: invalid 29-bit CAN ID 0x%X", frame.ID)
	}
	// NMEA 2000 always uses a four-byte, 29-bit extended CAN identifier.
	buf := []byte{
		0xaa,
		0xE0 | frame.Length,
		byte(frame.ID),
		byte(frame.ID >> 8),
		byte(frame.ID >> 16),
		byte(frame.ID >> 24),
	}
	// Append the data payload bytes and the end-of-frame marker.
	buf = append(buf, frame.Data[0:frame.Length]...)
	buf = append(buf, 0x55)

	c.mu.Lock()
	port := c.port
	c.mu.Unlock()
	if port == nil {
		return fmt.Errorf("WriteFrame: USB-CAN port is not open")
	}
	c.writeMu.Lock()
	o, err := port.Write(buf)
	c.writeMu.Unlock()
	if o != len(buf) {
		return fmt.Errorf("WriteFrame sent %d of %d bytes", o, len(buf))
	}
	if err != nil {
		return err
	}

	return nil
}

// Ready closes once Run has opened the serial adapter for writes.
func (c *usbCANChannel) Ready() <-chan struct{} { return c.ready }

// NewUSBCAN creates a USB-CAN channel for the given serial port.
// The channel is not opened, and no handler is registered, until Run() is called.
func NewUSBCAN(log *slog.Logger, port string) *usbCANChannel {
	return newUSBCANChannel(log, usbCANChannelOptions{
		SerialPortName: port,
		SerialBaudRate: 2000000,
	})
}

// RunUSBCAN creates a USB-CAN channel for the given serial port and runs it,
// calling handler for each received CAN frame. Blocks until error or context done.
func RunUSBCAN(ctx context.Context, log *slog.Logger, port string, handler func(can.Frame)) error {
	ch := newUSBCANChannel(log, usbCANChannelOptions{
		SerialPortName: port,
		SerialBaudRate: 2000000,
	})

	return ch.Run(ctx, handler)
}
