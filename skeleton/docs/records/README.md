# docs/records/ — Dimensi Records

> Anchor EKA: Knowledge Layer — dimensi **records** (EKA 8).

## Tujuan

Dimensi records mewadahi catatan kronologis yang tidak berubah: release record, catatan peristiwa, dan jejak audit proyek. Berbeda dengan dokumen kerja yang terus berkembang, catatan di sini bersifat faktu al — ia merekam apa yang terjadi.

## Yang Ada di Sini

| Token | Tipe | Format nama |
|---|---|---|
| `rel-` | Release Record | `rel-<id>.md` |

## State Vector

| Tipe | Domain state yang dimiliki |
|---|---|
| `rel-` | `content-state`, `existence-state` |

Nilai `content-state`: `draft → review → approved → amended`. Nilai `existence-state`: `active → archived → retired`. Field lain = N/A. Setelah release, `rel-` praktis `approved` dan tidak diedit lagi (perubahan fakta = instance baru dengan `amends`).

## Struktur Konten yang Baik

Struktur wajib (keluarga dokumen pengetahuan):

- `## Purpose` — release/peristiwa apa yang dicatat.
- `## Content` — ringkasan eksekusi dan rilis.

## Release Record = Agregat Eksekusi + Release Gates

`rel-` adalah **agregat** dari hasil eksekusi (work item yang selesai, sesi yang terjadi) dan **release gates** (gate persetujuan dan kesiapan yang lulus sebelum rilis). Keduanya dirujuk dengan relasi (`derives-from`/`validates`), bukan dikutip ulang — agar catatan tetap ringkas dan jejaknya tertelusur.

## Konvensi Nama

`rel-<id>.md`, kebab-case, unik dalam (namespace, type). Tanpa akhiran `-v<nn>`. Contoh: `rel-2026-08-mvp.md`.

## Kepemilikan

| Peran | Tanggung jawab |
|---|---|
| DevOps | pemilik `rel-` (eksekusi release) |
| Tech Lead | penandatangan release gate teknis |
| Product Owner | penandatangan release gate produk |

## Terkait

- [operations/](../operations/) — prosedur release yang dijalankan (`run-`).
- [quality/](../quality/) — release gates diverifikasi oleh `rvw-`.
- [decisions/](../decisions/) — keputusan yang dipicu release tercatat sebagai `dec-`.
- [operating/work-items/](../operating/work-items/) — hasil eksekusi yang diagregasi `rel-`.
