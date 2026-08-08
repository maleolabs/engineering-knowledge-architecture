package machine

import (
	"encoding/json"
	"sort"

	"github.com/maleolabs/engineering-knowledge-architecture/exchange"
)

// Collection is the machine projection of a domain query: every unit of
// one Engineering Domain of the project, sorted by canonical form.
type Collection struct {
	Schema     string      `json:"schema"`
	Collection string      `json:"collection"` // "domain"
	Domain     string      `json:"domain"`
	Count      int         `json:"count"`
	Units      []*Document `json:"units"` // sorted by canonical form
}

// NewCollection maps the units of one domain query to a machine
// Collection. The units are sorted by canonical form regardless of the
// input order (determinism contract: sorted inputs by canonical form);
// an empty result is an empty unit list, never null. Domain carries the
// canonical Engineering Domain name (e.g. "Execution").
func NewCollection(domain string, units []*exchange.Unit) (*Collection, error) {
	docs := make([]*Document, 0, len(units))
	for _, u := range units {
		d, err := NewDocument(u)
		if err != nil {
			return nil, err
		}
		docs = append(docs, d)
	}
	// Canonical-form order is the collection contract (deterministic;
	// lexicographic on the canonical identity form — instance versions
	// are compared textually, matching the RSF canonical key ordering).
	sort.Slice(docs, func(i, j int) bool {
		return docs[i].CanonicalForm < docs[j].CanonicalForm
	})
	return &Collection{
		Schema:     Schema,
		Collection: "domain",
		Domain:     domain,
		Count:      len(docs),
		Units:      docs,
	}, nil
}

// Marshal serializes the Collection deterministically: the same
// formatting as Document.Marshal (two-space indent, trailing newline).
func (c *Collection) Marshal() ([]byte, error) {
	out, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return nil, err
	}
	return append(out, '\n'), nil
}
