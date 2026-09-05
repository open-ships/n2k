## Change Log for open-ships/n2k

### v1.3.0 — 2026-09-05 — bounded, cancellable, epoch-safe reliability

- Made application writes immutable snapshots at admission, with `WriteContext`,
  bounded physical write deadlines, and `WriteError` reporting partial-transfer
  uncertainty. Protocol frames can proceed during application ISO transfers.
- Bound requests, queued writes, reassembly, discovery, and reconnect work to
  connection and claim epochs. Close and terminal failures invalidate pending
  operations immediately; old work is never replayed after reconnect.
- Added owned deep message cloning for subscribers and registry snapshots, a
  bounded replay capture with loss counters, and readiness/queue status.
- Fixed ISO transport timer races, echoed announcement collisions, indefinite
  CTS holds, and shutdown/reset cancellation of BAM and RTS/CTS transfers.
- Made Actisense startup, writes, and restoration cancellable on their exact
  connection. Shared pollable serial I/O and Linux SocketCAN descriptors make
  physical shutdown interruptible, including serial backpressure.
- Corrected string boundaries, physical precision and sentinel availability,
  transactional decode, specific-before-fallback dispatch, and unresolved field
  diagnostics. All typed variants now participate in value roundtrip tests.
- Made PGN and enum generation offline and checksum-pinned. Published separate
  decode/encode completeness, limitations, provenance, and hardware-verification
  metadata; typed coverage is not a completeness or certification claim.
- Added executable conformance evidence, scheduled long fuzzing/resource soaks,
  and independent reliability regressions. Missing/skipped required software
  evidence fails the gate; hardware and licensed certification stay separate.
- Intentionally refined the still-evolving v1 API: physical metadata uses
  `float64`, enum constants are type-prefixed, and broadcast providers accept
  context. `Broadcast`/`BroadcastPGN` return `(stop, error)`, cap active workers
  at 64, and require providers to honor cancellation. No v2 module path is used.

### v1.2.0 — 2026-08-20 — complete Actisense NMEA 2000 control plane

- Replaced gateway-owned Actisense `Client` writes with the honest
  `ActisenseGatewaySession` contract. BST-93/94 sends now run under the
  physical gateway identity, while `NewClient` accepts only source-authoritative
  BST-95 or CAN-ASCII representations. This is an intentional v1 contract
  refinement while the API is documented as in flux.
- Added one typed BEM command Module covering every compiled NMEA 2000-facing
  verb: product/CAN/port information, Rx/Tx state and multi-reply lists,
  lifecycle, echo/time, and explicit EEPROM, flash, default, and reinitialize
  operations. Sends never mutate Tx lists; explicit batched changes snapshot,
  activate once, roll back, and best-effort restore.
- Added remote Actisense device control through the exact variable-length
  addressed PGN-126720 envelope. Correlation binds remote source, local
  destination, response group, verb, connection epoch, and claim epoch;
  duplicate keys, pending tables, and response trains are bounded.
- Hardened the sole-reader session with response-group/origin correlation,
  scoped negative acknowledgements, immediate send failure, typed partial
  results, safe reply #257 overflow, complete Port Inventory accumulation, and
  NGX F2 completion that requires both standard slots and proprietary bitmaps.
- Added configurable serial settings, Actisense CAN ASCII mode 6 and NMEA 2000
  Type-A ASCII parsing/emission, a custom byte-stream session seam, EBL writing
  and live wire tracing, exact unframed replay evidence, typed diagnostics,
  and cumulative per-layer transport/BST/BEM/latency/reset metrics.
- Added an SDK-commit-pinned independent corpus with request and response bytes
  for all 23 solicited commands, remote/ASCII/EBL vectors, and an opt-in NGT/NGX
  hardware matrix runner. Known SDK defects remain deliberate compatibility
  deviations rather than copied behavior.

### v1.1.0 — 2026-08-20 — Actisense raw CAN, serial, and capture support

- Added a bounded Actisense protocol Module for BST-93/94/95/D0, Type-2
  lengths, acknowledged BEM mode and Tx-PGN management, typed F0/F1/F2/F4
  diagnostics, direct 115200 8N1 serial access and enumeration, and EBL replay.
  `FormatActisenseRaw` is a source-authoritative BST-95 `Bus`; the existing
  `FormatActisense` route remains the gateway-owned v1 message session.
