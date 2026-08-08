package store

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

// This file implements the immutable payload store: the content-
// addressed archive of Engineering Knowledge Objects. Every payload
// row is keyed by its content-derived hash — SHA-256(unit.json ||
// content), the same per-unit digest as the RSF — and is written at
// most once: there is deliberately NO update path for payload rows in
// this package. History is the accumulation of payload rows, chained
// per reference through prev_hash (first-writer wins; a payload's
// lineage is fixed the moment it is inserted). SQLite is only the
// persistence layer; the source of immutability is the content hash,
// never a database-generated value.
//
// All SQL is parameterized.

// PutUnit stores one immutable unit payload and points its reference
// at it, in one transaction:
//
//  1. hash = SHA-256(unitJSON || content).
//  2. The payload row is inserted when absent (INSERT ... DO NOTHING):
//     a payload with this hash already in the archive is left
//     untouched (first-writer wins — prev_hash is never modified). For
//     a NEW payload, prev_hash is the object_hash of the reference's
//     current payload ("" when the form has no reference yet): the
//     lineage edge "this object supersedes <prev> within its
//     reference".
//  3. The object_refs row is upserted (the reference is the mutable
//     part — it may move from one immutable payload to another).
//
// PutUnit returns the stored hash. It never updates an existing
// payload row.
func (s *Store) PutUnit(unitJSON, content []byte, r Ref) (string, error) {
	hash := hashUnit(unitJSON, content)

	tx, err := s.db.Begin()
	if err != nil {
		return "", fmt.Errorf("store: cannot begin put: %w", err)
	}
	defer tx.Rollback()

	// Look up the reference's current payload hash BEFORE the insert:
	// it becomes the new payload's lineage predecessor.
	prevHash := ""
	var current string
	err = tx.QueryRow(`SELECT object_hash FROM object_refs WHERE form = ?`, r.Form).Scan(&current)
	switch {
	case err == nil:
		prevHash = current
	case isNoRows(err):
		// No reference yet: a root payload ("" lineage).
	default:
		return "", fmt.Errorf("store: cannot read reference %s: %w", r.Form, err)
	}

	// Immutable insert: a payload with this hash already exists -> no-op
	// (first-writer wins; prev_hash and bytes stay untouched).
	if _, err := tx.Exec(`INSERT INTO object_payloads (object_hash, unit_json, content, prev_hash, created_at)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(object_hash) DO NOTHING`,
		hash, unitJSON, content, prevHash, time.Now().UTC().Format(time.RFC3339)); err != nil {
		return "", fmt.Errorf("store: cannot insert payload %s: %w", hash, err)
	}

	// The reference is the only mutable part: upsert all index columns
	// from the caller's ref (derived from the payload at insert).
	if _, err := tx.Exec(`INSERT INTO object_refs (
		form, object_hash, project_id, source_repo, namespace, type, id,
		instance_version, revision, dimension, domain, phase, updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(form) DO UPDATE SET
		object_hash = excluded.object_hash,
		project_id = excluded.project_id,
		source_repo = excluded.source_repo,
		namespace = excluded.namespace,
		type = excluded.type,
		id = excluded.id,
		instance_version = excluded.instance_version,
		revision = excluded.revision,
		dimension = excluded.dimension,
		domain = excluded.domain,
		phase = excluded.phase,
		updated_at = excluded.updated_at`,
		r.Form, hash, r.ProjectID, r.SourceRepo, r.Namespace, r.Type, r.ID,
		r.InstanceVersion, r.Revision, r.Dimension, r.Domain, r.Phase, r.UpdatedAt); err != nil {
		return "", fmt.Errorf("store: cannot upsert reference %s: %w", r.Form, err)
	}

	if err := tx.Commit(); err != nil {
		return "", fmt.Errorf("store: cannot commit put of %s: %w", hash, err)
	}
	return hash, nil
}

// hashUnit computes the content address of one unit: SHA-256(unit.json
// bytes || content bytes) — the same byte concatenation and digest as
// the RSF per-unit digest (exchange/serialize.go). The concatenation is
// built into an explicit fresh buffer: append() against the caller's
// unitJSON slice could alias its spare capacity.
func hashUnit(unitJSON, content []byte) string {
	buf := make([]byte, 0, len(unitJSON)+len(content))
	buf = append(buf, unitJSON...)
	buf = append(buf, content...)
	sum := sha256.Sum256(buf)
	return hex.EncodeToString(sum[:])
}

// Payload returns the immutable unit_json and content bytes of one
// payload; an error when no payload carries this hash.
func (s *Store) Payload(hash string) (unitJSON, content []byte, err error) {
	var unit, body []byte
	err = s.db.QueryRow(`SELECT unit_json, content FROM object_payloads WHERE object_hash = ?`, hash).
		Scan(&unit, &body)
	if err != nil {
		if isNoRows(err) {
			return nil, nil, fmt.Errorf("store: payload %s not found", hash)
		}
		return nil, nil, fmt.Errorf("store: cannot read payload %s: %w", hash, err)
	}
	return unit, body, nil
}

// PayloadCount returns the number of stored immutable payloads.
func (s *Store) PayloadCount() (int, error) {
	return s.count("object_payloads")
}

// PayloadRow is one immutable payload of the archive.
type PayloadRow struct {
	// ObjectHash is the content-derived key (SHA-256(unit.json ||
	// content)).
	ObjectHash string
	// UnitJSON is the exact canonical RSF unit entry bytes.
	UnitJSON []byte
	// Content is the representation payload bytes (unit content).
	Content []byte
}

// AllPayloads returns every stored payload, sorted by object_hash
// (deterministic order for integrity scans).
func (s *Store) AllPayloads() ([]PayloadRow, error) {
	rows, err := s.db.Query(`SELECT object_hash, unit_json, content FROM object_payloads ORDER BY object_hash`)
	if err != nil {
		return nil, fmt.Errorf("store: cannot query payloads: %w", err)
	}
	defer rows.Close()
	var out []PayloadRow
	for rows.Next() {
		var p PayloadRow
		if err := rows.Scan(&p.ObjectHash, &p.UnitJSON, &p.Content); err != nil {
			return nil, fmt.Errorf("store: cannot scan payload row: %w", err)
		}
		out = append(out, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("store: cannot read payloads: %w", err)
	}
	return out, nil
}
