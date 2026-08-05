# EKA Reference Serialization Format (RSF) v1.0

| Field | Value |
|---|---|
| **Title** | EKA Reference Serialization Format (RSF) v1.0 |
| **Status** | Reference (not normative) |
| **Anchor** | EKA Exchange Specification v1.0 (Ratified) |
| **Zone** | reference/ |
| **Scope** | One canonical serialization projection of the Exchange Package object model; an implementation of the Exchange Contract, not an extension of it |

**Reading this document:** capitalized terms refer to canonical definitions (EKA Core Specification v1.0 Section 3 + glossary) or to Exchange-Spec subordinate terms (EKA Exchange Specification v1.0 §4.1). Terms marked "(RSF)" are defined in Section 3; they are subordinate to the Exchange Specification and never canonical. "must" = binding requirement within the RSF. "should" = recommendation within the bounds of the RSF. "may" = option within the bounds of the RSF.

## 1. Purpose

1.1 The **EKA Reference Serialization Format (RSF) v1.0** is one canonical implementation of the Exchange Contract: a concrete serialization projection of the Canonical Exchange Package Object Model (EKA Exchange Specification v1.0 §4.4). It makes the object model — Contract Header, Manifest, Exchange Units, Declarations, integrity data — serializable in a deterministic, lossless, machine-parseable form.

1.2 The RSF is **not normative**. It lives in the Reference Zone as meta-documentation of the Reference Implementation (EKA Naming and Terminology Specification v1.0 §6.1). It is not an extension of the Exchange Specification, does not amend it, and is not a new Specification Family (Naming §5.6). The Exchange Specification remains the sole normative contract; the RSF is one of its possible projections (Exchange §2.3, §4.4).

1.3 The RSF is the serialization target for future `eka export` and `eka import` capabilities (Naming §7.1), for EKA SDKs (`eka-sdk-<language>`, Naming §7.4), for interoperability testing between implementations, and for conformance testing of the Exchange Contract.

1.4 Other implementations remain free to serialize the object model differently while conformant to the Exchange Specification. The RSF claims no privileged status at the contract level: a package in any conformant serialization is judged by the same Conformance Rules (Exchange §14) and round-trip guarantees (§15).

1.5 **Core invariant.** The RSF is a bijective projection of the Exchange Package object model: every object model element maps to exactly one RSF element, and every RSF element that realizes an object model element maps back to exactly one. Nothing extra, nothing missing. RSF-subordinate support elements (Section 3) exist outside the object model: they realize nothing and are realized by nothing, and their presence never changes the projection (Appendix A).

1.6 The RSF serializes **engineering knowledge**, not repository implementation. It must not assume Git, Markdown storage, databases, Atrium, MCP, or a specific programming language.

1.7 **Two serialization layers, kept distinct.** The reference implementation repository uses Markdown + frontmatter as its *repository serialization* (Skeleton Zone, `skeleton/docs/README.md`). The RSF is a different layer: *exchange package serialization*. The RSF must not couple Content to Markdown (Section 6). It may, however, define its own canonical default content representation for the reference implementation while allowing other representations (HTML, JSON, plain text, diagram sources) to be carried with declared Representation Identifiers (Section 6). The two layers share no serialization dependency: a repository may store its Content however it chooses and still export conformant RSF packages.

## 2. Scope

### 2.1 In scope
- conceptual package identity and canonical package structure (Section 4);
- the serialization model of one Exchange Unit, the Unit Entry (Section 5);
- the Content representation model: representation-tagged payloads (Section 6);
- the Attachment model for supporting resources (Section 7);
- the Manifest model (Section 8);
- serialization principles: deterministic ordering, canonical identifiers, encoding conventions, integrity, extensibility (Section 9);
- round-trip mapping to Exchange §15 invariants (Section 10);
- compatibility strategy for the RSF itself (Section 11);
- reference examples, illustrative only (Section 12);
- implementation recommendations for future import/export engines (Section 13).

