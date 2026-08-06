---
namespace: feather
type: tkt
id: bug-empty-title-crash
instance-version: 1
revision: 1
author: Lukas Weber
created: 2026-08-01
updated: 2026-08-02
supersedes: []
derives-from:
  - ctr:wave-7
  - bug:empty-title-crash
depends-on: []
amends: []
validates: []
---

> Generated — State Projection. Do NOT edit state here; refresh on read.

# Ticket — empty-title-crash

## Commands

- `eka view ticket tkt-bug-empty-title-crash` — render this projection (status derived from `bug:empty-title-crash`).
- `eka view execution` — see the active container board (Wave 7).
- `go test ./internal/api/ -run EmptyTitle` — run the regression test for the empty-title publish case.
- `curl -s -X POST localhost:8080/api/posts/p-7f3a/publish` with an empty title — expect HTTP 400, process alive.

## Projected Status

Projected from `bug:empty-title-crash` (owner state, read on refresh): **done**.
