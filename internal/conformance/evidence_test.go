package conformance

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func fixtureIndex(evidence ...Evidence) Index {
	return Index{SchemaVersion: 2, Requirements: []Requirement{{ID: "FIX-01", Category: "fixture", Behavior: "executable evidence", Evidence: evidence}}}
}

func TestResolveRequiresPackageQualifiedExecutableMatches(t *testing.T) {
	index := fixtureIndex(Evidence{Package: Module, Pattern: "^TestRequired$", Kind: "software"})
	if _, err := Resolve(index, Catalog{Module + "/other": {"TestRequired"}}); err == nil {
		t.Fatal("a name in a different package must not satisfy the requirement")
	}
	if _, err := Resolve(index, Catalog{Module: {"TestRequired"}}); err != nil {
		t.Fatal(err)
	}
	index.Requirements[0].Evidence[0].Pattern = "TestRequired*"
	if err := Validate(index); err == nil {
		t.Fatal("legacy unanchored/glob references must be rejected")
	}
	index.Requirements[0].Evidence[0].Pattern = "^Test.*$"
	if _, err := Resolve(index, Catalog{Module: {"TestRequired", "TestActisenseHardwareMatrix"}}); err == nil {
		t.Fatal("a broad software reference must not select hardware")
	}
}

func TestExecutionSummaryDistinguishesFailureMissingAndSkippedBranches(t *testing.T) {
	for _, expected := range []string{"pass", "fail", "skip", "missing", "skipped-child"} {
		t.Run(expected, func(t *testing.T) {
			index := fixtureIndex(Evidence{Package: Module, Pattern: "^TestRequired$", Kind: "software"})
			results, err := Resolve(index, Catalog{Module: {"TestRequired"}})
			if err != nil {
				t.Fatal(err)
			}
			if expected == "skipped-child" {
				Record(results, Event{Package: Module, Test: "TestRequired/device", Action: "skip"})
				Record(results, Event{Package: Module, Test: "TestRequired", Action: "pass"})
			} else if expected != "missing" {
				Record(results, Event{Package: Module, Test: "TestRequired", Action: expected})
			}
			err = Summarize(results)
			if (err == nil) != (expected == "pass") {
				t.Fatalf("summary error = %v for %s", err, expected)
			}
			if expected == "skipped-child" {
				expected = "skip"
			}
			if results[0].Status != expected {
				t.Fatalf("status = %s, want %s", results[0].Status, expected)
			}
		})
	}
}

func TestOptionalEvidenceCannotBeSelectedBySoftwareOrHideFailure(t *testing.T) {
	for _, kind := range []string{"hardware", "soak"} {
		t.Run(kind, func(t *testing.T) {
			index := fixtureIndex(
				Evidence{Package: Module, Pattern: "^Test.*$", Kind: "software"},
				Evidence{Package: Module, Pattern: "^TestExternalCapability$", Kind: kind},
			)
			catalog := Catalog{Module: {"TestRequired", "TestExternalCapability"}}
			if _, err := Resolve(index, catalog); err == nil {
				t.Fatal("software selection must reject any test classified as opt-in")
			}
			index.Requirements[0].Evidence[0].Pattern = "^TestRequired$"
			results, err := Resolve(index, catalog)
			if err != nil {
				t.Fatal(err)
			}
			Record(results, Event{Package: Module, Test: "TestRequired", Action: "fail"})
			if err := Summarize(results); err == nil || results[0].Status != "fail" {
				t.Fatalf("optional evidence must not conceal a software failure: results=%+v error=%v", results, err)
			}
		})
	}
}

func TestRunnerUsesGoJSONAndNeverExecutesConfiguredHardware(t *testing.T) {
	directory := t.TempDir()
	const pkg = Module + "/evidencefixture"
	if err := os.WriteFile(filepath.Join(directory, "go.mod"), []byte("module "+pkg+"\n\ngo 1.26.6\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	fixture := `package evidencefixture
import "testing"
func TestPass(t *testing.T) {}
func TestSkipped(t *testing.T) { t.Skip("fixture missing capability") }
func TestActisenseHardwareMatrix(t *testing.T) { t.Fatal("hardware must never execute") }
`
	if err := os.WriteFile(filepath.Join(directory, "fixture_test.go"), []byte(fixture), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("GOWORK", "off")
	t.Setenv("N2K_ACTISENSE_HARDWARE_CONFIG", "must-not-be-used")
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	index := fixtureIndex(Evidence{Package: pkg, Pattern: "^TestPass$", Kind: "software"})
	index.Requirements = append(index.Requirements, Requirement{ID: "HW-01", Category: "hardware", Behavior: "not executed in software CI", Evidence: []Evidence{{Package: pkg, Pattern: "^TestActisenseHardwareMatrix$", Kind: "hardware"}}})
	var artifact bytes.Buffer
	report, err := Run(ctx, directory, index, &artifact)
	if err != nil {
		t.Fatal(err)
	}
	if report.Status != "pass" || report.Requirements[1].Status != "not-run" {
		t.Fatalf("unexpected evidence status: %+v", report)
	}
	if !strings.Contains(artifact.String(), `"Test":"TestPass"`) || strings.Contains(artifact.String(), `"Test":"TestActisenseHardwareMatrix"`) {
		t.Fatalf("unexpected executable event stream: %s", artifact.String())
	}
	index.Requirements[0].Evidence[0].Pattern = "^TestSkipped$"
	report, err = Run(ctx, directory, index, nil)
	if err == nil || report.Requirements[0].Status != "skip" {
		t.Fatalf("skipped evidence must fail the gate: report=%+v error=%v", report, err)
	}
}
