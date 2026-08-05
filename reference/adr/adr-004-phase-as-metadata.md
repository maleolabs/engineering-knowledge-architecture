---
namespace: eka-ref-impl
type: adr
id: 004-phase-as-metadata
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
  - 002-state-vector-encoding
change-log:
  - date: 2026-08-05
    domain: content-state
    from: proposed
    to: accepted
    by: Engineering Architecture
---

# ADR-004 — Phase sebagai Metadata (context attribute), bukan folder

## Context

Implementasi awal menyandikan phase sebagai **folder** (`docs/mvp/`, dst.): fase produk menjadi bagian lokasi, dan karena Identity lama disandikan via lokasi, fase ikut menjadi bagian representasi Identity. Ini melanggar EKA 11.2 (Phase adalah **context attribute** pada Artifact planning/scope — bukan kategori, bukan State Domain) dan P3 (Identity independen dari lokasi). EKA 7.1 juga menolak "Lifecycle State" sebagai domain tunggal: yang bertahan adalah Existence State (domain) + Phase (context).

## Decision

1. **Field `phase` di frontmatter**, hanya pada artifact `scp-` (Scope Definition) dan `plan-` (Plan) — sesuai EKA 11.2 (phase menempel pada Scope Definition / Plan sebagai context attribute).
2. **Value set**: `discovery | mvp | milestone | release | growth | maturity | sunset` (nilai EKA 11.2–11.3 dalam lowercase).
3. **Phase change adalah context update** yang diotorisasi **gate kesiapan** (EKA 11.2), dievaluasi atas agregat State owner (release-ready = semua work item Done ∧ seluruh container Completed ∧ plan Immutable ∧ gate review lulus ∧ gate persetujuan Content lulus).
4. **Pencatatan**: setiap phase change direkam di `change-log` dengan `domain: phase` (mis. `from: discovery, to: mvp`).
5. **Tidak ada folder phase** — artifact tidak berpindah lokasi saat phase berubah; Identity tidak tersentuh (P3).

## Consequences

- **Positif**: Identity decoupled dari phase (P3) — Scope Definition tetap identitas yang sama saat produk berpindah fase (EKA 11.2).
- **Positif**: phase change menjadi operasi metadata yang diaudit (change-log), bukan operasi pemindahan file.
- **Positif**: glob per fase tidak lagi diperlukan; query fase menjadi query field.
- **Negatif (disengaja)**: folder `mvp/` hilang; tooling yang membaca fase dari path putus (`breaking-changes.md` #12).

## Alternatives Considered

- **Folder per phase** (status quo legacy) — ditolak: phase menjadi bagian lokasi → bagian representasi Identity; melanggar EKA 6.4, P3.
- **Phase sebagai State Domain** — ditolak: EKA 7.1 menolak Lifecycle State; phase tidak memiliki transisi protocol, ia context yang berubah via gate (EKA 11.2).
- **Phase sebagai klasifikasi (Knowledge Dimension)** — ditolak: klasifikasi adalah properti retrieval (P15), bukan konteks waktu; phase tidak stabil sebagai sumbu klasifikasi.

## References

- EKA 3 (definisi Phase), 7.1 (Lifecycle State ditolak), 7.5 (gate kesiapan → phase change), 11.2 (Phase sebagai konteks), 11.3 (lifecycle produk)
- Prinsip P3 (Stable Identity), P9 (Structure as Projection of State)
- Terkait: [ADR-001](adr-001-identity-serialization.md), [ADR-002](adr-002-state-vector-encoding.md), [ADR-005](adr-005-dimension-layout.md)
