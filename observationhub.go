package n2k

import (
	"errors"
	"sync"

	"github.com/open-ships/n2k/raw"
)

// ErrObservationOverflow reports that an observation subscriber could not
// keep up. Only that subscriber is closed; protocol processing continues.
var ErrObservationOverflow = errors.New("n2k: observation buffer overflow")

type observationHub struct {
	mu          sync.Mutex
	buffer      int
	subscribers map[*observationSubscription]struct{}
	backlog     []raw.Observation
	missed      bool
	closed      bool
	terminalErr error
}

type observationSubscription struct {
	hub  *observationHub
	ch   chan raw.Observation
	once sync.Once

	mu  sync.Mutex
	err error
}

func newObservationHub(buffer int) *observationHub {
	if buffer <= 0 {
		buffer = defaultReceiveBuffer
	}
	return &observationHub{
		buffer:      buffer,
		subscribers: make(map[*observationSubscription]struct{}),
	}
}

func (h *observationHub) publish(observation raw.Observation) {
	observation = observation.Clone()
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.closed {
		return
	}
	if len(h.subscribers) == 0 {
		if len(h.backlog) == h.buffer {
			copy(h.backlog, h.backlog[1:])
			h.backlog[len(h.backlog)-1] = observation
			h.missed = true
			return
		}
		h.backlog = append(h.backlog, observation)
		return
	}
	for sub := range h.subscribers {
		select {
		case sub.ch <- observation.Clone():
		default:
			sub.failLocked(ErrObservationOverflow)
			delete(h.subscribers, sub)
		}
	}
}

func (h *observationHub) subscribe() *observationSubscription {
	h.mu.Lock()
	defer h.mu.Unlock()
	sub := &observationSubscription{hub: h, ch: make(chan raw.Observation, h.buffer)}
	for _, observation := range h.backlog {
		sub.ch <- observation.Clone()
	}
	h.backlog = nil
	if h.missed {
		h.missed = false
		sub.failLocked(ErrObservationOverflow)
		return sub
	}
	if h.closed {
		sub.failLocked(h.terminalErr)
		return sub
	}
	h.subscribers[sub] = struct{}{}
	return sub
}

func (h *observationHub) close(err error) {
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

func (h *observationHub) subscriberCount() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return len(h.subscribers)
}

func (s *observationSubscription) failLocked(err error) {
	s.mu.Lock()
	s.err = err
	s.mu.Unlock()
	s.once.Do(func() { close(s.ch) })
}

func (s *observationSubscription) unsubscribe() {
	if s == nil || s.hub == nil {
		return
	}
	s.hub.mu.Lock()
	delete(s.hub.subscribers, s)
	s.once.Do(func() { close(s.ch) })
	s.hub.mu.Unlock()
}

func (s *observationSubscription) terminalError() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.err
}
