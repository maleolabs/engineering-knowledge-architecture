# docs/research/ — Dimensi Research

> Anchor EKA: Knowledge Layer — dimensi **research** (EKA 8), ekstensi token `fnd-` (EKA 14.1).

## Tujuan

Dimensi research mewadahi temuan investigasi dan riset: eksperimen, benchmark, penelusuran pustaka, dan temuan teknis yang belum (atau tidak) menjadi keputusan. Research adalah dimensi **pengetahuan** — bukan pengganti keputusan.

## Yang Ada di Sini

| Token | Tipe | Format nama |
|---|---|---|
| `fnd-` | Research Finding | `fnd-<id>.md` |

Ekstensi sesuai EKA 14.1: token `fnd-` ditambahkan pada token table standar untuk mewadahi temuan riset, dengan `dimension: research` dan folder rumah `research/`.

## State Vector

| Tipe | Domain state yang dimiliki |
|---|---|
| `fnd-` | `content-state`, `existence-state` |

Nilai `content-state`: `draft → review → approved → amended`. Nilai `existence-state`: `active → archived → retired`. Field lain = N/A.

## Struktur Konten yang Baik

Struktur wajib (keluarga dokumen pengetahuan, dengan ekstensi riset):

- `## Purpose` — pertanyaan riset yang dijawab.
- `## Content` — uraian temuan.
- `## Investigation Summary` — ringkasan metode dan bukti.
- `## Conclusion` — simpulan dan rekomendasi.

## Distilasi Wajib (EKA 11.4)

Temuan riset **tidak boleh berhenti** sebagai `fnd-`. Saat sebuah temuan memengaruhi arah proyek, temuan itu wajib terdistilasi ke dimensi decisions:

1. Simpulan yang diadopsi → `dec-`/`adr-` baru (atau amendemen keputusan yang relevan), dengan `derives-from: [fnd:<id>]`.
2. `fnd-` asli tetap utuh sebagai jejak bukti; baru boleh diarsipkan setelah distilasi selesai.

Temuan yang tidak diadopsi cukup diarsipkan dengan catatan alasan penolakan.

## Konvensi Nama

`fnd-<id>.md`, kebab-case, unik dalam (namespace, type). Tanpa akhiran `-v<nn>`. Contoh: `fnd-pendekatan-lock-plan.md`.

## Kepemilikan

| Peran | Tanggung jawab |
|---|---|
| Engineers | penulis `fnd-` (investigasi teknis) |
| Tech Lead | peninjau dan penyalur distilasi ke keputusan |
| Semua peran | boleh membuka riset baru |

## Terkait

- [decisions/](../decisions/) — **muara wajib** distilasi riset (EKA 11.4).
- [specifications/](../specifications/) — temuan yang diadopsi menjadi `spec-`.
- [operating/work-items/spikes/](../operating/work-items/spikes/) — spike menghasilkan bahan riset; `fnd-` dan `spk-` saling merujuk.
- [vocabulary/](../vocabulary/) — istilah baru dari riset didaftarkan di `gls-`.
