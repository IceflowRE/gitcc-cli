package cli

import (
	"github.com/spf13/cobra"

	"github.com/IceflowRE/gitcc/v3/standalone/gitcc"
)

type Cli struct {
	rootCmd *cobra.Command
}

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

func (cli *Cli) Execute() error {
	return cli.rootCmd.Execute()
}
