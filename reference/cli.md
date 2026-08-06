# EKA CLI — Official Interface (`eka`)

> Convention document, not an artifact. Meta-documentation of the `reference/` zone.
> Implementation: Go (`cmd/` + `bootstrap/` + `conformance/`), module `github.com/maleolabs/engineering-knowledge-architecture`.
> Related: [`conformance-notes.md`](conformance-notes.md) (interpretation decisions + rule traceability), [`../skeleton/docs/exchange/validation.md`](../skeleton/docs/exchange/validation.md) (Conformance Rules R0–R12), [`eka-reference-serialization-format-v1.1.md`](eka-reference-serialization-format-v1.1.md) (the serialization format used by `eka export` and `eka import`).

## CLI philosophy

`eka` is the **canonical executable form of the Conformance Rules** (Naming and Terminology Specification §7.5) — the official interface of Engineering Knowledge Architecture for humans and agents (Naming and Terminology Specification v1.1 §7). Current roles:

- **`eka init`** — the official Repository Bootstrapper: analyzes the workspace, composes a bootstrap plan, configures the project interactively, generates an EKA repository from the Reference Skeleton, then validates it.
- **`eka export`** — the first practical implementation of the Exchange Specification: exports engineering knowledge as an **Exchange Package** per the Reference Serialization Format (RSF).
- **`eka import`** — the inverse of export: consumes an Exchange Package and integrates knowledge into an existing EKA repository — implementing the import semantics of Exchange Specification §11.
- **`eka validate`** — the conformance validator: repository conformance must not rest on manual review alone — rules R0–R12 in `skeleton/docs/exchange/validation.md` are designed to be mechanical, and this validator is their canonical implementation (P16: enforcement mechanisms vary, invariants stay identical).
- **`eka view`** — the Knowledge Projection Engine: read-only projections of the Engineering Knowledge Model (the five domain projections `discovery` / `architecture` / `planning` / `execution` / `operations` + the `ticket` projection), rendered as per-domain visualizations — Kanban board (execution), roadmap (planning), dependency tree (architecture), information cards (discovery), release timeline (operations), detail card (ticket) — the canonical executable form of the State Projection semantics (Core Specification §11), relationship-derived, never markdown-rendered.

New to EKA? Start with the [Engineering Operating Guide](../skeleton/docs/workflow-guide.md) — the primary onboarding document (mental model, lifecycle, domains, workflows).

Consequences of this philosophy:

