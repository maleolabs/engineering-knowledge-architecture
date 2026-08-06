# Glossary — EKA v1.1

Alphabetical glossary of all capitalized canonical terms from the EKA v1.1 specification. Definitions are reproduced from the canonical text (`standard/eka-specification-v1.1.md`) — not paraphrased. Section references are provided for navigation.

## A

### Artifact
an engineering knowledge entity that has Identity, Content, a State Vector (the State domains it **owns**), and Relationship. The basic unit of the model. *(Section 3)*

### Artifact Instance
one version of a Line's existence: Line + `InstanceVersion`. *(Section 3)*

### Artifact Line
the enduring Identity entity: one `(Namespace, Type, ID)`. *(Section 3)*

## C

### Change Log
the chronological record of State transitions on an Artifact: domain, old value, new value, time, authority. Mandatory for all State Domains. *(Section 3)*

### Command
a deterministic execution instruction consumed by an executor (human or agent). Command Content is Content (owned by the Knowledge Layer); its execution is governed by Protocol (Operating Layer). *(Section 3)*

### Consumed Knowledge
the knowledge a domain consumes from other strata: the Artifact types of higher-stratum domains that its Artifacts derive from or depend on. *(Section 8.2)*

### Container State
the Execution Container's State Domain: values `Active → Completed`; responsibility "Execution Container open/closed; concurrency"; owner Operating Layer. Transition rules: "Completed is a derived transition: triggered by the Execution State aggregate (all work items Done); exactly-one-Active (mutual exclusion)." *(Section 7.2)*

### Content
the semantic payload of an Artifact: intent, decisions, design, constraints, procedures, notes. Belongs to the Knowledge Layer. *(Section 3)*

### Content State
Content maturity State Domain: values `Draft → Review → Approved`; post-approval terminal: `Amended | Superseded`; responsibility "Content maturity as knowledge; governance channel"; owner "Knowledge Layer (gate: owner/approver)". Transition rules: "Forward-only; approval gate; changes after Approved only via amendment/supersession (P8)." *(Section 7.2)*

## D

### Distillation
transformation of ephemeral knowledge (working context, review findings) into durable knowledge (decisions, ADRs, Records). *(Section 3)*

## E

### Engineering Domain
the home classification of an Artifact: the stratum-aligned category of engineering knowledge to which the Artifact belongs. Classification property (P15 — reclassification never touches Identity); derived from the Artifact's Knowledge Dimension and token family, or declared by an extension; never Identity, never part of the State Vector. *(Section 8.1)*

### Engineering Domain Registry
the canonical vocabulary of Engineering Domains: the five domains, their stratum positions, and the knowledge they produce and consume. Refinement is taxonomy governance (§14.2). *(Section 8.2)*

### Exchange Layer
the third architectural layer: "transformational boundary: serialization, validation, import/export, mediation of external systems". Owns "exchange contracts, round-trip rules, conformance validation"; does not own "Content and State (never becomes an owner)". *(Section 4.1)*

### Execution Container
an execution Artifact that wraps work items and carries a concurrency convention (exactly-one-active). Its State Domain: Container State. Example implementation: sprint. *(Section 3)*

### Execution State
work item progress State Domain: values `Planned → Todo → In Progress → In Review → Done`; responsibility "Work item progress through Protocol"; owner "Operating Layer (single-writer: work item artifact)". Transition rules: "Strictly sequential; never skip; never revert; one initial; one terminal; every transition recorded in the Change Log." *(Section 7.2)*

### Existence State
Artifact presence State Domain: values `Active → Archived → Retired`; responsibility "Artifact presence in active process vs preservation"; owner "Operating Layer (transitions); Knowledge Layer (preservation principle, P12)". Transition rules: "Forward-only; **Archived** = reference-only (no other State transitions, no Content mutation, still available in retrieval); **Retired** = terminal preservation (not surfaced in normal retrieval, Content immutable, Identity unchanged). Applies to all Artifact types." *(Section 7.2)*

## G

### Gate
a condition that must be satisfied before a transition or execution may occur (approval gate, readiness gate, review gate). *(Section 3)*

## I

### Identity
the property of an Artifact that permanently distinguishes it from all other Artifacts: `(Namespace, Type, ID[, InstanceVersion])`. See Section 6. *(Section 3)*

### Identity Registry
a Knowledge Layer function that guarantees Identity uniqueness, Identity resolution to Artifacts, and referential integrity. *(Section 3)*

### InstanceVersion
part of the instance Identity; the instance discriminator within one Line. "InstanceVersion (the version pointing to a different instance) — **part of the instance Identity**. A new instance is created deliberately (example: plan v2 after v1 is locked). Instance change = instance Identity change; the Line remains." *(Sections 6.1, 6.3)*

## K

