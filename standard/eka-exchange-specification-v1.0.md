# Engineering Knowledge Architecture — Exchange Specification Version 1

| Field | Value |
|---|---|
| **Version** | 1.0 |
| **Status** | Ratification candidate — milestone 16.1.2 of EKA v1.0 |
| **Anchor** | Engineering Knowledge Architecture (EKA) Canonical Specification v1.0 (Ratified) |
| **Zone** | standard/ |
| **Scope** | Conceptual exchange contract — not serialization, not implementation |

**Reading this document:** capitalized terms refer to canonical definitions (EKA v1.0 Section 3 + glossary). "must" = binding requirement (harus/wajib). "should" = recommendation within the bounds of the contract. "may" = implementation option within the bounds of the contract.

## 1. Purpose

1.1 This document is the **Exchange Contract v1** mandated by EKA v1.0 Section 16.1.2: the definition of round-trip, idempotency, referential integrity, schema versioning, and the conformance suite for the Exchange Layer.

1.2 It operationalizes canonical Section 13 (Import/Export Model) and global invariant 5.4.7 (lossless round-trip, P13) into a precise, testable contract that every implementation (repository, database, graph store, Knowledge OS) must obey when exchanging knowledge.

1.3 It is the seam for milestone 16.1.4 (Knowledge OS integration): any system that produces or consumes a conformant Exchange Package participates in the EKA ecosystem without additional negotiation.

1.4 Relationship to the Canonical Specification

- This specification operationalizes canonical Section 13 (Import/Export Model) and global invariant 5.4.7 (lossless round-trip, P13) into a testable contract.
- It introduces Exchange-Spec subordinate terms, defined in Section 4.1.
- It does not amend the canonical specification.

## 2. Scope

### 2.1 In scope
- the smallest exchange unit and the unit hierarchy;
- the contract for Identity, Relationship, State, and versioning representation;
- the conceptual composition of the Exchange Package;
- import, export, and synchronization semantics;
- conformance validation and round-trip guarantees;
- compatibility and evolution strategy.

### 2.2 Non-goals (explicitly out of scope)
- serialization formats (JSON, YAML, protobuf, XML, markdown, or any other) — mechanism, per canonical 12.4;
- CLI, REST API, MCP, or any interface design;
- storage layout, file naming, directory structure, addressing (canonical 12.3);
- transport protocols;
- cryptographic algorithms (mechanism = implementation decision; see Section 17);
- database schemas or index design.

2.3 This specification establishes only the contract and its invariants. Any concrete serialization of this contract is an implementation (canonical 1.3, 12.4) and must be validated against it.

## 3. Design Principles

| # | Principle | Source | Exchange implication |
|---|---|---|---|
| D1 | Storage independence | 12.2–12.4 | The contract applies to any medium; no location/path/format dependencies anywhere |
| D2 | Lossless exchange | P13, 5.4.7 | Identity, State, Content, Relationship, classification survive; nothing lost, nothing duplicated |
| D3 | Deterministic ordering | P11, 12.3 | Identical repositories produce identical exports; ordering contract in 6.3 / 10.5 |
| D4 | Versioned contract | 13.2.6 | Every package declares exchange format version + specification version; the importer enforces |
| D5 | Minimum canonical core | P14 | Exchange core = unit, identity, state, relationship, round-trip; the rest = declared extensions |
| D6 | Enforcement varies, invariants don't | P16 | Conformance verdicts are identical across implementations; mechanisms (validators, constraints) differ |
| D7 | Validation before commit | 13.2.5, transfer.md 1.5 | No write before full validation; the target is revalidated after import |
| D8 | Explicit over implicit | P2, P9 | Nothing derived is exchanged as fact; projections never travel as truth |

## 4. Canonical Exchange Model

### 4.1 Concepts

