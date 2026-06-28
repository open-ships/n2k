package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"go/format"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

const canboatURL = "https://raw.githubusercontent.com/canboat/canboat/refs/heads/master/docs/canboat.json"
const generatedDefinitionsPath = "pgn/canboat_generated.go"
const generatedConcreteRegistryPath = "pgn/concrete_generated.go"
const generatedConcreteGlob = "pgn/*_concrete_generated.go"

type canboatFile struct {
	SchemaVersion string       `json:"SchemaVersion"`
	Version       string       `json:"Version"`
	PGNs          []canboatPGN `json:"PGNs"`
}

type canboatPGN struct {
	PGN                          uint32         `json:"PGN"`
	Id                           string         `json:"Id"`
	Description                  string         `json:"Description"`
	Explanation                  string         `json:"Explanation"`
	URL                          string         `json:"URL"`
	Type                         string         `json:"Type"`
	Complete                     bool           `json:"Complete"`
	Fallback                     bool           `json:"Fallback"`
	Missing                      []string       `json:"Missing"`
	FieldCount                   int            `json:"FieldCount"`
	Length                       *int           `json:"Length"`
	MinLength                    *int           `json:"MinLength"`
	Priority                     *uint8         `json:"Priority"`
	TransmissionInterval         *int           `json:"TransmissionInterval"`
	TransmissionIrregular        *bool          `json:"TransmissionIrregular"`
	RepeatingFieldSet1StartField *int           `json:"RepeatingFieldSet1StartField"`
	RepeatingFieldSet1CountField *int           `json:"RepeatingFieldSet1CountField"`
	RepeatingFieldSet1Size       *int           `json:"RepeatingFieldSet1Size"`
	RepeatingFieldSet2StartField *int           `json:"RepeatingFieldSet2StartField"`
	RepeatingFieldSet2CountField *int           `json:"RepeatingFieldSet2CountField"`
	RepeatingFieldSet2Size       *int           `json:"RepeatingFieldSet2Size"`
	Fields                       []canboatField `json:"Fields"`
}

type canboatField struct {
	Order                               int      `json:"Order"`
	Id                                  string   `json:"Id"`
	Name                                string   `json:"Name"`
	Description                         text     `json:"Description"`
	BitLength                           *uint16  `json:"BitLength"`
	BitLengthField                      *int     `json:"BitLengthField"`
	BitLengthVariable                   bool     `json:"BitLengthVariable"`
	BitOffset                           *uint16  `json:"BitOffset"`
	BitStart                            *uint16  `json:"BitStart"`
	Resolution                          *float64 `json:"Resolution"`
	Signed                              bool     `json:"Signed"`
	Unit                                string   `json:"Unit"`
	FieldType                           string   `json:"FieldType"`
	PhysicalQuantity                    string   `json:"PhysicalQuantity"`
	LookupEnumeration                   string   `json:"LookupEnumeration"`
	LookupBitEnumeration                string   `json:"LookupBitEnumeration"`
	LookupIndirectEnumeration           string   `json:"LookupIndirectEnumeration"`
	LookupIndirectEnumerationFieldOrder *int     `json:"LookupIndirectEnumerationFieldOrder"`
	LookupFieldTypeEnumeration          string   `json:"LookupFieldTypeEnumeration"`
	Match                               *int     `json:"Match"`
	RangeMin                            *float64 `json:"RangeMin"`
	RangeMax                            *float64 `json:"RangeMax"`
	Offset                              *float64 `json:"Offset"`
	OutOfRangeValue                     *int64   `json:"OutOfRangeValue"`
	PartOfPrimaryKey                    *bool    `json:"PartOfPrimaryKey"`
	ReservedValue                       *int64   `json:"ReservedValue"`
	UnknownValue                        *int64   `json:"UnknownValue"`
	Condition                           string   `json:"Condition"`
}

type text string

