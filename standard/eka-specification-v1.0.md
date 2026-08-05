# Engineering Knowledge Architecture (EKA) — Canonical Specification

| Field | Value |
|---|---|
| **Version** | 1.0 |
| **Status** | Ratified |
| **Otoritas** | Standar kanonik model pengetahuan engineering |
| **Cakupan** | Konsep, invariant, dan kontrak — bukan mekanisme implementasi |

**Pembacaan dokumen ini:** istilah dengan huruf kapital awal merujuk definisi di Section 3. Kata "harus/wajib" menandakan requirement yang mengikat seluruh implementasi. Kata "dapat/boleh" menandakan pilihan implementasi dalam batas kontrak.

---

## 1. Engineering Knowledge Architecture Overview

### 1.1 Definisi standard

EKA adalah model konseptual kanonik untuk pengetahuan engineering: definisi Artifact, Identity, State, taksonomi pengetahuan, arsitektur lapisan, dan kontrak pertukaran antar sistem. Standard ini lahir dari dua tanggung jawab yang terbukti hidup berdampingan pada implementasi awal: **Knowledge Base** (penyimpanan dan preservasi pengetahuan) dan **Engineering Operating System** (eksekusi deterministik, koordinasi agent, governance). Keduanya adalah dua lapisan dari satu sistem, diikat oleh Identity artifact.

### 1.2 Ruang lingkup

Standard ini menetapkan:
- Konsep fundamental (Section 3) dan prinsip (Section 2);
- Arsitektur lapisan dan kontraknya (Sections 4–5);
- Identity Model (Section 6) dan State Taxonomy (Section 7);
- Taksonomi pengetahuan, eksekusi, dan artifact (Sections 8–10);
- Model lifecycle konseptual (Section 11);
- Kontrak penyimpanan-independen dan pertukaran (Sections 12–13);
- Model ekstensi dan evolusi (Sections 14–16).

Standard ini **tidak** menetapkan: format serialisasi, layout penyimpanan, struktur direktori, template dokumen, skema penamaan, bahasa query, atau mekanisme enforcement spesifik. Seluruhnya adalah keputusan implementasi yang wajib mematuhi kontrak di dokumen ini.

### 1.3 Hubungan dengan implementasi

**Satu implementasi adalah satu serialisasi dari arsitektur ini — bukan arsitekturnya.** Implementasi dapat berupa: repositori git, database relasional, graph database, object store, AI-native knowledge store, Knowledge OS, atau pipeline import/export. Setiap implementasi wajib: (a) memenuhi seluruh invariant (Section 5.4), (b) menyediakan resolusi Identity, (c) mendukung exchange lossless melalui kontrak Exchange Layer (Section 13).

### 1.4 Posisi jangka panjang

EKA adalah model pengetahuan engineering kanonik yang dapat diintegrasikan dengan Knowledge OS masa depan melalui mekanisme import/export terstandar. Arsitektur bertahan terhadap pergantian media penyimpanan karena dibangun di atas konsep — Identity, State, Layer, kontrak — bukan di atas mekanisme.

---

## 2. Architectural Principles

| # | Prinsip | Rasional |
|---|---|---|
| P1 | **Separation of Concerns** | Tanggung jawab pengetahuan dan eksekusi dipisahkan menjadi lapisan dengan kontrak eksplisit. Fusion keduanya dalam satu hierarki adalah sumber konflasi. |
| P2 | **Explicit State** | State selalu eksplisit sebagai metadata ber-owner — tidak pernah implisit di struktur. Struktur boleh menjadi proyeksi State, tetapi bukan fakta independen. |
| P3 | **Stable Identity** | Identity immutable dan independen dari lokasi, penyimpanan, State, dan klasifikasi. Referensi selalu memakai Identity, tidak pernah lokasi. |
| P4 | **Protocol vs Content Distinction** | Protocol adalah properti Operating Layer; Content adalah properti Knowledge Layer. Setiap artifact yang melayani keduanya mendefinisikan keduanya secara terpisah. |
| P5 | **Layer Independence** | Setiap lapisan dapat berevolusi tanpa merombak yang lain: taksonomi berubah tanpa menyentuh protocol; protocol diperkuat tanpa menggeser home pengetahuan. |
| P6 | **Single Writer** | Setiap field State memiliki tepat satu owner. Tampilan lain adalah State Projection — di-generate atau divalidasi, tidak pernah diedit sebagai fakta independen. |
| P7 | **Forward-Only Transitions** | Seluruh State Domain bergerak maju tanpa regresi. Koreksi dilakukan dengan instance baru + Relationship, bukan mutasi. |
| P8 | **Approved-Content Immutability** | Content yang telah melewati gate persetujuan tidak dimutasi diam-diam. Perubahan hanya melalui kanal governance. Prasyarat preservasi. |
| P9 | **Structure as Projection of State** | Organisasi/posisi struktural diturunkan dari State dan Identity — struktur tidak pernah menjadi fakta kedua yang dapat melenceng dari State. |
| P10 | **Two Change Channels** | Kanal Content (governance) dan kanal State (protocol) terpisah. Mencampur keduanya adalah pelanggaran. |
| P11 | **Determinism by Protocol** | Urutan eksekusi didefinisikan protocol: "apa selanjutnya" selalu terjawab. Enforcement adalah kapabilitas implementasi; requirement-nya adalah standard. |
| P12 | **Preservation Over Deletion** | Sejarah adalah pengetahuan. Keputusan yang salah pun dipertahankan. Superseded/Archived adalah Record, bukan sampah. |
| P13 | **Lossless Exchange** | Pertukaran antar sistem tidak boleh kehilangan atau menduplikasi Identity, State, Content, atau Relationship. |
| P14 | **Minimum Canonical Core** | Standard menetapkan konsep dan kontrak; implementasi memilih mekanisme. Semakin kecil inti, semakin besar umur standard. |
| P15 | **Classification is Property, Not Identity** | Knowledge Dimension adalah properti artifact; perubahan klasifikasi tidak pernah memutus referensi. |
| P16 | **Enforcement Capability Varies, Invariants Don't** | Mekanisme enforcement berbeda antar implementasi (constraint struktural, constraint database, validasi); invariant yang harus ditegakkan identik. |

