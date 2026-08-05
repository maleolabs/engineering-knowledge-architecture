# stories/ — Story (`sto-`)

> Anchor EKA: Operating Layer — state domain Execution State; subtipe work item `sto-`.

## Tujuan

Story adalah unit kerja yang memberikan nilai pengguna yang dapat diamati, dijelaskan dari sudut pandang pemangku kepentingan.

## Token & State Vector

| Token | Folder | State yang dimiliki |
|---|---|---|
| `sto-` | `work-items/stories/` | `execution-state`, `existence-state` |

Nilai `execution-state`: `planned → todo → in-progress → in-review → done` (forward-only, never skip, never revert). Nilai `existence-state`: `active → archived → retired`.

## Struktur Konten Wajib

- `## Description` — kebutuhan pengguna dan nilai yang diberikan.
- `## Acceptance Criteria` — kondisi terukur yang memenuhi definisi selesai.

## Konvensi Nama

`sto-<id>.md`, kebab-case, unik dalam (namespace, type). Tanpa akhiran `-v<nn>`. Contoh: `sto-login-email.md`.

## Kepemilikan

| Peran | Tanggung jawab |
|---|---|
| Product Owner | definisi nilai dan kriteria penerimaan |
| Engineer (implementer) | pemilik tunggal state; eksekusi |
| Tech Lead | peninjau teknis pada `in-review` |

## Terkait

- [requirements/](../../../requirements/) — story mewujudkan `req-`.
- [containers/](../containers/) — story dimiliki kontainer aktif.
- [projections/](../projections/) — status story diproyeksikan ke `tkt-`.
