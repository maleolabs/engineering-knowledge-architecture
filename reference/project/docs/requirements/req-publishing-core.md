---
namespace: feather
type: req
id: publishing-core
instance-version: 1
revision: 1
content-state: approved
existence-state: active
dimension: requirements
author: Maya Patel
created: 2026-05-12
updated: 2026-05-25
supersedes: []
derives-from: []
depends-on:
  - vis:feather-vision
  - str:feather-2026
amends: []
validates: []
change-log:
  - date: 2026-05-12
    domain: existence-state
    from: "-"
    to: active
    by: Maya Patel
  - date: 2026-05-12
    domain: content-state
    from: "-"
    to: draft
    by: Maya Patel
  - date: 2026-05-19
    domain: content-state
    from: draft
    to: review
    by: Maya Patel
  - date: 2026-05-25
    domain: content-state
    from: review
    to: approved
    by: Maya Patel
---

# Requirement — Publishing Core

## Purpose

The core publishing flow of Feather: the minimal set of authoring behaviors every writer needs, derived from the vision (plain files, author experience first) and the 2026 strategy.

## Content

**Statement.** An author can create a post as a Markdown draft, edit it, preview it, and publish it to the public site; published posts can be edited again, and every save while writing is persisted automatically (autosave).

Acceptance conditions:

1. **Draft → publish.** A draft can be published with one action; publishing makes the post appear on the public site immediately.
2. **Publish → edit.** A published post can be edited; saving an edit updates the public post and records a new revision.
3. **Autosave.** Content typed in the editor is saved without an explicit save action; a page reload never loses more than the last keystroke.
4. **Content as files.** Every post is stored as a Markdown file; the file is the source of truth (details in `adr:content-storage`).

Fulfillment is verified by `rvw:publishing-core-review` and the work items of the publishing-core epic (`epc:authoring-experience`).
