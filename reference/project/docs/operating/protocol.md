# Protocol — Operating Manual

> Anchor EKA: Operating Layer (OS). Convention document — not an artifact (no `type`/`id` in frontmatter).
> Standard: EKA v1.1, dated 2026-08-05.

This document is the operating manual of the EKA serialization: the ordering chain, State Domains, transition rules, and the obligations every state writer must obey.

## 1. Ordering Chain

Every value must be born in this order — steps must not be skipped:

```
requirement → scope → capability → plan → container → work item → session → review
```

| Step | Token | Folder | Notes |
|---|---|---|---|
| requirement | `req-` | `../requirements/` | agreed need |
| scope | `scp-` | `../planning/` | phased context (carries `phase`) |
| capability | `epc-` | `../planning/` | capability realizing the scope |
| plan | `plan-` | `../planning/` | execution plan; locks atomically with container |
| container | `ctr-` | `containers/` | active execution; exactly one at a time |
| work item | `sto-`/`ts-`/`bug-`/`td-`/`ch-`/`spk-` | `work-items/` | executed unit of work |
| session | `ses-` | `sessions/` | execution notes of work |
| review | `rvw-` | `../quality/` | result verification |

## 2. Five State Domains

Each domain has valid values, one initial value, one terminal value, and **forward-only** transitions (no skipping, no reverting).

| Domain | Valid values | Owned by | Initial → Terminal |
|---|---|---|---|
| Content State | `draft, review, approved, amended` — ADR variant: `proposed, accepted, superseded`; decision variant: `draft, accepted, superseded` | knowledge & planning artifacts (`vis-`, `str-`, `req-`, `scp-`, `epc-`, `plan-`, `trc-`, `arc-`, `adr-`, `dec-`, `spec-`, `std-`, `run-`, `rvw-`, `rel-`, `gls-`, `fnd-`) | draft/proposed → terminal (approved/accepted/amended/superseded) |
| Execution State | `planned, todo, in-progress, in-review, done` | `sto-`, `ts-`, `bug-`, `td-`, `ch-`, `spk-` | planned → done |
| Planning State | `draft, approved, immutable` | `plan-` | draft → immutable |
| Container State | `active, completed` | `ctr-` | active → completed |
| Existence State | `active, archived, retired` | universal — every type that carries state | active → archived → retired |

Transition rules (apply to all domains):

- **Forward-only:** only toward values ahead in the value order.
- **Never skip:** every sequential step must be passed (e.g. Execution State must not jump from `todo` straight to `done`; it must pass through `in-progress` and `in-review`).
- **Never revert:** values already passed must not be returned to (e.g. `approved` does not return to `draft`).
- **Single initial:** one valid initial value per domain.
- **Single terminal:** one valid terminal value per domain.
- **Recorded:** every transition is recorded in `change-log` by its single owner (see §6).

## 3. Exactly One Active Container

Only **one** `ctr-` may have `container-state: active` at a time (mutual exclusion). A new container may be created only after the previous one reaches `completed`. A newly born active container locks its plan (see §4).

## 4. Lock-Atomic-With-Generation

- `plan-` moves `draft → approved → immutable`.
- The transition of `plan-` to `immutable` happens **atomically** with the creation of `ctr-`: both are one operation; there must be no immutable plan without a container, and no container without a locked plan.
- Once locked, **any change to the plan does not edit that instance** — instead, a new generation is created: `plan-<id>-v<instance-version + 1>.md`. Each instance is a snapshot of the truth of its generation.
- A new generation must pass through the same chain (draft → approved → immutable).

## 5. Gates

Gates are evaluated against **owner state** (state owned by the artifact itself), not projections:

| Gate | Where | Pass condition |
|---|---|---|
| Approval gate | Knowledge Layer & planning | artifact reaches `content-state: approved` (or `accepted`) by its content owner |
| Readiness gate | work item → execution | work item is in the correct Execution State (e.g. `in-progress`) and its dependencies are satisfied |
| Review gate | before completion | work item is in `in-review` and the review (`rvw-`) validating it is approved |

If a gate evaluates a projection (e.g. a container table) and the result conflicts with owner state, **owner state wins**.

## 6. Change-Log Rules

- Every state-owning artifact has a `change-log:` field — an array of `{date, domain, from, to, by}` objects.
- Written by **one writer** (the state owner) on **every** transition — no transition without a record.
- The last entry per domain **must equal** the current field value (validated by exchange/validation.md).
- `domain` uses the State Domain name (e.g. `execution-state`, `phase`).

## 7. Two Change Channels

| Channel | Mechanism | Writes to |
|---|---|---|
| Content governance | content governance (amendment, content review, supersession) | artifact content, `revision` |
| State protocol | state protocol (transitions, `change-log`) | state fields, `instance-version` |

**The two channels must not be mixed:** changing state is not a way to change content, and changing content is not a way to change state. Content changes to a locked versioned artifact require a new generation (not direct editing).

## 8. Distillation Obligations

- Every session (`ses-`) producing findings **must be distilled** before archival: findings affecting direction → `dec-`/`adr-`; technical findings not yet decided → `fnd-` (research); proven procedures → `run-` (operations).
- Spike work items (`spk-`) must include a link to their knowledge distillation destination (research/decisions) in `## Conclusion`.
- Archiving without distillation violates the protocol (EKA 11.4).
