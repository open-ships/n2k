// Command roundtrip replays a candump log (candump -L / -l format) through
// the n2k decode pipeline and verifies that every decoded message re-encodes
// back to the bytes that were on the wire.
//
// Usage:
//
//	go run ./cmd/roundtrip capture.log
//	candump -L vcan0 | go run ./cmd/roundtrip -
//
// Each CAN frame is fed through the same adapter the live pipeline uses
// (including fast-packet assembly), then every assembled payload is decoded
// with pgn.DecodeMessage and re-encoded with pgn.EncodeMessage. By default
// every message is printed as a wire / re-encode byte pair (with a marker
// line under any differing bytes) followed by a summary; -quiet suppresses
// the per-message output. The result is classified per message:
//
//	exact       re-encoded bytes match the wire bytes exactly
//	padded      bytes match except for trailing 0xFF padding on one side
//	value-equal bytes differ, but decoding the re-encoded payload yields an
//	            identical struct (usually reserved-bit or padding conventions)
//	MISMATCH    re-encoded payload decodes to different field values
//
// The whole log is always processed and the exit status is always 0 (2 for
// usage or I/O errors). Pass -report <file> to additionally collect every
// imperfect message (anything not exact/padded) plus the closing summary into
// a file for later review.
package main

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/brutella/can"
	"github.com/open-ships/n2k/internal/adapter"
	"github.com/open-ships/n2k/internal/decoder"
	"github.com/open-ships/n2k/pgn"
)

type outcome int

const (
	outcomeExact outcome = iota
	outcomePadded
	outcomeValueEqual
	outcomeMismatch
	outcomeDecodeFail
	outcomeEncodeFail
	outcomeUnknownPGN
	numOutcomes
)

var outcomeNames = [numOutcomes]string{
	"exact", "padded", "value-equal", "MISMATCH", "decode-fail", "encode-fail", "unknown-pgn",
}

// statKey groups counters by PGN and the struct variant it decoded into,
// since one PGN number can dispatch to several proprietary variants.
type statKey struct {
	pgn        uint32
	structName string
}

type example struct {
	key       statKey
	reason    string
	original  []byte
	reencoded []byte
}

// checker implements adapter.PacketHandler: it receives every assembled
// packet, round-trips it through the codec, and tallies the outcome.
type checker struct {
	stats       map[statKey]*[numOutcomes]int
	examples    []example
	maxExamples int
	messages    int
	quiet       bool

	// report, when non-nil, receives every message that isn't a clean byte
	// match (streamed as encountered, so memory stays flat on large logs).
	report io.Writer
	issues int
}

func (c *checker) Decode(p decoder.Packet) {
	c.messages++

	if len(p.ParseErrors) > 0 {
		c.record(statKey{p.Info.PGN, "(no metadata)"}, outcomeUnknownPGN, p.Data, nil, joinErrors(p.ParseErrors))
		return
	}

	msg, err := pgn.DecodeMessage(p.Info, p.Data)
	if err != nil {
		c.record(statKey{p.Info.PGN, "(undecoded)"}, outcomeDecodeFail, p.Data, nil, err.Error())
		return
	}
	key := statKey{p.Info.PGN, structName(msg)}

	reencoded, err := pgn.EncodeMessage(msg)
	if err != nil {
		c.record(key, outcomeEncodeFail, p.Data, nil, err.Error())
		return
	}

	if bytes.Equal(p.Data, reencoded) {
		c.record(key, outcomeExact, p.Data, reencoded, "")
		return
	}
	if paddingOnlyDiff(p.Data, reencoded) {
		c.record(key, outcomePadded, p.Data, reencoded, "")
		return
	}

	// Bytes differ beyond padding: decide whether any information was lost by
	// decoding the re-encoded payload and comparing structs.
	redecoded, err := pgn.DecodeMessage(p.Info, reencoded)
	switch {
	case err != nil:
		c.record(key, outcomeMismatch, p.Data, reencoded, "re-encoded payload failed to decode: "+err.Error())
	case structName(redecoded) != key.structName:
		c.record(key, outcomeMismatch, p.Data, reencoded,
			"re-encoded payload decoded as different variant "+structName(redecoded))
	case reflect.DeepEqual(msg, redecoded):
		c.record(key, outcomeValueEqual, p.Data, reencoded, "bytes differ but field values survive")
	default:
		c.record(key, outcomeMismatch, p.Data, reencoded, "field values changed after re-encode")
	}
}

