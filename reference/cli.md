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
- **`eka view`** — the Knowledge Projection Engine: read-only projections of the Engineering Knowledge Model (the five domain projections `discovery` / `architecture` / `planning` / `execution` / `operations` + the `ticket` projection), rendered as per-domain visualizations — Kanban board (execution), roadmap (planning), dependency tree (architecture), information cards (discovery), release timeline (operations), detail card (ticket) — the canonical executable form of the State Projection semantics (Core Specification §11), relationship-derived, never markdown-rendered. Reads **Canonical Knowledge Objects from the EKA workspace canonical store via the runtime services** (`Knowledge.UnitsByProject` → `exchange.DecodeUnit`, the Runtime API — ADR-014) and projects over `exchange.Unit` — the repository must be registered and synced first (`eka sync`); the projection covers the complete knowledge of the project, the union of its registered repositories (ADR-012, ADR-013, ADR-014).
- **`eka watch`** — the realtime projection viewer: the same Knowledge Projections as `eka view`, refreshed in place by polling the canonical store — TTY-only, read-only, live refusal frames (repository not registered in the workspace), Ctrl-C to stop.
- **`eka get`** — the Machine Retrieval Interface: retrieves Engineering Knowledge as **canonical JSON generated directly from Canonical Knowledge Objects** (`exchange.Unit`) via the Runtime API (`Knowledge.Search`, `Resolver.Resolve`, `Workspace.FindRepo`) — schema **`eka-cko-v1`**, deterministic, stable across minor releases. It **never renders for readability, never parses Markdown, never queries SQLite, never reuses projection renderers**: where `eka view` derives presentational meaning, `eka get` preserves complete Engineering Knowledge semantics with no presentational transformation (ADR-015). stdout carries **only** the JSON document; errors go to stderr; exit codes 0 (success) / 1 (workspace or repository-state refusal, mirroring `eka view`) / 2 (usage, unknown identity, internal). Read-only end to end: never creates a workspace, never writes the store. The machine path (`cmd/get.go` → `machine/` → Runtime API) and the projection path (`cmd/view.go`, `cmd/watch.go` → `view/` → Runtime API) share only the Runtime API and the CKO model (ADR-015 §Decision 5).
- **`eka sync`** — the Knowledge Runtime synchronization command: pull (verify the repository's Knowledge Snapshot and seed the EKA Workspace canonical store, or seed from the docs tree via the conformance gate when no snapshot exists) then push (assemble the repository's canonical objects into a deterministic snapshot at `exchange/snapshots/`). Idempotent; deletions never applied.
- **`eka project`** — the project/repository registry of the EKA Workspace: `register` (bind a repository to a project; same `--name` = same project) and `list`.
- **`eka status`** — the workspace status probe: path, schema version, canonical store totals (Objects = references, Payloads = immutable objects, Attachments), per-repository last sync. Read-only; never creates the workspace.
- **`eka integrity`** — the store integrity verifier: `eka integrity check` recomputes every content-derived hash, strict-decodes every payload, verifies every reference (target existence + derived index columns), recomputes attachment digests, and checks the repository registry — the canonical store's immutability guarantee, verified independent of the storage engine. Manual database modification is detected, not prevented.

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