---

## 3. Core Concepts

Definisi presisi — setiap istilah dipakai dengan makna tunggal di seluruh dokumen ini.

- **Artifact** — entitas pengetahuan engineering yang memiliki Identity, Content, State Vector (domain State yang **dimiliki**), dan Relationship. Unit dasar model.
- **Content** — muatan semantik Artifact: intent, keputusan, desain, kendala, prosedur, catatan. Milik Knowledge Layer.
- **Well-formed Content** — Content yang mematuhi struktur yang ditetapkan untuk tipe Artifact-nya, sehingga dapat diparse dan dieksekusi secara deterministik.
- **Identity** — properti Artifact yang membedakannya dari semua Artifact lain secara permanen: `(Namespace, Type, ID[, InstanceVersion])`. Lihat Section 6.
- **Artifact Line** — entitas Identity yang bertahan: satu `(Namespace, Type, ID)`.
- **Artifact Instance** — satu versi eksistensi dari sebuah Line: Line + `InstanceVersion`.
- **State** — fakta tentang posisi Artifact dalam proses tertentu.
- **State Vector** — tuple State Domain **yang dimiliki** oleh sebuah Artifact. State yang diproyeksikan (State Projection) bukan bagian State Vector.
- **State Domain** — dimensi State independen dengan semantik, owner, dan aturan transisinya sendiri (Section 7). Domain bersifat orthogonal.
- **State Projection** — tampilan State yang diturunkan dari owner (contoh: agregat State work item → status Execution Container). Proyeksi tidak memiliki State sendiri; proyeksi tidak pernah menjadi writer.
- **Projection Semantics** — aturan komputasi State Projection: dari owner mana, agregasi bagaimana, apa yang ditampilkan.
- **Projection Refresh** — mekanisme dan waktu validasi State Projection terhadap owner (on-read dan/atau event).
- **Knowledge Dimension** — sumbu klasifikasi pengetahuan (Section 8). Properti Artifact, bukan Identity.
- **Protocol** — aturan deterministik milik Operating Layer: urutan, transisi State, kuncian, gate, perintah eksekusi.
- **Layer** — lapisan arsitektur dengan kontrak eksplisit: Knowledge Layer, Operating Layer, Exchange Layer (Section 4).
- **Namespace** — ruang Identity yang memisahkan domain pengelolaan (produk, organisasi, sistem).
- **Relationship** — relasi eksplisit antar Artifact yang direferensikan by Identity: supersedes, amends, derives-from, depends-on, validates.
- **Gate** — kondisi yang harus dipenuhi sebelum transisi atau eksekusi boleh terjadi (gate persetujuan, gate kesiapan, gate review).
- **Command** — instruksi eksekusi deterministik yang dikonsumsi oleh eksekutor (manusia atau agent). Content Command adalah Content (milik Knowledge Layer); eksekusinya diatur Protocol (Operating Layer).
- **Execution Container** — Artifact eksekusi yang membungkus work item dan membawa konvensi konkurensi (exactly-one-active). Domain State-nya: Container State. Contoh implementasi: sprint.
- **Phase** — konteks produk/scope sepanjang waktu (Discovery, MVP, Milestone, Release). Phase adalah **context attribute** pada Artifact planning/scope: bukan kategori, bukan State Domain. **Phase change** adalah context update yang diotorisasi Gate kesiapan (Sections 7.5, 11.2).
- **Record** — Artifact yang dipreservasi sebagai sejarah (Superseded, Archived, Retired, release record) dengan Content immutable.
- **Distillation** — transformasi pengetahuan ephemeral (konteks kerja, temuan review) menjadi pengetahuan durable (keputusan, ADR, Record).
- **Change Log** — catatan kronologis transisi State pada Artifact: domain, nilai lama, nilai baru, waktu, otoritas. Wajib untuk seluruh State Domain.
- **Identity Registry** — fungsi Knowledge Layer yang menjamin keunikan Identity, resolusi Identity ke Artifact, dan integritas referensi.
- **Trigger** — hubungan antar domain/operasi di mana suatu peristiwa memicu validasi atau transisi pada domain lain (definisi formal interaksi, Section 7.5).
- **Knowledge OS** — platform eksekusi pengetahuan masa depan yang mengonsumsi dan memproduksi Artifact EKA melalui Exchange Layer. Bukan bagian standard; konsumen standard.

---

## 4. Layer Model

### 4.1 Tiga lapisan

| Layer | Peran | Memiliki | Tidak memiliki |
|---|---|---|---|
| **Knowledge Layer (KB)** | Store pengetahuan: Content, klasifikasi, preservasi, referensi | Content, klasifikasi (Knowledge Dimensions), Relationship, history/Records, administrasi Identity (Identity Registry) | State proses, Protocol eksekusi |
| **Operating Layer (OS)** | State machine & Protocol eksekusi | State Domain (Execution, Planning, Container, Existence), urutan, konkurensi, kuncian, Gate, Command | Content (tidak pernah mengedit Content) |
| **Exchange Layer (EX)** | Batas transformasional: serialisasi, validasi, import/export, mediasi sistem eksternal | Kontrak pertukaran, aturan round-trip, validasi kepatuhan | Content dan State (tidak pernah menjadi owner) |

### 4.2 Keputusan: Exchange Layer sebagai lapisan ketiga

Exchange Layer wajib ada sebagai lapisan terpisah karena:

1. **Invariant berbeda**: exchange memiliki invariant sendiri (round-trip lossless, idempotensi, referential integrity) yang tidak dimiliki KB maupun OS.
2. **Arah interaksi berbeda**: KB dan OS berinteraksi secara internal; EX berinteraksi dengan sistem eksternal — mediasi membutuhkan kontrak batas eksplisit.
3. **Visi Knowledge OS**: integrasi masa depan membutuhkan seam yang didefinisikan di level standard, bukan di level implementasi.

EX tidak memiliki Content maupun State — ia adalah lapisan batas yang memvalidasi dan mentransformasikan representasi.

