# Contributing

Issues and focused pull requests are welcome. For protocol changes, include a
capture, standard reference, or minimal frame sequence that demonstrates the
behavior. Keep generated PGN output separate from hand-written runtime changes
when practical.

## Development

Go 1.23 or newer and `just` are required. Run `just setup` once to install the
pinned development tools.

Before submitting a change, run:

```sh
go test ./...
just pgn-sync-check
just lint
just secure
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

Architecture vocabulary and invariants live in [CONTEXT.md](CONTEXT.md).
Design rationale lives in [docs/adr](docs/adr).
