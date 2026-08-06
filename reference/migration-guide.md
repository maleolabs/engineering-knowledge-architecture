# Migration Guide — Legacy Structure → EKA v1.0

Migration guide from the legacy documentation structure to the EKA v1.0 serialization. All changes are **breaking by design**: legacy consumers must stop working so they cannot read Identity/State from location (EKA 6.4, P9).

## Part A — Complete mapping table (legacy → new)

| Legacy element | New home | Action | Compatibility | Rationale |
|---|---|---|---|---|
| `docs/README.md` | `skeleton/docs/README.md` | transform | breaking | New serialization conventions |
| `docs/manifesto/` | `docs/intent/` as `vis-` | move + rename | breaking | Product Intent dimension |
| `docs/prd/` | `docs/requirements/` as `req-`; amendments as instances with `amends` | move + rename + transform | breaking | Requirements dimension |
| `docs/mvp/` | `docs/planning/` as `scp-…-v<n>` with `phase: mvp` | transform | breaking | Phase as context (EKA 11.2); Identity decoupled from phase (P3); `mvp-nnn` collisions resolved (EKA 6.4) |
| `docs/epics/` | `docs/planning/` as `epc-` | move + rename | breaking | Planning Knowledge dimension |
| `docs/architecture/` | `docs/architecture/` as `arc-` | rename | breaking | Type token |
| `docs/adr/` | `docs/decisions/` as `adr-` | move + rename | breaking | Single Decisions dimension |
| `docs/decisions/` | `docs/decisions/` as `dec-` | move + rename | breaking | Single Decisions dimension |
| `docs/roadmap/` | `docs/planning/` as `plan-…-v<n>` | move + rename + transform | breaking | Misnomer corrected; Planning State preserved: approved→approved, immutable→immutable |
| `docs/sprints/` | `docs/operating/containers/` as `ctr-` | move + rename + transform | breaking | Execution Container (EKA 10); Container State; tables = projections (EKA 7.4) |
| `docs/tickets/` | `docs/operating/projections/` as `tkt-` | move + rename + transform | breaking | Empty State Vector, `derives-from` |
| `docs/work-items/{stories,technical-stories,bugs,tech-debt,chores,spikes}/` | `docs/operating/work-items/<subtype>/` as `sto-`, `ts-`, `bug-`, `td-`, `ch-`, `spk-` | move + rename | breaking | Single-writer Execution State |
| `docs/work-items/planning/` | deprecated | — | breaking | Catch-all dissolved; content → `planning/`, planning work → the correct work item type |
| `docs/reviews/` | `docs/quality/` as `rvw-` with `validates` | move + rename | breaking | Governance & Quality dimension |
| `docs/sessions/` | `docs/operating/sessions/` as `ses-` | move + rename | breaking | Existence State; Distillation mandatory (EKA 11.4) |
| `docs/operations/` | split: `docs/operations/` as `run-` (procedures) + `docs/standards/` as `std-` (conventions) | — | breaking | Operational vs Standards (EKA 8) |
| `docs/planning/` | `docs/planning/` as `trc-`/`plan-` | move + rename | breaking | Catch-all dissolved |
| `docs/specification-corpus/` | `docs/vocabulary/` as `gls-`; actual specs → `docs/specifications/` as `spec-` | move + rename | breaking | Misnomer corrected; Vocabulary ≠ Specifications (EKA 8) |
| `documentation-guide.md` | split: `reference-architecture.md` + `skeleton/docs/README.md` + `operating/protocol.md` | — | breaking | Standard ≠ serialization (EKA 1.3) |
| `README.md` | new root `README.md` | — | breaking | New repository identity |
| 3-way status sync | single-writer frontmatter + projections | — | breaking | P6, 7.4 |
| metadata table (Status/Author/…) | frontmatter | — | breaking | D2.8: status → state domains; version split into instance-version + revision |

## Part B — Step-by-step migration strategy

1. **Snapshot** — Freeze active sprints, record the final status of all artifacts (work items, containers, plans), commit a baseline before migration so every condition can be restored from git history.
2. **Create the new structure** — Copy `skeleton/docs/` into the project; set `namespace` on all artifacts per project; read `skeleton/docs/README.md` as the source of truth for the structure.
3. **Migrate knowledge artifacts in dependency order**:
   - `intent` (vision/manifesto → `vis-`, strategy → `str-`);
   - `requirements` (amendments → new instances + `amends` relationship);
   - `decisions` (`adr/` + `decisions/` merged into `decisions/`; status mapping: Draft→draft, Review→review, Approved→approved, Accepted→accepted, Superseded→superseded, Amended→amended);
   - `architecture` (`arc-`);
   - `specifications` (new — `spec-`);
   - `standards` (extracted from operations — `std-`);
   - `operations` (procedures only — `run-`);
   - `quality` (`rvw-` with `validates`);
   - `vocabulary` (`gls-`);
   - `planning` (`scp-` with `phase`, `epc-`, `plan-` with `planning-state`, `trc-`).
4. **Migrate operating artifacts**:
   - **Work items first** — content + `execution-state` from legacy status; `change-log` rebuilt from the legacy Change Log; the work item file is now the single writer;
   - **Containers** (`ctr-`) — `container-state: active/completed` from sprint state; container tables = regenerated projections;
   - **Tickets** (`tkt-`) — + `derives-from` + projected status;
   - **Sessions** (`ses-`) — with mandatory Distillation before Archived (EKA 11.4).
5. **Rebuild relationships** — Scan legacy content; encode `amends`/`supersedes`/`derives-from`/`depends-on`/`validates` in frontmatter; references always by Identity.
6. **Validate** — Run the `exchange/validation.md` checklist (9 rules, ADR-006); fix all findings before proceeding.
7. **Archive legacy elements** — Restore from git history when needed; **do not** migrate duplicated legacy status as authoritative — single-writer (P6).
8. **Update tooling** — Legacy status/location consumers migrate to Identity/State in frontmatter; legacy globs (`mvp-*`, `sp-*`, etc.) are deliberately broken (see `breaking-changes.md`).
