# docs/specifications/ — Dimensi Specifications

> Anchor EKA: Knowledge Layer — dimensi **specifications** (EKA 8).

## Tujuan

Dimensi specifications mewadahi spesifikasi yang dapat diimplementasikan: rincian perilaku, antarmuka, format, dan aturan teknis yang cukup presisi untuk menjadi dasar implementasi dan pengujian. Spesifikasi menjembatani kebutuhan (`req-`) menuju implementasi.

## Yang Ada di Sini

| Token | Tipe | Format nama |
|---|---|---|
| `spec-` | Specification | `spec-<id>.md` |

## State Vector

| Tipe | Domain state yang dimiliki |
|---|---|
| `spec-` | `content-state`, `existence-state` |

Nilai `content-state`: `draft → review → approved → amended`. Nilai `existence-state`: `active → archived → retired`. Field lain = N/A.

## Struktur Konten yang Baik

Struktur wajib (keluarga dokumen pengetahuan):

- `## Purpose` — bagian sistem apa yang dispesifikasikan.
- `## Content` — spesifikasi rinci: perilaku, input/output, batasan, format.

## Konvensi Nama

`spec-<id>.md`, kebab-case, unik dalam (namespace, type). Tanpa akhiran `-v<nn>`. Contoh: `spec-ticket-projection.md`.

## Kepemilikan

| Peran | Tanggung jawab |
|---|---|
| Tech Lead | pemilik tunggal spesifikasi teknis |
| Engineers | kontributor dan penerap |
| Product Owner | peninjau kesesuaian dengan kebutuhan |

## Terkait

- [requirements/](../requirements/) — spesifikasi menurunkan `req-`.
- [architecture/](../architecture/) — spesifikasi tunduk pada batasan `arc-`.
- [vocabulary/](../vocabulary/) — istilah dalam spesifikasi wajib terdefinisi (Vocabulary ≠ Specifications).
- [quality/](../quality/) — `rvw-` memvalidasi spesifikasi lewat `validates`.
- [standards/](../standards/) — spesifikasi mengikuti `std-`.
