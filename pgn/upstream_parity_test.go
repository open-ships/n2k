package pgn

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"math"
	"net/http"
	"os"
	"sort"
	"strings"
	"testing"
	"time"
)

var upstreamParityURL = upstreamSchemaURL()

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

func TestUpstreamParity(t *testing.T) {
	if os.Getenv("UPSTREAM_PARITY") == "" {
		t.Skip("set UPSTREAM_PARITY=1 to fetch current upstream schema and run parity checks")
	}

	raw, err := fetchUpstreamSchema()
	if err != nil {
		t.Fatalf("fetch upstream schema: %v", err)
	}

	var upstream sourceFile
	if err := json.Unmarshal(raw, &upstream); err != nil {
		t.Fatalf("parse upstream JSON: %v", err)
	}

	metadata := metadataEntries(PgnInfoLookup)
	unseen := metadataEntries(UnseenLookup)
	known := make(map[string]*PgnInfo, len(metadata)+len(unseen))
	for key, info := range metadata {
		known[key] = info
	}
	for key, info := range unseen {
		if _, exists := known[key]; !exists {
			known[key] = info
		}
	}

	upstreamKeys := make(map[string]sourcePGN, len(upstream.PGNs))
	missingMetadataComplete := make([]string, 0)
	missingKnown := make([]string, 0)
	fieldMismatches := make([]string, 0)
	upstreamFieldTypes := make(map[string]struct{})

	for _, def := range upstream.PGNs {
		key := upstreamKey(def.PGN, def.Description)
		upstreamKeys[key] = def
		for _, field := range def.Fields {
			upstreamFieldTypes[field.FieldType] = struct{}{}
		}

		if _, ok := known[key]; !ok {
			missingKnown = append(missingKnown, formatUpstreamPGN(def))
		}
		local, metadataOK := metadata[key]
		if def.Complete && !metadataOK {
			missingMetadataComplete = append(missingMetadataComplete, formatUpstreamPGN(def))
		}
		if metadataOK {
			fieldMismatches = append(fieldMismatches, compareUpstreamPGN(def, local)...)
		}
	}

	staleMetadata := make([]string, 0)
	for key, local := range metadata {
		if _, ok := upstreamKeys[key]; !ok {
			staleMetadata = append(staleMetadata, fmt.Sprintf("%d %s", local.PGN, local.Description))
		}
	}
	sort.Strings(staleMetadata)

	unsupportedFieldTypes := missingFieldTypes(upstreamFieldTypes, localFieldTypes(metadata))

	t.Logf("upstream source %s schema %s: %d PGN entries", upstream.Version, upstream.SchemaVersion, len(upstream.PGNs))
	if len(unsupportedFieldTypes) > 0 {
		t.Errorf("unsupported upstream source field types: %s", strings.Join(unsupportedFieldTypes, ", "))
	}
	if len(missingMetadataComplete) > 0 {
		t.Errorf("complete upstream PGNs missing generated structs (%d):\n%s",
			len(missingMetadataComplete), sampleLines(missingMetadataComplete, 25))
	}
	if len(missingKnown) > 0 {
		t.Errorf("upstream PGNs missing from generated metadata (%d):\n%s",
			len(missingKnown), sampleLines(missingKnown, 25))
	}
	if len(staleMetadata) > 0 {
		t.Errorf("generated PGNs not present in current upstream source by PGN+description (%d):\n%s",
			len(staleMetadata), sampleLines(staleMetadata, 25))
	}
	if len(fieldMismatches) > 0 {
		t.Errorf("generated PGNs with field/layout mismatches (%d):\n%s",
			len(fieldMismatches), sampleLines(fieldMismatches, 25))
	}
}

