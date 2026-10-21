package pgn

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// Basic RDS G0, IEC 62106 table E.1 (ETSI TS 101 756 V1.6.1, annex C).
// The later DAB repertoire differs and must not be used for FM radio text.
// Undefined positions retain their byte value.
var rdsExtended = [...]rune{
	0x00e1, 0x00e0, 0x00e9, 0x00e8, 0x00ed, 0x00ec, 0x00f3, 0x00f2, 0x00fa, 0x00f9, 0x00d1, 0x00c7, 0x015e, 0x00df, 0x00a1, 0x0132,
	0x00e2, 0x00e4, 0x00ea, 0x00eb, 0x00ee, 0x00ef, 0x00f4, 0x00f6, 0x00fb, 0x00fc, 0x00f1, 0x00e7, 0x015f, 0x011f, 0x0131, 0x0133,
	0x00aa, 0x03b1, 0x00a9, 0x2030, 0x011e, 0x011b, 0x0148, 0x0151, 0x03c0, 0x20ac, 0x00a3, 0x0024, 0x2190, 0x2191, 0x2192, 0x2193,
	0x00ba, 0x00b9, 0x00b2, 0x00b3, 0x00b1, 0x0130, 0x0144, 0x0171, 0x00b5, 0x00bf, 0x00f7, 0x00b0, 0x00bc, 0x00bd, 0x00be, 0x00a7,
	0x00c1, 0x00c0, 0x00c9, 0x00c8, 0x00cd, 0x00cc, 0x00d3, 0x00d2, 0x00da, 0x00d9, 0x0158, 0x010c, 0x0160, 0x017d, 0x00d0, 0x013f,
	0x00c2, 0x00c4, 0x00ca, 0x00cb, 0x00ce, 0x00cf, 0x00d4, 0x00d6, 0x00db, 0x00dc, 0x0159, 0x010d, 0x0161, 0x017e, 0x0111, 0x0140,
	0x00c3, 0x00c5, 0x00c6, 0x0152, 0x0177, 0x00dd, 0x00d5, 0x00d8, 0x00de, 0x014a, 0x0154, 0x0106, 0x015a, 0x0179, 0x0166, 0x00f0,
	0x00e3, 0x00e5, 0x00e6, 0x0153, 0x0175, 0x00fd, 0x00f5, 0x00f8, 0x00fe, 0x014b, 0x0155, 0x0107, 0x015b, 0x017a, 0x0167, 0x00ff,
}

func rdsCodepoint(value byte) rune {
	switch value {
	case 0x24:
		return '¤'
	case 0x5e:
		return '―'
	case 0x60:
		return '║'
	case 0x7e:
		return '¯'
	}
	if value >= 0x80 {
		return rdsExtended[value-0x80]
	}
	return rune(value)
}

var rdsBytes = func() map[rune]byte {
	result := make(map[rune]byte, 256)
	for value := 0; value < 256; value++ {
		result[rdsCodepoint(byte(value))] = byte(value)
	}
	return result
}()

// UTF-8 takes precedence, matching the source decoder and radios that already
// convert RDS text. Otherwise the entire byte sequence uses the G0 repertoire.
func decodeRDSText(value string) string {
	if utf8.ValidString(value) {
		return value
	}
	var decoded strings.Builder
	for i := 0; i < len(value); i++ {
		decoded.WriteRune(rdsCodepoint(value[i]))
	}
	return decoded.String()
}

func encodeRDSText(value string) (string, error) {
	if !utf8.ValidString(value) {
		return "", fmt.Errorf("RDS text must be valid UTF-8")
	}
	var encoded strings.Builder
	for _, char := range value {
		b, ok := rdsBytes[char]
		if !ok {
			return value, nil
		}
		encoded.WriteByte(b)
	}
	wire := encoded.String()
	// Some G0 byte sequences are also valid UTF-8 with a different meaning.
	// Use UTF-8 for those and for text outside G0, preserving the decoder's
	// precedence rule. Unchanged decoded messages replay their original bytes.
	if decodeRDSText(wire) != value {
		return value, nil
	}
	return wire, nil
}