| Exchange-Spec term | Definition |
|---|---|
| **Exchange Package** | self-contained, versioned carrier containing Exchange Units + Contract Header + manifest + declarations; produced by one exporter in one operation |
| **Exchange Unit** | the smallest exchange unit: one Artifact Instance with complete Identity, Content, State Vector (all owned domains), Change Log, Relationship, and classification |
| **Contract Header** | declarative package block: exchange format version, specification version, exporter identity, export scope, extension declarations, data integrity |
| **Export Scope** | declared type of export selection (Section 5.3) |
| **Closure** | smallest set of Artifact Instances that contains the seed set and is closed under traversal of declared Relationships (Section 12.2) |
| **External Reference** | Relationship target whose Identity is not carried by the package; must be declared (Section 12.3) |
| **Import Manifest** | record of import decisions during validation: per-unit verdict accepted / rejected / warned / duplicate (no-op) / re-namespaced, each with a reason |

4.2 These are Exchange-Spec terms subordinate to the canonical glossary — they extend, never replace, canonical terms.

4.3 The Exchange Layer has this model and nothing more: it has no Content and no State (4.1); it never writes State transitions (8.5); it never mutates Content (P8, P10).

## 5. Exchange Units

### 5.1 Unit hierarchy

| Level | Composition | Typical use |
|---|---|---|
| Artifact Instance | one instance: Identity + Content + State Vector + Change Log + Relationship | smallest unit; per-instance transfer |
| Artifact Line | all instances of one (Namespace, Type, ID), ordered by InstanceVersion | history transfer, supersession chains |
| Collection | explicit set of Lines selected by Identity (e.g., a set of ADRs) | curated transfer |
| Graph | Closure of a seed set over Relationship (12.2) | full dependency transfer |
| Repository | all instances of all Lines | full mirror, backup, migration |

### 5.2 Decision: the smallest exchange unit is the Artifact Instance

The smallest exchange unit **must** be the Artifact Instance.

Justification:
1. It is the smallest entity that is a complete Artifact per the canonical definition (Section 3): it carries Identity, Content, State Vector, and Relationship. An Artifact Line is only an Identity grouping (Section 3) — it has no Content and no State of its own, so it cannot be exchanged as a complete Artifact.
2. InstanceVersion is part of Identity (canonical 6.1); instance granularity is the exact semantic granularity of Identity.
3. Line-level exchange forces transfer of the full instance history when only one instance is needed — violating Minimum Canonical Core (P14).
4. All unit invariants (State validity, Change Log consistency, referential integrity) are instance-local or instance-declared; Line/Collection only add selection semantics, not new invariants.

Alternatives evaluated and rejected: whole Artifact (ambiguous — Line or Instance? it repeats the same question); Knowledge Object (not a canonical glossary term; using a new core concept violates Section 3 discipline); Package (composite, not atomic); Collection (selection semantics, not an atomic unit).

5.3 Export Scope (selection, not unit): single instance; single Line; Collection; Graph (Closure); Repository. Any scope **must** resolve to a set of Artifact Instances; the smallest unit in every package remains the Artifact Instance (5.2).

## 6. Identity Representation

6.1 Every Exchange Unit **must** carry the complete instance Identity tuple: `(Namespace, Type, ID, InstanceVersion)` (canonical 6.1). No component may be omitted, derived, or defaulted.

6.2 Canonical serialization contract (6.2.6): Identity **must** be serialized:
- **losslessly** — the tuple round-trips without information loss;
- **unambiguously** — one representation ↔ exactly one Identity; no two representations of one Identity;
- **machine-parseably** — parseable without human interpretation;
- **independently** — must not depend on location, file name, path, process stage, or representation conventions (canonical 6.4: Identity must not be encoded in location).

6.3 Canonical ordering key: Identity tuples are compared component-by-component in the order `(Namespace, Type, ID, InstanceVersion)`, with a total, stable comparison over the canonical serialization of each component. All ordered collections within a package (manifest, units, Relationship lists, External Reference declarations) **must** be ordered by this key. Consequence: identical repository state **must** produce identical packages (10.5).

6.4 Revision: travels as metadata on the unit — never part of Identity (canonical 6.1, 6.3). Revision changes must not affect references or Identity equality. Revision **must not** be used as an identity ordering key.

