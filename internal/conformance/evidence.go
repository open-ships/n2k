// Package conformance links public requirement claims to discoverable Go tests
// and records their executable outcomes. It never runs physical hardware tests.
package conformance

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"slices"
	"strings"
	"time"
)

const Module = "github.com/open-ships/n2k"

type Index struct {
	SchemaVersion int           `json:"schemaVersion"`
	Requirements  []Requirement `json:"requirements"`
}

type Requirement struct {
	ID       string     `json:"id"`
	Category string     `json:"category"`
	Behavior string     `json:"behavior"`
	Evidence []Evidence `json:"evidence"`
}

type Evidence struct {
	Package string `json:"package"`
	Pattern string `json:"pattern"`
	Kind    string `json:"kind"`
}

type Catalog map[string][]string

type Event struct {
	Action  string  `json:"Action"`
	Package string  `json:"Package"`
	Test    string  `json:"Test"`
	Output  string  `json:"Output"`
	Elapsed float64 `json:"Elapsed"`
}

type Execution struct {
	Test           string   `json:"test"`
	Status         string   `json:"status"`
	ElapsedSeconds float64  `json:"elapsedSeconds"`
	Skipped        []string `json:"skippedSubtests,omitempty"`
}

type EvidenceResult struct {
	Evidence
	Status string      `json:"status"`
	Tests  []Execution `json:"tests,omitempty"`
}

type RequirementResult struct {
	ID       string           `json:"id"`
	Status   string           `json:"status"`
	Evidence []EvidenceResult `json:"evidence"`
}

type Report struct {
	SchemaVersion      int                 `json:"schemaVersion"`
	CreatedAt          time.Time           `json:"createdAt"`
	Commit             string              `json:"commit"`
	WorkingTreeDirty   bool                `json:"workingTreeDirty"`
	GoVersion          string              `json:"goVersion"`
	RequirementsSHA256 string              `json:"requirementsSHA256"`
	Status             string              `json:"status"`
	Requirements       []RequirementResult `json:"requirements"`
	Errors             []string            `json:"errors,omitempty"`
}

func Validate(index Index) error {
	if index.SchemaVersion != 2 || len(index.Requirements) == 0 {
		return errors.New("conformance index requires schemaVersion 2 and nonempty requirements")
	}
	seen := make(map[string]bool)
	for _, requirement := range index.Requirements {
		if requirement.ID == "" || requirement.Category == "" || requirement.Behavior == "" || len(requirement.Evidence) == 0 || seen[requirement.ID] {
			return fmt.Errorf("invalid or duplicate requirement %q", requirement.ID)
		}
		seen[requirement.ID] = true
		for _, evidence := range requirement.Evidence {
			if evidence.Package != Module && !strings.HasPrefix(evidence.Package, Module+"/") {
				return fmt.Errorf("%s: evidence package must be fully qualified within %s", requirement.ID, Module)
			}
			if strings.Contains(evidence.Package, "..") || strings.ContainsAny(evidence.Package, " \t\n") {
				return fmt.Errorf("%s: invalid package %q", requirement.ID, evidence.Package)
			}
			if evidence.Kind != "software" && evidence.Kind != "hardware" && evidence.Kind != "soak" {
				return fmt.Errorf("%s: unsupported evidence kind %q", requirement.ID, evidence.Kind)
			}
			if !strings.HasPrefix(evidence.Pattern, "^Test") || !strings.HasSuffix(evidence.Pattern, "$") || strings.Contains(evidence.Pattern, "/") {
				return fmt.Errorf("%s: evidence must use an anchored top-level test pattern", requirement.ID)
			}
			if _, err := regexp.Compile(evidence.Pattern); err != nil {
				return fmt.Errorf("%s: invalid pattern: %w", requirement.ID, err)
			}
		}
	}
	return nil
}

func packages(index Index) []string {
	var result []string
	for _, requirement := range index.Requirements {
		for _, evidence := range requirement.Evidence {
			result = append(result, evidence.Package)
		}
	}
	slices.Sort(result)
	return slices.Compact(result)
}

