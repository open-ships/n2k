package pgn

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"testing"
)

type pgnManifest struct {
	SchemaVersion int                   `json:"schemaVersion"`
	Package       string                `json:"package"`
	Source        string                `json:"source"`
	Description   string                `json:"description"`
	Notes         []string              `json:"notes"`
	Counts        pgnManifestCounts     `json:"counts"`
	Categories    []pgnManifestCategory `json:"categories"`
}

type pgnManifestCounts struct {
	Categories          int `json:"categories"`
	SupportedVariants   int `json:"supportedVariants"`
	UniqueSupportedPgns int `json:"uniqueSupportedPgns"`
	TypedVariants       int `json:"typedVariants"`
	UniqueTypedPgns     int `json:"uniqueTypedPgns"`
}

type pgnManifestCategory struct {
	ID   string             `json:"id"`
	Name string             `json:"name"`
	PGNs []pgnManifestEntry `json:"pgns"`
}

type pgnManifestEntry struct {
	ID             string                  `json:"id"`
	SourceID       string                  `json:"sourceId,omitempty"`
	PGN            uint32                  `json:"pgn"`
	Description    string                  `json:"description"`
	Implementation string                  `json:"implementation"`
	Complete       bool                    `json:"complete"`
	Fallback       bool                    `json:"fallback,omitempty"`
	Missing        []string                `json:"missing,omitempty"`
	FastPacket     bool                    `json:"fastPacket"`
	ManufacturerID int                     `json:"manufacturerId"`
	Proprietary    bool                    `json:"proprietary"`
	SourceFile     string                  `json:"sourceFile"`
	PayloadShape   pgnManifestPayloadShape `json:"payloadShape"`
	Example        map[string]any          `json:"example"`
}

type pgnManifestPayloadShape struct {
	FieldCount int                `json:"fieldCount"`
	Fields     []pgnManifestField `json:"fields"`
}

type pgnManifestField struct {
	Index          int     `json:"index"`
	Name           string  `json:"name"`
	BitLength      uint16  `json:"bitLength"`
	BitOffset      uint16  `json:"bitOffset"`
	VariableLength bool    `json:"variableLength"`
	FieldType      string  `json:"fieldType"`
	GoType         string  `json:"goType"`
	Resolution     float32 `json:"resolution"`
	Signed         bool    `json:"signed"`
	Unit           string  `json:"unit,omitempty"`
	BitLookupName  string  `json:"bitLookupName,omitempty"`
	Match          *int    `json:"match,omitempty"`
}

type pgnManifestStructField struct {
	Type     string
	JSONName string
}

func TestPGNManifestMatchesRegistry(t *testing.T) {
	expectedManifest := buildExpectedManifest(t)
	if os.Getenv("UPDATE_PGN_MANIFEST") == "1" {
		raw, err := json.MarshalIndent(expectedManifest, "", "  ")
		if err != nil {
			t.Fatalf("marshal manifest: %v", err)
		}
		raw = append(raw, '\n')
		if err := os.WriteFile("manifest.json", raw, 0o644); err != nil {
			t.Fatalf("write manifest: %v", err)
		}
		return
	}

	raw, err := os.ReadFile("manifest.json")
	if err != nil {
		t.Fatalf("read manifest: %v", err)
	}

	var manifest pgnManifest
	if err := json.Unmarshal(raw, &manifest); err != nil {
		t.Fatalf("parse manifest: %v", err)
	}
	if manifest.SchemaVersion != 3 {
		t.Fatalf("schemaVersion = %d, want 3", manifest.SchemaVersion)
	}

	actual := make(map[string]pgnManifestEntry)
	actualCategories := make(map[string]string)
	for _, category := range manifest.Categories {
		for _, entry := range category.PGNs {
			if entry.SourceFile == "" {
				t.Errorf("%s has empty sourceFile", entry.ID)
			}
			if _, exists := actual[entry.ID]; exists {
				t.Fatalf("duplicate manifest id %q", entry.ID)
			}
			actual[entry.ID] = entry
			actualCategories[entry.ID] = category.ID
		}
	}

	if manifest.Counts.Categories != len(manifest.Categories) {
		t.Errorf("manifest category count = %d, want %d", manifest.Counts.Categories, len(manifest.Categories))
	}
	if !reflect.DeepEqual(manifest.Counts, expectedManifest.Counts) {
		t.Errorf("manifest counts mismatch\n got: %#v\nwant: %#v", manifest.Counts, expectedManifest.Counts)
	}

	expected := manifestEntryMap(expectedManifest)
	for id, want := range expected {
		got, exists := actual[id]
		if !exists {
			t.Errorf("manifest missing %s", id)
			continue
		}
		if gotCategory := actualCategories[id]; gotCategory != functionCategory(want.SourceFile) {
			t.Errorf("manifest entry %s is in category %q, want %q", id, gotCategory, functionCategory(want.SourceFile))
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("manifest entry %s mismatch\n got: %#v\nwant: %#v", id, got, want)
		}
		delete(actual, id)
	}
	if len(actual) > 0 {
		var extra []string
		for id := range actual {
			extra = append(extra, id)
		}
		sort.Strings(extra)
		t.Fatalf("manifest contains extra entries: %s", strings.Join(extra, ", "))
	}
}

