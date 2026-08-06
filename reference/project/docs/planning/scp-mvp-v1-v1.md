---
namespace: feather
type: scp
id: mvp-v1
instance-version: 1
revision: 1
content-state: approved
existence-state: active
phase: mvp
dimension: planning
author: Maya Patel
created: 2026-06-28
updated: 2026-07-03
supersedes: []
derives-from:
  - req:publishing-core
  - adr:content-storage
depends-on: []
amends: []
validates: []
change-log:
  - date: 2026-06-28
    domain: existence-state
    from: "-"
    to: active
    by: Maya Patel
  - date: 2026-06-28
    domain: content-state
    from: "-"
    to: draft
    by: Maya Patel
  - date: 2026-06-30
    domain: content-state
    from: draft
    to: review
    by: Maya Patel
  - date: 2026-07-03
    domain: content-state
    from: review
    to: approved
    by: Maya Patel
  - date: 2026-06-28
    domain: phase
    from: "-"
    to: mvp
    by: Maya Patel
---

# Scope — MVP v1

## Objective

Define the bounded scope of the Feather MVP: the smallest shippable product that satisfies the publishing-core requirement, under the phase `mvp`. Approved scope becomes the basis of the roadmap plan (`plan:roadmap-v1`).

## Scope

- Publishing core: create draft, autosave, preview, publish, edit published post (epic `epc:authoring-experience`).
- Markdown rendering with syntax highlighting and the syntax extensions approved by `spk:markdown-syntax-extension` (tables, footnotes).
- Search via SQLite FTS5 over title + body (`adr:search-sqlite-fts`).
- Single-author deployment behind Caddy (`dec:reverse-proxy`), with the deploy and backup runbooks (`run:deploy-feather`, `run:backup-feather`).
- Toolchain hygiene: Go 1.24 (chore `ch:update-go-version`), N+1 query cleanup (`td:reduce-query-count`).

## Out of Scope

- Comments (`req:comments-phase2` — stays draft for a later phase).
- Distribution beyond the site itself: RSS and sitemap (epic `epc:distribution`) are planned after the MVP, not in it.
- Multi-author workflows, roles, and permissions.
- Plugin model (`adr:plugin-model-deferred`), themes beyond the single default template.
- Social features, analytics, newsletters, media uploads.
