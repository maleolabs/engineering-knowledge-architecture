package bootstrap

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
)

// This file implements Stage 2 of the bootstrap model: Bootstrap Planning.
// BuildPlan derives a deterministic action list from discovery results and
// wizard answers. The plan is the single source of truth for both the
// dry-run preview and the generator; nothing writes outside the plan.

// ActionKind is the kind of one planned action.
type ActionKind string

// Action kinds. Each kind has one deterministic rendering (see String).
const (
	// ActionCreateDir creates a directory (mode 0755).
	ActionCreateDir ActionKind = "create-dir"
	// ActionCreateFile writes one embedded skeleton file verbatim.
	ActionCreateFile ActionKind = "create-file"
	// ActionGenerateReadme writes README.md with the heading substituted.
	ActionGenerateReadme ActionKind = "generate-readme"
	// ActionReuse leaves an existing file untouched (identical content).
	ActionReuse ActionKind = "reuse"
	// ActionOverwriteConfirm marks an existing file whose content differs
	// from the generated one: it needs explicit confirmation (interactive)
	// or is skipped (non-interactive).
	ActionOverwriteConfirm ActionKind = "overwrite-confirm"
	// ActionSkip records a deliberately skipped step (e.g. git init not
	// wanted, README not requested).
	ActionSkip ActionKind = "skip"
	// ActionGitInit runs `git init` in the target.
	ActionGitInit ActionKind = "git-init"
	// ActionGitSkip records that git init will not run and why.
	ActionGitSkip ActionKind = "git-skip"
	// ActionValidate runs conformance validation over the target.
	ActionValidate ActionKind = "validate"
)

// Action is one deterministic plan step.
type Action struct {
	// Kind selects the step.
	Kind ActionKind
	// Path is the affected path: the target itself for git/validate, a
	// forward-slash relative path (e.g. "docs/README.md") for files and
	// skeleton directories. The target-dir creation action uses the
	// sentinel "." with the target name in Detail.
	Path string
	// Detail is optional context rendered in parentheses.
	Detail string
	// Content is the exact bytes the action will write, resolved at plan
	// time so the plan is self-contained (dry-run preview and generation
	// always agree). Only file-writing actions carry it.
	Content []byte
}

// String renders the action as a stable plan line.
func (a Action) String() string {
	switch a.Kind {
	case ActionCreateDir:
		// The target-dir action uses the sentinel path "." with the target
		// name in Detail (a real target name could collide with a skeleton
		// directory name such as "docs").
		p := a.Path
		if p == "." && a.Detail != "" {
			p = a.Detail
		}
		return "create dir: " + p
	case ActionCreateFile:
		return "create file: " + a.Path + " (from skeleton)"
	case ActionGenerateReadme:
		return "generate file: " + a.Path + " (from skeleton)"
	case ActionReuse:
		if a.Detail != "" {
			return "reuse: " + a.Path + " (" + a.Detail + ")"
		}
		return "reuse: " + a.Path
	case ActionOverwriteConfirm:
		return "overwrite confirm: " + a.Path
	case ActionSkip:
		return "skip: " + a.Path + " (" + a.Detail + ")"
	case ActionGitInit:
		return "git init: " + a.Path
	case ActionGitSkip:
		return "git init: " + a.Detail
	case ActionValidate:
		if a.Detail != "" {
			return "validate: " + a.Path + " " + a.Detail
		}
		return "validate: " + a.Path
	default:
		return string(a.Kind) + ": " + a.Path
	}
}

