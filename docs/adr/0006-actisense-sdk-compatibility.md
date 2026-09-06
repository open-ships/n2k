# ADR 0006: Actisense SDK compatibility scope and gateway remote identity

Status: accepted

## Context

The published SDK supports remote BEM through a gateway-owned Session as well
as generic BST/raw writes without changing the existing operating mode. n2k
previously required a source-authoritative Client for remote commands, which
excluded NGT gateways that cannot run mode 5. The August 25 SDK revision also
corrected the documented Rx-mask and Tx-rate sentinel semantics.

## Decision

Pin software compatibility to SDK commit
`ed2268a6e8db0645f75e4ef17eed2e937d025040`. Keep NMEA 0183 and `!PARLB` outside
scope. Add typed support for the complete documented BEM-14 and legacy F1
contracts. Preserve the compiled BEM-42 layout where the Markdown conflicts
with it. Do not invent analogue, legacy baud-code, or firmware-upload layouts
absent from the SDK.

Keep one remote correlation manager shared by Client and gateway sessions.
Gateway requests use BST-94 and the gateway's source. A random remote Echo
challenge verifies the live return address before each request. BEM CAN
configuration contains a stored address, which cannot establish this identity.
Registration binds both addresses and connection/identity epochs. Changes
cancel pending work and physical writes; reconnect never retries a command.
Only one probe runs at a time and Echo commands cannot overlap a probe to the
same source. A device change with no notification is discovered by the next
probe; mismatched destinations are rejected meanwhile.

Preserving the current operating mode is opt-in. Readiness still requires an
acknowledged getter. This option sends no implicit mode setter or restoration;
explicit mode setters update status. PGN writes require mode 1 or 2. Generic
BST and raw writes share the existing sole writer, copy inputs, have fixed size
limits and command deadlines, and fail immediately when disconnected.

## Consequences

Gateway remote operations add one Echo round trip. They work with the
gateway's actual identity without creating a virtual node or mutating PGN
lists. Tests cover wrong destinations, address changes, disconnects, mode
preservation, framing, and cancellation. Hardware runs record firmware,
runner revision, JSON outcomes, and hashed EBL captures separately.

The two misleading old mask/rate constants remain deprecated aliases to avoid
silently changing their bytes. Invalid Rx masks now fail locally. This is a
documented v1 contract correction. Software parity supports a scoped
compatibility claim; it cannot establish vendor approval or universal hardware
coverage.
