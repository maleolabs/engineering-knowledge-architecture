# Conformance Traceability Matrix — EKA v1.1

| Property | Value |
|---|---|
| **Status** | Ratified — official EKA v1.0 component |
| **Version** | v1.0 |
| **Document type** | Convention document (not an artifact) — zone `reference/` |
| **Related** | [`conformance-notes.md`](conformance-notes.md) (interpretations + gaps), [`../skeleton/docs/exchange/validation.md`](../skeleton/docs/exchange/validation.md) (Conformance Rules R0–R12), [`../standard/eka-specification-v1.1.md`](../standard/eka-specification-v1.1.md) (canonical standard) |

> **Governance rule (formal).** This matrix is the **single source of truth** for EKA conformance coverage. The matrix **MUST be updated in the same Pull Request** as changes to the specification, the Conformance Rules (`validation.md`), the validator implementation (`conformance/`, `cmd/eka/`), or tests — and conversely, the matrix is never edited without a related change. See [`../CONTRIBUTING.md`](../CONTRIBUTING.md).

> **Domain-aware rules (EKA v1.1).** This matrix covers the **full rule surface R0–R12**. The domain-aware rules **R10–R12** (Core v1.1 §8.1 — Engineering Domains and Knowledge Stratification) arrived with the EKA v1.1 **taxonomy evolution** (Core §14.2.3 taxonomy-extension governance: core closed, taxonomy open; documented in the [migration report](migration-report-engineering-domains-v1.1.md)); their rule rows are added in Section 2 below. They are implemented in `conformance/rules_domain.go` (+ `conformance/domain.go`) with automated coverage (`conformance/rules_domain_test.go`, `conformance/domain_test.go`, `conformance/validate_test.go`). Semantics and verdicts: `skeleton/docs/exchange/validation.md` + `conformance-notes.md`.

---

## 1. Layer model

This matrix traces conformance coverage through **5 layers**, from requirements to automated evidence:

| # | Layer | Identifier | Example | Source of truth |
|---|---|---|---|---|
| 1 | **Engineering Requirement** | `REQ-nnn` (REQ-001..REQ-019) | REQ-002 Identity uniqueness | Section 3 of this document |
| 2 | **Specification** | Anchor `§` (section/principle number) | §6.2.2, P3 | `standard/eka-specification-v1.1.md` |
| 3 | **Conformance Rule** | `Rn` (R0, R1–R12) | R1 | `skeleton/docs/exchange/validation.md` (R1–R12) + R0 (structural, defined in `conformance/`); R10–R12 per Core v1.1 §8.1 |
| 4 | **Implementation** | `file:func` (package-relative) | `conformance/rules.go:rule1` | Go code `conformance/` + `cmd/eka/` |
| 5 | **Automated Test** | function name `TestXxx` | `TestRule2ExactCounts` | `*_test.go` in `conformance/` + `cmd/eka/` |

**Identifier conventions (deterministic, automation-ready):**

- **REQ ID** — `REQ-<3 digits>`; stable, **never reused**. Removed/replaced requirements keep their IDs reserved; new requirements get the next ID.
- **Rule ID** — R0 (structural) + R1–R12, fixed from `validation.md` (R10–R12 defined by Core v1.1 §8.1); no renumbering.
- **Implementation** — `path/file.go:funcName` relative to the repo root; helper functions serving one rule are written in parentheses after their main function.
- **Test** — Go function name (`TestXxx`); one function = one row in Section 4.
- **Coverage Status** — fixed enumeration values: `Enforced (tested)` (rule implemented + tested), `Governance-only (uncovered)` (normative in the spec, not mechanically enforced), `Partially enforced` (part of the surface enforced).

---

## 2. Main matrix

One row per Conformance Rule. The `Enforced (tested)` status applies to all 13 rules: every rule has a Go implementation and automated test coverage.

