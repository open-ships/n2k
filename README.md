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

// Write a message. The struct knows its own PGN number.
// Priority defaults to 6, destination defaults to broadcast (255).
h := uint64(15708)
heading := &pgn.VesselHeading{
    Heading: &h,
}
result := client.Write(heading)
if err := result.Wait(); err != nil {
    log.Printf("write failed: %v", err)
}

// Explicitly set priority and destination
heading2 := &pgn.VesselHeading{
    Info:    pgn.MessageInfo{Priority: pgn.Priority(2), TargetId: pgn.Target(42)},
    Heading: &h,
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

### Address Claiming

Every device that transmits on NMEA 2000 must claim a unique bus address (1–253) using the ISO 11783 address claim protocol (PGN 60928). `NewClient` handles this automatically — it broadcasts an address claim, waits for contention, and only returns once a valid address is secured.

**How contention works:** Each device has a 64-bit NAME. When two devices claim the same address, the lower NAME wins and keeps the address; the loser must yield. The client supports two modes:

```go
// Auto mode (default) — starts at address 253 and negotiates downward on
// contention. If all addresses are exhausted, NewClient returns an error.
client, err := n2k.NewClient(ctx, n2k.CAN("can0"))

// Explicit mode — uses a fixed address. If another device with a lower NAME
// contests it, NewClient returns an error instead of retrying.
client, err := n2k.NewClient(ctx,
    n2k.CAN("can0"),
    n2k.WithSourceAddress(42),
)
```

**Device NAME:** The NAME determines who wins contention. Lower NAME = higher priority. Customize it to control your device's identity and arbitration priority on the bus:

```go
client, err := n2k.NewClient(ctx,
    n2k.CAN("can0"),
    n2k.WithName(n2k.DeviceName{
        IndustryGroup:    4,     // 3 bits: 4 = Marine
        ManufacturerCode: 2000,  // 11 bits: unassigned/experimental range
        DeviceClass:      25,    // 7 bits: 25 = Internetwork Device
        DeviceFunction:   130,   // 8 bits: 130 = PC Gateway
        DeviceInstance:   0,     // 8 bits
        SystemInstance:   0,     // 4 bits
        IdentityNumber:   12345, // 21 bits: unique per physical device
    }),
)
```

When `WithName` is not set, `DefaultDeviceName()` is used — it randomizes the identity number so multiple clients from the same binary can coexist on one bus.

**Claim timeout:** `NewClient` blocks for up to 1500ms (the default) to allow the network to respond to the initial claim. On heavily contested buses, increase it:

```go
client, err := n2k.NewClient(ctx,
    n2k.CAN("can0"),
    n2k.WithClaimTimeout(3 * time.Second),
)
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
    n2k.Filter("pgn == 127250"),
) { ... }

// Filter on decoded fields
for msg, err := range n2k.Receive(ctx,
    n2k.CAN("can0"),
    n2k.Filter("pgn == 127250 && msg.Heading > 3.14"),
) { ... }

// Filter by source address
for msg, err := range n2k.Receive(ctx,
    n2k.CAN("can0"),
    n2k.Filter("source == 3"),
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

When writing, `Priority` and `TargetId` default to 6 and 255 respectively when nil. When reading, they are populated from the wire. Use the helpers `pgn.Priority(v)` and `pgn.Target(v)` for concise literal construction:

```go
info := pgn.MessageInfo{
    PGN:      126996,
    SourceId:  3,
    Priority: pgn.Priority(2),
    TargetId: pgn.Target(42),
}
```

## Physical Values

Decoded numeric fields on PGN structs (`*uint64`/`*int64`) hold **raw wire ticks**, not physical
quantities. A field's metadata `Resolution` (and, for some fields, an additive `Offset`) is what
converts a tick count into a physical value in the field's `Unit` -- for example, `VesselHeading`'s
`Heading` field is transmitted in units of 0.0001 radian, so a decoded value of `15708` means
`1.5708` rad, not `15708` rad:

```go
h := uint64(15708)
heading := &pgn.VesselHeading{Heading: &h}

// field order 2 is Heading -- see the struct's n2k tag or its PgnInfo metadata
v, unit, ok, err := pgn.PhysicalValue(heading, 2)
if err != nil {
    panic(err)
}
if ok {
    fmt.Printf("heading: %.4f %s\n", v, unit) // heading: 1.5708 rad
}
```

`PhysicalValue` looks up the field by its metadata source order (not by struct field name), applies
`raw*Resolution + Offset`, and returns the field's unit label. `ok` is `false` (with a `nil` error)
when the field's decoded value is `nil` -- the wire sent the field's null/out-of-range sentinel, or
the payload ended before reaching it. It returns an error for an unknown struct or field order, a
non-numeric field (strings, binary data, `FLOAT` fields, Match-selector fields), or a field that
only exists inside a repeating group (those aren't addressable by a bare source order -- decode the
group slice and inspect its elements instead).

## Unit Types

The `units` package provides type-safe quantity wrappers (`Distance`, `Velocity`, `Pressure`, ...)
with built-in unit conversion, but it is a **standalone library**: nothing in the decode path
constructs or returns `units` values today. `pgn.PhysicalValue` (above) returns a plain `float64` in
the field's native metadata unit. Wrapping that into a `units` type is left to the caller; wiring
`units` directly into decoding would require switching PGN struct fields from raw-tick
`*uint64`/`*int64` to float-based quantities, which is a larger, deliberately deferred change (see
the `units` package doc comment).

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

## License

MIT -- see LICENSE.

## Acknowledgments

### [boatkit-io/n2k](https://github.com/boatkit-io/n2k/)

This project is inspired by [boatkit-io/n2k](https://github.com/boatkit-io/n2k/).
