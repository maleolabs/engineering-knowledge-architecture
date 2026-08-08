package store

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"path/filepath"
	"strings"
	"testing"

	"github.com/maleolabs/engineering-knowledge-architecture/exchange"
)

// openTest opens a fresh store in a temp directory for one test.
func openTest(t *testing.T) *Store {
	t.Helper()
	s, err := Open(t.TempDir())
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() { s.Close() })
	return s
}

// unitJSON builds a canonical unit.json for one identity with the given
// content payload, via the exchange serializer (the same bytes a
// package carries).
func unitJSON(t *testing.T, ns, typ, id string, version, revision int) []byte {
	t.Helper()
	u := testUnit(ns, typ, id, version, revision)
	data, err := exchange.MarshalUnit(u)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

// testUnit builds a minimal exchange.Unit for one identity.
func testUnit(ns, typ, id string, version, revision int) *exchange.Unit {
	u := &exchange.Unit{
		Identity: exchange.Identity{Namespace: ns, Type: typ, ID: id, InstanceVersion: version},
		Revision: revision,
		StateVector: exchange.StateVector{
			ContentState:   "draft",
			ExistenceState: "active",
		},
		ChangeLog:     []exchange.ChangeLogEntry{},
		Relationships: []exchange.Relationship{},
		Classification: exchange.Classification{
			Dimension: "decisions",
			Domain:    "Architecture",
		},
		Phase:          "wave-1",
		Content:        exchange.ContentRef{Representation: "eka/structured-text/1", File: "content"},
		ContentPayload: []byte("body"),
	}
	u.CanonicalIdentityForm = u.Identity.CanonicalForm()
	return u
}

// ref builds the reference of one unit for the given provenance.
func ref(form, project, repo string) Ref {
	return Ref{
		Form:            form,
		ProjectID:       project,
		SourceRepo:      repo,
		Namespace:       "acme",
		Type:            "sto",
		ID:              "x",
		InstanceVersion: 1,
		Revision:        1,
		Dimension:       "decisions",
		Domain:          "Architecture",
		Phase:           "wave-1",
		UpdatedAt:       "2026-08-07T00:00:00Z",
	}
}

func TestOpenCreatesDatabaseAndSchemaVersion(t *testing.T) {
	dir := t.TempDir()
	s, err := Open(dir)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer s.Close()

	// eka.db must exist on disk.
	if _, err := filepath.Glob(filepath.Join(dir, "eka.db*")); err != nil {
		t.Fatal(err)
	}
	meta, err := s.Meta()
	if err != nil {
		t.Fatal(err)
	}
	if meta["schema_version"] != "2" {
		t.Errorf("schema_version = %q, want 2", meta["schema_version"])
	}
	v, err := s.SchemaVersion()
	if err != nil {
		t.Fatal(err)
	}
	if v != 2 {
		t.Errorf("SchemaVersion = %d, want 2", v)
	}
	// Fresh databases go straight to the v2 schema: no v1 tables.
	for _, table := range []string{"objects", "relationships", "change_log"} {
		if hasTableRaw(t, s, table) {
			t.Errorf("fresh database must not contain the v1 table %s", table)
		}
	}
	for _, table := range []string{"object_payloads", "object_refs"} {
		if !hasTableRaw(t, s, table) {
			t.Errorf("fresh database must contain the v2 table %s", table)
		}
	}
}

func TestOpenIsIdempotent(t *testing.T) {
	dir := t.TempDir()
	s1, err := Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	s1.Close()
	s2, err := Open(dir)
	if err != nil {
		t.Fatalf("second Open: %v", err)
	}
	defer s2.Close()
	v, err := s2.SchemaVersion()
	if err != nil {
		t.Fatal(err)
	}
	if v != 2 {
		t.Errorf("second open schema version = %d, want 2", v)
	}
}

func TestSetMetaRoundtrip(t *testing.T) {
	s := openTest(t)
	if err := s.SetMeta("k1", "v1"); err != nil {
		t.Fatal(err)
	}
	if err := s.SetMeta("k2", "v2"); err != nil {
		t.Fatal(err)
	}
	// Upsert semantics: same key replaces.
	if err := s.SetMeta("k1", "v1b"); err != nil {
		t.Fatal(err)
	}
	meta, err := s.Meta()
	if err != nil {
		t.Fatal(err)
	}
	if meta["k1"] != "v1b" || meta["k2"] != "v2" {
		t.Errorf("meta = %v, want k1=v1b k2=v2", meta)
	}
}

// TestPutUnitImmutablePayload: the core immutability contract — the
// same content always maps to the same hash, a payload row is written
// once and never modified, and a changed content produces a NEW payload
// row chained through prev_hash while the reference moves.
func TestPutUnitImmutablePayload(t *testing.T) {
	s := openTest(t)
	form := "acme/sto:x:1"
	uj1 := unitJSON(t, "acme", "sto", "x", 1, 1)
	uj2 := unitJSON(t, "acme", "sto", "x", 1, 2)

	// First insert: root payload (prev_hash "").
	h1, err := s.PutUnit(uj1, []byte("body v1"), ref(form, "p", "r"))
	if err != nil {
		t.Fatal(err)
	}
	if h1 == "" {
		t.Fatal("hash must be non-empty")
	}
	got, ok, err := s.Ref(form)
	if err != nil || !ok {
		t.Fatalf("Ref = %v, %v", ok, err)
	}
	if got.ObjectHash != h1 {
		t.Errorf("reference points at %s, want %s", got.ObjectHash, h1)
	}
	prev := payloadPrevHash(t, s, h1)
	if prev != "" {
		t.Errorf("first payload prev_hash = %q, want root \"\"", prev)
	}

	// Same content, same form: same hash, payload row untouched, one
	// row total. The reference upsert is idempotent.
	h1b, err := s.PutUnit(uj1, []byte("body v1"), ref(form, "p", "r"))
	if err != nil {
		t.Fatal(err)
	}
	if h1b != h1 {
		t.Errorf("same content must return the same hash: %s vs %s", h1b, h1)
	}
	if n := payloadRowCount(t, s, h1); n != 1 {
		t.Errorf("payload rows for %s = %d, want 1 (immutable, no duplicates)", h1, n)
	}
	if prev := payloadPrevHash(t, s, h1); prev != "" {
		t.Errorf("re-insert must not modify prev_hash: got %q", prev)
	}
	unit, content, err := s.Payload(h1)
	if err != nil {
		t.Fatal(err)
	}
	if string(unit) != string(uj1) || string(content) != "body v1" {
		t.Error("re-insert must not modify payload bytes")
	}

	// Different content: NEW payload row; the reference moves to it;
	// the new payload's prev_hash is the old hash (lineage within the
	// reference). The old payload stays in the archive (history).
	h2, err := s.PutUnit(uj2, []byte("body v2"), ref(form, "p", "r"))
	if err != nil {
		t.Fatal(err)
	}
	if h2 == h1 {
		t.Error("different content must produce a different hash")
	}
	if prev := payloadPrevHash(t, s, h2); prev != h1 {
		t.Errorf("second payload prev_hash = %q, want %q (lineage)", prev, h1)
	}
	got, ok, err = s.Ref(form)
	if err != nil || !ok {
		t.Fatalf("Ref = %v, %v", ok, err)
	}
	if got.ObjectHash != h2 {
		t.Errorf("reference must move to %s, still at %s", h2, got.ObjectHash)
	}
	if n, err := s.PayloadCount(); err != nil || n != 2 {
		t.Errorf("PayloadCount = %d, %v; want 2 (history retained)", n, err)
	}

	// Third insert: the lineage chain continues (prev = h2).
	h3, err := s.PutUnit(unitJSON(t, "acme", "sto", "x", 1, 3), []byte("body v3"), ref(form, "p", "r"))
	if err != nil {
		t.Fatal(err)
	}
	if prev := payloadPrevHash(t, s, h3); prev != h2 {
		t.Errorf("third payload prev_hash = %q, want %q", prev, h2)
	}
}

// TestPutUnitPayloadBytesNeverUpdated: even when a reference moves to a
// new payload, every existing payload row keeps its original bytes.
// The store has no path to update object_payloads: all writes are
// inserts, and a hash conflict is a no-op.
func TestPutUnitPayloadBytesNeverUpdated(t *testing.T) {
	s := openTest(t)
	form := "acme/sto:x:1"
	uj := unitJSON(t, "acme", "sto", "x", 1, 1)
	h, err := s.PutUnit(uj, []byte("body v1"), ref(form, "p", "r"))
	if err != nil {
		t.Fatal(err)
	}
	// The same payload shared by a second form: the payload row must
	// not be duplicated or rewritten.
	other := "acme/sto:y:1"
	ref2 := ref(other, "p", "r")
	ref2.Type = "sto"
	ref2.ID = "y"
	ref2.Form = other
	h2, err := s.PutUnit(uj, []byte("body v1"), ref2)
	if err != nil {
		t.Fatal(err)
	}
	if h2 != h {
		t.Errorf("shared payload must return the same hash: %s vs %s", h2, h)
	}
	if n := payloadRowCount(t, s, h); n != 1 {
		t.Errorf("payload rows for %s = %d, want 1", h, n)
	}
	unit, content, err := s.Payload(h)
	if err != nil {
		t.Fatal(err)
	}
	if string(unit) != string(uj) || string(content) != "body v1" {
		t.Error("payload bytes changed after a second reference pointed at it")
	}
}

// TestRefsFilteredByProvenance: refs are listed per provenance pair;
// the same basename under a different project never leaks.
func TestRefsFilteredByProvenance(t *testing.T) {
	s := openTest(t)
	refs := []Ref{
		ref("acme/sto:a:1", "proj-a", "repo-a"),
		ref("acme/sto:m:1", "proj-a", "repo-a"),
		ref("acme/sto:z:1", "proj-a", "repo-a"),
	}
	refs[0].ID, refs[1].ID, refs[2].ID = "a", "m", "z"
	for i := range refs {
		refs[i].Form = "acme/sto:" + refs[i].ID + ":1"
	}
	for _, r := range refs {
		if _, err := s.PutUnit(unitJSON(t, "acme", "sto", r.ID, 1, 1), []byte("body"), r); err != nil {
			t.Fatal(err)
		}
	}
	// A second project with the same basename must stay invisible.
	other := ref("acme/sto:q:1", "proj-b", "repo-a")
	other.ID = "q"
	if _, err := s.PutUnit(unitJSON(t, "acme", "sto", "q", 1, 1), []byte("body"), other); err != nil {
		t.Fatal(err)
	}
	got, err := s.Refs("proj-a", "repo-a")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 3 || got[0].Form != "acme/sto:a:1" || got[1].Form != "acme/sto:m:1" || got[2].Form != "acme/sto:z:1" {
		t.Errorf("Refs order = %v, want sorted by form (a, m, z)", got)
	}
	for _, r := range got {
		if r.ProjectID != "proj-a" || r.SourceRepo != "repo-a" {
			t.Errorf("ref %s leaked provenance %s/%s", r.Form, r.ProjectID, r.SourceRepo)
		}
	}
}

// TestRefAbsent: an unknown form has no reference.
func TestRefAbsent(t *testing.T) {
	s := openTest(t)
	got, ok, err := s.Ref("acme/nope:x:1")
	if err != nil {
		t.Fatal(err)
	}
	if ok || got != nil {
		t.Errorf("Ref = %+v, %v; want nil, false", got, ok)
	}
}

// TestPayloadRoundtrip: Payload returns exactly the stored bytes.
func TestPayloadRoundtrip(t *testing.T) {
	s := openTest(t)
	uj := unitJSON(t, "acme", "sto", "x", 1, 1)
	h, err := s.PutUnit(uj, []byte("payload body"), ref("acme/sto:x:1", "p", "r"))
	if err != nil {
		t.Fatal(err)
	}
	unit, content, err := s.Payload(h)
	if err != nil {
		t.Fatal(err)
	}
	if string(unit) != string(uj) || string(content) != "payload body" {
		t.Errorf("Payload roundtrip mismatch: %q / %q", unit, content)
	}
	// Absent hash errors.
	if _, _, err := s.Payload(strings.Repeat("0", 64)); err == nil {
		t.Error("Payload of an absent hash must error")
	}
}

func TestCounts(t *testing.T) {
	s := openTest(t)
	if n, err := s.PayloadCount(); err != nil || n != 0 {
		t.Errorf("PayloadCount = %d, %v; want 0", n, err)
	}
	if n, err := s.RefCount(); err != nil || n != 0 {
		t.Errorf("RefCount = %d, %v; want 0", n, err)
	}
	if n, err := s.AttachmentCount(); err != nil || n != 0 {
		t.Errorf("AttachmentCount = %d, %v; want 0", n, err)
	}

	if _, err := s.PutUnit(unitJSON(t, "acme", "sto", "x", 1, 1), []byte("body"), ref("acme/sto:x:1", "p", "r")); err != nil {
		t.Fatal(err)
	}
	if err := s.UpsertAttachment("p", "r", "docs/diagram.txt", "d1", []byte("data")); err != nil {
		t.Fatal(err)
	}
	if n, err := s.RefCount(); err != nil || n != 1 {
		t.Errorf("RefCount = %d, %v; want 1", n, err)
	}
	if n, err := s.PayloadCount(); err != nil || n != 1 {
		t.Errorf("PayloadCount = %d, %v; want 1", n, err)
	}
	if n, err := s.AttachmentCount(); err != nil || n != 1 {
		t.Errorf("AttachmentCount = %d, %v; want 1", n, err)
	}
}

// TestAllPayloadsSortedByHash: the integrity scan order is
// deterministic.
func TestAllPayloadsSortedByHash(t *testing.T) {
	s := openTest(t)
	// Insert in non-sorted form order; the payload hashes are content-
	// derived so their order is not the insert order.
	for i, id := range []string{"z", "a", "m"} {
		r := ref("acme/sto:"+id+":1", "p", "r")
		r.ID = id
		if _, err := s.PutUnit(unitJSON(t, "acme", "sto", id, 1, i+1), []byte("body "+id), r); err != nil {
			t.Fatal(err)
		}
	}
	rows, err := s.AllPayloads()
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 3 {
		t.Fatalf("AllPayloads = %d rows, want 3", len(rows))
	}
	for i := 1; i < len(rows); i++ {
		if rows[i-1].ObjectHash >= rows[i].ObjectHash {
			t.Errorf("payloads not sorted by hash at %d", i)
		}
	}
}

// TestForeignKeysEnforced: a reference pointing at a missing payload is
// rejected (foreign_keys pragma on).
func TestForeignKeysEnforced(t *testing.T) {
	s := openTest(t)
	if _, err := s.db.Exec(`INSERT INTO object_refs (form, object_hash, project_id, source_repo, namespace, type, id, instance_version, revision, dimension, domain, phase, updated_at)
		VALUES ('ghost/x:1:1', '` + strings.Repeat("0", 64) + `', 'p', 'r', 'ghost', 'x', '1', 1, 1, '', '', '', '2026-08-07T00:00:00Z')`); err == nil {
		t.Error("foreign key violation must be rejected")
	}
}

// --- attachments (unchanged behavior) ---

func TestAttachmentsUpsert(t *testing.T) {
	s := openTest(t)
	if err := s.UpsertAttachment("p", "r", "docs/a.txt", "digest1", []byte("one")); err != nil {
		t.Fatal(err)
	}
	if err := s.UpsertAttachment("p", "r", "docs/a.txt", "digest2", []byte("two")); err != nil {
		t.Fatal(err)
	}
	if err := s.UpsertAttachment("p", "other", "docs/a.txt", "digest9", []byte("nine")); err != nil {
		t.Fatal(err)
	}
	if n, err := s.AttachmentCount(); err != nil || n != 2 {
		t.Errorf("AttachmentCount = %d, %v; want 2", n, err)
	}
	atts, err := s.Attachments("p", "r")
	if err != nil {
		t.Fatal(err)
	}
	if len(atts) != 1 || atts[0].Digest != "digest2" || string(atts[0].Data) != "two" {
		t.Errorf("Attachments(p, r) = %+v, want the updated row only", atts)
	}
}

// --- sync log (unchanged behavior) ---

func TestSyncLog(t *testing.T) {
	s := openTest(t)
	entries := []SyncEntry{
		{ProjectID: "p", Repo: "r", Direction: "pull", SnapshotDigest: "d1", Units: 2, At: "2026-08-01T00:00:00Z"},
		{ProjectID: "p", Repo: "r", Direction: "push", SnapshotDigest: "d2", Units: 2, At: "2026-08-01T00:01:00Z"},
		{ProjectID: "p", Repo: "r", Direction: "pull", SnapshotDigest: "d3", Units: 0, At: "2026-08-01T00:02:00Z"},
		{ProjectID: "p", Repo: "other", Direction: "pull", SnapshotDigest: "d9", Units: 1, At: "2026-08-01T00:03:00Z"},
	}
	for _, e := range entries {
		if err := s.RecordSync(e); err != nil {
			t.Fatal(err)
		}
	}
	recent, err := s.RecentSyncs("p", "r", 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(recent) != 2 || recent[0].Seq != 3 || recent[1].Seq != 2 {
		t.Errorf("RecentSyncs = %+v, want [seq3 seq2] (desc)", recent)
	}
	digest, ok, err := s.LastPullDigest("p", "r")
	if err != nil {
		t.Fatal(err)
	}
	if !ok || digest != "d3" {
		t.Errorf("LastPullDigest = %q, %v; want d3, true", digest, ok)
	}
	if _, ok, err := s.LastPullDigest("p", "other2"); err != nil || ok {
		t.Errorf("LastPullDigest for unknown repo = %v, %v; want false", ok, err)
	}
}

// --- helpers ---

// payloadPrevHash reads the prev_hash of one payload row.
func payloadPrevHash(t *testing.T, s *Store, hash string) string {
	t.Helper()
	var prev string
	if err := s.db.QueryRow(`SELECT prev_hash FROM object_payloads WHERE object_hash = ?`, hash).Scan(&prev); err != nil {
		t.Fatalf("read prev_hash of %s: %v", hash, err)
	}
	return prev
}

// payloadRowCount counts payload rows with one hash.
func payloadRowCount(t *testing.T, s *Store, hash string) int {
	t.Helper()
	var n int
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM object_payloads WHERE object_hash = ?`, hash).Scan(&n); err != nil {
		t.Fatal(err)
	}
	return n
}

// hasTableRaw reports whether a table exists in the store's database.
func hasTableRaw(t *testing.T, s *Store, name string) bool {
	t.Helper()
	var n int
	err := s.db.QueryRow(`SELECT COUNT(*) FROM sqlite_master WHERE type = 'table' AND name = ?`, name).Scan(&n)
	if err != nil {
		t.Fatal(err)
	}
	return n > 0
}

// --- schema v1 -> v2 migration ---

// createV1DB builds a v1-schema database (the experimental v0.2.0
// layout) at dir/eka.db with the given objects, relationships and
// change-log rows, and records schema_version = 1.
func createV1DB(t *testing.T, dir string, objs []map[string]any, rels [][3]string, logs [][6]string) {
	t.Helper()
	db, err := sql.Open("sqlite", "file:"+filepath.Join(dir, "eka.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	statements := []string{
		`CREATE TABLE meta (key TEXT PRIMARY KEY, value TEXT NOT NULL)`,
		`CREATE TABLE projects (id TEXT PRIMARY KEY, name TEXT NOT NULL, created TEXT NOT NULL)`,
		`CREATE TABLE repos (project_id TEXT NOT NULL REFERENCES projects(id), name TEXT NOT NULL, path TEXT NOT NULL, created TEXT NOT NULL, PRIMARY KEY (project_id, name))`,
		`CREATE TABLE objects (
			form TEXT PRIMARY KEY, project_id TEXT NOT NULL, namespace TEXT NOT NULL,
			type TEXT NOT NULL, id TEXT NOT NULL, instance_version INTEGER NOT NULL,
			revision INTEGER NOT NULL, author TEXT NOT NULL, created TEXT NOT NULL,
			updated TEXT NOT NULL, content_representation TEXT NOT NULL, content BLOB,
			state_content TEXT NOT NULL, state_execution TEXT NOT NULL,
			state_planning TEXT NOT NULL, state_container TEXT NOT NULL,
			state_existence TEXT NOT NULL, phase TEXT NOT NULL, dimension TEXT NOT NULL,
			domain TEXT NOT NULL, source_repo TEXT NOT NULL, digest TEXT NOT NULL,
			dimensions_secondary TEXT NOT NULL)`,
		`CREATE TABLE relationships (form TEXT NOT NULL REFERENCES objects(form), rel_type TEXT NOT NULL, target TEXT NOT NULL, PRIMARY KEY (form, rel_type, target))`,
		`CREATE TABLE change_log (form TEXT NOT NULL REFERENCES objects(form), seq INTEGER NOT NULL, date TEXT NOT NULL, domain TEXT NOT NULL, from_val TEXT NOT NULL, to_val TEXT NOT NULL, by TEXT NOT NULL, PRIMARY KEY (form, seq))`,
		`CREATE TABLE attachments (project_id TEXT NOT NULL, source_repo TEXT NOT NULL, id TEXT NOT NULL, digest TEXT NOT NULL, data BLOB NOT NULL, PRIMARY KEY (project_id, source_repo, id))`,
		`CREATE TABLE sync_log (seq INTEGER PRIMARY KEY AUTOINCREMENT, project_id TEXT NOT NULL, repo TEXT NOT NULL, direction TEXT NOT NULL, snapshot_digest TEXT NOT NULL, units INTEGER NOT NULL, at TEXT NOT NULL)`,
	}
	for _, stmt := range statements {
		if _, err := db.Exec(stmt); err != nil {
			t.Fatalf("create v1 table: %v", err)
		}
	}
	if _, err := db.Exec(`INSERT INTO meta (key, value) VALUES ('schema_version', '1')`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`INSERT INTO projects (id, name, created) VALUES ('p', 'p', '2026-08-01')`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`INSERT INTO repos (project_id, name, path, created) VALUES ('p', 'r', '/tmp/r', '2026-08-01')`); err != nil {
		t.Fatal(err)
	}
	for _, o := range objs {
		cols := []string{"form", "project_id", "namespace", "type", "id", "instance_version", "revision", "author", "created", "updated", "content_representation", "content", "state_content", "state_execution", "state_planning", "state_container", "state_existence", "phase", "dimension", "domain", "source_repo", "digest", "dimensions_secondary"}
		vals := make([]any, 0, len(cols))
		for _, c := range cols {
			v, ok := o[c]
			if !ok {
				v = ""
			}
			vals = append(vals, v)
		}
		ph := strings.TrimSuffix(strings.Repeat("?,", len(cols)), ",")
		if _, err := db.Exec(`INSERT INTO objects (`+strings.Join(cols, ",")+`) VALUES (`+ph+`)`, vals...); err != nil {
			t.Fatalf("insert v1 object: %v", err)
		}
	}
	for _, r := range rels {
		if _, err := db.Exec(`INSERT INTO relationships (form, rel_type, target) VALUES (?, ?, ?)`, r[0], r[1], r[2]); err != nil {
			t.Fatal(err)
		}
	}
	for i, e := range logs {
		if _, err := db.Exec(`INSERT INTO change_log (form, seq, date, domain, from_val, to_val, by) VALUES (?, ?, ?, ?, ?, ?, ?)`,
			e[0], i+1, e[1], e[2], e[3], e[4], e[5]); err != nil {
			t.Fatal(err)
		}
	}
}

// TestMigrationV1ToV2: an experimental v0.2.0 database (v1 schema)
// migrates on Open: every object becomes an immutable payload with the
// RECOMPUTED content hash, change-log transitions survive inside the
// unit.json payload, the v1 tables are dropped, and the result passes
// VerifyIntegrity.
func TestMigrationV1ToV2(t *testing.T) {
	dir := t.TempDir()
	createV1DB(t, dir, []map[string]any{
		{
			"form": "acme/adr:001:1", "project_id": "p", "namespace": "acme", "type": "adr", "id": "001",
			"instance_version": 1, "revision": 2, "author": "Eng", "created": "2026-08-01", "updated": "2026-08-02",
			"content_representation": "eka/structured-text/1", "content": []byte("# ADR 001\n"),
			"state_content": "stable", "state_execution": "", "state_planning": "", "state_container": "",
			"state_existence": "active", "phase": "", "dimension": "decisions", "domain": "Architecture",
			"source_repo": "r", "digest": "trust-me-not", "dimensions_secondary": `["architecture"]`,
		},
		{
			"form": "acme/sto:login-email:1", "project_id": "p", "namespace": "acme", "type": "sto", "id": "login-email",
			"instance_version": 1, "revision": 1, "author": "", "created": "2026-08-01", "updated": "2026-08-01",
			"content_representation": "eka/structured-text/1", "content": []byte("# Store\n"),
			"state_content": "draft", "state_execution": "", "state_planning": "", "state_container": "",
			"state_existence": "active", "phase": "wave-1", "dimension": "", "domain": "Architecture",
			"source_repo": "r", "digest": "also-untrusted", "dimensions_secondary": "",
		},
	}, [][3]string{
		{"acme/adr:001:1", "depends-on", "acme/sto:login-email:1"},
		{"acme/adr:001:1", "amends", "acme/adr:000-before:1"},
	}, [][6]string{
		{"acme/adr:001:1", "2026-08-01", "existence-state", "-", "active", "Eng"},
		{"acme/adr:001:1", "2026-08-02", "content-state", "draft", "stable", "Eng"},
	})

	s, err := Open(dir)
	if err != nil {
		t.Fatalf("Open on v1 database: %v", err)
	}
	defer s.Close()

	v, err := s.SchemaVersion()
	if err != nil {
		t.Fatal(err)
	}
	if v != 2 {
		t.Errorf("schema version after migration = %d, want 2", v)
	}
	// v1 tables are gone.
	for _, table := range []string{"objects", "relationships", "change_log"} {
		if hasTableRaw(t, s, table) {
			t.Errorf("v1 table %s must be dropped by the migration", table)
		}
	}

	// References: two forms, derived columns from the migrated rows.
	refs, err := s.Refs("p", "r")
	if err != nil {
		t.Fatal(err)
	}
	if len(refs) != 2 {
		t.Fatalf("refs after migration = %d, want 2", len(refs))
	}
	adr := refs[0]
	if adr.Form != "acme/adr:001:1" || adr.Namespace != "acme" || adr.Type != "adr" ||
		adr.ID != "001" || adr.InstanceVersion != 1 || adr.Revision != 2 ||
		adr.Dimension != "decisions" || adr.Domain != "Architecture" || adr.Phase != "" {
		t.Errorf("migrated adr ref = %+v", adr)
	}

	// The object hash is RECOMPUTED from the reconstructed bytes — the
	// v1 digest column is never trusted.
	unitJSON, content, err := s.Payload(adr.ObjectHash)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "# ADR 001\n" {
		t.Errorf("migrated content = %q", content)
	}
	u, err := exchange.DecodeUnit(unitJSON, content)
	if err != nil {
		t.Fatalf("migrated unit.json must decode: %v", err)
	}
	// Change-log transitions survive inside the payload, in occurrence
	// order.
	if len(u.ChangeLog) != 2 || u.ChangeLog[0].Domain != "existence-state" ||
		u.ChangeLog[1].Domain != "content-state" || u.ChangeLog[1].To != "stable" {
		t.Errorf("migrated change log = %+v, want both transitions in order", u.ChangeLog)
	}
	// Relationships survive, ordered by (rel_type, target).
	if len(u.Relationships) != 2 || u.Relationships[0].Type != "amends" ||
		u.Relationships[1].Type != "depends-on" {
		t.Errorf("migrated relationships = %+v, want sorted (amends, depends-on)", u.Relationships)
	}
	if u.Revision != 2 || u.StateVector.ContentState != "stable" ||
		len(u.Classification.DimensionsSecondary) != 1 ||
		u.Classification.DimensionsSecondary[0] != "architecture" {
		t.Errorf("migrated unit fields mismatch: %+v", u)
	}

	// Hash = SHA-256(marshal(unit) || content), recomputed in the test
	// from the reconstructed unit — agreement proves the migration
	// serializes the same canonical bytes.
	expect := sha256.Sum256(append(unitJSON, content...))
	if adr.ObjectHash != hex.EncodeToString(expect[:]) {
		t.Errorf("migrated hash %s != SHA-256(unit.json || content) %s", adr.ObjectHash, hex.EncodeToString(expect[:]))
	}
	// prev_hash = "" for every migrated payload (v1 had no payload
	// history; the change log is inside the payload).
	if prev := payloadPrevHash(t, s, adr.ObjectHash); prev != "" {
		t.Errorf("migrated payload prev_hash = %q, want \"\"", prev)
	}

	// The migrated database is fully consistent.
	report, err := s.VerifyIntegrity()
	if err != nil {
		t.Fatal(err)
	}
	if len(report.Violations) != 0 {
		t.Errorf("migrated database must pass integrity: %+v", report.Violations)
	}
}

// TestMigrationFreshDBIsV2Direct: opening a database that has no v1
// tables (even a stray meta table) creates the v2 schema directly.
func TestMigrationFreshDBIsV2Direct(t *testing.T) {
	dir := t.TempDir()
	s, err := Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	if !hasTableRaw(t, s, "object_payloads") || !hasTableRaw(t, s, "object_refs") {
		t.Error("fresh database must create the v2 tables")
	}
	if hasTableRaw(t, s, "objects") {
		t.Error("fresh database must not create v1 tables")
	}
}

// --- integrity verification ---

// seedIntegrityStore builds a store with one payload + ref (with
// classification), one attachment and one registered repo/project, all
// consistent.
func seedIntegrityStore(t *testing.T) *Store {
	t.Helper()
	s := openTest(t)
	if _, err := s.PutUnit(unitJSON(t, "acme", "sto", "x", 1, 1), []byte("body"), ref("acme/sto:x:1", "p", "r")); err != nil {
		t.Fatal(err)
	}
	digest := sha256.Sum256([]byte("diagram"))
	if err := s.UpsertAttachment("p", "r", "docs/diagram.txt", hex.EncodeToString(digest[:]), []byte("diagram")); err != nil {
		t.Fatal(err)
	}
	if _, err := s.db.Exec(`INSERT INTO projects (id, name, created) VALUES ('p', 'p', '2026-08-01')`); err != nil {
		t.Fatal(err)
	}
	if _, err := s.db.Exec(`INSERT INTO repos (project_id, name, path, created) VALUES ('p', 'r', '/tmp/r', '2026-08-01')`); err != nil {
		t.Fatal(err)
	}
	return s
}

func TestVerifyIntegrityClean(t *testing.T) {
	s := seedIntegrityStore(t)
	report, err := s.VerifyIntegrity()
	if err != nil {
		t.Fatal(err)
	}
	if report.PayloadsChecked != 1 || report.RefsChecked != 1 || report.AttachmentsChecked != 1 {
		t.Errorf("checked counts = %d/%d/%d, want 1/1/1",
			report.PayloadsChecked, report.RefsChecked, report.AttachmentsChecked)
	}
	if report.OrphanPayloads != 0 {
		t.Errorf("orphans = %d, want 0", report.OrphanPayloads)
	}
	if len(report.Violations) != 0 {
		t.Errorf("clean store must have no violations: %+v", report.Violations)
	}
}

// TestVerifyIntegrityPayloadHash: tampering with payload bytes is
// detected as a payload-hash violation (and, when the tampered bytes
// are no longer a valid unit, as a payload-decode violation too).
func TestVerifyIntegrityPayloadHash(t *testing.T) {
	s := seedIntegrityStore(t)
	rows, err := s.AllPayloads()
	if err != nil {
		t.Fatal(err)
	}
	hash := rows[0].ObjectHash
	if _, err := s.db.Exec(`UPDATE object_payloads SET content = ? WHERE object_hash = ?`, []byte("tampered"), hash); err != nil {
		t.Fatal(err)
	}
	report, err := s.VerifyIntegrity()
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, v := range report.Violations {
		if v.Kind == "payload-hash" && v.Subject == hash {
			found = true
		}
	}
	if !found {
		t.Errorf("payload-hash violation missing: %+v", report.Violations)
	}
}

// TestVerifyIntegrityReferenceIndex: an index column diverging from the
// payload is detected as a reference-index violation.
func TestVerifyIntegrityReferenceIndex(t *testing.T) {
	s := seedIntegrityStore(t)
	if _, err := s.db.Exec(`UPDATE object_refs SET namespace = 'evil' WHERE form = 'acme/sto:x:1'`); err != nil {
		t.Fatal(err)
	}
	report, err := s.VerifyIntegrity()
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, v := range report.Violations {
		if v.Kind == "reference-index" && v.Subject == "acme/sto:x:1" &&
			strings.Contains(v.Detail, "namespace") {
			found = true
		}
	}
	if !found {
		t.Errorf("reference-index violation missing: %+v", report.Violations)
	}
}

// TestVerifyIntegrityReferenceTarget: a reference pointing at a
// nonexistent payload is detected as a reference-target violation
// (crafted with foreign keys off — the pragma is the trust boundary the
// check exists for).
func TestVerifyIntegrityReferenceTarget(t *testing.T) {
	s := seedIntegrityStore(t)
	if _, err := s.db.Exec(`PRAGMA foreign_keys = OFF`); err != nil {
		t.Fatal(err)
	}
	if _, err := s.db.Exec(`INSERT INTO object_refs (form, object_hash, project_id, source_repo, namespace, type, id, instance_version, revision, dimension, domain, phase, updated_at)
		VALUES ('ghost/sto:q:1', '` + strings.Repeat("0", 64) + `', 'p', 'r', 'ghost', 'sto', 'q', 1, 1, '', '', '', '2026-08-07T00:00:00Z')`); err != nil {
		t.Fatal(err)
	}
	report, err := s.VerifyIntegrity()
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, v := range report.Violations {
		if v.Kind == "reference-target" && v.Subject == "ghost/sto:q:1" {
			found = true
		}
	}
	if !found {
		t.Errorf("reference-target violation missing: %+v", report.Violations)
	}
}

// TestVerifyIntegrityOrphansAreNotViolations: deleting a reference
// turns its payload into an orphan (retained history archive) — counted
// but never a violation.
func TestVerifyIntegrityOrphansAreNotViolations(t *testing.T) {
	s := seedIntegrityStore(t)
	if _, err := s.db.Exec(`DELETE FROM object_refs WHERE form = 'acme/sto:x:1'`); err != nil {
		t.Fatal(err)
	}
	report, err := s.VerifyIntegrity()
	if err != nil {
		t.Fatal(err)
	}
	if report.OrphanPayloads != 1 {
		t.Errorf("orphans = %d, want 1", report.OrphanPayloads)
	}
	for _, v := range report.Violations {
		t.Errorf("orphans must not produce violations, got %+v", v)
	}
}

// TestVerifyIntegrityAttachmentHash: tampering with attachment bytes is
// detected as an attachment-hash violation.
func TestVerifyIntegrityAttachmentHash(t *testing.T) {
	s := seedIntegrityStore(t)
	if _, err := s.db.Exec(`UPDATE attachments SET data = ? WHERE id = 'docs/diagram.txt'`, []byte("changed")); err != nil {
		t.Fatal(err)
	}
	report, err := s.VerifyIntegrity()
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, v := range report.Violations {
		if v.Kind == "attachment-hash" && v.Subject == "docs/diagram.txt" {
			found = true
		}
	}
	if !found {
		t.Errorf("attachment-hash violation missing: %+v", report.Violations)
	}
}

// TestVerifyIntegrityRegistry: a repo referencing a missing project is
// detected as a registry violation (crafted with foreign keys off).
func TestVerifyIntegrityRegistry(t *testing.T) {
	s := seedIntegrityStore(t)
	if _, err := s.db.Exec(`PRAGMA foreign_keys = OFF`); err != nil {
		t.Fatal(err)
	}
	if _, err := s.db.Exec(`DELETE FROM projects WHERE id = 'p'`); err != nil {
		t.Fatal(err)
	}
	report, err := s.VerifyIntegrity()
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, v := range report.Violations {
		if v.Kind == "registry" && v.Subject == "p/r" {
			found = true
		}
	}
	if !found {
		t.Errorf("registry violation missing: %+v", report.Violations)
	}
}

// TestVerifyIntegrityDeterministicOrder: violations are sorted by
// (Kind, Subject), so identical corruption produces identical reports.
func TestVerifyIntegrityDeterministicOrder(t *testing.T) {
	s := seedIntegrityStore(t)
	rows, err := s.AllPayloads()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.db.Exec(`UPDATE object_payloads SET content = ? WHERE object_hash = ?`, []byte("tampered"), rows[0].ObjectHash); err != nil {
		t.Fatal(err)
	}
	if _, err := s.db.Exec(`UPDATE object_refs SET dimension = 'wrong' WHERE form = 'acme/sto:x:1'`); err != nil {
		t.Fatal(err)
	}
	run := func() []IntegrityViolation {
		report, err := s.VerifyIntegrity()
		if err != nil {
			t.Fatal(err)
		}
		return report.Violations
	}
	first := run()
	second := run()
	if len(first) != len(second) {
		t.Fatalf("violation count differs between runs: %d vs %d", len(first), len(second))
	}
	for i := range first {
		if first[i] != second[i] {
			t.Errorf("violations differ between runs at %d: %+v vs %+v", i, first[i], second[i])
		}
	}
	for i := 1; i < len(first); i++ {
		if first[i-1].Kind > first[i].Kind ||
			(first[i-1].Kind == first[i].Kind && first[i-1].Subject > first[i].Subject) {
			t.Errorf("violations not sorted by (Kind, Subject) at %d: %+v", i, first)
		}
	}
}

// TestOpenRefusesNewerSchema: a database recorded at a schema version
// newer than this implementation is refused (downgrade protection).
func TestOpenRefusesNewerSchema(t *testing.T) {
	dir := t.TempDir()
	s, err := Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.SetMeta("schema_version", "99"); err != nil {
		t.Fatal(err)
	}
	s.Close()
	if _, err := Open(dir); err == nil {
		t.Fatal("a database newer than this implementation must be refused")
	}
}

// TestVerifyIntegrityPayloadDecode: tampering with the unit.json bytes
// is detected as BOTH a payload-hash violation (bytes changed) and a
// payload-decode violation (bytes no longer strict-decodable) — the
// payload-decode kind is asserted here (the content-tamper test above
// covers payload-hash alone).
func TestVerifyIntegrityPayloadDecode(t *testing.T) {
	s := seedIntegrityStore(t)
	rows, err := s.AllPayloads()
	if err != nil {
		t.Fatal(err)
	}
	hash := rows[0].ObjectHash
	if _, err := s.db.Exec(`UPDATE object_payloads SET unit_json = ? WHERE object_hash = ?`, []byte("not json"), hash); err != nil {
		t.Fatal(err)
	}
	report, err := s.VerifyIntegrity()
	if err != nil {
		t.Fatal(err)
	}
	kinds := map[string]bool{}
	for _, v := range report.Violations {
		if v.Subject == hash {
			kinds[v.Kind] = true
		}
	}
	if !kinds["payload-hash"] {
		t.Errorf("payload-hash violation missing for tampered unit_json: %+v", report.Violations)
	}
	if !kinds["payload-decode"] {
		t.Errorf("payload-decode violation missing for tampered unit_json: %+v", report.Violations)
	}
}
