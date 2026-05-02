# n2k

[![Test](https://github.com/open-ships/n2k/actions/workflows/test.yaml/badge.svg)](https://github.com/open-ships/n2k/actions/workflows/test.yaml)
[![Lint](https://github.com/open-ships/n2k/actions/workflows/lint.yml/badge.svg)](https://github.com/open-ships/n2k/actions/workflows/lint.yml)
[![Secure](https://github.com/open-ships/n2k/actions/workflows/security.yaml/badge.svg)](https://github.com/open-ships/n2k/actions/workflows/security.yaml)

`n2k` is a Go library for reading and writing NMEA 2000 marine network messages from CAN bus hardware into strongly-typed Go structs.

## Installation

```bash
go get github.com/open-ships/n2k
```


## Using `n2k`

### Reading and Writing

`Client` provides read and write access to NMEA 2000. Use it when you need to transmit messages in addition to receiving them.

```go
ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
defer stop()

client, err := n2k.NewClient(ctx,
    n2k.CAN("can0"),
    n2k.WithSourceAddress(42),
)
if err != nil {
    panic(err)
}
defer client.Close()

// Write a message — the struct knows its own PGN number.
// Priority defaults to 6, destination defaults to broadcast (255).
heading := &pgn.VesselHeading{
    Heading: ptrFloat32(1.5708),
}
result := client.Write(heading)
if err := result.Wait(); err != nil {
    log.Printf("write failed: %v", err)
}

// Explicitly set priority and destination
heading2 := &pgn.VesselHeading{
    Info:    pgn.MessageInfo{Priority: pgn.Priority(2), TargetId: pgn.Target(42)},
    Heading: ptrFloat32(1.5708),
}
client.Write(heading2)

// Read messages (same as top-level API)
for msg, err := range client.Receive() {
    if err != nil {
        panic(err)
    }
    fmt.Printf("Msg: %v\n", msg)
}
```

### Read-only

#### Iterator API:

```go

ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
defer stop()

for msg, err := range n2k.Receive(ctx, n2k.CAN("can0")) {
    if err != nil {
        panic(err)
    }
    fmt.Printf("Msg: %v\n", msg)
}
```

#### Scanner API:

```go
s := n2k.NewScanner(ctx, n2k.CAN("can0"))
for s.Next() {
    fmt.Printf("Msg: %v\n", msg)
}
if err := s.Err(); err != nil {
    ...
}
```

### Multiple Networks

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

### Filter Messages using Common Expression Language

Filter messages using [CEL](https://github.com/google/cel-go) expressions.

`n2k` automatically optimizes filters for max performance -- metadata-only expressions skip decoding entirely.

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
| `n2k.WithSourceAddress(addr)` | Explicit source address for writes (contention is fatal) |
| `n2k.WithName(name)` | ISO 11783 device NAME for address claiming |

### Testing with Replay

```go
frames := []can.Frame{
    {ID: 0x09F10D01, Length: 8, Data: [8]uint8{1, 2, 3, 4, 5, 6, 7, 8}},
}

for msg, err := range n2k.Receive(ctx, n2k.Replay(frames)) {
    // test your message handling
}
```

## PGN Types

All decoded messages implement the `pgn.Message` interface and are pointers to typed structs in the `pgn` package. Use a type switch to handle specific message types. PGN structs are organized across category files — `system.go`, `navigation.go`, `engine.go`, etc.

```go
type Message interface {
    PGNNumber() uint32
}
```

Every struct carries a `pgn.MessageInfo` field with wire metadata:

```go
type MessageInfo struct {
    Timestamp time.Time
    Priority  *uint8
    PGN       uint32
    SourceId  uint8
    TargetId  *uint8
}
```

When writing, `Priority` and `TargetId` default to 6 and 255 respectively when nil. When reading, they are populated from the wire.

## Unit Types

Physical quantities use type-safe wrappers from the `units` package with built-in conversion methods.

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

## Known Limitations

- Cross-field validation is not yet implemented (stubs exist for future work).
- One physical bus per client.
- Address claiming uses a 1500ms default timeout; on heavily contested buses, increase via `WithClaimTimeout`.
- Transport Protocol receive through the `Client` read API delivers only the first 8 bytes of reassembled payloads (TP send works fully).

## License

MIT -- see LICENSE.

## Acknowledgments

### [canboat](https://github.com/canboat/)

The PGN definitions and decoders at the core of this library are generated from the canboat project's open-source NMEA 2000 database. canboat reverse-engineered the NMEA 2000 protocol through network observation and public sources, producing the comprehensive PGN catalog that makes libraries like this one possible. For deeper understanding of NMEA 2000 message semantics, field definitions, and manufacturer-specific PGNs, refer to the canboat documentation.

### [boatkit-io/n2k](https://github.com/boatkit-io/n2k/)

This project is inspired by [boatkit-io/n2k](https://github.com/boatkit-io/n2k/).
