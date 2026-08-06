# EKA Naming and Terminology Specification v1.0

| Field | Value |
|---|---|
| **Version** | 1.0 |
| **Status** | Ratified |
| **Anchor** | EKA Core Specification v1.0 (Ratified) |
| **Zone** | standard/ |
| **Scope** | Official naming and terminology for the EKA ecosystem: product identity, specification families, reference components, tooling, repository naming, canonical term list. Meta-specification — governs how the standard names things; not an architecture-domain Family. |

**Reading this document:** capitalized terms refer to canonical definitions (EKA Core Specification v1.0 Section 3 + glossary) or to terms defined in this document (Section 2.3). "must" = binding requirement. "should" = recommendation within the bounds of the contract. "may" = option within the bounds of the contract. The EKA Core Specification remains the authoritative conceptual model; this document governs naming and terminology only.

## 1. Purpose

1.1 This document is the single normative source for how the EKA ecosystem names things: the official product identity, the naming pattern for specification Families, the official terms for reference components and tooling, the repository naming convention, and the canonical term list with deprecated aliases.

1.2 It exists because terminology drift is contract drift. Identity, State Vector, Exchange Layer, Knowledge Dimension, and every other canonical term are referenced by implementation code, validation rules, traceability matrices, and future Knowledge OS consumers. A stable vocabulary is a stable contract (Section 3.1).

1.3 It resolves, at the standard level, inconsistencies already present in the ecosystem: orphan labels, three-way rule naming, zone names that were only implied, and mixed-language term usage.

1.4 This document does not amend the EKA Core Specification. It names things; the Core Specification defines them.

## 2. Scope

### 2.1 In scope

- official product identity and branding conventions (Section 4);
- specification Family naming, versioning, and registration (Section 5);
- official terms for reference components (Section 6);
- tooling naming philosophy (Section 7);
- repository naming convention (Section 8);
- canonical term table and deprecated terminology (Sections 9–10);
- migration of existing documentation (Section 11);
- governance and evolution of the vocabulary itself (Section 12).

### 2.2 Out of scope

- any change to the conceptual model, invariants, or contracts of the EKA Core Specification or the EKA Exchange Specification;
- serialization formats, directory layout, or file naming inside per-project repositories (canonical 12.3);
- the content of the Conformance Rules themselves (R0–R9); this document names them, it does not redefine them;
- naming of Artifacts inside user repositories (their Identity is governed by canonical Sections 6 and 10).

### 2.3 Definitions of new terms

| Term | Definition |
|---|---|
| **Family** | an architectural domain of the EKA standard expressed as one specification document lineage: "EKA \<Family\> Specification". Families are domains, not documents (Section 5.1). |
| **Family Registry** | the normative registry of Family names, types, anchors, and dependencies, maintained in this document (Section 5.6). |
| **Meta-specification** | a specification that governs other specifications (naming, governance, process) rather than an architecture domain. This document is a meta-specification. |
| **Anchor Reference** | the formal name of the established marker convention (previously "Anchor EKA"): a line in a convention document declaring its anchor into the EKA standard (Section 6.4). |
| **Conformance Rules** | the canonical rule set R0–R9: R1–R9 are the nine numbered rules of the EKA Exchange Specification §14.2; R0 is the structural rule defined by the Reference Validator. |
| **Conformance Suite** | the Conformance Rules together with their automated tests and Test Fixtures; the executable form of conformance requirements. |
| **Reference Validator** | the conformance-checking tool: the `eka` command-line tool and its `conformance/` engine. |
| **Test Fixture** | a sample conformant or non-conformant input used to verify the Reference Validator. |
| **Example Repository** | a sample populated repository used for onboarding (future component). |
| **canonical text** | the binding text of a specification (e.g., the verbatim copy held in the Standard Zone). Adjective usage; not a specification name. |
| **EKA ecosystem** | the official EKA standard, its Families, reference components, tooling, and repositories governed by this document. |

## 3. Naming Principles

