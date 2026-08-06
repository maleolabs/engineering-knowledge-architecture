# EKA Representation Alias Registry v1.1

| Field | Value |
|---|---|
| **Title** | EKA Representation Alias Registry |
| **Version** | 1.1 |
| **Status** | Ratified (convention layer) |
| **Anchor** | EKA Core Specification v1.1 §8.1/§8.2 |
| **Zone** | standard/ |
| **Scope** | The registry of methodology representation names mapped onto canonical Engineering Domains and Artifact types; convention-layer document, never part of the Core standard. |

**Reading this document:** capitalized terms refer to canonical definitions (EKA Core Specification v1.1 Section 3 + glossary). "must" = binding requirement. "should" = recommendation. "may" = option. The Engineering Domain Registry (Core §8.2) remains the authoritative domain vocabulary; this registry maps methodology terms onto it and never amends it. Representation Alias is a registered term (Naming and Terminology Specification v1.1 §2.3, §9.1; Core §8.1).

## 1. Purpose

1.1 This document is the canonical register of methodology representation names — PRD, ADR, RFC, Epic, Sprint, Ticket, Release, and similar — mapped onto canonical tokens + Engineering Domains. It is the document named by Core §8.2: Representation Aliases live here, not in the Engineering Domain Registry.

1.2 Canonical domains remain authoritative. Representations remain conventions: an alias is never Identity, never an Artifact type, never a frontmatter value, and never a validation target beyond R9 per token. The canonical token and Artifact type govern Identity and State.

1.3 It exists because methodology vocabulary is unbounded while the canonical model is closed. Without a register, methodology terms would drift into frontmatter, Identity, and validation. The registry keeps the boundary explicit and the extension path controlled.

## 2. Canonical vs Methodology

**Canonical.** The 26 tokens and five Engineering Domains (Core §8.1) are the model. Tokens are type values; domains are classification. Both are normative; both are closed under taxonomy governance (§14.2).

**Methodology.** Representation names are convention vocabulary: human-facing terms from Scrum, Kanban, Shape Up, or internal workflows. They map onto the canonical model; they never add to it.

| Engineering Domain | Token families | Representation examples |
|---|---|---|
| **Discovery** | vis-, str-, req-, fnd- | PRD, Vision Document, Product Brief, Business Requirement, Research Note |
| **Architecture** | arc-, adr-, dec-, spec-, std-, gls- | ADR, RFC, Decision Record, Design Doc, Specification, Standard, Glossary |
| **Planning** | scp-, epc-, plan-, trc- | Epic, Initiative, Feature, Backlog, Roadmap, Iteration Plan |
| **Execution** | rvw-, ctr-, tkt-, sto-, ts-, bug-, td-, ch-, spk-, ses- | Sprint, Iteration, Cycle, Wave, Task, Work Item, Ticket, Card, Bug, Defect, Incident, Chore, Spike, Session |
| **Operations** | run-, rel- | Release, Incident Report, Maintenance Log, Changelog, Runbook, Playbook, SOP, Post-Mortem |

**Full alias → token → domain mapping.** Rows marked *canonical* are the Core §8.1 alias table, reproduced verbatim. Unmarked rows are extensions registered here.

| Representation Alias | Canonical token | Engineering Domain |
|---|---|---|
| PRD | req- | Discovery |
| Vision Document | vis- | Discovery |
| Product Brief | req- | Discovery |
| Business Requirement | req- | Discovery |
| Research Note | fnd- | Discovery |
| ADR | adr- | Architecture |
| RFC | adr- | Architecture |
| RFC | spec- | Architecture |
| Decision Record | dec- | Architecture |
| Design Doc | arc- | Architecture |
| Specification | spec- | Architecture |
| Standard | std- | Architecture |
| Glossary | gls- | Architecture |
| Epic | epc- | Planning |
| Initiative | scp- | Planning |
| Feature | epc- | Planning |
| Backlog | plan- | Planning |
| Roadmap | plan- | Planning |
| Iteration Plan | plan- | Planning |
| Sprint | ctr- | Execution |
| Iteration | ctr- | Execution |
| Cycle | ctr- | Execution |
| Wave | ctr- | Execution |
| Story | sto- | Execution |
| Task | sto- | Execution |
| Work Item | sto- | Execution |
| Ticket | tkt- | Execution |
| Card | tkt- | Execution |
| Bug | bug- | Execution |
| Defect | bug- | Execution |
| Incident | bug- | Execution |
| Chore | ch- | Execution |
| Spike | spk- | Execution |
| Session | ses- | Execution |
| Release | rel- | Operations |
| Changelog | rel- | Operations |
| Runbook | run- | Operations |
| Playbook | run- | Operations |
| SOP | run- | Operations |
| Post-Mortem | run- | Operations |
| Incident Report | run- | Operations |
| Maintenance Log | run- | Operations |

**Mapping notes.**

