// Command conformance executes public software evidence. Hardware and soak
// claims remain explicitly not-run and cannot be enabled by this command.
package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/open-ships/n2k/internal/conformance"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	indexPath := flag.String("index", "conformance/requirements.json", "requirement index")
	artifacts := flag.String("artifacts", "conformance-artifacts", "directory for report and raw test events")
	check := flag.Bool("check", false, "validate executable references without running tests")
	timeout := flag.Duration("timeout", 20*time.Minute, "total discovery/execution deadline")
	flag.Parse()
	data, err := os.ReadFile(*indexPath)
	if err != nil {
		return err
	}
	var index conformance.Index
	if err := json.Unmarshal(data, &index); err != nil {
		return err
	}
	if *timeout <= 0 {
		return fmt.Errorf("timeout must be positive")
	}
	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()
	if *check {
		catalog, err := conformance.Discover(ctx, ".", index)
		if err != nil {
			return err
		}
		_, err = conformance.Resolve(index, catalog)
		return err
	}
	if err := os.MkdirAll(*artifacts, 0o750); err != nil {
		return err
	}
	events, err := os.Create(filepath.Join(*artifacts, "local-test-events.jsonl"))
	if err != nil {
		return err
	}
	report, runErr := conformance.Run(ctx, ".", index, events)
	closeErr := events.Close()
	if runErr != nil && len(report.Errors) == 0 {
		report.Errors = append(report.Errors, runErr.Error())
	}
	report.GoVersion = runtime.Version()
	commit, commitErr := exec.CommandContext(ctx, "git", "rev-parse", "HEAD").Output()
	if commitErr == nil {
		report.Commit = strings.TrimSpace(string(commit))
	}
	status, statusErr := exec.CommandContext(ctx, "git", "status", "--porcelain").Output()
	report.WorkingTreeDirty = statusErr != nil || len(status) != 0
	hash := sha256.Sum256(data)
	report.RequirementsSHA256 = hex.EncodeToString(hash[:])
	encoded, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(*artifacts, "local-evidence.json"), append(encoded, '\n'), 0o600); err != nil {
		return err
	}
	fmt.Printf("software conformance: %s; hardware and soak evidence not run; report: %s\n", report.Status, filepath.Join(*artifacts, "local-evidence.json"))
	if runErr != nil {
		return runErr
	}
	return closeErr
}
