# Engineering Knowledge Lifecycle

> Anchor EKA: conceptual lifecycle (EKA 11) — how knowledge is born, organized, validated, projected, exchanged, and consumed.
> Convention document, not an artifact (no `type`/`id`). Companion to [README.md](README.md) (what lives where) and [operating/protocol.md](operating/protocol.md) (how state moves).

## Purpose

The lifecycle is **how knowledge flows**; the Engineering Domain is **where it sits**; the Knowledge Stratum is **its authority**. These are three orthogonal axes (Core v1.1 §8.1):

- **Lifecycle** — the movement of State over time: one artifact moves through its State Domains while its domain never changes.
- **Engineering Domain** — the home classification of an artifact: one of five canonical categories (Discovery, Architecture, Planning, Execution, Operations), derived from the token family and Knowledge Dimension.
- **Knowledge Stratum** — the authority level of the domain: a fixed position in the strict order **Discovery → Architecture → Planning → Execution → Operations** (stratum 1 highest → 5). Always derived from the domain; never declared by an artifact.

An artifact is born in a domain, moves through its lifecycle, and is consumed — it never changes stratum. This document describes the whole flow, step by step.

## Produce

Knowledge enters the repository at the point where it is created:

- **Ephemeral knowledge** — session notes (`ses-`), spike investigation (`spk-`), research findings (`fnd-`): captured while the work happens, low ceremony, `content-state: draft` where applicable.
- **Durable knowledge** — requirements (`req-`), architecture (`arc-`), decisions (`adr-`), specifications (`spec-`), standards (`std-`): started as drafts and matured through the governance channel.

**Distillation** (EKA 11.4) is the bridge from ephemeral to durable: findings from sessions, spikes, and reviews are distilled into the durable dimension they affect — direction → `adr-`/`dec-`, undecided technical findings → `fnd-`, proven procedures → `run-`. Archiving without distillation violates the protocol.

## Organize

Every artifact is classified so that it can be found, related, and validated:

- **Identity frontmatter** — `namespace`, `type`, `id`, `instance-version`, `revision`; the filename is only a projection.
- **Classification** — `dimension` (must equal the home folder, R6) and `dimensions-secondary`; the Engineering Domain is **derived** from the token family — you may declare it explicitly with `domain:` frontmatter, but it must then match the derived home domain (R11).
- **Relationships to higher strata** — `derives-from`/`depends-on` chains give an artifact its stratification traceability: a non-Discovery artifact must reach a strictly higher stratum through resolvable references, directly or transitively (R10). Supersession (`supersedes`/`amends`) may only target the same or a lower stratum — never upward (R12).
- **Change-log** — every state transition recorded by its single owner; the last entry per domain must equal the current value (R7).

## Validate

Before commit, run the mechanical checklist in [exchange/validation.md](exchange/validation.md) — the full rule set **R0–R12**:

| Group | Rules | Verdict |
|---|---|---|
| Structural | R0 (frontmatter, artifact rule) | blocking |
| Exchange (EKA Exchange Specification §14.2) | R1–R9 (identity, filename, state, ownership, references, dimension, change-log, projections, well-formedness) | blocking |
| Domain-aware (Core v1.1 §8.1) | R10 stratification traceability | **warning** |
| Domain-aware (Core v1.1 §8.1) | R11 domain coherence, R12 cross-stratum supersession prohibition | blocking |

**Draft tolerance:** `content-state: draft` relaxes what must already resolve — dangling references are warnings (R5), and draft knowledge artifacts are exempt from R10. Drafts are the sandbox of the lifecycle; approval moves them into the governed world.

## Project

`eka view` projects the repository without writing anything:

- **`eka view sprint`** — the active Execution Container's work items by Execution State columns.
- **`eka view wave`** — the active container's tickets with projected status.
- **`eka view ticket <id>`** — one ticket's projected status, derived from owner state.

Sprint, wave, and ticket are **Execution projections** — read-only views over Execution-domain artifacts (`ctr-`, `tkt-`, work items). The context header carries a `Domain: Execution` row. **A projection never writes** (P6): it has no State of its own, and refresh-on-read is the only policy. Projections are how Execution knowledge is *seen*; they are not how it is *stored*.

## Exchange

Import/export moves knowledge between repositories without loss (P13):

- **Export** — the package carries Identity, State (full change-log), Content, Relationships, and Classification, including each unit's Engineering Domain.
- **Import** — classification round-trips: the `dimension` and optional `domain` values come back exactly; the Engineering Domain is **derived at import** from the token family (never taken as an opaque value), so the domain mapping cannot drift between systems.

Exchange packages cross the same validation pipeline (R0–R12) as the repository itself. See [exchange/transfer.md](exchange/transfer.md) for the round-trip contract.

## Consume

Knowledge is consumed by humans and agents alike:

- **Runbooks (`run-`)** — executed as procedures during operations; the terminal, most actionable stratum.
- **Specifications and ADRs (`spec-`, `adr-`, `arc-`)** — read as the authority for what is built and why; they sit high in the stratum order and outrank lower-stratum content that contradicts them.
- **Tickets (`tkt-`) and work items** — read as the current execution plan: what to do, in what order, in what state.

Consumers read by Identity and by domain, never by location. A runbook that contradicts its approved specification violates the Stratum Authority Invariant and must be fixed at the runbook (lower stratum), not by overriding the specification.

## Lifecycle × Domain

Each lifecycle step is dominated by one Engineering Domain — the domain whose knowledge is produced or consumed at that step:

| Lifecycle step | Dominant domain | What it produces/consumes |
|---|---|---|
| **Produce** | Discovery (+ Execution) | needs and findings (`req-`, `fnd-`); ephemeral session/spike knowledge |
| **Organize** | Architecture | durable form: `arc-`, `adr-`, `spec-`, `std-`, `gls-` |
| **Validate** | all (Exchange Layer) | the R0–R12 verdict over every artifact |
| **Project** | Execution | sprint/wave/ticket views — read-only |
| **Exchange** | all (Exchange Layer) | packages carrying every domain |
| **Consume** | Operations (+ Execution) | `run-` executed; tickets read |

## Mental Model Shift

Folders are the **physical structure**; the lifecycle is the **mental model**. The same knowledge is born (Produce), shaped (Organize), checked (Validate), seen (Project), moved (Exchange), and used (Consume) — wherever its folder sits.

The Engineering Domain is **orientation, not location**: it tells you what kind of knowledge an artifact is and how much authority it carries, independent of where the file lives. The `dimension == folder` rule (R6) keeps the two aligned, but the domain — derived from the token — is what the standard and the validator reason about.

## Related

- [README.md](README.md) — serialization conventions: identity, state, phase, relationships, classification, projections.
- [operating/protocol.md](operating/protocol.md) — the Operating Manual: ordering chain, State Domains, gates, change-log, distillation obligations.
- [exchange/validation.md](exchange/validation.md) — the mechanical checklist (R0–R12) behind the Validate step.
- [exchange/transfer.md](exchange/transfer.md) — import/export conventions behind the Exchange step.
