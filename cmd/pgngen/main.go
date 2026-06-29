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

const definitionsPath = "pgn/upstream_definitions.go"
const dispatchPath = "pgn/dispatch.go"

var upstreamURL = upstreamSchemaURL()

func upstreamSchemaURL() string {
	project := "can" + "boat"
	return "https://raw.githubusercontent.com/" + project + "/" + project + "/refs/heads/master/docs/" + project + ".json"
}

type sourceFile struct {
	SchemaVersion string      `json:"SchemaVersion"`
	Version       string      `json:"Version"`
	PGNs          []sourcePGN `json:"PGNs"`
}

type sourcePGN struct {
	PGN                          uint32        `json:"PGN"`
	Id                           string        `json:"Id"`
	Description                  string        `json:"Description"`
	Explanation                  string        `json:"Explanation"`
	URL                          string        `json:"URL"`
	Type                         string        `json:"Type"`
	Complete                     bool          `json:"Complete"`
	Fallback                     bool          `json:"Fallback"`
	Missing                      []string      `json:"Missing"`
	FieldCount                   int           `json:"FieldCount"`
	Length                       *int          `json:"Length"`
	MinLength                    *int          `json:"MinLength"`
	Priority                     *uint8        `json:"Priority"`
	TransmissionInterval         *int          `json:"TransmissionInterval"`
	TransmissionIrregular        *bool         `json:"TransmissionIrregular"`
	RepeatingFieldSet1StartField *int          `json:"RepeatingFieldSet1StartField"`
	RepeatingFieldSet1CountField *int          `json:"RepeatingFieldSet1CountField"`
	RepeatingFieldSet1Size       *int          `json:"RepeatingFieldSet1Size"`
	RepeatingFieldSet2StartField *int          `json:"RepeatingFieldSet2StartField"`
	RepeatingFieldSet2CountField *int          `json:"RepeatingFieldSet2CountField"`
	RepeatingFieldSet2Size       *int          `json:"RepeatingFieldSet2Size"`
	Fields                       []sourceField `json:"Fields"`
}