`init`, `export` and `import` render a progressive tree; `validate` renders the report as the body; `view` renders the projection as the body — each projection is a per-domain visualization (board, roadmap, tree, cards, timeline, detail card); `watch` renders the same per-domain visualizations as `view`, refreshed in place. The runtime commands (`sync`, `project`, `status`, `integrity`) render a context header + report + summary — single-operation reports, no progressive tree. Every command ends with a summary block; `watch` is the interactive exception — it runs until Ctrl-C, then clears the screen and exits `0`. `eka get` is the machine exception — it renders **nothing**: no context header, no workflow body, no summary; stdout carries only the JSON document (see [`eka get` — Machine Retrieval Interface](#eka-get--machine-retrieval-interface)).

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

One interaction model across `init` (5-stage tree), `export` (6-stage tree), `import` (7-stage tree), `validate` (header + report + summary) and `view` (header + projection + summary). Future commands must follow the same model (see [Contribution guide](#contribution-guide-adding-a-command)). `eka get` is the deliberate exception — a machine command, not a rendering command: no header, no tree, no summary, JSON only (see its section).

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
eka watch <projection> [target] [--interval N]
eka get <target>
eka validate [path]
eka sync [path]
eka sync pull [path] [--from-docs]
eka sync push [path]
eka project register [path] [--name NAME]
eka project list
eka status
eka integrity check
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
| `1` | **Blocking violation** — validation found errors, or the repository produced by `eka init` failed validation; also **workspace/repository-state refusal** for the runtime commands (`eka view` / `eka get`: unregistered repository, missing workspace) | At least one error (severity `error`); `eka view` on a repository not registered in the EKA workspace |
| `2` | **Usage/internal error** — the command did not run | Unknown command, unknown flag, too many arguments, invalid path |

Warning semantics: **warnings never affect the exit code**. A repository with warnings still exits `0`. The `1` refusal class of `eka view` / `eka get` is documented per command in their sections (deterministic refusal message + hint; no output produced).

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

## `eka sync` — Knowledge Runtime Synchronization

```
eka sync [path]
eka sync pull [path] [--from-docs]
eka sync push [path]
```

`eka sync` operates the **EKA Knowledge Runtime** (milestone v0.2.0, ADR-009/ADR-010): the bidirectional transport between a registered repository and the EKA Workspace canonical store. The workspace (`~/.eka/` or `$EKA_HOME`) is canonical; the transport is the snapshot directory `<repo>/exchange/snapshots` — an RSF package in directory layout, verified byte-exact on every read and written deterministically.

### Synchronization philosophy

- **The workspace is the store; the repository is the transport.** Immutable knowledge payloads, their mutable references, and attachments live in the workspace database; the repository carries a deterministic snapshot of its own knowledge. Knowledge objects are content-addressed (`object_hash` = the RSF per-unit digest), so snapshot digests and store hashes agree by construction.
- **`eka sync [path]` runs the full cycle: pull, then push.** `eka sync pull` and `eka sync push` run one side.
- **Auto-registration:** an unregistered repository is registered on first sync (project name = repository basename). Explicit registration with `--name` exists for multi-repository projects (`eka project register`).
- **Git-native:** the snapshot is ordinary repository content — commit, review, push, merge through normal Git workflows. No hooks, no wrappers (ADR-010 §Decision 4); synchronization is explicit by design.
- **Deterministic:** identical state → identical pull result, byte-identical snapshot, identical report.

### Pull semantics

| Mode | When | Behavior |
|---|---|---|
| **Snapshot** (default) | a snapshot exists at `<repo>/exchange/snapshots` and `--from-docs` is not set | the package is verified byte-exact (structure, strict JSON, SHA-256 integrity, self-consistency — the same verification path as `eka import`), then its units are stored as **immutable payloads** (the canonical `unit.json` entry inserted verbatim, keyed by its content hash — which equals the snapshot's per-unit digest), its references are upserted, and its attachments are stored id-keyed with digest verification, all attributed with `source_repo` = the repository name. An unchanged snapshot digest reports `unchanged` and skips the work (idempotent). A corrupt snapshot is always an error — never silently skipped. |
| **Docs** (migration) | no snapshot exists yet | the repository's `docs/` tree is validated against the conformance rules (R0–R12); blocking violations **refuse the pull** (exit `1`, full report printed). The package is then assembled exactly as `eka export` would assemble it and seeded the same way — the migration path for existing repositories. |
| **Docs** (`--from-docs`) | forced re-seed | the same docs-mode path, independent of snapshot presence — the reconciliation tool when the docs tree and the snapshot drift. |

**Deletions are never applied in v0.2:** units missing from a new pull stay in the canonical store (additive transport; future deletion protocol reserved).

### Push semantics

The repository's references in the canonical store (`source_repo` = this repository's name) are resolved to their immutable payloads and assembled into a deterministic RSF package, written **atomically**: entries are staged in `<repo>/exchange/.snapshots-tmp`, then the old snapshot directory is removed and the staging directory renamed into place — a failed push leaves the previous snapshot untouched. A repository with no stored references is a **no-op** (nothing written). Namespace resolution for the package label: the existing snapshot's header namespace when a snapshot exists, else the most common namespace among the objects (ties → lexicographically smallest), else an error.

Every pull and push run is recorded in the `sync_log` table (project, repo, direction, snapshot digest, unit count, timestamp) — the backing data of `eka status` and of the idempotent-pull check.

### Exit codes

| Code | Meaning |
|---|---|
| `0` | Sync succeeded — pulled (seeded or unchanged), pushed, or no-op |
| `1` | **Validation failure** (docs gate refused the repository — full report printed) or **integrity failure** (snapshot package refused as corrupt) |
| `2` | Usage or internal error — workspace resolution, registry failure, filesystem failure |

### Determinism

The report is deterministic and machine-readable: Runtime context header (Workspace, Project, Repository) + summary (Repository, Project, Status, Pull, Push, Snapshot) + optional warning bullets. Status values: `registered (new)` / `unchanged` / `synced`. Pull detail: `not run` / `unchanged (snapshot up to date)` / `<source>: <n> units, <m> attachments` where `<source>` is `snapshot` or `docs`. Push detail: `no-op (no stored objects)` or `<n> units, <m> attachments`. Snapshot detail: package label + 12-hex digest (e.g. `rsf-repo-feather-1.1 (3f9a2c1d0b4e)`), or `none`. Warnings render as `• no changes: snapshot already up to date` on unchanged runs.

### Examples

```sh
eka sync                 # full cycle on the current repository (pull + push)
eka sync /path/to/repo   # full cycle on another repository
eka sync pull            # pull only (snapshot mode)
eka sync pull --from-docs # force re-seed from the docs tree
eka sync push            # push only (store → snapshot)
```

Illustrative first-sync output (repository with no snapshot yet — docs-mode migration):

```
Runtime
Workspace  /home/user/.eka
Project    myproj
Repository myproj
↓ Sync

Summary:
└── Repository: myproj
└── Project: myproj
└── Status: registered (new)
└── Pull: docs: 37 units, 0 attachments
└── Push: 37 units, 0 attachments
└── Snapshot: rsf-repo-myproj-1.1 (7e9b3c1d2a4f)
```

A second run reports `Status: unchanged`, `Pull: unchanged (snapshot up to date)`, and the `• no changes: snapshot already up to date` warning.

### Error handling

| Condition | Behavior |
|---|---|
| Repository non-conformant (docs gate, no snapshot / `--from-docs`) | Pull refused; full validation report printed; `sync pull refused: repository validation failed with N blocking error(s); no knowledge seeded` on stderr; exit `1` |
| Snapshot corrupt (digest, JSON, structure, self-consistency) | `sync pull failed: snapshot package refused: ...`; exit `1` |
| Workspace resolution failure (`EKA_HOME` relative, no home dir) | Error; exit `2` |
| Path not accessible / filesystem failure | Error; exit `2` |

## `eka project` — Workspace Project Registry

```
eka project register [path] [--name NAME]
eka project list
```

`eka project` manages the projects and repositories registered in the EKA Workspace (default `~/.eka`, or `$EKA_HOME`). A project groups one or more repositories; the repository name is the basename of its normalized absolute path (the unit key of the canonical store); the project name is the `--name` flag value or the same basename.

### `eka project register [path] [--name NAME]`

Registers the EKA repository at `path` (default: the current directory). Registering the same repository again is a **no-op** (the stored path is refreshed; the report says `already registered`). A repository registered under one project always resolves to that project — the first registration owns it. Multi-repository projects use the same `--name`:

```sh
eka project register                       # current directory, project = basename
eka project register /path/to/repo         # project = "repo"
eka project register /path/to/repo --name atrium   # project "atrium", repo "repo"
```

| Code | Meaning |
|---|---|
| `0` | Registration succeeded (new or already registered) |
| `2` | Usage or internal error — workspace resolution, registry failure, inaccessible path |

Output: `Project` context header (Project, Repository, Path) + summary (Project, Repository, Path, Status: `registered` / `already registered`).

### `eka project list`

Lists every project with its repositories (name and path), sorted deterministically — projects by id, repositories by name. A workspace with no registered projects prints `No projects registered yet. Run 'eka project register' to add one.` and exits `0`.

| Code | Meaning |
|---|---|
| `0` | Success (list printed, or informational empty message) |
| `2` | Usage or internal error |

## `eka status` — Workspace Status

```
eka status
```

`eka status` prints the EKA Workspace overview: workspace path, schema version, workspace id, created date, canonical store totals (**Objects** = the number of references, **Payloads** = the number of immutable objects including retained history, **Attachments**), registered projects and repositories, and the most recent sync per repository (`[direction digest at time]` from the sync log). There is no Relationships count: relationships live inside the immutable payloads and are parsed at read time.

`eka status` is a **read-only probe**: without a workspace it prints `No EKA workspace at <home> yet. Run 'eka project register' to create it.` and exits `0` — it never initializes the workspace, never creates files.

| Code | Meaning |
|---|---|
| `0` | Status shown (or no workspace yet) |
| `2` | Usage or internal error |

## `eka integrity check` — Canonical Store Integrity

```
eka integrity check
```

`eka integrity check` verifies the integrity of the EKA Workspace canonical store (default `~/.eka`, or `$EKA_HOME`) — the executable form of the Immutable Engineering Knowledge Model's verification contract (ADR-011 §Decision 4). Engineering Knowledge Objects are immutable and content-addressed: every payload row is keyed by `SHA-256(unit.json ‖ content)`, written once, and never updated. The command recomputes every content-derived value and compares it with the stored state, independent of the storage engine.

### Verification levels

The scan is **read-only** (all SQL parameterized; no writes) and checks:

| Level | Check |
|---|---|
| **1. Payload integrity** | recompute `SHA-256(unit.json ‖ content)` for every `object_payloads` row and compare with `object_hash` |
| **2. Payload decode** | every `unit_json` payload strict-decodes (unknown fields rejected — the same decode path as the exchange layer, RSF §9.5) |
| **3. Reference integrity** | every `object_refs.object_hash` exists in `object_payloads`; the refs' derived index columns and `form` equal the referenced payload's identity fields (namespace, type, id, instance version, revision, dimension, domain, phase) |
| **4. Workspace integrity** | registry foreign keys (`repos` → `projects`); attachment digests recomputed against `attachments.data` |

Violation kinds in the report: `payload-hash`, `payload-decode`, `reference-target`, `reference-index`, `attachment-hash`, `registry` — sorted by (kind, subject) for deterministic output.

### Detected, not prevented

Manual modification of the SQLite file is **DETECTED, not prevented** — the runtime verifies and reports inconsistencies; it does not pretend a hand-edited database is trusted. SQLite is a persistence layer, not a trust boundary. Unreferenced payloads are the **immutable history archive**: they are counted in the report (`History payloads`) and **never** reported as violations.

### Exit codes

| Code | Meaning |
|---|---|
| `0` | clean — no violations (history payloads may be present) |
| `1` | violations found — the report is printed; stderr carries `eka: integrity check found N violation(s)` |
| `2` | Usage or internal error — workspace resolution or store read failure |

### Output

Deterministic, non-TTY-safe: `Runtime` context header (Workspace, Schema) + summary (Payloads checked, References checked, Attachments checked, History payloads, Violations) + optional violation lines. Example (clean store):

```
Runtime
Workspace  /home/user/.eka
Schema     v2

Summary:
└── Payloads checked: 37
└── References checked: 37
└── Attachments checked: 0
└── History payloads: 0
└── Violations: 0
```

### Error handling

| Condition | Behavior |
|---|---|
| Workspace resolution failure (`EKA_HOME` relative, no home dir) | `eka: <error>`, exit `2` |
| Store read failure during the scan | `eka: integrity check failed: <error>`, exit `2` |

## `eka view` — Knowledge Projections

```
eka view [projection] [target]
```

`eka view` projects the **Engineering Knowledge Model** of the project owning the repository rooted at the current directory: read-only views over the **complete Engineering Knowledge of the project** — the union of every registered repository's units — and their relationships. Projections are named by **Engineering Domain** — `discovery`, `architecture`, `planning`, `execution`, `operations` — plus the `ticket` and `board` projections. The `target` argument is required by the ticket projection only (a bare ticket id, `tkt-<id>` or `tkt:<id>`); domain projections ignore it. With no arguments, the available projections are listed and the command exits `0`.

The projection is built from **Canonical Knowledge Objects (CKO) read from the EKA workspace canonical store via the runtime services**: `Knowledge.UnitsByProject` (the Runtime API — ADR-014) resolves every reference of the project to its immutable payload (`object_payloads`), and `exchange.DecodeUnit` strict-decodes each into an `exchange.Unit`, ordered by canonical form (ADR-013; ADR-012; ADR-014; the CKO schema is [`cko-specification.md`](cko-specification.md)). The repository must be **registered and synced** first: `eka sync` compiles the authoring tree through the Knowledge Compiler (conformance-gated) and seeds the store. At projection time **no Markdown is read and no conformance gate runs** — the projection input is exactly the synced canonical state. Authoring UX: **write Markdown → `eka sync` → `eka view`**.

### Projection philosophy — Knowledge Projection

- The viewing layer is **generated from the Engineering Knowledge Model** — artifact identity, State fields, and relationships — **not from markdown rendering**. The CLI never depends on markdown formatting; no file text is parsed for view content (ticket bodies are never read, container `## Work Items` tables are never parsed).
- Markdown remains the **editing interface**; EKA is the **canonical model**; the projection is the **view**. A projection has no State of its own and never becomes a writer (P6, Core Specification §11) — `eka view` is the canonical executable form of the State Projection semantics.
- The projection engine is **pure data in, pure data out**: it contains no terminal knowledge, no command framework, and no output.

### Snapshot interaction model

`eka view` is a single synchronous snapshot — no TUI, no navigation, no editing, no background process. One run produces one projection, then exits:

```
run → resolve → read → construct → render → exit
```

1. **Resolve** — the EKA workspace is resolved and the current directory is resolved to a registered repository (`Workspace.FindRepo`, the Runtime API — ADR-014). A repository **not registered** in the workspace is **refused** (exit `1`) with the deterministic message: `eka: view refused: repository <abs> is not registered in the EKA workspace; run 'eka sync' (auto-registers) or 'eka project register' first`. No projection is produced.
2. **Read** — the project's canonical units are read from the workspace canonical store via the runtime Knowledge service: `Knowledge.UnitsByProject` resolves every reference to its immutable payload, `exchange.DecodeUnit` strict-decodes each (canonical-form order). No Markdown is read and no conformance gate runs at projection time — authoring validation ran inside `eka sync` (the compile gate). A registered project with **no synced knowledge** renders an empty projection with the informational note `no synced knowledge for project <id>; run 'eka sync' after editing docs`, exit `0`.
3. **Construct** — one Knowledge Graph over the decoded units (the `view` package), then one projection build (the builder).
4. **Render** — the projection renderer applies the projection's visualization (board, roadmap, tree, cards, timeline, detail card), deterministically.
5. **Exit** — mapped per the exit-code contract below.

The engine is synchronous and stateless: a future loading state can wrap the whole call without restructuring.

### Exit codes

| Code | Meaning |
|---|---|
| `0` | Projection produced — including empty projections (no active container, no domain artifacts, no synced knowledge) |
| `1` | Repository not registered in the EKA workspace — no projection is produced (deterministic refusal message + hint printed) |
| `2` | Usage or internal error — unknown projection, missing or unknown ticket target, workspace/store failure |

Warnings never block a projection.

### Projections

The canonical projections — one per Engineering Domain, plus the ticket and board projections:

| Command | Engineering Domain | Content | Rendered as |
|---|---|---|---|
| `eka view discovery` | Discovery | `vis-` / `str-` / `req-` / `fnd-` artifacts — vision, strategy, requirements, research findings | Information cards |
| `eka view architecture` | Architecture | `adr-` / `dec-` / `arc-` / `spec-` / `std-` / `gls-` artifacts — decisions, architecture descriptions, specifications, standards, vocabulary; per-node **Content State** (`draft` / `review` / `approved` / post-approval terminals, incl. the ADR/decision variants) shown on the node | Dependency tree |
| `eka view planning` | Planning | `scp-` / `epc-` / `plan-` / `trc-` artifacts — scope definitions, epics, plans, traceability artifacts, sequenced by phase; per-entry **Planning State** (`draft` / `approved` / `immutable`) | Roadmap / timeline |
| `eka view execution` | Execution | The **active Execution Container**'s work items on the fixed five-column board (`planned` / `todo` / `in-progress` / `in-review` / `done`; empty columns keep their heading), each item a two-line card — name over `type · container` (same tag rule as the board — an item shared across containers shows all of them). The ticket list is **not rendered as a block** — tickets project individually via `eka view ticket`. | Kanban board |
| `eka view board` | Execution | **Every work item in the project** — across all execution containers (active and completed) and outside any container — on the fixed five-column board, each item a two-line card — name over `type · container` (`wave-7`); items not referenced by any ticket container show `unassigned` | Kanban board |
| `eka view operations` | Operations | `run-` / `rel-` artifacts — runbooks, operational guides, release records | Release summary / activity timeline |
| `eka view ticket <id>` | Execution | One ticket's detail: **projected status derived from the referenced work item's owner Execution State — never from the ticket's own text** — plus the container it derives from and its `derives-from` references | Detail card |

**Visualization by domain.** Each projection renders as a purpose-built console for its domain's primary question:

| Engineering Domain | Visualization style | Primary question it answers |
|---|---|---|
| Discovery | Information cards | What are we building? |
| Architecture | Dependency tree | How is the system structured? |
| Planning | Roadmap / timeline | What are we planning next? |
| Execution | Kanban board (five columns) | What is currently being worked on? |
| Board (Execution) | Kanban board (five columns, per-item container tags) | What is the total work in the project? |
| Operations | Release summary / activity timeline | What has been delivered? |
| Ticket (Execution) | Detail card | What is the state of this ticket? |

The visualization is **read-only presentation of the model** — it never becomes new state (P6). Every shape renders deterministically and non-TTY-safe: plain text plus UTF-8 box drawing, no ANSI escapes. The kanban boards are **width-adaptive on a TTY**: columns scale with the terminal width (capped, never below the minimum) so the grid fits the screen. Every work item renders as a **two-line card** — the item name on the first line, the `type · container` context on the second — separated by a blank gap row so cards stay visually distinct. The type renders as a **badge** (`[sto]`, `[bug]`, …) colored per type — story/story-alias blue, technical story cyan, bug/defect red, tech debt amber, chore/spike gray, unknown tokens neutral gray (never an error) — while the container tag keeps the **execution-state color** of its column, so the badge and the tag are separate colored segments. The badge and the container tag survive narrow terminals in that order: the name truncates first, then the badge is dropped, and the container tag is the last thing to go. Non-TTY output (pipes, CI, tests) uses the fixed maximum column width and stays byte-identical.

**Engineering Domain context header.** Every projection carries a `Domain: <domain>` row in its context header alongside `Knowledge` — the Engineering Domain the projection reads (Core v1.1 §8.1). Domain projections are named after the domain they select; `ticket` and `board` are Execution-domain projections. The projections never write (P6); the Engineering Domain itself is derived from the token family and is never part of a projection's state.

**CLI-level aliases.** `eka view sprint` and `eka view wave` are accepted and resolve to the **`execution` projection** with **identical output** — they are CLI-surface aliases only, never projections of their own. This follows the Representation Alias Registry philosophy ([`standard/representation-alias-registry-v1.1.md`](../standard/representation-alias-registry-v1.1.md)): methodology terms are convention-layer names mapped onto canonical model objects; the CLI accepts them for familiarity while the projection model stays canonical and closed. Unlike artifact aliases, CLI aliases need **no registry row** — they are command-line conveniences, not knowledge-model vocabulary — but the registry's governance note documents the relationship (see Section 3 there).

**Membership derivation rule** (relationships and the token mapping only — never file text):

- **Execution membership** — a work item belongs to an execution container iff a ticket (`tkt-`) of that container derives from it. A ticket belongs to a container iff one of its `derives-from` references resolves to the container's identity line. Container `## Work Items` tables are **not** parsed; the ticket is never parsed beyond its frontmatter relationships.
- **Board membership** — the board selects **every work item line** whose token family owns the Execution State domain, regardless of container membership: items of the active container, of completed containers, and items no ticket references (`unassigned`). Each item's container tag follows the execution membership rule above (the containers whose tickets reference it), so the board keeps container context per item — the board never merges containers into an anonymous aggregate.
- **Domain projection membership** — `discovery`, `architecture`, `planning`, `execution`, and `operations` select every artifact whose token family is homed in that Engineering Domain (the token → domain mapping, Core v1.1 §8.1). Selection is relationship-only — no markdown parsing, no content inspection; grouping inside a projection uses the domain's own State (Content State, Planning State, Execution State).

- No active container → the execution projection is a valid empty projection (`No active container.`), exit `0`; a domain projection on a repository with no artifacts in that domain projects empty, exit `0`; a board with no work items projects empty (`No work items.`), exit `0`.
- Several active containers (invalid state) → the lexicographically smallest canonical identity is shown with a warning line.
- Ticket target forms: `eka view ticket <id>`, `eka view ticket tkt-<id>`, `eka view ticket tkt:<id>` (the prefix is stripped; `tkt-` and `tkt:` are equivalent).

### Projection architecture

One pipeline, five stages, one dependency direction — the **Projection Renderer** is the presentation stage of the projection pipeline:

```
Workspace canonical store (project units)
    ↓
Store Read             Knowledge.UnitsByProject (Runtime API; store.UnitsByProject behind the Kernel)
                       — object_refs → object_payloads → exchange.DecodeUnit
                       (canonical-form order; seeded by eka sync)
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
| Projection engine | `view/` (public package) | Knowledge Graph over the store-decoded CKO set (identity index, relationship resolution, membership helpers) + independent projection builders (one per projection: `discovery`, `architecture`, `planning`, `execution`, `operations`, `ticket`, `board`) + the projection registry. Pure data in, pure data out. |
| Terminal rendering | `cmd/view.go` + `cmd/ui` | Argument validation, the runtime service reads (repository resolution via `Workspace.FindRepo`, project units via `Knowledge.UnitsByProject` — the Runtime API, ADR-014), dispatch to the per-domain projection renderer (`cmd/ui`), exit-code mapping. No projection logic. |

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

The conformance gate moved into `eka sync`. Projections no longer run a pre-render conformance check: authoring validation (R0–R12) runs inside `eka sync` when the store is seeded (the compiler's authoring-validation stage; blocking violations refuse the pull with the full report, exit `1`). The projection reads the synced canonical state; the correctness of the projected units is the store's integrity contract — `eka integrity check` (payload hash, decode, references, attachments, workspace). Warnings never block a projection.

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
│ Planned (2) │ Todo (1)    │ In Progress │ In Review   │ Done (3)    │
│             │             │ (1)         │ (1)         │             │
├─────────────┼─────────────┼─────────────┼─────────────┼─────────────┤
│ ▸ sto-alpha │ ▸ sto-gamma │ ▸ sto-delta │ ▸ sto-      │ ▸ sto-zeta  │
│ ▸ sto-beta  │             │             │   epsilon   │ ▸ sto-eta   │
│             │             │             │             │ ▸ sto-theta │
└─────────────┴─────────────┴─────────────┴─────────────┴─────────────┘

Summary:
└── Container: eka-cli/ctr:sprint-12
└── Active Work: 2
└── Completed Work: 3
└── Review Queue: 1
└── Overall Progress: ████░░░░░░ 3/8 (38%)
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

## `eka get` — Machine Retrieval Interface

```
eka get <target>
```

`eka get` is the **machine interface** of the EKA Runtime (ADR-015): it retrieves Engineering Knowledge as **canonical JSON generated directly from Canonical Knowledge Objects** (`exchange.Unit`) via the Runtime API — `Knowledge.Search`, `Resolver.Resolve`, `Workspace.FindRepo`. It **never renders for readability, never parses Markdown, never queries SQLite, never reuses projection renderers**. Where `eka view` derives presentational meaning (projected status, boards, dependency trees), `eka get` preserves **complete Engineering Knowledge semantics** — identity, state, relationships, change-log, content, integrity metadata — with no presentational transformation. The consumers it serves: MCP, Atrium, VS Code extensions, AI agents, automation and plain scripts — anything that wants the knowledge itself, machine-readable, deterministic, stable.

The repository must be **registered and synced** first (`eka sync`). Authoring UX: **write Markdown → `eka sync` → `eka get`**.

### Human vs machine — the separation

```
Human:   store → Runtime API → projection model (view/)  → renderer → terminal
Machine: store → Runtime API → CKO (machine/)            → canonical JSON → stdout
```

The two commands share **only** the Runtime API and the CKO model (ADR-015 §Decision 5):

```
cmd/view.go, cmd/watch.go → view/ → Runtime API          (projection path)
cmd/get.go                → machine/ → Runtime API       (machine path)
```

- The machine path imports `{runtime, exchange}`; it never imports `view/`. The projection path never imports `machine/`.
- No projection code is reachable from the machine path, no machine code from the projection path — a rendering change in `view/` can never touch the machine contract, and the JSON is independent of projection rendering.
- Projections derive presentational meaning; `eka get` performs no presentational transformation. The two paths cannot drift because both are pure functions of the same objects.

### Query model

`eka get` takes exactly **one target, no flags in v0.2**. Two target shapes, discriminated syntactically:

| Target shape | Meaning | Result |
|---|---|---|
| contains `:` | **Identity lookup** — resolved via `Resolver.Resolve`: the RSF **canonical form** (`<ns>/<type>:<id>:<v>`, the exact instance) or the **qualified line form** (`<ns>/<type>:<id>`, the lowest instance-version of the line) | one **Document** |
| no `:` | **Domain query** — one of the five Engineering Domains: `discovery` / `architecture` / `planning` / `execution` / `operations` | one **Collection** |

- **The namespace is required** on identity lookups — unqualified forms (`<type>:<id>`) are refused as ambiguous: the Runtime resolves globally, with no referrer context (the refusal message names the required forms). Tickets and containers are identity lookups on their type tokens — `eka get acme/tkt:EK-123`, `eka get acme/ctr:core` — no storage concepts appear anywhere in the query surface.
- A domain query returns the project's units whose **derived domain** equals the token (classification, else type token), selected via `Knowledge.Search` over the project union, **sorted by canonical form**:

```json
{
  "schema": "eka-cko-v1",
  "collection": "domain",
  "domain": "Planning",
  "count": 4,
  "units": [ /* Documents, sorted by canonical form */ ]
}
```

`collection` is the discriminator: a Collection carries it, a single-object Document does not. `domain` carries the canonical Engineering Domain name (e.g. `Planning`, `Execution`).

### JSON contract (schema `eka-cko-v1`)

One CKO = one **Document**. The Document serializes the unit's fields in the fixed declared order of the table — the **serialization order**, which machine consumers may rely on (ADR-015 §Decision 2; implemented by the `machine/` package):

| Field | Source (CKO) | Notes |
|---|---|---|
| `schema` | contract constant | `"eka-cko-v1"` — the stable schema string every Document carries |
| `identity` | `Unit.Identity` | the complete identity tuple: `namespace`, `type`, `id`, `instance_version` (RSF unit.json naming) |
| `canonical_form` | `Identity.CanonicalForm()` | `<ns>/<type>:<id>:<v>` — the RSF Canonical Identity Form |
| `engineering_domain` | derived | the Engineering Domain: `Classification.Domain` when present, else `conformance.DomainForToken` on the type token (ADR-008) |
| `stratum` | derived | `conformance.Stratum(domain)` — the Knowledge Stratum, 1 highest → 5 lowest (ADR-008) |
| `revision` | `Unit.Revision` | unit metadata — never identity, never an ordering key; omitted when `0` |
| `author` / `created` / `updated` | `Unit.Author` / `Created` / `Updated` | omitted when empty (`""` on the source) |
| `state_vector` | `Unit.StateVector` | the five owned domains in canonical declared order, omitted when empty (an empty vector renders `{}`) — `content-state`, `execution-state`, `planning-state`, `container-state`, `existence-state`, **same field naming as the RSF unit.json** |
| `phase` | `Unit.Phase` | phase metadata (`scp-` / `plan-` artifacts); omitted when empty |
| `classification` | `Unit.Classification` | `dimension`, `dimensions_secondary`, `domain`; omitted when the CKO carries none |
| `relationships` | `Unit.Relationships` | the stored order — the CKO's canonical (type, target) order; each entry `{type, target}` |
| `change_log` | `Unit.ChangeLog` | occurrence order; each entry `{date, domain, from, to, by}` |
| `content` | `ContentRef` + payload | `{representation, text}` — `representation` is the Representation Identifier (`eka/structured-text/1` today), `text` is the **opaque representation payload carried verbatim** |
| `object_hash` | payload digest | the content-derived digest of the immutable payload the unit was decoded from: `SHA-256(unit.json ‖ content)`, byte-identical to the RSF per-unit digest (ADR-011) |

**Content is never parsed or re-structured.** Markdown is one representation (`eka/structured-text/1`); the JSON carries it verbatim as CKO content — it does **not** serialize Markdown structure (headings, frontmatter, files). The RSF's `file` indirection is package-layout-specific and has no place in the machine document; the payload travels inline as `text`. A metadata-only consumer pays the payload size; the documented remedy is a future content-filter flag (ADR-015 §Consequences).

**Determinism and stability.** Fixed struct field order (the table above is the order), stable schema string, collections sorted by canonical form (`relationships` by the CKO's stored (type, target) order, `change_log` in occurrence order), **no timestamps, no host-dependent values**. Identical synced store state → **byte-identical JSON**. `eka-cko-v1` is **stable across minor releases** — future tooling can depend on it. Changes are **additive** (new fields appended) or **schema-versioned** (a breaking change bumps the schema string; it never mutates the contract under existing consumers).

### Output contract and exit codes

| Channel | Contract |
|---|---|
| `stdout` | **only** the JSON document (+ one trailing newline) — one Document for an identity lookup, one Collection for a domain query. No banners, no progress, no decorations — the output is machine-parseable as-is. |
| `stderr` | errors and refusals — deterministic messages with hints, a single `eka: ...` line |
| exit `0` | success — the Document or Collection was emitted |
| exit `1` | **workspace or repository-state refusal** — mirrors `eka view`'s unregistered-repository refusal class (ADR-013 §Decision 3): a missing workspace (detached `runtime.Open` — `eka get` **never creates a workspace**; the detached state is a refusal, not an initialization) and an unregistered repository (`Workspace.FindRepo` misses) both refuse with a deterministic message + hint |
| exit `2` | usage (unparseable target, unknown domain token, unqualified identity), **unknown identity**, internal error |

`eka get` is **read-only end to end**: `runtime.Open` (the read-style entry, never `runtime.Ensure`), resolve, query, serialize, emit. It never registers, never syncs, never writes the store.

### Examples

```sh
eka get feather/adr:001-identity-serialization:1   # identity lookup — canonical form (exact instance)
eka get feather/sto:publish-post                   # identity lookup — qualified line form (lowest instance)
eka get acme/tkt:EK-123                            # tickets are identity lookups on their type token
eka get acme/ctr:core                              # containers likewise
eka get architecture                               # domain query — every Architecture unit as a Collection
eka get execution                                  # domain query — every Execution unit as a Collection
```

Illustrative Document (identities, dates and content vary per repository; the shape and field order are fixed). `...` marks elided content — real output is complete:

```json
{
  "schema": "eka-cko-v1",
  "identity": {
    "namespace": "feather",
    "type": "adr",
    "id": "001-identity-serialization",
    "instance_version": 1
  },
  "canonical_form": "feather/adr:001-identity-serialization:1",
  "engineering_domain": "Architecture",
  "stratum": 2,
  "revision": 1,
  "author": "Engineering Architecture",
  "created": "2026-01-15",
  "updated": "2026-01-15",
  "state_vector": {
    "content-state": "accepted",
    "existence-state": "active"
  },
  "classification": {
    "dimension": "decisions",
    "domain": "Architecture"
  },
  "relationships": [
    { "type": "derives-from", "target": "feather/req:identity:1" },
    { "type": "depends-on", "target": "feather/adr:002-state-vector-encoding:1" }
  ],
  "change_log": [
    { "date": "2026-01-15", "domain": "content-state", "from": "proposed", "to": "accepted", "by": "Engineering Architecture" }
  ],
  "content": {
    "representation": "eka/structured-text/1",
    "text": "# ADR-001 — Identity Serialization\n\n..."
  },
  "object_hash": "3f9a2c1d0b4e..."
}
```

### Future extensions (not implemented)

Relationship traversal, knowledge graph, timeline/history, semantic/vector search, metadata filtering, context generation and content filtering all grow as `eka get` targets/flags over the **same Runtime API** (Resolver, Relations, Timeline, `Knowledge.Search` boundaries, ADR-014). None are built in this milestone, and none change the JSON contract — everything is additive; the schema stays `eka-cko-v1`-compatible (ADR-015 §Decision 6).

## `eka watch` — Realtime Projections

```
eka watch <projection> [target] [--interval N]
```

`eka watch` is the **realtime projection viewer**: it renders the same projections as `eka view` — `discovery`, `architecture`, `planning`, `execution`, `operations`, `ticket`, plus the `sprint` / `wave` CLI aliases of `execution` — and refreshes them in place while the repository changes. Like `eka view`, it is read-only: a projection has no State of its own and never becomes a writer (P6, Core Specification §11). Like `eka view`, each refresh reads the project's **Canonical Knowledge Objects from the workspace canonical store via the runtime services** (`Knowledge.UnitsByProject` per cycle, the Runtime API — ADR-014) and projects over `exchange.Unit` (ADR-013) — an `eka sync` run in another terminal is picked up on the next tick without a restart. The [projections table](#projections) above defines what each projection shows; the ticket target argument behaves exactly as in `eka view`.

### Interaction model

- **No keyboard navigation** — the viewer refreshes by polling; Ctrl-C (SIGINT) stops it. No paging, no cursor movement, no editing.
- **Polling refresh, no filesystem watchers** — the project's units are re-read from the canonical store at a fixed interval (`--interval`); there is no fsnotify, no compile, no Markdown read at watch time.
- **TTY-only** — a terminal is required. On a non-TTY (pipe, redirect, CI) the command exits `2` with the deterministic error `requires a terminal` — the non-TTY determinism contract (byte-identical, no ANSI) applies to this error path.

### Refresh model

| Aspect | Behavior |
|---|---|
| `--interval N` | refresh period in seconds — default `2`, minimum `1` |
| Redraw | only when the frame changed — identical frames are not redrawn |
| Clear screen | on open and on exit |

### Live refusal handling

While the current directory is **not registered** in the EKA workspace, the projection is replaced by an **unregistered-repository refusal frame** — the deterministic refusal message with the sync hint — instead of the projection:

- watching **keeps running** and **auto-recovers**: the first tick that finds the repository registered and synced (e.g. `eka sync` run in another terminal) renders the projection.

The old compile-failure frame is gone: authoring validation (R0–R12) runs inside `eka sync`, never at watch time — a registered repository with no synced knowledge renders the empty projection, not a failure frame.

### Frame footer

Every frame closes with the footer line `watching — Ctrl-C to stop (interval Ns)` — e.g. `watching — Ctrl-C to stop (interval 2s)`.

### Exit codes

| Code | Meaning |
|---|---|
| `0` | clean stop — Ctrl-C (SIGINT): screen cleared, exit |
| `2` | usage or internal error — unknown projection, missing/unknown ticket target, non-TTY invocation |

Refusal states are rendered as frames, never as exit codes: an unregistered repository does not stop the viewer, and `eka watch` never exits `1`.

### Determinism

`eka watch` is TTY-only by design, so the non-TTY determinism contract applies to its error path: non-TTY output is exactly one deterministic error line, exit `2`, no ANSI. On a TTY, each frame is deterministic for its snapshot — identical repository state produces an identical frame, which is why unchanged frames are not redrawn; the only time-dependent behavior is the refresh cadence itself.

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
Root: . — 142 .md files, 52 artifacts, 0 errors, 15 warnings

Results (sorted by file, then rule):
  [warning] R10 reference/decisions/adr-001-identity-serialization.md: no resolvable derives-from/depends-on chain reaches a stratum above Architecture (stratum 2); stratification traceability is missing (chains must reach one of: Discovery)
  [warning] R10 reference/decisions/adr-002-state-vector-encoding.md: no resolvable derives-from/depends-on chain reaches a stratum above Architecture (stratum 2); stratification traceability is missing (chains must reach one of: Discovery)
  [warning] R10 reference/decisions/adr-003-projection-model.md: no resolvable derives-from/depends-on chain reaches a stratum above Architecture (stratum 2); stratification traceability is missing (chains must reach one of: Discovery)
  [warning] R10 reference/decisions/adr-004-phase-as-metadata.md: no resolvable derives-from/depends-on chain reaches a stratum above Architecture (stratum 2); stratification traceability is missing (chains must reach one of: Discovery)
  [warning] R10 reference/decisions/adr-005-dimension-layout.md: no resolvable derives-from/depends-on chain reaches a stratum above Architecture (stratum 2); stratification traceability is missing (chains must reach one of: Discovery)
  [warning] R10 reference/decisions/adr-006-exchange-conventions.md: no resolvable derives-from/depends-on chain reaches a stratum above Architecture (stratum 2); stratification traceability is missing (chains must reach one of: Discovery)
  [warning] R10 reference/decisions/adr-007-extension-research-finding.md: no resolvable derives-from/depends-on chain reaches a stratum above Architecture (stratum 2); stratification traceability is missing (chains must reach one of: Discovery)
  [warning] R10 reference/decisions/adr-008-engineering-domain-model.md: no resolvable derives-from/depends-on chain reaches a stratum above Architecture (stratum 2); stratification traceability is missing (chains must reach one of: Discovery)
  [warning] R10 reference/decisions/adr-009-knowledge-runtime-architecture.md: no resolvable derives-from/depends-on chain reaches a stratum above Architecture (stratum 2); stratification traceability is missing (chains must reach one of: Discovery)
  [warning] R10 reference/decisions/adr-010-synchronization-model.md: no resolvable derives-from/depends-on chain reaches a stratum above Architecture (stratum 2); stratification traceability is missing (chains must reach one of: Discovery)
  [warning] R10 reference/decisions/adr-011-immutable-engineering-knowledge-model.md: no resolvable derives-from/depends-on chain reaches a stratum above Architecture (stratum 2); stratification traceability is missing (chains must reach one of: Discovery)
  [warning] R10 reference/decisions/adr-012-canonical-knowledge-object-runtime.md: no resolvable derives-from/depends-on chain reaches a stratum above Architecture (stratum 2); stratification traceability is missing (chains must reach one of: Discovery)
  [warning] R10 reference/decisions/adr-013-store-backed-projections.md: no resolvable derives-from/depends-on chain reaches a stratum above Architecture (stratum 2); stratification traceability is missing (chains must reach one of: Discovery)
  [warning] R10 reference/decisions/adr-014-runtime-interface-architecture.md: no resolvable derives-from/depends-on chain reaches a stratum above Architecture (stratum 2); stratification traceability is missing (chains must reach one of: Discovery)
  [warning] R10 reference/decisions/adr-015-machine-retrieval-interface.md: no resolvable derives-from/depends-on chain reaches a stratum above Architecture (stratum 2); stratification traceability is missing (chains must reach one of: Discovery)

Verdict: PASS
Summary:
└── Artifacts: 52
└── Errors: 0
└── Warnings: 15
└── Status: Repository conforms to EKA v1.1
```

> Note: the `.md files` count is a snapshot — it grows with every new convention document added; the output format stays fixed. The artifact, error, and warning counts are the contract (52 artifacts; error > 0 ⇒ FAIL). The 15 R10 warnings are the repository's own stratification traceability gap: the 15 Implementation ADRs (Architecture, stratum 2) carry no `derives-from`/`depends-on` chain reaching a higher stratum. The Reference Project (`reference/project/`) demonstrates the best practice: **37 artifacts, 0 errors, 0 warnings** (every artifact carries a stratification chain). R10 is a warning — it never blocks the verdict; the exit code stays `0`.

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

The EKA repository **passes its own conformance suite**: all `.md` files scanned, 52 artifacts (37 in the Reference Project `reference/project/` + 15 Implementation ADRs), **0 errors, 15 warnings (R10 stratification traceability), exit 0** (see example output above). R10 warnings never block the verdict.

This milestone is codified as the automated test `TestReferenceImplementationConforms` in `conformance/self_validation_test.go`: the test locates the repository root, runs `Validate` over the whole repository, and asserts 0 blocking errors. Any conformance regression (e.g., a new ADR violating a rule) therefore breaks the test suite before it can reach a commit.

## CLI architecture

The CLI is organized as **two layers + one entry point**:

| Layer | Location | Role |
|---|---|---|
| Command layer | `cmd/` (package `cmd` + `cmd/ui`) | **Only** Cobra command definitions and presentation rendering: registration, flags, help, argument validation, dispatch to services, `cmd/ui` output. No domain logic. |
| Application layer | `runtime/` (public package) | The **Runtime Kernel** (ADR-014): the internal **Runtime API** (Workspace, Knowledge, Resolver, Relations, Timeline, Snapshot, Integrity — concrete service types, no Go interfaces) + the **Authoring API** (`runtime.Authoring`: Validate/Compile/Sync); the one sanctioned entry point every consumer — the CLI, future Context Engine/MCP/Atrium — talks to. `store/`, `workspace/`, `sync/` and `compile/` are kernel-internal behind it; production `cmd/` imports only `runtime/`. |
| Application layer | `bootstrap/` (public package) | Repository Bootstrapper: discovery, planning, wizard, generation, validation — reusable without the CLI. |
| Application layer | `exchange/` (public package) | Import/export engine (Exchange Spec + RSF): discovery, loading, model building, serialization, deserialization, identity/relationship resolver, conflict analyzer, integration engine (staged commit + rollback), package writer — reusable without the CLI. |
| Application layer | `conformance/` (public package) | Validation engine: scanning, artifact classification, rules R0–R12, result model (`Report`); also provides `Scan` and `ParseReference` for other consumers. |
| Application layer | `view/` (public package) | Knowledge Projection Engine: Knowledge Graph (identity index, relationship resolution, membership helpers) + independent projection builders + projection registry — reusable without the CLI. |
| Application layer | `machine/` (public package) | The **machine interface serializer** (ADR-015): CKO → canonical JSON — `Document` (schema `eka-cko-v1`, fixed field order, derived domain/stratum, content payload verbatim, object hash) + `Collection` (domain query envelope, units sorted by canonical form). Pure CKO in, canonical JSON out: never renders, never parses Markdown, never touches storage — reusable by future machine consumers (MCP, Atrium) without the CLI. |
| Application layer | `workspace/` (public package, kernel-internal since ADR-014) | EKA Workspace: workspace resolution (`EKA_HOME` / `~/.eka`), `workspace.json` metadata, project/repository registry, canonical store handle (`Store()` — kernel-internal accessor). Consumers reach it through the runtime services. |
| Application layer | `store/` (public package, kernel-internal since ADR-014) | Canonical store: SQLite schema v2 (`object_payloads` — immutable, content-addressed payloads, insert-only; `object_refs` — mutable references with derived index columns; `attachments`, `sync_log`, `meta`), in-place v1→v2 schema migration, payload insert, reference upserts, sync log, integrity verification (`VerifyIntegrity`). Private persistence behind the Kernel. |
| Application layer | `sync/` (public package, kernel-internal since ADR-014) | Knowledge Runtime synchronization engine: pull (snapshot verification + upsert / docs-mode conformance gate + seed) and push (deterministic snapshot assembly + atomic swap). Orchestration beneath the Authoring API (`Authoring.Sync`). |
| Entry point | `cmd/eka/main.go` | Thin: `os.Exit(cmd.Execute(...))`. Executable name: `eka`. |

```
cmd/                package cmd — Cobra command definitions (command layer)
  root.go           root command + Execute(args, stdin, stdout, stderr) int
  validate.go       validate command
  init.go           init command
  export.go         export command
  import.go         import command
  view.go           view command (renders projections via cmd/ui)
  watch.go          watch command
  get.go            get command (machine retrieval — emits canonical CKO JSON verbatim)
  sync.go           sync command tree (sync / sync pull / sync push)
  project.go        project command tree (project register / project list)
  status.go         status command
  integrity.go      integrity command tree (integrity check)
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
runtime/            public package — the Runtime Kernel (ADR-014): internal
                    Runtime API (Workspace, Knowledge, Resolver, Relations,
                    Timeline, Snapshot, Integrity) + Authoring API
                    (runtime.Authoring: Validate/Compile/Sync); the one
                    sanctioned entry point for consumers
bootstrap/          public package — eka init engine (application layer)
exchange/           public package — import/export engine (application layer)
conformance/        public package — validation engine (application layer)
view/               public package — knowledge projection engine (application layer)
machine/            public package — the machine interface serializer (application layer):
                    CKO → canonical JSON, schema eka-cko-v1 (ADR-015)
workspace/          public package — EKA workspace + registry (application layer;
                    kernel-internal since ADR-014)
store/              public package — canonical store, SQLite schema v2
                    (application layer; kernel-internal since ADR-014)
sync/               public package — synchronization engine (application layer;
                    kernel-internal since ADR-014)
skeletonembed.go    root package — //go:embed skeleton (canonical Reference Skeleton)
```

Principles:

- **Cobra is an adapter, not the architecture.** The framework (currently Cobra) is an implementation detail of the command interface. Business logic lives in `runtime/` (the Runtime Kernel — ADR-014), `bootstrap/`, `exchange/`, `conformance/`, `view/` and `machine/` — public packages imported as `github.com/maleolabs/engineering-knowledge-architecture/runtime`, `.../bootstrap`, `.../conformance`, etc., **with no dependency on `cmd/`**. The kernel-internal packages (`workspace/`, `store/`, `sync/`, `compile/`) sit behind the runtime services: production code outside `runtime/` must not import them — enforced structurally by the import graph.
- **Human and machine paths are separated, formally.** The projection path (`cmd/view.go`, `cmd/watch.go` → `view/`) and the machine path (`cmd/get.go` → `machine/`) share only the Runtime API and the CKO model: the machine path imports `{runtime, exchange}`, never `view/`; the projection path never imports `machine/` (ADR-015 §Decision 5). A rendering change can never touch the machine contract.
- **The command layer calls services, not the other way around.** Future tooling (import/export, graph query, SDKs, Knowledge OS integration, MCP servers) imports the runtime services (and the exchange/export packages) without being tied to Cobra — and without learning storage internals.
- **No `internal/` or `pkg/`** — there is no second internal consumer; the application packages are already public API. The Kernel boundary is enforced by import discipline (ADR-014), not by a directory convention. Adding a directory without an immediate purpose is speculative abstraction (forbidden).
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
| `eka view` | **Implemented** | Knowledge projections (execution / planning / architecture / discovery / operations / ticket; CLI aliases `sprint`, `wave` → execution): read-only views derived from the Engineering Knowledge Model — relationships + State, never markdown text — over Canonical Knowledge Objects read from the workspace canonical store via the runtime services (`Knowledge.UnitsByProject` → `exchange.DecodeUnit`, project union; ADR-013, ADR-014); requires a registered + synced repository (`eka sync` first), unregistered refused (exit 1); rendered as per-domain visualizations (Kanban board, roadmap, dependency tree, information cards, release timeline, detail card). Deterministic, exit codes 0/1/2. |
| `eka watch` | **Implemented** | Realtime projection viewer: same projections as `eka view` (incl. `sprint` / `wave` aliases); TTY-only, polling refresh of the canonical store (`--interval`, default 2s, min 1s), redraw on change only, live refusal frames (repository not registered) with auto-recovery, Ctrl-C to stop (exit `0`). |
| `eka get` | **Implemented** | Machine retrieval interface (ADR-015): canonical CKO JSON (schema `eka-cko-v1`) emitted on stdout — identity lookup (`<ns>/<type>:<id>[:<v>]` via `Resolver.Resolve`, namespace required) or Engineering Domain query (`discovery`/`architecture`/`planning`/`execution`/`operations` via `Knowledge.Search` → domain Collection sorted by canonical form); generated directly from Canonical Knowledge Objects via the Runtime API (`machine/` serializer), never rendered, never parsed from Markdown, never touching SQLite; stdout carries only the JSON; exit 0 / 1 (no workspace, unregistered repository — mirrors `eka view`) / 2 (usage, unknown identity, internal); read-only, never creates a workspace; requires a registered + synced repository (`eka sync` first). |
| `eka validate` | **Implemented** | Full conformance validator (R0–R12: R1–R9 + structural R0 + domain-aware R10–R12). |
| `eka sync` | **Implemented** | Knowledge Runtime synchronization (v0.2.0): pull (snapshot mode: verify + idempotent upsert; docs mode: conformance-gated seed) + push (deterministic snapshot, atomic swap); auto-registration; deletions never applied; exit codes 0/1/2. |
| `eka project` | **Implemented** | Workspace project registry: `register [path] [--name NAME]` (no-op re-registration) + `list` (deterministic, sorted). |
| `eka status` | **Implemented** | Workspace status probe: path, schema version, id, store totals (Objects = references, Payloads = immutable objects, Attachments), per-repository last sync; read-only, never creates the workspace. |
| `eka integrity check` | **Implemented** | Store integrity verification (Immutable Engineering Knowledge Model, ADR-011): recompute payload hashes, strict-decode payloads, verify reference targets + derived index columns, recompute attachment digests, check the repository registry; unreferenced payloads counted as history, never violations; manual modification detected, not prevented; exit codes 0/1/2. |
| `eka version` | **Implemented** | CLI build version + EKA standard version (currently `EKA standard 1.1`). |
| `eka completion` | **Implemented** | bash/zsh/fish/powershell completion scripts (provided by Cobra). |
| `eka diagnose` | Not implemented | Repository diagnostics — future candidate. |
| `eka graph` | Not implemented | Query/knowledge graph over artifacts. |
| Runtime evolution | Future | Deletion protocol, cloud sync, hooks/automation, event-driven watch, wire protocols / MCP integration (wrapping the machine interface), Atrium. |

Future commands are added following the [Contribution guide](#contribution-guide-adding-a-command) — without architectural refactoring.
