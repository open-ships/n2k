package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"

	"github.com/open-ships/n2k"
)

func main() {
	iface := flag.String("i", "can0", "CAN interface name")
	usb := flag.String("u", "", "USB-CAN serial port (e.g., /dev/ttyUSB0)")
	expr := flag.String("f", "", "CEL filter expression")
	unknown := flag.Bool("unknown", false, "include unknown PGNs")
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	var opts []n2k.Option
	if *usb != "" {
		opts = append(opts, n2k.USB(*usb))
	} else {
		opts = append(opts, n2k.CAN(*iface))
	}
	if *expr != "" {
		opts = append(opts, n2k.Filter(*expr))
	}
	if *unknown {
		opts = append(opts, n2k.IncludeUnknown())
	}
	opts = append(opts, n2k.WithLogger(slog.Default()))

	enc := json.NewEncoder(os.Stdout)
	for msg, err := range n2k.Receive(ctx, opts...) {
		if err != nil {
			log.Fatal(err)
		}
		if err := enc.Encode(msg); err != nil {
			fmt.Fprintf(os.Stderr, "encode error: %v\n", err)
		}
	}
}
