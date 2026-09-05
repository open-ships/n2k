# Agent Requirements

- Read `CONTEXT.md` before architecture or runtime changes; preserve its
  invariants and update it when domain language or module ownership changes.
- Before pushing, opening a PR, or marking work complete, run and pass:
  - `go test ./...`
  - `just pgn-sync-check`
  - `just format-check`
  - `just lint`
  - `just secure`
  - `just release-check`
- When auto-written PGN files, PGN metadata, or `pgn/manifest.json` need updating, regenerate them with `just pgn-sync`.
- Files in `pgn/` must not include `generated` in their names. Use category-style names such as `navigation.go`, plus `dispatch.go` and `upstream_definitions.go`.
- PGN struct names must not start with `Pgn`; use names such as `VesselHeading` and `AisClassAPositionReport`.
- Do not introduce the upstream source project's name as a contiguous token in file names, identifiers, comments, tests, generated output, or manifest fields. Exception: `README.md` names it in the "Why n2k" comparison and the acknowledgments — that usage is owner-approved; do not extend it elsewhere.
- Do not treat lint or security failures as optional. Fix them before handoff unless an external blocker prevents a local pass.

## Release discipline

- Every change merged into `main` is release-bearing. Update `VERSION` and add
  a completed top-level `CHANGELOG.md` release section in the same pull request,
  even for documentation, CI, or maintenance-only work.
- `VERSION` contains the exact semantic version to publish, without a leading
  `v`. The first `###` section in `CHANGELOG.md` must be
  `### v<VERSION> — YYYY-MM-DD — <summary>`; do not merge an `Unreleased`
  section or allow already-published notes to accumulate there.
- Choose the increment according to semantic versioning. Exported compatible
  functionality requires at least a minor increment; compatible fixes and
  maintenance use a patch increment. While the owner-designated v1 API remains
  in flux, intentional breaking contract refinements may use a minor increment
  without a new module path; document every affected API and migration in the
  completed release notes. Once compatibility is promised, breaking changes
  require a new major module path.
- Run `just release-check` after fetching current tags. It must prove that the
  proposed version is newer than every published tag, or that the matching
  annotated tag already points at the exact commit being verified. Never rely
  on release automation to invent a patch version.
- A successful `main` build must publish the annotated `v<VERSION>` tag and
  GitHub release for that exact commit. A failed, missing, or mismatched release
  is a blocking regression: fix and rerun the release before starting or
  merging further feature work, then verify both the tag and release.
