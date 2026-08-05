---
namespace: eka-ref-impl
type: adr
id: 005-dimension-layout
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

# ADR-005 — Layout Dimensi: 12 folder knowledge = 12 Knowledge Dimension, 1:1

## Context

Taksonomi legacy tercampur: 16 folder dengan isi bercampur dimensi (operations bercampur standards, planning sebagai catch-all, work-items/planning, specification-corpus sebagai misnomer vocabulary). EKA 8 menetapkan 12 Knowledge Dimension dengan pemisahan tegas (Operational ≠ Standards, Vocabulary ≠ Specifications, dst.); P15 menetapkan klasifikasi sebagai properti Artifact; P1 menetapkan Separation of Concerns antara pengetahuan dan eksekusi. Folder legacy mencampurkan ketiganya: klasifikasi, pipeline, dan lapisan.

## Decision

1. **12 folder knowledge = 12 Knowledge Dimension 1:1** (EKA 8): `intent`, `requirements`, `architecture`, `decisions`, `specifications`, `standards`, `operations`, `quality`, `planning`, `records`, `research`, `vocabulary`.
2. **`operating/`** — folder lapisan OS (containers, work-items, projections, sessions, protocol): state machine & protocol hidup di sini, bukan di folder dimensi (P1, EKA 4.1).
3. **`exchange/`** — folder lapisan EX (`validation.md`, `transfer.md`): kontrak pertukaran (EKA 13).
4. **Aturan lokasi**: artifact knowledge hidup di folder dimensinya; validasi menegakkan **`dimension == folder`** (field `dimension` di frontmatter harus sama dengan folder tempat file berada) — P15, P9.
5. **Artifact operating dikecualikan** dari aturan `dimension == folder` (work item berdimensi Requirements/Records/Research namun hidup di `operating/work-items/` — dimensi OS menentukan home-nya, bukan Knowledge Dimension-nya).
6. **Catch-all dibubarkan**: tidak ada folder yang menampung campuran dimensi (EKA 14.2, P15); content didistribusikan ke dimensi yang tepat.

## Consequences

- **Positif**: klasifikasi stabil dan dapat divalidasi mekanis — `dimension == folder` (P15); reklasifikasi tidak memutus referensi karena referensi by Identity (P3).
- **Positif**: pemisahan lapisan tegas (P1): knowledge di folder dimensi, eksekusi di `operating/`, exchange di `exchange/`.
- **Positif**: catch-all legacy hilang; setiap artifact memiliki home taksonomi yang jelas.
- **Negatif (disengaja)**: seluruh path legacy berubah (`breaking-changes.md` #1, #8, #9); migrasi mengikuti `migration-guide.md`.

## Alternatives Considered

- **Mempertahankan 16 folder legacy** — ditolak: pencampuran dimensi dan lapisan; melanggar P1, P15, EKA 8.
- **Subfolder per dimensi di bawah satu folder knowledge** — ditolak: 1:1 folder↔dimensi dipilih agar validasi `dimension == folder` sederhana dan deterministik (P16: enforcement bervariasi, invariant identik).
- **Klasifikasi hanya di frontmatter, folder bebas** — ditolak: folder adalah proyeksi klasifikasi (P9); tanpa proyeksi, navigasi dan validasi melemah.

## References

- EKA 8 (Knowledge Taxonomy — 12 dimensi), 14.2 (ekstensi; core closed, taxonomy open), 4.1 (Layer Model)
- Prinsip P1 (Separation of Concerns), P9 (Structure as Projection of State), P15 (Classification is Property, Not Identity)
- Terkait: [ADR-001](adr-001-identity-serialization.md), [ADR-006](adr-006-exchange-conventions.md), [ADR-007](adr-007-extension-research-finding.md)
