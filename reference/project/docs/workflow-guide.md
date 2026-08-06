# Engineering Operating Guide

> Anchor EKA: primary onboarding — how engineering knowledge is produced, organized, validated, projected, exchanged, and consumed.
> Convention document, not an artifact (no `type`/`id`). Companion to [README.md](README.md) (serialization conventions), [lifecycle.md](lifecycle.md) (the concise lifecycle reference), [operating/protocol.md](operating/protocol.md) (how state moves).

This guide is the starting point for working with Engineering Knowledge Architecture (EKA). It is written for an experienced software engineer who has never seen EKA: no prior knowledge of the standard, and no prior knowledge of Scrum, is assumed. The ideas come before the terminology, and canonical terms are introduced with plain-language glosses on first use.

The guide is a journey in twelve parts:

| Part | What it answers |
|---|---|
| 1. Why Engineering Knowledge Matters | Why does any of this exist? |
| 2. The Mental Model | What is EKA, underneath all tools? |
| 3. Engineering Knowledge Lifecycle | How does knowledge flow? |
| 4. Engineering Domains | Where does knowledge sit, and who outranks whom? |
| 5. Knowledge Evolution | How does an idea become operations? |
| 6. Methodology Independence | What does EKA have to do with Scrum? |
| 7. Daily Workflow | What does an engineer actually do? |
| 8. AI Workflow | How do AI agents fit in? |
| 9. CLI Workflow | What does the `eka` command-line tool do? |
| 10. Knowledge Projection | What are the read-only views? |
| 11. Common Mistakes | What do people get wrong? |
| 12. Complete Lifecycle Example | One project, end to end |

## 1. Why Engineering Knowledge Matters

Software projects rarely fail because of the code. They fail because the knowledge around the code fragments.

Three failure modes are familiar to every experienced engineer:

- **Decisions get lost.** "Why did we choose this database?" was answered once in a chat thread, and the answer died with the thread. Six months later, someone reopens the question at ten times the cost.
- **State gets duplicated.** The same work item has a status in the ticket tool, a spreadsheet, and a slide deck. Nobody can say which one is true, and keeping all three in sync is a full-time job.
- **Knowledge dies with people.** The person who understood the billing system leaves, and the understanding leaves with them. The docs, such as they are, describe what the system does — not why it does it.

EKA — **Engineering Knowledge Architecture** — exists to make knowledge survive. The word "documentation" gets in the way of that sentence, so it is worth being precise.

EKA is **not a documentation system**. A documentation system organizes documents: folders, templates, formatting. EKA is a **model** — a way of representing engineering knowledge that does not depend on any particular tool, file format, or process. The unit of that model is the **Artifact** (glossary: a piece of engineering knowledge — a requirement, a decision, a plan, a work item, a runbook) that carries four things:

- **Identity** — a permanent name that never changes, no matter where the knowledge lives.
- **State** — the status fields the artifact itself owns, with their history.
- **Content** — what the knowledge actually says.
- **Relationships** — explicit links to other artifacts.

Because knowledge is modeled this way, it can move between tools without losing its meaning. That is the whole point: **the model survives the tools.** The [EKA Core Specification](../standard/eka-specification-v1.1.md) defines the model; this repository is one concrete implementation of it (Git + Markdown). Other storage — relational databases, graph stores, future platforms — is equally valid as long as it honors the same contracts.

## 2. The Mental Model

The entire architecture reduces to one chain:

```
Engineering Knowledge
        │
        ▼
     Identity       a permanent name for a piece of knowledge
        │
        ▼
      State         the status fields the artifact owns
        │
        ▼
  Relationships     explicit links to other artifacts, by identity
        │
        ▼
    Lifecycle       how knowledge flows over time
```

Read the chain slowly: every other idea in EKA hangs off it.

