# Protocol — Operating Manual

> Anchor EKA: Operating Layer (OS). Dokumen konvensi — bukan artefak (tanpa `type`/`id` di frontmatter).
> Standar: EKA v1.0, tanggal 2026-08-05.

Dokumen ini adalah manual pengoperasian serialisasi EKA: rantai pemesanan, state domains, aturan transisi, dan kewajiban yang harus dipatuhi oleh setiap penulis state.

## 1. Rantai Ordering

Setiap nilai harus lahir dalam urutan ini — langkah tidak boleh dilompati:

```
requirement → scope → capability → plan → container → work item → session → review
```

| Langkah | Token | Folder | Catatan |
|---|---|---|---|
| requirement | `req-` | `../requirements/` | kebutuhan yang disepakati |
| scope | `scp-` | `../planning/` | konteks berfase (memuat `phase`) |
| capability | `epc-` | `../planning/` | kapabilitas yang mewujudkan scope |
| plan | `plan-` | `../planning/` | rencana eksekusi; mengunci atomik dengan container |
| container | `ctr-` | `containers/` | eksekusi aktif; tepat satu pada satu waktu |
| work item | `sto-`/`ts-`/`bug-`/`td-`/`ch-`/`spk-` | `work-items/` | unit kerja yang dieksekusi |
| session | `ses-` | `sessions/` | catatan eksekusi kerja |
| review | `rvw-` | `../quality/` | verifikasi hasil |

## 2. Lima State Domain

Setiap domain memiliki nilai valid, satu nilai awal, satu nilai akhir, dan transisi **forward-only** (tidak boleh melompat, tidak boleh mundur).

| Domain | Nilai valid | Dimiliki oleh | Awal → Akhir |
|---|---|---|---|
| Content State | `draft, review, approved, amended` — varian ADR: `proposed, accepted, superseded`; varian decision: `draft, accepted, superseded` | artefak pengetahuan & planning (`vis-`, `str-`, `req-`, `scp-`, `epc-`, `plan-`, `trc-`, `arc-`, `adr-`, `dec-`, `spec-`, `std-`, `run-`, `rvw-`, `rel-`, `gls-`, `fnd-`) | draft/proposed → terminal (approved/accepted/amended/superseded) |
| Execution State | `planned, todo, in-progress, in-review, done` | `sto-`, `ts-`, `bug-`, `td-`, `ch-`, `spk-` | planned → done |
| Planning State | `draft, approved, immutable` | `plan-` | draft → immutable |
| Container State | `active, completed` | `ctr-` | active → completed |
| Existence State | `active, archived, retired` | universal — semua tipe yang menyimpan state | active → archived → retired |

Aturan transisi (berlaku untuk semua domain):

- **Forward-only:** hanya menuju nilai yang berada di depan dalam urutan nilai.
- **Never skip:** setiap langkah berurutan dilalui (mis. Execution State tidak boleh lompat dari `todo` langsung ke `done`; harus melewati `in-progress` dan `in-review`).
- **Never revert:** nilai yang sudah dilewati tidak boleh kembali (mis. `approved` tidak kembali ke `draft`).
- **Single initial:** satu nilai awal yang valid per domain.
- **Single terminal:** satu nilai akhir yang valid per domain.
- **Recorded:** setiap transisi dicatat di `change-log` oleh pemilik tunggalnya (lihat §6).

## 3. Tepat Satu Kontainer Aktif

Hanya ada **satu** `ctr-` dengan `container-state: active` pada satu waktu (mutual exclusion). Pembuatan kontainer baru hanya diizinkan setelah kontainer lama mencapai `completed`. Kontainer aktif yang baru lahir mengunci plan (lihat §4).

## 4. Lock-Atomic-With-Generation

- `plan-` bergerak `draft → approved → immutable`.
- Transisi `plan-` menjadi `immutable` terjadi **atomik** dengan pembuatan `ctr-`: keduanya adalah satu operasi; tidak boleh ada plan immutable tanpa kontainer, dan tidak boleh ada kontainer tanpa plan terkunci.
- Setelah terkunci, **perubahan apa pun terhadap plan tidak mengedit instance itu** — melainkan membuat generasi baru: `plan-<id>-v<instance-version + 1>.md`. Setiap instance adalah snapshot kebenaran pada generasinya.
- Generasi baru harus melalui rantai yang sama (draft → approved → immutable).

## 5. Gate

Gate dievaluasi berdasarkan **owner state** (state milik artefak itu sendiri), bukan proyeksi:

| Gate | Dimana | Syarat lulus |
|---|---|---|
| Approval gate | Knowledge Layer & planning | artefak mencapai `content-state: approved` (atau `accepted`) oleh pemilik kontennya |
| Readiness gate | work item → eksekusi | work item berada di Execution State yang benar (mis. `in-progress`) dan dependensinya terpenuhi |
| Review gate | sebelum selesai | work item berada di `in-review` dan review (`rvw-`) yang memvalidasinya sudah disetujui |

Jika sebuah gate mengevaluasi proyeksi (mis. tabel kontainer) dan hasilnya bertentangan dengan owner state, **owner state yang menang**.

## 6. Aturan Change-Log

- Setiap artefak pemilik state memiliki field `change-log:` — larik objek `{date, domain, from, to, by}`.
- Ditulis oleh **satu penulis** (pemilik state) pada **setiap** transisi — tidak ada transisi tanpa catatan.
- Entri terakhir per domain **harus sama** dengan nilai field saat ini (divalidasi oleh exchange/validation.md).
- `domain` memakai nama domain state (mis. `execution-state`, `phase`).

## 7. Dua Kanal Perubahan

| Kanal | Mekanisme | Menulis ke |
|---|---|---|
| Content governance | tata kelola konten (amandemen, review konten, supersesi) | konten artefak, `revision` |
| State protocol | protokol state (transisi, `change-log`) | field state, `instance-version` |

**Kedua kanal tidak boleh dicampur:** mengubah state bukan cara mengubah konten, dan mengubah konten bukan cara mengubah state. Perubahan konten pada artefak berversi yang sudah terkunci mewajibkan generasi baru (bukan edit langsung).

## 8. Kewajiban Distilasi

- Setiap sesi (`ses-`) yang menghasilkan temuan **wajib didistilasi** sebelum diarsipkan: temuan yang memengaruhi arah → `dec-`/`adr-`; temuan teknis yang belum memutuskan → `fnd-` (research); prosedur yang terbukti → `run-` (operations).
- Work item tipe spike (`spk-`) wajib memuat tautan ke tempat distilasi pengetahuannya (research/decisions) di `## Conclusion`.
- Arsip tanpa distilasi adalah pelanggaran protokol (EKA 11.4).