6.5 Re-namespace: the only Identity transformation allowed during exchange (13.2.4). It **must**: be declared explicitly before commit; be applied to the complete Identity tuple of every affected unit; be applied consistently to every reference to those units inside and outside the package; be recorded in the Import Manifest. Partial or silent re-namespacing is **forbidden**.

## 7. Relationship Representation

7.1 All Relationships **must** be expressed by Identity — never by location, path, display name, or classification (P3, 6.2.3).

7.2 Canonical Relationship types: `supersedes`, `amends`, `derives-from`, `depends-on`, `validates` (Section 3). Relationships are recorded on the referring artifact (one-directional, per validation rule 5).

7.3 Extensibility: new Relationship types are lightweight extensions (canonical 14.1). Packages using non-canonical types **must** declare them in the Contract Header extension declarations; an importer may reject undeclared or unknown types (16.3).

7.4 Resolution on import **must** follow this order:
1. **local** — the target Identity is carried by the package;
2. **global** — the target Identity already exists in the target repository;
3. **external** — declared External Reference: **must** resolve to an artifact in the target repository, or the referring artifact has Content State Draft (draft tolerance, validation rule 5).

Failure of all three: blocking (11.1 phase 5), except under draft tolerance.

7.5 Post-exchange validity: references remain valid because Identity is immutable and location-independent (P3, 12.2). Import **must not** rewrite references into paths/locations; it re-resolves them.

## 8. State Representation

8.1 **Equal-exchange decision**: every Exchange Unit **must** carry its complete owned State Vector — all State Domains owned for its type, with exact values (canonical 7.4; type→state binding per canonical Section 10) — plus the full Change Log. Selection is permitted only at unit granularity; partial state exchange (a subset of owned domains) is **forbidden**.

Justification: partial state breaks Change Log consistency (validation rule 7) and forward-only verification (P7); lossless exchange (P13) requires the full vector; unit-level selection already covers all scoping needs (5.3).

8.2 Change Log: the full chronological record of transitions — domain, old value, new value, time, authority (Section 3). **Must** be exchanged in occurrence order, append-only; import **must not** fabricate, truncate, or reorder entries.

8.3 Existence State: preserved exactly; transfer mechanics **must not** change it (transfer.md section 2). Archived and Retired artifacts may be exchanged and **must** remain Archived/Retired at the target.

8.4 Phase: a context attribute on planning/scope artifacts (canonical 11.2) — not a State Domain, never part of the owned State Vector. The current Phase value **must** travel with planning/scope artifacts, together with the history of context updates when the exporter maintains it. The reference implementation records Phase in change-log consistency checks (validation.md rule 7); that is a serialization-level detail and must not be read as making Phase a State Domain.

8.5 Projections: never exchanged as a source of truth (transfer.md section 3). Artifacts with an empty State Vector (e.g., Ticket = (∅), canonical 7.4) are exchanged with an empty vector; their state is regenerated at the target via Projection Refresh from owner state. The Exchange Layer never writes State (4.1) — import copies values and entries verbatim; transitions remain the exclusive right of the owning layer (transfer.md section 3).

## 9. Versioning Model

### 9.1 Three version dimensions

| Dimension | Applies to | Responsibility | Part of Identity? |
|---|---|---|---|
| Artifact version — **InstanceVersion** | Artifact Instance | distinguishes instances within one Line; created deliberately (canonical 6.3) | Yes |
| Artifact version — **Revision** | Content of one instance | tracks Content evolution; changes on every edit (canonical 6.1, 6.3) | No |
| **Specification version** | the knowledge being exchanged | version of the EKA contract the artifact complies with; determines applicable taxonomy, State Domains, and variants | No |
| **Exchange format version** | the package itself | version of the serialization contract; determines package structure and validation rules | No |

### 9.2 Negotiation rules

