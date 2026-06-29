package pgn

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

type CanboatDefinition struct {
	PGN                          uint32
	StructName                   string
	CanboatId                    string
	Description                  string
	Explanation                  string
	URL                          string
	Type                         string
	Complete                     bool
	Fallback                     bool
	Missing                      []string
	FieldCount                   int
	Length                       *int
	MinLength                    *int
	Priority                     *uint8
	TransmissionInterval         *int
	TransmissionIrregular        *bool
	RepeatingFieldSet1StartField *int
	RepeatingFieldSet1CountField *int
	RepeatingFieldSet1Size       *int
	RepeatingFieldSet2StartField *int
	RepeatingFieldSet2CountField *int
	RepeatingFieldSet2Size       *int
	Fields                       []CanboatFieldDefinition
}

type CanboatFieldDefinition struct {
	Order                               int
	CanboatId                           string
	Name                                string
	Description                         string
	BitLength                           *uint16
	BitLengthField                      *int
	BitOffset                           *uint16
	BitStart                            *uint16
	Resolution                          *float64
	Signed                              bool
	Unit                                string
	FieldType                           string
	PhysicalQuantity                    string
	LookupEnumeration                   string
	LookupBitEnumeration                string
	LookupIndirectEnumeration           string
	LookupIndirectEnumerationFieldOrder *int
	LookupFieldTypeEnumeration          string
	Match                               *int
	RangeMin                            *float64
	RangeMax                            *float64
	Offset                              *float64
	OutOfRangeValue                     *int64
	PartOfPrimaryKey                    *bool
	ReservedValue                       *int64
	UnknownValue                        *int64
	Condition                           string
	BitLengthVariable                   bool
}

func init() {
	initPgnInfoLookup()
}

func initPgnInfoLookup() {
	if len(canboatDefinitions) == 0 {
		return
	}

	PgnInfoLookup = make(map[uint32][]*PgnInfo)
	UnseenLookup = make(map[uint32][]*PgnInfo)

	infos := make([]*PgnInfo, 0, len(canboatDefinitions))
	for _, def := range canboatDefinitions {
		info := pgnInfoFromCanboat(def)
		infos = append(infos, info)
	}

	sort.SliceStable(infos, func(i, j int) bool {
		if infos[i].PGN != infos[j].PGN {
			return infos[i].PGN < infos[j].PGN
		}
		return infos[i].Description < infos[j].Description
	})

	for _, info := range infos {
		PgnInfoLookup[info.PGN] = append(PgnInfoLookup[info.PGN], info)
		if !info.Complete || len(info.Missing) > 0 {
			UnseenLookup[info.PGN] = append(UnseenLookup[info.PGN], info)
		}
	}
}

func pgnInfoFromCanboat(def CanboatDefinition) *PgnInfo {
	fields := make(map[int]*FieldDescriptor, len(def.Fields))
	var manID ManufacturerCodeConst
	for _, field := range def.Fields {
		fd := fieldDescriptorFromCanboat(field)
		fields[field.Order] = fd
		if field.Name == "Manufacturer Code" && field.Match != nil {
			manID = ManufacturerCodeConst(*field.Match)
		}
	}

	return &PgnInfo{
		CanboatId:                    def.CanboatId,
		Id:                           firstNonEmpty(def.StructName, exportedID(def.CanboatId, def.Description)),
		PGN:                          def.PGN,
		Description:                  def.Description,
		Explanation:                  def.Explanation,
		URL:                          def.URL,
		Fast:                         def.Type == "Fast",
		Type:                         def.Type,
		Complete:                     def.Complete,
		Fallback:                     def.Fallback,
		Missing:                      append([]string(nil), def.Missing...),
		Length:                       def.Length,
		MinLength:                    def.MinLength,
		Priority:                     def.Priority,
		TransmissionInterval:         def.TransmissionInterval,
		TransmissionIrregular:        def.TransmissionIrregular,
		RepeatingFieldSet1StartField: def.RepeatingFieldSet1StartField,
		RepeatingFieldSet1CountField: def.RepeatingFieldSet1CountField,
		RepeatingFieldSet1Size:       def.RepeatingFieldSet1Size,
		RepeatingFieldSet2StartField: def.RepeatingFieldSet2StartField,
		RepeatingFieldSet2CountField: def.RepeatingFieldSet2CountField,
		RepeatingFieldSet2Size:       def.RepeatingFieldSet2Size,
		ManId:                        manID,
		Fields:                       fields,
	}
}

