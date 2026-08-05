---
namespace: eka-ref-impl
type: adr
id: 001-identity-serialization
instance-version: 1
revision: 1
content-state: accepted
existence-state: active
dimension: decisions
author: Engineering Architecture
created: 2026-08-05
updated: 2026-08-05
supersedes: []
derives-from: []
depends-on: []
change-log:
  - date: 2026-08-05
    domain: content-state
    from: proposed
    to: accepted
    by: Engineering Architecture
---

# ADR-001 — Serialisasi Identity: token tipe eksplisit + Identity penuh di frontmatter

## Context

Implementasi awal menyandikan empat tipe Artifact (scope definition, plan, Execution Container, ticket) dalam **satu ruang ID bersama** dengan prefiks yang sama (`mvp-nnn`). Akibatnya, Type tidak dapat dibedakan secara deterministik dari representasi Identity: representasi yang sama dapat dibaca sebagai scope definition maupun plan. Ini adalah studi kasus pelanggaran EKA 6.4: aturan 6.2.1–6.2.2 dilanggar (Type tidak tegas → ID tidak unik per `(Namespace, Type)`), dan aturan 6.2.3 dilanggar secara konseptual (Identity disandikan melalui lokasi dan konvensi representasi, bukan properti Artifact). Pelajaran mengikat: **Identity tidak boleh disandikan dalam lokasi, tahap proses, atau konvensi representasi** (EKA 6.4).

## Decision

Serialisasi Identity pada repositori ini:

1. **Identity lengkap hidup di frontmatter**: `namespace`, `type`, `id`, `instance-version`, `revision` (EKA 6.4, P3, P9). Frontmatter adalah lokasi kebenaran Identity; referensi selalu by Identity.
2. **Filename adalah proyeksi Identity**, bukan Identity itu sendiri (P9): pola `<type-token>-<id>[-v<nn>].md`, dengan `<type-token>` = token tipe eksplisit, `<id>` = ID unik dalam `(Namespace, Type)`. Akhiran `-v<nn>` **wajib** untuk tipe berversi (`scp-`, `plan-`) — selalu, termasuk v1 — dan **dilarang** untuk tipe lain.
3. **Tabel 26 token tipe** (token bebas-ambiguitas: tidak ada token yang menjadi prefiks token lain; anti-prefix yang dikoreksi: `sto-`/`str-`, `spk-`/`spec-`):

| Token | Tipe Artifact | Dimensi |
|---|---|---|
| `vis-` | Vision / Manifesto | Product Intent |
| `str-` | Strategy | Product Intent |
| `req-` | Requirement (PRD) | Requirements |
| `scp-` | Scope Definition | Planning + Requirements |
| `epc-` | Epic | Planning Knowledge |
| `plan-` | Plan (roadmap) | Planning Knowledge |
| `ctr-` | Execution Container | — (OS) |
| `tkt-` | Ticket | — (OS, proyeksi) |
| `sto-` | Work Item: Story | Requirements / Records / Research |
| `ts-` | Work Item: Technical Story | Requirements / Records / Research |
| `bug-` | Work Item: Bug | Requirements / Records / Research |
| `td-` | Work Item: Tech Debt | Requirements / Records / Research |
| `ch-` | Work Item: Chore | Requirements / Records / Research |
| `spk-` | Work Item: Spike | Requirements / Records / Research |
| `ses-` | Session | — (OS) |
| `rvw-` | Review | Governance & Quality |
| `adr-` | ADR | Decisions |
| `dec-` | Decision Record | Decisions |
| `arc-` | Architecture Description | Architecture |
| `spec-` | Specification | Specifications |
| `std-` | Standard / Guideline | Standards & Guidelines |
| `run-` | Runbook / Operational Guide | Operational Knowledge |
| `rel-` | Release Record | Records |
| `gls-` | Glossary / Term | Vocabulary |
| `trc-` | Traceability / Relationship Artifact | Planning Knowledge |
| `fnd-` | Research Finding (ekstensi — ADR-007) | Research |

4. **Parsing deterministik**: representasi Identity dapat diparse tanpa ambiguitas dari frontmatter; filename divalidasi konsisten dengan frontmatter, bukan sebaliknya.

## Consequences

- **Positif**: parsing Identity deterministik; kolisi `mvp-nnn` terselesaikan (ID unik per `(Namespace, Type)` — EKA 6.2.2); Identity decoupled dari lokasi, tahap proses, dan phase (P3, P9).
- **Positif**: navigasi manusia tetap mudah (token tipe di filename) tanpa mengorbankan kebenaran Identity.
- **Negatif (disengaja)**: seluruh pola penamaan legacy putus (`mvp-*`, `sp-*`, dst.) — konsumen lama wajib bermigrasi (lihat `reference/breaking-changes.md`).
- **Negatif**: setiap file kini memuat frontmatter Identity yang wajib dijaga konsistensinya — ditutup oleh validasi mekanis (`dimension == folder`, token valid, dst., ADR-005/006).

## Alternatives Considered

- **Prefiks bersama + pembeda folder** (status quo legacy) — ditolak: Identity tetap disandikan via lokasi; melanggar P3/P9 dan EKA 6.4.
- **Type sebagai suffix** (mis. `id-type.md`) — ditolak: parsing ambigu dengan ID kebab-case multi-kata; glob per tipe sulit.
- **Type hanya di frontmatter, filename bebas** — ditolak: navigasi dan validasi manusia melemah; filename sebagai proyeksi yang konsisten membantu determinisme tanpa menjadi sumber kebenaran.

## References

- EKA 6.1 (komposisi Identity), 6.2 (aturan Identity), 6.3 (semantik versi), 6.4 (studi kasus kolisi)
- Prinsip P3 (Stable Identity), P9 (Structure as Projection of State)
- Terkait: [ADR-002](adr-002-state-vector-encoding.md) (state di frontmatter), [ADR-005](adr-005-dimension-layout.md) (folder = dimensi), [ADR-007](adr-007-extension-research-finding.md) (token `fnd-`)
