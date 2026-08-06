---
namespace: feather
type: tkt
id: ts-markdown-renderer
instance-version: 1
revision: 1
author: Lukas Weber
created: 2026-07-29
updated: 2026-08-04
supersedes: []
derives-from:
  - ctr:wave-7
  - ts:markdown-renderer
depends-on: []
amends: []
validates: []
---

> Generated — State Projection. Do NOT edit state here; refresh on read.

# Ticket — markdown-renderer

## Commands

- `eka view ticket tkt-ts-markdown-renderer` — render this projection (status derived from `ts:markdown-renderer`).
- `eka view execution` — see the active container board (Wave 7).
- `go test ./internal/render/...` — run the golden-file renderer tests.
- `go test ./internal/render/ -run XSS` — run the safety corpus (script injection probes).

## Projected Status

Projected from `ts:markdown-renderer` (owner state, read on refresh): **in-review**.
