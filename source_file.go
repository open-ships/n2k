package n2k

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/open-ships/n2k/internal/candump"
	"github.com/open-ships/n2k/raw"
)

// fileSource replays CAN frames from a candump -L / -l log file.
type fileSource struct {
	path           string
	originalTiming bool
}

type captureOptions struct {
	originalTiming bool
}

// FileOption configures a candump File or EBL capture source.
type FileOption interface {
	applyCapture(*captureOptions)
}

type fileOptionFunc func(*captureOptions)

func (f fileOptionFunc) applyCapture(options *captureOptions) { f(options) }

// OriginalTiming replays frames or messages paced by the capture's own
// timestamps instead of as fast as they can be read. Records without
// timestamps are delivered immediately.
func OriginalTiming() FileOption {
	return fileOptionFunc(func(options *captureOptions) {
		options.originalTiming = true
	})
}

func (s *fileSource) run(ctx context.Context, log *slog.Logger, handler func(raw.Observation)) error {
	f, err := os.Open(s.path) // #nosec G304 -- reading the user-supplied capture log is this source's purpose.
	if err != nil {
		return fmt.Errorf("n2k: opening log file: %w", err)
	}
	defer func() { _ = f.Close() }()

	var lastTS time.Time
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		record, ok := candump.ParseRecord(scanner.Text())
		if !ok {
			continue
		}
		frame, ts := record.Frame, record.Timestamp

		if s.originalTiming && !ts.IsZero() && !lastTS.IsZero() {
			if delta := ts.Sub(lastTS); delta > 0 {
				timer := time.NewTimer(delta)
				select {
				case <-timer.C:
				case <-ctx.Done():
					timer.Stop()
					return ctx.Err()
				}
			}
		}
		if !ts.IsZero() {
			lastTS = ts
		}

		now := time.Now()
		timestamp := ts
		if timestamp.IsZero() {
			timestamp = now
		}
		networkID := record.Network
		if networkID == "" {
			networkID = s.path
		}
		handler(normalizeObservation(raw.Observation{
			Kind:       raw.KindFrame,
			Timestamp:  timestamp,
			ReceivedAt: now,
			AdapterID:  "file:" + s.path,
			NetworkID:  networkID,
			Direction:  raw.DirectionReceived,
			Frame:      &frame,
		}))
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("n2k: reading log file: %w", err)
	}
	return nil
}