// TestEnumConstantsCoverLookupEnumerations checks that every LookupEnumeration
// name referenced by a PgnInfoLookup field descriptor has a corresponding
// <CamelCase>Const type declared in enums.go.
//
// It always runs offline: enums and metadata share the pinned schema input
// consumed by just pgn-sync, and missing lookup families are ordinary failures.
//
// Mapping rule (derived by checking enums.go's actual type names against the
// LookupEnumeration names in upstream_definitions.go): split the
// LookupEnumeration name on '_', title-case each segment (uppercase the
// first rune, lowercase the rest -- including segments that are pure digits,
// e.g. "16" stays "16"), concatenate the segments, and append "Const". For
// example "DIRECTION_REFERENCE" -> "DirectionReferenceConst" and
// "SEATALK_PILOT_MODE_16" -> "SeatalkPilotMode16Const". Verified against
// every LookupEnumeration name currently in the generated metadata: every
// name that has a matching enums.go type resolves correctly under this rule
// with zero exceptions, so no alias map is needed today; add one here if a
// future case needs it.
//
// The same naming rule applies to bit, indirect, and field-type lookups.
func TestEnumConstantsCoverLookupEnumerations(t *testing.T) {
	constTypes := enumConstTypeNames(t)

	names := make(map[string]struct{})
	for _, infos := range PgnInfoLookup {
		for _, info := range infos {
			for _, field := range info.Fields {
				for _, name := range []string{field.LookupEnumeration, field.LookupBitEnumeration, field.LookupIndirectEnumeration, field.LookupFieldTypeEnumeration} {
					if name != "" {
						names[name] = struct{}{}
					}
				}
			}
		}
	}

	missing := make([]string, 0, len(names))
	for name := range names {
		want := lookupEnumerationConstName(name)
		if _, ok := constTypes[want]; !ok {
			missing = append(missing, fmt.Sprintf("%s -> %s", name, want))
		}
	}
	if len(missing) > 0 {
		t.Errorf("LookupEnumeration names with no matching enums.go Const type (%d):\n%s",
			len(missing), sampleLines(missing, 50))
	}
}

// lookupEnumerationConstName converts a LookupEnumeration name to the Go
// Const type name enums.go would declare for it under the naming rule
// documented on TestEnumConstantsCoverLookupEnumerations.
func lookupEnumerationConstName(name string) string {
	var b strings.Builder
	for _, part := range strings.Split(name, "_") {
		if part == "" {
			continue
		}
		b.WriteString(strings.ToUpper(part[:1]))
		b.WriteString(strings.ToLower(part[1:]))
	}
	b.WriteString("Const")
	return b.String()
}

// enumConstTypeNames AST-parses enums.go and returns the set of declared
// type names, following the same parse-package-source pattern
// metadata_consistency_test.go uses for the PGN struct files.
func enumConstTypeNames(t *testing.T) map[string]struct{} {
	t.Helper()
	fset := token.NewFileSet()
	parsed, err := parser.ParseFile(fset, "enums.go", nil, 0)
	if err != nil {
		t.Fatalf("parse enums.go: %v", err)
	}
	types := make(map[string]struct{})
	for _, decl := range parsed.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			types[typeSpec.Name.Name] = struct{}{}
		}
	}
	return types
}

func fetchUpstreamSchema() ([]byte, error) {
	client := http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(upstreamParityURL)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("GET %s: %s", upstreamParityURL, resp.Status)
	}
	return io.ReadAll(resp.Body)
}

func metadataEntries(lookup map[uint32][]*PgnInfo) map[string]*PgnInfo {
	entries := make(map[string]*PgnInfo)
	for _, infos := range lookup {
		for _, info := range infos {
			entries[upstreamKey(info.PGN, info.Description)] = info
		}
	}
	return entries
}

func upstreamKey(pgn uint32, description string) string {
	return fmt.Sprintf("%d\x00%s", pgn, description)
}

func formatUpstreamPGN(def sourcePGN) string {
	return fmt.Sprintf("%d %s", def.PGN, def.Description)
}

func localFieldTypes(metadata map[string]*PgnInfo) map[string]struct{} {
	types := make(map[string]struct{})
	for _, info := range metadata {
		for _, field := range info.Fields {
			types[field.SourceType] = struct{}{}
		}
	}
	return types
}

func missingFieldTypes(upstream map[string]struct{}, local map[string]struct{}) []string {
	missing := make([]string, 0)
	for typ := range upstream {
		if _, ok := local[typ]; ok {
			continue
		}
		missing = append(missing, typ)
	}
	sort.Strings(missing)
	return missing
}

