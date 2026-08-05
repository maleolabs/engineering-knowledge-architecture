# Ratification Notes — EKA Exchange Specification v1.0

> Anchor EKA: Exchange Layer — kontrak pertukaran (EKA 13, 16.1 milestone 2). Dokumen konvensi, bukan artefak.
> Standar: EKA v1.0, tanggal 2026-08-05.

## Status

**Ratified** — EKA Exchange Specification v1.0 (refinement pass selesai; arsitektur identik dengan versi sebelumnya, tidak ada keputusan arsitektur yang dibuka ulang).

## Ringkasan Refinement

Refinement pass pra-ratifikasi — bukan redesign. Invariant, aturan, dan keputusan arsitektur Exchange Specification v1.0 **tidak berubah**. Perubahan bersifat aditif dan editorial:

| Area | Perubahan | Sifat |
|---|---|---|
| Exchange Object Model (§4.4 baru) | Model objek kanonik Exchange Package: hierarki kontainmen (Contract Header, Manifest, Exchange Units, Declarations, integrity data), komposisi Exchange Unit, kardinalitas, determinisme, dan prinsip proyeksi ("setiap serialisasi adalah proyeksi dari model ini") | Aditif — abstraksi baru, tidak mengubah aturan §10 |
| Capability Declaration (§4.5 baru + §10.1) | Blok deklaratif opsional di Contract Header: keluarga spesifikasi didukung, kelas ekstensi, tipe Relationship ekstensi, varian state, export scopes, kapabilitas protokol masa depan; warning-not-blocking untuk mismatch | Aditif — opsional, backward compatible |
| Lokasi deklarasi ditetapkan tunggal (§7.3, §10.1–10.2, §10.4, §18.3) | Deklarasi Closure/External Reference/Extension = elemen Declarations level paket; Header hanya mengumumkan; Manifest hanya daftar unit + scope type | Konsistensi — menghapus kontradiksi lokasi |
| Idempotency dipertegas (§11.1 fase 6–7, §11.2) | Unit identik (Identity + payload sama) → duplicate/no-op fase 7; konflik fase 6 hanya untuk Identity sama + payload berbeda | Konsistensi — selaras §15.5 |
| Coverage Conformance Rules lengkap (§11.1 fase 3, 10) | Rule R6 (classification) masuk pipeline fase 3; R8 (single-writer/projection) dievaluasi via revalidasi fase 10 terhadap repository state | Konsistensi — tidak ada rule tanpa fase |
| Referensi "validation rule N" → "rule N (§14.2)" | Konsisten dengan Naming and Terminology Specification v1.0 | Terminologi |
| Tabel Conformance Rules diberi label R1–R9 (§14.2) + catatan R0, ADR-superseded, order change-log | Selaras pemetaan Aturan 1–9 ↔ R1–R9 (Naming §9.3) | Terminologi |
| Perbaikan editorial (#13, #19, #21–25) | Prosa fase 10 rollback, "permissible differences" (15.4), dedupe §1.4, wording §9.1, reading note (gloss + subordinate terms), footer resmi | Editorial |

## Terminologi

Verdict: **delapan istilah inti dipertahankan** — Exchange Package, Exchange Unit, Contract Header, Manifest, Export Scope, Import Manifest, Closure, External Reference. Tidak ada rename: seluruh istilah sudah terdaftar (Naming and Terminology §12.1), tidak ambigu, dan rename = perubahan kontrak (N1). Penambahan istilah baru: **Capability Declaration**, **Closure Declaration**, **Collection**, **Graph** — didaftarkan di tabel konsep §4.1 (defined-before-use, per tata kelola terminologi).

## Model Konseptual Baru

1. **Canonical Exchange Package Object Model (§4.4)** — abstraksi kanonik yang menjadi acuan derivasi format serialisasi masa depan. Properti kunci: kardinalitas eksplisit per elemen, prinsip proyeksi, determinisme, dan kesetaraan round-trip didefinisikan sebagai kesetaraan proyeksi model yang sama.
2. **Capability Declaration (§4.5)** — deklarasi dukungan implementasi, opsional, tidak pernah syarat validitas paket; mismatch = warning di Import Manifest; jalur rejection tetap milik §9.2 (versi) dan §16.3 (ekstensi tak dikenal).

## Konsistensi

Audit konsistensi final (alex-qa): **0 Critical, 4 Major, 16 Minor, 5 Editorial** — seluruhnya difiksasi pada level teks. Verifikasi positif: 19 istilah glossary konsisten, seluruh cross-reference internal + canonical resolve, pemetaan R1–R9 ↔ Aturan 1–9 defensibel, draft tolerance konsisten di 4 titik, idempotency/round-trip konsisten, kepatuhan Naming and Terminology Specification (pola H1, Anchor, Status, slug, bahasa Inggris).

## Kesiapan Arsitektural

**Arsitektur dinyatakan siap ratifikasi.** Tidak ada blocker. Perubahan berikutnya dibatasi pada perbaikan editorial dan klarifikasi terminologi; setiap perubahan aditif di versi mendatang mengikuti aturan minor version (Naming and Terminology §5.3).

---

*End of Ratification Notes — EKA Exchange Specification v1.0.*
