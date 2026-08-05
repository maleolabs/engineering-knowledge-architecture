package bootstrap

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// --- BuildPlan ----------------------------------------------------------

func TestBuildPlanNewDirectory(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "fresh")
	d, err := Discover(dir)
	if err != nil {
		t.Fatal(err)
	}
	plan := BuildPlan(dir, d, DefaultAnswers(d))

	// Grouping order: dirs, files, git, validation.
	lastDir, lastFile, lastGit := -1, -1, -1
	validateIdx := -1
	for i, a := range plan {
		switch a.Kind {
		case ActionCreateDir:
			lastDir = i
		case ActionCreateFile, ActionGenerateReadme:
			if i < lastDir {
				t.Errorf("file action %s before dir actions complete", a.Path)
			}
			lastFile = i
		case ActionGitInit, ActionGitSkip:
			if i < lastFile {
				t.Errorf("git action before file actions complete")
			}
			lastGit = i
		case ActionValidate:
			if i < lastGit {
				t.Errorf("validate action before git action")
			}
			validateIdx = i
		}
	}
	if validateIdx != len(plan)-1 {
		t.Errorf("validate must be the final action, got index %d of %d", validateIdx, len(plan))
	}

	// The target dir itself is created (sentinel path ".", target in Detail).
	if plan[0].Kind != ActionCreateDir || plan[0].Path != "." || plan[0].Detail != dir {
		t.Errorf("first action must create the target dir: %+v", plan[0])
	}
	if plan[0].String() != "create dir: "+dir {
		t.Errorf("target-dir action renders wrong: %s", plan[0].String())
	}
	// README is generated, not raw-copied.
	for _, a := range plan {
		if a.Path == "README.md" && a.Kind != ActionGenerateReadme {
			t.Errorf("README.md must be generated, got %s", a.Kind)
		}
	}
	// Every file action carries content.
	for _, a := range plan {
		switch a.Kind {
		case ActionCreateFile, ActionGenerateReadme, ActionOverwriteConfirm:
			if len(a.Content) == 0 {
				t.Errorf("action %s has no content", a.String())
			}
		}
	}
}

func TestBuildPlanDeterministic(t *testing.T) {
	dir := t.TempDir()
	d, err := Discover(dir)
	if err != nil {
		t.Fatal(err)
	}
	a := DefaultAnswers(d)
	p1 := BuildPlan(dir, d, a)
	p2 := BuildPlan(dir, d, a)
	if len(p1) != len(p2) {
		t.Fatalf("plan lengths differ: %d vs %d", len(p1), len(p2))
	}
	for i := range p1 {
		if p1[i].String() != p2[i].String() {
			t.Errorf("plan differs at %d: %s vs %s", i, p1[i], p2[i])
		}
	}
}

func TestBuildPlanExistingEkaRepo(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "docs", "operating"), 0o755)
	os.MkdirAll(filepath.Join(dir, "docs", "exchange"), 0o755)
	os.WriteFile(filepath.Join(dir, "docs", "exchange", "validation.md"), []byte("v"), 0o644)
	os.WriteFile(filepath.Join(dir, "docs", "exchange", "transfer.md"), []byte("t"), 0o644)
	d, err := Discover(dir)
	if err != nil {
		t.Fatal(err)
	}
	plan := BuildPlan(dir, d, DefaultAnswers(d))
	if len(plan) != 2 {
		t.Fatalf("existing EKA repo must plan reuse+validate only, got %d actions", len(plan))
	}
	if plan[0].Kind != ActionReuse {
		t.Errorf("first action must be reuse, got %+v", plan[0])
	}
	if plan[1].Kind != ActionValidate {
		t.Errorf("second action must be validate, got %+v", plan[1])
	}
}

func TestBuildPlanOverwriteConfirm(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "docs", "exchange"), 0o755)
	os.WriteFile(filepath.Join(dir, "docs", "exchange", "validation.md"), []byte("custom"), 0o644)
	d, err := Discover(dir)
	if err != nil {
		t.Fatal(err)
	}
	plan := BuildPlan(dir, d, DefaultAnswers(d))
	found := false
	for _, a := range plan {
		if a.Path == "docs/exchange/validation.md" && a.Kind == ActionOverwriteConfirm {
			found = true
		}
	}
	if !found {
		t.Error("differing existing file must be planned as overwrite-confirm")
	}
}

func TestBuildPlanReuseIdenticalFile(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "docs", "exchange"), 0o755)
	content, err := readSkeletonFile("docs/exchange/validation.md")
	if err != nil {
		t.Fatal(err)
	}
	os.WriteFile(filepath.Join(dir, "docs", "exchange", "validation.md"), content, 0o644)
	d, err := Discover(dir)
	if err != nil {
		t.Fatal(err)
	}
	plan := BuildPlan(dir, d, DefaultAnswers(d))
	for _, a := range plan {
		if a.Path == "docs/exchange/validation.md" && a.Kind != ActionReuse {
			t.Errorf("identical file must be planned as reuse, got %s", a.Kind)
		}
	}
}

