package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/go-git/go-git/v6/plumbing/object"
	"github.com/spf13/cobra"

	"github.com/IceflowRE/gitcc/v3/standalone/gitcc"
)

type messageCmd struct {
	*validateBaseCmd

	file string
}

func newMessageCmd(ctx *validationContext) (cmd *messageCmd) {
	cmd = &messageCmd{
		validateBaseCmd: newValidateBaseCmd(&cobra.Command{
			Use:   "message message",
			Args:  cobra.MaximumNArgs(1),
			Short: "Validate a given message.",
		}, ctx),
	}
	cmd.Flags().StringVarP(&cmd.file, "file", "", "",
		"Path to a text file to validate.")
	cmd.RunE = cmd.runE
	cmd.PreRunE = func(_ *cobra.Command, args []string) error {
		if len(args) != 0 && cmd.file != "" {
			return errors.New("cannot specify both a message argument and --file flag")
		}
		if len(args) == 0 && cmd.file == "" {
			return errors.New("must specify either a message argument or --file flag")
		}

		return nil
	}

	return cmd
}

func (cmd *messageCmd) runE(_ *cobra.Command, args []string) error {
	return cmd.validate(func(validator gitcc.Validator) error {
		if cmd.file != "" {
			content, err := os.ReadFile(cmd.file)
			if err != nil {
				return fmt.Errorf("failed to read file: %w", err)
			}
			args = []string{string(content)}
		}

		commit := &object.Commit{
			Message: args[0],
		}

		res := validator.Validate(commit)
		fmt.Fprintln(cmd.OutOrStdout(), res.String())

		cmd.SilenceErrors = true
		cmd.SilenceUsage = true

		return getExitErrorFromStatus(res.Status)
	})
}