- **Identity**: a permanent name for a piece of knowledge that never changes, no matter where it lives. Not a file path, not a folder — a name. You refer to knowledge by Identity, so renames and reorganizations can never silently break your references.
- **State**: a fact about where an artifact is in a process — is this requirement approved? is this work item in progress? State is owned: each field has exactly one writer, and every change is recorded in a **Change Log** (the chronological record of transitions).
- **Content**: the semantic payload — what the knowledge actually says. Content and State are separate channels: content matures through review and approval; state moves through the execution protocol. The two never mix.
- **Relationships**: explicit links to other artifacts by Identity — *derives from*, *depends on*, *supersedes*, *amends*, *validates*. A specification derives from a requirement; a work item depends on the plan.
- **Lifecycle**: the movement of State over time — produce, organize, validate, project, exchange, consume (Section 3).

Together, Identity + State + Content + Relationships make an **Artifact**. The **State Vector** is the set of state fields an artifact owns, plus their history.

Here is the idea that takes most people a moment: **Markdown is one possible representation of this model — not the model itself.** The same knowledge wears three faces:

```
        ┌──────────────────────────┐
        │  Markdown files          │   authoring representation
        └────────────┬─────────────┘
                     │
        ┌────────────▼─────────────┐
        │   Engineering            │
        │   Knowledge Model        │
        │   (Identity · State ·    │
        │    Content · Relations)  │
        └───────┬──────────┬───────┘
                │          │
   ┌────────────▼─────┐  ┌─▼─────────────────────┐
   │  Projections     │  │  Knowledge Packages   │
   │  (viewing)       │  │  (exchange)           │
   └──────────────────┘  └───────────────────────┘
```

- You **author** knowledge in Markdown.
- You **view** it through projections — read-only renderings derived from the model, never from file text.
- You **exchange** it as Knowledge Packages — portable, self-contained representations that move between systems without loss.

Three faces, one model underneath. If you keep one sentence from this guide, keep this one: **the model is the source of truth; every file, view, and package is a representation of it.** The concise lifecycle reference is [lifecycle.md](lifecycle.md); the serialization conventions of this repository live in [README.md](README.md).

## 3. Engineering Knowledge Lifecycle

The **Engineering Knowledge Lifecycle** is the six-stage flow of knowledge through the model: **Produce → Organize → Validate → Project → Exchange → Consume**. An artifact is born, shaped, checked, seen, moved, and used:

```
Produce ──▶ Organize ──▶ Validate ──▶ Project ──▶ Exchange ──▶ Consume
 (born)     (shaped)     (checked)    (seen)     (moved)      (used)
```

The lifecycle is one of three orthogonal axes (the other two arrive in Section 4): an artifact moves through the lifecycle while its domain and its stratum never change.

### 3.1 Produce — capture knowledge before it evaporates

**What happens.** Knowledge enters the repository at the point of creation, in two kinds:

- **Ephemeral** — session notes, spike investigations, research findings: captured while the work happens, low ceremony, usually drafts.
- **Durable** — requirements, architecture descriptions, decisions, specifications, standards: started as drafts and matured through governance.

**Distillation** is the bridge from ephemeral to durable: findings from sessions, spikes, and reviews are distilled into the durable artifact they affect — a direction becomes a decision, a proven procedure becomes a runbook. Archiving ephemeral knowledge without distilling it first violates the protocol (see [operating/protocol.md](operating/protocol.md)).

**What you do.** Write. Capture the finding in the moment, in the dimension it belongs to. When a session produces a durable insight, create the durable artifact it implies. Producing is a writing step — no ceremony, no templates.

**Value for software delivery.** Intent is captured at the cheapest possible moment. A requirement written while the conversation is fresh costs minutes; the same requirement reconstructed six months later costs a design cycle.

### 3.2 Organize — classify knowledge so it can be found and governed

**What happens.** Every artifact is classified: identity fields are recorded, the artifact is placed in the dimension that matches its content, and relationships are wired to the higher-stratum knowledge it builds on. Every state transition is recorded in the Change Log by its single owner.

**What you do.** Fill in the classification truthfully and wire the relationships: a specification derives from a requirement; a work item depends on the plan. The Engineering Domain (Section 4) is derived from the artifact's type — you rarely need to declare it.

