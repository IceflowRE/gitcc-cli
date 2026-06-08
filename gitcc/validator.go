package gitcc

import (
	"github.com/go-git/go-git/v6/plumbing/object"
)

type ValidatorConstructor func(options map[string]string) (Validator, error)

// Validator defines the interface for validating commit messages.
// SetOptions is called once before any validation.
type Validator interface {
	Validate(commit *object.Commit) Result
}
