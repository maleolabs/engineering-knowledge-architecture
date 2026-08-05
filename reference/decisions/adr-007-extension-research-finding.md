---
namespace: eka-ref-impl
type: adr
id: 007-extension-research-finding
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
  - 005-dimension-layout
  - 006-exchange-conventions
change-log:
  - date: 2026-08-05
    domain: content-state
    from: proposed
    to: accepted
    by: Engineering Architecture
---

# ADR-007 — Ekstensi: tipe Artifact `fnd-` (Research Finding)

## Context

Dimensi **Research** (EKA 8) — "Temuan investigasi, hasil riset teknis; wajib jalur Distillation ke dimensi durable" — **tidak memiliki tipe Artifact** di Artifact Taxonomy EKA 10. Pada implementasi awal, hasil investigasi (spike, riset teknis) tidak memiliki home durable: temuan hilang atau terpendam di session ephemeral. EKA 14.1 menyediakan mekanisme ekstensi ringan untuk tipe Artifact baru; EKA 14.2.6 mewajibkan tipe baru mendeklarasikan State Vector owned **lengkap** (tidak ada pewarisan default implisit).

## Decision

Daftarkan tipe extension **Research Finding** (`fnd-`) melalui mekanisme ekstensi EKA 14.1:

| Aspek | Deklarasi |
|---|---|
| Token tipe | `fnd-` (masuk tabel 26 token, ADR-001) |
| Knowledge Dimension | Research |
| State Vector owned (lengkap) | `(Content State, Existence State)` — domain lain not-applicable |
| Aturan Identity | Line + instance; ID unik dalam `(namespace, fnd)` |
| Folder | `research/` (aturan `dimension == folder`, ADR-005) |
| Relationship | `derives-from` (mis. dari spike `spk-`); output Distillation menuju dimensi durable (keputusan, ADR, Record) |

Konsekuensi governance: ekstensi terdaftar sebagai bagian standard (EKA 14.2.5: proposal → review → acceptance) dan **dapat di-exchange** (EKA 14.2.4) — tercakup schema versioning konvensi exchange (ADR-006).

## Consequences

- **Positif**: jalur Distillation spike → pengetahuan durable menjadi eksplisit (EKA 11.4): temuan riset (`fnd-`) didistilasi menjadi keputusan/ADR/Record, tidak menguap di session.
- **Positif**: dimensi Research (EKA 8) kini memiliki home artifact; hasil investigasi terpreservasi (P12).
- **Positif**: ekstensi sah dan exchangeable — tidak menyimpang dari invariant (EKA 14.2.1); backward compatible (14.2.2).
- **Negatif**: satu tipe baru wajib dijaga kepatuhannya — deklarasi State Vector owned lengkap (Content, Existence) tidak boleh berubah implisit (14.2.6).

## Alternatives Considered

- **Menggunakan tipe existing (`spec-`/`std-`/`rel-`)** — ditolak: Research ≠ Specifications/Standards/Records (EKA 8); Research akumulatif, bukan immutable, dan mewajibkan jalur Distillation.
- **Tanpa tipe khusus (temuan hanya di session)** — ditolak: session ephemeral by design (EKA 10); melanggar P12 (Preservation Over Deletion) dan EKA 11.4.
- **Research sebagai dimensi tanpa artifact type** — ditolak: dimensi tanpa tipe tidak dapat diproduksi/diekspor sebagai Artifact (EKA 14.1, 13.1).

## References

- EKA 8 (dimensi Research), 10 (Artifact Taxonomy), 11.4 (Distillation lifecycle), 14.1 (titik ekstensi), 14.2 (aturan ekstensi)
- Prinsip P12 (Preservation Over Deletion), P15 (Classification is Property, Not Identity)
- Terkait: [ADR-001](adr-001-identity-serialization.md), [ADR-005](adr-005-dimension-layout.md), [ADR-006](adr-006-exchange-conventions.md)
