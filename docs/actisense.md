# Actisense NMEA 2000 support

The v2 API separates two identities that the wire protocol cannot combine:

- `Client` is a source-authoritative NMEA 2000 node. With Actisense hardware
  it uses BST-95 binary raw CAN or mode-6 CAN ASCII, both of which carry the
  complete 29-bit CAN identifier.
- `ActisenseGatewaySession` sends assembled BST-94 PGNs under the physical
  gateway's own source address. It owns local BEM control and never claims a
  virtual NMEA 2000 address.

`FormatActisense` and `FormatActisenseN2KASCII` remain passive read formats.
Passing either to `NewClient` returns
`ErrActisenseGatewaySessionRequired`; this prevents a gateway-owned source from
silently impersonating the Client's claimed identity.

## Gateway-owned session

```go
session, err := n2k.NewActisenseTCPSession(ctx, "10.0.0.5:2000",
    n2k.WithActisenseSessionMode(n2k.ActisenseModeTransferReceiveAll),
    n2k.WithActisenseCommandTimeout(5*time.Second),
)
if err != nil { ... }
defer session.Close()

product, err := session.GetProductInfo(ctx)
if err != nil { ... }
fmt.Println(product.Model, product.SoftwareVersion)

// The gateway supplies the on-wire source address. No Tx list is changed.
err = session.SendRawPGN(ctx, 126992, 3, 255, payload)
```

Use `NewActisenseSerialSession` for a direct serial connection or
`NewActisenseGatewaySession` for a custom full-duplex byte stream. Serial
settings are configurable and default to 115200 8N1:

```go
serial := n2k.ActisenseSerialConfig{
    BaudRate: 38400,
    DataBits: 8,
    Parity: n2k.ActisenseParityNone,
    StopBits: n2k.ActisenseStopBitsOne,
}
session, err := n2k.NewActisenseSerialSession(ctx, "/dev/ttyUSB0", serial)
```

`Observations` emits owned assembled messages and gateway/transport evidence.
`Diagnostics` emits typed startup, error, system, and negative-ack records.
`Status` reports connection epoch, operating/receive role, known model
capabilities, subscriber counts, trace failure, and cumulative metrics.

## Source-authoritative Client

```go
client, err := n2k.NewClient(ctx,
    n2k.ActisenseTCP("10.0.0.5:2000"),
    n2k.WithPreferredAddress(80),
)
```

`ActisenseTCP` and `ActisenseSerial` select acknowledged BST-95 mode for
`NewClient`. `TCP(..., FormatActisenseCANASCII)` and
`Serial(..., FormatActisenseCANASCII)` select acknowledged mode 6. A rejected,
timed-out, or stripped BEM setup fails startup; there is no fallback to
gateway-owned writes. `Client.Status().Actisense` contains the raw gateway's
transport, framing, BEM, reset, and reconnect counters.

## Local and remote commands

`ActisenseGatewaySession` and `ActisenseRemoteDevice` embed the same typed
`ActisenseDevice` command Interface. A remote handle uses the exact Actisense
manufacturer envelope (`11 99`) in addressed PGN 126720 at priority 3:

```go
remote, err := client.ActisenseRemoteDevice(35)
if err != nil { ... }
product, err := remote.GetProductInfo(ctx)
echoed, err := remote.Echo(ctx, []byte("probe"))
metrics := remote.Metrics()
```

Remote responses must come from source 35 and target the Client's current
claimed address. Connection and claim epochs are part of the pending key;
reconnect or address change completes old requests with
`ErrActisenseRemoteEpochChanged`. Pending requests, duplicate keys, and
multi-reply trains are bounded.

## Compiled command ledger

Baseline: Actisense SDK commit `9de7343`. “Implemented” means the command is
available locally and remotely through the shared typed Module unless noted.

