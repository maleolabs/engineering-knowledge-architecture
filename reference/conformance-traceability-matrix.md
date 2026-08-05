# Conformance Traceability Matrix — EKA v1.0

| Properti | Nilai |
|---|---|
| **Status** | Ratified — komponen resmi EKA v1.0 |
| **Versi** | v1.0 |
| **Jenis dokumen** | Dokumen konvensi (bukan artefak) — zona `reference/` |
| **Terkait** | [`conformance-notes.md`](conformance-notes.md) (interpretasi + gap), [`../skeleton/docs/exchange/validation.md`](../skeleton/docs/exchange/validation.md) (9 aturan konformitas), [`../standard/eka-specification-v1.0.md`](../standard/eka-specification-v1.0.md) (standard kanonik) |

> **Aturan governance (formal).** Matriks ini adalah **single source of truth** cakupan konformansi EKA v1.0. Matriks **WAJIB diperbarui dalam Pull Request yang sama** dengan perubahan pada spesifikasi, aturan konformitas (`validation.md`), implementasi validator (`conformance/`, `cmd/eka/`), atau test — dan sebaliknya, matriks tidak pernah diedit tanpa perubahan terkait. Lihat [`../CONTRIBUTING.md`](../CONTRIBUTING.md).

---

## 1. Layer model

Matriks ini menelusuri cakupan konformansi melalui **5 lapisan**, dari kebutuhan hingga bukti otomatis:

| # | Layer | Identifikasi | Contoh | Sumber kebenaran |
|---|---|---|---|---|
| 1 | **Engineering Requirement** | `REQ-nnn` (REQ-001..REQ-016) | REQ-002 Identity uniqueness | Section 3 dokumen ini |
| 2 | **Specification** | Anchor `§` (nomor section/principle) | §6.2.2, P3 | `standard/eka-specification-v1.0.md` |
| 3 | **Conformance Rule** | `Rn` (R0, R1–R9) | R1 | `skeleton/docs/exchange/validation.md` (R1–R9) + R0 (struktural, didefinisikan `conformance/`) |
| 4 | **Implementation** | `file:func` (relatif paket) | `conformance/rules.go:rule1` | Kode Go `conformance/` + `cmd/eka/` |
| 5 | **Automated Test** | nama fungsi `TestXxx` | `TestRule2ExactCounts` | `*_test.go` di `conformance/` + `cmd/eka/` |

**Konvensi identifier (deterministik, siap konsumsi otomatis):**

- **REQ ID** — `REQ-<3 digit>`; stabil, **tidak pernah digunakan ulang**. Requirement yang dihapus/diganti tetap dicadangkan; requirement baru mendapat ID berikutnya.
- **Rule ID** — R0 (struktural) + R1–R9, tetap dari `validation.md`; tidak ada penomoran ulang.
- **Implementation** — `path/file.go:funcName` relatif terhadap root repo; fungsi pendukung yang dilayani satu rule ditulis dalam kurung setelah fungsi utamanya.
- **Test** — nama fungsi Go (`TestXxx`); satu fungsi = satu baris di Section 4.
- **Coverage Status** — nilai enumerasi tetap: `Enforced (tested)` (rule diimplementasikan + dites), `Governance-only (uncovered)` (normatif di spec, tidak di-enforce mekanis), `Partially enforced` (sebagian permukaan di-enforce).

---

## 2. Matriks utama

Satu baris per Conformance Rule. Status `Enforced (tested)` berlaku untuk seluruh 10 rule: setiap rule memiliki implementasi Go dan cakupan test otomatis.

