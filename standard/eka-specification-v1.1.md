# Engineering Knowledge Architecture (EKA) — Canonical Specification

| Field | Value |
|---|---|
| **Version** | 1.1 |
| **Revision** | v1.1 — additive taxonomy evolution (Engineering Domains + Knowledge Stratification); v1.0 contract unchanged |
| **Status** | Ratified |
| **Authority** | The canonical standard of the engineering knowledge model |
| **Scope** | Concepts, invariants, and contracts — not implementation mechanisms |

**Reading this document:** capitalized terms refer to the definitions in Section 3. The word "must" denotes a requirement binding on all implementations. The word "may" denotes an implementation option within the bounds of the contract.

---

## 1. Engineering Knowledge Architecture Overview

### 1.1 Definition of the standard

EKA is the canonical conceptual model for engineering knowledge: the definition of Artifact, Identity, State, knowledge taxonomies, the layer architecture, and exchange contracts between systems. This standard originated from two responsibilities that proved to coexist in the initial implementation: the **Knowledge Base** (storage and preservation of knowledge) and the **Engineering Operating System** (deterministic execution, agent coordination, governance). Both are two layers of one system, bound together by Artifact Identity.

### 1.2 Scope

This standard establishes:
- fundamental concepts (Section 3) and principles (Section 2);
- the layer architecture and its contracts (Sections 4–5);
- the Identity Model (Section 6) and State Taxonomy (Section 7);
- knowledge, execution, and artifact taxonomies (Sections 8–10);
- the conceptual lifecycle model (Section 11);
- storage-independent and exchange contracts (Sections 12–13);
- the extension and evolution model (Sections 14–16).

This standard does **not** establish: serialization formats, storage layout, directory structure, document templates, naming schemes, query languages, or specific enforcement mechanisms. All of these are implementation decisions that must comply with the contract in this document.

### 1.3 Relationship to implementations

**An implementation is one serialization of this architecture — not the architecture itself.** An implementation may be: a git repository, a relational database, a graph database, an object store, an AI-native knowledge store, a Knowledge OS, or an import/export pipeline. Every implementation must: (a) satisfy all invariants (Section 5.4), (b) provide Identity resolution, (c) support lossless exchange through the Exchange Layer contract (Section 13).

### 1.4 Long-term position

EKA is the canonical engineering knowledge model that can be integrated with a future Knowledge OS through standardized import/export mechanisms. The architecture survives changes of storage medium because it is built on concepts — Identity, State, Layer, contracts — not on mechanisms.

---

## 2. Architectural Principles

| # | Principle | Rationale |
|---|---|---|
| P1 | **Separation of Concerns** | Knowledge and execution responsibilities are separated into layers with explicit contracts. Fusing both into a single hierarchy is the source of conflation. |
| P2 | **Explicit State** | State is always explicit as owned metadata — never implicit in structure. Structure may be a projection of State, but never an independent fact. |
| P3 | **Stable Identity** | Identity is immutable and independent of location, storage, State, and classification. References always use Identity, never location. |
| P4 | **Protocol vs Content Distinction** | Protocol is a property of the Operating Layer; Content is a property of the Knowledge Layer. Every artifact serving both defines the two separately. |
| P5 | **Layer Independence** | Each layer may evolve without overhauling the others: taxonomy changes without touching protocol; protocol is strengthened without shifting the home of knowledge. |
| P6 | **Single Writer** | Every State field has exactly one owner. Any other view is a State Projection — generated or validated, never edited as an independent fact. |
| P7 | **Forward-Only Transitions** | All State Domains move forward without regression. Corrections are made with a new instance + Relationship, not by mutation. |
| P8 | **Approved-Content Immutability** | Content that has passed the approval gate is not silently mutated. Changes occur only through the governance channel. A prerequisite for preservation. |
| P9 | **Structure as Projection of State** | Structural organization/position is derived from State and Identity — structure never becomes a second fact that can drift from State. |
| P10 | **Two Change Channels** | The Content channel (governance) and the State channel (protocol) are separate. Mixing them is a violation. |
| P11 | **Determinism by Protocol** | Execution order is defined by protocol: "what next" is always answered. Enforcement is an implementation capability; the requirement is the standard. |
| P12 | **Preservation Over Deletion** | History is knowledge. Even wrong decisions are preserved. Superseded/Archived are Records, not garbage. |
| P13 | **Lossless Exchange** | Exchange between systems must not lose or duplicate Identity, State, Content, or Relationship. |
| P14 | **Minimum Canonical Core** | The standard establishes concepts and contracts; implementations choose mechanisms. The smaller the core, the longer the life of the standard. |
| P15 | **Classification is Property, Not Identity** | Knowledge Dimension is an artifact property; a classification change never breaks references. |
| P16 | **Enforcement Capability Varies, Invariants Don't** | Enforcement mechanisms differ across implementations (structural constraints, database constraints, validation); the invariants to be enforced are identical. |

---

## 3. Core Concepts

Precise definitions — every term is used with a single meaning throughout this document.

- **Artifact** — an engineering knowledge entity that has Identity, Content, a State Vector (the State domains it **owns**), and Relationship. The basic unit of the model.
- **Content** — the semantic payload of an Artifact: intent, decisions, design, constraints, procedures, notes. Belongs to the Knowledge Layer.
- **Well-formed Content** — Content that conforms to the structure established for its Artifact type, so that it can be parsed and executed deterministically.
- **Identity** — the property of an Artifact that permanently distinguishes it from all other Artifacts: `(Namespace, Type, ID[, InstanceVersion])`. See Section 6.
- **Artifact Line** — the enduring Identity entity: one `(Namespace, Type, ID)`.
- **Artifact Instance** — one version of a Line's existence: Line + `InstanceVersion`.
- **State** — a fact about the position of an Artifact within a given process.
- **State Vector** — the tuple of State Domains **owned** by an Artifact. Projected State (State Projection) is not part of the State Vector.
- **State Domain** — an independent State dimension with its own semantics, owner, and transition rules (Section 7). Domains are orthogonal.
- **State Projection** — a State view derived from the owner (example: aggregate of work item State → Execution Container status). A projection has no State of its own; a projection never becomes a writer.
- **Projection Semantics** — the computation rules of a State Projection: from which owner, how aggregation works, what is displayed.
- **Projection Refresh** — the mechanism and timing of State Projection validation against the owner (on-read and/or event).
- **Knowledge Dimension** — an axis of knowledge classification (Section 8). An Artifact property, not Identity.
- **Protocol** — deterministic rules owned by the Operating Layer: ordering, State transitions, locking, gates, execution commands.
- **Layer** — an architectural layer with explicit contracts: Knowledge Layer, Operating Layer, Exchange Layer (Section 4).
- **Namespace** — an Identity space that separates management domains (products, organizations, systems).
- **Relationship** — an explicit relation between Artifacts referenced by Identity: supersedes, amends, derives-from, depends-on, validates.
- **Gate** — a condition that must be satisfied before a transition or execution may occur (approval gate, readiness gate, review gate).
- **Command** — a deterministic execution instruction consumed by an executor (human or agent). Command Content is Content (owned by the Knowledge Layer); its execution is governed by Protocol (Operating Layer).
- **Execution Container** — an execution Artifact that wraps work items and carries a concurrency convention (exactly-one-active). Its State Domain: Container State. Example implementation: sprint.
- **Phase** — product/scope context over time (Discovery, MVP, Milestone, Release). Phase is a **context attribute** on planning/scope Artifacts: not a category, not a State Domain. **Phase change** is a context update authorized by a readiness Gate (Sections 7.5, 11.2).
- **Record** — an Artifact preserved as history (Superseded, Archived, Retired, release record) with immutable Content.
- **Distillation** — transformation of ephemeral knowledge (working context, review findings) into durable knowledge (decisions, ADRs, Records).
- **Change Log** — the chronological record of State transitions on an Artifact: domain, old value, new value, time, authority. Mandatory for all State Domains.
- **Identity Registry** — a Knowledge Layer function that guarantees Identity uniqueness, Identity resolution to Artifacts, and referential integrity.
- **Trigger** — a relationship between domains/operations in which an event triggers validation or a transition in another domain (formal definition of interaction, Section 7.5).
- **Knowledge OS** — a future knowledge execution platform that consumes and produces EKA Artifacts through the Exchange Layer. Not part of the standard; a consumer of the standard.

