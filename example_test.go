package n2k_test

import (
	"context"
	"fmt"
	"log"

	"github.com/open-ships/n2k"
	"github.com/open-ships/n2k/pgn"
)

// Replay the bundled capture file -- no hardware required. Swap
// n2k.File(...) for n2k.CAN("can0"), n2k.TCP("192.168.4.1:1457",
// n2k.FormatYDRaw), or n2k.USB("/dev/ttyUSB0") to read a live bus.
func ExampleReceive() {
	ctx := context.Background()

	for msg, err := range n2k.Receive(ctx, n2k.File("testdata/sample.log")) {
		if err != nil {
			log.Fatal(err)
		}
		if heading, ok := msg.(*pgn.VesselHeading); ok {
			if rad, present := heading.HeadingValue(); present {
				fmt.Printf("heading: %.4f rad\n", rad)
				return
			}
		}
	}
	// Output: heading: 1.9624 rad
}

// Filter messages with a CEL expression; metadata-only expressions skip
// decoding entirely.
func ExampleReceive_filter() {
	ctx := context.Background()

	for msg, err := range n2k.Receive(ctx,
		n2k.File("testdata/sample.log"),
		n2k.Filter("pgn == 128267"),
	) {
		if err != nil {
			log.Fatal(err)
		}
		if depth, ok := msg.(*pgn.WaterDepth); ok {
			if meters, present := depth.DepthValue(); present {
				fmt.Printf("water depth: %.2f m\n", meters)
				return
			}
		}
	}
	// Output: water depth: 2.70 m
}

// The Scanner API is an alternative to the iterator.
func ExampleNewScanner() {
	ctx := context.Background()

	scanner := n2k.NewScanner(ctx, n2k.File("testdata/sample.log"))
	messages := 0
	for scanner.Next() {
		messages++
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("decoded over 900 messages: %t\n", messages > 900)
	// Output: decoded over 900 messages: true
}

// NewClient claims a bus address and provides write access. Not runnable
// without CAN hardware, so this example is compile-only.
func ExampleNewClient() {
	ctx := context.Background()

	client, err := n2k.NewClient(ctx, n2k.CAN("can0"))
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = client.Close() }()

	heading := &pgn.VesselHeading{}
	heading.SetHeadingValue(1.5708) // radians in, raw wire ticks underneath
	if err := client.Write(heading).Wait(); err != nil {
		log.Fatal(err)
	}
}