// Discover uses compiled go test -list output, not source-text matching. Listing
// does not execute any test, including the optional hardware and soak entries.
func Discover(ctx context.Context, directory string, index Index) (Catalog, error) {
	if err := Validate(index); err != nil {
		return nil, err
	}
	catalog := make(Catalog)
	args := append([]string{"test", "-json", "-list", "^Test"}, packages(index)...)
	err := runGo(ctx, directory, args, nil, func(event Event) {
		if event.Action == "output" {
			name := strings.TrimSpace(event.Output)
			if strings.HasPrefix(name, "Test") && !strings.ContainsAny(name, " \t\n/") {
				catalog[event.Package] = append(catalog[event.Package], name)
			}
		}
	})
	if err != nil {
		return nil, fmt.Errorf("discover executable evidence: %w", err)
	}
	return catalog, nil
}

// Resolve fails every unmatched reference, including optional evidence, and
// prevents broad software patterns from selecting hardware or soak tests.
func Resolve(index Index, catalog Catalog) ([]RequirementResult, error) {
	if err := Validate(index); err != nil {
		return nil, err
	}
	optIn := make(map[string]map[string]bool)
	for _, requirement := range index.Requirements {
		for _, evidence := range requirement.Evidence {
			if evidence.Kind == "software" {
				continue
			}
			compiled, err := regexp.Compile(evidence.Pattern)
			if err != nil {
				return nil, err
			}
			if optIn[evidence.Package] == nil {
				optIn[evidence.Package] = make(map[string]bool)
			}
			for _, name := range catalog[evidence.Package] {
				if compiled.MatchString(name) {
					optIn[evidence.Package][name] = true
				}
			}
		}
	}
	var results []RequirementResult
	var failures []error
	for _, requirement := range index.Requirements {
		result := RequirementResult{ID: requirement.ID, Status: "missing"}
		for _, evidence := range requirement.Evidence {
			compiled, err := regexp.Compile(evidence.Pattern)
			if err != nil {
				return nil, err
			}
			item := EvidenceResult{Evidence: evidence, Status: "missing"}
			for _, name := range catalog[evidence.Package] {
				if !compiled.MatchString(name) {
					continue
				}
				if evidence.Kind == "software" && (optIn[evidence.Package][name] || name == "TestActisenseHardwareMatrix" || name == "TestReliabilitySoak") {
					failures = append(failures, fmt.Errorf("%s: software pattern selects opt-in test %s", requirement.ID, name))
					continue
				}
				status := "missing"
				if evidence.Kind != "software" {
					status, item.Status = "not-run", "not-run"
				}
				item.Tests = append(item.Tests, Execution{Test: name, Status: status})
			}
			if len(item.Tests) == 0 {
				failures = append(failures, fmt.Errorf("%s: no executable test matches %s %s", requirement.ID, evidence.Package, evidence.Pattern))
			}
			result.Evidence = append(result.Evidence, item)
		}
		results = append(results, result)
	}
	return results, errors.Join(failures...)
}

// Record retains top-level outcomes and any skipped descendants. A passed
// parent cannot disguise skipped branches of its claimed evidence.
func Record(results []RequirementResult, event Event) {
	if event.Test == "" || (event.Action != "pass" && event.Action != "fail" && event.Action != "skip") {
		return
	}
	for i := range results {
		for j := range results[i].Evidence {
			evidence := &results[i].Evidence[j]
			if evidence.Kind != "software" || evidence.Package != event.Package {
				continue
			}
			for k := range evidence.Tests {
				execution := &evidence.Tests[k]
				if event.Test == execution.Test {
					execution.Status, execution.ElapsedSeconds = event.Action, event.Elapsed
				} else if event.Action == "skip" && strings.HasPrefix(event.Test, execution.Test+"/") {
					execution.Skipped = append(execution.Skipped, event.Test)
				}
			}
		}
	}
}

