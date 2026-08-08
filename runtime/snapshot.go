package runtime

import (
	"github.com/maleolabs/engineering-knowledge-architecture/exchange"
)

// This file implements the SnapshotService: the verified package read
// side of the Runtime. Writing snapshots happens through the Authoring
// API (Authoring.Sync — push side); this service is the read side:
// full byte-exact verification of an RSF package wherever it lives on
// disk.

// SnapshotService reads and verifies RSF packages. Concrete and
// documented — no interface type.
type SnapshotService struct{ rt *Runtime }

// Read opens and fully verifies the package at path — a single-file
// .ekapkg ZIP container or the equivalent directory layout — and
// returns the deserialized RSF object model. The verification is
// byte-exact: entry structure, strict JSON decoding (unknown fields
// rejected), SHA-256 integrity (package, per-unit, per-attachment
// digests) and manifest self-consistency. Any failure is a
// *exchange.PackageError (rejected package class).
//
// Writing snapshots happens through Authoring.Sync (the push side of
// the sync engine) — this service is read-only.
func (s *SnapshotService) Read(path string) (*exchange.Package, error) {
	return exchange.LoadPackage(path)
}
