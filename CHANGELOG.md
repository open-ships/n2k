## Change Log for open-ships/n2k

### Unreleased — candump round-trip tool

- `cmd/roundtrip` — replays a candump `-L` log through fast-packet assembly, decodes each message, re-encodes it, and compares against the wire bytes. Prints a before/after hex view per message (with markers under differing bytes) and a per-PGN summary; `-report <file>` collects every issue plus the summary into a file. Always processes the whole log and exits 0. Reads a log file or stdin (`candump -L vcan0 | roundtrip -`).

### Unreleased — Architecture deepening

Seven refactors from an architecture review: one read pipeline, one CAN-ID codec, a public bus seam, a testable transport state machine, and a metadata-driven PGN codec with round-trip coverage for every PGN.

**Breaking: public API**

- `WithBus` now takes the new exported `n2k.Bus` interface (was `internal/canbus.Interface`, which external callers could not implement). `Bus.Run(ctx, handler)` receives the frame handler; implement `Bus` to inject fake hardware in tests.
- Repeating field sets are now slice fields: PGN structs with repeating groups carry `Repeating1 []<Type>Repeating1` (and `Repeating2`) instead of flattened single-element scalar fields. Struct shapes, JSON output (`"repeating1": [...]`), and CEL filter expressions referencing the old inline names all change.
- `NewClient` now rejects an out-of-range `WithName` device NAME (`DeviceName.Validate` is enforced; fields were previously silently bit-masked) and surfaces CEL filter compile errors eagerly for replay clients too.

**Breaking: decode/encode semantics**

- Numeric fields (≥ 4 bits, non-lookup) decode NMEA 2000 null sentinels to `nil` instead of a non-nil pointer holding the sentinel value.
- Variable-length strings (STRING_LAU) encode as length = len+2 with no trailing NUL (was len+3 plus NUL); decode strips trailing 0x00/0xFF filler. Strings now round-trip.
- Variable-length binary fields decode their real payload driven by the referenced length field (previously decoded as empty). `BitLengthField` references are widths in bits; dynamic length fields are widths in bytes.
- Fixed-width binary fields are bit-exact in both directions (no byte-alignment padding); non-byte-aligned fields such as PGN 130061's selection mask change wire bytes and no longer misalign subsequent fields.
- On encode, repeating-group count fields and variable-length length fields derive from the actual data, overriding user-set values.
- Fast-packet PGNs with payloads ≤ 8 bytes (e.g. empty repeating groups) are now correctly fast-framed on write; previously they were sent as raw single frames the reader could not parse.
- Transport-protocol (ISO-TP) reassembled messages now pass through the pre-decode filter like every other message.

**New**