1. Every package **must** declare in the Contract Header: the exchange format version and the specification version (13.2.6; transfer.md 1.6).
2. The importer **must** reject packages with an unsupported exchange format version.
3. The importer **must** reject packages whose specification version cannot be validated.
4. All unit validation uses the taxonomy and state variants of the declared specification version (e.g., Content State variants: standard/ADR/decision, canonical 7.2).

### 9.3 Backward compatibility strategy

- Exchange format v1 **must** accept packages written against v1.0 of this specification.
- Contract additions **must** be optional and introduced via declared extensions (Section 16) — never silently change core semantics.
- Packages without extension declarations **must** be importable by every conformant importer (minimum canonical core, P14).

## 10. Exchange Package Structure

Conceptual composition — no serialization is mandated.

### 10.1 Contract Header
Declares: exchange format version; specification version; exporter identity; export scope (5.3); creation metadata (timestamp — package metadata, not artifact knowledge); extension declarations; data integrity (17.1).

### 10.2 Manifest
Ordered list (6.3) of all Exchange Units; declares the scope type; closure declaration (12.2); External Reference declarations (12.3); extension declarations (16.3).

### 10.3 Exchange Units
Each unit **must** carry:
- Identity: canonical `(Namespace, Type, ID, InstanceVersion)` (6.1–6.2);
- Revision metadata (6.4);
- State Vector: all owned domains, exact values (8.1);
- Change Log: full, ordered (8.2);
- Content: complete, Well-formed for its type (validation rule 9);
- Relationship: all, by Identity (7.1–7.2);
- Classification: primary Knowledge Dimension (+ optional secondary) (canonical 8.1; P15).

### 10.4 Declarations
External Reference Declarations: every Relationship target not carried by the package (12.3). Extension Declarations: every non-canonical artifact type, Relationship type, or classification use (16.3).

### 10.5 Deterministic ordering contract
- manifest, units, external references: canonical Identity key (6.3);
- Change Log entries: occurrence order, with deterministic tie-break (entry order);
- Relationship lists: (source Identity key, Relationship type, target Identity key);
- extension and external declarations: canonical key.

Consequence: two exports from identical repository state **must** be identical up to package metadata (15.4).

### 10.6 Package integrity
The manifest **must** correspond 1:1 with the units — no phantom entries, no missing units. A package failing self-consistency **must** be rejected before any validation phase (11.1).

## 11. Import Semantics

### 11.1 Ordered phases
Phases **must** execute in order. No write before phase 9. A blocking failure in any phase aborts the import before commit.

| # | Phase | Checks | Blocking (0) | Warning (W) |
|---|---|---|---|---|
| 1 | Contract validation | header well-formed; exchange format version supported; specification version declared and validatable | unsupported format; invalid header | — |
| 2 | Identity validation | every unit Identity canonical (6.2); unique within package | duplicate/non-canonical Identity | — |
| 3 | State validation | values ∈ domain value sets (canonical 7.2, incl. variants); owned-set compliance per type (rule 4); forward-only + Change Log consistency (rule 7) | invalid values; foreign domain fields; log inconsistency | — |
| 4 | Structural validation | Content Well-formed per type (rule 9) | malformed content | — |
| 5 | Referential validation | resolution per 7.4; External References declared; draft tolerance (rule 5) | unresolved non-draft references; undeclared externals | unresolved references on Draft artifacts |
| 6 | Conflict detection | Identity already exists at target with same InstanceVersion | conflict (policy: reject by default) | — |
| 7 | Duplicate detection | package (or its units) already imported — idempotency (13.2.2) | — | entire import = no-op, no write |
| 8 | Dependency resolution | commit order: target Relationships precede referrers; closure completeness | unresolved dependencies | — |
| 9 | Commit | atomic write of all accepted units | any failure → no partial commit | — |
| 10 | Post-commit revalidation | target revalidated after import (13.2.5, transfer.md 1.5) | violation | — |

### 11.2 Conflict policy (13.2.4)
- Identity conflict: **reject by default**. Explicit re-namespace is the only alternative (6.5). Silent merge is **forbidden**.
- Duplicate detection: re-importing an identical package **must** be a no-op — no new units, no state changes, no fabricated Change Log entries, no conflict errors.

