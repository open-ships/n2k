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
	"regexp"
	"strconv"
	"strings"
)

const defaultCanboatJSONURL = "https://raw.githubusercontent.com/canboat/canboat/master/docs/canboat.json"

type canboatCatalog struct {
	PGNs []canboatPGN `json:"PGNs"`
}

type canboatPGN struct {
	PGN         uint32         `json:"PGN"`
	ID          string         `json:"Id"`
	Description string         `json:"Description"`
	Type        string         `json:"Type"`
	Fields      []canboatField `json:"Fields"`
}

type canboatField struct {
	Order                int         `json:"Order"`
	Name                 string      `json:"Name"`
	BitLength            int         `json:"BitLength"`
	BitOffset            int         `json:"BitOffset"`
	BitLengthVariable    bool        `json:"BitLengthVariable"`
	FieldType            string      `json:"FieldType"`
	Resolution           json.Number `json:"Resolution"`
	Signed               bool        `json:"Signed"`
	Unit                 string      `json:"Unit"`
	LookupEnumeration    string      `json:"LookupEnumeration"`
	LookupBitEnumeration string      `json:"LookupBitEnumeration"`
	Match                *int        `json:"Match"`
}

type typedEntry struct {
	ID          string
	PGN         uint32
	Description string
}

func main() {
	registryPath := flag.String("registry", "pgn/registry.go", "path to pgn registry.go")
	inputPath := flag.String("input", "", "local canboat.json file; defaults to --url")
	url := flag.String("url", defaultCanboatJSONURL, "canboat.json URL")
	check := flag.Bool("check", false, "fail if registry.go is not already synced")
	flag.Parse()

	if err := run(*registryPath, *inputPath, *url, *check); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(registryPath string, inputPath string, url string, check bool) error {
	registry, err := os.ReadFile(registryPath)
	if err != nil {
		return fmt.Errorf("read registry: %w", err)
	}

	catalog, err := loadCanboatCatalog(inputPath, url)
	if err != nil {
		return err
	}

	typedEntries, err := parseTypedEntries(string(registry))
	if err != nil {
		return err
	}
	catalogOnly := catalogOnlyEntries(catalog.PGNs, typedEntries)

	updated, err := replaceCatalogOnlyList(string(registry), renderCatalogOnlyList(catalogOnly))
	if err != nil {
		return err
	}
	formatted, err := format.Source([]byte(updated))
	if err != nil {
		return fmt.Errorf("format registry: %w", err)
	}

	if check {
		current, err := format.Source(registry)
		if err != nil {
			return fmt.Errorf("format current registry: %w", err)
		}
		if !bytes.Equal(current, formatted) {
			return fmt.Errorf("%s is not synced with canboat.json", registryPath)
		}
		return nil
	}

	if !bytes.Equal(registry, formatted) {
		if err := os.WriteFile(registryPath, formatted, 0o644); err != nil {
			return fmt.Errorf("write registry: %w", err)
		}
	}
	return nil
}

func loadCanboatCatalog(inputPath string, url string) (canboatCatalog, error) {
	var raw []byte
	var err error
	if inputPath != "" {
		raw, err = os.ReadFile(inputPath)
		if err != nil {
			return canboatCatalog{}, fmt.Errorf("read canboat json: %w", err)
		}
	} else {
		raw, err = readURL(url)
		if err != nil {
			return canboatCatalog{}, err
		}
	}

	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.UseNumber()
	var catalog canboatCatalog
	if err := decoder.Decode(&catalog); err != nil {
		return canboatCatalog{}, fmt.Errorf("parse canboat json: %w", err)
	}
	return catalog, nil
}

func readURL(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetch canboat json: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("fetch canboat json: %s", resp.Status)
	}
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read canboat json response: %w", err)
	}
	return raw, nil
}

func parseTypedEntries(registry string) ([]typedEntry, error) {
	start := strings.Index(registry, "var typedPgnList = []PgnInfo{")
	if start < 0 {
		return nil, fmt.Errorf("typedPgnList not found")
	}
	end := strings.Index(registry[start:], "var pgnList = ")
	if end < 0 {
		return nil, fmt.Errorf("pgnList not found after typedPgnList")
	}
	block := registry[start : start+end]

	entryRE := regexp.MustCompile(`(?s)\{\s*Id:\s*"((?:[^"\\]|\\.)*)",\s*PGN:\s*(\d+),\s*Description:\s*"((?:[^"\\]|\\.)*)"`)
	matches := entryRE.FindAllStringSubmatch(block, -1)
	entries := make([]typedEntry, 0, len(matches))
	for _, match := range matches {
		id, err := strconv.Unquote(`"` + match[1] + `"`)
		if err != nil {
			return nil, fmt.Errorf("unquote typed id: %w", err)
		}
		pgn, err := strconv.ParseUint(match[2], 10, 32)
		if err != nil {
			return nil, fmt.Errorf("parse typed pgn: %w", err)
		}
		description, err := strconv.Unquote(`"` + match[3] + `"`)
		if err != nil {
			return nil, fmt.Errorf("unquote typed description: %w", err)
		}
		entries = append(entries, typedEntry{
			ID:          id,
			PGN:         uint32(pgn),
			Description: description,
		})
	}
	return entries, nil
}

