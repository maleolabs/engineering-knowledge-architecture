# Engineering Workflow Guide

> Anchor EKA: primary onboarding — the mental model behind EKA: how engineering knowledge is produced, organized, validated, projected, exchanged, and consumed (Core v1.1 §8.1, EKA 11).
> Convention document, not an artifact (no `type`/`id`). Companion to [README.md](README.md) (what lives where), [lifecycle.md](lifecycle.md) (the concise lifecycle reference), [operating/protocol.md](operating/protocol.md) (how state moves).

## 1. Purpose

This guide is the starting point for working with Engineering Knowledge Architecture (EKA). It does not teach document templates — there are none to memorize. It establishes the **mental model**: what engineering knowledge is, where it sits, how much authority it carries, and how it flows through a repository.

Three ideas carry the whole model:

1. **EKA models engineering knowledge, not documentation traditions.** The unit is the Artifact: a stable identity, owned state, content, and relationships. Folders and file types are just one representation of that model.
2. **Folders and Markdown are one representation.** EKA is serialization-independent (Core v1.1 §12): Git + Markdown is what this repository uses, but the model survives other storage — databases, graph stores, future platforms.
3. **Methodologies are conventions.** Scrum, Kanban, and Shape Up are convention layers over EKA. They constrain how you work; they never define what knowledge is.

The rest of this guide walks the model in one page, expands each lifecycle step, then shows how users, AI, and methodologies fit into it.

## 2. Engineering Knowledge in One Page

The unit of EKA is the **Artifact** (Core v1.1 §3): an engineering knowledge entity with four parts:

| Part | What it is |
|---|---|
| **Identity** | `(namespace, type, id, instance-version)` — permanent and immutable; everything references knowledge by Identity, never by file location |
| **State Vector** | the State Domains the artifact owns (content, execution, planning, container, existence); each field has a single writer and transitions forward only |
| **Content** | the semantic payload — what the knowledge actually says |
| **Relationships** | explicit links to other artifacts by Identity: `derives-from`, `depends-on`, `supersedes`, `amends`, `validates` |

Three orthogonal axes give every artifact its position and movement (Core v1.1 §8.1):

```
   WHERE IT SITS              ITS AUTHORITY               HOW IT FLOWS
   Engineering Domain         Knowledge Stratum          Lifecycle
   ─────────────────          ──────────────────         ──────────────────
   Discovery                  stratum 1 (highest)        Produce   (born)
   Architecture               stratum 2                  Organize  (shaped)
   Planning                   stratum 3                  Validate  (checked)
   Execution                  stratum 4                  Project   (seen)
   Operations                 stratum 5 (lowest)         Exchange  (moved)
                                                         Consume   (used)
```

- **Engineering Domain — where knowledge sits.** One of five canonical categories, derived from the artifact's token family and Knowledge Dimension. An artifact is born in its domain and never changes it.
- **Knowledge Stratum — its authority.** A fixed position in the strict order **Discovery → Architecture → Planning → Execution → Operations**. Always derived from the domain, never declared by an artifact. Lower-stratum knowledge must not contradict higher-stratum knowledge in force (Stratum Authority Invariant, §8.1).
- **Lifecycle — how it flows.** The six-step movement of State over time: **Produce → Organize → Validate → Project → Exchange → Consume**. The domain never changes while the artifact moves through its lifecycle.

Together: the domain says *what kind of knowledge it is*, the stratum says *how much authority it carries*, and the lifecycle position says *where it is right now*. The concise reference is [lifecycle.md](lifecycle.md); the rest of this guide expands each step with what happens, what you do, and how the CLI and AI participate.

## 3. The Engineering Knowledge Lifecycle

### 3.1 Produce — create knowledge

**What happens.** Knowledge enters the repository at the point of creation, in two kinds:

- **Ephemeral** — session notes (`ses-`), spike investigations (`spk-`), research findings (`fnd-`): captured while the work happens, low ceremony, usually drafts.
- **Durable** — requirements (`req-`), architecture descriptions (`arc-`), decisions (`adr-`, `dec-`), specifications (`spec-`), standards (`std-`): started as drafts and matured through governance.