func buildExpectedManifest(t *testing.T) pgnManifest {
	t.Helper()

	uniquePGNs := make(map[uint32]struct{})
	uniqueTypedPGNs := make(map[uint32]struct{})
	categories := make(map[string]*pgnManifestCategory)
	var categoryOrder []string
	var typedVariants int
	registryInfos := manifestRegistryInfos()

	for _, info := range registryInfos {
		uniquePGNs[info.PGN] = struct{}{}
		entry := manifestEntry(t, info)
		typedVariants++
		uniqueTypedPGNs[entry.PGN] = struct{}{}
		categoryID := functionCategory(entry.SourceFile)
		category, ok := categories[categoryID]
		if !ok {
			categories[categoryID] = &pgnManifestCategory{
				ID:   categoryID,
				Name: manifestCategoryName(categoryID),
			}
			category = categories[categoryID]
			categoryOrder = append(categoryOrder, categoryID)
		}
		category.PGNs = append(category.PGNs, entry)
	}

	manifestCategories := make([]pgnManifestCategory, 0, len(categoryOrder))
	for _, categoryID := range categoryOrder {
		manifestCategories = append(manifestCategories, *categories[categoryID])
	}

	return pgnManifest{
		SchemaVersion: 3,
		Package:       "github.com/open-ships/n2k/pgn",
		Source:        "Generated Pgn* structs and PgnInfoLookup metadata",
		Description:   "NMEA 2000 PGN variants decoded and encoded by this package.",
		Notes: []string{
			"Each entry corresponds to one Pgn* struct with PGNNumber, MessageInfo, SetMessageInfo, DecodePayload, and EncodePayload methods.",
			"Categories are grouped by the pgn package source file that owns each PGN struct.",
			"Duplicate PGN numbers are preserved when manufacturer-specific or group-function variants share the same PGN.",
			"Payload shapes are derived from each PgnInfo.Fields map and sorted by field index.",
			"Examples are deterministic sample JSON objects; they illustrate shape and are not captured bus values.",
			"Incomplete upstream definitions are included when a generated PGN struct exists for them.",
		},
		Counts: pgnManifestCounts{
			Categories:          len(manifestCategories),
			SupportedVariants:   len(registryInfos),
			UniqueSupportedPgns: len(uniquePGNs),
			TypedVariants:       typedVariants,
			UniqueTypedPgns:     len(uniqueTypedPGNs),
		},
		Categories: manifestCategories,
	}
}

func manifestRegistryInfos() []*PgnInfo {
	var infos []*PgnInfo
	for _, pgnInfos := range PgnInfoLookup {
		infos = append(infos, pgnInfos...)
	}
	sort.SliceStable(infos, func(i, j int) bool {
		if infos[i].PGN != infos[j].PGN {
			return infos[i].PGN < infos[j].PGN
		}
		if infos[i].Description != infos[j].Description {
			return infos[i].Description < infos[j].Description
		}
		return infos[i].Id < infos[j].Id
	})
	return infos
}