### Knowledge Dimension
an axis of knowledge classification (Section 8). An Artifact property, not Identity. *(Section 3)*

### Knowledge Layer
the first architectural layer: "knowledge store: Content, classification, preservation, references". Owns "Content, classification (Knowledge Dimensions), Relationship, history/Records, Identity administration (Identity Registry)"; does not own "process State, execution Protocol". *(Section 4.1)*

### Knowledge OS
a future knowledge execution platform that consumes and produces EKA Artifacts through the Exchange Layer. Not part of the standard; a consumer of the standard. *(Section 3)*

### Knowledge Stratum
the authority level of an Engineering Domain: a fixed position in the strict linear order Discovery → Architecture → Planning → Execution → Operations (stratum 1 highest → 5). Always derived from the Engineering Domain — never declared by an Artifact, never part of the State Vector. One Artifact has exactly one Engineering Domain and therefore exactly one stratum. *(Section 8.1)*

## L

### Layer
an architectural layer with explicit contracts: Knowledge Layer, Operating Layer, Exchange Layer (Section 4). *(Section 3)*

## N

### Namespace
an Identity space that separates management domains (products, organizations, systems). *(Section 3)*

## O

### Operating Layer
the second architectural layer: "State machine & execution Protocol". Owns "State Domains (Execution, Planning, Container, Existence), ordering, concurrency, locking, Gates, Command"; does not own "Content (never edits Content)". *(Section 4.1)*

## P

### Phase
product/scope context over time (Discovery, MVP, Milestone, Release). Phase is a **context attribute** on planning/scope Artifacts: not a category, not a State Domain. **Phase change** is a context update authorized by a readiness Gate (Sections 7.5, 11.2). *(Section 3)*

### Planning State
plan commitment State Domain: values `Draft → Approved → Immutable`; responsibility "Plan commitment: tentative → committed → locked"; owner Operating Layer. Transition rules: "Forward-only; Approved = ready for execution; **Immutable is reached atomically with the Execution Container creation event** (lock-atomic-with-generation); post-lock changes = new instance (InstanceVersion)." *(Section 7.2)*

### Produced Knowledge
the knowledge an Engineering Domain produces: the Artifact types homed in the domain. *(Section 8.2)*

### Protocol
deterministic rules owned by the Operating Layer: ordering, State transitions, locking, gates, execution commands. *(Section 3)*

### Projection Refresh
the mechanism and timing of State Projection validation against the owner (on-read and/or event). *(Section 3)*

### Projection Semantics
the computation rules of a State Projection: from which owner, how aggregation works, what is displayed. *(Section 3)*

## R

### Record
an Artifact preserved as history (Superseded, Archived, Retired, release record) with immutable Content. *(Section 3)*

### Relationship
an explicit relation between Artifacts referenced by Identity: supersedes, amends, derives-from, depends-on, validates. *(Section 3)*

### Representation Alias
a methodology term mapped onto a canonical token + Engineering Domain (e.g., PRD → req-, Sprint/Iteration → ctr-); never a frontmatter value, never an Artifact type in its own right. *(Section 8.1)*

### Representation Alias Registry
the separate convention-layer document mapping methodology terms onto canonical tokens + Engineering Domain; Representation Aliases are not part of the Engineering Domain Registry. *(Section 8.2)*

### Revision
tracking the Content evolution of the same instance — **not part of Identity**. "Revision (tracking the Content evolution of the same instance) — **not part of Identity**. Revision changes on every edit and must not break references." *(Sections 6.1, 6.3)*

## S

### State
a fact about the position of an Artifact within a given process. *(Section 3)*

### State Domain
an independent State dimension with its own semantics, owner, and transition rules (Section 7). Domains are orthogonal. *(Section 3)*

### State Projection
a State view derived from the owner (example: aggregate of work item State → Execution Container status). A projection has no State of its own; a projection never becomes a writer. *(Section 3)*

### State Vector
the tuple of State Domains **owned** by an Artifact. Projected State (State Projection) is not part of the State Vector. *(Section 3)*

## T

### Trigger
a relationship between domains/operations in which an event triggers validation or a transition in another domain (formal definition of interaction, Section 7.5). *(Section 3)*

## W

### Well-formed Content
Content that conforms to the structure established for its Artifact type, so that it can be parsed and executed deterministically. *(Section 3)*

---

*Canonical EKA v1.1 glossary — definitions are binding; see `eka-specification-v1.1.md` for full context.*

**Terminology governance note:** this glossary is the source of Core term definitions. Official ecosystem naming (product identity, Specification Families, reference components, tooling, repository naming) and the deprecated terminology list are governed by the **EKA Naming and Terminology Specification v1.1** (`eka-naming-and-terminology-specification-v1.1.md`). New terms follow the terminology governance there — terms are extended, never forked.