**Value for software delivery.** Findable, related knowledge is auditable knowledge. "Why does this work item exist?" gets a one-link answer instead of an archaeology dig.

### 3.3 Validate — the conformance gate

**What happens.** A mechanical rule set (R0–R12) runs over every artifact before commit: identity uniqueness, state validity, referential integrity, classification coherence, change-log consistency. Drafts get tolerance: dangling references are warnings, and draft knowledge is exempt from stratification traceability. The full checklist lives in [exchange/validation.md](exchange/validation.md).

**What you do.** Run the validator before committing. Treat warnings as quality signals, errors as blockers. When a lower-stratum artifact contradicts higher-stratum knowledge in force, fix the lower one — never the higher one (the Stratum Authority Invariant, Section 4).

**Value for software delivery.** Conformance is mechanical, so quality does not depend on memory or mood. The same gate that runs on your laptop can run in CI on every commit — non-conformant knowledge never ships silently.

### 3.4 Project — see the current picture

**What happens.** The repository is projected into read-only views: one per Engineering Domain, plus a per-ticket view. A projection is derived from the model — identity, state, relationships — never from file text. Projections are refreshed on read, never edited (Section 10).

**What you do.** Read projections to answer "what is planned, in progress, approved, done" — then edit the source of truth, not the view.

**Value for software delivery.** One shared picture for standups, planning, and reporting. Every reader sees the same current state, because every reader reads the same model.

### 3.5 Exchange — move knowledge between systems

**What happens.** Export builds a Knowledge Package carrying Identity, State (with full change log), Content, and Relationships. Import integrates the package **atomically** with a **conservative merge**: new artifacts are written, identical duplicates are skipped, any difference conflicts and aborts — no silent overwrites. The round-trip contract is defined in the [EKA Exchange Specification](../standard/eka-exchange-specification-v1.1.md).

**What you do.** Export to share knowledge with another repository or team; import to receive it. The package is how knowledge leaves and enters a system without loss.

**Value for software delivery.** Knowledge moves between teams, vendors, and platforms without a lossy copy-paste step. The round-trip is guaranteed by contract, not by hope.

### 3.6 Consume — run your work from knowledge

**What happens.** Knowledge is read and executed: runbooks are executed as procedures; specifications and decisions are read as the authority for what is built and why; tickets and work items are read as the current execution plan.

**What you do.** Read by Identity and by domain, never by location. Execute the runbook; honor the specification. When a runbook contradicts its approved specification, fix the runbook — it sits on a lower stratum.

**Value for software delivery.** Knowledge stops being a record of the past and becomes the operating instructions of the present: runbooks run, specifications constrain, plans steer.

Each stage maps directly to engineering value: **Produce** captures intent before it evaporates; **Consume** gives you knowledge you can actually run your work from. The stage-by-stage expansion is [lifecycle.md](lifecycle.md) — the concise reference this guide complements.

## 4. Engineering Domains

Every artifact belongs to exactly one **Engineering Domain** — the stratum-aligned category of engineering knowledge it holds. There are exactly five, in a strict authority order:

```
Discovery (1 — highest authority)
     ▼
Architecture (2)
     ▼
Planning (3)
     ▼
Execution (4)
     ▼
Operations (5 — lowest authority)
```

Two canonical terms carry this section:

- **Engineering Domain** — where the knowledge sits: one of five categories, derived from the artifact's type, never declared by hand.
- **Knowledge Stratum** — its authority: a fixed position in the order above, always derived from the domain. Stratum 1 has the highest authority; stratum 5 the lowest.

**Stratum Authority Invariant**, in plain words: *what you decided earlier and approved has authority over what you do later; lower strata refine, never redefine, higher strata.* If a runbook (Operations) contradicts an approved specification (Architecture), the runbook is wrong — you fix it by writing a corrected version, never by quietly overriding the specification. And a lower-stratum artifact never supersedes or amends a higher-stratum one.

### 4.1 Discovery — stratum 1 (highest authority)

