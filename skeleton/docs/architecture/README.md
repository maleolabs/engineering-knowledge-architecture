# docs/architecture/ — Dimensi Architecture

> Anchor EKA: Knowledge Layer — dimensi **architecture** (EKA 8).

## Tujuan

Dimensi architecture mewadahi deskripsi arsitektur sistem: struktur komponen, interaksi, batasan teknis, dan alasan keputusan arsitektural. Dokumen `arc-` menggambarkan "bagaimana sistem disusun", sementara "mengapa disusun demikian" hidup di `adr-` pada dimensi decisions.

## Yang Ada di Sini

| Token | Tipe | Format nama |
|---|---|---|
| `arc-` | Architecture Description | `arc-<id>.md` |

## State Vector

| Tipe | Domain state yang dimiliki |
|---|---|
| `arc-` | `content-state`, `existence-state` |

Nilai `content-state`: `draft → review → approved → amended`. Nilai `existence-state`: `active → archived → retired`. Field lain = N/A.

## Struktur Konten yang Baik

Struktur wajib (keluarga dokumen pengetahuan):

- `## Purpose` — cakupan arsitektur yang dideskripsikan.
- `## Content` — deskripsi komponen, interaksi, dan batasan (diagram boleh dirujuk lewat path).

## Konvensi Nama

`arc-<id>.md`, kebab-case, unik dalam (namespace, type). Tanpa akhiran `-v<nn>`. Contoh: `arc-identitas-namespace.md`.

## Kepemilikan

| Peran | Tanggung jawab |
|---|---|
| Tech Lead | pemilik tunggal konten architecture |
| Engineers | kontributor segmen arsitektur |
| DevOps | kontributor aspek deployment-infrastruktur |

## Terkait

- [decisions/](../decisions/) — `adr-` menjelaskan keputusan di balik `arc-`.
- [specifications/](../specifications/) — `arc-` diwujudkan menjadi spesifikasi rinci.
- [standards/](../standards/) — `std-` mengikat gaya dan konvensi teknis.
- [requirements/](../requirements/) — arsitektur memenuhi kebutuhan.
