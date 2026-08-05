# CLI EKA — Antarmuka Resmi (`eka`)

> Dokumen konvensi, bukan artefak. Meta-dokumentasi zona `reference/`.
> Implementasi: Go (`cmd/` + `bootstrap/` + `conformance/`), module `github.com/maleolabs/engineering-knowledge-architecture`.
> Terkait: [`conformance-notes.md`](conformance-notes.md) (keputusan interpretasi + traceability aturan), [`../skeleton/docs/exchange/validation.md`](../skeleton/docs/exchange/validation.md) (9 aturan konformitas), [`eka-reference-serialization-format-v1.0.md`](eka-reference-serialization-format-v1.0.md) (target serialisasi `export`/`import` masa depan).

## Filosofi CLI

`eka` adalah **bentuk executable dari spesifikasi EKA** — antarmuka resmi Engineering Knowledge Architecture bagi manusia dan agent (Naming and Terminology Specification v1.0 §7). Dua peran saat ini:

- **`eka init`** — Repository Bootstrapper resmi: menganalisis workspace, menyusun rencana bootstrap, mengonfigurasi proyek secara interaktif, membangkitkan repositori EKA dari Reference Skeleton, lalu memvalidasinya.
- **`eka validate`** — Validator konformitas: konformitas repositori tidak boleh bergantung hanya pada review manual — aturan R1–R9 di `skeleton/docs/exchange/validation.md` dirancang mekanis, dan validator ini adalah implementasi kanoniknya (P16: mekanisme enforcement bervariasi, invariant tetap identik).

Konsekuensi dari filosofi ini:

