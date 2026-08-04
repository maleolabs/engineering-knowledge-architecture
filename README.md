# Project Documentation Template

Standardized documentation structure for engineering projects. This template provides a complete, opinionated folder layout and governance rules for managing product, architecture, planning, and implementation documentation.

## What Is This?

A reusable `docs/` structure that enforces:

- Consistent folder placement — every document has a home
- Kebab-case naming — no ambiguous filenames
- Document lifecycle management — Draft → Review → Approved → Superseded
- Clear ownership — each folder has a defined owner role
- Workflow alignment — PRD → MVP → Epics → Roadmap → Sprint → Tickets → Work Items → Sessions → Reviews

## How To Use

1. Copy the `docs/` folder into your project root.
2. Read [`docs/README.md`](docs/README.md) — it is the single source of truth for structure, workflow, and rules.
3. Adapt folder contents to your product's MVP and PRD.
4. Remove folders you don't need, or add new ones following the naming conventions.

## Naming Conventions

| Rule | Example |
|---|---|
| All filenames use kebab-case | `mvp-001-anvil-v1.md` |
| No full capslock filenames | WRONG: `ANVIL_MANIFESTO.md` / RIGHT: `anvil-manifesto.md` |
| No underscores in filenames | WRONG: `my_document.md` / RIGHT: `my-document.md` |
| No PascalCase or camelCase | WRONG: `MyDocument.md` / RIGHT: `my-document.md` |
| README.md is the only exception | `README.md` (conventional uppercase) |

## Document Lifecycle

Living documents follow a status lifecycle:

```
Draft → Review → Approved → Amended (via amendment docs)
```

Architecture Decision Records (ADRs):

```
Proposed → Accepted → Superseded (by new ADR)
```

Work items follow per-type lifecycles — see `docs/work-items/README.md`.

## Structure Overview

```
docs/
├── prd/                    # Product Requirements Documents
├── mvp/                    # Minimum Viable Product definitions
├── epics/                  # Epic definitions
├── architecture/           # Architecture documents
├── adr/                    # Architecture Decision Records
├── decisions/              # Operational/reversible decisions
├── roadmap/                # Sprint roadmaps (immutable once sprints exist)
├── sprints/                # Sprint execution documents (one active at a time)
├── tickets/                # Wave/ticket execution documents
├── work-items/             # Work items (stories, bugs, tech debt, chores, spikes)
├── reviews/                # Code and architecture reviews
├── sessions/               # Implementation session contexts
├── operations/             # Deployment guides, exit codes, server setup
├── planning/               # Cross-cutting planning artifacts
├── manifesto/              # Product manifesto and principles
└── specification-corpus/   # Canonical definitions and shared terminology
```

See [`docs/README.md`](docs/README.md) for the full structure guide, workflow, and governance rules.
