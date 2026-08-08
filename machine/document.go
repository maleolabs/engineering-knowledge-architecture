// Package machine implements the machine interface of the EKA Runtime:
// it serializes Canonical Knowledge Objects (CKO) to the deterministic
// canonical JSON consumed by `eka get` and future machine consumers
// (MCP, Atrium, VS Code, AI agents, scripts).
//
// The machine interface is the machine-readable counterpart of the
// projection engine (view/): it never renders for readability, never
// parses Markdown, never touches storage, and never reuses projection
// renderers — pure CKO in, canonical JSON out.
//
// Determinism contract: fixed struct field order (declaration order),
// a stable schema string ("eka-cko-v1", stable across minor releases),
// and sorted inputs by canonical form (collections).
package machine

import (
	"encoding/json"
	"fmt"

	"github.com/maleolabs/engineering-knowledge-architecture/conformance"
	"github.com/maleolabs/engineering-knowledge-architecture/exchange"
)

// Schema is the canonical JSON schema identifier of the machine
// interface: "eka-cko-v1". Stable across minor releases.
const Schema = "eka-cko-v1"

// Document is the canonical machine projection of one Canonical
// Knowledge Object (one exchange.Unit). Field order is the fixed
// serialization order of the schema.
type Document struct {
	Schema            string           `json:"schema"`
	Identity          Identity         `json:"identity"`
	CanonicalForm     string           `json:"canonical_form"`
	EngineeringDomain string           `json:"engineering_domain"`
	Stratum           int              `json:"stratum"`
	Revision          int              `json:"revision,omitempty"`
	Author            string           `json:"author,omitempty"`
	Created           string           `json:"created,omitempty"`
	Updated           string           `json:"updated,omitempty"`
	StateVector       StateVector      `json:"state_vector"`
	Phase             string           `json:"phase,omitempty"`
	Classification    *Classification  `json:"classification,omitempty"`
	Relationships     []Relationship   `json:"relationships,omitempty"`
	ChangeLog         []ChangeLogEntry `json:"change_log,omitempty"`
	Content           Content          `json:"content"`
	ObjectHash        string           `json:"object_hash"`
}

// Identity is the complete identity tuple of the CKO, in the fixed
// declared order (identical to the RSF unit.json naming).
type Identity struct {
	Namespace       string `json:"namespace"`
	Type            string `json:"type"`
	ID              string `json:"id"`
	InstanceVersion int    `json:"instance_version"`
}

// StateVector carries the five owned state domains with the canonical
// RSF unit.json naming. Values are never empty strings in a conformant
// repository, so omitempty loses nothing.
type StateVector struct {
	ContentState   string `json:"content-state,omitempty"`
	ExecutionState string `json:"execution-state,omitempty"`
	PlanningState  string `json:"planning-state,omitempty"`
	ContainerState string `json:"container-state,omitempty"`
	ExistenceState string `json:"existence-state,omitempty"`
}

// Relationship is one recorded relationship by Identity (stored order
// preserved).
type Relationship struct {
	Type   string `json:"type"`
	Target string `json:"target"`
}

// ChangeLogEntry is one recorded transition in occurrence order.
type ChangeLogEntry struct {
	Date   string `json:"date"`
	Domain string `json:"domain"`
	From   string `json:"from"`
	To     string `json:"to"`
	By     string `json:"by"`
}

// Classification carries the primary Knowledge Dimension, at most one
// secondary, and the Engineering Domain — omitted entirely (nil) when
// the CKO carries none.
type Classification struct {
	Dimension           string   `json:"dimension,omitempty"`
	DimensionsSecondary []string `json:"dimensions_secondary,omitempty"`
	Domain              string   `json:"domain,omitempty"`
}

// Content is the representation-tagged knowledge payload: the opaque
// representation payload of the CKO, never parsed or re-structured.
// Text may be large — that is the knowledge content; it is never
// truncated.
type Content struct {
	Representation string `json:"representation"`
	Text           string `json:"text"`
}

// NewDocument maps one Canonical Knowledge Object (exchange.Unit) to
// its machine Document:
//
//   - identity tuple and canonical form pass through;
//   - engineering_domain is Classification.Domain when non-empty, else
//     derived from the artifact type token via conformance.DomainForToken
//     (the single shared source of truth) — an unknown token is a
//     deterministic error;
//   - stratum is the authority stratum of the engineering domain
//     (conformance.Stratum);
//   - state vector, classification, relationships (stored order),
//     change log (occurrence order) pass through;
//   - content is {representation, text} where text is the opaque
//     representation payload — never parsed, never truncated;
//   - object_hash is the CKO digest ("" for hand-built units, kept
//     as-is).
func NewDocument(u *exchange.Unit) (*Document, error) {
	if u == nil {
		return nil, fmt.Errorf("machine: cannot build a document from a nil unit")
	}
	domain := u.Classification.Domain
	if domain == "" {
		d, ok := conformance.DomainForToken(u.Identity.Type)
		if !ok {
			return nil, fmt.Errorf("machine: unknown artifact type %q", u.Identity.Type)
		}
		domain = string(d)
	}
	doc := &Document{
		Schema:            Schema,
		Identity:          Identity{Namespace: u.Identity.Namespace, Type: u.Identity.Type, ID: u.Identity.ID, InstanceVersion: u.Identity.InstanceVersion},
		CanonicalForm:     u.CanonicalIdentityForm,
		EngineeringDomain: domain,
		Stratum:           conformance.Stratum(conformance.Domain(domain)),
		Revision:          u.Revision,
		Author:            u.Author,
		Created:           u.Created,
		Updated:           u.Updated,
		StateVector: StateVector{
			ContentState:   u.StateVector.ContentState,
			ExecutionState: u.StateVector.ExecutionState,
			PlanningState:  u.StateVector.PlanningState,
			ContainerState: u.StateVector.ContainerState,
			ExistenceState: u.StateVector.ExistenceState,
		},
		Phase:      u.Phase,
		Content:    Content{Representation: u.Content.Representation, Text: string(u.ContentPayload)},
		ObjectHash: u.Digest,
	}
	if u.Classification.Dimension != "" || len(u.Classification.DimensionsSecondary) > 0 || u.Classification.Domain != "" {
		c := Classification{
			Dimension:           u.Classification.Dimension,
			DimensionsSecondary: u.Classification.DimensionsSecondary,
			Domain:              u.Classification.Domain,
		}
		doc.Classification = &c
	}
	for _, r := range u.Relationships {
		doc.Relationships = append(doc.Relationships, Relationship{Type: r.Type, Target: r.Target})
	}
	for _, e := range u.ChangeLog {
		doc.ChangeLog = append(doc.ChangeLog, ChangeLogEntry{Date: e.Date, Domain: e.Domain, From: e.From, To: e.To, By: e.By})
	}
	return doc, nil
}

// Marshal serializes the Document deterministically: json.MarshalIndent
// with two-space indentation (fixed struct field order = declaration
// order) plus a single trailing newline.
func (d *Document) Marshal() ([]byte, error) {
	out, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		return nil, err
	}
	return append(out, '\n'), nil
}

// MarshalUnit is the convenience entry point: NewDocument + Marshal in
// one call. It is the function the machine consumers (eka get, MCP,
// Atrium, ...) call per CKO.
func MarshalUnit(u *exchange.Unit) ([]byte, error) {
	doc, err := NewDocument(u)
	if err != nil {
		return nil, err
	}
	return doc.Marshal()
}
