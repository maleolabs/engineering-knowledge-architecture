# Implementation Terminology Glossary

Glossary of **implementation-level** terms — terms born from this repository's serialization conventions (Git + Markdown). Terms defined by the standard (Artifact, State Vector, Identity, etc.) are **not redefined here**; they live in [`standard/glossary.md`](../standard/glossary.md) and are only linked.

## A

### artifact rule
Rule for distinguishing an Artifact from a convention document: **a file is an Artifact iff its frontmatter carries `type` AND `id`**. See `reference-architecture.md` §3. Artifact concept: [standard/glossary.md](../standard/glossary.md).

## D

### derived condition
A condition that is not a state domain value but a transition-triggering condition (example: "Completed" on containers/sessions, derived from the aggregated state of work items). Applies in EKA 7.2; in serialization, derived conditions are never written as state values in frontmatter — they are computed.

### derived-from
The frontmatter relationship field (`derives-from`) stating that a State Projection is derived from the owner artifact (example: `tkt-` → `ctr-` → work item). References by Identity, not location. Relationship concept: [standard/glossary.md](../standard/glossary.md).

### dimension == folder rule
Location rule: knowledge artifacts must live in their Knowledge Dimension folder (`dimension == folder`), enforced by validation; operating artifacts are exempt. See ADR-005 and EKA 8.

## E

### exactly-one-active
Operating Layer concurrency convention: at most one Execution Container (`ctr-`) with `container-state: active` at any time; the next creation waits. Concepts: EKA 5.2, 9.

## F

### frontmatter
YAML metadata block at the head of a Markdown file — the serialization location of Identity, State Vector, Relationship, and change-log. Frontmatter is the only place state is written (single-writer); the file body is Content.

## I

### instance-version vs revision
Two version meanings in frontmatter: `instance-version` is part of Identity (a new instance = a new instance Identity, e.g., `plan-x-v2`); `revision` tracks Content evolution of the same instance and is **not** part of Identity. Concepts: EKA 6.3, [standard/glossary.md](../standard/glossary.md).

## L

### lock-atomic-with-generation
Operating Layer invariant: the Execution Container creation event locks the plan (Planning State → Immutable) and creates the container atomically — no gap between lock and generation. Concepts: EKA 5.2, 9.

## N

### namespace
Frontmatter field separating Identity spaces (product/organization/system). In this repository, example artifacts use the project namespace; meta-documentation uses `eka-ref-impl`. Concept: [standard/glossary.md](../standard/glossary.md).

## O

### on-read refresh
Default Projection Refresh policy of this serialization: a State Projection is validated against its owner **when read** (not event-driven). The "projections never write" invariant remains absolute. Concepts: EKA 5.5, 15.5, [standard/glossary.md](../standard/glossary.md).

## P

### projected state
State the artifact does not own (not part of the State Vector) but derives from the owner through Projection Semantics and is validated via Projection Refresh. The serialization marks it with the "Generated — State Projection" header on generated artifacts. Concept: [standard/glossary.md](../standard/glossary.md).

## S

### State Vector owned
The part of the State Vector the artifact **owns** (per its type), serialized as state fields in frontmatter; domains not owned = absence (not-applicable). Projections are not part of the owned State Vector. Concept: EKA 7.4, [standard/glossary.md](../standard/glossary.md).

## T

### type token
Filename prefix `<type-token>-<id>[-v<nn>]` marking the artifact type (26 ambiguity-free tokens, see `reference-architecture.md` §2.1). The token is an Identity projection for human navigation + validation; the true Identity is in frontmatter.

## W

### well-formed content
Content complying with the per-artifact-type structure (defined per-folder in the skeleton) so it can be parsed and executed deterministically. Concept: EKA 3, [standard/glossary.md](../standard/glossary.md).

## Z

### zone
Top-level division of the repository: `standard/` (canonical texts, pre-layer), `skeleton/` (copyable project structure — serialization), `reference/` (meta-documentation of this implementation). Zones are a repository organization concept, not a standard concept.
