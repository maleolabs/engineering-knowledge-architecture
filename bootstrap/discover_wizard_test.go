package bootstrap

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// --- Workspace Discovery ------------------------------------------------

func TestDiscoverNonexistentTarget(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "nope")
	d, err := Discover(dir)
	if err != nil {
		t.Fatal(err)
	}
	if d.Exists {
		t.Error("target must not exist")
	}
	if d.BaseName != "nope" {
		t.Errorf("BaseName = %q, want %q", d.BaseName, "nope")
	}
}

func TestDiscoverEmptyDir(t *testing.T) {
	dir := t.TempDir()
	d, err := Discover(dir)
	if err != nil {
		t.Fatal(err)
	}
	if !d.Exists || !d.IsDir {
		t.Error("temp dir must exist as a directory")
	}
	if d.IsGitRepo || d.HasReadme || d.HasDocs || d.IsEkaRepo {
		t.Errorf("empty dir must have no features: %+v", d)
	}
	if len(d.ConfigFiles) != 0 {
		t.Errorf("empty dir must have no config files: %v", d.ConfigFiles)
	}
}

func TestDiscoverGitRepoInTarget(t *testing.T) {
	dir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dir, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}
	d, err := Discover(dir)
	if err != nil {
		t.Fatal(err)
	}
	if !d.IsGitRepo {
		t.Error("target with .git must be detected as a git repository")
	}
	if d.GitRoot != dir {
		t.Errorf("GitRoot = %q, want %q", d.GitRoot, dir)
	}
}

func TestDiscoverGitRepoInAncestor(t *testing.T) {
	parent := t.TempDir()
	if err := os.Mkdir(filepath.Join(parent, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}
	child := filepath.Join(parent, "sub", "deeper")
	if err := os.MkdirAll(child, 0o755); err != nil {
		t.Fatal(err)
	}
	d, err := Discover(child)
	if err != nil {
		t.Fatal(err)
	}
	if !d.IsGitRepo {
		t.Error(".git in an ancestor must be detected")
	}
	if d.GitRoot != parent {
		t.Errorf("GitRoot = %q, want %q", d.GitRoot, parent)
	}
}

func TestDiscoverGitWorktreeFile(t *testing.T) {
	dir := t.TempDir()
	// A worktree has .git as a file, not a directory.
	if err := os.WriteFile(filepath.Join(dir, ".git"), []byte("gitdir: /x"), 0o644); err != nil {
		t.Fatal(err)
	}
	d, err := Discover(dir)
	if err != nil {
		t.Fatal(err)
	}
	if !d.IsGitRepo {
		t.Error("a .git file must also count as a git repository")
	}
}

func TestDiscoverReadme(t *testing.T) {
	for _, name := range []string{"README.md", "README"} {
		dir := t.TempDir()
		if err := os.WriteFile(filepath.Join(dir, name), []byte("# hi"), 0o644); err != nil {
			t.Fatal(err)
		}
		d, err := Discover(dir)
		if err != nil {
			t.Fatal(err)
		}
		if !d.HasReadme || d.ReadmePath != name {
			t.Errorf("name %s: HasReadme=%v ReadmePath=%q", name, d.HasReadme, d.ReadmePath)
		}
	}
}

func TestDiscoverEkaRepoMarkers(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "docs", "operating"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(dir, "docs", "exchange"), 0o755); err != nil {
		t.Fatal(err)
	}
	os.WriteFile(filepath.Join(dir, "docs", "exchange", "validation.md"), []byte("v"), 0o644)
	os.WriteFile(filepath.Join(dir, "docs", "exchange", "transfer.md"), []byte("t"), 0o644)

	d, err := Discover(dir)
	if err != nil {
		t.Fatal(err)
	}
	if !d.HasDocs || !d.IsEkaRepo {
		t.Errorf("full marker set must be detected: %+v", d)
	}

	// Remove one marker: no longer an EKA repo.
	os.Remove(filepath.Join(dir, "docs", "exchange", "transfer.md"))
	d, err = Discover(dir)
	if err != nil {
		t.Fatal(err)
	}
	if d.IsEkaRepo {
		t.Error("missing transfer.md must disqualify the EKA marker")
	}
}

