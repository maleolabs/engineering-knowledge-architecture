# docs/ — Source of Truth for EKA Serialization

> Anchor EKA: overall serialization — Knowledge Layer (content), Operating Layer (state), Exchange Layer (validation & transfer).
> Standard: EKA v1.1, dated 2026-08-05.

## Status of This Folder

This folder is the **EKA v1.1 serialization** (EKA projected to Git + Markdown), **not EKA architecture itself**. The architecture — three layers (Knowledge/Operating/Exchange), five State Domains, the Protocol, and the five Engineering Domains (Core v1.1 §8.1) — is conceptual; this folder is only its file representation. In other words: changing the folder does not change EKA, and EKA can be serialized to other media without losing meaning.

## Navigation Map (17 Entries)

| Entry | Anchor EKA | Content | Engineering Domain |
|---|---|---|---|
| [README.md](.) (this file) | serialization | source of truth + convention summary | — |
| [workflow-guide.md](workflow-guide.md) | onboarding | Engineering Workflow Guide — primary onboarding: mental model, lifecycle, CLI & AI participation | — |
| [intent/](intent/) | intent dimension | `vis-` Vision/Manifesto, `str-` Strategy | Discovery |
| [requirements/](requirements/) | requirements dimension | `req-` Requirement | Discovery |
| [architecture/](architecture/) | architecture dimension | `arc-` Architecture Description | Architecture |
| [decisions/](decisions/) | decisions dimension | `adr-` ADR, `dec-` Decision Record | Architecture |
| [specifications/](specifications/) | specifications dimension | `spec-` Specification | Architecture |
| [standards/](standards/) | standards dimension | `std-` Standard/Guideline | Architecture |
| [operations/](operations/) | operations dimension | `run-` Runbook | Operations |
| [quality/](quality/) | quality dimension | `rvw-` Review | Execution |
| [planning/](planning/) | planning dimension | `scp-`, `epc-`, `plan-`, `trc-` | Planning |
| [records/](records/) | records dimension | `rel-` Release Record | Operations |
| [research/](research/) | research dimension | `fnd-` Research Finding (EKA 14.1) | Discovery |
| [vocabulary/](vocabulary/) | vocabulary dimension | `gls-` Glossary/Term | Architecture |
| [lifecycle.md](lifecycle.md) | lifecycle | Engineering Knowledge Lifecycle (produce → consume) | — |
| [operating/](operating/) | Operating Layer | state, protocol, work items, containers, sessions, projections | Execution |
| [exchange/](exchange/) | Exchange Layer | validation, transfer | — |

## Engineering Domains

Every folder (and every artifact type) belongs to exactly one of the five canonical **Engineering Domains** (Core v1.1 §8.1) — the stratum-aligned category of engineering knowledge it holds:

| Engineering Domain | Stratum | Folders / dimensions | Token families |
|---|---|---|---|
| **Discovery** | 1 (highest authority) | intent, requirements, research | vis-, str-, req-, fnd- |
| **Architecture** | 2 | architecture, decisions, specifications, standards, vocabulary | arc-, adr-, dec-, spec-, std-, gls- |
| **Planning** | 3 | planning | scp-, epc-, plan-, trc- |
| **Execution** | 4 | quality (+ operating/) | rvw-, ctr-, tkt-, sto-, ts-, bug-, td-, ch-, spk-, ses- |
| **Operations** | 5 | operations, records | run-, rel- |

The domains form a strict authority chain: **Discovery → Architecture → Planning → Execution → Operations**. **Stratum Authority Invariant:** knowledge in a lower stratum must not contradict knowledge in a higher stratum that is in force — resolve contradictions by changing the lower-stratum knowledge (new instance + relationship, forward-only), never by superseding or amending upward.

Methodology note: **PRD**, **ADR/RFC**, **Epic**, **Initiative**, **Sprint/Iteration**, **Ticket**, **Release**, **Incident**, **Runbook** are **Representation Aliases** — methodology terms mapped onto a canonical token + Engineering Domain (e.g. PRD → `req-`, Sprint → `ctr-`). They are never frontmatter values and never artifact types of their own. The full alias catalog and its extension governance live in the [Representation Alias Registry](../../standard/representation-alias-registry-v1.1.md). Methodologies (Scrum, Kanban, Shape Up, …) are **convention layers over EKA**, not part of the Core standard.

New to EKA? Start with the [Engineering Workflow Guide](workflow-guide.md) — the primary onboarding document. How knowledge flows through these domains: [lifecycle.md](lifecycle.md).

## Serialization Conventions Summary

### Identity
- Identity = `(namespace, type, id, instance-version)`; **lives in frontmatter**, filename is only a projection.
- Fields: `namespace` (default: product name), `type` (token, must match filename), `id` (kebab-case, unique within (namespace, type)), `instance-version` (int, default 1; required for `scp-`/`plan-`), `revision` (int, default 1; only for content edits).
- Filename pattern: `<type-token>-<id>.md`; for versioned types (`scp-`, `plan-`): `<type-token>-<id>-v<instance-version>.md` (always, including v1). The `-v<nn>` suffix is forbidden for other types. Details: [planning/README.md](planning/README.md).

### State
- Five owned State Domains: Content State, Execution State, Planning State, Container State, Existence State. Absence of a field = not applicable (N/A) for that type.
- Each state field is written by only **one owner** (single-writer, P6); projections never write state.
- Transitions are **forward-only**; every transition is recorded in `change-log`.
- Values and transition rules: [operating/protocol.md](operating/protocol.md).

### Phase
- `phase` is a **context attribute** on `scp-`/`plan-` artifacts, **not a folder**. Values: `discovery | mvp | milestone | release | growth | maturity | sunset`.
- Phase changes are recorded in `change-log` with `domain: phase`.

### Relationships
- Fields: `amends`, `supersedes`, `derives-from`, `depends-on`, `validates` — reference lists in `<type-token>:<id>[:<instance-version>]` form (cross-namespace: `<namespace>/<type-token>:<id>`).
- Written only on the referring artifact; references must resolve (see validation).

### Classification
- `dimension` (primary) + `dimensions-secondary` (list) — for knowledge artifacts; `dimension` must equal the home folder. Operational artifacts (work items) use `dimension` informationally only; `ctr-`/`tkt-`/`ses-` carry no `dimension`.

### Projections
- Container and ticket tables carry the header `> Generated — State Projection. Do NOT edit state here; refresh on read.`
- Projections are refreshed on read, never edited manually.

## Artifacts vs Convention Documents

- **Artifacts:** files whose frontmatter contains `type` **and** `id`. Examples: `req-login.md`, `ctr-wave-1.md`, `plan-release-1-v1.md`.
- **Convention documents:** files with neither — all READMEs, `operating/protocol.md`, `exchange/validation.md`, `exchange/transfer.md`. Convention documents explain the rules and carry no state.

## Validation

Before committing a new artifact or a state change, run the mechanical checklist in [exchange/validation.md](exchange/validation.md). For import/export between repositories: [exchange/transfer.md](exchange/transfer.md).