### 2.2 Out of scope
- transport protocols and delivery mechanisms;
- CLI, REST API, MCP, or any interface design — Naming §7 governs tooling naming; the RSF defines no commands and no APIs;
- programming language structures, language bindings, or API shapes — SDKs implement the projection; they are not specified here;
- ZIP, TAR, or any container-internal byte layout — the RSF defines logical structure only and states the logical equivalence between single-file and directory containers (Section 4.4);
- storage layout, file naming, and directory structure inside repositories (canonical 12.3);
- database schemas and index design;
- cryptographic algorithms beyond the single integrity digest mechanism chosen in Section 9.4 — the Exchange Specification stays mechanism-neutral (§17.1); the RSF, being an implementation, mandates one;
- redefinition of the Conformance Rules R0–R9 (Exchange §14.2) or the addition of new EKA concepts;
- any amendment to the Exchange Specification.

2.3 A package that conforms to the RSF conforms to the Exchange Contract. The RSF adds RSF-level requirements on top of Exchange requirements; it removes none.

## 3. Definitions — RSF Subordinate Terms

The terms below are defined by this document at RSF level. They are subordinate to the canonical glossary and to the Exchange-Spec terms (Exchange §4.1); they extend neither. They must be capitalized when used as terms and never used as canonical terms.

| Term | Definition (RSF-level) |
|---|---|
| **Package Header (RSF)** | the RSF realization of the Contract Header (Exchange §10.1): the contract facts plus the Serialization Version and the Package Identity Label |
| **Package Identity Label (RSF)** | deterministic, conceptual identity of a package; derived from Export Scope + declared source namespace + Serialization Version; not required to be a filename (Section 4.1) |
| **Unit Entry (RSF)** | the serialized realization of exactly one Exchange Unit; carries the full unit composition (Section 5.1) |
| **Attachment (RSF)** | a supporting resource (diagram, image, binary, PDF, asset) carried in the Attachments Collection; not an Exchange Unit, never part of Content semantics (Section 7) |
| **Attachment ID (RSF)** | deterministic, package-scoped identifier of an Attachment (Section 7.2) |
| **Attachments Collection (RSF)** | package-level collection of Attachments (Section 7.1) |
| **Attachment Reference Set (RSF)** | RSF-subordinate set on a Unit Entry referencing Attachments by Attachment ID (Section 7.3) |
| **Representation Identifier (RSF)** | declared token identifying the representation of a Content payload (Section 6.2) |
| **Canonical Identity Form (RSF)** | the deterministic string form of the Identity tuple: `<namespace>/<type>:<id>:<instance-version>` (Section 5.2) |
| **EKA Structured Text (RSF)** | the canonical default Content representation: a metadata block plus structured body sections (Section 6.3) |
| **Declarations Block (RSF)** | the RSF realization of the package-level Declarations element (Exchange §10.4) |
| **Integrity Block (RSF)** | the RSF realization of package integrity data (Exchange §17.1): package-level, per-unit, and per-Attachment digests (Section 9.4) |
| **Serialization Version (RSF)** | the version of the RSF itself (v1); distinct from the exchange format version and the specification version (Section 11.2) |

## 4. Package Model

### 4.1 Package identity

4.1.1 Every package has a conceptual identity — the **Package Identity Label (RSF)** — derived deterministically from: the declared Export Scope (Exchange §5.3); the declared source namespace (the Namespace the exporter declares as primary; for multi-namespace packages, the lexicographically smallest Namespace); and the Serialization Version (RSF). Illustrative pattern only: `rsf-<scope>-<namespace>-<serialization-version>`. The label is an identifier, not a filename requirement: identity is conceptual and location-independent. The independence discipline of Exchange §6.2 applies to unit Identity; the RSF applies the same discipline to the package label by convention.

4.1.2 Determinism: identical export operations on identical repository state must yield the same Package Identity Label.

### 4.2 Package extension

4.2.1 The RSF recommends the extension `.ekapkg` for single-file containers (illustrative: `rsf-repo-acme-1.ekapkg`). Directory-layout packages carry no extension. The extension is a **convention, not a contract**: importers must not rely on it to classify a package; the declared Package Header facts are authoritative.

### 4.3 Package metadata (contract facts)