func TestBuildPlanGitSkipReasons(t *testing.T) {
	cases := []struct {
		name string
		d    *Discovery
		a    Answers
		want string
	}{
		{"existing repo", &Discovery{BaseName: "x", IsGitRepo: true}, Answers{}, "skipped (already a git repository)"},
		{"no git binary", &Discovery{BaseName: "x", GitAvailable: false}, Answers{InitGit: true}, "skipped (git not available)"},
		{"declined", &Discovery{BaseName: "x", GitAvailable: true}, Answers{InitGit: false, Interactive: true}, "skipped (declined)"},
		{"non-interactive", &Discovery{BaseName: "x", GitAvailable: true}, Answers{InitGit: false}, "skipped (non-interactive mode)"},
	}
	for _, tc := range cases {
		plan := BuildPlan("t", tc.d, tc.a)
		gitLine := ""
		for _, a := range plan {
			if a.Kind == ActionGitSkip {
				gitLine = a.Detail
			}
		}
		if gitLine != tc.want {
			t.Errorf("%s: git detail = %q, want %q", tc.name, gitLine, tc.want)
		}
	}
}

func TestBuildPlanGitInitPlanned(t *testing.T) {
	d := &Discovery{BaseName: "x", GitAvailable: true}
	a := Answers{InitGit: true, Interactive: true}
	plan := BuildPlan("t", d, a)
	found := false
	for _, act := range plan {
		if act.Kind == ActionGitInit {
			found = true
		}
	}
	if !found {
		t.Error("git init must be planned when requested and possible")
	}
}

func TestActionString(t *testing.T) {
	cases := []struct {
		a    Action
		want string
	}{
		{Action{Kind: ActionCreateDir, Path: "docs/"}, "create dir: docs/"},
		{Action{Kind: ActionCreateFile, Path: "docs/x.md"}, "create file: docs/x.md (from skeleton)"},
		{Action{Kind: ActionGenerateReadme, Path: "README.md"}, "generate file: README.md (from skeleton)"},
		{Action{Kind: ActionReuse, Path: "docs/README.md"}, "reuse: docs/README.md"},
		{Action{Kind: ActionOverwriteConfirm, Path: "README.md"}, "overwrite confirm: README.md"},
		{Action{Kind: ActionGitInit, Path: "t"}, "git init: t"},
		{Action{Kind: ActionGitSkip, Detail: "skipped (already a git repository)"}, "git init: skipped (already a git repository)"},
		{Action{Kind: ActionValidate, Path: "t", Detail: "after generation"}, "validate: t after generation"},
	}
	for _, tc := range cases {
		if got := tc.a.String(); got != tc.want {
			t.Errorf("Action.String() = %q, want %q", got, tc.want)
		}
	}
}

// --- generatedReadme ----------------------------------------------------

func TestGeneratedReadmeHeading(t *testing.T) {
	got := string(generatedReadme("my-project"))
	lines := strings.Split(got, "\n")
	if lines[0] != "# my-project" {
		t.Errorf("first line = %q, want %q", lines[0], "# my-project")
	}
	if !strings.Contains(got, "Tanggal templat: 2026-08-05") {
		t.Error("template date must be preserved verbatim")
	}
	// Byte-identical to the template apart from the heading line.
	want, err := readSkeletonFile("README.md")
	if err != nil {
		t.Fatal(err)
	}
	orig := strings.Split(string(want), "\n")
	orig[0] = "# my-project"
	if got != strings.Join(orig, "\n") {
		t.Error("generated README must differ from the template only in the heading line")
	}
}

// --- Apply --------------------------------------------------------------

func applyOpts(interactive bool) ApplyOptions {
	return ApplyOptions{
		Interactive: interactive,
		Stdin:       strings.NewReader(""),
		Stdout:      &bytes.Buffer{},
		Stderr:      &bytes.Buffer{},
	}
}

func TestApplyNonInteractiveSkipsOverwrites(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "docs", "exchange"), 0o755)
	os.WriteFile(filepath.Join(dir, "docs", "exchange", "validation.md"), []byte("custom"), 0o644)
	d, err := Discover(dir)
	if err != nil {
		t.Fatal(err)
	}
	plan := BuildPlan(dir, d, DefaultAnswers(d))
	res, err := Apply(dir, plan, applyOpts(false))
	if err != nil {
		t.Fatal(err)
	}
	content, _ := os.ReadFile(filepath.Join(dir, "docs", "exchange", "validation.md"))
	if string(content) != "custom" {
		t.Error("non-interactive apply must never replace an existing file")
	}
	found := false
	for _, p := range res.SkippedFiles {
		if p == "docs/exchange/validation.md" {
			found = true
		}
	}
	if !found {
		t.Errorf("skipped file must be reported, got %v", res.SkippedFiles)
	}
}

