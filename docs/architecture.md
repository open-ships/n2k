# Architecture

`n2k` is split into deep modules with narrow seams. The public package owns
lifecycle and policy; internal packages own wire protocols.

```text
CAN / USB / TCP / UDP / serial / capture
              |
              v
       Bus or read-only source
              |
              v
      owned Observation Adapter
              |
      +-------+------------------+----------------+
      |                          |                |
      v                          v                v
system protocol path       user read pipeline   raw observation stream
claim / ISO request        metadata filter      frame / message / error
ISO transport              fast-packet assembly adapter / network / time
group functions            typed PGN decode     direction / owned bytes
registry / correlator      field filter
      |                          |                |
      +-----------+--------------+----------------+
                  v
          bounded subscriptions
```

## Runtime ownership

One `Client` owns one writable `Bus`, one address claimer, one ISO transport
manager, one system router, an application job worker, a protocol job worker,
and one physical wire writer. Separate bounded application, required, and
advisory queues preserve admission and priority. It also owns fixed discovery
and rejoin workers, up to 64 scheduled providers, and live typed and observation
subscriptions. `Receive`, `Scanner`, and `Observations` subscriptions are
independent.
When one subscriber exceeds its configured buffer, only that subscriber ends
with `ErrReceiveOverflow`.
Each subscriber receives its own deep message clone, including pointer fields,
repeating fields, diagnostics, and retained wire bookkeeping. Registry ingress
and snapshots use the same generated clone implementation.

The observation Module is the ownership Seam between transport Adapters and
consumers. Its Interface preserves `AdapterID`, `NetworkID`, source time, host
receipt time, relative gateway time, and direction. Payloads and frames are
copied before publication, keeping failure Locality at the subscriber that
falls behind (`ErrObservationOverflow`). This Depth supports faithful capture,
multi-network diagnostics, and bridging without coupling the codec to a
specific transport Implementation.

Actisense transports first cross a bounded BDTP/BST Module. BST-95 becomes a
source-authoritative Frame; BST-93 and Type-A ASCII remain gateway-owned
assembled Messages; BST-D0 becomes an assembled Message directly. BEM records
stay in a control path and surface as typed diagnostics as well as gateway or
transport-error Observations. The session is the sole reader and serialized
writer for one connection epoch, so response correlation, mode readiness,
reset handling, exact wire tracing, and disconnect cancellation retain
Locality.

The public `ActisenseTCP` and `ActisenseSerial` Adapters apply role policy at
the existing source/Bus seam. Their `source.run` Interface passively receives
gateway-owned messages; their `busBacked.newBus` Interface, reached by
`NewClient`, requires source-authoritative raw CAN. Explicit gateway-owned
formats have no Bus implementation and are rejected by `NewClient`.

`ActisenseGatewaySession` is the separate Interface for BST-93/94 transmission.
It exposes observations, diagnostics, status, raw/typed PGN sends, and the full
typed local BEM command Module without pretending that the caller controls the
gateway's CAN source. The same command Module sits behind
`ActisenseRemoteDevice`; only its addressed PGN-126720 envelope Adapter and
correlation origin differ.

The system path receives frames first. It handles protocol traffic regardless
of a user's CEL filter or whether the application is currently reading. This
ordering prevents user code from breaking address defense, request correlation,
transport reassembly, or device discovery.

Commanded Address (PGN 65240) crosses a narrow transport-to-claiming seam. The
transport manager must first complete a broadcast ISO transport transfer. The
client retains the transfer metadata and exact payload long enough to require
nine bytes, compare the complete 64-bit NAME, and reject addressed or
fast-packet lookalikes. The claiming module owns the actual state transition,
contention-history reset, notification, and immediate Address Claim response.
Application writes wait through a fresh contention window after the move.

## Write selection

`Client.Write` encodes and owns a snapshot before returning. Bounded FIFO
admission captures both connection and claim epochs. `WriteContext` also binds
the job to caller cancellation. Framing chooses:

| Payload and PGN | Wire mechanism |
|---|---|
| Non-fast, at most 8 bytes | One CAN frame |
| Fast PGN, at most 223 bytes | NMEA 2000 fast packet |
| Larger broadcast, at most 1785 bytes | ISO TP BAM |
| Larger addressed, at most 1785 bytes | ISO TP RTS/CTS |

Message-oriented custom gateways may implement `MessageWriter` only when their
Bus contract can still honor the Client's source identity. Buses that open
asynchronously may implement `ReadyBus`; the client will not begin address
claiming until `Ready` closes. Actisense BST-93/94 is deliberately outside that
seam because its wire source is owned by the gateway. Actisense BST-95 and CAN
ASCII mode 6 are true `Bus` Adapters.

Automatic protocol writes enter a dedicated protocol-transmission Module.
Required traffic (heartbeat, claims, ISO and group-function responses) has a
bounded high-priority lane; enumeration and information probes have a bounded
advisory lane with admission retry. Application writes cannot starve either
lane while waiting for ISO pacing, CTS, or acknowledgement. The sole physical
writer selects protocol records between application frames; a framed gateway
record itself is indivisible. Every physical write has a deadline (default one
second, configurable with `WithWriteTimeout`). Queue saturation, encoding,
transport failure, retry policy, and metrics now share one implementation and
test seam.

