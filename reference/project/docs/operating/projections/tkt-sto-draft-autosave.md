---
namespace: feather
type: tkt
id: sto-draft-autosave
instance-version: 1
revision: 1
author: Lukas Weber
created: 2026-07-29
updated: 2026-08-05
supersedes: []
derives-from:
  - ctr:wave-7
  - sto:draft-autosave
depends-on: []
amends: []
validates: []
---

> Generated — State Projection. Do NOT edit state here; refresh on read.

# Ticket — draft-autosave

## Commands

- `eka view ticket tkt-sto-draft-autosave` — render this projection (status derived from `sto:draft-autosave`).
- `eka view execution` — see the active container board (Wave 7).
- `go test ./internal/editor/...` — run the autosave/debounce tests.
- `./feather serve --addr :8080` — manual check: type, reload, expect the buffer restored.

## Projected Status

Projected from `sto:draft-autosave` (owner state, read on refresh): **in-progress**.