---

## 4. Layer Model

### 4.1 The three layers

| Layer | Role | Owns | Does not own |
|---|---|---|---|
| **Knowledge Layer (KB)** | Knowledge store: Content, classification, preservation, references | Content, classification (Knowledge Dimensions), Relationship, history/Records, Identity administration (Identity Registry) | process State, execution Protocol |
| **Operating Layer (OS)** | State machine & execution Protocol | State Domains (Execution, Planning, Container, Existence), ordering, concurrency, locking, Gates, Command | Content (never edits Content) |
| **Exchange Layer (EX)** | Transformational boundary: serialization, validation, import/export, mediation of external systems | exchange contracts, round-trip rules, conformance validation | Content and State (never becomes an owner) |

### 4.2 Decision: the Exchange Layer as a third layer

The Exchange Layer must exist as a separate layer because:

1. **Different invariants**: exchange has its own invariants (lossless round-trip, idempotency, referential integrity) that neither KB nor OS has.
2. **Different interaction direction**: KB and OS interact internally; EX interacts with external systems — mediation requires an explicit boundary contract.
3. **Knowledge OS vision**: future integration requires a seam defined at the standard level, not at the implementation level.

EX owns neither Content nor State — it is a boundary layer that validates and transforms representations.

### 4.3 Layer independence

- KB may change its taxonomy without changing OS Protocol (P5).
- OS may add protocol variants without changing KB classification.
- EX may add serialization formats without changing KB/OS.
- All three are bound by: **Identity** (Section 6), **global invariants** (Section 5.4), and **inter-layer contracts** (Section 5.3).

---

## 5. Layer Contracts

### 5.1 Knowledge Layer — contract

- **Responsibilities**: storing Content; classifying Artifacts into Knowledge Dimensions; maintaining cross-references (referential integrity); preserving history (P8, P12); providing retrieval semantics; administering Identity (Identity Registry).
- **Ownership**: Content; classification; Relationship records; history; Identity.
- **Allowed interactions**: serving Content reads to OS and EX; accepting State writes from OS (as metadata on the Artifact); accepting new Artifacts from OS through the creation protocol (Identity issued by the Identity Registry); executing Content State transitions through the approval Gate.
- **Invariants**: Approved/Immutable Content is not mutated (P8); references are always valid (no dangling references); classification changes without changing Identity (P15); Identity is immutable (P3).
- **Synchronization boundaries**: only Content State is owned here. Other State Domains are only reflected — KB reads State for queries, it does not write it.
- **Extension points**: new Knowledge Dimension; new Artifact type; new Relationship type (Section 14).

### 5.2 Operating Layer — contract

- **Responsibilities**: running the execution Protocol (ordering, State transitions, Gates); managing concurrency (exactly-one-active); managing locking/immutability (lock before consumption); defining Commands; providing what-next discoverability.
- **Ownership**: Execution State, Planning State, Container State, Existence State transitions; Projection Semantics (container view, ticket view).
- **Allowed interactions**: reading Content from KB (execution: Command reads instructions; Execution Container is generated from the plan); writing State to Artifacts (single-writer per field, P6); creating new Artifacts (Execution Container, ticket, session) through the creation protocol with Identity from the Identity Registry; validating State Projections against the owner.
- **Invariants**: forward-only transitions (P7); one writer per State field (P6); exactly-one-active Execution Container; **lock-atomic-with-generation** — the Execution Container creation event locks the plan and creates the container atomically; every transition is recorded in the Change Log; OS never changes Content — only State.
- **Synchronization boundaries**: State Domains are owned here; State Projections are generated/validated here; Content is only read.
- **Extension points**: new protocol variant, new Command type, new Gate type (Section 14).

### 5.3 Inter-layer contracts — coupling points

| Coupling | Direction | Contract content |
|---|---|---|
| **OS reads Content** | KB → OS | Content of the executed Artifact must be Well-formed and deterministic; Identity resolvable. |
| **OS writes State** | OS → KB | State is written as metadata on the Artifact; does not change Content; recorded in the Change Log. |
| **OS produces Artifact** | OS → KB | Creation protocol: Identity is requested from the Identity Registry; new Artifact Content is valid per the taxonomy. |
| **KB keeps Content consistent with State** | KB ↔ OS | Locked plan Content must not change (the lock gates the Content channel); State Gates allow/forbid Content transitions. |
| **EX validates & transfers** | EX ↔ KB/OS | Import/export passes through the standard contract; EX validates conformance before commit; never writes State or Content directly. |

### 5.4 Global invariants

1. **Identity immutable** — Identity is not changed by State, location, classification, or Content revision (P3).
2. **Single owner per State field** — one writer; all other views are State Projections (P6).
3. **Structure as projection of State** — structural position/representation is derived from State; never an independent fact (P9).
4. **State changes only via Protocol** — no State change path exists other than Protocol transitions (P11).
5. **Two separate change channels** — the Content channel (governance) and the State channel (Protocol) never mix (P10).
6. **Approved Content immutable** — preservation (P8, P12).
7. **Round-trip lossless** — exchange loses/duplicates nothing (P13).

### 5.5 Synchronization semantics

