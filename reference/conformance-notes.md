# Conformance Implementation Notes — EKA Validator (`conformance/`)

> Convention document, not an artifact. Meta-documentation of the `reference/` zone.
> Related: [`cli.md`](cli.md) (CLI documentation), [`../skeleton/docs/exchange/validation.md`](../skeleton/docs/exchange/validation.md) (Conformance Rules R0–R12), [`../standard/eka-specification-v1.1.md`](../standard/eka-specification-v1.1.md) (canonical standard).

> **Traceability consolidation.** The rule traceability tables (specification ↔ implementation) of this document have been consolidated into the [`conformance-traceability-matrix.md`](conformance-traceability-matrix.md) — the single source of truth for conformance coverage (REQ→Spec→Rule→Impl→Test). This document now holds only **interpretation decisions (29 items)** and **known gaps**; update that matrix, not the tables here, when coverage changes.

> **Domain-aware rules (EKA v1.1).** The validator now executes the full rule set **R0–R12**. The interpretation decisions below cover R0–R9; the domain-aware rules R10–R12 (stratification traceability — warning; domain coherence — blocking; cross-stratum supersession prohibition — blocking; Core v1.1 §8.1) are implemented in `conformance/rules_domain.go` with their semantics documented in `skeleton/docs/exchange/validation.md` — no separate interpretation decisions were needed.

## Purpose

This document explains **how the `eka` CLI executes the EKA specification mechanically**: artifact vs convention document classification, rules R0–R12, and — most importantly — every **interpretation decision** taken when rule text is not precise enough for machine execution. The goal is rule-by-rule traceability: from standard text → `validation.md` rule text → validator behavior → Go implementation location.

