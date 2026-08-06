---
namespace: feather
type: spk
id: markdown-syntax-extension
instance-version: 1
revision: 1
execution-state: planned
existence-state: active
dimension: research
author: Jonas Berg
created: 2026-08-05
updated: 2026-08-05
supersedes: []
derives-from: []
depends-on:
  - ctr:wave-7
  - plan:roadmap-v1:1
amends: []
validates: []
change-log:
  - date: 2026-08-05
    domain: existence-state
    from: "-"
    to: active
    by: Jonas Berg
  - date: 2026-08-05
    domain: execution-state
    from: "-"
    to: planned
    by: Jonas Berg
---

# Spike — Markdown Syntax Extensions (Tables, Footnotes)

## Description

Time-box (3 days): decide whether the renderer enables GFM table syntax and footnote extensions on top of the base Goldmark set. The question matters because enabling extensions is a one-line flag, but disabling them later is a content migration. Scope: evaluate tables, footnotes, strikethrough, task lists — the four most-requested GFM additions — for rendering quality, security surface, and round-trip fidelity in the editor preview.

## Investigation Notes

- Reference corpus: 30 real-world posts (tables in 12, footnotes in 4, strikethrough in 7, task lists in 2).
- Evaluate each extension on: HTML correctness, XSS surface (attribute/URL handling), preview-vs-public parity, and Goldmark's extension maturity (maintenance status, known issues).
- Record output samples for each extension in the finding; no decision is made by this spike — the decision belongs to the distillation artifact.

## Conclusion

Not yet concluded — investigation pending. **Distillation destination (EKA 11.4):** the findings will be written up as a research finding (`fnd-` in `docs/research/`); the adopted extension set will be recorded by a follow-up decision (`adr-` or `dec-` in `docs/decisions/`) that the renderer story (`ts:markdown-renderer`) can depend on. This spike must not be archived until the distillation artifact exists.