- **Projection Refresh**: a State Projection is validated against the owner at two points — when read (on-read) and when the owner changes (event). The preference policy between the two is an open question (Section 15); the invariant "a projection never becomes a writer" is absolute.
- **Trigger**: certain events (State transitions, OS operations) trigger validation or a transition in another domain. Triggers are defined per domain pair (Section 7.5).

---

## 6. Identity Model

### 6.1 Identity composition

**Identity instance** = `(Namespace, Type, ID, InstanceVersion)`.
**Identity line** = `(Namespace, Type, ID)`.

| Property | Part of Identity? | Reason |
|---|---|---|
| **Namespace** | **Yes** | Separates Identity spaces (products, organizations, systems). Two Artifacts with the same ID in different Namespaces are two Artifacts. |
| **Artifact Type** | **Yes** | Type determines the applicable State Domains and the binding Protocol. Type is an Identity qualifier, not a classification. |
| **ID** | **Yes** | Unique token within `(Namespace, Type)`. |
| **InstanceVersion** | **Yes (instance discriminator)** | Distinguishes instances within one Line (plan v1 vs v2). See 6.3. |
| **Knowledge Dimension (classification)** | **No** | Retrieval property. Reclassification must not break references (P15). |
| **Engineering Domain (classification)** | **No** | Classification property (P15): derived or declared, never Identity; reclassification never changes Identity (§8.1). |
| **Location / organization** | **No** | Location is a projection; Identity is independent of location (P9). |
| **Storage backend** | **No** | Identity survives across storage (Section 12). |
| **State (all domains)** | **No** | State changes; Identity never does (P3). |
| **Revision (Content edit history)** | **No** | Revision tracks the evolution of the same instance's Content; changing Revision does not change Identity. |

### 6.2 Identity rules

1. Identity is established once at creation and never changed.
2. ID is unique within `(Namespace, Type)`; InstanceVersion is unique within a Line.
3. References are always by Identity — never by location, display name, or classification.
4. Two Artifacts are the same iff their Identity is the same; Relationship never changes Identity.
5. Artifact Type determines the applicable State Vector (Section 10) — the type→state binding is part of the standard, not an implementation choice.
6. Identity must be serializable **losslessly, unambiguously, and machine-parseably** in all implementations. The serialization mechanism is an implementation decision.
7. Supersession, amendment, and derivation are expressed as **Relationships between Identities**, not Identity changes.

### 6.3 Version semantics

Two meanings of version, two answers:

- **InstanceVersion** (the version pointing to a different instance) — **part of the instance Identity**. A new instance is created deliberately (example: plan v2 after v1 is locked). Instance change = instance Identity change; the Line remains.
- **Revision** (tracking the Content evolution of the same instance) — **not part of Identity**. Revision changes on every edit and must not break references.

Rule: **Line Identity never changes; Instance Identity changes only when a new instance is deliberately created; Revision never touches Identity.** Supersession is a Relationship between two Lines — not an Identity replacement.

### 6.4 Violation case study: Identity space collision

The initial implementation encoded four Artifact types (scope definition, plan, Execution Container, ticket) in one shared ID space with the same prefix, so Type could not be distinguished deterministically from its Identity representation: one representation could be read as either a scope definition or a plan. Violation analysis:

- Rules 6.2.1–6.2.2 are violated: Type is not explicit → ID is not unique per `(Namespace, Type)`.
- Rule 6.2.3 is violated conceptually: Identity is encoded through location and representation conventions, not as an Artifact property.
- Root cause: Identity is derived from structure, not from Artifact properties; and the process stage (pipeline) became part of the Identity representation.

Binding lesson for all implementations: **Identity must not be encoded in location, process stage, or representation conventions.** Identity is a first-class property established by the Identity Registry.

---

## 7. State Taxonomy

### 7.1 Candidate domain evaluation

| Candidate | Decision | Reason |
|---|---|---|
| **Artifact State (unified)** | **Reject** | A State monolith is the root cause of status duplication: one Artifact has many independent State dimensions. State is always a State Vector. |
| **Execution State** | **Accept** | Work item progress through Protocol. |
| **Lifecycle State (generic)** | **Reject as a single domain** | "Lifecycle" is a composition of State Domains over time; product Phase is context, not a State Domain. What survives: Existence State (domain) + Phase (context). |
| **Governance/Content State** | **Accept** (named Content State) | Content maturity through the approval gate. |
| **Review State** | **Reject as a domain** | Review is a Gate on transitions (the Review stage in Content State; the "In Review" value in Execution State) + an Artifact type (review record). No independent State remains after that separation. |
| **Release State** | **Reject as a domain** | Release is a Phase context + a readiness Gate evaluated over the aggregate of other State Domains. |
| **Planning State** | **Accept** | Plan commitment level: Draft → Approved → Immutable. |
| **Container State** | **Accept** | Execution Container open/closed + concurrency. |
| **Existence State** | **Accept** | Artifact presence: active vs preserved. |

### 7.2 Formal domains

General domain rules: every domain has exactly one **initial state** and exactly one **terminal state** (State with no outgoing transitions); all transitions are forward-only (P7); corrections are made with a new instance + Relationship, not regression. **Domain values cover only owned State; derived conditions — for example "Completed" in a working context — are not domain values but transition-triggering conditions.**

| Domain | Values | Responsibility | Owner | Transition rules |
|---|---|---|---|---|
| **Content State** | Draft → Review → Approved; post-approval terminal: Amended \| Superseded | Content maturity as knowledge; governance channel | Knowledge Layer (gate: owner/approver) | Forward-only; approval gate; changes after Approved only via amendment/supersession (P8). **Variants**: standard (terminal Amended), decision record (Approval stage named Accepted; supersession optional), ADR (Approval stage named Accepted; supersession must point to a successor). Variants may rename stages and eliminate optional stages, but must preserve the semantic position: pre-approval → approval gate → post-approval terminal. |
| **Execution State** | Planned → Todo → In Progress → In Review → Done | Work item progress through Protocol | Operating Layer (single-writer: work item artifact) | Strictly sequential; never skip; never revert; one initial; one terminal; every transition recorded in the Change Log. |
| **Planning State** | Draft → Approved → Immutable | Plan commitment: tentative → committed → locked | Operating Layer | Forward-only; Approved = ready for execution; **Immutable is reached atomically with the Execution Container creation event** (lock-atomic-with-generation); post-lock changes = new instance (InstanceVersion). |
| **Container State** | Active → Completed | Execution Container open/closed; concurrency | Operating Layer | **Completed is a derived transition**: triggered by the Execution State aggregate (all work items Done); exactly-one-Active (mutual exclusion). |
| **Existence State** | Active → Archived → Retired | Artifact presence in active process vs preservation | Operating Layer (transitions); Knowledge Layer (preservation principle, P12) | Forward-only; **Archived** = reference-only (no other State transitions, no Content mutation, still available in retrieval); **Retired** = terminal preservation (not surfaced in normal retrieval, Content immutable, Identity unchanged). Applies to all Artifact types. |

