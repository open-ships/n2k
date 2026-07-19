# n2k Context

This is the short entry point for maintainers and coding agents. The README
explains the public product; [docs/architecture.md](docs/architecture.md)
explains the runtime design.

## Domain language

- **Frame**: one classical CAN 2.0 frame, with a 29-bit identifier and at most
  eight data bytes.
- **PGN**: an NMEA 2000 message type identified by its Parameter Group Number.
- **Message**: a typed `pgn.Message`; a message may occupy one frame, a fast
  packet, or an ISO transport transfer.
- **Source / destination**: dynamic 8-bit bus addresses. Addresses 0–251 are
  claimable; 252–253 are reserved, 254 means unable to claim, and 255 is
  broadcast/global.
- **NAME**: stable 64-bit ISO 11783 device identity used for address claiming.
- **Commanded Address**: PGN 65240, a nine-byte broadcast ISO transport
  transfer that assigns a claimable address to the node with an exact NAME.
- **Pipeline**: frame metadata, fast-packet assembly, decode, filtering, and
  delivery.
- **Observation**: an owned record at the frame, assembled-message, or decode-
  error layer carrying source/network identity, capture and receipt time, and
  receive/transmit direction.
- **System router**: protocol-only decode path for claiming, requests, group
  functions, correlation, and the device registry. User filters never affect it.
- **Protocol transmission**: the required and advisory priority lanes for
  automatic bus-citizenship traffic. Application writes cannot consume them.
- **Connection epoch**: one successful gateway connection and its NMEA network
  readiness cycle. A new epoch must reclaim before ordinary traffic resumes.
- **Bus**: the extension seam for physical or virtual read/write transports.
- **Adapter**: converts a transport representation into owned CAN frames.

## Non-negotiable invariants

1. Protocol handling runs before user delivery and cannot be stalled by an
   unread subscriber.
2. Every queue and reassembly table is bounded. Overflow is explicit; it is
   never silent backpressure or unbounded memory growth.
3. Each active transport receive session is unambiguous by source and
   destination because TP.DT packets do not carry the transported PGN.
4. Unmodified decoded PGNs re-encode to their original bytes. A field mutation
   invalidates that replay path and encodes the current fields.
5. Runtime bus failure reaches readers, writers, `Client.Err`, and
   `Client.Status`; a client never remains half-alive after `Bus.Run` exits.
6. Caller-owned frames, payloads, and registry snapshots are copied at
   ownership boundaries.
7. Generated PGN files and metadata change only through `just pgn-sync`.
8. Commanded Address changes state only after an exact nine-byte BAM transfer,
   an exact 64-bit NAME match, and a requested address in 0–251; the node then
   immediately reclaims that address before application writes resume.
9. Required protocol traffic has bounded priority admission; rejection,
   encoding failure, or transport failure is terminal and observable.
10. A reconnect creates a new connection epoch, clears stale topology,
    reclaims the address, waits through contention, then restarts enumeration
    and scheduled transmissions.
11. Schema-declared conditions, signed widths, ranges, and sentinels are codec
    semantics, not documentation-only metadata.

## Where to make changes

- Public lifecycle and writes: `client.go`, `status.go`, `writeresult.go`
- Read pipeline and fan-out: `pipeline.go`, `scanner.go`, `messagehub.go`,
  `observation.go`, `observationhub.go`, `raw/`
- Protocol transmission policy: `protocoltx.go`, `client.go`, protocol writers
- Connection epochs: `client.go`, `internal/gateway/tcpbus.go`, `registry.go`
- Hardware/network seams: `bus.go`, `source*.go`, `internal/canbus/`,
  `internal/gateway/`
- Fast-packet assembly: `internal/adapter/`
- ISO transport protocol: `internal/transport/`
- Wire codec and generated messages: `pgn/`
- Address claiming and Commanded Address: `internal/claiming/`, `client.go`
- Device discovery: `registry.go`
- CLI command tree, help, and completion: `cmd/n2k/`

## Required verification

Run `go test ./...`, `just pgn-sync-check`, `just lint`, and `just secure`
before handoff. Run `just test-race` for concurrency changes and `just
fuzz-smoke` for parser, codec, or framing changes. `just conformance-local`
produces the reproducible protocol evidence described in
[docs/conformance.md](docs/conformance.md); it is not a substitute for the
licensed NMEA certification run.
