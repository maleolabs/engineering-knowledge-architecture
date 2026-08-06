package view

import (
	"reflect"
	"testing"
)

func TestArchitectureProjection(t *testing.T) {
	g := loadFixture(t, "valid")
	p, err := Build("architecture", g, "")
	if err != nil {
		t.Fatalf("Build(architecture): %v", err)
	}
	arch, ok := p.(*ArchitectureProjection)
	if !ok {
		t.Fatalf("Build(architecture) = %T, want *ArchitectureProjection", p)
	}

	// Fixed group order.
	wantOrder := []string{"Decisions", "Architecture Descriptions", "Specifications", "Standards & Guidelines", "Vocabulary"}
	gotOrder := make([]string, len(arch.Groups))
	for i, gr := range arch.Groups {
		gotOrder[i] = gr.Name
	}
	if !reflect.DeepEqual(gotOrder, wantOrder) {
		t.Errorf("group order = %v, want %v", gotOrder, wantOrder)
	}

	// Decisions merges adr- and dec- in canonical identity order
	// (adr:001 < adr:002 < adr:003 < dec:001), with the ADR and decision
	// content-state variants.
	wantGroups := map[string][][4]string{
		"Decisions": {
			{validForm + "adr:001-login-serialization", "accepted", "", ""},
			{validForm + "adr:002-session-encoding", "superseded", "", ""},
			{validForm + "adr:003-token-format", "accepted", "", ""},
			{validForm + "dec:001-api-shape", "accepted", "", ""},
		},
		"Architecture Descriptions": {
			{validForm + "arc:system-architecture", "approved", "", ""},
		},
		"Specifications": {
			{validForm + "spec:auth-flow", "draft", "", ""},
		},
		"Standards & Guidelines": {
			{validForm + "std:gofmt", "review", "", ""},
		},
		"Vocabulary": {
			{validForm + "gls:domain-terms", "amended", "", ""},
		},
	}
	for _, gr := range arch.Groups {
		want := wantGroups[gr.Name]
		if got := groupArtifactStates(gr); !reflect.DeepEqual(got, want) {
			t.Errorf("group %q = %+v, want %+v", gr.Name, got, want)
		}
	}
}

// TestArchitectureProjectionEmptyDomain: no Architecture artifacts —
// empty projection, fixed group shape preserved.
func TestArchitectureProjectionEmptyDomain(t *testing.T) {
	g := loadFixture(t, "no-active")
	p, err := Build("architecture", g, "")
	if err != nil {
		t.Fatalf("Build(architecture): %v", err)
	}
	arch := p.(*ArchitectureProjection)
	if len(arch.Groups) != 5 {
		t.Fatalf("groups = %d, want the fixed five", len(arch.Groups))
	}
	for _, gr := range arch.Groups {
		if len(gr.Artifacts) != 0 {
			t.Errorf("group %q must be empty, got %v", gr.Name, gr.Artifacts)
		}
	}
}
