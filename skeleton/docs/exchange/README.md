# docs/exchange/ — Exchange Layer (EX)

> Anchor EKA: Exchange Layer — validasi + transfer. EX tidak memiliki konten dan tidak memiliki state.

## Peran

Exchange Layer mengelola **batas** serialisasi: memvalidasi konformitas artefak terhadap kontrak EKA dan mengatur impor/ekspor antar repositori. EX adalah lapisan yang "tidak memiliki keduanya" — tidak memiliki konten (milik Knowledge Layer) dan tidak memiliki state (milik Operating Layer).

## Yang Dimiliki EX (Kontrak)

| Kontrak | Dimana |
|---|---|
| Aturan validasi konformitas (9 aturan mekanis) | [validation.md](validation.md) |
| Konvensi impor/ekspor (round-trip, konflik identitas, skema) | [transfer.md](transfer.md) |
| Tabel token, format referensi, dan format identitas | diserialisasi di README proyek; EX menegakkannya saat validasi |

## Yang Tidak Dimiliki EX

- **Konten** — isi artefak adalah milik pemilik dimensi (Knowledge Layer). EX hanya memeriksa bentuknya.
- **State** — nilai state dan transisinya adalah milik pemilik state (Operating Layer). EX hanya memeriksa validitas nilainya, bukan mengubahnya.

## Alur Pemakaian

1. Penulis membuat/mengubah artefak.
2. Sebelum commit: jalankan checklist [validation.md](validation.md).
3. Untuk pertukaran lintas repositori: ikuti [transfer.md](transfer.md).
4. Pelanggaran = tolak commit sampai konform.

## Terkait

- [../operating/protocol.md](../operating/protocol.md) — aturan yang divalidasi EX lahir dari OS protocol.
- [../README.md](../README.md) — konvensi identitas, state, phase, relasi, klasifikasi, proyeksi.