- Added `ActisenseTCP` and `ActisenseSerial` constructors. Read-only callers use
  compatible Actisense receive behavior, while `NewClient` requires
  acknowledged mode 5 and never falls back to gateway-owned writes. Explicit
  stream formats remain available as compatibility escape hatches;
  `ActisenseModeError` makes negotiation failures available to `errors.As`.
- Actisense parser and gateway failures now surface as owned transport/gateway
  observations and status counters instead of disappearing silently. D0
  payloads enter the Pipeline as assembled Messages up to the NMEA transport
  maximum without fabricated CAN frames.
- Pinned PGN generation to the canboat v7.1.0 source revision represented by
  the stable public API; the separate upstream-parity check continues to track
  canboat master so schema upgrades are reviewed explicitly.
- Updated the pinned Go and CI toolchain to 1.26.6 for current standard-library
  security fixes.

### v1.0.2 — 2026-07-25 — Go toolchain update

- Updated the pinned Go and CI toolchain to 1.26.5 and refreshed security-tool
  versions while preserving CGO-free Linux, macOS, and Windows builds.

### v1.0.1 — 2026-07-25 — library-only distribution

- Removed the in-repository `n2k` command, its tests and demo, Cobra
  dependencies, build recipe, and GoReleaser packaging. Command-line tooling
  now lives in the separate `open-ships/n2k-cli` repository.
- Releases now publish the Go module's source tag without prebuilt command
  archives.

### v1.0.0 — 2026-07-18 — schema-safe codec, resilient network sessions, and raw observability

- The shared PGN codec now executes schema-declared conditional fields,
  including the proprietary header in PGN 126208, and validates signed-width
  representability, declared physical ranges, nulls, and special sentinels.
- Automatic bus-citizenship traffic has bounded required and advisory priority
  lanes. Heartbeat, address claim, ISO request, group-function, enumeration,
  and device-information failures can no longer disappear as discarded write
  results or log-only errors.
- TCP reconnect starts a new NMEA network epoch: writes are gated while the
  client reclaims its address, waits through contention, clears stale device
  topology, re-enumerates, and restarts scheduled traffic.
- Added owned raw observations with Adapter/network identity, source and host
  timestamps, gateway-relative time, direction, frames, assembled payloads,
  and decode errors. `Observe`, `Client.Observations`, and
  `ReplayObservations` use bounded, non-blocking ownership seams.
- The Cobra-based CLI now supports `record`, `replay`, `validate`, `devices`,
  and `pgn` in addition to typed `sniff` and version output. Nested help, typo
  suggestions, argument validation, and generated Bash, Zsh, Fish, and
  PowerShell completion cover source paths, formats, and known PGNs. Text
  output names concrete PGN types and renders generated physical values with
  SI units; JSON retains exact wire values. `just build` produces `bin/n2k`.
- `Client.Status` now exposes connection epochs, protocol lanes, typed/raw
  subscriber counts, write counters, and observation counters.
- Added a machine-readable requirement index and local evidence families for
  claiming, heartbeat, ISO requests, group functions, fast packet, BAM,
  RTS/CTS, malformed traffic, saturation, reconnect, timing, codec semantics,
  and raw observation. Licensed review, hardware-in-the-loop execution, and
  official NMEA certification remain external product-release activities.
- The module now declares a stable v1 baseline. Exported v1 compatibility
  follows semantic versioning; the release workflow publishes the `v1.0.0`
  tag after the required gates pass on the release commit.

### v0.3.0 — 2026-07-18 — bounded lifecycle, wire fidelity, and hostile-input hardening

- Commanded Address (PGN 65240) is now a first-class address-claim transition:
  only an exact nine-byte broadcast ISO transport transfer for this node's
  complete NAME and a claimable address is accepted, followed immediately by
  a new Address Claim and contention window. Adversarial tests cover malformed,
  mismatched, fast-packet, addressed, and special-address commands.
- Added `just conformance-local`, a current standards/evidence guide, and a
  machine-readable certification evidence template. These make local protocol
  regression evidence reproducible without misrepresenting it as the licensed
  NMEA certification run.
- Live reads now use independent bounded subscriptions; slow consumers receive
  `ErrReceiveOverflow` without stalling protocol handling. Writes use a bounded
  queue and report `ErrWriteQueueFull`.
