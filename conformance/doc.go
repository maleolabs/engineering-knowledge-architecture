// Package conformance implements the official EKA conformance validator:
// the mechanical checklist defined in
// skeleton/docs/exchange/validation.md (Rules R1-R9) extended with the
// Engineering Domain rules R10-R12 of the EKA v1.1 standard, for the
// Engineering Knowledge Architecture (EKA) serialization.
//
// The package is deliberately independent from the CLI layer (cmd/eka) so
// that future tooling can import it directly as
// github.com/maleolabs/engineering-knowledge-architecture/conformance.
//
// The single entry point is Validate, which recursively scans a directory
// tree for .md files, classifies them as Artifacts or Convention Documents,
// and applies the twelve conformance rules. The result is a deterministic,
// sorted Report of errors and warnings; warnings never block a pass.
//
// Interpretation decisions that go beyond the literal rule text are
// documented in comments at the relevant check and collected in the
// "Interpretation decisions" section of the CLI/package documentation.
package conformance
