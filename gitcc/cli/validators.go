package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/IceflowRE/gitcc/v3/standalone/gitcc/validators"
)

type validatorsCmd struct {
	*cobra.Command

	ctx *validationContext
}

func newValidatorsCmd(ctx *validationContext) *validatorsCmd {
	cmd := &validatorsCmd{
		Command: &cobra.Command{
			Use:   "validator",
			Short: "Manage validator.",
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
		Args:  cobra.ExactArgs(2),
		RunE:  cmd.compileValidator,
	})
	cmd.AddCommand(&cobra.Command{
		Use:   "cache",
		Short: "print cache directory",
		RunE:  cmd.printCacheDir,
	})

	return cmd
}

func (cmd *validatorsCmd) listValidators(*cobra.Command, []string) (err error) {
	err = cmd.initDB()
	if err != nil {
		return fmt.Errorf("failed to initialize validator database: %w", err)
	}

	for _, name := range cmd.ctx.db.AvailableNames() {
		fmt.Fprintln(cmd.OutOrStdout(), name)
	}

	return nil
}

func (cmd *validatorsCmd) removeValidators(_ *cobra.Command, args []string) (err error) {
	err = cmd.initDB()
	if err != nil {
		return fmt.Errorf("failed to initialize validator database: %w", err)
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
		fmt.Fprintf(cmd.OutOrStdout(), "deleted versions\n")
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
		return fmt.Errorf("failed to initialize validator database: %w", err)
	}

	name := args[0]
	path := args[1]

	_, err = cmd.ctx.db.CompileCustom(name, path, "")
	if err != nil {
		return fmt.Errorf("failed to compile validator: %w", err)
	}

	fmt.Fprintln(cmd.OutOrStdout(), name)

	return nil
}

func (cmd *validatorsCmd) printCacheDir(_ *cobra.Command, _ []string) (err error) {
	path, err := validators.GetGitccCacheDir()
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
