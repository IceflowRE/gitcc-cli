package gitcc

import (
	"github.com/go-git/go-git/v6/plumbing/object"
)

type Validator interface {
	SetOptions(options map[string]string) error
	Validate(commit *object.Commit) Result
}