### 7.3 Mapping of initial implementation status machines → domains

| Initial implementation status machine | Domain |
|---|---|
| Living documents Draft→Review→Approved→Amended | Content State (standard variant) |
| ADR Proposed→Accepted→Superseded | Content State (ADR variant) |
| Decisions Draft→Accepted→Superseded (optional) | Content State (decision variant) |
| Roadmap Draft→Approved→Immutable | Planning State |
| Sprints Active→Completed | Container State |
| Sessions Active→Completed→Archived | Existence State (Completed = derived condition of the referenced work item) |
| Work items Planned→Todo→In Progress→In Review→Done | Execution State |

Empirical evidence of domain independence: the initial implementation needed seven different status machines because one machine cannot express seven semantics. The State Taxonomy organizes them into domains with explicit rules.

### 7.4 State Vector

Every Artifact carries a **State Vector** = the tuple of State Domains it **owns** per its type; domains that do not apply are marked not-applicable. Examples: work item = `(Execution State, Existence State)` — Content State does not apply; plan = `(Content State, Planning State, Existence State)`; Execution Container = `(Container State, Existence State)`; ADR = `(Content State, Existence State)`.

**An Artifact whose entire State is projected has an empty State Vector** (example: ticket = `(∅)`; its State is a State Projection over the referenced work item). This formally resolves status duplication: status on derived representations (container view, ticket view) is a State Projection of the owner, validated through Projection Refresh — replacing the writer-less "must always agree" invariant.

### 7.5 Interactions between domains

| Interaction | Source → Target | Semantics |
|---|---|---|
| Execution Container creation (generation event) | Operating Layer operation → Planning State | The container creation event from a plan triggers the plan's transition to Immutable; atomic with the creation (lock-atomic-with-generation). |
| Work item State aggregate | Execution State → Container State | All work items in the container Done triggers container Completed (derived transition). |
| Plan lock → Content channel | Planning State → Content State | An Immutable plan gates Content changes on that instance (governance channel locked). |
| Content readiness → commitment | Content State → Planning State | The approval Gate on a plan requires Content to be ready (Approved in Planning State composes Content maturity). |
| Container Completed → preservation | Container State → Existence State | A finished container may be archived. |
| Release readiness Gate | Execution + Planning + Container (+ review Gate, approval Gate) → **phase change** | Aggregate State evaluation to authorize a Phase change (context update; not a state transition). See 11.2. |

### 7.6 Independence vs unification — the decision

State Domains **remain independent** (not unified). Reasons: (1) empirical evidence from the initial implementation — different semantics require different machines; (2) different owners (KB vs OS); (3) different rates of change (Content changes slowly via governance; execution State changes fast via Protocol); (4) unification reintroduces the status duplication problem. Interactions are managed explicitly via Triggers and Gates (7.5), not by merging domains.

---

## 8. Knowledge Taxonomy

Knowledge Layer classification dimensions. One Artifact has **one primary dimension** + optional secondary dimension(s). Classification is a retrieval property — reclassification does not change Identity (P15).

| Dimension | Contents | Stability | Notes |
|---|---|---|---|
| **Product Intent** | Vision, strategy, principles, non-negotiables | Very stable | Strategy is a dimension not previously represented explicitly. |
| **Requirements** | Requirement documents — "what must be built and why" — with their amendments | Stable (evolves via governance) | |
| **Architecture** | System descriptions, domain models, component boundaries | Stable | |
| **Decisions** | Irreversible decisions (ADR) and reversible ones (decision log); including product decisions | Stable, cumulative | |
| **Specifications** | Functional and non-functional (NFR) specifications, API, data | Stable when Approved | Dimension not previously represented explicitly; distinct from Vocabulary. |
| **Standards & Guidelines** | Engineering standards, conventions, definition of done | Stable, high governance | Previously mixed with operational knowledge. |
| **Operational Knowledge** | Runbooks, deployment, migration, checklists | Medium (changes per environment) | Separate from Standards & Guidelines. |
| **Governance & Quality** | Review findings, audits, quality gates | Cumulative | Quality was previously only a role, not a dimension. |
| **Planning Knowledge** | Plan Content (roadmap, milestone definition) and artifact relationships (traceability) | Medium (commitment) | Traceability is a Relationship artifact. |
| **Records** | Release record, change log, historical snapshots | Immutable | Release record was not previously represented. |
| **Research** | Investigation findings, technical research results | Cumulative | Distillation path to durable dimensions mandatory. |
| **Vocabulary** | Glossary, canonical terms, lifecycle model | Very stable | Not Specifications — this separation must be maintained. |

### 8.1 Engineering Domains and Knowledge Stratification

**Engineering Domain.** An **Engineering Domain** is the home classification of an Artifact: the stratum-aligned category of engineering knowledge to which the Artifact belongs. Engineering Domain is a **classification property** (P15 extension — reclassification never touches Identity): it is **derived** from the Artifact's Knowledge Dimension and token family, or **declared** by an extension; it is **never Identity** (Section 6.1) and never part of the State Vector.

The five canonical Engineering Domains, in stratum order (stratum 1 highest → 5):

| Engineering Domain | Stratum | Token families | Knowledge Dimensions |
|---|---|---|---|
| **Discovery** | 1 (highest) | vis-, str-, req-, fnd- | intent, requirements, research |
| **Architecture** | 2 | arc-, adr-, dec-, spec-, std-, gls- | architecture, decisions, specifications, standards, vocabulary |
| **Planning** | 3 | scp-, epc-, plan-, trc- | planning |
| **Execution** | 4 | rvw-, ctr-, tkt-, sto-, ts-, bug-, td-, ch-, spk-, ses- | quality + operating tokens (informational dimension) |
| **Operations** | 5 | run-, rel- | operations, records |

Full token mapping (26 tokens):

| Token | Knowledge Dimension | Engineering Domain | Stratum |
|---|---|---|---|
| vis- | intent | Discovery | 1 |
| str- | intent | Discovery | 1 |
| req- | requirements | Discovery | 1 |
| fnd- | research | Discovery | 1 |
| arc- | architecture | Architecture | 2 |
| adr- | decisions | Architecture | 2 |
| dec- | decisions | Architecture | 2 |
| spec- | specifications | Architecture | 2 |
| std- | standards | Architecture | 2 |
| gls- | vocabulary | Architecture | 2 |
| scp- | planning | Planning | 3 |
| epc- | planning | Planning | 3 |
| plan- | planning | Planning | 3 |
| trc- | planning | Planning | 3 |
| rvw- | quality | Execution | 4 |
| ctr- | — (operating token; informational dimension) | Execution | 4 |
| tkt- | — (operating token; informational dimension) | Execution | 4 |
| sto- | — (operating token; informational dimension) | Execution | 4 |
| ts- | — (operating token; informational dimension) | Execution | 4 |
| bug- | — (operating token; informational dimension) | Execution | 4 |
| td- | — (operating token; informational dimension) | Execution | 4 |
| ch- | — (operating token; informational dimension) | Execution | 4 |
| spk- | — (operating token; informational dimension) | Execution | 4 |
| ses- | — (operating token; informational dimension) | Execution | 4 |
| run- | operations | Operations | 5 |
| rel- | records | Operations | 5 |

