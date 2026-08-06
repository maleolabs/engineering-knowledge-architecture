# Validation — Conformance Checklist (R0–R12)

> Anchor EKA: Exchange Layer — conformance validation. Convention document, not an artifact.
> Standard: EKA v1.1, dated 2026-08-05.

> Term mapping (Naming and Terminology Specification v1.1 §9.3): the "Conformance Rules 1–9" below = **Conformance Rules R1–R9** (EKA Exchange Specification §14.2). R0 (structural artifact-rule) is defined by the Reference Validator (`conformance/`), not part of these nine rules. **R10–R12** are the domain-aware rules of EKA Core Specification v1.1 §8.1 (Engineering Domains and Knowledge Stratification). The full rule set is **R0–R12** (thirteen rules): R0 structural, R1–R9 exchange (§14.2), R10–R12 domain-aware (Core v1.1 §8.1).

The following mechanical checklist runs **before commit** for every new artifact or state change. All rules are mechanical (automatable). Format: 1 = pass, 0 = fail (blocking), W = warning.

## Conformance Rule 1 — Identity Uniqueness

No duplicate `(namespace, type, id, instance-version)` may exist anywhere in the repository.
- [ ] No other artifact with the same combination.

## Conformance Rule 2 — Filename Consistency

- [ ] Filename token == `type` value in frontmatter.
- [ ] Version suffix on the filename (if any) == `instance-version` value.
- [ ] The `-v<nn>` suffix is **only** allowed for `scp-`/`plan-`; other types must not carry it.
- [ ] `scp-`/`plan-` **must** carry the `-v<nn>` suffix (including v1).

## Conformance Rule 3 — State Value Validity

Every state field value must be a member of its domain's value set:

| Domain | Value set |
|---|---|
| Content State (general) | draft, review, approved, amended |
| Content State (ADR) | proposed, accepted, superseded |
| Content State (decision) | draft, accepted, superseded |
| Execution State | planned, todo, in-progress, in-review, done |
| Planning State | draft, approved, immutable |
| Container State | active, completed |
| Existence State | active, archived, retired |
| Phase (context, scp-/plan-) | discovery, mvp, milestone, release, growth, maturity, sunset |

- [ ] All state field values ∈ their value sets.
- [ ] `phase` value ∈ phase set (only on `scp-`/`plan-`).

## Conformance Rule 4 — Owned-Set Conformance

The state fields present in a file must **exactly** match the owned domain set for its type (absence = N/A):

| Type | Owned set |
|---|---|
| vis-, str-, req-, scp-, epc-, trc-, arc-, adr-, dec-, spec-, std-, run-, rvw-, rel-, gls-, fnd- | content-state, existence-state |
| plan- | content-state, planning-state, existence-state |
| sto-, ts-, bug-, td-, ch-, spk- | execution-state, existence-state |
| ctr- | container-state, existence-state |
| ses- | existence-state |
| tkt- | (none — empty state vector) |

- [ ] No state fields owned by other types in the file (e.g. `container-state` on a work item = violation).
- [ ] `tkt-` carries no state fields at all.

## Conformance Rule 5 — Referential Integrity

All references (`amends`, `supersedes`, `derives-from`, `depends-on`, `validates`) must resolve to an existing artifact.

- [ ] Every reference points to an existing artifact (format `<type>:<id>[:<instance-version>]`, cross-namespace `<ns>/<type>:<id>`).
- [ ] References are written only on the referring artifact (not bidirectional).
- [ ] `content-state: draft` → unresolved references allowed (**W** warning).
- [ ] Non-draft `content-state` → unresolved reference = **0** (error).

## Conformance Rule 6 — Dimension == Folder

- [ ] Knowledge artifacts: `dimension` value == home folder (e.g. a file in `docs/specifications/` must have `dimension: specifications`).
- [ ] Operational artifacts (work items) may carry informational `dimension` — not evaluated.
- [ ] `ctr-`, `tkt-`, `ses-` **must not** carry `dimension`.

## Conformance Rule 7 — Change-Log Consistency

- [ ] Every owned domain's last `change-log` entry == current field value (per domain: content-state, execution-state, planning-state, container-state, existence-state, and `phase` on scp-/plan-).
- [ ] No transition without a `change-log` entry.

## Conformance Rule 8 — Single-Writer & Projections

