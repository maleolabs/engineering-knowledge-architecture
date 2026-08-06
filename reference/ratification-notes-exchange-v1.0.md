# Ratification Notes — EKA Exchange Specification v1.0

> Anchor EKA: Exchange Layer — exchange contract (EKA 13, 16.1 milestone 2). Convention document, not an artifact.
> Standard: EKA v1.0, dated 2026-08-05.

## Status

**Ratified** — EKA Exchange Specification v1.0 (refinement pass complete; architecture identical to the previous version, no architectural decision reopened).

## Refinement Summary

The pre-ratification refinement pass — not a redesign. Invariants, rules, and architectural decisions of the Exchange Specification v1.0 **are unchanged**. Changes are additive and editorial:

| Area | Change | Nature |
|---|---|---|
| Exchange Object Model (new §4.4) | Canonical Exchange Package object model: containment hierarchy (Contract Header, Manifest, Exchange Units, Declarations, integrity data), Exchange Unit composition, cardinalities, determinism, and the projection principle ("every serialization is a projection of this model") | Additive — new abstraction, does not change §10 rules |
| Capability Declaration (new §4.5 + §10.1) | Optional declarative block in the Contract Header: supported specification families, extension classes, Relationship type extensions, state variants, export scopes, future protocol capabilities; warning-not-blocking for mismatches | Additive — optional, backward compatible |
| Single declaration location (§7.3, §10.1–10.2, §10.4, §18.3) | Closure/External Reference/Extension declarations = package-level Declarations element; the Header only announces; the Manifest is only the unit list + scope type | Consistency — removes location contradictions |
| Idempotency clarified (§11.1 phases 6–7, §11.2) | Identical units (same Identity + payload) → duplicate/no-op at phase 7; phase 6 conflicts only for same Identity + different payload | Consistency — aligned with §15.5 |
| Complete Conformance Rule coverage (§11.1 phases 3, 10) | Rule R6 (classification) enters the pipeline at phase 3; R8 (single-writer/projection) evaluated via phase 10 revalidation against repository state | Consistency — no rule without a phase |
| "validation rule N" references → "rule N (§14.2)" | Consistent with EKA Naming and Terminology Specification v1.0 | Terminology |
| Conformance Rules table labeled R1–R9 (§14.2) + R0 note, ADR-superseded, change-log order | Aligned with the rule 1–9 ↔ R1–R9 mapping (Naming §9.3) | Terminology |
| Editorial fixes (#13, #19, #21–25) | Phase 10 rollback prose, "permissible differences" (15.4), §1.4 dedupe, §9.1 wording, reading note (glossary + subordinate terms), official footer | Editorial |

## Terminology

Verdict: **eight core terms retained** — Exchange Package, Exchange Unit, Contract Header, Manifest, Export Scope, Import Manifest, Closure, External Reference. No renames: all terms were already registered (Naming and Terminology §12.1), unambiguous, and a rename = contract change (N1). New terms added: **Capability Declaration**, **Closure Declaration**, **Collection**, **Graph** — registered in the §4.1 concept table (defined-before-use, per terminology governance).

## New Conceptual Model

1. **Canonical Exchange Package Object Model (§4.4)** — the canonical abstraction serving as the reference for deriving future serialization formats. Key properties: explicit cardinality per element, the projection principle, determinism, and round-trip equality defined as equality of projections of the same model.
2. **Capability Declaration (§4.5)** — declaration of implementation support, optional, never a package validity condition; mismatch = warning in the Import Manifest; the rejection path stays with §9.2 (versions) and §16.3 (unknown extensions).

## Consistency

Final consistency audit (alex-qa): **0 Critical, 4 Major, 16 Minor, 5 Editorial** — all fixed at the text level. Positive verification: 19 glossary terms consistent, all internal + canonical cross-references resolve, the R1–R9 ↔ rule 1–9 mapping defensible, draft tolerance consistent at 4 points, idempotency/round-trip consistent, EKA Naming and Terminology Specification compliance (H1 pattern, Anchor, Status, slug, English language).

## Architectural Readiness

**Architecture declared ready for ratification.** No blockers. Further changes are limited to editorial fixes and terminology clarification; any additive change in future versions follows the minor version rule (Naming and Terminology §5.3).

---

*End of Ratification Notes — EKA Exchange Specification v1.0.*
