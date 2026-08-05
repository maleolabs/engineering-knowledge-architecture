---
namespace: eka-ref-impl
type: adr
id: 003-projection-model
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
  - 002-state-vector-encoding
change-log:
  - date: 2026-08-05
    domain: content-state
    from: proposed
    to: accepted
    by: Engineering Architecture
---

# ADR-003 — Model Proyeksi: tabel container dan ticket adalah State Projection

## Context

Tabel sprint dan dokumen ticket (wave docs) pada implementasi awal menyimpan status yang **diduplikasi** dari work item: status di tabel sprint, status di dokumen ticket, dan status di file work item — tiga salinan, dipertahankan sinkron secara manual tanpa penulis tunggal. EKA 7.4 menyelesaikan ini secara formal: representasi turunan adalah **State Projection** (tampilan derived, tidak memiliki State sendiri, tidak pernah menjadi writer — P6, P9), dan Artifact yang seluruh statenya diproyeksikan memiliki **State Vector kosong** (contoh: Ticket = `(∅)`). Kebijakan waktu validasi proyeksi masih merupakan open question (EKA 15.5: event-driven vs on-read).

## Decision

1. **Tabel work item pada container dan artifact ticket adalah State Projection** (EKA 7.4, 9 — Projection Semantics): statusnya diturunkan dari owner (work item yang direferensikan), divalidasi melalui Projection Refresh, dan **tidak pernah diedit sebagai fakta independen**.
2. **Ticket (`tkt-`) memiliki State Vector kosong**: tidak ada field state di frontmatter; seluruh state ticket adalah proyeksi dari work item yang direferensikan (EKA 10 — Ticket = `∅`).
3. **Relasi ditulis by Identity**: `derives-from: [ctr:<id>]` di frontmatter ticket; container menunjuk work item; rantai proyeksi selalu berujung pada owner (EKA 6.2.7).
4. **Artifact generated membawa header eksplisit**:

   > Generated — State Projection. Do NOT edit state here; refresh on read.

5. **Kebijakan refresh default: on-read** — proyeksi divalidasi terhadap owner saat dibaca (EKA 15.5); invariant "proyeksi tidak pernah menjadi writer" bersifat absolut (EKA 5.5).

## Consequences

- **Positif**: single-writer terjaga — tidak ada writer kedua pada status (P6); duplikasi status formal terhapus.
- **Positif**: tabel container/ticket dapat di-regenerate kapan saja dari owner tanpa kehilangan informasi.
- **Negatif**: pembaca proyeksi dapat melihat status basi sampai refresh on-read dilakukan — konsekuensi kebijakan yang dipilih, dikompensasi header peringatan.
- **Negatif**: tooling lama yang menulis status ke tabel/ticket putus secara disengaja (`breaking-changes.md` #4–5).

## Alternatives Considered

- **Tabel/ticket sebagai sumber status otoritatif** — ditolak: dua writer per field state; melanggar P6; mengulang duplikasi legacy.
- **Skrip sinkronisasi otomatis state owner → proyeksi** — ditolak: enforcement adalah kapabilitas implementasi (EKA 12.3, P16), tetapi kebenaran tetap di owner; proyeksi tetap proyeksi, bukan fakta kedua.
- **Ticket memiliki owned state sendiri** — ditolak: EKA 10 menetapkan Ticket = `(∅)`; status ticket adalah proyeksi atas work item yang direferensikan (resolusi ratifikasi Issue #1).

## References

- EKA 7.4 (State Vector; State Projection), 7.5 (interaksi), 9 (Execution Taxonomy — Projection Semantics), 10 (Artifact Taxonomy — Ticket)
- Prinsip P6 (Single Writer), P9 (Structure as Projection of State)
- EKA 15.5 (open question: kebijakan Projection Refresh)
- Terkait: [ADR-002](adr-002-state-vector-encoding.md)
