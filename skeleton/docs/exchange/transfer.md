# Transfer — Konvensi Impor/Ekspor

> Anchor EKA: Exchange Layer — transfer (import/export). Dokumen konvensi, bukan artefak.
> Standar: EKA v1.0, tanggal 2026-08-05.

Dokumen ini mengatur pertukaran artefak EKA antar repositori (mis. repositori induk/komponen/arsip). Semua konvensi berlaku dua arah (impor = ekspor yang dibalik).

## 1. Persyaratan Round-Trip

### 1.1 Lossless (tanpa kehilangan)

Transfer harus mempertahankan **semua** informasi identitas dan state:

- [ ] Identitas penuh: `namespace`, `type`, `id`, `instance-version`, `revision`.
- [ ] State penuh: seluruh owned state vector **beserta** riwayat `change-log`.
- [ ] Konten yang well-formed (bagian wajib utuh).
- [ ] Relasi **berdasarkan identitas** (referensi dipertahankan sebagai referensi, bukan diubah menjadi teks).
- [ ] Klasifikasi: `dimension`, `dimensions-secondary`.
- [ ] Status preservasi (`existence-state`) tidak berubah oleh transfer.

### 1.2 Idempotent

- [ ] Impor ulang paket yang sama = **no-op** (tidak ada duplikasi), atau deklarasi *clean replace* yang eksplisit.
- [ ] Impor ulang tidak pernah membuat artefak duplikat.

### 1.3 Integritas Referensial

- [ ] Tidak ada referensi menggantung (dangling): artefak yang dirujuk wajib ikut dipindahkan, sudah ada di target, atau diizinkan sebagai peringatan karena status target `draft`.

### 1.4 Kebijakan Konflik Identitas

Saat identitas `(namespace, type, id, instance-version)` sudah ada di target:

| Opsi | Syarat |
|---|---|
| **Tolak** (default) | konflik dilaporkan; transfer dibatalkan |
| **Re-namespace eksplisit** | seluruh identitas artefak dipindahkan ke namespace baru secara deklaratif (dan semua referensinya diperbarui konsisten) |

- [ ] **Tidak pernah** dilakukan *silent merge* — menggabungkan dua artefak berbeda identitas secara diam-diam dilarang.

### 1.5 Validasi Sebelum Commit

- [ ] Paket yang akan ditransfer lolos seluruh checklist [validation.md](validation.md).
- [ ] Setelah impor, target divalidasi ulang sebelum commit.

### 1.6 Versi Kontrak

- [ ] Setiap paket transfer mendeklarasikan versi kontrak exchange yang dipakainya (mis. `eka-exchange: 1.0`).
- [ ] Impor menolak paket dengan versi kontrak yang tidak didukung.

## 2. Yang Dipreservasi

| Aspek | Keterangan |
|---|---|
| Identitas penuh | `(namespace, type, id, instance-version)` tidak berubah selama transfer |
| State + change-log | seluruh riwayat transisi domain yang dimiliki |
| Konten well-formed | bagian wajib per keluarga tipe tetap utuh |
| Relasi by identity | referensi tetap mengikuti identitas (bukan path file) |
| Klasifikasi | `dimension`/`dimensions-secondary` dipertahankan |
| Status preservasi | `existence-state` tidak diubah oleh mekanika transfer |

## 3. Batasan EX pada Transfer

- EX **tidak menilai kebenaran konten** — hanya konformitas dan integritas.
- EX **tidak mengubah state** — transisi state tetap hanya boleh oleh pemilik state; transfer hanya menyalin nilainya.
- Proyeksi (tiket/tabel) tidak ditransfer sebagai sumber kebenaran; setelah impor, proyeksi di-refresh ulang dari owner state di target.
