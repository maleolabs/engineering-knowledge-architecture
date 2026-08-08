package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/maleolabs/engineering-knowledge-architecture/machine"
	"github.com/maleolabs/engineering-knowledge-architecture/runtime"
	"github.com/spf13/cobra"
)

// newGetCommand builds the `eka get` command: the machine interface of
// the EKA Runtime. It retrieves Canonical Knowledge Objects of the
// project owning the repository rooted at the current directory and
// emits them as the deterministic canonical JSON of the machine
// package (schema "eka-cko-v1"). No rendering, no Markdown, no
// banners: stdout carries ONLY the JSON document followed by a single
// trailing newline — machine consumers parse stdout verbatim.
//
// Query model (knowledge-shaped):
//
//	target containing ":"  identity lookup — the RSF canonical form
//	                       ("<ns>/<type>:<id>:<v>", exact instance) or
//	                       the qualified line form
//	                       ("<ns>/<type>:<id>", lowest instance). The
//	                       namespace is required (the Runtime resolves
//	                       globally; unqualified forms are refused).
//	target without ":"     domain query — one of the five Engineering
//	                       Domain tokens (discovery|architecture|
//	                       planning|execution|operations): the
//	                       "domain" collection of every matching unit.
//
// Exit codes:
//
//	0  JSON document produced
//	1  workspace/repository-state refusal (no workspace, repository
//	   not registered)
//	2  usage or internal error (invalid target or domain, unknown
//	   identity, resolver/store failure)
func newGetCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "get <target>",
		Short: "Retrieve knowledge as machine-readable CKO JSON",
		Long: `Retrieve Engineering Knowledge as machine-readable Canonical
Knowledge Object (CKO) JSON — the machine interface of the EKA
Runtime. Where 'eka view' renders human projections for reading,
'eka get' emits the deterministic canonical JSON consumed by
scripts, MCP, Atrium, VS Code and AI agents. The JSON schema is
"eka-cko-v1" and stable across minor releases; output is
deterministic (fixed field order, units sorted by canonical form,
no timestamps, no host-dependent values).

The target is either a knowledge identity or an Engineering Domain:

  identity  <ns>/<type>:<id>[:<instance-version>]
            the RSF canonical form (the exact instance) or the
            qualified line form (the lowest instance-version of the
            line). The namespace is required: the Runtime resolves
            globally, so unqualified forms are refused.
  domain    one of the five Engineering Domain tokens:
            discovery | architecture | planning | execution |
            operations
            the response is a "domain" collection of every matching
            unit of the project owning the current repository,
            sorted by canonical form, carrying the canonical
            Engineering Domain name (e.g. "Execution").

The repository must be registered in the EKA workspace and synced
first ('eka sync'). Repository resolution is exact-path: run this
command inside the repository root.

Output contract: stdout carries ONLY the JSON document — one
document for an identity lookup, one collection for a domain query
— followed by a single trailing newline. No banners, no
informational lines: machine consumers parse stdout verbatim.
Errors go to stderr, one 'eka: ...' line per error (the bare
'eka get' usage summary is the exception and also goes to stderr).

Exit codes:
  0  JSON document produced
  1  workspace/repository-state refusal (no EKA workspace,
     repository not registered in the workspace)
  2  usage or internal error (invalid target or domain, unknown
     identity, resolver failure)

Boundaries (documented, not implemented in this milestone): the
query surface stays minimal by design — future flags may add
traversal (upstream/downstream), timeline queries and result
filters; the schema is designed to stay stable while these extend
the query surface.`,
		Example: `  eka get feather/adr:001-identity-serialization:1
  eka get feather/sto:publish-post
  eka get architecture
  eka get execution`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				// Machine commands never print banners: the no-argument
				// case is a usage error with the query-model summary on
				// stderr, exit 2.
				fmt.Fprintln(cmd.ErrOrStderr(), "eka: get: usage: eka get <target>")
				fmt.Fprintln(cmd.ErrOrStderr(), "eka: get:   identity  <ns>/<type>:<id>[:<instance-version>]  (canonical form or qualified line form)")
				fmt.Fprintln(cmd.ErrOrStderr(), "eka: get:   domain    discovery | architecture | planning | execution | operations")
				fmt.Fprintln(cmd.ErrOrStderr(), "eka: get: run 'eka get --help' for the full reference")
				return &exitError{code: exitUsage}
			}
			// The resolution prologue: open (never create) the Runtime,
			// then gate on workspace and repository state.
			r, err := runtime.Open()
			if err != nil {
				return err // Exit 2: workspace resolution.
			}
			defer r.Close()
			if !r.Exists() {
				// Workspace-state refusal: `eka get` never creates a
				// workspace — deterministic message, exit 1.
				fmt.Fprintf(cmd.ErrOrStderr(), "eka: get refused: no EKA workspace at %s; run 'eka sync' first\n", r.Path())
				return &exitError{code: exitFail}
			}
			abs, err := filepath.Abs(".")
			if err != nil {
				return fmt.Errorf("get failed: %w", err)
			}
			repo, found, err := r.Workspace.FindRepo(abs)
			if err != nil {
				return fmt.Errorf("get failed: %w", err) // Exit 2: registry failure.
			}
			if !found {
				// Repository-state refusal: deterministic message and
				// exit 1 — no JSON is produced.
				fmt.Fprintf(cmd.ErrOrStderr(), "eka: get refused: repository %s is not registered in the EKA workspace; run 'eka sync' (auto-registers) or 'eka project register' first\n", abs)
				return &exitError{code: exitFail}
			}
			target := args[0]
			var out []byte
			if strings.Contains(target, ":") {
				// Identity lookup: canonical form (exact instance) or
				// qualified line form (lowest instance) — the Resolver
				// contract. Unqualified forms are refused by the
				// Resolver with the expected forms listed.
				unit, ok, err := r.Resolver.Resolve(target)
				if err != nil {
					return fmt.Errorf("get: %w", err) // Exit 2: usage.
				}
				if !ok {
					return fmt.Errorf("get: no knowledge object matches %q", target) // Exit 2.
				}
				out, err = machine.MarshalUnit(unit)
				if err != nil {
					return fmt.Errorf("get failed: %w", err) // Exit 2: internal.
				}
			} else {
				// Domain query: one of the five Engineering Domain
				// tokens; the collection JSON carries the canonical
				// domain name.
				name, ok := domainTokenName(target)
				if !ok {
					return fmt.Errorf("get: unknown domain %q — available domains: %s", target, strings.Join(domainTokenList(), ", "))
				}
				units, err := r.Knowledge.Search(runtime.SearchQuery{ProjectID: repo.ProjectID, Domain: name})
				if err != nil {
					return fmt.Errorf("get failed: %w", err) // Exit 2: store failure.
				}
				col, err := machine.NewCollection(name, units)
				if err != nil {
					return fmt.Errorf("get failed: %w", err) // Exit 2: internal.
				}
				out, err = col.Marshal()
				if err != nil {
					return fmt.Errorf("get failed: %w", err) // Exit 2: internal.
				}
			}
			// Output contract: stdout carries ONLY the JSON document
			// plus its single trailing newline (Marshal emits it) —
			// written verbatim, never re-rendered.
			if _, err := cmd.OutOrStdout().Write(out); err != nil {
				return fmt.Errorf("get failed: %w", err)
			}
			return nil
		},
	}
}

// domainTokens are the five Engineering Domain query tokens in stratum
// order, mapped to the canonical Engineering Domain names (the values
// carried by Classification.Domain of stored units and the machine
// JSON). Deterministic — never derived from map iteration.
var domainTokens = []struct {
	token string
	name  string
}{
	{"discovery", "Discovery"},
	{"architecture", "Architecture"},
	{"planning", "Planning"},
	{"execution", "Execution"},
	{"operations", "Operations"},
}

// domainTokenName maps a query token to its canonical Engineering
// Domain name; the second return value is false for unknown tokens.
func domainTokenName(token string) (string, bool) {
	for _, d := range domainTokens {
		if d.token == token {
			return d.name, true
		}
	}
	return "", false
}

// domainTokenList renders the five query tokens as the deterministic
// "a | b | c" usage list of the domain usage error.
func domainTokenList() []string {
	out := make([]string, 0, len(domainTokens))
	for _, d := range domainTokens {
		out = append(out, d.token)
	}
	return out
}
