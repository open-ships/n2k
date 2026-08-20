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

// upstreamRevision pins generation to the schema version represented by the
// checked-in public API. The separate upstream-parity test intentionally
// follows canboat master so maintainers can review schema upgrades explicitly.
const upstreamRevision = "8f737e93ba2a1dd8684d4bc267126c55173c3ee5" // canboat schema 7.1.0

var upstreamURL = upstreamSchemaURL()

func upstreamSchemaURL() string {
	project := "can" + "boat"
	return "https://raw.githubusercontent.com/" + project + "/" + project + "/" + upstreamRevision + "/docs/" + project + ".json"
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
	Type                         string        `json:"Type"`
	Complete                     bool          `json:"Complete"`
	Fallback                     bool          `json:"Fallback"`
	Missing                      []string      `json:"Missing"`
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

type pgnStructField struct {
	Source sourceField
	Name   string
	Type   string
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
	// The definitions are assembled from per-chunk functions rather than
	// declared as one literal: a single ~6k-line composite literal lowers
	// into one enormous initializer function, which the race-instrumented
	// compiler cannot build ("NewBulk too big"). A var initializer (not
	// init()) keeps Go's dependency-ordered initialization for the metadata
	// tables built from these definitions.
	names := uniquePGNTypeNames(source.PGNs)
	const chunkSize = 100
	chunks := (len(source.PGNs) + chunkSize - 1) / chunkSize
	fmt.Fprintf(&b, "var sourceDefinitions = buildSourceDefinitions()\n\n")
	fmt.Fprintf(&b, "func buildSourceDefinitions() []SourceDefinition {\n")
	fmt.Fprintf(&b, "var defs []SourceDefinition\n")
	for chunk := 0; chunk < chunks; chunk++ {
		fmt.Fprintf(&b, "defs = append(defs, sourceDefinitionsChunk%d()...)\n", chunk)
	}
	fmt.Fprintf(&b, "return defs\n}\n\n")
	for chunk := 0; chunk < chunks; chunk++ {
		start := chunk * chunkSize
		end := start + chunkSize
		if end > len(source.PGNs) {
			end = len(source.PGNs)
		}
		fmt.Fprintf(&b, "func sourceDefinitionsChunk%d() []SourceDefinition {\n", chunk)
		fmt.Fprintf(&b, "return []SourceDefinition{\n")
		for i := start; i < end; i++ {
			writePGN(&b, source.PGNs[i], names[i])
		}
		fmt.Fprintf(&b, "}\n}\n\n")
	}
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
	fmt.Fprintf(&b, "}\n\n")
	fmt.Fprintf(&b, "// structTypeRegistry instantiates a PGN struct by its type name.\n")
	fmt.Fprintf(&b, "var structTypeRegistry = map[string]func() PGN{\n")
	sortedNames := append([]string(nil), names...)
	sort.Strings(sortedNames)
	for _, name := range sortedNames {
		fmt.Fprintf(&b, "%q: func() PGN { return &%s{} },\n", name, name)
	}
	fmt.Fprintf(&b, "}\n")
	return format.Source(b.Bytes())
}

func generatePGNCategory(source sourceFile, category string, entries []pgnOutput) ([]byte, error) {
	var body bytes.Buffer
	needsMath := false
	for _, entry := range entries {
		if writePGNStruct(&body, entry.PGN, entry.TypeName) {
			needsMath = true
		}
	}
	var b bytes.Buffer
	fmt.Fprintf(&b, "// Code generated by PGN generator; DO NOT EDIT.\n")
	fmt.Fprintf(&b, "// Source schema %s version %s.\n\n", source.SchemaVersion, source.Version)
	fmt.Fprintf(&b, "package pgn\n\n")
	if needsMath {
		fmt.Fprintf(&b, "import \"math\"\n\n")
	}
	b.Write(body.Bytes())
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
		base := pgnStructTypeName(firstNonEmpty(pgn.Id, pgn.Description))
		count := used[base]
		used[base] = count + 1
		if count > 0 {
			base = fmt.Sprintf("%s%d", base, count+1)
		}
		names[i] = base
	}
	return names
}