- Runtime bus failures propagate to readers, writers, `Client.Err`, and
  `Client.Status`. Scanner, bus, scheduler, and write shutdown paths are
  cancellation-safe and race-tested.
- ISO transport and fast-packet assembly validate lengths, packet ranges,
  session ownership, timeouts, and resource bounds. Missing/out-of-order
  packets cannot be accepted as a complete transfer.
- Unmodified decoded messages preserve their original payload exactly; field
  changes encode current values. String sentinel and fixed-padding edge cases
  now round-trip consistently.
- USB-CAN startup configures 250 kbit/s extended-frame operation before the
  bus becomes writable, and stream parsers resynchronize after corrupt input.
- NMEA 2000 source-address claiming now uses the valid 0–251 range. New
  `WithPreferredAddress` lets applications restore a persisted prior address
  while retaining automatic contention handling.
- Added health/configuration hooks, adversarial and fuzz tests, cross-platform
  and race CI, pinned lint/security tools, architecture context, ADRs,
  contribution guidance, and private security-reporting instructions.

### v0.2.2 — 2026-07-18 — clean CLI shutdown

`n2k sniff` now handles termination signals through context cancellation and
exits successfully after its input pipeline shuts down.

### v0.2.1 — 2026-07-16 — panics in a misbehaving Bus no longer crash the process

A `Bus` implementation that panics (rather than returning an error) from `WriteFrame`, `Run`, or a
message handler previously crashed the whole host process, because n2k runs those on its own
goroutines where an escaping panic is unrecoverable by the caller. Every n2k-owned goroutine now
recovers panics, logs them with a stack trace, and degrades locally instead:

- The initial address claim runs on an n2k goroutine (since `WithReconnect`); a panic there now
  surfaces as a `NewClient` error (`n2k: panic during address claim: ...`) rather than a crash.
- The write loop recovers per job, completing that write with an error and continuing to serve —
  a panicking write no longer wedges every later `Write` blocked on the queue.
- The bus read loop recovers per frame, so a fault while handling one frame (including the defensive
  address claim the claimer writes inline in response to a contender or ISO request) is logged and
  skipped without tearing down the reader.
- The heartbeat and system-router dispatch goroutines are guarded the same way.

### v0.2.0 — 2026-07-15 — writable TCP gateway buses

`n2k.TCP(addr, format)` now works with `NewClient`, not just `Receive`/`NewScanner` — the full bus
stack over the boat's WiFi gateway:

- **Yacht Devices RAW** (`FormatYDRaw`) is frame-level in both directions: transmitted frames go
  out as `msgid b0..b7<CR><LF>` lines and the gateway echoes them back with direction T, so
  address claiming, heartbeats, and group functions behave exactly as on CAN hardware.
- **Actisense format** (`FormatActisense`) is message-level: sends use the binary transmit command
  (priority, PGN, destination, length, data in DLE/STX framing) after a gateway initialization
  command at connect. The gateway performs fast-packet fragmentation and stamps its own claimed
  source address (the protocol carries none on sends).
- New exported `MessageWriter` interface: a `Bus` that also implements it receives whole
  assembled payloads (≤223 bytes) from the client's write path instead of CAN frames; ISO-TP
  transfers and protocol frames still use `WriteFrame`. Custom `WithBus` transports can opt in.
- Writers issued before the connection is up (the address claimer fires immediately) block until
  the dial completes rather than failing.
- **Auto-reconnect** (`WithReconnect(ReconnectPolicy{InitialBackoff, MaxBackoff})`, defaulting to
  500ms/30s): a dropped TCP connection is re-dialed with exponential backoff instead of ending the
  read loop. TCP sources only; the initial connection must still succeed within the claim timeout,
  and a transparent transport reconnect does not re-run NMEA 2000 address claiming. Writers parked
  during an outage are released by reconnect or by `Close`, and a write is only retried on a fresh
  connection when nothing was sent, so a retry never duplicates a partial frame.
- Generator: `sourceDefinitions` is now assembled from per-chunk functions instead of one ~6k-line
  composite literal, fixing `go build -race ./pgn` ("NewBulk too big" internal compiler error) —
  `just test-race` works again.
- `UDP` and `File` sources remain read-only.

### v0.1.0 — 2026-07-14 — semver releases, installable CLI, typed physical accessors

