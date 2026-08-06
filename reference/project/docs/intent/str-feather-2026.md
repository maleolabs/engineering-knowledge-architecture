---
namespace: feather
type: str
id: feather-2026
instance-version: 1
revision: 1
content-state: approved
existence-state: active
dimension: intent
author: Maya Patel
created: 2026-05-06
updated: 2026-05-20
supersedes: []
derives-from:
  - vis:feather-vision
depends-on: []
amends: []
validates: []
change-log:
  - date: 2026-05-06
    domain: existence-state
    from: "-"
    to: active
    by: Maya Patel
  - date: 2026-05-06
    domain: content-state
    from: "-"
    to: draft
    by: Maya Patel
  - date: 2026-05-13
    domain: content-state
    from: draft
    to: review
    by: Maya Patel
  - date: 2026-05-20
    domain: content-state
    from: review
    to: approved
    by: Maya Patel
---

# Strategy 2026 — Feather

## Purpose

The strategy for Feather's 2026 cycle: where the product is heading and what is deliberately deprioritized. It operationalizes the vision into an order of investment.

## Content

**Author experience first, distribution second.**

1. **Win the writing loop (Q2–Q3).** Draft → publish → edit with autosave and live preview must feel instant. This is the moat: writers stay for the editor, not the hosting.
2. **Distribution later, and only the cheap forms (Q4).** RSS and sitemap ship before any social or newsletter integration. No engagement features in 2026.
3. **One contributor profile.** Feather optimizes for the individual writer and a small team of authors; multi-site and role management are out of scope for the year.
4. **No plugin model in v1.** Extensibility is a future problem; a plugin runtime would consume the whole engineering budget.

Strategic bets derived from this document: the publishing core requirement (`req:publishing-core`) and the deferred comments requirement (`req:comments-phase2`).