func pgnStructTypeName(value string) string {
	name := exportedGoIdentifier(value)
	if strings.HasPrefix(strings.ToLower(strings.TrimSpace(value)), "pgn") && strings.HasPrefix(name, "Pgn") {
		name = "ParameterGroupNumber" + strings.TrimPrefix(name, "Pgn")
	}
	if strings.HasPrefix(name, "PGN") {
		name = "ParameterGroupNumber" + strings.TrimPrefix(name, "PGN")
	}
	name = strings.TrimPrefix(name, "Pgn")
	if name == "" {
		name = "Message"
	}
	if name[0] >= '0' && name[0] <= '9' {
		name = "Message" + name
	}
	return name
}

func writePGNStruct(b *bytes.Buffer, p sourcePGN, typeName string) bool {
	fixed, groups := partitionPGNStructFields(p)

	fmt.Fprintf(b, "type %s struct {\n", typeName)
	fmt.Fprintf(b, "Info MessageInfo `json:\"info\"`\n")
	for _, field := range fixed {
		fmt.Fprintf(b, "%s %s `json:\"%s,omitempty\" n2k:\"%d\"`\n",
			field.Name, field.Type, sourceJSONFieldName(field.Source), field.Source.Order)
	}
	for setIndex, group := range groups {
		if len(group) == 0 {
			continue
		}
		set := setIndex + 1
		fmt.Fprintf(b, "Repeating%d []%sRepeating%d `json:\"repeating%d,omitempty\" n2k:\"rep%d\"`\n",
			set, typeName, set, set, set)
	}
	fmt.Fprintf(b, "}\n\n")

	for setIndex, group := range groups {
		if len(group) == 0 {
			continue
		}
		fmt.Fprintf(b, "type %sRepeating%d struct {\n", typeName, setIndex+1)
		for _, field := range group {
			fmt.Fprintf(b, "%s %s `json:\"%s,omitempty\" n2k:\"%d\"`\n",
				field.Name, field.Type, sourceJSONFieldName(field.Source), field.Source.Order)
		}
		fmt.Fprintf(b, "}\n\n")
	}

	fmt.Fprintf(b, "func (m *%s) PGNNumber() uint32 { return %d }\n", typeName, p.PGN)
	fmt.Fprintf(b, "func (m *%s) MessageInfo() MessageInfo { return m.Info }\n", typeName)
	fmt.Fprintf(b, "func (m *%s) SetMessageInfo(info MessageInfo) { m.Info = info }\n", typeName)
	fmt.Fprintf(b, "func (m *%s) DecodePayload(payload []uint8) error { return decodeFields(m, payload) }\n", typeName)
	fmt.Fprintf(b, "func (m *%s) EncodePayload() ([]uint8, error) { return encodeFields(m) }\n\n", typeName)

	wroteAccessors := writePhysicalAccessors(b, typeName, fixed, groups)

	for setIndex, group := range groups {
		if len(group) == 0 {
			continue
		}
		elementType := fmt.Sprintf("%sRepeating%d", typeName, setIndex+1)
		if writePhysicalAccessors(b, elementType, group, [2][]pgnStructField{}) {
			wroteAccessors = true
		}
	}
	return wroteAccessors
}

// physicalAccessor describes one generated <Field>Value / Set<Field>Value
// method pair: a plain numeric raw-ticks field whose metadata gives it a
// physical interpretation.
type physicalAccessor struct {
	field      pgnStructField
	resolution float64
	offset     float64
	unit       string
}

// physicalAccessors selects the fields of a struct that warrant typed
// physical-value accessors: plain numeric raw-ticks fields (not lookups,
// match selectors, floats, strings, or binary data) whose metadata carries a
// unit label, a non-unity resolution, or an additive offset. Accessors whose
// method names would collide with a struct field or another method are
// skipped.
func physicalAccessors(fields []pgnStructField, reserved map[string]bool) []physicalAccessor {
	var accessors []physicalAccessor
	for _, field := range fields {
		if field.Type != "*uint64" && field.Type != "*int64" {
			continue
		}
		src := field.Source
		if src.Match != nil || src.LookupEnumeration != "" || src.LookupBitEnumeration != "" ||
			src.LookupIndirectEnumeration != "" || src.LookupFieldTypeEnumeration != "" {
			continue
		}
		resolution := 1.0
		if src.Resolution != nil {
			resolution = *src.Resolution
		}
		offset := 0.0
		if src.Offset != nil {
			offset = *src.Offset
		}
		if src.Unit == "" && resolution == 1 && offset == 0 {
			continue
		}
		if resolution == 0 {
			continue
		}
		getter, setter := field.Name+"Value", "Set"+field.Name+"Value"
		if reserved[getter] || reserved[setter] {
			continue
		}
		reserved[getter] = true
		reserved[setter] = true
		accessors = append(accessors, physicalAccessor{
			field:      field,
			resolution: resolution,
			offset:     offset,
			unit:       src.Unit,
		})
	}
	return accessors
}

