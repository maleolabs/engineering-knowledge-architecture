# ADR (Architecture Decision Records)

Architecture Decision Records capture irreversible technical decisions with their context, rationale, and consequences. ADRs are binding once accepted.

## What Goes Here

- Architecture Decision Records
- Superseded ADRs (kept for historical reference)

## What Does NOT Go Here

- Reversible operational decisions → `../decisions/`
- Architecture descriptions → `../architecture/`
- Product decisions → `../prd/`

## Naming Convention

| Type | Pattern | Example |
|---|---|---|
| ADR | `adr-nnn-<decision-topic>.md` | `adr-003-runtime-and-release-lifecycle.md` |

## ADR Template

Every ADR should follow this structure:

```markdown
# ADR-NNN: <Decision Title>

| Field    | Value       |
|----------|-------------|
| Status   | Proposed    |
| Date     | YYYY-MM-DD  |
| Author   | <name>      |

## Context

What is the issue that we're seeing that is motivating this decision?

## Decision

What is the change that we're proposing and/or doing?

## Consequences

What becomes easier or harder because of this change?

## Alternatives Considered

What other options were evaluated?

## References

Links to related documents, discussions, or evidence.
```

## Status Lifecycle

```
Proposed → Accepted → Superseded
```

- **Proposed**: ADR is being discussed; not yet binding.
- **Accepted**: ADR is finalized and binding on the team.
- **Superseded**: ADR has been replaced by a new ADR (link to the replacement).

## Ownership

| Role | Responsibility |
|---|---|
| Tech Lead | Creates ADRs for major technical decisions |
| Engineers | Propose ADRs, participate in review |
| Architect | Approves ADRs |

## Related Folders

- `../architecture/` — System design and component descriptions
- `../decisions/` — Reversible operational decisions
