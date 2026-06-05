package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/IceflowRE/gitcc-cli/v3/gitcc/validators"
)

var errFailedToInitializeDB = errors.New("failed to initialize validator database")

type validatorsCmd struct {
	*cobra.Command

	ctx *validationContext
}

func newValidatorsCmd(ctx *validationContext) *validatorsCmd {
	cmd := &validatorsCmd{
		Command: &cobra.Command{
			Use:   "validator",
			Short: "Manage validators.",
		},
		ctx: ctx,
	}
	cmd.AddCommand(&cobra.Command{
		Use:   "ls",
		Short: "list available validators",
		RunE:  cmd.listValidators,
	})
	cmd.AddCommand(&cobra.Command{
		Use:   "rm name...",
		Short: "remove a validator",
		Args:  cobra.MinimumNArgs(1),
		RunE:  cmd.removeValidators,
	})
	cmd.AddCommand(&cobra.Command{
		Use:   "prune",
		Short: "remove older versions of validators",
		RunE:  cmd.pruneValidators,
	})
	cmd.AddCommand(&cobra.Command{
		Use:   "compile name path",
		Short: "compile a validator from a go file",
		Args:  cobra.ExactArgs(2), //nolint:mnd
		RunE:  cmd.compileValidator,
	})
	cmd.AddCommand(&cobra.Command{
		Use:   "dir",
		Short: "print validator directory",
		RunE:  cmd.printValidatorDir,
	})

	return cmd
}

func (cmd *validatorsCmd) listValidators(*cobra.Command, []string) (err error) {
	err = cmd.initDB()
	if err != nil {
		return fmt.Errorf("%w: %w", errFailedToInitializeDB, err)
	}

	for _, name := range cmd.ctx.db.AvailableNames() {
		fmt.Fprintln(cmd.OutOrStdout(), name)
	}

	return nil
}

func (cmd *validatorsCmd) removeValidators(_ *cobra.Command, args []string) (err error) {
	err = cmd.initDB()
	if err != nil {
		return fmt.Errorf("%w: %w", errFailedToInitializeDB, err)
	}

	for _, name := range args {
		val := cmd.ctx.db.GetCustomByName(name)
		if val == "" {
			fmt.Fprintf(cmd.ErrOrStderr(), "validator %s does not exist\n", name)

			continue
		}
		err := os.Remove(val)
		if err != nil {
			return fmt.Errorf("failed to remove validator %s: %w", name, err)
		}
		fmt.Fprintln(cmd.OutOrStdout(), name)
	}

	return nil
}

func (cmd *validatorsCmd) pruneValidators(_ *cobra.Command, _ []string) (err error) {
	deleted, err := validators.PruneValidators()

	if len(deleted) > 0 {
		fmt.Fprint(cmd.OutOrStdout(), "removed\n")

		for _, name := range deleted {
			fmt.Fprintf(cmd.OutOrStdout(), "%s\n", name)
		}
	}

	if err != nil {
		return fmt.Errorf("failed to prune validators: %w", err)
	}

	return nil
}

func (cmd *validatorsCmd) compileValidator(_ *cobra.Command, args []string) (err error) {
	err = cmd.initDB()
	if err != nil {
		return fmt.Errorf("%w: %w", errFailedToInitializeDB, err)
	}

	name := args[0]
	path := args[1]

	_, err = cmd.ctx.db.CompileCustom(path, name, "")
	if err != nil {
		return fmt.Errorf("failed to compile validator: %w", err)
	}

	fmt.Fprintln(cmd.OutOrStdout(), name)

	return nil
}

func (cmd *validatorsCmd) printValidatorDir(_ *cobra.Command, _ []string) (err error) {
	path, err := validators.GetValidatorCacheDir()
	if err != nil {
		return fmt.Errorf("failed to get cache directory: %w", err)
	}
	fmt.Fprintln(cmd.OutOrStdout(), path)

	return nil
}

func (cmd *validatorsCmd) initDB() (err error) {
	if cmd.ctx.db == nil {
		cmd.ctx.db, err = validators.NewDB()
		if err != nil {
			return err
		}
	}

	return nil
}
