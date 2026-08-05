package conformance

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// This file implements the validation engine entry point:
//
//	report, err := conformance.Validate(root)
//
// The engine walks the root recursively, classifies every .md file, and runs
// Rules R1-R9 over the collected artifacts. It never imports anything from
// cmd/ and is fully reusable by future tooling.
//
// Scan policy (documented interpretation):
//   - Only files ending in .md are examined; every other file is ignored.
//   - Convention documents (no frontmatter, or frontmatter without both
//     `type` and `id`) are counted but skipped.
//   - Directories named "testdata" and directories starting with "." (e.g.
//     .git) are not descended into: they hold test fixtures and VCS
//     metadata, not repository knowledge content. Without this, Go test
//     fixtures under conformance/testdata would be validated as if they were
//     part of the knowledge base.
//   - Symlinks are not followed (filepath.WalkDir behavior).
//   - An unreadable .md file aborts the run with an error (a scan that
//     cannot see every file cannot certify compliance).
func Validate(root string) (*Report, error) {
	info, err := os.Stat(root)
	if err != nil {
		return nil, fmt.Errorf("cannot access scan root: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("scan root is not a directory: %s", root)
	}

	report := &Report{Root: filepath.Clean(root)}
	e := &engine{report: report}

	var paths []string
	err = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if path != root {
				name := d.Name()
				if strings.HasPrefix(name, ".") || name == "testdata" {
					return filepath.SkipDir
				}
			}
			return nil
		}
		if strings.HasSuffix(d.Name(), ".md") {
			paths = append(paths, path)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("scan failed: %w", err)
	}
	report.FilesScanned = len(paths)

	// Parse phase: classify every file. Parse failures become R0 results.
	for _, path := range paths {
		rel, relErr := filepath.Rel(root, path)
		if relErr != nil {
			rel = path
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("cannot read %s: %w", rel, err)
		}
		artifact, results := analyzeFile(rel, path, data)
		report.Results = append(report.Results, results...)
		if artifact != nil {
			e.artifacts = append(e.artifacts, artifact)
		}
	}
	report.Artifacts = len(e.artifacts)

	// Identity index for reference resolution and the R9 supersession check.
	e.buildIndex()

	// Rules. R0 results were collected during parsing; R1 runs over the
	// whole set; R2-R9 run per artifact.
	e.rule1()
	for _, a := range e.artifacts {
		e.rule2(a)
		e.rule3(a)
		e.rule4(a)
		e.rule5(a)
		e.rule6(a)
		e.rule7(a)
		e.rule8(a)
		e.rule9(a)
	}

	return report, nil
}

// engine carries the state shared by all rule checks.
type engine struct {
	report    *Report
	artifacts []*Artifact
	// byLine indexes artifacts by identity line (namespace, type, id);
	// each bucket holds all instances sorted by instance-version.
	byLine map[string][]*Artifact
}

// identityKey builds the identity line key for the index.
func identityKey(ns, typeToken, id string) string {
	return ns + "\x00" + typeToken + "\x00" + id
}

func (e *engine) buildIndex() {
	e.byLine = make(map[string][]*Artifact)
	for _, a := range e.artifacts {
		key := identityKey(a.Namespace, a.Type, a.ID)
		e.byLine[key] = append(e.byLine[key], a)
	}
	for _, bucket := range e.byLine {
		sort.Slice(bucket, func(i, j int) bool {
			return bucket[i].InstanceVersion < bucket[j].InstanceVersion
		})
	}
}

// add appends a result for an artifact.
func (e *engine) add(a *Artifact, rule string, sev Severity, format string, args ...any) {
	e.report.Results = append(e.report.Results, Result{
		File:     a.RelPath,
		Rule:     rule,
		Severity: sev,
		Message:  fmt.Sprintf(format, args...),
	})
}

// addFile appends a result not tied to a parsed artifact.
func (e *engine) addFile(rel, rule string, sev Severity, format string, args ...any) {
	e.report.Results = append(e.report.Results, Result{
		File:     rel,
		Rule:     rule,
		Severity: sev,
		Message:  fmt.Sprintf(format, args...),
	})
}