func manifestEntry(t *testing.T, info *PgnInfo) pgnManifestEntry {
	t.Helper()
	sourceFile := manifestStructSourceFiles(t)[info.Id]
	if sourceFile == "" {
		t.Fatalf("missing source file for %s", info.Id)
	}
	return pgnManifestEntry{
		ID:             info.Id,
		SourceID:       info.CanboatId,
		PGN:            info.PGN,
		Description:    info.Description,
		Implementation: "typed",
		Complete:       info.Complete,
		Fallback:       info.Fallback,
		Missing:        append([]string(nil), info.Missing...),
		FastPacket:     info.Fast,
		ManufacturerID: int(info.ManId),
		Proprietary:    IsProprietaryPGN(info.PGN),
		SourceFile:     sourceFile,
		PayloadShape:   manifestPayloadShape(info.Fields),
		Example:        manifestExample(t, info),
	}
}

func manifestEntryMap(manifest pgnManifest) map[string]pgnManifestEntry {
	entries := make(map[string]pgnManifestEntry)
	for _, category := range manifest.Categories {
		for _, entry := range category.PGNs {
			entries[entry.ID] = entry
		}
	}
	return entries
}

func functionCategory(sourceFile string) string {
	category := strings.TrimSuffix(sourceFile, filepath.Ext(sourceFile))
	return strings.TrimSuffix(category, "_pgn_generated")
}

func manifestCategoryName(categoryID string) string {
	switch categoryID {
	case "ais":
		return "AIS"
	case "bep_marine":
		return "BEP Marine"
	case "bg":
		return "B&G"
	case "diverse_yacht_services":
		return "Diverse Yacht Services"
	case "sea_recovery":
		return "Sea Recovery"
	case "seatalk":
		return "SeaTalk"
	case "simnet":
		return "SimNet"
	case "sonichub":
		return "SonicHub"
	default:
		parts := strings.Split(categoryID, "_")
		for i, part := range parts {
			if part == "" {
				continue
			}
			parts[i] = strings.ToUpper(part[:1]) + part[1:]
		}
		return strings.Join(parts, " ")
	}
}

func manifestPayloadShape(fields map[int]*FieldDescriptor) pgnManifestPayloadShape {
	indexes := make([]int, 0, len(fields))
	for index := range fields {
		indexes = append(indexes, index)
	}
	sort.Ints(indexes)

	manifestFields := make([]pgnManifestField, 0, len(indexes))
	for _, index := range indexes {
		field := fields[index]
		manifestField := pgnManifestField{
			Index:          index,
			Name:           field.Name,
			BitLength:      field.BitLength,
			BitOffset:      field.BitOffset,
			VariableLength: field.BitLengthVariable,
			FieldType:      field.CanboatType,
			GoType:         field.GolangType,
			Resolution:     field.Resolution,
			Signed:         field.Signed,
			Unit:           field.Unit,
			BitLookupName:  field.BitLookupName,
		}
		if field.Match != nil {
			match := *field.Match
			manifestField.Match = &match
		}
		manifestFields = append(manifestFields, manifestField)
	}

	return pgnManifestPayloadShape{
		FieldCount: len(manifestFields),
		Fields:     manifestFields,
	}
}

func manifestExample(t *testing.T, info *PgnInfo) map[string]any {
	t.Helper()
	return manifestTypedExample(t, info.Id, info.PGN)
}

func manifestTypedExample(t *testing.T, structName string, pgn uint32) map[string]any {
	t.Helper()
	fields := manifestStructFields(t, structName)
	example := make(map[string]any, len(fields))
	for _, field := range fields {
		example[field.JSONName] = manifestExampleValue(t, field.Type, field.JSONName, pgn)
	}
	return example
}

func manifestExampleValue(t *testing.T, goType string, jsonName string, pgn uint32) any {
	t.Helper()
	goType = strings.TrimSpace(goType)
	goType = strings.TrimPrefix(goType, "*")

	switch goType {
	case "MessageInfo":
		return map[string]any{
			"timestamp": "2026-01-01T00:00:00Z",
			"priority":  float64(6),
			"pgn":       float64(pgn),
			"sourceId":  float64(1),
			"targetId":  float64(255),
		}
	case "units.Distance":
		return manifestUnitExample(0)
	case "units.Flow":
		return manifestUnitExample(1)
	case "units.Pressure":
		return manifestUnitExample(2)
	case "units.Temperature":
		return manifestUnitExample(3)
	case "units.Velocity":
		return manifestUnitExample(4)
	case "units.Volume":
		return manifestUnitExample(5)
	case "string":
		return "example"
	case "bool":
		return true
	case "float32", "float64":
		return 1.23
	}

	if strings.HasPrefix(goType, "[]") {
		elementType := strings.TrimSpace(strings.TrimPrefix(goType, "[]"))
		if elementType == "uint8" || elementType == "byte" {
			return []any{float64(1), float64(2), float64(3)}
		}
		return []any{manifestTypedExample(t, elementType, pgn)}
	}

	if isManifestIntegerType(goType) {
		if jsonName == "pgn" {
			return float64(pgn)
		}
		return float64(1)
	}

	if _, ok := manifestStructFieldCache(t)[goType]; ok {
		return manifestTypedExample(t, goType, pgn)
	}

	// Enum aliases marshal as their numeric underlying value.
	return float64(1)
}