### 11.3 Import Manifest
Phases 1–8 **must** produce an Import Manifest recording the per-unit verdict: accepted / rejected / warned / duplicate (no-op) / re-namespaced, each with a reason. The Import Manifest is the authoritative record of the import; validation failures (0) block, warnings (W) permit commit with notes.

### 11.4 Validation-before-commit invariant
No unit may be written before all phases pass for it (13.2.5); the target **must** pass revalidation after commit (phase 10).

## 12. Export Semantics

### 12.1 Selection
The exporter **must** declare its Export Scope (5.3). Selection resolves to a set of Artifact Instances; resolution **must** be deterministic (identical selection criteria → identical set).

### 12.2 Closure computation
For Graph scope, Closure is the smallest set of instances that contains the seed set and is closed under Relationship traversal: for every included instance, every target of a `depends-on`, `derives-from`, `validates`, `supersedes`, or `amends` Relationship **must** be included; for `supersedes`/`amends`, every artifact that references the instance is also included (so supersession/amendment chains are complete — history links, canonical 13.1). Computation is a fixed point over a finite Identity set: it terminates and is deterministic for a given seed.

### 12.3 External Reference detection
After closure, every Relationship target not carried by the package **must** be declared as an External Reference. Rules:
1. declaration is mandatory — an undeclared out-of-package reference makes the package invalid (10.6);
2. on import, an External Reference **must** resolve to an artifact in the target repository, or the referring artifact meets draft tolerance (7.4);
3. dangling after non-draft import is **forbidden** (13.2.3) — the package **must not** silently drop or rewrite External References.

### 12.4 Dependency integrity
The package **must** be referentially closed for non-draft units: every unit with non-Draft Content State **must** have all Relationship targets carried or declared external (validation rule 5). Draft units may carry unresolved references (warning only).

### 12.5 Package integrity
The exporter **must** validate the complete package against the Section 14 conformance rules before release (analogous to validation-before-commit, 13.2.5). The manifest↔unit 1:1 correspondence **must** hold (10.6). Exports **must** be deterministic (10.5).

## 13. Synchronization Model

Conceptual strategy — no algorithm is specified.

### 13.1 Strategies

| Strategy | Defined behavior | Typical use |
|---|---|---|
| **Replace** | full package replaces the matching scope at the target; declared clean replace (13.2.2); post-sync, the scope's identity set == the package's identity set; all replaced units carry state + full Change Log | mirror, restore |
| **Merge** | Identity lines join: units whose Identity is absent at the target are added; overlapping Identities go through the conflict policy (13.2); Change Log appended in original order; existing content never silently replaced | two-way sync |
| **Patch** | targeted delta: an explicit unit set (or declared per-unit operations) with declared semantics; validation pipeline identical to full import | incremental sync |

### 13.2 Merge behavior — defined
- Identity lines: disjoint lines are added; overlapping lines → conflict policy; never silent merge.
- Change Log: append-only; imported entries appended in recorded order; never rewritten, truncated, or reordered.
- Content: same Identity, different Content → content conflict → governance channel only; never auto-merged (P8, P10).
- State: forward-only reconciliation — imported values are accepted only if reachable forward from the target value within the domain chain (P7); target ahead of import → regression → rejected; equal values → no-op.

### 13.3 Conflict handling

| Conflict | Policy |
|---|---|
| Identity | reject (default) or explicit re-namespace (6.5); silent merge forbidden (13.2.4) |
| State | forward-only reconciliation (13.2); regression or divergent chains rejected |
| Content | governance channel only; never auto-merge (P8, P10) |

### 13.4 Synchronization invariants
1. no Identity duplication;
2. no dangling references after sync;
3. no silent State regression;
4. no Content mutation outside the governance channel;
5. Change Log append-only;
6. idempotent sync — re-synchronizing equivalent state = no-op;
7. re-namespace applied consistently to all references.

## 14. Conformance Requirements