| Rule | Requirement ID(s) | Spec Anchor | Implementation | Automated Tests | Coverage Status | Notes |
|---|---|---|---|---|---|---|
| **R0** (struktural) | REQ-001 | §3 (Artifact), §5.1 | `conformance/artifact.go:analyzeFile` | `TestAnalyzeNoFrontmatterIsConventionDoc`, `TestAnalyzeFrontmatterWithoutTypeIDIsConventionDoc`, `TestAnalyzeUnterminatedFrontmatter`, `TestAnalyzeBrokenYAML`, `TestAnalyzeTypeXorID`, `TestAnalyzeValidArtifact`, `TestAnalyzeMissingIdentityFields`, `TestAnalyzeNonIntVersion`, `TestAnalyzeInvalidDate`, `TestAnalyzeUnknownType`, `TestAnalyzeChangeLogNotList`, `TestAnalyzeMalformedChangeLogEntry`, `TestUnknownTypeIsStructural` | Enforced (tested) | Bucket struktural sebelum aturan bernomor (bukan salah satu dari 9 aturan). Artifact rule: `type` + `id` ⇒ Artifact; `type` XOR `id` = malformed. Token tipe tak dikenal = error struktural (interpretasi #25); `instance-version` non-integer = error (#26). Artifact gagal klasifikasi → aturan R2, R3, R4, R6, R7 dilewati untuk file itu; R1 tetap berlaku (identitas tetap diindeks); R5 tetap memeriksa referensi dan melaporkan token tak dikenal; R8/R9 tidak berlaku untuk tipe tak dikenal. |
| **R1** | REQ-002 | §6.2.2, P3 | `conformance/rules.go:rule1` (+ `conformance/validate.go:identityKey`/`buildIndex`) | — (unit test khusus tidak ada; dicover via `TestInvalidFixtures` → fixture `invalid-dup-identity`) | Enforced (tested) | Duplikat `(namespace, type, id, instance-version)` = error. ID unik dalam `(namespace, type)`; InstanceVersion unik dalam line. |
| **R2** | REQ-003 | §6.4, P9 | `conformance/rules.go:rule2` (+ `conformance/filename.go:parseFilename`) | `TestRule2ExactCounts` (+ infra: `TestParseFilename`, `TestParseFilenameEmpty`; fixture `invalid-filename`) | Enforced (tested) | Token filename == `type` frontmatter; akhiran `-v<nn>` == `instance-version`; `-v<nn>` hanya `scp-`/`plan-` dan WAJIB ada (termasuk v1). Jumlah digit bebas (`-v1`/`-v01` valid, #16). **Gap terdokumentasi:** bagian id pada filename TIDAK dicocokkan dengan `id` frontmatter (#17) — filename adalah proyeksi (ADR-001), Identity sejati di frontmatter. |
| **R3** | REQ-004 | §7.2 | `conformance/rules.go:rule3` (+ `conformance/state.go:contentStateVariant`/`domainValues`) | `TestPhaseValidation` (+ infra: `TestPhaseValueSet`; fixture `invalid-state-value`) | Enforced (tested) | Nilai field state ∈ value set domain. Varian content-state per keluarga tipe: living (`draft/review/approved/amended`), ADR (`proposed/accepted/superseded`), decision (`draft/accepted/superseded`). `phase` hanya pada `scp-`/`plan-` dan ∈ phase value set. |
| **R4** | REQ-005 | §7.4, §10 | `conformance/rules.go:rule4` | — (unit test khusus tidak ada; dicover via `TestInvalidFixtures` → fixture `invalid-ownership`) (+ infra: `TestOwnedSets`) | Enforced (tested) | Field state pada file == owned set tipe-nya. Field owned absen = error (interpretasi #2); field non-owned hadir = error; `tkt-` state vector kosong. Binding type→state dari §10. |
| **R5** | REQ-006 | §6.2.7, §13.2.3 | `conformance/rules.go:rule5` (+ `parseReference`, `resolve`) | `TestRule5DraftSeverity`, `TestRule5CrossNamespaceAndVersionResolution`, `TestRule5VersionedReferenceToMissingInstance` (+ infra: `TestParseReference`, `TestParseReferenceCrossNamespace`; fixture `invalid-reference`) | Enforced (tested) | Referensi malformed = error selalu; unresolved pada `content-state: draft` = warning, selain itu (termasuk tanpa content-state) = error (interpretasi #10). Self-reference = error. **Bare-id diterima** sebagai line reference dalam namespace+type perujuk (#9 — 6/7 ADR nyata memakai bentuk ini). Format lintas-namespace `<ns>/<type>:<id>[:<ver>]` (#12). **Tidak di-enforce:** konvensi referensi dua arah (#11, gap terdokumentasi). |
| **R6** | REQ-007 | §8, P15 | `conformance/rules.go:rule6` (+ `dimensionFolderFor`, `dimensionList`) | `TestDimensionFolderResolution` (+ infra: `TestDimensionTokens`; fixture `invalid-dimension`) | Enforced (tested) | Artifact pengetahuan: `dimension` == folder rumah; folder rumah = **ancestor terdekat** yang namanya ∈ 12 dimensi (#13); tidak ditemukan = error. `dimension` wajib pada artifact pengetahuan (#14), dilarang pada `ctr-`/`tkt-`/`ses-`; `dimensions-secondary` divalidasi token saja (#15). |
| **R7** | REQ-008 | §5.2, P7 | `conformance/rules.go:rule7` (+ `entriesForDomain`, `indexOfEntry`; `conformance/state.go:isLegalTransition`) | `TestRule7ExactCounts` (+ infra: `TestIsLegalTransition`; fixture `invalid-changelog`) | Enforced (tested) | Setiap domain owned punya entri change-log; entri terakhir == nilai field saat ini; transisi legal; initial `from: "-"` legal untuk semua domain, `to: "-"` = invalid (#5). **Tidak di-enforce:** kontiguitas rantai `entry[i].from == entry[i-1].to` (#4 — snapshot statis). Execution State strictly adjacent; domain lain forward-only tanpa adjacency wajib (#6). Entri `phase`: value-set + well-formedness saja (#7). `tkt-` (state vector kosong) boleh tanpa change-log (#8). |
| **R8** | REQ-009 | §7.4, P6 | `conformance/rules.go:rule8` (+ `workItemsTable`, `compareWorkItemsTable`, `resolveWorkItemCell`, `hasProjectionHeader`, `splitTableRow`) | `TestRule8TicketHeaderAndDerivation`, `TestRule8ContainerTableMismatchIsWarningOnly` (+ fixture `invalid-projection`) | Enforced (tested) | `tkt-`: state vector kosong + `derives-from` resolve ke artifact `ctr-` (#19) + header proyeksi wajib ada (posisi bebas, #18). Tabel `## Work Items` pada `ctr-`: format GFM; mismatch dengan owner state = **warning** (owner state = sumber kebenaran); tabel tak ter-parse / state hilang / row tak ter-resolve = warning (#20–21). |
| **R9** | REQ-010 | §10, §14.2.6 | `conformance/rules.go:rule9` (+ `requiredSectionsFor`, `headingMatches`, `hasReplacement`) | `TestRule9SupersededADRWithReplacement`, `TestRule9VersionedReplacementMustNameInstance` (+ fixture `invalid-sections`, `invalid-adr-superseded`) | Enforced (tested) | Bagian wajib per keluarga tipe (validation.md Aturan 9). **`fnd-` wajib 4 section** (Purpose, Content, Investigation Summary, Conclusion — #22). Pencocokan heading level 2: `## Name` eksak atau `## Name <suffix>` (#23). ADR `superseded` wajib dirujuk ≥1 artifact lain via `supersedes` yang resolve ke identity line-nya; referensi berversi wajib menunjuk instance eksak (#24). |

---

## 3. Requirement index

### 3(a) Enforced requirements (di-enforce oleh rule)

| Req ID | Requirement | Spec anchor | Enforced by | Status | Notes |
|---|---|---|---|---|---|
| REQ-001 | **Artifact identity rule** — setiap artifact membawa `type` + `id`; frontmatter dengan hanya salah satunya, token tipe tak dikenal, atau field identity invalid (namespace, `instance-version` non-integer, tanggal invalid) = malformed (struktural) | §3, §5.1 | R0 | Enforced (tested) | Gagal struktural di bucket R0; aturan bernomor tidak dijalankan untuk artifact tersebut (interpretasi #25) |
| REQ-002 | **Identity uniqueness** — tidak ada duplikat `(namespace, type, id, instance-version)` di seluruh repositori; ID unik dalam `(namespace, type)`, InstanceVersion unik dalam line | §6.2.2, P3 | R1 | Enforced (tested) | Diperiksa lintas seluruh repositori (rule atas seluruh set artifact) |
| REQ-003 | **Filename as projection of identity** — token filename == `type`; akhiran `-v<nn>` == `instance-version`; akhiran hanya `scp-`/`plan-` dan wajib ada (termasuk v1) | §6.4, P9 | R2 | Enforced (tested) | Bagian id filename vs `id` frontmatter tidak diperiksa (gap terdokumentasi, interpretasi #17) |
| REQ-004 | **State value validity** — setiap nilai field state ∈ value set domain-nya (termasuk varian content-state living/ADR/decision); `phase` ∈ phase value set dan hanya pada `scp-`/`plan-` | §7.2 | R3 | Enforced (tested) | Value set domain dari tabel §7.2 + varian per keluarga tipe |
| REQ-005 | **Owned state vector** — field state pada file == owned set tipe-nya; `tkt-` state vector kosong; tidak ada field state milik tipe lain | §7.4, §10 | R4 | Enforced (tested) | Owned set dari §10 (binding type→state adalah bagian standard) |
| REQ-006 | **Referential integrity** — semua referensi (`amends`, `supersedes`, `derives-from`, `depends-on`, `validates`) resolve ke artifact yang ada; malformed = error; unresolved pada draft = warning, selain itu error; self-reference = error | §6.2.7, §13.2.3 | R5 | Enforced (tested) | Bare-id diterima per interpretasi #9; konvensi dua arah tidak di-enforce (#11) |
| REQ-007 | **Classification-location consistency** — artifact pengetahuan: `dimension` == folder rumah (ancestor terdekat ∈ 12 dimensi); `ctr-`/`tkt-`/`ses-` tidak boleh membawa `dimension` | §8, P15 | R6 | Enforced (tested) | Klasifikasi = properti, bukan identity (P15) |
| REQ-008 | **Change-log consistency** — setiap domain owned punya entri; entri terakhir == nilai saat ini; transisi legal; initial `from: "-"` | §5.2, P7 | R7 | Enforced (tested) | Kontiguitas rantai tidak di-enforce (interpretasi #4) |
| REQ-009 | **Single-writer & projection discipline** — `tkt-` state vector kosong + `derives-from` ctr- + header proyeksi; tabel Work Items kontainer divalidasi terhadap owner state | §7.4, P6 | R8 | Enforced (tested) | Mismatch tabel vs owner state = warning (owner = sumber kebenaran) |
| REQ-010 | **Well-formed content** — bagian wajib per keluarga tipe ada; ADR `superseded` wajib dirujuk penggantinya | §10, §14.2.6 | R9 | Enforced (tested) | `fnd-` wajib 4 section (interpretasi #22); heading level 2 (#23) |

### 3(b) Governance-only requirements (normatif di spec, TIDAK di-enforce rule mana pun)

| Req ID | Requirement | Spec anchor | Enforced by | Status | Notes |
|---|---|---|---|---|---|
| REQ-011 | **Concurrency control** — tepat satu Execution Container aktif (exactly-one-active); pembuatan berikutnya menunggu | §9 (Concurrency Control), §7.5, §5.2 | — | Governance-only (uncovered) | Invariant Operating Layer; validasi mekanis membutuhkan observasi lintas artifact + waktu. Gap #29 di `conformance-notes.md` — kandidat aturan masa depan |
| REQ-012 | **Plan immutability / lock-atomic-with-generation** — pembuatan Execution Container mengunci plan secara atomik; Content plan locked tidak berubah; perubahan pasca-lock = instance baru | §5.2, §9 (Versioning/Immutability) | — | Governance-only (uncovered) | Transisi yang tidak dapat diobservasi dari snapshot statis tunggal; membutuhkan riwayat lintas instance |
| REQ-013 | **Approved-content immutability** — Content yang telah melewati gate persetujuan tidak dimutasi diam-diam; perubahan hanya via kanal governance (amend/supersede) | P8, §5.4 (invariant 6) | — | Governance-only (uncovered) | Membutuhkan bukti non-mutasi lintas waktu; snapshot tidak dapat membuktikannya |
| REQ-014 | **Phase change via readiness gate** — perubahan phase hanya diotorisasi Gate kesiapan: release-ready = (semua work item Done) ∧ (container Completed) ∧ (plan locked) ∧ (gate review) ∧ (gate persetujuan) | §11.2, §7.5 | — | Governance-only (uncovered) | Evaluasi agregat State lintas artifact; phase adalah context attribute, bukan state domain |
| REQ-015 | **Distillation before archive** — Session/Review wajib didistilasi ke artifact durable (ADR/decision) sebelum Archived | §11.4 | — | Governance-only (uncovered) | Aturan lifecycle; membutuhkan semantik relationship lintas tipe |
| REQ-016 | **Identity storage-independence** — Identity tidak disandikan dalam lokasi/penyimpanan; setiap Identity resolvable ke tepat satu artifact | §6.4, §12.2 | R2, R6 (sebagian) | Partially enforced | R2 (filename sebagai proyeksi) dan R6 (`dimension` == folder) meng-enforce permukaan serialisasi dari invariant ini; storage-independence sendiri tidak dapat dicek mekanis di dalam satu repositori |

---

## 4. Test coverage index

Setiap fungsi test **diklasifikasikan** tepat satu kali di indeks ini; fungsi integrasi (seperti `TestInvalidFixtures`) dapat dikutip ulang sebagai referensi silang pada baris rule yang diaturnya. Total: **54 test** (46 `conformance` + 8 `cmd/eka`).

### 4(a) Rule tests (unit, per rule)

| Rule | Test functions |
|---|---|
| R0 | `TestAnalyzeNoFrontmatterIsConventionDoc`, `TestAnalyzeFrontmatterWithoutTypeIDIsConventionDoc`, `TestAnalyzeUnterminatedFrontmatter`, `TestAnalyzeBrokenYAML`, `TestAnalyzeTypeXorID`, `TestAnalyzeValidArtifact`, `TestAnalyzeMissingIdentityFields`, `TestAnalyzeNonIntVersion`, `TestAnalyzeInvalidDate`, `TestAnalyzeUnknownType`, `TestAnalyzeChangeLogNotList`, `TestAnalyzeMalformedChangeLogEntry`, `TestUnknownTypeIsStructural` (13) |
| R1 | — (tidak ada unit test khusus; dicover via `TestInvalidFixtures` → `invalid-dup-identity`) |
| R2 | `TestRule2ExactCounts` (1) |
| R3 | `TestPhaseValidation` (1) |
| R4 | — (tidak ada unit test khusus; dicover via `TestInvalidFixtures` → `invalid-ownership`) |
| R5 | `TestRule5DraftSeverity`, `TestRule5CrossNamespaceAndVersionResolution`, `TestRule5VersionedReferenceToMissingInstance` (3) |
| R6 | `TestDimensionFolderResolution` (1) |
| R7 | `TestRule7ExactCounts` (1) |
| R8 | `TestRule8TicketHeaderAndDerivation`, `TestRule8ContainerTableMismatchIsWarningOnly` (2) |
| R9 | `TestRule9SupersededADRWithReplacement`, `TestRule9VersionedReplacementMustNameInstance` (2) |

### 4(b) Infrastructure tests (tidak terikat satu rule tertentu)

Diklasifikasikan eksplisit sebagai `infrastructure` — melindungi prasyarat pemindaian, model hasil, parsing grammar, tabel state, dan CLI.

| Grup | Test functions |
|---|---|
| Model hasil & determinisme (`report.go`/`validate.go`) | `TestReportCounts`, `TestReportPassSemantics`, `TestRelPathIsRootRelative`, `TestReportDeterminism`, `TestSortedResultsOrder` (5) |
| Kebijakan pemindaian (`validate.go`) | `TestScanSkipsTestdataAndHiddenDirs`, `TestValidateInputErrors`, `TestConventionDocumentsAreSkipped` (3) |
| Parsing filename (`filename.go`) | `TestParseFilename`, `TestParseFilenameEmpty` (2) |
| Parsing reference grammar (`rules.go:parseReference`) | `TestParseReference`, `TestParseReferenceCrossNamespace` (2) |
| Tabel state & taksonomi (`state.go`) | `TestTypeTokenCount`, `TestOwnedSets`, `TestContentStateVariant`, `TestDimensionTokens`, `TestIsLegalTransition`, `TestPhaseValueSet` (6) |
| Self-conformance | `TestReferenceImplementationConforms`, `TestFindRepoRoot` (2) |
| **CLI layer** (`cmd/eka/main_test.go`) | `TestExitCodeUsage`, `TestExitCodeBadPath`, `TestHelpExitsZero`, `TestValidateValidRepoExitsZero`, `TestValidateInvalidRepoExitsOne`, `TestWarningsDoNotAffectExitCode`, `TestDefaultPathIsCurrentDirectory`, `TestOutputIsDeterministic` (8) |

### 4(c) Fixture-based integration tests (`conformance/testdata/`)

| Test functions | Cakupan |
|---|---|
| `TestValidFixtureRepo` | Repo valid (`valid/`, 6 artifact) — seluruh rule lulus tanpa error |
| `TestInvalidFixtures` | 11 direktori skenario invalid: `invalid-malformed` → R0; `invalid-dup-identity` → R1; `invalid-filename` → R2; `invalid-state-value` → R3; `invalid-ownership` → R4; `invalid-reference` → R5; `invalid-dimension` → R6; `invalid-changelog` → R7; `invalid-projection` → R8; `invalid-sections` + `invalid-adr-superseded` → R9 |

---

## 5. Coverage analysis

### Covered

- **Rule coverage: 10/10** — R0 + R1–R9 seluruhnya `Enforced (tested)`: setiap rule memiliki implementasi Go yang dipanggil dari `conformance/validate.go:Validate` dan cakupan test (unit, infrastructure, atau fixture).
- **Test coverage: 54 test** (46 `conformance` + 8 `cmd/eka`), seluruhnya terpetakan di Section 4.
- **Self-conformance PASS** — `go run ./cmd/eka validate .` pada repositori ini: 7 artifact, 0 error, 0 warning, exit 0 (dikodifikasi sebagai `TestReferenceImplementationConforms`).

### Uncovered specification sections

Bagian spec yang hanya tercover sebagai requirement normatif (REQ-011..REQ-016, Section 3(b)), dengan status `Governance-only (uncovered)` atau `Partially enforced`:

| Req ID | Spec anchor | Status | Rekomendasi |
|---|---|---|---|
| REQ-011 (exactly-one-active) | §9, §7.5, §5.2 | Governance-only (uncovered) | Kandidat aturan masa depan — **follow-up, tidak diusulkan sekarang** |
| REQ-012 (lock-atomic-with-generation) | §5.2, §9 | Governance-only (uncovered) | Kandidat aturan masa depan — **follow-up** |
| REQ-013 (approved-content immutability) | P8, §5.4 | Governance-only (uncovered) | Kandidat aturan masa depan — **follow-up** |
| REQ-014 (phase readiness gate) | §11.2, §7.5 | Governance-only (uncovered) | Kandidat aturan masa depan — **follow-up** |
| REQ-015 (distillation before archive) | §11.4 | Governance-only (uncovered) | Kandidat aturan masa depan — **follow-up** |
| REQ-016 (identity storage-independence) | §6.4, §12.2 | Partially enforced (R2/R6) | Sisa permukaan tidak dapat dicek mekanis dalam satu repo — **follow-up** |

### Orphan implementations

**Tidak ada.** Verifikasi dengan membaca kode: `conformance/validate.go:Validate` memanggil `analyzeFile` (R0) pada fase parse dan `rule1`–`rule9` pada fase rule; seluruh fungsi pendukung (Section 2, kolom Implementation) direferensikan oleh rule yang melayaninya. `conformance/report.go` dan `cmd/eka/main.go` adalah infrastruktur engine/CLI, bukan rule.

### Orphan tests

**Tidak ada.** Seluruh 54 fungsi test terpetakan di Section 4 — setiap fungsi muncul tepat satu kali, tidak ada duplikasi dan tidak ada yang tidak terpetakan.

### Known gaps (dari `conformance-notes.md`)

| Gap | Detail | Referensi |
|---|---|---|
| R2 — id filename vs id frontmatter | Bagian id pada filename tidak dicocokkan dengan `id` frontmatter | [conformance-notes.md](conformance-notes.md) — interpretasi #17, tabel Gap |
| Exactly-one-active container | Invariant "tepat satu Execution Container aktif" tidak divalidasi | [conformance-notes.md](conformance-notes.md) — interpretasi #29, tabel Gap |
| R5 — referensi dua arah | Konvensi "referensi hanya pada artefak perujuk" tidak di-enforce | [conformance-notes.md](conformance-notes.md) — interpretasi #11, tabel Gap |

Ketiga gap tercermin pada kolom Notes Section 2 (R2, R5) dan status REQ-011 di Section 3(b). Gap tidak ditutup pada versi ini — perubahan perilaku validator adalah perubahan ruang lingkup, di luar governance v1.0.

---

## 6. Summary

| Metrik | Nilai |
|---|---|
| Total Requirements | **16** (10 enforced + 6 governance-only) |
| Total Conformance Rules | **10** (R0 + R1–R9) |
| Total Implementations | **10 rule implementations** (`analyzeFile` + `rule1`–`rule9`) + **20 fungsi pendukung** (`parseFilename`, `identityKey`, `buildIndex`, `contentStateVariant`, `domainValues`, `isLegalTransition`, `parseReference`, `resolve`, `dimensionFolderFor`, `dimensionList`, `entriesForDomain`, `indexOfEntry`, `workItemsTable`, `compareWorkItemsTable`, `resolveWorkItemCell`, `hasProjectionHeader`, `splitTableRow`, `requiredSectionsFor`, `headingMatches`, `hasReplacement`) + engine/report/CLI (`validate.go`, `report.go`, `cmd/eka/main.go`) |
| Total Test Suites | **2 paket**: `conformance` 46 test + `cmd/eka` 8 test = **54 test** |
| Coverage saat ini | **10/10 rule enforced & tested (100% rule coverage)**; requirement coverage = **10 enforced dari 16 total**; spec-section coverage — requirement enforced memetakan ke §3, §5, §6, §7, §8, §10, §13, §14; requirement governance-only memetakan ke §5, §7, §9, §11, §12 |
| Self-conformance | `eka validate .` = 7 artifact, 0 error, 0 warning, exit 0 |
| Gap teridentifikasi | REQ-011..REQ-016 tidak di-enforce (6 requirement governance-only) + 3 gap terdokumentasi (filename-id, exactly-one-active, referensi dua arah) |
| Follow-up yang direkomendasikan | (1) Kandidat aturan masa depan untuk REQ-011..REQ-016 — **tidak diusulkan sekarang**; (2) otomasi konsumsi matriks ini (parser deterministik atas struktur tabel) — **tidak sekarang**, matriks adalah living document markdown |
