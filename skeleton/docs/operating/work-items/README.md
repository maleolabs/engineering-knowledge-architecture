# docs/operating/work-items/ — Work Items

> Anchor EKA: Operating Layer — state domain **Execution State**.

## Tujuan

Work items adalah unit kerja terkecil yang dieksekusi dalam kontainer. Enam subtipe menampung seluruh jenis pekerjaan; semuanya berbagi satu state domain yang dimiliki: Execution State.

## Execution State

| Nilai | Makna |
|---|---|
| `planned` | terdaftar, belum dijadwalkan |
| `todo` | siap dikerjakan |
| `in-progress` | sedang dikerjakan |
| `in-review` | menunggu review gate |
| `done` | selesai dan tervalidasi |

Aturan: **never skip** (wajib melewati setiap nilai berurutan), **never revert** (tidak boleh kembali ke nilai sebelumnya). Transisi forward-only dan tercatat di `change-log`.

## Single-Writer

Setiap work item memiliki **satu penulis state** (pelaksana/implementer-nya). Hanya penulis itu yang mengubah `execution-state` dan `existence-state`. Pelaksana menulis `change-log` pada setiap transisi; peninjau tidak pernah menulis state work item.

## Enam Subtipe

| Token | Subtipe | Folder |
|---|---|---|
| `sto-` | Story | [stories/](stories/) |
| `ts-` | Technical Story | [technical-stories/](technical-stories/) |
| `bug-` | Bug | [bugs/](bugs/) |
| `td-` | Tech Debt | [tech-debt/](tech-debt/) |
| `ch-` | Chore | [chores/](chores/) |
| `spk-` | Spike | [spikes/](spikes/) |

## Struktur Umum yang Baik

Semua work item wajib memuat:

- `## Description` — apa yang dikerjakan.
- Kriteria verifikasi sesuai subtipe (`## Acceptance Criteria`, `## Impact`, dst. — lihat README subtipe).
- Frontmatter identitas: `namespace`, `type`, `id`; state: `execution-state`, `existence-state`; `change-log` untuk setiap transisi.
- `dimension` boleh diisi informasional (mis. `requirements`), tetapi **tidak** menjadi penentu folder rumah.

## Terkait

- [containers/](../containers/) — work item hidup dalam `ctr-` yang dirujuk.
- [projections/](../projections/) — `tkt-` memproyeksikan work item ke tabel/status.
- [sessions/](../sessions/) — eksekusi work item dicatat dalam `ses-`.
- [../quality/](../../quality/) — hasil diverifikasi `rvw-`.
