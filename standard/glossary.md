# Glossary — EKA v1.0

Glosarium alfabetis seluruh istilah kanonik berhuruf kapital dari spesifikasi EKA v1.0. Definisi direproduksi dari teks kanonik (`standard/eka-specification-v1.0.md`) — bukan parafrase. Rujukan section dicantumkan untuk navigasi.

## A

### Artifact
entitas pengetahuan engineering yang memiliki Identity, Content, State Vector (domain State yang **dimiliki**), dan Relationship. Unit dasar model. *(Section 3)*

### Artifact Instance
satu versi eksistensi dari sebuah Line: Line + `InstanceVersion`. *(Section 3)*

### Artifact Line
entitas Identity yang bertahan: satu `(Namespace, Type, ID)`. *(Section 3)*

## C

### Change Log
catatan kronologis transisi State pada Artifact: domain, nilai lama, nilai baru, waktu, otoritas. Wajib untuk seluruh State Domain. *(Section 3)*

### Command
instruksi eksekusi deterministik yang dikonsumsi oleh eksekutor (manusia atau agent). Content Command adalah Content (milik Knowledge Layer); eksekusinya diatur Protocol (Operating Layer). *(Section 3)*

### Container State
Domain State Execution Container: nilai `Active → Completed`; responsibility "Execution Container terbuka/tertutup; konkurensi"; owner Operating Layer. Aturan transisi: "Completed adalah transisi derived: dipicu agregat Execution State (semua work item Done); exactly-one-Active (mutual exclusion)." *(Section 7.2)*

### Content
muatan semantik Artifact: intent, keputusan, desain, kendala, prosedur, catatan. Milik Knowledge Layer. *(Section 3)*

### Content State
Domain State kematangan Content: nilai `Draft → Review → Approved`; terminal pasca-persetujuan: `Amended | Superseded`; responsibility "Kematangan Content sebagai pengetahuan; kanal governance"; owner "Knowledge Layer (gate: owner/approver)". Aturan transisi: "Forward-only; gate persetujuan; perubahan pasca-Approved hanya via amendmen/supersesi (P8)." *(Section 7.2)*

## D

### Distillation
transformasi pengetahuan ephemeral (konteks kerja, temuan review) menjadi pengetahuan durable (keputusan, ADR, Record). *(Section 3)*

## E

### Exchange Layer
Lapisan arsitektur ketiga: "Batas transformasional: serialisasi, validasi, import/export, mediasi sistem eksternal". Memiliki "Kontrak pertukaran, aturan round-trip, validasi kepatuhan"; tidak memiliki "Content dan State (tidak pernah menjadi owner)". *(Section 4.1)*

### Execution Container
Artifact eksekusi yang membungkus work item dan membawa konvensi konkurensi (exactly-one-active). Domain State-nya: Container State. Contoh implementasi: sprint. *(Section 3)*

### Execution State
Domain State progress work item: nilai `Planned → Todo → In Progress → In Review → Done`; responsibility "Progress work item melalui Protocol"; owner "Operating Layer (single-writer: artifact work item)". Aturan transisi: "Strictly sequential; never skip; never revert; satu initial; satu terminal; tiap transisi tercatat di Change Log." *(Section 7.2)*

### Existence State
Domain State presensi Artifact: nilai `Active → Archived → Retired`; responsibility "Presensi Artifact dalam proses aktif vs preservasi"; owner "Operating Layer (transisi); Knowledge Layer (prinsip preservasi, P12)". Aturan transisi: "Forward-only; **Archived** = reference-only (tanpa transisi State lain, tanpa mutasi Content, tetap tersedia di retrieval); **Retired** = preservasi terminal (tidak disurface di retrieval normal, Content immutable, Identity tidak berubah). Berlaku untuk semua tipe Artifact." *(Section 7.2)*

## G

### Gate
kondisi yang harus dipenuhi sebelum transisi atau eksekusi boleh terjadi (gate persetujuan, gate kesiapan, gate review). *(Section 3)*

## I

### Identity
properti Artifact yang membedakannya dari semua Artifact lain secara permanen: `(Namespace, Type, ID[, InstanceVersion])`. Lihat Section 6. *(Section 3)*

### Identity Registry
fungsi Knowledge Layer yang menjamin keunikan Identity, resolusi Identity ke Artifact, dan integritas referensi. *(Section 3)*

