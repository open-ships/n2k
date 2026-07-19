package n2k

import (
	"context"
	"errors"
	"log/slog"
	"sync"

	"github.com/open-ships/n2k/pgn"
)

// Scanner reads decoded NMEA 2000 messages one at a time.
// Call Next() to advance, Message() to get the current message, and Err() for errors.
type Scanner struct {
	ctx    context.Context
	cancel context.CancelFunc
	cfg    config
	msg    pgn.Message
	err    error
	ch     chan pgn.Message
	once   sync.Once
	sub    *messageSubscription
}

// NewScanner creates a Scanner that reads from the configured sources.
func NewScanner(ctx context.Context, opts ...Option) *Scanner {
	cfg := config{}
	var optionErr error
	for _, o := range opts {
		if o == nil {
			optionErr = errors.New("n2k: nil Option")
			break
		}
		o.apply(&cfg)
	}
	cfg.applyReconnect()
	if cfg.logger == nil {
		cfg.logger = slog.Default()
	}

	scannerCtx, cancel := context.WithCancel(ctx)
	receiveBuffer := defaultReceiveBuffer
	if cfg.receiveBuffer != nil && *cfg.receiveBuffer > 0 {
		receiveBuffer = *cfg.receiveBuffer
	}
	s := &Scanner{
		ctx:    scannerCtx,
		cancel: cancel,
		cfg:    cfg,
		ch:     make(chan pgn.Message, receiveBuffer),
	}
	// Validate eagerly so a misconfigured Scanner (e.g. no sources) fails on
	// the very first Next() call rather than only after the goroutine starts.
	s.err = optionErr
	if s.err == nil {
		s.err = cfg.validate()
	}
	return s
}

// Next advances the scanner to the next message. Returns false when no more messages
// are available (source exhausted or error occurred). Check Err() after Next returns false.
func (s *Scanner) Next() bool {
	if s.sub != nil {
		msg, ok := <-s.sub.ch
		if !ok {
			s.err = s.sub.terminalError()
			return false
		}
		s.msg = msg
		return true
	}
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

// Close stops the scanner and releases its source or live-client
// subscription. It is safe to call more than once.
func (s *Scanner) Close() error {
	if s.cancel != nil {
		s.cancel()
	}
	if s.sub != nil {
		s.sub.unsubscribe()
	}
	return nil
}

func (s *Scanner) run() {
	defer close(s.ch)
	p, err := newReadPipeline(s.ctx, s.cfg, channelEmitter(s.ctx, s.ch))
	if err != nil {
		s.err = err
		return
	}
	if err := runSources(s.ctx, s.cfg.logger, s.cfg.sources, p.HandleFrame); err != nil && s.ctx.Err() == nil {
		s.err = err
	}
}
