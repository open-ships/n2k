# Actisense NMEA 2000 support

n2k implements the published NMEA 2000 and binary control interface of the
Actisense SDK at commit `ed2268a6e8db0645f75e4ef17eed2e937d025040`, plus the
documented port duplicate-delete and legacy F1 commands. NMEA 0183 and `!PARLB`
are excluded. Software compatibility does not establish Actisense endorsement
or verification of every hardware model and firmware revision.

The v1 API separates two identities that the wire protocol cannot combine:

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

`Observations` emits owned frames, assembled messages, and gateway/transport evidence.
`Diagnostics` emits typed startup, error, system, and negative-ack records.
`Status` reports connection epoch, operating/receive role, known model
capabilities, subscriber counts, trace failure, and cumulative metrics.

For device control without changing the current mode, pass
`WithActisensePreserveOperatingMode()`. Each connection queries and acknowledges
the existing mode; startup and Close send no mode setter. An explicit
`SetOperatingMode` updates `Status().OperatingMode`. The last mode/preserve
option wins. `SendPGN`, `SendRawPGN`, and gateway remote requests require mode
1 or 2. A preserved raw-CAN mode still does not make the session a `Client`.

`SendBST(ctx, bytes)` accepts exactly one checksum-free BST record (ID, length,
payload) and adds its checksum and BDTP framing. `SendRaw(ctx, bytes)` writes
verbatim bytes. These match the SDK's two generic send modes, share the sole
writer with BEM, and accept unknown records. BST records are bounded to 1800
bytes; raw writes to 65536 bytes. The caller chooses bytes appropriate for the
device's mode. Both copy their inputs and honor the command timeout. No send
waits for a disconnected session to reconnect or retries a partial write.
`WithActisenseSessionReconnect` applies to TCP, serial, and custom streams.

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

The same handle is available from `session.ActisenseRemoteDevice(35)`, including
on NGT hardware that cannot support a raw-CAN Client. Before each request the
session sends a random 16-byte remote Echo challenge to verify the gateway's
current return address. It then sends the requested command through BST-94.
Replies must match the remote source, verified destination, response group,
verb, and connection/identity epochs. Probes are serialized; an outstanding
Echo to the same device rejects another probe with a request-in-flight error.
Other distinct commands can remain in flight concurrently.

`Status().GatewaySourceAddress` is nil until verified and clears on disconnect
or mode change. `IdentityEpoch` changes when this identity is invalidated or a
probe discovers a different address; pending work is canceled immediately.
The address is never inferred from `CANConfig.SourceAddress`, which is a stored
preferred/previous address. Arbitration can move the live address. A change
that firmware does not report is discovered by the next probe; a reply to an
unverified destination cannot complete a pending request.

## Compiled command ledger

Baseline: Actisense SDK commit `ed2268a`. “Implemented” means the command is
available locally and remotely through the shared typed Module unless noted.

| BEM | Capability | Classification |
|---:|---|---|
| `00` | Reinitialize main application | Implemented; explicit disruptive call only |
| `01` | Commit session settings to EEPROM | Implemented; explicit persistence call only |
| `02` | Commit session settings to flash | Implemented; explicit persistence call only |
| `11` | Get/set operating mode | Implemented; setup and restore verify acknowledgement |
| `13` | Get/set port P-Codes | Implemented, including no-change sentinel |
| `14` | Get/set port duplicate deletion | Implemented from the published wire contract; explicit setter persists immediately |
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
| `48` / `49` | Legacy Format-1 enable lists | Implemented explicit reads; two/four ordered parts, maximum 50 entries, partial results retained |
| `4A` | Delete session PGN lists | Implemented; explicit mutation |
| `4B` | Activate session PGN lists | Implemented; explicit mutation |
| `4C` | Restore default PGN lists | Implemented; explicit disruptive/persistent firmware operation |
| `4D` | PGN-list parameters/status | Implemented |
| `4E` / `4F` | Rx/Tx Format-2 enable lists | Implemented multi-reply reads; NGX completion requires both every standard slot and the proprietary bitmap |
| `F0` / `F1` / `F2` / `F4` | Startup, error, system status, negative acknowledgement | Implemented typed unsolicited diagnostics |

`GetRxPGNEnableListF1` and `GetTxPGNEnableListF1` are opt-in legacy methods.
Actisense discontinued F1 at firmware 2.500; NGX/W2K firmware never implemented
it. Current applications should use the F2 methods. F1 results expose
`PartsReceived` so partial lists cannot be mistaken for complete results.

The published SDK does not supply usable contracts for these features:

- Deprecated `12`/`16` baud-code arrays vary by model and have no complete wire
  definition or compiled implementation. Use modern `17` baudrate control.
- `60`–`62` and `F3` have incomplete/TBD payload contracts and no usable
  compiled NMEA 2000 operation.
- Firmware uploading and additional EMU/PRO configuration require vendor
  specifications beyond this SDK. They are not implemented by this change.
- NMEA 0183 and `!PARLB` are outside this library's Actisense scope.

`RawRequest` and `RawRequestMulti` expose bounded local/remote BEM access for
caller-defined commands. They are escape hatches, not typed support for an
unknown payload. BEM `42` follows the compiled SDK's nine-byte NAME and stored
address layout; its Markdown page describes a conflicting layout. That
discrepancy requires device evidence or vendor clarification, not guessing.

## Rx masks and Tx rates

The four accepted 32-bit Rx masks are `ActisenseRxPGNMaskPGN` (`03FFFF00`),
`ActisenseRxPGNMaskPDUFormat` (`03FF0000`), `ActisenseRxPGNMaskPDUNibble`
(`03F00000`), and `ActisenseRxPGNMaskDataPage` (`03000000`). They match PGN
ranges; they cannot filter source addresses. A nil mask or
`ActisenseRxPGNMaskDefault` selects the firmware default;
`ActisenseRxPGNMaskNoChange` preserves the existing mask. Replies report the
effective mask, not a sentinel. Arbitrary masks are rejected before sending.

Tx rates `1..65534` are milliseconds, and `ActisenseTxPGNRateEvent` is zero.
A nil rate or any value at least `65535` means no change, including
`FFFFFFFE` and `FFFFFFFF`. Use `ActisenseTxPGNRateNoChange`; this command has
no restore-default rate sentinel. The old `ActisenseRxPGNMaskAcceptAll` and
`ActisenseTxPGNRateDefault` names remain deprecated aliases preserving their
wire values; both actually mean no change. Sentinel acknowledgements are
validated against their effective semantics.

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
`SetCANConfig`, `SetCANInfoField`, `SetPortDuplicateDelete`, and `DefaultPGNLists` may persist state or
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
the three additional documented commands, remote wrapping, ASCII, and EBL.
These are protocol fixtures, not physical captures. Local CI exercises malformed
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

The machine-readable [compatibility scope](../conformance/actisense-compatibility.json)
distinguishes implemented software, absent specifications, deliberate exclusions,
and hardware evidence. The supported claim is compatibility with this published
NMEA 2000/binary SDK surface. “Officially approved by Actisense” and “all things
Actisense” require evidence that the codebase and local tests cannot provide.
