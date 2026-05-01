package decoder

import (
	"fmt"

	"github.com/open-ships/n2k/pgn"
)

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
	if len(pkt.Decoders) > 0 {
		for _, dec := range pkt.Decoders {
			stream := pgn.NewPgnDataStream(pkt.Data)
			ret, err := dec(pkt.Info, stream)
			if err != nil {
				pkt.ParseErrors = append(pkt.ParseErrors, err)
				continue
			} else {
				ps.pgnReady(ret)
				return
			}
		}

		ps.pgnReady(pkt.UnknownPGN())
	} else {
		pkt.ParseErrors = append(pkt.ParseErrors, fmt.Errorf("no matching decoder"))
		ps.pgnReady(pkt.UnknownPGN())
	}
}

func (ps *Decoder) pgnReady(fullPGN pgn.Message) {
	if ps.handler != nil {
		ps.handler.HandleStruct(fullPGN)
	}
}
