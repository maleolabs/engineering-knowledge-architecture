# Perubahan Breaking — Struktur Legacy → EKA v1.0

Ringkasan 14 perubahan breaking terhadap struktur lama. Setiap perubahan menyertakan: old → new, rasional arsitektur (anchor EKA), dampak, dan mitigasi.

| # | Old → New | Rasional (anchor EKA) | Dampak | Mitigasi |
|---|---|---|---|---|
| 1 | 16 folder campur dimensi → 12 folder dimensi + `operating/` + `exchange/` | Knowledge Taxonomy 1:1 (EKA 8); P1 | Seluruh path berubah | `migration-guide.md` |
| 2 | Ruang ID bersama `mvp-*` (4 tipe) → token tipe eksplisit `scp-`/`plan-`/`ctr-`/`tkt-` | Identity (EKA 6.4); P3 | Glob legacy `mvp-*` putus | `adr-001` + tabel token |
| 3 | `roadmap/` (misnomer) → `plan-` di `planning/` | Artifact Plan + Planning State (EKA 7.2/10) | Konsumen roadmap putus | `adr-002` + mapping status |
| 4 | `sprints/` + 3-way status → `ctr-` + single-writer proyeksi | P6, EKA 7.4 | Tooling tabel sprint harus membaca frontmatter work item | `adr-003` + regenerasi proyeksi |
| 5 | `tickets/` wave docs → `tkt-` proyeksi dengan State Vector kosong | EKA 7.4 | Command eksekusi membaca format `tkt-` baru | `adr-003` |
| 6 | `adr/` + `decisions/` terpisah → folder `decisions/` tunggal, dua tipe | Dimensi Decisions (EKA 8) | Path berubah | peta pemetaan |
| 7 | `specification-corpus/` (misnomer) → `vocabulary/` (`gls-`); spec asli → `specifications/` (`spec-`) | Vocabulary ≠ Specifications (EKA 8) | Glossary pindah; konten spec pindah | peta pemetaan |
| 8 | `planning/` catch-all + `work-items/planning/` → dibubarkan | EKA 14.2; P15 | Konten didistribusikan ke `trc-`/`plan-`/tipe work item yang tepat | `adr-005` |
| 9 | `operations/` campur → `operations/` (prosedur) + `standards/` (konvensi) | Operational vs Standards (EKA 8) | Exit code/output conventions pindah | peta pemetaan |
| 10 | Metadata table (Status/Author/Created/Updated/Version) → frontmatter state domains | Explicit State (P2); D2.8 | Tooling yang membaca `Status:` putus | `adr-002` + mapping nilai |
| 11 | Nilai status "In Progress"/"In Review" → `in-progress`/`in-review` | Value sets EKA 7.2 | String lama tidak valid | mapping nilai di `validation.md` |
| 12 | Folder phase → field `phase` pada `scp-`/`plan-` saja | Phase sebagai konteks (EKA 11.2) | Folder `mvp/` hilang | `adr-004` |
| 13 | `documentation-guide.md` tunggal → standard/serialization/protocol terpisah | EKA 1.3 | Pembaca guide membaca 3 dokumen | `reference-architecture.md` |
| 14 | `sp-` (spike) → `spk-` (token anti-prefix) | EKA 6.4; token bebas-ambiguitas | Glob `sp-*` putus | tabel token |

## Catatan tooling

- **Seluruh perubahan breaking bersifat disengaja (intentional)**: konsumen legacy **harus** putus agar tidak membaca Identity/state dari lokasi (EKA 6.4, P9). Kompatibilitas diam-diam akan mengulang pelanggaran yang sama: Identity disandikan di struktur.
- **Tooling baru membaca frontmatter saja**: Identity (`namespace`, `type`, `id`, `instance-version`, `revision`), State (`content-state`, `execution-state`, `planning-state`, `container-state`, `existence-state`), dan Relationship (`supersedes`, `amends`, `derives-from`, `depends-on`, `validates`) — semuanya dari frontmatter, tidak pernah dari path.
- **Filename untuk navigasi manusia + validasi konsistensi**: pola `<type-token>-<id>[-v<nn>]` memudahkan browsing dan menegakkan determinisme, tetapi bukan sumber kebenaran; validasi mekanis memeriksa konsistensi filename ↔ frontmatter.
- Referensi migrasi detail: [`migration-guide.md`](migration-guide.md); konvensi lengkap: [`reference-architecture.md`](reference-architecture.md).
