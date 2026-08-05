# Arsitektur Referensi — Serialisasi EKA v1.0 dalam Git + Markdown

Dokumen ini menjelaskan **bagaimana repositori ini menserialisasi EKA v1.0**: pemetaan zona → lapisan, konvensi serialisasi yang diimplementasikan, aturan artifact, dan keputusan implementasi kunci.

Repositori ini adalah **satu serialisasi (Git+Markdown) dari standard — bukan arsitekturnya** (EKA 1.3).

## 1. Pemetaan zona → lapisan

| Zona | Isi | Lapisan EKA | Catatan |
|---|---|---|---|
| `standard/` | Salinan teks kanonik EKA v1.0 + glosarium kanonik | **Pra-lapisan** (pre-layer) | Standard mendefinisikan lapisan; ia bukan artifact proyek mana pun. |
| `skeleton/docs/` — 12 folder dimensi knowledge | intent, requirements, architecture, decisions, specifications, standards, operations, quality, planning, records, research, vocabulary | **Knowledge Layer (KB)** | Content, klasifikasi, Relationship, Records, Identity. |
| `skeleton/docs/operating/` | containers (`ctr-`), work-items (`sto-`/`ts-`/`bug-`/`td-`/`ch-`/`spk-`), projections (`tkt-`), sessions (`ses-`), protocol | **Operating Layer (OS)** | State Domain (Execution, Planning, Container, Existence), Protocol, Gate, Command. |
| `skeleton/docs/exchange/` | `validation.md`, `transfer.md` | **Exchange Layer (EX)** | Kontrak round-trip, validasi kepatuhan, schema versioning. |
| `reference/` | Meta-dokumentasi implementasi ini | — | Dokumentasi konvensi, bukan bagian serialisasi proyek. |

Ketiga lapisan terikat oleh Identity `(Namespace, Type, ID, InstanceVersion)` dan 7 invariant global (EKA 5.4).

## 2. Konvensi serialisasi

### 2.1 Encoding Identity

- **Lokasi kebenaran**: frontmatter artifact — `namespace`, `type`, `id`, `instance-version`, `revision` (EKA 6.4, P3, P9).
- **Filename adalah proyeksi**: pola `<type-token>-<id>[-v<nn>].md` — token tipe eksplisit, ID bebas kolisi per `(Namespace, Type)`. Akhiran `-v<nn>` **wajib** untuk tipe berversi (`scp-`, `plan-`) — selalu, termasuk v1 — dan **dilarang** untuk tipe lain. Filename hanya proyeksi Identity untuk navigasi manusia + validasi konsistensi; Identity sejati hidup di frontmatter.
- **Tabel 26 token tipe** (bebas ambiguitas: tidak ada token yang menjadi prefiks token lain; pasangan anti-prefix: `sto-`/`str-`, `spk-`/`spec-`):

| Token | Tipe Artifact | Dimensi | Token | Tipe Artifact | Dimensi |
|---|---|---|---|---|---|
| `vis-` | Vision / Manifesto | Product Intent | `ses-` | Session | — (OS) |
| `str-` | Strategy | Product Intent | `rvw-` | Review | Governance & Quality |
| `req-` | Requirement (PRD) | Requirements | `adr-` | ADR | Decisions |
| `scp-` | Scope Definition | Planning + Requirements | `dec-` | Decision Record | Decisions |
| `epc-` | Epic | Planning Knowledge | `arc-` | Architecture Description | Architecture |
| `plan-` | Plan (roadmap) | Planning Knowledge | `spec-` | Specification | Specifications |
| `ctr-` | Execution Container | — (OS) | `std-` | Standard / Guideline | Standards & Guidelines |
| `tkt-` | Ticket | — (OS, proyeksi) | `run-` | Runbook / Operational Guide | Operational Knowledge |
| `sto-` | Work Item: Story | Requirements / Records / Research | `rel-` | Release Record | Records |
| `ts-` | Work Item: Technical Story | Requirements / Records / Research | `gls-` | Glossary / Term | Vocabulary |
| `bug-` | Work Item: Bug | Requirements / Records / Research | `trc-` | Traceability / Relationship | Planning Knowledge |
| `td-` | Work Item: Tech Debt | Requirements / Records / Research | `fnd-` | Research Finding (ekstensi, ADR-007) | Research |
| `ch-` | Work Item: Chore | Requirements / Records / Research | | | |
| `spk-` | Work Item: Spike | Requirements / Records / Research | | | |

### 2.2 Encoding State Vector

