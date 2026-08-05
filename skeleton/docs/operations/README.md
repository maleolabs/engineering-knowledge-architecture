# docs/operations/ — Dimensi Operations

> Anchor EKA: Knowledge Layer — dimensi **operations** (EKA 8).

## Tujuan

Dimensi operations mewadahi pengetahuan operasional untuk menjalankan dan memelihara sistem: runbook, prosedur deployment, pemulihan, dan tugas rutin. Dokumen di sini adalah **prosedur** ("bagaimana melakukan"), dipisahkan dari **standar** ("aturan yang harus diikuti") yang hidup di dimensi standards.

## Yang Ada di Sini

| Token | Tipe | Format nama |
|---|---|---|
| `run-` | Runbook | `run-<id>.md` |

## State Vector

| Tipe | Domain state yang dimiliki |
|---|---|
| `run-` | `content-state`, `existence-state` |

Nilai `content-state`: `draft → review → approved → amended`. Nilai `existence-state`: `active → archived → retired`. Field lain = N/A.

## Struktur Konten yang Baik

Struktur wajib (keluarga dokumen pengetahuan):

- `## Purpose` — situasi/prosedur apa yang dijelaskan.
- `## Content` — langkah-langkah prosedur, prasyarat, dan hasil yang diharapkan.

## Konvensi Nama

`run-<id>.md`, kebab-case, unik dalam (namespace, type). Tanpa akhiran `-v<nn>`. Contoh: `run-deploy-staging.md`.

## Prosedur vs Standar

Pemisahan dilakukan pada EKA 8: **standards** (`std-`) menetapkan aturan dan konvensi; **operations** (`run-`) menjelaskan prosedur eksekusi. Runbook mengikuti standar, bukan menggantikannya. Jika sebuah runbook menetapkan aturan baru, aturan itu harus diangkat menjadi `std-`.

## Kepemilikan

| Peran | Tanggung jawab |
|---|---|
| DevOps | pemilik tunggal runbook |
| Tech Lead | peninjau runbook yang menyentuh arsitektur |
| Engineers | penulis runbook untuk komponen yang dibangunnya |

## Terkait

- [standards/](../standards/) — prosedur di sini tunduk pada `std-`.
- [records/](../records/) — `rel-` mencatat release yang dieksekusi dengan runbook.
- [quality/](../quality/) — efektivitas prosedur diverifikasi lewat `rvw-`.
- [sessions/](../sessions/) — temuan operasional dapat terdistilasi menjadi `run-` baru.
