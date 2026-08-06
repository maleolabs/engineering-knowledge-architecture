# sessions/ — Session (`ses-`)

> Anchor EKA: Operating Layer — State Domain **Existence State**; session artifact.

## Purpose

A session is an ephemeral record of one working opportunity: what was done, the context at the time, implementation notes, and verification performed. Sessions **carry no permanent progress status** — progress is held by other State Domains (work item Execution State).

## Token & State Vector

| Token | Folder | Owned state |
|---|---|---|
| `ses-` | `operating/sessions/` | `existence-state` |

Only one domain is owned: `existence-state` (`active → archived → retired`). Other state fields (content-state, execution-state, planning-state, container-state) are **not applicable** (N/A) to sessions.

## "Completed" = Derived Condition

Sessions have no "completed" value in frontmatter. "Completed" is a **derived condition** — a session is considered complete when the work it records is represented on a work item (moving toward `done`) and its verification is recorded. After that, the session only moves `existence-state` to `archived` (or `retired`).

## Ephemeral Content

Sessions are ephemeral — written fast, read fast, and never a permanent reference. Required structure:

- `## Context` — what is being worked on and its context (reference the `ctr-`/work item).
- `## Notes` — execution notes, micro-decisions, constraints.
- `## Verification` — evidence of verification performed (build, test, manual check).

## Mandatory Distillation Before Archival (EKA 11.4)

Before a `ses-` is archived, worthwhile findings **must be distilled**: decisions → `dec-`/`adr-`; research findings → `fnd-`; proven procedures → `run-`; new work items → created in the active container. Archiving a session with undistilled findings violates the protocol.

## Naming Conventions

`ses-<id>.md`, kebab-case, unique within (namespace, type). No `-v<nn>` suffix. Example: `ses-2026-08-05-projection-implementation.md`.

## Ownership

| Role | Responsibility |
|---|---|
| Engineer (session author) | single owner of the session's `existence-state` |
| All roles | write sessions for the work they do |

## Related

- [work-items/](../work-items/) — sessions serve work item execution.
- [containers/](../containers/) — sessions happen in the context of the active `ctr-`.
- [decisions/](../../decisions/) and [research/](../../research/) — session distillation destinations.
