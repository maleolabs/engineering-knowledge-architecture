# Sessions

Implementation session contexts — ephemeral working documents from implementation sessions. Sessions capture the context, notes, and verification of a specific implementation effort.

## What Goes Here

- Implementation session folders
- Session context files
- Session notes and verification records

## Folder Structure

Each session is a folder: `impl-<work-items>-<date>/`

Inside each session folder:

| File | Purpose |
|---|---|
| `context.md` | Session context — what is being implemented and why |
| `notes.md` | Working notes — observations, decisions, blockers |
| `verification.md` | Verification results — how the work was validated |

### Example

```
impl-ts-012-001-st-012-001-2025-01-15/
├── context.md
├── notes.md
└── verification.md
```

## Naming Convention

| Type | Pattern | Example |
|---|---|---|
| Session folder | `impl-<work-items>-<date>/` | `impl-ts-012-001-2025-01-15/` |
| Context file | `context.md` | (fixed name) |
| Notes file | `notes.md` | (fixed name) |
| Verification file | `verification.md` | (fixed name) |

All filenames inside sessions are **lowercase**.

## Status Lifecycle

```
Active → Completed → Archived
```

- **Active**: Session is in progress.
- **Completed**: Session work is done; verification recorded.
- **Archived**: Session is preserved for reference.

## Ownership

| Role | Responsibility |
|---|---|
| Engineers | Create and maintain session documents |
| Tech Lead | Reviews session outputs |

## Related Folders

- `../work-items/` — Sessions implement work items
- `../reviews/` — Sessions are validated by reviews