func compareUpstreamPGN(def sourcePGN, local *PgnInfo) []string {
	mismatches := make([]string, 0)
	prefix := fmt.Sprintf("%d %s", def.PGN, def.Description)
	if local.SourceID != def.Id {
		mismatches = append(mismatches, fmt.Sprintf("%s: id local=%s upstream=%s", prefix, local.SourceID, def.Id))
	}
	if local.Explanation != def.Explanation {
		mismatches = append(mismatches, fmt.Sprintf("%s: explanation mismatch", prefix))
	}
	if local.Type != def.Type {
		mismatches = append(mismatches, fmt.Sprintf("%s: type local=%s upstream=%s", prefix, local.Type, def.Type))
	}
	if local.Fast != (def.Type == "Fast") {
		mismatches = append(mismatches, fmt.Sprintf("%s: fast local=%v upstream=%s", prefix, local.Fast, def.Type))
	}
	if local.Complete != def.Complete {
		mismatches = append(mismatches, fmt.Sprintf("%s: complete local=%v upstream=%v", prefix, local.Complete, def.Complete))
	}
	if local.Fallback != def.Fallback {
		mismatches = append(mismatches, fmt.Sprintf("%s: fallback local=%v upstream=%v", prefix, local.Fallback, def.Fallback))
	}
	if !equalStringSlices(local.Missing, def.Missing) {
		mismatches = append(mismatches, fmt.Sprintf("%s: missing local=%v upstream=%v", prefix, local.Missing, def.Missing))
	}
	mismatches = append(mismatches, compareIntPtr(prefix, "length", local.Length, def.Length)...)
	mismatches = append(mismatches, compareIntPtr(prefix, "minLength", local.MinLength, def.MinLength)...)
	mismatches = append(mismatches, compareUint8Ptr(prefix, "priority", local.Priority, def.Priority)...)
	mismatches = append(mismatches, compareIntPtr(prefix, "transmissionInterval", local.TransmissionInterval, def.TransmissionInterval)...)
	mismatches = append(mismatches, compareBoolPtr(prefix, "transmissionIrregular", local.TransmissionIrregular, def.TransmissionIrregular)...)
	mismatches = append(mismatches, compareIntPtr(prefix, "repeatingFieldSet1StartField", local.RepeatingFieldSet1StartField, def.RepeatingFieldSet1StartField)...)
	mismatches = append(mismatches, compareIntPtr(prefix, "repeatingFieldSet1CountField", local.RepeatingFieldSet1CountField, def.RepeatingFieldSet1CountField)...)
	mismatches = append(mismatches, compareIntPtr(prefix, "repeatingFieldSet1Size", local.RepeatingFieldSet1Size, def.RepeatingFieldSet1Size)...)
	mismatches = append(mismatches, compareIntPtr(prefix, "repeatingFieldSet2StartField", local.RepeatingFieldSet2StartField, def.RepeatingFieldSet2StartField)...)
	mismatches = append(mismatches, compareIntPtr(prefix, "repeatingFieldSet2CountField", local.RepeatingFieldSet2CountField, def.RepeatingFieldSet2CountField)...)
	mismatches = append(mismatches, compareIntPtr(prefix, "repeatingFieldSet2Size", local.RepeatingFieldSet2Size, def.RepeatingFieldSet2Size)...)

	if local.Id == "" {
		mismatches = append(mismatches, fmt.Sprintf("%s: struct id is empty", prefix))
	}
	if len(local.Fields) != def.FieldCount {
		mismatches = append(mismatches, fmt.Sprintf("%s: field count local=%d upstream=%d", prefix, len(local.Fields), def.FieldCount))
	}

	for _, upstreamField := range def.Fields {
		localField := local.Fields[upstreamField.Order]
		if localField == nil {
			mismatches = append(mismatches, fmt.Sprintf("%d %s field %d: missing local field", def.PGN, def.Description, upstreamField.Order))
			continue
		}
		fieldPrefix := fmt.Sprintf("%d %s field %d %s", def.PGN, def.Description, upstreamField.Order, upstreamField.Name)
		if localField.SourceID != upstreamField.Id {
			mismatches = append(mismatches, fmt.Sprintf("%s: id local=%s upstream=%s", fieldPrefix, localField.SourceID, upstreamField.Id))
		}
		if localField.Name != upstreamField.Name {
			mismatches = append(mismatches, fmt.Sprintf("%s: name local=%s upstream=%s", fieldPrefix, localField.Name, upstreamField.Name))
		}
		if localField.Description != string(upstreamField.Description) {
			mismatches = append(mismatches, fmt.Sprintf("%s: description mismatch", fieldPrefix))
		}
		if upstreamField.BitLength != nil && localField.BitLength != *upstreamField.BitLength {
			mismatches = append(mismatches, fmt.Sprintf("%s: bitLength local=%d upstream=%d", fieldPrefix, localField.BitLength, *upstreamField.BitLength))
		}
		expectedVariable := upstreamField.BitLengthVariable || upstreamField.BitLength == nil
		if localField.BitLengthVariable != expectedVariable {
			mismatches = append(mismatches, fmt.Sprintf("%s: bitLengthVariable local=%v upstream=%v", fieldPrefix, localField.BitLengthVariable, expectedVariable))
		}
		mismatches = append(mismatches, compareIntPtr(fieldPrefix, "bitLengthField", localField.BitLengthField, upstreamField.BitLengthField)...)
		if upstreamField.BitOffset != nil && localField.BitOffset != *upstreamField.BitOffset {
			mismatches = append(mismatches, fmt.Sprintf("%s: bitOffset local=%d upstream=%d", fieldPrefix, localField.BitOffset, *upstreamField.BitOffset))
		}
		if upstreamField.BitStart != nil && localField.BitStart != *upstreamField.BitStart {
			mismatches = append(mismatches, fmt.Sprintf("%s: bitStart local=%d upstream=%d", fieldPrefix, localField.BitStart, *upstreamField.BitStart))
		}
		if localField.SourceType != upstreamField.FieldType {
			mismatches = append(mismatches, fmt.Sprintf("%s: type local=%s upstream=%s", fieldPrefix, localField.SourceType, upstreamField.FieldType))
		}
		if upstreamField.Resolution != nil && math.Abs(float64(localField.Resolution)-*upstreamField.Resolution) > 1e-6 {
			mismatches = append(mismatches, fmt.Sprintf("%s: resolution local=%g upstream=%g", fieldPrefix, localField.Resolution, *upstreamField.Resolution))
		}
		if localField.Signed != upstreamField.Signed {
			mismatches = append(mismatches, fmt.Sprintf("%s: signed local=%v upstream=%v", fieldPrefix, localField.Signed, upstreamField.Signed))
		}
		if localField.Unit != upstreamField.Unit {
			mismatches = append(mismatches, fmt.Sprintf("%s: unit local=%s upstream=%s", fieldPrefix, localField.Unit, upstreamField.Unit))
		}
		if localField.PhysicalQuantity != upstreamField.PhysicalQuantity {
			mismatches = append(mismatches, fmt.Sprintf("%s: physicalQuantity local=%s upstream=%s", fieldPrefix, localField.PhysicalQuantity, upstreamField.PhysicalQuantity))
		}
		if localField.LookupEnumeration != upstreamField.LookupEnumeration {
			mismatches = append(mismatches, fmt.Sprintf("%s: lookupEnumeration local=%s upstream=%s", fieldPrefix, localField.LookupEnumeration, upstreamField.LookupEnumeration))
		}
		if localField.LookupBitEnumeration != upstreamField.LookupBitEnumeration {
			mismatches = append(mismatches, fmt.Sprintf("%s: lookupBitEnumeration local=%s upstream=%s", fieldPrefix, localField.LookupBitEnumeration, upstreamField.LookupBitEnumeration))
		}
		if localField.LookupIndirectEnumeration != upstreamField.LookupIndirectEnumeration {
			mismatches = append(mismatches, fmt.Sprintf("%s: lookupIndirectEnumeration local=%s upstream=%s", fieldPrefix, localField.LookupIndirectEnumeration, upstreamField.LookupIndirectEnumeration))
		}
		mismatches = append(mismatches, compareIntPtr(fieldPrefix, "lookupIndirectEnumerationFieldOrder", localField.LookupIndirectEnumerationFieldOrder, upstreamField.LookupIndirectEnumerationFieldOrder)...)
		if localField.LookupFieldTypeEnumeration != upstreamField.LookupFieldTypeEnumeration {
			mismatches = append(mismatches, fmt.Sprintf("%s: lookupFieldTypeEnumeration local=%s upstream=%s", fieldPrefix, localField.LookupFieldTypeEnumeration, upstreamField.LookupFieldTypeEnumeration))
		}
		if localField.Condition != upstreamField.Condition {
			mismatches = append(mismatches, fmt.Sprintf("%s: condition local=%s upstream=%s", fieldPrefix, localField.Condition, upstreamField.Condition))
		}
		mismatches = append(mismatches, compareIntPtr(fieldPrefix, "match", localField.Match, upstreamField.Match)...)
		mismatches = append(mismatches, compareFloatPtr(fieldPrefix, "rangeMin", localField.RangeMin, upstreamField.RangeMin)...)
		mismatches = append(mismatches, compareFloatPtr(fieldPrefix, "rangeMax", localField.RangeMax, upstreamField.RangeMax)...)
		mismatches = append(mismatches, compareFloatPtr(fieldPrefix, "offset", localField.Offset, upstreamField.Offset)...)
		mismatches = append(mismatches, compareInt64Ptr(fieldPrefix, "outOfRangeValue", localField.OutOfRangeValue, upstreamField.OutOfRangeValue)...)
		mismatches = append(mismatches, compareInt64Ptr(fieldPrefix, "reservedValue", localField.ReservedValue, upstreamField.ReservedValue)...)
		mismatches = append(mismatches, compareInt64Ptr(fieldPrefix, "unknownValue", localField.UnknownValue, upstreamField.UnknownValue)...)
		mismatches = append(mismatches, compareBoolPtr(fieldPrefix, "partOfPrimaryKey", localField.PartOfPrimaryKey, upstreamField.PartOfPrimaryKey)...)
	}
	return mismatches
}