func fieldDescriptorFromCanboat(field CanboatFieldDefinition) *FieldDescriptor {
	bitLength := uint16(0)
	if field.BitLength != nil {
		bitLength = *field.BitLength
	}
	bitOffset := uint16(0)
	if field.BitOffset != nil {
		bitOffset = *field.BitOffset
	}
	bitStart := uint16(0)
	if field.BitStart != nil {
		bitStart = *field.BitStart
	}
	resolution := float32(1)
	if field.Resolution != nil {
		resolution = float32(*field.Resolution)
	}
	return &FieldDescriptor{
		CanboatId:                           field.CanboatId,
		Name:                                field.Name,
		Description:                         field.Description,
		BitLength:                           bitLength,
		BitLengthField:                      field.BitLengthField,
		BitOffset:                           bitOffset,
		BitLengthVariable:                   field.BitLengthVariable || field.BitLength == nil,
		CanboatType:                         field.FieldType,
		BitStart:                            bitStart,
		PhysicalQuantity:                    field.PhysicalQuantity,
		LookupEnumeration:                   field.LookupEnumeration,
		LookupBitEnumeration:                field.LookupBitEnumeration,
		LookupIndirectEnumeration:           field.LookupIndirectEnumeration,
		LookupIndirectEnumerationFieldOrder: field.LookupIndirectEnumerationFieldOrder,
		LookupFieldTypeEnumeration:          field.LookupFieldTypeEnumeration,
		Condition:                           field.Condition,
		Offset:                              field.Offset,
		RangeMin:                            field.RangeMin,
		RangeMax:                            field.RangeMax,
		OutOfRangeValue:                     field.OutOfRangeValue,
		ReservedValue:                       field.ReservedValue,
		UnknownValue:                        field.UnknownValue,
		PartOfPrimaryKey:                    field.PartOfPrimaryKey,
		GolangType:                          pgnFieldGoType(field),
		Resolution:                          resolution,
		Signed:                              field.Signed,
		Unit:                                field.Unit,
		Match:                               field.Match,
	}
}

func pgnFieldGoType(field CanboatFieldDefinition) string {
	switch field.FieldType {
	case "RESERVED", "SPARE":
		return ""
	case "STRING_FIX", "STRING_LZ", "STRING_LAU":
		return "string"
	case "BINARY", "VARIABLE", "DYNAMIC_FIELD_VALUE":
		return "[]uint8"
	case "FLOAT":
		return "*float32"
	default:
		if field.Signed {
			return "*int64"
		}
		return "*uint64"
	}
}

func decodeStructPayload(info *PgnInfo, target PGN, payload []uint8) error {
	if info == nil {
		return fmt.Errorf("missing PGN metadata for %T", target)
	}
	if target == nil {
		return fmt.Errorf("nil PGN target for %s", info.Description)
	}
	stream := NewPgnDataStream(payload)
	if !info.MatchesData(stream.data) {
		return fmt.Errorf("match failed for %s", info.Description)
	}
	values, err := decodeCanboatFieldValues(info, stream)
	if err != nil {
		return err
	}
	return setStructFields(target, values)
}

func encodeStructPayload(info *PgnInfo, message PGN) ([]byte, error) {
	if info == nil {
		return nil, fmt.Errorf("missing PGN metadata for %T", message)
	}
	values, err := structFieldValues(message)
	if err != nil {
		return nil, err
	}
	return encodeCanboatFieldValues(info, values)
}

func decodeCanboatFieldValues(info *PgnInfo, stream *PGNDataStream) (map[int]any, error) {
	fields := orderedFieldDescriptors(info)
	values := make(map[int]any, len(fields))
	for _, field := range fields {
		value, hasValue, err := decodeCanboatFieldValue(field.descriptor, stream)
		if err != nil {
			return values, err
		}
		if hasValue {
			values[field.index] = value
		}
		if stream.isEOF() {
			return values, nil
		}
	}
	return values, nil
}