func TestApplyInteractiveOverwrite(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "docs", "exchange"), 0o755)
	os.WriteFile(filepath.Join(dir, "docs", "exchange", "validation.md"), []byte("custom"), 0o644)
	d, err := Discover(dir)
	if err != nil {
		t.Fatal(err)
	}
	plan := BuildPlan(dir, d, DefaultAnswers(d))
	var asked []string
	opts := applyOpts(true)
	opts.ConfirmOverwrite = func(path string) (bool, error) {
		asked = append(asked, path)
		return true, nil
	}
	res, err := Apply(dir, plan, opts)
	if err != nil {
		t.Fatal(err)
	}
	if len(asked) == 0 {
		t.Fatal("overwrite confirmation must be requested interactively")
	}
	content, _ := os.ReadFile(filepath.Join(dir, "docs", "exchange", "validation.md"))
	if bytes.Equal(content, []byte("custom")) {
		t.Error("confirmed overwrite must replace the file")
	}
	found := false
	for _, p := range res.OverwrittenFiles {
		if p == "docs/exchange/validation.md" {
			found = true
		}
	}
	if !found {
		t.Errorf("overwritten file must be reported, got %v", res.OverwrittenFiles)
	}
}

func TestApplyInteractiveDecline(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "docs", "exchange"), 0o755)
	os.WriteFile(filepath.Join(dir, "docs", "exchange", "validation.md"), []byte("custom"), 0o644)
	d, err := Discover(dir)
	if err != nil {
		t.Fatal(err)
	}
	plan := BuildPlan(dir, d, DefaultAnswers(d))
	opts := applyOpts(true)
	opts.ConfirmOverwrite = func(path string) (bool, error) { return false, nil }
	res, err := Apply(dir, plan, opts)
	if err != nil {
		t.Fatal(err)
	}
	content, _ := os.ReadFile(filepath.Join(dir, "docs", "exchange", "validation.md"))
	if string(content) != "custom" {
		t.Error("declined overwrite must preserve the file")
	}
	if len(res.OverwrittenFiles) != 0 {
		t.Errorf("no file may be reported overwritten, got %v", res.OverwrittenFiles)
	}
}

func TestApplyGitInit(t *testing.T) {
	dir := t.TempDir()
	plan := []Action{
		{Kind: ActionGitInit, Path: dir},
		{Kind: ActionValidate, Path: dir},
	}
	called := false
	opts := applyOpts(false)
	opts.GitInit = func(d string, stdout, stderr io.Writer) error {
		called = true
		if d != dir {
			t.Errorf("git init dir = %q, want %q", d, dir)
		}
		return nil
	}
	res, err := Apply(dir, plan, opts)
	if err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Error("git init must be executed when planned")
	}
	if res.GitStatus != "initialized" {
		t.Errorf("GitStatus = %q, want initialized", res.GitStatus)
	}
}

func TestApplyGitInitFailureIsWarning(t *testing.T) {
	dir := t.TempDir()
	plan := []Action{{Kind: ActionGitInit, Path: dir}}
	var errb bytes.Buffer
	opts := applyOpts(false)
	opts.Stderr = &errb
	opts.GitInit = func(d string, stdout, stderr io.Writer) error {
		return errors.New("git exploded")
	}
	res, err := Apply(dir, plan, opts)
	if err != nil {
		t.Fatalf("a failed git init must not fail the generation: %v", err)
	}
	if !strings.Contains(res.GitStatus, "failed") {
		t.Errorf("GitStatus = %q, want failed", res.GitStatus)
	}
	if !strings.Contains(errb.String(), "warning") {
		t.Errorf("stderr must carry a warning, got %q", errb.String())
	}
}

func TestApplyGitStatusExisting(t *testing.T) {
	plan := []Action{{Kind: ActionGitSkip, Detail: "skipped (already a git repository)"}}
	res, err := Apply(t.TempDir(), plan, applyOpts(false))
	if err != nil {
		t.Fatal(err)
	}
	if res.GitStatus != "existing" {
		t.Errorf("GitStatus = %q, want existing", res.GitStatus)
	}
}

func TestApplyCreateDirForFileConflict(t *testing.T) {
	dir := t.TempDir()
	// "docs" exists as a file: creating the directory must fail cleanly.
	os.WriteFile(filepath.Join(dir, "docs"), []byte("x"), 0o644)
	plan := []Action{{Kind: ActionCreateDir, Path: "docs"}}
	if _, err := Apply(dir, plan, applyOpts(false)); err == nil {
		t.Error("conflicting file must produce an error")
	}
}
