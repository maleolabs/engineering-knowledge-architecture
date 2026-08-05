# Matriks Traceability — Elemen Repositori → Anchor EKA v1.0

Matriks lengkap: setiap elemen repositori (root README, seluruh file `standard/`, seluruh file `reference/`, 7 ADR, dan — mereferensikan zona sibling — setiap elemen `skeleton/docs/`) dipetakan ke anchor EKA (section/prinsip/taksonomi) beserta rasionalnya.

Legenda tipe: **artifact** = file ber-frontmatter `type` + `id` (artifact rule); **konvensi** = dokumen konvensi (tanpa `type`/`id`).

## Elemen zona Standard (A)

| Elemen | Zona | Tipe | Anchor EKA | Rasional |
|---|---|---|---|---|
| `standard/README.md` | standard | konvensi | 1.2, 4 | Menjelaskan zona pra-lapisan; standard mendefinisikan lapisan, bukan artifact proyek. |
| `standard/eka-specification-v1.0.md` | standard | konvensi (teks kanonik) | seluruh dokumen (1–16) | Salinan verbatim standard kanonik; referensi konformasi. |
| `standard/eka-exchange-specification-v1.0.md` | standard | konvensi (teks kanonik) | 13, 5.4 (invariant 7), P13, 16.1 (milestone 2) | Exchange Contract v1 (Ratified): mengoperasionalkan Section 13 + invariant 5.4.7 menjadi kontrak round-trip, idempotensi, referential integrity, schema versioning, conformance suite (R1–R9); refinement pass menambah Exchange Package Object Model (§4.4), Capability Declaration (§4.5), lokasi deklarasi tunggal, coverage rule lengkap; istilah subordinate didefinisikan di §4.1 tanpa mengamandemen teks kanonik. |
| `reference/ratification-notes-exchange-v1.0.md` | reference | konvensi | 16.1 (milestone 2), P13 | Ratification report Exchange Specification v1.0: ringkasan refinement, terminologi, model konseptual baru, konfirmasi kesiapan arsitektural. |
| `reference/eka-reference-serialization-format-v1.0.md` | reference | konvensi | Exchange Spec 4.4, 10, 15, 17.1; Naming 6.1, 7 | RSF v1.0 — serialisasi proyeksi kanonik Exchange Package Object Model (referensi, bukan normatif): Package Model, Unit Entry, Content Representation Model, Attachment Model, Manifest, prinsip serialisasi, round-trip mapping, compatibility, contoh konseptual, rekomendasi implementasi import/export. |
| `standard/eka-naming-and-terminology-specification-v1.0.md` | standard | konvensi (teks kanonik) | 3, 14.2, P14 | Meta-specification (Ratified): penamaan resmi ekosistem — product identity EKA, pola Family "EKA \<Family\> Specification", reference components, tooling (`eka` + subcommands verb), repository naming (`eka-<component>`), tabel istilah kanonik + deprecated, migrasi; tidak mengamandemen teks kanonik. |
| `standard/glossary.md` | standard | konvensi | 3 | Definisi eksak istilah Core Concepts; tidak diparafrase. |

## Elemen zona Reference (C)

