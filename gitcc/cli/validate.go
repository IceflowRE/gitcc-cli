package cli

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"

	"github.com/IceflowRE/gitcc/v3/standalone/gitcc"
	"github.com/IceflowRE/gitcc/v3/standalone/gitcc/validators"
)

type validationContext struct {
	// validator to use
	validator gitcc.Validator
	// db is not loaded if validator is set
	db *validators.DB
}

func newValidationContext(validator gitcc.Validator) *validationContext {
	return &validationContext{
		validator: validator,
	}
}

type validateBaseCmd struct {
	*cobra.Command

	ctx *validationContext

	validatorName    string
	validatorPath    string
	compileIfMissing bool
	options          map[string]string
}

func newValidateBaseCmd(ccmd *cobra.Command, ctx *validationContext) *validateBaseCmd {
	cmd := &validateBaseCmd{
		Command: ccmd,
		ctx:     ctx,
	}

	cmd.PersistentFlags().StringVarP(&cmd.validatorName, "name", "", "",
		"Name of the validator.")
	cmd.PersistentFlags().StringVarP(&cmd.validatorPath, "path", "", "",
		"Path to a go file with a validator implementation.")
	cmd.PersistentFlags().BoolVarP(&cmd.compileIfMissing, "compile", "c", false,
		"Compile the validator if it is outdated or missing. This flag is only applicable when --path is specified. Note that this will compile and execute code from the specified path.")
	cmd.PersistentFlags().StringToStringVarP(&cmd.options, "options", "o", nil,
		"Options to pass to the validator in the format --options key=value. This flag can be specified multiple times for multiple options.")
	if ctx.validator == nil {
		cmd.PersistentFlags().MarkHidden("compile")
		cmd.PersistentFlags().MarkHidden("name")
		cmd.PersistentFlags().MarkHidden("path")
		cmd.PreRunE = cmd.preRunE
	}

	return cmd
}

func (cmd *validateBaseCmd) preRunE(*cobra.Command, []string) (err error) {
	if cmd.validatorName == "" && cmd.validatorPath == "" {
		return errors.New("must specify either --name or --path")
	}
	cmd.ctx.db, err = validators.NewDB()
	if err != nil {
		return fmt.Errorf("failed to initialize validator database: %w", err)
	}

	// check builtins
	if cmd.validatorPath == "" {
		cmd.ctx.validator, err = cmd.ctx.db.GetBuiltin(cmd.validatorName)
		if err != nil {
			return err
		}
	}

	return nil
}

func (cmd *validateBaseCmd) validate(action func(validator gitcc.Validator) error) error {
	if cmd.ctx.validator == nil {
		return cmd.execExternal()
	}

	err := cmd.ctx.validator.SetOptions(cmd.options)
	if err != nil {
		return fmt.Errorf("failed to set options: %w", err)
	}

	return action(cmd.ctx.validator)
}

func (cmd *validateBaseCmd) execExternal() (err error) {
	validatorExecutable, err := cmd.getExternalValidator()
	if err != nil {
		return err
	}

	valCmd := exec.Command(validatorExecutable, os.Args[1:]...)
	valCmd.Stdin = os.Stdin
	valCmd.Stdout = os.Stdout
	valCmd.Stderr = os.Stderr

	err = valCmd.Run()
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return &ExitError{exitErr.ExitCode()}
	}
	if err != nil {
		return err
	}

	return nil
}

func (cmd *validateBaseCmd) getExternalValidator() (validatorExecutable string, err error) {
	switch {
	case cmd.validatorPath != "" && cmd.compileIfMissing:
		validatorExecutable, err = cmd.ctx.db.GetOrCompileCustom(cmd.validatorPath, cmd.validatorName)
		if err != nil {
			return "", fmt.Errorf("failed to get or compile validator: %w", err)
		}
	case cmd.validatorPath != "":
		validatorExecutable = cmd.ctx.db.GetCustom(cmd.validatorPath)
		if validatorExecutable == "" {
			return "", fmt.Errorf("validator does not exist, compile it first")
		}
	case cmd.validatorName != "":
		validatorExecutable = cmd.ctx.db.GetCustomByName(cmd.validatorName)
		if validatorExecutable == "" {
			return "", fmt.Errorf("validator '%s' not found", cmd.validatorName)
		}
	default:
	}

	return validatorExecutable, nil
}