**Knowledge Stratum.** A **Knowledge Stratum** is the authority level of an Engineering Domain: a fixed position in the strict linear order **Discovery → Architecture → Planning → Execution → Operations** (stratum 1 highest authority → 5). The stratum is always **derived** from the Engineering Domain — never declared by an Artifact, never part of the State Vector. One Artifact has exactly one Engineering Domain and therefore exactly one stratum.

**Stratum Authority Invariant.** Engineering Domains form a strict linear order by Knowledge Stratum. Knowledge in a lower stratum must not contradict knowledge in a higher stratum that is in force (Content State Approved or beyond, or Planning State Immutable); where a contradiction is discovered, it is resolved by changing the lower-stratum knowledge through the governance channel (new instance + Relationship, forward-only), never by overriding the higher stratum. A lower-stratum artifact must not supersede or amend a higher-stratum artifact.

**Stratum vs lifecycle vs Phase.** Stratum, lifecycle, and Phase are three distinct axes. **Stratum** is an authority property: static per Engineering Domain, fixed at classification. **Lifecycle** is the movement of State over time (Section 7): the same Artifact moves through its State Domains while its stratum does not move. **Phase** is product context (Sections 3, 11.2): a context attribute on planning/scope Artifacts. One Artifact has a fixed stratum while its lifecycle moves and its Phase context may change.

**Methodology independence.** PRD, ADR, RFC, Epic, Initiative, Sprint, Iteration, Ticket, Release, Incident, and Runbook are **Representation Aliases**: methodology terms mapped onto a canonical token + Engineering Domain. They are never frontmatter values and never Artifact types in their own right — the canonical token and Artifact type govern Identity and State.

| Representation Alias | Canonical token | Engineering Domain |
|---|---|---|
| PRD | req- | Discovery |
| ADR / RFC | adr- | Architecture |
| Epic | epc- | Planning |
| Initiative | scp- | Planning |
| Sprint / Iteration | ctr- | Execution |
| Ticket | tkt- | Execution |
| Release | rel- | Operations |
| Incident | bug- | Execution |
| Runbook | run- | Operations |

Methodologies (Scrum, Kanban, Shape Up, and similar) are **convention layers over EKA**: they may map onto tokens and domains through Representation Aliases, but they are never part of the Core standard.

### 8.2 Engineering Domain Registry

**Engineering Domain Registry.** The **Engineering Domain Registry** is the canonical vocabulary of Engineering Domains: the five domains, their stratum positions, and the knowledge they produce and consume. The Registry is normative: classification, stratum derivation, and the Conformance Rules R10–R12 operate on this vocabulary. The token→domain and dimension→domain mappings remain in §8.1; the Registry defines domain semantics, not the mapping. **Produced Knowledge** is the knowledge an Engineering Domain produces: the Artifact types homed in the domain. **Consumed Knowledge** is the knowledge a domain consumes from other strata: the Artifact types of higher-stratum domains that its Artifacts derive from or depend on.

| Engineering Domain | Stratum | Produced | Consumed | Constrains | Aliases |
|---|---|---|---|---|---|
| **Discovery** | 1 | Vision, Strategy, Requirement, research findings | — (nothing above stratum 1) | Architecture, Planning, Execution, Operations | PRD (req-) |
| **Architecture** | 2 | Architecture Description, ADR, Decision Record, Specification, Standard, Glossary Term | Requirement, research, intent (Discovery) | Planning, Execution, Operations | ADR / RFC (adr-) |
| **Planning** | 3 | Scope Definition, Epic, Plan, Traceability Artifact | Requirement, intent (Discovery); Decision Record, Specification, Standard (Architecture) | Execution, Operations | Epic (epc-), Initiative (scp-) |
| **Execution** | 4 | Review, Execution Container, Ticket, Work Item, Session | Plan, Scope (Planning); Specification, Standard, Decision Record (Architecture); Requirement (Discovery) | Operations | Sprint / Iteration (ctr-), Ticket (tkt-), Incident (bug-) |
| **Operations** | 5 | Runbook, Release Record | execution aggregate, review/session findings (Execution); Standard (Architecture) | — (lowest authority) | Release (rel-), Runbook (run-) |

**Discovery — Stratum 1 (highest authority).**

- **Purpose:** home of intent, requirements, and research knowledge.
- **Responsibilities:** owns intent, requirements, and research knowledge: Vision/Manifesto, Strategy, Requirement (PRD), research findings (fnd-).
- **Produced Knowledge:** Vision/Manifesto (vis-), Strategy (str-), Requirement/PRD (req-), research findings (fnd-).
- **Consumed Knowledge:** none from other strata — nothing ranks above stratum 1; consumption is within-domain (research feeds requirements).
- **Relationships:** root of the authority chain (derives from nothing); constrains Architecture, Planning, Execution, Operations (representation alias: PRD).
- **Knowledge Stratum:** 1.
- **Constraints:** must not be contradicted by lower-stratum knowledge in force (§8.1); must not be superseded or amended by lower-stratum Artifacts (R12).

**Architecture — Stratum 2.**

- **Purpose:** home of decisions, specifications, standards, and vocabulary — the durable design knowledge that binds planning and execution.
- **Responsibilities:** owns architecture, decisions, specifications, standards, and vocabulary knowledge: Architecture Description, ADR, Decision Record, Specification, Standard/Guideline, Glossary/Term.
- **Produced Knowledge:** Architecture Description (arc-), ADR (adr-), Decision Record (dec-), Specification (spec-), Standard/Guideline (std-), Glossary/Term (gls-).
- **Consumed Knowledge:** Discovery — Requirement (req-), research findings (fnd-), intent (vis-, str-): decisions and specifications derive from requirements.
- **Relationships:** derives from Discovery; constrains Planning, Execution, Operations (representation alias: ADR / RFC).
- **Knowledge Stratum:** 2.
- **Constraints:** must not contradict Discovery knowledge in force; must not be redefined by lower strata (Execution must not redefine Architecture); must not be superseded or amended by lower-stratum Artifacts (R12).

**Planning — Stratum 3.**