| Elemen | Zona | Tipe | Anchor EKA | Rasional |
|---|---|---|---|---|
| `README.md` (root) | root | konvensi | 1.3, 16.2 | Identitas repositori: satu serialisasi, bukan arsitektur. |
| `reference/README.md` | reference | konvensi | 1.3 | Indeks meta-dokumentasi implementasi. |
| `reference/reference-architecture.md` | reference | konvensi | 1.3, 12.4, 6.4 | Menjelaskan serialisasi: zona→lapisan, konvensi, artifact rule. |
| `reference/migration-guide.md` | reference | konvensi | 16.2, 6.4 | Peta legacy→baru + strategi; konformasi serialisasi. |
| `reference/philosophy.md` | reference | konvensi | 1.1, 16.3 | Narasi dual-layer, protocol first-class, phase-as-context, single-writer. |
| `reference/terminology-glossary.md` | reference | konvensi | 3 | Istilah level implementasi; istilah kanonik ditautkan, tidak didefinisikan ulang. |
| `reference/breaking-changes.md` | reference | konvensi | 6.4, P9 | 14 perubahan breaking disengaja agar konsumen legacy tidak membaca Identity dari lokasi. |
| `reference/adr-summary.md` | reference | konvensi | 14.2.5 | Indeks governance keputusan implementasi (accepted). |
| `reference/traceability-matrix.md` | reference | konvensi | 12.2, 12.4 | Bukti konformasi: setiap elemen → anchor standard. |
| `reference/ratification-notes.md` | reference | konvensi | 16.1 (milestone 1) | Catatan ratifikasi EKA v1.0 (stabilization pass, verbatim). |

### Implementation ADR (`reference/decisions/`) — seluruhnya artifact (`type: adr`, `dimension: decisions`, status accepted)

| Elemen | Zona | Tipe | Anchor EKA | Rasional |
|---|---|---|---|---|
| `adr-001-identity-serialization.md` | reference/decisions | artifact | 6.2, 6.4, 6.3, P3, P9 | Identity di frontmatter; filename proyeksi; 26 token bebas-ambiguitas. |
| `adr-002-state-vector-encoding.md` | reference/decisions | artifact | 7.1, 7.2, 7.3, 7.4, 5.2, P2, P6, P7 | 5 field state frontmatter per domain owned; absence = not-applicable; change-log; mapping legacy. |
| `adr-003-projection-model.md` | reference/decisions | artifact | 7.4, 7.5, 9, 10, 15.5, P6, P9 | Ticket/tabel container = State Projection; State Vector kosong; refresh on-read. |
| `adr-004-phase-as-metadata.md` | reference/decisions | artifact | 3, 7.1, 7.5, 11.2, 11.3, P3 | Phase = field frontmatter pada scp-/plan-; phase change = context update via gate. |
| `adr-005-dimension-layout.md` | reference/decisions | artifact | 8, 4.1, 14.2, P1, P9, P15 | 12 folder = 12 dimensi 1:1 + operating/ + exchange/; `dimension == folder`. |
| `adr-006-exchange-conventions.md` | reference/decisions | artifact | 13.1, 13.2, 13.3, P13, P16 | Seam exchange: validation.md (9 aturan) + transfer.md (round-trip, konflik Identity, idempotensi, schema versioning). |
| `adr-007-extension-research-finding.md` | reference/decisions | artifact | 8, 10, 11.4, 14.1, 14.2, P12 | Ekstensi tipe `fnd-`: State Vector owned (Content, Existence) lengkap; exchangeable. |

### Elemen tooling — implementasi Go (validator EKA)

