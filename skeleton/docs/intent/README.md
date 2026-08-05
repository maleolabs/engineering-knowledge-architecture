# docs/intent/ — Dimensi Intent

> Anchor EKA: Knowledge Layer — dimensi **intent** (EKA 8).

## Tujuan

Dimensi intent mewadahi arahan dan alasan keberadaan produk/proyek: visi, manifesto, dan strategi. Dokumen di sini menjawab pertanyaan "mengapa proyek ini ada" dan "ke arah mana ia bergerak", dan menjadi rujukan pembenaran (justification) bagi dimensi lain.

## Yang Ada di Sini

| Token | Tipe | Format nama |
|---|---|---|
| `vis-` | Vision/Manifesto | `vis-<id>.md` |
| `str-` | Strategy | `str-<id>.md` |

## State Vector

| Tipe | Domain state yang dimiliki |
|---|---|
| `vis-` | `content-state`, `existence-state` |
| `str-` | `content-state`, `existence-state` |

Nilai `content-state`: `draft → review → approved → amended`. Nilai `existence-state`: `active → archived → retired`. Field lain yang tidak terdaftar = tidak berlaku (N/A). Setiap transisi dicatat di `change-log` oleh pemilik tunggal.

## Struktur Konten yang Baik

Struktur wajib (keluarga dokumen pengetahuan):

- `## Purpose` — tujuan dokumen ini.
- `## Content` — isi visi/manifesto atau strategi.

## Konvensi Nama

`vis-<id>.md` dan `str-<id>.md`, kebab-case, unik dalam (namespace, type). Tanpa akhiran `-v<nn>` (tipe tidak berversi). Contoh: `vis-inti-produk.md`, `str-masuk-pasar-2026.md`.

## Kepemilikan

| Peran | Tanggung jawab |
|---|---|
| Product Owner | pemilik tunggal konten intent |
| Tech Lead | kontributor untuk `str-` teknis |
| Semua peran | boleh mengusulkan perubahan; perubahan konten yang disetujui menjadi instance baru dengan `amends` |

## Terkait

- [requirements/](../requirements/) — intent diturunkan menjadi kebutuhan (`req-`).
- [vocabulary/](../vocabulary/) — istilah kunci intent harus terdefinisi di `gls-`.
- [decisions/](../decisions/) — keputusan strategis (`dec-`) merujuk balik ke intent.
- [planning/](../planning/) — `scp-` menjabarkan konteks berfase dari strategi.