**Purpose.** The home of intent, requirements, and research — why the product exists and what it must do.

**Responsibilities.** Owns vision, strategy, requirements, and research findings.

**Knowledge produced.** Vision, strategy, requirement documents, research findings.

**Knowledge consumed.** Nothing from other strata — nothing ranks above it.

**Relationship to other domains.** Root of the authority chain; constrains every domain below.

**Knowledge Stratum.** 1.

**Example.** "We are building a billing feature and customers must be able to download invoices" — *what we're building and why*.

### 4.2 Architecture — stratum 2

**Purpose.** The home of durable design knowledge: decisions, specifications, standards, vocabulary — the knowledge that binds planning and execution.

**Responsibilities.** Owns architecture descriptions, decisions, specifications, standards, glossary terms.

**Knowledge produced.** Architecture descriptions, decision records, specifications, standards, glossary terms.

**Knowledge consumed.** Discovery — requirements, research, intent: decisions and specifications derive from requirements.

**Relationship to other domains.** Derives from Discovery; constrains Planning, Execution, Operations.

**Knowledge Stratum.** 2.

**Example.** "We will use provider X for billing, because..." — a decision about which database to use, which provider to pick, how the system is shaped.

### 4.3 Planning — stratum 3

**Purpose.** The home of commitment knowledge: scope, epics, and plans that convert requirements and decisions into an executable commitment.

**Responsibilities.** Owns scope definitions, epics, plans, traceability artifacts.

**Knowledge produced.** Scope definitions, epics, plans, traceability.

**Knowledge consumed.** Discovery (requirements, intent) and Architecture (decisions, specifications, standards): scope derives from requirements; plans honor approved decisions.

**Relationship to other domains.** Derives from Discovery and Architecture; constrains Execution and Operations.

**Knowledge Stratum.** 3.

**Example.** "Release 1 ships invoice download; release 2 ships usage reports" — the roadmap that turns decisions into a sequence of work.

### 4.4 Execution — stratum 4

**Purpose.** The home of work in motion: containers, tickets, work items, reviews, sessions — knowledge that carries work through the execution protocol.

**Responsibilities.** Owns execution containers, tickets, work items, reviews, sessions.

**Knowledge produced.** Execution containers, tickets, work items, reviews, sessions.

**Knowledge consumed.** Planning (plans, scope), Architecture (specifications, standards, decisions), Discovery (requirements): execution follows the plan and the specifications it must satisfy.

**Relationship to other domains.** Derives from Planning, transitively from Architecture and Discovery; constrains Operations.

**Knowledge Stratum.** 4.

**Example.** The sprint's ticket list and the work items behind the tickets — where statuses actually live.

### 4.5 Operations — stratum 5 (lowest authority)

**Purpose.** The home of operational knowledge and records: runbooks and release records that preserve what ran and how it runs.

**Responsibilities.** Owns runbooks and release records.

**Knowledge produced.** Runbooks, release records.

**Knowledge consumed.** Execution (the execution aggregate, review and session findings) and Architecture (standards): release records preserve what was executed; runbooks comply with standards.

**Relationship to other domains.** Derives from Execution, transitively from all higher strata; constrains nothing.

**Knowledge Stratum.** 5.

**Example.** The runbook you execute during a billing outage — the most actionable stratum.

The full registry of the five domains — produced and consumed knowledge, constraints, aliases — is the Engineering Domain Registry in the [EKA Core Specification](../standard/eka-specification-v1.1.md) (§8.2); the authoritative token→domain mapping is §8.1.

## 5. Knowledge Evolution

The lifecycle (Section 3) describes how one artifact flows. Knowledge evolution describes how a whole *idea* flows through the strata:

```
Idea ──▶ Discovery ──▶ Architecture ──▶ Planning ──▶ Execution ──▶ Operations
          what & why        how          commitment    work done     running
```

A concrete story: a support complaint about slow search results.

