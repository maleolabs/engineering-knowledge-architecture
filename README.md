# Project Docs Template — Reference Implementation EKA v1.0

Repositori ini adalah **Reference Implementation** dari **Engineering Knowledge Architecture (EKA) v1.0** — model konseptual kanonik untuk pengetahuan engineering (Artifact, Identity, State, lapisan, kontrak pertukaran).

> **Satu serialisasi, bukan arsitekturnya.** Repositori ini adalah ONE serialization (Git + Markdown) dari standard EKA — bukan standard itu sendiri, dan bukan arsitekturnya (EKA 1.3). Standard hidup utuh di `standard/`; yang dapat disalin ke proyek ada di `skeleton/`.

## Orientasi tiga zona

| Zona | Isi | Peran |
|---|---|---|
| [`standard/`](standard/README.md) | Teks kanonik EKA v1.0 + glosarium istilah kanonik | **Standard itu sendiri** (pra-lapisan: mendefinisikan lapisan, bukan artifact proyek) |
| [`skeleton/`](skeleton/README.md) | Struktur proyek yang dapat disalin — 12 folder dimensi (KB) + `operating/` (OS) + `exchange/` (EX) | **Serialisasi** (format Git+Markdown yang diimplementasikan repositori ini) |
| [`reference/`](reference/README.md) | Meta-dokumentasi implementasi ini: arsitektur serialisasi, migrasi, filosofi, ADR, traceability | **Dokumentasi** keputusan dan konvensi implementasi |

## Jalur adopsi

1. Salin `skeleton/docs/` ke root proyek.
2. Set `namespace` pada frontmatter seluruh artifact sesuai proyek.
3. Baca `skeleton/docs/README.md` — sumber kebenaran struktur, workflow, dan aturan.
4. Gunakan [`standard/`](standard/README.md) untuk konformasi dan [`reference/`](reference/README.md) untuk konvensi serialisasi.

## Dokumen kunci

- **Spesifikasi**: [`standard/eka-specification-v1.0.md`](standard/eka-specification-v1.0.md) — teks kanonik lengkap (16 section).
- **Arsitektur serialisasi**: [`reference/reference-architecture.md`](reference/reference-architecture.md) — zona → lapisan, konvensi Identity/State, tabel 26 token, artifact rule.
- **Panduan migrasi**: [`reference/migration-guide.md`](reference/migration-guide.md) — peta legacy → baru + strategi langkah-demi-langkah.
- **Perubahan breaking**: [`reference/breaking-changes.md`](reference/breaking-changes.md) — 14 perubahan disengaja + catatan tooling.
- **Keputusan implementasi**: [`reference/adr-summary.md`](reference/adr-summary.md) — 7 Implementation ADR (accepted).

## Status

| Item | Status |
|---|---|
| EKA v1.0 | **Ratified** (2026-08-05; lihat [`reference/ratification-notes.md`](reference/ratification-notes.md)) |
| Referensi implementasi | Aktif (zona standard / skeleton / reference lengkap) |
| Serialisasi | Git + Markdown (satu-satunya serialisasi repositori ini) |