`WriteError.CompletedRecords` counts accepted physical records (frames, or whole
messages for a custom MessageWriter) and retains transmission uncertainty;
neither success nor a count asserts that a remote application accepted the data.
`WaitContext` cancels only waiting. Use `WriteContext` to cancel the operation.
An interrupted physical record can leave a partial transfer, so the client does
not resend it. ISO sessions additionally have an absolute 30-second deadline.

Unix serial transports share `internal/serialio`, which combines library line
configuration with an owned pollable writer. Linux CAN sockets are opened in
nonblocking mode and owned by Go's poller. Custom legacy buses must make Close
interrupt blocked I/O; context-aware buses can expose `ContextBus` and
`ContextMessageWriter` without closing an unrelated connection epoch.

## Connection epochs

Reconnect-capable buses implement `ConnectionLifecycleBus`. A connection is
not published to writers until Client installs a new readiness gate. After a
drop, stale registry topology is cleared. On the next epoch the client sends a
fresh Address Claim, waits through contention, then reopens writes, enumerates
the bus, and restarts the heartbeat. The network-session lifecycle Module
therefore owns both TCP connectivity and NMEA network citizenship instead of
leaving distributed reconnect checks at call sites.

`network_session.go` owns operation invalidation and bounded background work;
the Client serializes identity transitions with system-message delivery. Each
address change advances the claim epoch. Disconnect/claim changes immediately
cancel pending requests and queued/in-progress writes, clear partial fast
packets and ISO sessions, and reject queued system messages from older epochs.
Requests are capped at 64 and match source, destination, and both epochs.
Application admission fails with `ErrNotReady` until the current gate opens.
Historical user observations retain their epochs for diagnostics.

The stale-epoch gate belongs to a live Client, not to every read pipeline.
Standalone scanners and replay may combine independent networks or historical
connection epochs in any order. Their bounded fast-packet table includes the
network and both epochs in each key, so no network is suppressed by another's
larger epoch and partial packets cannot cross an epoch boundary.

For Actisense, an epoch is published only after the gateway acknowledges the
requested operating mode. Raw mode fails closed when BEM is rejected or
stripped. Volatile mode changes are best-effort restored on clean close. PGN
list changes are explicit batched transactions with one activation and a
same-epoch restore; sends never change the list. EEPROM/flash commits and
reinitialization exist only as explicit caller methods. The role-aware public
constructors do not fall back from raw CAN to a gateway-owned message session.

## Failure model

Bus termination, address-claim failure, and receive overflow are data, not log
side effects. They propagate through iterators/scanners and write results.
`Client.Status` supplies a polling seam for health checks and metrics. Raw
Actisense Clients and gateway sessions add cumulative transport byte/call,
BDTP failure, per-BST, BEM correlation/error/timeout, latency, reconnect, and
gateway-reset metrics. `slog` supplies structured diagnostics. `Close` cancels
owned work, closes the bus to release blocked I/O, stops scheduled protocol
work, fails queued writes, and waits for owned goroutines.

## Codec fidelity

Decode stores both the original payload and a canonical encoding of decoded
fields. Encode compares current fields with the canonical bytes. If unchanged,
it returns a copy of the original payload, retaining reserved bits and trailing
bytes. If changed, it encodes the current fields. See
[ADR 0002](adr/0002-wire-fidelity.md).

The schema-semantic codec Module also owns field conditions, signed
representability, physical ranges, declared sentinel values, and malformed
versus truncated errors. Unknown condition expressions fail during plan
compilation instead of silently shifting subsequent fields. Centralizing these
rules across all generated PGNs is the primary codec Leverage point.
Decode uses a fresh value and commits only on success; malformed input cannot
leave partially updated fields or stale variable-length content. Raw
canonicalization preserves unchanged invalid measurements without describing
them as valid physical values. Accessors use float64 scale/offset metadata and
report unavailable for sentinel or out-of-range ticks.

## Conformance lab

[`conformance/requirements.json`](../conformance/requirements.json) is the
machine-readable local evidence index. `just conformance-local` runs the
claiming, heartbeat, ISO request, group-function, fast-packet, BAM, RTS/CTS,
hostile-input, saturation, reconnect, timing, codec, and observation Seams.
Licensed review, hardware-in-the-loop evidence, and formal NMEA certification
remain external release activities described in
[`docs/conformance.md`](conformance.md).

Actisense additionally uses the SDK-pinned independent corpus in
[`conformance/actisense-golden.json`](../conformance/actisense-golden.json).
`just actisense-hardware <config>` runs the opt-in NGT/NGX matrix; a skipped
hardware test is never represented as a pass.

## Generated boundary

Category files under `pgn/`, `pgn/dispatch.go`, metadata, and the manifest are
generated from the checked-in schema snapshot. Change generator inputs or the
generator, then run `just pgn-sync`; never repair generated output by hand.
Generation verifies the snapshot's pinned checksum and needs no network.
All lookup families produce enums, with type-prefixed constants. Manifest
provenance, `DecodeComplete`, `EncodeComplete`, `CodecLimitations`, and
`HardwareVerified` keep generated coverage separate from verified support.
Runtime `MessageInfo.DecodeIssues` identifies unresolved conditional widths.
See [ADR 0005](adr/0005-reliability-boundaries.md) for the v1 contract changes.