- The validator is the **single source of mechanical truth**; where rule text is ambiguous, interpretation decisions are documented (see `conformance-notes.md`) — no behavior is invented without a basis.
- The EKA repository itself must pass its own validator (see [Repository conformance](#repository-conformance)) — a prerequisite for the validator to be trusted by other repositories.
- Deterministic CLI behavior: identical input produces identical output; non-TTY output is byte-identical and free of ANSI control sequences (see [CLI UX](#cli-ux)).
- The CLI is an **adapter** — business logic (bootstrap, validation) lives in reusable application packages, independent of the CLI framework (see [CLI architecture](#cli-architecture)).

## CLI UX

All command output is rendered by `cmd/ui` — a presentation subpackage of `cmd` with no business logic. Every renderer is a pure function of (data, style); the only time-dependent output is the TTY-only spinner animation, which always ends in a deterministic final state.

### Communication philosophy

Calm, professional, and unadorned:

- No exclamation marks, no ALL-CAPS emphasis, no banners, no ASCII logos.
- Meaning is carried by words and icons; color is decoration, never the message.
- Trust through clarity: every line answers "what is happening" and "what was the outcome".

### Interaction model

All commands share one three-part hierarchy:

1. **Context header** — identifies the object being processed.
2. **Workflow body** — the operation's stages (progressive tree) or, for single-operation commands, the report.
3. **Summary** — the outcome as facts.

`init`, `export` and `import` render a progressive tree; `validate` renders the report as the body; `view` renders the projection as the body — each projection is a per-domain visualization (board, roadmap, tree, cards, timeline, detail card). Every command ends with a summary block.

### Context header

The first lines of every command orient the reader on the object, not the action:

```
Repository
Name        myproj
Namespace   eka-cli
Knowledge   EKA v1
↓ Bootstrap
```

- First line: object kind (`Repository`, `Knowledge Package`) — the dynamic identity anchor.
- Identity rows: aligned label/value pairs; each command shows the rows that identify its object (`Name`, `Namespace`, `Path`, `Package`, `Scope`, `Output`, `Format`, `Knowledge`).
- Pipeline separator: `↓ <pipeline>` — `↓ Bootstrap`, `↓ Export`, `↓ Import`, `↓ View`, `↓ Validate`.

### Object-centric execution

Every workflow states the object it processes before the actions: `init` bootstraps a **Repository**, `export`/`import` move a **Knowledge Package**, `validate` checks a **Repository**. The object kind opens the output and is repeated as the tree root.

### Progressive workflow tree

Multi-step commands render stages as a tree with deterministic `[i/n]` prefixes (`[1/6] Discover repository`):

- **TTY** — the tree redraws in place: completed steps show `✓`, the active step shows the spinner, pending steps stay dim; `Finish` leaves the cursor on a fresh line.
- **Non-TTY** — steps emit deterministic sequential lines as they complete (`├── <label>`, `│   ✓ <detail>`); no redraw, no animation.
- Failure renders `failed: <detail>` under the leaf — the word, never color alone.

### Contextual loading

Loading states always describe the work in progress ("Loading Engineering Knowledge...") — a bare spinner never appears. Export's load phase is the current example: the message prints once on non-TTY; on a TTY the animation ends in the deterministic line `✓ <message>`. Loading states are only shown for operations taking roughly a second or more.

### Color semantics

A soft, hand-rolled 256-color palette — the only colors the CLI may emit:

| Role | ANSI SGR | Usage |
|---|---|---|
| Info | `38;5;75` | labels, headings (soft blue) |
| Success | `38;5;114` | completed steps, `✓`, PASS (soft green) |
| Warning | `38;5;214` | warnings (amber) |
| Error | `38;5;167` | failures, FAIL (muted red) |
| Progress | `38;5;80` | active spinner frame (soft cyan) |
| Dim | `38;5;245` | secondary/detail lines (gray) |

Color auto-disables when the writer is not a terminal, `NO_COLOR` is set, or `TERM=dumb`. Color is never the sole carrier of meaning: a failed step prints `failed: ...`, and PASS/FAIL are always printed as words.

### Icon semantics

A deliberately minimal Unicode set — no emojis:

| Icon | Meaning |
|---|---|
| `✓` | completed step, finished loading |
| `⠋ ⠙ ⠹ ⠸ ⠼ ⠴ ⠦ ⠧ ⠇ ⠏` | Braille spinner cycle (TTY only, progress color) |
| `├──` / `└──` / `│` | tree connectors (non-TTY tree lines, summary items) |
| `•` | detail-list bullet (verbose sections) |
| `→` | relationship direction (verbose export external references) |
| `↓` | pipeline separator in the context header |

All glyphs are valid UTF-8; icons decorate, text carries meaning.

Projection renderers additionally use Unicode box-drawing characters (`┌ ─ ┬ ┐ ├ ┼ ┤ └ ┴ ┘`) for visualization frames — Kanban board, cards, tree. They carry layout, never meaning; they are plain UTF-8 text on non-TTY output, never ANSI.

### Summaries

Every successful command ends with a `Summary:` block: a tree-style list of `└── label: value` outcome facts — project name, namespace, git status, validation verdict, counts, digests — never a log replay. On failure, diagnostics go to stderr and stdout stays empty (import contract).

### Verbose mode

`-v` / `--verbose` (persistent, presentation-only) adds detail sections between the tree and the summary: plan actions and per-unit lists for `init`, per-unit identities/attachments/external references for `export`, imported/skipped/warning lists for `import`. It never changes the interaction model, the exit codes, or determinism. Default output is concise.

### Determinism

The contract: **identical input → identical output**.

- Non-TTY output (pipes, CI, tests) is byte-identical across runs: plain text plus UTF-8 icons, no ANSI escapes, no `\r`, no spinner frames.
- TTY-only affordances — animation, in-place redraw, color — never leak into piped or redirected output.
- The final rendered tree is deterministic on both TTY and non-TTY; the only time-dependent bytes are the intermediate TTY animation frames.

### Accessibility

The CLI is usable without color, over SSH, in CI, and in basic terminals: the plain non-TTY output is the complete output, not a degraded fallback. `NO_COLOR` and `TERM=dumb` disable color on a TTY without changing anything else.

### Consistency

One interaction model across `init` (5-stage tree), `export` (6-stage tree), `import` (7-stage tree), `validate` (header + report + summary) and `view` (header + projection + summary). Future commands must follow the same model (see [Contribution guide](#contribution-guide-adding-a-command)).

## Installation

Prerequisites: **Go 1.24+**.

### From the module (after publication)

```sh
go install github.com/maleolabs/engineering-knowledge-architecture/cmd/eka@latest
```

The `eka` binary installs to `$GOBIN` (default `$GOPATH/bin`).

### From source (this repo)

```sh
cd <root-repo-eka>
go build -o eka ./cmd/eka
```

The build result is a **standalone, portable binary** — no runtime dependencies besides the binary itself; it can be copied to another machine (same architecture/OS) and run directly.

## Usage

```
eka [--verbose|-v]
eka version
eka init [project-name] [--dry-run]
eka export [target...] [-o|--output path]
eka import <package-path>
eka view [projection] [target]
eka validate [path]
eka completion [bash|zsh|fish|powershell]
eka help [command]
```

- `eka` without a subcommand shows the **product landing** — a calm orientation (description, compact command overview, help and version pointers) — and exits `0`. No banner, no decoration.
- `eka version` prints the CLI build version (default `dev`; set at build time via `-ldflags "-X .../cmd.version=v1.2.3"`) and the EKA standard version implemented — current output: `dev (EKA standard 1.1)`.
- `eka -h` / `eka --help` / `eka help` prints the command reference and exits `0`.
- `eka help <command>` prints command help.

### Exit codes

| Code | Meaning | Example |
|---|---|---|
| `0` | **Success** — compliant validation (warnings allowed) / initialization completed / dry-run / help / completion | Repository passes R0–R12; `eka init` completed and validated |
| `1` | **Blocking violation** — validation found errors, or the repository produced by `eka init` failed validation | At least one error (severity `error`) |
| `2` | **Usage/internal error** — the command did not run | Unknown command, unknown flag, too many arguments, invalid path |

Warning semantics: **warnings never affect the exit code**. A repository with warnings still exits `0`.

## `eka init` — Repository Bootstrapper

`eka init` is not a template generator. It is the official **Repository Bootstrapper** for Engineering Knowledge Architecture knowledge repositories: it analyzes the workspace first, composes a plan, asks only what is not yet known, generates from the Reference Skeleton, and validates the result.

### Initialization philosophy

Four principles:

1. **Understand first, change later** — initialization runs in five stages: *Workspace Discovery → Bootstrap Planning → Interactive Configuration → Repository Generation → Validation*. The workspace is understood before any modification; no blind generation.
2. **Adaptive** — automatic discovery drives the wizard: the user is never asked a question whose answer is already known.
3. **Idempotent** — `eka init` is safe to run repeatedly; existing artifacts are detected and reused, skipped, or explicitly confirmed before replacement. User content is never silently overwritten.
4. **Automatically tested** — initialization counts as successful only if the generated repository passes `eka validate`.

### Five-stage workflow

```
Workspace Discovery
        ↓
Bootstrap Planning
        ↓
Interactive Configuration
        ↓
Repository Generation
        ↓
Validation
```

1. **Workspace Discovery** — inspects the target: presence/absence of the directory, Git directory (at the target or its ancestors), README, `docs/` directory, existing EKA repository (markers: `docs/operating/` + `docs/exchange/validation.md` + `docs/exchange/transfer.md`), and configuration files (informational: `.gitignore`, `.editorconfig`, `.eka.*`, `eka.*`).
2. **Bootstrap Planning** — composes a deterministic plan: directories to create, files to generate, resources to reuse, Git status, validation steps. The plan carries the content of the files to be written — dry-run and execution cannot diverge.
3. **Interactive Configuration** — adaptive wizard (see below).
4. **Repository Generation** — copies the Reference Skeleton (embedded from `skeleton/`, see Architecture) to the target: `docs/**` verbatim + `README.md` with the title replaced by the project name. Existing files with identical content are reused; differing ones are confirmed (interactive) or skipped (non-interactive). `git init` only if planned (see Git detection).
5. **Validation** — `conformance.Validate(target)` runs; results are printed. Initialization succeeds only if the repository passes.

### Modes

| Command | Behavior |
|---|---|
| `eka init` | Initializes the current directory |
| `eka init .` | Equivalent to `eka init` |
| `eka init <project-name>` | Creates a new project directory and initializes it as an EKA repository |
| `eka init [name] --dry-run` | Prints the bootstrap plan without writing anything; exits `0` |

### Adaptive wizard

The wizard only asks relevant questions:

| Question | Asked only if |
|---|---|
| Project Name | Target directory name is empty/unusable |
| Namespace | Always (required; lowercase, digits, hyphens; no `/`, `:`, spaces) |
| Project Description | Always (optional; may be empty) |
| Generate README | No README exists yet |
| Initialize Git | No Git exists **and** the `git` binary is available |

There is no "Methodology" question — EKA v1 has no canonical methodology taxonomy. If answers cannot be read (stdin is not a terminal or EOF), the wizard uses deterministic defaults and **never** runs `git init`.

### Smart Git detection

- Existing Git repository (at the target or an ancestor) → Git question skipped, `git init` not executed.
- New directory / no Git → initialization offered; if accepted: `git init`; if declined: proceed normally.
- Non-interactive (pipe, `/dev/null`, file, CI) → Git never initialized.
- `git init` failure (e.g., binary missing) → warning, initialization continues.

### Idempotency

Running `eka init` twice on the same target never damages the repository:

- Existing EKA repository → detected, all resources reused, only validation runs (output: "already initialized").
- Files to be generated already exist:
  - identical content → reused (no write);
  - different content → explicit confirmation (interactive; default reject) or skipped + reported (non-interactive).

### Dry-run

`eka init --dry-run` prints the bootstrap plan — directories to create, files to generate, resources to reuse, Git status, validation plan — **without writing a single file**. Deterministic: two dry-runs on the same state produce identical output.

### Output summary

On completion, a concise summary is printed: Project Name, Namespace, Repository Type (new / existing-dir / existing-eka), Git Status, Knowledge Standard Version (EKA v1.1), Validation Result (PASS/FAIL + error/warning counts), and suggested next steps.

## `eka export` — Knowledge Package Export

`eka export` is the first practical implementation of the **EKA Exchange Specification** and the **Reference Serialization Format (RSF)**. It is **not** a repository archive and **not** a ZIP utility — it builds a canonical Exchange Package.

### Export philosophy

Export is treated as a **knowledge transformation**, not file copying:

```
Repository
    ↓
Engineering Knowledge Model
    ↓
Reference Serialization Format
    ↓
Exchange Package
```

- The package represents **Engineering Knowledge** (Identity, State, Content, Relationship, Classification), not the repository layout — no repository paths inside the package; all identities canonical.
- The package is a bijective projection of the Exchange Package Object Model (Exchange Spec §4.4) onto the RSF.
- Identical repositories always produce **byte-identical** packages.

### Workflow

1. **Repository validation** — `conformance.Validate(root)` runs automatically (equivalent to `eka validate`). On failure: **export stops, no package produced** (exit `1`). Only conformant repositories can be exported.
2. **Discovery** — repository identity (namespace across all artifacts), specification version (EKA v1.1), scope, package metadata.
3. **Loading** — artifacts loaded via `conformance.Scan` (same scan policy as the validator: `.md` only, skip `testdata`/dot-dirs/symlinks); content bodies extracted byte-exact.
4. **Scope resolution** — select units per scope (see below).
5. **Model construction** — build the Exchange object model: Header, Manifest, Units, Declarations, Integrity.
6. **Serialization (RSF)** — deterministic projection: JSON blocks + content payloads + attachments.
7. **Write** — package `.ekapkg` (ZIP) or directory layout.

### Export scope

| Argument | Scope | Package contents |
|---|---|---|
| (no target) | **Repository** (default) | all artifacts of all Lines |
| `<type>:<id>` | **Line** | all instances of that Line |
| `<type>:<id>:<instance-version>` | **Instance** | one instance |
| multiple targets | **Collection** | union of the requested Lines/instances |

References leaving the package → **External Reference Declaration** in `declarations.json` (Exchange §12.3) — dependency integrity is preserved without closure traversal in v1. Dangling references on Draft artifacts are tolerated (draft tolerance, R5) and recorded; dangling references on non-Draft artifacts block validation (and therefore export).

### Package contents (reference RSF projection)

| Entry | Contents |
|---|---|
| `header.json` | Package Header: serialization version `1`, exchange format version `1`, specification version `1.0`, exporter `eka`, package identity label, scope, namespace. No creation timestamp (byte determinism). |
| `manifest.json` | Manifest: ordered unit list (canonical identity form), per-unit + package digests, counts, closure declaration (scope + seed). |
| `units/<ns>/<type>-<id>-v<nn>/unit.json` | Unit metadata: full identity, revision, author/created/updated, exact state vector, Change Log in occurrence order, ordered Relationships by Identity, classification, phase (scp-/plan-). |
| `units/<ns>/<type>-<id>-v<nn>/content` | Content body (representation `eka/structured-text/1`), byte-exact. |
| `attachments/<path>` | Non-`.md` files under `docs/` (diagrams, images, etc.), byte-exact. Attachment ID = relative path. v1: no unit→attachment references. |
| `declarations.json` | Declarations Block: closure declaration, external reference declarations, extension declarations (empty in v1). |
| `integrity.json` | SHA-256 digests: package-level (over all entries except `manifest.json` + `integrity.json`), per-unit (`unit.json ‖ content`), per-attachment. |

Encoding: UTF-8 without BOM, LF, fixed-order JSON structs, ordered ZIP entries with zero timestamps. Package name: `rsf-<scope>-<namespace>-1.ekapkg` (RSF §4.1–4.2).

### Documented RSF deviations (v1)

1. No creation timestamp in the header — byte determinism (RSF §4.3 permits differing metadata).
2. The header does not announce integrity/declarations — the presence of `integrity.json`/`declarations.json` is deterministic.
3. Content carried byte-exact without LF normalization — losslessness; the declared canonicalization = byte-exact payload (RSF §9.3/§6.3.3).
4. Attachment ID = repository-relative path, not "referring unit identity + resource name" (RSF §7.2 recommended rule).
5. The package digest does not cover `manifest.json` — avoids self-reference (the manifest carries the package digest).

### Output

- Default: `<label>.ekapkg` in the current directory.
- `-o <file>.ekapkg` — custom file path.
- `-o <dir>` / path ending in a separator — directory layout (same logical structure, without ZIP).

Both modes deterministic. The same package from identical repositories → identical.

### Determinism

- All collections ordered by canonical identity key; Change Log in occurrence order; relationships ordered (type, target); ZIP entries ordered; JSON fields in fixed order.
- No timestamps, no host-dependent values, no absolute paths inside the package.
- Digests computed over canonical bytes.

### Error handling

| Condition | Behavior |
|---|---|
| Repository non-conformant | Stops, validation report printed, **no package**, exit `1` |
| Target missing / invalid syntax / ambiguous | Error listing the available artifacts, exit `2` |
| Identity component violates charset (RSF §5.2.3) | Export refused (security: path traversal prevention), exit `2` |
| Serialization/fs failure (output not writable, etc.) | Error, exit `2` |

## `eka import` — Knowledge Integration

`eka import <package-path>` consumes an Exchange Package (`.ekapkg` or directory layout) and integrates engineering knowledge into the EKA repository in the current directory. It is **not** an archive extraction — it is a knowledge integration pipeline implementing the import semantics of Exchange Specification §11.

### Import philosophy

```
Exchange Package
    ↓
Reference Deserialization
    ↓
Exchange Model
    ↓
Repository Integration
    ↓
Repository Validation
```

- Identity = the canonical lookup mechanism — never filesystem paths.
- Integration is **atomic**: either everything succeeds, or nothing changes.
- v1 strategy = **conservative merge**: only new artifacts are written; identical duplicates = no-op; any payload difference = conflict → abort. No overwrite, no delete, no replace strategy (future).

### Workflow (Exchange §11 pipeline)

1. **Repository discovery + gate** — the target must be an EKA repository (markers `docs/operating/` + `docs/exchange/`); `conformance.Validate` runs before import. Invalid target → stop, exit `1`.
2. **Package validation** — package integrity (SHA-256 digests: package-level over all entries except `manifest.json` + `integrity.json`, per-unit over `unit.json ‖ content`, per-attachment), well-formed JSON, manifest↔units 1:1, unknown fields/entries rejected (RSF §9.5), version compatibility (serialization `1`, exchange format `1`, specification `1.0` — incompatibilities rejected with a "found vs supported" diagnostic).
3. **Phases 1–8 (analysis, no writes)** — contract → identity (charset RSF §5.2.3, unique within the package) → state (valid values, owned-set, change-log consistency) → structural (well-formedness per type family) → referential (local → global → external per §7.4; dangling non-draft → abort; draft → warning) → conflict → duplicate → dependency order.
4. **Phase 9 (commit)** — staged writes (temp + rename, atomic per file); any failure → rollback (created files and directories removed), repository unchanged.
5. **Phase 10 (revalidation)** — `conformance.Validate` after integration; failure → rollback, exit `1`.

### Conflicts

| Condition | Behavior |
|---|---|
| Same Identity + identical payload (content, state, change-log, relationship, classification, metadata) | Duplicate → no-op (skipped) |
| Same Identity + different payload | **Conflict → import aborted**; per-identity summary lists the differences |
| Undeclared reference / unresolved external (non-draft) | **Abort** — referential integrity cannot be maintained |
| Dangling reference on a Draft artifact | Tolerated (draft tolerance, R5), recorded as a warning |
| Target attachment already exists, different content | Conflict → abort; identical → skipped |

No partial integration: all analysis happens before a single file is written.

### Rollback guarantees

- Per-file staged commit; commit failure → all created files and directories removed (deepest-first), pre-existing files untouched.
- Post-import revalidation failure → the entire import result is rolled back, the repository returns to its pre-import state.
- If rollback itself fails (e.g., filesystem turned read-only), the error carries an explicit warning that the repository may be partially modified.

### Compatibility

- The package must declare serialization version `1`, exchange format version `1`, specification version `1.0` — supported.
- Unsupported combinations are rejected with a diagnostic (found vs supported values).
- Packages without a Capability Declaration are accepted (optional per Exchange §4.5); packages with unknown extensions are explicitly rejected.

### v1 limitations

- Cross-namespace references only for v1 instances (repository format `<ns>/<type>:<id>`; instance > 1 → clear error).
- No forward-only state reconciliation (Exchange §13.2) — v1 is conservative: identical = no-op, different = conflict.
- Advanced replace/merge strategies: future.

### Error handling

| Condition | Behavior |
|---|---|
| Target not an EKA repository / non-conformant | Stops, exit `1` |
| Identity/state/metadata conflict, relationship failure, revalidation failure | Stops + rollback, exit `1` |
| Invalid package (digest, JSON, missing entries, unsupported version, unknown fields) | Stops, repository unchanged, exit `2` |
| Filesystem / usage failure | Error, exit `2` |

## `eka view` — Knowledge Projections

```
eka view [projection] [target]
```

`eka view` projects the **Engineering Knowledge Model** of the repository rooted at the current directory: read-only views over the repository's artifacts and their relationships. Projections are named by **Engineering Domain** — `discovery`, `architecture`, `planning`, `execution`, `operations` — plus the `ticket` projection. The `target` argument is required by the ticket projection only (a bare ticket id, `tkt-<id>` or `tkt:<id>`); domain projections ignore it. With no arguments, the available projections are listed and the command exits `0`.

### Projection philosophy — Knowledge Projection

- The viewing layer is **generated from the Engineering Knowledge Model** — artifact identity, State fields, and relationships — **not from markdown rendering**. The CLI never depends on markdown formatting; no file text is parsed for view content (ticket bodies are never read, container `## Work Items` tables are never parsed).
- Markdown remains the **editing interface**; EKA is the **canonical model**; the projection is the **view**. A projection has no State of its own and never becomes a writer (P6, Core Specification §11) — `eka view` is the canonical executable form of the State Projection semantics.
- The projection engine is **pure data in, pure data out**: it contains no terminal knowledge, no command framework, and no output.

### Snapshot interaction model

`eka view` is a single synchronous snapshot — no TUI, no navigation, no editing, no background process. One run produces one projection, then exits:

```
run → validate → load → construct → render → exit
```

1. **Validate** — the conformance gate (R0–R12) runs first; a non-conformant repository is refused (exit `1`).
2. **Load** — one `conformance.Scan` of the repository.
3. **Construct** — one Knowledge Graph, then one projection build (the builder).
4. **Render** — the projection renderer applies the projection's visualization (board, roadmap, tree, cards, timeline, detail card), deterministically.
5. **Exit** — mapped per the exit-code contract below.

The engine is synchronous and stateless: a future loading state can wrap the whole call without restructuring.

### Exit codes

| Code | Meaning |
|---|---|
| `0` | Projection produced — including empty projections (no active container, no domain artifacts) |
| `1` | Repository validation failed — no projection is produced (the full report is printed) |
| `2` | Usage or internal error — unknown projection, missing or unknown ticket target, unreadable root |

Warnings never block a projection.

### Projections

The canonical projections — one per Engineering Domain, plus the ticket projection:

| Command | Engineering Domain | Content | Rendered as |
|---|---|---|---|
| `eka view discovery` | Discovery | `vis-` / `str-` / `req-` / `fnd-` artifacts — vision, strategy, requirements, research findings | Information cards |
| `eka view architecture` | Architecture | `adr-` / `dec-` / `arc-` / `spec-` / `std-` / `gls-` artifacts — decisions, architecture descriptions, specifications, standards, vocabulary; per-node **Content State** (`draft` / `review` / `approved` / post-approval terminals, incl. the ADR/decision variants) shown on the node | Dependency tree |
| `eka view planning` | Planning | `scp-` / `epc-` / `plan-` / `trc-` artifacts — scope definitions, epics, plans, traceability artifacts, sequenced by phase; per-entry **Planning State** (`draft` / `approved` / `immutable`) | Roadmap / timeline |
| `eka view execution` | Execution | The **active Execution Container**'s work items on the fixed five-column board (`planned` / `todo` / `in-progress` / `in-review` / `done`; empty columns keep their heading). The ticket list is **not rendered as a block** — tickets project individually via `eka view ticket`. | Kanban board |
| `eka view operations` | Operations | `run-` / `rel-` artifacts — runbooks, operational guides, release records | Release summary / activity timeline |
| `eka view ticket <id>` | Execution | One ticket's detail: **projected status derived from the referenced work item's owner Execution State — never from the ticket's own text** — plus the container it derives from and its `derives-from` references | Detail card |

**Visualization by domain.** Each projection renders as a purpose-built console for its domain's primary question:

| Engineering Domain | Visualization style | Primary question it answers |
|---|---|---|
| Discovery | Information cards | What are we building? |
| Architecture | Dependency tree | How is the system structured? |
| Planning | Roadmap / timeline | What are we planning next? |
| Execution | Kanban board (five columns) | What is currently being worked on? |
| Operations | Release summary / activity timeline | What has been delivered? |
| Ticket (Execution) | Detail card | What is the state of this ticket? |

The visualization is **read-only presentation of the model** — it never becomes new state (P6). Every shape renders deterministically and non-TTY-safe: plain text plus UTF-8 box drawing, no ANSI escapes.

**Engineering Domain context header.** Every projection carries a `Domain: <domain>` row in its context header alongside `Knowledge` — the Engineering Domain the projection reads (Core v1.1 §8.1). Domain projections are named after the domain they select; `ticket` is an Execution-domain projection. The projections never write (P6); the Engineering Domain itself is derived from the token family and is never part of a projection's state.

**CLI-level aliases.** `eka view sprint` and `eka view wave` are accepted and resolve to the **`execution` projection** with **identical output** — they are CLI-surface aliases only, never projections of their own. This follows the Representation Alias Registry philosophy ([`standard/representation-alias-registry-v1.1.md`](../standard/representation-alias-registry-v1.1.md)): methodology terms are convention-layer names mapped onto canonical model objects; the CLI accepts them for familiarity while the projection model stays canonical and closed. Unlike artifact aliases, CLI aliases need **no registry row** — they are command-line conveniences, not knowledge-model vocabulary — but the registry's governance note documents the relationship (see Section 3 there).

**Membership derivation rule** (relationships and the token mapping only — never file text):

- **Execution membership** — a work item belongs to an execution container iff a ticket (`tkt-`) of that container derives from it. A ticket belongs to a container iff one of its `derives-from` references resolves to the container's identity line. Container `## Work Items` tables are **not** parsed; the ticket is never parsed beyond its frontmatter relationships.
- **Domain projection membership** — `discovery`, `architecture`, `planning`, `execution`, and `operations` select every artifact whose token family is homed in that Engineering Domain (the token → domain mapping, Core v1.1 §8.1). Selection is relationship-only — no markdown parsing, no content inspection; grouping inside a projection uses the domain's own State (Content State, Planning State, Execution State).

- No active container → the execution projection is a valid empty projection (`No active container.`), exit `0`; a domain projection on a repository with no artifacts in that domain projects empty, exit `0`.
- Several active containers (invalid state) → the lexicographically smallest canonical identity is shown with a warning line.
- Ticket target forms: `eka view ticket <id>`, `eka view ticket tkt-<id>`, `eka view ticket tkt:<id>` (the prefix is stripped; `tkt-` and `tkt:` are equivalent).

### Projection architecture

One pipeline, five stages, one dependency direction — the **Projection Renderer** is the presentation stage of the projection pipeline:

```
Repository
    ↓
Knowledge Loader       one conformance.Scan — the repository's artifacts
    ↓
Knowledge Graph        identity index, relationship resolution, membership helpers
    ↓
Projection Builder     INFORMATION — what the projection shows (artifacts, order, state groups)
    ↓
Projection Renderer    PRESENTATION — how the projection looks (board, roadmap, tree, cards, timeline)
    ↓
Terminal Output
```

Two layers, one dependency direction:

| Layer | Location | Role |
|---|---|---|
| Projection engine | `view/` (public package) | Knowledge Graph (identity index, relationship resolution, membership helpers) + independent projection builders (one per projection: `discovery`, `architecture`, `planning`, `execution`, `operations`, `ticket`) + the projection registry. Pure data in, pure data out. |
| Terminal rendering | `cmd/view.go` + `cmd/ui` | Argument validation, the conformance gate, dispatch to the per-domain projection renderer (`cmd/ui`), exit-code mapping. No projection logic. |

- **Builder and renderer are independent responsibilities.** The builder defines the projection's **information** — which artifacts, in which order, grouped into which state sets; the renderer defines its **presentation** — the visual shape (Kanban board, roadmap, dependency tree, information cards, timeline, detail card) and the framing. A renderer can attach to any builder's output without touching the builder.
- The **renderer does not know the repository layout**; the **builders do not know terminals**.
- The registry is the closed set of named projections; a future projection is added by registering an independent **builder + renderer** pair — new renderers attach without pipeline redesign. The sprint/wave aliases are resolved at the command layer (they never enter the registry).
- Determinism contract: all ordering is canonical — artifacts by canonical line identity form, state groups in the fixed value order of their domain (execution states, content-state variants, planning-state), tickets by canonical identity, references in file order. The renderer preserves that ordering; it adds none of its own.

### Visualization principles

The renderers follow six principles:

- **Information before metadata** — the artifact's substance leads; state and identity follow (cards open with content, tree nodes carry state as a tag, never as the headline).
- **Visualization before serialization** — the view's shape is decided by the renderer, never by the repository's file layout; no markdown structure leaks into the projection.
- **Hierarchy before verbosity** — nested structure (tree branches, board columns, timeline phases) over flat lists.
- **Whitespace as design** — spacing groups and separates; it is layout, not wasted width.
- **One focal point** — each projection answers one question (see the visualization table); supporting detail stays secondary.
- **State understood in seconds** — a glance answers the projection's question: the board shows what is in progress, the roadmap shows what is next, the tree shows how the system is shaped.

### Validation

An automatic **pre-render conformance check**: `conformance.Validate` runs before any projection is built. A repository with blocking violations produces **no projection** — the full report is printed and the command exits `1`. Warnings never block a projection.

### Determinism and UX

- **Context header** — object kind (`Discovery`, `Architecture`, `Planning`, `Execution`, `Operations`, `Ticket`) + identity rows (`Container` / `Ticket`, `Repository`) + `Knowledge EKA v1` + `Domain <domain>` (the projection's Engineering Domain), closed by the `↓ View` pipeline.
- **State colors** — per-domain state semantics:
  - **Execution states** — `planned` dim, `todo` info, `in-progress` progress, `in-review` warning, `done` success; `unresolved` reads as warning.
  - **Content-state variants** (architecture, planning, discovery, operations) — `draft` dim, `review` / `proposed` warning, `approved` / `accepted` success, `superseded` / `amended` dim (terminal records).
  - **Planning state** — `draft` dim, `approved` info, `immutable` success (locked plan in force).
  Icons decorate (`✓` done/approved, `→` in progress, `•` everything else); the state word carries the meaning.
- **Summary block** — every projection closes with a `Summary:` of outcome facts (container, counts, status). The execution projection closes with engineering insights instead of raw counts: **Active Work** (not yet done), **Completed Work** (done), **Review Queue** (in review), **Overall Progress** (done / total).
- **Calm tone** — no banners, no ALL-CAPS, color is never the sole carrier of meaning.
- **Non-TTY deterministic** — piped/CI output is byte-identical plain text with UTF-8 icons, no ANSI escapes; color auto-disables on non-TTY, `NO_COLOR` or `TERM=dumb`.

### Examples

Illustrative non-TTY output sketches — one per projection (identities and counts vary per repository; the layout is fixed). All shapes are deterministic plain text with UTF-8 box drawing; color is never part of the layout.

`eka view discovery`:

```
Discovery
Repository  .
Knowledge   EKA v1
Domain      Discovery
↓ View

┌───────────────────────────────────────────────────┐
│ The canonical executable form of the Conformance  │
│ Rules.                                            │
│ vis:eka-cli · approved                            │
└───────────────────────────────────────────────────┘
┌───────────────────────────────────────────────────┐
│ Identities are permanent, canonical, and never    │
│ location-bound.                                   │
│ req:identity · approved                           │
└───────────────────────────────────────────────────┘
┌───────────────────────────────────────────────────┐
│ Exchange must be lossless and round-trip-safe.    │
│ req:exchange · draft                              │
└───────────────────────────────────────────────────┘
Summary:
└── Artifacts: 3
└── Approved: 2
└── Draft: 1
```

`eka view architecture`:

```
Architecture
Repository  .
Knowledge   EKA v1
Domain      Architecture
↓ View

arc:eka-cli (approved)
├── adr:identity-serialization (accepted)
│   └── req:identity (approved)
├── adr:state-vector-encoding (accepted)
├── dec:projection-model (accepted)
├── spec:exchange-format (in review)
├── std:terminology (approved)
└── gls:eka-terms (approved)
Summary:
└── Artifacts: 7
└── Accepted: 3
└── In review: 1
└── Approved: 3
```

`eka view planning`:

```
Planning
Repository  .
Knowledge   EKA v1
Domain      Planning
↓ View

Roadmap
release-1                 in force (plan:release-1 approved)
  • scp:onboarding        draft
  • plan:release-1        approved
  • plan:release-1-v2     immutable
release-2                 planned
  • epc:login-ux          draft
Summary:
└── Artifacts: 4
└── Draft: 2
└── Approved: 1
└── Immutable: 1
```

`eka view execution` (the aliases `eka view sprint` and `eka view wave` produce identical output):

```
Execution
Container   eka-cli/ctr:sprint-12
Repository  .
Knowledge   EKA v1
Domain      Execution
↓ View

┌─────────────┬─────────────┬─────────────┬─────────────┬─────────────┐
│ Planned     │ Todo        │ In Progress │ In Review   │ Done        │
├─────────────┼─────────────┼─────────────┼─────────────┼─────────────┤
│ • sto-alpha │ • sto-gamma │ → sto-delta │ • sto-      │ ✓ sto-zeta  │
│ • sto-beta  │             │             │   epsilon   │ ✓ sto-eta   │
│             │             │             │             │ ✓ sto-theta │
└─────────────┴─────────────┴─────────────┴─────────────┴─────────────┘

Summary:
└── Container: eka-cli/ctr:sprint-12
└── Active work: 5 (2 planned, 1 todo, 1 in-progress, 1 in-review)
└── Completed work: 3
└── Review queue: 1
└── Overall progress: 3 of 8 work items done
└── Status: active
```

`eka view operations`:

```
Operations
Repository  .
Knowledge   EKA v1
Domain      Operations
↓ View

Release summary
  rel:v090 (approved) — v0.9.0, publishing core from ctr:wave-6
Activity
  run:deploy-feather (approved)
  run:backup-feather (approved)
Summary:
└── Releases: 1
└── Runbooks: 2
```

`eka view ticket sto-alpha`:

```
Ticket
Ticket      eka-cli/tkt:tkt-sto-alpha
Repository  .
Knowledge   EKA v1
Domain      Execution
↓ View

┌─ tkt-sto-alpha ───────────────────── planned ─┐
│ Work item   sto-alpha (planned)               │
│ Container   ctr:sprint-12                     │
│ Derives     sto-alpha, ctr:sprint-12          │
└───────────────────────────────────────────────┘
Summary:
└── Ticket: tkt-sto-alpha
└── Work item: eka-cli/sto-alpha (planned)
└── Container: eka-cli/ctr:sprint-12
└── Status: planned
```

## `eka validate` — Conformance Validator

```
eka validate [path]
```

- `path` — optional; repository root to validate. Default: **current directory** (`.`).

### Example output (the EKA repository itself)

```
Repository
Path       .
Knowledge   EKA v1.1
↓ Validate

Repository validation
Root: . — 131 .md files, 45 artifacts, 0 errors, 8 warnings

Results (sorted by file, then rule):
  [warning] R10 reference/decisions/adr-001-identity-serialization.md: no resolvable derives-from/depends-on chain reaches a stratum above Architecture (stratum 2); stratification traceability is missing (chains must reach one of: Discovery)
  [warning] R10 reference/decisions/adr-002-state-vector-encoding.md: no resolvable derives-from/depends-on chain reaches a stratum above Architecture (stratum 2); stratification traceability is missing (chains must reach one of: Discovery)
  [warning] R10 reference/decisions/adr-003-projection-model.md: no resolvable derives-from/depends-on chain reaches a stratum above Architecture (stratum 2); stratification traceability is missing (chains must reach one of: Discovery)
  [warning] R10 reference/decisions/adr-004-phase-as-metadata.md: no resolvable derives-from/depends-on chain reaches a stratum above Architecture (stratum 2); stratification traceability is missing (chains must reach one of: Discovery)
  [warning] R10 reference/decisions/adr-005-dimension-layout.md: no resolvable derives-from/depends-on chain reaches a stratum above Architecture (stratum 2); stratification traceability is missing (chains must reach one of: Discovery)
  [warning] R10 reference/decisions/adr-006-exchange-conventions.md: no resolvable derives-from/depends-on chain reaches a stratum above Architecture (stratum 2); stratification traceability is missing (chains must reach one of: Discovery)
  [warning] R10 reference/decisions/adr-007-extension-research-finding.md: no resolvable derives-from/depends-on chain reaches a stratum above Architecture (stratum 2); stratification traceability is missing (chains must reach one of: Discovery)
  [warning] R10 reference/decisions/adr-008-engineering-domain-model.md: no resolvable derives-from/depends-on chain reaches a stratum above Architecture (stratum 2); stratification traceability is missing (chains must reach one of: Discovery)

Verdict: PASS
Summary:
└── Artifacts: 45
└── Errors: 0
└── Warnings: 8
└── Status: Repository conforms to EKA v1.1
```

> Note: the `.md files` count is a snapshot — it grows with every new convention document added; the output format stays fixed. The artifact, error, and warning counts are the contract (45 artifacts; error > 0 ⇒ FAIL). The 8 R10 warnings are the repository's own stratification traceability gap: the 8 Implementation ADRs (Architecture, stratum 2) carry no `derives-from`/`depends-on` chain reaching a higher stratum. The Reference Project (`reference/project/`) demonstrates the best practice: **37 artifacts, 0 errors, 0 warnings** (every artifact carries a stratification chain). R10 is a warning — it never blocks the verdict; the exit code stays `0`.

Output structure:

1. **Context header** — the object being validated (`Repository`), its identity rows (`Path`, `Knowledge`), and the `↓ Validate` pipeline.
2. **Repository validation** — scanned root, number of `.md` files, number of artifacts, number of errors, number of warnings.
3. **Results** — each violation on one line `[severity] rule file: message`; deterministically sorted by file, then rule (R0, R1–R12), then severity, then message. If no violations, `(no violations found)` is printed.
4. **Verdict + summary** — `Verdict: PASS` if no blocking errors, `FAIL` otherwise; the summary block closes with Artifacts, Errors, Warnings and the conformance Status.

### Example output with violations

```
  [error] R4 docs/decisions/adr-002-state-vector-encoding.md: missing owned state field existence-state on type "adr"
  [warning] R5 docs/decisions/adr-003-projection-model.md: unresolved reference "sto-x" in `depends-on` (allowed while content-state is draft)
  [warning] R10 docs/architecture/arc-001-system.md: no resolvable derives-from/depends-on chain reaches a stratum above Architecture (stratum 2); stratification traceability is missing (chains must reach one of: Discovery)
  [error] R11 docs/requirements/req-001-auth.md: declared domain "Architecture" does not match the home domain "Discovery" of type "req"
  [error] R12 docs/planning/plan-release-1-v1.md: supersedes targets eka-cli/req:login-email:1 (Discovery, stratum 1), which is strictly higher than Planning (stratum 3); cross-stratum supersession is prohibited
```

### Validation process

The `eka validate [path]` flow:

1. **Recursive scan** — the entire directory tree under `path` is traversed.
2. **Classification** — each `.md` file's frontmatter is inspected:
   - Frontmatter carries **`type` AND `id`** → **Artifact**, evaluated against R1–R12.
   - Otherwise (README, `protocol.md`, `validation.md`, `transfer.md`, canonical standard texts, etc.) → **Convention Document**, counted but skipped.
   - Frontmatter carries **exactly one** of `type`/`id` → malformed, reported as a structural error (R0).
3. **Execution of the 13 Conformance Rules (R0–R12)** — exactly as defined in [`skeleton/docs/exchange/validation.md`](../skeleton/docs/exchange/validation.md):
   - R1 Identity uniqueness, R2 filename consistency, R3 state value validity, R4 owned-set compliance, R5 referential integrity, R6 dimension == folder, R7 change-log consistency, R8 single-writer & projection, R9 well-formedness.
   - R10 stratification traceability (warning), R11 domain coherence (blocking), R12 cross-stratum supersession prohibition (blocking) — the domain-aware rules of Core v1.1 §8.1.
   - Pre-rule structural errors (invalid frontmatter, artifact rule violation, missing/corrupt identity fields, unknown type token) are grouped as **R0**.
   - Mechanical interpretation of each rule: [`conformance-notes.md`](conformance-notes.md).
4. **Deterministic reporting** — results sorted (file → rule → severity → message), so output is stable across machines and runs.

### Scan scope

- Only files ending in `.md` are inspected; other files are ignored.
- Directories named `testdata` and directories starting with a dot (e.g., `.git`) are **not descended into** — test fixtures and VCS metadata are not knowledge base content.
- Symlinks are not followed.
- Unreadable `.md` file → `eka validate` fails with an error (exit `2`): a scan that cannot see all files cannot assert conformance.

## Shell completion

`eka completion [bash|zsh|fish|powershell]` prints the completion script for the selected shell (provided by the Cobra framework). Use, for example:

```sh
source <(eka completion bash)
```

## Repository conformance

The EKA repository **passes its own conformance suite**: all `.md` files scanned, 45 artifacts (37 in the Reference Project `reference/project/` + 8 Implementation ADRs), **0 errors, 8 warnings (R10 stratification traceability), exit 0** (see example output above). R10 warnings never block the verdict.

This milestone is codified as the automated test `TestReferenceImplementationConforms` in `conformance/self_validation_test.go`: the test locates the repository root, runs `Validate` over the whole repository, and asserts 0 blocking errors. Any conformance regression (e.g., a new ADR violating a rule) therefore breaks the test suite before it can reach a commit.

## CLI architecture

The CLI is organized as **two layers + one entry point**:

| Layer | Location | Role |
|---|---|---|
| Command layer | `cmd/` (package `cmd` + `cmd/ui`) | **Only** Cobra command definitions and presentation rendering: registration, flags, help, argument validation, dispatch to services, `cmd/ui` output. No domain logic. |
| Application layer | `bootstrap/` (public package) | Repository Bootstrapper: discovery, planning, wizard, generation, validation — reusable without the CLI. |
| Application layer | `exchange/` (public package) | Import/export engine (Exchange Spec + RSF): discovery, loading, model building, serialization, deserialization, identity/relationship resolver, conflict analyzer, integration engine (staged commit + rollback), package writer — reusable without the CLI. |
| Application layer | `conformance/` (public package) | Validation engine: scanning, artifact classification, rules R0–R12, result model (`Report`); also provides `Scan` and `ParseReference` for other consumers. |
| Application layer | `view/` (public package) | Knowledge Projection Engine: Knowledge Graph (identity index, relationship resolution, membership helpers) + independent projection builders + projection registry — reusable without the CLI. |
| Entry point | `cmd/eka/main.go` | Thin: `os.Exit(cmd.Execute(...))`. Executable name: `eka`. |

```
cmd/                package cmd — Cobra command definitions (command layer)
  root.go           root command + Execute(args, stdin, stdout, stderr) int
  validate.go       validate command
  init.go           init command
  export.go         export command
  import.go         import command
  view.go           view command (renders projections via cmd/ui)
  execute_test.go   CLI tests (exit codes, help, completion, modes)
cmd/ui/             package ui — presentation primitives (no business logic)
  ui.go             Style: color/TTY/verbose context per execution
  header.go         context header (object kind, identity rows, pipeline)
  tree.go           progressive workflow tree (TTY redraw / plain lines)
  spinner.go        contextual loading (message + Braille frame, TTY only)
  summary.go        closing summary block + verbose bullet lists
  color.go          soft 256-color palette (the only colors the CLI may emit)
  icon.go           icon set (✓ • → ↓, tree connectors, spinner frames)
  step.go           deterministic [i/n] step prefix
cmd/eka/
  main.go           thin: os.Exit(cmd.Execute(...))
bootstrap/          public package — eka init engine (application layer)
exchange/           public package — import/export engine (application layer)
conformance/        public package — validation engine (application layer)
view/               public package — knowledge projection engine (application layer)
skeletonembed.go    root package — //go:embed skeleton (canonical Reference Skeleton)
```

Principles:

- **Cobra is an adapter, not the architecture.** The framework (currently Cobra) is an implementation detail of the command interface. Business logic lives in `bootstrap/`, `exchange/` and `conformance/` — public packages imported as `github.com/maleolabs/engineering-knowledge-architecture/bootstrap` and `.../conformance`, **with no dependency on `cmd/`**.
- **The command layer calls services, not the other way around.** Future tooling (import/export, graph query, SDKs, Knowledge OS integration) can import the application packages without being tied to Cobra.
- **No `internal/` or `pkg/`** — there is no second internal consumer; `bootstrap/`, `exchange/` and `conformance/` are already public API. Adding a directory without an immediate purpose is speculative abstraction (forbidden).
- **Reference Skeleton embedded** (`skeletonembed.go`): `eka init` generates repositories from the canonical `skeleton/` source, not from a hardcoded directory. The standalone binary can still generate the structure without a repository checkout.
- **Deterministic exit codes** (0/1/2) mapped in `cmd/root.go`; all errors go through a single output path `eka: <message>`.

## Contribution guide: adding a command

A new command is added without architectural refactoring:

1. **Service first** — implement business logic in an application package (`bootstrap/`, `conformance/`, `exchange/`, or a new public package), complete with tests. The CLI must not contain domain logic.
2. **Define the command** — new file `cmd/<name>.go`, package `cmd`: `Use` (verb, Naming §7.1), `Short` (one line), `Long` (details), `Example`, flags (with descriptions), argument validation via Cobra (e.g., `MaximumNArgs`), `RunE` calls the service then renders output.
3. **Register on the root** — add to `rootCmd.AddCommand(...)` in `cmd/root.go`.
4. **Exit codes** — success → `nil` (0); blocking violation → `*exitError{code: 1}` (or an equivalent sentinel); usage/internal error → plain error (mapped to 2).
5. **Test** — add cases in `cmd/execute_test.go` (exit codes, help, determinism) + service tests in the application package.
6. **Document** — update this document (command tables, examples) and the traceability matrix.
7. **Naming follows Naming and Terminology Specification v1.1 §7** — subcommands are verbs (`validate`, `init`, `diagnose`, `import`, `export`, `sync`, `format`, `graph`); do not introduce new patterns.

## CLI roadmap

| Command | Status | Notes |
|---|---|---|
| `eka init` | **Implemented** | Repository Bootstrapper (5 stages, adaptive wizard, dry-run, idempotent, post-generation validation). |
| `eka export` | **Implemented** | Exchange Package export (RSF v1.1): repo/line/instance/collection scope, automatic validation, deterministic, external reference declaration, attachments, SHA-256 digests. |
| `eka import` | **Implemented** | Exchange Package import (RSF v1.1 + Exchange §11): package + integrity validation, identity/relationship resolution, conflict → abort, atomic staged commit, rollback, post-import revalidation. |
| `eka view` | **Implemented** | Knowledge projections (execution / planning / architecture / discovery / operations / ticket; CLI aliases `sprint`, `wave` → execution): read-only views derived from the Engineering Knowledge Model — relationships + State, never markdown text — rendered as per-domain visualizations (Kanban board, roadmap, dependency tree, information cards, release timeline, detail card). Conformance-gated, deterministic, exit codes 0/1/2. |
| `eka validate` | **Implemented** | Full conformance validator (R0–R12: R1–R9 + structural R0 + domain-aware R10–R12). |
| `eka version` | **Implemented** | CLI build version + EKA standard version (currently `EKA standard 1.1`). |
| `eka completion` | **Implemented** | bash/zsh/fish/powershell completion scripts (provided by Cobra). |
| `eka diagnose` | Not implemented | Repository diagnostics — future candidate. |
| `eka graph` | Not implemented | Query/knowledge graph over artifacts. |

Future commands are added following the [Contribution guide](#contribution-guide-adding-a-command) — without architectural refactoring.
