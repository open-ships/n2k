# Agent Requirements

- Before pushing, opening a PR, or marking work complete, run and pass:
  - `go test ./...`
  - `just lint`
  - `just secure`
- When auto-written PGN files, PGN metadata, or `pgn/manifest.json` change, also run and pass `go run ./cmd/pgngen --check`.
- Files in `pgn/` must not include `generated` in their names. Use category-style names such as `navigation.go`, plus `dispatch.go` and `upstream_definitions.go`.
- Do not introduce the upstream source project's name as a contiguous token in file names, identifiers, comments, tests, docs, generated output, or manifest fields.
- Do not treat lint or security failures as optional. Fix them before handoff unless an external blocker prevents a local pass.
