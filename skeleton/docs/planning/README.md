# docs/planning/ — Dimensi Planning

> Anchor EKA: Knowledge Layer — dimensi **planning** (EKA 8) + state domain **Planning State**.

## Tujuan

Dimensi planning mewadahi artefak perencanaan: definisi scope, kapabilitas, rencana eksekusi, dan artefak relasi. Inilah satu-satunya dimensi tempat `phase` (discovery | mvp | milestone | release | growth | maturity | sunset) dipakai — sebagai **atribut konteks** pada `scp-`/`plan-`, bukan folder.

## Yang Ada di Sini

| Token | Tipe | Format nama | Berversi |
|---|---|---|---|
| `scp-` | Scope Definition | `scp-<id>-v<instance-version>.md` | ya (fase) |
| `epc-` | Epic | `epc-<id>.md` | tidak |
| `plan-` | Plan | `plan-<id>-v<instance-version>.md` | ya (fase, Planning State) |
| `trc-` | Traceability/Relationship artifact | `trc-<id>.md` | tidak |

## State Vector

| Tipe | Domain state yang dimiliki |
|---|---|
| `scp-`, `epc-`, `trc-` | `content-state`, `existence-state` |
| `plan-` | `content-state`, `planning-state`, `existence-state` |

Nilai `content-state`: `draft → review → approved → amended`. Nilai `planning-state`: `draft → approved → immutable`. Nilai `existence-state`: `active → archived → retired`. Field lain = N/A.

## Struktur Konten yang Baik

Struktur wajib (keluarga planning artifact):

- `## Objective` — tujuan artefak ini.
- `## Scope` — cakupan yang termasuk.
- `## Out of Scope` — cakupan yang sengaja dikecualikan.

## Konvensi Nama

Berversi (selalu, termasuk v1): `scp-<id>-v<instance-version>.md` dan `plan-<id>-v<instance-version>.md`. Tidak berversi: `epc-<id>.md`, `trc-<id>.md`. `instance-version` wajib di frontmatter untuk `scp-`/`plan-`.

## Catatan Khusus per Tipe

### `scp-` — Scope Definition
Membawa `phase` sebagai atribut konteks. Perubahan phase dicatat di `change-log` dengan `domain: phase`. Scope yang disetujui menjadi dasar kontainer eksekusi.

### `plan-` — Plan
Membawa `phase` dan **Planning State** (`draft → approved → immutable`). **Lock-atomic-with-generation:** saat `plan-` menjadi `immutable`, perubahan apa pun — termasuk perbaikan — tidak boleh mengedit instance itu; buat `instance-version` baru (`plan-<id>-v<nn+1>.md`). Transisi ke `immutable` terjadi atomik dengan pembuatan kontainer (lihat [operating/containers/](../operating/containers/) dan [operating/protocol.md](../operating/protocol.md)).

### `epc-` — Epic
Kapabilitas yang mewujudkan scope. Tidak membawa `phase`; rujuk `scp-` induknya dengan `derives-from` bila perlu.

### `trc-` — Traceability/Relationship artifact
Artefak yang **hanya membawa relasi** (mis. matriks kebutuhan→spesifikasi→work item, peta dimensi). Konten utamanya adalah daftar referensi yang harus dapat di-resolve; ia tidak menggantikan relasi yang ditulis pada artefak perujuk.

## Kepemilikan

| Peran | Tanggung jawab |
|---|---|
| Product Owner | pemilik `scp-`, `epc-`, `plan-` (lingkup & prioritas) |
| Tech Lead | peninjau kelayakan dan pemilik `plan-` teknis |
| Engineers | kontributor estimasi dan detail plan |

## Terkait

- [requirements/](../requirements/) — scope menyeleksi `req-`.
- [intent/](../intent/) — `scp-` menjabarkan strategi ke konteks berfase.
- [operating/containers/](../operating/containers/) — `plan-` terkunci atomik dengan lahirnya `ctr-`.
- [operating/projections/](../operating/projections/) — tiket mewakili work item dalam kontainer.
- [records/](../records/) — `rel-` mencatat hasil eksekusi plan.