func decodeCanboatFieldValue(field *FieldDescriptor, stream *PGNDataStream) (any, bool, error) {
	switch field.CanboatType {
	case "RESERVED", "SPARE":
		stream.skipBits(field.BitLength)
		return nil, false, nil
	case "STRING_FIX":
		v, err := stream.readFixedString(field.BitLength)
		if err != nil {
			return nil, false, err
		}
		return v, true, nil
	case "STRING_LZ":
		v, err := stream.readStringWithLength()
		if err != nil {
			return nil, false, err
		}
		return v, true, nil
	case "STRING_LAU":
		v, err := stream.readStringWithLengthAndControl()
		if err != nil {
			return nil, false, err
		}
		return v, true, nil
	case "BINARY", "VARIABLE", "DYNAMIC_FIELD_VALUE":
		v, err := stream.readBinaryData(byteAlignedBitLength(field))
		if err != nil {
			return nil, false, err
		}
		return v, true, nil
	case "FLOAT":
		v, err := stream.readFloat32()
		if err != nil {
			return nil, false, err
		}
		return v, true, nil
	case "DECIMAL", "DURATION", "PGN", "ISO_NAME", "DYNAMIC_FIELD_KEY", "DYNAMIC_FIELD_LENGTH":
		v, err := stream.getNumberRaw(field.BitLength)
		if err != nil {
			return nil, false, err
		}
		return v, true, nil
	default:
		v, err := stream.getNumberRaw(field.BitLength)
		if err != nil {
			return nil, false, err
		}
		if field.Signed {
			return signExtend(v, field.BitLength), true, nil
		}
		return v, true, nil
	}
}

func encodeCanboatFieldValues(info *PgnInfo, values map[int]any) ([]byte, error) {
	return encodeCanboatFields(info, func(order int, field *FieldDescriptor) (any, bool) {
		value, ok := values[order]
		return value, ok
	})
}

func encodeCanboatFields(info *PgnInfo, valueFor func(int, *FieldDescriptor) (any, bool)) ([]byte, error) {
	writer := NewPGNDataStreamWriter()
	for _, field := range orderedFieldDescriptors(info) {
		value, hasValue := valueFor(field.index, field.descriptor)
		switch field.descriptor.CanboatType {
		case "RESERVED":
			writer.writeReservedBits(field.descriptor.BitLength)
		case "SPARE":
			writer.writeSpareBits(field.descriptor.BitLength)
		case "STRING_FIX":
			if s, ok := value.(string); ok {
				writer.writeFixedString(s, field.descriptor.BitLength)
			} else {
				writer.writeFixedString("", field.descriptor.BitLength)
			}
		case "STRING_LZ":
			if s, ok := value.(string); ok {
				writer.writeStringWithLength(s)
			} else {
				writer.writeStringWithLength("")
			}
		case "STRING_LAU":
			if s, ok := value.(string); ok {
				writer.writeStringWithLengthAndControl(s)
			} else {
				writer.writeStringWithLengthAndControl("")
			}
		case "BINARY", "VARIABLE", "DYNAMIC_FIELD_VALUE":
			if data, ok := genericBytes(value); ok {
				bitLength := field.descriptor.BitLength
				if bitLength == 0 {
					bitLength = uint16(len(data) * 8)
				}
				writer.writeBinaryData(data, bitLength)
			} else {
				writer.writeBinaryData(nil, byteAlignedBitLength(field.descriptor))
			}
		case "FLOAT":
			if f, ok := genericFloat32(value); ok {
				writer.writeFloat32(&f)
			} else {
				writer.writeFloat32(nil)
			}
		default:
			if raw, ok := genericUint64(value); ok && hasValue {
				writer.setErr(writer.putNumberRaw(raw, field.descriptor.BitLength))
			} else if raw, ok := genericInt64(value); ok && hasValue {
				writer.setErr(writer.putSignedNumber(raw, field.descriptor.BitLength))
			} else if field.descriptor.Match != nil {
				writer.writeLookupField(uint64(*field.descriptor.Match), field.descriptor.BitLength)
			} else if field.descriptor.Signed {
				writer.setErr(writer.putNullSigned(field.descriptor.BitLength))
			} else {
				writer.setErr(writer.putNullUnsigned(field.descriptor.BitLength))
			}
		}
	}
	return writer.Bytes(), writer.Err()
}