- **RFC → adr-** (canonical, proposal form) or **spec-** (ratified form). An RFC proposing a decision is an ADR; an RFC accepted as normative description of behavior is a Specification. The proposal form is the Core §8.1 mapping.
- **Task → sto-.** Generic work item. When the nature is known, use the specific work item token (ts-, bug-, td-, ch-, spk-) instead; the alias only names the default member.
- **Work Item → sto-.** Generic family term covering the work item token family (sto-, ts-, bug-, td-, ch-, spk-); registered against the default member.
- **Changelog → rel-.** A release-oriented record of changes; distinct from the canonical Change Log (the state-transition record, Core §5.2). The two never share content.
- **Incident Report → run-.** Operational knowledge documenting handling and root cause; distinct from the Incident alias (bug-), which names the tracked defect itself.
- **Card → tkt-.** Kanban card: the ticket projection. No State Vector of its own (Core §7.4).
- **Cycle → ctr-.** Shape Up time-boxed work block; same container semantics as Sprint/Iteration.

## 3. Governance

Extension of the registry is the only path for new representation names. It is **methodology-pack governance**: a methodology (Scrum pack, Kanban pack, Shape Up pack, internal workflow) proposes aliases; the registry accepts them; the Core standard is never touched.

**Process.**

1. **Proposal** — alias name, target token + Engineering Domain, methodology origin. The target token must already exist in the canonical 26-token table (§8.1); a proposal for a new token is taxonomy governance (§14.2), not registry governance.
2. **Review** — three checks: (a) does it collide with an existing alias or token? (b) does the mapping respect the domain semantics — does the term name knowledge homed in the target domain? (c) does it weaken the canonical model — would accepting it blur a domain boundary, an Artifact type, or the token mapping?
3. **Acceptance** — the alias is registered in this document (Section 5). No Core modification needed.
4. **Publication** — the methodology pack documents its terms with their registered aliases.

**Hard rules.**

- An alias never creates a new Artifact type.
- An alias never changes Identity, State, or validation semantics of the mapped Artifact type.
- Aliases are never frontmatter values: unknown aliases in frontmatter are rejected by R0 — the type token is the only valid type value.
- A rejected proposal is rejected for a stated reason; resubmission requires the reason to be addressed.

## 4. Methodology Conventions

Scrum, Kanban, Shape Up, and internal workflows are **convention layers over EKA** (Core §8.1). They map their terminology onto this registry; they may add aliases through the extension process (Section 3); they never extend the canonical domain model.

- A methodology pack names things its own way; the registered alias is the interface to the canonical model.
- A methodology pack may define alias subsets (e.g., Kanban pack: Card, Defect, Chore, Cycle) without affecting other packs.
- No methodology pack is part of the Core standard; no pack may redefine a canonical token, domain, stratum, or Artifact type.
- Terminology governance follows Naming and Terminology Specification v1.1 §12: terms are extended, never forked; deprecation is the only retirement path.

## 5. Registry Table

Complete current registry. Canonical rows (Core §8.1) are marked; all rows are valid aliases for documentation, conversation, and methodology packs — never for frontmatter.

| Alias | Token | Engineering Domain | Notes |
|---|---|---|---|
| PRD | req- | Discovery | canonical (Core §8.1) |
| Vision Document | vis- | Discovery | |
| Product Brief | req- | Discovery | requirements form of intent |
| Business Requirement | req- | Discovery | |
| Research Note | fnd- | Discovery | research findings |
| ADR | adr- | Architecture | canonical (Core §8.1) |
| RFC | adr- | Architecture | canonical (Core §8.1); proposal form |
| RFC | spec- | Architecture | ratified form: normative specification |
| Decision Record | dec- | Architecture | reversible decisions |
| Design Doc | arc- | Architecture | architecture description |
| Specification | spec- | Architecture | |
| Standard | std- | Architecture | |
| Glossary | gls- | Architecture | vocabulary |
| Epic | epc- | Planning | canonical (Core §8.1) |
| Initiative | scp- | Planning | canonical (Core §8.1) |
| Feature | epc- | Planning | sized planning unit |
| Backlog | plan- | Planning | ordered work plan |
| Roadmap | plan- | Planning | |
| Iteration Plan | plan- | Planning | plan for one iteration (ctr-) |
| Sprint | ctr- | Execution | canonical (Core §8.1) |
| Iteration | ctr- | Execution | canonical (Core §8.1) |
| Cycle | ctr- | Execution | Shape Up container |
| Wave | ctr- | Execution | |
| Story | sto- | Execution | |
| Task | sto- | Execution | generic; prefer specific work item token |
| Work Item | sto- | Execution | generic family term |
| Ticket | tkt- | Execution | canonical (Core §8.1); projection |
| Card | tkt- | Execution | Kanban; projection |
| Bug | bug- | Execution | |
| Defect | bug- | Execution | synonym of Bug |
| Incident | bug- | Execution | canonical (Core §8.1) |
| Chore | ch- | Execution | |
| Spike | spk- | Execution | |
| Session | ses- | Execution | |
| Release | rel- | Operations | canonical (Core §8.1); record |
| Changelog | rel- | Operations | record; ≠ canonical Change Log (§5.2) |
| Runbook | run- | Operations | canonical (Core §8.1) |
| Playbook | run- | Operations | operational guide |
| SOP | run- | Operations | standard operating procedure |
| Post-Mortem | run- | Operations | operational knowledge |
| Incident Report | run- | Operations | ≠ Incident (bug-) |
| Maintenance Log | run- | Operations | |

---

*End of Representation Alias Registry — v1.1 (Ratified, convention layer). Anchored to EKA Core Specification v1.1 §8.1/§8.2; Representation Alias per Naming and Terminology Specification v1.1 §2.3.*