4.3.1 The Package Header (RSF) declares: exchange format version; specification version; Serialization Version (RSF); exporter identity; Export Scope; creation timestamp. These are contract facts and package metadata only (Exchange §15.4): they are excluded from round-trip equality. The Header also announces integrity data and the presence of declarations (Exchange §10.1).

### 4.4 Canonical package structure

4.4.1 Conceptually, a package contains exactly:

| Element | Cardinality | Realizes |
|---|---|---|
| Package Header | exactly one | Contract Header (Exchange §10.1) |
| Manifest | exactly one | Manifest (§10.2) |
| Unit Entries | zero or more | Exchange Units (§10.3) |
| Attachments Collection | zero or more | — (RSF support element) |
| Declarations Block | exactly one, may be empty | Declarations (§10.4) |
| Integrity Block | exactly one | integrity data (§17.1) |

4.4.2 A single-file container and a directory layout are the **same logical structure**: the logical elements above map one-to-one onto archive entries (single-file) or onto files/directories (directory layout). The RSF does not specify archive internals (ZIP, TAR, or any other): any container that preserves the logical structure and supports integrity verification is acceptable (Section 2.2).

## 5. Artifact Serialization Model

### 5.1 Unit Entry composition

One Unit Entry (RSF) per Exchange Unit (Exchange §5.2). The composition is exactly the object-model unit composition (Exchange §4.4), plus the RSF-subordinate Attachment Reference Set:

| Element | Cardinality | Rule |
|---|---|---|
| Identity | exactly one, complete tuple | canonical tuple `(Namespace, Type, ID, InstanceVersion)`; never omitted, derived, or defaulted (Exchange §6.1); serialized in the Canonical Identity Form (5.2) |
| Revision metadata | exactly one | unit metadata; never part of Identity (Exchange §6.4); never an ordering key |
| State Vector | exactly one, complete | all owned domains of the unit's Type, exact values (Exchange §8.1); empty-vector Types carry an empty State block — present but empty (Exchange §8.5, R4) |
| Change Log | exactly one, full and ordered | all transitions in occurrence order, append-only: domain, old value, new value, time, authority (Exchange §8.2) |
| Content | exactly one, complete, Well-formed for the Type | representation-tagged payload (Section 6) |
| Relationship set | zero or more | all Relationships by Identity, recorded on the referring unit (Exchange §7.1–7.2); ordered by (source Identity key, Relationship type, target Identity key) (Exchange §10.5) |
| Classification | exactly one primary Knowledge Dimension; at most one secondary | canonical Section 8; P15; classification never touches Identity (R6) |
| Phase context | zero or one | current Phase value plus context-update history when the exporter maintains it; planning/scope artifacts only (Exchange §8.4); never a State Domain, never part of the State Vector |
| Attachment Reference Set (RSF) | zero or more | references to Attachments by Attachment ID (Section 7.3); RSF support element, outside the object model |

5.1.1 **Empty State Vector.** An empty-vector Type (e.g., Ticket) is serialized with an empty State block, never an absent one: absence is not a value and would be ambiguous (Exchange §8.5, R4). State is copied verbatim; the RSF never writes State transitions (Exchange §4.1, §8.5; transfer.md §3).

5.1.2 **Change Log.** Entries keep occurrence order with deterministic tie-break (entry order) (Exchange §8.2, §10.5). Import must not fabricate, truncate, or reorder entries.

### 5.2 Canonical Identity Form (RSF)

5.2.1 The Identity tuple is serialized as the deterministic string:

`<namespace>/<type>:<id>:<instance-version>`

5.2.2 **Decision and justification.** The reference implementation's repository serialization uses the reference forms `<type>:<id>[:<instance-version>]` (same-namespace) and `<ns>/<type>:<id>` (cross-namespace) (`skeleton/docs/exchange/validation.md`, Aturan 5). The RSF form keeps the same component order and separators and extends the cross-namespace form with the instance-version suffix in the same `:<nn>` style. Two differences are deliberate:

1. the Namespace is always present — the RSF cannot rely on package context to supply it, because Identity must be serialized independently (Exchange §6.2) and no component may be omitted, derived, or defaulted (§6.1);
2. the InstanceVersion is always present — in the reference form, omission means "instance 1" by repository convention; a package serialization must not default identity components (§6.1).

