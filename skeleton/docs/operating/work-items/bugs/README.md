# bugs/ — Bug (`bug-`)

> Anchor EKA: Operating Layer — state domain Execution State; subtipe work item `bug-`.

## Tujuan

Bug adalah unit kerja untuk perilaku yang menyimpang dari yang diharapkan: kegagalan fungsi, ketidaksesuaian dengan spesifikasi, atau regresi. Setiap bug harus dapat direproduksi dan ditelusuri ke artefak yang melanggarnya.

## Token & State Vector

| Token | Folder | State yang dimiliki |
|---|---|---|
| `bug-` | `work-items/bugs/` | `execution-state`, `existence-state` |

Nilai `execution-state`: `planned → todo → in-progress → in-review → done` (forward-only, never skip, never revert). Nilai `existence-state`: `active → archived → retired`.

## Struktur Konten Wajib

- `## Description` — gejala, langkah reproduksi, dan perilaku yang diharapkan.
- `## Impact` — dampak terhadap pengguna/sistem dan tingkat keparahan.
- `## Root Cause` (opsional) — akar masalah bila sudah teridentifikasi.

## Konvensi Nama

`bug-<id>.md`, kebab-case, unik dalam (namespace, type). Tanpa akhiran `-v<nn>`. Contoh: `bug-tiket-stale-setelah-refresh.md`.

## Kepemilikan

| Peran | Tanggung jawab |
|---|---|
| Pelapor | mendeskripsikan gejala dan dampak |
| Engineer (implementer) | pemilik tunggal state; diagnosis dan perbaikan |
| Tech Lead | peninjau perbaikan pada `in-review` |

## Terkait

- [specifications/](../../../specifications/) — perilaku yang benar dirujuk dari `spec-`.
- [quality/](../../../quality/) — perbaikan diverifikasi oleh `rvw-`.
- [containers/](../containers/) — bug dimiliki kontainer aktif.
