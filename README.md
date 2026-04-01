# n2k

A Go library for decoding NMEA 2000 marine network messages from CAN bus hardware into strongly-typed Go structs.

## Install

```bash
go get github.com/open-ships/n2k/pkg/n2k
```

## Usage

### Iterator API

```go
package main

import (
    "context"
    "fmt"
    "os"
    "os/signal"

    "github.com/open-ships/n2k/pkg/n2k"
    "github.com/open-ships/n2k/pkg/pgn"
)

func main() {
    ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
    defer stop()

    for msg, err := range n2k.Receive(ctx, n2k.CAN("can0")) {
        if err != nil {
            panic(err)
        }
        switch m := msg.(type) {
        case *pgn.VesselHeading:
            fmt.Printf("Heading: %v\n", m.Heading)
        case *pgn.WindData:
            fmt.Printf("Wind Speed: %v\n", m.WindSpeed)
        }
    }
}
```

### Scanner API

```go
s := n2k.NewScanner(ctx, n2k.CAN("can0"))
for s.Next() {
    switch msg := s.Message().(type) {
    case *pgn.VesselHeading:
        fmt.Printf("Heading: %v\n", msg.Heading)
    }
}
if err := s.Err(); err != nil {
    panic(err)
}
```

### Multiple Sources

Read from multiple CAN interfaces simultaneously:

```go
for msg, err := range n2k.Receive(ctx,
    n2k.CAN("can0"),
    n2k.CAN("can1"),
    n2k.USB("/dev/ttyUSB0"),
) {
    // messages from all sources, interleaved by arrival
}
```

### CEL Filtering

Filter messages using [CEL](https://github.com/google/cel-go) expressions. The library automatically optimizes filters -- metadata-only expressions skip decoding entirely.

```go
// Only vessel heading messages
for msg, err := range n2k.Receive(ctx,
    n2k.CAN("can0"),
    n2k.Filter(`pgn == 127250`),
) { ... }

// Filter on decoded fields
for msg, err := range n2k.Receive(ctx,
    n2k.CAN("can0"),
    n2k.Filter(`pgn == 127250 && msg.Heading > 3.14`),
) { ... }

// Filter by source address
for msg, err := range n2k.Receive(ctx,
    n2k.CAN("can0"),
    n2k.Filter(`source == 3`),
) { ... }
```

**Filter variables:**

| Variable | Type | Description |
|----------|------|-------------|
| `pgn` | `int` | Parameter Group Number |
| `source` | `int` | Source address (0-252) |
| `priority` | `int` | Message priority (0-7) |
| `destination` | `int` | Destination address (255 = broadcast) |
| `msg.<field>` | varies | Decoded struct field (case-insensitive) |

### Options

| Option | Description |
|--------|-------------|
| `n2k.CAN(iface)` | SocketCAN source (e.g., `"can0"`) |
| `n2k.USB(port)` | USB-CAN serial source (e.g., `"/dev/ttyUSB0"`) |
| `n2k.Replay(frames)` | Replay source for testing |
| `n2k.Filter(expr)` | CEL filter expression |
| `n2k.IncludeUnknown()` | Include undecodable messages as `*pgn.UnknownPGN` |
| `n2k.WithLogger(l)` | Override default `slog.Logger` |

### Testing with Replay

```go
frames := []can.Frame{
    {ID: 0x09F10D01, Length: 8, Data: [8]uint8{1, 2, 3, 4, 5, 6, 7, 8}},
}

for msg, err := range n2k.Receive(ctx, n2k.Replay(frames)) {
    // test your message handling
}
```

## Sniffer CLI

Print decoded NMEA 2000 messages as JSON:

```bash
# Read from SocketCAN
go run ./cmd/sniffer.go -i can0

# Read from USB-CAN
go run ./cmd/sniffer.go -u /dev/ttyUSB0

# With CEL filter
go run ./cmd/sniffer.go -i can0 -f 'pgn == 127250'

# Include unknown PGNs
go run ./cmd/sniffer.go -i can0 -unknown

# Pipe to jq
go run ./cmd/sniffer.go -i can0 | jq .
```

## PGN Types

All decoded messages are pointers to generated structs in the `pkg/pgn` package. Use a type switch to handle specific message types. See `pkg/pgn/pgninfo_generated.go` for the full list.

Every struct embeds `pgn.MessageInfo`:

```go
type MessageInfo struct {
    Timestamp time.Time
    Priority  uint8
    PGN       uint32
    SourceId  uint8
    TargetId  uint8
}
```

## Unit Types

Physical quantities use type-safe wrappers from the `pkg/units` package with built-in conversion methods.

## Hardware

Tested with:
- **SocketCAN**: MCP2515 (SPI), PEAK PCAN-USB
- **USB-CAN**: USB-CAN Analyzer dongles (2 Mbaud serial)

Both use the NMEA 2000 standard bitrate of 250 kbps.
