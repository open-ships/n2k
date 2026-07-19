// Command n2k is the NMEA 2000 command-line tool: decode live or recorded
// bus traffic into JSON, no code required.
//
//	n2k sniff -tcp 192.168.4.1:1457        # Yacht Devices WiFi gateway
//	n2k sniff -i can0                      # SocketCAN (Linux)
//	n2k sniff -file capture.log            # candump -L/-l replay
//
// Install with: go install github.com/open-ships/n2k/cmd/n2k@latest
package main

import (
	"bufio"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/open-ships/n2k"
	"github.com/open-ships/n2k/pgn"
)

// Version metadata, overridden at release time via -ldflags.
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

const rootUsage = `n2k -- NMEA 2000 command-line tools

Usage:
  n2k <command> [flags]

Commands:
  sniff      Decode traffic to typed JSON lines
  record     Capture owned raw observations as candump or JSON lines
  replay     Decode a candump capture, optionally with original timing
  validate   Check a source or capture for undecodable messages
  devices    Enumerate devices on a writable live network
  pgn        Describe or list the package's typed PGN support
  version    Print version information

Run "n2k <command> -h" for command flags.
`

func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, rootUsage)
		os.Exit(2)
	}
	switch os.Args[1] {
	case "sniff":
		if err := sniff(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "n2k sniff: %v\n", err)
			os.Exit(1)
		}
	case "record":
		if err := record(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "n2k record: %v\n", err)
			os.Exit(1)
		}
	case "replay":
		if err := replay(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "n2k replay: %v\n", err)
			os.Exit(1)
		}
	case "validate":
		if err := validate(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "n2k validate: %v\n", err)
			os.Exit(1)
		}
	case "devices":
		if err := devices(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "n2k devices: %v\n", err)
			os.Exit(1)
		}
	case "pgn":
		if err := pgnCommand(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "n2k pgn: %v\n", err)
			os.Exit(1)
		}
	case "version", "-v", "--version":
		fmt.Printf("n2k %s (commit %s, built %s)\n", version, commit, date)
	case "-h", "--help", "help":
		fmt.Print(rootUsage)
	default:
		fmt.Fprintf(os.Stderr, "n2k: unknown command %q\n\n%s", os.Args[1], rootUsage)
		os.Exit(2)
	}
}

func sniff(args []string) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	go func() {
		<-ctx.Done()
		stop()
	}()

	return sniffContext(ctx, args)
}

func sniffContext(ctx context.Context, args []string) error {
	flags := flag.NewFlagSet("sniff", flag.ContinueOnError)
	sourceFlags := addSourceFlags(flags)
	expr := flags.String("f", "", "CEL filter expression (e.g. 'pgn == 127250')")
	unknown := flags.Bool("unknown", false, "include undecodable messages as unknown PGNs")
	flags.Usage = func() {
		_, _ = fmt.Fprintf(flags.Output(), "Usage: n2k sniff <source> [flags]\n\nExactly one source is required: -i, -u, -file, -tcp, or -udp.\n\n")
		flags.PrintDefaults()
	}
	if err := flags.Parse(args); err != nil {
		return err
	}

	source, err := sourceFlags.option()
	if err != nil {
		flags.Usage()
		return err
	}

	opts := []n2k.Option{source}
	if *expr != "" {
		opts = append(opts, n2k.Filter(*expr))
	}
	if *unknown {
		opts = append(opts, n2k.IncludeUnknown())
	}

	enc := json.NewEncoder(os.Stdout)
	for msg, err := range n2k.Receive(ctx, opts...) {
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			return err
		}
		if err := enc.Encode(msg); err != nil {
			fmt.Fprintf(os.Stderr, "encode error: %v\n", err)
		}
	}
	return nil
}

type sourceFlagValues struct {
	iface  *string
	usb    *string
	file   *string
	tcp    *string
	udp    *string
	format *string
	timing *bool
}

func addSourceFlags(flags *flag.FlagSet) sourceFlagValues {
	return sourceFlagValues{
		iface:  flags.String("i", "", "SocketCAN interface (Linux, e.g. can0)"),
		usb:    flags.String("u", "", "USB-CAN serial port (e.g. /dev/ttyUSB0)"),
		file:   flags.String("file", "", "candump -L/-l capture file"),
		tcp:    flags.String("tcp", "", "TCP gateway address (e.g. 192.168.4.1:1457)"),
		udp:    flags.String("udp", "", "UDP listen address (e.g. :1457)"),
		format: flags.String("format", "raw", "stream format for -tcp/-udp: raw or actisense"),
		timing: flags.Bool("timing", false, "with -file: pace by capture timestamps"),
	}
}

func (values sourceFlagValues) option() (n2k.Option, error) {
	return sourceOption(*values.iface, *values.usb, *values.file, *values.tcp, *values.udp, *values.format, *values.timing)
}

