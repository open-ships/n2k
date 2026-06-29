# Agent Requirements

- Before pushing, opening a PR, or marking work complete, run and pass:
  - `go test ./...`
  - `just lint`
  - `just secure`
- When auto-written PGN files, PGN metadata, or `pgn/manifest.json` change, also run and pass `go run ./cmd/canboatgen --check`.
- Files in `pgn/` must not include `generated` in their names. Use category-style names such as `navigation.go`, plus `dispatch.go` and `canboat_definitions.go`.
- Do not treat lint or security failures as optional. Fix them before handoff unless an external blocker prevents a local pass.
