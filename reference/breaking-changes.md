# Breaking Changes — Legacy Structure → EKA v1.0

Summary of the 14 breaking changes against the legacy structure. Each change lists: old → new, architectural rationale (EKA anchor), impact, and mitigation.

| # | Old → New | Rationale (EKA anchor) | Impact | Mitigation |
|---|---|---|---|---|
| 1 | 16 mixed-dimension folders → 12 dimension folders + `operating/` + `exchange/` | Knowledge Taxonomy 1:1 (EKA 8); P1 | All paths change | `migration-guide.md` |
| 2 | Shared `mvp-*` ID space (4 types) → explicit type tokens `scp-`/`plan-`/`ctr-`/`tkt-` | Identity (EKA 6.4); P3 | Legacy globs `mvp-*` break | `adr-001` + token table |
| 3 | `roadmap/` (misnomer) → `plan-` in `planning/` | Plan Artifact + Planning State (EKA 7.2/10) | Roadmap consumers break | `adr-002` + status mapping |
| 4 | `sprints/` + 3-way status → `ctr-` + single-writer projections | P6, EKA 7.4 | Sprint-table tooling must read work item frontmatter | `adr-003` + projection regeneration |
| 5 | `tickets/` wave docs → `tkt-` projections with empty State Vector | EKA 7.4 | Execution commands read the new `tkt-` format | `adr-003` |
| 6 | Separate `adr/` + `decisions/` → single `decisions/` folder, two types | Decisions dimension (EKA 8) | Paths change | mapping table |
| 7 | `specification-corpus/` (misnomer) → `vocabulary/` (`gls-`); actual specs → `specifications/` (`spec-`) | Vocabulary ≠ Specifications (EKA 8) | Glossary moves; spec content moves | mapping table |
| 8 | `planning/` catch-all + `work-items/planning/` → dissolved | EKA 14.2; P15 | Content distributed to `trc-`/`plan-`/correct work item types | `adr-005` |
| 9 | Mixed `operations/` → `operations/` (procedures) + `standards/` (conventions) | Operational vs Standards (EKA 8) | Exit code/output conventions move | mapping table |
| 10 | Metadata table (Status/Author/Created/Updated/Version) → frontmatter state domains | Explicit State (P2); D2.8 | Tooling reading `Status:` breaks | `adr-002` + value mapping |
| 11 | Status values "In Progress"/"In Review" → `in-progress`/`in-review` | Value sets EKA 7.2 | Old strings invalid | value mapping in `validation.md` |
| 12 | Phase folders → `phase` field on `scp-`/`plan-` only | Phase as context (EKA 11.2) | `mvp/` folder gone | `adr-004` |
| 13 | Single `documentation-guide.md` → separate standard/serialization/protocol | EKA 1.3 | Guide readers now read 3 documents | `reference-architecture.md` |
| 14 | `sp-` (spike) → `spk-` (anti-prefix token) | EKA 6.4; ambiguity-free tokens | Glob `sp-*` breaks | token table |

## Tooling notes

- **All breaking changes are intentional**: legacy consumers **must** break so they cannot read Identity/state from location (EKA 6.4, P9). Silent compatibility would repeat the same violation: Identity encoded in structure.
- **New tooling reads frontmatter only**: Identity (`namespace`, `type`, `id`, `instance-version`, `revision`), State (`content-state`, `execution-state`, `planning-state`, `container-state`, `existence-state`), and Relationship (`supersedes`, `amends`, `derives-from`, `depends-on`, `validates`) — all from frontmatter, never from paths.
- **Filename for human navigation + consistency validation**: the `<type-token>-<id>[-v<nn>]` pattern aids browsing and enforces determinism, but is not the source of truth; mechanical validation checks filename ↔ frontmatter consistency.
- Detailed migration references: [`migration-guide.md`](migration-guide.md); full conventions: [`reference-architecture.md`](reference-architecture.md).