| Elemen | Zona | Tipe | Anchor EKA | Rasional |
|---|---|---|---|---|
| `go.mod` | root (tooling) | konvensi (build) | P16 | Modul `github.com/maleolabs/engineering-knowledge-architecture` (Go 1.24+); dependensi: `gopkg.in/yaml.v3`, `spf13/cobra` (CLI adapter), `golang.org/x/term` (deteksi TTY wizard); pintu masuk `go install`/`go build` CLI. |
| `cmd/` (command layer) | root (tooling) | konvensi (CLI) | 13.3, P16, Naming §7 | Definisi perintah Cobra murni: `root.go` (root command + `Execute(args, stdin, stdout, stderr) int` + mapping exit 0/1/2), `validate.go`, `init.go`; tanpa logika domain; help + completion standar Cobra. |
| `cmd/eka/` | root (tooling) | konvensi (CLI) | 13.3, P16 | Entry point tipis: `os.Exit(cmd.Execute(...))`; nama executable `eka`. |
| `bootstrap/` | root (tooling) | konvensi (engine) | Naming §7, Skeleton | Engine `eka init` (application layer, package publik): Workspace Discovery, Bootstrap Planning, Interactive Wizard (adaptif, non-interaktif deterministik via `x/term`), Repository Generation dari Reference Skeleton, Validasi pasca-generasi. |
| `skeletonembed.go` | root (tooling) | konvensi (embed) | Naming §6.1 | `//go:embed skeleton` — Reference Skeleton kanonik ter-embed untuk `eka init` (binari standalone, tanpa hardcode direktori). |
| `conformance/` | root (tooling) | konvensi (engine) | 13, P16 | Implementasi kanonik Conformance Rules (validation.md): engine publik reusable, independen dari CLI; entry `Validate(root) (*Report, error)`, `Scan(root) ([]Artifact, error)`, `ParseReference` (additif untuk konsumen exchange). |
| `exchange/` | root (tooling) | konvensi (engine) | Exchange Spec 4.4, 10, 12, 15; RSF | Engine ekspor (application layer, package publik): discovery, loading (via `conformance.Scan`), scope resolution (repo/line/instance/collection), external reference declaration, proyeksi RSF deterministik (header/manifest/unit/content/attachments/declarations/integrity), writer ZIP + direktori, guard charset identity (RSF §5.2.3, pencegahan path traversal). |
| `reference/cli.md` | reference | konvensi | 13, P16, Naming §7 | Dokumentasi CLI resmi: filosofi, instalasi, `eka init` (5 tahap, wizard adaptif, idempotensi, dry-run, validasi pasca-generasi), `eka validate`, exit codes, shell completion, arsitektur CLI (Cobra adapter + application layer), panduan kontribusi perintah baru, roadmap. |
| `reference/conformance-notes.md` | reference | konvensi | 13, P16 | Rekam traceability aturan R0–R9 → anchor EKA → lokasi implementasi + 29 keputusan interpretasi (kebijakan: didokumentasikan sebelum implementasi). |
| `.gitignore` | root (tooling) | konvensi (hygiene) | — | Hygiene implementasi: biner `eka` hasil build tidak pernah masuk VCS. |

## Elemen zona Skeleton (B) — zona sibling, direferensikan

