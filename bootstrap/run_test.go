package bootstrap

import (
	"bytes"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/maleolabs/engineering-knowledge-architecture"
	"github.com/maleolabs/engineering-knowledge-architecture/conformance"
)

// runOpts builds Options with captured buffers and a non-interactive stdin.
func runOpts(target string, stdin string) (Options, *bytes.Buffer, *bytes.Buffer) {
	var out, errb bytes.Buffer
	return Options{
		Target: target,
		Stdin:  strings.NewReader(stdin),
		Stdout: &out,
		Stderr: &errb,
	}, &out, &errb
}

// contains reports whether xs contains s.
func contains(xs []string, s string) bool {
	for _, x := range xs {
		if x == s {
			return true
		}
	}
	return false
}

// skeletonFiles returns the sorted relative paths of the embedded skeleton.
func skeletonFiles(t *testing.T) []string {
	t.Helper()
	sub, err := fs.Sub(skeletonembed.FS, "skeleton")
	if err != nil {
		t.Fatal(err)
	}
	var files []string
	fs.WalkDir(sub, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files
}

// walkFiles returns sorted relative file paths under dir.
func walkFiles(t *testing.T, dir string) []string {
	t.Helper()
	var files []string
	filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			rel, err := filepath.Rel(dir, path)
			if err != nil {
				return err
			}
			files = append(files, filepath.ToSlash(rel))
		}
		return nil
	})
	return files
}

// Scenario 1: empty directory init. Answers are piped via strings.Reader,
// but a non-terminal stdin means the wizard is skipped: the piped answers
// must be ignored and deterministic defaults used.
func TestRunEmptyDirectory(t *testing.T) {
	dir := t.TempDir()
	// Answers that would change the outcome if they were consumed.
	opts, _, _ := runOpts(dir, "piped-name\npiped-ns\npiped description\ny\ny\n")
	outcome, err := Run(opts)
	if err != nil {
		t.Fatal(err)
	}
	if outcome.ProjectName != filepath.Base(dir) {
		t.Errorf("ProjectName = %q, want %q (piped answers must be ignored)", outcome.ProjectName, filepath.Base(dir))
	}
	if outcome.Namespace != filepath.Base(dir) {
		t.Errorf("Namespace = %q, want %q", outcome.Namespace, filepath.Base(dir))
	}
	if outcome.RepoType != "existing-dir" {
		t.Errorf("RepoType = %q, want existing-dir", outcome.RepoType)
	}
	if outcome.Report == nil || !outcome.Report.Pass() {
		t.Fatalf("generated repo must validate, report: %+v", outcome.Report)
	}
	want := skeletonFiles(t)
	if len(outcome.CreatedFiles) != len(want) {
		t.Errorf("created %d files, want %d", len(outcome.CreatedFiles), len(want))
	}
	if len(outcome.CreatedDirs) != len(wantDirs(t)) {
		t.Errorf("created %d dirs, want %d", len(outcome.CreatedDirs), len(wantDirs(t)))
	}
	if outcome.GitStatus != "skipped (non-interactive mode)" {
		t.Errorf("GitStatus = %q, want skipped (non-interactive mode)", outcome.GitStatus)
	}
	if len(outcome.SkippedFiles) != 0 || len(outcome.OverwrittenFiles) != 0 {
		t.Errorf("empty dir must not skip or overwrite anything: %+v", outcome)
	}
}

func wantDirs(t *testing.T) []string {
	t.Helper()
	sub, err := fs.Sub(skeletonembed.FS, "skeleton")
	if err != nil {
		t.Fatal(err)
	}
	var dirs []string
	fs.WalkDir(sub, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path != "." && d.IsDir() {
			dirs = append(dirs, path)
		}
		return nil
	})
	return dirs
}

