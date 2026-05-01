## Change Log for open-ships/n2k

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
    Heading: ptrFloat32(1.5708),
}
```

After:
```go
heading := &pgn.VesselHeading{
    Heading: ptrFloat32(1.5708),
}
```

With explicit priority/target:
```go
heading := &pgn.VesselHeading{
    Info:    pgn.MessageInfo{Priority: pgn.Priority(2), TargetId: pgn.Target(42)},
    Heading: ptrFloat32(1.5708),
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
- Encoder functions for all PGN structs that have decoders
- `EncoderLookup` map in `pgn/registry.go` for PGN-to-encoder dispatch
- `PgnInfo.Encoder` field on PGN metadata

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

- Reorganized PGN structs from single generated file (`pgninfo_generated.go`) into category files: `ais.go`, `communication.go`, `electrical.go`, `engine.go`, `entertainment.go`, `environmental.go`, `lighting.go`, `navigation.go`, `other.go`, `propulsion.go`, `sensors.go`, `system.go`
- `pgn/registry.go` — consolidated `pgnList`, `PgnInfoLookup`, `EncoderLookup`, and `init()`
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