**Distillation** (Core v1.1 §11.4) is the bridge from ephemeral to durable: findings from sessions, spikes, and reviews are distilled into the durable artifact they affect — a direction becomes a decision, proven procedures become a runbook. Archiving ephemeral knowledge without distilling it first violates the protocol.

**What you do.** Write. Capture findings in the moment, in the folder of the knowledge dimension they belong to. When a session or review produces a durable insight, create the durable artifact it implies. Produce is a writing step — no ceremony, no templates.

**CLI participation.** None. Producing is plain editing; the CLI never writes content for you.

**AI participation.** Drafting and distillation support: AI can draft requirements and decisions from session notes, summarize research findings, and suggest what ephemeral knowledge should be distilled — and into which artifact type. The human stays the author; AI accelerates the capture.

### 3.2 Organize — classify knowledge

**What happens.** Every artifact is classified so it can be found, related, and validated:

- **Identity frontmatter** — `namespace`, `type`, `id`, `instance-version`, `revision`; the filename is only a projection of Identity.
- **Classification** — `dimension` must equal the home folder (R6); the Engineering Domain is derived from the token family. Declaring `domain:` explicitly is optional — if declared, it must match the derived home domain (R11).
- **Relationships to higher strata** — `derives-from` / `depends-on` chains give stratification traceability (R10): a non-Discovery artifact must reach a strictly higher stratum, directly or transitively. `supersedes` / `amends` may only target the same or a lower stratum — never upward (R12).
- **Change-log** — every state transition recorded by its single owner; the last entry must equal the current value (R7).

**What you do.** Fill in frontmatter truthfully. Let the domain derive from the token — you rarely need to declare it. Wire `derives-from` / `depends-on` toward the higher-stratum knowledge your artifact builds on: a specification derives from requirements; a work item depends on the plan. Record transitions as they happen.

**CLI participation.** `eka validate` gives structural feedback on frontmatter, classification, relationships, and change-log consistency — the mechanical check of the Organize step.

**AI participation.** Frontmatter generation and relationship suggestions: AI can produce identity fields from content, propose `derives-from` / `depends-on` targets, and flag classification mismatches before validation runs.

### 3.3 Validate — the conformance gate

**What happens.** Before commit, the mechanical checklist **R0–R12** runs over every artifact:

| Group | Rules | Verdict |
|---|---|---|
| Structural | R0 (frontmatter, artifact rule) | blocking |
| Exchange (Exchange Spec §14.2) | R1–R9 (identity, filename, state, ownership, references, dimension, change-log, projections, well-formedness) | blocking |
| Domain-aware (Core v1.1 §8.1) | R10 stratification traceability | **warning** |
| Domain-aware (Core v1.1 §8.1) | R11 domain coherence, R12 cross-stratum supersession prohibition | blocking |

**Draft tolerance:** `content-state: draft` relaxes the rules — dangling references are warnings (R5), and draft knowledge artifacts are exempt from R10. Drafts are the sandbox of the lifecycle; approval moves them into the governed world.

**What you do.** Run `eka validate` before committing. Treat warnings as quality signals, errors as blockers. When a lower-stratum artifact contradicts higher-stratum knowledge in force, fix the lower one — never the higher one (Stratum Authority Invariant).

**CLI participation.** `eka validate` is the canonical mechanical implementation of R0–R12: deterministic output, exit codes 0/1/2, warnings never affect the exit code. Run it locally before commit and in CI to block non-conformant merges.

**AI participation.** AI-written artifacts pass the same gate. Enforcement mechanisms vary by implementation; the invariants are identical (P16) — an AI that edits markdown directly still produces artifacts that must validate.

### 3.4 Project — read-only views

