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

## Status Lifecycle

```
Draft → Ready → In Progress → In Review → Verified → Done
```

## Ownership

| Role | Responsibility |
|---|---|
| Engineers | Execute spikes |
| Tech Lead | Defines spike scope and time-box |

## Related Folders

- `../decisions/` — Spike conclusions may produce decisions
- `../../adr/` — Spike findings may lead to ADRs
- `../../sessions/` — Spike work may use implementation sessions
