package pgn

import "fmt"

func (field *planField) keepDecimalZero(raws map[int]uint64) bool {
	if field.trailingZeroIsDigitField == nil || field.trailingZeroIsDigitValue == nil {
		return false
	}
	value, ok := raws[*field.trailingZeroIsDigitField]
	return ok && value == uint64(*field.trailingZeroIsDigitValue)
}

// DECIMAL bytes contain digit pairs 0..99, most significant pair first.
// They are neither BCD nor a little-endian binary integer.
func readDecimal(stream *PGNDataStream, bits uint16, keepZero bool) (string, error) {
	if stream.getBitOffset()+uint32(bits) > uint32(len(stream.data))*8 {
		stream.byteOffset = uint16(len(stream.data))
		stream.bitOffset = 0
		return "", nil
	}
	data, err := stream.readBinaryData(bits)
	if err != nil {
		return "", err
	}
	digits := make([]byte, 0, 2*len(data))
	for _, pair := range data {
		if pair > 99 {
			return "", nil
		}
		digits = append(digits, '0'+pair/10, '0'+pair%10)
	}
	if !keepZero && digits[len(digits)-1] == '0' {
		digits = digits[:len(digits)-1]
	}
	return string(digits), nil
}

func writeDecimal(writer *PGNDataStreamWriter, digits string, bits uint16, keepZero bool) {
	if digits == "" {
		writer.writeReservedBits(bits)
		return
	}
	width := int(bits/8) * 2
	if !keepZero && len(digits) == width-1 {
		digits += "0"
	} else if len(digits) != width || (!keepZero && digits[len(digits)-1] == '0') {
		writer.setErr(fmt.Errorf("decimal address %q has an invalid width or ambiguous trailing zero", digits))
		return
	}
	data := make([]byte, width/2)
	for i := range data {
		a, b := digits[i*2], digits[i*2+1]
		if a < '0' || a > '9' || b < '0' || b > '9' {
			writer.setErr(fmt.Errorf("decimal address %q contains a non-digit", digits))
			return
		}
		data[i] = (a-'0')*10 + b - '0'
	}
	writer.writeBinaryData(data, bits)
}