| # | Principle | Rationale |
|---|---|---|
| N1 | **Naming is architecture** | Stable vocabulary = stable contracts. Renaming a term ripples through Identity, validation, traceability, and external consumers. Every name change is a contract change. |
| N2 | **Clarity over cleverness** | Names must be descriptive and boring. No wordplay, no coined acronyms, no metaphor that needs explanation. |
| N3 | **Minimal abbreviations** | The only abbreviations are "EKA" and the Family tokens. No per-feature acronyms; no capitalized initialisms inside terms. |
| N4 | **Technology-neutral** | Canonical names must not name a storage medium, language, or format. The only permitted technology suffix is the language token on SDKs (Section 7.4). |
| N5 | **Human and AI-agent readable** | Names must be deterministic, unambiguous, and greppable. The "eka" prefix is reserved so that every official name is uniquely identifiable in any namespace. |
| N6 | **Ten-year horizon** | Vocabulary evolves additively: new terms are registered, old terms are deprecated with mapping, never silently redefined. Versioning is additive (Section 5.3). |
| N7 | **Naming precedence** | New specification → Family pattern + registration (Section 5). New tool → `eka <verb>` or `eka-sdk-<language>` (Section 7). New repository → `eka-<component>` (Section 8). New term → terminology review; extend the glossary, never fork it (Section 12). |

## 4. Official Product Identity (EKA)

### 4.1 Official name and abbreviation

| Aspect | Rule |
|---|---|
| Official name | **Engineering Knowledge Architecture** |
| Official abbreviation | **EKA** — capitalized, no periods, no decoration |
| Versioned form | "EKA v1.0" (version applies to the standard corpus) |
| Ecosystem form | "EKA ecosystem" |
| Forbidden forms | "Eka", "E.K.A.", "EKA!", lowercase "eka" alone in prose |

4.2 "EKA" must not be written with periods ("E.K.A.") or in mixed case ("Eka"). The lowercase form "eka" collides with the Indonesian word *eka* ("one") and with the repository directory name `eka-specification`; in prose, the abbreviation is always capitalized.

4.3 Branding conventions:

- **Sentence case in prose**: "Engineering Knowledge Architecture" and "EKA" are written in normal case inside running text.
- **Capitalized in titles**: headings and titles may use the full name or "EKA" as given.
- **EKA as noun modifier**: official text must use "EKA Artifact", "EKA Core Specification", "EKA ecosystem" — never the possessive "EKA's artifact".
- **Pronunciation**: not specified. The written form is canonical; no phonetic guidance is given.

4.4 In localized prose (e.g., Indonesian convention documents), "EKA" is used unchanged — abbreviations are language-neutral. The full name may be translated contextually only after the official form has appeared (Section 9.4).

## 5. Specification Families

### 5.1 Family concept

A **Family** is an architectural domain of the EKA standard, expressed as one specification document lineage. Families are domains, not documents: a Family may contain exactly one ratified specification per version at a time, and a single document belongs to exactly one Family.

This document is the exception: it is a **meta-specification** (Section 2.3), not a Family. It governs naming for all Families but does not itself constitute a knowledge domain.

### 5.2 Family naming pattern

- Title: **EKA \<Family\> Specification v\<major\>.\<minor\>** — e.g., "EKA Core Specification v1.0", "EKA Exchange Specification v1.0".
- The Family token is a single descriptive word (or hyphenated compound): `Core`, `Exchange`, `Conformance`, `Runtime`, `Interoperability`.
- "Specification" is part of the title; the document is a Specification, the whole corpus is the standard (Section 9.5).
- Shorthand: after first use of the full name, "the \<Family\> Specification" is acceptable (e.g., "the Exchange Specification").

### 5.3 Version notation

| Change | Version rule |
|---|---|
| Additive change within a Family (new sections, clarifications, subordinate terms) | **minor** bump: v1.0 → v1.1 |
| Contract-breaking change (amends Core invariants, changes verdict semantics, redefines a canonical term) | **major** bump: v1.x → v2.0 |
| Family version vs Core version | Families version independently, but every Family version declares the Core version it anchors to (5.5) |

### 5.4 Document identity pattern

- File slug: `eka-<family>-specification-v<major>.<minor>.md`.
- The Core Family omits the token: `eka-specification-v1.0.md`.
- The meta-specification follows the same pattern: `eka-naming-and-terminology-specification-v1.0.md`.

