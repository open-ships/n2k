# Architecture

`n2k` is split into deep modules with narrow seams. The public package owns
lifecycle and policy; internal packages own wire protocols.

```text
CAN / USB / TCP / UDP / capture
              |
              v
       Bus or read-only source
              |
              v
      frame ownership boundary
              |
      +-------+------------------+
      |                          |
      v                          v
system protocol path       user read pipeline
claim / ISO request        metadata filter
ISO transport              fast-packet assembly
group functions            typed PGN decode
registry / correlator      field filter
      |                          |
      +-----------+--------------+
                  v
          bounded subscriptions
```

## Runtime ownership

One `Client` owns one writable `Bus`, one address claimer, one ISO transport
manager, one system router, one serialized write queue, and zero or more live
read subscriptions. `Receive` and `Scanner` subscriptions are independent.
When one subscriber exceeds its configured buffer, only that subscriber ends
with `ErrReceiveOverflow`.

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

## Generated boundary

Category files under `pgn/`, `pgn/dispatch.go`, metadata, and the manifest are
generated from the checked-in schema snapshot. Change generator inputs or the
generator, then run `just pgn-sync`; never repair generated output by hand.
