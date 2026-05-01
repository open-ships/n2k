package pgn

type UnknownPGN struct {
	Info MessageInfo `json:"info"`
	Data             []uint8               `json:"data"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst     `json:"industryCode"`
	Reason           error                 `json:"reason"`
	WasUnseen        bool                  `json:"wasUnseen"`
}

func (u *UnknownPGN) PGNNumber() uint32 {
	return u.Info.PGN
}