- **Purpose:** home of commitment knowledge — scope, epics, and plans that convert requirements and decisions into an executable commitment.
- **Responsibilities:** owns planning knowledge: Scope Definition, Epic, Plan (roadmap), Traceability/Relationship Artifact.
- **Produced Knowledge:** Scope Definition (scp-), Epic (epc-), Plan (plan-), Traceability Artifact (trc-).
- **Consumed Knowledge:** Discovery — Requirement (req-), intent (vis-, str-); Architecture — Decision Record/ADR (dec-, adr-), Specification (spec-), Standard (std-): scope derives from requirements; plans honor approved decisions and specifications.
- **Relationships:** derives from Discovery and Architecture; constrains Execution, Operations (representation aliases: Epic, Initiative).
- **Knowledge Stratum:** 3.
- **Constraints:** must not contradict Discovery or Architecture knowledge in force; must not supersede or amend higher-stratum Artifacts (R12).

**Execution — Stratum 4.**

- **Purpose:** home of execution knowledge — work items, containers, tickets, reviews, and sessions that carry work through Protocol.
- **Responsibilities:** owns execution knowledge and the operating tokens: Review, Execution Container, Ticket, Work Item (story, technical story, bug, tech debt, chore, spike), Session.
- **Produced Knowledge:** Review (rvw-), Execution Container (ctr-), Ticket (tkt-), Work Items (sto-, ts-, bug-, td-, ch-, spk-), Session (ses-).
- **Consumed Knowledge:** Planning — Plan (plan-), Scope (scp-); Architecture — Specification (spec-), Standard (std-), Decision Record (dec-, adr-); Discovery — Requirement (req-): execution follows the plan and the specifications it must satisfy.
- **Relationships:** derives from Planning (transitively Architecture, Discovery); constrains Operations (representation aliases: Sprint/Iteration, Ticket, Incident).
- **Knowledge Stratum:** 4.
- **Constraints:** must not redefine Architecture; must not contradict a locked plan (Planning State Immutable) or higher-stratum knowledge in force; must not supersede or amend higher-stratum Artifacts (R12).

**Operations — Stratum 5 (lowest authority).**

- **Purpose:** home of operational knowledge and records — runbooks and release records that preserve what ran and how it runs.
- **Responsibilities:** owns operations and records knowledge: Runbook/Operational Guide, Release Record.
- **Produced Knowledge:** Runbook/Operational Guide (run-), Release Record (rel-).
- **Consumed Knowledge:** Execution — execution aggregate and review/session findings (ctr-, tkt-, sto-, ts-, bug-, td-, ch-, spk-, rvw-, ses-): the Release Record carries the execution aggregate; Architecture — Standard (std-): runbooks comply with standards.
- **Relationships:** derives from Execution (transitively all higher strata); constrains nothing (representation aliases: Release, Runbook).
- **Knowledge Stratum:** 5.
- **Constraints:** must not contradict any higher-stratum knowledge in force; must not supersede or amend higher-stratum Artifacts (R12); must not rewrite executed history (release records preserve the execution aggregate, P12).

**Registry governance.** The Engineering Domain Registry is canonical vocabulary. Refinement — domain addition, stratum reassignment, or boundary change — is taxonomy governance (§14.2): proposal → review → acceptance; a refinement must not weaken the Stratum Authority Invariant (§16.3). Representation Aliases are **not** part of the Registry: they live in the **Representation Alias Registry**, a separate convention-layer document mapping methodology terms onto canonical tokens + Engineering Domain (§8.1 table).

**Refinement decision — Vocabulary → Architecture (kept from v1.1).** The Vocabulary dimension (gls-) is homed in the Architecture domain. Rationale: vocabulary is very-stable, high-governance knowledge whose authority class matches Architecture. Alternative considered and rejected: a dedicated Vocabulary domain — a sixth stratum with no authority gap.

### 8.3 Knowledge Stratification Governance

**Terminology.** **Knowledge Stratum** is the canonical term for the authority level of an Engineering Domain. The term was chosen over "Knowledge Layer" and "Knowledge Level": (1) collision avoidance — "Knowledge Layer" already names the architectural layer (Section 4), and stratification must never be confused with the Layer model; (2) no ranking connotation — "level" implies a value or competence ranking; the stratum order is an authority order, not a quality ranking; (3) the geological metaphor — strata are horizontal bands that rest on and are constrained by the bands above them, exactly the authority semantics of the Stratum Authority Invariant. The architectural Layer (KB/OS/EX) remains the only "Layer" usage in the standard; "Knowledge Layer" or "Knowledge Level" in a stratification context is non-canonical.

**Meaning.** A Knowledge Stratum is the authority level of an Engineering Domain: a fixed position in the strict linear order Discovery → Architecture → Planning → Execution → Operations (stratum 1 highest → 5). The stratum is always derived from the Engineering Domain — never declared by an Artifact, never part of Identity (Section 6.1), never part of the State Vector; reclassification never changes Identity (P15).

**Governance.** The Stratum Authority Invariant (§8.1) is a taxonomy invariant — not one of the global invariants (Section 5.4). It can only be strengthened, never weakened, by extensions (§16.3; §14.2.1). Any change to the domain order or to a stratum assignment is taxonomy governance (§14.2): proposal → review → acceptance, registered as part of the standard.

**Validation implications.** The Conformance Rules set R0–R12 (Naming and Terminology Specification v1.1 §9.3) operationalize stratification:

- **R10 — stratification traceability (warning):** every Artifact not in stratum 1 must have a resolvable derives-from/depends-on chain (direct or transitive) reaching a strictly higher stratum. Exempt: tkt-/ses- tokens and draft knowledge Artifacts (work items own no Content State and are never exempt via the draft clause). Stratification is a structural quality signal, never a commit blocker.
- **R11 — domain coherence (blocking):** a declared Engineering Domain must be one of the five canonical domains and equal the Type's home domain; absent = derived, no check.
- **R12 — cross-stratum supersession prohibition (blocking):** supersedes/amends must never target an Artifact in a strictly higher stratum; durable content moves down the authority chain, never up.

**Relationship to Engineering Domains.** Stratification maps 1:1 onto Engineering Domains: stratum(D) is the derived function of the domain order — the domain order **is** the stratum order (§8.2). One Artifact has exactly one Engineering Domain and therefore exactly one stratum; no Artifact declares a stratum.

---

## 9. Execution Taxonomy

Classification of Operating Layer Protocol elements:

