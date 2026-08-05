# Zona Standard — EKA v1.0

Zona ini berisi **standard itu sendiri**: salinan kanonik dari Engineering Knowledge Architecture (EKA) v1.0 beserta glosarium istilah kanoniknya dan kontrak pertukarannya (Exchange Specification v1).

## Sifat zona ini

- Zona ini adalah **pra-lapisan** (pre-layer): standard mendefinisikan lapisan (Knowledge Layer, Operating Layer, Exchange Layer), tetapi standard sendiri **bukan artifact dari proyek mana pun**.
- Isi zona ini **bukan bagian dari serialisasi proyek** — ia adalah teks kanonik yang disalin untuk referensi, onboarding, dan konformasi.
- Berbeda dengan artifact di `skeleton/docs/`, dokumen di zona ini tidak membawa frontmatter Identity/State: dokumen di sini adalah salinan standard, bukan Artifact yang dikelola Operating Layer.

## Isi

| File | Isi |
|---|---|
| [`eka-specification-v1.0.md`](eka-specification-v1.0.md) | Teks kanonik lengkap EKA v1.0 (16 section): prinsip (P1–P16), Core Concepts, Layer Model, kontrak lapisan, Identity Model, State Taxonomy, Knowledge Taxonomy, Execution Taxonomy, Artifact Taxonomy, lifecycle konseptual, storage independence, import/export, ekstensi, open questions, evolusi. |
| [`eka-exchange-specification-v1.0.md`](eka-exchange-specification-v1.0.md) | Exchange Contract v1 (Ratified, milestone 16.1.2): unit exchange terkecil (Artifact Instance), representasi Identity/Relationship/State, tiga dimensi versioning, **Exchange Package Object Model** (§4.4), **Capability Declaration** (§4.5), semantik import/export/sinkronisasi, konformitas (R1–R9), jaminan round-trip, kompatibilitas, keamanan, evolusi. Konseptual — bebas format serialisasi. |
| [`eka-naming-and-terminology-specification-v1.0.md`](eka-naming-and-terminology-specification-v1.0.md) | Meta-specification (Ratified): penamaan resmi ekosistem EKA — product identity, pola naming Specification Families, reference components, tooling, repository naming, tabel istilah kanonik, daftar terminologi deprecated, migrasi. |
| [`glossary.md`](glossary.md) | Glosarium alfabetis seluruh istilah kanonik berhuruf kapital, dengan definisi eksak dari teks kanonik. |

## Rujukan lain

- Dokumentasi meta implementasi: [`../reference/README.md`](../reference/README.md)
- Struktur serialisasi yang dapat disalin: [`../skeleton/README.md`](../skeleton/README.md)