func (c *checker) record(key statKey, o outcome, original, reencoded []byte, reason string) {
	counters := c.stats[key]
	if counters == nil {
		counters = &[numOutcomes]int{}
		c.stats[key] = counters
	}
	counters[o]++

	if c.report != nil && reason != "" {
		c.issues++
		writeMessage(c.report, c.messages, key, o, original, reencoded, reason)
	}

	if !c.quiet {
		writeMessage(os.Stdout, c.messages, key, o, original, reencoded, reason)
		return
	}

	// Quiet mode: keep the first few examples of anything that isn't a clean
	// byte match for the summary.
	if reason == "" || len(c.examples) >= c.maxExamples {
		return
	}
	c.examples = append(c.examples, example{
		key:       key,
		reason:    fmt.Sprintf("%s: %s", outcomeNames[o], reason),
		original:  append([]byte(nil), original...),
		reencoded: append([]byte(nil), reencoded...),
	})
}

// printer emits formatted output and deliberately drops write errors:
// console writes are best-effort, and the report file's errors surface at
// Flush/Close in main.
type printer struct{ w io.Writer }

func (p printer) f(format string, args ...any) { _, _ = fmt.Fprintf(p.w, format, args...) }

// writeMessage writes one message's before/after byte view, e.g.
//
//	#3 pgn 127250 VesselHeading [value-equal] bytes differ but field values survive
//	    wire       ff 5c 3d ff 7f ff 7f 03
//	    re-encode  ff 5c 3d ff 7f ff 7f ff
//	    diff                            ^^
func writeMessage(w io.Writer, idx int, key statKey, o outcome, wire, reencoded []byte, reason string) {
	p := printer{w}
	p.f("#%d pgn %d %s [%s]", idx, key.pgn, key.structName, outcomeNames[o])
	if reason != "" {
		p.f(" %s", reason)
	}
	p.f("\n")
	p.f("    %-11s% x\n", "wire", wire)
	if reencoded != nil {
		p.f("    %-11s% x\n", "re-encode", reencoded)
		if !bytes.Equal(wire, reencoded) {
			p.f("    %-11s%s\n", "diff", diffMarkers(wire, reencoded))
		}
	}
	p.f("\n")
}

// diffMarkers returns a "^^" marker aligned under every byte column where a
// and b differ (including columns present in only one of them).
func diffMarkers(a, b []byte) string {
	n := max(len(a), len(b))
	var sb strings.Builder
	for i := range n {
		if i > 0 {
			sb.WriteByte(' ')
		}
		if i < len(a) && i < len(b) && a[i] == b[i] {
			sb.WriteString("  ")
		} else {
			sb.WriteString("^^")
		}
	}
	return strings.TrimRight(sb.String(), " ")
}

func structName(msg pgn.PGN) string {
	t := reflect.TypeOf(msg)
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	return t.Name()
}

func joinErrors(errs []error) string {
	parts := make([]string, len(errs))
	for i, err := range errs {
		parts[i] = err.Error()
	}
	return strings.Join(parts, "; ")
}

// paddingOnlyDiff reports whether a and b are identical except that the
// longer one carries extra trailing 0xFF bytes (the NMEA 2000 fill value used
// to pad single-frame payloads to 8 bytes).
func paddingOnlyDiff(a, b []byte) bool {
	short, long := a, b
	if len(short) > len(long) {
		short, long = long, short
	}
	if !bytes.Equal(long[:len(short)], short) {
		return false
	}
	for _, v := range long[len(short):] {
		if v != 0xFF {
			return false
		}
	}
	return true
}

// parseLine parses one candump -L / -l log line, e.g.
//
//	(1720000000.000000) can0 09F50E7F#00FFFFFFFFFFFFFF
//
// It finds the ID#DATA field so extra or missing columns are tolerated.
// Returns ok=false for lines that carry no classic CAN data frame (RTR
// frames, CAN FD "##" frames, comments, malformed lines).
func parseLine(line string) (can.Frame, bool) {
	var idData string
	for _, field := range strings.Fields(line) {
		if strings.Contains(field, "#") {
			idData = field
			break
		}
	}
	parts := strings.SplitN(idData, "#", 2)
	if len(parts) != 2 || parts[0] == "" {
		return can.Frame{}, false
	}
	payload := parts[1]
	if strings.HasPrefix(payload, "#") || strings.HasPrefix(payload, "R") {
		return can.Frame{}, false
	}

	id, err := strconv.ParseUint(parts[0], 16, 32)
	if err != nil {
		return can.Frame{}, false
	}
	data, err := hex.DecodeString(payload)
	if err != nil || len(data) > 8 {
		return can.Frame{}, false
	}

	f := can.Frame{ID: uint32(id), Length: uint8(len(data))}
	copy(f.Data[:], data)
	return f, true
}