### 4.3 Independensi lapisan

- KB dapat mengubah taksonomi tanpa mengubah Protocol OS (P5).
- OS dapat menambah protocol variant tanpa mengubah klasifikasi KB.
- EX dapat menambah format serialisasi tanpa mengubah KB/OS.
- Ketiganya terikat oleh: **Identity** (Section 6), **invariant global** (Section 5.4), dan **kontrak antarlapisan** (Section 5.3).

---

## 5. Layer Contracts

### 5.1 Knowledge Layer — kontrak

- **Responsibilities**: menyimpan Content; mengklasifikasikan Artifact ke Knowledge Dimensions; memelihara cross-reference (referential integrity); mempreservasi history (P8, P12); menyediakan semantik retrieval; mengadministrasikan Identity (Identity Registry).
- **Ownership**: Content; klasifikasi; Relationship records; history; Identity.
- **Allowed interactions**: melayani pembacaan Content ke OS dan EX; menerima penulisan State dari OS (sebagai metadata pada Artifact); menerima Artifact baru dari OS melalui creation protocol (Identity diberikan oleh Identity Registry); menjalankan transisi Content State melalui Gate persetujuan.
- **Invariants**: Content Approved/Immutable tidak dimutasi (P8); referensi selalu valid (tidak ada dangling reference); klasifikasi berubah tanpa mengubah Identity (P15); Identity immutable (P3).
- **Synchronization boundaries**: hanya Content State yang dimiliki di sini. State Domain lain hanya direfleksikan — KB membaca State untuk query, tidak menulisnya.
- **Extension points**: Knowledge Dimension baru; tipe Artifact baru; Relationship type baru (Section 14).

### 5.2 Operating Layer — kontrak

- **Responsibilities**: menjalankan Protocol eksekusi (urutan, transisi State, Gate); mengelola konkurensi (exactly-one-active); mengelola kuncian/immutability (lock sebelum konsumsi); mendefinisikan Command; menyediakan what-next discoverability.
- **Ownership**: Execution State, Planning State, Container State, transisi Existence State; Projection Semantics (container view, ticket view).
- **Allowed interactions**: membaca Content dari KB (eksekusi: Command membaca instruksi; Execution Container di-generate dari plan); menulis State ke Artifact (single-writer per field, P6); membuat Artifact baru (Execution Container, ticket, session) melalui creation protocol dengan Identity dari Identity Registry; memvalidasi State Projection terhadap owner.
- **Invariants**: transisi forward-only (P7); satu writer per State field (P6); exactly-one-active Execution Container; **lock-atomic-with-generation** — peristiwa pembuatan Execution Container mengunci plan dan membuat container secara atomik; setiap transisi tercatat di Change Log; OS tidak pernah mengubah Content — hanya State.
- **Synchronization boundaries**: State Domain dimiliki di sini; State Projection di-generate/divalidasi di sini; Content hanya dibaca.
- **Extension points**: protocol variant baru, Command type baru, Gate type baru (Section 14).

### 5.3 Kontrak antarlapisan — titik kopling

| Kopling | Arah | Konten kontrak |
|---|---|---|
| **OS membaca Content** | KB → OS | Content Artifact yang dieksekusi harus Well-formed dan deterministik; Identity resolvable. |
| **OS menulis State** | OS → KB | State ditulis sebagai metadata pada Artifact; tidak mengubah Content; tercatat di Change Log. |
| **OS memproduksi Artifact** | OS → KB | Creation protocol: Identity diminta dari Identity Registry; Content Artifact baru valid menurut taksonomi. |
| **KB menjaga konsistensi Content dengan State** | KB ↔ OS | Content plan yang locked tidak boleh berubah (kuncian mem-gate kanal Content); Gate State mem-bolehkan/melarang transisi Content. |
| **EX memvalidasi & mentransfer** | EX ↔ KB/OS | Import/export melewati kontrak standard; EX memvalidasi kepatuhan sebelum commit; tidak pernah menulis State atau Content secara langsung. |

### 5.4 Invariant global

1. **Identity immutable** — Identity tidak berubah oleh State, lokasi, klasifikasi, atau revisi Content (P3).
2. **Single owner per State field** — satu penulis; semua tampilan lain adalah State Projection (P6).
3. **Structure as projection of State** — posisi/representasi struktural diturunkan dari State; tidak pernah fakta independen (P9).
4. **State changes only via Protocol** — tidak ada jalur perubahan State selain transisi Protocol (P11).
5. **Two separate change channels** — kanal Content (governance) dan kanal State (Protocol) tidak pernah bercampur (P10).
6. **Approved Content immutable** — preservasi (P8, P12).
7. **Round-trip lossless** — exchange tidak kehilangan/menduplikasi apa pun (P13).

### 5.5 Semantik sinkronisasi

- **Projection Refresh**: State Projection divalidasi terhadap owner pada dua titik — saat dibaca (on-read) dan saat owner berubah (event). Kebijakan preferensi antara keduanya adalah open question (Section 15); invariant "proyeksi tidak pernah menjadi writer" bersifat absolut.
- **Trigger**: peristiwa tertentu (transisi State, operasi OS) memicu validasi atau transisi di domain lain. Trigger didefinisikan per pasangan domain (Section 7.5).

---

## 6. Identity Model

### 6.1 Komposisi Identity

**Identity instance** = `(Namespace, Type, ID, InstanceVersion)`.
**Identity line** = `(Namespace, Type, ID)`.

| Properti | Bagian dari Identity? | Alasan |
|---|---|---|
| **Namespace** | **Ya** | Memisahkan ruang Identity (produk, organisasi, sistem). Dua Artifact dengan ID sama di Namespace berbeda adalah dua Artifact. |
| **Artifact Type** | **Ya** | Type menentukan State Domain yang berlaku dan Protocol yang mengikat. Type adalah kualifikasi Identity, bukan klasifikasi. |
| **ID** | **Ya** | Token unik dalam `(Namespace, Type)`. |
| **InstanceVersion** | **Ya (diskriminator instance)** | Membedakan instance dalam satu Line (plan v1 vs v2). Lihat 6.3. |
| **Knowledge Dimension (klasifikasi)** | **Tidak** | Properti retrieval. Reklasifikasi tidak boleh memutus referensi (P15). |
| **Location / organisasi** | **Tidak** | Lokasi adalah proyeksi; Identity independen dari lokasi (P9). |
| **Storage backend** | **Tidak** | Identity bertahan lintas penyimpanan (Section 12). |
| **State (semua domain)** | **Tidak** | State berubah; Identity tidak pernah (P3). |
| **Revision (riwayat edit Content)** | **Tidak** | Revision melacak evolusi Content instance yang sama; mengubah revision tidak mengubah Identity. |

