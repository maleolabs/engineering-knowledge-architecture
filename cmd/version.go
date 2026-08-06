package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// version is the CLI build version. It defaults to "dev" and is
// overridable at build time:
//
//	go build -ldflags "-X github.com/maleolabs/engineering-knowledge-architecture/cmd.version=v1.2.3" ./cmd/eka
//
// The version identifies the CLI implementation, never the standard: the
// EKA standard version is standardVersion and is fixed by the ratified
// specifications.
var version = "dev"

// standardVersion is the EKA standard corpus version this CLI implements
// (Core + Exchange + Naming and Terminology v1.0).
const standardVersion = "1.0"

// newVersionCommand builds the `eka version` command: prints the CLI
// build version and the EKA standard version. Deterministic output.
func newVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Long: `Print the CLI build version and the EKA standard version this
CLI implements.

The CLI version is set at build time (default "dev"):
  go build -ldflags "-X .../cmd.version=v1.2.3" ./cmd/eka

The standard version is fixed by the ratified specifications (EKA v1.0).`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			s := styleFor(cmd)
			fmt.Fprintf(s.W, "eka %s\n", version)
			fmt.Fprintf(s.W, "EKA standard %s\n", standardVersion)
			return nil
		},
	}
}
