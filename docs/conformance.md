# Protocol conformance

This repository keeps two claims deliberately separate:

1. **Local protocol evidence** is public, executable, and runs in CI with
   `just conformance-local` plus the normal test, race, lint, and security
   gates.
2. **Formal product certification** requires licensed standards, the official
   NMEA certification hardware/software, assigned manufacturer and product
   codes, representative product hardware, and NMEA validation of the tool's
   encrypted result. The public test suite does not claim or imply that status.

## Controlled baseline

Baseline metadata was checked against the publishers on 2026-07-18:

| Scope | Controlled reference | Why it matters |
|---|---|---|
| CAN-frame mapping and transport/network protocol | [ISO 11783-3:2026, edition 5](https://www.iso.org/standard/89949.html), published 2026-03 | Current ISO transport baseline, including the transport layer used to carry PGN 65240. |
| Source-address and NAME management | [ISO 11783-5:2019, edition 3](https://www.iso.org/standard/74366.html), confirmed current in 2024 | Network-management baseline for address changes and subsequent claims. |
| Marine application and certification | [NMEA 2000 Version 3.000 with amendments](https://www.nmea.org/nmea-2000.html) | Product-facing NMEA requirements and the official certification process. |

The ISO and NMEA documents and official test vectors are licensed material.
They are not copied into this repository. A licensed reviewer must maintain a
private requirement crosswalk and record only non-confidential requirement
identifiers, results, and artifact hashes here.

## Reproducible local gate

Run:

```sh
just conformance-local
go test ./...
just test-race
just pgn-sync-check
just lint
just secure
```

`conformance-local` gives the externally visible protocol cases a stable,
executable evidence run and reruns the deep protocol modules for CAN framing,
fast-packet discrimination, ISO transport, address claiming, gateway parsing,
reconnect, saturation, and raw observation. The full test suite remains the
authoritative local regression gate.

The machine-readable evidence index is
[`conformance/requirements.json`](../conformance/requirements.json). Its schema,
unique IDs, and evidence presence are checked by
`TestConformanceRequirementIndex`.
Schema v2 references package-qualified, anchored test patterns. Discovery uses
compiled test names, not source-text searching. `cmd/conformance` then consumes
`go test -json`: missing matches, failed tests, and skipped required tests
(including skipped child cases) fail the software evidence gate. It preserves
`conformance-artifacts/local-test-events.jsonl` and `local-evidence.json`.
Hardware and soak entries are explicitly `not-run` in this command; ordinary
unit test success cannot turn either into a pass.

### Sustained reliability

```sh
just reliability-soak 1h 70m
just fuzz-long 2m
```

The software soak cycles bounded subscribers and writes, malformed traffic,
reconnect/reclaim, terminal failure/recreation, and replay. Its JSON artifact
records actual duration, cycle counts, overflows, and bounded recent post-GC
heap/goroutine samples, platform, and the exact executed test binary's SHA-256.
It tolerates runtime noise but fails sustained growth
over declared thresholds; it is not a proof of absence of leaks or a hardware
benchmark. Set `N2K_SOAK_ARTIFACT_DIR` when running concurrent qualification jobs
so they do not overwrite one another's evidence.

CI runs a short smoke; the scheduled reliability workflow runs the required
gates, a one-hour software soak, race detection, and longer fuzz campaigns.
`just fuzz-smoke` runs each of the five harnesses with a 10,000-execution budget
and a two-minute hard timeout. Count-based completion avoids the timed fuzz
coordinator's reported [cancellation race](https://github.com/golang/go/issues/75804).
All failures still fail CI, and any failing input corpus is uploaded for replay.
The longer, time-based campaigns remain separate from this PR smoke gate.
Hardware qualification remains a separate 24-hour lab activity using the
declared adapter/firmware, nominal and saturated traffic, reconnects, and power
cycling. Record exact commit, configuration, duration, failures, and capture
hashes. If hardware is absent, report `not-run`, not pass.

The generated manifest separately records typed coverage, decode/encode
completeness, known codec limitations, schema revision/checksum, and hardware
verification. v1.3 includes 599 typed variants over 348 PGNs; 108 variants are
marked complete for encode and decode, and none are hardware-verified.

`main` requires the aggregate GitHub Actions `release-gate` status with strict
up-to-date checking. The desired rule is versioned in
`.github/main-ruleset.json`; changing a workflow file alone does not install or
update GitHub repository protection. Publishing still follows the exact
VERSION/tag/release discipline in `AGENTS.md`.

### Requirement families

| Evidence ID | Local behavior family | Representative executable evidence |
|---|---|---|
| AC-01 | Address claim, defense, retry, cannot-claim, and surfaced failures | `internal/claiming` suite |
| HB-01 | Heartbeat cadence, retiming, disable, and sequence rollover | `TestHeartbeat_*` |
| IR-01 | ISO request targeting, responses, and acknowledgements | `TestISORequest_*` |
| GF-01 | Supported group-function requests, explicit unsupported-command/parameter acknowledgements, and scheduling | `TestGroupFunction_*` |
| FP-01 | Fast-packet limits, assembly, ownership, and malformed input | `internal/adapter`; `internal/framer` suites |
| BAM-01 / RTS-01 | BAM and RTS/CTS transmit/receive state machines | `internal/transport` suite |
| MAL-01 | Malformed frames, announcements, checksums, and streams | transport and gateway hostile-input tests |
| SAT-01 | Bounded application, protocol, subscription, and reassembly state | queue, hub, adapter, and transport tests |
| REC-01 | Reconnect epoch reclaims before protocol restart | `TestTCPClientReconnectReclaimsBeforeRestartingProtocolTraffic` |
| TIM-01 | Heartbeat cadence, assembly expiry, transport timeout, and bounded startup waits | deterministic clock checks where provided, plus bounded real-time waits |
| CODEC-01 | Conditional fields, signed widths, ranges, and sentinels | codec and data-stream writer tests |
| OBS-01 | Source/network identity, timestamps, direction, and ownership | observation and gateway tests |
| ACT-01 | Honest gateway-session versus source-authoritative Client roles, acknowledged readiness, explicit Tx-list mutation | `TestGatewaySessionDoesNotMutateTransmitListImplicitly`; `TestNewClient_TCPActisense_RequiresGatewaySession`; raw reconnect/source tests |
| ACT-02 | Complete compiled solicited BEM surface, independent local request/reply bytes, bounded correlation, typed partials | `conformance/actisense-golden.json`; `TestActisenseGolden*`; `internal/actisense` session tests |
| ACT-03 | Remote BEM PGN-126720 envelope, address/epoch correlation, cancellation, errors, and metrics | `TestConformanceActisenseRemoteGoldenEnvelope`; `TestActisenseRemote*` |
| ACT-04 | CAN ASCII mode 6 and N2K Type-A ASCII parse/emit fidelity and bounds | `TestConformanceActisenseASCIIGoldenVectors`; `internal/gateway` ASCII tests |
| EBL-01 | Bounded EBL read/write, live wire trace, timestamps, direction, exact invalid/unframed evidence, Type-2 replay | `TestConformanceActisenseEBLGoldenVersionRecord`; `internal/ebl`; public trace/replay tests |

These IDs are public repository evidence identifiers. They are not licensed
ISO or NMEA clause numbers, and they do not replace the private crosswalk.

### Actisense parity evidence

[`conformance/actisense-golden.json`](../conformance/actisense-golden.json)
pins the reference implementation at Actisense SDK commit `9de7343` and stores
independent bytes for every compiled solicited BEM command and response, the
remote PGN-126720 envelope, CAN/N2K ASCII, and EBL. Tests consume this corpus
directly. The ledger in [Actisense support](actisense.md) classifies compiled,
deprecated, documentation-only, and intentionally superseded behavior.

Physical evidence is opt-in because public CI has no NGT/NGX hardware. Copy
[`conformance/actisense-hardware.example.json`](../conformance/actisense-hardware.example.json),
replace the example endpoints, then run:

```sh
just actisense-hardware conformance/actisense-hardware.local.json
```

The runner exercises local product/mode/echo commands, selected supported/F2
lists and port inventory, live EBL trace, source-authoritative raw Client
startup, and optional remote BEM. Record hardware model, serial-safe identity,
firmware, runner commit, result, and capture hash outside git. An absent or
skipped lab configuration is not a hardware pass and does not support a 100%
hardware-parity claim.

### Commanded Address evidence map

| Evidence ID | Expected behavior | Executable evidence |
|---|---|---|
| CA-01 | Accept PGN 65240 only after a complete broadcast ISO transport transfer. | `TestConformanceCommandedAddressMatchingBAMReclaimsAddress`; `TestConformanceCommandedAddressRejectsFastPacketAndAddressedTP` |
| CA-02 | Require an exact nine-byte payload. | `TestConformanceCommandedAddressRequiresExactTransferAndTarget` |
| CA-03 | Compare all 64 NAME bits exactly. | `TestConformanceCommandedAddressRequiresExactTransferAndTarget` (the negative vector differs only at bit 63) |
| CA-04 | Accept only claimable destination addresses 0–251. | `TestConformanceCommandedAddressRequiresExactTransferAndTarget`; `TestClaimer_CommandedAddressRejectsNonTargetsAndSpecialAddresses` |
| CA-05 | Ignore a command for the current address. | `TestClaimer_CommandedAddressRejectsNonTargetsAndSpecialAddresses` |
| CA-06 | Change the runtime source and immediately broadcast Address Claim with the same NAME. | `TestConformanceCommandedAddressMatchingBAMReclaimsAddress` |
| CA-07 | Start a fresh contention window before application traffic resumes. | `TestConformanceCommandedAddressMatchingBAMReclaimsAddress` |
| CA-08 | Preserve automatic versus explicit contention policy after the command. | `TestClaimer_CommandedAddressReclaimsImmediately`; `TestClaimer_CommandedAddressPreservesExplicitContentionPolicy` |
| CA-09 | Surface failure to transmit the required new claim as a terminal protocol failure. | `TestClaimer_CommandedAddressClaimFailureIsReturned`; `TestConformanceCommandedAddressClaimFailureTerminatesClient` |

These IDs are repository evidence identifiers, not ISO or NMEA clause numbers.
An authorized reviewer must map them to the licensed editions above and add or
change cases wherever the normative requirements differ.

## Licensed ISO review

For each release candidate that changes framing, transport, claiming, timing,
or protocol responses:

1. Pin the exact implementation commit and Go toolchain.
2. Have an authorized reviewer cross-check the implementation and tests against
   ISO 11783-3:2026 and the current network-management reference. Record the
   reviewer, edition, review date, result, and private matrix SHA-256.
3. Run the complete suite on classical 250 kbit/s extended CAN hardware under
   nominal, saturated, reordered, duplicated, malformed, timeout, contention,
   disconnect, and restart conditions relevant to the product.
4. Store licensed matrices and raw lab output outside git. Put hashes and the
   non-confidential result metadata in an evidence record.

Until that licensed review is recorded, describe the result as “local protocol
regression passed,” never “ISO conformant.”

## Official NMEA certification run

NMEA states that transmitting/receiving products must be certified through its
official process. The certification tool produces an encrypted result for NMEA
validation, and product certification requires assigned manufacturer and
product codes. The tool, license, product identity, and physical unit under
test are not present in this repository, so this step cannot run in public CI.

For a release candidate:

1. Obtain the current standard/tool directly from NMEA and record their exact
   versions and amendments.
2. Configure the real product with its assigned stable NAME, manufacturer code,
   product code, production firmware, transceiver, and power behavior.
3. Run every mandatory official test and the product's declared transmit and
   receive PGN set. Do not substitute mocks or a community test suite.
4. Hash the encrypted output, submit it to NMEA, and record NMEA's validation
   result and date. Keep licensed tool output in controlled storage.
5. Copy [the evidence example](../conformance/evidence.example.json) to the
   ignored `conformance-artifacts/` directory and fill it without embedding
   licensed content, credentials, serial secrets, or raw test vectors.

Formal release evidence is complete only when local gates pass at the recorded
commit, the licensed ISO crosswalk is signed, the official tool run passes, and
NMEA has validated its encrypted result.
