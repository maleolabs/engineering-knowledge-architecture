# Contributing — Repositori EKA v1.0

Panduan kontribusi dan aturan governance untuk repositori **Reference Implementation EKA v1.0** (spec + reference implementation + validator + tooling).

> Dokumen konvensi, bukan artefak. Terkait: [`reference/conformance-traceability-matrix.md`](reference/conformance-traceability-matrix.md) (single source of truth cakupan konformansi), [`reference/conformance-notes.md`](reference/conformance-notes.md) (keputusan interpretasi), [`skeleton/docs/exchange/validation.md`](skeleton/docs/exchange/validation.md) (aturan konformitas R1–R9), [`standard/eka-specification-v1.0.md`](standard/eka-specification-v1.0.md) (standard kanonik).

## Prinsip

- **Satu proyek kanonik.** Repositori ini adalah satu kesatuan: standard (`standard/`), reference implementation (`skeleton/`), validator (`conformance/` + `cmd/eka/`), dan dokumentasi (`reference/`). Perubahan tidak boleh memperlakukan satu bagian sebagai proyek terpisah.
- **Kontribusi lengkap = perubahan spec → validator → test → matriks dalam Pull Request yang SAMA.** Tidak ada PR parsial yang mengubah perilaku konformansi tanpa menutup seluruh rantai traceability.
- **Dokumentasi adalah bagian dari produk.** Setiap keputusan interpretasi didokumentasikan sebelum implementasi; tidak ada perilaku yang diciptakan tanpa dasar spec (lihat kebijakan interpretasi di `reference/conformance-notes.md`).

## Definisi "implementasi tidak lengkap"

Suatu kontribusi dianggap **tidak lengkap** jika mengubah spec, aturan konformitas, validator, atau test **TANPA memperbarui** [`reference/conformance-traceability-matrix.md`](reference/conformance-traceability-matrix.md). PR semacam itu ditolak sampai matriks disinkronkan.

## Workflow kontribusi

1. **Pahami EKA.** Baca `standard/` (standard kanonik — teks normatif), `skeleton/docs/` (struktur dan aturan serialisasi), dan `reference/` (konvensi implementasi). Mulai dari `README.md` root.
2. **Jika mengubah conformance behavior**, perbarui dalam PR yang sama:
   - teks aturan: `skeleton/docs/exchange/validation.md` (R1–R9) atau definisi R0 di `conformance/`;
   - implementasi: `conformance/` (rule/helper) atau `cmd/eka/` (CLI);
   - test: `*_test.go` di `conformance/` + `cmd/eka/`;
   - matriks: `reference/conformance-traceability-matrix.md` (Section 2–4);
   - catatan interpretasi: `reference/conformance-notes.md` — tambahkan keputusan interpretasi **baru** dengan nomor berikutnya; keputusan lama tidak dihapus.
3. **Kualitas wajib** — semua harus lulus sebelum PR:
   ```sh
   go test ./...
   go vet ./...
   gofmt -l .
   go build ./...
   go run ./cmd/eka validate .
   ```
   `go run ./cmd/eka validate .` (self-validation) harus PASS: 0 error, exit 0.
4. **Matrix maintenance rules**:
   - **Aturan baru** → baris baru di Section 2 (matriks utama) + requirement `REQ-nnn` baru di Section 3. REQ ID tidak pernah digunakan ulang; nomor baru = lanjutan.
   - **Rule dimodifikasi** → update kolom Spec Anchor / Implementation / Automated Tests / Notes pada baris rule terkait; jika semantik requirement berubah, update frasa requirement di Section 3.
   - **Test baru** → update Section 4 (test coverage index). Setiap fungsi test muncul tepat satu kali — jangan duplikasi, jangan tinggalkan orphan.
   - **Coverage status berubah** → update Section 5 (coverage analysis) dan Section 6 (summary) dengan angka yang konsisten.
5. **Dokumentasi wajib konsisten**: jika CLI berubah, update `reference/cli.md`; jika interpretasi berubah, update `reference/conformance-notes.md`; jika dokumen baru ditambahkan ke zona `reference/`, tambahkan baris indeks di `reference/README.md`.

## Governance rule formal

- [`reference/conformance-traceability-matrix.md`](reference/conformance-traceability-matrix.md) adalah **single source of truth** cakupan konformansi: Engineering Requirement → Specification → Conformance Rule → Implementation → Automated Test → Coverage Status → Notes.
- Matriks **WAJIB diperbarui dalam Pull Request yang sama** dengan perubahan spec/rule/implementasi/test — dan **sebaliknya**: matriks tidak pernah diedit tanpa perubahan terkait (tidak ada PR "matriks saja").
- ID rule (R0–R9) dan ID requirement (REQ-nnn) bersifat stabil dan tidak pernah digunakan ulang.
- Perubahan perilaku validator di luar R0–R9 (aturan baru) **tidak** diusulkan melalui implementasi langsung: aturan baru adalah keputusan governance (lihat gap yang diketahui di `reference/conformance-notes.md` dan rekomendasi follow-up di matriks Section 5–6).

## Review checklist (untuk reviewer)

Sebelum menyetujui PR yang menyentuh conformance:

- [ ] Matriks konsisten dengan perubahan: Section 2 (rule), Section 3 (requirement), Section 4 (test), Section 5–6 (analisis & summary)?
- [ ] Tidak ada orphan: setiap rule R0–R9 muncul tepat satu kali di Section 2; setiap test function muncul tepat satu kali di Section 4?
- [ ] Tidak ada rule/requirement ID yang digunakan ulang?
- [ ] Angka summary valid: jumlah test di Section 6 == jumlah baris test di Section 4?
- [ ] `go test ./...`, `go vet ./...`, `gofmt`, `go build`, dan `go run ./cmd/eka validate .` lulus?
- [ ] Interpretasi baru terdokumentasi di `reference/conformance-notes.md` (nomor lanjutan), jika ada?
- [ ] Dokumentasi terkait (README, `reference/README.md`, `reference/cli.md`) diperbarui jika nama/struktur berubah?
