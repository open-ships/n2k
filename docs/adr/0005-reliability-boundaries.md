# ADR 0005: Owned messages and cancellable network lifetimes

Status: Accepted — 2026-09-05

## Context

Independent probes found failure paths not covered by successful roundtrips:
string overflow, false physical availability, stale decode state, shared typed
messages, ISO timer races, stalled device writes, and requests surviving their
network session. Priority admission alone did not let protocol responses run
while an application transfer waited for pacing or its peer.

## Decision

- Encode outgoing payloads and clone headers synchronously at admission. Every
  subscriber, correlator reply, and registry snapshot owns a generated deep clone.
- Bind work to connection and claim epochs. A lifecycle transition cancels old
  operations and resets partial assembly. Application calls made while not ready
  are rejected, not retained for a later identity.
- Separate application/protocol jobs from physical wire serialization. Select
  protocol records between application frames, apply physical write deadlines,
  and never retry an uncertain write across reconnect.
- Give ISO sessions one completion owner, identity-checked timers, absolute
  deadlines, bounded state, and cancellable pacing and writes.
- Make real device I/O interruptible using pollable descriptors. Share serial
  ownership between USB-CAN and Actisense; bind gateway cleanup to its exact epoch.
- Bound request/discovery/rejoin/replay/scheduler state. Scheduled providers
  accept context, run in an owned worker, and must cooperate with cancellation.
- Generate from an offline checksummed snapshot, make decode transactional, and
  report support limitations independently of generated type counts.
- Execute the evidence index against compiled, actually passing tests. Missing
  or skipped required software evidence fails; hardware and certification remain
  explicit, separate qualification states.

## Consequences and v1 migration

The owner explicitly keeps the API on the existing v1 module path while it is
still evolving. v1.3 intentionally changes enum constant names to type-prefixed
names, physical metadata to float64, and broadcast providers to
`func(context.Context) pgn.Message`. Broadcast creation now returns `(stop, err)`.
Providers must observe context cancellation and must not call their own blocking
stop function or Client.Close. Schedule count is limited to 64.

`Write` now does encoding work before returning, but does not wait for queue
space or transmission. The caller may mutate its message after return, not
concurrently with the call. `WriteContext` cancels queued/in-progress work;
`WriteResult.WaitContext` cancels only the wait. Inspect `WriteError` before
deciding whether an uncertain transfer should be attempted again.

Replay capture retains the latest 4096 frames; callers needing complete history
must drain observations or write a capture, not use WrittenFrames indefinitely.
Device/codec breadth is not a promise of hardware completeness: the manifest
records zero hardware-verified PGNs until actual lab evidence exists.

These changes add copy/encode cost and can surface overload sooner. That cost
buys explicit ownership, bounded resource use, and failures that callers can
reason about without depending on internal scheduling or old network identity.
