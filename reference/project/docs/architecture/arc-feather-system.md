---
namespace: feather
type: arc
id: feather-system
instance-version: 1
revision: 1
content-state: approved
existence-state: active
dimension: architecture
author: Lukas Weber
created: 2026-06-15
updated: 2026-06-22
supersedes: []
derives-from:
  - req:publishing-core
depends-on: []
amends: []
validates: []
change-log:
  - date: 2026-06-15
    domain: existence-state
    from: "-"
    to: active
    by: Lukas Weber
  - date: 2026-06-15
    domain: content-state
    from: "-"
    to: draft
    by: Lukas Weber
  - date: 2026-06-18
    domain: content-state
    from: draft
    to: review
    by: Lukas Weber
  - date: 2026-06-22
    domain: content-state
    from: review
    to: approved
    by: Lukas Weber
---

# Architecture — Feather System

## Purpose

The system architecture of Feather: the component structure, data flow, and technical constraints that bind all implementation below. It fulfills the publishing-core requirement with a deliberately boring, single-binary design.

## Content

**Components:**

| Component | Responsibility |
|---|---|
| **Feather server** (Go, single binary) | HTTP API, markdown rendering (Goldmark), publishing commands, search via SQLite FTS5 |
| **Content store** (`content/posts/`) | Post Markdown files in Git — the source of truth for post content |
| **SQLite database** (`feather.db`) | Metadata index: post titles, slugs, status, timestamps, FTS5 search index, revisions |
| **Static assets** (`public/`) | Theme CSS/JS; rendered pages are served directly |
| **Caddy** | Reverse proxy, TLS, static file serving in production |

**Data flow (publish):**

```
author → editor (autosave → content/posts/<slug>.md)
      → POST /api/posts/{id}/publish
      → server updates SQLite (status, published_at) + FTS index
      → Caddy serves the public page
```

**Constraints:**

- One process, one database file, one content directory. No queues, no workers, no external services.
- The Markdown file is the source of truth for content (see `adr:content-storage`); the database is a derived index and can be rebuilt from the files.
- Deployment is a file copy plus a restart behind Caddy (see `run:deploy-feather`).

The decisions behind this shape: `adr:content-storage`, `adr:search-sqlite-fts`, `adr:plugin-model-deferred`, `dec:reverse-proxy`.
