package machine

import (
	"encoding/json"
	"strconv"
	"strings"
	"testing"

	"github.com/maleolabs/engineering-knowledge-architecture/exchange"
)

// goldenUnit is the hand-built CKO the golden tests marshal: one
// fully-populated adr unit (states, classification, relationships,
// change log, metadata, digest).
func goldenUnit() *exchange.Unit {
	return &exchange.Unit{
		Identity: exchange.Identity{
			Namespace:       "feather",
			Type:            "adr",
			ID:              "001-serialization",
			InstanceVersion: 1,
		},
		CanonicalIdentityForm: "feather/adr:001-serialization:1",
		Revision:              2,
		Author:                "Engineering",
		Created:               "2026-08-05",
		Updated:               "2026-08-06",
		StateVector: exchange.StateVector{
			ContentState:   "accepted",
			ExistenceState: "active",
		},
		ChangeLog: []exchange.ChangeLogEntry{
			{Date: "2026-08-05", Domain: "content-state", From: "proposed", To: "accepted", By: "Engineering"},
		},
		Relationships: []exchange.Relationship{
			{Type: "depends-on", Target: "feather/sto:publish-post:1"},
		},
		Classification: exchange.Classification{
			Dimension: "decisions",
			Domain:    "Architecture",
		},
		Content: exchange.ContentRef{Representation: "eka/structured-text/1", File: "units/feather/adr-001-serialization-v1/content.md"},
		ContentPayload: []byte("# ADR-001 — Login serialization\n\n" +
			"## Context\n\nContext body.\n"),
		Digest: "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
	}
}

// TestGoldenFieldOrder pins the exact serialized bytes of a fully
// populated document: fixed field order (struct declaration order),
// stable schema string, two-space indent, single trailing newline.
func TestGoldenFieldOrder(t *testing.T) {
	want := `{
  "schema": "eka-cko-v1",
  "identity": {
    "namespace": "feather",
    "type": "adr",
    "id": "001-serialization",
    "instance_version": 1
  },
  "canonical_form": "feather/adr:001-serialization:1",
  "engineering_domain": "Architecture",
  "stratum": 2,
  "revision": 2,
  "author": "Engineering",
  "created": "2026-08-05",
  "updated": "2026-08-06",
  "state_vector": {
    "content-state": "accepted",
    "existence-state": "active"
  },
  "classification": {
    "dimension": "decisions",
    "domain": "Architecture"
  },
  "relationships": [
    {
      "type": "depends-on",
      "target": "feather/sto:publish-post:1"
    }
  ],
  "change_log": [
    {
      "date": "2026-08-05",
      "domain": "content-state",
      "from": "proposed",
      "to": "accepted",
      "by": "Engineering"
    }
  ],
  "content": {
    "representation": "eka/structured-text/1",
    "text": "# ADR-001 — Login serialization\n\n## Context\n\nContext body.\n"
  },
  "object_hash": "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
}
`
	got, err := MarshalUnit(goldenUnit())
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != want {
		t.Errorf("serialized document differs:\ngot:\n%s\nwant:\n%s", got, want)
	}
}

// TestOmitemptySemantics: a unit without states, classification,
// relationships, change log or metadata serializes with those keys
// absent. The state_vector block is a value struct (never omitted by
// encoding/json) and serializes as an empty object — the RSF unit.json
// empty-vector behavior of §5.1.1.
func TestOmitemptySemantics(t *testing.T) {
	u := &exchange.Unit{
		Identity: exchange.Identity{
			Namespace:       "acme",
			Type:            "run",
			ID:              "backup",
			InstanceVersion: 1,
		},
		CanonicalIdentityForm: "acme/run:backup:1",
		Content:               exchange.ContentRef{Representation: "eka/structured-text/1"},
		ContentPayload:        []byte("# Runbook\n"),
	}
	doc, err := NewDocument(u)
	if err != nil {
		t.Fatal(err)
	}
	if doc.Classification != nil {
		t.Errorf("classification must be nil for a unit without classification, got %+v", doc.Classification)
	}
	if len(doc.Relationships) != 0 || doc.Relationships != nil {
		t.Errorf("relationships must stay nil, got %v", doc.Relationships)
	}
	if len(doc.ChangeLog) != 0 || doc.ChangeLog != nil {
		t.Errorf("change_log must stay nil, got %v", doc.ChangeLog)
	}
	got, err := doc.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	s := string(got)
	for _, absent := range []string{
		`"revision"`, `"author"`, `"created"`, `"updated"`, `"phase"`,
		`"classification"`, `"relationships"`, `"change_log"`,
	} {
		if strings.Contains(s, absent) {
			t.Errorf("serialized document must not carry %s:\n%s", absent, s)
		}
	}
	if !strings.Contains(s, `"state_vector": {}`) {
		t.Errorf("state_vector must serialize as an empty object (RSF §5.1.1 empty-vector behavior):\n%s", s)
	}
	// The derived engineering domain: Classification.Domain is empty,
	// so it comes from the type token ("run" -> Operations, stratum 5).
	if !strings.Contains(s, `"engineering_domain": "Operations"`) || !strings.Contains(s, `"stratum": 5`) {
		t.Errorf("domain must be derived from the type token:\n%s", s)
	}
	// object_hash stays as-is for hand-built units ("" kept as-is).
	if !strings.Contains(s, `"object_hash": ""`) {
		t.Errorf("object_hash must be kept as-is (empty for hand-built units):\n%s", s)
	}
}

