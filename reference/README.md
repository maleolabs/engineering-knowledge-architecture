# Zona Reference — Dokumentasi Meta Implementasi EKA v1.0

Zona ini berisi **meta-dokumentasi** dari Reference Implementation EKA v1.0: bagaimana repositori ini menserialisasi standard, mengapa keputusan implementasi diambil, apa yang berubah dari struktur lama, dan jejak keputusan (ADR).

## Status dokumentasi

| Dokumen | Status | Catatan |
|---|---|---|
| [`reference-architecture.md`](reference-architecture.md) | Aktif | Arsitektur serialisasi repositori (zona → lapisan, konvensi serialisasi, artifact rule). |
| [`migration-guide.md`](migration-guide.md) | Aktif | Peta migrasi lengkap struktur legacy → struktur EKA + strategi langkah-demi-langkah. |
| [`philosophy.md`](philosophy.md) | Aktif | Narasi mengapa EKA ada dan mengapa repositori ini disusun demikian. |
| [`terminology-glossary.md`](terminology-glossary.md) | Aktif | Glosarium istilah level implementasi (istilah kanonik ada di `standard/glossary.md`). |
| [`breaking-changes.md`](breaking-changes.md) | Aktif | Ringkasan 14 perubahan breaking terhadap struktur legacy. |
| [`adr-summary.md`](adr-summary.md) | Aktif | Indeks 7 Implementation ADR (zona `decisions/`). |
| [`traceability-matrix.md`](traceability-matrix.md) | Aktif | Matriks penelusuran: setiap elemen repositori → anchor EKA. |
| [`ratification-notes.md`](ratification-notes.md) | Aktif | Catatan ratifikasi EKA v1.0 (verbatim dari stabilization pass). |
| [`cli.md`](cli.md) | Aktif | Dokumentasi CLI resmi `eka` (validator): instalasi, penggunaan, exit codes, proses validasi, arsitektur tooling, roadmap. |
| [`conformance-notes.md`](conformance-notes.md) | Aktif | Catatan implementasi konformansi: 29 keputusan interpretasi + gap yang diketahui (traceability tabel telah dikonsolidasi ke Conformance Traceability Matrix). |
| [`conformance-traceability-matrix.md`](conformance-traceability-matrix.md) | Aktif | Single source of truth cakupan konformansi: REQ→Spec→Rule→Impl→Test→Coverage→Notes (R0–R9, 54 test, 16 requirement). |
| [`decisions/`](decisions/) | Aktif | 7 Implementation ADR (keputusan arsitektur serialisasi). |

> Governance kontribusi + definisi "implementasi tidak lengkap" berada di [`../CONTRIBUTING.md`](../CONTRIBUTING.md) (root repo).

## Rujukan lintas zona

| Zona | Peran | Masuk lewat |
|---|---|---|
| **A — `../standard/`** | Standard kanonik EKA v1.0 (pra-lapisan) | [`../standard/README.md`](../standard/README.md) |
| **B — `../skeleton/`** | Struktur proyek yang dapat disalin (serialisasi Git+Markdown) | [`../skeleton/README.md`](../skeleton/README.md) |
| **C — `../reference/`** | Meta-dokumentasi implementasi ini | dokumen ini |
