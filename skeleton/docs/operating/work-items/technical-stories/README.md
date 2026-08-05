# technical-stories/ — Technical Story (`ts-`)

> Anchor EKA: Operating Layer — state domain Execution State; subtipe work item `ts-`.

## Tujuan

Technical Story adalah unit kerja yang menghasilkan nilai teknis internal (infrastruktur, refactor, integrasi sistem) yang tidak langsung terlihat pengguna tetapi diperlukan untuk kualitas dan keberlanjutan sistem.

## Token & State Vector

| Token | Folder | State yang dimiliki |
|---|---|---|
| `ts-` | `work-items/technical-stories/` | `execution-state`, `existence-state` |

Nilai `execution-state`: `planned → todo → in-progress → in-review → done` (forward-only, never skip, never revert). Nilai `existence-state`: `active → archived → retired`.

## Struktur Konten Wajib

- `## Description` — pekerjaan teknis dan alasannya.
- `## Acceptance Criteria` — kondisi terukur yang memenuhi definisi selesai (termasuk verifikasi teknis).

## Konvensi Nama

`ts-<id>.md`, kebab-case, unik dalam (namespace, type). Tanpa akhiran `-v<nn>`. Contoh: `ts-migrasi-format-frontmatter.md`.

## Kepemilikan

| Peran | Tanggung jawab |
|---|---|
| Tech Lead | definisi pekerjaan teknis |
| Engineer (implementer) | pemilik tunggal state; eksekusi |
| DevOps | kolaborator untuk pekerjaan infrastruktur |

## Terkait

- [architecture/](../../../architecture/) — pekerjaan teknis mengejawantahkan `arc-`.
- [standards/](../../../standards/) — hasil harus mematuhi `std-`.
- [containers/](../containers/) — technical story dimiliki kontainer aktif.