| Rule | Requirement ID(s) | Spec Anchor | Implementation | Automated Tests | Coverage Status | Notes |
|---|---|---|---|---|---|---|
| **R0** (structural) | REQ-001 | §3 (Artifact), §5.1 | `conformance/artifact.go:analyzeFile` | `TestAnalyzeNoFrontmatterIsConventionDoc`, `TestAnalyzeFrontmatterWithoutTypeIDIsConventionDoc`, `TestAnalyzeUnterminatedFrontmatter`, `TestAnalyzeBrokenYAML`, `TestAnalyzeTypeXorID`, `TestAnalyzeValidArtifact`, `TestAnalyzeMissingIdentityFields`, `TestAnalyzeNonIntVersion`, `TestAnalyzeInvalidDate`, `TestAnalyzeUnknownType`, `TestAnalyzeChangeLogNotList`, `TestAnalyzeMalformedChangeLogEntry`, `TestUnknownTypeIsStructural` | Enforced (tested) | Structural bucket before the numbered rules (not one of the 9 rules). Artifact rule: `type` + `id` ⇒ Artifact; `type` XOR `id` = malformed. Unknown type token = structural error (interpretation #25); non-integer `instance-version` = error (#26). An artifact that fails classification → rules R2, R3, R4, R6, R7 skipped for that file; R1 still applies (identity still indexed); R5 still checks references and reports unknown tokens; R8/R9 not applicable to unknown types. |
| **R1** | REQ-002 | §6.2.2, P3 | `conformance/rules.go:rule1` (+ `conformance/validate.go:identityKey`/`buildIndex`) | — (no dedicated unit test; covered via `TestInvalidFixtures` → fixture `invalid-dup-identity`) | Enforced (tested) | Duplicate `(namespace, type, id, instance-version)` = error. ID unique within `(namespace, type)`; InstanceVersion unique within a line. |
| **R2** | REQ-003 | §6.4, P9 | `conformance/rules.go:rule2` (+ `conformance/filename.go:parseFilename`) | `TestRule2ExactCounts` (+ infra: `TestParseFilename`, `TestParseFilenameEmpty`; fixture `invalid-filename`) | Enforced (tested) | Filename token == frontmatter `type`; `-v<nn>` suffix == `instance-version`; `-v<nn>` only on `scp-`/`plan-` and MANDATORY (including v1). Digit count free (`-v1`/`-v01` valid, #16). **Documented gap:** the id part of the filename is NOT matched against the frontmatter `id` (#17) — the filename is a projection (ADR-001), the true Identity is in frontmatter. |
| **R3** | REQ-004 | §7.2 | `conformance/rules.go:rule3` (+ `conformance/state.go:contentStateVariant`/`domainValues`) | `TestPhaseValidation` (+ infra: `TestPhaseValueSet`; fixture `invalid-state-value`) | Enforced (tested) | State field values ∈ domain value set. Content-state variants per type family: living (`draft/review/approved/amended`), ADR (`proposed/accepted/superseded`), decision (`draft/accepted/superseded`). `phase` only on `scp-`/`plan-` and ∈ phase value set. |
| **R4** | REQ-005 | §7.4, §10 | `conformance/rules.go:rule4` | — (no dedicated unit test; covered via `TestInvalidFixtures` → fixture `invalid-ownership`) (+ infra: `TestOwnedSets`) | Enforced (tested) | State fields on the file == its type's owned set. Absent owned field = error (interpretation #2); present non-owned field = error; `tkt-` carries an empty state vector. Type→state binding from §10. |
| **R5** | REQ-006 | §6.2.7, §13.2.3 | `conformance/rules.go:rule5` (+ `parseReference`, `resolve`) | `TestRule5DraftSeverity`, `TestRule5CrossNamespaceAndVersionResolution`, `TestRule5VersionedReferenceToMissingInstance` (+ infra: `TestParseReference`, `TestParseReferenceCrossNamespace`; fixture `invalid-reference`) | Enforced (tested) | Malformed reference = always error; unresolved on `content-state: draft` = warning, otherwise (including without content-state) = error (interpretation #10). Self-reference = error. **Bare-id accepted** as a line reference in the referrer's namespace+type (#9 — 6/7 real ADRs use this form). Cross-namespace format `<ns>/<type>:<id>[:<ver>]` (#12). **Not enforced:** the bidirectional reference convention (#11, documented gap). |
| **R6** | REQ-007 | §8, P15 | `conformance/rules.go:rule6` (+ `dimensionFolderFor`, `dimensionList`) | `TestDimensionFolderResolution` (+ infra: `TestDimensionTokens`; fixture `invalid-dimension`) | Enforced (tested) | Knowledge artifact: `dimension` == home folder; home folder = **nearest ancestor** whose name ∈ 12 dimensions (#13); not found = error. `dimension` mandatory on knowledge artifacts (#14), forbidden on `ctr-`/`tkt-`/`ses-`; `dimensions-secondary` validated as token only (#15). |
| **R7** | REQ-008 | §5.2, P7 | `conformance/rules.go:rule7` (+ `entriesForDomain`, `indexOfEntry`; `conformance/state.go:isLegalTransition`) | `TestRule7ExactCounts` (+ infra: `TestIsLegalTransition`; fixture `invalid-changelog`) | Enforced (tested) | Every owned domain has a change-log entry; last entry == current field value; transitions legal; initial `from: "-"` legal for all domains, `to: "-"` = invalid (#5). **Not enforced:** chain contiguity `entry[i].from == entry[i-1].to` (#4 — static snapshot). Execution State strictly adjacent; other domains forward-only without required adjacency (#6). `phase` entries: value-set + well-formedness only (#7). `tkt-` (empty state vector) may omit the change-log (#8). |
| **R8** | REQ-009 | §7.4, P6 | `conformance/rules.go:rule8` (+ `workItemsTable`, `compareWorkItemsTable`, `resolveWorkItemCell`, `hasProjectionHeader`, `splitTableRow`) | `TestRule8TicketHeaderAndDerivation`, `TestRule8ContainerTableMismatchIsWarningOnly` (+ fixture `invalid-projection`) | Enforced (tested) | `tkt-`: empty state vector + `derives-from` resolving to a `ctr-` artifact (#19) + projection header mandatory (position free, #18). `## Work Items` table on `ctr-`: GFM format; mismatch with owner state = **warning** (owner state = source of truth); unparseable table / missing state / unresolvable row = warning (#20–21). |
| **R9** | REQ-010 | §10, §14.2.6 | `conformance/rules.go:rule9` (+ `requiredSectionsFor`, `headingMatches`, `hasReplacement`) | `TestRule9SupersededADRWithReplacement`, `TestRule9VersionedReplacementMustNameInstance` (+ fixture `invalid-sections`, `invalid-adr-superseded`) | Enforced (tested) | Required sections per type family (validation.md rule 9). **`fnd-` requires 4 sections** (Purpose, Content, Investigation Summary, Conclusion — #22). Level-2 heading matching: `## Name` exact or `## Name <suffix>` (#23). ADR `superseded` must be referenced by ≥1 other artifact via `supersedes` resolving to its identity line; versioned references must point to the exact instance (#24). |
| **R10** | REQ-017 | §8.1, §8.3 | `conformance/rules_domain.go:rule10` (+ `conformance/rules_domain.go:domainNames`; `conformance/domain.go:DomainForToken`/`Stratum`/`StrataAbove`) | `TestRule10Stratification` (+ infra: `TestStratumOrdering`; fixtures `domain-valid`, `domain-invalid`) | Enforced (tested) | Stratification traceability (**warning**, never a commit blocker): every artifact below stratum 1 (Discovery) must have a resolvable `derives-from`/`depends-on` chain (direct or transitive) reaching a strictly higher stratum; deterministic BFS over resolvable edges, cycles/self-references harmless (Core v1.1 §8.3). Exempt: `tkt-`/`ses-` tokens and draft knowledge artifacts (work-item tokens own no `content-state` and are never exempt via the draft clause). **v1.1 addition** (REQ-017) — arrived with the EKA v1.1 taxonomy evolution (Core §14.2.3), see migration report. |
| **R11** | REQ-018 | §8.1, §8.3; Exchange v1.1 §14.2 R6 | `conformance/rules_domain.go:rule11` (+ `conformance/domain.go:DomainForToken`/`IsDomain`/`DomainNames`; `conformance/artifact.go:analyzeFile` — `Artifact.Domain` parse) | `TestRule11DomainCoherence`, `TestRule11NonStringDomainIsStructural` (+ fixtures `domain-valid`, `domain-invalid`) | Enforced (tested) | Domain coherence (**blocking**): the optional `domain` frontmatter field, when present, must be one of the five canonical Engineering Domains AND equal the artifact's home domain (`DomainForToken`); absent = OK (domain is derived — never Identity, never part of the State Vector, P15). Non-string `domain` value = R0 structural error, not R11. Declared-domain contract per Exchange v1.1 §14.2 R6. **v1.1 addition** (REQ-018). |
| **R12** | REQ-019 | §8.1 (Stratum Authority Invariant), §8.3 | `conformance/rules_domain.go:rule12` (+ `conformance/rules_domain.go:artifactIdentityForm`; `conformance/domain.go:DomainForToken`/`Stratum`) | `TestRule12SupersessionProhibition` (+ fixtures `domain-valid`, `domain-invalid`) | Enforced (tested) | Cross-stratum supersession prohibition (**blocking**): `supersedes`/`amends` may never target an artifact in a strictly higher stratum (smaller stratum number) — durable content moves down the authority chain, never up (Stratum Authority Invariant). Same-stratum and lower-stratum targets pass; unresolvable targets are left to R5 (R12 evaluates resolvable targets only). **v1.1 addition** (REQ-019). |

---

## 3. Requirement index

### 3(a) Enforced requirements (enforced by rules)

| Req ID | Requirement | Spec anchor | Enforced by | Status | Notes |
|---|---|---|---|---|---|
| REQ-001 | **Artifact identity rule** — every artifact carries `type` + `id`; frontmatter with only one of them, unknown type token, or invalid identity fields (namespace, non-integer `instance-version`, invalid date) = malformed (structural) | §3, §5.1 | R0 | Enforced (tested) | Structural failure goes to the R0 bucket; numbered rules not run for that artifact (interpretation #25) |
| REQ-002 | **Identity uniqueness** — no duplicate `(namespace, type, id, instance-version)` across the repository; ID unique within `(namespace, type)`, InstanceVersion unique within a line | §6.2.2, P3 | R1 | Enforced (tested) | Checked across the whole repository (rule over the full artifact set) |
| REQ-003 | **Filename as projection of identity** — filename token == `type`; `-v<nn>` suffix == `instance-version`; suffix only on `scp-`/`plan-` and mandatory (including v1) | §6.4, P9 | R2 | Enforced (tested) | Filename id part vs frontmatter `id` not checked (documented gap, interpretation #17) |
| REQ-004 | **State value validity** — every state field value ∈ its domain value set (including content-state variants living/ADR/decision); `phase` ∈ phase value set and only on `scp-`/`plan-` | §7.2 | R3 | Enforced (tested) | Domain value sets from the §7.2 table + per-type-family variants |
| REQ-005 | **Owned state vector** — state fields on the file == its type's owned set; `tkt-` carries an empty state vector; no state fields owned by other types | §7.4, §10 | R4 | Enforced (tested) | Owned sets from §10 (type→state binding is part of the standard) |
| REQ-006 | **Referential integrity** — all references (`amends`, `supersedes`, `derives-from`, `depends-on`, `validates`) resolve to an existing artifact; malformed = error; unresolved on draft = warning, otherwise error; self-reference = error | §6.2.7, §13.2.3 | R5 | Enforced (tested) | Bare-id accepted per interpretation #9; bidirectional convention not enforced (#11) |
| REQ-007 | **Classification-location consistency** — knowledge artifact: `dimension` == home folder (nearest ancestor ∈ 12 dimensions); `ctr-`/`tkt-`/`ses-` must not carry `dimension` | §8, P15 | R6 | Enforced (tested) | Classification is a property, not identity (P15) |
| REQ-008 | **Change-log consistency** — every owned domain has an entry; last entry == current value; transitions legal; initial `from: "-"` | §5.2, P7 | R7 | Enforced (tested) | Chain contiguity not enforced (interpretation #4) |
| REQ-009 | **Single-writer & projection discipline** — `tkt-` empty state vector + `derives-from` to ctr- + projection header; container Work Items tables validated against owner state | §7.4, P6 | R8 | Enforced (tested) | Table vs owner state mismatch = warning (owner = source of truth) |
| REQ-010 | **Well-formed content** — required sections per type family present; superseded ADR must be referenced by its replacement | §10, §14.2.6 | R9 | Enforced (tested) | `fnd-` requires 4 sections (interpretation #22); level-2 headings (#23) |
| REQ-017 | **Stratification traceability** — every artifact whose Engineering Domain is not Discovery (stratum 1) has a resolvable `derives-from`/`depends-on` chain (direct or transitive) reaching a strictly higher stratum; exempt: `tkt-`/`ses-` tokens and draft knowledge artifacts | §8.1, §8.3 | R10 | Enforced (tested) | Warning severity (never blocks a commit); **v1.1 addition** (new REQ entry for R10, taxonomy evolution Core §14.2.3) |
| REQ-018 | **Domain coherence** — the optional `domain` frontmatter field, when present, is one of the five canonical Engineering Domains and equals the artifact's home domain; absent = derived, no check | §8.1, §8.3; Exchange v1.1 §14.2 R6 | R11 | Enforced (tested) | Blocking severity; classification property, never Identity (P15); **v1.1 addition** (new REQ entry for R11) |
| REQ-019 | **Cross-stratum supersession prohibition** — `supersedes`/`amends` never target an artifact in a strictly higher stratum; durable content moves down the authority chain | §8.1, §8.3 | R12 | Enforced (tested) | Blocking severity; unresolvable targets left to R5; **v1.1 addition** (new REQ entry for R12) |

### 3(b) Governance-only requirements (normative in the spec, NOT enforced by any rule)

| Req ID | Requirement | Spec anchor | Enforced by | Status | Notes |
|---|---|---|---|---|---|
| REQ-011 | **Concurrency control** — exactly one active Execution Container (exactly-one-active); the next creation waits | §9 (Concurrency Control), §7.5, §5.2 | — | Governance-only (uncovered) | Operating Layer invariant; mechanical validation requires cross-artifact + temporal observation. Gap #29 in `conformance-notes.md` — future rule candidate |
| REQ-012 | **Plan immutability / lock-atomic-with-generation** — Execution Container creation locks the plan atomically; locked plan Content unchanged; post-lock changes = new instance | §5.2, §9 (Versioning/Immutability) | — | Governance-only (uncovered) | Transition unobservable from a single static snapshot; requires cross-instance history |
| REQ-013 | **Approved-content immutability** — Content that passed the approval gate is not silently mutated; changes only via the governance channel (amend/supersede) | P8, §5.4 (invariant 6) | — | Governance-only (uncovered) | Requires non-mutation evidence over time; a snapshot cannot prove it |
| REQ-014 | **Phase change via readiness gate** — phase changes authorized only by the readiness Gate: release-ready = (all work items Done) ∧ (container Completed) ∧ (plan locked) ∧ (review gate) ∧ (approval gate) | §11.2, §7.5 | — | Governance-only (uncovered) | Aggregate State evaluation across artifacts; phase is a context attribute, not a state domain |
| REQ-015 | **Distillation before archive** — Session/Review must be distilled into durable artifacts (ADR/decision) before Archived | §11.4 | — | Governance-only (uncovered) | Lifecycle rule; requires cross-type relationship semantics |
| REQ-016 | **Identity storage-independence** — Identity is not encoded in location/storage; every Identity resolves to exactly one artifact | §6.4, §12.2 | R2, R6 (partially) | Partially enforced | R2 (filename as projection) and R6 (`dimension` == folder) enforce the serialization surface of this invariant; storage-independence itself cannot be checked mechanically within a single repository |

---

## 4. Test coverage index

Every test function is **classified exactly once** in this index; integration functions (such as `TestInvalidFixtures`) may be re-quoted as cross-references on the rule rows they serve. Total: **67 tests** (59 `conformance` + 8 `cmd/eka`).

### 4(a) Rule tests (unit, per rule)

| Rule | Test functions |
|---|---|
| R0 | `TestAnalyzeNoFrontmatterIsConventionDoc`, `TestAnalyzeFrontmatterWithoutTypeIDIsConventionDoc`, `TestAnalyzeUnterminatedFrontmatter`, `TestAnalyzeBrokenYAML`, `TestAnalyzeTypeXorID`, `TestAnalyzeValidArtifact`, `TestAnalyzeMissingIdentityFields`, `TestAnalyzeNonIntVersion`, `TestAnalyzeInvalidDate`, `TestAnalyzeUnknownType`, `TestAnalyzeChangeLogNotList`, `TestAnalyzeMalformedChangeLogEntry`, `TestUnknownTypeIsStructural` (13) |
| R1 | — (no dedicated unit test; covered via `TestInvalidFixtures` → `invalid-dup-identity`) |
| R2 | `TestRule2ExactCounts` (1) |
| R3 | `TestPhaseValidation` (1) |
| R4 | — (no dedicated unit test; covered via `TestInvalidFixtures` → `invalid-ownership`) |
| R5 | `TestRule5DraftSeverity`, `TestRule5CrossNamespaceAndVersionResolution`, `TestRule5VersionedReferenceToMissingInstance` (3) |
| R6 | `TestDimensionFolderResolution` (1) |
| R7 | `TestRule7ExactCounts` (1) |
| R8 | `TestRule8TicketHeaderAndDerivation`, `TestRule8ContainerTableMismatchIsWarningOnly` (2) |
| R9 | `TestRule9SupersededADRWithReplacement`, `TestRule9VersionedReplacementMustNameInstance` (2) |
| R10 | `TestRule10Stratification` (1) |
| R11 | `TestRule11DomainCoherence`, `TestRule11NonStringDomainIsStructural` (2) |
| R12 | `TestRule12SupersessionProhibition` (1) |

### 4(b) Infrastructure tests (not bound to one specific rule)

Explicitly classified as `infrastructure` — protects scan prerequisites, the result model, grammar parsing, state tables, and the CLI.

| Group | Test functions |
|---|---|
| Result model & determinism (`report.go`/`validate.go`) | `TestReportCounts`, `TestReportPassSemantics`, `TestRelPathIsRootRelative`, `TestReportDeterminism`, `TestSortedResultsOrder` (5) |
| Scan policy (`validate.go`) | `TestScanSkipsTestdataAndHiddenDirs`, `TestValidateInputErrors`, `TestConventionDocumentsAreSkipped` (3) |
| Filename parsing (`filename.go`) | `TestParseFilename`, `TestParseFilenameEmpty` (2) |
| Reference grammar parsing (`rules.go:parseReference`) | `TestParseReference`, `TestParseReferenceCrossNamespace` (2) |
| State tables & taxonomy (`state.go`) | `TestTypeTokenCount`, `TestOwnedSets`, `TestContentStateVariant`, `TestDimensionTokens`, `TestIsLegalTransition`, `TestPhaseValueSet` (6) |
| Domain ontology (`domain.go`) | `TestDomainForTokenComplete`, `TestDomainForTokenValues`, `TestDomainForDimensionComplete`, `TestDomainForDimensionValues`, `TestStratumOrdering`, `TestDomainNamesSorted`, `TestIsDomain` (7) |
| Self-conformance | `TestReferenceImplementationConforms`, `TestFindRepoRoot` (2) |
| **CLI layer** (`cmd/eka/main_test.go`) | `TestExitCodeUsage`, `TestExitCodeBadPath`, `TestHelpExitsZero`, `TestValidateValidRepoExitsZero`, `TestValidateInvalidRepoExitsOne`, `TestWarningsDoNotAffectExitCode`, `TestDefaultPathIsCurrentDirectory`, `TestOutputIsDeterministic` (8) |

### 4(c) Fixture-based integration tests (`conformance/testdata/`)

| Test functions | Coverage |
|---|---|
| `TestValidFixtureRepo` | Valid repo (`valid/`, 6 artifacts) — all rules pass without errors |
| `TestInvalidFixtures` | 11 invalid scenario directories: `invalid-malformed` → R0; `invalid-dup-identity` → R1; `invalid-filename` → R2; `invalid-state-value` → R3; `invalid-ownership` → R4; `invalid-reference` → R5; `invalid-dimension` → R6; `invalid-changelog` → R7; `invalid-projection` → R8; `invalid-sections` + `invalid-adr-superseded` → R9 |
| `TestDomainValidFixtureRepo`, `TestDomainInvalidFixtureRepo` | Engineering Domain fixtures: `domain-valid/` (7 artifacts — all rules pass, 0 errors / 0 warnings); `domain-invalid/` → R11 (unknown domain + home-domain mismatch), R12 (upward `supersedes` + `amends`), R10 warnings on the isolated `ctr-`/`sto-` artifacts (draft spec + ticket exempt) |

---

## 5. Coverage analysis

### Covered

- **Rule coverage: 13/13** — R0 + R1–R12 all `Enforced (tested)`: every rule has a Go implementation called from `conformance/validate.go:Validate` and test coverage (unit, infrastructure, or fixture).
- **Test coverage: 67 tests** (59 `conformance` + 8 `cmd/eka`), all mapped in Section 4.
- **Self-conformance PASS** — `go run ./cmd/eka validate .` on this repository: 8 artifacts, 0 errors, 8 warnings (R10 on the 8 Architecture-stratum ADRs), exit 0 (codified as `TestReferenceImplementationConforms`).

### Uncovered specification sections

Spec sections covered only as normative requirements (REQ-011..REQ-016, Section 3(b)), with status `Governance-only (uncovered)` or `Partially enforced`:

| Req ID | Spec anchor | Status | Recommendation |
|---|---|---|---|
| REQ-011 (exactly-one-active) | §9, §7.5, §5.2 | Governance-only (uncovered) | Future rule candidate — **follow-up, not proposed now** |
| REQ-012 (lock-atomic-with-generation) | §5.2, §9 | Governance-only (uncovered) | Future rule candidate — **follow-up** |
| REQ-013 (approved-content immutability) | P8, §5.4 | Governance-only (uncovered) | Future rule candidate — **follow-up** |
| REQ-014 (phase readiness gate) | §11.2, §7.5 | Governance-only (uncovered) | Future rule candidate — **follow-up** |
| REQ-015 (distillation before archive) | §11.4 | Governance-only (uncovered) | Future rule candidate — **follow-up** |
| REQ-016 (identity storage-independence) | §6.4, §12.2 | Partially enforced (R2/R6) | Remaining surface not mechanically checkable within one repo — **follow-up** |

### Orphan implementations

**None.** Verified by reading the code: `conformance/validate.go:Validate` calls `analyzeFile` (R0) in the parse phase and `rule1`–`rule12` in the rule phase; all helper functions (Section 2, Implementation column) are referenced by the rule they serve. `conformance/report.go` and `cmd/eka/main.go` are engine/CLI infrastructure, not rules.

### Orphan tests

**None.** All 67 test functions are mapped in Section 4 — every function appears exactly once, no duplication, nothing unmapped.

### Known gaps (from `conformance-notes.md`)

| Gap | Detail | Reference |
|---|---|---|
| R2 — filename id vs frontmatter id | The id part of the filename is not matched against the frontmatter `id` | [conformance-notes.md](conformance-notes.md) — interpretation #17, Gap table |
| Exactly-one-active container | The "exactly one active Execution Container" invariant is not validated | [conformance-notes.md](conformance-notes.md) — interpretation #29, Gap table |
| R5 — bidirectional references | The "references only on the referring artifact" convention is not enforced | [conformance-notes.md](conformance-notes.md) — interpretation #11, Gap table |

All three gaps are reflected in the Notes column of Section 2 (R2, R5) and the REQ-011 status in Section 3(b). Gaps are not closed in this version — changing validator behavior is a scope change, outside v1.0 governance.

---

## 6. Summary

| Metric | Value |
|---|---|
| Total Requirements | **19** (13 enforced + 6 governance-only; REQ-017–REQ-019 added as v1.1 entries) |
| Total Conformance Rules | **13** (R0 + R1–R12; R10–R12 per Core v1.1 §8.1) |
| Total Implementations | **13 rule implementations** (`analyzeFile` + `rule1`–`rule12`) + **28 helper functions** (`parseFilename`, `identityKey`, `buildIndex`, `contentStateVariant`, `domainValues`, `isLegalTransition`, `parseReference`, `resolve`, `dimensionFolderFor`, `dimensionList`, `entriesForDomain`, `indexOfEntry`, `workItemsTable`, `compareWorkItemsTable`, `resolveWorkItemCell`, `hasProjectionHeader`, `splitTableRow`, `requiredSectionsFor`, `headingMatches`, `hasReplacement`, `domainNames`, `artifactIdentityForm`, `DomainForToken`, `DomainForDimension`, `Stratum`, `StrataAbove`, `DomainNames`, `IsDomain`) + engine/report/CLI (`validate.go`, `report.go`, `cmd/eka/main.go`) |
| Total Test Suites | **2 packages**: `conformance` 59 tests + `cmd/eka` 8 tests = **67 tests** |
| Current coverage | **13/13 rules enforced & tested (100% rule coverage)**; requirement coverage = **13 enforced of 19 total**; spec-section coverage — enforced requirements map to §3, §5, §6, §7, §8 (incl. §8.1 Engineering Domains, §8.3 Stratification Governance), §10, §13, §14 (+ Exchange v1.1 §14.2 R6); governance-only requirements map to §5, §7, §9, §11, §12 |
| Self-conformance | `eka validate .` = 8 artifacts, 0 errors, 8 warnings (R10 on the 8 Architecture-stratum ADRs), exit 0 |
| Identified gaps | REQ-011..REQ-016 not enforced (6 governance-only requirements) + 3 documented gaps (filename-id, exactly-one-active, bidirectional references) |
| Recommended follow-ups | (1) Future rule candidates for REQ-011..REQ-016 — **not proposed now**; (2) automation of matrix consumption (deterministic parser over the table structure) — **not now**, the matrix is a living markdown document |
