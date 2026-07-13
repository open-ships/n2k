package n2k

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/brutella/can"
	"github.com/open-ships/n2k/internal/candump"
)

// fileSource replays CAN frames from a candump -L / -l log file.
type fileSource struct {
	path           string
	originalTiming bool
}

// FileOption configures a File source.
type FileOption interface {
	applyFile(*fileSource)
}

type fileOptionFunc func(*fileSource)

func (f fileOptionFunc) applyFile(s *fileSource) { f(s) }

// OriginalTiming replays frames paced by the log's own timestamps instead of
// as fast as they can be read. Frames without timestamps are delivered
// immediately.
func OriginalTiming() FileOption {
	return fileOptionFunc(func(s *fileSource) {
		s.originalTiming = true
	})
}

func (s *fileSource) run(ctx context.Context, log *slog.Logger, handler func(can.Frame)) error {
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

		frame, ts, ok := candump.Parse(scanner.Text())
		if !ok {
			continue
		}

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

		handler(frame)
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("n2k: reading log file: %w", err)
	}
	return nil
}
