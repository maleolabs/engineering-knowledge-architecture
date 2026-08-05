# docs/vocabulary/ — Dimensi Vocabulary

> Anchor EKA: Knowledge Layer — dimensi **vocabulary** (EKA 8).

## Tujuan

Dimensi vocabulary mewadahi istilah dan definisinya: istilah produk, teknis, dan domain. Satu istilah = satu `gls-`; definisi bersifat kanonis dan dirujuk oleh semua dimensi lain agar percakapan antarperan tidak ambigu.

## Yang Ada di Sini

| Token | Tipe | Format nama |
|---|---|---|
| `gls-` | Glossary/Term | `gls-<id>.md` |

## State Vector

| Tipe | Domain state yang dimiliki |
|---|---|
| `gls-` | `content-state`, `existence-state` |

Nilai `content-state`: `draft → review → approved → amended`. Nilai `existence-state`: `active → archived → retired`. Field lain = N/A.

## Struktur Konten yang Baik

Struktur wajib (keluarga dokumen pengetahuan):

- `## Purpose` — istilah apa yang didefinisikan.
- `## Content` — definisi, sinonim, non-istilah (apa yang bukan maknanya), dan contoh pemakaian.

## Konvensi Nama

`gls-<id>.md`, kebab-case, unik dalam (namespace, type). Tanpa akhiran `-v<nn>`. Contoh: `gls-namespace.md`, `gls-artefak.md`.

## Vocabulary ≠ Specifications

Dimensi vocabulary **mendefinisikan makna istilah**; dimensi specifications **menetapkan perilaku dan format teknis**. Aturan perilaku bukan kosakata — dan definisi kosakata bukan spesifikasi. `gls-` menjawab "apa artinya X"; `spec-` menjawab "bagaimana X bekerja". Jangan menulis spesifikasi ke dalam `gls-` (dan sebaliknya).

## Kepemilikan

| Peran | Tanggung jawab |
|---|---|
| Product Owner | pemilik istilah produk/domain |
| Tech Lead | pemilik istilah teknis |
| Semua peran | mengusulkan istilah baru |

## Terkait

- [specifications/](../specifications/) — spesifikasi memakai istilah yang terdefinisi di sini.
- [intent/](../intent/) — istilah kunci visi/strategi terdefinisi di sini.
- [standards/](../standards/) — standar memakai kosakata baku.
- [research/](../research/) — istilah baru dari riset didaftarkan di sini.