func catalogOnlyEntries(upstream []canboatPGN, typedEntries []typedEntry) []canboatPGN {
	typedIDs := make(map[string]struct{}, len(typedEntries))
	typedKeys := make(map[string]struct{}, len(typedEntries))
	for _, entry := range typedEntries {
		typedIDs[strings.ToLower(entry.ID)] = struct{}{}
		typedKeys[pgnDescriptionKey(entry.PGN, entry.Description)] = struct{}{}
	}

	catalogOnly := make([]canboatPGN, 0, len(upstream))
	for _, entry := range upstream {
		if _, exists := typedIDs[strings.ToLower(entry.ID)]; exists {
			continue
		}
		if _, exists := typedKeys[pgnDescriptionKey(entry.PGN, entry.Description)]; exists {
			continue
		}
		catalogOnly = append(catalogOnly, entry)
	}
	return catalogOnly
}

func pgnDescriptionKey(pgn uint32, description string) string {
	return fmt.Sprintf("%d\x00%s", pgn, description)
}

func replaceCatalogOnlyList(registry string, renderedCatalog string) (string, error) {
	start := strings.Index(registry, "// catalogOnlyList contains")
	if start < 0 {
		return "", fmt.Errorf("catalogOnlyList comment not found")
	}
	end := strings.Index(registry[start:], "var typedPgnList = []PgnInfo{")
	if end < 0 {
		return "", fmt.Errorf("typedPgnList not found after catalogOnlyList")
	}
	return registry[:start] + renderedCatalog + "\n" + registry[start+end:], nil
}

func renderCatalogOnlyList(entries []canboatPGN) string {
	var b strings.Builder
	b.WriteString("// catalogOnlyList contains current canboat catalog entries that do not have a typed decoder in this package yet.\n")
	b.WriteString("// They are still indexed into pgnList and PgnInfoLookup so known catalog PGNs are not treated as missing data.\n")
	b.WriteString("var catalogOnlyList = []PgnInfo{\n")
	for _, entry := range entries {
		renderCatalogEntry(&b, entry)
	}
	b.WriteString("}\n")
	return b.String()
}

func renderCatalogEntry(b *strings.Builder, entry canboatPGN) {
	entryID := entry.ID
	if entryID == "" {
		entryID = fmt.Sprintf("pgn%d", entry.PGN)
	}

	fmt.Fprintf(b, "\t{\n")
	fmt.Fprintf(b, "\t\tId:          %s,\n", goString(entryID))
	fmt.Fprintf(b, "\t\tPGN:         %d,\n", entry.PGN)
	fmt.Fprintf(b, "\t\tDescription: %s,\n", goString(entry.Description))
	fmt.Fprintf(b, "\t\tFast:        %t,\n", entry.Type == "Fast")
	fmt.Fprintf(b, "\t\tManId:       %d,\n", manufacturerID(entry))
	if len(entry.Fields) > 0 {
		fmt.Fprintf(b, "\t\tFields: map[int]*FieldDescriptor{\n")
		for _, field := range entry.Fields {
			renderField(b, field)
		}
		fmt.Fprintf(b, "\t\t},\n")
	}
	fmt.Fprintf(b, "\t},\n")
}

func renderField(b *strings.Builder, field canboatField) {
	fmt.Fprintf(b, "\t\t\t%d: {\n", field.Order)
	fmt.Fprintf(b, "\t\t\t\tName:              %s,\n", goString(field.Name))
	fmt.Fprintf(b, "\t\t\t\tBitLength:         %d,\n", field.BitLength)
	fmt.Fprintf(b, "\t\t\t\tBitOffset:         %d,\n", field.BitOffset)
	fmt.Fprintf(b, "\t\t\t\tBitLengthVariable: %t,\n", field.BitLengthVariable)
	fmt.Fprintf(b, "\t\t\t\tCanboatType:       %s,\n", goString(field.FieldType))
	fmt.Fprintf(b, "\t\t\t\tGolangType:        \"\",\n")
	fmt.Fprintf(b, "\t\t\t\tResolution:        %s,\n", goNumber(field.Resolution, "1"))
	fmt.Fprintf(b, "\t\t\t\tSigned:            %t,\n", field.Signed)
	if field.Unit != "" {
		fmt.Fprintf(b, "\t\t\t\tUnit:              %s,\n", goString(field.Unit))
	}
	if field.LookupBitEnumeration != "" {
		fmt.Fprintf(b, "\t\t\t\tBitLookupName:     %s,\n", goString(field.LookupBitEnumeration))
	}
	if field.Match != nil {
		fmt.Fprintf(b, "\t\t\t\tMatch:             ptrInt(%d),\n", *field.Match)
	}
	fmt.Fprintf(b, "\t\t\t},\n")
}

func manufacturerID(entry canboatPGN) int {
	for _, field := range entry.Fields {
		isManufacturer := field.Name == "Manufacturer Code" || field.LookupEnumeration == "MANUFACTURER_CODE"
		if isManufacturer && field.Match != nil {
			return *field.Match
		}
	}
	return 0
}

func goString(value string) string {
	return strconv.QuoteToASCII(value)
}

func goNumber(value json.Number, defaultValue string) string {
	if value.String() == "" {
		return defaultValue
	}
	return value.String()
}