- **Lima field frontmatter**, satu per domain state yang **dimiliki** (EKA 7.4): `content-state`, `execution-state`, `planning-state`, `container-state`, `existence-state`.
- **Absence = not-applicable**: field hanya hadir untuk domain yang dimiliki tipe artifact tersebut (contoh: ADR = `content-state` + `existence-state`; work item = `execution-state` + `existence-state`; ticket = tidak ada field state sama sekali).
- **Value sets** (EKA 7.2, nilai lowercase):
  - `content-state` (varian per tipe, EKA 7.2): living `draft | review | approved | amended`; ADR `proposed | accepted | superseded`; decision record `draft | accepted | superseded`
  - `execution-state`: `planned | todo | in-progress | in-review | done`
  - `planning-state`: `draft | approved | immutable`
  - `container-state`: `active | completed` (completed = transisi derived)
  - `existence-state`: `active | archived | retired`
- **Single-writer per field** (P6): setiap field state memiliki tepat satu owner; tampilan lain adalah proyeksi.
- **`change-log`**: array `{date, domain, from, to, by}` — catatan wajib seluruh transisi state (EKA 5.2).

### 2.3 Phase sebagai konteks

- Field `phase` pada artifact `scp-`/`plan-` saja, nilai `discovery|mvp|milestone|release|growth|maturity|sunset` (EKA 11.2, ADR-004).
- Phase change = context update yang diotorisasi gate kesiapan; dicatat di `change-log` dengan `domain: phase`. Tidak ada folder phase.

### 2.4 Relationship

- Relasi disandikan by Identity di frontmatter: `supersedes`, `amends`, `derives-from`, `depends-on`, `validates` (EKA 6.2.7, 13.2.3). Referensi tidak pernah by lokasi (P3).

### 2.5 Klasifikasi

- Field `dimension` di frontmatter; artifact knowledge hidup di folder dimensinya — aturan `dimension == folder` ditegakkan validasi (EKA 8, P15, ADR-005). Artifact operating (`operating/`) dikecualikan.

### 2.6 Proyeksi

- Tabel container dan ticket (`tkt-`) adalah State Projection (EKA 7.4): ticket ber-State Vector kosong + `derives-from`, artifact generated ber-header "Generated — State Projection. Do NOT edit state here; refresh on read."; kebijakan refresh default on-read (EKA 15.5, ADR-003).

### 2.7 Well-formed content

- Content mengikuti struktur per tipe artifact (skeleton per folder) sehingga machine-parseable dan deterministik (EKA 3, 5.3).

## 3. Aturan artifact vs dokumen konvensi

> **Suatu file adalah Artifact iff frontmatter-nya memuat `type` DAN `id`.**

- **Artifact**: membawa Identity lengkap + State Vector sesuai tipenya; dikelola Operating Layer; dapat di-exchange.
- **Dokumen konvensi** (contoh: `README.md`, `operating/protocol.md`, `exchange/validation.md`, `exchange/transfer.md`): dokumen yang **menjelaskan** konvensi — tidak memiliki `type`/`id`, tidak membawa Identity, bukan bagian state machine, tidak di-exchange sebagai Artifact.

Dokumen konvensi dapat dikenali dari ketiadaan pasangan `type`+`id` di frontmatter.

## 4. Ringkasan keputusan implementasi kunci

| Keputusan | Anchor EKA | ADR |
|---|---|---|
| Identity di frontmatter; filename = proyeksi; token eksplisit | 6.4, P3, P9 | [ADR-001](decisions/adr-001-identity-serialization.md) |
| State = 5 field frontmatter per domain owned; change-log | 7.2, 5.2, P6 | [ADR-002](decisions/adr-002-state-vector-encoding.md) |
| Ticket/tabel container = proyeksi; refresh on-read | 7.4, 15.5, P6 | [ADR-003](decisions/adr-003-projection-model.md) |
| Phase = metadata frontmatter, bukan folder | 11.2, P3 | [ADR-004](decisions/adr-004-phase-as-metadata.md) |
| 12 folder = 12 dimensi 1:1 + operating/ + exchange/ | 8, P15 | [ADR-005](decisions/adr-005-dimension-layout.md) |
| Seam exchange = validation.md + transfer.md | 13, P13 | [ADR-006](decisions/adr-006-exchange-conventions.md) |
| Tipe extension `fnd-` (Research Finding) | 14.1, 14.2, 11.4 | [ADR-007](decisions/adr-007-extension-research-finding.md) |

## Referensi

- Standard kanonik: [`../standard/eka-specification-v1.0.md`](../standard/eka-specification-v1.0.md)
- Struktur yang dapat disalin: [`../skeleton/README.md`](../skeleton/README.md)
- Peta migrasi: [`migration-guide.md`](migration-guide.md)
- Perubahan breaking: [`breaking-changes.md`](breaking-changes.md)