type sourceField struct {
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

type pgnOutput struct {
	PGN      sourcePGN
	TypeName string
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
	check := flag.Bool("check", false, "fail if PGN outputs are not current")
	flag.Parse()

	raw, err := fetch(upstreamURL)
	if err != nil {
		fatal(err)
	}
	var source sourceFile
	if err := json.Unmarshal(raw, &source); err != nil {
		fatal(err)
	}
	definitions, err := generateDefinitions(source)
	if err != nil {
		fatal(err)
	}
	pgns, err := generatePGNs(source)
	if err != nil {
		fatal(err)
	}
	pgns[definitionsPath] = definitions
	if *check {
		if err := checkPGNFiles(pgns); err != nil {
			fatal(err)
		}
		return
	}
	if err := writePGNFiles(pgns); err != nil {
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

func generateDefinitions(source sourceFile) ([]byte, error) {
	sort.SliceStable(source.PGNs, func(i, j int) bool {
		if source.PGNs[i].PGN != source.PGNs[j].PGN {
			return source.PGNs[i].PGN < source.PGNs[j].PGN
		}
		return source.PGNs[i].Description < source.PGNs[j].Description
	})

	var b bytes.Buffer
	fmt.Fprintf(&b, "// Code generated by PGN generator; DO NOT EDIT.\n")
	fmt.Fprintf(&b, "// Source schema %s version %s.\n\n", source.SchemaVersion, source.Version)
	fmt.Fprintf(&b, "package pgn\n\n")
	fmt.Fprintf(&b, "var sourceDefinitions = []SourceDefinition{\n")
	names := uniquePGNTypeNames(source.PGNs)
	for i, pgn := range source.PGNs {
		writePGN(&b, pgn, names[i])
	}
	fmt.Fprintf(&b, "}\n")
	return format.Source(b.Bytes())
}

func generatePGNs(source sourceFile) (map[string][]byte, error) {
	sort.SliceStable(source.PGNs, func(i, j int) bool {
		if source.PGNs[i].PGN != source.PGNs[j].PGN {
			return source.PGNs[i].PGN < source.PGNs[j].PGN
		}
		return source.PGNs[i].Description < source.PGNs[j].Description
	})

	names := uniquePGNTypeNames(source.PGNs)
	categories := make(map[string][]pgnOutput)
	var categoryOrder []string
	for i, pgn := range source.PGNs {
		category := pgnCategory(pgn)
		if _, ok := categories[category]; !ok {
			categoryOrder = append(categoryOrder, category)
		}
		categories[category] = append(categories[category], pgnOutput{
			PGN:      pgn,
			TypeName: names[i],
		})
	}

	outputs := make(map[string][]byte, len(categories)+1)
	dispatch, err := generatePGNDispatch(source, source.PGNs, names)
	if err != nil {
		return nil, err
	}
	outputs[dispatchPath] = dispatch

	for _, category := range categoryOrder {
		out, err := generatePGNCategory(source, category, categories[category])
		if err != nil {
			return nil, err
		}
		outputs[pgnCategoryPath(category)] = out
	}
	return outputs, nil
}

func generatePGNDispatch(source sourceFile, pgns []sourcePGN, names []string) ([]byte, error) {
	var b bytes.Buffer
	fmt.Fprintf(&b, "// Code generated by PGN generator; DO NOT EDIT.\n")
	fmt.Fprintf(&b, "// Source schema %s version %s.\n\n", source.SchemaVersion, source.Version)
	fmt.Fprintf(&b, "package pgn\n\n")
	fmt.Fprintf(&b, "import (\n")
	fmt.Fprintf(&b, "%q\n", "errors")
	fmt.Fprintf(&b, "%q\n", "fmt")
	fmt.Fprintf(&b, ")\n\n")
	fmt.Fprintf(&b, "func DecodePayload(info MessageInfo, payload []uint8) (PGN, error) {\n")
	fmt.Fprintf(&b, "switch info.PGN {\n")
	byPGN := make(map[uint32][]string)
	var pgnOrder []uint32
	for i, pgn := range pgns {
		if _, ok := byPGN[pgn.PGN]; !ok {
			pgnOrder = append(pgnOrder, pgn.PGN)
		}
		byPGN[pgn.PGN] = append(byPGN[pgn.PGN], names[i])
	}
	sort.Slice(pgnOrder, func(i, j int) bool { return pgnOrder[i] < pgnOrder[j] })
	for _, pgn := range pgnOrder {
		fmt.Fprintf(&b, "case %d:\n", pgn)
		fmt.Fprintf(&b, "return decodePGNCandidates(info, payload,\n")
		for _, name := range byPGN[pgn] {
			fmt.Fprintf(&b, "&%s{Info: info},\n", name)
		}
		fmt.Fprintf(&b, ")\n")
	}
	fmt.Fprintf(&b, "default:\n")
	fmt.Fprintf(&b, "return nil, fmt.Errorf(\"no PGN struct for %%d\", info.PGN)\n")
	fmt.Fprintf(&b, "}\n")
	fmt.Fprintf(&b, "}\n\n")
	fmt.Fprintf(&b, "func decodePGNCandidates(info MessageInfo, payload []uint8, candidates ...PGN) (PGN, error) {\n")
	fmt.Fprintf(&b, "var errs []error\n")
	fmt.Fprintf(&b, "for _, candidate := range candidates {\n")
	fmt.Fprintf(&b, "candidate.SetMessageInfo(info)\n")
	fmt.Fprintf(&b, "if err := candidate.DecodePayload(payload); err != nil {\n")
	fmt.Fprintf(&b, "errs = append(errs, err)\n")
	fmt.Fprintf(&b, "continue\n")
	fmt.Fprintf(&b, "}\n")
	fmt.Fprintf(&b, "return candidate, nil\n")
	fmt.Fprintf(&b, "}\n")
	fmt.Fprintf(&b, "if len(errs) == 0 {\n")
	fmt.Fprintf(&b, "return nil, fmt.Errorf(\"no PGN struct for %%d\", info.PGN)\n")
	fmt.Fprintf(&b, "}\n")
	fmt.Fprintf(&b, "return nil, fmt.Errorf(\"no matching PGN struct for %%d: %%w\", info.PGN, errors.Join(errs...))\n")
	fmt.Fprintf(&b, "}\n")
	return format.Source(b.Bytes())
}

func generatePGNCategory(source sourceFile, category string, entries []pgnOutput) ([]byte, error) {
	var b bytes.Buffer
	fmt.Fprintf(&b, "// Code generated by PGN generator; DO NOT EDIT.\n")
	fmt.Fprintf(&b, "// Source schema %s version %s.\n\n", source.SchemaVersion, source.Version)
	fmt.Fprintf(&b, "package pgn\n\n")
	for _, entry := range entries {
		writePGNStruct(&b, entry.PGN, entry.TypeName)
	}
	return format.Source(b.Bytes())
}

func writePGNFiles(outputs map[string][]byte) error {
	if err := removeStalePGNFiles(outputs); err != nil {
		return err
	}
	paths := sortedPGNPaths(outputs)
	for _, path := range paths {
		if err := os.WriteFile(path, outputs[path], 0o600); err != nil {
			return err
		}
	}
	return nil
}

func checkPGNFiles(outputs map[string][]byte) error {
	if err := checkStalePGNFiles(outputs); err != nil {
		return err
	}
	paths := sortedPGNPaths(outputs)
	for _, path := range paths {
		if err := validateGeneratedPath(path); err != nil {
			return err
		}
		current, err := os.ReadFile(path) // #nosec G304 -- path is restricted to known PGN output files under pgn/.
		if err != nil {
			return err
		}
		if !bytes.Equal(current, outputs[path]) {
			return fmt.Errorf("%s is not synced with %s", path, upstreamURL)
		}
	}
	return nil
}

func validateGeneratedPath(path string) error {
	clean := filepath.Clean(path)
	if path != clean || filepath.IsAbs(clean) || clean == ".." || strings.HasPrefix(clean, ".."+string(os.PathSeparator)) {
		return fmt.Errorf("invalid PGN output path %q", path)
	}
	if filepath.Dir(clean) != "pgn" {
		return fmt.Errorf("PGN output path %q must be under pgn", path)
	}
	base := filepath.Base(clean)
	if strings.Contains(base, "generated") {
		return fmt.Errorf("PGN output path %q must not contain generated", path)
	}
	if base == filepath.Base(definitionsPath) || base == filepath.Base(dispatchPath) {
		return nil
	}
	if strings.HasSuffix(base, ".go") && !strings.HasSuffix(base, "_test.go") {
		return nil
	}
	return fmt.Errorf("unexpected PGN output path %q", path)
}

func removeStalePGNFiles(outputs map[string][]byte) error {
	paths, err := existingPGNPaths()
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

func checkStalePGNFiles(outputs map[string][]byte) error {
	paths, err := existingPGNPaths()
	if err != nil {
		return err
	}
	for _, path := range paths {
		if _, ok := outputs[path]; !ok {
			return fmt.Errorf("%s is stale PGN output", path)
		}
	}
	return nil
}

func existingPGNPaths() ([]string, error) {
	var paths []string
	legacyMatches, err := filepath.Glob("pgn/*generated*.go")
	if err != nil {
		return nil, err
	}
	paths = append(paths, legacyMatches...)
	for _, path := range []string{definitionsPath, dispatchPath} {
		if _, err := os.Stat(path); err == nil {
			paths = append(paths, path)
		} else if err != nil && !os.IsNotExist(err) {
			return nil, err
		}
	}
	sort.Strings(paths)
	return paths, nil
}

func sortedPGNPaths(outputs map[string][]byte) []string {
	paths := make([]string, 0, len(outputs))
	for path := range outputs {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	return paths
}

func uniquePGNTypeNames(pgns []sourcePGN) []string {
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

func writePGNStruct(b *bytes.Buffer, p sourcePGN, typeName string) {
	fmt.Fprintf(b, "type %s struct {\n", typeName)
	fmt.Fprintf(b, "Info MessageInfo `json:\"info\"`\n")

	fieldNames := make(map[string]int)
	for _, field := range p.Fields {
		fieldType, ok := pgnFieldType(field)
		if !ok {
			continue
		}
		fieldName := uniqueFieldName(fieldNames, exportedGoIdentifier(firstNonEmpty(field.Id, field.Name)))
		fmt.Fprintf(b, "%s %s `json:\"%s,omitempty\" n2k:\"%d\"`\n",
			fieldName, fieldType, sourceJSONFieldName(field), field.Order)
	}
	fmt.Fprintf(b, "}\n\n")

	fmt.Fprintf(b, "func (m *%s) PGNNumber() uint32 { return %d }\n", typeName, p.PGN)
	fmt.Fprintf(b, "func (m *%s) MessageInfo() MessageInfo { return m.Info }\n", typeName)
	fmt.Fprintf(b, "func (m *%s) SetMessageInfo(info MessageInfo) { m.Info = info }\n", typeName)
	fmt.Fprintf(b, "func (m *%s) DecodePayload(payload []uint8) error {\n", typeName)
	fmt.Fprintf(b, "return decodeStructPayload(lookupPgnInfo(%d, %q), m, payload)\n", p.PGN, p.Description)
	fmt.Fprintf(b, "}\n")
	fmt.Fprintf(b, "func (m *%s) EncodePayload() ([]uint8, error) {\n", typeName)
	fmt.Fprintf(b, "return encodeStructPayload(lookupPgnInfo(%d, %q), m)\n", p.PGN, p.Description)
	fmt.Fprintf(b, "}\n\n")
}

func pgnFieldType(field sourceField) (string, bool) {
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

func sourceJSONFieldName(field sourceField) string {
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

func pgnCategoryPath(category string) string {
	return fmt.Sprintf("pgn/%s.go", category)
}

func pgnCategory(p sourcePGN) string {
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

func writePGN(b *bytes.Buffer, p sourcePGN, typeName string) {
	fmt.Fprintf(b, "{PGN:%d, StructName:%q, SourceID:%q, Description:%q, Type:%q, Complete:%t, Fallback:%t, FieldCount:%d",
		p.PGN, typeName, p.Id, p.Description, p.Type, p.Complete, p.Fallback, p.FieldCount)
	if p.Explanation != "" {
		fmt.Fprintf(b, ", Explanation:%q", p.Explanation)
	}
	if len(p.Missing) > 0 {
		fmt.Fprintf(b, ", Missing:[]string{%s}", quotedStrings(p.Missing))
	}
	if p.Length != nil {
		fmt.Fprintf(b, ", Length:sourceInt(%d)", *p.Length)
	}
	if p.MinLength != nil {
		fmt.Fprintf(b, ", MinLength:sourceInt(%d)", *p.MinLength)
	}
	if p.Priority != nil {
		fmt.Fprintf(b, ", Priority:sourceUint8(%d)", *p.Priority)
	}
	if p.TransmissionInterval != nil {
		fmt.Fprintf(b, ", TransmissionInterval:sourceInt(%d)", *p.TransmissionInterval)
	}
	if p.TransmissionIrregular != nil {
		fmt.Fprintf(b, ", TransmissionIrregular:sourceBool(%t)", *p.TransmissionIrregular)
	}
	writeOptionalInt(b, "RepeatingFieldSet1StartField", p.RepeatingFieldSet1StartField)
	writeOptionalInt(b, "RepeatingFieldSet1CountField", p.RepeatingFieldSet1CountField)
	writeOptionalInt(b, "RepeatingFieldSet1Size", p.RepeatingFieldSet1Size)
	writeOptionalInt(b, "RepeatingFieldSet2StartField", p.RepeatingFieldSet2StartField)
	writeOptionalInt(b, "RepeatingFieldSet2CountField", p.RepeatingFieldSet2CountField)
	writeOptionalInt(b, "RepeatingFieldSet2Size", p.RepeatingFieldSet2Size)
	if len(p.Fields) > 0 {
		fmt.Fprintf(b, ", Fields:[]SourceFieldDefinition{\n")
		for _, field := range p.Fields {
			writeField(b, field)
		}
		fmt.Fprintf(b, "}")
	}
	fmt.Fprintf(b, "},\n")
}

func writeField(b *bytes.Buffer, f sourceField) {
	fmt.Fprintf(b, "{Order:%d, SourceID:%q, Name:%q, Signed:%t, Unit:%q, FieldType:%q",
		f.Order, f.Id, f.Name, f.Signed, f.Unit, f.FieldType)
	if f.Description != "" {
		fmt.Fprintf(b, ", Description:%q", string(f.Description))
	}
	if f.BitLength != nil {
		fmt.Fprintf(b, ", BitLength:sourceUint16(%d)", *f.BitLength)
	}
	if f.BitLengthField != nil {
		fmt.Fprintf(b, ", BitLengthField:sourceInt(%d)", *f.BitLengthField)
	}
	if f.BitLengthVariable {
		fmt.Fprintf(b, ", BitLengthVariable:true")
	}
	if f.BitOffset != nil {
		fmt.Fprintf(b, ", BitOffset:sourceUint16(%d)", *f.BitOffset)
	}
	if f.BitStart != nil {
		fmt.Fprintf(b, ", BitStart:sourceUint16(%d)", *f.BitStart)
	}
	if f.Resolution != nil {
		fmt.Fprintf(b, ", Resolution:sourceFloat64(%s)", floatLiteral(*f.Resolution))
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
		fmt.Fprintf(b, ", LookupIndirectEnumerationFieldOrder:sourceInt(%d)", *f.LookupIndirectEnumerationFieldOrder)
	}
	if f.LookupFieldTypeEnumeration != "" {
		fmt.Fprintf(b, ", LookupFieldTypeEnumeration:%q", f.LookupFieldTypeEnumeration)
	}
	if f.Match != nil {
		fmt.Fprintf(b, ", Match:sourceMatch(%d)", *f.Match)
	}
	if f.RangeMin != nil {
		fmt.Fprintf(b, ", RangeMin:sourceFloat64(%s)", floatLiteral(*f.RangeMin))
	}
	if f.RangeMax != nil {
		fmt.Fprintf(b, ", RangeMax:sourceFloat64(%s)", floatLiteral(*f.RangeMax))
	}
	if f.Offset != nil {
		fmt.Fprintf(b, ", Offset:sourceFloat64(%s)", floatLiteral(*f.Offset))
	}
	if f.OutOfRangeValue != nil {
		fmt.Fprintf(b, ", OutOfRangeValue:sourceInt64(%d)", *f.OutOfRangeValue)
	}
	if f.PartOfPrimaryKey != nil {
		fmt.Fprintf(b, ", PartOfPrimaryKey:sourceBool(%t)", *f.PartOfPrimaryKey)
	}
	if f.ReservedValue != nil {
		fmt.Fprintf(b, ", ReservedValue:sourceInt64(%d)", *f.ReservedValue)
	}
	if f.UnknownValue != nil {
		fmt.Fprintf(b, ", UnknownValue:sourceInt64(%d)", *f.UnknownValue)
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
	fmt.Fprintf(b, ", %s:sourceInt(%d)", name, *value)
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
