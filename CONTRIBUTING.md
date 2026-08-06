# Contributing — EKA v1.1 Repository

Contribution guide and governance rules for the **EKA v1.1 Reference Implementation** repository (spec + Reference Implementation + validator + tooling).

> Convention document, not an artifact. Related: [`reference/conformance-traceability-matrix.md`](reference/conformance-traceability-matrix.md) (single source of truth for conformance coverage), [`reference/conformance-notes.md`](reference/conformance-notes.md) (interpretation decisions), [`skeleton/docs/exchange/validation.md`](skeleton/docs/exchange/validation.md) (Conformance Rules R1–R9), [`standard/eka-specification-v1.1.md`](standard/eka-specification-v1.1.md) (canonical standard).

## Principles

- **One canonical project.** This repository is a single whole: the standard (`standard/`), the Reference Implementation (`skeleton/`), the validator (`conformance/` + `cmd/eka/`), and the documentation (`reference/`). Changes must not treat one part as a separate project.
- **A complete contribution = spec change → validator → tests → matrix in the SAME Pull Request.** No partial PR that changes conformance behavior without closing the entire traceability chain.
- **Documentation is part of the product.** Every interpretation decision is documented before implementation; no behavior is invented without a spec basis (see the interpretation policy in `reference/conformance-notes.md`).

## Definition of "incomplete implementation"

A contribution is considered **incomplete** if it changes the spec, the Conformance Rules, the validator, or the tests **WITHOUT updating** [`reference/conformance-traceability-matrix.md`](reference/conformance-traceability-matrix.md). Such a PR is rejected until the matrix is synchronized.

## Contribution workflow

1. **Understand EKA.** Read `standard/` (canonical standard — normative text), `skeleton/docs/` (structure and serialization rules), and `reference/` (implementation conventions). Start from the root `README.md`.
2. **If changing conformance behavior**, update in the same PR:
   - rule text: `skeleton/docs/exchange/validation.md` (R1–R9) or the R0 definition in `conformance/`;
   - implementation: `conformance/` (rule/helper) or `cmd/eka/` (CLI);
   - tests: `*_test.go` in `conformance/` + `cmd/eka/`;
   - matrix: `reference/conformance-traceability-matrix.md` (Sections 2–4);
   - interpretation notes: `reference/conformance-notes.md` — add **new** interpretation decisions with the next number; old decisions are never removed.
3. **Mandatory quality** — all must pass before a PR:
   ```sh
   go test ./...
   go vet ./...
   gofmt -l .
   go build ./...
   go run ./cmd/eka validate .
   ```
   `go run ./cmd/eka validate .` (self-validation) must PASS: 0 errors, exit 0.
4. **Matrix maintenance rules**:
   - **New rule** → new row in Section 2 (main matrix) + new `REQ-nnn` requirement in Section 3. REQ IDs are never reused; new numbers continue the sequence.
   - **Rule modified** → update the Spec Anchor / Implementation / Automated Tests / Notes columns on the affected rule row; if the requirement semantics change, update the requirement phrasing in Section 3.
   - **New test** → update Section 4 (test coverage index). Every test function appears exactly once — no duplication, no orphans.
   - **Coverage status changed** → update Section 5 (coverage analysis) and Section 6 (summary) with consistent numbers.
5. **Documentation must stay consistent**: if the CLI changes, update `reference/cli.md`; if an interpretation changes, update `reference/conformance-notes.md`; if a new document is added to the `reference/` zone, add an index row in `reference/README.md`.

## Formal governance rules

- [`reference/conformance-traceability-matrix.md`](reference/conformance-traceability-matrix.md) is the **single source of truth** for conformance coverage: Engineering Requirement → Specification → Conformance Rule → Implementation → Automated Test → Coverage Status → Notes.
- The matrix **MUST be updated in the same Pull Request** as spec/rule/implementation/test changes — and **conversely**: the matrix is never edited without a related change (no "matrix-only" PRs).
- Rule IDs (R0–R12) and requirement IDs (REQ-nnn) are stable and never reused.
- Validator behavior changes outside R0–R12 (new rules) are **not** proposed through direct implementation: new rules are a governance decision (see the known gaps in `reference/conformance-notes.md` and the follow-up recommendations in matrix Sections 5–6).

## Review checklist (for reviewers)

Before approving a conformance-touching PR:

- [ ] Matrix consistent with the changes: Section 2 (rule), Section 3 (requirement), Section 4 (test), Sections 5–6 (analysis & summary)?
- [ ] No orphans: every rule R0–R12 appears exactly once in Section 2; every test function appears exactly once in Section 4?
- [ ] No rule/requirement ID reused?
- [ ] Summary numbers valid: test count in Section 6 == number of test rows in Section 4?
- [ ] `go test ./...`, `go vet ./...`, `gofmt`, `go build`, and `go run ./cmd/eka validate .` pass?
- [ ] New interpretations documented in `reference/conformance-notes.md` (continuing numbers), if any?
- [ ] Related documentation (README, `reference/README.md`, `reference/cli.md`) updated if names/structure changed?
