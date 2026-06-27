package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/format"
	"io"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

const canboatURL = "https://raw.githubusercontent.com/canboat/canboat/refs/heads/master/docs/canboat.json"

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
	raw, err := fetch(canboatURL)
	if err != nil {
		fatal(err)
	}
	var source canboatFile
	if err := json.Unmarshal(raw, &source); err != nil {
		fatal(err)
	}
	out, err := generate(source)
	if err != nil {
		fatal(err)
	}
	if err := os.WriteFile("pgn/canboat_generated.go", out, 0o644); err != nil {
		fatal(err)
	}
}

func fetch(url string) ([]byte, error) {
	client := http.Client{Timeout: 60 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("GET %s: %s", url, resp.Status)
	}
	return io.ReadAll(resp.Body)
}

func generate(source canboatFile) ([]byte, error) {
	sort.SliceStable(source.PGNs, func(i, j int) bool {
		if source.PGNs[i].PGN != source.PGNs[j].PGN {
			return source.PGNs[i].PGN < source.PGNs[j].PGN
		}
		return source.PGNs[i].Description < source.PGNs[j].Description
	})

	var b bytes.Buffer
	fmt.Fprintf(&b, "// Code generated by go run ./cmd/canboatgen; DO NOT EDIT.\n")
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
