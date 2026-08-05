# Validation — Checklist Konformitas (9 Aturan)

> Anchor EKA: Exchange Layer — validasi konformitas. Dokumen konvensi, bukan artefak.
> Standar: EKA v1.0, tanggal 2026-08-05.

Checklist mekanis berikut dijalankan **sebelum commit** setiap artefak baru atau perubahan state. Semua aturan bersifat mekanis (dapat diotomasi). Format: 1 = lolos, 0 = gagal (blocking), W = peringatan.

## Aturan 1 — Keunikan Identitas

Tidak boleh ada duplikat `(namespace, type, id, instance-version)` di seluruh repositori.
- [ ] Tidak ada artefak lain dengan kombinasi yang sama.

## Aturan 2 — Konsistensi Filename

- [ ] Token pada filename == nilai `type` di frontmatter.
- [ ] Akhiran versi pada filename (jika ada) == nilai `instance-version`.
- [ ] Akhiran `-v<nn>` **hanya** diizinkan untuk `scp-`/`plan-`; tipe lain dilarang membawanya.
- [ ] `scp-`/`plan-` **wajib** memuat akhiran `-v<nn>` (termasuk v1).

## Aturan 3 — Validitas Nilai State

Setiap nilai field state harus anggota himpunan nilai domainnya:

| Domain | Himpunan nilai |
|---|---|
| Content State (umum) | draft, review, approved, amended |
| Content State (ADR) | proposed, accepted, superseded |
| Content State (decision) | draft, accepted, superseded |
| Execution State | planned, todo, in-progress, in-review, done |
| Planning State | draft, approved, immutable |
| Container State | active, completed |
| Existence State | active, archived, retired |
| Phase (konteks, scp-/plan-) | discovery, mvp, milestone, release, growth, maturity, sunset |

- [ ] Semua nilai field state ∈ himpunan nilainya.
- [ ] Nilai `phase` ∈ himpunan phase (hanya pada `scp-`/`plan-`).

## Aturan 4 — Kepatuhan Owned-Set

Field state yang ada pada sebuah file harus **persis** sama dengan himpunan owned domain untuk tipe-nya (absen = N/A):

| Tipe | Owned set |
|---|---|
| vis-, str-, req-, scp-, epc-, trc-, arc-, adr-, dec-, spec-, std-, run-, rvw-, rel-, gls-, fnd- | content-state, existence-state |
| plan- | content-state, planning-state, existence-state |
| sto-, ts-, bug-, td-, ch-, spk- | execution-state, existence-state |
| ctr- | container-state, existence-state |
| ses- | existence-state |
| tkt- | (tidak ada — state vector kosong) |

- [ ] Tidak ada field state milik tipe lain pada file (mis. `container-state` pada work item = pelanggaran).
- [ ] `tkt-` tidak memuat field state apa pun.

## Aturan 5 — Integritas Referensial

Semua referensi (`amends`, `supersedes`, `derives-from`, `depends-on`, `validates`) harus dapat di-resolve ke artefak yang ada.

- [ ] Setiap referensi menunjuk artefak yang ada (format `<type>:<id>[:<instance-version>]`, lintas-namespace `<ns>/<type>:<id>`).
- [ ] Referensi hanya ditulis pada artefak perujuk (bukan dua arah).
- [ ] `content-state: draft` → referensi yang belum ter-resolve diizinkan (**W** peringatan).
- [ ] `content-state` non-draft → referensi yang belum ter-resolve = **0** (error).

## Aturan 6 — Dimensi == Folder

- [ ] Artefak pengetahuan: nilai `dimension` == folder rumahnya (mis. file di `docs/specifications/` harus `dimension: specifications`).
- [ ] Artefak operasional (work items) boleh memuat `dimension` informasional — tidak dievaluasi.
- [ ] `ctr-`, `tkt-`, `ses-` **tidak boleh** memuat `dimension`.

## Aturan 7 — Konsistensi Change-Log

- [ ] Setiap domain yang dimiliki memiliki entri `change-log` terakhir == nilai field saat ini (per domain: content-state, execution-state, planning-state, container-state, existence-state, dan `phase` pada scp-/plan-).
- [ ] Tidak ada transisi tanpa entri `change-log`.

## Aturan 8 — Single-Writer & Proyeksi

- [ ] Proyeksi (`tkt-`, tabel work item dalam `ctr-`) tidak memuat field owned state milik artefak lain.
- [ ] Tabel work item pada kontainer sesuai dengan owner state work item yang bersangkutan — **divalidasi pada saat baca** (jika tidak cocok: peringatan W; sumber kebenaran tetap owner state).
- [ ] Header proyeksi (`> Generated — State Projection. Do NOT edit state here; refresh on read.`) ada pada file proyeksi.

## Aturan 9 — Well-Formedness

Struktur konten wajib ada sesuai keluarga tipe:

| Keluarga | Bagian wajib |
|---|---|
| Planning artifact (scp-, epc-, plan-, trc-) | `## Objective`, `## Scope`, `## Out of Scope` |
| Work item (sto-, ts-, ch-) | `## Description`, `## Acceptance Criteria` |
| Bug (bug-) | `## Description`, `## Impact` |
| Tech Debt (td-) | `## Description`, `## Acceptance Criteria`, `## Debt Rationale` |
| Spike (spk-) | `## Description`, `## Investigation Notes`, `## Conclusion` (memuat tautan distilasi) |
| Decision record (adr-, dec-) | `## Context`, `## Decision`, `## Consequences`, `## Alternatives Considered` |
| Knowledge doc (vis-, str-, req-, arc-, spec-, std-, run-, rel-, gls-, fnd-) | `## Purpose`, `## Content` |
| Review (rvw-) | `## Purpose`, `## Content`, `## Findings`, `## Action Items` |
| Research Finding (fnd-) | `## Purpose`, `## Content`, `## Investigation Summary`, `## Conclusion` |
| Container (ctr-) | `## Objective`, `## Work Items`, `## Change Log` |
| Ticket (tkt-) | `## Commands`, `## Projected Status` |
| Session (ses-) | `## Context`, `## Notes`, `## Verification` |

- [ ] Semua bagian wajib ada untuk tipe masing-masing.
- [ ] Supersesi ADR: `adr-` dengan `content-state: superseded` wajib dirujuk penggantinya (0 jika tidak).

## Hasil

- [ ] Semua aturan 1–9 lolos → commit diizinkan.
- [ ] Ada nilai **0** → perbaiki dahulu, jangan commit.
- [ ] Hanya **W** → commit diizinkan dengan catatan peringatan.