Result: one form ↔ exactly one Identity; the form is lossless, unambiguous, machine-parseable, and independent (Exchange §6.2).

5.2.3 **Component charset.** Namespace, Type, and ID components must not contain `/`, `:`, or whitespace; InstanceVersion is a non-negative integer. This keeps the form parseable without human interpretation (§6.2).

5.2.4 **Ordering key.** All ordered collections keyed by Identity use component-wise comparison `(Namespace, Type, ID, InstanceVersion)` — the canonical ordering key (Exchange §6.3).

### 5.3 Relationships, Classification, Phase

5.3.1 Relationships use the five canonical types (`supersedes`, `amends`, `derives-from`, `depends-on`, `validates`) and are expressed by Identity only — never by location, path, or display name (Exchange §7.1). Non-canonical Relationship types must be declared in the Declarations Block as extensions (Exchange §7.3, §16.3).

5.3.2 Classification carries the primary Knowledge Dimension plus at most one secondary; classification changes never affect Identity (Exchange §14.2 R6).

5.3.3 Phase is a context attribute (canonical 11.2; Exchange §8.4) on planning/scope artifacts only. The RSF serializes it outside the State Vector; the reference implementation records Phase in change-log consistency checks — that is a serialization-level detail and must not be read as making Phase a State Domain (Exchange §8.4).

### 5.4 Revision and unit metadata

5.4.1 Revision travels as unit metadata, never as part of Identity (Exchange §6.4). It is serialized with the Unit Entry, is never used as an ordering key, and never affects reference resolution or Identity equality (Exchange §6.4).

5.4.2 Unit metadata is distinct from package metadata (Section 4.3): package metadata is excluded from round-trip equality (§15.4); unit metadata is not — Revision and recorded Change Log entries are part of the object model and must round-trip.

## 6. Content Representation Model

### 6.1 Chain

knowledge → Artifact → representation → serialization. The Artifact's Content is serialized as a payload plus a declared Representation Identifier (RSF).

### 6.2 Representation Identifiers

- Content = payload + Representation Identifier (RSF).
- The RSF registers exactly one representation in v1: `eka/structured-text/1` (conceptual token), defined in 6.3.
- Other representations (HTML, JSON, plain text, diagram sources, and others) may be carried as opaque payloads with declared Representation Identifiers.
- An unknown Representation Identifier must be declared as an extension in the Declarations Block; without a declaration the importer must reject the package (Exchange §16.3 alignment: unknown without declaration → reject; declared → accept if implemented, reject explicitly if not — never silent).

### 6.3 EKA Structured Text (RSF) — the canonical default

6.3.1 The v1 canonical default representation is **EKA Structured Text**: a metadata block (ordered field:value lines carrying non-Content unit metadata) followed by a structured body (level-1 section headings with prose). It is deliberately defined without naming Markdown, YAML, or any file format: the RSF does not couple Content to a specific representation technology. The reference implementation's repository serialization happens to use Markdown + frontmatter; the RSF borrows the convention's shape (metadata block + structured body sections) and defines it as the EKA Structured Text representation — the same tension resolution stated in Section 1.7, applied at the Content level.

6.3.2 **Well-formedness.** Well-formedness per artifact type family remains a conformance concern (Exchange §14.2 R9) and is evaluated on the canonical representation: the v1 required-section table follows the reference implementation's body conventions (`skeleton/docs/exchange/validation.md`, Aturan 9 — e.g., planning artifacts require Objective/Scope/Out of Scope sections; decision records require Context/Decision/Consequences/Alternatives Considered). Enforcement position is a serialization detail; verdict semantics are defined by R9 (Exchange §14.2 note).

6.3.3 **Content equality.** Content equality (Exchange §15.2) is decidable because the RSF declares its canonicalization: equality is evaluated on the canonical representation after canonicalization per Section 9.3 (normalized encoding, deterministic field order, normalized line endings). For non-canonical representations, equality is byte-level over the payload as carried — the implementation-declared canonicalization allowed by Exchange §15.2 (Appendix A, open question 3).

### 6.4 Non-canonical representations

