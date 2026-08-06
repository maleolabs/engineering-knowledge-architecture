---
namespace: feather
type: std
id: definition-of-done
instance-version: 1
revision: 1
content-state: approved
existence-state: active
dimension: standards
author: Maya Patel
created: 2026-06-10
updated: 2026-06-17
supersedes: []
derives-from: []
depends-on:
  - req:publishing-core
amends: []
validates: []
change-log:
  - date: 2026-06-10
    domain: existence-state
    from: "-"
    to: active
    by: Maya Patel
  - date: 2026-06-10
    domain: content-state
    from: "-"
    to: draft
    by: Maya Patel
  - date: 2026-06-13
    domain: content-state
    from: draft
    to: review
    by: Maya Patel
  - date: 2026-06-17
    domain: content-state
    from: review
    to: approved
    by: Maya Patel
---

# Standard — Definition of Done

## Purpose

The definition of done (DoD) for Feather work items: the conditions a work item must satisfy before its Execution State may move to `done`. It binds behavior and quality for every contributor, and it is the baseline the review artifact (`rvw-`) checks against.

## Content

A work item is **done** only when **all** of the following hold:

1. **Behavior** — the acceptance criteria of the work item are demonstrably met (manual or automated check recorded in the session).
2. **Tests** — unit tests cover the changed behavior; the full suite passes (`go test ./...`).
3. **Conformance** — `eka validate` reports zero errors for the repository; no new warnings without a recorded note.
4. **Documentation** — user-visible behavior changes are reflected in the relevant runbook or specification; terms used are registered in the glossary.
5. **Review** — the change passed the review gate: `execution-state: in-review` was held until the validating review (`rvw-`) recorded findings, and blocking findings are resolved.
6. **Release readiness** — no known open bug with severity above `minor` on the changed path (see `bug:empty-title-crash` as the severity baseline).

Exceptions require a written deviation note in the work item's session before the state transition is recorded.
