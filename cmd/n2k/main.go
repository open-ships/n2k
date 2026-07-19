// Command n2k is the NMEA 2000 command-line tool: decode, capture, replay,
// validate, and inspect NMEA 2000 traffic without writing Go code.
//
//	n2k sniff --tcp 192.168.4.1:1457    # Yacht Devices WiFi gateway
//	n2k sniff -i can0                   # SocketCAN (Linux)
//	n2k replay capture.log              # candump capture
//
// Install with: go install github.com/open-ships/n2k/cmd/n2k@latest
package main

import (
	"bufio"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
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
	"github.com/spf13/cobra"
)

// Version metadata, overridden at release time via -ldflags.
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	root := newRootCommand(os.Stdin, os.Stdout, os.Stderr)
	if err := root.ExecuteContext(ctx); err != nil {
		if ctx.Err() != nil {
			return
		}
		_, _ = fmt.Fprintf(root.ErrOrStderr(), "n2k: %v\n", err)
		os.Exit(1)
	}
}

func newRootCommand(in io.Reader, out, errOut io.Writer) *cobra.Command {
	root := &cobra.Command{
		Use:           "n2k",
		Short:         "NMEA 2000 capture, diagnostics, and schema tools",
		SilenceErrors: true,
		SilenceUsage:  true,
		Version:       fmt.Sprintf("%s (commit %s, built %s)", version, commit, date),
		Example: commandExamples(
			"n2k sniff --tcp 192.168.4.1:1457",
			"n2k record -i can0 --out capture.log",
			"n2k replay capture.log --timing=false",
			"n2k pgn 127250",
		),
		RunE: func(command *cobra.Command, _ []string) error {
			return command.Help()
		},
	}
	root.SetIn(in)
	root.SetOut(out)
	root.SetErr(errOut)
	root.SetVersionTemplate("n2k {{.Version}}\n")
	root.CompletionOptions.DisableDefaultCmd = false
	root.CompletionOptions.HiddenDefaultCmd = false

	root.AddCommand(
		newSniffCommand(),
		newRecordCommand(),
		newReplayCommand(),
		newValidateCommand(),
		newDevicesCommand(),
		newPGNCommand(),
		newVersionCommand(),
	)
	root.InitDefaultCompletionCmd()
	return root
}

type sourceFlagValues struct {
	iface  string
	usb    string
	file   string
	tcp    string
	udp    string
	format string
	timing bool
}

func addSourceFlags(command *cobra.Command, values *sourceFlagValues) {
	flags := command.Flags()
	flags.StringVarP(&values.iface, "interface", "i", "", "SocketCAN interface (Linux, e.g. can0)")
	flags.StringVarP(&values.usb, "usb", "u", "", "USB-CAN serial port (e.g. /dev/ttyUSB0)")
	flags.StringVar(&values.file, "file", "", "candump -L/-l capture file")
	flags.StringVar(&values.tcp, "tcp", "", "TCP gateway address (e.g. 192.168.4.1:1457)")
	flags.StringVar(&values.udp, "udp", "", "UDP listen address (e.g. :1457)")
	flags.StringVar(&values.format, "format", "raw", "stream format for --tcp/--udp: raw or actisense")
	flags.BoolVar(&values.timing, "timing", false, "with --file: pace by capture timestamps")

	mustConfigure(command.MarkFlagFilename("file"))
	mustConfigure(command.MarkFlagFilename("usb"))
	mustConfigure(command.RegisterFlagCompletionFunc("format", cobra.FixedCompletions(
		[]string{"raw\tYacht Devices RAW ASCII", "actisense\tActisense binary framing"},
		cobra.ShellCompDirectiveNoFileComp,
	)))
	mustConfigure(command.RegisterFlagCompletionFunc("interface", cobra.NoFileCompletions))
	mustConfigure(command.RegisterFlagCompletionFunc("tcp", cobra.NoFileCompletions))
	mustConfigure(command.RegisterFlagCompletionFunc("udp", cobra.NoFileCompletions))
}

func mustConfigure(err error) {
	if err != nil {
		panic(err)
	}
}

func commandExamples(examples ...string) string {
	return "  " + strings.Join(examples, "\n  ")
}

func (values sourceFlagValues) option() (n2k.Option, error) {
	return sourceOption(values.iface, values.usb, values.file, values.tcp, values.udp, values.format, values.timing)
}

