package n2k

import (
	"context"
	"log/slog"
	"runtime/debug"
	"sync"

	"github.com/brutella/can"
	"github.com/open-ships/n2k/pgn"
)

// systemPGNs are always decoded by the system router: product information
// (126996), configuration information (126998), and group functions (126208).
var systemPGNs = []uint32{126996, 126998, 126208}

// systemRouter is the client's protocol-message path. It decodes the PGNs the
// client itself must react to (product info requests, group functions,
// request/response correlation, device registry updates) through a dedicated
// read pipeline that is independent of any user-configured filter, so user
// filters can never break protocol behavior.
type systemRouter struct {
	ctx      context.Context
	log      *slog.Logger
	pipeline *readPipeline
	ch       chan pgn.Message

	mu       sync.Mutex
	accept   map[uint32]int // PGN → reference count
	handlers []func(pgn.Message)
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
		ctx:      ctx,
		log:      cfg.logger,
		pipeline: p,
		ch:       ch,
		accept:   make(map[uint32]int),
	}
	for _, n := range systemPGNs {
		r.accept[n] = 1
	}
	return r, nil
}

// handleFrame feeds one raw CAN frame into the system pipeline when its PGN
// is in the accept set.
func (r *systemRouter) handleFrame(frame can.Frame, pgnNum uint32) {
	r.mu.Lock()
	_, ok := r.accept[pgnNum]
	r.mu.Unlock()
	if ok {
		r.pipeline.HandleFrame(frame)
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

// run drains the pipeline output and dispatches each message to every
// registered handler. It exits when the context is canceled.
func (r *systemRouter) run() {
	for {
		select {
		case msg := <-r.ch:
			r.mu.Lock()
			handlers := make([]func(pgn.Message), len(r.handlers))
			copy(handlers, r.handlers)
			r.mu.Unlock()
			for _, h := range handlers {
				r.dispatch(h, msg)
			}
		case <-r.ctx.Done():
			return
		}
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