| Elemen | Zona | Tipe | Anchor EKA | Rasional |
|---|---|---|---|---|
| `skeleton/docs/README.md` | skeleton | konvensi | 1.3 | Entry point serialisasi proyek; sumber kebenaran struktur. |
| `skeleton/docs/intent/` (`vis-`, `str-`) | skeleton | KB (folder) | 8 (Product Intent), 10 | Home visi/manifesto + strategi (tipe baru). |
| `skeleton/docs/requirements/` (`req-`) | skeleton | KB (folder) | 8 (Requirements), 10 | Requirement + amendmen sebagai instance dengan `amends`. |
| `skeleton/docs/architecture/` (`arc-`) | skeleton | KB (folder) | 8 (Architecture), 10 | Architecture Description. |
| `skeleton/docs/decisions/` (`adr-`, `dec-`) | skeleton | KB (folder) | 8 (Decisions), 10 | Dimensi Decisions tunggal: ADR + decision record. |
| `skeleton/docs/specifications/` (`spec-`) | skeleton | KB (folder) | 8 (Specifications), 10 | Dimensi baru; terpisah dari Vocabulary. |
| `skeleton/docs/standards/` (`std-`) | skeleton | KB (folder) | 8 (Standards & Guidelines), 10 | Konvensi, terpisah dari operasional. |
| `skeleton/docs/operations/` (`run-`) | skeleton | KB (folder) | 8 (Operational Knowledge), 10 | Prosedur saja. |
| `skeleton/docs/quality/` (`rvw-`) | skeleton | KB (folder) | 8 (Governance & Quality), 10 | Review + findings; `validates`. |
| `skeleton/docs/planning/` (`scp-`, `epc-`, `plan-`, `trc-`) | skeleton | KB (folder) | 8 (Planning Knowledge), 10, 11.2 | Scope (dengan `phase`), epic, plan (Planning State), traceability. |
| `skeleton/docs/records/` (`rel-`) | skeleton | KB (folder) | 8 (Records), 10 | Release record; immutable. |
| `skeleton/docs/research/` (`fnd-`) | skeleton | KB (folder) | 8 (Research), 11.4, 14.1 | Tipe ekstensi; jalur Distillation wajib. |
| `skeleton/docs/vocabulary/` (`gls-`) | skeleton | KB (folder) | 8 (Vocabulary), 10 | Glossary; bukan Specifications. |
| `skeleton/docs/operating/` | skeleton | OS (folder) | 4.1, 5.2 | Lapisan OS: state machine & protocol. |
| `skeleton/docs/operating/containers/` (`ctr-`) | skeleton | OS (folder) | 10, 7.2, 7.5, 9 | Execution Container; Container State; exactly-one-active. |
| `skeleton/docs/operating/work-items/` (`sto-`, `ts-`, `bug-`, `td-`, `ch-`, `spk-`) | skeleton | OS (folder) | 10, 7.2, 9 | Single-writer Execution State. |
| `skeleton/docs/operating/projections/` (`tkt-`) | skeleton | OS (folder) | 7.4, 10, P6 | Ticket: State Vector kosong, proyeksi. |
| `skeleton/docs/operating/sessions/` (`ses-`) | skeleton | OS (folder) | 10, 11.4 | Existence State; Distillation wajib sebelum Archived. |
| `skeleton/docs/operating/protocol.md` | skeleton | konvensi | 9, 11, 5.2 | Definisi protocol, urutan, gate, command. |
| `skeleton/docs/exchange/validation.md` | skeleton | konvensi | 13, P16 | 9 aturan validasi kepatuhan. |
| `skeleton/docs/exchange/transfer.md` | skeleton | konvensi | 13.2, P13 | Round-trip, kebijakan konflik Identity, idempotensi, schema versioning. |

## Konvensi serialisasi (aturan yang diimplementasikan)

| Konvensi | Anchor EKA | Rasional |
|---|---|---|
| Identity encoding (frontmatter + filename `<type-token>-<id>[-v<nn>]`) | 6.2, 6.4, P3 | Identity canonical, unambiguous, machine-parseable; decoupled dari lokasi. |
| State encoding (5 field per domain owned; absence = not-applicable) | 7.2, P6 | State eksplisit, single-writer; domain tidak berlaku tidak diserialisasi. |
| Phase sebagai metadata | 11.2 | Context attribute, bukan kategori/state; diotorisasi gate kesiapan. |
| Relationships by Identity (`supersedes`, `amends`, `derives-from`, `depends-on`, `validates`) | 6.2.7, 13.2.3 | Referensi tidak pernah by lokasi; referential integrity lintas sistem. |
| Klasifikasi `dimension == folder` | 8, P15 | Klasifikasi properti artifact; reklasifikasi tidak memutus referensi. |
| Proyeksi (ticket, tabel container) | 7.4, P6, P9 | Proyeksi bukan writer; refresh terhadap owner. |
| Artifact rule (`type` + `id` ⇒ Artifact) | 3, 5.1 | Pembeda artifact vs dokumen konvensi; identitas kepemilikan lapisan. |
| Change-log (`{date, domain, from, to, by}`) | 5.2 | Catatan wajib seluruh transisi state. |
| Exactly-one-active (container) | 9 | Konkurensi mutual exclusion Execution Container. |
| Lock-atomic-with-generation | 5.2 | Pembuatan container atomik dengan lock plan (Planning State → Immutable). |
| On-read refresh | 15.5 | Kebijakan default Projection Refresh; invariant proyeksi bukan writer absolut. |
| Exchange validation (9 aturan) | 13 | Validasi sebelum commit; kepatuhan terhadap standard. |
| Ekstensi `fnd-` (Research Finding) | 14.1 | Tipe baru via mekanisme ekstensi; State Vector owned lengkap (14.2.6). |