- `n2k.Bus` — public bus seam; anyone can fake a bus for tests.
- `pgn.PhysicalValue(msg, fieldOrder)` — unit-scaled value of a decoded field (resolution + offset + unit label), documenting the raw-ticks contract of decoded fields.
- Round-trip test harness over all 599 PGN structs: 599/599 null-encoding fixpoint, 591/599 value round-trip (8 documented skips for fields whose width lives outside their own PGN's metadata).
- `pgn.UnknownPGN` implements `MessageInfo()`/`SetMessageInfo()`.

**Internal**

- ~80k lines of generated per-struct encode/decode bodies replaced by one metadata-driven interpreter (`pgn/codec.go`); generated files shrank 78%.
- One `readPipeline` shared by `Scanner`, `Receive`, and `Client` (previously two duplicated decode paths).
- CAN-ID and fast-packet header codecs consolidated in `internal/framer` (four scattered implementations removed).
- Transport manager: injectable timers/sleep for deterministic tests; completion callbacks and acks moved outside the manager lock.
- Manifest schema v4: examples now include payload fields and nested repeating groups (the previous walker silently emitted metadata-only entries); upstream schema refreshed 2.4.0 → 2.5.0 (599 PGN definitions, up from 579).
- Opt-in (`UPSTREAM_PARITY=1`) drift check for `enums.go` lookup enumerations; provenance of `enums.go` documented.

---

### 2026-06-29 - Unified PGN structs

- Replaced multiple PGN implementation paths with one `Pgn*` struct per PGN shape.
- Added the `pgn.PGN` interface: `PGNNumber`, `MessageInfo`, `SetMessageInfo`, `DecodePayload`, and `EncodePayload`.
- Decode dispatch now returns `pgn.PGN` values directly through the PGN switch.
- Removed the function-pointer registry and moved encode/decode behavior onto the structs themselves.
- PGN files are grouped by category using filenames such as `navigation.go` and `system.go`.

---

### 2026-05-01 — Message interface and type-safe API

Introduces the `pgn.Message` interface and eliminates redundant PGN specification when constructing messages. Users no longer need to pass the PGN number in `MessageInfo` — each struct knows its own PGN via its `PGNNumber()` method.

**New: `pgn.Message` interface**

- Single-method interface: `PGNNumber() uint32`
- All 278 PGN structs (including `UnknownPGN`) implement `Message`
- `Client.Write`, `Client.Receive`, `Receive`, and `Scanner` all use `Message` instead of `any`
- Encoder and decoder function signatures updated from `any` to `Message` throughout

**Changed: `MessageInfo` fields**

- `Priority` changed from `uint8` to `*uint8` — nil defaults to 6 (standard priority) on write
- `TargetId` changed from `uint8` to `*uint8` — nil defaults to 255 (broadcast) on write
- `SourceId` remains `uint8`; `Timestamp` remains `time.Time`

**New: convenience helpers**

- `pgn.Priority(v uint8) *uint8` — construct priority for `MessageInfo`
- `pgn.Target(v uint8) *uint8` — construct target address for `MessageInfo`

**Simplified write API**

Before:
```go
heading := &pgn.VesselHeading{
    Info:    pgn.MessageInfo{PGN: 127250, Priority: 2},
    Heading: ptrUint64(15708),
}
```

After:
```go
heading := &pgn.VesselHeading{
    Heading: ptrUint64(15708),
}
```

With explicit priority/target:
```go
heading := &pgn.VesselHeading{
    Info:    pgn.MessageInfo{Priority: pgn.Priority(2), TargetId: pgn.Target(42)},
    Heading: ptrUint64(15708),
}
```

**Cleanup**

- Removed stale `pgngen` / code generator references from comments
- Updated README with new API examples and `Message` interface docs

---

### 2026-04-25 — Write capabilities, encoders, and transport protocol

Adds the `Client` API for bidirectional NMEA 2000 communication, PGN encoders, and multi-frame transport protocol support.

**New: `Client` API (`client.go`)**

- `NewClient(ctx, opts...)` — creates a client for reading and writing PGN messages
- `Client.Write(msg)` — asynchronous, FIFO-ordered message transmission
- `Client.Receive()` / `Client.Scanner()` — delegates to the top-level read APIs
- `Client.WrittenFrames()` — inspect captured frames in replay/test mode
- `Client.Close()` — graceful shutdown with write drain
- `WriteResult` with `Wait()` / `Err()` for write completion tracking

**New: PGN encoders**

- `pgn/pgndatastream_writer.go` — `PGNDataStreamWriter` for encoding fields back to wire format
- `EncodePayload` methods on PGN structs

**New: frame encoding (`internal/framer/`)**

- `BuildCANID` — construct 29-bit CAN identifiers from PGN, priority, source, destination
- `FrameSingle` — single-frame (<=8 byte) CAN framing
- `FrameFastPacket` — multi-frame fast-packet framing with sequence IDs

**New: transport protocol (`internal/transport/`)**

- BAM (Broadcast Announce Message) for multi-frame broadcasts >223 bytes
- RTS/CTS (Request to Send / Clear to Send) for addressed multi-frame messages
- `Manager` coordinates send/receive state machines with timeouts

**New: address claiming (`internal/claiming/`)**

- ISO 11783 address claim protocol implementation
- `DeviceName` builder with industry code, manufacturer, device class, etc.

**New: PGN struct files**

- Reorganized PGN structs into category files.
- `pgn/enums.go` — NMEA 2000 enumeration constants

**New: options**

- `WithSourceAddress(addr)` — explicit source address for writes
- `WithName(name)` — ISO 11783 device NAME for address claiming

---

### 2026-04-15 — CI hardening and cleanup

- Added GitHub Actions security scanning workflow (`security.yaml`)
- Added CI badges to README (Test, Lint, Secure)
- Bumped Go version and tooling
- Unexported CAN bus channel types (`USBCANChannel` → `usbCANChannel`, `SocketCANChannel` → `socketCANChannel`)
- Removed SocketCAN auto-setup from CI (not available in GitHub Actions runners)
- Fixed `reflect.Ptr` → `reflect.Pointer` deprecation in `scanner.go`
- Cleaned up lint issues across the codebase

---

### 2026-03-31

Initial setup
