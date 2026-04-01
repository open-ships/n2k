package n2k

import (
	"context"
	"log/slog"
	"sync"

	"github.com/brutella/can"
	"github.com/open-ships/n2k/pkg/n2k/internal/canbus"
)

type source interface {
	run(ctx context.Context, log *slog.Logger, handler func(can.Frame)) error
}

type socketCANSource struct {
	iface string
}

func (s *socketCANSource) run(ctx context.Context, log *slog.Logger, handler func(can.Frame)) error {
	return canbus.RunSocketCAN(ctx, log, s.iface, handler)
}

type usbCANSource struct {
	port string
}

func (s *usbCANSource) run(ctx context.Context, log *slog.Logger, handler func(can.Frame)) error {
	return canbus.RunUSBCAN(ctx, log, s.port, handler)
}

type replaySource struct {
	frames []can.Frame
}

func (s *replaySource) run(ctx context.Context, log *slog.Logger, handler func(can.Frame)) error {
	for _, f := range s.frames {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			handler(f)
		}
	}
	return nil
}

func runSources(ctx context.Context, log *slog.Logger, sources []source, handler func(can.Frame)) error {
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

	safeHandler := func(f can.Frame) {
		mu.Lock()
		defer mu.Unlock()
		handler(f)
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
