# docs/operating/work-items/ — Work Items

> Anchor EKA: Operating Layer — State Domain **Execution State**.

## Purpose

Work items are the smallest unit of work executed within a container. Six subtypes cover all kinds of work; all share one owned State Domain: Execution State.

## Execution State

| Value | Meaning |
|---|---|
| `planned` | registered, not yet scheduled |
| `todo` | ready to work |
| `in-progress` | being worked |
| `in-review` | awaiting the review gate |
| `done` | finished and validated |

Rules: **never skip** (each value must be passed in sequence), **never revert** (no returning to a previous value). Transitions are forward-only and recorded in `change-log`.

## Single-Writer

Every work item has **one state writer** (its implementer). Only that writer changes `execution-state` and `existence-state`. The implementer writes `change-log` on every transition; reviewers never write work item state.

## Six Subtypes

| Token | Subtype | Folder |
|---|---|---|
| `sto-` | Story | [stories/](stories/) |
| `ts-` | Technical Story | [technical-stories/](technical-stories/) |
| `bug-` | Bug | [bugs/](bugs/) |
| `td-` | Tech Debt | [tech-debt/](tech-debt/) |
| `ch-` | Chore | [chores/](chores/) |
| `spk-` | Spike | [spikes/](spikes/) |

## Good Common Structure

All work items must contain:

- `## Description` — what is being done.
- Verification criteria per subtype (`## Acceptance Criteria`, `## Impact`, etc. — see subtype READMEs).
- Identity frontmatter: `namespace`, `type`, `id`; state: `execution-state`, `existence-state`; `change-log` for every transition.
- `dimension` may be set informationally (e.g. `requirements`), but **does not** determine the home folder.

## Related

- [containers/](../containers/) — work items live in the referenced `ctr-`.
- [projections/](../projections/) — `tkt-` projects work items into tables/status.
- [sessions/](../sessions/) — work item execution is recorded in `ses-`.
- [../quality/](../../quality/) — results are verified by `rvw-`.
