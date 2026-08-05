package conformance

// This file encodes the EKA v1.0 type/state taxonomy used by the rules:
// the 26 artifact type tokens, the 12 knowledge dimensions, the five owned
// state domains (+ phase context attribute), the value sets per domain, and
// the forward-only transition tables.
//
// Grounding:
//   - skeleton/docs/exchange/validation.md (Rules 3 and 4 tables)
//   - skeleton/docs/operating/protocol.md §2 (value sets, transition rules)
//   - reference/reference-architecture.md §2.1 (26-token table)
//   - standard/eka-specification-v1.0.md §7 (State Taxonomy) and §8 (Knowledge Taxonomy)

// State domain names, exactly as they appear in frontmatter.
const (
	DomainContentState   = "content-state"
	DomainExecutionState = "execution-state"
	DomainPlanningState  = "planning-state"
	DomainContainerState = "container-state"
	DomainExistenceState = "existence-state"
	DomainPhase          = "phase"
)

// stateFields lists the five owned state domains (phase is a context
// attribute, not a state domain).
var stateFields = []string{
	DomainContentState,
	DomainExecutionState,
	DomainPlanningState,
	DomainContainerState,
	DomainExistenceState,
}

// relationshipFields lists the relationship fields validated by Rule 5.
var relationshipFields = []string{
	"amends",
	"supersedes",
	"derives-from",
	"depends-on",
	"validates",
}

// versionedTypes are the artifact types that MUST carry a -v<nn> filename
// suffix (Rule 2).
var versionedTypes = map[string]bool{"scp": true, "plan": true}

// projectionTypes may never carry a dimension (Rule 6); ticket is also the
// empty-state-vector type (Rule 4).
var projectionTypes = map[string]bool{"ctr": true, "tkt": true, "ses": true}

// workItemTypes are the operating-layer work items whose dimension is
// informational (Rule 6) and which own the Execution State domain.
var workItemTypes = map[string]bool{"sto": true, "ts": true, "bug": true, "td": true, "ch": true, "spk": true}

// TypeInfo describes one of the 26 artifact types.
type TypeInfo struct {
	// Token is the frontmatter `type` value, without the trailing dash
	// used in filenames (e.g. "adr").
	Token string
	// Owned lists the state domains this type owns (Rule 4). Absence of a
	// domain here means the field must not appear on the artifact.
	Owned []string
	// IsKnowledge reports whether the type is a knowledge artifact whose
	// `dimension` must equal its home folder (Rule 6). Work items are
	// exempt (dimension is informational).
	IsKnowledge bool
}

// typeTokens is the canonical 26-token table (reference-architecture.md §2.1,
// validation.md Rule 4). The owned sets follow validation.md Rule 4 exactly.
var typeTokens = map[string]TypeInfo{
	"vis":  {Token: "vis", Owned: []string{DomainContentState, DomainExistenceState}, IsKnowledge: true},
	"str":  {Token: "str", Owned: []string{DomainContentState, DomainExistenceState}, IsKnowledge: true},
	"req":  {Token: "req", Owned: []string{DomainContentState, DomainExistenceState}, IsKnowledge: true},
	"scp":  {Token: "scp", Owned: []string{DomainContentState, DomainExistenceState}, IsKnowledge: true},
	"epc":  {Token: "epc", Owned: []string{DomainContentState, DomainExistenceState}, IsKnowledge: true},
	"plan": {Token: "plan", Owned: []string{DomainContentState, DomainPlanningState, DomainExistenceState}, IsKnowledge: true},
	"ctr":  {Token: "ctr", Owned: []string{DomainContainerState, DomainExistenceState}, IsKnowledge: false},
	"tkt":  {Token: "tkt", Owned: []string{}, IsKnowledge: false},
	"sto":  {Token: "sto", Owned: []string{DomainExecutionState, DomainExistenceState}, IsKnowledge: false},
	"ts":   {Token: "ts", Owned: []string{DomainExecutionState, DomainExistenceState}, IsKnowledge: false},
	"bug":  {Token: "bug", Owned: []string{DomainExecutionState, DomainExistenceState}, IsKnowledge: false},
	"td":   {Token: "td", Owned: []string{DomainExecutionState, DomainExistenceState}, IsKnowledge: false},
	"ch":   {Token: "ch", Owned: []string{DomainExecutionState, DomainExistenceState}, IsKnowledge: false},
	"spk":  {Token: "spk", Owned: []string{DomainExecutionState, DomainExistenceState}, IsKnowledge: false},
	"ses":  {Token: "ses", Owned: []string{DomainExistenceState}, IsKnowledge: false},
	"rvw":  {Token: "rvw", Owned: []string{DomainContentState, DomainExistenceState}, IsKnowledge: true},
	"adr":  {Token: "adr", Owned: []string{DomainContentState, DomainExistenceState}, IsKnowledge: true},
	"dec":  {Token: "dec", Owned: []string{DomainContentState, DomainExistenceState}, IsKnowledge: true},
	"arc":  {Token: "arc", Owned: []string{DomainContentState, DomainExistenceState}, IsKnowledge: true},
	"spec": {Token: "spec", Owned: []string{DomainContentState, DomainExistenceState}, IsKnowledge: true},
	"std":  {Token: "std", Owned: []string{DomainContentState, DomainExistenceState}, IsKnowledge: true},
	"run":  {Token: "run", Owned: []string{DomainContentState, DomainExistenceState}, IsKnowledge: true},
	"rel":  {Token: "rel", Owned: []string{DomainContentState, DomainExistenceState}, IsKnowledge: true},
	"gls":  {Token: "gls", Owned: []string{DomainContentState, DomainExistenceState}, IsKnowledge: true},
	"trc":  {Token: "trc", Owned: []string{DomainContentState, DomainExistenceState}, IsKnowledge: true},
	"fnd":  {Token: "fnd", Owned: []string{DomainContentState, DomainExistenceState}, IsKnowledge: true},
}

