# Glosarium Terminologi Implementasi

Glosarium istilah **level implementasi** — istilah yang lahir dari konvensi serialisasi repositori ini (Git + Markdown). Istilah yang didefinisikan standard (Artifact, State Vector, Identity, dst.) **tidak didefinisikan ulang di sini**; istilah tersebut hidup di [`standard/glossary.md`](../standard/glossary.md) dan hanya ditautkan.

## A

### artifact rule
Aturan untuk membedakan Artifact dari dokumen konvensi: **suatu file adalah Artifact iff frontmatter-nya memuat `type` DAN `id`**. Lihat `reference-architecture.md` §3. Konsep Artifact: [standard/glossary.md](../standard/glossary.md).

## D

### derived condition
Kondisi yang bukan nilai domain state, melainkan kondisi pemicu transisi (contoh: "Completed" pada container/session, diturunkan dari agregat state work item). Berlaku di EKA 7.2; dalam serialisasi, kondisi derived tidak pernah ditulis sebagai nilai state di frontmatter — ia dihitung.

### derived-from
Field relasi frontmatter (`derives-from`) yang menyatakan State Projection diturunkan dari artifact owner (contoh: `tkt-` → `ctr-` → work item). Referensi by Identity, bukan lokasi. Konsep Relationship: [standard/glossary.md](../standard/glossary.md).

### dimension == folder rule
Aturan lokasi: artifact knowledge harus hidup di folder Knowledge Dimension-nya (`dimension == folder`), ditegakkan validasi; artifact operating dikecualikan. Lihat ADR-005 dan EKA 8.

## E

### exactly-one-active
Konvensi konkurensi Operating Layer: paling banyak satu Execution Container (`ctr-`) berstatus `container-state: active` pada satu waktu; pembuatan berikutnya menunggu. Konsep: EKA 5.2, 9.

## F

### frontmatter
Blok metadata YAML di kepala file Markdown — lokasi serialisasi Identity, State Vector, Relationship, dan change-log. Frontmatter adalah satu-satunya tempat state ditulis (single-writer); body file adalah Content.

## I

### instance-version vs revision
Dua makna versi dalam frontmatter: `instance-version` adalah bagian Identity (instance baru = Identity instance baru, contoh `plan-x-v2`); `revision` melacak evolusi Content instance yang sama dan **bukan** bagian Identity. Konsep: EKA 6.3, [standard/glossary.md](../standard/glossary.md).

## L

### lock-atomic-with-generation
Invariant Operating Layer: peristiwa pembuatan Execution Container mengunci plan (Planning State → Immutable) dan membuat container secara atomik — tidak ada celah antara lock dan generation. Konsep: EKA 5.2, 9.

## N

### namespace
Field frontmatter yang memisahkan ruang Identity (produk/organisasi/sistem). Pada repositori ini artifact contoh menggunakan namespace proyek; dokumentasi meta menggunakan `eka-ref-impl`. Konsep: [standard/glossary.md](../standard/glossary.md).

## O

### on-read refresh
Kebijakan Projection Refresh default serialisasi ini: State Projection divalidasi terhadap owner **saat dibaca** (bukan event-driven). Invariant "proyeksi tidak pernah menjadi writer" tetap absolut. Konsep: EKA 5.5, 15.5, [standard/glossary.md](../standard/glossary.md).

## P

### projected state
State yang tidak dimiliki artifact (bukan bagian State Vector), melainkan diturunkan dari owner melalui Projection Semantics dan divalidasi via Projection Refresh. Serialisasi menandainya dengan header "Generated — State Projection" pada artifact generated. Konsep: [standard/glossary.md](../standard/glossary.md).

## S

### State Vector owned
Bagian State Vector yang **dimiliki** artifact (sesuai tipenya) dan diserialisasi sebagai field state di frontmatter; domain yang tidak dimiliki = absence (not-applicable). Proyeksi bukan bagian State Vector owned. Konsep: EKA 7.4, [standard/glossary.md](../standard/glossary.md).

## T

### type token
Awalan filename `<type-token>-<id>[-v<nn>]` yang menandai tipe artifact (26 token bebas-ambiguitas, lihat `reference-architecture.md` §2.1). Token adalah proyeksi Identity untuk navigasi manusia + validasi; Identity sejati ada di frontmatter.

## W

### well-formed content
Content yang mematuhi struktur per tipe artifact (didefinisikan skeleton per folder) sehingga dapat diparse dan dieksekusi secara deterministik. Konsep: EKA 3, [standard/glossary.md](../standard/glossary.md).

## Z

### zone
Pembagian tingkat atas repositori: `standard/` (teks kanonik, pra-lapisan), `skeleton/` (struktur proyek yang dapat disalin — serialisasi), `reference/` (meta-dokumentasi implementasi ini). Zona adalah konsep organisasi repositori, bukan konsep standard.
