package gitcc

import (
	"github.com/go-git/go-git/v6/plumbing/object"
)

// Validator defines the interface for validating commit messages.
// SetOptions is called once before any validation.
type Validator interface {
	SetOptions(options map[string]string) error
	Validate(commit *object.Commit) Result
}
