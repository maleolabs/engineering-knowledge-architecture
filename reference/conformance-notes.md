# Catatan Implementasi Konformansi — Validator EKA (`conformance/`)

> Dokumen konvensi, bukan artefak. Meta-dokumentasi zona `reference/`.
> Terkait: [`cli.md`](cli.md) (dokumentasi CLI), [`../skeleton/docs/exchange/validation.md`](../skeleton/docs/exchange/validation.md) (9 aturan konformitas), [`../standard/eka-specification-v1.0.md`](../standard/eka-specification-v1.0.md) (standard kanonik).

> **Konsolidasi traceability.** Tabel traceability aturan (spesifikasi ↔ implementasi) pada dokumen ini telah dikonsolidasi ke [`conformance-traceability-matrix.md`](conformance-traceability-matrix.md) — single source of truth cakupan konformansi (REQ→Spec→Rule→Impl→Test). Dokumen ini kini hanya memegang **keputusan interpretasi (29 item)** dan **gap yang diketahui**; perbarui matriks tersebut, bukan tabel di sini, saat cakupan berubah.

## Tujuan

Dokumen ini menjelaskan **bagaimana CLI `eka` menjalankan spesifikasi EKA secara mekanis**: klasifikasi artifact vs dokumen konvensi, aturan R0–R9, dan — yang terpenting — setiap **keputusan interpretasi** yang diambil ketika teks aturan tidak cukup presisi untuk dieksekusi mesin. Tujuannya adalah traceability aturan-demi-aturan: dari teks standard → teks aturan `validation.md` → perilaku validator → lokasi implementasi Go.

**Kebijakan interpretasi:** jika spesifikasi ambigu, keputusan **didokumentasikan sebelum implementasi**; tidak ada perilaku yang diciptakan tanpa dasar. Setiap keputusan di bawah merujuk dasar spesifikasinya (teks aturan, ADR, README dimensi, atau kenyataan repositori sendiri). Keputusan yang sama juga dicatat sebagai komentar `Interpretation (documented)` di sumber kode terkait.

## Keputusan interpretasi (29 item)

Nomor #1–#28 adalah keputusan yang diambil selama implementasi; #29 adalah gap yang diketahui dan sengaja **tidak** divalidasi (dokumentasi lengkap di bagian [Gap yang diketahui](#gap-yang-diketahui)).

