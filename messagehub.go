package n2k

import (
	"errors"
	"sync"

	"github.com/open-ships/n2k/pgn"
)

// ErrReceiveOverflow reports that a live Client subscriber could not keep up
// with bus traffic. The subscriber is closed rather than allowing application
// backpressure to stall address claiming and other protocol processing.
var ErrReceiveOverflow = errors.New("n2k: receive buffer overflow")

const defaultReceiveBuffer = 64

// messageHub fans live messages out to independent subscriptions. Publishing
// never blocks: a slow subscription is failed explicitly while other
// subscribers and the bus protocol engine continue making progress.
type messageHub struct {
	mu          sync.Mutex
	buffer      int
	subscribers map[*messageSubscription]struct{}
	backlog     []pgn.Message
	missed      bool
	closed      bool
	terminalErr error
}

type messageSubscription struct {
	hub  *messageHub
	ch   chan pgn.Message
	once sync.Once

	mu  sync.Mutex
	err error
}

func newMessageHub(buffer int) *messageHub {
	if buffer <= 0 {
		buffer = defaultReceiveBuffer
	}
	return &messageHub{
		buffer:      buffer,
		subscribers: make(map[*messageSubscription]struct{}),
	}
}

func (h *messageHub) publish(msg pgn.Message) {
	if msg == nil {
		return
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.closed {
		return
	}
	if len(h.subscribers) == 0 {
		if len(h.backlog) == h.buffer {
			copy(h.backlog, h.backlog[1:])
			h.backlog[len(h.backlog)-1] = msg
			h.missed = true
			return
		}
		h.backlog = append(h.backlog, msg)
		return
	}
	for sub := range h.subscribers {
		select {
		case sub.ch <- msg:
		default:
			sub.failLocked(ErrReceiveOverflow)
			delete(h.subscribers, sub)
		}
	}
}

func (h *messageHub) subscribe() *messageSubscription {
	h.mu.Lock()
	defer h.mu.Unlock()
	sub := &messageSubscription{hub: h, ch: make(chan pgn.Message, h.buffer)}
	for _, msg := range h.backlog {
		sub.ch <- msg
	}
	h.backlog = nil
	if h.missed {
		h.missed = false
		sub.failLocked(ErrReceiveOverflow)
		return sub
	}
	if h.closed {
		sub.failLocked(h.terminalErr)
		return sub
	}
	h.subscribers[sub] = struct{}{}
	return sub
}

func (h *messageHub) close(err error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.closed {
		return
	}
	h.closed = true
	h.terminalErr = err
	for sub := range h.subscribers {
		sub.failLocked(err)
		delete(h.subscribers, sub)
	}
}

func (h *messageHub) subscriberCount() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return len(h.subscribers)
}

func (s *messageSubscription) failLocked(err error) {
	s.mu.Lock()
	s.err = err
	s.mu.Unlock()
	s.once.Do(func() { close(s.ch) })
}

func (s *messageSubscription) unsubscribe() {
	if s == nil || s.hub == nil {
		return
	}
	s.hub.mu.Lock()
	delete(s.hub.subscribers, s)
	s.once.Do(func() { close(s.ch) })
	s.hub.mu.Unlock()
}

func (s *messageSubscription) terminalError() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.err
}
