// Package gitcc provides the core functionality for gitcc.
package gitcc

import (
	"github.com/go-git/go-git/v6/plumbing/object"
)

// ValidatorConstructor is a function type that constructs a new Validator with the given options.
type ValidatorConstructor func(options map[string]string) (Validator, error)

// Validator defines the interface for validating commit messages.
// SetOptions is called once before any validation.
type Validator interface {
	Validate(commit *object.Commit) Result
}
