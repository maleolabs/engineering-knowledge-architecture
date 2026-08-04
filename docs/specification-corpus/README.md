# Specification Corpus

Specification vocabulary — canonical definitions, lifecycle models, and shared terminology. This folder contains the authoritative references for terms and concepts used across the project.

## What Goes Here

- Canonical term definitions
- Lifecycle models
- Shared vocabulary
- Glossary
- Specification references

## What Does NOT Go Here

- Architecture descriptions → `../architecture/`
- Product requirements → `../prd/`
- Architecture decisions → `../adr/`

## Naming Convention

| Type | Pattern | Example |
|---|---|---|
| Vocabulary | Descriptive kebab-case | `vocabulary.md` |
| Lifecycle model | Descriptive kebab-case | `lifecycle-model.md` |
| Glossary | Descriptive kebab-case | `glossary.md` |

Use descriptive names that clearly indicate the document's content.

## Ownership

| Role | Responsibility |
|---|---|
| Tech Lead | Maintains canonical definitions |
| All contributors | Reference and use shared terminology |

## Related Folders

- `../architecture/` — Architecture uses specification terms
- `../manifesto/` — Manifesto defines product principles
- `../adr/` — ADRs may define new terminology
