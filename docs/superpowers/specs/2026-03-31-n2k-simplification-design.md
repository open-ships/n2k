# N2K Package Simplification Design

## Goal

Simplify the n2k package so users can read CAN frames (from SocketCAN or USB) and receive decoded Go structs through a clean, idiomatic API. Collapse the internal pipeline, eliminate handler chaining, and add CEL-based filtering.

## Public API

### Iterator (primary)

```go
for msg, err := range n2k.Receive(ctx,
    n2k.CAN("can0"),
    n2k.CAN("can1"),
    n2k.USB("/dev/cu.usbserial-1234"),
    n2k.Filter(`pgn == 127250 && msg.Heading > 3.14`),
    n2k.IncludeUnknown(),
    n2k.WithLogger(myLogger),
) {
    switch m := msg.(type) {
    case *pgn.VesselHeading:
        fmt.Println(m.Heading)
    case *pgn.WindData:
        fmt.Println(m.WindSpeed)
    case *pgn.UnknownPGN:
        fmt.Println(m.Reason)
    }
}
```

### Scanner (alternative)

```go
s := n2k.NewScanner(ctx,
    n2k.CAN("can0"),
    n2k.Filter(`source == 3`),
)
for s.Next() {
    switch msg := s.Message().(type) {
    case *pgn.VesselHeading:
        // ...
    }
}
if err := s.Err(); err != nil {
    log.Fatal(err)
}
```

### Function Signatures

- `Receive(ctx context.Context, opts ...Option) iter.Seq2[any, error]`
- `NewScanner(ctx context.Context, opts ...Option) *Scanner`
- `Scanner.Next() bool`
- `Scanner.Message() any`
- `Scanner.Err() error`

### Options

| Function | Description |
|----------|-------------|
| `CAN(iface string)` | Add a SocketCAN source (e.g., `"can0"`) |
| `USB(port string)` | Add a USB-CAN source (e.g., `"/dev/cu.usbserial-1234"`) |
| `Replay(frames []can.Frame)` | Add a replay source for testing |
| `Filter(expr string)` | CEL filter expression (auto-partitioned) |
| `IncludeUnknown()` | Include `*pgn.UnknownPGN` in the stream (default: drop and log at debug) |
| `WithLogger(l *slog.Logger)` | Override logger (default: `slog.Default()`) |

At least one source (`CAN`, `USB`, or `Replay`) is required. Multiple sources are fan-in'd internally, interleaved by arrival order.

## CEL Filtering

A single `Filter()` expression is automatically partitioned into two stages for efficiency.

### Stage 1: Pre-decode filter

Runs on every frame before decoding. Available variables:

| Variable | Type | Description |
|----------|------|-------------|
| `pgn` | `int` | Parameter Group Number |
| `source` | `int` | Source address (0-252) |
| `priority` | `int` | Message priority (0-7) |
| `destination` | `int` | Destination address (255 = broadcast) |

If the expression only references these variables, decoding is skipped entirely for non-matches.

### Stage 2: Post-decode filter

Runs only on decoded structs that passed stage 1. Available variables:

| Variable | Type | Description |
|----------|------|-------------|
| `msg` | `map[string]any` | Decoded struct fields as a flat map |

Struct fields are accessible by Go field name or lowercase equivalent. Both `msg.Heading` and `msg.heading` work. Unit types serialize as their numeric value for comparison.

### Partitioning Logic

At `Filter()` call time, the CEL expression is compiled and the AST is walked:

- All variables are metadata fields -> pure pre-filter, no decoding needed
- Any `msg.*` field referenced -> split at AND boundaries into pre-filter (metadata predicates) and post-filter (field predicates)
- Only `msg.*` fields -> decode everything, filter after
- OR that mixes metadata and field references -> cannot split, decode all, evaluate full expression

### Examples

```
"pgn == 127250"                          -> pre-filter only, skip decode
"pgn == 127250 && msg.Heading > 3.14"   -> pre-filter on pgn, decode, post-filter on Heading
"msg.WindSpeed > 10.0"                   -> decode all, post-filter on WindSpeed
"source == 3 || msg.Heading > 1.0"       -> can't split (OR), decode all, eval full expression
```

