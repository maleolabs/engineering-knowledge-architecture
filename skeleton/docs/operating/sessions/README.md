# sessions/ — Session (`ses-`)

> Anchor EKA: Operating Layer — state domain **Existence State**; artefak sesi.

## Tujuan

Session adalah catatan ephemeral dari satu kesempatan kerja: apa yang dikerjakan, konteks saat itu, catatan implementasi, dan verifikasi yang dilakukan. Sesi **tidak membawa status kemajuan permanen** — kemajuan dipegang oleh state domain lain (Execution State work item).

## Token & State Vector

| Token | Folder | State yang dimiliki |
|---|---|---|
| `ses-` | `operating/sessions/` | `existence-state` |

Hanya satu domain yang dimiliki: `existence-state` (`active → archived → retired`). Field state lain (content-state, execution-state, planning-state, container-state) **tidak berlaku** (N/A) untuk sesi.

## "Completed" = Kondisi Turunan

Sesi tidak memiliki nilai "completed" di frontmatter. "Completed" adalah **kondisi turunan** — sesi dianggap selesai ketika kerja yang dicatatnya sudah terwakili pada work item (menuju `done`) dan verifikasinya tercatat. Setelah itu sesi hanya memindahkan `existence-state` ke `archived` (atau `retired`).

## Konten Ephemeral

Sesi bersifat ephemeral — ditulis cepat, dibaca cepat, dan tidak menjadi rujukan permanen. Struktur wajib:

- `## Context` — apa yang sedang dikerjakan dan konteksnya (rujuk `ctr-`/work item).
- `## Notes` — catatan eksekusi, keputusan mikro, kendala.
- `## Verification` — bukti verifikasi yang dilakukan (build, test, manual check).

## Distilasi Wajib Sebelum Arsip (EKA 11.4)

Sebelum `ses-` diarsipkan, temuan yang layak **wajib didistilasi**: keputusan → `dec-`/`adr-`; temuan riset → `fnd-`; prosedur yang terbukti → `run-`; work item baru → dibuat di kontainer aktif. Sesi yang diarsipkan dengan temuan yang belum terdistilasi melanggar protokol.

## Konvensi Nama

`ses-<id>.md`, kebab-case, unik dalam (namespace, type). Tanpa akhiran `-v<nn>`. Contoh: `ses-2026-08-05-implementasi-proyeksi.md`.

## Kepemilikan

| Peran | Tanggung jawab |
|---|---|
| Engineer (penulis sesi) | pemilik tunggal `existence-state` sesi |
| Semua peran | menulis sesi untuk kerja yang dilakukannya |

## Terkait

- [work-items/](../work-items/) — sesi melayani eksekusi work item.
- [containers/](../containers/) — sesi terjadi dalam konteks `ctr-` aktif.
- [decisions/](../../decisions/) dan [research/](../../research/) — muara distilasi sesi.
