# Operations

Operational references — deployment guides, exit codes, output conventions, and server setup. This folder contains documents that support the operational side of the project.

## What Goes Here

- Deployment guides
- Exit code documentation
- Server setup instructions
- Migration guides
- Checklists
- Output conventions

## What Does NOT Go Here

- Architecture decisions → `../adr/`
- Architecture descriptions → `../architecture/`
- Product requirements → `../prd/`

## Naming Convention

| Type | Pattern | Example |
|---|---|---|
| Operations doc | Descriptive kebab-case | `deployment.md` |
| Migration guide | Descriptive kebab-case | `migration-guide-v2.md` |
| Checklist | Descriptive kebab-case | `release-checklist.md` |

Use descriptive names that clearly indicate the document's content.

## Ownership

| Role | Responsibility |
|---|---|
| DevOps | Creates deployment and infrastructure docs |
| Engineers | Contribute operational references |
| Tech Lead | Reviews operational documents |

## Related Folders

- `../architecture/` — Architecture informs operations
- `../adr/` — ADRs may affect operational decisions
- `../sessions/` — Sessions may reference operational guides