func TestDiscoverConfigFiles(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{".gitignore", ".editorconfig", ".eka.json", "eka.yaml", "notes.txt"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	d, err := Discover(dir)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{".editorconfig", ".eka.json", ".gitignore", "eka.yaml"}
	if !reflect.DeepEqual(d.ConfigFiles, want) {
		t.Errorf("ConfigFiles = %v, want %v", d.ConfigFiles, want)
	}
}

func TestDiscoverTargetIsFile(t *testing.T) {
	file := filepath.Join(t.TempDir(), "f")
	if err := os.WriteFile(file, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	d, err := Discover(file)
	if err != nil {
		t.Fatal(err)
	}
	if !d.Exists || d.IsDir {
		t.Errorf("file target: Exists=%v IsDir=%v", d.Exists, d.IsDir)
	}
}

// --- Wizard: namespace rules -------------------------------------------

func TestIsValidNamespace(t *testing.T) {
	valid := []string{"eka", "my-project", "a1b2", "x-y-9", "-edge"}
	for _, s := range valid {
		if !isValidNamespace(s) {
			t.Errorf("isValidNamespace(%q) = false, want true", s)
		}
	}
	invalid := []string{"", "My-Project", "my_project", "my/project", "my:project", "my project", "my.project"}
	for _, s := range invalid {
		if isValidNamespace(s) {
			t.Errorf("isValidNamespace(%q) = true, want false", s)
		}
	}
}

func TestSanitizeNamespace(t *testing.T) {
	cases := map[string]string{
		"My Project":     "my-project",
		"Foo/Bar:Baz":    "foo-bar-baz",
		"  spaced  out ": "spaced-out",
		"already-good":   "already-good",
		"UPPER":          "upper",
		"a--b":           "a-b",
		"---":            "",
		"":               "",
		"123":            "123",
	}
	for in, want := range cases {
		if got := sanitizeNamespace(in); got != want {
			t.Errorf("sanitizeNamespace(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestDefaultNamespace(t *testing.T) {
	if got := defaultNamespace("my-project"); got != "my-project" {
		t.Errorf("valid name must pass through, got %q", got)
	}
	if got := defaultNamespace("My Project"); got != "my-project" {
		t.Errorf("invalid name must be sanitized, got %q", got)
	}
	if got := defaultNamespace("///"); got != fallbackName {
		t.Errorf("unusable name must fall back, got %q", got)
	}
}

// --- Wizard: adaptivity -------------------------------------------------

func TestNeededQuestions(t *testing.T) {
	base := &Discovery{BaseName: "my-project", HasReadme: true, IsGitRepo: true, GitAvailable: true}
	qs := NeededQuestions(base)
	for _, q := range qs {
		if q.Kind == QProjectName || q.Kind == QReadme || q.Kind == QGit || q.Kind == QNamespace {
			t.Errorf("question %s must not be asked for an initialized-looking dir", q.Kind)
		}
	}
	// Description is always asked interactively.
	if len(qs) != 1 || qs[0].Kind != QDescription {
		t.Errorf("only the description question expected, got %v", qs)
	}
}

func TestNeededQuestionsProjectName(t *testing.T) {
	for _, base := range []string{"", "/"} {
		qs := NeededQuestions(&Discovery{BaseName: base, HasReadme: true, IsGitRepo: true})
		found := false
		for _, q := range qs {
			if q.Kind == QProjectName {
				found = true
			}
		}
		if !found {
			t.Errorf("BaseName %q: project name question must be asked", base)
		}
	}
}

func TestNeededQuestionsNamespace(t *testing.T) {
	// Base name is not a valid namespace → namespace question is asked.
	qs := NeededQuestions(&Discovery{BaseName: "My Project", HasReadme: true, IsGitRepo: true})
	found := false
	for _, q := range qs {
		if q.Kind == QNamespace {
			found = true
		}
	}
	if !found {
		t.Error("invalid base name must trigger the namespace question")
	}
}

func TestNeededQuestionsReadme(t *testing.T) {
	with := NeededQuestions(&Discovery{BaseName: "x", HasReadme: true, IsGitRepo: true})
	for _, q := range with {
		if q.Kind == QReadme {
			t.Error("readme question must be absent when README exists")
		}
	}
	without := NeededQuestions(&Discovery{BaseName: "x", HasReadme: false, IsGitRepo: true})
	found := false
	for _, q := range without {
		if q.Kind == QReadme {
			found = true
		}
	}
	if !found {
		t.Error("readme question must be asked when README is absent")
	}
}

func TestNeededQuestionsGit(t *testing.T) {
	// No question when already a git repository.
	if has := hasQuestion(NeededQuestions(&Discovery{BaseName: "x", HasReadme: true, IsGitRepo: true, GitAvailable: true}), QGit); has {
		t.Error("git question must be absent inside an existing git repository")
	}
	// No question when git is not available.
	if has := hasQuestion(NeededQuestions(&Discovery{BaseName: "x", HasReadme: true, IsGitRepo: false, GitAvailable: false}), QGit); has {
		t.Error("git question must be absent when git is unavailable")
	}
	// Question asked only when both conditions allow it.
	if !hasQuestion(NeededQuestions(&Discovery{BaseName: "x", HasReadme: true, IsGitRepo: false, GitAvailable: true}), QGit) {
		t.Error("git question must be asked when git is available and no repo exists")
	}
}

func hasQuestion(qs []Question, kind QuestionKind) bool {
	for _, q := range qs {
		if q.Kind == kind {
			return true
		}
	}
	return false
}

func TestDefaultAnswers(t *testing.T) {
	d := &Discovery{BaseName: "my-project", HasReadme: false, IsGitRepo: false, GitAvailable: true}
	a := DefaultAnswers(d)
	if a.ProjectName != "my-project" || a.Namespace != "my-project" {
		t.Errorf("defaults: %+v", a)
	}
	if !a.GenerateReadme {
		t.Error("README must be generated when absent")
	}
	if a.InitGit {
		t.Error("non-interactive runs must never init git")
	}
	if a.Interactive {
		t.Error("defaults are non-interactive")
	}

	// Unusable base name falls back deterministically.
	a = DefaultAnswers(&Discovery{BaseName: "", HasReadme: true, IsGitRepo: true})
	if a.ProjectName != fallbackName || a.Namespace != fallbackName {
		t.Errorf("fallback defaults: %+v", a)
	}
	if a.GenerateReadme {
		t.Error("README must not be generated when one exists")
	}
}

// --- Wizard: interactive Ask --------------------------------------------

func TestAskUsesPipedAnswers(t *testing.T) {
	d := &Discovery{BaseName: "My Project", HasReadme: false, IsGitRepo: false, GitAvailable: true}
	// "My Project" is a usable display name, so only the namespace,
	// description, README and git questions are asked. Answers: namespace
	// (first invalid, then valid), description, README no, git no.
	input := "bad_namespace\ncool-product\nA description\nn\nn\n"
	var out strings.Builder
	a, err := Ask(d, strings.NewReader(input), &out)
	if err != nil {
		t.Fatal(err)
	}
	// The base name is usable, so the project name defaults silently.
	if a.ProjectName != "My Project" {
		t.Errorf("ProjectName = %q, want %q", a.ProjectName, "My Project")
	}
	if a.Namespace != "cool-product" {
		t.Errorf("Namespace = %q, want cool-product (invalid answer must be re-prompted)", a.Namespace)
	}
	if a.Description != "A description" {
		t.Errorf("Description = %q", a.Description)
	}
	if a.GenerateReadme {
		t.Error("GenerateReadme must honor the 'n' answer")
	}
	if a.InitGit {
		t.Error("InitGit must honor the 'n' answer")
	}
	if !strings.Contains(out.String(), "invalid namespace") {
		t.Errorf("re-prompt message expected in output:\n%s", out.String())
	}
}

func TestAskProjectNameQuestion(t *testing.T) {
	// An unusable base name triggers the project name question; a project
	// name that is not a valid namespace then triggers the namespace
	// question.
	d := &Discovery{BaseName: "", HasReadme: true, IsGitRepo: true, GitAvailable: true}
	input := "Cool Product\ncool-product\n\n"
	var out strings.Builder
	a, err := Ask(d, strings.NewReader(input), &out)
	if err != nil {
		t.Fatal(err)
	}
	if a.ProjectName != "Cool Product" {
		t.Errorf("ProjectName = %q, want %q", a.ProjectName, "Cool Product")
	}
	if a.Namespace != "cool-product" {
		t.Errorf("Namespace = %q, want cool-product", a.Namespace)
	}
}

func TestAskEofFallsBackToDefaults(t *testing.T) {
	d := &Discovery{BaseName: "my-project", HasReadme: true, IsGitRepo: true, GitAvailable: true}
	a, err := Ask(d, strings.NewReader(""), &strings.Builder{})
	if err != nil {
		t.Fatal(err)
	}
	want := DefaultAnswers(d)
	if a.ProjectName != want.ProjectName || a.Namespace != want.Namespace {
		t.Errorf("EOF must yield defaults, got %+v", a)
	}
}
