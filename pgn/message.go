package pgn

type Message interface {
	PGNNumber() uint32
}
