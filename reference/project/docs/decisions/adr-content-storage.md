---
namespace: feather
type: adr
id: content-storage
instance-version: 1
revision: 1
content-state: accepted
existence-state: active
dimension: decisions
author: Lukas Weber
created: 2026-06-16
updated: 2026-06-24
supersedes: []
derives-from:
  - req:publishing-core
depends-on:
  - fnd:markdown-editor-options
amends: []
validates: []
change-log:
  - date: 2026-06-16
    domain: existence-state
    from: "-"
    to: active
    by: Lukas Weber
  - date: 2026-06-16
    domain: content-state
    from: "-"
    to: proposed
    by: Lukas Weber
  - date: 2026-06-24
    domain: content-state
    from: proposed
    to: accepted
    by: Lukas Weber
---

# ADR — Post Storage: Markdown Files with SQLite Index

## Context

The publishing-core requirement demands drafts, publishing, editing, revisions, and autosave. The storage choice determines how those behave: where content lives, how revisions are kept, and what must be rebuilt if the index is lost. The research on editor approaches (`fnd:markdown-editor-options`) established that authors write raw Markdown, so storage can treat files as the canonical form.

## Decision

Store every post as a **Markdown file** under `content/posts/<slug>.md`; keep a **SQLite database** holding only metadata (slug, title, status, timestamps, revision count) plus the FTS5 search index. The file is the source of truth for content; the database is a derived, rebuildable index.

## Consequences

- Posts are diffable, reviewable, and portable in Git; a writer can edit a post with any tool.
- Revisions are Git history for content plus a revision counter in the index; the database can be rebuilt from the files (`feather reindex`) at any time.
- Autosave writes files, so the editor needs debounced writes and a write lock to avoid concurrent-clobber (handled in `sto:draft-autosave`).
- Publish is a metadata flip (file → published), which keeps the file and the public state separable.
- Cost: metadata queries (listing, search) always touch SQLite; there is no alternative query path.

## Alternatives Considered

- **Database-only storage** — rejected: content locked in a schema, no Git diffing, contradicts the vision of plain files.
- **Full CMS (WYSIWYG + DB)** — rejected by `fnd:markdown-editor-options` and by the single-binary constraint.
