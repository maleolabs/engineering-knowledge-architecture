# Engineering Knowledge Architecture (EKA) v1.0 — Reference Implementation

Repositori ini adalah **Reference Implementation** dari **Engineering Knowledge Architecture (EKA) v1.0** — model konseptual kanonik untuk pengetahuan engineering (Artifact, Identity, State, lapisan, kontrak pertukaran).

> **Satu serialisasi, bukan arsitekturnya.** Repositori ini adalah ONE serialization (Git + Markdown) dari standard EKA — bukan standard itu sendiri, dan bukan arsitekturnya (EKA 1.3). Standard hidup utuh di `standard/`; yang dapat disalin ke proyek ada di `skeleton/`.

## Orientasi tiga zona

| Zona | Isi | Peran |
|---|---|---|
| [`standard/`](standard/README.md) | Teks kanonik EKA v1.0 (Core + Exchange + Naming and Terminology) + glosarium istilah kanonik | **Standard itu sendiri** (pra-lapisan: mendefinisikan lapisan, bukan artifact proyek) |
| [`skeleton/`](skeleton/README.md) | Struktur proyek yang dapat disalin — 12 folder dimensi (KB) + `operating/` (OS) + `exchange/` (EX) | **Serialisasi** (format Git+Markdown yang diimplementasikan repositori ini) |
| [`reference/`](reference/README.md) | Meta-dokumentasi implementasi ini: arsitektur serialisasi, migrasi, filosofi, ADR, traceability | **Dokumentasi** keputusan dan konvensi implementasi |

## Jalur adopsi

1. Salin `skeleton/docs/` ke root proyek.
2. Set `namespace` pada frontmatter seluruh artifact sesuai proyek.
3. Baca `skeleton/docs/README.md` — sumber kebenaran struktur, workflow, dan aturan.
4. Gunakan [`standard/`](standard/README.md) untuk konformasi dan [`reference/`](reference/README.md) untuk konvensi serialisasi.

## Dokumen kunci

- **Spesifikasi**: [`standard/eka-specification-v1.0.md`](standard/eka-specification-v1.0.md) — teks kanonik lengkap (16 section).
- **Exchange Specification**: [`standard/eka-exchange-specification-v1.0.md`](standard/eka-exchange-specification-v1.0.md) — Exchange Contract v1 (milestone 16.1.2): unit exchange, Identity/State/Relationship transport, versioning, import/export/sync semantics, round-trip guarantees. Konseptual — independen dari Git, Markdown, database, CLI, MCP.
- **Naming and Terminology Specification**: [`standard/eka-naming-and-terminology-specification-v1.0.md`](standard/eka-naming-and-terminology-specification-v1.0.md) — (Ratified) penamaan resmi ekosistem EKA: product identity, specification Families, reference components, tooling, repository naming, daftar terminologi deprecated.
- **Reference Serialization Format (RSF)**: [`reference/eka-reference-serialization-format-v1.0.md`](reference/eka-reference-serialization-format-v1.0.md) — (referensi, bukan normatif) satu serialisasi proyeksi kanonik dari Exchange Package Object Model; target serialisasi `eka export`/`eka import` + SDK masa depan.
- **Arsitektur serialisasi**: [`reference/reference-architecture.md`](reference/reference-architecture.md) — zona → lapisan, konvensi Identity/State, tabel 26 token, artifact rule.
- **Panduan migrasi**: [`reference/migration-guide.md`](reference/migration-guide.md) — peta legacy → baru + strategi langkah-demi-langkah.
- **Perubahan breaking**: [`reference/breaking-changes.md`](reference/breaking-changes.md) — 14 perubahan disengaja + catatan tooling.
- **Keputusan implementasi**: [`reference/adr-summary.md`](reference/adr-summary.md) — 7 Implementation ADR (accepted).
- **Dokumentasi CLI**: [`reference/cli.md`](reference/cli.md) — instalasi, penggunaan, exit codes, arsitektur tooling `eka`.
- **Catatan konformansi**: [`reference/conformance-notes.md`](reference/conformance-notes.md) — 29 keputusan interpretasi + gap yang diketahui (matriks traceability-nya telah dikonsolidasi ke Conformance Traceability Matrix).
- **Conformance Traceability Matrix**: [`reference/conformance-traceability-matrix.md`](reference/conformance-traceability-matrix.md) — **single source of truth** cakupan konformansi: Engineering Requirement → Specification → Conformance Rule → Implementation → Automated Test → Coverage Status → Notes.
- **Panduan kontribusi**: [`CONTRIBUTING.md`](CONTRIBUTING.md) — governance kontribusi + definisi "implementasi tidak lengkap" + review checklist.

## Tooling

**CLI EKA** (`eka`) adalah antarmuka resmi Engineering Knowledge Architecture (Naming §7). Framework: Cobra (adapter murni; logika bisnis di `bootstrap/` + `conformance/`).

- **`eka init`** — Repository Bootstrapper: menganalisis workspace, wizard adaptif, membangkitkan repositori EKA dari Reference Skeleton, lalu memvalidasi. Idempoten; dukung `--dry-run`.
- **`eka validate`** — validator konformitas kanonik: menjalankan 9 aturan R1–R9 dari `skeleton/docs/exchange/validation.md` secara mekanis (P16). Konformitas repositori tidak bergantung hanya pada review manual — `eka validate` memutuskan kepatuhan dengan deterministik, termasuk exit code untuk integrasi CI.
- **`eka completion`** — script completion bash/zsh/fish/powershell.

```sh
go build -o eka ./cmd/eka && ./eka validate .
```

Repositori ini **lolos suite konformansinya sendiri** (7 artifact, 0 error, exit 0) — dikodifikasi sebagai test `TestReferenceImplementationConforms`. Dokumentasi lengkap: [`reference/cli.md`](reference/cli.md); catatan interpretasi dan traceability aturan: [`reference/conformance-notes.md`](reference/conformance-notes.md).

**Governance.** Cakupan konformansi dikelola lewat [`reference/conformance-traceability-matrix.md`](reference/conformance-traceability-matrix.md) — single source of truth yang menelusuri setiap Conformance Rule (R0–R9) dari requirement spec hingga implementasi Go dan test otomatisnya. Matriks **wajib diperbarui dalam Pull Request yang sama** dengan setiap perubahan spec/rule/implementasi/test, dan sebaliknya. Seluruh aturan kontribusi dan definisi "implementasi tidak lengkap" ada di [`CONTRIBUTING.md`](CONTRIBUTING.md).

## Status

| Item | Status |
|---|---|
| EKA v1.0 | **Ratified** (2026-08-05; lihat [`reference/ratification-notes.md`](reference/ratification-notes.md)) |
| EKA Exchange Specification v1.0 | **Ratified** (refinement pass + ratification report: [`reference/ratification-notes-exchange-v1.0.md`](reference/ratification-notes-exchange-v1.0.md)) |
| EKA Naming and Terminology Specification v1.0 | **Ratified** (meta-specification — penamaan resmi ekosistem) |
| Referensi implementasi | Aktif (zona standard / skeleton / reference lengkap) |
| Serialisasi | Git + Markdown (satu-satunya serialisasi repositori ini) |