**Typed physical-value accessors** — every numeric PGN struct field with a physical
interpretation (a unit, non-unity resolution, or offset) now has generated accessors, on
repeating-group element structs too: `<Field>Value() (float64, bool)` returns the physical value
in the field's schema unit, `Set<Field>Value(float64)` writes one (rounded to the nearest wire
tick). Raw-tick fields are unchanged underneath, so round-trip fidelity is preserved.
`pgn.PhysicalValue` remains for dynamic, metadata-driven access.

**Named CLI** — `cmd/sniffer.go` is replaced by `cmd/n2k`, a subcommand CLI
(`go install github.com/open-ships/n2k/cmd/n2k@latest`):

- `n2k sniff` decodes to JSON lines from `-i` (SocketCAN), `-u` (USB serial), `-file`
  (candump replay, `-timing` for real-speed), `-tcp`/`-udp` gateways (`-format raw|actisense`),
  with `-f` CEL filters and `-unknown`.
- `n2k version` reports the release version (set by goreleaser).

**Semver releases with prebuilt binaries** — releases are tagged `v0.x.y`, and goreleaser
publishes archives for Linux/macOS/Windows (amd64/arm64). Release names carry the CalVer date.
A commented Homebrew-tap config is in `.goreleaser.yaml`, pending an `open-ships/homebrew-tap`
repo and token.

**No-hardware quickstart** — `testdata/sample.log` is a real six-second capture (position, AIS,
and identity PGNs removed; see `testdata/README.md`), used by new runnable `Example*` tests and
the README quickstart, with an animated sniffer demo at `.github/demo.svg`.

### 2026.07.13-2 — bus citizenship, device registry, and new sources

**Bus citizenship** — bus clients now meet the protocol obligations of a real NMEA 2000 device:

- Heartbeat (PGN 126993) transmitted every 60s once the address claim completes; `WithHeartbeatInterval(d)` retunes it, `0` disables it.
- Product information (126996) and configuration information (126998) are answered when requested; set them with `WithProductInfo` / `WithConfigInfo` (sane software-gateway defaults otherwise).
- ISO requests (59904) now honor the destination address (previously claim requests addressed to other nodes were answered too). Addressed requests for unsupported PGNs are NAKed with an ISO acknowledgement (59392); broadcast requests for unsupported PGNs are ignored per ISO 11783-3.
- Group functions (126208): request group functions can transmit, retime, pause (interval 0), or restore-default (0xFFFFFFFE) the heartbeat and any scheduled broadcast; heartbeat cadence clamps to [1s, 60s]. Commands, read-fields, write-fields, and parameterized requests are refused with acknowledge group functions carrying the proper PGN error codes.
- `Request[T pgn.Message](ctx, client, target)` — typed ISO request/response: sends a request for T's PGN and awaits the matching reply (from `target`, or any device when target is 255) with a 1250ms default timeout.
- `client.Broadcast(interval, provide)` — periodic transmission scheduler; the provider is called per tick and may return nil to skip. Returns a stop function; group functions can retime it.
- All protocol reactions run on a dedicated system decode path, independent of user `Filter` expressions.

**Device registry**

- `client.Devices()` / `client.DeviceAt(source)` — a live, NAME-keyed map of the bus built from address claims, product info, and configuration info. Tracks address moves, last-seen times, and requests product/config info from newly seen devices once per NAME. At startup the client broadcasts an address-claim request to enumerate the bus.

**New sources** (read-only: `Receive`/`NewScanner`)

- `File(path, ...)` — candump `-L`/`-l` log files; `OriginalTiming()` paces replay by the log's timestamps.
- `TCP(addr, format)` / `UDP(listenAddr, format)` — network gateway streams: `FormatYDRaw` (Yacht Devices RAW ASCII) and `FormatActisense` (Actisense binary framing; assembled messages are re-framed through the normal decode pipeline).
- candump parsing moved from the roundtrip tool into `internal/candump` (now with timestamp extraction).

### 2026.07.12 — candump round-trip tool

- `cmd/roundtrip` — replays a candump `-L` log through fast-packet assembly, decodes each message, re-encodes it, and compares against the wire bytes. Prints a before/after hex view per message (with markers under differing bytes) and a per-PGN summary; `-report <file>` collects every issue plus the summary into a file. Always processes the whole log and exits 0. Reads a log file or stdin (`candump -L vcan0 | roundtrip -`).

### 2026.07.08 — Architecture deepening

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
