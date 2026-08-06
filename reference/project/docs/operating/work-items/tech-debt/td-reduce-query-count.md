---
namespace: feather
type: td
id: reduce-query-count
instance-version: 1
revision: 1
execution-state: planned
existence-state: active
dimension: architecture
author: Lukas Weber
created: 2026-08-03
updated: 2026-08-03
supersedes: []
derives-from: []
depends-on:
  - ctr:wave-7
  - plan:roadmap-v1:1
amends: []
validates: []
change-log:
  - date: 2026-08-03
    domain: existence-state
    from: "-"
    to: active
    by: Lukas Weber
  - date: 2026-08-03
    domain: execution-state
    from: "-"
    to: planned
    by: Lukas Weber
---

# Tech Debt — Reduce N+1 Queries in Post Listing

## Description

The post listing endpoint (`GET /api/posts`) issues one query per post row: the base list query fetches index rows, then a per-post query resolves each post's file metadata (title from frontmatter, word count, updated-at). With the MVP's post counts this is invisible, but the pattern scales linearly and the FTS/listing path is exactly where the index should be doing the work.

Location: `internal/posts/list.go` (list handler + `PostSummary` loader).

## Acceptance Criteria

- `GET /api/posts` issues a constant number of queries regardless of post count (single join or indexed lookup).
- Post counts below 10,000 remain correct; the FTS search path is unaffected.
- The change is covered by the existing API tests plus a query-count assertion (count queries via `sqlite` trace hook).
- `eka validate` stays clean; no public API change.

## Debt Rationale

The N+1 pattern was accepted deliberately in Wave 6 to ship the storage layer fast (foundation wave, no listing perf requirement in scope). The debt is now registered because the listing path becomes a hot path once Wave 7's publish flow goes live and the site gets real traffic. Repaying it is a small, contained change with a measurable assertion — the cheapest point to repay is before the next container plans around this code.