6.4.1 A package may carry Content in a non-canonical representation only if the Representation Identifier is declared as an extension (6.2). The RSF defines conformance evaluation (R9) for the canonical representation only; for non-canonical representations, well-formedness verification is performed by importers that declare support for the representation, and the package is rejected explicitly by those that do not (Exchange §16.3: explicit, never silent).

### 6.5 Representation registry

6.5.1 The set of registered Representation Identifiers is an RSF-level registry. RSF v1 registers exactly one entry: `eka/structured-text/1`. New entries are registered through minor Serialization Versions (Section 11.2) and follow the terminology governance of the EKA ecosystem (Naming §12.1): proposals, impact checks, acceptance, registration — never parallel registries.

6.5.2 A package must not invent Representation Identifiers that collide with reserved RSF tokens; unknown identifiers without a matching declaration are rejected per 6.2.

## 7. Attachment Model

### 7.1 Role

Attachments (RSF) are supporting resources: diagrams, images, binaries, PDFs, other assets. They are never part of Content semantics — Content remains complete and Well-formed without them; attachments are supporting resources referenced from units, never inlined into Content. They carry no Identity, State Vector, Change Log, Classification, or Phase: they are not Exchange Units and realize no object model element (Appendix A).

### 7.2 Attachment identity

Each Attachment carries an Attachment ID (RSF): deterministic and unique within the package. The exporter derives it from the referring unit's Canonical Identity Form and the resource name (recommended rule); identical repository state must yield identical Attachment IDs. The Attachments Collection is ordered by Attachment ID.

### 7.3 Referencing

Units reference Attachments through the Attachment Reference Set (RSF) on the referring unit — a reference by Attachment ID, never an inline payload. The discipline mirrors Exchange §7.1–7.2: one-directional, recorded on the referrer, by identifier.

### 7.4 Referential integrity

- every referenced Attachment must be present in the package, or declared external in the Declarations Block (mirroring the External Reference discipline, Exchange §12.3);
- an undeclared dangling Attachment reference makes the package invalid (mirroring the self-consistency rule, Exchange §10.6);
- Attachment handling never changes exchange-conformance verdicts (R1–R9 apply to units); import decisions on Attachments are recorded in the Import Manifest (Exchange §11.3).

### 7.5 Integrity and metadata

- each Attachment carries a digest in the Integrity Block (Section 9.4);
- metadata: media type label and size (conceptual, informational) — never semantic.

### 7.6 Round-trip

Attachments are carried verbatim. They are outside the Exchange §15.2 invariant set, but the RSF determinism rule still applies: identical repository state → identical Attachments Collection (Section 9.1).

## 8. Manifest Model

### 8.1 Responsibilities

The Manifest realizes the Exchange Manifest (Exchange §10.2, §10.6):

| Responsibility | Rule |
|---|---|
| ordered unit list | every Unit Entry exactly once, ordered by Canonical Identity Form key (Exchange §6.3); 1:1 with units — no phantom entries, no missing units (§10.6) |
| scope type | declared Export Scope (Exchange §5.3) |
| version echo | Serialization Version, exchange format version, specification version — must equal the Package Header values (self-consistency, §10.6 discipline) |
| Package Identity Label | echo of the label |
| dependency summary | counts: units, Attachments, External References, extensions |
| integrity information | package-level digest and per-unit digests (Section 9.4) |
| closure declaration summary | echo of the Closure Declaration (Exchange §12.2): Export Scope and seed set reference |

### 8.2 Self-consistency

A package failing the Manifest ↔ unit 1:1 correspondence, or the Header ↔ Manifest version consistency, must be rejected before any validation phase (Exchange §10.6, §11.1).

## 9. Serialization Principles

### 9.1 Deterministic ordering

- Manifest, Unit Entries, External Reference declarations: canonical Identity key (Exchange §6.3);
- Change Log entries: occurrence order, deterministic tie-break = entry order (Exchange §8.2, §10.5);
- Relationship lists: (source Identity key, Relationship type, target Identity key) (Exchange §10.5);
- Attachments Collection: Attachment ID; Declarations records: canonical key.