### 14.1 Validator guarantees
A conformant validator **must**:
1. be mechanical — applies only the rules below, no human judgment;
2. be deterministic — identical input → identical verdict;
3. be read-only — validation never modifies the repository or the package;
4. agree with all other conformant validators (P16): verdicts derive only from this contract, never from layout, enforcement, or implementation tooling (canonical 12.3);
5. validate before commit (13.2.5) and after import (11.1 phase 10).

### 14.2 Nine conformance rules (mechanical)

| # | Rule | Violation verdict |
|---|---|---|
| 1 | **Identity uniqueness**: no duplicate (Namespace, Type, ID, InstanceVersion) in the repository (6.2.2) | blocking |
| 2 | **Identity canonical**: serialized losslessly, unambiguously, machine-parseably (6.2.6); never encoded in location/stage/convention (canonical 6.4) | blocking |
| 3 | **State value validity**: every owned domain value ∈ its domain value set, incl. declared variants (canonical 7.2) | blocking |
| 4 | **Owned-set compliance**: present state fields == the artifact type's owned domains (canonical 7.4, Section 10); no foreign fields; empty-vector types carry no fields (canonical 7.4) | blocking |
| 5 | **Referential integrity**: every reference resolves per 7.4; draft tolerance for unresolved references on Draft artifacts; a superseded ADR must point to its successor | blocking (non-draft); warning (draft) |
| 6 | **Classification**: primary Knowledge Dimension declared and valid; secondary optional; classification changes never touch Identity (P15) | blocking |
| 7 | **Change Log consistency**: every transition has an entry; last entry per domain == current value; no transitions without entries; order preserved (P7) | blocking |
| 8 | **Single-writer & projection non-writer**: projections carry no owned state of other artifacts; projections never write (P6) | blocking |
| 9 | **Well-formedness**: Content structure matches the artifact type family (Section 3) | blocking |

### 14.3 Verdict semantics
- **Conformant**: all rules pass, no warnings.
- **Conformant with warnings**: no blocking violations; warnings recorded.
- **Non-conformant**: at least one blocking violation → no commit.

Verdicts apply to the exchanged knowledge and to repository state alike: a unit that fails a rule is non-conformant both as exchanged content and as committed repository state. Implementation-specific enforcement mechanisms (filename conventions, folder layout, database constraints, tooling) **must not** change verdicts (P16, canonical 12.3). In the reference implementation, the filename/dimension-folder rules (validation.md rules 2 and 6) are enforcement mechanisms of one serialization — the invariants they serve are rules 2 and 6 above.

## 15. Round-Trip Guarantees

### 15.1 Definition
Lossless exchange (P13) iff for any repository R: export E(R) → import into an empty repository T → export E(T) **equals** E(R) up to permissible differences (15.4). Idempotency (13.2.2): importing E into a repository that already contains its result = no-op.

### 15.2 Invariant properties (must hold across export→import→export)

| Property | Requirement |
|---|---|
| Identity set equality | (Namespace, Type, ID, InstanceVersion) sets identical; no loss, no duplication |
| State vector equality | owned domain values identical per Identity, including Existence State |
| Change Log equality | identical entries, identical order: domain, old value, new value, time, authority |
| Content equality | content equal at the semantic level; byte-level equality is meaningful only relative to an implementation-defined canonical content representation; implementations must declare their canonicalization so that equality is decidable; Well-formed per type; no content format is mandated |
| Relationship equality | (source, type, target) sets identical |
| Classification equality | primary and secondary Knowledge Dimensions identical |
| Line integrity | instance set per Line identical; Line Identity unchanged |

### 15.3 Projection exclusion
Projections are not exchanged (8.5); they are regenerated by Projection Refresh at the target. Projection differences are therefore never a loss.

### 15.4 Permissible differences
The following may differ between exports of equivalent repositories:
- physical order/layout in storage (P3, P9);
- storage location and addressing (canonical 12.3);
- projection refresh timestamps (8.5);
- package metadata: creation timestamp, exporter identity in the Contract Header (10.1).

**Must NOT** differ: anything listed in 15.2.