// BuildPlan derives the deterministic plan for target from discovery d and
// answers a. Ordering is stable: directories, files, git, validation; each
// group is sorted by path. A target that is already an EKA repository
// yields a reuse + validate plan only — nothing is ever planned to be
// overwritten silently.
func BuildPlan(target string, d *Discovery, a Answers) []Action {
	if d.Exists && d.IsEkaRepo {
		return []Action{
			{Kind: ActionReuse, Path: target, Detail: "existing EKA repository (already initialized)"},
			{Kind: ActionValidate, Path: target},
		}
	}

	plan := []Action{}

	// Directories.
	if !d.Exists {
		// The target dir uses the sentinel path "." (see Action.String);
		// joining it onto the target always yields the target itself.
		plan = append(plan, Action{Kind: ActionCreateDir, Path: ".", Detail: target})
	}
	skelDirs, skelFiles, err := skeletonTree()
	if err != nil {
		// The skeleton is embedded at compile time; failure is impossible
		// in practice. Keep the plan build honest anyway.
		plan = append(plan, Action{Kind: ActionSkip, Path: "skeleton", Detail: "unreadable: " + err.Error()})
		return plan
	}
	for _, rel := range skelDirs {
		if !pathExists(filepath.Join(d.AbsTarget, rel)) {
			plan = append(plan, Action{Kind: ActionCreateDir, Path: rel})
		}
	}

	// Files: README first, then skeleton docs files, all sorted by path.
	if a.GenerateReadme {
		existing := readmePath(d)
		switch {
		case existing == "":
			plan = append(plan, Action{Kind: ActionGenerateReadme, Path: "README.md", Content: generatedReadme(a.ProjectName)})
		case fileMatches(filepath.Join(d.AbsTarget, existing), generatedReadme(a.ProjectName)):
			plan = append(plan, Action{Kind: ActionReuse, Path: existing})
		default:
			plan = append(plan, Action{Kind: ActionOverwriteConfirm, Path: existing, Content: generatedReadme(a.ProjectName)})
		}
	} else if readmePath(d) == "" {
		plan = append(plan, Action{Kind: ActionSkip, Path: "README.md", Detail: "not requested"})
	}

	for _, rel := range skelFiles {
		if rel == "README.md" {
			continue // handled above; the skeleton copy is never written raw
		}
		content, err := readSkeletonFile(rel)
		if err != nil {
			plan = append(plan, Action{Kind: ActionSkip, Path: rel, Detail: "unreadable: " + err.Error()})
			continue
		}
		targetPath := filepath.Join(d.AbsTarget, rel)
		switch {
		case !pathExists(targetPath):
			plan = append(plan, Action{Kind: ActionCreateFile, Path: rel, Content: content})
		case fileMatches(targetPath, content):
			plan = append(plan, Action{Kind: ActionReuse, Path: rel})
		default:
			plan = append(plan, Action{Kind: ActionOverwriteConfirm, Path: rel, Content: content})
		}
	}

	// Git.
	switch {
	case d.IsGitRepo:
		plan = append(plan, Action{Kind: ActionGitSkip, Detail: "skipped (already a git repository)"})
	case !d.GitAvailable:
		plan = append(plan, Action{Kind: ActionGitSkip, Detail: "skipped (git not available)"})
	case !a.InitGit:
		if a.Interactive {
			plan = append(plan, Action{Kind: ActionGitSkip, Detail: "skipped (declined)"})
		} else {
			plan = append(plan, Action{Kind: ActionGitSkip, Detail: "skipped (non-interactive mode)"})
		}
	default:
		plan = append(plan, Action{Kind: ActionGitInit, Path: target})
	}

	// Validation.
	plan = append(plan, Action{Kind: ActionValidate, Path: target, Detail: "after generation"})
	return plan
}

// readmePath returns the name of an existing README ("README.md" or
// "README"), or "" when none exists.
func readmePath(d *Discovery) string {
	return d.ReadmePath
}

// pathExists reports whether path exists (any type).
func pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// fileMatches reports whether the file at path exists and contains exactly
// the given bytes.
func fileMatches(path string, want []byte) bool {
	got, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	return bytes.Equal(got, want)
}

// generatedReadme returns the README content for the given project name:
// the embedded template with the exact heading line "# <Nama Produk>"
// replaced by "# <name>". Every other line (including the template date)
// stays verbatim.
func generatedReadme(name string) []byte {
	data, err := readSkeletonFile("README.md")
	if err != nil {
		return nil
	}
	lines := strings.Split(string(data), "\n")
	if len(lines) > 0 && lines[0] == "# <Nama Produk>" {
		lines[0] = "# " + name
	}
	return []byte(strings.Join(lines, "\n"))
}
