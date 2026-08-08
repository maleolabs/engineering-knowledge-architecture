package store

import (
	"database/sql"
	"errors"
	"fmt"
)

// This file implements the schema migration driver: the meta table
// (schema bookkeeping) plus the schema version steps of the EKA
// workspace database.
//
// Migration contract: Open runs the current schemaVersion's steps
// against a database in a transaction; every step is idempotent in
// practice so re-opening an up-to-date database is a no-op. Future
// schema versions append steps and bump schemaVersion; a database
// newer than this implementation is refused (downgrade protection).
//
// Schema v2 (the Immutable Engineering Knowledge Model): Engineering
// Knowledge Objects are immutable and content-addressed. The store
// keeps two tables:
//
//   - object_payloads: the immutable payload archive, keyed by the
//     content-derived hash SHA-256(unit.json || content) (== the RSF
//     per-unit digest). Payload rows are written once and never
//     updated; history is the accumulation of payload rows, chained
//     per reference through prev_hash. SQLite is only the persistence
//     layer — the source of immutability is the content hash, never a
//     database-generated value.
//   - object_refs: the mutable reference table (resolver key = RSF
//     canonical identity form), with derived index columns filled from
//     the payload at insert.
//
// Relationships and change logs are serialized inside unit.json (the
// payload); the v1 objects/relationships/change_log tables are gone.
// The remaining v1 tables (projects, repos, attachments, sync_log,
// meta) are unchanged. The v1 -> v2 migration (migrate.go)
// reconstructs each v1 object's unit payload, recomputes the content
// hash, and drops the v1 tables.

