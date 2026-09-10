package pgn

import "errors"

// UnknownPGN retains an undecodable payload, header, and decode reason.
// It implements Message, but not PGN: callers must handle its raw Data explicitly
// instead of treating undecoded content as a supported typed transmission.
type UnknownPGN struct {
	Info             MessageInfo           `json:"info"`
	Data             []uint8               `json:"data"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst     `json:"industryCode"`
	Reason           error                 `json:"reason"`
	WasUnseen        bool                  `json:"wasUnseen"`
}

// PGNNumber returns the raw message number from Info.
func (u *UnknownPGN) PGNNumber() uint32 {
	return u.Info.PGN
}

// MessageInfo returns the metadata by value. Pointer and slice fields alias
// Info; use MessageInfo.Clone for an independent copy.
func (u *UnknownPGN) MessageInfo() MessageInfo {
	return u.Info
}

// SetMessageInfo replaces Info by value. The caller retains ownership of any
// referenced pointers and slices; pass info.Clone to transfer an independent copy.
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
