package bootstrap

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// This file implements Stage 3 of the bootstrap model: the Interactive
// Wizard. The wizard is adaptive: it only asks questions whose answers are
// not already known from discovery. When stdin is not a terminal, the
// wizard is skipped entirely and DefaultAnswers provides deterministic
// discovery-derived values.
//
// EKA v1 has no methodology taxonomy, so no methodology question exists:
// there is no canonical set of methodologies to choose from. The question
// is intentionally absent and this fact is documented in the spec.

// QuestionKind identifies one wizard question.
type QuestionKind string

// The five wizard questions, asked in this fixed order.
const (
	// QProjectName asks for the display name used in the README heading.
	// Asked only when the target directory name is empty or unusable
	// (filesystem root, whitespace-only).
	QProjectName QuestionKind = "project-name"
	// QNamespace asks for the frontmatter namespace. Asked only when the
	// project name is not itself a valid namespace.
	QNamespace QuestionKind = "namespace"
	// QDescription asks for an optional free-text description. Collected
	// for the record; the EKA v1 skeleton has no template slot for it, so
	// it is not written anywhere in v1.
	QDescription QuestionKind = "description"
	// QReadme asks whether to generate README.md from the skeleton
	// template. Asked only when no README exists yet.
	QReadme QuestionKind = "readme"
	// QGit asks whether to run `git init`. Asked only when the target is
	// not already a git repository and a git executable is available.
	QGit QuestionKind = "git"
)

// Question is one wizard prompt: what to ask, how to word it, and the
// default answer offered.
type Question struct {
	Kind    QuestionKind
	Prompt  string
	Default string
}

// Answers are the wizard outcomes consumed by the planner.
type Answers struct {
	// ProjectName is the display name for the README heading.
	ProjectName string
	// Namespace is the frontmatter namespace (validated).
	Namespace string
	// Description is optional free text; unused by the v1 skeleton.
	Description string
	// GenerateReadme requests README.md generation from the template.
	GenerateReadme bool
	// InitGit requests `git init` in the target.
	InitGit bool
	// Interactive reports whether the answers came from an interactive
	// session (affects plan wording for skipped git init).
	Interactive bool
}

// fallbackName is used when no usable name can be derived from the target
// directory (e.g. bootstrapping the filesystem root).
const fallbackName = "eka-project"

// validProjectName reports whether a directory base name is usable as a
// project display name: non-empty, not the filesystem root, not
// whitespace-only.
func validProjectName(name string) bool {
	return name != "" && name != "/" && strings.TrimSpace(name) != ""
}

// isValidNamespace reports whether s is a valid EKA namespace: non-empty
// and containing only lowercase letters, digits and hyphens (no '/', ':',
// whitespace or any other character).
func isValidNamespace(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9', r == '-':
		default:
			return false
		}
	}
	return true
}

// sanitizeNamespace derives a valid namespace from arbitrary text:
// lowercase, every run of invalid characters becomes a single hyphen,
// leading/trailing hyphens are trimmed. Returns "" when nothing usable
// remains.
func sanitizeNamespace(s string) string {
	var b strings.Builder
	dashPending := false
	for _, r := range strings.ToLower(s) {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			if dashPending && b.Len() > 0 {
				b.WriteByte('-')
			}
			dashPending = false
			b.WriteRune(r)
		case r == '-':
			dashPending = true
		default:
			// Whitespace, '/', ':' and everything else act as separators.
			dashPending = true
		}
	}
	return strings.Trim(b.String(), "-")
}

// defaultNamespace derives the namespace default from a project name.
func defaultNamespace(name string) string {
	if isValidNamespace(name) {
		return name
	}
	if ns := sanitizeNamespace(name); ns != "" {
		return ns
	}
	return fallbackName
}

// DefaultAnswers returns the deterministic non-interactive answers derived
// from discovery: project name = target base name (sanitized/fallback when
// unusable), namespace = sanitized project name, no description, README
// generated when absent, git never initialized. Non-interactive runs never
// run `git init`: it is a side effect beyond file writes and would break
// determinism guarantees.
func DefaultAnswers(d *Discovery) Answers {
	project := d.BaseName
	if !validProjectName(project) {
		project = defaultNamespace(d.BaseName)
	}
	return Answers{
		ProjectName:    project,
		Namespace:      defaultNamespace(project),
		Description:    "",
		GenerateReadme: !d.HasReadme,
		InitGit:        false,
		Interactive:    false,
	}
}