| # | Aturan | Ambiguitas | Keputusan | Rasional |
|---|---|---|---|---|
| 1 | Artifact rule | Frontmatter `date: 2026-08-05` tanpa tanda kutip ter-parse sebagai timestamp oleh yaml.v3, bukan string | Field tanggal menerima string **atau** node timestamp; dinormalisasi ke `YYYY-MM-DD` | ADR kanonik repositori menulis tanggal tanpa tanda kutip; menolaknya akan mematahkan artifact repositori sendiri |
| 2 | R4 | "persis sama (absen = N/A)" — apakah field owned yang **tidak hadir** boleh? | Field owned yang absen = **ERROR** (`missing owned state field`) | "persis sama" + ADR-002 "field hanya hadir untuk domain yang dimiliki" + seluruh artifact nyata membawa semua field owned-nya; R7 juga membutuhkan nilai saat ini untuk dibandingkan |
| 3 | R4 | Field non-owned yang absen — bagaimana dievaluasi? | Absen field **non-owned** = N/A, tidak dievaluasi (bukan error) | Semantik eksplisit "absen = N/A" pada teks aturan |
| 4 | R7 | "Tidak ada transisi tanpa entri change-log" (bullet 2) tidak dapat diverifikasi dari snapshot | Di-enforce: keberadaan entri per domain owned, `last-entry == nilai saat ini`, well-formedness entri, legalitas transisi. **Tidak** di-enforce: kontiguitas rantai `entry[i].from == entry[i-1].to` | Snapshot statis tidak mengobservasi state antara; ADR kanonik memulai riwayat content-state di tengah rantai (mis. `proposed -> accepted` tanpa entri awal) |
| 5 | R7 | Makna `from: "-"` dan `to: "-"` tidak didefinisikan teks aturan | `from: "-"` = penanda state awal, legal untuk semua domain; `to: "-"` = **invalid** | Konvensi yang dipakai seluruh ADR repositori; `to` harus nilai nyata domain |
| 6 | R7 | Seberapa ketat "forward-only"? | Execution State: **strictly adjacent** (tidak boleh skip). Domain lain: forward-only (indeks naik, adjacency tidak wajib) | EKA §7.2 "strictly sequential" berlaku hanya untuk Execution State; domain lain cukup tidak regresi (P7) |
| 7 | R7 | Entri change-log untuk `phase` — aturan urutan? | Entri phase: validasi value-set + well-formedness **saja**; tanpa constraint urutan | Phase adalah context attribute (EKA 11.2), bukan state domain |
| 8 | R7 | Artifact tanpa change-log sama sekali — error? | Error hanya jika artifact memiliki ≥1 domain owned atau phase; `tkt-` (state vector kosong) boleh tanpa change-log | Tidak ada state yang perlu dicatat transisinya |
| 9 | R5 | Referensi bare-id (`001-identity-serialization`) tidak cocok dengan grammar `<type>:<id>` | Bare-id diterima = line reference yang di-resolve dalam namespace + type perujuk | 6 dari 7 ADR nyata memakai bentuk de-facto ini; menolaknya mematahkan repositori sendiri |
| 10 | R5 | Referensi malformed vs unresolved — severity? | Malformed (tidak ter-parse) = **error selalu**; unresolved = warning bila `content-state: draft`, error selain itu (termasuk artifact tanpa content-state) | Teks aturan hanya membahas unresolved; malformed adalah pelanggaran struktural yang tidak pernah dapat di-resolve |
| 11 | R5 | "Referensi hanya ditulis pada artefak perujuk (bukan dua arah)" — di-enforce? | **Tidak** di-enforce secara mekanis | Bukan pemeriksaan mekanis yang jelas (membutuhkan semantik arah per relationship); tugas implementasi melarang menambah aturan |
| 12 | R5 | Format referensi lintas-namespace | `<ns>/<type>:<id>`; versi opsional `<ns>/<type>:<id>:<ver>` | Perluasan grammar teks aturan yang konsisten dengan `<type>:<id>[:<instance-version>]` |
| 13 | R6 | "Folder rumahnya" untuk file di luar `docs/` (mis. `reference/decisions/`) | Folder rumah = ancestor terdekat yang namanya ∈ 12 dimensi; tidak ditemukan → error | Membuat `reference/decisions/` dan `docs/decisions/` sama-sama terpetakan ke `decisions`; tetap bekerja untuk artifact bersarang di bawah folder dimensi |
| 14 | R6 | Artifact pengetahuan tanpa field `dimension` | = **ERROR** | Klasifikasi adalah properti wajib artifact (P15; reference-architecture.md §2.5) |
| 15 | R6 | `dimensions-secondary` — divalidasi sejauh apa? | Validasi token saja (harus ∈ 12 dimensi); dilarang pada `ctr-`/`tkt-`/`ses-` | Properti klasifikasi sekunder; tidak ada aturan teks yang lebih dalam |
| 16 | R2 | Jumlah digit akhiran `-v<nn>` | 1+ digit diterima (`-v1` dan `-v01` sama-sama valid) | "termasuk v1" (docs/README.md, validation.md Aturan 2) |
| 17 | R2 | Bagian id filename vs id frontmatter | **Tidak** dicocokkan (gap terdokumentasi) | Aturan hanya menuntut konsistensi token + versi; filename adalah proyeksi (ADR-001), bukan Identity |
| 18 | R8 | Posisi header proyeksi dalam file | Kehadiran di mana pun dalam file (pencocokan eksak, whitespace trailing diabaikan); posisi bebas | Teks aturan hanya menuntut header "ada pada file proyeksi" |
| 19 | R8 | `tkt-` harus menunjuk apa via `derives-from`? | ≥1 referensi `derives-from` yang resolve ke artifact `ctr-`; referensi ke work item tidak diwajibkan | ADR-003 §3; README proyeksi menampilkan keduanya, teks aturan hanya kontainer |
| 20 | R8 | Format tabel `## Work Items` tidak dicontohkan di containers/README.md | Didefinisikan: tabel GFM (baris header + baris separator); kolom pertama = id work item atau `<type>:<id>`; kolom execution-state dikenali dari varian header (execution-state / execution state / execution_state / status); tabel tak ter-parse / kolom state hilang / row tak ter-resolve / nilai proyeksi invalid = **WARNING** | Teks aturan tidak memberi contoh; hasil perbandingan tabel bersifat warning-oriented (owner state = sumber kebenaran) |
| 21 | R8 | Row tabel tanpa sel state / row pendek | Row tanpa sel state dilewati diam-diam; row pendek (sel < header) = warning | Row tanpa state tidak dapat dibandingkan; row pendek menandakan tabel rusak |
| 22 | R9 | `fnd-` muncul di dua baris tabel validation.md (Knowledge doc vs Research Finding) | `fnd-` wajib **4 section**: Purpose, Content, Investigation Summary, Conclusion | Baris khusus (Research Finding) menang atas baris umum (Knowledge doc); research/README.md mengonfirmasi struktur 4 section |
| 23 | R9 | Pencocokan heading section wajib | `## Name` eksak atau `## Name <sesuatu>`; `###` tidak dihitung | Heading level 2 adalah konvensi struktur konten; varian dengan suffix diizinkan |
| 24 | R9 | ADR superseded — siapa penggantinya? | ≥1 artifact lain yang `supersedes`-nya resolve ke identity line ADR tersebut; referensi berversi wajib menunjuk instance eksak | ADR diganti per instance (identity line), bukan per line |
| 25 | R0 | Token tipe tidak dikenal | = error struktural (R0); aturan R2, R3, R4, R6, R7 dilewati untuk artifact tersebut; R5 tetap memeriksa referensi (token tak dikenal dilaporkan); R8/R9 tidak berlaku | Aturan bernomor tidak bermakna tanpa token yang dikenal; pengecualian R5 terverifikasi secara empiris |
| 26 | R0 | `instance-version: "1"` (dengan tanda kutip) | = error (bukan integer) | Spesifikasi mendefinisikan field sebagai integer; ADR kanonik menulis tanpa kutip |
| 27 | Scan | Direktori `testdata/` dan berawalan titik (`.git`) | **Tidak dituruni** | Fixture pengujian Go bukan konten knowledge base; tanpa ini fixture akan mematahkan self-validation |
| 28 | Scan | File non-.md, symlink, file tak terbaca | Non-.md diabaikan; symlink tidak diikuti; `.md` tak terbaca → `Validate` error | Pemindaian yang tidak melihat seluruh file tidak dapat menyatakan kepatuhan |
| 29 | Gap | Exactly-one-active container (protocol.md §3) tidak divalidasi | **Tidak** divalidasi; dicatat sebagai gap / kandidat aturan masa depan | Bukan bagian dari R1–R9; tugas implementasi melarang menambah aturan |