// Scenario 2: existing non-empty directory without EKA is adopted; custom
// files survive and no existing file is replaced silently.
func TestRunAdoptsExistingDirectory(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "notes.txt"), []byte("keep me"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(dir, "docs", "exchange"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "docs", "exchange", "validation.md"), []byte("custom"), 0o644); err != nil {
		t.Fatal(err)
	}

	opts, _, _ := runOpts(dir, "")
	outcome, err := Run(opts)
	if err != nil {
		t.Fatal(err)
	}
	if outcome.RepoType != "existing-dir" {
		t.Errorf("RepoType = %q, want existing-dir", outcome.RepoType)
	}
	if data, _ := os.ReadFile(filepath.Join(dir, "notes.txt")); string(data) != "keep me" {
		t.Error("custom file must be preserved")
	}
	if data, _ := os.ReadFile(filepath.Join(dir, "docs", "exchange", "validation.md")); string(data) != "custom" {
		t.Error("differing existing file must not be replaced in non-interactive mode")
	}
	if !contains(outcome.SkippedFiles, "docs/exchange/validation.md") {
		t.Errorf("skipped file must be reported, got %v", outcome.SkippedFiles)
	}
	if outcome.Report == nil || !outcome.Report.Pass() {
		t.Errorf("adopted repo must still validate: %+v", outcome.Report)
	}
}

// Scenario 3: existing git repository → no git question, no git init.
func TestRunExistingGitRepository(t *testing.T) {
	dir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dir, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}
	gitCalled := false
	opts, _, _ := runOpts(dir, "")
	opts.GitInit = func(d string, stdout, stderr io.Writer) error {
		gitCalled = true
		return nil
	}
	outcome, err := Run(opts)
	if err != nil {
		t.Fatal(err)
	}
	if gitCalled {
		t.Error("git init must not run inside an existing git repository")
	}
	if outcome.GitStatus != "existing" {
		t.Errorf("GitStatus = %q, want existing", outcome.GitStatus)
	}
	// The wizard must not offer the git question.
	d, err := Discover(dir)
	if err != nil {
		t.Fatal(err)
	}
	if hasQuestion(NeededQuestions(d), QGit) {
		t.Error("git question must be absent inside an existing git repository")
	}
	if outcome.Report == nil || !outcome.Report.Pass() {
		t.Error("repo must validate")
	}
}

// Scenario 4: existing README is preserved; an identical one is reused.
func TestRunPreservesExistingReadme(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("# Mine\n\ncustom\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	opts, _, _ := runOpts(dir, "")
	outcome, err := Run(opts)
	if err != nil {
		t.Fatal(err)
	}
	if data, _ := os.ReadFile(filepath.Join(dir, "README.md")); string(data) != "# Mine\n\ncustom\n" {
		t.Error("existing README must be preserved")
	}
	// The README exists, so nothing is planned for it: no write, no skip.
	for _, p := range append(append([]string{}, outcome.CreatedFiles...), outcome.OverwrittenFiles...) {
		if p == "README.md" {
			t.Errorf("existing README must not be written, got in created/overwritten: %v", p)
		}
	}
	// The README question must not be offered either.
	d, err := Discover(dir)
	if err != nil {
		t.Fatal(err)
	}
	if hasQuestion(NeededQuestions(d), QReadme) {
		t.Error("readme question must be absent when README exists")
	}
}

func TestRunReusesIdenticalReadme(t *testing.T) {
	dir := t.TempDir()
	generated := generatedReadme(filepath.Base(dir))
	if err := os.WriteFile(filepath.Join(dir, "README.md"), generated, 0o644); err != nil {
		t.Fatal(err)
	}
	opts, _, _ := runOpts(dir, "")
	outcome, err := Run(opts)
	if err != nil {
		t.Fatal(err)
	}
	if data, _ := os.ReadFile(filepath.Join(dir, "README.md")); !bytes.Equal(data, generated) {
		t.Error("identical README must be left untouched")
	}
	if len(outcome.CreatedFiles) != len(skeletonFiles(t))-1 {
		t.Errorf("all skeleton docs files must still be created, got %d (want %d)", len(outcome.CreatedFiles), len(skeletonFiles(t))-1)
	}
}

// Scenario 5: existing docs directory is merged — skeleton files added,
// custom files kept.
func TestRunMergesExistingDocs(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "docs", "custom"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "docs", "custom", "mine.md"), []byte("# mine"), 0o644); err != nil {
		t.Fatal(err)
	}
	opts, _, _ := runOpts(dir, "")
	outcome, err := Run(opts)
	if err != nil {
		t.Fatal(err)
	}
	if data, _ := os.ReadFile(filepath.Join(dir, "docs", "custom", "mine.md")); string(data) != "# mine" {
		t.Error("custom docs file must be preserved")
	}
	if _, err := os.Stat(filepath.Join(dir, "docs", "exchange", "validation.md")); err != nil {
		t.Error("skeleton files must be merged into the existing docs tree")
	}
	if outcome.Report == nil || !outcome.Report.Pass() {
		t.Error("merged repo must validate")
	}
}

