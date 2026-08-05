# containers/ — Execution Container (`ctr-`)

> Anchor EKA: Operating Layer — state domain **Container State**.

## Tujuan

Container adalah unit eksekusi yang menaungi sekumpulan work item terhadap satu plan terkunci. Container memegang relasi agregasi (work item mana yang masuk gelombang ini) dan memproyeksikannya ke dalam tabel.

## Token & State Vector

| Token | Folder | State yang dimiliki |
|---|---|---|
| `ctr-` | `operating/containers/` | `container-state`, `existence-state` |

Nilai `container-state`: `active → completed`. Nilai `existence-state`: `active → archived → retired`. Work item tidak menyimpan state container; container tidak menyimpan Execution State work item (itu dimiliki work item).

## Tepat Satu Kontainer Aktif

Mutual exclusion: hanya **satu** `ctr-` dengan `container-state: active` pada satu waktu. Kontainer baru lahir hanya setelah yang lama `completed`. Lahirnya kontainer **atomik** dengan penguncian plan pendukungnya (`plan-` → `immutable`; lihat [protocol.md](../protocol.md) §4).

## `completed` = Transisi Turunan

`container-state: completed` bukan nilai yang ditulis sembarangan — ini **transisi turunan** (derived transition) yang dipicu oleh agregat Execution State seluruh work item di dalamnya: semua work item berstatus `done`. Saat agregat terpenuhi, pemilik container menulis transisi `active → completed` ke `change-log`.

## Tabel Work Item = Proyeksi

Bagian `## Work Items` adalah **proyeksi** — snapshot status work item pemiliknya:

```
> Generated — State Projection. Do NOT edit state here; refresh on read.
```

- Proyeksi di-refresh saat dibaca; jangan mengedit state di tabel ini.
- Konflik antara tabel dan owner state: **owner state yang menang** (lihat validasi, aturan 8).

## Snapshot Semantics

Container merekam snapshot konteks pada pembuatannya: plan terkunci (beserta `instance-version`-nya), daftar work item awal, dan ruang lingkup. Snapshot tidak berubah selama container hidup; perubahan ke depan terjadi pada artefak/instance berikutnya.

## Struktur Konten Wajib

- `## Objective` — tujuan eksekusi gelombang ini.
- `## Work Items` — tabel proyeksi work item (token, id, ringkasan, Execution State).
- `## Change Log` — catatan transisi state container.

## Konvensi Nama

`ctr-<id>.md`, kebab-case, unik dalam (namespace, type). Tanpa akhiran `-v<nn>`. Contoh: `ctr-gelombang-1.md`.

## Kepemilikan

| Peran | Tanggung jawab |
|---|---|
| Tech Lead | pemilik tunggal state container; penulis transisi `completed` |
| Engineers | pelaksana work item dalam container |

## Terkait

- [planning/](../../planning/) — container mengunci `plan-` (lock-atomic-with-generation).
- [work-items/](../work-items/) — unit kerja yang diagregasi.
- [projections/](../projections/) — `tkt-` memproyeksikan work item ke status per gelombang.
- [sessions/](../sessions/) — eksekusi dalam container dicatat per sesi.
