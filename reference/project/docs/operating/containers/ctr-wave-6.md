---
namespace: feather
type: ctr
id: wave-6
instance-version: 1
revision: 1
container-state: completed
existence-state: active
author: Lukas Weber
created: 2026-07-16
updated: 2026-07-28
supersedes: []
derives-from: []
depends-on:
  - plan:roadmap-v1:1
amends: []
validates: []
change-log:
  - date: 2026-07-16
    domain: existence-state
    from: "-"
    to: active
    by: Lukas Weber
  - date: 2026-07-16
    domain: container-state
    from: "-"
    to: active
    by: Lukas Weber
  - date: 2026-07-28
    domain: container-state
    from: active
    to: completed
    by: Lukas Weber
---

# Container — Wave 6: Foundation

## Objective

Lay the foundation under the roadmap plan (`plan:roadmap-v1:1`): repository bootstrap, CI gate, storage layout, and the post scaffolding command. Everything the publishing core in Wave 7 builds on.

## Work Items

Delivered as prose (wave completed; work items closed out and archived per protocol):

- Repository bootstrap from the EKA skeleton; namespace `feather` established.
- CI pipeline running `eka validate` + `go test ./...` on every push.
- Storage layout per `adr:content-storage`: `content/posts/` + SQLite index with migration scaffolding.
- `feather new` scaffolding command (slug generation, frontmatter skeleton).

## Change Log

- 2026-07-16: container created (active), locking `plan:roadmap-v1:1`.
- 2026-07-28: aggregate of the wave's work items reached done; container marked completed by its owner. Wave 7 then became active (exactly-one-active preserved).