Consequence: identical repository state must produce identical packages (Exchange §10.5).

### 9.2 Canonical identifiers

- the Canonical Identity Form (5.2) is the only identity spelling;
- content canonicalization: EKA Structured Text canonical form (6.3.3) for equality and digest purposes.

### 9.3 Encoding conventions (v1)

- UTF-8, no BOM;
- normalized line endings (LF);
- deterministic ordering of all maps and dictionaries — unordered maps must not appear in serialized output; fields serialize in declared fixed order or sorted order, never insertion order;
- the EKA Structured Text metadata block serializes fields in a fixed declared order.

### 9.4 Integrity expectations

- the package digest (SHA-256) covers everything in the package except the Integrity Block itself;
- per-unit digests over each Unit Entry's canonical serialization; per-Attachment digests;
- importers must verify integrity before Exchange §11.1 phase 1; verification failure → reject (Exchange §17.1);
- **mechanism choice: SHA-256.** Justification: the Exchange Specification must stay mechanism-neutral (§17.1), but the RSF is an implementation and may mandate a mechanism. SHA-256 is deterministic, widely available, collision-resistant for integrity purposes, and keeps v1 tooling simple. Future Serialization Versions may upgrade the mechanism (Section 11.2).

### 9.5 Future extensibility

- the extension area of the package is realized by the Declarations Block (Exchange §16.3);
- **unknown-field policy: reject by default.** A field unknown to a v1 implementer is a forward-compatibility violation and must be rejected explicitly — never silently ignored — unless covered by a declared extension (Exchange §16.3: forward compatibility is bounded by declarations);
- declared extensions may carry their own fields, parsed only by implementers that declare support.

## 10. Round-Trip Mapping

### 10.1 Invariant mapping

RSF round-trip = Exchange §15 guarantees applied to this projection. The invariant properties (Exchange §15.2) map to RSF elements as follows:

| Invariant (§15.2) | RSF element |
|---|---|
| Identity set equality | Manifest ordered unit list + Unit Entry Identities (Canonical Identity Form) |
| State vector equality | State blocks: exact values, all owned domains |
| Change Log equality | Change Log entries: domain, old value, new value, time, authority, in occurrence order |
| Content equality | Content payloads + declared Representation Identifier; equality decidable via EKA Structured Text canonicalization (6.3.3) |
| Relationship equality | Relationship sets: (source, type, target) by Identity |
| Classification equality | Classification fields: primary + secondary Knowledge Dimensions |
| Line integrity | instance grouping: Unit Entries grouped by (Namespace, Type, ID) across InstanceVersions |

### 10.2 Permissible differences

Per Exchange §15.4, the following may differ between exports of equivalent repositories: physical order/layout in storage; storage location and addressing; projection refresh timestamps; package metadata — creation timestamp and exporter identity in the Package Header. Nothing else may differ.

### 10.3 Projection exclusion

Projections are never exchanged (Exchange §8.5; transfer.md §3): ticket tables and container tables are regenerated at the target via Projection Refresh. Projection differences are therefore never a loss (Exchange §15.3). The RSF realizes this by excluding projection artifacts from Repository-scope exports except where the projection is itself an owned Artifact (e.g., a Ticket unit carried with its empty State Vector — Exchange §8.5).

### 10.4 Guarantees

Export → import into an empty target → export again yields the identical package up to 10.2 (Exchange §15.1). Re-importing the same package is a no-op — no new units, no State changes, no fabricated Change Log entries, no conflict errors (Exchange §15.5).

## 11. Compatibility Strategy

### 11.1 The RSF among projections

The RSF is one of many possible projections of the object model. At the contract level, an importer that accepts RSF packages must not reject other conformant serializations: the Exchange Contract is serialization-independent (Exchange §2.3, §4.4). RSF support is an implementation capability, not a contract obligation.

### 11.2 RSF versioning

The RSF versions itself with the Serialization Version (RSF), distinct from the exchange format version and the specification version (Exchange §9.1):

| Change | Version rule |
|---|---|
| additive: new optional elements, new registered Representation Identifiers, new extensions | minor: v1.0 → v1.1 |
| contract-breaking: structural change, semantic change to existing elements | major: v1.x → v2.0 |