func setStructFields(target PGN, values map[int]any) error {
	rv := reflect.ValueOf(target)
	if !rv.IsValid() || rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("expected non-nil pointer target, got %T", target)
	}
	rv = rv.Elem()
	if rv.Kind() != reflect.Struct {
		return fmt.Errorf("expected pointer to struct target, got %T", target)
	}

	rt := rv.Type()
	for i := 0; i < rv.NumField(); i++ {
		structField := rt.Field(i)
		orderText := structField.Tag.Get("n2k")
		if orderText == "" {
			continue
		}
		order, err := strconv.Atoi(orderText)
		if err != nil {
			return fmt.Errorf("%s has invalid n2k tag %q", structField.Name, orderText)
		}
		value, ok := values[order]
		if !ok {
			continue
		}
		if err := setStructField(rv.Field(i), value); err != nil {
			return fmt.Errorf("%s: %w", structField.Name, err)
		}
	}
	return nil
}

func setStructField(dst reflect.Value, value any) error {
	if !dst.CanSet() || value == nil {
		return nil
	}
	src := reflect.ValueOf(value)
	if !src.IsValid() {
		return nil
	}
	if src.Type().AssignableTo(dst.Type()) {
		dst.Set(src)
		return nil
	}
	if src.Type().ConvertibleTo(dst.Type()) {
		dst.Set(src.Convert(dst.Type()))
		return nil
	}
	if dst.Kind() == reflect.Ptr {
		if src.Kind() == reflect.Ptr {
			if src.Type().AssignableTo(dst.Type()) {
				dst.Set(src)
				return nil
			}
			if src.Type().Elem().ConvertibleTo(dst.Type().Elem()) {
				converted := reflect.New(dst.Type().Elem())
				converted.Elem().Set(src.Elem().Convert(dst.Type().Elem()))
				dst.Set(converted)
				return nil
			}
		}
		if src.Type().AssignableTo(dst.Type().Elem()) || src.Type().ConvertibleTo(dst.Type().Elem()) {
			converted := reflect.New(dst.Type().Elem())
			if src.Type().AssignableTo(dst.Type().Elem()) {
				converted.Elem().Set(src)
			} else {
				converted.Elem().Set(src.Convert(dst.Type().Elem()))
			}
			dst.Set(converted)
			return nil
		}
	}
	return fmt.Errorf("cannot assign %s to %s", src.Type(), dst.Type())
}

func structFieldValues(message PGN) (map[int]any, error) {
	rv := reflect.ValueOf(message)
	if !rv.IsValid() || rv.Kind() != reflect.Ptr || rv.IsNil() {
		return nil, fmt.Errorf("expected non-nil pointer message, got %T", message)
	}
	rv = rv.Elem()
	if rv.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected pointer to struct message, got %T", message)
	}

	values := make(map[int]any)
	rt := rv.Type()
	for i := 0; i < rv.NumField(); i++ {
		structField := rt.Field(i)
		orderText := structField.Tag.Get("n2k")
		if orderText == "" {
			continue
		}
		order, err := strconv.Atoi(orderText)
		if err != nil {
			return nil, fmt.Errorf("%s has invalid n2k tag %q", structField.Name, orderText)
		}
		field := rv.Field(i)
		if field.Kind() == reflect.Ptr && field.IsNil() {
			continue
		}
		values[order] = field.Interface()
	}
	return values, nil
}

func lookupPgnInfo(pgn uint32, description string) *PgnInfo {
	for _, info := range PgnInfoLookup[pgn] {
		if info.Description == description {
			return info
		}
	}
	return nil
}

func genericBytes(value any) ([]uint8, bool) {
	switch v := value.(type) {
	case []byte:
		return []uint8(v), true
	case string:
		return []uint8(v), true
	default:
		return nil, false
	}
}

