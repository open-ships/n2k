package n2k

import (
	"context"
	"log/slog"
	"time"

	"github.com/brutella/can"
	"github.com/open-ships/n2k/internal/adapter"
	"github.com/open-ships/n2k/internal/decoder"
	"github.com/open-ships/n2k/pgn"
	"github.com/open-ships/n2k/raw"
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
	observe        func(raw.Observation)
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
	p.adapter.SetOutput(p)
	return p, nil
}

// HandleFrame is the single entry point for raw CAN frames (pre-filter + assembly + decode).
func (p *readPipeline) HandleFrame(frame can.Frame) {
	now := time.Now()
	p.HandleObservation(normalizeObservation(raw.Observation{
		Kind:       raw.KindFrame,
		Timestamp:  now,
		ReceivedAt: now,
		Direction:  raw.DirectionReceived,
		Frame:      &frame,
	}))
}

// HandleObservation is the source-aware pipeline entry point.
func (p *readPipeline) HandleObservation(observation raw.Observation) {
	observation = normalizeObservation(observation)
	if observation.Kind == raw.KindMessage {
		info := messageInfoForObservation(observation)
		if p.filter != nil && !p.filter.evalPre(info) {
			return
		}
		packet := decoder.NewPacket(info, observation.Payload)
		packet.Complete = true
		packet.FilterCandidates()
		p.decoder.Decode(*packet)
		return
	}
	if observation.Kind != raw.KindFrame || observation.Frame == nil {
		return
	}
	frame := *observation.Frame
	info := messageInfoForObservation(observation)
	if p.filter != nil {
		if !p.filter.evalPre(info) {
			return
		}
	}
	p.adapter.HandleMessageWithInfo(&frame, info)
}

// Decode implements adapter.PacketHandler. It exposes the assembled owned
// payload before typed decoding, then delegates to the decoder.
func (p *readPipeline) Decode(packet decoder.Packet) {
	if p.observe != nil {
		p.observe(observationFromPacket(raw.KindMessage, packet, ""))
	}
	p.decoder.Decode(packet)
}

func (p *readPipeline) setObservationOutput(observe func(raw.Observation)) {
	p.observe = observe
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
	p.Decode(*packet)
}

// HandleStruct implements decoder.Handler: unknown-PGN policy, post-filter,
// and delivery.
func (p *readPipeline) HandleStruct(msg pgn.Message) {
	if msg == nil {
		return
	}
	if u, ok := msg.(*pgn.UnknownPGN); ok && !p.includeUnknown {
		if p.observe != nil {
			packet := decoder.Packet{Info: u.Info, Data: append([]byte(nil), u.Data...)}
			reason := "unknown PGN"
			if u.Reason != nil {
				reason = u.Reason.Error()
			}
			p.observe(observationFromPacket(raw.KindDecodeError, packet, reason))
		}
		p.log.Debug("dropping unknown PGN", "pgn", u.Info.PGN, "reason", u.Reason)
		return
	}
	if u, ok := msg.(*pgn.UnknownPGN); ok && p.observe != nil {
		packet := decoder.Packet{Info: u.Info, Data: append([]byte(nil), u.Data...)}
		reason := "unknown PGN"
		if u.Reason != nil {
			reason = u.Reason.Error()
		}
		p.observe(observationFromPacket(raw.KindDecodeError, packet, reason))
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