func (t *text) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*t = ""
		return nil
	}
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		*t = text(s)
		return nil
	}
	*t = text(string(data))
	return nil
}

func main() {
	check := flag.Bool("check", false, "fail if generated CANboat metadata is not current")
	flag.Parse()

	raw, err := fetch(canboatURL)
	if err != nil {
		fatal(err)
	}
	var source canboatFile
	if err := json.Unmarshal(raw, &source); err != nil {
		fatal(err)
	}
	definitions, err := generateDefinitions(source)
	if err != nil {
		fatal(err)
	}
	concrete, err := generateConcrete(source)
	if err != nil {
		fatal(err)
	}
	if *check {
		currentDefinitions, err := os.ReadFile(generatedDefinitionsPath)
		if err != nil {
			fatal(err)
		}
		if !bytes.Equal(currentDefinitions, definitions) {
			fatal(fmt.Errorf("%s is not synced with %s", generatedDefinitionsPath, canboatURL))
		}
		if err := checkConcreteFiles(concrete); err != nil {
			fatal(err)
		}
		return
	}
	if err := os.WriteFile(generatedDefinitionsPath, definitions, 0o600); err != nil {
		fatal(err)
	}
	if err := writeConcreteFiles(concrete); err != nil {
		fatal(err)
	}
}

func fetch(url string) ([]byte, error) {
	client := http.Client{Timeout: 60 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("GET %s: %s", url, resp.Status)
	}
	return io.ReadAll(resp.Body)
}

func generateDefinitions(source canboatFile) ([]byte, error) {
	sort.SliceStable(source.PGNs, func(i, j int) bool {
		if source.PGNs[i].PGN != source.PGNs[j].PGN {
			return source.PGNs[i].PGN < source.PGNs[j].PGN
		}
		return source.PGNs[i].Description < source.PGNs[j].Description
	})

	var b bytes.Buffer
	fmt.Fprintf(&b, "// Code generated by PGN catalog generator; DO NOT EDIT.\n")
	fmt.Fprintf(&b, "// Source: %s\n", canboatURL)
	fmt.Fprintf(&b, "// CANboat schema %s version %s.\n\n", source.SchemaVersion, source.Version)
	fmt.Fprintf(&b, "package pgn\n\n")
	fmt.Fprintf(&b, "var canboatGeneratedDefinitions = []CanboatDefinition{\n")
	for _, pgn := range source.PGNs {
		writePGN(&b, pgn)
	}
	fmt.Fprintf(&b, "}\n")
	return format.Source(b.Bytes())
}

type concretePGN struct {
	PGN      canboatPGN
	TypeName string
}

func generateConcrete(source canboatFile) (map[string][]byte, error) {
	sort.SliceStable(source.PGNs, func(i, j int) bool {
		if source.PGNs[i].PGN != source.PGNs[j].PGN {
			return source.PGNs[i].PGN < source.PGNs[j].PGN
		}
		return source.PGNs[i].Description < source.PGNs[j].Description
	})

	names := uniqueConcreteTypeNames(source.PGNs)
	categories := make(map[string][]concretePGN)
	var categoryOrder []string
	for i, pgn := range source.PGNs {
		category := concreteCategory(pgn)
		if _, ok := categories[category]; !ok {
			categoryOrder = append(categoryOrder, category)
		}
		categories[category] = append(categories[category], concretePGN{
			PGN:      pgn,
			TypeName: names[i],
		})
	}

	outputs := make(map[string][]byte, len(categories)+1)
	registry, err := generateConcreteRegistry(source, categoryOrder)
	if err != nil {
		return nil, err
	}
	outputs[generatedConcreteRegistryPath] = registry

	for _, category := range categoryOrder {
		out, err := generateConcreteCategory(source, category, categories[category])
		if err != nil {
			return nil, err
		}
		outputs[concreteCategoryPath(category)] = out
	}
	return outputs, nil
}