func genericFloat32(value any) (float32, bool) {
	switch v := value.(type) {
	case float32:
		return v, true
	case *float32:
		if v == nil {
			return 0, false
		}
		return *v, true
	case float64:
		return float32(v), true
	case *float64:
		if v == nil {
			return 0, false
		}
		return float32(*v), true
	default:
		return 0, false
	}
}

func genericUint64(value any) (uint64, bool) {
	switch v := value.(type) {
	case uint:
		return uint64(v), true
	case uint8:
		return uint64(v), true
	case uint16:
		return uint64(v), true
	case uint32:
		return uint64(v), true
	case uint64:
		return v, true
	case int:
		return uint64(v), v >= 0
	case int8:
		return uint64(v), v >= 0
	case int16:
		return uint64(v), v >= 0
	case int32:
		return uint64(v), v >= 0
	case int64:
		return uint64(v), v >= 0
	default:
		rv := reflect.ValueOf(value)
		if rv.IsValid() && rv.Kind() == reflect.Ptr && !rv.IsNil() {
			return genericUint64(rv.Elem().Interface())
		}
		return 0, false
	}
}

func genericInt64(value any) (int64, bool) {
	switch v := value.(type) {
	case int:
		return int64(v), true
	case int8:
		return int64(v), true
	case int16:
		return int64(v), true
	case int32:
		return int64(v), true
	case int64:
		return v, true
	default:
		rv := reflect.ValueOf(value)
		if rv.IsValid() && rv.Kind() == reflect.Ptr && !rv.IsNil() {
			return genericInt64(rv.Elem().Interface())
		}
		return 0, false
	}
}

func (info *PgnInfo) MatchesData(data []uint8) bool {
	if info == nil {
		return false
	}
	for _, field := range info.Fields {
		if field.Match == nil {
			continue
		}
		if field.BitLength == 0 {
			return false
		}
		value, err := readRawAt(data, field.BitOffset, field.BitLength)
		if err != nil {
			return false
		}
		if int(value) != *field.Match {
			return false
		}
	}
	return true
}

func FilterMatchingPgnInfos(candidates []*PgnInfo, data []uint8) []*PgnInfo {
	if len(candidates) == 0 {
		return nil
	}
	matched := make([]*PgnInfo, 0, len(candidates))
	for _, candidate := range candidates {
		if candidate.MatchesData(data) {
			matched = append(matched, candidate)
		}
	}
	if len(matched) == 0 {
		return candidates
	}
	return matched
}

type orderedCanboatField struct {
	index      int
	descriptor *FieldDescriptor
}

func orderedFieldDescriptors(info *PgnInfo) []orderedCanboatField {
	indexes := make([]int, 0, len(info.Fields))
	for index := range info.Fields {
		indexes = append(indexes, index)
	}
	sort.Ints(indexes)
	fields := make([]orderedCanboatField, 0, len(indexes))
	for _, index := range indexes {
		fields = append(fields, orderedCanboatField{index: index, descriptor: info.Fields[index]})
	}
	return fields
}

func readRawAt(data []uint8, bitOffset uint16, bitLength uint16) (uint64, error) {
	stream := NewPgnDataStream(data)
	stream.skipBits(bitOffset)
	return stream.getNumberRaw(bitLength)
}

func byteAlignedBitLength(field *FieldDescriptor) uint16 {
	if field.BitLength == 0 {
		return 0
	}
	return (field.BitLength + 7) &^ 0x7
}

func signExtend(value uint64, bitLength uint16) int64 {
	if bitLength == 0 {
		return 0
	}
	mask := uint64(1 << (bitLength - 1))
	if value&mask == 0 {
		return int64(value)
	}
	value ^= mask
	return -int64(mask) + int64(value)
}

func exportedID(canboatID string, description string) string {
	if canboatID != "" {
		return canboatID
	}
	return strings.ReplaceAll(description, " ", "")
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func canboatUint8(v uint8) *uint8       { return &v }
func canboatInt(v int) *int             { return &v }
func canboatUint16(v uint16) *uint16    { return &v }
func canboatFloat64(v float64) *float64 { return &v }
func canboatBool(v bool) *bool          { return &v }
func canboatInt64(v int64) *int64       { return &v }

func canboatMatch(v int) *int { return &v }
