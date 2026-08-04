# Spikes

Research and investigation tasks. Spikes are time-boxed explorations to reduce technical uncertainty before implementation begins.

## What Goes Here

- Spike definitions and scope
- Spike findings and conclusions
- Technical investigation results

## What Does NOT Go Here

- Implementation work → `../stories/`, `../technical-stories/`
- Architecture decisions → `../../adr/`
- Bug investigations → `../bugs/`

## Naming Convention

| Type | Pattern | Example |
|---|---|---|
| Spike | `sp-nnn-<slug>.md` | `sp-001-evaluate-cache-strategies.md` |

## Status

Every work item uses one of the following statuses. No other status is allowed.

```
Planned → Todo → In Progress → In Review → Done
```

| Status | Meaning | Condition |
|---|---|---|
| `Planned` | Default status when assigned to a sprint | Sprint generation (initial state) |
| `Todo` | Will be worked on — the item has entered an execution wave | Wave becomes active |
| `In Progress` | Implementation is underway | Assignee starts working |
| `In Review` | Pull Request has been created | PR opened |
| `Done` | Implementation is merged | PR merged to the target branch |

Status transitions are strictly sequential. Never skip a status. Never revert to a previous status.

## Ownership

| Role | Responsibility |
|---|---|
| Engineers | Execute spikes |
| Tech Lead | Defines spike scope and time-box |

## Related Folders

- `../decisions/` — Spike conclusions may produce decisions
- `../../adr/` — Spike findings may lead to ADRs
- `../../sessions/` — Spike work may use implementation sessions