// TestContentEscaping: content with quotes, newlines and non-ASCII
// characters must serialize as valid JSON and round-trip losslessly.
func TestContentEscaping(t *testing.T) {
	body := "He said \"quoted\", tab\there,\nand unicode: café — 日本語 ✓\n"
	u := &exchange.Unit{
		Identity: exchange.Identity{
			Namespace: "acme", Type: "spec", ID: "x", InstanceVersion: 1,
		},
		CanonicalIdentityForm: "acme/spec:x:1",
		Content:               exchange.ContentRef{Representation: "eka/structured-text/1"},
		ContentPayload:        []byte(body),
	}
	out, err := MarshalUnit(u)
	if err != nil {
		t.Fatal(err)
	}
	if !json.Valid(out) {
		t.Fatalf("output must be valid JSON:\n%s", out)
	}
	var doc Document
	if err := json.Unmarshal(out, &doc); err != nil {
		t.Fatalf("round-trip unmarshal failed: %v", err)
	}
	if doc.Content.Text != body {
		t.Errorf("round-trip text = %q, want %q", doc.Content.Text, body)
	}
}

// TestUnknownTypeToken: a unit whose type token has no home domain and
// no declared Classification.Domain is a deterministic error.
func TestUnknownTypeToken(t *testing.T) {
	u := &exchange.Unit{
		Identity: exchange.Identity{
			Namespace: "acme", Type: "wat", ID: "x", InstanceVersion: 1,
		},
		CanonicalIdentityForm: "acme/wat:x:1",
	}
	_, err := NewDocument(u)
	if err == nil {
		t.Fatal("NewDocument must fail for an unknown artifact type")
	}
	if got := err.Error(); got != `machine: unknown artifact type "wat"` {
		t.Errorf("error = %q, want deterministic message", got)
	}
}

// TestCollectionEmpty: an empty unit list yields a collection with
// count 0 and an empty unit list (never null), with the stable
// collection shape.
func TestCollectionEmpty(t *testing.T) {
	c, err := NewCollection("Execution", nil)
	if err != nil {
		t.Fatal(err)
	}
	if c.Count != 0 || len(c.Units) != 0 {
		t.Errorf("count = %d, units = %d, want 0 and 0", c.Count, len(c.Units))
	}
	got, err := c.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	want := `{
  "schema": "eka-cko-v1",
  "collection": "domain",
  "domain": "Execution",
  "count": 0,
  "units": []
}
`
	if string(got) != want {
		t.Errorf("empty collection differs:\ngot:\n%s\nwant:\n%s", got, want)
	}
}

// TestCollectionSortedByCanonicalForm: the collection sorts its units
// by canonical form regardless of the input order (determinism
// contract). Documents carry no ordering metadata beyond the pinned
// field order, so the assertion reads the serialized units array.
func TestCollectionSortedByCanonicalForm(t *testing.T) {
	mk := func(ns, typ, id string, v int) *exchange.Unit {
		return &exchange.Unit{
			Identity:              exchange.Identity{Namespace: ns, Type: typ, ID: id, InstanceVersion: v},
			CanonicalIdentityForm: ns + "/" + typ + ":" + id + ":" + itoa(v),
			Content:               exchange.ContentRef{Representation: "eka/structured-text/1"},
			ContentPayload:        []byte("# " + id + "\n"),
		}
	}
	units := []*exchange.Unit{
		mk("zeta", "sto", "late", 1),
		mk("alpha", "adr", "first", 1),
		mk("alpha", "adr", "first", 2),
		mk("alpha", "sto", "zz", 1),
	}
	c, err := NewCollection("Execution", units)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{
		"alpha/adr:first:1",
		"alpha/adr:first:2",
		"alpha/sto:zz:1",
		"zeta/sto:late:1",
	}
	for i, w := range want {
		if c.Units[i].CanonicalForm != w {
			t.Errorf("units[%d] = %q, want %q", i, c.Units[i].CanonicalForm, w)
		}
	}
}

// TestDeterminism: two marshals of the same input are byte-identical.
func TestDeterminism(t *testing.T) {
	a, err := MarshalUnit(goldenUnit())
	if err != nil {
		t.Fatal(err)
	}
	b, err := MarshalUnit(goldenUnit())
	if err != nil {
		t.Fatal(err)
	}
	if string(a) != string(b) {
		t.Error("two marshals of the same unit must be byte-identical")
	}
	c, err := NewCollection("Architecture", []*exchange.Unit{goldenUnit(), goldenUnit()})
	if err != nil {
		t.Fatal(err)
	}
	ca, err := c.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	cb, err := c.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if string(ca) != string(cb) {
		t.Error("two collection marshals must be byte-identical")
	}
}

func itoa(v int) string {
	return strconv.Itoa(v)
}
