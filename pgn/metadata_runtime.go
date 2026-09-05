package pgn

import (
	"fmt"
	"sort"
	"strings"
)

type SourceDefinition struct {
	PGN                          uint32
	StructName                   string
	SourceID                     string
	Description                  string
	Explanation                  string
	Type                         string
	Complete                     bool
	Fallback                     bool
	Missing                      []string
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
	Fields                       []SourceFieldDefinition
}

type SourceFieldDefinition struct {
	Order                               int
	SourceID                            string
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

// structInfoLookup maps generated struct type names (PgnInfo.Id) to their
// metadata. The metadata-driven codec (pgn/codec.go) uses it to compile a
// decode/encode plan for each PGN struct type.
var structInfoLookup map[string]*PgnInfo

func init() {
	initPgnInfoLookup()
}

func initPgnInfoLookup() {
	if len(sourceDefinitions) == 0 {
		return
	}

	PgnInfoLookup = make(map[uint32][]*PgnInfo)
	UnseenLookup = make(map[uint32][]*PgnInfo)
	structInfoLookup = make(map[string]*PgnInfo, len(sourceDefinitions))

	infos := make([]*PgnInfo, 0, len(sourceDefinitions))
	for _, def := range sourceDefinitions {
		info := pgnInfoFromSource(def)
		infos = append(infos, info)
		structInfoLookup[info.Id] = info
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

func pgnInfoFromSource(def SourceDefinition) *PgnInfo {
	fields := make(map[int]*FieldDescriptor, len(def.Fields))
	var manID ManufacturerCodeConst
	for _, field := range def.Fields {
		fd := fieldDescriptorFromSource(field)
		fields[field.Order] = fd
		if field.Name == "Manufacturer Code" && field.Match != nil {
			manID = ManufacturerCodeConst(*field.Match)
		}
	}

	info := &PgnInfo{
		SourceID:                     def.SourceID,
		Id:                           firstNonEmpty(def.StructName, exportedID(def.SourceID, def.Description)),
		PGN:                          def.PGN,
		Description:                  def.Description,
		Explanation:                  def.Explanation,
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
	info.DecodeComplete, info.EncodeComplete, info.CodecLimitations = codecSupport(def)
	return info
}

func codecSupport(def SourceDefinition) (bool, bool, []string) {
	var limitations []string
	if !def.Complete {
		limitations = append(limitations, "source schema is incomplete")
	}
	for _, field := range def.Fields {
		if field.FieldType == "VARIABLE" {
			limitations = append(limitations, fmt.Sprintf("field %d (%s) requires an unambiguous fixed-width commanded PGN field", field.Order, field.Name))
		}
	}
	if def.StructName == "NavicoUdbDatabaseObjectDump" || def.StructName == "NavicoDataTypeSourceDirectory" {
		limitations = append(limitations, "record lengths include a sub-record header absent from the source field layout")
	}
	return len(limitations) == 0, len(limitations) == 0, limitations
}

func fieldDescriptorFromSource(field SourceFieldDefinition) *FieldDescriptor {
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
	resolution := float64(1)
	if field.Resolution != nil {
		resolution = *field.Resolution
	}
	return &FieldDescriptor{
		SourceID:                            field.SourceID,
		Name:                                field.Name,
		Description:                         field.Description,
		BitLength:                           bitLength,
		BitLengthField:                      field.BitLengthField,
		BitOffset:                           bitOffset,
		BitLengthVariable:                   field.BitLengthVariable || field.BitLength == nil,
		SourceType:                          field.FieldType,
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

func pgnFieldGoType(field SourceFieldDefinition) string {
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

func readRawAt(data []uint8, bitOffset uint16, bitLength uint16) (uint64, error) {
	stream := NewPgnDataStream(data)
	stream.skipBits(bitOffset)
	return stream.getNumberRaw(bitLength)
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

func exportedID(upstreamID string, description string) string {
	if upstreamID != "" {
		return upstreamID
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

func sourceUint8(v uint8) *uint8       { return &v }
func sourceInt(v int) *int             { return &v }
func sourceUint16(v uint16) *uint16    { return &v }
func sourceFloat64(v float64) *float64 { return &v }
func sourceBool(v bool) *bool          { return &v }
func sourceInt64(v int64) *int64       { return &v }

func sourceMatch(v int) *int { return &v }