### InstanceVersion
bagian dari Identity instance; diskriminator instance dalam satu Line. "InstanceVersion (versi yang menunjuk instance berbeda) — **bagian dari Identity instance**. Instance baru diciptakan dengan sengaja (contoh: plan v2 setelah v1 terkunci). Perubahan instance = perubahan Identity instance; Line tetap." *(Sections 6.1, 6.3)*

## K

### Knowledge Dimension
sumbu klasifikasi pengetahuan (Section 8). Properti Artifact, bukan Identity. *(Section 3)*

### Knowledge Layer
Lapisan arsitektur pertama: "Store pengetahuan: Content, klasifikasi, preservasi, referensi". Memiliki "Content, klasifikasi (Knowledge Dimensions), Relationship, history/Records, administrasi Identity (Identity Registry)"; tidak memiliki "State proses, Protocol eksekusi". *(Section 4.1)*

### Knowledge OS
platform eksekusi pengetahuan masa depan yang mengonsumsi dan memproduksi Artifact EKA melalui Exchange Layer. Bukan bagian standard; konsumen standard. *(Section 3)*

## L

### Layer
lapisan arsitektur dengan kontrak eksplisit: Knowledge Layer, Operating Layer, Exchange Layer (Section 4). *(Section 3)*

## N

### Namespace
ruang Identity yang memisahkan domain pengelolaan (produk, organisasi, sistem). *(Section 3)*

## O

### Operating Layer
Lapisan arsitektur kedua: "State machine & Protocol eksekusi". Memiliki "State Domain (Execution, Planning, Container, Existence), urutan, konkurensi, kuncian, Gate, Command"; tidak memiliki "Content (tidak pernah mengedit Content)". *(Section 4.1)*

## P

### Phase
konteks produk/scope sepanjang waktu (Discovery, MVP, Milestone, Release). Phase adalah **context attribute** pada Artifact planning/scope: bukan kategori, bukan State Domain. **Phase change** adalah context update yang diotorisasi Gate kesiapan (Sections 7.5, 11.2). *(Section 3)*

### Planning State
Domain State commitment plan: nilai `Draft → Approved → Immutable`; responsibility "Commitment plan: tentative → committed → locked"; owner Operating Layer. Aturan transisi: "Forward-only; Approved = siap dieksekusi; **Immutable dicapai atomik dengan peristiwa pembuatan Execution Container** (lock-atomic-with-generation); perubahan pasca-lock = instance baru (InstanceVersion)." *(Section 7.2)*

### Protocol
aturan deterministik milik Operating Layer: urutan, transisi State, kuncian, gate, perintah eksekusi. *(Section 3)*

### Projection Refresh
mekanisme dan waktu validasi State Projection terhadap owner (on-read dan/atau event). *(Section 3)*

### Projection Semantics
aturan komputasi State Projection: dari owner mana, agregasi bagaimana, apa yang ditampilkan. *(Section 3)*

## R

### Record
Artifact yang dipreservasi sebagai sejarah (Superseded, Archived, Retired, release record) dengan Content immutable. *(Section 3)*

### Relationship
relasi eksplisit antar Artifact yang direferensikan by Identity: supersedes, amends, derives-from, depends-on, validates. *(Section 3)*

### Revision
pelacakan evolusi Content instance yang sama — **bukan bagian dari Identity**. "Revision (pelacakan evolusi Content instance yang sama) — **bukan bagian dari Identity**. Revision berubah setiap edit dan tidak boleh memutus referensi." *(Sections 6.1, 6.3)*

## S

### State
fakta tentang posisi Artifact dalam proses tertentu. *(Section 3)*

### State Domain
dimensi State independen dengan semantik, owner, dan aturan transisinya sendiri (Section 7). Domain bersifat orthogonal. *(Section 3)*

### State Projection
tampilan State yang diturunkan dari owner (contoh: agregat State work item → status Execution Container). Proyeksi tidak memiliki State sendiri; proyeksi tidak pernah menjadi writer. *(Section 3)*

### State Vector
tuple State Domain **yang dimiliki** oleh sebuah Artifact. State yang diproyeksikan (State Projection) bukan bagian State Vector. *(Section 3)*

## T

### Trigger
hubungan antar domain/operasi di mana suatu peristiwa memicu validasi atau transisi pada domain lain (definisi formal interaksi, Section 7.5). *(Section 3)*

## W

### Well-formed Content
Content yang mematuhi struktur yang ditetapkan untuk tipe Artifact-nya, sehingga dapat diparse dan dieksekusi secara deterministik. *(Section 3)*

---

*Glosarium kanonik EKA v1.0 — definisi mengikat; lihat `eka-specification-v1.0.md` untuk konteks penuh.*
