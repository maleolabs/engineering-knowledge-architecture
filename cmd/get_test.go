package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/maleolabs/engineering-knowledge-architecture/runtime"
)

// seedGetRepo seeds a fresh workspace (EKA_HOME) with a copy of the
// view "valid" fixture through the Runtime (the store-backed setup of
// the get path: runtime.Ensure + Authoring.Sync), optionally adding
// extra authoring docs before the sync, and chdirs into the repo copy.
// Returns the repo path.
func seedGetRepo(t *testing.T, extra func(repo string)) string {
	t.Helper()
	t.Setenv("EKA_HOME", t.TempDir())
	repo := copyFixture(t, filepath.Join("..", "view", "testdata", "valid"))
	if extra != nil {
		extra(repo)
	}
	r, err := runtime.Ensure()
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	if _, err := runtime.Authoring.Sync(r, repo, runtime.SyncOptions{Pull: true, Push: true}); err != nil {
		t.Fatal(err)
	}
	chdirInto(t, repo)
	return repo
}

// getDocument parses the stdout of a successful get run as a machine
// document.
func getDocument(t *testing.T, args ...string) map[string]any {
	t.Helper()
	code, out, errText := runIn(args)
	if code != 0 {
		t.Fatalf("exit = %d, want 0\nstdout: %s\nstderr: %s", code, out, errText)
	}
	var doc map[string]any
	if err := json.Unmarshal([]byte(out), &doc); err != nil {
		t.Fatalf("stdout must be a single parseable JSON document: %v\nstdout: %q", err, out)
	}
	return doc
}

// TestGetIdentityCanonicalFormGolden: identity lookup by the RSF
// canonical form produces the exact pinned machine JSON (byte-compare)
// and nothing else on stdout.
func TestGetIdentityCanonicalFormGolden(t *testing.T) {
	seedGetRepo(t, nil)
	code, out, errText := runIn([]string{"get", "eka-view-fixture/adr:001-login-serialization:1"})
	if code != 0 {
		t.Fatalf("exit = %d, want 0\nstdout: %s\nstderr: %s", code, out, errText)
	}
	if errText != "" {
		t.Errorf("stderr must be empty on success, got %q", errText)
	}
	want := `{
  "schema": "eka-cko-v1",
  "identity": {
    "namespace": "eka-view-fixture",
    "type": "adr",
    "id": "001-login-serialization",
    "instance_version": 1
  },
  "canonical_form": "eka-view-fixture/adr:001-login-serialization:1",
  "engineering_domain": "Architecture",
  "stratum": 2,
  "revision": 1,
  "author": "Engineering Architecture",
  "created": "2026-08-05",
  "updated": "2026-08-05",
  "state_vector": {
    "content-state": "accepted",
    "existence-state": "active"
  },
  "classification": {
    "dimension": "decisions",
    "domain": "Architecture"
  },
  "change_log": [
    {
      "date": "2026-08-05",
      "domain": "existence-state",
      "from": "-",
      "to": "active",
      "by": "Engineering Architecture"
    },
    {
      "date": "2026-08-05",
      "domain": "content-state",
      "from": "proposed",
      "to": "accepted",
      "by": "Engineering Architecture"
    }
  ],
  "content": {
    "representation": "eka/structured-text/1",
    "text": "\n# ADR-001 — Login serialization\n\n## Context\n\nContext body.\n\n## Decision\n\nDecision body.\n\n## Consequences\n\nConsequences body.\n\n## Alternatives Considered\n\nAlternatives body.\n"
  },
  "object_hash": "d7a3985a0eedb95d065257369856652ef25b1e59d814a142485d2a9f7165e938"
}
`
	if out != want {
		t.Errorf("stdout must be the pinned golden document:\ngot:\n%s\nwant:\n%s", out, want)
	}
}

