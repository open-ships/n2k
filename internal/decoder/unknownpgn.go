package decoder

import (
	"fmt"
	"strings"

	"github.com/open-ships/n2k/pgn"
)

func buildUnknownPGN(p *Packet) *pgn.UnknownPGN {
	ret := &pgn.UnknownPGN{
		Data:   p.Data,
		Reason: fmt.Errorf("%s", mergeErrorStrings(p.ParseErrors)),
	}
	ret.Info = p.Info

	if pgn.IsProprietaryPGN(ret.Info.PGN) {
		if p.Manufacturer != 0 {
			ret.ManufacturerCode = p.Manufacturer
		} else {
			ret.ManufacturerCode, ret.IndustryCode, _ = pgn.GetProprietaryInfo(p.Data)
		}
	}

	ret.WasUnseen = pgn.SearchUnseenList(ret.Info.PGN)
	return ret
}

func mergeErrorStrings(errs []error) string {
	sErrs := make([]string, 0, len(errs))
	for i := range errs {
		sErrs = append(sErrs, errs[i].Error())
	}
	return strings.Join(sErrs, ", ")
}