**What happens.** `eka view` projects the repository without writing anything: the **domain projections** — `discovery`, `architecture`, `planning`, `execution`, `operations` — plus the **ticket** projection. Each domain projection is a read-only view over the artifacts of one Engineering Domain (Core v1.1 §8.1), derived from the model (identity, state, relationships), never from markdown text. `execution` shows the active container's tickets and work items by Execution State; `planning` groups scope/epic/plan/traceability artifacts by Planning State and phase; `architecture` groups decisions, specifications, and standards by Content State; `ticket <id>` projects one ticket's status from its owner work item.

**What you do.** Read projections to see the current picture — what is planned, in progress, approved, done. Never edit a projection: projections are refreshed on read, not edited (P6). A projection has no State of its own and never becomes a writer.

**CLI participation.** `eka view execution` / `eka view planning` / `eka view architecture` / `eka view discovery` / `eka view operations` / `eka view ticket <id>` — the Knowledge Projection Engine: conformance-gated, deterministic, read-only. The methodology names still work: `eka view sprint` and `eka view wave` are CLI-level aliases that resolve to the `execution` projection with identical output — the same convention-layer philosophy as the Representation Alias Registry: methodology terms stay convention-layer names, the projection model stays canonical.

**AI participation.** Projections give AI structured context: deterministic, relationship-derived views of engineering knowledge by domain, free of markdown noise — a reliable input for agents that need to know the current state of work.

### 3.5 Exchange — portable knowledge packages

**What happens.** Import/export moves knowledge between repositories without loss (P13). Export builds a package carrying Identity, State (full change-log), Content, Relationships, and Classification — including each unit's Engineering Domain. Import integrates **atomically** with a **conservative merge**: new artifacts are written, identical duplicates are skipped, any difference conflicts and aborts.

**What you do.** Export to share knowledge with another repository or team; import to receive it. The classification round-trips exactly: `dimension` and the optional `domain` come back unchanged, and the Engineering Domain is derived at import from the token family — never taken as an opaque value — so the domain mapping cannot drift between systems.

**CLI participation.** `eka export` (deterministic, validated before export — a non-conformant repository exports nothing) and `eka import` (atomic, conservative, rollback on failure).

**AI participation.** Exchange Packages are consumable knowledge context: AI can read a package's units as structured knowledge — identity, state, content, relationships, classification — without parsing repository layout.

### 3.6 Consume — knowledge in use

**What happens.** Knowledge is read and executed by humans and agents:

- **Runbooks (`run-`)** — executed as procedures during operations; the most actionable stratum.
- **Specifications and ADRs (`spec-`, `adr-`, `arc-`)** — read as the authority for what is built and why; they outrank lower-stratum content that contradicts them.
- **Tickets (`tkt-`) and work items** — read as the current execution plan: what to do, in what order, in what state.

**What you do.** Read by Identity and by domain, never by location. Execute runbooks; honor higher-stratum knowledge. When a runbook contradicts its approved specification, fix the runbook (lower stratum) — never override the specification.

**CLI participation.** The CLI itself consumes the model on every command: `validate` reads every artifact, `view` reads state and relationships, `export` reads the whole repository. Knowledge in EKA is queryable and operable, not just storable.

**AI participation.** AI consumes the same knowledge through the same model — frontmatter, state, relationships — and can execute commands whose content is governed by protocol.

## 4. How Users Work with EKA

A typical session, end to end:

1. **Initialize** — `eka init my-project` analyzes the workspace, generates the repository from the Reference Skeleton, and validates the result.
2. **Write your first artifacts** — a requirement under `docs/requirements/`, a decision under `docs/decisions/`, a runbook under `docs/operations/`, following the conventions in [docs/README.md](README.md): identity frontmatter, classification, relationships.
3. **Validate before commit** — `eka validate` gives the verdict; fix errors, consider warnings.
4. **See the execution picture** — `eka view execution` shows the active container's work by execution state (`eka view sprint` and `eka view wave` are aliases for it); `eka view ticket <id>` shows one ticket's projected status.
5. **Exchange when needed** — `eka export` produces a portable package; another repository receives it with `eka import`.