| BEM | Capability | v2 classification |
|---:|---|---|
| `00` | Reinitialize main application | Implemented; explicit disruptive call only |
| `01` | Commit session settings to EEPROM | Implemented; explicit persistence call only |
| `02` | Commit session settings to flash | Implemented; explicit persistence call only |
| `11` | Get/set operating mode | Implemented; setup and restore verify acknowledgement |
| `13` | Get/set port P-Codes | Implemented, including no-change sentinel |
| `15` | Get/set total operating time | Implemented, including passkey form |
| `17` | Get/set port baudrate | Implemented, including session/store rates and all sentinels |
| `18` | Echo | Implemented; local maximum 222 bytes, remote maximum 206 bytes |
| `1B` | Port inventory | Implemented as a true multi-reply slot accumulator |
| `40` | Supported PGN list | Implemented as a transfer-ID checked chunk walk |
| `41` | Product information | Implemented for Format 2 and legacy five-reply devices; partial typed data survives timeout |
| `42` | Get/set CAN configuration/NAME | Implemented |
| `43` / `44` | Get/set installation descriptions | Implemented with bounded ASCII encoding |
| `45` | Get manufacturer CAN information | Implemented read-only |
| `46` | Get/set Rx PGN state and mask | Implemented |
| `47` | Get/set Tx PGN state and rate | Implemented; response retains timeout and priority |
| `4A` | Delete session PGN lists | Implemented; explicit mutation |
| `4B` | Activate session PGN lists | Implemented; explicit mutation |
| `4C` | Restore default PGN lists | Implemented; explicit disruptive/persistent firmware operation |
| `4D` | PGN-list parameters/status | Implemented |
| `4E` / `4F` | Rx/Tx Format-2 enable lists | Implemented multi-reply reads; NGX completion requires both every standard slot and the proprietary bitmap |
| `F0` / `F1` / `F2` / `F4` | Startup, error, system status, negative acknowledgement | Implemented typed unsolicited diagnostics |

The following declarations are not parity targets:

- `48`/`49` Format-1 lists are deprecated, removed from the compiled SDK, and
  negatively acknowledged by current firmware.
- `14` port duplicate delete is documentation-only and absent from the SDK
  implementation.
- `60`–`62` and `F3` have incomplete/TBD payload contracts and no usable
  compiled NMEA 2000 operation.
- NMEA 0183 and `!PARLB` are outside this library's Actisense scope.

## PGN-list lifecycle and persistence

`SendPGN` and `SendRawPGN` never query, enable, activate, or persist a PGN.
When a volatile transmit list needs staging, batch it deliberately:

```go
rate := uint32(1000)
err := session.ConfigureTransmitPGNs(ctx, []n2k.ActisenseTxPGNConfiguration{
    {PGN: 126992, Flag: n2k.ActisensePGNEnabled, Rate: &rate},
    {PGN: 126996, Flag: n2k.ActisensePGNRespondMode},
})
```

The transaction snapshots each distinct PGN, stages every change, activates
once, rolls back on failure, and best-effort restores the first snapshot on a
clean close in the same connection epoch. It does not commit persistence.

`CommitEEPROM`, `CommitFlash`, `Reinitialize`, `SetPortBaudrate`,
`SetCANConfig`, `SetCANInfoField`, and `DefaultPGNLists` may persist state or
interrupt the device depending on firmware. They exist for capability parity
but are never invoked by connection setup, sending, rollback, or close.

## Formats, capture, and metrics

| Representation | Parse | Emit | Writable role |
|---|:---:|:---:|---|
| BST-93 receive / BST-94 send | yes | yes | `ActisenseGatewaySession` only |
| BST-95 raw CAN | yes | yes | source-authoritative `Client` |
| BST-D0 assembled transport message | yes | n/a | receive/capture |
| CAN-frame ASCII mode 6 | yes | yes | source-authoritative `Client` |
| NMEA 2000 Type-A ASCII | yes | yes | passive source; gateway-owned identity |
| EBL | yes | yes | capture/replay and session wire trace |

`NewEBLWriter` writes SDK-compatible TimeUTC/Version/description, direction,
and BSTRaw records. `WithActisenseWireTrace` attaches a live trace to a gateway
session. Valid BDTP receives become checksum-stripped BSTRaw records; invalid,
partial, boot-banner, and other unframed bytes remain exact raw evidence.

Metrics cover transport calls/bytes/errors, valid datagrams, unframed bytes,
framing/checksum/length/oversize failures, per-BST counts, BEM requests,
responses, correlation misses, duplicates, timeouts, device errors, negative
acks, response-train overflow, in-flight high-water mark, and latency
minimum/average/maximum. Gateway session totals survive reconnect epochs.

## Compatibility evidence and deliberate deviations

[`conformance/actisense-golden.json`](../conformance/actisense-golden.json)
contains independent request and reply bytes for every compiled solicited verb,
remote wrapping, ASCII, and EBL. Local CI exercises the corpus, malformed
replies, timeouts, cross-group traffic, negative-ack scoping, partial results,
reconnects, and response #257 without blocking the sole reader.

`n2k` intentionally does not reproduce these reference defects:

- valid BDTP/BST-D0 payloads are not truncated to the SDK's 512-byte cap;
- duplicate requests are rejected instead of overwriting a callback;
- immediate write failure completes the typed request;
- callbacks/subscribers cannot stall the receive thread;
- F2 list completion cannot finish with missing standard slots;
- Port Inventory aggregates every advertised slot.

Run `just actisense-hardware <config>` for the opt-in NGT/NGX matrix described
in [Protocol conformance](conformance.md). The checked-in suite provides local
compatibility evidence; it does not claim that unavailable physical hardware
passed.
