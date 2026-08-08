# Feather

**Feather is a markdown blogging platform that gets out of the way.** Write in Markdown, publish with one command, and let readers subscribe over RSS — no database to administer, no admin panel to learn, no build pipeline to babysit. Feather stores every post as a plain Markdown file in Git, keeps a lightweight SQLite index for search and listings, and serves a fast static site through Caddy. It is a small Go codebase with a deliberately boring architecture, so the interesting work stays in the writing.

## Quick Start

```sh
# Build the server (Go 1.24+)
go build -o feather ./cmd/feather

# Run locally — serves http://localhost:8080, content/posts/ as the post store
./feather serve --addr :8080

# Create a post as a Markdown file and publish it
./feather new "Hello, world"        # writes content/posts/hello-world.md
./feather publish hello-world       # moves it through draft -> published
```

## Knowledge Base

This repository is an **EKA v1.1 serialization** (Git + Markdown). The engineering knowledge behind Feather — vision, requirements, architecture, decisions, planning, execution, and operations — lives in [docs/](docs/README.md), which is the source of truth for the structure. Read [docs/workflow-guide.md](docs/workflow-guide.md) first if you are new to the conventions.

```sh
# Validate the knowledge base against the conformance rules (R0–R12)
eka validate

# Seed the EKA workspace canonical store (registers this repository)
eka sync

# See the current execution picture: active container, work items by state
# (projections read the workspace store, never the docs tree directly)
eka view execution
```

## Status

Feather is under active engineering effort (May–August 2026), roughly halfway to the MVP: the publishing core (draft → publish → edit, autosave) is in flight in the active execution container [Wave 7](docs/operating/containers/ctr-wave-7.md). See the [roadmap](docs/planning/plan-roadmap-v1-v1.md) for the committed plan and [the traceability matrix](docs/planning/trc-feather-traceability.md) for how the knowledge hangs together.

## Contributing

Work is tracked as EKA work items inside the active container. Every change to knowledge artifacts must pass `eka validate` before commit — see [docs/exchange/validation.md](docs/exchange/validation.md) for the mechanical checklist.
