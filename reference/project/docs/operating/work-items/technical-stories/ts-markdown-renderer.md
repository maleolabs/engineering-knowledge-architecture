---
namespace: feather
type: ts
id: markdown-renderer
instance-version: 1
revision: 1
execution-state: in-review
existence-state: active
dimension: architecture
author: Jonas Berg
created: 2026-07-29
updated: 2026-08-04
supersedes: []
derives-from: []
depends-on:
  - ctr:wave-7
  - plan:roadmap-v1:1
amends: []
validates: []
change-log:
  - date: 2026-07-29
    domain: existence-state
    from: "-"
    to: active
    by: Jonas Berg
  - date: 2026-07-29
    domain: execution-state
    from: "-"
    to: planned
    by: Jonas Berg
  - date: 2026-07-31
    domain: execution-state
    from: planned
    to: todo
    by: Jonas Berg
  - date: 2026-08-01
    domain: execution-state
    from: todo
    to: in-progress
    by: Jonas Berg
  - date: 2026-08-04
    domain: execution-state
    from: in-progress
    to: in-review
    by: Jonas Berg
---

# Technical Story — Markdown Renderer with Syntax Highlighting

## Description

Replace the placeholder renderer with a production one: Goldmark with syntax highlighting (chroma) and safe-by-default HTML (no raw HTML unless explicitly allowed). The renderer is a pure function `(markdown, template) → html` — no state, no plugins (per `adr:plugin-model-deferred`) — used by both the preview pane and the public page so preview and production always agree.

## Acceptance Criteria

- Rendering is deterministic: identical input → identical HTML, verified by a golden-file test corpus.
- Fenced code blocks get syntax highlighting; unknown languages degrade to plain `<pre>`.
- Raw HTML in post bodies is escaped by default; an explicit per-post flag can allow it.
- Rendering is safe against malicious input: no script injection via attributes or links (test corpus includes XSS probes).
- Performance: rendering the largest test post (< 100 KB) completes in < 50 ms; results cached per revision hash.
- Review gate passed (`in-review` held until `rvw:publishing-core-review` records the verdict).
