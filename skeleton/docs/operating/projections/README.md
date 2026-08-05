# projections/ — Ticket (`tkt-`)

> Anchor EKA: Operating Layer — proyeksi murni; **tidak memiliki** state domain apa pun.

## Tujuan

Ticket adalah proyeksi satu work item ke bentuk ringkas: perintah (commands) yang dapat dieksekusi dan status yang diproyeksikan. Ticket **bukan pemilik state** — ia bayangan dari owner state work item di belakangnya.

## Token & State Vector

| Token | Folder | State yang dimiliki |
|---|---|---|
| `tkt-` | `operating/projections/` | **tidak ada (state vector kosong)** |

Ticket tidak memuat field state apa pun (`content-state`, `execution-state`, `planning-state`, `container-state`, `existence-state` semuanya N/A). Status yang terlihat di ticket hanyalah salinan untuk dibaca.

## Relasi

Ticket mewakili satu work item: `derives-from: [ctr:<id>, <type-work-item>:<id>]` — menunjuk kontainer aktif dan work item sumbernya. Referensi harus dapat di-resolve; jika work item sumber berubah status, ticket wajib di-refresh (lihat di bawah).

## Konten

- `## Commands` — perintah deterministic per jenis work item (mis. "jalankan test untuk X", "perbarui changelog") — tetap sama setiap refresh.
- `## Projected Status` — status yang diproyeksikan dari owner state work item (Execution State, dsb.) — diisi ulang setiap refresh.

Header wajib di bagian atas:

```
> Generated — State Projection. Do NOT edit state here; refresh on read.
```

## Refresh on Read

Proyeksi di-refresh setiap kali dibaca (default). Mekanisme event-driven otomatis (proyeksi diperbarui saat state berubah) diserahkan pada tooling masa depan; sampai tooling itu ada, **proyeksi adalah hasil pembacaan saat itu juga** — bukan sumber kebenaran.

## Jangan Edit State di Proyeksi

Mengubah status di ticket **tidak mengubah** work item dan merupakan pelanggaran single-writer (P6). Perubahan status hanya dilakukan oleh pemilik state work item; setelah itu ticket di-refresh.

## Konvensi Nama

`tkt-<id>.md`, kebab-case, unik dalam (namespace, type). Tanpa akhiran `-v<nn>`. Contoh: `tkt-sto-login-email.md`.

## Kepemilikan

| Peran | Tanggung jawab |
|---|---|
| OS (proyeksi) | "penulis" isi ticket — dihasilkan ulang, bukan diedit manual |
| Semua peran | hanya membaca; lapor anomali (ticket ≠ owner state) ke pemilik state |

## Terkait

- [work-items/](../work-items/) — sumber owner state yang diproyeksikan.
- [containers/](../containers/) — tabel work item kontainer adalah proyeksi sejenis.
- [validation.md](../../exchange/validation.md) — aturan 8: proyeksi tidak memuat state milik artefak lain.
