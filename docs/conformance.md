# Protocol conformance

This repository keeps two claims deliberately separate:

1. **Local protocol evidence** is public, deterministic, and runs in CI with
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

`conformance-local` gives the externally visible Commanded Address cases a
stable, verbose test run and reruns the deep protocol modules for CAN framing,
fast-packet discrimination, ISO transport, and address claiming. The full test
suite remains the authoritative local regression gate.

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
