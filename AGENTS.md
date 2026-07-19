# Agent Requirements

- Read `CONTEXT.md` before architecture or runtime changes; preserve its
  invariants and update it when domain language or module ownership changes.
- Before pushing, opening a PR, or marking work complete, run and pass:
  - `go test ./...`
  - `just pgn-sync-check`
  - `just lint`
  - `just secure`
- When auto-written PGN files, PGN metadata, or `pgn/manifest.json` need updating, regenerate them with `just pgn-sync`.
- Files in `pgn/` must not include `generated` in their names. Use category-style names such as `navigation.go`, plus `dispatch.go` and `upstream_definitions.go`.
- PGN struct names must not start with `Pgn`; use names such as `VesselHeading` and `AisClassAPositionReport`.
- Do not introduce the upstream source project's name as a contiguous token in file names, identifiers, comments, tests, generated output, or manifest fields. Exception: `README.md` names it in the "Why n2k" comparison and the acknowledgments — that usage is owner-approved; do not extend it elsewhere.
- Do not treat lint or security failures as optional. Fix them before handoff unless an external blocker prevents a local pass.
