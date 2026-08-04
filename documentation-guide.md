# Documentation Guide

> **This document is the consolidated merge of all README.md files in the `docs/` structure.** For the per-folder granular version, see the individual README.md files.

This guide serves as the primary context document for understanding the entire `docs/` structure, workflow, and rules — in one place.

---

## Table of Contents

- [1. Overview](#1-overview)
- [2. Documentation Structure & Governance](#2-documentation-structure--governance)
  - [A. Documentation Structure](#a-documentation-structure)
  - [B. Naming Conventions](#b-naming-conventions)
  - [C. Folder Placement Rules](#c-folder-placement-rules)
  - [D. Documentation Workflow](#d-documentation-workflow)
  - [E. Active Sprint Rule](#e-active-sprint-rule)
  - [F. Roadmap Rules](#f-roadmap-rules)
  - [G. Document Status Lifecycle](#g-document-status-lifecycle)
  - [H. Metadata Convention](#h-metadata-convention)
  - [I. Adding New Folders](#i-adding-new-folders)
  - [J. References](#j-references)
- [3. Folder Reference](#3-folder-reference)
  - [3.1 prd/](#31-prd)
  - [3.2 mvp/](#32-mvp)
  - [3.3 epics/](#33-epics)
  - [3.4 architecture/](#34-architecture)
  - [3.5 adr/](#35-adr)
  - [3.6 decisions/](#36-decisions)
  - [3.7 roadmap/](#37-roadmap)
  - [3.8 sprints/](#38-sprints)
  - [3.9 tickets/](#39-tickets)
  - [3.10 work-items/](#310-work-items)
    - [3.10.1 bugs/](#3101-bugs)
    - [3.10.2 stories/](#3102-stories)
    - [3.10.3 technical-stories/](#3103-technical-stories)
    - [3.10.4 tech-debt/](#3104-tech-debt)
    - [3.10.5 chores/](#3105-chores)
    - [3.10.6 spikes/](#3106-spikes)
    - [3.10.7 planning/](#3107-planning)
  - [3.11 reviews/](#311-reviews)
  - [3.12 sessions/](#312-sessions)
  - [3.13 operations/](#313-operations)
  - [3.14 planning/](#314-planning)
  - [3.15 manifesto/](#315-manifesto)
  - [3.16 specification-corpus/](#316-specification-corpus)
- [4. Quick Reference](#4-quick-reference)
  - [4.1 Naming Conventions Summary](#41-naming-conventions-summary)
  - [4.2 Status Model Summary](#42-status-model-summary)
  - [4.3 Documentation Workflow Summary](#43-documentation-workflow-summary)
- [5. Source Files](#5-source-files)

---

## 1. Overview

*Source: `README.md` (repo-level)*

### What Is This?

A reusable `docs/` structure that enforces:

- Consistent folder placement — every document has a home
- Kebab-case naming — no ambiguous filenames
- Document lifecycle management — Draft → Review → Approved → Superseded
- Clear ownership — each folder has a defined owner role
- Workflow alignment — PRD → MVP → Epics → Roadmap → Sprint → Tickets → Work Items → Sessions → Reviews

### How To Use

1. Copy the `docs/` folder into your project root.
2. Read this document — it is the single source of truth for structure, workflow, and rules.
3. Adapt folder contents to your product's MVP and PRD.
4. Remove folders you don't need, or add new ones following the naming conventions.

### Naming Conventions

| Rule | Example |
|---|---|
| All filenames use kebab-case | `mvp-001-anvil-v1.md` |
| No full capslock filenames | WRONG: `ANVIL_MANIFESTO.md` / RIGHT: `anvil-manifesto.md` |
| No underscores in filenames | WRONG: `my_document.md` / RIGHT: `my-document.md` |
| No PascalCase or camelCase | WRONG: `MyDocument.md` / RIGHT: `my-document.md` |
| README.md is the only exception | `README.md` (conventional uppercase) |

### Document Lifecycle

Living documents follow a status lifecycle:

```
Draft → Review → Approved → Amended (via amendment docs)
```

Architecture Decision Records (ADRs):

```
Proposed → Accepted → Superseded (by new ADR)
```

Work items follow per-type lifecycles — see [Section 3.10](#310-work-items).

### Structure Overview

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

---

## 2. Documentation Structure & Governance

*Source: `docs/README.md` — the single source of truth for the documentation structure, naming conventions, workflow, and governance rules.*

### A. Documentation Structure

| Folder | Purpose | File Naming | Owner |
|---|---|---|---|
| `prd/` | Product Requirements Documents — what the product must do and why | `prd-nnn-<product-name>.md` | Product Owner |
| `mvp/` | Minimum Viable Product definitions — scope boundaries per release milestone | `mvp-nnn-<product-name>-<version>.md` | Product Owner |
| `epics/` | Epic definitions — capability areas that break MVPs into work streams | `epic-nnn-<capability-area>.md` | Product Owner / Tech Lead |
| `architecture/` | Architecture documents — domain models, system design, component boundaries | `nnn-<topic>.md` (numbered) | Tech Lead |
| `adr/` | Architecture Decision Records — irreversible technical decisions with rationale | `adr-nnn-<decision-topic>.md` | Tech Lead / Engineers |
| `decisions/` | Operational/review decisions — reversible decisions from reviews, spikes, or discussions | `nnn-<decision-topic>.md` | Any contributor |
| `roadmap/` | Sprint roadmaps — planned work breakdown derived from epic planning + active MVP. Immutable once sprints are generated from them | `mvp-nnn-v<version>.md` | Product Owner / Tech Lead |
| `sprints/` | Sprint documents — execution snapshots from the latest roadmap. One active sprint at a time | `mvp-nnn-s<nn>.md` | Scrum Master / Tech Lead |
| `tickets/` | Wave/ticket execution documents — deterministic wave decomposition of active sprints, consumed by execution commands | `<sprint-id>-wave-tickets.md` | Tech Lead |
| `work-items/` | Work item definitions — Technical Stories, Stories, Bugs, Tech Debt, Chores, Spikes. Organized by type in subfolders | Per type (see below) | Engineers |
| `reviews/` | Review documents — code reviews, architecture reviews, audit findings | `nnn-<review-topic>.md` | Reviewers |
| `sessions/` | Implementation session contexts — ephemeral working documents from implementation sessions | `impl-<work-items>-<date>/` folders | Engineers |
| `operations/` | Operational references — deployment guides, exit codes, output conventions, server setup | Descriptive kebab-case | DevOps / Engineers |
| `planning/` | Cross-cutting planning artifacts — traceability matrices, transition plans | Descriptive kebab-case | Tech Lead |
| `manifesto/` | Product manifesto — vision, principles, non-negotiables | `<product>-manifesto.md` | Product Owner |
| `specification-corpus/` | Specification vocabulary — canonical definitions, lifecycle models, shared terminology | Descriptive kebab-case | Tech Lead |

### B. Naming Conventions

#### Strict Rules

1. **All filenames MUST use kebab-case**: lowercase letters, numbers, hyphens between words.
2. **NEVER use full capslock filenames**: WRONG `ANVIL_MANIFESTO.md`, RIGHT `anvil-manifesto.md`.
3. **NEVER use underscores in filenames** (except in document IDs inside file content).
4. **NEVER use PascalCase or camelCase in filenames**.
5. **README.md is the ONLY exception** (conventional uppercase).

#### Examples

| Correct | Incorrect |
|---|---|
| `mvp-001-anvil-v1.md` | `MVP_001_ANVIL_V1.md` |
| `adr-003-runtime-and-release-lifecycle.md` | `ADR_003_Runtime_and_Release_Lifecycle.md` |
| `bug-001-description-field-null.md` | `BUG_001_Description_Field_Null.md` |
| `anvil-manifesto.md` | `ANVIL_MANIFESTO.md` |
| `traceability-matrix.md` | `Traceability_Matrix.md` |

#### Work Item File Naming

Work item files use their ID prefix in lowercase:

| Type | Pattern | Example |
|---|---|---|
| Bug | `bug-nnn-<slug>.md` | `bug-002-description-field-null.md` |
| Story | `st-nnn-nnn-<slug>.md` | `st-012-001-user-login-flow.md` |
| Technical Story | `ts-nnn-nnn-<slug>.md` | `ts-012-001-database-migration.md` |
| Tech Debt | `td-nnn-<slug>.md` | `td-002-remove-legacy-auth.md` |
| Chore | `ch-nnn-nnn-<slug>.md` | `ch-001-001-update-dependencies.md` |
| Spike | `sp-nnn-<slug>.md` | `sp-001-evaluate-cache-strategies.md` |

#### Session File Naming

Session folders: `impl-<work-items>-<date>/`

Context files inside sessions (all lowercase):
- `context.md`
- `notes.md`
- `verification.md`

### C. Folder Placement Rules

#### Strict Rules

1. **EVERY file MUST be inside a specific subfolder.** No loose files in `docs/` root.
2. If a document doesn't fit any existing folder, create a new folder with a `README.md` or place it in the most relevant existing folder.
3. Migration guides, checklists, and project-level docs go in the most relevant subfolder:
   - Migration guides → `operations/` (see [Section 3.13](#313-operations))
   - Project structure docs → `architecture/` (see [Section 3.4](#34-architecture))
   - Checklists → `operations/` or the relevant folder

#### Common Mistakes

| Mistake | Correct Placement |
|---|---|
| `docs/project-structure.md` (loose file) | `docs/architecture/project-structure.md` |
| `docs/migration-guide-v2.md` (loose file) | `docs/operations/migration-guide-v2.md` |
| `docs/ANVIL_MANIFESTO.md` (capslock) | `docs/manifesto/anvil-manifesto.md` |
| `docs/CONTEXT.md` (capslock) | `docs/sessions/<session>/context.md` |

### D. Documentation Workflow

Documents flow through the system in this order:

```
PRD → MVP → Epic Planning → Roadmap → Sprint → Ticket → Work Items → Sessions
                                                                         ↓
                                                                     Reviews
```

#### Workflow Steps

| Step | Input | Output | Description |
|---|---|---|---|
| PRD | Business requirements | `prd/prd-nnn-*.md` | Defines what the product must do and why |
| MVP | PRD | `mvp/mvp-nnn-*.md` | Defines release scope boundaries from PRD |
| Epic Planning | MVP | `epics/epic-nnn-*.md` | Breaks MVP into capability areas (epics) |
| Roadmap | Epic Planning + MVP | `roadmap/mvp-nnn-v*.md` | Planned work breakdown; one per MVP version |
| Sprint | Roadmap | `sprints/mvp-nnn-s*.md` | Execution snapshot from roadmap; one active at a time |
| Ticket | Sprint | `tickets/<sprint-id>-wave-tickets.md` | Deterministic wave decomposition of active sprint |
| Work Items | Sprint/Ticket | `work-items/<type>/*.md` | Implementation units within sprints |
| Sessions | Work Items | `sessions/impl-*/` | Ephemeral implementation contexts |
| Reviews | Completed work | `reviews/nnn-*.md` | Validation of completed work |

### E. Active Sprint Rule

**ONLY ONE active sprint at a time.**

- A sprint becomes active when generated from the roadmap.
- The next sprint CANNOT be generated until the active sprint is completed (all work items Done).
- Sprint completion is recorded in the sprint's Change Log.
- If a sprint is partially completed, deferred items roll to the next sprint.

#### Sprint Lifecycle

```
Active → Completed
```

- **Active**: Sprint is in progress; work items are being worked on.
- **Completed**: All work items are Done; sprint is closed.

### F. Roadmap Rules

1. **Roadmaps are derived from epic planning + actual work items registered in epic planning.**
2. **Roadmaps must align with the active MVP.**
3. **Once a sprint is generated from a roadmap, the roadmap becomes IMMUTABLE.**
4. Planning changes after an immutable roadmap require a new roadmap version.
5. Never modify a roadmap that has generated sprints.

#### Roadmap Lifecycle

```
Draft → Approved → Immutable
```

- **Draft**: Roadmap is being planned.
- **Approved**: Roadmap is finalized and ready for sprint generation.
- **Immutable**: Sprints have been generated; roadmap cannot be changed.

### G. Document Status Lifecycle

#### Living Documents (PRD, MVP, Epic, Architecture)

```
Draft → Review → Approved → Amended
```

- **Draft**: Document is being written.
- **Review**: Document is under review.
- **Approved**: Document is finalized and accepted.
- **Amended**: Document has been updated via a separate amendment document.

#### Architecture Decision Records (ADR)

```
Proposed → Accepted → Superseded
```

- **Proposed**: ADR is being discussed.
- **Accepted**: ADR is finalized and binding.
- **Superseded**: ADR has been replaced by a new ADR.

#### Work Items

Per-type lifecycles — see [Section 3.10](#310-work-items) for details. The unified 5-status model is documented in [Section 4.2](#42-status-model-summary).

#### Roadmaps

```
Draft → Approved → Immutable
```

#### Sprints

```
Active → Completed
```

#### Sessions

```
Active → Completed → Archived
```

- **Active**: Session is in progress.
- **Completed**: Session work is done.
- **Archived**: Session is preserved for reference.

### H. Metadata Convention

Documents that support metadata use a YAML frontmatter or Markdown table at the top of the file:

```markdown
| Field       | Value                    |
|-------------|--------------------------|
| Status      | Draft                    |
| Author      | Jane Doe                 |
| Created     | 2025-01-15               |
| Updated     | 2025-01-20               |
| Version     | 1.0                      |
```

Status values should match the lifecycle defined for that document type.

### I. Adding New Folders

If you need a new documentation folder:

1. Create the folder under `docs/`.
2. Add a `README.md` explaining purpose, naming conventions, ownership, and related folders.
3. Update this document's structure table (Section A).
4. Follow kebab-case naming for the folder itself.

### J. References

- [Section 1](#1-overview) — what this template is and how to use it
- [Section 3.5](#35-adr) — Architecture Decision Record format
- [Section 3.10](#310-work-items) — work item types and lifecycles
- [Section 3.12](#312-sessions) — implementation session structure

---

## 3. Folder Reference

### 3.1 prd/

*Source: `docs/prd/README.md`*

Product Requirements Documents define what the product must do and why. PRDs are the starting point of the documentation workflow.

#### What Goes Here

- Product requirement documents
- PRD amendments (changes to approved PRDs)

#### What Does NOT Go Here

- MVP definitions → [Section 3.2](#32-mvp)
- Architecture decisions → [Section 3.5](#35-adr)
- Implementation details → [Section 3.10](#310-work-items)

#### Naming Convention

| Type | Pattern | Example |
|---|---|---|
| PRD | `prd-nnn-<product-name>.md` | `prd-001-anvil.md` |
| Amendment | `prd-nnn-amendment-nn-<slug>.md` | `prd-001-amendment-01-scope-change.md` |

#### Status Lifecycle

```
Draft → Review → Approved → Amended
```

#### Ownership

| Role | Responsibility |
|---|---|
| Product Owner | Creates and maintains PRDs |
| Tech Lead | Reviews technical feasibility |
| Stakeholders | Approve PRD content |

#### Related Folders

- [Section 3.2](#32-mvp) — MVP definitions derive from PRDs
- [Section 3.3](#33-epics) — Epics break down PRD requirements
- [Section 3.15](#315-manifesto) — Product vision and principles

---

### 3.2 mvp/

*Source: `docs/mvp/README.md`*

Minimum Viable Product definitions establish scope boundaries per release milestone. MVPs translate PRD requirements into achievable release increments.

#### What Goes Here

- MVP scope definitions
- MVP amendments (changes to approved MVPs)
- Release milestone boundaries

#### What Does NOT Go Here

- Product requirements → [Section 3.1](#31-prd)
- Epic breakdowns → [Section 3.3](#33-epics)
- Sprint planning → [Section 3.7](#37-roadmap)

#### Naming Convention

| Type | Pattern | Example |
|---|---|---|
| MVP | `mvp-nnn-<product-name>-<version>.md` | `mvp-001-anvil-v1.md` |
| Amendment | `mvp-nnn-amendment-nn-<slug>.md` | `mvp-001-amendment-01-scope-adjustment.md` |

#### Status Lifecycle

```
Draft → Review → Approved → Amended
```

#### Ownership

| Role | Responsibility |
|---|---|
| Product Owner | Defines MVP scope |
| Tech Lead | Validates technical feasibility |
| Stakeholders | Approve MVP boundaries |

#### Related Folders

- [Section 3.1](#31-prd) — PRDs define what the product must do
- [Section 3.3](#33-epics) — Epics break MVPs into capability areas
- [Section 3.7](#37-roadmap) — Roadmaps plan work within an MVP

---

### 3.3 epics/

*Source: `docs/epics/README.md`*

Epic definitions represent capability areas that break MVPs into work streams. Each epic groups related work items under a coherent theme.

#### What Goes Here

- Epic definitions
- Epic scope and acceptance criteria
- Epic planning artifacts

#### What Does NOT Go Here

- Individual work items → [Section 3.10](#310-work-items)
- Sprint planning → [Section 3.7](#37-roadmap)
- Product requirements → [Section 3.1](#31-prd)

#### Naming Convention

| Type | Pattern | Example |
|---|---|---|
| Epic | `epic-nnn-<capability-area>.md` | `epic-001-user-authentication.md` |

#### Status Lifecycle

```
Draft → Review → Approved → Amended
```

#### Ownership

| Role | Responsibility |
|---|---|
| Product Owner | Defines epic scope |
| Tech Lead | Breaks epics into work items |
| Engineers | Implement epic work items |

#### Related Folders

- [Section 3.2](#32-mvp) — MVPs define which epics are in scope
- [Section 3.10](#310-work-items) — Work items implement epic capabilities
- [Section 3.7](#37-roadmap) — Roadmaps schedule epic work items

---

### 3.4 architecture/

*Source: `docs/architecture/README.md`*

Architecture documents describe domain models, system design, component boundaries, and technical structure. This folder is the primary reference for how the system is built.

#### What Goes Here

- Domain models
- System design documents
- Component boundary definitions
- Numbered architecture documents
- Subfolders for specific domains (e.g., `configuration/`, `database/`)

#### What Does NOT Go Here

- Architecture decisions → [Section 3.5](#35-adr) (for irreversible decisions)
- Operational guides → [Section 3.13](#313-operations)
- Product requirements → [Section 3.1](#31-prd)

#### Naming Convention

| Type | Pattern | Example |
|---|---|---|
| Architecture doc | `nnn-<topic>.md` | `001-domain-model.md` |
| Domain subfolder | `<domain-name>/` | `configuration/` |

Architecture documents use sequential numbering with descriptive slugs.

#### Status Lifecycle

```
Draft → Review → Approved → Amended
```

#### Ownership

| Role | Responsibility |
|---|---|
| Tech Lead | Creates and maintains architecture docs |
| Engineers | Contribute domain-specific architecture |
| Architect | Reviews and approves |

#### Related Folders

- [Section 3.5](#35-adr) — Irreversible architecture decisions
- [Section 3.6](#36-decisions) — Reversible operational decisions
- [Section 3.16](#316-specification-corpus) — Canonical definitions and terminology

---

### 3.5 adr/

*Source: `docs/adr/README.md`*

Architecture Decision Records capture irreversible technical decisions with their context, rationale, and consequences. ADRs are binding once accepted.

#### What Goes Here

- Architecture Decision Records
- Superseded ADRs (kept for historical reference)

#### What Does NOT Go Here

- Reversible operational decisions → [Section 3.6](#36-decisions)
- Architecture descriptions → [Section 3.4](#34-architecture)
- Product decisions → [Section 3.1](#31-prd)

#### Naming Convention

| Type | Pattern | Example |
|---|---|---|
| ADR | `adr-nnn-<decision-topic>.md` | `adr-003-runtime-and-release-lifecycle.md` |

#### ADR Template

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

#### Status Lifecycle

```
Proposed → Accepted → Superseded
```

- **Proposed**: ADR is being discussed; not yet binding.
- **Accepted**: ADR is finalized and binding on the team.
- **Superseded**: ADR has been replaced by a new ADR (link to the replacement).

#### Ownership

| Role | Responsibility |
|---|---|
| Tech Lead | Creates ADRs for major technical decisions |
| Engineers | Propose ADRs, participate in review |
| Architect | Approves ADRs |

#### Related Folders

- [Section 3.4](#34-architecture) — System design and component descriptions
- [Section 3.6](#36-decisions) — Reversible operational decisions

---

### 3.6 decisions/

*Source: `docs/decisions/README.md`*

Operational and review decisions — reversible decisions from reviews, spikes, or discussions. Lighter weight than ADRs; used for decisions that can be changed without formal supersession.

#### What Goes Here

- Operational decisions
- Review outcomes
- Spike conclusions
- Discussion summaries that result in a decision

#### What Does NOT Go Here

- Irreversible architecture decisions → [Section 3.5](#35-adr)
- Architecture descriptions → [Section 3.4](#34-architecture)
- Product requirements → [Section 3.1](#31-prd)

#### Naming Convention

| Type | Pattern | Example |
|---|---|---|
| Decision | `nnn-<decision-topic>.md` | `001-use-eslint-for-linting.md` |

Sequential numbering with descriptive slugs.

#### Status Lifecycle

```
Draft → Accepted → Superseded (optional)
```

Decisions in this folder are inherently reversible. Supersession is informal — update or archive as needed.

#### Ownership

| Role | Responsibility |
|---|---|
| Any contributor | Can create decision records |
| Tech Lead | Reviews and accepts |

#### Related Folders

- [Section 3.5](#35-adr) — For irreversible architecture decisions
- [Section 3.11](#311-reviews) — Reviews that may produce decisions
- [Section 3.10.6](#3106-spikes) — Spikes that may produce decisions

---

### 3.7 roadmap/

*Source: `docs/roadmap/README.md`*

Sprint roadmaps — planned work breakdown derived from epic planning and the active MVP. One roadmap per MVP version. Roadmaps become immutable once sprints are generated from them.

#### What Goes Here

- Roadmap documents (one per MVP version)
- Work breakdown planning from epics

#### What Does NOT Go Here

- Sprint execution documents → [Section 3.8](#38-sprints)
- Epic definitions → [Section 3.3](#33-epics)
- MVP definitions → [Section 3.2](#32-mvp)

#### Naming Convention

| Type | Pattern | Example |
|---|---|---|
| Roadmap | `mvp-nnn-v<version>.md` | `mvp-001-v1.md` |

One roadmap per MVP version.

#### Status Lifecycle

```
Draft → Approved → Immutable
```

- **Draft**: Roadmap is being planned.
- **Approved**: Roadmap is finalized and ready for sprint generation.
- **Immutable**: Sprints have been generated from this roadmap; it cannot be changed.

#### Rules

1. Roadmaps are derived from epic planning + actual work items registered in epic planning.
2. Roadmaps must align with the active MVP.
3. Once a sprint is generated from a roadmap, the roadmap becomes IMMUTABLE.
4. Planning changes after an immutable roadmap require a new roadmap version.
5. Never modify a roadmap that has generated sprints.

#### Ownership

| Role | Responsibility |
|---|---|
| Product Owner | Defines roadmap scope from MVP |
| Tech Lead | Plans work breakdown from epics |

#### Related Folders

- [Section 3.2](#32-mvp) — MVP defines the release scope
- [Section 3.3](#33-epics) — Epics provide capability area breakdown
- [Section 3.8](#38-sprints) — Sprints are generated from roadmaps

---

### 3.8 sprints/

*Source: `docs/sprints/README.md`*

Sprint documents — execution snapshots from the latest roadmap. Only one sprint can be active at a time. Sprints are generated from roadmaps and drive ticket execution.

#### What Goes Here

- Sprint execution documents
- Sprint change logs
- Sprint completion records

#### What Does NOT Go Here

- Roadmap planning → [Section 3.7](#37-roadmap)
- Wave/ticket decomposition → [Section 3.9](#39-tickets)
- Individual work items → [Section 3.10](#310-work-items)

#### Naming Convention

| Type | Pattern | Example |
|---|---|---|
| Sprint | `mvp-nnn-s<nn>.md` | `mvp-001-s01.md` |

#### Work Item Status

Every work item in the Sprint's Work Items table uses the 5-status model defined in [Section 4.2](#42-status-model-summary). No other status is allowed.

#### Status Lifecycle

```
Active → Completed
```

- **Active**: Sprint is in progress; tickets are being worked on.
- **Completed**: All tickets are Done; sprint is closed.

#### Active Sprint Rule

**ONLY ONE active sprint at a time.**

- A sprint becomes active when generated from the roadmap.
- The next sprint CANNOT be generated until the active sprint is completed (all work items Done).
- Sprint completion is recorded in the sprint's Change Log.
- If a sprint is partially completed, deferred items roll to the next sprint.

#### Ownership

| Role | Responsibility |
|---|---|
| Scrum Master | Manages sprint execution |
| Tech Lead | Generates sprints from roadmap |
| Engineers | Execute sprint work items |

#### Related Folders

- [Section 3.7](#37-roadmap) — Sprints are generated from roadmaps
- [Section 3.9](#39-tickets) — Tickets are generated from sprints
- [Section 3.10](#310-work-items) — Work items are the implementation units

---

### 3.9 tickets/

*Source: `docs/tickets/README.md`*

Wave/ticket execution documents — deterministic wave decomposition of active sprints. Tickets are consumed by execution commands and track the granular execution of sprint work.

#### What Goes Here

- Wave ticket documents
- Ticket execution tracking
- Deterministic decomposition of sprint work

#### What Does NOT Go Here

- Sprint definitions → [Section 3.8](#38-sprints)
- Individual work items → [Section 3.10](#310-work-items)
- Roadmap planning → [Section 3.7](#37-roadmap)

#### Naming Convention

| Type | Pattern | Example |
|---|---|---|
| Wave tickets | `<sprint-id>-wave-tickets.md` | `mvp-001-s01-wave-tickets.md` |

#### Ticket Status

Every ticket uses the 5-status model defined in [Section 4.2](#42-status-model-summary). No other status is allowed. Ticket status mirrors the sprint work item status — when a work item changes status, its corresponding ticket changes too.

#### Status Rules

- Status transitions are strictly sequential: `Planned` → `Todo` → `In Progress` → `In Review` → `Done`.
- Never skip a status. Never revert to a previous status.
- `Planned` is the only valid initial status.
- `Done` is the only valid terminal status.
- Ticket status and its corresponding sprint work item status must always agree.
- Status changes are recorded in the Ticket document's Change Log.

#### Ownership

| Role | Responsibility |
|---|---|
| Tech Lead | Generates tickets from sprints |
| Engineers | Execute tickets |

#### Related Folders

- [Section 3.8](#38-sprints) — Tickets are generated from sprints
- [Section 3.10](#310-work-items) — Work items are referenced by tickets

---

### 3.10 work-items/

*Source: `docs/work-items/README.md`*

Work item definitions — the implementation units within sprints. Work items are organized by type in subfolders. Each type has its own naming convention and lifecycle.

#### Subfolders

| Subfolder | Purpose | File Naming | Example |
|---|---|---|---|
| `bugs/` | Bug reports and fixes | `bug-nnn-<slug>.md` | `bug-001-description-field-null.md` |
| `stories/` | User stories | `st-nnn-nnn-<slug>.md` | `st-012-001-user-login-flow.md` |
| `technical-stories/` | Technical implementation stories | `ts-nnn-nnn-<slug>.md` | `ts-012-001-database-migration.md` |
| `tech-debt/` | Technical debt items | `td-nnn-<slug>.md` | `td-002-remove-legacy-auth.md` |
| `chores/` | Maintenance and operational tasks | `ch-nnn-nnn-<slug>.md` | `ch-001-001-update-dependencies.md` |
| `spikes/` | Research and investigation tasks | `sp-nnn-<slug>.md` | `sp-001-evaluate-cache-strategies.md` |
| `planning/` | Planning-related work items | Descriptive kebab-case | `planning-epic-breakdown.md` |

#### General Naming Rules

1. Work item files use their ID prefix in **lowercase**.
2. All filenames use **kebab-case**.
3. The ID prefix identifies the work item type.
4. Sequential numbering within each type.

#### Work Item Status

Every work item uses the 5-status model defined in [Section 4.2](#42-status-model-summary). No other status is allowed.

#### Ownership

| Role | Responsibility |
|---|---|
| Tech Lead | Creates and prioritizes work items |
| Engineers | Implement and update work items |
| QA | Verifies completed work items |

#### Related Folders

- [Section 3.8](#38-sprints) — Work items are assigned to sprints
- [Section 3.9](#39-tickets) — Tickets reference work items
- [Section 3.12](#312-sessions) — Implementation sessions work on work items
- [Section 3.11](#311-reviews) — Reviews validate completed work items

---

#### 3.10.1 bugs/

*Source: `docs/work-items/bugs/README.md`*

Bug reports and fixes. Each bug document describes a defect, its impact, and the resolution.

##### What Goes Here

- Bug reports
- Bug fix documentation
- Root cause analysis

##### What Does NOT Go Here

- Technical debt → [Section 3.10.4](#3104-tech-debt)
- Feature requests → [Section 3.10.2](#3102-stories)
- Chores → [Section 3.10.5](#3105-chores)

##### Naming Convention

| Type | Pattern | Example |
|---|---|---|
| Bug | `bug-nnn-<slug>.md` | `bug-001-description-field-null.md` |

##### Status

Every work item uses the 5-status model defined in [Section 4.2](#42-status-model-summary). Status transitions are strictly sequential. Never skip a status. Never revert to a previous status.

##### Ownership

| Role | Responsibility |
|---|---|
| Engineers | Report and fix bugs |
| QA | Verifies bug fixes |
| Tech Lead | Prioritizes bugs |

##### Related Folders

- [Section 3.10.2](#3102-stories) — User stories that may introduce bugs
- [Section 3.10.4](#3104-tech-debt) — Technical debt that may cause bugs

---

#### 3.10.2 stories/

*Source: `docs/work-items/stories/README.md`*

User stories — work items that deliver user-facing value. Stories describe functionality from the user's perspective.

##### What Goes Here

- User stories
- Acceptance criteria
- Story implementation notes

##### What Does NOT Go Here

- Technical implementation details → [Section 3.10.3](#3103-technical-stories)
- Bug reports → [Section 3.10.1](#3101-bugs)
- Technical debt → [Section 3.10.4](#3104-tech-debt)

##### Naming Convention

| Type | Pattern | Example |
|---|---|---|
| Story | `st-nnn-nnn-<slug>.md` | `st-012-001-user-login-flow.md` |

The double-number pattern (`st-nnn-nnn`) allows grouping stories under epics or themes.

##### Status

Every work item uses the 5-status model defined in [Section 4.2](#42-status-model-summary). Status transitions are strictly sequential. Never skip a status. Never revert to a previous status.

##### Ownership

| Role | Responsibility |
|---|---|
| Product Owner | Defines story requirements |
| Engineers | Implement stories |
| QA | Verifies story acceptance criteria |

##### Related Folders

- [Section 3.10.3](#3103-technical-stories) — Technical implementation stories
- [Section 3.10.1](#3101-bugs) — Bugs that may be found in stories
- [Section 3.3](#33-epics) — Stories are grouped under epics

---

#### 3.10.3 technical-stories/

*Source: `docs/work-items/technical-stories/README.md`*

Technical implementation stories — work items that deliver technical value without direct user-facing impact. Technical stories describe implementation details, infrastructure changes, and system improvements.

##### What Goes Here

- Technical implementation stories
- Infrastructure work items
- System improvement stories
- Non-functional requirement implementations

##### What Does NOT Go Here

- User-facing stories → [Section 3.10.2](#3102-stories)
- Bug reports → [Section 3.10.1](#3101-bugs)
- Technical debt → [Section 3.10.4](#3104-tech-debt)
- Research spikes → [Section 3.10.6](#3106-spikes)

##### Naming Convention

| Type | Pattern | Example |
|---|---|---|
| Technical Story | `ts-nnn-nnn-<slug>.md` | `ts-012-001-database-migration.md` |

The double-number pattern (`ts-nnn-nnn`) allows grouping technical stories under epics or themes.

##### Status

Every work item uses the 5-status model defined in [Section 4.2](#42-status-model-summary). Status transitions are strictly sequential. Never skip a status. Never revert to a previous status.

##### Ownership

| Role | Responsibility |
|---|---|
| Tech Lead | Defines technical stories |
| Engineers | Implement technical stories |
| QA | Verifies technical acceptance criteria |

##### Related Folders

- [Section 3.10.2](#3102-stories) — User stories that may depend on technical stories
- [Section 3.10.6](#3106-spikes) — Spikes that may produce technical stories
- [Section 3.4](#34-architecture) — Architecture docs that inform technical stories

---

#### 3.10.4 tech-debt/

*Source: `docs/work-items/tech-debt/README.md`*

Technical debt items — work items that address accumulated technical compromises. Tech debt items track shortcuts, deprecated patterns, and refactoring needs.

##### What Goes Here

- Technical debt documentation
- Refactoring plans
- Deprecated code tracking
- Migration work items

##### What Does NOT Go Here

- Bug reports → [Section 3.10.1](#3101-bugs)
- Feature stories → [Section 3.10.2](#3102-stories)
- Maintenance chores → [Section 3.10.5](#3105-chores)

##### Naming Convention

| Type | Pattern | Example |
|---|---|---|
| Tech Debt | `td-nnn-<slug>.md` | `td-002-remove-legacy-auth.md` |

##### Status

Every work item uses the 5-status model defined in [Section 4.2](#42-status-model-summary). Status transitions are strictly sequential. Never skip a status. Never revert to a previous status.

##### Ownership

| Role | Responsibility |
|---|---|
| Engineers | Identify and document tech debt |
| Tech Lead | Prioritizes tech debt items |
| Engineers | Resolve tech debt items |

##### Related Folders

- [Section 3.10.5](#3105-chores) — Maintenance tasks
- [Section 3.10.2](#3102-stories) — Feature work that may introduce tech debt
- [Section 3.5](#35-adr) — Architecture decisions that may create tech debt

---

#### 3.10.5 chores/

*Source: `docs/work-items/chores/README.md`*

Maintenance and operational tasks. Chores are work items that don't add user-facing features but keep the project healthy.

##### What Goes Here

- Dependency updates
- Build process maintenance
- Configuration changes
- Operational tasks

##### What Does NOT Go Here

- Feature work → [Section 3.10.2](#3102-stories)
- Bug fixes → [Section 3.10.1](#3101-bugs)
- Technical debt → [Section 3.10.4](#3104-tech-debt)

##### Naming Convention

| Type | Pattern | Example |
|---|---|---|
| Chore | `ch-nnn-nnn-<slug>.md` | `ch-001-001-update-dependencies.md` |

##### Status

Every work item uses the 5-status model defined in [Section 4.2](#42-status-model-summary). Status transitions are strictly sequential. Never skip a status. Never revert to a previous status.

##### Ownership

| Role | Responsibility |
|---|---|
| Engineers | Execute chores |
| Tech Lead | Prioritizes chores |

##### Related Folders

- [Section 3.10.4](#3104-tech-debt) — Technical debt items
- [Section 3.13](#313-operations) — Operational references

---

#### 3.10.6 spikes/

*Source: `docs/work-items/spikes/README.md`*

Research and investigation tasks. Spikes are time-boxed explorations to reduce technical uncertainty before implementation begins.

##### What Goes Here

- Spike definitions and scope
- Spike findings and conclusions
- Technical investigation results

##### What Does NOT Go Here

- Implementation work → [Section 3.10.2](#3102-stories), [Section 3.10.3](#3103-technical-stories)
- Architecture decisions → [Section 3.5](#35-adr)
- Bug investigations → [Section 3.10.1](#3101-bugs)

##### Naming Convention

| Type | Pattern | Example |
|---|---|---|
| Spike | `sp-nnn-<slug>.md` | `sp-001-evaluate-cache-strategies.md` |

##### Status

Every work item uses the 5-status model defined in [Section 4.2](#42-status-model-summary). Status transitions are strictly sequential. Never skip a status. Never revert to a previous status.

##### Ownership

| Role | Responsibility |
|---|---|
| Engineers | Execute spikes |
| Tech Lead | Defines spike scope and time-box |

##### Related Folders

- [Section 3.6](#36-decisions) — Spike conclusions may produce decisions
- [Section 3.5](#35-adr) — Spike findings may lead to ADRs
- [Section 3.12](#312-sessions) — Spike work may use implementation sessions

---

#### 3.10.7 planning/

*Source: `docs/work-items/planning/README.md`*

Planning-related work items. This subfolder contains work items that support the planning process itself.

##### What Goes Here

- Epic breakdown work items
- Planning facilitation tasks
- Cross-cutting planning artifacts

##### What Does NOT Go Here

- Implementation work items → [Section 3.10.2](#3102-stories), [Section 3.10.3](#3103-technical-stories)
- Spike investigations → [Section 3.10.6](#3106-spikes)
- Planning documents → [Section 3.14](#314-planning)

##### Naming Convention

| Type | Pattern | Example |
|---|---|---|
| Planning item | Descriptive kebab-case | `planning-epic-breakdown.md` |

##### Ownership

| Role | Responsibility |
|---|---|
| Tech Lead | Creates planning work items |
| Product Owner | Reviews planning items |

##### Related Folders

- [Section 3.14](#314-planning) — Cross-cutting planning artifacts
- [Section 3.3](#33-epics) — Epic definitions
- [Section 3.7](#37-roadmap) — Roadmap planning

---

### 3.11 reviews/

*Source: `docs/reviews/README.md`*

Review documents — code reviews, architecture reviews, and audit findings. Reviews validate completed work and capture feedback.

#### What Goes Here

- Code review summaries
- Architecture review documents
- Audit findings
- Review outcomes and action items

#### What Does NOT Go Here

- Architecture decisions → [Section 3.5](#35-adr)
- Operational decisions → [Section 3.6](#36-decisions)
- Work item reviews (status updates) → [Section 3.10](#310-work-items)

#### Naming Convention

| Type | Pattern | Example |
|---|---|---|
| Review | `nnn-<review-topic>.md` | `001-auth-module-review.md` |

Sequential numbering with descriptive slugs.

#### Ownership

| Role | Responsibility |
|---|---|
| Reviewers | Create review documents |
| Tech Lead | Reviews findings and assigns actions |
| Engineers | Address review feedback |

#### Related Folders

- [Section 3.6](#36-decisions) — Reviews may produce decisions
- [Section 3.10](#310-work-items) — Reviews validate completed work items
- [Section 3.12](#312-sessions) — Reviews validate session outputs

---

### 3.12 sessions/

*Source: `docs/sessions/README.md`*

Implementation session contexts — ephemeral working documents from implementation sessions. Sessions capture the context, notes, and verification of a specific implementation effort.

#### What Goes Here

- Implementation session folders
- Session context files
- Session notes and verification records

#### Folder Structure

Each session is a folder: `impl-<work-items>-<date>/`

Inside each session folder:

| File | Purpose |
|---|---|
| `context.md` | Session context — what is being implemented and why |
| `notes.md` | Working notes — observations, decisions, blockers |
| `verification.md` | Verification results — how the work was validated |

##### Example

```
impl-ts-012-001-st-012-001-2025-01-15/
├── context.md
├── notes.md
└── verification.md
```

#### Naming Convention

| Type | Pattern | Example |
|---|---|---|
| Session folder | `impl-<work-items>-<date>/` | `impl-ts-012-001-2025-01-15/` |
| Context file | `context.md` | (fixed name) |
| Notes file | `notes.md` | (fixed name) |
| Verification file | `verification.md` | (fixed name) |

All filenames inside sessions are **lowercase**.

#### Status Lifecycle

```
Active → Completed → Archived
```

- **Active**: Session is in progress.
- **Completed**: Session work is done; verification recorded.
- **Archived**: Session is preserved for reference.

#### Ownership

| Role | Responsibility |
|---|---|
| Engineers | Create and maintain session documents |
| Tech Lead | Reviews session outputs |

#### Related Folders

- [Section 3.10](#310-work-items) — Sessions implement work items
- [Section 3.11](#311-reviews) — Sessions are validated by reviews

---

### 3.13 operations/

*Source: `docs/operations/README.md`*

Operational references — deployment guides, exit codes, output conventions, and server setup. This folder contains documents that support the operational side of the project.

#### What Goes Here

- Deployment guides
- Exit code documentation
- Server setup instructions
- Migration guides
- Checklists
- Output conventions

#### What Does NOT Go Here

- Architecture decisions → [Section 3.5](#35-adr)
- Architecture descriptions → [Section 3.4](#34-architecture)
- Product requirements → [Section 3.1](#31-prd)

#### Naming Convention

| Type | Pattern | Example |
|---|---|---|
| Operations doc | Descriptive kebab-case | `deployment.md` |
| Migration guide | Descriptive kebab-case | `migration-guide-v2.md` |
| Checklist | Descriptive kebab-case | `release-checklist.md` |

Use descriptive names that clearly indicate the document's content.

#### Ownership

| Role | Responsibility |
|---|---|
| DevOps | Creates deployment and infrastructure docs |
| Engineers | Contribute operational references |
| Tech Lead | Reviews operational documents |

#### Related Folders

- [Section 3.4](#34-architecture) — Architecture informs operations
- [Section 3.5](#35-adr) — ADRs may affect operational decisions
- [Section 3.12](#312-sessions) — Sessions may reference operational guides

---

### 3.14 planning/

*Source: `docs/planning/README.md`*

Cross-cutting planning artifacts — traceability matrices, transition plans, and other planning documents that span multiple work streams.

#### What Goes Here

- Traceability matrices
- Transition plans
- Cross-cutting planning documents
- Planning analysis artifacts

#### What Does NOT Go Here

- Epic definitions → [Section 3.3](#33-epics)
- Roadmap planning → [Section 3.7](#37-roadmap)
- Work item planning → [Section 3.10.7](#3107-planning)

#### Naming Convention

| Type | Pattern | Example |
|---|---|---|
| Planning doc | Descriptive kebab-case | `traceability-matrix.md` |
| Transition plan | Descriptive kebab-case | `v2-transition-plan.md` |

Use descriptive names that clearly indicate the document's content.

#### Ownership

| Role | Responsibility |
|---|---|
| Tech Lead | Creates cross-cutting planning artifacts |
| Product Owner | Reviews planning alignment with MVP |

#### Related Folders

- [Section 3.3](#33-epics) — Epics inform planning
- [Section 3.7](#37-roadmap) — Roadmaps are informed by planning
- [Section 3.10.7](#3107-planning) — Planning-specific work items

---

### 3.15 manifesto/

*Source: `docs/manifesto/README.md`*

Product manifesto — vision, principles, and non-negotiables. The manifesto defines the product's identity and guiding philosophy.

#### What Goes Here

- Product manifesto
- Vision statement
- Core principles
- Non-negotiable values

#### What Does NOT Go Here

- Product requirements → [Section 3.1](#31-prd)
- Architecture decisions → [Section 3.5](#35-adr)
- MVP definitions → [Section 3.2](#32-mvp)

#### Naming Convention

| Type | Pattern | Example |
|---|---|---|
| Manifesto | `<product>-manifesto.md` | `anvil-manifesto.md` |

#### Ownership

| Role | Responsibility |
|---|---|
| Product Owner | Defines product vision and principles |
| Founders/Leadership | Approve manifesto |

#### Related Folders

- [Section 3.1](#31-prd) — PRDs are guided by the manifesto
- [Section 3.2](#32-mvp) — MVP scope aligns with manifesto principles
- [Section 3.16](#316-specification-corpus) — Shared terminology and definitions

---

### 3.16 specification-corpus/

*Source: `docs/specification-corpus/README.md`*

Specification vocabulary — canonical definitions, lifecycle models, and shared terminology. This folder contains the authoritative references for terms and concepts used across the project.

#### What Goes Here

- Canonical term definitions
- Lifecycle models
- Shared vocabulary
- Glossary
- Specification references

#### What Does NOT Go Here

- Architecture descriptions → [Section 3.4](#34-architecture)
- Product requirements → [Section 3.1](#31-prd)
- Architecture decisions → [Section 3.5](#35-adr)

#### Naming Convention

| Type | Pattern | Example |
|---|---|---|
| Vocabulary | Descriptive kebab-case | `vocabulary.md` |
| Lifecycle model | Descriptive kebab-case | `lifecycle-model.md` |
| Glossary | Descriptive kebab-case | `glossary.md` |

Use descriptive names that clearly indicate the document's content.

#### Ownership

| Role | Responsibility |
|---|---|
| Tech Lead | Maintains canonical definitions |
| All contributors | Reference and use shared terminology |

#### Related Folders

- [Section 3.4](#34-architecture) — Architecture uses specification terms
- [Section 3.15](#315-manifesto) — Manifesto defines product principles
- [Section 3.5](#35-adr) — ADRs may define new terminology

---

## 4. Quick Reference

### 4.1 Naming Conventions Summary

*Source: `docs/README.md` Section B*

#### Strict Rules

1. **All filenames MUST use kebab-case**: lowercase letters, numbers, hyphens between words.
2. **NEVER use full capslock filenames**: WRONG `ANVIL_MANIFESTO.md`, RIGHT `anvil-manifesto.md`.
3. **NEVER use underscores in filenames** (except in document IDs inside file content).
4. **NEVER use PascalCase or camelCase in filenames**.
5. **README.md is the ONLY exception** (conventional uppercase).

#### Work Item File Naming

| Type | Pattern | Example |
|---|---|---|
| Bug | `bug-nnn-<slug>.md` | `bug-002-description-field-null.md` |
| Story | `st-nnn-nnn-<slug>.md` | `st-012-001-user-login-flow.md` |
| Technical Story | `ts-nnn-nnn-<slug>.md` | `ts-012-001-database-migration.md` |
| Tech Debt | `td-nnn-<slug>.md` | `td-002-remove-legacy-auth.md` |
| Chore | `ch-nnn-nnn-<slug>.md` | `ch-001-001-update-dependencies.md` |
| Spike | `sp-nnn-<slug>.md` | `sp-001-evaluate-cache-strategies.md` |
| PRD | `prd-nnn-<product-name>.md` | `prd-001-anvil.md` |
| MVP | `mvp-nnn-<product-name>-<version>.md` | `mvp-001-anvil-v1.md` |
| Epic | `epic-nnn-<capability-area>.md` | `epic-001-user-authentication.md` |
| ADR | `adr-nnn-<decision-topic>.md` | `adr-003-runtime-and-release-lifecycle.md` |
| Decision | `nnn-<decision-topic>.md` | `001-use-eslint-for-linting.md` |
| Roadmap | `mvp-nnn-v<version>.md` | `mvp-001-v1.md` |
| Sprint | `mvp-nnn-s<nn>.md` | `mvp-001-s01.md` |
| Wave tickets | `<sprint-id>-wave-tickets.md` | `mvp-001-s01-wave-tickets.md` |
| Review | `nnn-<review-topic>.md` | `001-auth-module-review.md` |
| Session folder | `impl-<work-items>-<date>/` | `impl-ts-012-001-2025-01-15/` |
| Manifesto | `<product>-manifesto.md` | `anvil-manifesto.md` |

#### Session Files (Fixed Names)

- `context.md`
- `notes.md`
- `verification.md`

### 4.2 Status Model Summary

*Source: `docs/work-items/README.md`, `docs/sprints/README.md`, `docs/tickets/README.md`*

#### Unified 5-Status Model (Work Items, Tickets, Sprints)

Every work item, ticket, and sprint work item entry uses one of the following statuses. No other status is allowed.

```
Planned → Todo → In Progress → In Review → Done
```

| Status | Meaning | Condition |
|---|---|---|
| `Planned` | Default status when assigned to a sprint | Sprint generation (initial state) |
| `Todo` | Will be worked on — the item has entered an execution wave | Item's wave becomes active in the ticket document |
| `In Progress` | Implementation is underway | Assignee starts working on the item |
| `In Review` | Pull Request has been created | PR opened for the item's implementation |
| `Done` | Implementation is merged | PR merged to the target branch |

#### Status Rules

- Status transitions are strictly sequential: `Planned` → `Todo` → `In Progress` → `In Review` → `Done`.
- Never skip a status. Never revert to a previous status.
- `Planned` is the only valid initial status — every work item starts here.
- `Done` is the only valid terminal status — there is no status beyond `Done`.
- Work item status, sprint table status, and ticket status must always agree.
- Status changes are recorded in the relevant Change Log.

#### Other Status Lifecycles

| Document Type | Lifecycle |
|---|---|
| Living documents (PRD, MVP, Epic, Architecture) | `Draft → Review → Approved → Amended` |
| ADR | `Proposed → Accepted → Superseded` |
| Decisions | `Draft → Accepted → Superseded (optional)` |
| Roadmaps | `Draft → Approved → Immutable` |
| Sprints | `Active → Completed` |
| Sessions | `Active → Completed → Archived` |

### 4.3 Documentation Workflow Summary

*Source: `docs/README.md` Section D*

#### End-to-End Flow

```
PRD → MVP → Epic Planning → Roadmap → Sprint → Ticket → Work Items → Sessions
                                                                         ↓
                                                                     Reviews
```

#### Workflow Steps

| Step | Input | Output | Description |
|---|---|---|---|
| PRD | Business requirements | `prd/prd-nnn-*.md` | Defines what the product must do and why |
| MVP | PRD | `mvp/mvp-nnn-*.md` | Defines release scope boundaries from PRD |
| Epic Planning | MVP | `epics/epic-nnn-*.md` | Breaks MVP into capability areas (epics) |
| Roadmap | Epic Planning + MVP | `roadmap/mvp-nnn-v*.md` | Planned work breakdown; one per MVP version |
| Sprint | Roadmap | `sprints/mvp-nnn-s*.md` | Execution snapshot from roadmap; one active at a time |
| Ticket | Sprint | `tickets/<sprint-id>-wave-tickets.md` | Deterministic wave decomposition of active sprint |
| Work Items | Sprint/Ticket | `work-items/<type>/*.md` | Implementation units within sprints |
| Sessions | Work Items | `sessions/impl-*/` | Ephemeral implementation contexts |
| Reviews | Completed work | `reviews/nnn-*.md` | Validation of completed work |

---

## 5. Source Files

Every README.md that was merged into this document:

| # | Path | Section |
|---|---|---|
| 1 | `README.md` | [Section 1](#1-overview) |
| 2 | `docs/README.md` | [Section 2](#2-documentation-structure--governance) |
| 3 | `docs/adr/README.md` | [Section 3.5](#35-adr) |
| 4 | `docs/architecture/README.md` | [Section 3.4](#34-architecture) |
| 5 | `docs/decisions/README.md` | [Section 3.6](#36-decisions) |
| 6 | `docs/epics/README.md` | [Section 3.3](#33-epics) |
| 7 | `docs/manifesto/README.md` | [Section 3.15](#315-manifesto) |
| 8 | `docs/mvp/README.md` | [Section 3.2](#32-mvp) |
| 9 | `docs/operations/README.md` | [Section 3.13](#313-operations) |
| 10 | `docs/planning/README.md` | [Section 3.14](#314-planning) |
| 11 | `docs/prd/README.md` | [Section 3.1](#31-prd) |
| 12 | `docs/reviews/README.md` | [Section 3.11](#311-reviews) |
| 13 | `docs/roadmap/README.md` | [Section 3.7](#37-roadmap) |
| 14 | `docs/sessions/README.md` | [Section 3.12](#312-sessions) |
| 15 | `docs/specification-corpus/README.md` | [Section 3.16](#316-specification-corpus) |
| 16 | `docs/sprints/README.md` | [Section 3.8](#38-sprints) |
| 17 | `docs/tickets/README.md` | [Section 3.9](#39-tickets) |
| 18 | `docs/work-items/README.md` | [Section 3.10](#310-work-items) |
| 19 | `docs/work-items/bugs/README.md` | [Section 3.10.1](#3101-bugs) |
| 20 | `docs/work-items/chores/README.md` | [Section 3.10.5](#3105-chores) |
| 21 | `docs/work-items/planning/README.md` | [Section 3.10.7](#3107-planning) |
| 22 | `docs/work-items/spikes/README.md` | [Section 3.10.6](#3106-spikes) |
| 23 | `docs/work-items/stories/README.md` | [Section 3.10.2](#3102-stories) |
| 24 | `docs/work-items/tech-debt/README.md` | [Section 3.10.4](#3104-tech-debt) |
| 25 | `docs/work-items/technical-stories/README.md` | [Section 3.10.3](#3103-technical-stories) |
