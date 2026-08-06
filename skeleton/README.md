# <Product Name>

> EKA project template — copied by adopting projects, not modified here.
> Template date: 2026-08-05

## Identity to Fill In

| Item | How to fill |
|---|---|
| **Product name** | replace the `<Product Name>` title above |
| **`namespace`** | decide once; becomes the default `namespace` in every project artifact frontmatter |

`namespace` is the first part of the `(Namespace, Type, ID, InstanceVersion)` identity. Choose a stable value, unique within the organization, e.g. `anvil`. Do not change it after the first artifact is created.

## EKA Serialization Conformance

This project uses EKA v1.1 serialization (Git+Markdown).

## Source of Truth

[docs/README.md](docs/README.md) is the single source of truth for this project's EKA serialization. Read it before writing the first artifact. Structure, identity/state/relationship conventions, and ownership rules are described there; operational details live in [docs/operating/protocol.md](docs/operating/protocol.md).

## Workflow Chain

Every delivered value follows the ordering chain:

```
requirement → scope → capability → plan → container → work item → session → review
```

| Step | Token | Folder |
|---|---|---|
| requirement — a need | `req-` | `docs/requirements/` |
| scope — phased context | `scp-` | `docs/planning/` |
| capability — a capability | `epc-` | `docs/planning/` |
| plan — a plan to be locked | `plan-` | `docs/planning/` |
| container — execution (exactly one active) | `ctr-` | `docs/operating/containers/` |
| work item — unit of work | `sto-`/`ts-`/`bug-`/`td-`/`ch-`/`spk-` | `docs/operating/work-items/` |
| session — execution notes | `ses-` | `docs/operating/sessions/` |
| review — verification | `rvw-` | `docs/quality/` |

## Brief Contribution Rules

| Folder | Content owner | Notes |
|---|---|---|
| `docs/intent/` … `docs/vocabulary/` | per role in each folder's README | knowledge artifacts (Knowledge Layer); `dimension` must equal folder |
| `docs/operating/` | state owner per type (single-writer) | containers/tickets are projections; do not edit state in projections |
| `docs/exchange/` | contract | where validation & transfer rules live; not a place for domain artifacts |

Core rules: identity lives in frontmatter (filename is only a projection), state is written only by its single owner, projections are re-read not edited, and every state transition is recorded in `change-log`. Details: [docs/operating/protocol.md](docs/operating/protocol.md).
