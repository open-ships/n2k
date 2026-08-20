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
manager, one system router, one serialized application-write queue, two
protocol-transmission queues, and zero or more live typed and observation
subscriptions. `Receive`, `Scanner`, and `Observations` subscriptions are
independent.
When one subscriber exceeds its configured buffer, only that subscriber ends
with `ErrReceiveOverflow`.

The observation Module is the ownership Seam between transport Adapters and
consumers. Its Interface preserves `AdapterID`, `NetworkID`, source time, host
receipt time, relative gateway time, and direction. Payloads and frames are
copied before publication, keeping failure Locality at the subscriber that
falls behind (`ErrObservationOverflow`). This Depth supports faithful capture,
multi-network diagnostics, and bridging without coupling the codec to a
specific transport Implementation.

Actisense transports first cross a bounded BDTP/BST Module. BST-95 becomes a
source-authoritative Frame, BST-93 retains its v1 synthetic-frame Adapter, and
BST-D0 becomes an assembled Message directly. BEM records stay in a
gateway-session control path and surface as gateway or transport-error
Observations. The session is the sole reader and serialized writer for one
connection epoch, so response correlation, mode readiness, reset handling,
and disconnect cancellation retain Locality.

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

`Client.Write` admits work to a bounded FIFO queue. Encoding then chooses:

| Payload and PGN | Wire mechanism |
|---|---|
| Non-fast, at most 8 bytes | One CAN frame |
| Fast PGN, at most 223 bytes | NMEA 2000 fast packet |
| Larger broadcast, at most 1785 bytes | ISO TP BAM |
| Larger addressed, at most 1785 bytes | ISO TP RTS/CTS |

Message-oriented gateways may implement `MessageWriter` and accept assembled
payloads up to 223 bytes. Buses that open asynchronously may implement
`ReadyBus`; the client will not begin address claiming until `Ready` closes.
The legacy Actisense BST-93/94 route is such a gateway-owned message session;
its wire source address is not controlled by the client. Actisense BST-95 raw
mode is the corresponding true `Bus` Adapter.

Automatic protocol writes enter a dedicated protocol-transmission Module.
Required traffic (heartbeat, claims, ISO and group-function responses) has a
bounded high-priority lane; enumeration and information probes have a bounded
advisory lane with admission retry. Application writes cannot starve either
lane. This narrow Interface gives high Leverage: queue saturation, encoding,
transport failure, retry policy, and metrics have one Implementation and one
test Seam.

## Connection epochs

Reconnect-capable buses implement `ConnectionLifecycleBus`. A connection is
not published to writers until Client installs a new readiness gate. After a
drop, stale registry topology is cleared. On the next epoch the client sends a
fresh Address Claim, waits through contention, then reopens writes, enumerates
the bus, and restarts the heartbeat. The network-session lifecycle Module
therefore owns both TCP connectivity and NMEA network citizenship instead of
leaving distributed reconnect checks at call sites.

For Actisense, an epoch is published only after the gateway acknowledges the
requested operating mode. Raw mode fails closed when BEM is rejected or
stripped. Volatile mode changes are best-effort restored on clean close and
are never committed to EEPROM or flash.

## Failure model

Bus termination, address-claim failure, and receive overflow are data, not log
side effects. They propagate through iterators/scanners and write results.
`Client.Status` supplies a polling seam for health checks and metrics. `slog`
supplies structured diagnostics. `Close` cancels owned work, closes the bus to
release blocked I/O, stops scheduled protocol work, fails queued writes, and
waits for owned goroutines.

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

## Conformance lab

[`conformance/requirements.json`](../conformance/requirements.json) is the
machine-readable local evidence index. `just conformance-local` runs the
claiming, heartbeat, ISO request, group-function, fast-packet, BAM, RTS/CTS,
hostile-input, saturation, reconnect, timing, codec, and observation Seams.
Licensed review, hardware-in-the-loop evidence, and formal NMEA certification
remain external release activities described in
[`docs/conformance.md`](conformance.md).

## Generated boundary

Category files under `pgn/`, `pgn/dispatch.go`, metadata, and the manifest are
generated from the checked-in schema snapshot. Change generator inputs or the
generator, then run `just pgn-sync`; never repair generated output by hand.
