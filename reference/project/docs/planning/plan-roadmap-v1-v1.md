---
namespace: feather
type: plan
id: roadmap-v1
instance-version: 1
revision: 1
content-state: approved
planning-state: approved
existence-state: active
phase: mvp
dimension: planning
author: Maya Patel
created: 2026-07-08
updated: 2026-07-14
supersedes: []
derives-from:
  - scp:mvp-v1:1
depends-on: []
amends: []
validates: []
change-log:
  - date: 2026-07-08
    domain: existence-state
    from: "-"
    to: active
    by: Maya Patel
  - date: 2026-07-08
    domain: content-state
    from: "-"
    to: draft
    by: Maya Patel
  - date: 2026-07-10
    domain: content-state
    from: draft
    to: review
    by: Maya Patel
  - date: 2026-07-14
    domain: content-state
    from: review
    to: approved
    by: Maya Patel
  - date: 2026-07-08
    domain: planning-state
    from: "-"
    to: draft
    by: Maya Patel
  - date: 2026-07-14
    domain: planning-state
    from: draft
    to: approved
    by: Maya Patel
  - date: 2026-07-08
    domain: phase
    from: "-"
    to: mvp
    by: Maya Patel
---

# Plan — Roadmap v1

## Objective

Sequence the MVP work into execution containers and record the committed order of delivery. This plan is approved; it locks atomically with the birth of its first container (`ctr:wave-7`), after which changes require a new instance (lock-atomic-with-generation, protocol §4).

## Scope

| Wave | Container | Delivers | Phase |
|---|---|---|---|
| Wave 6 | `ctr:wave-6` | Foundation: repository bootstrap, CI, storage layout, post scaffolding — completed | mvp |
| Wave 7 | `ctr:wave-7` | Publishing core: publish, autosave, renderer, bug fix, debt, chore, spike (active) | mvp |
| Wave 8 | (future) | Distribution: RSS + sitemap (epic `epc:distribution`) | mvp |
| Wave 9 | (future) | Comments phase 2 (`req:comments-phase2`) | milestone |

Ordering rationale: author experience first (strategy), then the cheap distribution forms, then comments.

## Out of Scope

- Multi-author, plugin model, themes, engagement features — outside the MVP phase (see `scp:mvp-v1`).
- Waves 8–9 are planned directions, not commitments; they are re-planned when the plan is next revised (new instance).
