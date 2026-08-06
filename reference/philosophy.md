# Philosophy — Why EKA Exists and Why This Repository Is Structured This Way

This document is a position narrative: the reasons EKA exists, the architectural insights underlying it, and their consequences for how this repository is structured. It is not part of the standard — it records why the standard and serialization decisions were made.

## Why EKA exists

The early implementation proved two responsibilities coexisting in one repository: storing knowledge (Knowledge Base) and running engineering work (Engineering Operating System). Both ran on the same structure — and that is where the conflation was born. Status was duplicated in three places with no single writer, seven state machines coexisted without shared rules, Identity was confused with location, and process pipelines were mixed into a folder taxonomy.

EKA is the answer to that conflation: not a new structure, but a **canonical conceptual model** that separates responsibilities explicitly, then lets each implementation choose its own mechanisms. This repository is one serialization of that model.

## The dual-layer insight: Knowledge Base and Operating System are two layers of one system

The core EKA insight: an engineering repository is not merely a place to store documents, and not merely an execution engine — it is **both at once, as two layers of one system**. The knowledge layer stores Content, classification, and history; the operating layer runs state machines, protocol, and governance. Both are bound by **Identity**: the same artifact is seen by the Knowledge Base as knowledge and by the Operating System as a stateful entity.

Direct consequence for this repository: the knowledge folders (12 dimensions) and the `operating/` folder are not two areas that happen to sit side by side — they are two layers whose contracts are defined by the standard (Sections 4–5). One artifact carries Content in the KB layer and State in the OS layer, without either writing into the other.

## Pipeline as first-class Protocol, not a taxonomy accident

In the legacy structure, the workflow (PRD → MVP → epics → roadmap → sprint → tickets → work items → sessions) looked like a folder hierarchy. That was a category error: a pipeline is an **execution order**, a property of the Operating Layer, not a property of classification. When the pipeline was encoded as folders, every process shift forced a location shift — and Identity shifted with it.

EKA moves the pipeline to its proper place: **Protocol** (EKA 3, 9). The order "requirement → scope → plan → container → work item → working context → validation" is defined as an Operating Layer property answering "what next" deterministically. Repository consequence: there are no pipeline folders; there are dimension folders (for knowledge) and operating folders (for execution), with order expressed through Relationship and State — not through folder position.

## Phase is context, not category

Product phases (Discovery, MVP, Milestone, Release) are **time-bound context**, not a permanent category and not state. The same product remains the same product when it changes phase — what changes is its context attribute. Encoding phase as folders made phase part of location, and location part of Identity: changing phase meant "moving" the artifact, while Identity must not shift (P3).

Repository consequence: `phase` is a frontmatter field on scope/plan artifacts only; phase changes are context updates authorized by the readiness gate (EKA 11.2) — not file moves. Identity is decoupled from phase.

## State is single-writer owned

Status duplication was the main disease of the legacy structure: status lived in document metadata tables, sprint tables, and ticket documents — three copies, no clear writer, with manual sync rituals that always lagged behind. EKA establishes: **every state field has exactly one owner** (P6). All other views are State Projections: derived, validated, and never writers.

Repository consequence: state is written only in the owning artifact's frontmatter (single-writer per domain), transitions are recorded in `change-log`, and derived representations (container tables, tickets) are explicitly labeled as projections refreshed on-read — not independent editable facts.

## The repository is one serialization, not the architecture

The final discipline guiding the whole structure: **the standard is canon; the repository is an example** (EKA 1.3, 16.2). EKA is defined over concepts — Identity, State, Layer, contracts — not over folders, naming patterns, or formats. This repository is merely one way to serialize those concepts into Git + Markdown. Consequences:

- The standard lives intact in `standard/` as canonical text unmixed with serialization decisions;
- Serialization conventions are documented explicitly (not hidden in habits) and maintained through ADRs;
- The copyable structure (`skeleton/`) is a byproduct — a serialization other projects can use;
- Every serialization decision can be accounted for against a standard section/principle (see `traceability-matrix.md`).

With this discipline, the repository remains useful six months from now — and sixty months, even when the storage medium changes: what is preserved is the concept, not the mechanism.