func commandContext() (context.Context, context.CancelFunc) {
	return signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
}

func record(args []string) error {
	ctx, stop := commandContext()
	defer stop()
	return recordContext(ctx, args)
}

func recordContext(ctx context.Context, args []string) error {
	flags := flag.NewFlagSet("record", flag.ContinueOnError)
	sourceFlags := addSourceFlags(flags)
	outputPath := flags.String("out", "-", "output path, or - for stdout")
	outputFormat := flags.String("output-format", "candump", "candump or jsonl")
	flags.Usage = func() {
		_, _ = fmt.Fprintf(flags.Output(), "Usage: n2k record <source> [-out capture.log] [-output-format candump|jsonl]\n\n")
		flags.PrintDefaults()
	}
	if err := flags.Parse(args); err != nil {
		return err
	}
	source, err := sourceFlags.option()
	if err != nil {
		flags.Usage()
		return err
	}
	if *outputFormat != "candump" && *outputFormat != "jsonl" {
		return fmt.Errorf("unknown -output-format %q: use candump or jsonl", *outputFormat)
	}

	writer, closeWriter, err := outputWriter(*outputPath)
	if err != nil {
		return err
	}
	defer closeWriter()
	buffered := bufio.NewWriter(writer)
	defer func() { _ = buffered.Flush() }()
	encoder := json.NewEncoder(buffered)
	for observation, observeErr := range n2k.Observe(ctx, source) {
		if observeErr != nil {
			if ctx.Err() != nil {
				return nil
			}
			return observeErr
		}
		if *outputFormat == "jsonl" {
			if err := encoder.Encode(observation); err != nil {
				return err
			}
			continue
		}
		if observation.Frame == nil {
			continue
		}
		if _, err := fmt.Fprintln(buffered, formatCandump(observation)); err != nil {
			return err
		}
	}
	return buffered.Flush()
}

func outputWriter(path string) (io.Writer, func(), error) {
	if path == "" || path == "-" {
		return os.Stdout, func() {}, nil
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600) // #nosec G304 -- the CLI writes the operator-selected capture path.
	if err != nil {
		return nil, func() {}, fmt.Errorf("opening output: %w", err)
	}
	return file, func() { _ = file.Close() }, nil
}

func formatCandump(observation n2k.Observation) string {
	timestamp := observation.Timestamp
	if timestamp.IsZero() {
		timestamp = observation.ReceivedAt
	}
	network := observation.NetworkID
	if network == "" || strings.ContainsAny(network, " \t") {
		network = "n2k"
	}
	frame := observation.Frame
	data := strings.ToUpper(hex.EncodeToString(frame.Data[:frame.Length]))
	return fmt.Sprintf("(%d.%06d) %s %08X#%s", timestamp.Unix(), timestamp.Nanosecond()/1000, network, frame.ID, data)
}

func replay(args []string) error {
	ctx, stop := commandContext()
	defer stop()
	flags := flag.NewFlagSet("replay", flag.ContinueOnError)
	file := flags.String("file", "", "candump capture file (or pass it positionally)")
	timing := flags.Bool("timing", true, "pace by capture timestamps")
	expr := flags.String("f", "", "CEL filter expression")
	unknown := flags.Bool("unknown", false, "include undecodable messages")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if *file == "" && flags.NArg() == 1 {
		*file = flags.Arg(0)
	}
	if *file == "" || flags.NArg() > 1 {
		return errors.New("a single candump capture path is required")
	}
	var fileOpts []n2k.FileOption
	if *timing {
		fileOpts = append(fileOpts, n2k.OriginalTiming())
	}
	opts := []n2k.Option{n2k.File(*file, fileOpts...)}
	if *expr != "" {
		opts = append(opts, n2k.Filter(*expr))
	}
	if *unknown {
		opts = append(opts, n2k.IncludeUnknown())
	}
	encoder := json.NewEncoder(os.Stdout)
	for message, receiveErr := range n2k.Receive(ctx, opts...) {
		if receiveErr != nil {
			if ctx.Err() != nil {
				return nil
			}
			return receiveErr
		}
		if err := encoder.Encode(message); err != nil {
			return err
		}
	}
	return nil
}

type validationSummary struct {
	Messages    int            `json:"messages"`
	Typed       int            `json:"typed"`
	Undecodable int            `json:"undecodable"`
	ByPGN       map[uint32]int `json:"byPgn"`
}