- [ ] Projections (`tkt-`, work item tables in `ctr-`) do not carry owned state fields of other artifacts.
- [ ] Container work item tables match the owner state of the respective work item — **validated on read** (on mismatch: warning W; the source of truth remains owner state).
- [ ] The projection header (`> Generated — State Projection. Do NOT edit state here; refresh on read.`) is present on projection files.

## Conformance Rule 9 — Well-Formedness

Required content structure per type family:

| Family | Required sections |
|---|---|
| Planning artifact (scp-, epc-, plan-, trc-) | `## Objective`, `## Scope`, `## Out of Scope` |
| Work item (sto-, ts-, ch-) | `## Description`, `## Acceptance Criteria` |
| Bug (bug-) | `## Description`, `## Impact` |
| Tech Debt (td-) | `## Description`, `## Acceptance Criteria`, `## Debt Rationale` |
| Spike (spk-) | `## Description`, `## Investigation Notes`, `## Conclusion` (contains distillation links) |
| Decision record (adr-, dec-) | `## Context`, `## Decision`, `## Consequences`, `## Alternatives Considered` |
| Knowledge doc (vis-, str-, req-, arc-, spec-, std-, run-, rel-, gls-, fnd-) | `## Purpose`, `## Content` |
| Review (rvw-) | `## Purpose`, `## Content`, `## Findings`, `## Action Items` |
| Research Finding (fnd-) | `## Purpose`, `## Content`, `## Investigation Summary`, `## Conclusion` |
| Container (ctr-) | `## Objective`, `## Work Items`, `## Change Log` |
| Ticket (tkt-) | `## Commands`, `## Projected Status` |
| Session (ses-) | `## Context`, `## Notes`, `## Verification` |

- [ ] All required sections present for the respective type.
- [ ] ADR supersession: an `adr-` with `content-state: superseded` must be referenced by its successor (0 otherwise).

## Conformance Rule 10 — Stratification Traceability (warning)

Every artifact whose Engineering Domain is not **Discovery** (stratum 1) must have a resolvable reference chain — `derives-from`/`depends-on`, direct or transitive — reaching an artifact in a **strictly higher stratum** (Discovery → Architecture → Planning → Execution → Operations; stratum 1 = highest). The chain is satisfied when any reached artifact's home domain has a smaller stratum number; cycles and self-references are harmless (the walk is bounded).

Exemptions:

- `tkt-` and `ses-` tokens (pure projections / operating records).
- Knowledge artifacts with `content-state: draft` (work-item tokens own no `content-state` and are never exempt via this clause — they require the chain like every other non-draft artifact).

A missing chain is a **W** (warning): stratification is a structural quality signal — it never blocks a commit.

- [ ] Non-Discovery artifact has a resolvable `derives-from`/`depends-on` chain (direct or transitive) reaching a strictly higher stratum (W if missing).
- [ ] Exemptions applied only for `tkt-`/`ses-` and `content-state: draft` knowledge artifacts.

## Conformance Rule 11 — Domain Coherence (blocking)

The optional `domain` frontmatter field, **when present**, must:

1. Be one of the five canonical Engineering Domains (`Discovery`, `Architecture`, `Planning`, `Execution`, `Operations`).
2. Equal the artifact's **home domain** — the domain derived from its token family (e.g. `adr-` → `Architecture`).

`domain` **absent** = OK (the domain is derived, never part of Identity or the State Vector). Any violation is a **0** (error, blocking).

- [ ] `domain` present → value is canonical **and** matches the token's home domain (0 otherwise).
- [ ] `domain` absent → no check.

## Conformance Rule 12 — Cross-Stratum Supersession Prohibition (blocking)

A `supersedes` or `amends` relationship may **never** target an artifact in a **strictly higher stratum** (smaller stratum number): durable content moves down the authority chain, never up (Stratum Authority Invariant, Core v1.1 §8.1). Same-stratum and lower-stratum targets pass. Unresolvable targets are left to R5 (dangling references) — R12 evaluates resolvable targets only. Any violation is a **0** (error, blocking).

- [ ] No `supersedes`/`amends` target resolves to an artifact in a strictly higher stratum (0 otherwise).

## Result

- [ ] All rules R1–R12 pass, with R0 clean → commit allowed.
- [ ] Any **0** → fix first, do not commit (R11/R12 findings are blocking like R1–R9).
- [ ] Only **W** → commit allowed with the warnings noted (R10 stratification warnings included).
