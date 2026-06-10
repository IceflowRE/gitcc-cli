package main

import (
	"github.com/go-git/go-git/v6/plumbing/object"

	"github.com/IceflowRE/gitcc-cli/v3/gitcc"
)

type Validator struct{}

func NewValidator(options map[string]string) (gitcc.Validator, error) {
	return &Validator{}, nil
}

func (v *Validator) Validate(commit *object.Commit) gitcc.Result {
	return gitcc.Result{
		Status: gitcc.Valid,
		// Messages are ignored for valid results.
		Message: "This is a dummy validator that always returns valid. Please implement your own validator.",
	}
}
