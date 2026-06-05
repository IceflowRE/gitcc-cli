// Package cli provides the commandline interface for gitcc.
package cli

import (
	"github.com/spf13/cobra"

	"github.com/IceflowRE/gitcc-cli/v3/gitcc"
)

// Cli is the commandline interface entrypoint.
type Cli struct {
	rootCmd *cobra.Command
}

// NewCli creates a new Cli instance.
func NewCli(validator gitcc.Validator) *Cli {
	cli := &Cli{
		rootCmd: &cobra.Command{},
	}
	cli.rootCmd.CompletionOptions.DisableDefaultCmd = true

	ctx := newValidationContext(validator)
	cli.rootCmd.AddCommand(
		newCommitCmd(ctx).Command,
		newHistoryCmd(ctx).Command,
		newMessageCmd(ctx).Command,
		newValidatorsCmd(ctx).Command,
	)

	return cli
}

// Execute runs the main application.
func (cli *Cli) Execute() error {
	return cli.rootCmd.Execute()
}