Rule: RSF v1 packages must remain importable by all RSF v1.x implementers (mirroring Exchange §9.3 backward compatibility).

### 11.3 Negotiation

- the Package Header declares the exchange format version, the specification version, and the Serialization Version;
- an importer must reject an unsupported Serialization Version before validation phases (analogous to Exchange §9.2.2, applied at the parse stage);
- exchange format and specification versions are negotiated exactly per Exchange §9.2.

### 11.4 Capability Declaration

- the Capability Declaration (Exchange §4.5) defines six capability classes in v1; the RSF must not invent a seventh;
- the Serialization Version is carried in the Package Header and the Manifest only; exporters must not claim "RSF v1" support through the Capability Declaration in this version;
- a future minor of the Exchange Specification may add a serialization capability class via class 6 (future protocol capabilities, Exchange §4.5, §18.4). The RSF records that possibility here and does not act on it.

## 12. Reference Examples (illustrative — not normative serialization)

The sketches below illustrate the object model applied to typical exports. They are annotated outlines, not serialized content; they define no syntax. Placeholder Identities (`acme/...`) are illustrative.

### 12.1 Example A — single Artifact

- **Export Scope:** single instance (Exchange §5.3).
- **Package contents:** Package Header (exchange format version, specification version, Serialization Version, exporter identity, scope, creation timestamp); Manifest (scope type; one unit); one Unit Entry: Identity, Revision metadata, State Vector, Change Log, Content (EKA Structured Text), Relationships, Classification, Phase context where applicable; Attachments Collection: empty; Declarations Block: empty; Integrity Block: package digest + one unit digest.
- **Integrity summary:** one package digest (SHA-256) over all elements except the Integrity Block; one per-unit digest.
- **Round-trip note:** import into an empty target, then export, yields the same package up to §15.4 (creation timestamp, exporter identity).

### 12.2 Example B — Epic with closure

- **Export Scope:** Graph — the Closure of a seed set containing the Epic (Exchange §12.2), traversing `depends-on`, `derives-from`, `validates`, `supersedes`, and `amends` in both directions so supersession and amendment chains are complete.
- **Package contents:** the Epic; its work items (stories, technical stories, chores) via `depends-on`; the ADRs the Epic derives from via `derives-from`; the Review records that validate it via `validates`; each with full unit composition; one Attachment (architecture diagram referenced by the Epic through its Attachment Reference Set); Declarations Block: Closure Declaration (scope + seed set), External Reference Declarations (e.g., a Standard in another Namespace referenced by a non-Draft unit), no extension declarations; Integrity Block: package digest over all units and Attachments, plus per-unit and per-Attachment digests.
- **Integrity summary:** closure completeness is verified at import (Exchange §11.1 phase 8); every External Reference must resolve in the target or the referring unit must have Content State Draft (Exchange §7.4, §12.3).
- **Round-trip note:** identical repository state → identical Closure → identical package up to §15.4; re-import is a no-op.

### 12.3 Example C — complete repository

- **Export Scope:** Repository — all instances of all Lines.
- **Package contents:** Package Header; Manifest: N units ordered by Canonical Identity Form, including every InstanceVersion per Line (Line integrity, §10.1) — e.g., a planning Line carrying both v1 and v2 instances; empty-vector units present with empty State blocks; projections excluded — never exchanged (Exchange §8.5); Attachments: all repository assets; Declarations Block: Closure Declaration (Repository scope, empty seed), External Reference Declarations if any cross-repository references exist; Integrity Block: package + per-unit + per-Attachment digests.
- **Integrity summary:** the package digest covers the full unit set; Manifest ↔ unit 1:1 correspondence must hold (Exchange §10.6).
- **Round-trip note:** re-export after import is identical up to §15.4; projections are regenerated at the target via Projection Refresh (Exchange §8.5).

## 13. Implementation Recommendations

Conceptual guidance for future engines behind `eka export` and `eka import` (Naming §7.1). Not CLI design, not APIs.

### 13.1 Export pipeline

