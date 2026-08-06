---
namespace: feather
type: epc
id: distribution
instance-version: 1
revision: 1
content-state: draft
existence-state: active
dimension: planning
author: Maya Patel
created: 2026-07-06
updated: 2026-07-06
supersedes: []
derives-from:
  - scp:mvp-v1:1
depends-on: []
amends: []
validates: []
change-log:
  - date: 2026-07-06
    domain: existence-state
    from: "-"
    to: active
    by: Maya Patel
  - date: 2026-07-06
    domain: content-state
    from: "-"
    to: draft
    by: Maya Patel
---

# Epic — Distribution

## Objective

Give readers a way to follow Feather sites without visiting them: RSS feeds and a sitemap. This epic is deliberately **post-MVP**: the strategy ranks author experience first and distribution second, so this capability is planned but not committed.

## Scope

- Atom RSS feed of published posts (`/feed.xml`), generated from the index.
- `sitemap.xml` for search engines.
- Conditional headers (ETag/Last-Modified) so feed readers poll cheaply.

## Out of Scope

- Newsletters, social posting, analytics, or engagement features.
- Anything requiring a background worker — generation stays on-read from the SQLite index.

This epic stays in `draft` until the MVP ships and the roadmap plan is revised; its work items will be created in a future container.