// migrate opens the schema and checks/advances the recorded version.
func (s *Store) migrate() error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("store: cannot begin migration: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.Exec(`CREATE TABLE IF NOT EXISTS meta (
		key   TEXT PRIMARY KEY,
		value TEXT NOT NULL
	)`); err != nil {
		return fmt.Errorf("store: cannot create meta table: %w", err)
	}

	var current int
	err = tx.QueryRow(`SELECT value FROM meta WHERE key = 'schema_version'`).Scan(&current)
	switch {
	case err == nil:
		if current > schemaVersion {
			return fmt.Errorf("store: database schema version %d is newer than this implementation (%d); upgrade the CLI", current, schemaVersion)
		}
	case isNoRows(err):
		current = 0
	default:
		return fmt.Errorf("store: cannot read schema_version: %w", err)
	}

	// Migration steps run in order; each step advances `current`. A
	// fresh database (current == 0, no v1 tables) goes straight to v2.
	// A database at schema v1 (or a partially migrated v1 database
	// whose tables exist but whose version was never recorded) runs the
	// v1 -> v2 conversion.
	switch {
	case current == 0 && !s.hasTable(tx, "objects"):
		if err := migrateToV2(tx); err != nil {
			return err
		}
		current = schemaVersion
	case current <= 1:
		if err := migrateV1toV2(tx); err != nil {
			return err
		}
		current = schemaVersion
	}

	if _, err := tx.Exec(`INSERT INTO meta (key, value) VALUES ('schema_version', ?)
		ON CONFLICT(key) DO UPDATE SET value = excluded.value`, fmt.Sprint(current)); err != nil {
		return fmt.Errorf("store: cannot record schema_version: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("store: cannot commit migration: %w", err)
	}
	return nil
}

// hasTable reports whether a table exists in the database (schema
// detection for the fresh-vs-v1 decision).
func (s *Store) hasTable(tx *sql.Tx, name string) bool {
	var n int
	err := tx.QueryRow(`SELECT COUNT(*) FROM sqlite_master WHERE type = 'table' AND name = ?`, name).Scan(&n)
	return err == nil && n > 0
}

// migrateToV2 creates the v2 schema on a fresh database: the shared
// registry tables (unchanged from v1) plus the immutable payload and
// reference tables. Steps run unconditionally (IF NOT EXISTS) so a
// partial database completes.
func migrateToV2(tx *sql.Tx) error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS projects (
			id      TEXT PRIMARY KEY,
			name    TEXT NOT NULL,
			created TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS repos (
			project_id TEXT NOT NULL REFERENCES projects(id),
			name       TEXT NOT NULL,
			path       TEXT NOT NULL,
			created    TEXT NOT NULL,
			PRIMARY KEY (project_id, name)
		)`,
		// A repository path is registered in exactly one project: the
		// first project to register the path owns it (registry
		// determinism; see workspace.RegisterRepo).
		`CREATE UNIQUE INDEX IF NOT EXISTS repos_path_uniq ON repos(path)`,
		// The immutable payload archive: one row per content-addressed
		// Engineering Knowledge Object. object_hash is SHA-256(unit.json
		// || content) — the RSF per-unit digest — and is never
		// database-generated. Payload rows are INSERT-only (no update
		// path exists in this package; the store API has no payload
		// mutation beyond first insert).
		`CREATE TABLE IF NOT EXISTS object_payloads (
			object_hash TEXT PRIMARY KEY,
			unit_json   BLOB NOT NULL,
			content     BLOB NOT NULL,
			prev_hash   TEXT NOT NULL,
			created_at  TEXT NOT NULL
		)`,
		// The mutable reference table: the resolver key is the RSF
		// canonical identity form; the row points at the current
		// immutable payload of that identity within its provenance
		// (project_id, source_repo). The index columns are derived from
		// the payload at insert (never stored independently in v1
		// columns). The FK keeps every reference pointing at a real
		// payload (foreign_keys pragma is on).
		`CREATE TABLE IF NOT EXISTS object_refs (
			form             TEXT PRIMARY KEY,
			object_hash      TEXT NOT NULL REFERENCES object_payloads(object_hash),
			project_id       TEXT NOT NULL,
			source_repo      TEXT NOT NULL,
			namespace        TEXT NOT NULL,
			type             TEXT NOT NULL,
			id               TEXT NOT NULL,
			instance_version INTEGER NOT NULL,
			revision         INTEGER NOT NULL,
			dimension        TEXT NOT NULL,
			domain           TEXT NOT NULL,
			phase            TEXT NOT NULL,
			updated_at       TEXT NOT NULL
		)`,
		`CREATE INDEX IF NOT EXISTS idx_refs_identity ON object_refs (namespace, type, id)`,
		`CREATE INDEX IF NOT EXISTS idx_refs_dimension ON object_refs (dimension)`,
		`CREATE INDEX IF NOT EXISTS idx_refs_domain ON object_refs (domain)`,
		`CREATE INDEX IF NOT EXISTS idx_refs_source ON object_refs (project_id, source_repo)`,
		// Reverse hash lookup (replication, future GC of the history
		// archive): which references point at a payload.
		`CREATE INDEX IF NOT EXISTS idx_refs_hash ON object_refs (object_hash)`,
		`CREATE TABLE IF NOT EXISTS attachments (
			project_id  TEXT NOT NULL,
			source_repo TEXT NOT NULL,
			id          TEXT NOT NULL,
			digest      TEXT NOT NULL,
			data        BLOB NOT NULL,
			PRIMARY KEY (project_id, source_repo, id)
		)`,
		`CREATE TABLE IF NOT EXISTS sync_log (
			seq             INTEGER PRIMARY KEY AUTOINCREMENT,
			project_id      TEXT NOT NULL,
			repo            TEXT NOT NULL,
			direction       TEXT NOT NULL,
			snapshot_digest TEXT NOT NULL,
			units           INTEGER NOT NULL,
			at              TEXT NOT NULL
		)`,
		`CREATE INDEX IF NOT EXISTS sync_log_repo_idx ON sync_log(project_id, repo, seq)`,
	}
	for _, stmt := range statements {
		if _, err := tx.Exec(stmt); err != nil {
			return fmt.Errorf("store: cannot create table: %w", err)
		}
	}
	return nil
}

// isNoRows reports whether err is sql.ErrNoRows (driver-agnostic).
func isNoRows(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}