1. **Idea.** Captured in a session note during the support call — ephemeral, low ceremony: "search is slow, users are complaining".
2. **Discovery.** Distilled into a requirement: "search results must return in under a second". A research finding compares search providers.
3. **Architecture.** A decision records the provider choice, deriving from the requirement. A specification defines the indexing contract.
4. **Planning.** A scope and a plan commit the work: release 1 covers the provider swap. The plan is approved, and later locked when execution starts.
5. **Execution.** A container opens; work items move from planned to done; reviews validate the result; tickets track status.
6. **Operations.** A release record captures what shipped; a runbook documents the rollout and rollback.

Each step *derives from* the one before it. That chain of relationships is what makes knowledge auditable: you can always ask "why does this work item exist?" and walk the links up to the requirement.

One rule governs correction, and it is worth stating plainly (P7 in the standard): **correction is forward-only.** When a decision changes, you do not quietly edit the old one — you write a new decision that supersedes it, and the old one stays as a record. The same holds for state: a work item moves forward — planned → in progress → done — and never backward; a mistaken status is corrected with a new instance and a relationship, never by silent edits. History is preserved, and knowledge never silently changes under you.

## 6. Methodology Independence

EKA is methodology-independent by design. Scrum, Kanban, Shape Up, and internal workflows are **convention layers**: they decide cadence, rituals, roles, and how work is scheduled — they never decide what knowledge is, where it sits, or how it is governed.

The bridge between methodology terms and the model is the **Representation Alias** (glossary: a methodology term mapped onto a canonical artifact type + Engineering Domain). The canonical catalog lives in the [Representation Alias Registry](../standard/representation-alias-registry-v1.1.md). A few examples:

| Methodology term | Alias of | Engineering Domain |
|---|---|---|
| PRD | `req-` (requirement) | Discovery |
| ADR / RFC | `adr-` (decision) | Architecture |
| Epic | `epc-` (epic) | Planning |
| Sprint / Iteration | `ctr-` (execution container) | Execution |
| Story | `sto-` (work item) | Execution |
| Ticket | `tkt-` (ticket) | Execution |
| Release | `rel-` (release record) | Operations |
| Runbook | `run-` (runbook) | Operations |

A team running Scrum still writes its PRD as a `req-` artifact in Discovery and its Sprint as a `ctr-` container in Execution. The methodology names the rituals; EKA names the knowledge.

Two consequences are worth remembering:

- **Aliases are never frontmatter values.** You never write `type: sprint`. The type token is canonical (`ctr-`); "sprint" is a word people use, not a value the model accepts.
- **The CLI speaks both languages.** `eka view sprint` and `eka view wave` are accepted as aliases for the canonical execution projection, with identical output.

In one sentence: **EKA does not standardize Scrum; it standardizes engineering knowledge.**

## 7. Daily Workflow

Strip away the model and the daily loop is small:

```
Create ──▶ Validate ──▶ Refine ──▶ Review ──▶ Project ──▶ Deliver ──▶ Maintain
  │           │           │          │           │           │           │
 write     check      respond     get it      see the    ship it    keep it
  │         structure  to gate    approved    picture   + record    alive
```

- **Create.** Write the knowledge your work produces: a session note during the investigation, a requirement when the scope is clear, a decision when a direction is chosen, a work item when the work is committed. The rule of thumb: *if you said it out loud, it probably belongs in the repository* — as ephemeral knowledge (session notes, findings) or durable knowledge (requirements, decisions), whichever fits.
- **Validate.** Before you commit, run the conformance check (hint: `eka validate`). Errors block; warnings inform.
- **Refine.** Fix what the gate flags: classification, relationships, change-log entries. This is where "organize" becomes a habit instead of a ceremony.
- **Review.** Move content through review toward approval. An approved requirement or decision is *in force* — it constrains everything below it.
- **Project.** Look at the picture, not the files: the execution view for "what is in progress", the planning view for "what is committed", the ticket view for "what is the status of this item" (hints: `eka view execution`, `eka view planning`, `eka view ticket <id>`).
- **Deliver.** When the work lands, record it: a release record for what shipped, a runbook for how to operate it, distillation of session findings into durable artifacts.
- **Maintain.** Supersede rather than rewrite: when a decision is outdated, write the new one and link it. Records stay; nothing is silently deleted.