**Interpretation policy:** if the specification is ambiguous, the decision is **documented before implementation**; no behavior is invented without a basis. Each decision below cites its specification basis (rule text, ADR, dimension README, or the repository's own reality). The same decision is also recorded as an `Interpretation (documented)` comment in the relevant source code.

## Interpretation decisions (29 items)

Numbers #1–#28 are decisions taken during implementation; #29 is a known gap deliberately **not** validated (full documentation in the [Known gaps](#known-gaps) section).

| # | Rule | Ambiguity | Decision | Rationale |
|---|---|---|---|---|
| 1 | Artifact rule | Unquoted frontmatter `date: 2026-08-05` parses as a timestamp by yaml.v3, not a string | Date fields accept a string **or** timestamp node; normalized to `YYYY-MM-DD` | The repository's canonical ADRs write dates unquoted; rejecting them would break the repository's own artifacts |
| 2 | R4 | "exactly the same (absence = N/A)" — may an **owned** field be **absent**? | Absent owned field = **ERROR** (`missing owned state field`) | "exactly the same" + ADR-002 "fields present only for owned domains" + all real artifacts carry all their owned fields; R7 also needs the current value to compare |
| 3 | R4 | Absent non-owned fields — how are they evaluated? | Absent **non-owned** field = N/A, not evaluated (not an error) | Explicit "absence = N/A" semantics in the rule text |
| 4 | R7 | "No transition without a change-log entry" (bullet 2) cannot be verified from a snapshot | Enforced: entry presence per owned domain, `last-entry == current value`, entry well-formedness, transition legality. **Not** enforced: chain contiguity `entry[i].from == entry[i-1].to` | A static snapshot does not observe intermediate state; canonical ADRs start content-state history mid-chain (e.g., `proposed -> accepted` without an initial entry) |
| 5 | R7 | Meaning of `from: "-"` and `to: "-"` not defined in rule text | `from: "-"` = initial state marker, legal for all domains; `to: "-"` = **invalid** | Convention used by all repository ADRs; `to` must be a real domain value |
| 6 | R7 | How strict is "forward-only"? | Execution State: **strictly adjacent** (no skipping). Other domains: forward-only (index increasing, adjacency not required) | EKA §7.2 "strictly sequential" applies only to Execution State; other domains merely must not regress (P7) |
| 7 | R7 | Change-log entries for `phase` — ordering rules? | Phase entries: value-set + well-formedness validation **only**; no ordering constraint | Phase is a context attribute (EKA 11.2), not a state domain |
| 8 | R7 | Artifact with no change-log at all — error? | Error only if the artifact has ≥1 owned domain or phase; `tkt-` (empty state vector) may omit the change-log | No state whose transitions need recording |
| 9 | R5 | Bare-id reference (`001-identity-serialization`) does not match the `<type>:<id>` grammar | Bare-id accepted = line reference resolved within the referrer's namespace + type | 6 of 7 real ADRs use this de-facto form; rejecting it would break the repository itself |
| 10 | R5 | Malformed vs unresolved references — severity? | Malformed (unparseable) = **always error**; unresolved = warning when `content-state: draft`, error otherwise (including artifacts without content-state) | Rule text only covers unresolved; malformed is a structural violation that can never resolve |
| 11 | R5 | "References are only written on the referring artifact (not bidirectional)" — enforced? | **Not** mechanically enforced | Not a clear mechanical check (requires per-relationship direction semantics); the implementation task forbids adding rules |
| 12 | R5 | Cross-namespace reference format | `<ns>/<type>:<id>`; optional version `<ns>/<type>:<id>:<ver>` | Grammar extension of the rule text, consistent with `<type>:<id>[:<instance-version>]` |
| 13 | R6 | "Home folder" for files outside `docs/` (e.g., `reference/decisions/`) | Home folder = nearest ancestor whose name ∈ 12 dimensions; not found → error | Makes `reference/decisions/` and `docs/decisions/` both map to `decisions`; still works for artifacts nested under dimension folders |
| 14 | R6 | Knowledge artifact without a `dimension` field | = **ERROR** | Classification is a mandatory artifact property (P15; reference-architecture.md §2.5) |
| 15 | R6 | `dimensions-secondary` — how far is it validated? | Token validation only (must ∈ 12 dimensions); forbidden on `ctr-`/`tkt-`/`ses-` | Secondary classification property; no deeper rule text |
| 16 | R2 | Digit count of the `-v<nn>` suffix | 1+ digits accepted (`-v1` and `-v01` both valid) | "including v1" (docs/README.md, validation.md rule 2) |
| 17 | R2 | Filename id part vs frontmatter id | **Not** matched (documented gap) | The rule only requires token + version consistency; the filename is a projection (ADR-001), not Identity |
| 18 | R8 | Position of the projection header in the file | Anywhere in the file (exact match, trailing whitespace ignored); position free | Rule text only requires the header to be "present in the projection file" |
| 19 | R8 | What must `tkt-` point to via `derives-from`? | ≥1 `derives-from` reference resolving to a `ctr-` artifact; references to work items not required | ADR-003 §3; the projection README shows both, rule text mentions containers only |
| 20 | R8 | `## Work Items` table format not exemplified in containers/README.md | Defined: GFM table (header row + separator row); first column = work item id or `<type>:<id>`; execution-state column recognized by header variants (execution-state / execution state / execution_state / status); unparseable table / missing state column / unresolvable row / invalid projection value = **WARNING** | Rule text gives no example; table comparison results are warning-oriented (owner state = source of truth) |
| 21 | R8 | Table row without a state cell / short row | Rows without a state cell silently skipped; short rows (cells < header) = warning | Rows without state cannot be compared; short rows indicate a broken table |
| 22 | R9 | `fnd-` appears in two rows of the validation.md table (Knowledge doc vs Research Finding) | `fnd-` requires **4 sections**: Purpose, Content, Investigation Summary, Conclusion | The specific row (Research Finding) wins over the general row (Knowledge doc); research/README.md confirms the 4-section structure |
| 23 | R9 | Matching of required section headings | `## Name` exact or `## Name <something>`; `###` not counted | Level-2 headings are the content structure convention; suffixed variants allowed |
| 24 | R9 | Superseded ADR — who replaces it? | ≥1 other artifact whose `supersedes` resolves to that ADR's identity line; versioned references must point to the exact instance | ADRs are replaced per instance (identity line), not per line |
| 25 | R0 | Unknown type token | = structural error (R0); rules R2, R3, R4, R6, R7 skipped for that artifact; R5 still checks references (unknown tokens reported); R8/R9 not applicable | Numbered rules are meaningless without a known token; the R5 exception is verified empirically |
| 26 | R0 | `instance-version: "1"` (quoted) | = error (not an integer) | The specification defines the field as an integer; canonical ADRs write it unquoted |
| 27 | Scan | `testdata/` directories and dot-prefixed (`.git`) | **Not descended into** | Go test fixtures are not knowledge base content; without this, fixtures would break self-validation |
| 28 | Scan | Non-.md files, symlinks, unreadable files | Non-.md ignored; symlinks not followed; unreadable `.md` → `Validate` error | A scan that cannot see all files cannot assert conformance |
| 29 | Gap | Exactly-one-active container (protocol.md §3) not validated | **Not** validated; recorded as a gap / future rule candidate | Not part of R1–R9; the implementation task forbids adding rules |

## Rule traceability matrix (specification ↔ implementation)

> **HISTORICAL — consolidated.** This traceability table has been moved to
> [`conformance-traceability-matrix.md`](conformance-traceability-matrix.md) as the
> **single source of truth** for conformance coverage (Engineering Requirement → Specification →
> Conformance Rule → Implementation → Automated Test). This document now holds only
> interpretation decisions (#1–#29) and known gaps. Conformance changes must
> update that matrix (see `CONTRIBUTING.md`).

## Known gaps

The following gaps are **deliberately not validated**; recorded so the decision is not lost and they can become future rule candidates:

| Gap | Detail | Rationale |
|---|---|---|
| R2 — filename id vs frontmatter id | The id part of the filename is not matched against the frontmatter `id` | The rule only requires token + version consistency; the filename is a projection (ADR-001), the true Identity is in frontmatter |
| Exactly-one-active container | The "exactly one active Execution Container" invariant (protocol.md §3) is not checked | Not part of R1–R9; the implementation task forbids adding rules — future rule candidate |
| R5 — bidirectional references | The "references are only written on the referring artifact (not bidirectional)" convention is not enforced | Not a clear mechanical check; requires per-relationship direction semantics |

## R10 Governance Roadmap

**Current state (v1.1).** R10 — stratification traceability — is a **warning** (Core v1.1 §8.3; `conformance/rules_domain.go`): every artifact below stratum 1 must have a resolvable `derives-from`/`depends-on` chain (direct or transitive) reaching a strictly higher stratum. Exempt: `tkt-`/`ses-` tokens and draft knowledge artifacts (work items own no content-state and are never exempt via the draft clause). Non-blocking by design: warnings never affect the verdict or exit code. Self-validation of this repository shows **8 R10 warnings** — 8 Implementation ADRs at Architecture stratum with no chain to Discovery (the count grew from 7 with ADR-008; see the [migration report §7](migration-report-engineering-domains-v1.1.md) and the [Verification statement](#verification-statement) below).

**Escalation path.** R10 severity evolves in three governed phases. Every phase change is taxonomy/validation governance (Core §14.2) and is announced in the release notes of the minor that ships it — never a silent change:

| Phase | Version | Severity | Notes |
|---|---|---|---|
| Guidance mode | v1.1 (now) | Warning | Non-blocking; repositories adapt at their own pace, adding chains where the knowledge warrants it; a warning may be accepted and documented instead of resolved. The escalation criteria below define when the next phase becomes eligible. |
| Blocking | v1.2 (next minor) | Blocking, same exemptions | Repositories have had one minor cycle to adapt. Migration note: stratum-isolated artifacts must add an upward chain — e.g., an Architecture artifact `derives-from` a Discovery artifact, or a traceability artifact `trc-` documents the linkage — or accept the exemption conditions. |
| Invariant | Future | Unconditional (exemptions re-evaluated) | Traceability elevated to a taxonomy invariant — strengthen-only, like the Stratum Authority Invariant (Core §8.1, §16.3); enforcement becomes unconditional. |

**Escalation criteria (all must hold before v1.2 makes R10 blocking):**

1. **Warning uptake** — a majority of conformant repositories show 0 R10 warnings.
2. **Tooling support** — `eka view` or the documentation provides a way to inspect chains (traceability projection).
3. **Guidance published** — this roadmap plus a workflow guide explain how to satisfy R10.
4. **Transition window** — one minor-version cycle (v1.1 → v1.2) has elapsed.

**No hidden breaking changes.** The roadmap is published now, in the spec-adjacent governance documentation; each escalation step is a governed change announced in the release notes of the minor that ships it.

**Rationale.** Warning-first avoids silent breaking for existing repositories; blocking gives the rule teeth; invariant status makes traceability a permanent architecture property.

## Verification statement

- **362 tests pass** (`go test ./...`): rule units, parsing, filenames, references, state, domain ontology (R102013R12), CLI (exit codes, output determinism), and self-validation.
- `go vet ./...` clean.
- **Self-validation PASS**: the EKA repository passes its own validator — 8 artifacts, 0 errors, 8 warnings (R10 stratification traceability), exit 0 — codified as `TestReferenceImplementationConforms` (`conformance/self_validation_test.go`).
