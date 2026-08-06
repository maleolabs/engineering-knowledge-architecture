---
namespace: feather
type: ch
id: update-go-version
instance-version: 1
revision: 1
execution-state: todo
existence-state: active
dimension: operations
author: Lukas Weber
created: 2026-08-04
updated: 2026-08-04
supersedes: []
derives-from: []
depends-on:
  - ctr:wave-7
  - plan:roadmap-v1:1
amends: []
validates: []
change-log:
  - date: 2026-08-04
    domain: existence-state
    from: "-"
    to: active
    by: Lukas Weber
  - date: 2026-08-04
    domain: execution-state
    from: "-"
    to: planned
    by: Lukas Weber
  - date: 2026-08-04
    domain: execution-state
    from: planned
    to: todo
    by: Lukas Weber
---

# Chore — Update Go Toolchain

## Description

Move the project to Go 1.24: bump `go.mod` (`go 1.24`), update the CI image (`golang:1.24`), and refresh dependencies (goldmark, chroma, SQLite driver) to versions that build under the new toolchain. No product behavior changes; the toolchain floor matters for the security fixes and the new standard-library APIs the renderer can use.

## Acceptance Criteria

- `go.mod` declares `go 1.24`; `go build ./...` and `go test ./...` pass locally and in CI.
- CI uses the 1.24 image; the `eka validate` gate still passes.
- Dependency bumps are limited to compatible releases; any breaking change is recorded as a note in the session.
- No runtime behavior change: the full test suite (including renderer golden files) passes unchanged.
