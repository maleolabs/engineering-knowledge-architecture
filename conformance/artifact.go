package conformance

import (
	"fmt"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// This file implements the artifact-vs-convention-document classification
// and the frontmatter parsing that feeds every rule:
//
//	A .md file is an Artifact iff its YAML frontmatter contains BOTH
//	`type` and `id`. (reference-architecture.md §3, docs/README.md)
//
// Files without both are Convention Documents (READMEs, protocol.md,
// validation.md, transfer.md, the canonical spec) and are skipped entirely.
// Files with exactly one of the two violate the artifact rule and are
// reported as malformed (R0 bucket, see report.go).
//
// Structural failures of the frontmatter block itself (unparseable YAML,
// unterminated block, non-mapping root, missing/invalid identity fields,
// invalid dates) are also reported in the R0 bucket because they are
// pre-conditions for the numbered rules rather than rule violations.

// Artifact is one parsed EKA artifact found under the scanned root.
type Artifact struct {
	// RelPath is the path relative to the scan root, used in results so
	// output is stable regardless of the machine.
	RelPath string
	// AbsPath is the absolute file path.
	AbsPath string

	// Identity (frontmatter is the source of truth).
	Namespace       string
	Type            string
	ID              string
	InstanceVersion int
	Revision        int
	Created         string
	Updated         string

	// States maps present state domain fields to their values. A domain
	// absent from this map is not present on the artifact.
	States map[string]string
	// Phase is the phase context attribute (scp-/plan- only).
	Phase       string
	HasPhase    bool
	HasPhaseKey bool

	// Classification (Rule 6).
	Dimension           string
	HasDimension        bool
	DimensionsSecondary []string

	// Relations maps relationship field name -> raw reference strings.
	Relations map[string][]string

	// ChangeLog holds the parsed change-log entries in file order.
	ChangeLog []ChangeLogEntry

	// BodyLines are the content lines after the frontmatter block.
	BodyLines []string
}

// ChangeLogEntry is one parsed {date, domain, from, to, by} entry.
type ChangeLogEntry struct {
	Date   string
	Domain string
	From   string
	To     string
	By     string
}

// projectHeader is the exact projection header line required by Rule 8
// (validation.md Rule 8, ADR-003 §4). The em dash is U+2014.
const projectHeader = "> Generated \u2014 State Projection. Do NOT edit state here; refresh on read."

// analyzeFile classifies and parses one .md file.
//
// It returns a nil artifact (and no results) for convention documents and
// for files without a frontmatter block. Non-nil results are R0 structural
// findings (malformed frontmatter, artifact-rule violations, missing or
// invalid identity fields).
func analyzeFile(relPath, absPath string, data []byte) (*Artifact, []Result) {
	lines := strings.Split(string(data), "\n")
	if strings.TrimRight(lines[0], "\r") != "---" {
		// No frontmatter block: convention document, skipped entirely.
		return nil, nil
	}

	closeIdx := -1
	for i := 1; i < len(lines); i++ {
		if strings.TrimRight(lines[i], "\r") == "---" {
			closeIdx = i
			break
		}
	}
	if closeIdx < 0 {
		return nil, []Result{{
			File:     relPath,
			Rule:     RuleStructural,
			Severity: SeverityError,
			Message:  "frontmatter block starts with --- but never closes",
		}}
	}

	var fm map[string]any
	if err := yaml.Unmarshal([]byte(strings.Join(lines[1:closeIdx], "\n")), &fm); err != nil {
		return nil, []Result{{
			File:     relPath,
			Rule:     RuleStructural,
			Severity: SeverityError,
			Message:  fmt.Sprintf("frontmatter is not valid YAML: %v", err),
		}}
	}
	if fm == nil {
		// An empty or comment-only block decodes to nil: treat it as a
		// document without type/id.
		fm = map[string]any{}
	}

	_, hasType := fm["type"]
	_, hasID := fm["id"]
	if hasType != hasID {
		return nil, []Result{{
			File:     relPath,
			Rule:     RuleStructural,
			Severity: SeverityError,
			Message:  "frontmatter contains exactly one of `type` and `id`; an artifact requires both (type XOR id is malformed)",
		}}
	}
	if !hasType {
		// Convention document: frontmatter without type AND without id.
		return nil, nil
	}

	a := &Artifact{
		RelPath:   relPath,
		AbsPath:   absPath,
		States:    map[string]string{},
		Relations: map[string][]string{},
	}
	var results []Result

	// --- Identity fields. ---
	a.Type, _ = asString(fm["type"])
	if _, ok := typeTokens[a.Type]; !ok {
		results = append(results, Result{
			File: relPath, Rule: RuleStructural, Severity: SeverityError,
			Message: fmt.Sprintf("unknown artifact type %q; expected one of the 26 EKA type tokens", a.Type),
		})
	}
	a.ID, _ = asString(fm["id"])
	if a.ID == "" {
		results = append(results, Result{
			File: relPath, Rule: RuleStructural, Severity: SeverityError,
			Message: "artifact id must be a non-empty string",
		})
	}

	a.Namespace, _ = asString(fm["namespace"])
	if a.Namespace == "" {
		results = append(results, Result{
			File: relPath, Rule: RuleStructural, Severity: SeverityError,
			Message: "missing required identity field `namespace`",
		})
	}

	if v, ok := fm["instance-version"]; ok {
		if n, valid := asInt(v); valid {
			a.InstanceVersion = n
		} else {
			results = append(results, Result{
				File: relPath, Rule: RuleStructural, Severity: SeverityError,
				Message: "`instance-version` must be an integer",
			})
		}
	} else {
		results = append(results, Result{
			File: relPath, Rule: RuleStructural, Severity: SeverityError,
			Message: "missing required identity field `instance-version`",
		})
	}
	if a.InstanceVersion < 1 {
		results = append(results, Result{
			File: relPath, Rule: RuleStructural, Severity: SeverityError,
			Message: "`instance-version` must be >= 1",
		})
	}

	if v, ok := fm["revision"]; ok {
		if n, valid := asInt(v); valid {
			a.Revision = n
		} else {
			results = append(results, Result{
				File: relPath, Rule: RuleStructural, Severity: SeverityError,
				Message: "`revision` must be an integer",
			})
		}
	} else {
		results = append(results, Result{
			File: relPath, Rule: RuleStructural, Severity: SeverityError,
			Message: "missing required identity field `revision`",
		})
	}
	if a.Revision < 1 {
		results = append(results, Result{
			File: relPath, Rule: RuleStructural, Severity: SeverityError,
			Message: "`revision` must be >= 1",
		})
	}

	for _, f := range []struct{ key, label string }{
		{"created", "`created`"},
		{"updated", "`updated`"},
	} {
		if v, ok := fm[f.key]; ok {
			s, ok := asDateString(v)
			if !ok || !validDate(s) {
				results = append(results, Result{
					File: relPath, Rule: RuleStructural, Severity: SeverityError,
					Message: fmt.Sprintf("%s must be a date in YYYY-MM-DD format", f.label),
				})
			} else if f.key == "created" {
				a.Created = s
			} else {
				a.Updated = s
			}
		}
	}

	// --- State fields (Rules 3 and 4 operate on these). ---
	for _, domain := range stateFields {
		if v, ok := fm[domain]; ok {
			if s, isStr := asString(v); isStr {
				a.States[domain] = s
			} else {
				results = append(results, Result{
					File: relPath, Rule: RuleStructural, Severity: SeverityError,
					Message: fmt.Sprintf("state field `%s` must be a string", domain),
				})
			}
		}
	}
	if v, ok := fm[DomainPhase]; ok {
		a.HasPhaseKey = true
		if s, isStr := asString(v); isStr {
			a.Phase = s
			a.HasPhase = true
		} else {
			results = append(results, Result{
				File: relPath, Rule: RuleStructural, Severity: SeverityError,
				Message: "`phase` must be a string",
			})
		}
	}

	// --- Classification (Rule 6). ---
	if v, ok := fm["dimension"]; ok {
		a.HasDimension = true
		if s, isStr := asString(v); isStr {
			a.Dimension = s
		} else {
			results = append(results, Result{
				File: relPath, Rule: RuleStructural, Severity: SeverityError,
				Message: "`dimension` must be a string",
			})
		}
	}
	if v, ok := fm["dimensions-secondary"]; ok {
		list, valid := asStringList(v)
		if !valid {
			results = append(results, Result{
				File: relPath, Rule: RuleStructural, Severity: SeverityError,
				Message: "`dimensions-secondary` must be a list of strings",
			})
		} else {
			a.DimensionsSecondary = list
		}
	}

	// --- Relationships (Rule 5). ---
	for _, field := range relationshipFields {
		if v, ok := fm[field]; ok {
			list, valid := asStringList(v)
			if !valid {
				results = append(results, Result{
					File: relPath, Rule: Rule5, Severity: SeverityError,
					Message: fmt.Sprintf("relationship field `%s` must be a list of references", field),
				})
			} else {
				a.Relations[field] = list
			}
		}
	}

	// --- Change-log (Rule 7). ---
	if v, ok := fm["change-log"]; ok {
		list, valid := asList(v)
		if !valid {
			results = append(results, Result{
				File: relPath, Rule: Rule7, Severity: SeverityError,
				Message: "`change-log` must be a list of {date, domain, from, to, by} entries",
			})
		} else {
			for i, item := range list {
				entry, err := parseChangeLogEntry(item)
				if err != nil {
					results = append(results, Result{
						File: relPath, Rule: Rule7, Severity: SeverityError,
						Message: fmt.Sprintf("change-log entry %d is malformed: %v", i+1, err),
					})
					continue
				}
				a.ChangeLog = append(a.ChangeLog, entry)
			}
		}
	}

	// --- Content. ---
	a.BodyLines = lines[closeIdx+1:]

	return a, results
}

// parseChangeLogEntry validates one change-log entry's shape. Domain
// ownership, value validity and transition legality are checked by Rule 7.
func parseChangeLogEntry(item any) (ChangeLogEntry, error) {
	m, ok := item.(map[string]any)
	if !ok {
		return ChangeLogEntry{}, fmt.Errorf("entry must be a mapping")
	}
	var e ChangeLogEntry
	required := []struct {
		key    string
		target *string
	}{
		{"date", &e.Date}, {"domain", &e.Domain}, {"from", &e.From}, {"to", &e.To}, {"by", &e.By},
	}
	for _, r := range required {
		v, ok := m[r.key]
		if !ok {
			return ChangeLogEntry{}, fmt.Errorf("missing required field %q", r.key)
		}
		if r.key == "date" {
			// The canonical ADRs write unquoted dates (e.g.
			// `date: 2026-08-05`), which yaml.v3 resolves as a
			// timestamp node; normalize it back to YYYY-MM-DD.
			s, isStr := asDateString(v)
			if !isStr {
				return ChangeLogEntry{}, fmt.Errorf("field %q must be a date", r.key)
			}
			e.Date = s
			continue
		}
		s, isStr := v.(string)
		if !isStr {
			return ChangeLogEntry{}, fmt.Errorf("field %q must be a string", r.key)
		}
		*r.target = s
	}
	if !validDate(e.Date) {
		return ChangeLogEntry{}, fmt.Errorf("`date` %q is not a valid YYYY-MM-DD date", e.Date)
	}
	return e, nil
}

// validDate reports whether s is a real calendar date in YYYY-MM-DD form.
func validDate(s string) bool {
	_, err := time.Parse("2006-01-02", s)
	return err == nil
}

// --- YAML value coercion helpers. ---

// asDateString coerces a frontmatter date value: plain strings pass through,
// and yaml.v3 timestamp nodes (unquoted YYYY-MM-DD values, used by the
// canonical ADRs) are normalized to YYYY-MM-DD.
func asDateString(v any) (string, bool) {
	if s, ok := v.(string); ok {
		return s, true
	}
	if t, ok := v.(time.Time); ok {
		return t.Format("2006-01-02"), true
	}
	return "", false
}

// asString coerces a scalar to string.
func asString(v any) (string, bool) {
	s, ok := v.(string)
	return s, ok
}

func asList(v any) ([]any, bool) {
	l, ok := v.([]any)
	return l, ok
}

func asStringList(v any) ([]string, bool) {
	l, ok := asList(v)
	if !ok {
		return nil, false
	}
	out := make([]string, 0, len(l))
	for _, item := range l {
		s, ok := item.(string)
		if !ok {
			return nil, false
		}
		out = append(out, s)
	}
	return out, true
}

// asInt coerces the integer forms produced by yaml.v3.
func asInt(v any) (int, bool) {
	switch n := v.(type) {
	case int:
		return n, true
	case int64:
		return int(n), true
	case uint64:
		return int(n), true
	case string:
		// Interpretation (documented): quoted integers such as
		// `instance-version: "1"` are rejected; the spec defines the
		// field as an integer and the canonical ADRs write it unquoted.
		return 0, false
	default:
		return 0, false
	}
}