### 6.2 Aturan Identity

1. Identity ditetapkan sekali saat creation, tidak pernah diubah.
2. ID unik dalam `(Namespace, Type)`; InstanceVersion unik dalam Line.
3. Referensi selalu by Identity — tidak pernah by lokasi, nama tampilan, atau klasifikasi.
4. Dua Artifact sama iff Identity-nya sama; Relationship tidak pernah mengubah Identity.
5. Type Artifact menentukan State Vector yang berlaku (Section 10) — binding type→state adalah bagian standard, bukan pilihan implementasi.
6. Identity harus dapat diserialisasi secara **lossless, unambiguous, dan machine-parseable** di semua implementasi. Mekanisme serialisasi adalah keputusan implementasi.
7. Supersession, amendment, dan derivation diekspresikan sebagai **Relationship antar Identity**, bukan perubahan Identity.

### 6.3 Semantik versi

Dua makna versi, dua jawaban:

- **InstanceVersion** (versi yang menunjuk instance berbeda) — **bagian dari Identity instance**. Instance baru diciptakan dengan sengaja (contoh: plan v2 setelah v1 terkunci). Perubahan instance = perubahan Identity instance; Line tetap.
- **Revision** (pelacakan evolusi Content instance yang sama) — **bukan bagian dari Identity**. Revision berubah setiap edit dan tidak boleh memutus referensi.

Aturan: **Line Identity tidak pernah berubah; Instance Identity berubah hanya saat instance baru sengaja diciptakan; Revision tidak pernah menyentuh Identity.** Supersession adalah Relationship antar dua Line — bukan pergantian Identity.

### 6.4 Studi kasus pelanggaran: kolisi ruang Identity

Implementasi awal menyandikan empat tipe Artifact (scope definition, plan, Execution Container, ticket) dalam satu ruang ID bersama dengan prefiks yang sama, sehingga Type tidak dapat dibedakan secara deterministik dari representasi Identity-nya: sebuah representasi dapat dibaca sebagai scope definition maupun plan. Analisis pelanggaran:

- Aturan 6.2.1–6.2.2 dilanggar: Type tidak tegas → ID tidak unik per `(Namespace, Type)`.
- Aturan 6.2.3 dilanggar secara konseptual: Identity disandikan melalui lokasi dan konvensi representasi, bukan properti Artifact.
- Akar masalah: Identity diturunkan dari struktur, bukan dari properti Artifact; dan tahap proses (pipeline) ikut menjadi bagian representasi Identity.

Pelajaran mengikat untuk semua implementasi: **Identity tidak boleh disandikan dalam lokasi, tahap proses, atau konvensi representasi.** Identity adalah properti pertama-class yang ditetapkan Identity Registry.

---

## 7. State Taxonomy

### 7.1 Evaluasi kandidat domain

| Kandidat | Keputusan | Alasan |
|---|---|---|
| **Artifact State (unified)** | **Tolak** | Monolit State adalah akar masalah duplikasi status: satu Artifact punya banyak dimensi State independen. State selalu berupa State Vector. |
| **Execution State** | **Terima** | Progress work item melalui Protocol. |
| **Lifecycle State (generic)** | **Tolak sebagai domain tunggal** | "Lifecycle" adalah komposisi State Domain sepanjang waktu; Phase produk adalah context, bukan State Domain. Yang bertahan: Existence State (domain) + Phase (context). |
| **Governance/Content State** | **Terima** (dinamai Content State) | Kematangan Content melalui gate persetujuan. |
| **Review State** | **Tolak sebagai domain** | Review adalah Gate pada transisi (stage Review di Content State; nilai "In Review" di Execution State) + tipe Artifact (review record). Tidak ada State independen yang tersisa setelah pemisahan itu. |
| **Release State** | **Tolak sebagai domain** | Release adalah Phase context + Gate kesiapan yang dievaluasi atas agregat State Domain lain. |
| **Planning State** | **Terima** | Commitment level plan: Draft → Approved → Immutable. |
| **Container State** | **Terima** | Execution Container terbuka/tertutup + konkurensi. |
| **Existence State** | **Terima** | Presensi Artifact: aktif vs dipreservasi. |

### 7.2 Domain formal

Aturan umum domain: setiap domain memiliki tepat satu **initial state** dan tepat satu **terminal state** (State tanpa transisi keluar); seluruh transisi forward-only (P7); koreksi dilakukan dengan instance baru + Relationship, bukan regresi. **Nilai domain hanya mencakup State yang dimiliki (owned); kondisi derived — misalnya "Completed" pada konteks kerja — bukan nilai domain, melainkan kondisi pemicu transisi.**

