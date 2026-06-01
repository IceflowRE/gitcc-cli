package regex

import (
	"fmt"
	"regexp"

	"github.com/go-git/go-git/v6/plumbing/object"

	"github.com/IceflowRE/gitcc/v3/standalone/gitcc"
)

const Name = "regex"

type Validator struct {
	gitcc.BaseValidator

	summaryRx     *regexp.Regexp
	descriptionRx *regexp.Regexp
}

func NewValidator() (gitcc.Validator, error) {
	return &Validator{}, nil
}

func (v *Validator) SetOptions(options map[string]string) error {
	v.Options = options

	if summaryRx, ok := options["summary"]; ok {
		rx, err := regexp.Compile(summaryRx)
		if err != nil {
			return fmt.Errorf("invalid summary regex: %w", err)
		}
		v.summaryRx = rx
	}

	if descriptionRx, ok := options["description"]; ok {
		rx, err := regexp.Compile(descriptionRx)
		if err != nil {
			return fmt.Errorf("invalid description regex: %w", err)
		}
		v.descriptionRx = rx
	}

	return nil
}

func (v *Validator) Validate(commit *object.Commit) gitcc.Result {
	return v.validateMessage(commit.Message)
}

func (v *Validator) validateMessage(message string) gitcc.Result {
	summary, description := gitcc.SplitCommitMessage(message)
	if v.summaryRx != nil && !v.summaryRx.MatchString(summary) {
		return gitcc.Result{
			Status:  gitcc.Invalid,
			Message: fmt.Sprintf("Summary does not match the pattern '%s'", v.summaryRx.String()),
		}
	}

	if v.descriptionRx != nil && !v.descriptionRx.MatchString(description) {
		return gitcc.Result{
			Status:  gitcc.Invalid,
			Message: fmt.Sprintf("Description does not match the pattern '%s'", v.descriptionRx.String()),
		}
	}

	return gitcc.Result{
		Status:  gitcc.Valid,
		Message: "",
	}
}
