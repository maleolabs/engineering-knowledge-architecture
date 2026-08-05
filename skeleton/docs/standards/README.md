# docs/standards/ — Dimensi Standards

> Anchor EKA: Knowledge Layer — dimensi **standards** (EKA 8).

## Tujuan

Dimensi standards mewadahi standar dan pedoman (guideline) yang mengikat perilaku, gaya, dan kualitas karya: konvensi kode, format dokumen, proses, dan kriteria mutu. Standar menetapkan **aturan yang harus diikuti**, berbeda dengan runbook yang menjelaskan **prosedur pelaksanaan**.

## Yang Ada di Sini

| Token | Tipe | Format nama |
|---|---|---|
| `std-` | Standard/Guideline | `std-<id>.md` |

## State Vector

| Tipe | Domain state yang dimiliki |
|---|---|
| `std-` | `content-state`, `existence-state` |

Nilai `content-state`: `draft → review → approved → amended`. Nilai `existence-state`: `active → archived → retired`. Field lain = N/A.

## Struktur Konten yang Baik

Struktur wajib (keluarga dokumen pengetahuan):

- `## Purpose` — area apa yang diatur.
- `## Content` — aturan/konvensi; boleh berbentuk checklist kepatuhan.

## Konvensi Nama

`std-<id>.md`, kebab-case, unik dalam (namespace, type). Tanpa akhiran `-v<nn>`. Contoh: `std-frontmatter-identitas.md`.

## Konvensi ≠ Prosedur

Standar menetapkan **konvensi** (apa yang boleh/tidak boleh). Prosedur langkah-demi-langkah dijalankan oleh runbook (`run-`) di dimensi operations; jika sebuah dokumen menjelaskan "cara melakukan", ia bukan standar.

## Kepemilikan

| Peran | Tanggung jawab |
|---|---|
| Tech Lead | pemilik tunggal standar teknis |
| DevOps | pemilik standar operasional/deployment |
| Semua peran | wajib mematuhi `std-` yang disetujui |

## Terkait

- [operations/](../operations/) — prosedur (`run-`) dipisahkan dari standar di sini.
- [quality/](../quality/) — kepatuhan standar diverifikasi oleh `rvw-`.
- [specifications/](../specifications/) — spesifikasi mengikuti standar.
- [vocabulary/](../vocabulary/) — istilah yang dimuat standar harus terdefinisi.