// dimensionTokens is the 12 knowledge dimension vocabulary (EKA §8, ADR-005).
var dimensionTokens = map[string]bool{
	"intent":         true,
	"requirements":   true,
	"architecture":   true,
	"decisions":      true,
	"specifications": true,
	"standards":      true,
	"operations":     true,
	"quality":        true,
	"planning":       true,
	"records":        true,
	"research":       true,
	"vocabulary":     true,
}

// Value sets per domain (validation.md Rule 3, protocol.md §2).

var executionStateValues = []string{"planned", "todo", "in-progress", "in-review", "done"}
var planningStateValues = []string{"draft", "approved", "immutable"}
var containerStateValues = []string{"active", "completed"}
var existenceStateValues = []string{"active", "archived", "retired"}
var phaseValues = []string{"discovery", "mvp", "milestone", "release", "growth", "maturity", "sunset"}

// contentStateVariants is the content-state value set per artifact family
// (validation.md Rule 3, EKA §7.2). "living" is the standard variant.
var contentStateVariants = map[string][]string{
	"living":   {"draft", "review", "approved", "amended"},
	"adr":      {"proposed", "accepted", "superseded"},
	"decision": {"draft", "accepted", "superseded"},
}

// contentStateVariant returns the content-state value set for an artifact
// type: adr uses the ADR variant, dec the decision variant, every other type
// the living variant.
func contentStateVariant(typeToken string) []string {
	switch typeToken {
	case "adr":
		return contentStateVariants["adr"]
	case "dec":
		return contentStateVariants["decision"]
	default:
		return contentStateVariants["living"]
	}
}

// domainValues returns the ordered value set for a state domain, selecting
// the content-state variant for the given artifact type.
func domainValues(domain, typeToken string) []string {
	switch domain {
	case DomainContentState:
		return contentStateVariant(typeToken)
	case DomainExecutionState:
		return executionStateValues
	case DomainPlanningState:
		return planningStateValues
	case DomainContainerState:
		return containerStateValues
	case DomainExistenceState:
		return existenceStateValues
	case DomainPhase:
		return phaseValues
	default:
		return nil
	}
}

// isLegalTransition reports whether the from -> to transition is legal for
// the domain on the given artifact type.
//
// Interpretation (documented): only Execution State is "strictly sequential"
// (EKA §7.2, protocol.md §2: never skip, never revert), so adjacency is
// required there. All other state domains are forward-only: the position
// index must strictly increase (never revert), but skipping is tolerated
// because the spec's forward-only rule does not demand adjacency for them.
// `from: "-"` is the initial-state marker established by the repository's
// own ADRs and is always a legal start. Phase is a context attribute, not a
// state domain: any transition between valid phase values is legal (no
// ordering constraint, EKA 11.2).
func isLegalTransition(domain, typeToken, from, to string) bool {
	if from == "-" || domain == DomainPhase {
		return true
	}
	values := domainValues(domain, typeToken)
	fromIdx := indexOf(values, from)
	toIdx := indexOf(values, to)
	if fromIdx < 0 || toIdx < 0 {
		return false
	}
	if domain == DomainExecutionState {
		return toIdx == fromIdx+1
	}
	return toIdx > fromIdx
}

// indexOf returns the position of v in values, or -1.
func indexOf(values []string, v string) int {
	for i, x := range values {
		if x == v {
			return i
		}
	}
	return -1
}

// contains reports whether values contains v.
func contains(values []string, v string) bool {
	return indexOf(values, v) >= 0
}
