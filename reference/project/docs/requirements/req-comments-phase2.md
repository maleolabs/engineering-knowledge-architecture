---
namespace: feather
type: req
id: comments-phase2
instance-version: 1
revision: 1
content-state: draft
existence-state: active
dimension: requirements
author: Maya Patel
created: 2026-05-28
updated: 2026-05-28
supersedes: []
derives-from: []
depends-on:
  - vis:feather-vision
amends: []
validates: []
change-log:
  - date: 2026-05-28
    domain: existence-state
    from: "-"
    to: active
    by: Maya Patel
  - date: 2026-05-28
    domain: content-state
    from: "-"
    to: draft
    by: Maya Patel
---

# Requirement — Comments (Phase 2)

## Purpose

A future-phase requirement for reader comments on published posts. It is deliberately **not** part of the MVP scope — it exists to record the intent and to keep the publishing core free of comment concerns.

This requirement stays in `draft` until the MVP ships: draft tolerance is the honest state for a deferred need, and it lets the requirement evolve without governance ceremony.

## Content

**Statement.** Readers can comment on published posts, and authors can moderate comments (approve, hide, delete).

Open questions to resolve before this moves to `review`:

- Storage: comments in the SQLite index, or as files alongside posts?
- Moderation model: pre-approval, post-publication, or none (phase 2 may ship without moderation).
- Spam posture: whether any filtering is required at all for a small-audience platform.

Constraint from the vision: the solution must not require a separate service; it has to fit the single-binary architecture.
