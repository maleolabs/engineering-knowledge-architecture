---
namespace: feather
type: gls
id: feather-terms
instance-version: 1
revision: 1
content-state: approved
existence-state: active
dimension: vocabulary
author: Maya Patel
created: 2026-06-11
updated: 2026-06-18
supersedes: []
derives-from: []
depends-on:
  - arc:feather-system
amends: []
validates: []
change-log:
  - date: 2026-06-11
    domain: existence-state
    from: "-"
    to: active
    by: Maya Patel
  - date: 2026-06-11
    domain: content-state
    from: "-"
    to: draft
    by: Maya Patel
  - date: 2026-06-15
    domain: content-state
    from: draft
    to: review
    by: Maya Patel
  - date: 2026-06-18
    domain: content-state
    from: review
    to: approved
    by: Maya Patel
---

# Glossary — Feather Terms

## Purpose

The canonical terms of the Feather domain. Definitions here are the reference for every other artifact; cross-role conversation (product, engineering, operations) stays unambiguous because the meaning of each term is fixed.

## Content

| Term | Definition | Non-term (what it is not) |
|---|---|---|
| **post** | A unit of published content: one Markdown file in `content/posts/` plus its index row. | Not a database row; the file is the source of truth. |
| **draft** | A post that exists as a file with index status `draft`; visible only to the author (and other authenticated users). | Not a branch of Git; drafts are regular files. |
| **publish** | The action of flipping a post's index status from `draft` to `published`, making it public. | Not a deployment; publishing never rebuilds the binary. |
| **revision** | A numbered version of a post's content; incremented on every successful save after the first. | Not the file mtime; revisions are explicit counters. |
| **autosave** | Background saving of the editor buffer without an explicit save action; guarantees at most one keystroke of loss. | Not version control; autosave is the safety net, Git is the history. |

Usage example: "The draft for `hello-world` was published as revision 3; the autosave kept the buffer safe while the title was being rewritten."
