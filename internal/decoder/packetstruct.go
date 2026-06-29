package decoder

import "github.com/open-ships/n2k/pgn"

type Handler interface {
	HandleStruct(pgn.Message)
}

type Decoder struct {
	handler Handler
}

func New() *Decoder {
	return &Decoder{}
}

func (ps *Decoder) SetOutput(sh Handler) {
	ps.handler = sh
}

func (ps *Decoder) Decode(pkt Packet) {
	msg, err := pgn.DecodePayload(pkt.Info, pkt.Data)
	if err != nil {
		pkt.ParseErrors = append(pkt.ParseErrors, err)
		ps.pgnReady(pkt.UnknownPGN())
		return
	}
	ps.pgnReady(msg)
}

func (ps *Decoder) pgnReady(fullPGN pgn.Message) {
	if ps.handler != nil {
		ps.handler.HandleStruct(fullPGN)
	}
}
