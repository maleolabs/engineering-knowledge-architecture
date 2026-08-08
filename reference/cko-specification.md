# Canonical Knowledge Object Specification — EKA Runtime v0.2.0

> Convention document, not an artifact (no `type`/`id`). Defines the **Canonical Knowledge Object (CKO)**: the one canonical internal representation of one Engineering Knowledge Object in the v0.2.0 Runtime.
> Authoritative decisions: [ADR-012](decisions/adr-012-canonical-knowledge-object-runtime.md) (this model), [ADR-011](decisions/adr-011-immutable-engineering-knowledge-model.md) (integrity and persistence), [ADR-008](decisions/adr-008-engineering-domain-model.md) (domain/stratum derivation). The schema *is* the `exchange.Unit` model: [`../../exchange/model.go`](../../exchange/model.go). This document summarizes and explains; the ADRs are normative for this implementation.
> Terminology status: the names *Canonical Knowledge Object*, *Knowledge Compiler*, and *representation* are **not finalized** — a future milestone may rename them. Everything here is experimental.

## 1. What a CKO is

A **Canonical Knowledge Object (CKO)** is the canonical internal representation of **one Engineering Knowledge Object** — one artifact instance with its identity, state, relationships, classification, and content — independent of the authoring format that produced it. It is the object the runtime stores, resolves, verifies, and projects; nothing above the authoring boundary reads anything else.

The CKO is not a new model: it is the `exchange.Unit` model of the Exchange layer (the RSF unit entry). One canonical representation means exactly one shape for one knowledge object — the same bytes whether produced by `eka export`, by the docs-mode pull, or by a future JSON authoring adapter. **Markdown is the authoring format; the CKO is the runtime representation.** The two are connected by the Knowledge Compiler (ADR-012 §Decision 2), never by the runtime itself.

## 2. Why: a representation-independent runtime

Before the CKO pivot, Markdown was the primary internal representation: the validator, the projections (`eka view` / `eka watch` built over `conformance.Artifact` via `conformance.Scan`), and the docs-mode pull pipeline all parsed Markdown (ADR-012 §Context). Every future subsystem — the Context Engine, Machine APIs, MCP, Atrium, the Knowledge Graph, Vector Search — would have inherited that coupling. The runtime must standardize **Engineering Knowledge**, not document formats.

The CKO completes the inversion begun by the immutable-model milestone: the store already persisted canonical objects (`object_payloads` rows are `unit.json` + representation-tagged content; ADR-011 §Decision 1); the pivot makes the **consumption side** canonical too. Consequences:

- The runtime never parses Markdown; swapping the authoring format touches only the authoring adapter, nothing below the compiler.
- Every future authoring interface — JSON, forms, visual editors, AI-generated knowledge — enters the runtime through the same compiler pipeline and inherits validation, normalization, and integrity verification.
- Every future consumer — terminal projections, MCP, Atrium — consumes the same CKO contract.

## 3. The schema

The CKO schema is the `Unit` struct (`exchange/model.go`), serialized as the RSF unit entry (`unit.json` + the representation payload). Field order in serialization is the fixed declared order of the struct (RSF §9.3); JSON field names below are the `unit.json` spellings.

### 3.1 Field reference

| Field | `unit.json` / model | Kind | Notes |
|---|---|---|---|
| **Identity** | `identity` (`Identity`) | object | the complete identity tuple `{namespace, type, id, instance_version}`; the RSF Canonical Identity Form `<namespace>/<type>:<id>:<instance-version>` is echoed as `canonical_identity_form` (`Identity.CanonicalForm()`) |
| **Revision** | `revision` | integer | unit metadata, never identity, never an ordering key |
| **Metadata** | `author`, `created`, `updated` | string | preserved losslessly; `""` when absent on the source |
| **State Vector** | `state_vector` (`StateVector`) | object | the five owned state domains in canonical declared order (`content-state`, `execution-state`, `planning-state`, `container-state`, `existence-state`); empty-vector types serialize as `{}` |
| **Change Log** | `change_log` (`ChangeLog`) | array | transition entries `{date, domain, from, to, by}` in occurrence order; never fabricated, truncated, or reordered |
| **Relationships** | `relationships` (`Relationship`) | array | `{type, target}` pairs ordered by (type, target); the source is the unit's own Identity |
| **Classification** | `classification` (`Classification`) | object | `{dimension, dimensions_secondary, domain}` — the primary Knowledge Dimension, the secondary dimensions list, and the Engineering Domain |
| **Phase** | `phase` | string | the context attribute on planning/scope artifacts; outside the State Vector, never a State Domain |
| **Content** | `content` (`ContentRef`) | object | `{representation, file}` — the representation-tagged payload reference; the raw payload bytes (`ContentPayload`) travel as the unit's content file, never inside `unit.json` |

### 3.2 Derived values — never authored

Two CKO values are **derived, never authored**; they come from the type token, not from the authoring representation:

| Derived | Derivation | Source of truth |
|---|---|---|
| **Engineering Domain** (`classification.domain`) | from the artifact type token | `conformance.DomainForToken` — the single shared source of truth (ADR-008; exposed as `Unit.Domain()`) |
| **Knowledge Stratum** | from the domain | `conformance.Stratum` (ADR-008; stratum 1 highest → 5) |

### 3.3 Integrity metadata

| Field | Value |
|---|---|
| **object_hash** | `SHA-256(unit.json bytes ‖ content bytes)` — the digest over the canonical unit serialization concatenated with the raw content payload; **byte-identical to the RSF per-unit digest** (ADR-011 §Decision 1; RSF §9.4) |