| Element | Definition | Responsibility | Invariant |
|---|---|---|---|
| **Ordering (chain)** | Stage relationship order: requirement → scope → capability → plan → container → work item → working context → validation | Answers "what next" deterministically; Protocol, not a location property | Order is explicitly defined; execution follows the order |
| **State Transitions** | Rules for moving values within a State Domain | Runs the Protocol per domain (Section 7.2) | Forward-only; never skip; never revert; recorded in the Change Log |
| **Concurrency Control** | Mutual exclusion locking over the Execution Container | Exactly-one-active container | One active container; subsequent creation waits |
| **Versioning / Immutability** | Lock the plan when execution starts; changes = new instance | Plan vs execution consistency | Lock-atomic-with-generation; locked Content does not change |
| **Gates** | Conditions before transition/execution: approval gate, readiness gate, review gate | Controls when a transition is legal | Gates are evaluated over owner State, not projections |
| **Commands** | Deterministic execution instructions consumed by executors | Translates Content Artifacts into execution actions | Command Content Well-formed; deterministic output |
| **Agent Coordination** | Semantics of agent–system interaction | Identity parseable, State explicit, Content deterministic, strict ordering | Agents read State/Identity without ambiguity |
| **Execution Containers** | Execution window wrapping work items | Aggregation, concurrency, snapshot | Container State derived from work items |
| **Projection Semantics** | Computation rules of State Projections (container view, ticket view): from which owner, how aggregation works | Provides views without adding writers | A projection never becomes a writer (P6, P9) |

---

## 10. Artifact Taxonomy

Conceptual Artifact type system — not a storage design. Every type defines: Knowledge Dimension, the State Domains it **owns**, and Identity/Relationship notes.

| Artifact type | Knowledge Dimension | Owned State Domains | Identity & relationship notes |
|---|---|---|---|
| Vision / Manifesto | Product Intent | Content, Existence | Single Line; amendments rare |
| Strategy | Product Intent | Content, Existence | New type (previously not represented) |
| Requirement (PRD) | Requirements | Content, Existence | Line + amendments as instance/Relationship |
| Scope Definition | Planning Knowledge + Requirements | Content, Existence | **Phase context** (Discovery/MVP/Milestone/Release) as attribute — not category |
| Epic | Planning Knowledge | Content, Existence | derives-from Scope relationship |
| Plan (roadmap) | Planning Knowledge | Content, Planning, Existence | InstanceVersion significant (v1, v2, ...); lock via generation |
| Execution Container (sprint) | — (Content = projection snapshot) | Container, Existence | Content derived; created by OS |
| Ticket | — (Content = execution instruction / Command) | **∅ (none; State is a State Projection over the referenced work item)** | Execution view: projection, not writer (P6) |
| Work Item (story, technical story, bug, tech debt, chore, spike) | Requirements / Records / Research | Execution, Existence | Single-writer Execution State; Content may be distilled into knowledge dimensions |
| Session | — (ephemeral by design) | Existence | Content ephemeral; Completed = derived condition; Distillation mandatory before Archived |
| Review | Governance & Quality | Content, Existence | Gate semantics; findings → Decisions |
| ADR | Decisions | Content (ADR variant), Existence | Supersession = Relationship to another Line |
| Decision Record | Decisions | Content (decision variant), Existence | Supersession optional |
| Architecture Description | Architecture | Content, Existence | |
| Specification | Specifications | Content, Existence | New type (previously not represented) |
| Standard / Guideline | Standards & Guidelines | Content, Existence | New type (previously not represented) |
| Runbook / Operational Guide | Operational Knowledge | Content, Existence | |
| Release Record | Records | Content, Existence | New type (previously not represented); carries execution aggregate + release gate |
| Glossary / Term | Vocabulary | Content, Existence | |
| Traceability / Relationship Artifact | Planning Knowledge | Content, Existence | Content = set of Relationships by Identity |

Rule: Type determines the State Vector (the type→state binding is part of the standard); new types are extensions (Section 14) with the obligation to declare the complete owned State Vector.

---

## 11. Conceptual Lifecycle

### 11.1 Two change channels over time

- **Content evolution**: Draft → Review → Approved → Amended/Superseded (Content State) — Content matures and is preserved; Approved Content is immutable (P8).
- **Execution progress**: Planned → … → Done (Execution State) — work items move forward; corrections via a new instance, not revert.

Both are orthogonal: an Artifact may have Approved Content and In Progress execution State simultaneously; the type determines which domains apply (Section 10).

### 11.2 Phase as context, not category

Discovery, MVP, Milestone, Release are **Phase context** on planning/scope Artifacts — an attribute, not a category and not a State Domain. Rules:

- Phase attaches to Scope Definition / Plan as a context attribute.
- **Phase change** is a context update authorized by a **readiness Gate**, evaluated over the State aggregate: **release-ready** = (all work items in scope Done) ∧ (all Execution Containers Completed) ∧ (plan locked / Immutable) ∧ (review gate passed) ∧ (Content approval gate passed).
- Phase never becomes part of Identity (P3) — a Scope Definition remains the same identity when the product changes phase; what changes is its context attribute.

### 11.3 Product lifecycle vs Artifact lifecycle

- **Artifact lifecycle**: the movement of State Domains per Artifact (Section 7) — Content matures, execution completes, plans are locked, Artifacts are archived or retired.
- **Product lifecycle**: the sequence of Phase contexts on the Scope Artifact — Discovery → MVP → Growth → Maturity → Sunset. Each phase produces its own Artifacts (new scope, new plan, new Release Record) with new Identity; the old phase remains as a Record (P12).

### 11.4 Distillation lifecycle

Ephemeral → durable: Session (Existence) and Review produce decisions/ADRs (Content State) through the **mandatory** Distillation path before Archived. Preservation: Superseded/Archived/Retired Artifacts remain as Records with immutable Content and intact supersession Relationships.

---

## 12. Storage Independence Model

### 12.1 Thought experiment

Assume the current storage medium disappears; knowledge is stored in a relational database, graph database, object store, Atrium, or future platform. What is tested: the Identity Model, State Taxonomy, Knowledge Taxonomy, layer contracts.

### 12.2 What survives (part of the standard)

| Concept | Why it survives |
|---|---|
| **Identity Model** | Identity = conceptual property `(Namespace, Type, ID[, InstanceVersion])`. In a relational database it becomes a key; in a graph, a node property; in an object store, a key. References by Identity hold in all of them. |
| **State Taxonomy** | State Domains are semantics, not storage. State Projection = view (relational), query (graph), computed (object store). Single-writer remains enforceable. |
| **Knowledge Taxonomy** | Classification is an Artifact property; any backend can index it. |
| **Layer Contracts & Global Invariants** | Contracts define behavior, not storage. All 7 invariants (5.4) remain in force. |
| **Exchange Contracts** | Round-trip, idempotency, referential integrity — independent of the medium. |

### 12.3 What belongs to the implementation

- Serialization formats, storage schemas, index structures.
- **Physical addressing** (path, key, URL) — the standard only requires: every Identity **resolvable** to exactly one Artifact within one system.
- **Retrieval**: the standard requires that declared query semantics be implementable; the query language is an implementation matter.
- **Enforcement capability**: structural constraints on file-based implementations; constraints in relational databases; a validation layer on graphs. **Invariant requirements are identical; enforcement mechanisms vary** (P16). The initial implementation's value of "structure as a state machine that is always in sync" is an implementation capability; the standard inherits its requirement (determinism, P11), not its mechanism.

