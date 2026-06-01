package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/IceflowRE/gitcc/v3/standalone/gitcc"
	"github.com/IceflowRE/gitcc/v3/standalone/gitcc/internal"
)

type historyCmd struct {
	*validateBaseCmd

	dir     string
	branch  string
	sha     string
	verbose bool
}

func newHistoryCmd(ctx *validationContext) *historyCmd {
	cmd := &historyCmd{
		validateBaseCmd: newValidateBaseCmd(&cobra.Command{
			Use:   "history",
			Args:  cobra.NoArgs,
			Short: "Validate the messages of all commits in the history of the current branch.",
		}, ctx),
	}
	cmd.RunE = cmd.runE
	cmd.Flags().StringVarP(&cmd.dir, "dir", "", "./",
		"Path to a git repository. If not specified, the current directory is used.")
	cmd.Flags().StringVarP(&cmd.branch, "branch", "", "",
		"Validate until the common ancestor of the current branch and the specified branch. If not specified, the entire history is validated.")
	cmd.Flags().StringVarP(&cmd.sha, "sha", "", "",
		"SHA of the commit to stop validating at. If not specified, the entire history is validated.")
	cmd.Flags().BoolVarP(&cmd.verbose, "verbose", "v", false,
		"Include also valid commits in output.")

	cmd.MarkFlagsMutuallyExclusive("branch", "sha")

	return cmd
}

func (cmd *historyCmd) runE(_ *cobra.Command, args []string) error {
	return cmd.validate(func(validator gitcc.Validator) error {
		repo, err := internal.LoadRepository(cmd.dir)
		if err != nil {
			return fmt.Errorf("failed to open repository: %w", err)
		}

		sha := cmd.sha
		if cmd.branch != "" {
			commit, err := internal.GetMergeBase(repo, cmd.branch)
			if err != nil {
				return fmt.Errorf("failed to get merge base: %w", err)
			}
			sha = commit.Hash.String()
		}

		results, err := internal.ValidateHistory(validator, repo, sha)
		if err != nil {
			return fmt.Errorf("validation failed: %w", err)
		}

		status := gitcc.Valid
		for _, res := range results {
			if cmd.verbose || res.Status != gitcc.Valid {
				fmt.Fprintln(cmd.OutOrStdout(), res.String())
			}
			if res.Status.Severity() > status.Severity() {
				status = res.Status
			}
		}

		cmd.SilenceErrors = true
		cmd.SilenceUsage = true

		return getExitErrorFromStatus(status)
	})
}