// TestGetIdentityQualifiedLineFormLowestInstance: the qualified line
// form resolves to the lowest instance-version of the line. The
// fixture copy gains a second instance of the plan:roadmap-2026 line
// before the sync; the line form must return instance 1.
func TestGetIdentityQualifiedLineFormLowestInstance(t *testing.T) {
	seedGetRepo(t, func(repo string) {
		v2 := `---
namespace: eka-view-fixture
type: plan
id: roadmap-2026
instance-version: 2
revision: 2
content-state: draft
planning-state: draft
existence-state: active
phase: release
dimension: planning
author: Engineering Architecture
created: 2026-08-06
updated: 2026-08-06
supersedes: []
derives-from: []
depends-on: []
change-log:
  - date: 2026-08-06
    domain: existence-state
    from: "-"
    to: active
    by: Engineering Architecture
  - date: 2026-08-06
    domain: content-state
    from: "-"
    to: draft
    by: Engineering Architecture
  - date: 2026-08-06
    domain: planning-state
    from: "-"
    to: draft
    by: Engineering Architecture
  - date: 2026-08-06
    domain: phase
    from: "-"
    to: release
    by: Engineering Architecture
---
# Plan — Roadmap 2026 (v2)

## Objective

Objective v2.

## Scope

Scope v2.

## Out of Scope

Out of scope v2.
`
		if err := os.WriteFile(filepath.Join(repo, "docs", "planning", "plan-roadmap-2026-v2.md"), []byte(v2), 0o644); err != nil {
			t.Fatal(err)
		}
	})
	doc := getDocument(t, "get", "eka-view-fixture/plan:roadmap-2026")
	if got := doc["canonical_form"]; got != "eka-view-fixture/plan:roadmap-2026:1" {
		t.Errorf("canonical_form = %v, want the lowest instance (v1)", got)
	}
	if got := doc["revision"]; got != float64(1) {
		t.Errorf("revision = %v, want 1", got)
	}
}

// TestGetUnregisteredRepoExitsOne: the repository-state gate runs
// first — a repository not registered in the EKA workspace is refused
// with the deterministic message and the sync hint, exit 1, no JSON.
// The workspace must exist (get never creates one): seed it with
// runtime.Ensure, then run from an unregistered directory.
func TestGetUnregisteredRepoExitsOne(t *testing.T) {
	t.Setenv("EKA_HOME", t.TempDir())
	r, err := runtime.Ensure()
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	chdirInto(t, t.TempDir())
	code, out, errText := runIn([]string{"get", "execution"})
	if code != 1 {
		t.Fatalf("exit = %d, want 1\nstdout: %s\nstderr: %s", code, out, errText)
	}
	if out != "" {
		t.Errorf("stdout must be empty (no JSON), got %q", out)
	}
	if !strings.Contains(errText, "eka: get refused: repository") {
		t.Errorf("stderr must carry the refusal, got %q", errText)
	}
	if !strings.Contains(errText, "not registered in the EKA workspace") {
		t.Errorf("stderr must explain the refusal, got %q", errText)
	}
	if !strings.Contains(errText, "eka sync") || !strings.Contains(errText, "eka project register") {
		t.Errorf("stderr must hint at sync/register, got %q", errText)
	}
}

// TestGetNoWorkspaceExitsOne: `eka get` never creates the workspace —
// a missing workspace.json is a deterministic refusal, exit 1.
func TestGetNoWorkspaceExitsOne(t *testing.T) {
	home := t.TempDir()
	t.Setenv("EKA_HOME", home)
	chdirInto(t, t.TempDir())
	code, out, errText := runIn([]string{"get", "execution"})
	if code != 1 {
		t.Fatalf("exit = %d, want 1\nstdout: %s\nstderr: %s", code, out, errText)
	}
	if out != "" {
		t.Errorf("stdout must be empty (no JSON), got %q", out)
	}
	if !strings.Contains(errText, "eka: get refused: no EKA workspace at") {
		t.Errorf("stderr must name the missing workspace, got %q", errText)
	}
	if !strings.Contains(errText, "run 'eka sync' first") {
		t.Errorf("stderr must hint at 'eka sync', got %q", errText)
	}
	// The refusal must not have created the workspace.
	if _, err := os.Stat(filepath.Join(home, "workspace.json")); !os.IsNotExist(err) {
		t.Error("get must never initialize the workspace")
	}
}

