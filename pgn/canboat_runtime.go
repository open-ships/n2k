package pgn

import (
	"fmt"
	"math"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

type CanboatDefinition struct {
	PGN                          uint32
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

type GenericMessage struct {
	Info        MessageInfo    `json:"info"`
	CanboatId   string         `json:"canboatId"`
	Description string         `json:"description"`
	Fields      map[string]any `json:"fields"`
	Raw         []uint8        `json:"raw,omitempty"`
}

func (m *GenericMessage) PGNNumber() uint32 {
	return m.Info.PGN
}

type generatedConcreteDefinition struct {
	Id      string
	Decoder func(MessageInfo, *PGNDataStream) (Message, error)
	Encoder func(Message) ([]byte, error)
}

func applyCanboatGeneratedOverlay() {
	if len(canboatGeneratedDefinitions) == 0 {
		return
	}

	existing := make(map[string]*PgnInfo)
	for _, infos := range PgnInfoLookup {
		for _, info := range infos {
			existing[pgnDescriptionKey(info.PGN, info.Description)] = info
		}
	}

	PgnInfoLookup = make(map[uint32][]*PgnInfo)
	UnseenLookup = make(map[uint32][]*PgnInfo)
	EncoderLookup = make(map[uint32]func(Message) ([]byte, error))

	generated := make([]*PgnInfo, 0, len(canboatGeneratedDefinitions))
	for _, def := range canboatGeneratedDefinitions {
		info := pgnInfoFromCanboat(def)
		if old := existing[pgnDescriptionKey(info.PGN, info.Description)]; old != nil && pgnInfoLayoutCompatible(old, info) {
			info.Id = old.Id
			info.Decoder = old.Decoder
			info.Encoder = old.Encoder
		}
		if concrete, ok := generatedConcreteDefinitions[pgnDescriptionKey(info.PGN, info.Description)]; ok && info.Decoder == nil {
			info.Id = concrete.Id
			info.Decoder = concrete.Decoder
			info.Encoder = concrete.Encoder
		}
		if info.Decoder == nil {
			info.Decoder = canboatGenericDecoder(info)
		}
		if info.Encoder == nil {
			info.Encoder = canboatGenericEncoder(info)
		}
		info.Self = info
		generated = append(generated, info)
	}

	sort.SliceStable(generated, func(i, j int) bool {
		if generated[i].PGN != generated[j].PGN {
			return generated[i].PGN < generated[j].PGN
		}
		return generated[i].Description < generated[j].Description
	})

	for _, info := range generated {
		PgnInfoLookup[info.PGN] = append(PgnInfoLookup[info.PGN], info)
		if !info.Complete || len(info.Missing) > 0 {
			UnseenLookup[info.PGN] = append(UnseenLookup[info.PGN], info)
		}
		if _, exists := EncoderLookup[info.PGN]; !exists {
			EncoderLookup[info.PGN] = info.Encoder
		}
	}
}

func pgnInfoLayoutCompatible(old *PgnInfo, current *PgnInfo) bool {
	if old == nil || current == nil {
		return false
	}
	if old.Fast != current.Fast {
		return false
	}
	if len(old.Fields) != len(current.Fields) {
		return false
	}
	for order, currentField := range current.Fields {
		oldField := old.Fields[order]
		if oldField == nil || currentField == nil {
			return false
		}
		if oldField.BitLength != currentField.BitLength ||
			oldField.BitOffset != currentField.BitOffset ||
			oldField.BitLengthVariable != currentField.BitLengthVariable ||
			oldField.Signed != currentField.Signed {
			return false
		}
		if !compatibleCanboatType(oldField.CanboatType, currentField.CanboatType) {
			return false
		}
		if math.Abs(float64(oldField.Resolution-currentField.Resolution)) > 1e-6 {
			return false
		}
	}
	return true
}

func compatibleCanboatType(oldType string, currentType string) bool {
	if oldType == currentType {
		return true
	}
	if numericCanboatType(oldType) && numericCanboatType(currentType) {
		return true
	}
	switch currentType {
	case "DYNAMIC_FIELD_KEY":
		return oldType == "FIELDTYPE_LOOKUP"
	case "DYNAMIC_FIELD_VALUE":
		return oldType == "KEY_VALUE"
	case "DYNAMIC_FIELD_LENGTH", "PGN", "ISO_NAME", "DECIMAL":
		return oldType == "NUMBER"
	case "DURATION":
		return oldType == "TIME" || oldType == "NUMBER"
	default:
		return false
	}
}

func numericCanboatType(fieldType string) bool {
	switch fieldType {
	case "BITLOOKUP", "DATE", "DECIMAL", "DURATION", "DYNAMIC_FIELD_KEY", "DYNAMIC_FIELD_LENGTH",
		"FIELD_INDEX", "INDIRECT_LOOKUP", "ISO_NAME", "LOOKUP", "MMSI", "NUMBER", "PGN", "TIME":
		return true
	default:
		return false
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
		Id:                           exportedID(def.CanboatId, def.Description),
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
		Resolution:                          resolution,
		Signed:                              field.Signed,
		Unit:                                field.Unit,
		Match:                               field.Match,
	}
}

func canboatGenericDecoder(info *PgnInfo) func(MessageInfo, *PGNDataStream) (Message, error) {
	return func(messageInfo MessageInfo, stream *PGNDataStream) (Message, error) {
		if !info.MatchesData(stream.data) {
			return nil, fmt.Errorf("match failed for %s", info.Description)
		}
		fields, err := decodeGenericFields(info, stream)
		if err != nil {
			return nil, err
		}
		return &GenericMessage{
			Info:        messageInfo,
			CanboatId:   info.CanboatId,
			Description: info.Description,
			Fields:      fields,
			Raw:         append([]uint8(nil), stream.data...),
		}, nil
	}
}

func decodeGeneratedConcrete(info *PgnInfo, target Message, stream *PGNDataStream) error {
	if info == nil {
		return fmt.Errorf("missing CANboat metadata for %T", target)
	}
	if target == nil {
		return fmt.Errorf("nil CANboat concrete target for %s", info.Description)
	}
	if !info.MatchesData(stream.data) {
		return fmt.Errorf("match failed for %s", info.Description)
	}
	values, err := decodeCanboatFieldValues(info, stream)
	if err != nil {
		return err
	}
	return setGeneratedConcreteFields(target, values)
}

func encodeGeneratedConcrete(info *PgnInfo, message Message) ([]byte, error) {
	if info == nil {
		return nil, fmt.Errorf("missing CANboat metadata for %T", message)
	}
	values, err := generatedConcreteFieldValues(message)
	if err != nil {
		return nil, err
	}
	return encodeCanboatFieldValues(info, values)
}

func canboatGenericEncoder(info *PgnInfo) func(Message) ([]byte, error) {
	return func(message Message) ([]byte, error) {
		generic, ok := message.(*GenericMessage)
		if !ok {
			return nil, fmt.Errorf("expected *GenericMessage for %s, got %T", info.Description, message)
		}
		if generic.PGNNumber() != info.PGN {
			return nil, fmt.Errorf("expected PGN %d, got %d", info.PGN, generic.PGNNumber())
		}
		if len(generic.Raw) > 0 {
			return append([]uint8(nil), generic.Raw...), nil
		}
		return encodeGenericFields(info, generic.Fields)
	}
}

func decodeGenericFields(info *PgnInfo, stream *PGNDataStream) (map[string]any, error) {
	fieldValues, err := decodeCanboatFieldValues(info, stream)
	if err != nil {
		return nil, err
	}
	fields := orderedFieldDescriptors(info)
	values := make(map[string]any, len(fieldValues))
	for _, field := range fields {
		value, ok := fieldValues[field.index]
		if ok {
			values[jsonFieldName(field.descriptor)] = value
		}
	}
	return values, nil
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

func encodeGenericFields(info *PgnInfo, values map[string]any) ([]byte, error) {
	return encodeCanboatFields(info, func(order int, field *FieldDescriptor) (any, bool) {
		return genericFieldValue(field, values)
	})
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

func setGeneratedConcreteFields(target Message, values map[int]any) error {
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
		if err := setGeneratedConcreteField(rv.Field(i), value); err != nil {
			return fmt.Errorf("%s: %w", structField.Name, err)
		}
	}
	return nil
}

func setGeneratedConcreteField(dst reflect.Value, value any) error {
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

func generatedConcreteFieldValues(message Message) (map[int]any, error) {
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

func generatedConcreteInfo(pgn uint32, description string) *PgnInfo {
	for _, info := range PgnInfoLookup[pgn] {
		if info.Description == description {
			return info
		}
	}
	return nil
}

func generatedConcreteTypeError(expected string, actual Message) error {
	return fmt.Errorf("expected *%s, got %T", expected, actual)
}

func genericFieldValue(field *FieldDescriptor, values map[string]any) (any, bool) {
	if len(values) == 0 {
		return nil, false
	}
	keys := []string{jsonFieldName(field), field.CanboatId, field.Name}
	for _, key := range keys {
		if key == "" {
			continue
		}
		value, ok := values[key]
		if ok {
			return value, true
		}
	}
	return nil, false
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

func orderedFields(info *PgnInfo) []*FieldDescriptor {
	ordered := orderedFieldDescriptors(info)
	fields := make([]*FieldDescriptor, 0, len(ordered))
	for _, field := range ordered {
		fields = append(fields, field.descriptor)
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

func jsonFieldName(field *FieldDescriptor) string {
	if field.CanboatId != "" {
		return field.CanboatId
	}
	name := strings.TrimSpace(field.Name)
	if name == "" {
		return "field"
	}
	parts := strings.FieldsFunc(name, func(r rune) bool {
		return (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && (r < '0' || r > '9')
	})
	for i := range parts {
		parts[i] = strings.ToLower(parts[i])
		if i > 0 && parts[i] != "" {
			parts[i] = strings.ToUpper(parts[i][:1]) + parts[i][1:]
		}
	}
	return strings.Join(parts, "")
}

func pgnDescriptionKey(pgn uint32, description string) string {
	return fmt.Sprintf("%d\x00%s", pgn, description)
}

func exportedID(canboatID string, description string) string {
	if canboatID != "" {
		return canboatID
	}
	return strings.ReplaceAll(description, " ", "")
}

func canboatUint8(v uint8) *uint8       { return &v }
func canboatInt(v int) *int             { return &v }
func canboatUint16(v uint16) *uint16    { return &v }
func canboatFloat64(v float64) *float64 { return &v }
func canboatBool(v bool) *bool          { return &v }
func canboatInt64(v int64) *int64       { return &v }

func canboatMatch(v int) *int { return &v }