| Document | Slug |
|---|---|
| EKA Core Specification v1.0 | `eka-specification-v1.0.md` |
| EKA Exchange Specification v1.0 | `eka-exchange-specification-v1.0.md` |
| EKA Naming and Terminology Specification v1.0 (meta) | `eka-naming-and-terminology-specification-v1.0.md` |
| Future: EKA Conformance Specification v1.0 | `eka-conformance-specification-v1.0.md` |

### 5.5 Anchor declaration

Every Family specification must declare its anchor to the Core in the front matter:

- an **Anchor** field naming the anchored document and version (the Exchange Specification pattern: "Engineering Knowledge Architecture (EKA) Canonical Specification v1.0 (Ratified)");
- a **Relationship to the EKA Core Specification** clause stating: what it operationalizes (anchored sections), that it does not amend the Core, and which subordinate terms it introduces.

A Family specification must not define or redefine canonical terms (EKA Core Specification Section 3 + glossary); it may only introduce subordinate terms, declared explicitly as such.

### 5.6 Family Registry and registration

| Family | Type | Status | Anchor | Dependencies |
|---|---|---|---|---|
| Core | domain | Ratified (v1.0) | self | — |
| Exchange | domain | Ratified (v1.0) | Core §13, §5.4.7, P13 | Core |
| Conformance | domain (candidate) | Not registered | Core (future) | Core, Exchange |
| Runtime | domain (candidate) | Not registered | Core (future) | Core |
| Interoperability | domain (candidate) | Not registered | Core (future) | Core, Exchange |
| Naming and Terminology | meta | Ratified (v1.0) | Core §3 | Core |

Registration process:

1. **Proposal** — any contributor may propose a Family by submitting: name (single token), type (domain or meta), scope statement, anchor into the Core, dependencies on other Families, and proposed slug.
2. **Review** — the EKA maintainers evaluate: does the domain justify a separate specification lineage, or does it belong in an existing Family? Does it weaken invariants (canonical 14.2.1)?
3. **Ratification** — accepted proposals are entered into the Family Registry above. Registered names are stable; a Family is never renamed, only versioned.

### 5.7 Status of "Canonical Specification"

- **EKA Core Specification** is the official name of the canonical conceptual model document (currently `standard/eka-specification-v1.0.md`).
- **"Canonical Specification"** is deprecated as a document name; it must not be used as a title or as a Family name.
- **"canonical"** remains in force as an adjective meaning "binding text": "the canonical text of the EKA Core Specification", "canonical terms", "canonical serialization" (canonical 6.2.6). Adjective usage is not a specification name.

## 6. Reference Components

### 6.1 Official terms

| Component | Official term | Role |
|---|---|---|
| A serialization demonstrating conformance | **Reference Implementation** | one concrete serialization of the standard (canonical 1.3, 16.2); this repository's role |
| The conformance-checking tool | **Reference Validator** | the `eka` command-line tool plus its `conformance/` engine; deterministic, read-only, mechanical (EKA Exchange Specification §14.1) |
| Rules + tests + fixtures | **Conformance Suite** | the executable form of conformance requirements: the Conformance Rules (R0–R9), their automated tests, and Test Fixtures |
| Sample conformant/non-conformant inputs | **Test Fixture** | inputs used to verify the Reference Validator (`conformance/testdata/`) |
| Sample populated repository for onboarding | **Example Repository** | future component; a populated, conformant repository demonstrating the Reference Implementation in use |

### 6.2 Relationships

- The Reference Implementation **ships with and uses** the Reference Validator; a repository that cannot pass its own validator is not a Reference Implementation.
- The Reference Validator is the **canonical executable implementation** of the Conformance Suite (P16: verdicts identical across implementations).
- The Conformance Suite is the **executable form of the Conformance Rules** (R0–R9); a future EKA Conformance Specification Family (5.6) will make this binding explicit at the standard level.
- Test Fixtures belong to the Conformance Suite; they are never normative inputs to any production repository.

### 6.3 Capitalization

Official terms are capitalized when used as terms: **Reference Implementation**, **Reference Validator**, **Conformance Suite**, **Test Fixture**, **Example Repository**. Lowercase is permitted only inside path, slug, or code tokens (`reference/`, `conformance/`, `eka-reference`). In localized prose the terms may be translated; first occurrence must carry the official English term (Section 9.4).

