package gitcc

import (
	"fmt"

	"github.com/go-git/go-git/v6/plumbing/object"
)

// Status represents the validation status of a commit.
type Status string

// Severity returns an integer representing the severity of the status, where higher values indicate more severe issues.
// This is used to sort.
func (s Status) Severity() int {
	switch s {
	case Valid:
		return 0
	case Warning:
		return 1
	case Invalid:
		return 2 //nolint:revive
	}

	return -1
}

const (
	// Invalid indicates that the commit message is invalid.
	Invalid Status = "Invalid"
	// Valid indicates that the commit message is valid.
	Valid Status = "OK"
	// Warning indicates that the commit message is valid, but does not meet the criteria completely.
	Warning Status = "Warning"
)

// Result represents the result of a validation.
type Result struct {
	Status  Status
	Message string
	Commit  *object.Commit
}

func (res *Result) String() string {
	msg := string(res.Status)
	if res.Commit != nil {
		msg = fmt.Sprintf("%s | %s | %s", msg, res.Commit.Hash.String(), MessageToSummary(res.Commit.Message))
	}
	if res.Message != "" {
		msg = fmt.Sprintf("%s\n    : %s", msg, res.Message)
	}

	return msg
}