- Validator adalah **satu-satunya sumber kebenaran mekanis**; jika aturan teks ambigu, keputusan interpretasi didokumentasikan (lihat `conformance-notes.md`) — tidak ada perilaku yang diciptakan tanpa dasar.
- Repositori EKA sendiri harus lolos validatornya sendiri (lihat [Conformance repositori](#conformance-repositori)) — prasyarat agar validator dapat dipercaya oleh repositori lain.
- Perilaku CLI deterministik: dua kali menjalankan perintah yang sama pada input yang sama menghasilkan keluaran yang identik byte-per-byte.
- CLI adalah **adapter** — logika bisnis (bootstrap, validasi) hidup di package aplikasi yang dapat digunakan ulang, independen dari framework CLI (lihat [Arsitektur CLI](#arsitektur-cli)).

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
eka init [project-name] [--dry-run]
eka validate [path]
eka completion [bash|zsh|fish|powershell]
eka help [command]
```

- `eka` tanpa subperintah mencetak usage dan keluar `2`.
- `eka -h` / `eka --help` / `eka help` mencetak help dan keluar `0`.
- `eka help <command>` mencetak help perintah.

### Exit codes

| Kode | Arti | Contoh |
|---|---|---|
| `0` | **Sukses** — validasi patuh (warning diizinkan) / inisialisasi selesai / dry-run / help / completion | Repositori lolos R1–R9; `eka init` selesai dan tervalidasi |
| `1` | **Pelanggaran blocking** — validasi menemukan error, atau repositori hasil `eka init` gagal validasi | Setidaknya satu error (severity `error`) |
| `2` | **Kesalahan penggunaan/internal** — perintah tidak berjalan | Perintah tidak dikenal, flag tidak dikenal, argumen berlebihan, path tidak valid |

Semantik warning: **warning tidak pernah memengaruhi exit code**. Repositori dengan warning tetap keluar `0`.

## `eka init` — Repository Bootstrapper

`eka init` bukan template generator. Ia adalah **Repository Bootstrapper resmi** untuk repositori pengetahuan Engineering Knowledge Architecture: menganalisis workspace terlebih dahulu, menyusun rencana, bertanya hanya yang belum diketahui, membangkitkan dari Reference Skeleton, dan memvalidasi hasilnya.

### Filosofi inisialisasi

Empat prinsip:

1. **Pahami dulu, ubah kemudian** — inisialisasi berjalan dalam lima tahap: *Workspace Discovery → Bootstrap Planning → Interactive Configuration → Repository Generation → Validation*. Workspace dipahami sebelum ada satu pun modifikasi; tidak ada generasi buta.
2. **Adaptif** — discovery otomatis memengaruhi wizard: pengguna tidak pernah ditanya hal yang jawabannya sudah diketahui.
3. **Idempoten** — `eka init` aman dijalankan berulang; artefak yang ada dideteksi, dipakai ulang, dilewati, atau dikonfirmasi eksplisit sebelum diganti. Konten pengguna tidak pernah ditimpa diam-diam.
4. **Teruji otomatis** — inisialisasi dianggap berhasil hanya jika repositori hasil generasi lolos `eka validate`.

### Alur kerja lima tahap

```
Workspace Discovery
        ↓
Bootstrap Planning
        ↓
Interactive Configuration
        ↓
Repository Generation
        ↓
Validation
```

1. **Workspace Discovery** — memeriksa target: ada/tidaknya direktori, direktori Git (di target maupun leluhurnya), README, direktori `docs/`, repositori EKA yang sudah ada (penanda: `docs/operating/` + `docs/exchange/validation.md` + `docs/exchange/transfer.md`), dan file konfigurasi (informasional: `.gitignore`, `.editorconfig`, `.eka.*`, `eka.*`).
2. **Bootstrap Planning** — menyusun rencana deterministik: direktori yang akan dibuat, file yang akan dibangkitkan, resource yang dipakai ulang, status Git, langkah validasi. Rencana membawa konten file yang akan ditulis — dry-run dan eksekusi tidak mungkin menyimpang.
3. **Interactive Configuration** — wizard adaptif (lihat di bawah).
4. **Repository Generation** — menyalin Reference Skeleton (embedded dari `skeleton/`, lihat Arsitektur) ke target: `docs/**` verbatim + `README.md` dengan judul diganti nama proyek. File yang sudah ada dengan konten identik dipakai ulang; yang berbeda dikonfirmasi (interaktif) atau dilewati (non-interaktif). `git init` hanya jika direncanakan (lihat deteksi Git).
5. **Validation** — `conformance.Validate(target)` dijalankan; hasil dicetak. Inisialisasi sukses hanya jika repositori lolos.

### Mode

| Perintah | Perilaku |
|---|---|
| `eka init` | Inisialisasi direktori saat ini |
| `eka init .` | Setara dengan `eka init` |
| `eka init <project-name>` | Membuat direktori proyek baru dan menginisialisasinya sebagai repositori EKA |
| `eka init [name] --dry-run` | Mencetak rencana bootstrap tanpa menulis apa pun; keluar `0` |

### Wizard adaptif

Wizard hanya menampilkan pertanyaan yang relevan:

| Pertanyaan | Muncul hanya jika |
|---|---|
| Project Name | Nama direktori target kosong/tidak dapat dipakai |
| Namespace | Selalu (wajib; huruf kecil, digit, tanda hubung; tanpa `/`, `:`, spasi) |
| Project Description | Selalu (opsional; boleh kosong) |
| Generate README | README belum ada |
| Initialize Git | Git belum ada **dan** binary `git` tersedia |

Tidak ada pertanyaan "Methodology" — EKA v1 tidak memiliki taksonomi methodology kanonik. Jika jawaban tidak dapat dibaca (stdin bukan terminal atau EOF), wizard memakai default deterministik dan **tidak pernah** menjalankan `git init`.

### Deteksi Git cerdas

- Sudah ada repositori Git (di target atau leluhur) → pertanyaan Git dilewati, `git init` tidak dieksekusi.
- Direktori baru / belum ada Git → ditawarkan inisialisasi; jika diterima: `git init`; jika ditolak: lanjut normal.
- Non-interaktif (pipe, `/dev/null`, file, CI) → Git tidak pernah diinisialisasi.
- Kegagalan `git init` (mis. binary tidak ada) → peringatan, inisialisasi tetap dilanjutkan.

### Idempotensi

`eka init` dua kali pada target yang sama tidak pernah merusak repositori:

- Repositori EKA yang sudah ada → terdeteksi, semua resource dipakai ulang, hanya validasi dijalankan (keluaran: "already initialized").
- File yang akan dibangkitkan sudah ada:
  - konten identik → dipakai ulang (tanpa tulis);
  - konten berbeda → konfirmasi eksplisit (interaktif; default tolak) atau dilewati + dilaporkan (non-interaktif).

### Dry-run

`eka init --dry-run` mencetak rencana bootstrap — direktori yang akan dibuat, file yang akan dibangkitkan, resource yang dipakai ulang, status Git, rencana validasi — **tanpa menulis satu pun file**. Deterministik: dua dry-run pada state yang sama menghasilkan keluaran identik.

### Ringkasan keluaran

Setelah selesai, dicetak ringkasan ringkas: Project Name, Namespace, Repository Type (new / existing-dir / existing-eka), Git Status, Knowledge Standard Version (EKA v1.0), Validation Result (PASS/FAIL + jumlah error/warning), dan langkah selanjutnya yang disarankan.

## `eka validate` — Validator Konformitas

```
eka validate [path]
```

- `path` — opsional; root repositori yang divalidasi. Default: **direktori saat ini** (`.`).

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

### Proses validasi

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

## Shell completion

`eka completion [bash|zsh|fish|powershell]` mencetak script completion untuk shell yang dipilih (disediakan oleh framework Cobra). Gunakan, misalnya:

```sh
source <(eka completion bash)
```

## Conformance repositori

Repositori EKA **lolos suite konformansinya sendiri**: seluruh file `.md` dipindai, 7 artifact (7 Implementation ADR di `reference/decisions/`), **0 error, 0 warning, exit 0** (lihat contoh keluaran di atas).

Milestone ini dikodifikasi sebagai pengujian otomatis `TestReferenceImplementationConforms` di `conformance/self_validation_test.go`: test tersebut menemukan root repositori, menjalankan `Validate` atas seluruh repositori, dan menegaskan 0 error blocking. Artinya, regresi konformansi apa pun (mis. ADR baru yang melanggar aturan) mematahkan test suite sebelum sempat masuk commit.

## Arsitektur CLI

CLI diorganisasi sebagai **dua lapisan + satu titik masuk**:

| Lapisan | Lokasi | Peran |
|---|---|---|
| Command layer | `cmd/` (package `cmd`) | **Hanya** definisi perintah Cobra: registrasi, flag, help, validasi argumen, dispatch ke layanan. Tidak ada logika domain. |
| Application layer | `bootstrap/` (package publik) | Repository Bootstrapper: discovery, planning, wizard, generasi, validasi — dapat digunakan ulang tanpa CLI. |
| Application layer | `conformance/` (package publik) | Mesin validasi: pemindaian, klasifikasi artifact, aturan R1–R9, model hasil (`Report`) — dapat digunakan ulang tanpa CLI. |
| Entry point | `cmd/eka/main.go` | Tipis: `os.Exit(cmd.Execute(...))`. Nama executable: `eka`. |

```
cmd/                package cmd — Cobra command definitions (command layer)
  root.go           root command + Execute(args, stdin, stdout, stderr) int
  validate.go       perintah validate
  init.go           perintah init
  execute_test.go   test CLI (exit codes, help, completion, mode)
cmd/eka/
  main.go           tipis: os.Exit(cmd.Execute(...))
bootstrap/          package publik — engine eka init (application layer)
conformance/        package publik — engine validasi (application layer)
skeletonembed.go    package root — //go:embed skeleton (Reference Skeleton kanonik)
```

Prinsip:

- **Cobra adalah adapter, bukan arsitektur.** Framework (saat ini Cobra) adalah detail implementasi antarmuka perintah. Logika bisnis hidup di `bootstrap/` dan `conformance/` — package publik yang diimpor sebagai `github.com/maleolabs/engineering-knowledge-architecture/bootstrap` dan `.../conformance`, **tanpa dependensi apa pun ke `cmd/`**.
- **Command layer memanggil layanan, bukan sebaliknya.** Tooling masa depan (import/export, graph query, SDK, integrasi Knowledge OS) cukup mengimpor package aplikasi tanpa menempel ke Cobra.
- **Tidak ada `internal/` atau `pkg/`** — tidak ada konsumen internal kedua; `bootstrap/` dan `conformance/` sudah merupakan public API. Menambah direktori tanpa tujuan langsung adalah speculative abstraction (dilarang).
- **Reference Skeleton ter-embed** (`skeletonembed.go`): `eka init` membangkitkan repositori dari sumber kanonik `skeleton/`, bukan dari direktori yang di-hardcode. Binari standalone tetap dapat membangkitkan struktur tanpa checkout repositori.
- **Exit codes deterministik** (0/1/2) dipetakan di `cmd/root.go`; semua error melewati satu jalur keluaran `eka: <pesan>`.

## Panduan kontribusi: menambah perintah

Perintah baru ditambahkan tanpa refactor arsitektur:

1. **Layanan dulu** — implementasi logika bisnis di package aplikasi (`bootstrap/`, `conformance/`, atau package publik baru), lengkap dengan test-nya. CLI tidak boleh berisi logika domain.
2. **Definisikan command** — file baru `cmd/<name>.go`, package `cmd`: `Use` (verb, Naming §7.1), `Short` (satu baris), `Long` (detail), `Example`, flag (dengan deskripsi), validasi argumen via Cobra (mis. `MaximumNArgs`), `RunE` memanggil layanan lalu merender keluaran.
3. **Daftarkan di root** — tambahkan ke `rootCmd.AddCommand(...)` di `cmd/root.go`.
4. **Exit codes** — sukses → `nil` (0); pelanggaran blocking → `*exitError{code: 1}` (atau sentinel setara); kesalahan penggunaan/internal → error biasa (dipetakan ke 2).
5. **Test** — tambahkan kasus di `cmd/execute_test.go` (exit codes, help, determinisme) + test layanan di package aplikasi.
6. **Dokumentasikan** — perbarui dokumen ini (tabel perintah, contoh) dan matriks traceability.
7. **Nama mengikuti Naming and Terminology Specification v1.0 §7** — subperintah adalah verb (`validate`, `init`, `diagnose`, `import`, `export`, `sync`, `format`, `graph`); jangan perkenalkan pola baru.

## Roadmap CLI

| Perintah | Status | Catatan |
|---|---|---|
| `eka init` | **Diimplementasikan** | Repository Bootstrapper (5 tahap, wizard adaptif, dry-run, idempoten, validasi pasca-generasi). |
| `eka validate` | **Diimplementasikan** | Validator konformitas penuh (R1–R9 + R0 struktural). |
| `eka completion` | **Diimplementasikan** | Script completion bash/zsh/fish/powershell (disediakan Cobra). |
| `eka diagnose` | Belum diimplementasikan | Diagnostik repositori — kandidat masa depan. |
| `eka import` | Belum diimplementasikan | Impor artifact eksternal (seam Exchange, Section 13). |
| `eka export` | Belum diimplementasikan | Ekspor artifact (seam Exchange, Section 13); target serialisasi: RSF v1.0. |
| `eka graph` | Belum diimplementasikan | Query/knowledge graph atas artifact. |

Perintah masa depan ditambahkan mengikuti [Panduan kontribusi](#panduan-kontribusi-menambah-perintah) — tanpa refactor arsitektur.
