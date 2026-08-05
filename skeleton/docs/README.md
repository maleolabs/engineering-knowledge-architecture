# docs/ — Sumber Kebenaran Serialisasi EKA

> Anchor EKA: serialisasi keseluruhan — Knowledge Layer (konten), Operating Layer (state), Exchange Layer (validasi & transfer).
> Standar: EKA v1.0, tanggal 2026-08-05.

## Status Folder Ini

Folder ini adalah **serialisasi EKA v1.0** (proyeksi EKA ke Git + Markdown), **bukan arsitektur EKA itu sendiri**. Arsitektur — tiga layer (Knowledge/Operating/Exchange), lima state domain, dan protokol — adalah konsep; folder ini hanyalah representasi file-nya. Artinya: mengubah folder tidak mengubah EKA, dan EKA dapat diserialisasikan ke media lain tanpa kehilangan makna.

## Peta Navigasi (15 Entri)

| Entri | Anchor EKA | Isi |
|---|---|---|
| [README.md](.) (ini) | serialisasi | sumber kebenaran + ringkasan konvensi |
| [intent/](intent/) | dimensi intent | `vis-` Vision/Manifesto, `str-` Strategy |
| [requirements/](requirements/) | dimensi requirements | `req-` Requirement |
| [architecture/](architecture/) | dimensi architecture | `arc-` Architecture Description |
| [decisions/](decisions/) | dimensi decisions | `adr-` ADR, `dec-` Decision Record |
| [specifications/](specifications/) | dimensi specifications | `spec-` Specification |
| [standards/](standards/) | dimensi standards | `std-` Standard/Guideline |
| [operations/](operations/) | dimensi operations | `run-` Runbook |
| [quality/](quality/) | dimensi quality | `rvw-` Review |
| [planning/](planning/) | dimensi planning | `scp-`, `epc-`, `plan-`, `trc-` |
| [records/](records/) | dimensi records | `rel-` Release Record |
| [research/](research/) | dimensi research | `fnd-` Research Finding (EKA 14.1) |
| [vocabulary/](vocabulary/) | dimensi vocabulary | `gls-` Glossary/Term |
| [operating/](operating/) | Operating Layer | state, protocol, work items, kontainer, sesi, proyeksi |
| [exchange/](exchange/) | Exchange Layer | validasi, transfer |

## Ringkasan Konvensi Serialisasi

### Identitas
- Identitas = `(namespace, type, id, instance-version)`; **hidup di frontmatter**, filename hanyalah proyeksi.
- Field: `namespace` (default: nama produk), `type` (token, harus cocok dengan filename), `id` (kebab-case, unik dalam pasangan (namespace, type)), `instance-version` (int, default 1; wajib untuk `scp-`/`plan-`), `revision` (int, default 1; hanya untuk edit konten).
- Pola filename: `<type-token>-<id>.md`; untuk tipe berversi (`scp-`, `plan-`): `<type-token>-<id>-v<instance-version>.md` (selalu, termasuk v1). Akhiran `-v<nn>` dilarang untuk tipe lain. Detail: [planning/README.md](planning/README.md).

### State
- Lima domain state yang dimiliki (OWNED): Content State, Execution State, Planning State, Container State, Existence State. Absennya sebuah field = tidak berlaku (N/A) untuk tipe itu.
- Setiap field state hanya ditulis oleh **satu pemilik** (single-writer, P6); proyeksi tidak pernah menulis state.
- Transisi **forward-only**; setiap transisi tercatat di `change-log`.
- Detail nilai dan aturan transisi: [operating/protocol.md](operating/protocol.md).

### Phase
- `phase` adalah **atribut konteks** pada artefak `scp-`/`plan-`, **bukan folder**. Nilai: `discovery | mvp | milestone | release | growth | maturity | sunset`.
- Perubahan phase dicatat di `change-log` dengan `domain: phase`.

### Relasi
- Field: `amends`, `supersedes`, `derives-from`, `depends-on`, `validates` — daftar referensi berformat `<type-token>:<id>[:<instance-version>]` (lintas-namespace: `<namespace>/<type-token>:<id>`).
- Hanya ditulis pada artefak yang mereferensikan; referensi harus dapat di-resolve (lihat validasi).

### Klasifikasi
- `dimension` (primer) + `dimensions-secondary` (daftar) — untuk artefak pengetahuan; `dimension` wajib sama dengan folder rumahnya. Artefak operasional (work items) memakai `dimension` informasional saja; `ctr-`/`tkt-`/`ses-` tidak membawa `dimension`.

### Proyeksi
- Tabel kontainer dan tiket membawa header `> Generated — State Projection. Do NOT edit state here; refresh on read.`
- Proyeksi di-refresh saat dibaca, bukan diedit manual.

## Artefak vs Dokumen Konvensi

- **Artefak:** file yang frontmatter-nya berisi `type` **dan** `id`. Contoh: `req-login.md`, `ctr-gelombang-1.md`, `plan-rilis-1-v1.md`.
- **Dokumen konvensi:** file tanpa keduanya — semua README, `operating/protocol.md`, `exchange/validation.md`, `exchange/transfer.md`. Dokumen konvensi menjelaskan aturan dan tidak membawa state.

## Validasi

Sebelum commit artefak baru atau perubahan state, jalankan checklist mekanis di [exchange/validation.md](exchange/validation.md). Untuk impor/ekspor antar repositori: [exchange/transfer.md](exchange/transfer.md).
