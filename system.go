package n2k

import (
	"context"
	"errors"
	"log/slog"
	"runtime/debug"
	"sync"

	"github.com/open-ships/n2k/pgn"
	"github.com/open-ships/n2k/raw"
)

// systemPGNs are always decoded by the system router: product information
// (126996), configuration information (126998), and group functions (126208).
var systemPGNs = []uint32{126996, 126998, 126208, 126720}

// systemRouter is the client's protocol-message path. It decodes the PGNs the
// client itself must react to (product info requests, group functions,
// request/response correlation, device registry updates) through a dedicated
// read pipeline that is independent of any user-configured filter, so user
// filters can never break protocol behavior.
type systemRouter struct {
	ctx        context.Context
	log        *slog.Logger
	pipeline   *readPipeline
	ch         chan pgn.Message
	done       chan struct{}
	dispatchMu *sync.Mutex
	current    func(pgn.MessageInfo) bool
	onError    func(error)

	mu                  sync.Mutex
	accept              map[uint32]int // PGN → reference count
	handlers            []func(pgn.Message)
	observationHandlers []func(raw.Observation)
}

// newSystemRouter builds the router and its dedicated pipeline. Call run in a
// goroutine to start dispatching decoded messages to handlers.
func newSystemRouter(ctx context.Context, cfg config) (*systemRouter, error) {
	// The system pipeline must see every protocol message: no user filter,
	// no unknown-PGN passthrough.
	sysCfg := cfg
	sysCfg.filterExpr = ""
	sysCfg.includeUnknown = false

	ch := make(chan pgn.Message, 64)
	p, err := newReadPipeline(ctx, sysCfg, channelEmitter(ctx, ch))
	if err != nil {
		return nil, err
	}

	r := &systemRouter{
		ctx:        ctx,
		log:        cfg.logger,
		pipeline:   p,
		ch:         ch,
		accept:     make(map[uint32]int),
		done:       make(chan struct{}),
		dispatchMu: &sync.Mutex{},
	}
	p.emit = func(msg pgn.Message) {
		select {
		case ch <- msg:
		default:
			if r.onError != nil {
				r.onError(errors.New("n2k: system protocol receive queue full"))
			}
		}
	}
	for _, n := range systemPGNs {
		r.accept[n] = 1
	}
	p.setObservationOutput(r.dispatchObservation)
	return r, nil
}

func (r *systemRouter) handleObservation(observation raw.Observation, pgnNum uint32) {
	r.mu.Lock()
	_, ok := r.accept[pgnNum]
	r.mu.Unlock()
	if ok {
		if observation.Kind == raw.KindMessage {
			// Already-assembled source records bypass readPipeline.Decode, whose
			// observation hook normally runs between assembly and typed decode.
			r.dispatchObservation(observation)
		}
		r.pipeline.HandleObservation(observation)
	}
}

// handleAssembled feeds an already-assembled payload (ISO-TP reassembly)
// through the same gate.
func (r *systemRouter) handleAssembled(info pgn.MessageInfo, data []byte) {
	r.mu.Lock()
	_, ok := r.accept[info.PGN]
	r.mu.Unlock()
	if ok {
		r.pipeline.InjectAssembled(info, data)
	}
}

// addPGN adds a PGN to the accept set (reference counted, so concurrent
// requests for the same PGN compose).
func (r *systemRouter) addPGN(pgnNum uint32) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.accept[pgnNum]++
}

// removePGN drops one reference to a PGN, removing it from the accept set
// when the last reference is gone.
func (r *systemRouter) removePGN(pgnNum uint32) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if n, ok := r.accept[pgnNum]; ok {
		if n <= 1 {
			delete(r.accept, pgnNum)
		} else {
			r.accept[pgnNum] = n - 1
		}
	}
}

// addHandler registers a dispatch target for decoded system messages.
// Handlers run sequentially on the dispatch goroutine, in registration order.
func (r *systemRouter) addHandler(h func(pgn.Message)) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.handlers = append(r.handlers, h)
}

// addObservationHandler registers for assembled, owned payloads before typed
// generated decoding. Manufacturer envelope Adapters use this Seam so decoder
// candidate ordering cannot hide a proprietary response.
func (r *systemRouter) addObservationHandler(handler func(raw.Observation)) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.observationHandlers = append(r.observationHandlers, handler)
}

func (r *systemRouter) dispatchObservation(observation raw.Observation) {
	r.mu.Lock()
	handlers := append([]func(raw.Observation){}, r.observationHandlers...)
	r.mu.Unlock()
	for _, handler := range handlers {
		func() {
			defer func() {
				if recovered := recover(); recovered != nil {
					r.log.Error("recovered panic in system observation handler", "panic", recovered, "stack", string(debug.Stack()))
				}
			}()
			handler(observation.Clone())
		}()
	}
}

// run drains the pipeline output and dispatches each message to every
// registered handler. It exits when the context is canceled.
func (r *systemRouter) run() {
	defer close(r.done)
	for {
		select {
		case msg := <-r.ch:
			r.deliver(msg)
		case <-r.ctx.Done():
			return
		}
	}
}

func (r *systemRouter) deliver(msg pgn.Message) {
	r.dispatchMu.Lock()
	defer r.dispatchMu.Unlock()
	if carrier, ok := msg.(infoCarrier); ok && r.current != nil && !r.current(carrier.MessageInfo()) {
		return
	}
	r.mu.Lock()
	handlers := append([]func(pgn.Message){}, r.handlers...)
	r.mu.Unlock()
	for _, handler := range handlers {
		owned, err := pgn.CloneMessage(msg)
		if err != nil {
			if r.onError != nil {
				r.onError(err)
			}
			return
		}
		r.dispatch(handler, owned)
	}
}

// dispatch delivers one message to one handler, recovering a panic so a faulty
// handler is logged and skipped rather than killing the dispatch loop — which
// runs on an n2k-owned goroutine, so an escaping panic would crash the process.
func (r *systemRouter) dispatch(h func(pgn.Message), msg pgn.Message) {
	defer func() {
		if rec := recover(); rec != nil {
			r.log.Error("recovered panic in system handler",
				"panic", rec, "stack", string(debug.Stack()))
		}
	}()
	h(msg)
}
