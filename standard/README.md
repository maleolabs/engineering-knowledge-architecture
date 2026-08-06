# Standard Zone — EKA v1.1

This zone contains **the standard itself**: the canonical copy of the Engineering Knowledge Architecture (EKA) v1.1, together with its glossary of canonical terms and its exchange contract (Exchange Specification v1).

## Nature of this zone

- This zone is a **pre-layer**: the standard defines the layers (Knowledge Layer, Operating Layer, Exchange Layer), but the standard itself is **not an artifact of any project**.
- The content of this zone is **not part of any project serialization** — it is canonical text copied for reference, onboarding, and conformance.
- Unlike the artifacts in `skeleton/docs/`, the documents in this zone carry no Identity/State frontmatter: they are copies of the standard, not Artifacts managed by the Operating Layer.

## Contents

| File | Contents |
|---|---|
| [`eka-specification-v1.1.md`](eka-specification-v1.1.md) | Complete canonical EKA v1.1 text (16 sections): principles (P1–P16), Core Concepts, Layer Model, layer contracts, Identity Model, State Taxonomy, Knowledge Taxonomy (incl. §8.1 Engineering Domains and Knowledge Stratification — five Engineering Domains, Stratum Authority Invariant, Representation Aliases), Execution Taxonomy, Artifact Taxonomy, conceptual lifecycle, storage independence, import/export, extensions, open questions, evolution. |
| [`eka-exchange-specification-v1.1.md`](eka-exchange-specification-v1.1.md) | Exchange Contract v1 (Ratified, milestone 16.1.2; revision v1.1 — additive: Engineering Domain in unit classification, Conformance Rules set R0–R12): the smallest exchange unit (Artifact Instance), Identity/Relationship/State representation, three versioning dimensions, **Exchange Package Object Model** (§4.4), **Capability Declaration** (§4.5), import/export/synchronization semantics, conformance (R1–R9; R0, R10–R12 complement — Naming §9.3), round-trip guarantees, compatibility, security, evolution. Conceptual — free of serialization formats. |
| [`eka-naming-and-terminology-specification-v1.1.md`](eka-naming-and-terminology-specification-v1.1.md) | Meta-specification (Ratified): official EKA ecosystem naming — product identity, Specification Families naming pattern, reference components, tooling, repository naming, canonical term table (incl. Engineering Domain, Knowledge Stratum, Representation Alias; "Domain" naming discipline), Conformance Rules R0–R12, deprecated terminology list, migration. |
| [`glossary.md`](glossary.md) | Alphabetical glossary of all capitalized canonical terms, with exact definitions from the canonical text. |
| [`representation-alias-registry-v1.1.md`](representation-alias-registry-v1.1.md) | Convention-layer register (Ratified): methodology representation names (PRD, ADR, RFC, Epic, Sprint, Ticket, Release, …) mapped onto canonical token + Engineering Domain (Core §8.1 alias table, extended); canonical vs methodology distinction, extension governance, methodology convention statements, full alias registry table. |

## Other references

- Implementation meta-documentation: [`../reference/README.md`](../reference/README.md)
- Copyable serialization structure: [`../skeleton/README.md`](../skeleton/README.md)