| Domain | Nilai | Responsibility | Owner | Aturan transisi |
|---|---|---|---|---|
| **Content State** | Draft → Review → Approved; terminal pasca-persetujuan: Amended \| Superseded | Kematangan Content sebagai pengetahuan; kanal governance | Knowledge Layer (gate: owner/approver) | Forward-only; gate persetujuan; perubahan pasca-Approved hanya via amendmen/supersesi (P8). **Varian**: standard (terminal Amended), decision record (stage Approval bernama Accepted; supersesi opsional), ADR (stage Approval bernama Accepted; supersesi wajib menunjuk pengganti). Varian boleh mengganti nama stage dan meniadakan stage opsional, tetapi wajib mempertahankan posisi semantik: pra-persetujuan → gate persetujuan → terminal pasca-persetujuan. |
| **Execution State** | Planned → Todo → In Progress → In Review → Done | Progress work item melalui Protocol | Operating Layer (single-writer: artifact work item) | Strictly sequential; never skip; never revert; satu initial; satu terminal; tiap transisi tercatat di Change Log. |
| **Planning State** | Draft → Approved → Immutable | Commitment plan: tentative → committed → locked | Operating Layer | Forward-only; Approved = siap dieksekusi; **Immutable dicapai atomik dengan peristiwa pembuatan Execution Container** (lock-atomic-with-generation); perubahan pasca-lock = instance baru (InstanceVersion). |
| **Container State** | Active → Completed | Execution Container terbuka/tertutup; konkurensi | Operating Layer | **Completed adalah transisi derived**: dipicu agregat Execution State (semua work item Done); exactly-one-Active (mutual exclusion). |
| **Existence State** | Active → Archived → Retired | Presensi Artifact dalam proses aktif vs preservasi | Operating Layer (transisi); Knowledge Layer (prinsip preservasi, P12) | Forward-only; **Archived** = reference-only (tanpa transisi State lain, tanpa mutasi Content, tetap tersedia di retrieval); **Retired** = preservasi terminal (tidak disurface di retrieval normal, Content immutable, Identity tidak berubah). Berlaku untuk semua tipe Artifact. |

### 7.3 Pemetaan mesin status implementasi awal → domain

| Mesin status implementasi awal | Domain |
|---|---|
| Living documents Draft→Review→Approved→Amended | Content State (varian standard) |
| ADR Proposed→Accepted→Superseded | Content State (varian ADR) |
| Decisions Draft→Accepted→Superseded (opsional) | Content State (varian decision) |
| Roadmap Draft→Approved→Immutable | Planning State |
| Sprints Active→Completed | Container State |
| Sessions Active→Completed→Archived | Existence State (Completed = kondisi derived dari work item yang direferensikan) |
| Work items Planned→Todo→In Progress→In Review→Done | Execution State |

Bukti empiris independensi domain: implementasi awal membutuhkan tujuh mesin status berbeda karena satu mesin tidak dapat mengekspresikan tujuh semantik. State Taxonomy mengorganisirnya menjadi domain dengan aturan eksplisit.

### 7.4 State Vector

Setiap Artifact membawa **State Vector** = tuple State Domain yang **dimiliki** (owned) sesuai tipenya; domain yang tidak berlaku ditandai not-applicable. Contoh: work item = `(Execution State, Existence State)` — Content State tidak berlaku; plan = `(Content State, Planning State, Existence State)`; Execution Container = `(Container State, Existence State)`; ADR = `(Content State, Existence State)`.

**Artifact yang seluruh statenya diproyeksikan memiliki State Vector kosong** (contoh: ticket = `(∅)`; State-nya adalah State Projection atas work item yang direferensikan). Ini menyelesaikan duplikasi status secara formal: status pada representasi turunan (container view, ticket view) adalah State Projection dari owner, divalidasi melalui Projection Refresh — menggantikan invariant "harus selalu setuju" tanpa penulis.

### 7.5 Interaksi antar domain

| Interaksi | Sumber → Target | Semantik |
|---|---|---|
| Pembuatan Execution Container (generation event) | Operasi Operating Layer → Planning State | Peristiwa pembuatan container dari plan memicu transisi plan ke Immutable; atomik dengan pembuatan (lock-atomic-with-generation). |
| Agregat State work item | Execution State → Container State | Semua work item dalam container Done memicu container Completed (transisi derived). |
| Lock plan → kanal Content | Planning State → Content State | Plan Immutable mem-gate perubahan Content pada instance itu (kanal governance terkunci). |
| Content readiness → commitment | Content State → Planning State | Gate Approval pada plan mensyaratkan Content siap (Approved di Planning State mengkomposisi kematangan Content). |
| Container Completed → preservasi | Container State → Existence State | Container selesai dapat di-archive. |
| Gate kesiapan release | Execution + Planning + Container (+ Gate review, Gate persetujuan) → **phase change** | Evaluasi agregat State untuk mengotorisasi perubahan Phase (context update; bukan state transition). Lihat 11.2. |

### 7.6 Independensi vs unifikasi — keputusan

State Domain **tetap independen** (tidak diunifikasi). Alasan: (1) bukti empiris implementasi awal — semantik berbeda membutuhkan mesin berbeda; (2) owner berbeda (KB vs OS); (3) laju perubahan berbeda (Content berubah lambat via governance; State eksekusi berubah cepat via Protocol); (4) unifikasi mengembalikan masalah duplikasi status. Interaksi dikelola eksplisit via Trigger dan Gate (7.5), bukan dengan menggabungkan domain.

---

## 8. Knowledge Taxonomy

Dimensi klasifikasi Knowledge Layer. Satu Artifact memiliki **satu dimensi primer** + dimensi sekunder opsional. Klasifikasi adalah properti retrieval — reklasifikasi tidak mengubah Identity (P15).

| Dimensi | Isi | Stabilitas | Catatan |
|---|---|---|---|
| **Product Intent** | Visi, strategi, prinsip, non-negotiables | Sangat stabil | Strategi adalah dimensi yang sebelumnya tidak terwakili secara eksplisit. |
| **Requirements** | Dokumen requirement — "apa yang harus dibangun dan mengapa" — beserta amendmen | Stabil (berevolusi via governance) | |
| **Architecture** | Deskripsi sistem, domain models, batas komponen | Stabil | |
| **Decisions** | Keputusan ireversibel (ADR) dan reversibel (decision log); termasuk keputusan produk | Stabil, akumulatif | |
| **Specifications** | Spesifikasi fungsional, non-fungsional (NFR), API, data | Stabil saat Approved | Dimensi yang sebelumnya tidak terwakili secara eksplisit; berbeda dari Vocabulary. |
| **Standards & Guidelines** | Standar engineering, konvensi, definisi done | Stabil, governance tinggi | Sebelumnya tercampur dengan pengetahuan operasional. |
| **Operational Knowledge** | Runbook, deployment, migration, checklist | Medium (berubah per environment) | Terpisah dari Standards & Guidelines. |
| **Governance & Quality** | Review findings, audit, quality gates | Akumulatif | Quality sebelumnya hanya sebagai peran, bukan dimensi. |
| **Planning Knowledge** | Content plan (roadmap, milestone definition) dan artifact relasi (traceability) | Medium (commitment) | Traceability adalah Relationship artifact. |
| **Records** | Release record, change log, snapshot historis | Immutable | Release record sebelumnya tidak terwakili. |
| **Research** | Temuan investigasi, hasil riset teknis | Akumulatif | Wajib jalur Distillation ke dimensi durable. |
| **Vocabulary** | Glossary, istilah kanonik, model lifecycle | Sangat stabil | Bukan Specifications — pemisahan ini wajib dijaga. |

