# spikes/ — Spike (`spk-`)

> Anchor EKA: Operating Layer — state domain Execution State; subtipe work item `spk-`.

## Tujuan

Spike adalah unit kerja investigasi terbatas-waktu untuk mengurangi ketidakpastian: membuktikan kelayakan, mengeksplorasi pendekatan, atau mengumpulkan data sebelum keputusan. Hasil spike **bukan** keputusan — hasilnya harus didistilasi ke tempat pengetahuan.

## Token & State Vector

| Token | Folder | State yang dimiliki |
|---|---|---|
| `spk-` | `work-items/spikes/` | `execution-state`, `existence-state` |

Nilai `execution-state`: `planned → todo → in-progress → in-review → done` (forward-only, never skip, never revert). Nilai `existence-state`: `active → archived → retired`.

## Struktur Konten Wajib

- `## Description` — pertanyaan yang diselidiki dan batas waktu/lingkup.
- `## Investigation Notes` — jejak penyelidikan: eksperimen, data, sumber.
- `## Conclusion` — simpulan **wajib** memuat tautan ke tempat distilasi pengetahuan (mis. `fnd-` di research/ atau `dec-`/`adr-` di decisions/) — distilasi sebelum arsip (EKA 11.4).

## Konvensi Nama

`spk-<id>.md`, kebab-case, unik dalam (namespace, type). Tanpa akhiran `-v<nn>`. Contoh: `spk-kelayakan-proyeksi-tiket.md`.

## Kepemilikan

| Peran | Tanggung jawab |
|---|---|
| Engineer (implementer) | pemilik tunggal state; penyelidikan |
| Tech Lead | peninjau kesimpulan dan jalur distilasi |

## Terkait

- [research/](../../../research/) — hasil riset distilasi ke `fnd-` (EKA 14.1).
- [decisions/](../../../decisions/) — simpulan yang diadopsi menjadi `dec-`/`adr-`.
- [specifications/](../../../specifications/) — temuan yang terbukti menjadi `spec-`.