### 6.4 Anchor Reference convention

The established marker line `> Anchor EKA: <target>` is formalized as the **Anchor Reference** convention:

- the convention is named **Anchor Reference**;
- the marker line text stays **`> Anchor EKA:`** — established, machine-greppable, and present in the Skeleton Zone; it is not renamed;
- every convention document should carry one Anchor Reference declaring the standard section or principle it anchors to.

## 7. Tooling Terminology

### 7.1 Naming philosophy

One root tool, one name: **`eka`**. Capabilities are subcommands. There must be no satellite binaries (`eka-import`, `eka-export`).

| Aspect | Rule |
|---|---|
| Root binary | `eka` |
| Subcommands | **verbs**: `validate`, `init`, `export`, `import` (implemented); `diagnose`, `sync`, `format`, `graph` (future candidates). `completion` (bash/zsh/fish/powershell) is provided by the Cobra framework — not a core EKA command |
| Philosophy | commands are verbs; tools are nouns (the tool is "the Reference Validator"; it runs `eka validate`) |
| Determinism | output byte-identical across runs; exit codes 0/1/2 as documented in `reference/cli.md` |
| Read-only | validation never modifies the repository (EKA Exchange Specification §14.1.3) |

7.2 The roadmap name `doctor` violates the verb rule. It is not implemented; it must be renamed **`diagnose`** before implementation. No compatibility cost.

7.3 Engine packages are named after the domain they implement, not after verbs: `conformance/` (validation engine), `exchange/` (import/export engine), etc. Package name = domain noun.

7.4 SDKs and language bindings use one pattern only: **`eka-sdk-<language>`** (`eka-sdk-go`, `eka-sdk-python`). A separate `eka-<language>` pattern must not exist: `eka-go` would be ambiguous between "Go SDK" and "Go component" under Section 8.

7.5 Canonical tool names are technology-neutral. The only permitted technology token is the language suffix on SDKs, because an SDK is inherently language-specific. The Reference Validator is described as "the canonical executable form of the Conformance Rules" — not "the executable form of the specification" in general.

## 8. Repository Naming Convention

### 8.1 Official repositories

All official EKA ecosystem repositories use the org-level prefix **`eka-`**:

| Repository | Content | Status |
|---|---|---|
| `eka-specifications` | the Standard Zone (all Families + this document) | future split candidate |
| `eka-reference` | the Reference Implementation (Skeleton Zone + Reference Zone + validator) | future split candidate |
| `eka-cli` | the Reference Validator as standalone tooling | future split candidate |
| `eka-sdk-go` | EKA SDK for Go | future |
| `eka-examples` | the Example Repository | future |
| `engineering-knowledge-architecture` | **this repository** (current monorepo: Standard + Skeleton + Reference + validator) | current |

### 8.2 This repository

- **Keep the current name** `engineering-knowledge-architecture`. It carries the full official name, matches the Go module path, and avoids churn while the ecosystem is a monorepo.
- **Optional future rename**: when (and only when) the monorepo splits per Section 8.1, the spec home may become `eka-specifications`. GitHub redirects the old name; module consumers are unaffected (8.3).
- The `eka-` prefix is reserved for official EKA ecosystem repositories. External, per-project repositories (end users serializing EKA) get **no naming rules** — their repositories are named by their own organizations. They must not claim the `eka-` prefix or the name "Reference Implementation"; that term is reserved for the official component (Section 6).

### 8.3 Go module path stability

The module path `github.com/maleolabs/engineering-knowledge-architecture` is **public API** and must remain stable regardless of repository renames. A future vanity import path (e.g., `eka.dev/...`) may be added as an alias via standard module tooling, but the current path must never break (`go install .../cmd/eka@latest` is a documented installation path, `reference/cli.md`).

## 9. Terminology Standard

### 9.1 Canonical term table

