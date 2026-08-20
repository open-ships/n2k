# ADR 0004: Actisense session and control-plane contract

Status: accepted

## Context

ADR 0003 correctly distinguished gateway-owned BST-93/94 messages from
source-authoritative BST-95 CAN frames, but retained a v1 compatibility path
that still represented the gateway-owned route as a `Bus`. It also coupled
BST-94 writes to implicit Tx-PGN list changes and made EEPROM/flash operations
unavailable. Those choices conflict with an honest `Client` identity and with
the compiled Actisense NMEA 2000 control surface.

The same typed BEM commands must work against a locally attached gateway and a
remote Actisense device over addressed PGN 126720. Multi-reply commands,
diagnostics, reconnects, address changes, and promiscuously observed traffic
make a sole-reader session and explicit response origin part of correctness,
not implementation detail.

## Decision

### Separate node and gateway identities

`NewClient` accepts only source-authoritative Actisense representations:
BST-95 binary raw CAN (`FormatActisenseRaw`) and mode-6 CAN ASCII
(`FormatActisenseCANASCII`). `ActisenseTCP` and `ActisenseSerial` select BST-95
when used with `NewClient`. Gateway-owned `FormatActisense` and
`FormatActisenseN2KASCII` remain read-only sources, and `NewClient` rejects
them with `ErrActisenseGatewaySessionRequired`.

Gateway-owned BST-93/94 transmission belongs to `ActisenseGatewaySession`.
That Interface sends assembled PGNs under the physical gateway's identity and
reports `SourceAuthoritative=false`; it never address-claims for a virtual
node.

### One deep BEM Module, two envelope Adapters

`internal/actisense` owns all compiled NMEA 2000-relevant BEM verbs, response
decoding, model capabilities, bounded multi-reply accumulation, diagnostics,
timeouts, and metrics. Local commands use BST-A1/A0 through the gateway
session. Remote commands reuse the same typed Interface through the Actisense
manufacturer envelope in addressed PGN 126720.

Local correlation includes response BST group, BEM verb, and origin. Remote
correlation additionally includes remote source, the local destination
snapshot, connection epoch, and address-claim epoch. Address or connection
changes cancel pending remote work. Duplicate correlation keys are rejected;
response trains and pending tables are bounded.

### Make gateway mutation deliberate and reversible

Sending a PGN never changes a Tx-PGN list. `ConfigureTransmitPGNs` snapshots
all affected entries, stages them, activates once, rolls back a failed
transaction, and best-effort restores the first snapshot on clean close in the
same epoch.

Potentially persistent or disruptive operations are public explicit methods:
`CommitEEPROM`, `CommitFlash`, `Reinitialize`, `SetPortBaudrate`,
`SetCANConfig`, `SetCANInfoField`, and `DefaultPGNLists`. The library never
calls them implicitly. This supersedes ADR 0003's absolute prohibition on
EEPROM/flash support while preserving safe defaults.

### Preserve evidence and capability honesty

TCP, serial, and custom byte-stream constructors share the same session
lifecycle. Serial settings are configurable and default to 115200 8N1.
Operating-mode readiness is acknowledged per epoch, and clean close
best-effort restores the prior mode.

Status exposes whether traffic is filtered, receive-all, or
source-authoritative; known model caveats are explicit. Per-layer transport,
BDTP/BST, BEM, latency, reset, and correlation metrics are cumulative across
reconnects. An optional EBL wire trace retains valid checksum-stripped BST
records and exact invalid or unframed bytes.

Independent golden vectors are pinned to Actisense SDK commit `9de7343`. Known
SDK defects—its 512-byte BDTP cap, duplicate callback overwrite, lossy send
failure, unsafe callback-thread behavior, and incomplete F2/port accumulation—
are recorded deviations rather than compatibility targets.

## Consequences

This deliberately changes v1 behavior while the library's public API is
documented as in flux: code that used `NewClient` with `FormatActisense` must
choose a source-authoritative raw format or move gateway-owned sends to
`ActisenseGatewaySession`. The distinction removes plausible but false
address-claim, heartbeat, request, and transport semantics.

The control plane is deeper but has one command model and two narrow envelope
Adapters. Persistence and reset capabilities exist without becoming hidden
side effects. Hardware parity remains an evidence claim: the public corpus and
opt-in NGT/NGX runner make it reproducible, but a hardware model is not marked
passed until a real lab run is recorded.