### Error Handling

CEL compilation and type-checking errors are caught at construction time. `Receive` yields `(nil, error)` immediately on first iteration. `NewScanner` fails on first `Next()` call. No runtime filter errors are possible.

## Package Structure

```
n2k/
  n2k.go                # Receive(), option types
  scanner.go             # NewScanner(), Scanner struct
  options.go             # CAN(), USB(), Replay(), Filter(), IncludeUnknown(), WithLogger()
  filter.go              # CEL compilation, AST partitioning, two-stage eval
  source.go              # Source interface, fan-in logic for multiple sources
  internal/
    canbus/
      socketcan.go       # SocketCAN via netlink + brutella/can
      usbcan.go          # USB-CAN serial protocol
    adapter/
      adapter.go         # CAN ID parsing (PGN, source, priority extraction)
      fastpacket.go      # Multi-frame assembly (MultiBuilder, sequence)
    decoder/
      decoder.go         # PGNDataStream, decode dispatch, UnknownPGN fallback
  pgn/                   # Exported: generated structs, decoders, registry
    pgninfo_generated.go
    pgntypes_generated.go
    ...
  units/                 # Exported: unit types with conversions
    distance.go
    temperature.go
    ...
```

### Migration from Current Structure

| Current | New | Notes |
|---------|-----|-------|
| `pkg/endpoint/` | Eliminated | Absorbed into `n2k.go` and `source.go` |
| `pkg/canbus/` | `internal/canbus/` | Made internal |
| `pkg/adapter/canadapter/` | `internal/adapter/` | Made internal |
| `pkg/pkt/` | `internal/decoder/` | Made internal |
| `pkg/pgn/` | `pgn/` | Stays exported, unchanged |
| `pkg/units/` | `units/` | Stays exported, unchanged |
| Handler interfaces | Eliminated | Replaced by internal channels/calls |
| `pkt.Packet` type | Internal to `internal/adapter/` | Still needed for fast-packet assembly, invisible to users |

## Error Handling

### Source errors (hardware failure, disconnection)

- `Receive` yields `(nil, error)` and the iterator terminates
- `Scanner.Next()` returns `false`, error available via `Scanner.Err()`
- If multiple sources: one failure terminates the entire stream (fail-fast). Partial data from a boat network is likely to cause confusion.

### Decode errors (malformed frame, unknown PGN)

- If `IncludeUnknown()` is set: yields `(*pgn.UnknownPGN, nil)` as a valid message
- If `IncludeUnknown()` is not set: silently dropped, logged at `slog.Debug` level
- Decode errors never terminate the stream

### Filter errors (invalid CEL expression)

- Caught at construction time, no runtime filter errors possible

## Logging

- Default: `slog.Default()`
- Override: `n2k.WithLogger(myLogger)`
- Unknown/dropped PGNs logged at `slog.Debug`

## Testing

### Unit tests

- `filter_test.go` — CEL compilation, AST partitioning (AND-splitting, OR fallback), pre/post stage assignment, case-insensitive field access
- `scanner_test.go` — Scanner lifecycle: Next/Message/Err, error termination, clean shutdown via context cancellation
- `source_test.go` — Fan-in from multiple sources, ordering, single-source failure terminates stream

### Internal tests

- `internal/adapter/` — CAN ID parsing, fast-packet assembly (out-of-order, duplicates, sparse frames)
- `internal/decoder/` — PGNDataStream bit-level reads, struct decoding against known payloads
- `internal/canbus/` — Minimal (hardware-dependent)

### Integration tests

- Replay captured CAN frames through `Receive` with various filters
- Verify correct structs with correct field values
- Uses `n2k.Replay(frames)` as the source

## Dependencies

### Existing

- `github.com/brutella/can` — SocketCAN I/O
- `go.bug.st/serial` — USB serial communication
- `github.com/vishvananda/netlink` — Linux netlink for CAN interface config

### New

- `github.com/google/cel-go` — CEL expression compilation and evaluation

### Exported sub-packages

- `pgn/` — Generated structs and decoders (users type-switch on these)
- `units/` — Unit types with conversion methods (users access these on struct fields)