// NeededQuestions returns the questions whose answers discovery cannot
// provide, in fixed order. It is pure (no I/O) so adaptivity is unit
// testable.
func NeededQuestions(d *Discovery) []Question {
	var qs []Question
	if !validProjectName(d.BaseName) {
		qs = append(qs, Question{
			Kind:    QProjectName,
			Prompt:  "Project name",
			Default: defaultNamespace(d.BaseName),
		})
	}
	project := d.BaseName
	if !validProjectName(project) {
		project = defaultNamespace(d.BaseName)
	}
	if !isValidNamespace(project) {
		qs = append(qs, Question{
			Kind:    QNamespace,
			Prompt:  "Namespace (lowercase letters, digits, hyphens)",
			Default: defaultNamespace(project),
		})
	}
	qs = append(qs, Question{Kind: QDescription, Prompt: "Project description (optional)"})
	if !d.HasReadme {
		qs = append(qs, Question{Kind: QReadme, Prompt: "Generate README.md from template?", Default: "y"})
	}
	if !d.IsGitRepo && d.GitAvailable {
		qs = append(qs, Question{Kind: QGit, Prompt: "Initialize git repository?", Default: "y"})
	}
	return qs
}

// Ask runs the interactive wizard: prints each needed question to w and
// reads answers from r. Answers are validated as they are entered; invalid
// namespaces are re-prompted. The wizard is sequentially adaptive: when the
// project name answer itself is not a valid namespace, the namespace
// question follows immediately. If the input stream ends early (EOF), the
// remaining answers fall back to their defaults so a closed pipe can never
// hang the run.
func Ask(d *Discovery, r io.Reader, w io.Writer) (Answers, error) {
	a := DefaultAnswers(d)
	a.Interactive = true
	sc := bufio.NewScanner(r)
	for _, q := range NeededQuestions(d) {
		switch q.Kind {
		case QProjectName:
			a.ProjectName = askLine(sc, w, q, a.ProjectName)
			// Sequential adaptivity: the answered name decides whether
			// the namespace question is still needed.
			if !isValidNamespace(a.ProjectName) {
				a.Namespace = askNamespace(sc, w, a.Namespace)
			}
		case QNamespace:
			a.Namespace = askNamespace(sc, w, a.Namespace)
		case QDescription:
			a.Description = askLine(sc, w, q, a.Description)
		case QReadme:
			a.GenerateReadme = askYesNo(sc, w, q, true)
		case QGit:
			a.InitGit = askYesNo(sc, w, q, true)
		}
	}
	return a, nil
}

// askNamespace prompts until a valid namespace is entered.
func askNamespace(sc *bufio.Scanner, w io.Writer, def string) string {
	q := Question{
		Kind:    QNamespace,
		Prompt:  "Namespace (lowercase letters, digits, hyphens)",
		Default: def,
	}
	for {
		answer := askLine(sc, w, q, def)
		if isValidNamespace(answer) {
			return answer
		}
		fmt.Fprintf(w, "invalid namespace %q — use lowercase letters, digits and hyphens only\n", answer)
	}
}

// askLine prints the prompt and reads one line. Empty input or a closed
// stream yields the default.
func askLine(sc *bufio.Scanner, w io.Writer, q Question, def string) string {
	prompt := q.Prompt
	if def != "" {
		prompt += " [" + def + "]"
	}
	fmt.Fprintf(w, "%s: ", prompt)
	if !sc.Scan() {
		fmt.Fprintln(w)
		return def
	}
	fmt.Fprintln(w)
	if answer := strings.TrimSpace(sc.Text()); answer != "" {
		return answer
	}
	return def
}

// askYesNo prints a y/n prompt and reads the answer. Empty input or a
// closed stream yields the default.
func askYesNo(sc *bufio.Scanner, w io.Writer, q Question, def bool) bool {
	suffix := " [y/N]"
	if def {
		suffix = " [Y/n]"
	}
	fmt.Fprintf(w, "%s%s: ", q.Prompt, suffix)
	if !sc.Scan() {
		fmt.Fprintln(w)
		return def
	}
	fmt.Fprintln(w)
	switch strings.ToLower(strings.TrimSpace(sc.Text())) {
	case "y", "yes":
		return true
	case "n", "no":
		return false
	default:
		return def
	}
}
