# docs/quality/ — Dimensi Quality

> Anchor EKA: Knowledge Layer — dimensi **quality** (EKA 8).

## Tujuan

Dimensi quality mewadahi hasil peninjauan dan verifikasi mutu: review kode, review arsitektur, dan review produk. Setiap `rvw-` memvalidasi satu atau lebih artefak lain melalui relasi `validates`, sehingga mutu dapat ditelusuri balik ke apa yang diverifikasi.

## Yang Ada di Sini

| Token | Tipe | Format nama |
|---|---|---|
| `rvw-` | Review | `rvw-<id>.md` |

## State Vector

| Tipe | Domain state yang dimiliki |
|---|---|
| `rvw-` | `content-state`, `existence-state` |

Nilai `content-state`: `draft → review → approved → amended`. Nilai `existence-state`: `active → archived → retired`. Field lain = N/A.

## Struktur Konten yang Baik

Struktur wajib (keluarga dokumen pengetahuan, dengan ekstensi review):

- `## Purpose` — objek dan lingkup review.
- `## Content` — uraian umum hasil review.
- `## Findings` — temuan: masalah, risiko, pelanggaran standar.
- `## Action Items` — tindak lanjut yang diperlukan.

Relasi `validates: [<type>:<id>]` harus menunjuk artefak yang direview (mis. `spec:login`, `std:frontmatter`).

## Konvensi Nama

`rvw-<id>.md`, kebab-case, unik dalam (namespace, type). Tanpa akhiran `-v<nn>`. Contoh: `rvw-serialisasi-frontmatter.md`.

## Kepemilikan

| Peran | Tanggung jawab |
|---|---|
| Engineers | peer review (kode/spesifikasi) |
| Tech Lead | review arsitektur; pemilik dimensi quality |
| DevOps | review operasional/deployment |
| Product Owner | review produk (pemenuhan kebutuhan) |

## Terkait

- [standards/](../standards/) — kepatuhan diverifikasi oleh `rvw-`.
- [specifications/](../specifications/) — `rvw-` memvalidasi `spec-`.
- [decisions/](../decisions/) — temuan review dapat melahirkan `dec-`.
- [operating/work-items/](../operating/work-items/) — action items dapat menjadi work item baru.