func Summarize(results []RequirementResult) error {
	var failures []error
	for i := range results {
		results[i].Status = "pass"
		for j := range results[i].Evidence {
			evidence := &results[i].Evidence[j]
			if evidence.Kind != "software" {
				results[i].Status = worstStatus(results[i].Status, "not-run")
				continue
			}
			evidence.Status = "pass"
			for k := range evidence.Tests {
				execution := &evidence.Tests[k]
				if execution.Status == "pass" && len(execution.Skipped) != 0 {
					execution.Status = "skip"
				}
				if execution.Status != "pass" {
					evidence.Status = worstStatus(evidence.Status, execution.Status)
					failures = append(failures, fmt.Errorf("%s: %s %s: %s", results[i].ID, evidence.Package, execution.Test, execution.Status))
				}
			}
			if len(evidence.Tests) == 0 {
				evidence.Status = "missing"
				failures = append(failures, fmt.Errorf("%s: evidence has no tests", results[i].ID))
			}
			if evidence.Status != "pass" {
				results[i].Status = worstStatus(results[i].Status, evidence.Status)
			}
		}
	}
	return errors.Join(failures...)
}

func worstStatus(current, next string) string {
	priority := map[string]int{"pass": 0, "not-run": 1, "skip": 2, "missing": 3, "fail": 4}
	if priority[next] > priority[current] {
		return next
	}
	return current
}

// Run executes only the exact software tests resolved from the index and
// streams the unmodified go test -json events to artifact. Hardware and soak
// entries stay explicitly not-run, regardless of environment configuration.
func Run(ctx context.Context, directory string, index Index, artifact io.Writer) (Report, error) {
	report := Report{SchemaVersion: 1, CreatedAt: time.Now().UTC(), Status: "fail"}
	catalog, err := Discover(ctx, directory, index)
	if err != nil {
		return report, err
	}
	report.Requirements, err = Resolve(index, catalog)
	if err != nil {
		return report, err
	}
	var failures []error
	for _, pkg := range packages(index) {
		names := make(map[string]bool)
		for _, requirement := range report.Requirements {
			for _, evidence := range requirement.Evidence {
				if evidence.Package == pkg && evidence.Kind == "software" {
					for _, execution := range evidence.Tests {
						names[execution.Test] = true
					}
				}
			}
		}
		if len(names) == 0 {
			continue
		}
		var patterns []string
		for name := range names {
			patterns = append(patterns, regexp.QuoteMeta(name))
		}
		slices.Sort(patterns)
		pattern := "^(" + strings.Join(patterns, "|") + ")$"
		args := []string{"test", "-json", "-count=1", "-timeout=10m", "-run", pattern, pkg}
		if err := runGo(ctx, directory, args, artifact, func(event Event) { Record(report.Requirements, event) }); err != nil {
			failures = append(failures, fmt.Errorf("%s: %w", pkg, err))
		}
	}
	if err := Summarize(report.Requirements); err != nil {
		failures = append(failures, err)
	}
	for _, err := range failures {
		report.Errors = append(report.Errors, err.Error())
	}
	if len(failures) == 0 {
		report.Status = "pass"
	}
	return report, errors.Join(failures...)
}

func runGo(ctx context.Context, directory string, args []string, artifact io.Writer, observe func(Event)) error {
	command := exec.CommandContext(ctx, "go", args...) // #nosec G204 -- fixed Go executable; validated local package/test arguments.
	command.Dir = directory
	for _, value := range os.Environ() {
		if !strings.HasPrefix(value, "N2K_ACTISENSE_HARDWARE_CONFIG=") && !strings.HasPrefix(value, "N2K_SOAK_DURATION=") {
			command.Env = append(command.Env, value)
		}
	}
	command.Stderr = os.Stderr
	stdout, err := command.StdoutPipe()
	if err != nil {
		return err
	}
	if err := command.Start(); err != nil {
		return err
	}
	var reader io.Reader = stdout
	if artifact != nil {
		reader = io.TeeReader(stdout, artifact)
	}
	decoder := json.NewDecoder(reader)
	var decodeErr error
	var diagnostic string
	for {
		var event Event
		if err := decoder.Decode(&event); err != nil {
			if !errors.Is(err, io.EOF) {
				decodeErr = err
				_ = command.Process.Kill()
			}
			break
		}
		if event.Action == "output" {
			diagnostic += event.Output
			if len(diagnostic) > 8192 {
				diagnostic = diagnostic[len(diagnostic)-8192:]
			}
		}
		observe(event)
	}
	if err := errors.Join(decodeErr, command.Wait()); err != nil {
		return fmt.Errorf("%w\n%s", err, diagnostic)
	}
	return nil
}