| Term | Definition (short) | Deprecated aliases | Notes |
|---|---|---|---|
| **Engineering Knowledge Architecture** | the EKA standard: canonical conceptual model for engineering knowledge — Artifact, Identity, State, taxonomy, layers, exchange contracts | "Project Docs Template", "Documentation Template", "Repository Standard" | official name, Section 4 |
| **EKA** | official abbreviation; capitalized, no periods | "Eka", "E.K.A.", "EKA!" | Section 4.2 |
| **EKA Core Specification** | the Family specification defining the canonical conceptual model | "Canonical Specification" (as a name) | "canonical" stays an adjective (5.7) |
| **EKA Exchange Specification** | the Family specification defining the Exchange Contract | "Exchange Specification" (as full name) | shorthand after first use is acceptable |
| **Exchange Contract** | the normative contract body **defined by** the EKA Exchange Specification: round-trip, idempotency, referential integrity, schema versioning, conformance | — | the specification is the document; the contract is its content (9.2) |
| **Exchange Contracts** | generic plural for the exchange contract family (canonical §12.2) | — | consistent with the singular term |
| **Conformance Rules** | the canonical rule set **R0–R9** | "validation rules"; "aturan validasi konformitas" (as a name) | R1–R9 = the nine numbered rules (Exchange §14.2 / Aturan 1–9); R0 = structural rule of the Reference Validator (9.3) |
| **Conformance Suite** | Conformance Rules + automated tests + Test Fixtures | — | Section 6 |
| **Reference Implementation** | a serialization demonstrating conformance; this repo's role | "Project Docs Template", "Documentation Template" | capitalized as a term (6.3) |
| **Reference Validator** | the conformance-checking tool (`eka` + `conformance/`) | "CLI EKA" (as a term) | "CLI" stays descriptive, not a term |
| **Test Fixture** | sample conformant/non-conformant input | — | Section 6 |
| **Example Repository** | sample populated repository for onboarding (future) | — | Section 6 |
| **Standard Zone** | the `standard/` zone: canonical texts of all specifications | "zona standard" (as official name) | directory stays `standard/` (9.4) |
| **Skeleton Zone** | the `skeleton/` zone: copyable project serialization | "zona skeleton" (as official name) | directory stays `skeleton/` |
| **Reference Zone** | the `reference/` zone: meta-documentation of the Reference Implementation | "zona reference" (as official name) | directory stays `reference/` |
| **Artifact** | canonical: entity with Identity, Content, State Vector, Relationship | — | glossary term, used verbatim |
| **Anchor Reference** | the formal name of the `> Anchor EKA:` marker convention | "Anchor EKA" (as the convention's name) | marker text unchanged (6.4) |
| **Conformance Rule** | one rule of the Conformance Rules set | "Aturan" (as canonical term) | Indonesian "aturan" remains the translation (9.4) |

### 9.2 Exchange Specification vs Exchange Contract

- **EKA Exchange Specification** = the document (the Family lineage).
- **Exchange Contract** = the normative contract body defined by that document (unit, identity representation, state representation, versioning, package structure, import/export/sync semantics, round-trip guarantees).
- "Exchange Contract v1" (exchange spec §1.1) is acceptable shorthand for "the v1 Exchange Contract defined by the EKA Exchange Specification v1.0"; it must not be used as the document's name.

### 9.3 Conformance Rules count

The Conformance Rules set is **R0–R9** (ten rules). The phrase "nine rules" refers to **R1–R9** only. R0 is the structural rule (artifact-rule parsing, frontmatter well-formedness) defined by the Reference Validator; it is not one of the nine rules of EKA Exchange Specification §14.2. Documentation must state which set it means; the mapping is:

| validation.md (Aturan) | Conformance Rule | Source |
|---|---|---|
| Aturan 1–9 | R1–R9 | EKA Exchange Specification §14.2 |
| — (structural bucket) | R0 | Reference Validator definition (`conformance/`) |

### 9.4 Language and localized prose rule

- All normative specifications and convention documents are written in **English**; English is the canonical language of the EKA ecosystem.
- Localized translations may exist as accessibility aids, with two obligations:
  1. the first occurrence of any defined term in a translation carries the official English term: "Conformance Rule (aturan konformitas)", "Standard Zone (zona standard)";
  2. a translation is never used as a canonical name in English text; the English document remains authoritative.

### 9.5 "Standard" vs "specification" usage

- **standard** (uncountable): the whole EKA normative corpus — "the EKA standard", "the standard mandates".
- **Specification** (countable): one Family document — "the EKA Exchange Specification".
- Rule: a single document is never "a standard"; the whole corpus is never "a specification".

## 10. Deprecated Terminology

| Deprecated term | Official term | Note |
|---|---|---|
| Project Docs Template | Reference Implementation | orphan label; root README title (drift #1); must not reappear |
| Documentation Template | Reference Implementation | not found in the current repository; forbidden |
| Repository Standard | — | not found in the current repository; previously used elsewhere; forbidden |
| Canonical Specification (as a name) | EKA Core Specification | "canonical text" adjective usage remains (5.7) |
| validation rules | Conformance Rules | "validation.md" file name stays (11.5) |
| Aturan 1–9 (as canonical term) | Conformance Rules R1–R9 | "Aturan" stays as Indonesian translation (9.4) |
| Exchange Specification (as full title) | EKA Exchange Specification | shorthand acceptable after first use |
| Anchor EKA (as convention name) | Anchor Reference | marker text `> Anchor EKA:` unchanged (6.4) |
| Knowledge Object | — | rejected in EKA Exchange Specification §5.2; must never reappear as a canonical term |
| CLI EKA (as a term) | Reference Validator | "CLI" remains descriptive prose |
| Eka / E.K.A. / EKA! | EKA | Section 4.2 |
| zona standard / zona skeleton / zona reference (as official names) | Standard Zone / Skeleton Zone / Reference Zone | translations acceptable in localized prose (9.4) |

## 11. Migration Recommendations

The following are recommendations, not mandates. File and directory renames in this repository are avoided unless strongly justified: links, traceability matrices, and validator paths depend on them.

| # | Action | Rationale | Cost |
|---|---|---|---|
| 1 | Root README title → "Engineering Knowledge Architecture (EKA) v1.0 — Reference Implementation" | removes the orphan "Project Docs Template" label | title-level only; zero link impact |
| 2 | `validation.md`: keep headings "Aturan 1–9"; add a mapping marker (Conformance Rule ↔ Aturan) in the header block | validator does not parse headings; labels stay localized per 9.4 | zero conformance impact |
| 3 | Prose zone names → official capitalized names on first use; directory paths unchanged | resolves the implied-zone naming gap | copy-level |
| 4 | GitHub repository rename: **defer**; revisit only at an ecosystem split (8.2) | module path and links are the contract; rename adds churn without value now | none now |
| 5 | No file or directory renames in this repository | stable links; traceability matrices reference exact paths | — |
| 6 | `glossary.md`: add a pointer that terminology governance lives in this document; do not fork definitions | glossary remains the source for Core concepts; this document is normative for naming terms via itself | one note |
| 7 | Do **not** migrate: `validation.md`/`transfer.md` filenames, the `> Anchor EKA:` marker text, Indonesian prose in skeleton/reference documents | churn without contract value; 9.4 already governs | — |

## 12. Governance and Future Evolution

### 12.1 Terminology review

- A new term enters official vocabulary only through terminology review: proposal → impact check (where is the term already used? which contracts mention it?) → acceptance → registration in Section 9 (or the glossary, for Core concepts).
- Terms are extended, never forked: a new term must be added to the canonical term table or the glossary; creating a parallel glossary is forbidden (N7).
- Deprecation is the only way to retire a term; deprecated terms keep their mapping (Section 10) indefinitely.

### 12.2 Evolution rules

- This document versions independently under the pattern it defines (5.3): minor = additive terminology changes; major = redefinition of a canonical term or a change to the naming pattern itself.
- A major revision of the EKA Core Specification triggers re-evaluation of this document's anchor (Section 5.5) and of every Family's anchor declaration.
- Naming changes that affect machine-consumed identifiers (rule IDs, REQ IDs, file slugs, module paths) must follow the stability rules of the artifacts they touch: R0–R9 and REQ-nnn IDs are stable and never reused (CONTRIBUTING.md); slugs change only via the Family registration process.

### 12.3 Conflict resolution

- On any naming conflict between documents, this document wins for naming; the EKA Core Specification wins for conceptual definition; the EKA Exchange Specification wins for exchange-contract interpretation.
- Ambiguity between the two is resolved by the EKA maintainers and recorded as an interpretation in the appropriate reference document.

---

*End of Naming and Terminology Specification — EKA Naming and Terminology v1.0 (Ratified, meta-specification).*
