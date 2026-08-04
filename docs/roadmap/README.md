# Roadmap

Sprint roadmaps — planned work breakdown derived from epic planning and the active MVP. One roadmap per MVP version. Roadmaps become immutable once sprints are generated from them.

## What Goes Here

- Roadmap documents (one per MVP version)
- Work breakdown planning from epics

## What Does NOT Go Here

- Sprint execution documents → `../sprints/`
- Epic definitions → `../epics/`
- MVP definitions → `../mvp/`

## Naming Convention

| Type | Pattern | Example |
|---|---|---|
| Roadmap | `mvp-nnn-v<version>.md` | `mvp-001-v1.md` |

One roadmap per MVP version.

## Status Lifecycle

```
Draft → Approved → Immutable
```

- **Draft**: Roadmap is being planned.
- **Approved**: Roadmap is finalized and ready for sprint generation.
- **Immutable**: Sprints have been generated from this roadmap; it cannot be changed.

## Rules

1. Roadmaps are derived from epic planning + actual work items registered in epic planning.
2. Roadmaps must align with the active MVP.
3. Once a sprint is generated from a roadmap, the roadmap becomes IMMUTABLE.
4. Planning changes after an immutable roadmap require a new roadmap version.
5. Never modify a roadmap that has generated sprints.

## Ownership

| Role | Responsibility |
|---|---|
| Product Owner | Defines roadmap scope from MVP |
| Tech Lead | Plans work breakdown from epics |

## Related Folders

- `../mvp/` — MVP defines the release scope
- `../epics/` — Epics provide capability area breakdown
- `../sprints/` — Sprints are generated from roadmaps
