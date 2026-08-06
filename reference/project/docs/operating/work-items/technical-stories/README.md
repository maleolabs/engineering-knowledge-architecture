# technical-stories/ — Technical Story (`ts-`)

> Anchor EKA: Operating Layer — State Domain Execution State; work item subtype `ts-`.

## Purpose

A technical story is a unit of work delivering internal technical value (infrastructure, refactoring, system integration) that is not directly visible to users but is required for system quality and sustainability.

## Token & State Vector

| Token | Folder | Owned state |
|---|---|---|
| `ts-` | `work-items/technical-stories/` | `execution-state`, `existence-state` |

`execution-state` values: `planned → todo → in-progress → in-review → done` (forward-only, never skip, never revert). `existence-state` values: `active → archived → retired`.

## Required Content Structure

- `## Description` — the technical work and its rationale.
- `## Acceptance Criteria` — measurable conditions satisfying the definition of done (including technical verification).

## Naming Conventions

`ts-<id>.md`, kebab-case, unique within (namespace, type). No `-v<nn>` suffix. Example: `ts-frontmatter-format-migration.md`.

## Ownership

| Role | Responsibility |
|---|---|
| Tech Lead | definition of technical work |
| Engineer (implementer) | single owner of state; execution |
| DevOps | collaborator for infrastructure work |

## Related

- [architecture/](../../../architecture/) — technical work embodies `arc-`.
- [standards/](../../../standards/) — results must conform to `std-`.
- [containers/](../../containers/) — technical stories are owned by the active container.
