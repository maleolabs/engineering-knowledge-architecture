// Package store implements the EKA Workspace database: the local SQLite
// canonical store backing the EKA Knowledge Runtime (milestone EKA
// v0.2.0). The workspace keeps the canonical projection of every
// registered repository's knowledge: immutable Engineering Knowledge
// Objects (content-addressed payloads), their mutable references,
// attachments, and the sync log.
//
// The Immutable Engineering Knowledge Model: knowledge objects are
// IMMUTABLE and content-addressed. Every change produces a new
// immutable payload; mutable state is limited to references (object_refs)
// and indexes. SQLite is only the persistence layer — it is never the
// source of immutability (the store API has no update path for payload
// rows; hashes are content-derived SHA-256, never database-generated).
// History emerges from immutable payloads + references (prev_hash
// lineage); there is no dedicated mutable history table. Relationships
// and change logs are serialized inside the payload's unit.json.
//
// The database lives at <workspace>/eka.db. It is opened with the
// driver modernc.org/sqlite (pure Go, no cgo) and the following
// pragmas:
//
//	journal_mode(WAL)   concurrent readers + one writer
//	busy_timeout(5000)  serialize writers instead of failing fast
//	foreign_keys(1)     referential integrity between tables
//
// Concurrency model (documented decision): the runtime assumes a single
// writer (the CLI is single-process). WAL plus busy_timeout make
// concurrent processes safe without file locking; no flock is
// implemented.
//
// All SQL is parameterized; values are never interpolated into
// statements.
package store

import (
	"database/sql"
	"fmt"
	"net/url"
	"path/filepath"

	_ "modernc.org/sqlite" // registers the "sqlite" driver on import.
)

// schemaVersion is the current database schema version. Future
// migrations bump it and append migration steps (schema.go,
// migrate.go); a database at an older version is migrated forward at
// Open.
const schemaVersion = 2

// Store is one opened EKA workspace database. It is safe for use by a
// single process; concurrent use is serialized by SQLite itself (WAL +
// busy_timeout).
type Store struct {
	db *sql.DB
}

// Open opens (creating when missing) the workspace database at
// dir/eka.db and migrates it to the current schema version. The parent
// directory must already exist. On any failure the database is closed
// and an error is returned.
func Open(dir string) (*Store, error) {
	if dir == "" {
		return nil, fmt.Errorf("store: empty directory")
	}
	path := filepath.Join(dir, "eka.db")
	dsn := "file:" + url.PathEscape(path) +
		"?_pragma=journal_mode(WAL)" +
		"&_pragma=busy_timeout(5000)" +
		"&_pragma=foreign_keys(1)"
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("store: cannot open %s: %w", path, err)
	}
	s := &Store{db: db}
	if err := s.migrate(); err != nil {
		db.Close()
		return nil, err
	}
	return s, nil
}

// Close closes the underlying database.
func (s *Store) Close() error {
	return s.db.Close()
}

// DB returns the underlying database handle. It is exposed for the
// workspace registry (projects/repos tables, which are workspace-level
// concerns); all queries against it must be parameterized.
func (s *Store) DB() *sql.DB { return s.db }