---

## 9. Execution Taxonomy

Klasifikasi elemen Protocol Operating Layer:

| Elemen | Definisi | Responsibility | Invariant |
|---|---|---|---|
| **Ordering (chain)** | Urutan hubungan stage: requirement → scope → capability → plan → container → work item → konteks kerja → validasi | Menjawab "apa selanjutnya" secara deterministik; Protocol, bukan properti lokasi | Urutan didefinisikan eksplisit; eksekusi mengikuti urutan |
| **State Transitions** | Aturan perpindahan nilai dalam State Domain | Menjalankan Protocol per domain (Section 7.2) | Forward-only; never skip; never revert; tercatat di Change Log |
| **Concurrency Control** | Kuncian mutual exclusion atas Execution Container | Exactly-one-active container | Satu container aktif; pembuatan berikutnya menunggu |
| **Versioning / Immutability** | Lock plan saat eksekusi dimulai; perubahan = instance baru | Konsistensi rencana vs eksekusi | Lock-atomic-with-generation; Content locked tidak berubah |
| **Gates** | Kondisi sebelum transisi/eksekusi: gate persetujuan, gate kesiapan, gate review | Mengontrol kapan transisi sah | Gate dievaluasi atas State owner, bukan proyeksi |
| **Commands** | Instruksi eksekusi deterministik yang dikonsumsi eksekutor | Menerjemahkan Content Artifact menjadi aksi eksekusi | Content Command Well-formed; output deterministik |
| **Agent Coordination** | Semantik interaksi agent dengan sistem | Identity parseable, State eksplisit, Content deterministik, urutan tegas | Agent membaca State/Identity tanpa ambiguitas |
| **Execution Containers** | Execution window yang membungkus work item | Agregasi, konkurensi, snapshot | State container derived dari work item |
| **Projection Semantics** | Aturan komputasi State Projection (container view, ticket view): dari owner mana, agregasi bagaimana | Menyediakan view tanpa menambah writer | Proyeksi tidak pernah menjadi writer (P6, P9) |

---

## 10. Artifact Taxonomy

Sistem tipe Artifact konseptual — bukan desain penyimpanan. Setiap tipe mendefinisikan: Knowledge Dimension, State Domain yang **dimiliki**, dan catatan Identity/Relationship.

| Tipe Artifact | Knowledge Dimension | State Domains dimiliki | Catatan Identity & relasi |
|---|---|---|---|
| Vision / Manifesto | Product Intent | Content, Existence | Line tunggal; amendmen jarang |
| Strategy | Product Intent | Content, Existence | Tipe baru (sebelumnya tidak terwakili) |
| Requirement (PRD) | Requirements | Content, Existence | Line + amendmen sebagai instance/Relationship |
| Scope Definition | Planning Knowledge + Requirements | Content, Existence | **Phase context** (Discovery/MVP/Milestone/Release) sebagai atribut — bukan kategori |
| Epic | Planning Knowledge | Content, Existence | Relasi derives-from Scope |
| Plan (roadmap) | Planning Knowledge | Content, Planning, Existence | InstanceVersion signifikan (v1, v2, ...); lock via generation |
| Execution Container (sprint) | — (Content = snapshot proyeksi) | Container, Existence | Content derived; dibuat oleh OS |
| Ticket | — (Content = instruksi eksekusi / Command) | **∅ (tidak ada; State adalah State Projection atas work item yang direferensikan)** | Execution view: proyeksi, bukan writer (P6) |
| Work Item (story, technical story, bug, tech debt, chore, spike) | Requirements / Records / Research | Execution, Existence | Single-writer Execution State; Content dapat didistilasi ke dimensi knowledge |
| Session | — (ephemeral by design) | Existence | Content ephemeral; Completed = kondisi derived; wajib Distillation sebelum Archived |
| Review | Governance & Quality | Content, Existence | Gate semantics; findings → Decisions |
| ADR | Decisions | Content (varian ADR), Existence | Supersession = Relationship ke Line lain |
| Decision Record | Decisions | Content (varian decision), Existence | Supersesi opsional |
| Architecture Description | Architecture | Content, Existence | |
| Specification | Specifications | Content, Existence | Tipe baru (sebelumnya tidak terwakili) |
| Standard / Guideline | Standards & Guidelines | Content, Existence | Tipe baru (sebelumnya tidak terwakili) |
| Runbook / Operational Guide | Operational Knowledge | Content, Existence | |
| Release Record | Records | Content, Existence | Tipe baru (sebelumnya tidak terwakili); berisi agregat eksekusi + gate release |
| Glossary / Term | Vocabulary | Content, Existence | |
| Traceability / Relationship Artifact | Planning Knowledge | Content, Existence | Content = kumpulan Relationship by Identity |

Aturan: Type menentukan State Vector (binding type→state adalah bagian standard); tipe baru adalah ekstensi (Section 14) dengan kewajiban mendeklarasikan State Vector owned lengkap.

---

## 11. Conceptual Lifecycle

### 11.1 Dua kanal perubahan sepanjang waktu

- **Evolusi Content**: Draft → Review → Approved → Amended/Superseded (Content State) — Content matang dan dipreservasi; Approved Content immutable (P8).
- **Progress eksekusi**: Planned → … → Done (Execution State) — work item bergerak maju; koreksi via instance baru, bukan revert.

