# Panduan Migrasi — Struktur Legacy → EKA v1.0

Panduan migrasi dari struktur dokumentasi lama menuju serialisasi EKA v1.0. Seluruh perubahan bersifat **breaking by design**: konsumen lama wajib berhenti bekerja agar tidak membaca Identity/State dari lokasi (EKA 6.4, P9).

## Bagian A — Peta pemetaan lengkap (legacy → baru)

| Elemen legacy | Home baru | Aksi | Kompatibilitas | Rasional |
|---|---|---|---|---|
| `docs/README.md` | `skeleton/docs/README.md` | transform | breaking | Konvensi serialisasi baru |
| `docs/manifesto/` | `docs/intent/` sebagai `vis-` | move + rename | breaking | Dimensi Product Intent |
| `docs/prd/` | `docs/requirements/` sebagai `req-`; amendmen sebagai instance dengan `amends` | move + rename + transform | breaking | Dimensi Requirements |
| `docs/mvp/` | `docs/planning/` sebagai `scp-…-v<n>` dengan `phase: mvp` | transform | breaking | Phase sebagai konteks (EKA 11.2); Identity decoupled dari phase (P3); kolisi `mvp-nnn` diselesaikan (EKA 6.4) |
| `docs/epics/` | `docs/planning/` sebagai `epc-` | move + rename | breaking | Dimensi Planning Knowledge |
| `docs/architecture/` | `docs/architecture/` sebagai `arc-` | rename | breaking | Type token |
| `docs/adr/` | `docs/decisions/` sebagai `adr-` | move + rename | breaking | Dimensi Decisions tunggal |
| `docs/decisions/` | `docs/decisions/` sebagai `dec-` | move + rename | breaking | Dimensi Decisions tunggal |
| `docs/roadmap/` | `docs/planning/` sebagai `plan-…-v<n>` | move + rename + transform | breaking | Misnomer dikoreksi; Planning State dipertahankan: approved→approved, immutable→immutable |
| `docs/sprints/` | `docs/operating/containers/` sebagai `ctr-` | move + rename + transform | breaking | Execution Container (EKA 10); Container State; tabel = proyeksi (EKA 7.4) |
| `docs/tickets/` | `docs/operating/projections/` sebagai `tkt-` | move + rename + transform | breaking | State Vector kosong, `derives-from` |
| `docs/work-items/{stories,technical-stories,bugs,tech-debt,chores,spikes}/` | `docs/operating/work-items/<subtype>/` sebagai `sto-`, `ts-`, `bug-`, `td-`, `ch-`, `spk-` | move + rename | breaking | Single-writer Execution State |
| `docs/work-items/planning/` | deprecated | — | breaking | Catch-all dibubarkan; konten → `planning/`, pekerjaan planning → tipe work item yang tepat |
| `docs/reviews/` | `docs/quality/` sebagai `rvw-` dengan `validates` | move + rename | breaking | Dimensi Governance & Quality |
| `docs/sessions/` | `docs/operating/sessions/` sebagai `ses-` | move + rename | breaking | Existence State; Distillation wajib (EKA 11.4) |
| `docs/operations/` | split: `docs/operations/` sebagai `run-` (prosedur) + `docs/standards/` sebagai `std-` (konvensi) | — | breaking | Operational vs Standards (EKA 8) |
| `docs/planning/` | `docs/planning/` sebagai `trc-`/`plan-` | move + rename | breaking | Catch-all dibubarkan |
| `docs/specification-corpus/` | `docs/vocabulary/` sebagai `gls-`; spec asli → `docs/specifications/` sebagai `spec-` | move + rename | breaking | Misnomer dikoreksi; Vocabulary ≠ Specifications (EKA 8) |
| `documentation-guide.md` | split: `reference-architecture.md` + `skeleton/docs/README.md` + `operating/protocol.md` | — | breaking | Standard ≠ serialisasi (EKA 1.3) |
| `README.md` | root `README.md` baru | — | breaking | Identitas repositori baru |
| 3-way status sync | single-writer frontmatter + proyeksi | — | breaking | P6, 7.4 |
| metadata table (Status/Author/…) | frontmatter | — | breaking | D2.8: status → state domains; version split menjadi instance-version + revision |

## Bagian B — Strategi migrasi langkah-demi-langkah

1. **Snapshot** — Bekukan sprint aktif, catat status final seluruh artifact (work item, container, plan), commit baseline sebelum migrasi agar semua kondisi dapat dipulihkan dari git history.
2. **Buat struktur baru** — Salin `skeleton/docs/` ke proyek; set `namespace` pada seluruh artifact sesuai proyek; baca `skeleton/docs/README.md` sebagai sumber kebenaran struktur.
3. **Migrasi artifact knowledge dalam urutan dependensi**:
   - `intent` (visi/manifesto → `vis-`, strategi → `str-`);
   - `requirements` (amendmen → instance baru + relasi `amends`);
   - `decisions` (`adr/` + `decisions/` digabung ke `decisions/`; mapping status: Draft→draft, Review→review, Approved→approved, Accepted→accepted, Superseded→superseded, Amended→amended);
   - `architecture` (`arc-`);
   - `specifications` (baru — `spec-`);
   - `standards` (diekstrak dari operations — `std-`);
   - `operations` (hanya prosedur — `run-`);
   - `quality` (`rvw-` dengan `validates`);
   - `vocabulary` (`gls-`);
   - `planning` (`scp-` dengan `phase`, `epc-`, `plan-` dengan `planning-state`, `trc-`).
4. **Migrasi artifact operating**:
   - **Work items dulu** — content + `execution-state` dari status legacy; `change-log` dibangun ulang dari Change Log legacy; file work item kini menjadi single-writer;
   - **Containers** (`ctr-`) — `container-state: active/completed` dari state sprint; tabel container = proyeksi yang di-regenerate;
   - **Tickets** (`tkt-`) — + `derives-from` + status diproyeksikan;
   - **Sessions** (`ses-`) — dengan Distillation wajib sebelum Archived (EKA 11.4).
5. **Bangun ulang relationship** — Scan konten legacy; encode `amends`/`supersedes`/`derives-from`/`depends-on`/`validates` di frontmatter; referensi selalu by Identity.
6. **Validasi** — Jalankan checklist `exchange/validation.md` (9 aturan, ADR-006); perbaiki seluruh temuan sebelum lanjut.
7. **Arsipkan elemen legacy** — Pulihkan dari git history bila diperlukan; **jangan** migrasikan status legacy yang terduplikasi sebagai otoritatif — single-writer (P6).
8. **Perbarui tooling** — Konsumen status/lokasi legacy bermigrasi ke Identity/State frontmatter; glob legacy (`mvp-*`, `sp-*`, dst.) sengaja diputus (lihat `breaking-changes.md`).