**Adopting into an existing project** is the same command: `eka init` detects the workspace (README, `docs/`, existing EKA markers), reuses what already exists, and only generates what is missing. It is idempotent and never overwrites user content — running it twice changes nothing.

You never manipulate the model directly: you edit Markdown, and the CLI keeps the repository conformant and queryable.

## 5. How AI Participates

AI participates in every lifecycle step, through the same model:

- **Same interface, same conventions.** AI reads and writes Markdown with the same frontmatter, state fields, and relationships as humans. There is no separate AI format and no separate AI lane.
- **Structured context.** Projections (`eka view`) give AI deterministic, relationship-derived views of engineering knowledge — per-domain views (`discovery`, `architecture`, `planning`, `execution`, `operations`) and per-ticket status; Exchange Packages give AI portable, self-contained knowledge units.
- **Distillation support.** AI can propose and draft distillation — turning session notes into a decision, findings into a runbook — but the distillation obligation and the governance channel remain.
- **Validation is never bypassed.** Whatever produces an artifact — human or AI — the R0–R12 gate applies identically (P16). AI is a faster author and reader, not an exception.

## 6. How Methodology Fits

EKA is methodology-agnostic by design (Core v1.1 §8.1):

- **Methodologies are convention layers.** Scrum, Kanban, Shape Up, and similar are convention layers over EKA — never part of the Core standard. They may map onto tokens and domains through Representation Aliases.
- **The Representation Alias Registry** maps methodology terms onto canonical tokens + Engineering Domain. Aliases are never frontmatter values and never artifact types of their own:

| Representation Alias | Canonical token | Engineering Domain |
|---|---|---|
| PRD | `req-` | Discovery |
| ADR / RFC | `adr-` | Architecture |
| Epic | `epc-` | Planning |
| Initiative | `scp-` | Planning |
| Sprint / Iteration | `ctr-` | Execution |
| Ticket | `tkt-` | Execution |
| Incident | `bug-` | Execution |
| Release | `rel-` | Operations |
| Runbook | `run-` | Operations |

A team running Scrum still writes a PRD as a `req-` artifact in Discovery and a Sprint as a `ctr-` container in Execution. The methodology decides cadence, rituals, and how work is scheduled; EKA decides what the knowledge is, where it sits, and how it is governed. **Methodologies constrain HOW you work; EKA models WHAT you know.**

## 7. The Mental Model Shift

The practical shift, in one sentence: **stop thinking in folders and document types; think in domains, strata, and lifecycle.**

| Old habit | EKA model |
|---|---|
| "The PRD lives in the requirements folder" | The PRD is a `req-` artifact — Discovery domain, stratum 1 |
| "The spec is outdated" | The spec is Approved content in force — higher stratum; lower-stratum knowledge must conform to it |
| "Move the doc to the archive folder" | Move the artifact through its lifecycle: supersede, archive as a Record |
| "Update the sprint board" | Projections refresh on read — edit the owner state, never the view |

Markdown is one representation of the model, not the model:

- The **filename** is a projection of Identity — the Identity lives in frontmatter.
- The **folder** is a projection of classification — `dimension` must equal the folder (R6), but the domain is derived from the token.
- **Projection tables** (container work-item lists, ticket status) are projections of state — refresh on read, never edited.

Because the model is separate from the representation, knowledge survives: references by Identity do not break on rename, classification can be reorganized without touching Identity (P15), and knowledge moves between systems without loss (P13).

## 8. Where to Go Next

- [README.md](README.md) — serialization conventions: identity, state, phase, relationships, classification, projections; the navigation map of `docs/`.
- [operating/protocol.md](operating/protocol.md) — the Operating Manual: ordering chain, State Domains, gates, change-log, distillation obligations.
- [lifecycle.md](lifecycle.md) — the concise lifecycle reference this guide expands.
- [../reference/cli.md](../reference/cli.md) — the CLI reference: `init`, `validate`, `export`, `import`, `view`, `version`, `completion`.
