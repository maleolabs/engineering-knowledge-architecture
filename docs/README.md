# Documentation Structure

This document is the single source of truth for the documentation structure, naming conventions, workflow, and governance rules.

---

## A. Documentation Structure

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

---

## B. Naming Conventions

### Strict Rules

1. **All filenames MUST use kebab-case**: lowercase letters, numbers, hyphens between words.
2. **NEVER use full capslock filenames**: WRONG `ANVIL_MANIFESTO.md`, RIGHT `anvil-manifesto.md`.
3. **NEVER use underscores in filenames** (except in document IDs inside file content).
4. **NEVER use PascalCase or camelCase in filenames**.
5. **README.md is the ONLY exception** (conventional uppercase).

### Examples

| Correct | Incorrect |
|---|---|
| `mvp-001-anvil-v1.md` | `MVP_001_ANVIL_V1.md` |
| `adr-003-runtime-and-release-lifecycle.md` | `ADR_003_Runtime_and_Release_Lifecycle.md` |
| `bug-001-description-field-null.md` | `BUG_001_Description_Field_Null.md` |
| `anvil-manifesto.md` | `ANVIL_MANIFESTO.md` |
| `traceability-matrix.md` | `Traceability_Matrix.md` |

### Work Item File Naming

Work item files use their ID prefix in lowercase:

| Type | Pattern | Example |
|---|---|---|
| Bug | `bug-nnn-<slug>.md` | `bug-002-description-field-null.md` |
| Story | `st-nnn-nnn-<slug>.md` | `st-012-001-user-login-flow.md` |
| Technical Story | `ts-nnn-nnn-<slug>.md` | `ts-012-001-database-migration.md` |
| Tech Debt | `td-nnn-<slug>.md` | `td-002-remove-legacy-auth.md` |
| Chore | `ch-nnn-nnn-<slug>.md` | `ch-001-001-update-dependencies.md` |
| Spike | `sp-nnn-<slug>.md` | `sp-001-evaluate-cache-strategies.md` |

### Session File Naming

Session folders: `impl-<work-items>-<date>/`

Context files inside sessions (all lowercase):
- `context.md`
- `notes.md`
- `verification.md`

---

## C. Folder Placement Rules

### Strict Rules

1. **EVERY file MUST be inside a specific subfolder.** No loose files in `docs/` root.
2. If a document doesn't fit any existing folder, create a new folder with a `README.md` or place it in the most relevant existing folder.
3. Migration guides, checklists, and project-level docs go in the most relevant subfolder:
   - Migration guides → `operations/`
   - Project structure docs → `architecture/`
   - Checklists → `operations/` or the relevant folder

### Common Mistakes

| Mistake | Correct Placement |
|---|---|
| `docs/project-structure.md` (loose file) | `docs/architecture/project-structure.md` |
| `docs/migration-guide-v2.md` (loose file) | `docs/operations/migration-guide-v2.md` |
| `docs/ANVIL_MANIFESTO.md` (capslock) | `docs/manifesto/anvil-manifesto.md` |
| `docs/CONTEXT.md` (capslock) | `docs/sessions/<session>/context.md` |

---

## D. Documentation Workflow

Documents flow through the system in this order:

```
PRD → MVP → Epic Planning → Roadmap → Sprint → Ticket → Work Items → Sessions
                                                                         ↓
                                                                     Reviews
```

### Workflow Steps

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

## E. Active Sprint Rule

**ONLY ONE active sprint at a time.**

- A sprint becomes active when generated from the roadmap.
- The next sprint CANNOT be generated until the active sprint is completed (all tickets Verified).
- Sprint completion is recorded in the sprint's Change Log.
- If a sprint is partially completed, deferred items roll to the next sprint.

### Sprint Lifecycle

```
Active → Completed
```

- **Active**: Sprint is in progress; tickets are being worked on.
- **Completed**: All tickets are Verified; sprint is closed.

---

## F. Roadmap Rules

1. **Roadmaps are derived from epic planning + actual work items registered in epic planning.**
2. **Roadmaps must align with the active MVP.**
3. **Once a sprint is generated from a roadmap, the roadmap becomes IMMUTABLE.**
4. Planning changes after an immutable roadmap require a new roadmap version.
5. Never modify a roadmap that has generated sprints.

### Roadmap Lifecycle

```
Draft → Approved → Immutable
```

- **Draft**: Roadmap is being planned.
- **Approved**: Roadmap is finalized and ready for sprint generation.
- **Immutable**: Sprints have been generated; roadmap cannot be changed.

---

## G. Document Status Lifecycle

### Living Documents (PRD, MVP, Epic, Architecture)

```
Draft → Review → Approved → Amended
```

- **Draft**: Document is being written.
- **Review**: Document is under review.
- **Approved**: Document is finalized and accepted.
- **Amended**: Document has been updated via a separate amendment document.

### Architecture Decision Records (ADR)

```
Proposed → Accepted → Superseded
```

- **Proposed**: ADR is being discussed.
- **Accepted**: ADR is finalized and binding.
- **Superseded**: ADR has been replaced by a new ADR.

### Work Items

Per-type lifecycles — see `work-items/README.md` for details.

### Roadmaps

```
Draft → Approved → Immutable
```

### Sprints

```
Active → Completed
```

### Sessions

```
Active → Completed → Archived
```

- **Active**: Session is in progress.
- **Completed**: Session work is done.
- **Archived**: Session is preserved for reference.

---

## H. Metadata Convention

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

---

## I. Adding New Folders

If you need a new documentation folder:

1. Create the folder under `docs/`.
2. Add a `README.md` explaining purpose, naming conventions, ownership, and related folders.
3. Update this document's structure table (Section A).
4. Follow kebab-case naming for the folder itself.

---

## J. References

- [Repo-level README](../README.md) — what this template is and how to use it
- [ADR Template](adr/README.md) — Architecture Decision Record format
- [Work Items](work-items/README.md) — work item types and lifecycles
- [Sessions](sessions/README.md) — implementation session structure
