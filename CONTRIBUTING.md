# Contributing

Issues and focused pull requests are welcome. For protocol changes, include a
capture, standard reference, or minimal frame sequence that demonstrates the
behavior. Keep generated PGN output separate from hand-written runtime changes
when practical.

## Development

Go 1.26.5 or newer and `just` are required. Run `just setup` once to install the
pinned development tools.

Before submitting a change, run:

```sh
go test ./...
just pgn-sync-check
just format-check
just lint
just secure
just release-check
```

Use `just test-race` for lifecycle/concurrency work and `just fuzz-smoke` for
parsers, codecs, framing, or transport protocols.

For network-management or transport changes, also run `just
conformance-local` and update the test/evidence mapping in
[docs/conformance.md](docs/conformance.md). This local gate is engineering
evidence, not a claim of NMEA certification. Licensed review and official-tool
records belong in the ignored `conformance-artifacts/` directory; commit only
non-confidential requirement identifiers and hashes.

PGN files and metadata are generated. Modify the generator or its inputs and
run `just pgn-sync`; do not hand-edit generated category files.

## Releases

[`VERSION`](VERSION) is the exact version the next successful `main` build must
release. Every pull request must increment it according to semantic versioning
and replace the first `CHANGELOG.md` section with a completed
`### v<VERSION> — YYYY-MM-DD — <summary>` entry, including documentation, CI,
and maintenance-only changes. `just release-check` rejects stale or mismatched
metadata and prevents the release workflow from silently inventing a patch
version.

Do not push release tags manually. After all required CI jobs pass, the release
workflow creates the annotated `v<VERSION>` tag at that exact `main` commit and
publishes its GitHub release. A missing or failed release blocks subsequent
merges until it is repaired and verified. The repository references an exact
semantic-version tag of the shared Open Ships release policy; that policy
publishes a deterministic source archive with checksums, an SBOM, and separate
build-provenance and SBOM attestations.

Architecture vocabulary and invariants live in [CONTEXT.md](CONTEXT.md).
Design rationale lives in [docs/adr](docs/adr).
