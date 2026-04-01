package n2k

import (
	"context"
	"log/slog"
	"sync"

	"github.com/brutella/can"
	"github.com/open-ships/n2k/internal/adapter"
	"github.com/open-ships/n2k/internal/decoder"
	"github.com/open-ships/n2k/pgn"
)

// Scanner reads decoded NMEA 2000 messages one at a time.
// Call Next() to advance, Message() to get the current message, and Err() for errors.
type Scanner struct {
	ctx  context.Context
	cfg  config
	msg  any
	err  error
	ch   chan any
	once sync.Once
}

// NewScanner creates a Scanner that reads from the configured sources.
func NewScanner(ctx context.Context, opts ...Option) *Scanner {
	cfg := config{}
	for _, o := range opts {
		o.apply(&cfg)
	}
	if cfg.logger == nil {
		cfg.logger = slog.Default()
	}

	return &Scanner{
		ctx: ctx,
		cfg: cfg,
		ch:  make(chan any, 64),
	}
}

// Next advances the scanner to the next message. Returns false when no more messages
// are available (source exhausted or error occurred). Check Err() after Next returns false.
func (s *Scanner) Next() bool {
	s.once.Do(func() {
		if err := s.cfg.validate(); err != nil {
			s.err = err
			close(s.ch)
			return
		}
		go s.run()
	})

	msg, ok := <-s.ch
	if !ok {
		return false
	}
	s.msg = msg
	return true
}

// Message returns the most recently scanned message.
func (s *Scanner) Message() any {
	return s.msg
}

// Err returns the first error encountered by the scanner.
func (s *Scanner) Err() error {
	return s.err
}

func (s *Scanner) run() {
	defer close(s.ch)

	a := adapter.NewCANAdapter()
	dec := decoder.New()

	dec.SetOutput(&scannerHandler{scanner: s})
	a.SetOutput(dec)

	err := runSources(s.ctx, s.cfg.logger, s.cfg.sources, func(f can.Frame) {
		a.HandleMessage(&f)
	})
	if err != nil {
		s.err = err
	}
}

type scannerHandler struct {
	scanner *Scanner
}

func (h *scannerHandler) HandleStruct(msg any) {
	if msg == nil {
		return
	}

	// Normalize value-type UnknownPGN to pointer for consistent downstream handling.
	if u, ok := msg.(pgn.UnknownPGN); ok {
		msg = &u
	}

	// Drop unknown PGNs unless IncludeUnknown is set
	if u, ok := msg.(*pgn.UnknownPGN); ok {
		if !h.scanner.cfg.includeUnknown {
			h.scanner.cfg.logger.Debug("dropping unknown PGN",
				"pgn", u.Info.PGN,
				"reason", u.Reason,
			)
			return
		}
	}

	select {
	case h.scanner.ch <- msg:
	case <-h.scanner.ctx.Done():
	}
}
