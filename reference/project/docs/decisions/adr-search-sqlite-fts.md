---
namespace: feather
type: adr
id: search-sqlite-fts
instance-version: 1
revision: 1
content-state: accepted
existence-state: active
dimension: decisions
author: Lukas Weber
created: 2026-06-18
updated: 2026-06-25
supersedes: []
derives-from:
  - req:publishing-core
depends-on:
  - fnd:search-approach
amends: []
validates: []
change-log:
  - date: 2026-06-18
    domain: existence-state
    from: "-"
    to: active
    by: Lukas Weber
  - date: 2026-06-18
    domain: content-state
    from: "-"
    to: proposed
    by: Lukas Weber
  - date: 2026-06-25
    domain: content-state
    from: proposed
    to: accepted
    by: Lukas Weber
---

# ADR — Search with SQLite FTS5

## Context

The public site needs full-text search over posts. The architecture constraint (single binary, no external services) rules out dedicated search engines unless the search requirement is strong enough to break it. The research finding `fnd:search-approach` measured FTS5 performance at the projected scale.

## Decision

Use **SQLite FTS5** for v1 search. A `posts_fts` virtual table is populated from the post index on publish/update and rebuilt by `feather reindex`; queries go through the FTS5 MATCH operator with ranking.

## Consequences

- Search runs inside the existing database: no new service, no new operational surface, backup stays a single file.
- Sub-10 ms queries at 10k posts (measured in `fnd:search-approach`); the ceiling is far above Feather's audience.
- FTS5 strengths and limits are accepted: prefix and phrase matching work; stemming is English-only; there is no typo tolerance. If search quality ever becomes a product complaint, the decision is revisited — but the index design (a derived, rebuildable table) makes such a swap cheap.

## Alternatives Considered

- **LIKE/GLOB scans** — rejected: linear degradation, no ranking, no stemming.
- **External engine (Meilisearch)** — rejected: second service to run and back up; violates the single-binary constraint at a scale that does not need it.
