package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/IceflowRE/gitcc/v3/standalone/gitcc"
	"github.com/IceflowRE/gitcc/v3/standalone/gitcc/internal"
)

type commitCmd struct {
	*validateBaseCmd

	dir string
}

func newCommitCmd(ctx *validationContext) *commitCmd {
	cmd := &commitCmd{
		validateBaseCmd: newValidateBaseCmd(&cobra.Command{
			Use:   "commit [sha256]",
			Args:  cobra.MaximumNArgs(1),
			Short: "Validate the message of a commit. By default, the HEAD commit is validated.",
		}, ctx),
	}
	cmd.Flags().StringVarP(&cmd.dir, "dir", "", "./",
		"Path to a git repository. If not specified, the current directory is used.")
	cmd.RunE = cmd.runE

	return cmd
}

func (cmd *commitCmd) runE(_ *cobra.Command, args []string) error {
	return cmd.validate(func(validator gitcc.Validator) error {
		repo, err := internal.LoadRepository(cmd.dir)
		if err != nil {
			return fmt.Errorf("failed to open repository: %w", err)
		}

		var res gitcc.Result
		if len(args) == 0 {
			res, err = internal.ValidateHead(validator, repo)
		} else {
			res, err = internal.ValidateCommit(validator, repo, args[0])
		}

		if err != nil {
			return fmt.Errorf("validation failed: %w", err)
		}

		fmt.Fprintln(cmd.OutOrStdout(), res.String())

		cmd.SilenceErrors = true
		cmd.SilenceUsage = true

		return getExitErrorFromStatus(res.Status)
	})
}