func newSniffCommand() *cobra.Command {
	var source sourceFlagValues
	var expression string
	var includeUnknown bool
	command := &cobra.Command{
		Use:   "sniff",
		Short: "Decode traffic to typed JSON lines",
		Args:  cobra.NoArgs,
		Example: commandExamples(
			"n2k sniff --tcp 192.168.4.1:1457",
			"n2k sniff -i can0 --filter 'pgn == 127250'",
			"n2k sniff --file capture.log --timing --unknown",
		),
		RunE: func(command *cobra.Command, _ []string) error {
			option, err := source.option()
			if err != nil {
				return err
			}
			return runSniff(command.Context(), command.OutOrStdout(), option, expression, includeUnknown)
		},
	}
	addSourceFlags(command, &source)
	command.Flags().StringVarP(&expression, "filter", "f", "", "CEL filter expression (e.g. 'pgn == 127250')")
	command.Flags().BoolVar(&includeUnknown, "unknown", false, "include undecodable messages as unknown PGNs")
	return command
}

func runSniff(ctx context.Context, out io.Writer, source n2k.Option, expression string, includeUnknown bool) error {
	opts := []n2k.Option{source}
	if expression != "" {
		opts = append(opts, n2k.Filter(expression))
	}
	if includeUnknown {
		opts = append(opts, n2k.IncludeUnknown())
	}

	encoder := json.NewEncoder(out)
	for message, receiveErr := range n2k.Receive(ctx, opts...) {
		if receiveErr != nil {
			if ctx.Err() != nil {
				return nil
			}
			return receiveErr
		}
		if err := encoder.Encode(message); err != nil {
			return fmt.Errorf("encoding message: %w", err)
		}
	}
	return nil
}

func newRecordCommand() *cobra.Command {
	var source sourceFlagValues
	var outputPath string
	var outputFormat string
	command := &cobra.Command{
		Use:   "record",
		Short: "Capture owned raw observations",
		Args:  cobra.NoArgs,
		Example: commandExamples(
			"n2k record -i can0 --out capture.log",
			"n2k record --tcp 192.168.4.1:1457 --out observations.jsonl --output-format jsonl",
		),
		RunE: func(command *cobra.Command, _ []string) error {
			option, err := source.option()
			if err != nil {
				return err
			}
			return runRecord(command.Context(), command.OutOrStdout(), option, outputPath, outputFormat)
		},
	}
	addSourceFlags(command, &source)
	command.Flags().StringVarP(&outputPath, "out", "o", "-", "output path, or - for stdout")
	command.Flags().StringVar(&outputFormat, "output-format", "candump", "output format: candump or jsonl")
	mustConfigure(command.MarkFlagFilename("out"))
	mustConfigure(command.RegisterFlagCompletionFunc("output-format", cobra.FixedCompletions(
		[]string{"candump\tReplayable candump text", "jsonl\tOwned observation JSON lines"},
		cobra.ShellCompDirectiveNoFileComp,
	)))
	return command
}

func runRecord(ctx context.Context, stdout io.Writer, source n2k.Option, outputPath, outputFormat string) error {
	if outputFormat != "candump" && outputFormat != "jsonl" {
		return fmt.Errorf("unknown output format %q: use candump or jsonl", outputFormat)
	}

	writer, closeWriter, err := outputWriter(outputPath, stdout)
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
		if outputFormat == "jsonl" {
			if err := encoder.Encode(observation); err != nil {
				return fmt.Errorf("encoding observation: %w", err)
			}
			continue
		}
		if observation.Frame == nil {
			continue
		}
		if _, err := fmt.Fprintln(buffered, formatCandump(observation)); err != nil {
			return fmt.Errorf("writing capture: %w", err)
		}
	}
	return buffered.Flush()
}