// writePhysicalAccessors emits physical-value getter/setter pairs for the
// eligible fields of one struct. It returns true when at least one setter
// was written (the setters use math.Round, so the file needs the math
// import).
func writePhysicalAccessors(b *bytes.Buffer, typeName string, fields []pgnStructField, groups [2][]pgnStructField) bool {
	reserved := reservedMethodNames(fields, groups)
	accessors := physicalAccessors(fields, reserved)
	for _, acc := range accessors {
		writePhysicalGetter(b, typeName, acc)
		writePhysicalSetter(b, typeName, acc)
	}
	return len(accessors) > 0
}

// reservedMethodNames collects the identifiers a generated accessor must not
// collide with: the struct's own field names (methods and fields share a
// namespace in Go) and the fixed interface methods every PGN struct carries.
func reservedMethodNames(fields []pgnStructField, groups [2][]pgnStructField) map[string]bool {
	reserved := map[string]bool{
		"Info":           true,
		"PGNNumber":      true,
		"MessageInfo":    true,
		"SetMessageInfo": true,
		"DecodePayload":  true,
		"EncodePayload":  true,
	}
	for _, field := range fields {
		reserved[field.Name] = true
	}
	for i, group := range groups {
		if len(group) > 0 {
			reserved[fmt.Sprintf("Repeating%d", i+1)] = true
		}
	}
	return reserved
}

func writePhysicalGetter(b *bytes.Buffer, typeName string, acc physicalAccessor) {
	name := acc.field.Name
	fmt.Fprintf(b, "// %sValue returns %s as a physical value%s (value = %s).\n",
		name, name, unitClause(acc.unit), physicalFormula(acc))
	fmt.Fprintf(b, "// The bool is false when %s is nil: the wire carried the field's null\n", name)
	fmt.Fprintf(b, "// sentinel or the payload ended before the field.\n")
	fmt.Fprintf(b, "func (m *%s) %sValue() (float64, bool) {\n", typeName, name)
	fmt.Fprintf(b, "if m.%s == nil {\nreturn 0, false\n}\n", name)
	fmt.Fprintf(b, "return %s, true\n", getterExpression(acc))
	fmt.Fprintf(b, "}\n\n")
}

func writePhysicalSetter(b *bytes.Buffer, typeName string, acc physicalAccessor) {
	name := acc.field.Name
	fmt.Fprintf(b, "// Set%sValue sets %s from a physical value%s, rounded to the nearest\n",
		name, name, unitClause(acc.unit))
	fmt.Fprintf(b, "// wire tick of %s.\n", floatLiteral(acc.resolution))
	fmt.Fprintf(b, "func (m *%s) Set%sValue(v float64) {\n", typeName, name)
	inner := "v"
	switch {
	case acc.offset > 0:
		inner = fmt.Sprintf("(v - %s)", floatLiteral(acc.offset))
	case acc.offset < 0:
		inner = fmt.Sprintf("(v + %s)", floatLiteral(-acc.offset))
	}
	rounded := fmt.Sprintf("math.Round(%s)", inner)
	if acc.resolution != 1 {
		rounded = fmt.Sprintf("math.Round(%s / %s)", inner, floatLiteral(acc.resolution))
	}
	castType := "uint64"
	if acc.field.Type == "*int64" {
		castType = "int64"
	}
	fmt.Fprintf(b, "raw := %s(%s)\n", castType, rounded)
	fmt.Fprintf(b, "m.%s = &raw\n", name)
	fmt.Fprintf(b, "}\n\n")
}

