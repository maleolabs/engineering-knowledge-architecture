package view

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

// TestProjectionsListsRegisteredNames verifies the registry surface:
// sorted, deterministic, closed over the three projections.
func TestProjectionsListsRegisteredNames(t *testing.T) {
	want := []string{"sprint", "ticket", "wave"}
	if got := Projections(); !reflect.DeepEqual(got, want) {
		t.Errorf("Projections() = %v, want %v", got, want)
	}
	for _, name := range want {
		if !IsProjection(name) {
			t.Errorf("IsProjection(%q) = false", name)
		}
	}
	if IsProjection("bogus") {
		t.Error("IsProjection(bogus) must be false")
	}
}

// TestUnknownProjection verifies Build's error contract for unregistered
// names.
func TestUnknownProjection(t *testing.T) {
	g := loadFixture(t, "valid")
	_, err := Build("bogus", g, "")
	if !errors.Is(err, ErrUnknownProjection) {
		t.Errorf("Build(bogus) error = %v, want ErrUnknownProjection", err)
	}
	if err == nil || !strings.Contains(err.Error(), "bogus") {
		t.Errorf("Build(bogus) error must name the projection, got %v", err)
	}
}

// TestTargetNotFound verifies the helpful TargetNotFoundError: the
// unresolved target plus the sorted available ticket list.
func TestTargetNotFound(t *testing.T) {
	g := loadFixture(t, "valid")
	_, err := Build("ticket", g, "tkt-ghost")
	var tne *TargetNotFoundError
	if !errors.As(err, &tne) {
		t.Fatalf("Build(ticket, tkt-ghost) error = %v, want *TargetNotFoundError", err)
	}
	if tne.Target != "tkt-ghost" {
		t.Errorf("TargetNotFoundError.Target = %q", tne.Target)
	}
	msg := tne.Error()
	if !strings.Contains(msg, "tkt-ghost") || !strings.Contains(msg, "available tickets") {
		t.Errorf("error must be helpful, got %q", msg)
	}
	want := []string{"bug-delta", "ch-epsilon", "sto-alpha", "sto-alpha-dup", "sto-beta", "sto-beta-multi", "sto-legacy", "ts-gamma", "unresolved"}
	if !reflect.DeepEqual(tne.Available, want) {
		t.Errorf("Available = %v, want %v", tne.Available, want)
	}

	// No tickets at all: distinct, still helpful message.
	empty := NewGraph(".", nil)
	_, err = Build("ticket", empty, "tkt-ghost")
	if !errors.As(err, &tne) {
		t.Fatalf("Build on empty graph error = %v, want *TargetNotFoundError", err)
	}
	if !strings.Contains(tne.Error(), "contains no tickets") {
		t.Errorf("empty-graph error must say so, got %q", tne.Error())
	}
}

// TestTargetRequired verifies the ticket projection refuses a missing
// target.
func TestTargetRequired(t *testing.T) {
	g := loadFixture(t, "valid")
	if _, err := Build("ticket", g, ""); err == nil {
		t.Error("Build(ticket, \"\") must fail")
	} else if !strings.Contains(err.Error(), "requires a target") {
		t.Errorf("error must explain the requirement, got %q", err.Error())
	}
}

// TestSprintIgnoresTarget verifies sprint/wave ignore the optional
// target argument.
func TestSprintIgnoresTarget(t *testing.T) {
	g := loadFixture(t, "valid")
	p, err := Build("sprint", g, "tkt-ghost")
	if err != nil {
		t.Fatalf("sprint with a target must not fail: %v", err)
	}
	if p.Name() != "sprint" {
		t.Errorf("Name() = %q, want sprint", p.Name())
	}
	if _, err := Build("wave", g, "anything"); err != nil {
		t.Errorf("wave with a target must not fail: %v", err)
	}
}
