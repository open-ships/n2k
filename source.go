package n2k

import (
	"context"
	"log/slog"
	"sync"

	"github.com/brutella/can"
	"github.com/open-ships/n2k/internal/canbus"
	"github.com/open-ships/n2k/raw"
)

type source interface {
	run(ctx context.Context, log *slog.Logger, handler func(raw.Observation)) error
}

// busBacked is implemented by sources that are backed by real hardware and
// can be opened as a read/write Bus (not just a frame replay).
type busBacked interface {
	newBus(log *slog.Logger) Bus
}

type socketCANSource struct {
	iface string
}

func (s *socketCANSource) run(ctx context.Context, log *slog.Logger, handler func(raw.Observation)) error {
	return canbus.RunSocketCAN(ctx, log, s.iface, func(frame can.Frame) {
		handler(frameObservation(frame, "socketcan:"+s.iface, s.iface, raw.DirectionReceived))
	})
}

func (s *socketCANSource) newBus(log *slog.Logger) Bus { return canbus.NewSocketCAN(log, s.iface) }

type usbCANSource struct {
	port string
}

func (s *usbCANSource) run(ctx context.Context, log *slog.Logger, handler func(raw.Observation)) error {
	return canbus.RunUSBCAN(ctx, log, s.port, func(frame can.Frame) {
		handler(frameObservation(frame, "usbcan:"+s.port, s.port, raw.DirectionReceived))
	})
}

func (s *usbCANSource) newBus(log *slog.Logger) Bus { return canbus.NewUSBCAN(log, s.port) }

type replaySource struct {
	observations []raw.Observation
}

func (s *replaySource) run(ctx context.Context, log *slog.Logger, handler func(raw.Observation)) error {
	for _, observation := range s.observations {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			handler(observation.Clone())
		}
	}
	return nil
}

func runSources(ctx context.Context, log *slog.Logger, sources []source, handler func(raw.Observation)) error {
	if len(sources) == 1 {
		return sources[0].run(ctx, log, handler)
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var (
		mu sync.Mutex
		wg sync.WaitGroup
	)
	errc := make(chan error, len(sources))

	safeHandler := func(observation raw.Observation) {
		mu.Lock()
		defer mu.Unlock()
		handler(observation)
	}

	for _, src := range sources {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := src.run(ctx, log, safeHandler); err != nil {
				cancel()
				errc <- err
			}
		}()
	}

	wg.Wait()
	close(errc)

	for err := range errc {
		if err != nil {
			return err
		}
	}
	return nil
}