Two timing questions answer themselves once the loop is familiar:

- **When do I write what?** In the moment, as the knowledge is produced — before it evaporates. Sessions and findings are cheap; decisions deserve a real artifact; work items exist when work is committed.
- **When do I check state?** Before every commit (validation) and whenever a status question comes up (projections). State lives in the source of truth; you check it through projections.

The full command reference is Section 9 and [reference/cli.md](../reference/cli.md); the state machine that governs transitions is [operating/protocol.md](operating/protocol.md).

## 8. AI Workflow

AI agents participate in the same model, on the same terms as humans. There is no separate AI format and no separate AI lane. Five roles are natural:

- **Generate initial knowledge.** Draft an artifact from a conversation: a session transcript becomes a requirement draft; a discussion becomes a decision draft. The human stays the author; the AI accelerates the capture.
- **Refine existing knowledge.** Distill sessions and reviews into durable artifacts: findings become decisions, proven procedures become runbooks. The distillation obligation (Section 3.1) does not change because an AI did the writing.
- **Review consistency.** Compare knowledge across strata and flag contradictions: does this work item contradict the approved specification? Does this runbook honor the standard? Contradictions resolve downward — the lower stratum changes.
- **Validate structure.** Check frontmatter, state, relationships, and change-log consistency before the mechanical gate runs — a fast feedback loop for the same rules the validator enforces.
- **Generate projections.** Summarize execution state for a standup, a report, or a decision — the same derivation rules the projections use, applied to the same model.

The invariant is one line (P16 in the standard): **enforcement mechanisms may vary; the invariants do not.** Whatever produces an artifact — human or AI — the conformance gate applies identically. AI output that fails validation is not special; it fixes its frontmatter like anyone else's. AI assists Engineering Knowledge; it never replaces it.

## 9. CLI Workflow

The CLI (`eka`) is introduced here, after the model, because it is best understood as a set of lifecycle tools — not as the point of EKA. Each command serves one stage:

| Command | Lifecycle stage | What it does for you |
|---|---|---|
| `eka init` | Organize | Bootstraps a repository from the reference skeleton; idempotent, never overwrites your content |
| `eka validate` | Validate | Runs the mechanical conformance gate (R0–R12) over the repository |
| `eka view` | Project | Renders read-only projections per domain, plus the ticket view (Section 10) |
| `eka export` | Exchange | Builds a Knowledge Package from the repository — validated first; a non-conformant repository exports nothing |
| `eka import` | Exchange | Integrates a package atomically with a conservative merge — new artifacts written, duplicates skipped, conflicts abort |
| `eka version` | — | Reports the CLI build and the EKA standard version implemented |

Two things the CLI deliberately does *not* do:

- **It never writes content for you.** Producing knowledge is plain editing. The CLI organizes, validates, projects, and exchanges — it does not author.
- **It does not manage people.** There is no workflow engine, no assignment, no notifications. The Operating Layer governs state transitions; the CLI makes them checkable and visible.

The UX conventions are worth knowing because they make the tool trustworthy:

- **Calm output.** No banners, no ALL-CAPS, no decoration. Every line answers "what is happening" and "what was the outcome".
- **Context header.** Every command opens by identifying the object it processes — repository, package, projection — before any action.
- **Deterministic.** Identical input produces identical output, in CI and in a pipe. Exit codes: 0 success, 1 blocking violation, 2 usage error; warnings never affect the exit code.

The full reference — installation, commands, output shapes — is [reference/cli.md](../reference/cli.md). This guide describes capabilities, not flags: the command set may evolve, but the lifecycle mapping above will not.

## 10. Knowledge Projection

A **Projection** is a read-only view derived from the Engineering Knowledge Model. Recall the three faces of Section 2:

```
Markdown files ──▶  the model  ──▶ Projections         (viewing)
                  (source of   ──▶ Knowledge Package   (exchange)
                   truth)
```

