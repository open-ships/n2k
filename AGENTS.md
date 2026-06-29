# Agent Requirements

- Before pushing, opening a PR, or marking work complete, run and pass:
  - `go test ./...`
  - `just lint`
  - `just secure`
- When generated PGN files, PGN metadata, or `pgn/manifest.json` change, also run and pass `go run ./cmd/canboatgen --check`.
- Do not treat lint or security failures as optional. Fix them before handoff unless an external blocker prevents a local pass.