The hash is **content-derived, never DB-generated**: the object *is* its address and its verification key. The same bytes hashed by the exchange layer, the store, or `eka integrity check` produce the same value, so object hashes agree with snapshot digests by construction. `content` is never `NULL` — a zero-length payload is represented, not absent; the digest covers `unit.json ‖ content` unconditionally.

## 4. Representation independence

The CKO schema contains **no Markdown syntax concepts** — nothing in it references frontmatter, headings, file layout, or any other authoring syntax:

- The **content payload is opaque and representation-tagged**: `content.representation` names the representation (`eka/structured-text/1` today), `content.file` locates the payload within the unit entry. The runtime never inspects payload internals — rendering is the presentation layer's job, extraction is future work (ADR-012 §Consequences).
- A future authoring representation (e.g., JSON authoring) stores its payload under its own representation identifier **without any schema change** (ADR-012 §Decision 5).
- Markdown awareness lives exclusively in the **authoring adapter** — the conformance package's `Scan`/`analyzeFile` (`conformance/scan.go`, `conformance/artifact.go`). The compiler is representation-independent by construction: it consumes what the adapter produces and never touches authoring syntax itself (`compile/compile.go`).

## 5. The Runtime Contract

Every runtime subsystem consumes the CKO; Markdown awareness belongs exclusively to the authoring adapter:

| Subsystem | Consumes | Never does |
|---|---|---|
| **Knowledge Compiler** (`compile/`) | authoring representations (via the adapter) | parses Markdown itself |
| **Runtime store** (`store/`) | CKO payloads: `object_payloads` (`unit.json` + content) + `object_refs` | caches or indexes Markdown |
| **Resolver** | the Canonical Identity Form → current immutable payload | re-reads authoring |
| **Projections** (`view/`, `eka view` / `eka watch`) | `exchange.Unit` (CKO), compiled on demand | parses Markdown (`view` imports no `conformance.Scan`/`Artifact`) |
| **`eka validate`** | the **authoring representation** — R0–R12, the compiler's stage-3 validation | validates CKO (that is `eka integrity check`'s scope) |
| **`eka integrity check`** | **CKO integrity**: payload hash, decode, references, attachments, registry | re-validates Markdown |

The runtime contract in one line: **authoring is an input; CKO is the model**. Two validators with two roles — authoring validation (`eka validate`, adapter layer) and runtime validation (`eka integrity check`, CKO layer) — are the explicit split (ADR-012 §Decision 4).

## 6. Relationship to the Reference Serialization Format

The CKO and the RSF share **one model, one canonical serialization**:

- `unit.json` **is** the serialization projection of the CKO. A CKO persisted in the store and the same knowledge transported in a Knowledge Snapshot are the same bytes — `object_hash` equals the RSF per-unit digest **by construction** (ADR-011 §Decision 6), so transport verification and store verification agree without translation.
- There is no second serialization format: introducing one would duplicate the canonical model and break digest identity with snapshots (ADR-012 §Alternatives).
- The runtime never serializes mutable state (registry, refs, `sync_log`) into packages; a snapshot contains only immutable CKO payloads and attachments.

## 7. References

- [ADR-012 — Canonical Knowledge Object Runtime](decisions/adr-012-canonical-knowledge-object-runtime.md) — the authoritative decision for this model
- [ADR-011 — Immutable Engineering Knowledge Model](decisions/adr-011-immutable-engineering-knowledge-model.md) — integrity (`object_hash` = RSF per-unit digest) and persistence (`object_payloads` / `object_refs`)
- [ADR-008 — Engineering Domain Model](decisions/adr-008-engineering-domain-model.md) — Engineering Domain and Knowledge Stratum derivation (`DomainForToken`, `Stratum`)
- [`runtime-architecture.md`](runtime-architecture.md) — the runtime document (schema summary, store, synchronization; pipeline overview in §2.1)
- [`migration-report-runtime-v0.2.0.md`](migration-report-runtime-v0.2.0.md) — what the runtime milestones change for existing repositories
- RSF v1.1 — unit serialization (§5), content representation model (§6), integrity and per-unit digest (§9.4): [`eka-reference-serialization-format-v1.1.md`](eka-reference-serialization-format-v1.1.md)
- CKO model (the schema of the canonical object): [`../../exchange/model.go`](../../exchange/model.go)
- CKO serialization — `MarshalUnit` (canonical `unit.json` bytes): [`../../exchange/emit.go`](../../exchange/emit.go); `DecodeUnit` (strict decode): [`../../exchange/decode.go`](../../exchange/decode.go)
- Knowledge Compiler — `compile.Compile` (the canonical authoring→CKO gateway): [`../../compile/compile.go`](../../compile/compile.go)
- Markdown authoring adapter — `Scan`/`analyzeFile`: [`../../conformance/scan.go`](../../conformance/scan.go), [`../../conformance/artifact.go`](../../conformance/artifact.go)
- Shared ontology helpers — `DomainForToken`/`Stratum`: [`../../conformance/domain.go`](../../conformance/domain.go); `OwnedDomains`/`DomainValues`: [`../../conformance/state.go`](../../conformance/state.go); `ParseReference`: [`../../conformance/rules.go`](../../conformance/rules.go)
- Conformance gate (R0–R12), the compiler's authoring-validation stage: [`../skeleton/docs/exchange/validation.md`](../skeleton/docs/exchange/validation.md)
