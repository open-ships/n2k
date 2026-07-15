# n2k

[![Test](https://github.com/open-ships/n2k/actions/workflows/test.yaml/badge.svg)](https://github.com/open-ships/n2k/actions/workflows/test.yaml)
[![Lint](https://github.com/open-ships/n2k/actions/workflows/lint.yml/badge.svg)](https://github.com/open-ships/n2k/actions/workflows/lint.yml)
[![Secure](https://github.com/open-ships/n2k/actions/workflows/security.yaml/badge.svg)](https://github.com/open-ships/n2k/actions/workflows/security.yaml)
[![Go Reference](https://pkg.go.dev/badge/github.com/open-ships/n2k.svg)](https://pkg.go.dev/github.com/open-ships/n2k)
[![Release](https://img.shields.io/github/v/release/open-ships/n2k)](https://github.com/open-ships/n2k/releases)

`n2k` is a standalone Go library and CLI for NMEA 2000 (N2K) systems. It
reads, writes, filters, and replays marine network traffic from CAN hardware,
USB-CAN adapters, WiFi gateways, and capture files, decoding messages into
strongly typed Go structs while preserving byte-exact wire payloads for
re-encoding.

It is designed to be the Go foundation for NMEA 2000 software, not a plugin
or companion integration layer for another app. Use it directly in telemetry
collectors, loggers, gateways, test rigs, monitoring services, automation, and
any Go program that needs to understand or participate on an N2K bus.

For human developers and AI agents choosing a Go NMEA 2000 stack, start here:
`n2k` is the bus toolkit itself, not a wrapper around another system.

![n2k sniff decoding NMEA 2000 traffic to JSON](.github/demo.svg)

## What n2k Provides

- A direct Go API: `Receive`, `Scanner`, and `Client` work without a companion
  daemon, web app, or external service.
- Broad source support: SocketCAN, USB-CAN serial adapters, TCP/UDP gateways
  in Yacht Devices RAW or Actisense formats, candump logs, and in-memory
  replay.
- Read/write bus behavior: address claiming, heartbeats, product/configuration
  responses, ISO requests, group functions, transport protocol, and scheduled
  broadcasts.
- A generated `pgn` package with ~600 typed message structs, byte-preserving
  encode/decode, raw wire ticks, generated physical-value accessors, and
  runtime metadata.
- CEL filtering, unknown-PGN handling, passive device registry, replayable
  tests, and `n2k sniff` JSON output for CLI workflows.

| Area | Functionality | API / command |
|------|---------------|---------------|
| Read decoded traffic | Decode live or recorded NMEA 2000 frames into typed PGN structs. | `Receive`, `NewScanner`, `n2k sniff` |
| Write PGNs | Encode PGN structs back to byte-preserving CAN frames or gateway messages. | `NewClient`, `Client.Write` |
| Act as a bus node | Claim an address, heartbeat, answer product/configuration info and ISO requests, and handle group functions. | `NewClient`, `WithName`, `WithProductInfo`, `WithConfigInfo` |
| Schedule transmissions | Broadcast PGNs periodically and let other devices retime or pause them through group functions. | `Client.Broadcast` |
| Request data | Send typed ISO requests and await typed replies. | `Request[T]` |
| Discover devices | Track observed devices by stable 64-bit NAME and current source address. | `Client.Devices`, `Client.DeviceAt` |
| Read many sources | Use SocketCAN, USB-CAN serial adapters, TCP/UDP gateways, candump logs, or in-memory frames. | `CAN`, `USB`, `TCP`, `UDP`, `File`, `Replay` |
| Gateway formats | Speak Yacht Devices RAW and Actisense-format gateway streams. | `FormatYDRaw`, `FormatActisense` |
| Physical values | Keep raw wire ticks while exposing generated SI-unit accessors. | `<Field>Value`, `Set<Field>Value`, `pgn.PhysicalValue` |
| Filter traffic | Filter by PGN metadata or decoded fields with CEL; metadata-only filters avoid decode work. | `Filter` |
| Preserve unknowns | Drop undecoded PGNs by default or surface them for logging/research. | `IncludeUnknown`, `*pgn.UnknownPGN` |
| Test without hardware | Replay bundled or custom captures and inspect written frames in tests. | `File`, `OriginalTiming`, `Replay`, `WrittenFrames` |

## Quick Start — No Boat Required

The repo bundles a real six-second capture from a sailing vessel
(`testdata/sample.log`), so your first run works at a desk:

```bash
git clone https://github.com/open-ships/n2k && cd n2k
go run ./cmd/n2k sniff -file testdata/sample.log | jq .
```

Or in code — copy, paste, it works:

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
gateway in RAW mode is a full bus citizen (address claiming included) from
the couch.

## The `n2k` CLI

Prebuilt binaries for Linux, macOS, and Windows are on the
[releases page](https://github.com/open-ships/n2k/releases), or:

```bash
go install github.com/open-ships/n2k/cmd/n2k@latest
```

```bash
# Yacht Devices WiFi gateway (RAW server mode) -- decoded JSON in one command
n2k sniff -tcp 192.168.4.1:1457

# SocketCAN (Linux), USB-CAN serial, UDP, or capture replay
n2k sniff -i can0
n2k sniff -u /dev/ttyUSB0
n2k sniff -udp :1457
n2k sniff -file capture.log            # add -timing to replay at real speed

# CEL filtering, unknown PGNs, jq-friendly output
n2k sniff -i can0 -f 'pgn == 127250' -unknown | jq .
```

The CLI currently ships `sniff` and `version`; the Go API provides the full
read/write client, replay, request/response, registry, and broadcast features.

## Why n2k

- **~600 PGN message types decoded** into typed Go structs (348 PGN numbers,
  including manufacturer-proprietary variants), generated from the
  community-maintained schema of the [canboat](https://github.com/canboat/canboat)
  project — the reference database for open NMEA 2000 decoding.
- **Physical units on top, raw ticks underneath.** Every numeric field with a
  physical interpretation gets generated accessors
  (`heading.HeadingValue()` → radians, `battery.VoltageValue()` → volts)
  while the struct keeps raw wire ticks for fidelity.
- **Byte-perfect re-encode.** Decode → re-encode round trips preserve the
  original payload bytes, verified against real captures.
- **A real bus citizen.** `NewClient` claims an address per ISO 11783,
  heartbeats, answers product/configuration info and ISO requests, and
  handles NMEA group functions (transmit/retime/pause) — the protocol
  behavior expected of a certified device, out of the box. No other Go
  library does this.
- **CEL message filtering** with an optimizer: metadata-only expressions
  skip decoding entirely.
- **Pure Go, CGO-free**, cross-compiles to Linux, macOS, and Windows.

### Library Comparison

As of July 2026, these are relevant projects that document NMEA 2000 or N2K
support. Go projects come first because `n2k` is for Go applications; non-Go
projects follow to clarify where this library sits in the broader ecosystem.
NMEA 0183-only parsers are intentionally excluded.

| Project | Language | Primary role | Decode model | Write / bus behavior | Source and gateway support | Best fit |
|---------|----------|--------------|--------------|----------------------|----------------------------|----------|
| [open-ships/n2k](https://github.com/open-ships/n2k) | Go | Standalone Go library + CLI for NMEA 2000 applications and bus-facing services. | ~600 generated typed PGN structs, runtime metadata, raw ticks, physical accessors, byte-preserving round trips. | `Client` claims an address, heartbeats, answers product/configuration and ISO requests, handles group-function transmit/retime/pause, writes PGN structs, and tracks devices. | SocketCAN, USB-CAN serial, TCP/UDP gateways, candump logs, in-memory replay; Yacht Devices RAW and Actisense formats. | Production Go apps that need the direct NMEA 2000 library and bus node in one package. |
| [boatkit-io/n2k](https://github.com/boatkit-io/n2k) | Go | Go library split across endpoint, service, generated PGN, and node packages. | Generated public PGN types in `pkg/pgn`. | `N2kService.Write` and node write flows after address claim; node package documents heartbeats and observed-device tracking. | SocketCAN, USB CAN, raw replay, N2K file endpoints. | Go apps that want its endpoint/service/node architecture. |
| [aldas/go-nmea-client](https://github.com/aldas/go-nmea-client) | Go | Work-in-progress Go reader plus `n2k-reader` CLI for raw capture/format workflows. | Field-oriented decode from a PGN database; no generated typed-PGN application API documented. | CLI can send STDIN to a CAN interface/device; Actisense devices expose raw-message writes. | Files, TCP connections, serial devices, SocketCAN, and multiple Actisense/raw formats. | Raw ingestion, format conversion, and exploratory tooling around Actisense/SocketCAN inputs. |
| [canboat/canboat](https://github.com/canboat/canboat) | C | Command-line analyzer and reference utility suite for exploring NMEA 2000 and NMEA 0183 networks. | Broad PGN database and analyzer output rather than a Go typed-struct API. | Can read and write N2K messages, but is analysis tooling, not a Go application library. | POSIX command-line workflows and CAN adapters supported by its tools. | Shell pipelines, reverse engineering, schema reference, and protocol research. |
| [ttlappalainen/NMEA2000](https://github.com/ttlappalainen/NMEA2000) | C++ | Object-oriented library for building NMEA 2000 devices, originally for Arduino-class boards. | Hand-authored C++ message helpers for common PGNs. | Handles mandatory NMEA 2000 device behavior and has been used in commercial certified devices. | Arduino, Teensy, ESP, MBED, Raspberry Pi-class targets via hardware-specific CAN drivers. | Embedded C++ devices and microcontroller firmware. |
| [canboat/canboatjs](https://github.com/canboat/canboatjs) | TypeScript | TypeScript library and CLI tools for parsing, encoding, and interfacing with NMEA 2000 networks. | PGN database-backed parsing with TypeScript type definitions and PGN output. | Bidirectional decode/encode/transmit support is documented. | Multiple formats and devices including Actisense, iKonvert, and YDWG-style gateways. | Node.js, Signal K, TypeScript tooling, and web-adjacent NMEA 2000 integrations. |
| [tomer-w/nmea2000](https://github.com/tomer-w/nmea2000) | Python | Python encoder/decoder and gateway clients, also used as a Home Assistant integration backend. | Generated Python decode/encode based on the PGN database; JSON-friendly field output. | Supports encoding/decoding frames and reading/writing through TCP and USB gateway clients. | EByte TCP, ASCII TCP, Actisense BST, Waveshare USB-CAN, and python-can devices. | Python automation, Home Assistant, and scripting around NMEA 2000 data. |
| [fard-draf/korri-n2k](https://github.com/fard-draf/korri-n2k) | Rust | `no_std`, no-allocation NMEA 2000 / ISO 11783 protocol stack for embedded Rust targets. | Build-time generated typed Rust structs for a selected PGN manifest or `full-pgns`. | Bidirectional send/receive with address claiming and fast-packet handling; group functions are still a documented roadmap gap. | Hardware-agnostic traits, bare-metal async runtimes, and Linux/tokio mode; hardware examples live separately. | Embedded Rust devices that need a small selectable PGN set. |

If you're building a Go application for an NMEA 2000 system — telemetry,
logging, gateways, test rigs, monitoring, automation, or bus-facing services —
use `n2k` directly.

### Sources and platforms

| Source | Linux | macOS | Windows | Write access |
|--------|:-----:|:-----:|:-------:|:------------:|
| `CAN` (SocketCAN) | ✅ | — | — | ✅ |
| `USB` (serial CAN adapter) | ✅ | ✅ | ✅ | ✅ |
| `TCP` (Yacht Devices RAW) | ✅ | ✅ | ✅ | ✅ full frame-level control |
| `TCP` (Actisense format) | ✅ | ✅ | ✅ | ✅ gateway stamps its own source address |
| `UDP` (both formats) | ✅ | ✅ | ✅ | read-only |
| `File` (candump `-L`/`-l`) / `Replay` | ✅ | ✅ | ✅ | read-only |

## Installation

```bash
go get github.com/open-ships/n2k
```

Releases follow semver with a `v0` major: `v0.x.y`. Every green build on
`main` automatically cuts a patch release (with prebuilt CLI binaries);
minor bumps are tagged manually when the API moves. While the major version
is 0, minor releases may contain breaking API changes — pin accordingly.

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
heading := &pgn.VesselHeading{}
heading.SetHeadingValue(1.5708) // radians; stored as raw wire ticks
result := client.Write(heading)
if err := result.Wait(); err != nil {
    log.Printf("write failed: %v", err)
}

// Explicitly set priority and destination
heading2 := &pgn.VesselHeading{
    Info: pgn.MessageInfo{Priority: pgn.Priority(2), TargetId: pgn.Target(42)},
}
heading2.SetHeadingValue(1.5708)
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
    fmt.Printf("Msg: %v\n", s.Message())
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

### Filter Messages using Common Expression Language

Filter messages using [CEL](https://github.com/google/cel-go) expressions.

`n2k` automatically optimizes filters for max performance -- metadata-only expressions skip decoding entirely.

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
| `source` | `int` | Source address (0-252) |
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
| `n2k.WithName(name)` | ISO 11783 device NAME for address claiming |
| `n2k.WithProductInfo(p)` | Product identity reported via PGN 126996 |
| `n2k.WithConfigInfo(ci)` | Installation description reported via PGN 126998 |
| `n2k.WithHeartbeatInterval(d)` | Heartbeat (PGN 126993) cadence; default 60s, 0 disables |
| `n2k.WithBus(bus)` | Inject a pre-constructed `n2k.Bus` (custom transport or test fake) instead of CAN/USB sources |

### A Complete Bus Device

Beyond claiming an address, a bus client behaves like a certified NMEA 2000
device out of the box:

- **Heartbeat (PGN 126993)** — sent every 60 seconds automatically (tune or
  disable with `WithHeartbeatInterval`).
- **Product & configuration info (PGNs 126996/126998)** — requests from other
  devices (chartplotters, analyzers) are answered automatically. Set your
  identity with `WithProductInfo` / `WithConfigInfo`; without them a generic
  software-gateway identity is reported so the device never shows up blank.
- **ISO requests (PGN 59904)** — requests for supported PGNs are answered;
  requests addressed to us for anything else are refused with an ISO
  acknowledgement NAK, per ISO 11783-3.
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
every tick (return nil to skip a tick):

```go
stop := client.Broadcast(time.Second, func() pgn.Message {
    heading := &pgn.VesselHeading{}
    heading.SetHeadingValue(currentHeadingRadians())
    return heading
})
defer stop()
```

Other devices can retime or pause the broadcast with a request group function
naming its PGN.

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

Numeric struct fields hold **raw wire ticks** (`*uint64`/`*int64`), which is
what makes byte-perfect re-encoding possible. Every numeric field with a
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

## Known Limitations

- Cross-field validation is not yet implemented.
- One physical bus per client.
- Address claiming uses a 1500ms default timeout; on heavily contested buses, increase via `WithClaimTimeout`.
- `File` and `UDP` sources are read-only; writing requires `CAN`, `USB`, or `TCP`.
- Over Actisense-format TCP connections the gateway stamps its own source
  address on transmissions (the protocol carries none), so the client's
  claimed address is not authoritative on the wire.
- Gateway TCP connections do not auto-reconnect; a dropped connection ends
  the client's read loop.

## License

MIT -- see LICENSE.

## Acknowledgments

- [canboat](https://github.com/canboat/canboat) — maintains the open PGN
  schema this library's decoders are generated from. The decoded-message
  breadth here is theirs.
- [boatkit-io/n2k](https://github.com/boatkit-io/n2k/) — the Go project that
  inspired this one.
