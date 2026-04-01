// Command sniffer reads NMEA 2000 messages from a SocketCAN interface and prints
// each decoded struct as a JSON object to stdout, one per line.
//
// Usage:
//
//	go run ./cmd/sniffer.go [-i can0]
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"

	"github.com/open-ships/n2k/pkg/adapter/canadapter"
	"github.com/open-ships/n2k/pkg/endpoint/socketcanendpoint"
	"github.com/open-ships/n2k/pkg/pkt"
)

type jsonPrinter struct {
	enc *json.Encoder
}

func (j *jsonPrinter) HandleStruct(v any) {
	if err := j.enc.Encode(v); err != nil {
		slog.Error(fmt.Sprintf("json encode: %v", err))
	}
}

func main() {
	iface := flag.String("i", "can0", "SocketCAN interface name")
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	log := slog.New(slog.NewTextHandler(os.Stderr, nil))

	printer := &jsonPrinter{enc: json.NewEncoder(os.Stdout)}
	decoder := pkt.NewPacketStruct()
	decoder.SetOutput(printer)

	adapter := canadapter.NewCANAdapter()
	adapter.SetOutput(decoder)

	ep := socketcanendpoint.NewSocketCANEndpoint(log, *iface)
	ep.SetOutput(adapter)

	log.Info("listening", "interface", *iface)
	if err := ep.Run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