func outputWriter(path string, stdout io.Writer) (io.Writer, func(), error) {
	if path == "" || path == "-" {
		return stdout, func() {}, nil
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

func newReplayCommand() *cobra.Command {
	var file string
	var timing bool
	var expression string
	var includeUnknown bool
	command := &cobra.Command{
		Use:   "replay [capture]",
		Short: "Decode a candump capture",
		Args:  cobra.MaximumNArgs(1),
		Example: commandExamples(
			"n2k replay capture.log",
			"n2k replay --file capture.log --timing=false --unknown",
		),
		RunE: func(command *cobra.Command, args []string) error {
			if file != "" && len(args) != 0 {
				return errors.New("pass a capture either positionally or with --file, not both")
			}
			if file == "" && len(args) == 1 {
				file = args[0]
			}
			if file == "" {
				return errors.New("a candump capture path is required")
			}
			return runReplay(command.Context(), command.OutOrStdout(), file, timing, expression, includeUnknown)
		},
	}
	command.Flags().StringVar(&file, "file", "", "candump capture file (or pass it positionally)")
	command.Flags().BoolVar(&timing, "timing", true, "pace by capture timestamps")
	command.Flags().StringVarP(&expression, "filter", "f", "", "CEL filter expression")
	command.Flags().BoolVar(&includeUnknown, "unknown", false, "include undecodable messages")
	mustConfigure(command.MarkFlagFilename("file"))
	return command
}

func runReplay(ctx context.Context, out io.Writer, file string, timing bool, expression string, includeUnknown bool) error {
	var fileOpts []n2k.FileOption
	if timing {
		fileOpts = append(fileOpts, n2k.OriginalTiming())
	}
	opts := []n2k.Option{n2k.File(file, fileOpts...)}
	if expression != "" {
		opts = append(opts, n2k.Filter(expression))
	}
	if includeUnknown {
		opts = append(opts, n2k.IncludeUnknown())
	}
	encoder := json.NewEncoder(out)
	for message, receiveErr := range n2k.Receive(ctx, opts...) {
		if receiveErr != nil {
			if ctx.Err() != nil {
				return nil
			}
			return receiveErr
		}
		if err := encoder.Encode(message); err != nil {
			return fmt.Errorf("encoding message: %w", err)
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

func newValidateCommand() *cobra.Command {
	var source sourceFlagValues
	var strict bool
	command := &cobra.Command{
		Use:   "validate",
		Short: "Check a source for undecodable messages",
		Args:  cobra.NoArgs,
		Example: commandExamples(
			"n2k validate --file capture.log",
			"n2k validate --file capture.log --strict",
		),
		RunE: func(command *cobra.Command, _ []string) error {
			option, err := source.option()
			if err != nil {
				return err
			}
			return runValidate(command.Context(), command.OutOrStdout(), option, strict)
		},
	}
	addSourceFlags(command, &source)
	command.Flags().BoolVar(&strict, "strict", false, "return a failure when undecodable messages are found")
	return command
}

func runValidate(ctx context.Context, out io.Writer, source n2k.Option, strict bool) error {
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
	encoder := json.NewEncoder(out)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(summary); err != nil {
		return fmt.Errorf("encoding summary: %w", err)
	}
	if strict && summary.Undecodable > 0 {
		return fmt.Errorf("%d undecodable messages", summary.Undecodable)
	}
	return nil
}

func newDevicesCommand() *cobra.Command {
	var source sourceFlagValues
	var wait time.Duration
	var claimTimeout time.Duration
	command := &cobra.Command{
		Use:   "devices",
		Short: "Enumerate devices on a writable live network",
		Args:  cobra.NoArgs,
		Example: commandExamples(
			"n2k devices --tcp 192.168.4.1:1457 --wait 5s",
			"n2k devices -i can0",
		),
		RunE: func(command *cobra.Command, _ []string) error {
			option, err := source.option()
			if err != nil {
				return err
			}
			return runDevices(command.Context(), command.OutOrStdout(), option, source.tcp != "", wait, claimTimeout)
		},
	}
	addWritableSourceFlags(command, &source)
	command.Flags().DurationVar(&wait, "wait", 3*time.Second, "discovery window")
	command.Flags().DurationVar(&claimTimeout, "claim-timeout", 2*time.Second, "address-claim startup timeout")
	return command
}

func addWritableSourceFlags(command *cobra.Command, values *sourceFlagValues) {
	flags := command.Flags()
	flags.StringVarP(&values.iface, "interface", "i", "", "SocketCAN interface (Linux, e.g. can0)")
	flags.StringVarP(&values.usb, "usb", "u", "", "USB-CAN serial port (e.g. /dev/ttyUSB0)")
	flags.StringVar(&values.tcp, "tcp", "", "TCP gateway address (e.g. 192.168.4.1:1457)")
	flags.StringVar(&values.format, "format", "raw", "stream format for --tcp: raw or actisense")

	mustConfigure(command.MarkFlagFilename("usb"))
	mustConfigure(command.RegisterFlagCompletionFunc("format", cobra.FixedCompletions(
		[]string{"raw\tYacht Devices RAW ASCII", "actisense\tActisense binary framing"},
		cobra.ShellCompDirectiveNoFileComp,
	)))
	mustConfigure(command.RegisterFlagCompletionFunc("interface", cobra.NoFileCompletions))
	mustConfigure(command.RegisterFlagCompletionFunc("tcp", cobra.NoFileCompletions))
}

func runDevices(ctx context.Context, out io.Writer, source n2k.Option, reconnect bool, wait, claimTimeout time.Duration) error {
	opts := []n2k.Option{source, n2k.WithClaimTimeout(claimTimeout)}
	if reconnect {
		opts = append(opts, n2k.WithReconnect(n2k.ReconnectPolicy{}))
	}
	client, err := n2k.NewClient(ctx, opts...)
	if err != nil {
		return err
	}
	defer func() { _ = client.Close() }()
	timer := time.NewTimer(wait)
	defer timer.Stop()
	select {
	case <-timer.C:
	case <-ctx.Done():
		return nil
	}
	if err := client.Err(); err != nil {
		return err
	}
	encoder := json.NewEncoder(out)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(client.Devices()); err != nil {
		return fmt.Errorf("encoding devices: %w", err)
	}
	return nil
}

func newPGNCommand() *cobra.Command {
	command := &cobra.Command{
		Use:               "pgn <number|list>",
		Short:             "Describe or list typed PGN support",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: completePGNs,
		Example: commandExamples(
			"n2k pgn 127250",
			"n2k pgn list",
		),
		RunE: func(command *cobra.Command, args []string) error {
			return runPGN(command.OutOrStdout(), args[0])
		},
	}
	return command
}

func completePGNs(_ *cobra.Command, _ []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	var completions []string
	if strings.HasPrefix("list", toComplete) {
		completions = append(completions, "list\tList every typed PGN variant")
	}
	numbers := sortedPGNNumbers()
	for _, number := range numbers {
		value := strconv.FormatUint(uint64(number), 10)
		if !strings.HasPrefix(value, toComplete) {
			continue
		}
		description := "typed PGN"
		if infos := pgn.PgnInfoLookup[number]; len(infos) > 0 && infos[0].Description != "" {
			description = strings.NewReplacer("\t", " ", "\n", " ").Replace(infos[0].Description)
		}
		completions = append(completions, value+"\t"+description)
	}
	return completions, cobra.ShellCompDirectiveNoFileComp
}

func sortedPGNNumbers() []uint32 {
	numbers := make([]uint32, 0, len(pgn.PgnInfoLookup))
	for number := range pgn.PgnInfoLookup {
		numbers = append(numbers, number)
	}
	sort.Slice(numbers, func(i, j int) bool { return numbers[i] < numbers[j] })
	return numbers
}

func runPGN(out io.Writer, value string) error {
	encoder := json.NewEncoder(out)
	encoder.SetIndent("", "  ")
	if value == "list" {
		for _, number := range sortedPGNNumbers() {
			for _, info := range pgn.PgnInfoLookup[number] {
				if err := encoder.Encode(info); err != nil {
					return fmt.Errorf("encoding PGN metadata: %w", err)
				}
			}
		}
		return nil
	}
	number, err := strconv.ParseUint(value, 0, 32)
	if err != nil {
		return fmt.Errorf("invalid PGN %q", value)
	}
	infos := pgn.PgnInfoLookup[uint32(number)]
	if len(infos) == 0 {
		return fmt.Errorf("PGN %d is not in the typed metadata", number)
	}
	if err := encoder.Encode(infos); err != nil {
		return fmt.Errorf("encoding PGN metadata: %w", err)
	}
	return nil
}

func newVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			_, err := fmt.Fprintf(command.OutOrStdout(), "n2k %s (commit %s, built %s)\n", version, commit, date)
			return err
		},
	}
}

// sourceOption converts the mutually exclusive source flags into the single
// n2k source Option they select.
func sourceOption(iface, usb, file, tcp, udp, format string, timing bool) (n2k.Option, error) {
	stream, err := streamFormat(format)
	if err != nil {
		return nil, err
	}
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
		if tcp != "" {
			sources = append(sources, n2k.TCP(tcp, stream))
		}
		if udp != "" {
			sources = append(sources, n2k.UDP(udp, stream))
		}
	}
	if len(sources) != 1 {
		return nil, errors.New("exactly one source is required: --interface, --usb, --file, --tcp, or --udp")
	}
	if timing && file == "" {
		return nil, errors.New("--timing only applies to --file sources")
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
		return 0, fmt.Errorf("unknown format %q: use raw or actisense", name)
	}
}
