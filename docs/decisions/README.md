# Decisions

Operational and review decisions — reversible decisions from reviews, spikes, or discussions. Lighter weight than ADRs; used for decisions that can be changed without formal supersession.

## What Goes Here

- Operational decisions
- Review outcomes
- Spike conclusions
- Discussion summaries that result in a decision

## What Does NOT Go Here

- Irreversible architecture decisions → `../adr/`
- Architecture descriptions → `../architecture/`
- Product requirements → `../prd/`

## Naming Convention

| Type | Pattern | Example |
|---|---|---|
| Decision | `nnn-<decision-topic>.md` | `001-use-eslint-for-linting.md` |

Sequential numbering with descriptive slugs.

## Status Lifecycle

```
Draft → Accepted → Superseded (optional)
```

Decisions in this folder are inherently reversible. Supersession is informal — update or archive as needed.

## Ownership

| Role | Responsibility |
|---|---|
| Any contributor | Can create decision records |
| Tech Lead | Reviews and accepts |

## Related Folders

- `../adr/` — For irreversible architecture decisions
- `../reviews/` — Reviews that may produce decisions
- `../work-items/spikes/` — Spikes that may produce decisions