func equalStringSlices(a []string, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func compareIntPtr(prefix string, name string, local *int, upstream *int) []string {
	if local == nil || upstream == nil {
		if local == upstream {
			return nil
		}
		return []string{fmt.Sprintf("%s: %s local=%v upstream=%v", prefix, name, formatPtr(local), formatPtr(upstream))}
	}
	if *local != *upstream {
		return []string{fmt.Sprintf("%s: %s local=%d upstream=%d", prefix, name, *local, *upstream)}
	}
	return nil
}

func compareInt64Ptr(prefix string, name string, local *int64, upstream *int64) []string {
	if local == nil || upstream == nil {
		if local == upstream {
			return nil
		}
		return []string{fmt.Sprintf("%s: %s local=%v upstream=%v", prefix, name, formatPtr(local), formatPtr(upstream))}
	}
	if *local != *upstream {
		return []string{fmt.Sprintf("%s: %s local=%d upstream=%d", prefix, name, *local, *upstream)}
	}
	return nil
}

func compareUint8Ptr(prefix string, name string, local *uint8, upstream *uint8) []string {
	if local == nil || upstream == nil {
		if local == upstream {
			return nil
		}
		return []string{fmt.Sprintf("%s: %s local=%v upstream=%v", prefix, name, formatPtr(local), formatPtr(upstream))}
	}
	if *local != *upstream {
		return []string{fmt.Sprintf("%s: %s local=%d upstream=%d", prefix, name, *local, *upstream)}
	}
	return nil
}

func compareBoolPtr(prefix string, name string, local *bool, upstream *bool) []string {
	if local == nil || upstream == nil {
		if local == upstream {
			return nil
		}
		return []string{fmt.Sprintf("%s: %s local=%v upstream=%v", prefix, name, formatPtr(local), formatPtr(upstream))}
	}
	if *local != *upstream {
		return []string{fmt.Sprintf("%s: %s local=%v upstream=%v", prefix, name, *local, *upstream)}
	}
	return nil
}

func compareFloatPtr(prefix string, name string, local *float64, upstream *float64) []string {
	if local == nil || upstream == nil {
		if local == upstream {
			return nil
		}
		return []string{fmt.Sprintf("%s: %s local=%v upstream=%v", prefix, name, formatPtr(local), formatPtr(upstream))}
	}
	if math.Abs(*local-*upstream) > 1e-9 {
		return []string{fmt.Sprintf("%s: %s local=%g upstream=%g", prefix, name, *local, *upstream)}
	}
	return nil
}

func formatPtr[T any](value *T) any {
	if value == nil {
		return nil
	}
	return *value
}

func sampleLines(lines []string, max int) string {
	sort.Strings(lines)
	if len(lines) <= max {
		return strings.Join(lines, "\n")
	}
	return strings.Join(lines[:max], "\n") + fmt.Sprintf("\n... %d more", len(lines)-max)
}