### 15.5 Idempotency guarantee
Re-importing the same package: no new units, no state changes, no fabricated Change Log entries, no conflict errors, no re-namespace. Re-export after import: identical package up to 15.4.

## 16. Compatibility Strategy

### 16.1 Version negotiation
- Packages declare the exchange format version + the specification version (9.2).
- Importers reject unsupported exchange format versions and unvalidatable specification versions (13.2.6).
- Validation always runs against the taxonomy of the declared specification version (9.2.4).

### 16.2 Backward compatibility
- An exchange format v1 importer **must** accept v1.0 packages (9.3).
- New contract features **must** be optional and declared extensions (9.3); the canonical core is closed (14.2.3).
- Extensions **must not** weaken invariants (14.2.1) and **must** remain exchangeable (14.2.4).

### 16.3 Forward compatibility limits
- Unknown exchange format version → reject (cannot be interpreted safely).
- Unknown artifact types or Relationship types without extension declarations → reject.
- Declared extensions: the importer may accept (if implemented) or reject; rejection is explicit, never silent.
- Backward compatibility is guaranteed for v1.x within this contract; forward compatibility is bounded by declarations.

## 17. Security Considerations

Conceptual only — no cryptographic algorithms are mandated (canonical 12.4, P14).

17.1 **Integrity verification**: every package **must** carry integrity data over its units (mechanism: implementation choice). Importers **must** verify before phase 1; verification failure → reject.

17.2 **Provenance**: the Contract Header records exporter identity. Provenance is recorded metadata; it is never treated as authority over state or content.

17.3 **Authority of state transitions**: Change Log entries record authority (Section 3); import copies entries verbatim — the Exchange Layer never writes transitions (4.1, 8.5). Imported state is copied, not re-authorized; the target's Protocol governs all subsequent transitions.

17.4 **Injection of malicious content**: conformance validation (Section 14) is not a security boundary. Imported content may be Well-formed yet untrusted; content governance (P8, P10) applies at the target. Re-namespace (6.5) is the explicit mechanism for importing into a separately trusted namespace.

17.5 **Trust boundaries**: import is a trust decision by the target. External References create cross-boundary links: their resolution is verified, but their target content is not validated by the importing package. Draft tolerance (7.4) is the only temporary unresolved state permitted.

## 18. Future Evolution

18.1 This specification is milestone 16.1.2: the exchange contract v1 with a conformance suite. Milestone 16.1.3 (reference implementations) **must** demonstrate conformance to Sections 14–15; milestone 16.1.4 (Knowledge OS) consumes this contract as its seam.

18.2 Conformance suite: a mechanical validator derived from Section 14.2, with verdict semantics fixed (14.3); validators **must** agree (P16). Suite extensions follow governance (14.2).

18.3 Extensions: new Artifact types, new Relationship types, new Knowledge Dimensions, and registered protocol variants go through governance 14.2 and **must** be declared in the Contract Header when used (10.4); extensions **must not** weaken invariants (14.2.1) and **must** remain exchangeable (14.2.4).

18.4 Future exchange format versions **must** be additive (16.2); round-trip guarantees (Section 15) and the invariant set (canonical 16.3) survive across versions.

18.5 Evolution never changes: Identity (P3), global invariants (5.4), the two-channel separation (P10), layer composition (canonical 4.1).

## Appendix A — Open Questions

The following questions remain open; this specification takes the stated positions.

1. **Offline/distributed Identity generation** (canonical open question 15.4) remains open. This specification assumes `(Namespace, Type, ID)` uniqueness is guaranteed by the producing system; collisions are handled at import via the conflict policy (11.2, 13.2.4).
2. **Semantic Line supersession** (canonical open question 15.2) affects the closure completeness of history chains. This specification chooses backward traversal of `supersedes`/`amends` for closure (12.2).
3. **Content canonicalization for equality checking** is implementation-declared (15.2); no content format is mandated.

---

*End of Exchange Specification — EKA Exchange v1.0 (milestone 16.1.2).*
