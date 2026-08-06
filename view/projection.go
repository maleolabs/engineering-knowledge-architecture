package view

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

// This file implements the projection registry: the closed set of
// named projections the engine can build. Each projection is an
// independent builder function; a future projection is added by
// registering a builder here — there is no pipeline to redesign.

// Projection is one view over the Knowledge Graph. The model types are
// plain data with deterministic ordering; rendering is the caller's
// concern.
type Projection interface {
	// Name returns the projection's registry name ("sprint", "wave",
	// "ticket").
	Name() string
}

// ErrUnknownProjection is returned by Build for a projection name that
// is not registered.
var ErrUnknownProjection = errors.New("unknown projection")

// TargetNotFoundError reports an unresolvable ticket target with the
// available targets, so the CLI can fail with a helpful message.
type TargetNotFoundError struct {
	// Projection is the projection name the target belongs to.
	Projection string
	// Target is the unresolved target as given.
	Target string
	// Available lists the resolvable targets, sorted by canonical
	// identity.
	Available []string
}

// Error renders the deterministic diagnostic.
func (e *TargetNotFoundError) Error() string {
	if len(e.Available) == 0 {
		return fmt.Sprintf("%s projection: target %q not found (the repository contains no tickets)",
			e.Projection, e.Target)
	}
	return fmt.Sprintf("%s projection: target %q not found — available tickets: %s",
		e.Projection, e.Target, strings.Join(e.Available, ", "))
}

// builder constructs one projection over the graph. target is the
// optional user-supplied target; projections that do not use it ignore
// it.
type builder func(g *Graph, target string) (Projection, error)

// registry maps projection names to their independent builders.
var registry = map[string]builder{
	"sprint": buildSprint,
	"wave":   buildWave,
	"ticket": buildTicket,
}

// Projections returns the registered projection names, sorted —
// deterministic regardless of map iteration order.
func Projections() []string {
	names := make([]string, 0, len(registry))
	for name := range registry {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// IsProjection reports whether name is a registered projection.
func IsProjection(name string) bool {
	_, ok := registry[name]
	return ok
}

// Build constructs the named projection over the graph. It returns
// ErrUnknownProjection (wrapped) for unregistered names and
// *TargetNotFoundError when the ticket target does not resolve. A
// missing target on a target-requiring projection is an error.
func Build(name string, g *Graph, target string) (Projection, error) {
	b, ok := registry[name]
	if !ok {
		return nil, fmt.Errorf("%w: %q", ErrUnknownProjection, name)
	}
	return b(g, target)
}