1. validate the source repository (Reference Validator, R0–R9) — Exchange §12.5;
2. resolve the Export Scope deterministically (Exchange §12.1);
3. compute the Closure (Exchange §12.2);
4. detect External References and require declarations (Exchange §12.3);
5. assemble the object model (Exchange §4.4);
6. project to the RSF (Sections 4–8);
7. compute integrity digests (Section 9.4);
8. emit atomically (single-file container or directory layout).

### 13.2 Import pipeline

1. verify package integrity (before Exchange §11.1 phase 1; §17.1);
2. run Exchange §11.1 phases 1–8 against the package;
3. map Unit Entries back to the object model (bijective inverse, Appendix A);
4. commit atomically (phase 9 — no partial writes);
5. post-commit revalidation (phase 10; rollback on failure).

### 13.3 Determinism checklist

- all collections ordered by canonical key (9.1);
- Canonical Identity Form used everywhere; no alternate spellings;
- normalized encoding (9.3) applied before digest computation;
- digests computed over canonical bytes only;
- no timestamps inside units — only package metadata and recorded Change Log entries;
- no environment-dependent or host-dependent values.

### 13.4 Failure handling principles

- no partial writes anywhere: export emits only complete packages; import aborts before commit on any blocking failure (Exchange §11.1) and rolls back on phase-10 failure;
- errors are reported explicitly, never silently swallowed (Exchange §16.3 explicit-rejection discipline);
- all import verdicts and warnings are recorded in the Import Manifest (Exchange §11.3).

### 13.5 Recommended conformance checks

- run the Reference Validator on the source repository before every export (13.1.1);
- validate produced packages against this document before release (Exchange §12.5);
- importers verify integrity before validation phases (9.4).

### 13.6 SDK guidance

- language bindings implement the RSF projection — parse and emit packages; they do not implement the Exchange Specification itself (Naming §7.4: `eka-sdk-<language>`);
- SDKs must not redefine the Conformance Rules; verdicts remain the Reference Validator's (P16);
- SDKs must keep package emission deterministic (13.3) and must reject unknown fields per 9.5.

## Appendix A — Bijective Projection Mapping

| Exchange object model element | RSF element |
|---|---|
| Contract Header (10.1) | Package Header |
| Manifest (10.2) | Manifest |
| Exchange Unit (10.3): Identity | Unit Entry: Identity (Canonical Identity Form) |
| Exchange Unit: Revision metadata | Unit Entry: Revision metadata |
| Exchange Unit: State Vector | Unit Entry: State block |
| Exchange Unit: Change Log | Unit Entry: Change Log entries |
| Exchange Unit: Content | Unit Entry: Content payload + Representation Identifier |
| Exchange Unit: Relationship set | Unit Entry: Relationship set |
| Exchange Unit: Classification | Unit Entry: Classification fields |
| Exchange Unit: Phase context (8.4) | Unit Entry: Phase context |
| Declarations (10.4) | Declarations Block |
| integrity data (17.1) | Integrity Block |
| — (no object model element) | Package Identity Label, Attachments Collection, Attachment Reference Set, Serialization Version (RSF support elements; Section 3) |

## Appendix B — RSF v1 Conventions Summary

| Convention | v1 value |
|---|---|
| Serialization Version | 1 (minor additions → 1.x; breaking changes → 2.0) |
| Package extension | `.ekapkg` single-file; none for directory layout; convention, not contract |
| Canonical Identity Form | `<namespace>/<type>:<id>:<instance-version>` — full tuple, no defaulting |
| Default Content representation | EKA Structured Text (`eka/structured-text/1`): metadata block + structured body |
| Integrity digest | SHA-256: package-level (excluding the Integrity Block), per-unit, per-Attachment |
| Encoding | UTF-8, no BOM, LF line endings, deterministic field order, no unordered maps |
| Unknown fields | reject by default; accepted only under a declared extension |
| Attachment referencing | Attachment Reference Set on the referring unit, by Attachment ID; external declaration mirrors Exchange §12.3 |
| Capability Declaration | not used by the RSF in v1; Serialization Version lives in Package Header and Manifest only |

---

*End of RSF — EKA Reference Serialization Format v1.0 (Reference, not normative).*