Keduanya orthogonal: sebuah Artifact dapat memiliki Content Approved dan State eksekusi In Progress secara simultan; tipe menentukan domain mana yang berlaku (Section 10).

### 11.2 Phase sebagai konteks, bukan kategori

Discovery, MVP, Milestone, Release adalah **Phase context** pada Artifact planning/scope — atribut, bukan kategori dan bukan State Domain. Aturan:

- Phase menempel pada Scope Definition / Plan sebagai context attribute.
- **Phase change** (perubahan phase) adalah context update yang diotorisasi **Gate kesiapan**, dievaluasi atas agregat State:
  **release-ready** = (semua work item dalam scope Done) ∧ (seluruh Execution Container Completed) ∧ (plan locked / Immutable) ∧ (gate review lulus) ∧ (gate persetujuan Content lulus).
- Phase tidak pernah menjadi bagian Identity (P3) — Scope Definition tetap identitas yang sama saat produk berpindah fase; yang berubah adalah context attribute-nya.

### 11.3 Lifecycle produk vs lifecycle Artifact

- **Artifact lifecycle**: pergerakan State Domain per Artifact (Section 7) — Content matang, eksekusi selesai, plan dikunci, Artifact di-archive atau di-retire.
- **Product lifecycle**: rangkaian Phase context pada Scope Artifact — Discovery → MVP → Growth → Maturity → Sunset. Setiap fase menghasilkan Artifact-nya sendiri (scope baru, plan baru, Release Record baru) dengan Identity baru; fase lama tetap sebagai Record (P12).

### 11.4 Distillation lifecycle

Ephemeral → durable: Session (Existence) dan Review menghasilkan keputusan/ADR (Content State) melalui jalur Distillation **wajib** sebelum Archived. Preservasi: Artifact Superseded/Archived/Retired tetap ada sebagai Record dengan Content immutable dan Relationship supersession yang utuh.

---

## 12. Storage Independence Model

### 12.1 Eksperimen pikiran

Asumsikan media penyimpanan saat ini hilang; pengetahuan disimpan di relational database, graph database, object store, Atrium, atau platform masa depan. Yang diuji: Identity Model, State Taxonomy, Knowledge Taxonomy, kontrak lapisan.

### 12.2 Yang bertahan (bagian standard)

| Konsep | Alasan bertahan |
|---|---|
| **Identity Model** | Identity = properti konseptual `(Namespace, Type, ID[, InstanceVersion])`. Di database relasional menjadi key; di graph menjadi node property; di object store menjadi key. Referensi by Identity berlaku di semua. |
| **State Taxonomy** | State Domain adalah semantik, bukan penyimpanan. State Projection = view (relasional), query (graph), computed (object store). Single-writer tetap enforceable. |
| **Knowledge Taxonomy** | Klasifikasi adalah properti Artifact; backend apa pun dapat mengindeksnya. |
| **Layer Contracts & Invariant global** | Kontrak mendefinisikan perilaku, bukan penyimpanan. Seluruh 7 invariant (5.4) tetap berlaku. |
| **Exchange Contracts** | Round-trip, idempotensi, referential integrity — independen dari medium. |

### 12.3 Yang termasuk implementasi

- Format serialisasi, skema penyimpanan, struktur indeks.
- **Addressing fisik** (path, key, URL) — standard hanya mensyaratkan: setiap Identity **resolvable** ke tepat satu Artifact dalam satu sistem.
- **Retrieval**: standard mensyaratkan semantik query yang dideklarasikan dapat diimplementasikan; bahasa query adalah implementasi.
- **Enforcement capability**: constraint struktural pada implementasi berbasis berkas; constraint pada database relasional; validation layer pada graph. **Requirement invariant identik; mekanisme enforcement bervariasi** (P16). Nilai "struktur sebagai state machine yang selalu sinkron" pada implementasi awal adalah kapabilitas implementasi; standard mewarisi requirement-nya (determinisme, P11), bukan mekanismenya.

### 12.4 Sikap standard terhadap serialisasi

Standard menetapkan kontrak (apa yang harus dipertahankan, aturan round-trip, validasi); implementasi menetapkan format. Serialisasi Identity harus canonical dan unambiguous di semua implementasi (aturan 6.2.6).

---

## 13. Import / Export Model

### 13.1 Yang harus dipertahankan dalam exchange

| Elemen | Requirement |
|---|---|
| **Identity** | Namespace, Type, ID, InstanceVersion — utuh; tanpa duplikasi; canonical. |
| **State** | Seluruh State Vector (semua domain dimiliki) dengan nilai eksak + history transisi (Change Log). |
| **Content** | Content lengkap, Well-formed sesuai tipe. |
| **Relationships** | Semua relasi by Identity (supersedes, amends, derives-from, depends-on, validates) — referential integrity lintas sistem. |
| **Classification** | Assignment Knowledge Dimension (primer + sekunder). |
| **History** | Link supersesi/amendmen, Change Log, status preservasi (Archived/Retired). |

### 13.2 Round-trip requirements

1. **Lossless**: tidak ada kehilangan atau duplikasi Identity/State/Content/Relationship.
2. **Idempotent**: re-import = no-op (atau replace bersih yang dideklarasikan) — tidak pernah menggandakan Artifact.
3. **Referential integrity**: tidak ada dangling reference setelah import; referensi lintas sistem ter-resolve atau ditolak eksplisit.
4. **Kebijakan konflik Identity**: import dengan Identity yang sudah ada = **tolak atau re-namespace eksplisit** — tidak pernah merge diam-diam.
5. **Validasi sebelum commit**: import memvalidasi kepatuhan terhadap standard (Identity unik, State valid, Content Well-formed) sebelum menulis.
6. **Schema versioning**: kontrak exchange sendiri berversi; import/export menyatakan versi kontrak yang dipatuhi.

### 13.3 Kontrak format serialisasi (bukan formatnya)

Format apa pun wajib: menyandikan Identity secara canonical; merepresentasikan State Vector lengkap; mengekspresikan Relationship by Identity; dapat divalidasi secara mekanis terhadap standard; dapat dibaca oleh Exchange Layer tanpa interpretasi ambigu. Format itu sendiri adalah keputusan implementasi.

