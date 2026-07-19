package n2k

import (
	"context"
	"log/slog"

	"github.com/brutella/can"
	"github.com/open-ships/n2k/internal/adapter"
	"github.com/open-ships/n2k/internal/decoder"
	"github.com/open-ships/n2k/pgn"
)

// infoCarrier is a local interface for types that expose MessageInfo through
// a typed method (eliminating the need for reflection). Both pgn.PGN and
// pgn.UnknownPGN satisfy this interface.
type infoCarrier interface {
	MessageInfo() pgn.MessageInfo
}

// readPipeline is the single decode path shared by Scanner, Receive, and
// Client: raw CAN frame -> pre-filter -> fast-packet assembly -> decode ->
// unknown-PGN policy -> post-filter -> out channel.
type readPipeline struct {
	ctx            context.Context
	log            *slog.Logger
	includeUnknown bool
	filter         *filter
	adapter        *adapter.CANAdapter
	decoder        *decoder.Decoder
	emit           func(pgn.Message)
}

// newReadPipeline compiles the filter eagerly and wires adapter -> decoder -> output stage.
// The returned pipeline delivers messages on out. It never closes out; the
// caller that owns the frame source closes it when the source ends.
func newReadPipeline(ctx context.Context, cfg config, emit func(pgn.Message)) (*readPipeline, error) {
	log := cfg.logger
	if log == nil {
		log = slog.Default()
	}
	var f *filter
	if cfg.filterExpr != "" {
		var err error
		f, err = compileFilter(cfg.filterExpr)
		if err != nil {
			return nil, err
		}
	}
	p := &readPipeline{
		ctx:            ctx,
		log:            log,
		includeUnknown: cfg.includeUnknown,
		filter:         f,
		adapter:        adapter.NewCANAdapter(),
		decoder:        decoder.New(),
		emit:           emit,
	}
	p.decoder.SetOutput(p)
	p.adapter.SetOutput(p.decoder)
	return p, nil
}

// HandleFrame is the single entry point for raw CAN frames (pre-filter + assembly + decode).
func (p *readPipeline) HandleFrame(frame can.Frame) {
	if p.filter != nil {
		info := adapter.NewPacketInfo(&frame)
		if !p.filter.evalPre(info) {
			return
		}
	}
	p.adapter.HandleMessage(&frame)
}

// InjectAssembled feeds an already-assembled payload (e.g. an ISO-TP reassembly)
// through candidate filtering and decode. Replaces the decoder-internals knowledge
// that previously lived in Client's transport OnComplete closure.
func (p *readPipeline) InjectAssembled(info pgn.MessageInfo, data []byte) {
	if p.filter != nil && !p.filter.evalPre(info) {
		return
	}
	packet := decoder.NewPacket(info, data)
	packet.Complete = true
	packet.FilterCandidates()
	p.decoder.Decode(*packet)
}

// HandleStruct implements decoder.Handler: unknown-PGN policy, post-filter,
// and delivery.
func (p *readPipeline) HandleStruct(msg pgn.Message) {
	if msg == nil {
		return
	}
	if u, ok := msg.(*pgn.UnknownPGN); ok && !p.includeUnknown {
		p.log.Debug("dropping unknown PGN", "pgn", u.Info.PGN, "reason", u.Reason)
		return
	}
	if p.filter != nil && p.filter.hasPost {
		carrier, ok := msg.(infoCarrier)
		if !ok {
			return
		}
		if !p.filter.evalPostWithInfo(carrier.MessageInfo(), structToFilterMap(msg)) {
			return
		}
	}
	if p.emit != nil && p.ctx.Err() == nil {
		p.emit(msg)
	}
}

func channelEmitter(ctx context.Context, out chan<- pgn.Message) func(pgn.Message) {
	return func(msg pgn.Message) {
		select {
		case out <- msg:
		case <-ctx.Done():
		}
	}
}
