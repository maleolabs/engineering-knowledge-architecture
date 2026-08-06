# Validation — Conformance Checklist (9 Rules)

> Anchor EKA: Exchange Layer — conformance validation. Convention document, not an artifact.
> Standard: EKA v1.0, dated 2026-08-05.

> Term mapping (Naming and Terminology Specification v1.0 §9.3): the "Conformance Rules 1–9" below = **Conformance Rules R1–R9** (EKA Exchange Specification §14.2). R0 (structural artifact-rule) is defined by the Reference Validator (`conformance/`), not part of these nine rules.

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

## Result

- [ ] All rules 1–9 pass → commit allowed.
- [ ] Any **0** → fix first, do not commit.
- [ ] Only **W** → commit allowed with the warnings noted.
