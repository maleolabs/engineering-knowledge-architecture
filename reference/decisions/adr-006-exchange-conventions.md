---
namespace: eka-ref-impl
type: adr
id: 006-exchange-conventions
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
  - 005-dimension-layout
change-log:
  - date: 2026-08-05
    domain: existence-state
    from: "-"
    to: active
    by: Engineering Architecture
  - date: 2026-08-05
    domain: content-state
    from: proposed
    to: accepted
    by: Engineering Architecture
---

# ADR-006 — Konvensi Exchange: validation.md + transfer.md sebagai seam lapisan EX

## Context

Visi Knowledge OS (EKA 1.4, 16.1) mensyaratkan seam pertukaran yang didefinisikan di level standard: repositori harus dapat diimpor/diekspor tanpa kehilangan Identity, State, Content, atau Relationship (EKA 13.1, P13). Tanpa konvensi exchange, repositori ini hanya gudang file — tidak siap menjadi konsumen/produsen Artifact bagi sistem eksternal.

## Decision

Seam exchange diwujudkan sebagai dua dokumen konvensi di `skeleton/docs/exchange/`:

1. **`validation.md`** — 9 aturan validasi kepatuhan (diekstrak dari kontrak standard):
   1. Identity lengkap & unik: `(namespace, type, id)` unik; `instance-version` unik per Line (6.2.2).
   2. Identity frontmatter canonical dan machine-parseable (6.2.6).
   3. Nilai state valid terhadap value set domain masing-masing (7.2).
   4. Transisi forward-only; `change-log` konsisten dengan state saat ini (P7, 5.2).
   5. Single-writer: tidak ada dua owner untuk satu field state (P6).
   6. Referensi by Identity; tidak ada dangling reference (6.2.3, 5.1).
   7. Klasifikasi: `dimension == folder` untuk artifact knowledge (8, P15).
   8. Phase: nilai valid + hanya pada `scp-`/`plan-` (11.2).
   9. Content Well-formed per tipe artifact (3, 5.3).
2. **`transfer.md`** — kontrak transfer, mengikuti EKA 13.2:
   - **Round-trip lossless** (13.2.1) dan **idempotent**: re-import = no-op (13.2.2);
   - **Referential integrity** lintas sistem (13.2.3);
   - **Kebijakan konflik Identity**: import dengan Identity yang sudah ada = **tolak atau re-namespace eksplisit** — tidak pernah merge diam-diam (13.2.4);
   - **Validasi sebelum commit** (13.2.5);
   - **Schema versioning**: kontrak exchange berversi; import/export menyatakan versi kontrak yang dipatuhi (13.2.6).

Kedua dokumen adalah dokumen konvensi (tanpa `type`/`id`) — bukan Artifact; mereka menjelaskan kontrak, tidak membawa state.

## Consequences

- **Positif**: repositori import/export-ready tanpa redesign — seam exchange sudah didefinisikan sejak awal (EKA 13).
- **Positif**: validator mekanis dapat dibangun dari 9 aturan `validation.md` (P16).
- **Positif**: konflik Identity tidak pernah diselesaikan diam-diam — invariant round-trip lossless terjaga (P13).
- **Negatif**: konvensi ini mengikat — setiap struktur baru harus lulus validasi; ekspor/import wajib menyatakan versi kontrak.

## Alternatives Considered

- **Tanpa seam exchange** — ditolak: EKA 1.3/13.1 mewajibkan dukungan exchange lossless pada setiap implementasi.
- **Eksporter/importer bespoke yang ditulis belakangan** — ditolak: seam harus didefinisikan di level kontrak sejak awal, bukan retrofit; integrasi Knowledge OS membutuhkan batas eksplisit (EKA 4.2).
- **Konvensi exchange di standard, bukan di repositori** — ditolak: standard menetapkan kontrak; repositori menetapkan serialisasi konkretnya (EKA 12.4, 13.3).

## References

- EKA 1.3, 1.4, 4.2, 13.1 (yang harus dipertahankan), 13.2 (round-trip requirements), 13.3 (kontrak format serialisasi)
- Prinsip P13 (Lossless Exchange), P16 (Enforcement Capability Varies, Invariants Don't)
- Terkait: [ADR-005](adr-005-dimension-layout.md), [ADR-007](adr-007-extension-research-finding.md)