func main() {
	quiet := flag.Bool("quiet", false, "suppress per-message before/after output; print only the summary")
	maxExamples := flag.Int("examples", 10, "max number of problem examples to print in the summary (-quiet mode)")
	reportPath := flag.String("report", "", "write every issue (anything not exact/padded) plus the summary to this file")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: roundtrip [flags] <candump.log | ->\n")
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}

	var in io.Reader = os.Stdin
	if name := flag.Arg(0); name != "-" {
		f, err := os.Open(name) // #nosec G304 -- reading the user-supplied capture log is this tool's purpose.
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
		defer func() { _ = f.Close() }()
		in = f
	}

	chk := &checker{stats: map[statKey]*[numOutcomes]int{}, maxExamples: *maxExamples, quiet: *quiet}

	var reportFile *os.File
	var reportBuf *bufio.Writer
	if *reportPath != "" {
		f, err := os.Create(*reportPath) // #nosec G304 -- the report destination is chosen by the user.
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
		reportFile = f
		reportBuf = bufio.NewWriter(f)
		chk.report = reportBuf
	}

	canAdapter := adapter.NewCANAdapter()
	canAdapter.SetOutput(chk)

	frames, skipped := 0, 0
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		frame, ok := parseLine(line)
		if !ok {
			skipped++
			continue
		}
		frames++
		canAdapter.HandleMessage(&frame)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	writeSummary(os.Stdout, chk, frames, skipped)
	printExamples(chk)

	if reportFile != nil {
		writeSummary(reportBuf, chk, frames, skipped)
		if err := reportBuf.Flush(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
		if err := reportFile.Close(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
		fmt.Printf("\nwrote %d issues to %s\n", chk.issues, reportFile.Name())
	}
}

func writeSummary(w io.Writer, chk *checker, frames, skipped int) {
	totals := [numOutcomes]int{}
	keys := make([]statKey, 0, len(chk.stats))
	for key, counters := range chk.stats {
		keys = append(keys, key)
		for o, n := range counters {
			totals[o] += n
		}
	}
	sort.Slice(keys, func(i, j int) bool {
		if keys[i].pgn != keys[j].pgn {
			return keys[i].pgn < keys[j].pgn
		}
		return keys[i].structName < keys[j].structName
	})

	p := printer{w}
	p.f("== summary ==\n")
	p.f("frames read: %d (skipped %d unparseable lines)\n", frames, skipped)
	p.f("messages assembled: %d\n", chk.messages)
	for o, name := range outcomeNames {
		p.f("  %-12s %d\n", name, totals[o])
	}

	p.f("\n%8s  %-50s %7s %7s %7s %7s %9s %6s %6s %8s\n",
		"PGN", "struct", "total", "exact", "padded", "valeq", "MISMATCH", "dec!", "enc!", "unknown")
	for _, key := range keys {
		counters := chk.stats[key]
		total := 0
		for _, n := range counters {
			total += n
		}
		p.f("%8d  %-50s %7d %7d %7d %7d %9d %6d %6d %8d\n",
			key.pgn, key.structName, total,
			counters[outcomeExact], counters[outcomePadded], counters[outcomeValueEqual],
			counters[outcomeMismatch], counters[outcomeDecodeFail], counters[outcomeEncodeFail],
			counters[outcomeUnknownPGN])
	}
}

// printExamples prints the capped problem examples collected in -quiet mode.
func printExamples(chk *checker) {
	if len(chk.examples) == 0 {
		return
	}
	fmt.Printf("\nexamples (first %d):\n", len(chk.examples))
	for _, ex := range chk.examples {
		fmt.Printf("- pgn %d (%s) %s\n", ex.key.pgn, ex.key.structName, ex.reason)
		if len(ex.original) > 0 {
			fmt.Printf("    wire:      % x\n", ex.original)
		}
		if len(ex.reencoded) > 0 {
			fmt.Printf("    re-encode: % x\n", ex.reencoded)
		}
	}
}
