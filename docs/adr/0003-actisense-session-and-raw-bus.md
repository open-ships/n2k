# ADR 0003: Separate Actisense gateway sessions from raw CAN buses

Status: superseded by [ADR 0004](0004-actisense-session-contract.md)

This record preserves the v1 decision history. Its compatibility exception for
using gateway-owned BST-93/94 as a `Bus`, its implicit per-PGN transmit-list
mutation, and its prohibition on exposing explicit persistence commands are no
longer active decisions.

## Context

Legacy BST-93/94 carries assembled NMEA 2000 messages, and BST-94 omits the
source address because the gateway transmits under its own identity. Treating
that route as an ordinary `Bus` overstates the client's address-claim and
protocol semantics. Actisense raw mode uses BST-95 and retains the complete CAN
identifier, but entering that mode requires a correlated, acknowledged BEM
exchange before any NMEA network traffic is written.

## Decision

One bounded Actisense protocol Module owns BDTP framing, BST datagrams, local
BEM correlation, and a sole-reader session for each connection epoch. TCP,
serial, EBL, and observation code are thin Adapters around that Module.
`FormatActisense` remains the v1-compatible gateway-owned BST-93/94 message
session. `FormatActisenseRaw` is the source-authoritative BST-95 `Bus` and is
the preferred Actisense route for `NewClient`.

The public `ActisenseTCP` and `ActisenseSerial` constructors select behavior at
the existing source/`Bus` seam. `Receive`, `Observe`, and `NewScanner` use
compatible receive behavior; `NewClient` always selects the
source-authoritative raw CAN `Bus`. The explicit formats remain available as
escape hatches. There is no automatic UDP constructor because a one-way
datagram source cannot negotiate or verify an operating mode.

Every writable session queries its prior mode, sets and verifies the requested
mode, and only then publishes readiness. Mode 2 retains an active Tx PGN list,
so the legacy writer queries/enables each first-used PGN and activates the
session list before sending. Configuration changes remain volatile: the
library never commits EEPROM or flash, and it best-effort restores the prior
operating mode on a clean close. A rejected, timed-out, or stripped BEM request
is an explicit readiness failure; raw mode never falls back silently.

## Consequences

Existing `Receive` callers keep BST-93 behavior, while BST-D0 can enter the
Pipeline as an assembled Message without synthetic CAN frames. Applications
that require independent node identity use the dedicated constructors or
`FormatActisenseRaw`. NGT-class hardware that does not support mode 5 can be
read through the dedicated constructors, but writable use must explicitly
select the gateway-owned message session. Reconnect creates fresh parser,
correlation, Tx-list cache, and operating-mode state before the next connection
epoch becomes writable.
