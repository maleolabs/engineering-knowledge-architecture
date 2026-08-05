# Filosofi — Mengapa EKA Ada dan Mengapa Repositori Ini Disusun Demikian

Dokumen ini adalah narasi posisi: alasan keberadaan EKA, wawasan arsitektur yang mendasarinya, dan konsekuensinya terhadap cara repositori ini disusun. Ini bukan bagian standard — ini catatan mengapa keputusan standard dan serialisasi diambil.

## Mengapa EKA ada

Implementasi awal membuktikan dua tanggung jawab yang hidup berdampingan dalam satu repositori: menyimpan pengetahuan (Knowledge Base) dan menjalankan pekerjaan engineering (Engineering Operating System). Keduanya berjalan di atas struktur yang sama — dan di situlah konflasi lahir. Status diduplikasi di tiga tempat tanpa penulis tunggal, tujuh mesin status hidup berdampingan tanpa aturan bersama, Identity dikacaukan dengan lokasi, dan pipeline proses tercampur menjadi taksonomi folder.

EKA adalah jawaban atas konflasi itu: bukan struktur baru, melainkan **model konseptual kanonik** yang memisahkan tanggung jawab secara eksplisit, lalu membiarkan setiap implementasi memilih mekanismenya sendiri. Repositori ini adalah salah satu serialisasi dari model itu.

## Wawasan dual-layer: Knowledge Base dan Operating System adalah dua lapisan dari satu sistem

Wawasan inti EKA: repositori engineering bukan sekadar tempat menyimpan dokumen, dan bukan sekadar mesin eksekusi — ia **keduanya sekaligus, sebagai dua lapisan dari satu sistem**. Lapisan pengetahuan menyimpan Content, klasifikasi, dan sejarah; lapisan operating menjalankan state machine, protocol, dan governance. Keduanya diikat oleh **Identity**: satu artifact yang sama dilihat Knowledge Base sebagai pengetahuan dan dilihat Operating System sebagai entitas ber-state.

Konsekuensi langsung bagi repositori ini: folder knowledge (12 dimensi) dan folder `operating/` bukan dua area yang kebetulan berdampingan — keduanya adalah dua lapisan yang kontraknya didefinisikan standard (Section 4–5). Satu artifact membawa Content di lapisan KB dan State di lapisan OS, tanpa saling menulis.

## Pipeline sebagai Protocol first-class, bukan kecelakaan taksonomi

Pada struktur lama, alur kerja (PRD → MVP → epics → roadmap → sprint → tickets → work items → sessions) tampak sebagai hierarki folder. Itu adalah kesalahan kategori: pipeline adalah **urutan eksekusi**, properti Operating Layer, bukan properti klasifikasi. Ketika pipeline disandikan sebagai folder, setiap pergeseran proses memaksa pergeseran lokasi — dan Identity ikut bergeser.

EKA memindahkan pipeline ke tempat yang benar: **Protocol** (EKA 3, 9). Urutan "requirement → scope → plan → container → work item → konteks kerja → validasi" didefinisikan sebagai properti Operating Layer yang menjawab "apa selanjutnya" secara deterministik. Konsekuensi repositori: tidak ada folder pipeline; yang ada adalah folder dimensi (untuk pengetahuan) dan folder operating (untuk eksekusi), dengan urutan diekspresikan lewat Relationship dan State — bukan lewat posisi folder.

## Phase adalah konteks, bukan kategori

Fase produk (Discovery, MVP, Milestone, Release) adalah **konteks sepanjang waktu**, bukan kategori permanen dan bukan state. Produk yang sama tetap produk yang sama saat berpindah fase — yang berubah adalah atribut konteksnya. Menyandikan phase sebagai folder membuat fase menjadi bagian lokasi, dan lokasi menjadi bagian Identity: pindah fase berarti "memindahkan" artifact, padahal Identity tidak boleh bergeser (P3).

Konsekuensi repositori: `phase` adalah field frontmatter pada artifact scope/plan saja; perpindahan fase adalah context update yang diotorisasi gate kesiapan (EKA 11.2) — bukan operasi pemindahan file. Identity decoupled dari phase.

## State dimiliki single-writer

Duplikasi status adalah penyakit utama struktur lama: status hidup di metadata tabel dokumen, di tabel sprint, di dokumen ticket — tiga salinan, tanpa penulis yang jelas, dengan ritual sinkronisasi manual yang selalu tertinggal. EKA menetapkan: **setiap field state memiliki tepat satu owner** (P6). Semua tampilan lain adalah State Projection: diturunkan, divalidasi, dan tidak pernah menjadi writer.

Konsekuensi repositori: state hanya ditulis di frontmatter artifact pemiliknya (single-writer per domain), transisi dicatat di `change-log`, dan representasi turunan (tabel container, ticket) diberi label eksplisit sebagai proyeksi yang di-refresh on-read — bukan fakta independen yang diedit.

## Repositori adalah satu serialisasi, bukan arsitekturnya

Disiplin terakhir yang menuntun seluruh struktur ini: **standard adalah kanon; repositori adalah contoh** (EKA 1.3, 16.2). EKA didefinisikan atas konsep — Identity, State, Layer, kontrak — bukan atas folder, pola penamaan, atau format. Repositori ini hanyalah satu cara menserialisasi konsep-konsep itu ke Git + Markdown. Konsekuensinya:

- Standard hidup utuh di `standard/` sebagai teks kanonik yang tidak tercampur keputusan serialisasi;
- Konvensi serialisasi didokumentasikan eksplisit (bukan disembunyikan di kebiasaan) dan dijaga lewat ADR;
- Struktur yang dapat disalin (`skeleton/`) adalah produk sampingan — serialisasi yang dapat dipakai proyek lain;
- Setiap keputusan serialisasi dapat dipertanggungjawabkan ke section/prinsip standard (lihat `traceability-matrix.md`).

Dengan disiplin ini, repositori tetap berguna enam bulan lagi — dan enam puluh bulan lagi, bahkan ketika media penyimpanan berganti: yang dipertahankan adalah konsep, bukan mekanisme.
