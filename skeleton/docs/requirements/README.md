# docs/requirements/ — Dimensi Requirements

> Anchor EKA: Knowledge Layer — dimensi **requirements** (EKA 8).

## Tujuan

Dimensi requirements mewadahi kebutuhan yang harus dipenuhi produk, diturunkan dari intent dan disepakati bersama pemangku kepentingan. Setiap `req-` adalah satu kebutuhan yang dapat diuji keterpenuhannya oleh `rvw-` atau work item.

## Yang Ada di Sini

| Token | Tipe | Format nama |
|---|---|---|
| `req-` | Requirement | `req-<id>.md` |

## State Vector

| Tipe | Domain state yang dimiliki |
|---|---|
| `req-` | `content-state`, `existence-state` |

Nilai `content-state`: `draft → review → approved → amended`. Nilai `existence-state`: `active → archived → retired`. Field lain yang tidak terdaftar = N/A.

## Struktur Konten yang Baik

Struktur wajib (keluarga dokumen pengetahuan):

- `## Purpose` — kebutuhan apa yang dijelaskan.
- `## Content` — pernyataan kebutuhan, kriteria penerimaan, dan konteksnya.

## Konvensi Nama

`req-<id>.md`, kebab-case, unik dalam (namespace, type). Tanpa akhiran `-v<nn>`. Contoh: `req-login-email.md`.

## Amendemen = Instance Baru

Perubahan terhadap kebutuhan yang sudah disetujui **tidak** mengedit dokumen lama secara membabi buta: buat instance baru `req-<id-baru>.md` dengan `content-state: amended` pada dokumen lama (atau arkivkan yang lama) dan field `amends: [req:<id-lama>]` pada yang baru. Rantai amendemen dapat ditelusuri melalui `amends`.

## Kepemilikan

| Peran | Tanggung jawab |
|---|---|
| Product Owner | pemilik tunggal konten requirements |
| Engineers | peninjau kelayakan teknis saat `review` |
| Semua peran | mengusulkan kebutuhan baru |

## Terkait

- [intent/](../intent/) — sumber penurunan kebutuhan.
- [specifications/](../specifications/) — `req-` dirinci menjadi `spec-`.
- [quality/](../quality/) — `rvw-` dapat memvalidasi pemenuhan kebutuhan.
- [planning/](../planning/) — `scp-` menyeleksi kebutuhan dalam cakupan.
