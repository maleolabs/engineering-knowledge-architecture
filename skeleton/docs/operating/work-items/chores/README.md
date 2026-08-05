# chores/ — Chore (`ch-`)

> Anchor EKA: Operating Layer — state domain Execution State; subtipe work item `ch-`.

## Tujuan

Chore adalah unit kerja administratif atau pemeliharaan yang tidak mengubah perilaku produk: pembaruan dependensi, konfigurasi, pembersihan, dan tugas rutin proyek.

## Token & State Vector

| Token | Folder | State yang dimiliki |
|---|---|---|
| `ch-` | `work-items/chores/` | `execution-state`, `existence-state` |

Nilai `execution-state`: `planned → todo → in-progress → in-review → done` (forward-only, never skip, never revert). Nilai `existence-state`: `active → archived → retired`.

## Struktur Konten Wajib

- `## Description` — tugas yang dikerjakan.
- `## Acceptance Criteria` — kondisi yang membuktikan tugas selesai.

## Konvensi Nama

`ch-<id>.md`, kebab-case, unik dalam (namespace, type). Tanpa akhiran `-v<nn>`. Contoh: `ch-perbarui-dependensi-ci.md`.

## Kepemilikan

| Peran | Tanggung jawab |
|---|---|
| Penugasan | siapa pun yang menugaskan (PO/Tech Lead/DevOps) |
| Engineer / DevOps (implementer) | pemilik tunggal state; eksekusi |

## Terkait

- [operations/](../../../operations/) — tugas rutin dapat merujuk `run-`.
- [containers/](../containers/) — chore dimiliki kontainer aktif.