// TestGetInvalidDomainExitsTwo: a target without ":" that is not one
// of the five Engineering Domain tokens is a usage error listing the
// valid domains — exit 2, no workspace access.
func TestGetInvalidDomainExitsTwo(t *testing.T) {
	seedGetRepo(t, nil)
	code, out, errText := runIn([]string{"get", "bogus"})
	if code != 2 {
		t.Fatalf("exit = %d, want 2\nstdout: %s\nstderr: %s", code, out, errText)
	}
	if out != "" {
		t.Errorf("stdout must be empty (no JSON), got %q", out)
	}
	if !strings.Contains(errText, `unknown domain "bogus"`) {
		t.Errorf("stderr must name the domain, got %q", errText)
	}
	for _, token := range []string{"discovery", "architecture", "planning", "execution", "operations"} {
		if !strings.Contains(errText, token) {
			t.Errorf("stderr must list the %s domain token, got %q", token, errText)
		}
	}
}

// TestGetUnknownIdentityExitsTwo: an identity that parses but does not
// exist is a deterministic usage-class error, exit 2.
func TestGetUnknownIdentityExitsTwo(t *testing.T) {
	seedGetRepo(t, nil)
	code, out, errText := runIn([]string{"get", "eka-view-fixture/sto:nonexistent"})
	if code != 2 {
		t.Fatalf("exit = %d, want 2\nstdout: %s\nstderr: %s", code, out, errText)
	}
	if out != "" {
		t.Errorf("stdout must be empty (no JSON), got %q", out)
	}
	if !strings.Contains(errText, `no knowledge object matches "eka-view-fixture/sto:nonexistent"`) {
		t.Errorf("stderr must name the missing identity, got %q", errText)
	}
}

// TestGetUnqualifiedIdentityExitsTwo: unqualified reference forms are
// refused with the expected forms listed (the Resolver contract) —
// exit 2. (A target without ":" is a domain query, not an identity:
// see TestGetInvalidDomainExitsTwo.)
func TestGetUnqualifiedIdentityExitsTwo(t *testing.T) {
	seedGetRepo(t, nil)
	for _, target := range []string{"sto:alpha", ":alpha"} {
		code, out, errText := runIn([]string{"get", target})
		if code != 2 {
			t.Errorf("target %q: exit = %d, want 2\nstdout: %s", target, code, out)
		}
		if out != "" {
			t.Errorf("target %q: stdout must be empty, got %q", target, out)
		}
		if !strings.Contains(errText, "<ns>/") {
			t.Errorf("target %q: stderr must list the expected qualified forms, got %q", target, errText)
		}
	}
}

// TestGetNoArgExitsTwo: `eka get` without a target is a usage error
// with the query-model summary on stderr (machine commands never print
// banners) — exit 2.
func TestGetNoArgExitsTwo(t *testing.T) {
	code, out, errText := runIn([]string{"get"})
	if code != 2 {
		t.Fatalf("exit = %d, want 2", code)
	}
	if out != "" {
		t.Errorf("stdout must be empty (no banner), got %q", out)
	}
	for _, want := range []string{
		"eka get <target>",
		"<ns>/<type>:<id>",
		"discovery | architecture | planning | execution | operations",
	} {
		if !strings.Contains(errText, want) {
			t.Errorf("stderr must carry the usage summary with %q, got %q", want, errText)
		}
	}
	// The usage error must not depend on a workspace.
	if strings.Contains(errText, "workspace") {
		t.Errorf("stderr must not touch the workspace, got %q", errText)
	}
}

// TestGetHelpExitsZero covers the help entry points: the long help
// documents the machine interface purpose, the query model, the stable
// schema, the stdout contract, the exit codes and the examples.
func TestGetHelpExitsZero(t *testing.T) {
	for _, args := range [][]string{{"get", "-h"}, {"get", "--help"}} {
		code, text, _ := runIn(args)
		if code != 0 {
			t.Errorf("args %v: exit = %d, want 0", args, code)
		}
		for _, want := range []string{
			"eka get",
			"machine-readable",
			"eka-cko-v1",
			"canonical form",
			"qualified line form",
			"discovery | architecture | planning | execution |",
			"operations",
			"stdout carries ONLY the JSON document",
			"Exit codes:",
			"eka view",
			"eka get feather/adr:001-identity-serialization:1",
			"eka get feather/sto:publish-post",
			"eka get architecture",
			"eka get execution",
		} {
			if !strings.Contains(text, want) {
				t.Errorf("args %v: help missing %q:\n%s", args, want, text)
			}
		}
	}
}