func generateConcreteRegistry(source canboatFile, categories []string) ([]byte, error) {
	var b bytes.Buffer
	fmt.Fprintf(&b, "// Code generated by PGN catalog generator; DO NOT EDIT.\n")
	fmt.Fprintf(&b, "// Source schema %s version %s.\n\n", source.SchemaVersion, source.Version)
	fmt.Fprintf(&b, "package pgn\n\n")
	fmt.Fprintf(&b, "var generatedConcreteDefinitions = mergeGeneratedConcreteDefinitions(\n")
	for _, category := range categories {
		fmt.Fprintf(&b, "%s,\n", concreteCategoryVarName(category))
	}
	fmt.Fprintf(&b, ")\n\n")
	fmt.Fprintf(&b, "func mergeGeneratedConcreteDefinitions(sources ...map[string]generatedConcreteDefinition) map[string]generatedConcreteDefinition {\n")
	fmt.Fprintf(&b, "merged := make(map[string]generatedConcreteDefinition)\n")
	fmt.Fprintf(&b, "for _, source := range sources {\n")
	fmt.Fprintf(&b, "for key, definition := range source {\n")
	fmt.Fprintf(&b, "merged[key] = definition\n")
	fmt.Fprintf(&b, "}\n")
	fmt.Fprintf(&b, "}\n")
	fmt.Fprintf(&b, "return merged\n")
	fmt.Fprintf(&b, "}\n")
	return format.Source(b.Bytes())
}

func generateConcreteCategory(source canboatFile, category string, entries []concretePGN) ([]byte, error) {
	var b bytes.Buffer
	fmt.Fprintf(&b, "// Code generated by PGN catalog generator; DO NOT EDIT.\n")
	fmt.Fprintf(&b, "// Source schema %s version %s.\n\n", source.SchemaVersion, source.Version)
	fmt.Fprintf(&b, "package pgn\n\n")
	fmt.Fprintf(&b, "var %s = map[string]generatedConcreteDefinition{\n", concreteCategoryVarName(category))
	for _, entry := range entries {
		pgn := entry.PGN
		typeName := entry.TypeName
		fmt.Fprintf(&b, "pgnDescriptionKey(%d, %q): {Id:%q, Decoder:Decode%s, Encoder:encode%sMsg},\n",
			pgn.PGN, pgn.Description, typeName, typeName, typeName)
	}
	fmt.Fprintf(&b, "}\n\n")
	for _, entry := range entries {
		writeConcretePGN(&b, entry.PGN, entry.TypeName)
	}
	return format.Source(b.Bytes())
}

func writeConcreteFiles(outputs map[string][]byte) error {
	if err := removeStaleConcreteFiles(outputs); err != nil {
		return err
	}
	paths := sortedConcretePaths(outputs)
	for _, path := range paths {
		if err := os.WriteFile(path, outputs[path], 0o600); err != nil {
			return err
		}
	}
	return nil
}

func checkConcreteFiles(outputs map[string][]byte) error {
	if err := checkStaleConcreteFiles(outputs); err != nil {
		return err
	}
	paths := sortedConcretePaths(outputs)
	for _, path := range paths {
		current, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if !bytes.Equal(current, outputs[path]) {
			return fmt.Errorf("%s is not synced with %s", path, canboatURL)
		}
	}
	return nil
}

func removeStaleConcreteFiles(outputs map[string][]byte) error {
	paths, err := existingConcretePaths()
	if err != nil {
		return err
	}
	for _, path := range paths {
		if _, ok := outputs[path]; ok {
			continue
		}
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			return err
		}
	}
	return nil
}

func checkStaleConcreteFiles(outputs map[string][]byte) error {
	paths, err := existingConcretePaths()
	if err != nil {
		return err
	}
	for _, path := range paths {
		if _, ok := outputs[path]; !ok {
			return fmt.Errorf("%s is stale generated concrete output", path)
		}
	}
	return nil
}

