package n2k

import (
	"context"
	"log/slog"
	"sync"

	"github.com/open-ships/n2k/pgn"
)

// Scanner reads decoded NMEA 2000 messages one at a time.
// Call Next() to advance, Message() to get the current message, and Err() for errors.
type Scanner struct {
	ctx  context.Context
	cfg  config
	msg  pgn.Message
	err  error
	ch   chan pgn.Message
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

	s := &Scanner{
		ctx: ctx,
		cfg: cfg,
		ch:  make(chan pgn.Message, 64),
	}
	// Validate eagerly so a misconfigured Scanner (e.g. no sources) fails on
	// the very first Next() call rather than only after the goroutine starts.
	s.err = cfg.validate()
	return s
}

// Next advances the scanner to the next message. Returns false when no more messages
// are available (source exhausted or error occurred). Check Err() after Next returns false.
func (s *Scanner) Next() bool {
	s.once.Do(func() {
		if s.err != nil {
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
func (s *Scanner) Message() pgn.Message {
	return s.msg
}

// Err returns the first error encountered by the scanner.
func (s *Scanner) Err() error {
	return s.err
}

func (s *Scanner) run() {
	defer close(s.ch)
	p, err := newReadPipeline(s.ctx, s.cfg, s.ch)
	if err != nil {
		s.err = err
		return
	}
	if err := runSources(s.ctx, s.cfg.logger, s.cfg.sources, p.HandleFrame); err != nil {
		s.err = err
	}
}