Each face serves one job: you *author* in Markdown, you *view* through projections, you *exchange* as packages. All three originate from the same model.

The rule that makes projections safe is P6, in plain words: **you never edit a projection; you edit the source of truth and re-render.** A projection has no state of its own and never becomes a writer. If the execution view says a work item is in progress and it is not, you do not edit the view — you change the work item's state in the source of truth, and the view refreshes on next read.

The canonical projections, one per Engineering Domain, plus the ticket view:

| Projection | What it shows |
|---|---|
| `eka view discovery` | Vision, strategy, requirements, research findings |
| `eka view architecture` | Decisions, architecture descriptions, specifications, standards, vocabulary — grouped by content state (draft / review / approved) |
| `eka view planning` | Scope, epics, plans, traceability — grouped by planning state and phase |
| `eka view execution` | The active container's tickets and work items, grouped by execution state (planned / todo / in-progress / in-review / done) |
| `eka view operations` | Runbooks and release records |
| `eka view ticket <id>` | One ticket's projected status — derived from its owner work item, never from ticket text |

Projections are derived from the model — identity, state, relationships — never from parsing markdown text. That is what makes them deterministic and trustworthy: two engineers, or an engineer and an agent, always see the same picture. The CLI-level aliases `sprint` and `wave` resolve to the `execution` projection with identical output — methodology vocabulary stays convention-layer, the projection model stays canonical.

## 11. Common Mistakes

Six misconceptions recur. Each is wrong for a specific reason; each has a specific truth behind it.

**1. "EKA is a markdown template."**
*Why it's wrong:* there is nothing to memorize — no required headings, no document skeleton, no template gallery. The model is Identity, State, Content, and Relationships; Markdown is just the serialization this repository uses.
*What's actually true:* EKA is a model that survives any storage. Folders and templates are conveniences of one implementation.

**2. "EKA replaces Scrum."**
*Why it's wrong:* EKA has no opinion on cadence, roles, or rituals. It does not schedule, assign, or estimate.
*What's actually true:* Scrum (or Kanban, or Shape Up) runs as a convention layer on top. The methodology maps its words onto the model through Representation Aliases; EKA standardizes the knowledge underneath.

**3. "EKA requires ADR."**
*Why it's wrong:* nothing forces you to record decisions in a particular form. A team that writes a decision as a paragraph in a session note and never promotes it has simply chosen ephemeral knowledge.
*What's actually true:* durable decisions belong in Architecture as `adr-`/`dec-` artifacts because that is where their authority sits — but the choice of what to promote is yours. Distillation is a protocol obligation, not a template mandate.

**4. "EKA requires PRD."**
*Why it's wrong:* the same logic. "PRD" is an alias for the requirement artifact type, and requirements are one of several kinds of Discovery knowledge.
*What's actually true:* intent, requirements, and research are the Discovery domain's business. You write a requirement when a requirement exists — the name you call it is convention.

**5. "EKA is a project management tool."**
*Why it's wrong:* it has no boards, no dashboards, no notifications, no permissions. It does not run your process.
*What's actually true:* EKA is the knowledge substrate your tools read from. A sprint board is a projection of the model; a project management tool can consume the same model through packages.

**6. "I must put status in my ticket."**
*Why it's wrong:* the ticket is a projection — it carries no state of its own. Its "status" is derived from the work item it references (P6, single-writer). Editing the ticket's text to say "done" changes nothing.
*What's actually true:* state lives in exactly one place — the owner. Work items own Execution State; tickets project it. Update the owner; the ticket follows on next read.

## 12. Complete Lifecycle Example

One walkthrough, end to end: **shipping a billing feature**. Names are illustrative — the point is the pattern, not the exact IDs. No frontmatter is shown; artifacts are referred to by their name pattern, which mirrors their Identity.

**Step 1 — the idea (Produce).** A support call becomes a session note (`ses-billing-idea`): customers keep asking for invoice downloads. Cheap, ephemeral, captured in the moment.