// Scenario 6: existing EKA repository → reuse + validate, zero writes.
func TestRunExistingEkaRepo(t *testing.T) {
	dir := t.TempDir()
	makeEkaRepo(t, dir)
	if err := os.Mkdir(filepath.Join(dir, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}
	opts, out, _ := runOpts(dir, "")
	outcome, err := Run(opts)
	if err != nil {
		t.Fatal(err)
	}
	if !outcome.AlreadyInitialized {
		t.Error("existing EKA repo must be detected as already initialized")
	}
	if outcome.RepoType != "existing-eka" {
		t.Errorf("RepoType = %q, want existing-eka", outcome.RepoType)
	}
	if len(outcome.Plan) != 2 {
		t.Errorf("plan must be reuse+validate only, got %d actions", len(outcome.Plan))
	}
	if len(outcome.CreatedDirs) != 0 || len(outcome.CreatedFiles) != 0 ||
		len(outcome.OverwrittenFiles) != 0 || len(outcome.SkippedFiles) != 0 {
		t.Errorf("existing EKA repo must receive zero writes: %+v", outcome)
	}
	if outcome.GitStatus != "existing" {
		t.Errorf("GitStatus = %q, want existing (discovered .git)", outcome.GitStatus)
	}
	if outcome.Report == nil || !outcome.Report.Pass() {
		t.Error("existing repo must validate")
	}
	// No wizard prompts may be emitted for an already-initialized repo.
	if out.Len() != 0 {
		t.Errorf("already-initialized repo must not prompt, got output:\n%s", out.String())
	}
}

func TestRunExistingEkaRepoWithoutGit(t *testing.T) {
	dir := t.TempDir()
	makeEkaRepo(t, dir)
	opts, _, _ := runOpts(dir, "")
	outcome, err := Run(opts)
	if err != nil {
		t.Fatal(err)
	}
	if outcome.GitStatus != "skipped (no git action planned)" {
		t.Errorf("GitStatus = %q, want skipped (no git action planned)", outcome.GitStatus)
	}
}

func makeEkaRepo(t *testing.T, dir string) {
	t.Helper()
	os.MkdirAll(filepath.Join(dir, "docs", "operating"), 0o755)
	os.MkdirAll(filepath.Join(dir, "docs", "exchange"), 0o755)
	os.WriteFile(filepath.Join(dir, "docs", "exchange", "validation.md"), []byte("# v"), 0o644)
	os.WriteFile(filepath.Join(dir, "docs", "exchange", "transfer.md"), []byte("# t"), 0o644)
}

// Scenario 7: repeated initialization — the second run is a no-op and the
// repository still validates.
func TestRunTwiceIsNoop(t *testing.T) {
	dir := t.TempDir()
	opts, _, _ := runOpts(dir, "")
	if _, err := Run(opts); err != nil {
		t.Fatal(err)
	}
	before := walkFiles(t, dir)
	second, err := Run(opts)
	if err != nil {
		t.Fatal(err)
	}
	if !second.AlreadyInitialized {
		t.Error("second run must detect the existing EKA repository")
	}
	if len(second.CreatedFiles) != 0 || len(second.OverwrittenFiles) != 0 || len(second.SkippedFiles) != 0 {
		t.Errorf("second run must write nothing: %+v", second)
	}
	after := walkFiles(t, dir)
	if !reflect.DeepEqual(before, after) {
		t.Errorf("second run changed the tree:\nbefore: %v\nafter:  %v", before, after)
	}
	if second.Report == nil || !second.Report.Pass() {
		t.Error("repo must still validate after the second run")
	}
}

// Scenario 8: dry-run mode — plan printed, nothing written, exit path 0.
func TestRunDryRun(t *testing.T) {
	dir := t.TempDir()
	opts, _, _ := runOpts(dir, "")
	opts.DryRun = true
	outcome, err := Run(opts)
	if err != nil {
		t.Fatal(err)
	}
	if !outcome.DryRun {
		t.Error("outcome must mark the dry run")
	}
	if outcome.Report != nil {
		t.Error("dry run must not validate")
	}
	if len(outcome.Plan) == 0 {
		t.Fatal("dry run must produce a plan")
	}
	// Stable ordering: the plan starts with dirs and ends with validate.
	if outcome.Plan[len(outcome.Plan)-1].Kind != ActionValidate {
		t.Errorf("plan must end with validate, got %+v", outcome.Plan[len(outcome.Plan)-1])
	}
	first := outcome.Plan[0]
	if first.Kind != ActionCreateDir && first.Kind != ActionGenerateReadme && first.Kind != ActionCreateFile {
		t.Errorf("plan must start with a dir/file action, got %+v", first)
	}
	// No writes: the target still has no files.
	if got := walkFiles(t, dir); len(got) != 0 {
		t.Errorf("dry run must not write anything, found: %v", got)
	}
}

// Scenario 9: failed validation. Unit level: injected failing validator
// surfaces through Run and through RunValidation.
func TestRunValidationFailingValidator(t *testing.T) {
	dir := t.TempDir()
	fail := conformance.Report{
		Root: dir,
		Results: []conformance.Result{{
			Severity: conformance.SeverityError,
			Rule:     conformance.RuleStructural,
			File:     "docs/x.md",
			Message:  "injected failure",
		}},
	}
	opts, _, _ := runOpts(dir, "")
	opts.Validate = func(root string) (*conformance.Report, error) { return &fail, nil }
	outcome, err := Run(opts)
	if err != nil {
		t.Fatalf("a failing validation must not be a Run error: %v", err)
	}
	if outcome.Report == nil || outcome.Report.Pass() {
		t.Error("outcome must carry the failing report")
	}

	// The stage component itself.
	report, err := RunValidation(dir, func(root string) (*conformance.Report, error) { return &fail, nil })
	if err != nil {
		t.Fatal(err)
	}
	if report.Pass() || report.ErrorCount() != 1 {
		t.Errorf("RunValidation must propagate the failing report: %+v", report)
	}
}

// Scenario 10: successful validation on the default path.
func TestRunSuccessfulValidation(t *testing.T) {
	dir := t.TempDir()
	opts, _, _ := runOpts(dir, "")
	outcome, err := Run(opts)
	if err != nil {
		t.Fatal(err)
	}
	if outcome.Report == nil {
		t.Fatal("validation must run after generation")
	}
	if !outcome.Report.Pass() {
		t.Errorf("fresh skeleton must validate, got %d errors", outcome.Report.ErrorCount())
	}
	if outcome.Report.FilesScanned != len(skeletonFiles(t)) {
		t.Errorf("scanned %d files, want %d", outcome.Report.FilesScanned, len(skeletonFiles(t)))
	}
	if outcome.Report.Artifacts != 0 {
		t.Errorf("fresh skeleton must contain no artifacts, got %d", outcome.Report.Artifacts)
	}
}

// --- Mode tests: eka init / eka init . / eka init <name> ----------------

func TestRunTargetCurrentDir(t *testing.T) {
	dir := t.TempDir()
	old, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(old)
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	for _, target := range []string{"", "."} {
		opts, _, _ := runOpts(target, "")
		outcome, err := Run(opts)
		if err != nil {
			t.Fatalf("target %q: %v", target, err)
		}
		if outcome.Report == nil || !outcome.Report.Pass() {
			t.Errorf("target %q: must validate", target)
		}
		// An empty target normalizes to ".".
		want := target
		if want == "" {
			want = "."
		}
		if outcome.Target != want {
			t.Errorf("Target = %q, want %q", outcome.Target, want)
		}
	}
}

func TestRunCreatesNewProjectDir(t *testing.T) {
	parent := t.TempDir()
	target := filepath.Join(parent, "myproj")
	opts, _, _ := runOpts(target, "")
	outcome, err := Run(opts)
	if err != nil {
		t.Fatal(err)
	}
	if outcome.RepoType != "new" {
		t.Errorf("RepoType = %q, want new", outcome.RepoType)
	}
	if _, err := os.Stat(filepath.Join(target, "docs", "README.md")); err != nil {
		t.Errorf("project dir must be created and bootstrapped: %v", err)
	}
	// The target dir itself must not be nested (eka-named/eka-named bug).
	if info, err := os.Stat(filepath.Join(target, "myproj")); err == nil && info.IsDir() {
		t.Error("target dir must not be nested inside itself")
	}
	if outcome.ProjectName != "myproj" {
		t.Errorf("ProjectName = %q, want myproj", outcome.ProjectName)
	}
}

func TestRunTargetIsFileError(t *testing.T) {
	file := filepath.Join(t.TempDir(), "f")
	os.WriteFile(file, []byte("x"), 0o644)
	opts, _, _ := runOpts(file, "")
	if _, err := Run(opts); err == nil {
		t.Error("file target must be an error")
	}
}

// TestRunTargetNamedDocs guards against the target-name/skeleton-dir
// collision: `eka init docs` must create docs/ with the skeleton INSIDE it,
// not misplace the tree at docs/docs.
func TestRunTargetNamedDocs(t *testing.T) {
	parent := t.TempDir()
	target := filepath.Join(parent, "docs")
	opts, _, _ := runOpts(target, "")
	outcome, err := Run(opts)
	if err != nil {
		t.Fatal(err)
	}
	if outcome.RepoType != "new" {
		t.Errorf("RepoType = %q, want new", outcome.RepoType)
	}
	if _, err := os.Stat(filepath.Join(target, "docs", "README.md")); err != nil {
		t.Errorf("skeleton must live at docs/docs for a target named docs: %v", err)
	}
	if _, err := os.Stat(filepath.Join(target, "README.md")); err != nil {
		t.Errorf("README must live at docs/README.md: %v", err)
	}
	if info, err := os.Stat(filepath.Join(target, "docs", "docs")); err == nil && info.IsDir() {
		t.Error("no double-nesting allowed")
	}
	if outcome.Report == nil || !outcome.Report.Pass() {
		t.Error("repo must validate")
	}
}

// --- Non-interactive determinism ----------------------------------------

func TestRunDeterministicAcrossPipedInput(t *testing.T) {
	dir := t.TempDir()
	runOnce := func(stdin string) string {
		opts, out, _ := runOpts(dir, stdin)
		if _, err := Run(opts); err != nil {
			t.Fatal(err)
		}
		return out.String()
	}
	// Different piped answers must not change non-interactive behavior.
	a := runOnce("first-answer\nsecond-answer\n")
	b := runOnce("")
	if a != b {
		t.Error("non-interactive runs must be byte-identical regardless of piped input")
	}
}

// TestRunStdinDevNullNonInteractive is the regression test for the
// char-device bug: /dev/null is a char device, so the old ModeCharDevice
// heuristic misclassified it as interactive, printed prompts to stdout and
// ran `git init` after EOF fell back to the default "y". A true terminal
// check (term.IsTerminal) must classify /dev/null as non-interactive: no
// prompts, deterministic defaults, no .git directory, exit path success.
func TestRunStdinDevNullNonInteractive(t *testing.T) {
	if runtime.GOOS == "windows" {
		// /dev/null does not exist on Windows; the bug only manifests on
		// Unix char devices. The pipe-based test below keeps coverage on
		// Windows-compatible platforms.
		t.Skip("os.Open(/dev/null) is Unix-only")
	}
	devnull, err := os.Open("/dev/null")
	if err != nil {
		t.Skipf("cannot open /dev/null: %v", err)
	}
	defer devnull.Close()

	dir := t.TempDir()
	var out bytes.Buffer
	opts := Options{
		Target: dir,
		Stdin:  devnull,
		Stdout: &out,
		Stderr: io.Discard,
	}
	outcome, err := Run(opts)
	if err != nil {
		t.Fatal(err)
	}
	if outcome.GitStatus != "skipped (non-interactive mode)" {
		t.Errorf("GitStatus = %q, want skipped (non-interactive mode): git init must never run for /dev/null stdin", outcome.GitStatus)
	}
	for _, a := range outcome.Plan {
		if a.Kind == ActionGitInit {
			t.Error("non-interactive run must not plan git init")
		}
	}
	if _, err := os.Stat(filepath.Join(dir, ".git")); !os.IsNotExist(err) {
		t.Errorf(".git must not exist after `eka init < /dev/null`, stat err = %v", err)
	}
	if out.Len() != 0 {
		t.Errorf("no prompts may be printed for /dev/null stdin, got:\n%s", out.String())
	}
	if outcome.Report == nil || !outcome.Report.Pass() {
		t.Fatalf("generated repo must validate, report: %+v", outcome.Report)
	}
}

// TestRunStdinPipeNonInteractive proves that a real *os.File pipe fd (the
// `echo | eka init` shape) stays non-interactive under the fd-based
// terminal check: no prompts, deterministic defaults, no git init.
func TestRunStdinPipeNonInteractive(t *testing.T) {
	pr, pw, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	// Answers that would flip the git default if they were consumed.
	if _, err := pw.WriteString("piped-name\npiped-ns\ndescription\ny\ny\n"); err != nil {
		t.Fatal(err)
	}
	pw.Close()
	defer pr.Close()

	dir := t.TempDir()
	var out bytes.Buffer
	opts := Options{
		Target: dir,
		Stdin:  pr,
		Stdout: &out,
		Stderr: io.Discard,
	}
	outcome, err := Run(opts)
	if err != nil {
		t.Fatal(err)
	}
	if outcome.GitStatus != "skipped (non-interactive mode)" {
		t.Errorf("GitStatus = %q, want skipped (non-interactive mode)", outcome.GitStatus)
	}
	if outcome.ProjectName != filepath.Base(dir) {
		t.Errorf("ProjectName = %q, want %q (piped answers must be ignored)", outcome.ProjectName, filepath.Base(dir))
	}
	if _, err := os.Stat(filepath.Join(dir, ".git")); !os.IsNotExist(err) {
		t.Errorf(".git must not exist after piped init, stat err = %v", err)
	}
	if out.Len() != 0 {
		t.Errorf("no prompts may be printed for a piped stdin, got:\n%s", out.String())
	}
	if outcome.Report == nil || !outcome.Report.Pass() {
		t.Fatalf("generated repo must validate, report: %+v", outcome.Report)
	}
}

// --- Skeleton integrity -------------------------------------------------

// TestSkeletonIntegrity verifies the generated tree matches the embedded
// skeleton file set: same paths, same bytes, README differing only in the
// substituted heading line.
func TestSkeletonIntegrity(t *testing.T) {
	dir := t.TempDir()
	opts, _, _ := runOpts(dir, "")
	if _, err := Run(opts); err != nil {
		t.Fatal(err)
	}
	got := walkFiles(t, dir)
	want := skeletonFiles(t)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("generated tree mismatch:\ngot:  %v\nwant: %v", got, want)
	}
	for _, rel := range want {
		data, err := os.ReadFile(filepath.Join(dir, filepath.FromSlash(rel)))
		if err != nil {
			t.Fatalf("missing generated file %s: %v", rel, err)
		}
		sub, err := fs.Sub(skeletonembed.FS, "skeleton")
		if err != nil {
			t.Fatal(err)
		}
		embedded, err := fs.ReadFile(sub, rel)
		if err != nil {
			t.Fatalf("reading skeleton %s: %v", rel, err)
		}
		if rel == "README.md" {
			// Only the heading line may differ.
			lines := strings.Split(string(data), "\n")
			if lines[0] == "" || !strings.HasPrefix(lines[0], "# ") {
				t.Errorf("README heading missing: %q", lines[0])
			}
			lines[0] = ""
			el := strings.Split(string(embedded), "\n")
			el[0] = ""
			if strings.Join(lines, "\n") != strings.Join(el, "\n") {
				t.Error("README must differ from the template only in the heading")
			}
		} else if !bytes.Equal(data, embedded) {
			t.Errorf("file %s differs from the embedded skeleton", rel)
		}
	}
}