// TestGetTooManyArgsExitsTwo: at most one target.
func TestGetTooManyArgsExitsTwo(t *testing.T) {
	code, _, _ := runIn([]string{"get", "execution", "architecture"})
	if code != 2 {
		t.Errorf("exit = %d, want 2", code)
	}
}

// TestGetStdoutPurity: the stdout of every success path is exactly one
// parseable JSON document — no extra lines, no banners.
func TestGetStdoutPurity(t *testing.T) {
	seedGetRepo(t, nil)
	for _, args := range [][]string{
		{"get", "eka-view-fixture/adr:001-login-serialization:1"},
		{"get", "eka-view-fixture/sto:alpha"},
		{"get", "architecture"},
		{"get", "execution"},
		{"get", "operations"},
	} {
		code, out, errText := runIn(args)
		if code != 0 {
			t.Fatalf("%v: exit = %d, want 0\nstderr: %s", args, code, errText)
		}
		// stdout carries ONLY the JSON document followed by a single
		// trailing newline: exactly one trailing newline, no banners,
		// no informational lines.
		if !strings.HasSuffix(out, "\n") || strings.HasSuffix(out, "\n\n") {
			t.Errorf("%v: stdout must carry exactly one trailing newline, got %q", args, out)
		}
		if !json.Valid([]byte(out)) {
			t.Errorf("%v: stdout must be valid JSON, got %q", args, out)
		}
		var doc map[string]any
		if err := json.Unmarshal([]byte(out), &doc); err != nil {
			t.Errorf("%v: stdout must parse as one JSON document: %v", args, err)
		}
		if doc["schema"] != "eka-cko-v1" {
			t.Errorf("%v: document must carry the stable schema, got %v", args, doc["schema"])
		}
	}
}

// TestGetDomainQueryCollection: a domain query produces the "domain"
// collection with the canonical domain name, the unit count and the
// units sorted by canonical form.
func TestGetDomainQueryCollection(t *testing.T) {
	seedGetRepo(t, nil)
	code, out, errText := runIn([]string{"get", "execution"})
	if code != 0 {
		t.Fatalf("exit = %d, want 0\nstdout: %s\nstderr: %s", code, out, errText)
	}
	var col struct {
		Schema     string           `json:"schema"`
		Collection string           `json:"collection"`
		Domain     string           `json:"domain"`
		Count      int              `json:"count"`
		Units      []map[string]any `json:"units"`
	}
	if err := json.Unmarshal([]byte(out), &col); err != nil {
		t.Fatal(err)
	}
	if col.Schema != "eka-cko-v1" || col.Collection != "domain" || col.Domain != "Execution" {
		t.Errorf("collection header = %+v, want schema eka-cko-v1 / collection domain / domain Execution", col)
	}
	// The fixture carries 17 Execution units: 2 containers, 9 ticket
	// projections, 6 work items.
	if col.Count != 17 || len(col.Units) != 17 {
		t.Errorf("count = %d (units %d), want 17", col.Count, len(col.Units))
	}
	// Sorted by canonical form.
	for i := 1; i < len(col.Units); i++ {
		prev := col.Units[i-1]["canonical_form"].(string)
		cur := col.Units[i]["canonical_form"].(string)
		if prev >= cur {
			t.Errorf("units not sorted at %d: %q >= %q", i, prev, cur)
		}
	}
	// Every unit carries the derived engineering domain and stratum.
	for _, u := range col.Units {
		if u["engineering_domain"] != "Execution" || u["stratum"] != float64(4) {
			t.Errorf("unit %v: engineering_domain/stratum wrong", u["canonical_form"])
		}
	}
}

// TestGetDeterministicCLI: two runs of each query produce
// byte-identical stdout.
func TestGetDeterministicCLI(t *testing.T) {
	seedGetRepo(t, nil)
	runOnce := func(args ...string) string {
		_, out, _ := runIn(args)
		return out
	}
	for _, args := range [][]string{
		{"get", "eka-view-fixture/adr:001-login-serialization:1"},
		{"get", "eka-view-fixture/sto:alpha"},
		{"get", "discovery"},
		{"get", "architecture"},
		{"get", "planning"},
		{"get", "execution"},
		{"get", "operations"},
	} {
		if a, b := runOnce(args...), runOnce(args...); a != b {
			t.Errorf("output differs between runs for %v", args)
		}
	}
}