**Step 2 — Discovery (Produce → Organize → Validate).** The session is distilled into a requirement (`req-invoice-self-service`): customers can download invoices as PDF. A research finding compares providers (`fnd-payment-providers`). The requirement derives from the product's vision and strategy; it is approved through review and is now in force.

**Step 3 — Architecture (Organize → Validate).** A decision records the provider choice (`adr-billing-provider`), deriving from the requirement. A specification defines the invoice format (`spec-invoice-format`). Both are approved; both now constrain everything below.

**Step 4 — Planning (Organize → Validate).** A scope definition commits release 1 (`scp-billing`), an epic groups the work (`epc-invoice-self-service`), and a plan sequences it (`plan-release-1`), deriving from the requirement and the decision. The plan is approved; when execution starts it locks — a later change would be a new instance, never an edit.

**Step 5 — Execution (Validate → Project → Consume).** A container opens (`ctr-sprint-12`). Work items are created under it: `sto-invoice-list`, `sto-pdf-generation`. Tickets project their status (`tkt-invoice-list`). The team reads the execution projection daily; work items move forward — planned → in-progress → in-review → done. A review records the verdict (`rvw-sprint-12`).

**Step 6 — Operations (Validate → Consume).** A release record captures what shipped and what the execution aggregate was (`rel-release-1`). A runbook documents the rollout and rollback (`run-invoice-rollout`), complying with the standards it touches. When the runbook contradicts the specification, the runbook is corrected — never the specification.

**Step 7 — Exchange (Exchange).** The knowledge is exported as a package and imported into the operations repository: same identities, same state, same relationships — no loss. Re-importing the same package is a no-op; a conflicting version aborts rather than silently merging.

The artifact trail, by domain:

| Engineering Domain | Artifacts (name pattern) | Tokens |
|---|---|---|
| Discovery | `req-invoice-self-service`, `fnd-payment-providers` | req-, fnd- |
| Architecture | `adr-billing-provider`, `spec-invoice-format` | adr-, spec- |
| Planning | `scp-billing`, `epc-invoice-self-service`, `plan-release-1` | scp-, epc-, plan- |
| Execution | `ctr-sprint-12`, `sto-invoice-list`, `sto-pdf-generation`, `tkt-invoice-list`, `rvw-sprint-12` | ctr-, sto-, tkt-, rvw- |
| Operations | `rel-release-1`, `run-invoice-rollout` | rel-, run- |

And the lifecycle stages the example passed through: **Produce** (session note, requirement), **Organize** (classification, relationships, change-log), **Validate** (the gate before every commit), **Project** (the execution and planning views the team read), **Exchange** (export to, and import into, the operations repository), **Consume** (runbook executed, release record read).

## Where to go next

- [lifecycle.md](lifecycle.md) — the concise lifecycle reference this guide expands.
- [README.md](README.md) — serialization conventions: identity, state, classification, relationships, projections; the navigation map of `docs/`.
- [operating/protocol.md](operating/protocol.md) — the Operating Manual: how state moves, transitions, gates, distillation.
- [reference/cli.md](../reference/cli.md) — the CLI reference: commands, UX conventions, exit codes.
- [EKA Core Specification](../standard/eka-specification-v1.1.md) — the canonical model; §8.1–8.3 define the Engineering Domains, the Domain Registry, and Stratification Governance.
- [Representation Alias Registry](../standard/representation-alias-registry-v1.1.md) — the canonical catalog of methodology-term aliases.
- [EKA Exchange Specification](../standard/eka-exchange-specification-v1.1.md) — round-trip guarantees and import semantics behind the Exchange stage.
- [Reference Serialization Format](../reference/eka-reference-serialization-format-v1.1.md) — the package format used by `eka export` and `eka import`.

### Onboarding recommendations

Suggested learning path: read this guide once, then run `eka init` in a scratch repository and write one artifact per domain — a requirement, a decision, a plan, a work item, a runbook — validating and viewing after each one, then finish with an export/import round-trip. An interactive tutorial and a video walkthrough would cover the same ground more gently; neither is promised yet.