### 12.4 The standard's stance on serialization

The standard establishes the contract (what must be preserved, round-trip rules, validation); the implementation establishes the format. Identity serialization must be canonical and unambiguous in all implementations (rule 6.2.6).

---

## 13. Import / Export Model

### 13.1 What must be preserved in exchange

| Element | Requirement |
|---|---|
| **Identity** | Namespace, Type, ID, InstanceVersion — intact; without duplication; canonical. |
| **State** | The complete State Vector (all owned domains) with exact values + transition history (Change Log). |
| **Content** | Complete Content, Well-formed per its type. |
| **Relationships** | All relationships by Identity (supersedes, amends, derives-from, depends-on, validates) — referential integrity across systems. |
| **Classification** | Knowledge Dimension assignment (primary + secondary). |
| **History** | Supersession/amendment links, Change Log, preservation status (Archived/Retired). |

### 13.2 Round-trip requirements

1. **Lossless**: no loss or duplication of Identity/State/Content/Relationship.
2. **Idempotent**: re-import = no-op (or a declared clean replace) — never duplicates Artifacts.
3. **Referential integrity**: no dangling references after import; cross-system references are resolved or explicitly rejected.
4. **Identity conflict policy**: importing an Identity that already exists = **reject or explicit re-namespace** — never a silent merge.
5. **Validation before commit**: import validates conformance to the standard (Identity unique, State valid, Content Well-formed) before writing.
6. **Schema versioning**: the exchange contract itself is versioned; import/export declares the contract version it complies with.

### 13.3 The serialization format contract (not the format)

Any format must: encode Identity canonically; represent the complete State Vector; express Relationships by Identity; be mechanically validatable against the standard; be readable by the Exchange Layer without ambiguous interpretation. The format itself is an implementation decision.

---

## 14. Extension Model

### 14.1 Extension points

| Point | Weight | Mechanism |
|---|---|---|
| New Artifact type | Light | Type definition: Knowledge Dimension + **complete owned State Vector** + Identity rules; registered in the taxonomy. |
| New Knowledge Dimension | Light–medium | New classification axis; must not break existing classifications (P15). |
| New Relationship type | Light | New relationship semantics between Identities. |
| New Protocol variant / Command / Gate | Medium | Protocol variation within existing invariants. |
| New State Domain | Heavy | Only if the State semantics are not covered by existing domains; full definition required (owner, rules, interactions). |
| New Phase vocabulary | Light | New Phase context values. |

### 14.2 Extension rules

1. Extensions **must not weaken invariants** (5.4).
2. **Backward compatibility**: all existing Artifacts remain valid under the extension.
3. **Core closed, taxonomy open**: Identity, layer contracts, and invariants are the core, closed to extension (changes = standard revision); taxonomies (types, dimensions, domains, protocol) are open under governance.
4. Extensions must be **exchangeable** (covered by schema versioning, Section 13).
5. Extension governance: proposal → review → acceptance (registered as part of the standard). This closes the extension-without-principles gap identified in the initial implementation.
6. New Artifact types **must declare the complete owned State Vector** — no implicit default inheritance.
7. New Artifact types and new Knowledge Dimensions **must declare their home Engineering Domain** (Section 8.1). Taxonomy extension governance covers Engineering Domains: the taxonomy is open to domain extension (item 3), and a new Engineering Domain follows the same governance — proposal → review → acceptance (item 5).

---

## 15. Open Questions

**Resolved during ratification:** (1) release is modeled as a Phase context + readiness Gate, not a State Domain; (2) the Exchange Layer is an architectural layer, not a cross-cutting concern; (3) new Artifact types must declare the complete owned State Vector (elevated to rule 14.2.6).

Questions that remain open, each with its own trade-offs:

1. **Planning State vs Content State (unification)** — Current decision: separate (commitment ≠ maturity). Trade-off: separation adds a domain; unification simplifies but mixes "mature Content" with "locked plan", which have different Triggers.
2. **Line semantics on supersession** — Is supersession two Lines with a Relationship (current decision) or one Line with instances? The current decision simplifies per-decision historiography; the alternative simplifies tracing decision chains but obscures that the new decision is a different Artifact.
3. **Depth of the query semantics contract** — How precisely must retrieval semantics be defined in the standard (for consistency across implementations) vs left to implementations (so innovation is not constrained)?
4. **Distributed/offline Identity generation** — Without a central Identity Registry, how do two systems generate collision-free IDs (global ID vs delegated Namespace vs registry)? Affects the exchange contract.
5. **Projection Refresh policy** — Event-driven (validate on every transition) vs on-read (validate when read): the trade-off of real-time consistency vs cost. The invariant "a projection is not a writer" is unchanged; the refresh mechanism is not yet locked.
6. **Depth of the Content structure contract (well-formedness)** — How strict must Content structure be per Artifact type (for machine-parseability) vs flexible (for expressiveness)? Too strict constrains; too loose breaks Command determinism.
7. **Multi-dimensional classification** — Current rule: one primary dimension + optional secondary. Conflicts between dimensions and the obligation of the secondary dimension are not fully specified.

---

## 16. Future Evolution

### 16.1 Conceptual milestones

1. **Ratification of this standard** as the canonical engineering knowledge model (contracts, invariants, base taxonomies).
2. **Exchange contract v1**: definition of round-trip, idempotency, referential integrity, schema versioning + conformance suite (conformance validator).
3. **Reference implementations**: (a) conforming the initial implementation to the standard; (b) a relational-database-based implementation; (c) a graph-based implementation — proving storage independence (Section 12).
4. **Knowledge OS integration**: the Exchange Layer becomes the seam — the Knowledge OS imports/exports Artifacts with intact Identity, State, and Relationship; knowledge becomes queryable and operable by external systems.
5. **Ecosystem**: validator, importer, agent tooling that comply with the contract; extensions registered via Section 14 governance.

### 16.2 The initial implementation's role going forward

The initial implementation changes status: from "the architecture itself" to **one reference serialization** that (a) demonstrates conformance to the standard, (b) serves as the onboarding baseline, (c) continues to function as an Engineering Operating System for projects that choose that medium. The standard becomes the canon; the implementation becomes the example.

### 16.3 Evolution invariants

Evolution of the standard never changes: Identity (P3), global invariants (5.4), the two-channel principle (P10), and layer composition (KB + OS + EX). What may evolve: taxonomies (dimensions, types, domains, protocol) through extension governance. Foundations that no iteration can negotiate: **knowledge base and operating system as two layers of one system, bound by Identity, State owned by the Operating Layer, pipeline as a first-class Protocol.**

The **Stratum Authority Invariant** (Section 8.1) is a **taxonomy invariant** — not one of the global invariants (Section 5.4); it can only be **strengthened, never weakened**, by extensions (Section 14).

---

*End of Canonical Specification — EKA v1.1 (Ratified).*