func manifestUnitExample(unitType float64) map[string]any {
	return map[string]any{
		"value":    1.23,
		"unit":     float64(0),
		"unitType": unitType,
	}
}

func isManifestIntegerType(goType string) bool {
	switch goType {
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "byte", "rune":
		return true
	default:
		return false
	}
}

func manifestStructFields(t *testing.T, structName string) []pgnManifestStructField {
	t.Helper()
	fields, ok := manifestStructFieldCache(t)[structName]
	if !ok {
		t.Fatalf("missing struct definition for %s", structName)
	}
	return fields
}

func manifestStructFieldCache(t *testing.T) map[string][]pgnManifestStructField {
	t.Helper()
	if manifestStructFieldsCache != nil {
		return manifestStructFieldsCache
	}

	files, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatalf("glob pgn source files: %v", err)
	}

	manifestStructFieldsCache = make(map[string][]pgnManifestStructField)
	for _, file := range files {
		if strings.HasSuffix(file, "_test.go") {
			continue
		}
		raw, err := os.ReadFile(file)
		if err != nil {
			t.Fatalf("read %s: %v", file, err)
		}
		for name, fields := range parseManifestStructFields(string(raw)) {
			manifestStructFieldsCache[name] = fields
		}
	}

	return manifestStructFieldsCache
}

func manifestStructSourceFiles(t *testing.T) map[string]string {
	t.Helper()
	if manifestStructSourceFilesCache != nil {
		return manifestStructSourceFilesCache
	}

	files, err := filepath.Glob("*_pgn_generated.go")
	if err != nil {
		t.Fatalf("glob PGN source files: %v", err)
	}

	manifestStructSourceFilesCache = make(map[string]string)
	for _, file := range files {
		raw, err := os.ReadFile(file)
		if err != nil {
			t.Fatalf("read %s: %v", file, err)
		}
		for name := range parseManifestStructFields(string(raw)) {
			manifestStructSourceFilesCache[name] = filepath.Base(file)
		}
	}

	return manifestStructSourceFilesCache
}

func parseManifestStructFields(source string) map[string][]pgnManifestStructField {
	structs := make(map[string][]pgnManifestStructField)
	structRE := regexp.MustCompile(`type\s+([A-Za-z_][A-Za-z0-9_]*)\s+struct\s+\{`)
	fieldRE := regexp.MustCompile("^\\s*[A-Za-z_][A-Za-z0-9_]*\\s+(.+?)\\s+`json:\"([^\"]+)\"`")

	matches := structRE.FindAllStringSubmatchIndex(source, -1)
	for _, match := range matches {
		name := source[match[2]:match[3]]
		bodyStart := match[1]
		bodyEnd := strings.Index(source[bodyStart:], "\n}")
		if bodyEnd < 0 {
			continue
		}
		body := source[bodyStart : bodyStart+bodyEnd]

		var fields []pgnManifestStructField
		for _, line := range strings.Split(body, "\n") {
			fieldMatch := fieldRE.FindStringSubmatch(line)
			if fieldMatch == nil {
				continue
			}
			jsonName := strings.Split(fieldMatch[2], ",")[0]
			if jsonName == "-" {
				continue
			}
			fields = append(fields, pgnManifestStructField{
				Type:     strings.TrimSpace(fieldMatch[1]),
				JSONName: jsonName,
			})
		}
		structs[name] = fields
	}

	return structs
}

var manifestStructFieldsCache map[string][]pgnManifestStructField
var manifestStructSourceFilesCache map[string]string
