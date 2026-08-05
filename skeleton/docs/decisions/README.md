# docs/decisions/ — Dimensi Decisions

> Anchor EKA: Knowledge Layer — dimensi **decisions** (EKA 8).

## Tujuan

Dimensi decisions mewadahi keputusan yang sudah diambil beserta konteks, alternatif, dan konsekuensinya — agar keputusan tidak "hilang" saat orang yang mengambilnya pergi. Dimensi ini memiliki dua varian artefak dengan varian state yang berbeda.

## Yang Ada di Sini

| Token | Tipe | Format nama |
|---|---|---|
| `adr-` | Architecture Decision Record | `adr-<id>.md` |
| `dec-` | Decision Record (keputusan umum/operasional) | `dec-<id>.md` |

## State Vector

| Tipe | Domain state yang dimiliki | Nilai `content-state` |
|---|---|---|
| `adr-` | `content-state`, `existence-state` | `proposed → accepted → superseded` |
| `dec-` | `content-state`, `existence-state` | `draft → accepted → superseded` |

Nilai `existence-state`: `active → archived → retired`. Catatan: varian `adr-` tidak menggunakan nilai `draft`/`review`/`approved`/`amended`; status keputusan langsung diekspresikan sebagai `proposed`/`accepted`/`superseded`.

## Struktur Konten yang Baik

Struktur wajib (keluarga decision record):

- `## Context` — latar belakang dan masalah yang memicu keputusan.
- `## Decision` — keputusan yang diambil.
- `## Consequences` — dampak positif dan negatif.
- `## Alternatives Considered` — alternatif yang dievaluasi dan alasan penolakannya.
- `## References` (opsional) — rujukan tambahan.

## Supersesi

- **`adr-`: supersesi wajib.** Saat sebuah ADR tidak lagi berlaku, buat ADR baru dan set `supersedes: [adr:<id-lama>]` pada yang baru serta `content-state: superseded` pada yang lama. ADR yang disupersesi tanpa rujukan pengganti adalah pelanggaran validasi.
- **`dec-`: supersesi opsional.** Keputusan operasional dapat disupersesi, tetapi tidak diwajibkan; keputusan cukup dipindahkan ke `existence-state: archived` bila tidak lagi relevan.

## Konvensi Nama

`adr-<id>.md` dan `dec-<id>.md`, kebab-case, unik dalam (namespace, type). Tanpa akhiran `-v<nn>`. Contoh: `adr-serialisasi-identitas.md`, `dec-git-workflow.md`.

## Kepemilikan

| Peran | Tanggung jawab |
|---|---|
| Tech Lead | pemilik `adr-` (keputusan arsitektur) |
| Product Owner | pemilik `dec-` (keputusan produk/lingkup) |
| Engineers / DevOps | kontributor `dec-` teknis-operasional |

## Terkait

- [architecture/](../architecture/) — `adr-` menjelaskan mengapa `arc-` demikian.
- [records/](../records/) — catatan kronologis keputusan saat release.
- [sessions/](../sessions/) — hasil distilasi sesi wajib bermuara ke `dec-`/`adr-`.
- [research/](../research/) — temuan riset mewajibkan jalur distilasi ke sini (EKA 11.4).
