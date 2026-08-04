# Architecture

Architecture documents describe domain models, system design, component boundaries, and technical structure. This folder is the primary reference for how the system is built.

## What Goes Here

- Domain models
- System design documents
- Component boundary definitions
- Numbered architecture documents
- Subfolders for specific domains (e.g., `configuration/`, `database/`)

## What Does NOT Go Here

- Architecture decisions → `../adr/` (for irreversible decisions)
- Operational guides → `../operations/`
- Product requirements → `../prd/`

## Naming Convention

| Type | Pattern | Example |
|---|---|---|
| Architecture doc | `nnn-<topic>.md` | `001-domain-model.md` |
| Domain subfolder | `<domain-name>/` | `configuration/` |

Architecture documents use sequential numbering with descriptive slugs.

## Status Lifecycle

```
Draft → Review → Approved → Amended
```

## Ownership

| Role | Responsibility |
|---|---|
| Tech Lead | Creates and maintains architecture docs |
| Engineers | Contribute domain-specific architecture |
| Architect | Reviews and approves |

## Related Folders

- `../adr/` — Irreversible architecture decisions
- `../decisions/` — Reversible operational decisions
- `../specification-corpus/` — Canonical definitions and terminology