---

## 14. Extension Model

### 14.1 Titik ekstensi

| Titik | Berat | Mekanisme |
|---|---|---|
| Artifact type baru | Ringan | Definisi tipe: Knowledge Dimension + **State Vector owned lengkap** + aturan Identity; terdaftar di taksonomi. |
| Knowledge Dimension baru | Ringan-sedang | Sumbu klasifikasi baru; tidak boleh memutus klasifikasi lama (P15). |
| Relationship type baru | Ringan | Semantik relasi baru antar Identity. |
| Protocol variant / Command / Gate baru | Sedang | Variasi Protocol dalam invariant yang ada. |
| State Domain baru | Berat | Hanya jika semantik State tidak tertampung domain yang ada; wajib definisi penuh (owner, aturan, interaksi). |
| Phase vocabulary baru | Ringan | Nilai Phase context baru. |

### 14.2 Aturan ekstensi

1. Ekstensi **tidak boleh melemahkan invariant** (5.4).
2. **Backward compatibility**: seluruh Artifact yang ada tetap valid di bawah ekstensi.
3. **Core closed, taxonomy open**: Identity, kontrak lapisan, dan invariant adalah core yang tertutup terhadap ekstensi (perubahan = revisi standard); taksonomi (tipe, dimensi, domain, protocol) terbuka dengan governance.
4. Ekstensi harus **dapat di-exchange** (tercover schema versioning, Section 13).
5. Governance ekstensi: proposal → review → acceptance (terdaftar sebagai bagian standard). Ini menutup celah ekstensi tanpa prinsip yang teridentifikasi pada implementasi awal.
6. Tipe Artifact baru **wajib mendeklarasikan State Vector owned lengkap** — tidak ada pewarisan default implisit.

---

## 15. Open Questions

**Resolved during ratification:** (1) release dimodelkan sebagai Phase context + Gate kesiapan, bukan State Domain; (2) Exchange Layer adalah lapisan arsitektur, bukan cross-cutting concern; (3) tipe Artifact baru wajib mendeklarasikan State Vector owned lengkap (diangkat menjadi aturan 14.2.6).

Pertanyaan yang tetap terbuka, dengan trade-off masing-masing:

1. **Planning State vs Content State (unifikasi)** — Keputusan saat ini: terpisah (commitment ≠ maturity). Trade-off: terpisah menambah domain; unifikasi menyederhanakan tetapi mencampur "Content matang" dengan "rencana terkunci" yang memiliki Trigger berbeda.
2. **Semantik Line pada supersession** — Apakah supersession adalah dua Line dengan Relationship (keputusan saat ini) atau satu Line dengan instance? Keputusan saat ini menyederhanakan historiografi per-keputusan; alternatif menyederhanakan penelusuran rantai keputusan tetapi mengaburkan bahwa keputusan baru adalah Artifact berbeda.
3. **Kedalaman kontrak query semantics** — Seberapa presisi semantik retrieval harus didefinisikan di standard (agar konsisten lintas implementasi) vs diserahkan ke implementasi (agar inovasi tidak terbatasi)?
4. **Generasi Identity terdistribusi/offline** — Tanpa Identity Registry pusat, bagaimana dua sistem menghasilkan ID bebas-kolisi (ID global vs Namespace terdelegasi vs registry)? Mempengaruhi kontrak exchange.
5. **Kebijakan Projection Refresh** — Event-driven (validasi tiap transisi) vs on-read (validasi saat dibaca): trade-off konsistensi real-time vs biaya. Invariant "proyeksi bukan writer" tidak berubah; mekanisme refresh belum dikunci.
6. **Kedalaman kontrak struktur Content (well-formedness)** — Seberapa ketat struktur Content per tipe Artifact (agar machine-parseable) vs fleksibel (agar ekspresif)? Terlalu ketat membatasi; terlalu longgar merusak determinisme Command.
7. **Klasifikasi multi-dimensi** — Aturan saat ini: satu dimensi primer + sekunder opsional. Konflik antar dimensi dan kewajiban dimensi sekunder belum sepenuhnya ditentukan.

---

## 16. Future Evolution

### 16.1 Milestone konseptual

1. **Ratifikasi standard ini** sebagai model kanonik engineering knowledge (kontrak, invariant, taksonomi dasar).
2. **Exchange contract v1**: definisi round-trip, idempotensi, referential integrity, schema versioning + conformance suite (validator kepatuhan).
3. **Referensi implementasi**: (a) konformasi implementasi awal terhadap standard; (b) implementasi berbasis database relasional; (c) implementasi berbasis graph — membuktikan storage independence (Section 12).
4. **Integrasi Knowledge OS**: Exchange Layer menjadi seam — Knowledge OS mengimpor/mengekspor Artifact dengan Identity, State, dan Relationship utuh; knowledge menjadi queryable dan operable oleh sistem eksternal.
5. **Ekosistem**: validator, importer, tooling agent yang mematuhi kontrak; ekstensi terdaftar via governance Section 14.

### 16.2 Peran implementasi awal ke depan

Implementasi awal berubah status: dari "arsitektur itu sendiri" menjadi **satu serialisasi referensi** yang (a) mendemonstrasikan kepatuhan terhadap standard, (b) menjadi baseline onboarding, (c) tetap berfungsi sebagai Engineering Operating System untuk proyek yang memilih medium tersebut. Standard menjadi kanon; implementasi menjadi contoh.

### 16.3 Invariant evolusi

Evolusi standard tidak pernah mengubah: Identity (P3), invariant global (5.4), prinsip dua kanal (P10), dan komposisi lapisan (KB + OS + EX). Yang boleh berevolusi: taksonomi (dimensi, tipe, domain, protocol) melalui governance ekstensi. Fondasi yang tidak dapat dinegosiasikan oleh iterasi mana pun: **knowledge base dan operating system sebagai dua lapisan dari satu sistem, diikat Identity, State dimiliki Operating Layer, pipeline sebagai Protocol first-class.**

---

*End of Canonical Specification — EKA v1.0 (Ratified).*
