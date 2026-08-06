# Migration Report — Engineering Domain Ontology v1.0 → v1.1

> Anchor EKA: Knowledge Taxonomy — Engineering Domains and Knowledge Stratification (Core v1.1 §8.1); taxonomy extension governance (Core §14.2.3). Convention document, not an artifact.
> Standard: EKA v1.1, dated 2026-08-05.

## 1. Purpose

This report records the evolution of the Engineering Domain ontology from EKA v1.0 to v1.1 and classifies every change for existing repositories. It is the required migration classification deliverable of the v1.1 taxonomy extension (Wave 1), produced under the **taxonomy-extension governance path** of Core §14.2.3: **core closed, taxonomy open**. Identity, layer contracts, and invariants are closed to extension; taxonomies (types, dimensions, domains, protocol) are open under governance — this evolution changes taxonomy only.

Target reading: repository owners and tooling maintainers holding a conformant v1.0 repository or v1.0 exchange packages. The verdict up front: **zero breaking changes; all changes are compatible or editorial.**

## 2. What changed

| # | Change | Where | Nature |
|---|---|---|---|
| 1 | **Engineering Domains + Knowledge Stratification** — five canonical Engineering Domains (Discovery → Architecture → Planning → Execution → Operations, stratum 1 highest → 5), the derived Knowledge Stratum, the **Stratum Authority Invariant**, and the Representation Alias table (methodology terms → canonical token + domain). | Core v1.1 §8.1 | Additive taxonomy |
| 2 | **Three new terms** registered — Engineering Domain, Knowledge Stratum, Representation Alias; no redefinition of any existing term. Conformance Rules set formalized as **R0–R12** (R1–R9 = Exchange §14.2; R0 = structural; R10–R12 = domain-aware, Core §8.1). | Naming v1.1 §2.3, §9.1, §9.3, §12.1 | Additive terminology |
| 3 | **Classification + Engineering Domain** — the Engineering Domain member enters unit classification (derived from the Type's home domain; optional explicit declaration must match); domain-aware rules R10–R12 complement R1–R9. | Exchange v1.1 §4.4, §14.2 (R6), §16 | Additive |
| 4 | **Serialization Version 1.1** — RSF adds the Engineering Domain member to unit Classification and the optional Representation metadata (Representation Alias (RSF)); legacy Serialization Version 1 remains importable with domain derivation at import. | RSF v1.1 §5.1, §5.4, §11.2 | Additive |
| 5 | **Validator R10–R12** — `conformance/domain.go` (single source of truth for token/dimension → domain mapping) + `conformance/rules_domain.go` (R10 warning, R11/R12 blocking). | `conformance/` | Additive implementation |
| 6 | **View header `Domain` row** — projection context headers carry an additive `Domain: <Engineering Domain>` row alongside `Knowledge`. | `cli.md` (view commands) | Additive presentation |
| 7 | **`eka version` → EKA standard 1.1** — the CLI reports the ratified standard corpus version (Core + Exchange + Naming v1.1). | `cmd/version.go` | Compatible |
| 8 | **Skeleton docs + lifecycle guide** — `skeleton/docs/lifecycle.md` (Engineering Domain as orientation, stratum as authority); zone README updates. | `skeleton/` | Editorial |
| 9 | **Terminology reframe** — sprint → Execution container alias, PRD → `req-` alias, ticket → `tkt-` alias; methodology terms are Representation Aliases, never frontmatter values. | Core §8.1, Naming §2.3 | Editorial |

## 3. Migration classification

### 3.1 BREAKING — target zero

There are **no breaking changes** in v1.1. The none-items are explicit:

| Would-be breaking area | Status | Why it does not break |
|---|---|---|
| Conformance Rules R0–R9 | **Unchanged** | The nine rules of Exchange §14.2 and the structural rule R0 keep their exact v1.0 semantics and verdicts; R10–R12 only complement them (Naming §9.3). |
| Identity Model | **Unchanged** | Identity tuple `(Namespace, Type, ID, InstanceVersion)` and the Identity invariants (Core 6, P3) are untouched; Engineering Domain is a classification property, never Identity (P15). |
| State Taxonomy | **Unchanged** | The five State Domains, their value sets, forward-only transitions, and the change-log contract (Core 7, P7) are untouched. |
| Layer Model + layer contracts | **Unchanged** | KB + OS + EX composition and the two-channel principle (P10) are untouched; the Exchange Contract v1 is unchanged (Exchange §1.1). |
| Global invariants | **Unchanged** | The seven global invariants (Core §5.4) are untouched. The Stratum Authority Invariant is a **taxonomy invariant** added alongside them — it can only be strengthened, never weakened (Core §8.1, §16.3). |
| Type tokens + Knowledge Dimensions | **Unchanged** | The 26 type tokens and 12 Knowledge Dimensions keep their meaning and folder mapping (R6 `dimension == folder` unchanged). |
| Exchange format version | **Unchanged** | Stays **1**. A v1.1 package is an exchange-format-v1 package with an additive Engineering Domain member (Exchange §9.3, §16.2). |
| CLI commands + exit codes | **Unchanged** | Command surface and the exit code contract (0 success — warnings allowed / 1 blocking violations / 2 usage) are unchanged. |
| Existing packages | **Still importable** | A v1.0 package carries no Engineering Domain member; the importer derives the domain from the Type token at import (RSF §11.2, Exchange §9.2.4). |

### 3.2 COMPATIBLE

| Change | Semantics | Compatibility guarantee |
|---|---|---|
| Optional `domain` frontmatter | R11: when present, `domain` must be one of the five canonical Engineering Domains **and** equal the artifact's home domain (from its token); **absent = OK** — the domain is derived. Unknown domain or mismatch = blocking. | A conformant v1.0 artifact carries no `domain` field; R11 is silent on it. Zero violations on conformant v1.0 repositories. |
| R11 + R12 blocking rules | R11 domain coherence (above); R12: `supersedes`/`amends` may never target an artifact in a **strictly higher stratum** (resolvable targets only; malformed/dangling targets stay with R5). | The rules forbid what conformant v1.0 content never does: declared-domain mismatch and upward supersession/amendment. Zero violations on conformant v1.0 repositories. |
| R10 warning | Every artifact below stratum 1 must have a resolvable `derives-from`/`depends-on` chain (direct or transitive) reaching a strictly higher stratum. Exempt: `tkt-`/`ses-` tokens; `content-state: draft` (knowledge artifacts only — work-item tokens own no content-state and are never exempt via this clause). | **Warning severity, never blocking**; warnings never affect the verdict or the exit code. Existing strata-isolated artifacts (e.g., ADRs with no chain to Discovery) produce warnings only. |
| Serialization Version 1 accepted at import | RSF v1.1 importer accepts Serialization Versions **1 and 1.1**; a v1.0 package derives the Engineering Domain at import and carries no Representation metadata. | Legacy packages import unchanged; re-export upgrades them to 1.1 with domain metadata. |
| Package label suffix `-1` / `-1.1` both accepted | The Package Identity Label carries the Serialization Version; importers accept both suffix values for the respective versions (RSF §4.1, §11.2). | No package relabeling required. |
| View header `Domain` row | Projection context headers gain an additive `Domain: <domain>` row (always `Execution` for sprint/wave/ticket projections) | Additive presentation row; consumers ignoring unknown rows are unaffected. |

### 3.3 EDITORIAL

| Change | Detail |
|---|---|
| Terminology reframe | Methodology terms become Representation Aliases mapped to canonical tokens + domains: sprint/iteration → `ctr-`, PRD → `req-`, ADR/RFC → `adr-`, epic → `epc-`, ticket → `tkt-`, release → `rel-`, incident → `bug-`, runbook → `run-` (Core §8.1 alias table; Naming §2.3). No frontmatter values change. |
| Docs rewrites | `skeleton/docs/lifecycle.md` (Engineering Domain as orientation, stratum as authority), zone READMEs, `reference-architecture.md` §2.8 (Engineering Domain), `terminology-glossary.md`. |
| `cli.md` example output | Self-validation example output now documents the repository's own **7 R10 warnings** (the 7 Implementation ADRs — Architecture, stratum 2, no chain to Discovery) with the note that R10 never blocks; verdict stays PASS. |
| Spec file renames v1.0 → v1.1 | `standard/eka-specification-v1.1.md`, `standard/eka-naming-and-terminology-specification-v1.1.md`, `standard/eka-exchange-specification-v1.1.md`, `reference/eka-reference-serialization-format-v1.1.md` with link updates across READMEs, traceability matrix, and ratifications notes. v1.0 content is unchanged inside the v1.1 files (additive revisions). |
| Rule-set phrasing | "nine rules" now always qualified: R1–R9 (Exchange §14.2) vs the full set R0–R12 (Naming §9.3). |

## 4. Migration strategy

Ordered steps for an existing v1.0 repository. Steps 1–2 are mandatory; 3–4 are optional improvements; 5 is a standing rule; 6 applies to package consumers.

1. **Upgrade the validator binary** — install the v1.1 `eka` (`go install github.com/maleolabs/engineering-knowledge-architecture/cmd/eka@latest`). The new binary enforces R0–R12; the old binary cannot see the domain-aware rules.
2. **Run `eka validate .`** — expect **0 errors and at most warnings** (R10 stratification traceability on strata-isolated artifacts). Exit code stays `0`; warnings never block. Any R11/R12 error indicates pre-existing non-conformance (declared-domain mismatch or upward supersession), which v1.0 validation should already have caught.
3. **Optionally add `domain` frontmatter** to artifacts. If declared, it must equal the token's home domain (R11, blocking) — e.g., `domain: Architecture` on `adr-…`. Absent is always correct; the domain is derived.
4. **Address R10 warnings where the knowledge warrants it** — add `derives-from`/`depends-on` chains to artifacts in a strictly higher stratum (e.g., an ADR deriving from a requirement or strategy). While an artifact is in `content-state: draft` the warning is suppressed (draft tolerance); `tkt-`/`ses-` are exempt outright. R10 is a quality signal: a warning may be accepted and documented instead of resolved.
5. **R12 standing rule: never supersede/amend upward** — `supersedes`/`amends` may never target a strictly higher stratum (blocking). Carry such changes **out in the higher stratum** (new instance + Relationship there, forward-only), per the Stratum Authority Invariant's governance channel (Core §8.1).
6. **Existing packages stay importable** — Serialization Version 1 packages import without change (domain derived at import). To upgrade: **re-export** the repository; the v1.1 exporter writes the derived Engineering Domain into unit Classification and optionally the Representation metadata; the new package declares Serialization Version 1.1.

## 5. Versioning rationale

v1.1, not v2.0, because:

- **Taxonomy extension, not contract change** — the evolution is exactly the case Core §14.2.3 opens: "Core closed, taxonomy open." Engineering Domains are taxonomy; Identity, layer contracts, and invariants are untouched.
- **No canonical term redefined** — Naming v1.1 registers three new terms (Engineering Domain, Knowledge Stratum, Representation Alias) under the §12.1 terminology governance; no existing term changes meaning (Naming §5.3 minor-version rule).
- **No invariant weakened** — the seven global invariants stand; the Stratum Authority Invariant is an **addition**, a taxonomy invariant that can only be strengthened, never weakened (Core §8.1, §16.3).
- **Exchange contract additive** — the Engineering Domain member is optional in unit classification: an absent member is derived from the Type token, so v1.0 packages remain conformant (Exchange §9.3, §16.2; R6).
- **RSF additive** — Serialization Version 1.1 accepts 1 and 1.1; a v1.0 package is not a forward-compatibility violation (RSF §11.2).
- A v2.0 bump would signal contract breakage, redefined terms, or weakened invariants — none of which occurred.

## 6. Governance follow-ups

| Follow-up | Status | Path |
|---|---|---|
| R10 escalation to blocking | **Future minor** | R10 is deliberately a warning today (stratification is a structural quality signal, never a commit blocker). Escalation to blocking severity, if ever, is a future minor change to the conformance suite under Core §14.2 governance — it would make stratification chains mandatory and must be announced as a minor revision. |
| Conformance-traceability-matrix extension for R10–R12 | **Pending governance work** | The current matrix covers the v1.0 rule surface R0–R9 with a note pointing at `conformance/rules_domain.go` coverage; the formal R10–R12 rows are the outstanding governance change that introduced the domain-aware rules. |
| Protocol-level stratum gates | **Future protocol variant** | Stratum awareness currently lives in the validator and the exchange classification. Protocol-level stratum gates (e.g., stratum-aware transfer policies) belong to a future registered protocol variant (Exchange §18.3) — never a silent core change. |
| Representation alias registry maintenance | **Standing** | The alias table (Core §8.1) and the RSF Representation Alias metadata (§5.4) must stay in sync; new methodology terms enter only via terminology review (Naming §12.1) — aliases are never frontmatter values and never Artifact types. |

## 7. Verification

- **362 tests green** — `go test ./...`: the full test surface (conformance suite incl. domain-aware rules, bootstrap, view, exchange, cmd) passes.
- **Self-validation 0 errors / 7 R10 warnings** — `eka validate .` on this repository: 7 artifacts, 0 errors, 7 R10 warnings (the 7 Implementation ADRs — Architecture stratum, no upward chain; accepted, documented, non-blocking), `Verdict: PASS`, exit 0. Codified as `TestReferenceImplementationConforms`. With ADR-008 (concurrent, see `reference/decisions/adr-008-engineering-domain-model.md`) the counts become 8 artifacts / 8 R10 warnings — same verdict.
- **Export/import round-trip byte-identical** — re-import of an identical package is a no-op; re-export after import yields the identical package up to the permissible differences of Exchange §15.4 (RSF §10); Engineering Domain derivation is deterministic, so equivalent repositories produce identical Classification payloads (Exchange §4.4).

---

*End of Migration Report — Engineering Domain Ontology v1.0 → v1.1.*
