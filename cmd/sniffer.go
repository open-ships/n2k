package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"

	"github.com/brutella/can"
	"github.com/open-ships/n2k/internal/adapter"
	"github.com/open-ships/n2k/internal/canbus"
	"github.com/open-ships/n2k/internal/decoder"
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
	dec := decoder.New()
	dec.SetOutput(printer)

	a := adapter.NewCANAdapter()
	a.SetOutput(dec)

	log.Info("listening", "interface", *iface)
	err := canbus.RunSocketCAN(ctx, log, *iface, func(frame can.Frame) {
		a.HandleMessage(&frame)
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
