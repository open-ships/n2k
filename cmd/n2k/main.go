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
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/open-ships/n2k"
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
  sniff     Decode NMEA 2000 traffic to JSON lines (live bus, gateway, or capture file)
  version   Print version information

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
	flags := flag.NewFlagSet("sniff", flag.ExitOnError)
	iface := flags.String("i", "", "SocketCAN interface (Linux, e.g. can0)")
	usb := flags.String("u", "", "USB-CAN serial port (e.g. /dev/ttyUSB0)")
	file := flags.String("file", "", "candump -L/-l capture file to replay")
	tcp := flags.String("tcp", "", "TCP gateway address (e.g. 192.168.4.1:1457)")
	udp := flags.String("udp", "", "UDP listen address (e.g. :1457)")
	format := flags.String("format", "raw", "stream format for -tcp/-udp: raw (Yacht Devices RAW) or actisense")
	timing := flags.Bool("timing", false, "with -file: pace replay by the log's own timestamps")
	expr := flags.String("f", "", "CEL filter expression (e.g. 'pgn == 127250')")
	unknown := flags.Bool("unknown", false, "include undecodable messages as unknown PGNs")
	flags.Usage = func() {
		_, _ = fmt.Fprintf(flags.Output(), "Usage: n2k sniff <source> [flags]\n\nExactly one source is required: -i, -u, -file, -tcp, or -udp.\n\n")
		flags.PrintDefaults()
	}
	if err := flags.Parse(args); err != nil {
		return err
	}

	source, err := sourceOption(*iface, *usb, *file, *tcp, *udp, *format, *timing)
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

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	enc := json.NewEncoder(os.Stdout)
	for msg, err := range n2k.Receive(ctx, opts...) {
		if err != nil {
			return err
		}
		if err := enc.Encode(msg); err != nil {
			fmt.Fprintf(os.Stderr, "encode error: %v\n", err)
		}
	}
	return nil
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