func existingConcretePaths() ([]string, error) {
	matches, err := filepath.Glob(generatedConcreteGlob)
	if err != nil {
		return nil, err
	}
	paths := append([]string(nil), matches...)
	if _, err := os.Stat(generatedConcreteRegistryPath); err == nil {
		paths = append(paths, generatedConcreteRegistryPath)
	} else if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	sort.Strings(paths)
	return paths, nil
}

func sortedConcretePaths(outputs map[string][]byte) []string {
	paths := make([]string, 0, len(outputs))
	for path := range outputs {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	return paths
}

func uniqueConcreteTypeNames(pgns []canboatPGN) []string {
	names := make([]string, len(pgns))
	used := make(map[string]int, len(pgns))
	for i, pgn := range pgns {
		base := "Pgn" + exportedGoIdentifier(firstNonEmpty(pgn.Id, pgn.Description))
		count := used[base]
		used[base] = count + 1
		if count > 0 {
			base = fmt.Sprintf("%s%d", base, count+1)
		}
		names[i] = base
	}
	return names
}

func writeConcretePGN(b *bytes.Buffer, p canboatPGN, typeName string) {
	fmt.Fprintf(b, "type %s struct {\n", typeName)
	fmt.Fprintf(b, "Info MessageInfo `json:\"info\"`\n")

	fieldNames := make(map[string]int)
	for _, field := range p.Fields {
		fieldType, ok := concreteFieldType(field)
		if !ok {
			continue
		}
		fieldName := uniqueFieldName(fieldNames, exportedGoIdentifier(firstNonEmpty(field.Id, field.Name)))
		fmt.Fprintf(b, "%s %s `json:\"%s,omitempty\" n2k:\"%d\"`\n",
			fieldName, fieldType, canboatJSONFieldName(field), field.Order)
	}
	fmt.Fprintf(b, "}\n\n")

	fmt.Fprintf(b, "func (m *%s) PGNNumber() uint32 { return %d }\n", typeName, p.PGN)
	fmt.Fprintf(b, "func Decode%s(Info MessageInfo, stream *PGNDataStream) (Message, error) {\n", typeName)
	fmt.Fprintf(b, "val := &%s{Info: Info}\n", typeName)
	fmt.Fprintf(b, "if err := decodeGeneratedConcrete(generatedConcreteInfo(%d, %q), val, stream); err != nil { return nil, err }\n", p.PGN, p.Description)
	fmt.Fprintf(b, "return val, nil\n")
	fmt.Fprintf(b, "}\n\n")

	fmt.Fprintf(b, "func Encode%s(val *%s) ([]byte, error) {\n", typeName, typeName)
	fmt.Fprintf(b, "return encodeGeneratedConcrete(generatedConcreteInfo(%d, %q), val)\n", p.PGN, p.Description)
	fmt.Fprintf(b, "}\n")
	fmt.Fprintf(b, "func encode%sMsg(v Message) ([]byte, error) {\n", typeName)
	fmt.Fprintf(b, "val, ok := v.(*%s)\n", typeName)
	fmt.Fprintf(b, "if !ok { return nil, generatedConcreteTypeError(%q, v) }\n", typeName)
	fmt.Fprintf(b, "return Encode%s(val)\n", typeName)
	fmt.Fprintf(b, "}\n\n")
}

func concreteFieldType(field canboatField) (string, bool) {
	switch field.FieldType {
	case "RESERVED", "SPARE":
		return "", false
	case "STRING_FIX", "STRING_LZ", "STRING_LAU":
		return "string", true
	case "BINARY", "VARIABLE", "DYNAMIC_FIELD_VALUE":
		return "[]uint8", true
	case "FLOAT":
		return "*float32", true
	default:
		if field.Signed {
			return "*int64", true
		}
		return "*uint64", true
	}
}

func uniqueFieldName(used map[string]int, name string) string {
	if name == "" {
		name = "Field"
	}
	if name == "Info" {
		name = "InfoField"
	}
	count := used[name]
	used[name] = count + 1
	if count == 0 {
		return name
	}
	return fmt.Sprintf("%s%d", name, count+1)
}

func exportedGoIdentifier(value string) string {
	parts := strings.FieldsFunc(value, func(r rune) bool {
		return (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && (r < '0' || r > '9')
	})
	var b strings.Builder
	for _, part := range parts {
		if part == "" {
			continue
		}
		b.WriteString(strings.ToUpper(part[:1]))
		if len(part) > 1 {
			b.WriteString(part[1:])
		}
	}
	out := b.String()
	if out == "" {
		return "Value"
	}
	if out[0] >= '0' && out[0] <= '9' {
		return "Pgn" + out
	}
	return out
}

func canboatJSONFieldName(field canboatField) string {
	if field.Id != "" {
		return field.Id
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

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return "value"
}

func concreteCategoryPath(category string) string {
	return fmt.Sprintf("pgn/%s_concrete_generated.go", category)
}

func concreteCategoryVarName(category string) string {
	return "generated" + exportedGoIdentifier(category) + "ConcreteDefinitions"
}

func concreteCategory(p canboatPGN) string {
	description := strings.TrimSpace(p.Description)
	if category, ok := vendorCategory(description); ok {
		return category
	}

	lowerDescription := strings.ToLower(description)
	id := strings.ToLower(p.Id)
	switch {
	case strings.HasPrefix(description, "AIS ") || strings.Contains(description, " AIS "):
		return "ais"
	case strings.HasPrefix(description, "DSC ") || strings.Contains(lowerDescription, "radio"):
		return "communication"
	case containsAny(lowerDescription, "wind", "weather", "meteorological", "temperature", "humidity", "pressure", "precipitation", "station data", "current station", "buoy", "salinity", "water depth", "tide", "set & drift"):
		return "environmental"
	case containsAny(lowerDescription, "position", "navigation", "route", "waypoint", "bearing", "course", "heading", "gnss", "gps", "glonass", "loran", "cog", "sog", "datum", "distance", "direction data", "vessel acceleration"):
		return "navigation"
	case containsAny(lowerDescription, "engine", "transmission", "trip fuel", "fuel economy"):
		return "engine"
	case containsAny(lowerDescription, "electric drive", "thruster", "rudder", "steering", "propulsion"):
		return "propulsion"
	case containsAny(lowerDescription, "ac ", "dc ", "battery", "batteries", "charger", "inverter", "generator", "utility", "power", "energy", "voltage", "current", "load controller", "switch bank", "electrical", "breaker", "lighting"):
		return "electrical"
	case containsAny(lowerDescription, "fluid level", "windlass", "watermaker", "elevator", "actuator", "payload mass"):
		return "equipment"
	case strings.HasPrefix(id, "ais"):
		return "ais"
	case p.PGN < 127000:
		return "system"
	case p.PGN >= 129000 && p.PGN < 130000:
		return "navigation"
	case p.PGN >= 130000 && p.PGN < 131000:
		return "environmental"
	default:
		return "system"
	}
}

func vendorCategory(description string) (string, bool) {
	prefixCategories := []struct {
		Prefix   string
		Category string
	}{
		{"Garmin AHRS ATT:", "garmin"},
		{"BEP Marine:", "bep_marine"},
		{"Diverse Yacht Services:", "diverse_yacht_services"},
		{"Sea Recovery:", "sea_recovery"},
		{"SonicHub:", "sonichub"},
		{"Seatalk1:", "seatalk"},
		{"Seatalk:", "seatalk"},
		{"SeaTalk:", "seatalk"},
		{"Simnet:", "simnet"},
		{"Simrad:", "simrad"},
		{"Lowrance:", "lowrance"},
		{"Lumishore:", "lumishore"},
		{"Maretron:", "maretron"},
		{"Mercury:", "mercury"},
		{"Xantrex:", "xantrex"},
		{"Victron:", "victron"},
		{"Webasto:", "webasto"},
		{"Yamaha:", "yamaha"},
		{"Yanmar:", "yanmar"},
		{"Airmar:", "airmar"},
		{"Carling:", "carling"},
		{"Chetco:", "chetco"},
		{"Fusion:", "fusion"},
		{"Furuno:", "furuno"},
		{"Garmin:", "garmin"},
		{"Honda:", "honda"},
		{"Navico:", "navico"},
		{"Suzuki:", "suzuki"},
		{"B&G:", "bg"},
	}
	for _, prefixCategory := range prefixCategories {
		if strings.HasPrefix(description, prefixCategory.Prefix) {
			return prefixCategory.Category, true
		}
	}
	return "", false
}

func containsAny(value string, needles ...string) bool {
	for _, needle := range needles {
		if strings.Contains(value, needle) {
			return true
		}
	}
	return false
}

func writePGN(b *bytes.Buffer, p canboatPGN) {
	fmt.Fprintf(b, "{PGN:%d, CanboatId:%q, Description:%q, Type:%q, Complete:%t, Fallback:%t, FieldCount:%d",
		p.PGN, p.Id, p.Description, p.Type, p.Complete, p.Fallback, p.FieldCount)
	if p.Explanation != "" {
		fmt.Fprintf(b, ", Explanation:%q", p.Explanation)
	}
	if p.URL != "" {
		fmt.Fprintf(b, ", URL:%q", p.URL)
	}
	if len(p.Missing) > 0 {
		fmt.Fprintf(b, ", Missing:[]string{%s}", quotedStrings(p.Missing))
	}
	if p.Length != nil {
		fmt.Fprintf(b, ", Length:canboatInt(%d)", *p.Length)
	}
	if p.MinLength != nil {
		fmt.Fprintf(b, ", MinLength:canboatInt(%d)", *p.MinLength)
	}
	if p.Priority != nil {
		fmt.Fprintf(b, ", Priority:canboatUint8(%d)", *p.Priority)
	}
	if p.TransmissionInterval != nil {
		fmt.Fprintf(b, ", TransmissionInterval:canboatInt(%d)", *p.TransmissionInterval)
	}
	if p.TransmissionIrregular != nil {
		fmt.Fprintf(b, ", TransmissionIrregular:canboatBool(%t)", *p.TransmissionIrregular)
	}
	writeOptionalInt(b, "RepeatingFieldSet1StartField", p.RepeatingFieldSet1StartField)
	writeOptionalInt(b, "RepeatingFieldSet1CountField", p.RepeatingFieldSet1CountField)
	writeOptionalInt(b, "RepeatingFieldSet1Size", p.RepeatingFieldSet1Size)
	writeOptionalInt(b, "RepeatingFieldSet2StartField", p.RepeatingFieldSet2StartField)
	writeOptionalInt(b, "RepeatingFieldSet2CountField", p.RepeatingFieldSet2CountField)
	writeOptionalInt(b, "RepeatingFieldSet2Size", p.RepeatingFieldSet2Size)
	if len(p.Fields) > 0 {
		fmt.Fprintf(b, ", Fields:[]CanboatFieldDefinition{\n")
		for _, field := range p.Fields {
			writeField(b, field)
		}
		fmt.Fprintf(b, "}")
	}
	fmt.Fprintf(b, "},\n")
}

func writeField(b *bytes.Buffer, f canboatField) {
	fmt.Fprintf(b, "{Order:%d, CanboatId:%q, Name:%q, Signed:%t, Unit:%q, FieldType:%q",
		f.Order, f.Id, f.Name, f.Signed, f.Unit, f.FieldType)
	if f.Description != "" {
		fmt.Fprintf(b, ", Description:%q", string(f.Description))
	}
	if f.BitLength != nil {
		fmt.Fprintf(b, ", BitLength:canboatUint16(%d)", *f.BitLength)
	}
	if f.BitLengthField != nil {
		fmt.Fprintf(b, ", BitLengthField:canboatInt(%d)", *f.BitLengthField)
	}
	if f.BitLengthVariable {
		fmt.Fprintf(b, ", BitLengthVariable:true")
	}
	if f.BitOffset != nil {
		fmt.Fprintf(b, ", BitOffset:canboatUint16(%d)", *f.BitOffset)
	}
	if f.BitStart != nil {
		fmt.Fprintf(b, ", BitStart:canboatUint16(%d)", *f.BitStart)
	}
	if f.Resolution != nil {
		fmt.Fprintf(b, ", Resolution:canboatFloat64(%s)", floatLiteral(*f.Resolution))
	}
	if f.PhysicalQuantity != "" {
		fmt.Fprintf(b, ", PhysicalQuantity:%q", f.PhysicalQuantity)
	}
	if f.LookupEnumeration != "" {
		fmt.Fprintf(b, ", LookupEnumeration:%q", f.LookupEnumeration)
	}
	if f.LookupBitEnumeration != "" {
		fmt.Fprintf(b, ", LookupBitEnumeration:%q", f.LookupBitEnumeration)
	}
	if f.LookupIndirectEnumeration != "" {
		fmt.Fprintf(b, ", LookupIndirectEnumeration:%q", f.LookupIndirectEnumeration)
	}
	if f.LookupIndirectEnumerationFieldOrder != nil {
		fmt.Fprintf(b, ", LookupIndirectEnumerationFieldOrder:canboatInt(%d)", *f.LookupIndirectEnumerationFieldOrder)
	}
	if f.LookupFieldTypeEnumeration != "" {
		fmt.Fprintf(b, ", LookupFieldTypeEnumeration:%q", f.LookupFieldTypeEnumeration)
	}
	if f.Match != nil {
		fmt.Fprintf(b, ", Match:canboatMatch(%d)", *f.Match)
	}
	if f.RangeMin != nil {
		fmt.Fprintf(b, ", RangeMin:canboatFloat64(%s)", floatLiteral(*f.RangeMin))
	}
	if f.RangeMax != nil {
		fmt.Fprintf(b, ", RangeMax:canboatFloat64(%s)", floatLiteral(*f.RangeMax))
	}
	if f.Offset != nil {
		fmt.Fprintf(b, ", Offset:canboatFloat64(%s)", floatLiteral(*f.Offset))
	}
	if f.OutOfRangeValue != nil {
		fmt.Fprintf(b, ", OutOfRangeValue:canboatInt64(%d)", *f.OutOfRangeValue)
	}
	if f.PartOfPrimaryKey != nil {
		fmt.Fprintf(b, ", PartOfPrimaryKey:canboatBool(%t)", *f.PartOfPrimaryKey)
	}
	if f.ReservedValue != nil {
		fmt.Fprintf(b, ", ReservedValue:canboatInt64(%d)", *f.ReservedValue)
	}
	if f.UnknownValue != nil {
		fmt.Fprintf(b, ", UnknownValue:canboatInt64(%d)", *f.UnknownValue)
	}
	if f.Condition != "" {
		fmt.Fprintf(b, ", Condition:%q", f.Condition)
	}
	fmt.Fprintf(b, "},\n")
}

func writeOptionalInt(b *bytes.Buffer, name string, value *int) {
	if value == nil {
		return
	}
	fmt.Fprintf(b, ", %s:canboatInt(%d)", name, *value)
}

func quotedStrings(values []string) string {
	quoted := make([]string, 0, len(values))
	for _, value := range values {
		quoted = append(quoted, strconv.Quote(value))
	}
	return strings.Join(quoted, ",")
}

func floatLiteral(v float64) string {
	return strconv.FormatFloat(v, 'g', -1, 64)
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