// getterExpression builds the raw-ticks-to-physical conversion expression for
// a getter body, omitting no-op scale and offset terms.
func getterExpression(acc physicalAccessor) string {
	expr := fmt.Sprintf("float64(*m.%s)", acc.field.Name)
	if acc.resolution != 1 {
		expr = fmt.Sprintf("%s * %s", expr, floatLiteral(acc.resolution))
	}
	switch {
	case acc.offset > 0:
		expr = fmt.Sprintf("%s + %s", expr, floatLiteral(acc.offset))
	case acc.offset < 0:
		expr = fmt.Sprintf("%s - %s", expr, floatLiteral(-acc.offset))
	}
	return expr
}

// physicalFormula renders the conversion as documentation text.
func physicalFormula(acc physicalAccessor) string {
	formula := "raw"
	if acc.resolution != 1 {
		formula = fmt.Sprintf("raw * %s", floatLiteral(acc.resolution))
	}
	switch {
	case acc.offset > 0:
		formula = fmt.Sprintf("%s + %s", formula, floatLiteral(acc.offset))
	case acc.offset < 0:
		formula = fmt.Sprintf("%s - %s", formula, floatLiteral(-acc.offset))
	}
	return formula
}

// unitClause renders a field unit for accessor documentation.
func unitClause(unit string) string {
	if unit == "" {
		return ""
	}
	return fmt.Sprintf(" in %s", unit)
}

// partitionPGNStructFields splits a PGN's supported fields into the parent
// struct's fixed fields and up to two repeating-set element field lists,
// using the source RepeatingFieldSet annotations. Field names are unique
// within each struct; the parent namespace reserves RepeatingN names for the
// slice fields of non-empty sets.
func partitionPGNStructFields(p sourcePGN) ([]pgnStructField, [2][]pgnStructField) {
	ranges := pgnRepeatingRanges(p)

	var fixedSources []sourceField
	var groupSources [2][]sourceField
	for _, field := range orderedSourceFields(p.Fields) {
		if _, ok := pgnFieldType(field); !ok {
			continue
		}
		assigned := false
		for i, r := range ranges {
			if r[1] > r[0] && field.Order >= r[0] && field.Order < r[1] {
				groupSources[i] = append(groupSources[i], field)
				assigned = true
				break
			}
		}
		if !assigned {
			fixedSources = append(fixedSources, field)
		}
	}

	fixedNames := make(map[string]int)
	for i := range groupSources {
		if len(groupSources[i]) > 0 {
			uniqueFieldName(fixedNames, fmt.Sprintf("Repeating%d", i+1))
		}
	}
	fixed := namedPGNStructFields(fixedSources, fixedNames)
	var groups [2][]pgnStructField
	for i := range groupSources {
		groups[i] = namedPGNStructFields(groupSources[i], make(map[string]int))
	}
	return fixed, groups
}

// pgnRepeatingRanges returns each repeating field set's [start, end) order
// range, or a zero range when the set is absent.
func pgnRepeatingRanges(p sourcePGN) [2][2]int {
	var ranges [2][2]int
	specs := [2]struct{ start, size *int }{
		{p.RepeatingFieldSet1StartField, p.RepeatingFieldSet1Size},
		{p.RepeatingFieldSet2StartField, p.RepeatingFieldSet2Size},
	}
	for i, spec := range specs {
		if spec.start == nil || spec.size == nil || *spec.size <= 0 {
			continue
		}
		ranges[i] = [2]int{*spec.start, *spec.start + *spec.size}
	}
	return ranges
}

func namedPGNStructFields(sources []sourceField, used map[string]int) []pgnStructField {
	fields := make([]pgnStructField, 0, len(sources))
	for _, source := range sources {
		fieldType, _ := pgnFieldType(source)
		fields = append(fields, pgnStructField{
			Source: source,
			Name:   uniqueFieldName(used, exportedGoIdentifier(firstNonEmpty(source.Id, source.Name))),
			Type:   fieldType,
		})
	}
	return fields
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

func orderedSourceFields(fields []sourceField) []sourceField {
	ordered := append([]sourceField(nil), fields...)
	sort.SliceStable(ordered, func(i, j int) bool {
		return ordered[i].Order < ordered[j].Order
	})
	return ordered
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
	fmt.Fprintf(b, "{PGN:%d, StructName:%q, SourceID:%q, Description:%q, Type:%q, Complete:%t, Fallback:%t",
		p.PGN, typeName, p.Id, p.Description, p.Type, p.Complete, p.Fallback)
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
