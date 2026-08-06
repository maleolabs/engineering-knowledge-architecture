---
namespace: feather
type: vis
id: feather-vision
instance-version: 1
revision: 1
content-state: approved
existence-state: active
dimension: intent
author: Maya Patel
created: 2026-05-04
updated: 2026-05-18
supersedes: []
derives-from: []
depends-on: []
amends: []
validates: []
change-log:
  - date: 2026-05-04
    domain: existence-state
    from: "-"
    to: active
    by: Maya Patel
  - date: 2026-05-04
    domain: content-state
    from: "-"
    to: draft
    by: Maya Patel
  - date: 2026-05-11
    domain: content-state
    from: draft
    to: review
    by: Maya Patel
  - date: 2026-05-18
    domain: content-state
    from: review
    to: approved
    by: Maya Patel
---

# Vision — Feather

## Purpose

The vision of the Feather product: why the project exists and the principles that constrain every decision below it. It is the root justification reference for all other dimensions.

## Content

Feather is a markdown blogging platform for people who think in plain text. Writing is thinking; the platform should get out of the way of both.

Three principles follow:

1. **Content is plain files.** Posts are Markdown files in Git — portable, diffable, versioned, and never locked in a database schema.
2. **Small and boring.** A single Go binary, SQLite for the index, static output behind Caddy. No plugin runtime, no admin panel, no queues.
3. **Author experience first.** The editor, autosave, and instant preview matter more than traffic features. Distribution (RSS, sitemap) is built later and kept simple.

The platform targets the solo writer and small teams — not enterprises, not agencies. Success means a writer never thinks about the platform at all.
