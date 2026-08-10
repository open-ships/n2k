# n2k

[![CI](https://github.com/open-ships/n2k/actions/workflows/test.yaml/badge.svg)](https://github.com/open-ships/n2k/actions/workflows/test.yaml)
[![Go Reference](https://pkg.go.dev/badge/github.com/open-ships/n2k.svg)](https://pkg.go.dev/github.com/open-ships/n2k)
[![Release](https://img.shields.io/github/v/release/open-ships/n2k)](https://github.com/open-ships/n2k/releases)

`n2k` is a standalone Go library for NMEA 2000 (N2K) — the CAN-based network
that connects marine instruments (GPS, depth, wind, engine, autopilot) on
boats. It reads, writes, filters, and replays bus traffic from CAN hardware,
USB-CAN adapters, WiFi gateways, and capture files, decoding messages into
strongly typed Go structs while preserving byte-exact wire payloads for
re-encoding.

It is the Go foundation for NMEA 2000 software — the bus toolkit itself, not a
plugin or companion layer for another app. Use it directly in telemetry
collectors, loggers, gateways, test rigs, monitoring services, automation, and
any Go program that needs to understand or participate on an N2K bus.

## Why n2k

- **~600 typed PGN message types.** A PGN (Parameter Group Number) is
  NMEA 2000's message-type identifier; `n2k` decodes ~600 of them into Go
  structs (348 distinct numbers, including manufacturer-proprietary variants),
  generated from the community-maintained
  [canboat](https://github.com/canboat/canboat) schema — the reference
  database for open NMEA 2000 decoding. Every numeric field with a physical
  interpretation gets a generated SI-unit accessor (`heading.HeadingValue()` →
  radians) over raw wire ticks.
- **Wire-faithful re-encode.** An unchanged decoded message retains its exact
  original payload, including reserved and trailing bytes. Mutating a field
  switches encoding to the current struct values, so edits are never silently
  discarded.
- **A real bus node, not just a decoder.** `NewClient` claims an address per
  ISO 11783, heartbeats, answers product/configuration info and ISO requests,
  accepts Commanded Address assignments, and handles NMEA group functions
  (transmit / retime / pause) — the protocol behavior NMEA 2000 requires of a
  transmitting device, handled for you.
- **Pure Go, CGO-free**, cross-compiles to Linux, macOS, and Windows.

## Quick Start — No Boat Required

The repo bundles a real six-second capture from a sailing vessel
(`testdata/sample.log`), so this program runs as-is at a desk:

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/open-ships/n2k"
	"github.com/open-ships/n2k/pgn"
)

func main() {
	ctx := context.Background()

	for msg, err := range n2k.Receive(ctx, n2k.File("testdata/sample.log")) {
		if err != nil {
			log.Fatal(err)
		}
		if heading, ok := msg.(*pgn.VesselHeading); ok {
			if rad, present := heading.HeadingValue(); present {
				fmt.Printf("heading: %.4f rad\n", rad)
			}
		}
	}
}
```

On a boat, only the source option changes:

```go
n2k.CAN("can0")                            // SocketCAN (Linux)
n2k.USB("/dev/ttyUSB0")                    // USB-CAN serial adapter
n2k.TCP("192.168.4.1:1457", n2k.FormatYDRaw)  // Yacht Devices WiFi gateway
n2k.UDP(":1457", n2k.FormatYDRaw)          // same gateway, UDP broadcast
n2k.TCP("10.0.0.5:2000", n2k.FormatActisense) // Actisense-format streams
```

The TCP/UDP sources mean you can develop on your laptop against your boat's
WiFi gateway — no CAN interface, no Linux, no cross-compiling until you
deploy. TCP works for the write path too: `NewClient` over a Yacht Devices
gateway in RAW mode is a full bus node, address claiming included, over the
WiFi gateway.

### Add it to your project

```bash
go get github.com/open-ships/n2k
```

For command-line capture and diagnostics, install
[`open-ships/n2k-cli`](https://github.com/open-ships/n2k-cli).

Releases follow semantic versioning. [`VERSION`](VERSION) declares the release
baseline; v1 provides the stable Go module path that pkg.go.dev and Go tooling
recognize as production-ready. The project preserves exported v1 behavior
across minor and patch releases, adds deprecations before removal, and reserves
breaking changes for a new major module path. Fully green release automation
publishes an annotated tag, deterministic source archive, checksums, SBOM, and
separate build-provenance and SBOM attestations through the shared Open Ships
release policy.

## Using `n2k`

Everything the library does, mapped to its API:

| Area | Functionality | API |
|------|---------------|-----|
| Read decoded traffic | Decode live or recorded NMEA 2000 frames into typed PGN structs. | `Receive`, `NewScanner` |
| Observe raw traffic | Retain owned frame/message/decode-error records with transport context. | `Observe`, `Client.Observations`, `ReplayObservations` |
| Write PGNs | Encode PGN structs back to byte-preserving CAN frames or gateway messages. | `NewClient`, `Client.Write` |
| Act as a bus node | Claim or accept a commanded address, heartbeat, answer product/configuration info and ISO requests, and handle group functions. | `NewClient`, `WithName`, `WithProductInfo`, `WithConfigInfo` |
| Schedule transmissions | Broadcast PGNs periodically and let other devices retime or pause them through group functions. | `Client.BroadcastPGN`, `Client.Broadcast` |
| Request data | Send typed ISO requests and await typed replies. | `Request[T]` |
| Discover devices | Track observed devices by stable 64-bit NAME and current source address. | `Client.Devices`, `Client.DeviceAt` |
| Read many sources | Use SocketCAN, USB-CAN serial adapters, TCP/UDP gateways, candump logs, or in-memory frames. | `CAN`, `USB`, `TCP`, `UDP`, `File`, `Replay` |
| Gateway formats | Speak Yacht Devices RAW and Actisense-format gateway streams. | `FormatYDRaw`, `FormatActisense` |
| Physical values | Keep raw wire ticks while exposing generated SI-unit accessors. | `<Field>Value`, `Set<Field>Value`, `pgn.PhysicalValue` |
| Filter traffic | Filter by PGN metadata or decoded fields with CEL; metadata-only filters avoid decode work. | `Filter` |
| Preserve unknowns | Drop undecoded PGNs by default or surface them for logging/research. | `IncludeUnknown`, `*pgn.UnknownPGN` |
| Test without hardware | Replay bundled or custom captures and inspect written frames in tests. | `File`, `OriginalTiming`, `Replay`, `WrittenFrames` |
| Inspect PGN support | Query complete runtime metadata, fields, ranges, and confidence. | `pgn.PgnInfoLookup` |
| Observe health | Export connection epochs, queue/counter metrics, address state, and structured logs. | `Client.Status`, `Client.Err`, `WithLogger` |
| Extend transports | Provide custom frame buses, observation buses, readiness, lifecycle, or assembled-message writers. | `Bus`, `ObservationBus`, `ReadyBus`, `ConnectionLifecycleBus`, `MessageWriter`, `WithBus` |

### Reading and Writing

`Client` provides read and write access to NMEA 2000. Use it when you need to transmit messages in addition to receiving them.

```go
ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
defer stop()

client, err := n2k.NewClient(ctx,
    n2k.CAN("can0"), // auto-claims an available address
)
if err != nil {
    panic(err)
}
defer client.Close()

// Write a message. The struct knows its own PGN number.
// Priority defaults to 6, destination defaults to broadcast (255).
heading := &pgn.VesselHeading{}
heading.SetHeadingValue(1.5708) // radians; stored as raw wire ticks
result := client.Write(heading)
if err := result.WaitContext(ctx); err != nil {
    log.Printf("write failed: %v", err)
}

// Explicitly set priority. VesselHeading is a broadcast (PDU2) PGN, so it
// cannot carry a destination address.
heading2 := &pgn.VesselHeading{
    Info: pgn.MessageInfo{Priority: pgn.Priority(2)},
}
heading2.SetHeadingValue(1.5708)
if err := client.Write(heading2).WaitContext(ctx); err != nil { ... }

// Addressed PGNs use TargetId. ISO Request is a PDU1 PGN.
requested := uint64(126996)
request := &pgn.IsoRequest{
    Info: pgn.MessageInfo{TargetId: pgn.Target(42)},
    Pgn:  &requested,
}
if err := client.Write(request).WaitContext(ctx); err != nil { ... }

// Read messages (same as top-level API)
for msg, err := range client.Receive() {
    if err != nil {
        panic(err)
    }
    fmt.Printf("Msg: %v\n", msg)
}
```

### Address Claiming

Every device that transmits on NMEA 2000 must claim a unique bus address (0–251) using the ISO 11783 address claim protocol (PGN 60928). `NewClient` handles this automatically — it broadcasts an address claim, waits for contention, and only returns once a valid address is secured.

**How contention works:** Each device has a 64-bit NAME. When two devices claim the same address, the lower NAME wins and keeps the address; the loser must yield. The client supports two modes:

```go
// Auto mode (default) — starts at address 251 and negotiates downward on
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
        ManufacturerCode: 123,   // 11 bits: replace with your assigned code
        DeviceClass:      25,    // 7 bits: 25 = Internetwork Device
        DeviceFunction:   130,   // 8 bits: 130 = PC Gateway
        DeviceInstance:   0,     // 8 bits
        SystemInstance:   0,     // 4 bits
        IdentityNumber:   12345, // 21 bits: unique per physical device
    }),
)
```

When `WithName` is not set, `DefaultDeviceName()` is used. That development default randomizes the identity number so multiple local processes can coexist on one bus. Production devices must provide their assigned manufacturer code and a stable, persisted identity with `WithName`; a random identity is not a provisioned product identity.

Persist the last claimed address and offer it on the next start; the client
will try it first and move if it is occupied:

```go
client, err := n2k.NewClient(ctx,
    n2k.CAN("can0"),
    n2k.WithPreferredAddress(loadLastAddress()),
)
// Save client.Status().Address after claiming or during orderly shutdown.
```

The client also handles Commanded Address (PGN 65240) without application
code. It accepts only an exact nine-byte broadcast ISO transport (BAM)
transfer whose full 64-bit NAME matches the client and whose requested address
is 0–251. A valid change immediately emits a new Address Claim and places
application writes behind a fresh contention window. Mismatched NAMEs,
fast-packet or addressed lookalikes, malformed lengths, special addresses, and
commands for the current address have no effect. Observe moves through
`client.Status().Address`; user filters and unread subscriptions cannot disable
this protocol behavior.

**Claim timeout:** `NewClient` blocks for up to 1500ms (the default) to allow the network to respond to the initial claim. On heavily contested buses, increase it:

```go
client, err := n2k.NewClient(ctx,
    n2k.CAN("can0"),
    n2k.WithClaimTimeout(3 * time.Second),
)
```

### Read-only

#### Iterator API

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

#### Scanner API

```go
s := n2k.NewScanner(ctx, n2k.CAN("can0"))
defer s.Close()
for s.Next() {
    fmt.Printf("Msg: %v\n", s.Message())
}
if err := s.Err(); err != nil {
    ...
}
```

#### Raw observations

Use `Observe` when typed decoding would discard context needed by a recorder,
bridge, latency monitor, or multi-network diagnostic:

```go
for observation, err := range n2k.Observe(ctx, n2k.File("capture.log")) {
    if err != nil {
        return err
    }
    fmt.Println(observation.AdapterID, observation.NetworkID,
        observation.Timestamp, observation.Direction, observation.Frame)
}
```

`client.Observations()` additionally publishes assembled-message and
decode-error events. Every record owns its frame and payload bytes. A slow
subscriber ends with `ErrObservationOverflow` without delaying address claims,
ISO responses, decoding, or other subscribers. `pgn.MessageInfo` carries the
same timing, Adapter, network, and direction context into typed messages.

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

### Log Files and Network Gateways

Frames don't have to come from local CAN hardware:

```go
// Replay a candump -L / -l capture, as fast as possible...
for msg, err := range n2k.Receive(ctx, n2k.File("testdata/sample.log")) { ... }

// ...or paced by the log's own timestamps.
for msg, err := range n2k.Receive(ctx, n2k.File("capture.log", n2k.OriginalTiming())) { ... }

// Yacht Devices YDWG-02 (RAW server mode) over TCP or UDP.
for msg, err := range n2k.Receive(ctx, n2k.TCP("192.168.4.1:1457", n2k.FormatYDRaw)) { ... }
for msg, err := range n2k.Receive(ctx, n2k.UDP(":1457", n2k.FormatYDRaw)) { ... }

// Actisense-format streams (W2K-1 gateways, or an NGT-1 behind a TCP
// bridge). Messages arrive pre-assembled and are re-framed internally so
// they flow through the same decode pipeline.
for msg, err := range n2k.Receive(ctx, n2k.TCP("10.0.0.5:2000", n2k.FormatActisense)) { ... }
```

`File` and `UDP` are read-only (use them with `Receive`/`NewScanner`). `TCP`
also works with `NewClient` for full read/write bus access:

```go
// A complete bus device over the boat's WiFi gateway. RAW mode is
// frame-level in both directions, so address claiming, heartbeats, and
// group functions behave exactly as on CAN hardware. The gateway echoes
// transmitted frames back, so the client also observes its own traffic.
client, err := n2k.NewClient(ctx, n2k.TCP("192.168.4.1:1457", n2k.FormatYDRaw))

// Writing over Actisense-format connections also works, with one caveat:
// the protocol is message-oriented and carries no source address on sends,
// so the gateway transmits under its own claimed address and does its own
// fast-packet fragmentation.
client, err := n2k.NewClient(ctx, n2k.TCP("10.0.0.5:2000", n2k.FormatActisense))
```

#### Source support by platform

| Source | Linux | macOS | Windows | Write access |
|--------|:-----:|:-----:|:-------:|:------------:|
| `CAN` (SocketCAN) | ✅ | ❌ | ❌ | ✅ |
| `USB` (serial CAN adapter) | ✅ | ✅ | ✅ | ✅ |
| `TCP` (Yacht Devices RAW) | ✅ | ✅ | ✅ | ✅ full frame-level control |
| `TCP` (Actisense format) | ✅ | ✅ | ✅ | ✅ gateway stamps its own source address |
| `UDP` (both formats) | ✅ | ✅ | ✅ | ❌ read-only |
| `File` (candump `-L`/`-l`) | ✅ | ✅ | ✅ | ❌ read-only |
| `Replay` (in-memory frames) | ✅ | ✅ | ✅ | ✅ writes captured via `WrittenFrames` |

### Filter Messages using Common Expression Language

Filter messages using [CEL](https://github.com/google/cel-go) expressions.

`n2k` automatically optimizes filters — metadata-only expressions skip decoding entirely.

```go
// Only vessel heading messages
for msg, err := range n2k.Receive(ctx,
    n2k.CAN("can0"),
    n2k.Filter("pgn == 127250"),
) { ... }

// Filter on decoded fields -- decoded numeric fields hold raw wire ticks,
// not physical units (see "Physical Values" below); Heading is in
// 0.0001-radian ticks, so 31416 here is > pi rad.
for msg, err := range n2k.Receive(ctx,
    n2k.CAN("can0"),
    n2k.Filter("pgn == 127250 && msg.Heading > 31416"),
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
| `source` | `int` | Source address (0-251) |
| `priority` | `int` | Message priority (0-7) |
| `destination` | `int` | Destination address (255 = broadcast) |
| `msg.<field>` | varies | Decoded struct field (case-insensitive), in raw wire ticks |

Repeating-group slice fields (`Repeating1`/`Repeating2`) are not addressable in filter expressions.

### Options

| Option | Description |
|--------|-------------|
| `n2k.CAN(iface)` | SocketCAN source (e.g., `"can0"`) |
| `n2k.USB(port)` | USB-CAN serial source (e.g., `"/dev/ttyUSB0"`) |
| `n2k.File(path, ...opts)` | candump `-L`/`-l` log file source (read-only); `n2k.OriginalTiming()` paces frames by log timestamps |
| `n2k.TCP(addr, format)` | Network gateway over TCP (read/write); format is `n2k.FormatYDRaw` or `n2k.FormatActisense` |
| `n2k.UDP(listenAddr, format)` | Network gateway datagrams (read-only), same formats |
| `n2k.Replay(frames)` | Replay source for testing |
| `n2k.Filter(expr)` | CEL filter expression |
| `n2k.IncludeUnknown()` | Include undecodable messages as `*pgn.UnknownPGN` |
| `n2k.WithLogger(l)` | Override default `slog.Logger` |
| `n2k.WithSourceAddress(addr)` | Explicit source address for writes (contention is fatal) |
| `n2k.WithPreferredAddress(addr)` | Starting address for automatic claiming; use a persisted prior address |
| `n2k.WithClaimTimeout(d)` | Initial address-claim deadline; default 1500ms |
| `n2k.WithName(name)` | ISO 11783 device NAME for address claiming |
| `n2k.WithProductInfo(p)` | Product identity reported via PGN 126996 |
| `n2k.WithConfigInfo(ci)` | Installation description reported via PGN 126998 |
| `n2k.WithHeartbeatInterval(d)` | Heartbeat (PGN 126993) cadence; default 60s, 0 disables |
| `n2k.WithReceiveBuffer(n)` | Per-subscription live receive buffer; default 64 |
| `n2k.WithWriteQueue(n)` | Pending asynchronous write capacity; default 64 |
| `n2k.WithReconnect(policy)` | Auto-reconnect dropped TCP gateway connections with exponential backoff (`ReconnectPolicy{InitialBackoff, MaxBackoff}`; zero values default to 500ms/30s) |
| `n2k.WithBus(bus)` | Inject a pre-constructed `n2k.Bus` (custom transport or test fake) instead of CAN/USB sources |

### Errors, backpressure, and health

Runtime failure is surfaced through the API. A bus disconnect ends live
iterators/scanners with that error, causes later writes to fail, and appears in
`Client.Err()` and `Client.Status()`.

Queues are deliberately bounded:

- `Write` never blocks on admission. When its queue is full, the returned
  result is already complete with `ErrWriteQueueFull`.
- Each live `Receive` or `Scanner` call has an independent subscription. A
  slow subscription ends with `ErrReceiveOverflow`; it cannot stall address
  claiming, transport handling, or another reader.
- Each raw observation subscription is independently bounded and ends with
  `ErrObservationOverflow` when slow.
- Required automatic protocol traffic uses a dedicated priority queue. Queue,
  encoding, and transport failures become terminal state instead of logs.

Use `WaitContext` when a caller has a deadline, tune bounds with
`WithWriteQueue` / `WithReceiveBuffer`, and expose `Status()` from health or
metrics endpoints:

```go
status := client.Status()
fmt.Println(status.Address, status.AddressClaimed, status.Connected,
    status.ConnectionEpoch, status.WriteQueueDepth,
    status.ProtocolRequiredQueueDepth, status.ReceiveSubscribers,
    status.ObservationSubscribers, status.TerminalError)
```

### Custom transports

`WithBus` is the extension seam for hardware and gateways not included here.
Implement `Bus` for frame-level read/write access. Optionally implement
`ObservationBus` to retain transport context, `ReadyBus` when opening is
asynchronous, `ConnectionLifecycleBus` for reconnect epochs, and
`MessageWriter` when the gateway accepts assembled PGN payloads instead of CAN
frames. The client still owns claiming, protocol routing, scheduling, and
shutdown.

### A Complete Bus Node

Beyond claiming an address, a bus client implements the protocol behavior
NMEA 2000 requires of a transmitting device:

- **Heartbeat (PGN 126993)** — sent every 60 seconds automatically (tune or
  disable with `WithHeartbeatInterval`).
- **Product & configuration info (PGNs 126996/126998)** — requests from other
  devices (chartplotters, analyzers) are answered automatically. Set your
  identity with `WithProductInfo` / `WithConfigInfo`; without them a generic
  software-gateway identity is reported so the device never shows up blank.
- **ISO requests (PGN 59904)** — requests for supported PGNs are answered;
  requests addressed to us for anything else are refused with an ISO
  acknowledgement NAK, per ISO 11783-3.
- **Commanded Address (PGN 65240)** — a matching nine-byte BAM assignment
  moves the node to a claimable address, immediately reclaims it, and opens a
  fresh contention window. Malformed and non-target commands are ignored.
- **Group functions (PGN 126208)** — request group functions can transmit,
  retime, pause (interval 0), or restore (interval `0xFFFFFFFE`) any PGN the
  client transmits, including scheduled broadcasts. Unsupported group
  functions (commands, read/write fields) are acknowledged with the proper
  error codes.

All of this runs on a dedicated decode path, so a `Filter(...)` expression
never breaks protocol behavior.

**Request/response** — ask another device for a PGN and await its typed reply
(default timeout 1250ms, the ISO 11783 response time):

```go
pi, err := n2k.Request[*pgn.ProductInformation](ctx, client, 0x23)
if err == nil {
    fmt.Println(pi.ModelId, pi.SoftwareVersionCode)
}
```

**Periodic broadcasts** — transmit a PGN on a schedule; the provider runs on
every tick (return nil to skip a tick). Prefer `BroadcastPGN`, which registers
the PGN immediately for deterministic replacement and group-function retiming:

```go
stop := client.BroadcastPGN(127250, time.Second, func() pgn.Message {
    heading := &pgn.VesselHeading{}
    heading.SetHeadingValue(currentHeadingRadians())
    return heading
})
defer stop()
```

Other devices can retime or pause the broadcast with a request group function
naming its PGN. Providers run outside the scheduler goroutine; panic is
contained and `Close` is not held hostage by a blocked provider. Providers
should still return promptly because Go cannot forcibly terminate a callback.

### Device Registry

Bus clients passively map the network: every address claim, product info, and
configuration info message updates a registry keyed by 64-bit NAME (addresses
are dynamic — the NAME is the stable identity). On first sight of a device
the client requests its product and configuration info, and at startup it
broadcasts an address-claim request so the whole bus announces itself.

```go
for _, d := range client.Devices() {
    model := "(no product info yet)"
    if d.ProductInfo != nil {
        model = d.ProductInfo.ModelId
    }
    fmt.Printf("addr %d: %s (manufacturer %d, last seen %s)\n",
        d.Address, model, d.Name.ManufacturerCode, d.LastSeen)
}

// Correlate a decoded message with the device that sent it.
if vh, ok := msg.(*pgn.VesselHeading); ok {
    if dev, found := client.DeviceAt(vh.Info.SourceId); found {
        fmt.Printf("heading from %016X\n", dev.RawName)
    }
}
```

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

Numeric struct fields hold **raw wire ticks** (`*uint64`/`*int64`). Exact
round trips are preserved by an owned original-payload snapshot; raw ticks
make field edits precise and predictable. Every numeric field with a
physical interpretation also gets **generated typed accessors** that do the
unit math — SI units in, SI units out, raw ticks underneath:

```go
heading := &pgn.VesselHeading{}
heading.SetHeadingValue(1.5708)      // radians -> stored as 15708 raw ticks

rad, ok := heading.HeadingValue()    // 1.5708, true
_ = heading.Heading                  // *uint64 raw ticks, still there

depth := decoded.(*pgn.WaterDepth)
meters, ok := depth.DepthValue()     // e.g. 2.70 m from raw 270
```

The accessor's `bool` is `false` when the field is nil — the wire sent the
field's null/out-of-range sentinel, or the payload ended before reaching it.
Each accessor documents its unit and conversion (`value = raw * resolution +
offset`); the units are the schema's SI units (`rad`, `m/s`, `K`, `V`, ...).

For dynamic, metadata-driven access — when the field is only known at
runtime — `pgn.PhysicalValue(msg, fieldOrder)` performs the same conversion
by field source order and also returns the unit label:

```go
v, unit, ok, err := pgn.PhysicalValue(heading, 2) // field order 2 = Heading
// v = 1.5708, unit = "rad"
```

`PhysicalValue` returns an error for unknown fields, non-numeric fields
(strings, binary, `FLOAT`, match selectors), and fields inside repeating
groups — for those, decode the group slice and use the element structs'
accessors instead (they're generated too).

## Unit Types

The `units` package provides type-safe quantity wrappers (`Distance`, `Velocity`, `Pressure`, ...)
with built-in unit conversion. It is independent of the PGN decode path:
decoded structs keep raw wire ticks, and the generated accessors above return
plain `float64` values in each field's native schema unit. Wrapping those
values into `units` types is left to the caller; wiring `units` directly into
decoding would require changing generated PGN struct fields from raw-tick
`*uint64`/`*int64` values to quantity types, which is deliberately deferred
(see the `units` package doc comment).

## Conformance Status

The public, reproducible protocol evidence is run with `just
conformance-local` and in Linux CI. Its controlled baseline is
[ISO 11783-3:2026, edition 5](https://www.iso.org/standard/89949.html) for
transport, [ISO 11783-5:2019, edition 3](https://www.iso.org/standard/74366.html)
for address management, and [NMEA 2000 Version 3.000 with
amendments](https://www.nmea.org/nmea-2000.html) for the marine application and
certification process. The test-to-behavior matrix is in [the conformance
guide](docs/conformance.md).

This evidence is not formal NMEA certification. The official licensed tool has
not been run against this repository alone: certification requires the actual
product hardware and firmware, assigned manufacturer and product codes, the
licensed tool, its encrypted result, and NMEA validation. The tracked
[evidence template](conformance/evidence.example.json) reports those external
steps as `not-run` until real records are supplied; the project does not infer
or simulate a pass.

## Comparison

Reviewed 2026-07-18, these are relevant projects that document NMEA 2000 support.
Go projects come first because `n2k` is for Go applications; non-Go projects
follow to show the broader ecosystem. NMEA 0183-only parsers are excluded.

Legend: ✅ first-class documented support · 🟡 partial, lower-level, or
workflow-specific · ❓ capability may exist but isn't documented · ❌ not
present or not applicable.

| Project | Lang | Go app API | Typed PGNs | Byte-preserving encode | Bus-node behavior | Group functions | Best fit |
|---------|------|------------|------------|------------------------|-------------------|-----------------|----------|
| [open-ships/n2k](https://github.com/open-ships/n2k) | Go | ✅ direct top-level API | ✅ ~600 generated structs | ✅ tested round trips | ✅ claim, heartbeat, info, requests | ✅ transmit / retime / pause | Go production apps, gateways, loggers, automation |
| [boatkit-io/n2k](https://github.com/boatkit-io/n2k) | Go | ✅ endpoint/service/node packages | ✅ generated structs | 🟡 documented encode/write | ✅ node claim + heartbeat | ❓ not documented | Go apps built around its service/node architecture |
| [aldas/go-nmea-client](https://github.com/aldas/go-nmea-client) | Go | 🟡 lower-level reader APIs | ❌ no typed structs | 🟡 raw/message oriented | 🟡 basic node mapping | ❓ not documented | Raw ingestion and format conversion |
| [canboat/canboat](https://github.com/canboat/canboat) | C | ❌ not Go | ❌ analyzer output only | 🟡 tooling-oriented | ❌ analysis suite | ❌ not a bus-node library | Shell analysis, reverse engineering, schema reference |
| [ttlappalainen/NMEA2000](https://github.com/ttlappalainen/NMEA2000) | C++ | ❌ not Go | 🟡 hand-authored helpers | 🟡 message helpers | ✅ mandatory device behavior | 🟡 device-oriented support | Embedded C++ NMEA 2000 devices |
| [canboat/canboatjs](https://github.com/canboat/canboatjs) | TypeScript | ❌ not Go | 🟡 TypeScript definitions | ✅ documented encode/transmit | 🟡 device integration | 🟡 plugin/provider-specific | Node.js, Signal K, TypeScript integrations |
| [tomer-w/nmea2000](https://github.com/tomer-w/nmea2000) | Python | ❌ not Go | 🟡 generated Python fields | ✅ generated encode/decode | ❓ not documented | ❓ not documented | Python automation and Home Assistant |
| [fard-draf/korri-n2k](https://github.com/fard-draf/korri-n2k) | Rust | ❌ not Go | ✅ generated structs | 🟡 generated serialization | ✅ address claiming + fast packet | ❌ roadmap gap | Embedded Rust with selectable PGN sets |

For a Go application, `open-ships/n2k` is designed to cover the full path in
one API: typed decoding, byte-preserving writes, bus-node behavior, group
functions, and gateway/file sources. Evaluate the details that matter to your
deployment rather than relying on the summary table alone.

## Known Limitations

- The codec implements the schema's current `PGNIsProprietary` field condition;
  a future schema condition expression must be added explicitly and otherwise
  fails closed during codec-plan compilation.
- Formal NMEA certification has not been run. Products that transmit on a real
  vessel network must follow the hardware/tool workflow in [Protocol
  conformance](docs/conformance.md) before making a certification claim.
- One physical bus per client.
- Address claiming uses a 1500ms default timeout; on heavily contested buses, increase via `WithClaimTimeout`.
- `File` and `UDP` sources are read-only; writing requires `CAN`, `USB`, `TCP`, or `Replay`.
- Over Actisense-format TCP connections the gateway stamps its own source
  address on transmissions (the protocol carries none), so the client's
  claimed address is not authoritative on the wire.
- Gateway TCP connections do not auto-reconnect by default; a dropped
  connection ends the read loop. Pass `WithReconnect` to re-dial dropped
  connections with backoff. Each new connection epoch reclaims the address,
  clears stale registry topology, waits through contention, re-enumerates, and
  resumes scheduled traffic.
- Live delivery is bounded by design. A reader that cannot keep up receives
  `ErrReceiveOverflow` or `ErrObservationOverflow`; increase
  `WithReceiveBuffer` or consume faster.

## License

MIT — see LICENSE.

## Development and security

See [CONTEXT.md](CONTEXT.md) for architecture vocabulary and invariants,
[the architecture guide](docs/architecture.md) for runtime ownership,
[the conformance guide](docs/conformance.md) for standards evidence,
[CONTRIBUTING.md](CONTRIBUTING.md) for required checks, and
[SECURITY.md](SECURITY.md) for private vulnerability reporting.

## Acknowledgments

- [canboat](https://github.com/canboat/canboat) — maintains the open PGN
  schema this library's decoders are generated from. The decoded-message
  breadth here is theirs.
- [boatkit-io/n2k](https://github.com/boatkit-io/n2k/) — the Go project that
  inspired this one.
