package main

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestPinnedSchemaLoadsOfflineWithOriginalAttribution(t *testing.T) {
	t.Chdir("../..")
	t.Setenv("PGN_REFRESH_SOURCE", "")
	t.Setenv("HTTP_PROXY", "http://127.0.0.1:1")
	t.Setenv("HTTPS_PROXY", "http://127.0.0.1:1")
	raw, err := loadSource()
	if err != nil {
		t.Fatal(err)
	}
	var notice struct {
		Copyright string
		License   string
	}
	if err := json.Unmarshal(raw, &notice); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(notice.Copyright, "Kees Verruijt") || !strings.Contains(notice.Copyright, "Licensed under the Apache License, Version 2.0") || notice.License != "Apache License Version 2.0" {
		t.Fatalf("source attribution or license lost: %+v", notice)
	}
}

func TestEnumGenerationScopesNamesAndPreservesLookupMeaning(t *testing.T) {
	source := sourceFile{
		LookupEnumerations: []sourceEnumeration{
			{Name: "FIRST", MaxValue: 255, EnumValues: []sourceEnumValue{{Name: "Unknown", Value: 255}}},
			{Name: "SECOND", MaxValue: 255, EnumValues: []sourceEnumValue{{Name: "Unknown", Value: 255}}},
		},
		LookupBitEnumerations: []sourceEnumeration{{Name: "FLAGS", MaxValue: 15,
			EnumBitValues: []sourceEnumValue{{Name: "Busy", Bit: 0}, {Name: "Failed", Bit: 15}}}},
		LookupIndirectEnumerations: []sourceEnumeration{{Name: "FUNCTION", MaxValue: 255,
			EnumValues: []sourceEnumValue{{Name: "Engine", Value1: 35, Value2: 140}, {Name: "Engine", Value1: 50, Value2: 140}}}},
	}
	generated, err := generateEnums(source)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"FirstUnknown FirstConst = 255", "SecondUnknown SecondConst = 255", "FlagsBusy", "FlagsFailed", "32768", "FunctionEngineClass35", "FunctionEngineClass50", "9100", "12940"} {
		if !strings.Contains(string(generated), want) {
			t.Errorf("missing independent enum expectation %q", want)
		}
	}
}