func validate(args []string) error {
	ctx, stop := commandContext()
	defer stop()
	flags := flag.NewFlagSet("validate", flag.ContinueOnError)
	sourceFlags := addSourceFlags(flags)
	strict := flags.Bool("strict", false, "return a failure when undecodable messages are found")
	if err := flags.Parse(args); err != nil {
		return err
	}
	source, err := sourceFlags.option()
	if err != nil {
		return err
	}
	summary := validationSummary{ByPGN: make(map[uint32]int)}
	for message, receiveErr := range n2k.Receive(ctx, source, n2k.IncludeUnknown()) {
		if receiveErr != nil {
			if ctx.Err() != nil {
				break
			}
			return receiveErr
		}
		summary.Messages++
		if unknown, ok := message.(*pgn.UnknownPGN); ok {
			summary.Undecodable++
			summary.ByPGN[unknown.Info.PGN]++
			continue
		}
		summary.Typed++
		if carrier, ok := message.(interface{ MessageInfo() pgn.MessageInfo }); ok {
			summary.ByPGN[carrier.MessageInfo().PGN]++
		}
	}
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(summary); err != nil {
		return err
	}
	if *strict && summary.Undecodable > 0 {
		return fmt.Errorf("%d undecodable messages", summary.Undecodable)
	}
	return nil
}

func devices(args []string) error {
	ctx, stop := commandContext()
	defer stop()
	flags := flag.NewFlagSet("devices", flag.ContinueOnError)
	sourceFlags := addSourceFlags(flags)
	wait := flags.Duration("wait", 3*time.Second, "discovery window")
	claimTimeout := flags.Duration("claim-timeout", 2*time.Second, "address-claim startup timeout")
	if err := flags.Parse(args); err != nil {
		return err
	}
	source, err := sourceFlags.option()
	if err != nil {
		return err
	}
	if *sourceFlags.file != "" || *sourceFlags.udp != "" {
		return errors.New("devices requires a writable -i, -u, or -tcp source")
	}
	opts := []n2k.Option{source, n2k.WithClaimTimeout(*claimTimeout)}
	if *sourceFlags.tcp != "" {
		opts = append(opts, n2k.WithReconnect(n2k.ReconnectPolicy{}))
	}
	client, err := n2k.NewClient(ctx, opts...)
	if err != nil {
		return err
	}
	defer func() { _ = client.Close() }()
	timer := time.NewTimer(*wait)
	select {
	case <-timer.C:
	case <-ctx.Done():
		timer.Stop()
	}
	if err := client.Err(); err != nil {
		return err
	}
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(client.Devices())
}

func pgnCommand(args []string) error {
	if len(args) != 1 {
		return errors.New("usage: n2k pgn <number|list>")
	}
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if args[0] == "list" {
		numbers := make([]uint32, 0, len(pgn.PgnInfoLookup))
		for number := range pgn.PgnInfoLookup {
			numbers = append(numbers, number)
		}
		sort.Slice(numbers, func(i, j int) bool { return numbers[i] < numbers[j] })
		for _, number := range numbers {
			for _, info := range pgn.PgnInfoLookup[number] {
				if err := encoder.Encode(info); err != nil {
					return err
				}
			}
		}
		return nil
	}
	number, err := strconv.ParseUint(args[0], 0, 32)
	if err != nil {
		return fmt.Errorf("invalid PGN %q", args[0])
	}
	infos := pgn.PgnInfoLookup[uint32(number)]
	if len(infos) == 0 {
		return fmt.Errorf("PGN %d is not in the typed metadata", number)
	}
	return encoder.Encode(infos)
}

// sourceOption converts the mutually exclusive source flags into the single
// n2k source Option they select.
func sourceOption(iface, usb, file, tcp, udp, format string, timing bool) (n2k.Option, error) {
	var sources []n2k.Option
	if iface != "" {
		sources = append(sources, n2k.CAN(iface))
	}
	if usb != "" {
		sources = append(sources, n2k.USB(usb))
	}
	if file != "" {
		var fileOpts []n2k.FileOption
		if timing {
			fileOpts = append(fileOpts, n2k.OriginalTiming())
		}
		sources = append(sources, n2k.File(file, fileOpts...))
	}
	if tcp != "" || udp != "" {
		streamFormat, err := streamFormat(format)
		if err != nil {
			return nil, err
		}
		if tcp != "" {
			sources = append(sources, n2k.TCP(tcp, streamFormat))
		}
		if udp != "" {
			sources = append(sources, n2k.UDP(udp, streamFormat))
		}
	}
	if len(sources) != 1 {
		return nil, fmt.Errorf("exactly one source is required: -i, -u, -file, -tcp, or -udp")
	}
	if timing && file == "" {
		return nil, fmt.Errorf("-timing only applies to -file sources")
	}
	return sources[0], nil
}

func streamFormat(name string) (n2k.StreamFormat, error) {
	switch name {
	case "raw":
		return n2k.FormatYDRaw, nil
	case "actisense":
		return n2k.FormatActisense, nil
	default:
		return 0, fmt.Errorf("unknown -format %q: use raw or actisense", name)
	}
}
