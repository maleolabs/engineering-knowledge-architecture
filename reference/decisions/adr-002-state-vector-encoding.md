---
namespace: eka-ref-impl
type: adr
id: 002-state-vector-encoding
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
depends-on:
  - 001-identity-serialization
change-log:
  - date: 2026-08-05
    domain: content-state
    from: proposed
    to: accepted
    by: Engineering Architecture
---

# ADR-002 — Encoding State Vector: lima field frontmatter + single-writer + change-log

## Context

Implementasi awal menjalankan **tujuh mesin status** (living documents, ADR, decisions, roadmap, sprints, sessions, work items — EKA 7.3) dengan sinkronisasi status tiga arah: metadata di tabel dokumen, status di tabel sprint, dan status di dokumen ticket — tanpa penulis tunggal. Ini melanggar P6 (Single Writer) dan menciptakan duplikasi status yang tidak pernah konsisten. EKA 7.1 menolak domain State terunifikasi; EKA 7.2 menetapkan lima domain State independen; EKA 7.4 menetapkan State Vector = tuple domain yang **dimiliki** per tipe.

## Decision

1. **Lima field frontmatter, satu per domain state yang dimiliki**: `content-state`, `execution-state`, `planning-state`, `container-state`, `existence-state` (EKA 7.2). Field hanya hadir untuk domain yang dimiliki tipe Artifact (EKA 7.4, 10); **absence = not-applicable**.
2. **Single-writer per field** (P6): setiap field state memiliki tepat satu owner (Content State → Knowledge Layer; Execution/Planning/Container/Existence → Operating Layer). Tidak ada jalur penulisan state lain.
3. **`change-log`**: array `{date, domain, from, to, by}` di frontmatter — catatan kronologis wajib seluruh transisi state (EKA 5.2).
4. **Nilai kanonik lowercase** (value sets EKA 7.2):
   - `content-state` (varian per tipe): living `draft | review | approved | amended`; ADR `proposed | accepted | superseded`; decision record `draft | accepted | superseded`
   - `execution-state`: `planned | todo | in-progress | in-review | done`
   - `planning-state`: `draft | approved | immutable`
   - `container-state`: `active | completed` (completed = transisi derived)
   - `existence-state`: `active | archived | retired`
5. **Pemetaan nilai legacy**:

| Mesin legacy | Nilai legacy | Nilai baru | Domain |
|---|---|---|---|
| Living documents | Draft | `draft` | content-state |
| Living documents | Review | `review` | content-state |
| Living documents | Approved | `approved` | content-state |
| Living documents | Amended | `amended` | content-state |
| ADR | Proposed | `proposed` | content-state |
| ADR | Accepted | `accepted` | content-state |
| ADR | Superseded | `superseded` | content-state |
| Decision Record | Draft | `draft` | content-state |
| Decision Record | Accepted | `accepted` | content-state |
| Decision Record | Superseded | `superseded` | content-state |
| Roadmap | Draft | `draft` | planning-state |
| Roadmap | Approved | `approved` | planning-state |
| Roadmap | Immutable | `immutable` | planning-state |
| Sprints | Active | `active` | container-state |
| Sprints | Completed | `completed` | container-state (derived) |
| Sessions | Active | `active` | existence-state |
| Sessions | Archived | `archived` | existence-state |
| Work items | Planned | `planned` | execution-state |
| Work items | Todo | `todo` | execution-state |
| Work items | In Progress | `in-progress` | execution-state |
| Work items | In Review | `in-review` | execution-state |
| Work items | Done | `done` | execution-state |

6. **Kondisi derived** (mis. container/session "Completed") bukan nilai domain — dihitung dari agregat state owner (EKA 7.2), tidak ditulis sebagai fakta (lihat ADR-003).

## Consequences

- **Positif**: duplikasi status tereliminasi — satu field, satu owner, satu sumber kebenaran (P6).
- **Positif**: validasi mekanis menjadi mungkin: nilai domain tervalidasi terhadap value set; `change-log` konsisten dengan state saat ini; transisi forward-only dapat diperiksa (P7).
- **Negatif (disengaja)**: konsumen status legacy (kolom `Status:`, sinkronisasi tiga arah) putus — lihat `breaking-changes.md` #10–11.
- **Negatif**: setiap transisi state kini wajib tercatat di `change-log` — disiplin baru bagi Operating Layer.

## Alternatives Considered

- **Domain State terunifikasi (Artifact State)** — ditolak: monolit State adalah akar duplikasi status (EKA 7.1); tujuh semantik tidak dapat diekspresikan satu mesin (EKA 7.3).
- **Mempertahankan sinkronisasi tiga arah** — ditolak: tidak ada writer; melanggar P6.
- **State sepenuhnya derived (tanpa field owned)** — ditolak: OS membutuhkan state yang dimiliki untuk menjalankan Protocol; proyeksi tanpa owner tidak memiliki sumber kebenaran.

## References

- EKA 7.1 (evaluasi kandidat domain), 7.2 (domain formal + value sets), 7.3 (pemetaan mesin legacy), 7.4 (State Vector), 5.2 (change-log)
- Prinsip P2 (Explicit State), P6 (Single Writer), P7 (Forward-Only Transitions)
- Terkait: [ADR-001](adr-001-identity-serialization.md), [ADR-003](adr-003-projection-model.md)
