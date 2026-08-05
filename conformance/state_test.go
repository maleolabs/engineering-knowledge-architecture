package conformance

import "testing"

func TestTypeTokenCount(t *testing.T) {
	if got := len(typeTokens); got != 26 {
		t.Fatalf("type token table has %d entries, want 26", got)
	}
}

func TestOwnedSets(t *testing.T) {
	// Expected owned sets per validation.md Rule 4 (type -> owned domains).
	want := map[string][]string{
		"vis":  {DomainContentState, DomainExistenceState},
		"str":  {DomainContentState, DomainExistenceState},
		"req":  {DomainContentState, DomainExistenceState},
		"scp":  {DomainContentState, DomainExistenceState},
		"epc":  {DomainContentState, DomainExistenceState},
		"plan": {DomainContentState, DomainPlanningState, DomainExistenceState},
		"ctr":  {DomainContainerState, DomainExistenceState},
		"tkt":  {},
		"sto":  {DomainExecutionState, DomainExistenceState},
		"ts":   {DomainExecutionState, DomainExistenceState},
		"bug":  {DomainExecutionState, DomainExistenceState},
		"td":   {DomainExecutionState, DomainExistenceState},
		"ch":   {DomainExecutionState, DomainExistenceState},
		"spk":  {DomainExecutionState, DomainExistenceState},
		"ses":  {DomainExistenceState},
		"rvw":  {DomainContentState, DomainExistenceState},
		"adr":  {DomainContentState, DomainExistenceState},
		"dec":  {DomainContentState, DomainExistenceState},
		"arc":  {DomainContentState, DomainExistenceState},
		"spec": {DomainContentState, DomainExistenceState},
		"std":  {DomainContentState, DomainExistenceState},
		"run":  {DomainContentState, DomainExistenceState},
		"rel":  {DomainContentState, DomainExistenceState},
		"gls":  {DomainContentState, DomainExistenceState},
		"trc":  {DomainContentState, DomainExistenceState},
		"fnd":  {DomainContentState, DomainExistenceState},
	}
	for tok, expected := range want {
		info, ok := typeTokens[tok]
		if !ok {
			t.Errorf("token %q missing from type table", tok)
			continue
		}
		if !sameStrings(info.Owned, expected) {
			t.Errorf("type %q owned set = %v, want %v", tok, info.Owned, expected)
		}
	}
}

func TestContentStateVariant(t *testing.T) {
	cases := []struct {
		typ  string
		want []string
	}{
		{"adr", []string{"proposed", "accepted", "superseded"}},
		{"dec", []string{"draft", "accepted", "superseded"}},
		{"vis", []string{"draft", "review", "approved", "amended"}},
		{"plan", []string{"draft", "review", "approved", "amended"}},
		{"sto", []string{"draft", "review", "approved", "amended"}},
	}
	for _, c := range cases {
		if got := contentStateVariant(c.typ); !sameStrings(got, c.want) {
			t.Errorf("contentStateVariant(%q) = %v, want %v", c.typ, got, c.want)
		}
	}
}

func TestDimensionTokens(t *testing.T) {
	want := []string{
		"intent", "requirements", "architecture", "decisions", "specifications",
		"standards", "operations", "quality", "planning", "records", "research",
		"vocabulary",
	}
	if len(dimensionTokens) != len(want) {
		t.Fatalf("dimension table has %d entries, want %d", len(dimensionTokens), len(want))
	}
	for _, d := range want {
		if !dimensionTokens[d] {
			t.Errorf("dimension %q missing", d)
		}
	}
}

func TestIsLegalTransition(t *testing.T) {
	cases := []struct {
		name   string
		domain string
		typ    string
		from   string
		to     string
		legal  bool
	}{
		// Execution State: strictly adjacent, forward only.
		{"exec adjacent forward", DomainExecutionState, "sto", "planned", "todo", true},
		{"exec full chain", DomainExecutionState, "sto", "in-review", "done", true},
		{"exec skip", DomainExecutionState, "sto", "planned", "in-progress", false},
		{"exec revert", DomainExecutionState, "sto", "in-progress", "todo", false},
		{"exec noop", DomainExecutionState, "sto", "todo", "todo", false},
		{"exec from done", DomainExecutionState, "sto", "done", "done", false},
		{"initial marker", DomainExecutionState, "sto", "-", "planned", true},
		// Other domains: forward-only, adjacency not required.
		{"living skip allowed", DomainContentState, "vis", "draft", "amended", true},
		{"living revert", DomainContentState, "vis", "approved", "draft", false},
		{"living adjacent", DomainContentState, "vis", "review", "approved", true},
		{"adr variant forward", DomainContentState, "adr", "proposed", "accepted", true},
		{"adr variant supersede", DomainContentState, "adr", "accepted", "superseded", true},
		{"adr variant revert", DomainContentState, "adr", "accepted", "proposed", false},
		{"decision variant", DomainContentState, "dec", "draft", "superseded", true},
		{"planning forward", DomainPlanningState, "plan", "draft", "immutable", true},
		{"planning revert", DomainPlanningState, "plan", "approved", "draft", false},
		{"container forward", DomainContainerState, "ctr", "active", "completed", true},
		{"container revert", DomainContainerState, "ctr", "completed", "active", false},
		{"existence skip", DomainExistenceState, "adr", "active", "retired", true},
		{"existence revert", DomainExistenceState, "adr", "archived", "active", false},
		// Phase: no transition ordering is defined.
		{"phase any move", DomainPhase, "plan", "discovery", "release", true},
		{"phase revert", DomainPhase, "plan", "maturity", "mvp", true},
	}
	for _, c := range cases {
		if got := isLegalTransition(c.domain, c.typ, c.from, c.to); got != c.legal {
			t.Errorf("%s: isLegalTransition(%s, %s, %q, %q) = %v, want %v",
				c.name, c.domain, c.typ, c.from, c.to, got, c.legal)
		}
	}
}

func TestPhaseValueSet(t *testing.T) {
	want := []string{"discovery", "mvp", "milestone", "release", "growth", "maturity", "sunset"}
	if !sameStrings(phaseValues, want) {
		t.Errorf("phaseValues = %v, want %v", phaseValues, want)
	}
}

func sameStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
