# CLI EKA — Validator Konformitas (`eka`)

> Dokumen konvensi, bukan artefak. Meta-dokumentasi zona `reference/`.
> Implementasi: Go (`cmd/eka` + `conformance/`), module `github.com/maleolabs/engineering-knowledge-architecture`.
> Terkait: [`conformance-notes.md`](conformance-notes.md) (keputusan interpretasi + traceability aturan), [`../skeleton/docs/exchange/validation.md`](../skeleton/docs/exchange/validation.md) (9 aturan konformitas).

## Filosofi CLI

`eka` adalah **bentuk executable dari spesifikasi EKA**. Konformitas repositori tidak boleh bergantung hanya pada review manual: aturan R1–R9 di `skeleton/docs/exchange/validation.md` dirancang mekanis (dapat diotomasi), dan CLI ini adalah implementasi kanoniknya. Dengan `eka validate`, kepatuhan terhadap standard menjadi sesuatu yang dapat diverifikasi mesin — sebelum commit, di CI, atau kapan pun diperlukan (P16: mekanisme enforcement bervariasi, invariant tetap identik).

Konsekuensi dari filosofi ini:

- Validator adalah **satu-satunya sumber kebenaran mekanis**; jika aturan teks ambigu, keputusan interpretasi didokumentasikan (lihat `conformance-notes.md`) — tidak ada perilaku yang diciptakan tanpa dasar.
- Repositori EKA sendiri harus lolos validatornya sendiri (lihat [Conformance repositori](#conformance-repositori)) — prasyarat agar validator dapat dipercaya oleh repositori lain.
- Perilaku CLI deterministik: dua kali menjalankan `eka validate` pada repositori yang sama menghasilkan keluaran yang identik byte-per-byte.

## Instalasi

Prasyarat: **Go 1.24+**.

### Dari module (setelah publikasi)

```sh
go install github.com/maleolabs/engineering-knowledge-architecture/cmd/eka@latest
```

Binari `eka` terpasang ke `$GOBIN` (default `$GOPATH/bin`).

### Dari sumber (repo ini)

```sh
cd <root-repo-eka>
go build -o eka ./cmd/eka
```

Hasil build adalah **binari standalone yang portabel** — tidak ada dependensi runtime selain binary itu sendiri; dapat disalin ke mesin lain (dengan arsitektur/OS yang sama) dan dijalankan langsung.

## Penggunaan

```
eka validate [path]
```

- `validate` — satu-satunya perintah yang diimplementasikan.
- `path` — opsional; root repositori yang divalidasi. Default: **direktori saat ini** (`.`).
- `eka -h` / `eka --help` / `eka help` mencetak usage dan keluar dengan kode 0.

### Contoh keluaran (repositori EKA itu sendiri)

```
EKA Conformance Validation
==========================
Root:      .
.md files: 51
Artifacts: 7
Errors:    0
Warnings:  0

Results (sorted by file, then rule):
  (no violations found)

Execution: PASS (0 errors, 0 warnings)
```

> Catatan: jumlah `.md files` adalah snapshot — bertambah setiap dokumen konvensi baru ditambahkan; format keluaran tetap. Jumlah artifact, error, dan warning adalah kontrak (7 artifact; error > 0 ⇒ FAIL).

Struktur keluaran:

1. **Ringkasan pemindaian** — root yang dipindai, jumlah file `.md`, jumlah artifact, jumlah error, jumlah warning.
2. **Hasil** — setiap pelanggaran dalam satu baris `[severity] rule file: pesan`; diurutkan deterministik berdasarkan file, lalu rule (R0, R1–R9), lalu severity, lalu pesan. Jika tidak ada pelanggaran, dicetak `(no violations found)`.
3. **Ringkasan eksekusi** — `PASS` jika tidak ada error blocking, `FAIL` jika ada; diikuti jumlah error dan warning.

### Contoh hasil dengan pelanggaran

```
  [error] R4 docs/decisions/adr-002-state-vector-encoding.md: missing owned state field existence-state on type "adr"
  [warning] R5 docs/decisions/adr-003-projection-model.md: unresolved reference "sto-x" in `depends-on` (allowed while content-state is draft)
```

## Exit codes

| Kode | Arti | Contoh |
|---|---|---|
| `0` | **Sepenuhnya patuh** — tidak ada error blocking; warning diperbolehkan | Repositori lolos semua aturan R1–R9 (0 error, berapa pun warning-nya) |
| `1` | **Ada pelanggaran blocking** — repositori tidak boleh di-commit sebelum diperbaiki | Setidaknya satu error (severity `error`) |
| `2` | **Kesalahan penggunaan/internal** — validasi tidak berjalan sama sekali | Perintah tidak dikenal, path tidak valid, root tidak dapat dibaca, `.md` tidak dapat dibaca |

Semantik warning: **warning tidak pernah memengaruhi exit code**. Sejalan dengan `validation.md` ("Hasil"): nilai `0` memblokir commit, nilai `W` hanya dicatat. Repositori dengan warning tetap keluar `0`.

Catatan penggunaan: `eka` tanpa argumen (tanpa subperintah) mencetak usage dan keluar `2` — bentuk kanonik adalah `eka validate [path]`.

## Proses validasi

Alur `eka validate [path]`:

1. **Pemindaian rekursif** — seluruh pohon direktori di bawah `path` ditelusuri.
2. **Klasifikasi** — setiap file `.md` diperiksa frontmatter-nya:
   - Frontmatter memuat **`type` DAN `id`** → **Artifact**, dievaluasi terhadap R1–R9.
   - Selain itu (README, `protocol.md`, `validation.md`, `transfer.md`, teks kanonik standard, dst.) → **Dokumen Konvensi**, dihitung tetapi dilewati.
   - Frontmatter memuat **tepat satu** dari `type`/`id` → malformed, dilaporkan sebagai error struktural (R0).
3. **Eksekusi 9 aturan konformitas (R1–R9)** — persis seperti didefinisikan di [`skeleton/docs/exchange/validation.md`](../skeleton/docs/exchange/validation.md):
   - R1 keunikan Identity, R2 konsistensi filename, R3 validitas nilai state, R4 kepatuhan owned-set, R5 integritas referensial, R6 dimensi == folder, R7 konsistensi change-log, R8 single-writer & proyeksi, R9 well-formedness.
   - Kesalahan struktural pra-aturan (frontmatter tidak valid, artifact rule dilanggar, field identity hilang/rusak, token tipe tidak dikenal) dikelompokkan sebagai **R0**.
   - Interpretasi mekanis setiap aturan: [`conformance-notes.md`](conformance-notes.md).
4. **Pelaporan deterministik** — hasil diurutkan (file → rule → severity → pesan), sehingga keluaran stabil antar mesin dan antar run.

### Scope pemindaian

- Hanya file berakhiran `.md` yang diperiksa; file lain diabaikan.
- Direktori bernama `testdata` dan direktori berawalan titik (mis. `.git`) **tidak dituruni** — fixture pengujian dan metadata VCS bukan konten knowledge base.
- Symlink tidak diikuti.
- File `.md` yang tidak dapat dibaca → `eka validate` gagal dengan error (exit `2`): pemindaian yang tidak dapat melihat seluruh file tidak dapat menyatakan kepatuhan.

## Conformance repositori

Repositori EKA **lolos suite konformansinya sendiri**: seluruh file `.md` dipindai, 7 artifact (7 Implementation ADR di `reference/decisions/`), **0 error, 0 warning, exit 0** (lihat contoh keluaran di atas).

Milestone ini dikodifikasi sebagai pengujian otomatis `TestReferenceImplementationConforms` di `conformance/self_validation_test.go`: test tersebut menemukan root repositori, menjalankan `Validate` atas seluruh repositori, dan menegaskan 0 error blocking. Artinya, regresi konformansi apa pun (mis. ADR baru yang melanggar aturan) mematahkan test suite sebelum sempat masuk commit.

## Arsitektur tooling

Validator dipisahkan menjadi dua lapisan:

| Lapisan | Lokasi | Peran |
|---|---|---|
| CLI | `cmd/eka/main.go` | Lapisan antarmuka: parsing argumen, usage, exit codes, format keluaran manusia. |
| Engine | `conformance/` (package publik) | Mesin validasi yang dapat digunakan ulang: pemindaian, klasifikasi artifact, aturan R1–R9, model hasil (`Report`). |

- Engine diimpor sebagai `github.com/maleolabs/engineering-knowledge-architecture/conformance`; **tidak memiliki dependensi apa pun ke CLI** (`cmd/`). CLI memanggil engine, bukan sebaliknya.
- Entry point engine: `Validate(root string) (*Report, error)` — `Report` berisi ringkasan pemindaian, daftar hasil, dan semantik `Pass()` (tidak ada error = lolos).
- Tooling masa depan (import/export, graph query, integrasi Knowledge OS) cukup mengimpor package `conformance` tanpa menempel ke antarmuka CLI.

## Roadmap CLI

Saat ini CLI **hanya menyediakan satu perintah: `validate`**. Perintah lain — `diagnose`, `import`, `export`, `graph` — adalah **arah masa depan dan BELUM diimplementasikan**; jangan menggunakannya. Setiap argumen perintah yang tidak dikenal ditolak dengan error yang menyebutkan scope validate-only (exit `2`).

| Perintah | Status | Catatan |
|---|---|---|
| `eka validate` | **Diimplementasikan** | Validator konformitas penuh (R1–R9 + R0 struktural). |
| `eka diagnose` | Belum diimplementasikan | Diagnostik repositori — kandidat masa depan. |
| `eka import` | Belum diimplementasikan | Impor artifact eksternal (seam Exchange, Section 13). |
| `eka export` | Belum diimplementasikan | Ekspor artifact (seam Exchange, Section 13). |
| `eka graph` | Belum diimplementasikan | Query/knowledge graph atas artifact. |