## Matriks traceability aturan (spesifikasi ↔ implementasi)

> **HISTORIS — dikonsolidasikan.** Tabel traceability ini dipindahkan ke
> [`conformance-traceability-matrix.md`](conformance-traceability-matrix.md) sebagai
> **single source of truth** cakupan konformansi (Engineering Requirement → Specification →
> Conformance Rule → Implementation → Automated Test). Dokumen ini kini hanya memegang
> keputusan interpretasi (#1–#29) dan gap yang diketahui. Perubahan conformance wajib
> memperbarui matriks tersebut (lihat `CONTRIBUTING.md`).

## Gap yang diketahui

Gap berikut **sengaja tidak divalidasi**; dicatat agar keputusan ini tidak hilang dan menjadi kandidat aturan masa depan:

| Gap | Detail | Rasional |
|---|---|---|
| R2 — id filename vs id frontmatter | Bagian id pada filename tidak dicocokkan dengan `id` frontmatter | Aturan hanya menuntut konsistensi token + versi; filename adalah proyeksi (ADR-001), Identity sejati di frontmatter |
| Exactly-one-active container | Invariant "tepat satu Execution Container aktif" (protocol.md §3) tidak diperiksa | Bukan bagian dari R1–R9; tugas implementasi melarang menambah aturan — kandidat aturan masa depan |
| R5 — referensi dua arah | Konvensi "referensi hanya ditulis pada artefak perujuk (bukan dua arah)" tidak di-enforce | Bukan pemeriksaan mekanis yang jelas; membutuhkan semantik arah per relationship |

## Pernyataan verifikasi

- **54 test** lulus (`go test ./...`): unit rule, parsing, filename, referensi, state, CLI (exit codes, determinisme keluaran), dan self-validation.
- `go vet ./...` bersih.
- **Self-validation PASS**: repositori EKA lolos validatornya sendiri — 7 artifact, 0 error, 0 warning, exit 0 — dikodifikasi sebagai `TestReferenceImplementationConforms` (`conformance/self_validation_test.go`).
