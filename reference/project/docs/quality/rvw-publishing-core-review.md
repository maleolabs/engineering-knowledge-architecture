---
namespace: feather
type: rvw
id: publishing-core-review
instance-version: 1
revision: 1
content-state: review
existence-state: active
dimension: quality
author: Lukas Weber
created: 2026-08-05
updated: 2026-08-05
supersedes: []
derives-from: []
depends-on:
  - plan:roadmap-v1:1
amends: []
validates:
  - sto:publish-post
change-log:
  - date: 2026-08-05
    domain: existence-state
    from: "-"
    to: active
    by: Lukas Weber
  - date: 2026-08-05
    domain: content-state
    from: "-"
    to: draft
    by: Lukas Weber
  - date: 2026-08-05
    domain: content-state
    from: draft
    to: review
    by: Lukas Weber
---

# Review — Publishing Core

## Purpose

Verify that the publishing core satisfies `req:publishing-core` (draft → publish → edit, autosave) and that the delivered work items meet the definition of done (`std:definition-of-done`). Scope: `sto:publish-post` (validated by this review), `sto:draft-autosave` and `ts:markdown-renderer` (in-flight — reviewed partially), plus the regression fix `bug:empty-title-crash`.

## Content

Reviewed against the specification (`spec:publishing-api`) and the storage decision (`adr:content-storage`):

- Publish flow: file remains the source of truth; index status flips; FTS index updated; public page renders. Verified end-to-end on a staging instance behind Caddy.
- Autosave: debounce, reload restore, and atomic writes verified in `ses:2026-08-04-authoring-session`; stale-`If-Match` reconciliation still open (work item `sto:draft-autosave`, in-progress).
- Renderer: golden-file corpus green; XSS probe corpus clean; preview/public parity holds.
- Empty-title regression: 400 + process alive; regression test committed.

## Findings

1. **Blocking (resolved in review):** none outstanding on `sto:publish-post` — publish flow meets acceptance criteria and DoD.
2. **Open (autosave):** stale-`If-Match` save loop (session finding) must be resolved before `sto:draft-autosave` can leave in-progress.
3. **Open (renderer):** unknown-language highlighting degrades to plain `<pre>` — acceptable, recorded as a note for the syntax-extension spike (`spk:markdown-syntax-extension`).
4. **Observation:** `GET /api/posts` N+1 pattern (tracked as `td:reduce-query-count`, planned) — confirmed not blocking at MVP scale.

## Action Items

- Move `sto:publish-post` to done if not already; sign the review gate for the publish path (already done by the owner).
- Track the autosave reconciliation finding on `sto:draft-autosave` (owner: Jonas Berg).
- Re-run this review when `sto:draft-autosave` and `ts:markdown-renderer` reach in-review (expected next week).
- This review moves to `approved` once the open findings are closed; it then serves as the release gate evidence for `rel:v090`.
