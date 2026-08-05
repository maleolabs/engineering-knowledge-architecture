# <Nama Produk>

> Templat proyek EKA — disalin oleh proyek pengadopsi, bukan diubah di sini.
> Tanggal templat: 2026-08-05

## Identitas yang Harus Diisi

| Item | Cara mengisi |
|---|---|
| **Nama produk** | ganti judul `<Nama Produk>` di atas |
| **`namespace`** | tentukan sekali; menjadi default `namespace` pada seluruh frontmatter artefak proyek |

`namespace` adalah bagian pertama dari identitas `(Namespace, Type, ID, InstanceVersion)`. Pilih nilai yang stabil dan unik dalam organisasi, misalnya `anvil`. Jangan ganti setelah artefak pertama dibuat.

## Kepatuhan Serialisasi EKA

Proyek ini menggunakan serialisasi EKA v1.0 (Git+Markdown).

## Sumber Kebenaran

[docs/README.md](docs/README.md) adalah sumber kebenaran tunggal serialisasi EKA proyek ini. Baca sebelum menulis artefak pertama. Struktur, konvensi identitas/state/relasi, dan aturan kepemilikan dijelaskan di sana; rincian operasional ada di [docs/operating/protocol.md](docs/operating/protocol.md).

## Rantai Alur Kerja

Setiap pengiriman nilai mengikuti rantai pemesanan (Ordering):

```
requirement → scope → capability → plan → container → work item → session → review
```

| Langkah | Token | Folder |
|---|---|---|
| requirement — kebutuhan | `req-` | `docs/requirements/` |
| scope — konteks berfase | `scp-` | `docs/planning/` |
| capability — kapabilitas | `epc-` | `docs/planning/` |
| plan — rencana yang akan dikunci | `plan-` | `docs/planning/` |
| container — eksekusi (tepat satu aktif) | `ctr-` | `docs/operating/containers/` |
| work item — unit kerja | `sto-`/`ts-`/`bug-`/`td-`/`ch-`/`spk-` | `docs/operating/work-items/` |
| session — catatan eksekusi | `ses-` | `docs/operating/sessions/` |
| review — verifikasi | `rvw-` | `docs/quality/` |

## Aturan Kontribusi Singkat

| Folder | Pemilik konten | Catatan |
|---|---|---|
| `docs/intent/` … `docs/vocabulary/` | per peran pada README tiap folder | artefak pengetahuan (Knowledge Layer); `dimension` wajib sama dengan folder |
| `docs/operating/` | pemilik state per tipe (single-writer) | kontainer/tiket adalah proyeksi; jangan mengedit state di proyeksi |
| `docs/exchange/` | kontrak (contract) | tempat aturan validasi & transfer; bukan tempat artefak domain |

Aturan inti: identitas hidup di frontmatter (filename hanyalah proyeksi), state hanya ditulis oleh pemilik tunggalnya, proyeksi dibaca ulang bukan diedit, dan semua transisi state tercatat di `change-log`. Detail: [docs/operating/protocol.md](docs/operating/protocol.md).
