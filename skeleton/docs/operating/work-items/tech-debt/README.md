# tech-debt/ — Tech Debt (`td-`)

> Anchor EKA: Operating Layer — state domain Execution State; subtipe work item `td-`.

## Tujuan

Tech Debt adalah unit kerja untuk utang teknis yang telah diidentifikasi: shortcut, komponen usang, dan ketidakkonsistenan yang membebani pengembangan berikutnya. Setiap `td-` harus mencatat alasan utang itu ada (rasional) agar pengambilannya dapat diprioritaskan secara sadar.

## Token & State Vector

| Token | Folder | State yang dimiliki |
|---|---|---|
| `td-` | `work-items/tech-debt/` | `execution-state`, `existence-state` |

Nilai `execution-state`: `planned → todo → in-progress → in-review → done` (forward-only, never skip, never revert). Nilai `existence-state`: `active → archived → retired`.

## Struktur Konten Wajib

- `## Description` — bentuk utang teknis dan lokasinya.
- `## Acceptance Criteria` — kondisi yang membuktikan utang terbayar.
- `## Debt Rationale` — mengapa utang ini diambil (keputusan/waktu), agar pengambilan keputusan di masa depan kontekstual.

## Konvensi Nama

`td-<id>.md`, kebab-case, unik dalam (namespace, type). Tanpa akhiran `-v<nn>`. Contoh: `td-migrasi-ke-single-writer.md`.

## Kepemilikan

| Peran | Tanggung jawab |
|---|---|
| Engineer (implementer) | pemilik tunggal state; eksekusi |
| Tech Lead | pengambil keputusan prioritas utang |
| Product Owner | peninjau dampak terhadap rencana |

## Terkait

- [decisions/](../../../decisions/) — `Debt Rationale` dapat merujuk `dec-` yang melahirkannya.
- [standards/](../../../standards/) — pelunasan utang mengembalikan kepatuhan `std-`.
- [containers/](../containers/) — tech debt dimiliki kontainer aktif.
