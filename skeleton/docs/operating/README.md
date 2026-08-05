# docs/operating/ — Operating Layer (OS)

> Anchor EKA: Operating Layer — state machine + protocol. OS memegang state; OS tidak pernah mengedit konten.

## Ringkasan

Folder ini menserialisasikan **state** dan **eksekusi** proyek: siapa pemilik state, bagaimana transisi terjadi, dan bagaimana state dibaca tanpa diedit.

- **State ownership (P6):** setiap state domain dimiliki oleh tepat satu tipe artefak (lihat State Vector di README tiap subfolder). Absennya field = N/A untuk tipe itu.
- **Single-writer:** hanya pemilik state yang menulis field statenya; setiap transisi dicatat di `change-log`. Tidak ada dua penulis untuk satu field.
- **Proyeksi tidak pernah menulis state:** tabel work item dalam kontainer dan `tkt-` hanyalah bayangan; di-refresh saat dibaca, bukan diedit.
- **Dua kanal perubahan:** tata kelola konten (content governance) dan protokol state (state protocol) adalah dua mekanisme terpisah yang tidak boleh dicampur.
- **Konten hidup di Knowledge Layer** (dimensi knowledge di `docs/`); OS hanya mengatur state artefak itu — OS tidak mengubah isinya.

## Dokumen

| File | Peran |
|---|---|
| [protocol.md](protocol.md) | Operating Manual — konvensi, bukan artefak (tanpa `type`/`id`) |

## Subfolder

| Folder | Isi | Anchor state domain |
|---|---|---|
| [work-items/](work-items/) | 6 subtipe unit kerja: `sto-`, `ts-`, `bug-`, `td-`, `ch-`, `spk-` | Execution State |
| [containers/](containers/) | `ctr-` Execution Container — tepat satu aktif | Container State |
| [sessions/](sessions/) | `ses-` Session — catatan eksekusi | Existence State |
| [projections/](projections/) | `tkt-` Ticket — state vector kosong (proyeksi murni) | — (tidak memiliki state) |

## Titik Kontak dengan Layer Lain

- Knowledge Layer: state konten (content-state) artefak pengetahuan dibaca OS saat menilai kesiapan, tetapi isinya dikelola oleh pemilik dimensi.
- Exchange Layer: validasi membaca state dari sini; lihat [../exchange/](../exchange/).
