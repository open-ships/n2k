package pgn

import "errors"

type UnknownPGN struct {
	Info             MessageInfo           `json:"info"`
	Data             []uint8               `json:"data"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst     `json:"industryCode"`
	Reason           error                 `json:"reason"`
	WasUnseen        bool                  `json:"wasUnseen"`
}

func (u *UnknownPGN) PGNNumber() uint32 {
	return u.Info.PGN
}

func (u *UnknownPGN) MessageInfo() MessageInfo {
	return u.Info
}

func (u *UnknownPGN) SetMessageInfo(info MessageInfo) {
	u.Info = info
}

// Clone returns an owned raw message and snapshots the diagnostic text.
func (u *UnknownPGN) Clone() Message {
	if u == nil {
		return (*UnknownPGN)(nil)
	}
	copy := *u
	copy.Info = u.Info.Clone()
	copy.Data = cloneSlice(u.Data)
	if u.Reason != nil {
		copy.Reason = errors.New(u.Reason.Error())
	}
	return &copy
}
