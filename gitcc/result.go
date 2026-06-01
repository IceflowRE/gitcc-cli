package gitcc

import (
	"fmt"

	"github.com/go-git/go-git/v6/plumbing/object"
)

type Status string

func (s Status) Severity() int {
	switch s {
	case Valid:
		return 0
	case Warning:
		return 1
	case Invalid:
		return 2
	}

	return -1
}

const (
	Invalid Status = "Invalid"
	Valid   Status = "OK"
	Warning Status = "Warning"
)

type Result struct {
	Status  Status
	Message string
	Commit  *object.Commit
}

func (res *Result) String() string {
	msg := string(res.Status)
	if res.Commit != nil {
		msg = fmt.Sprintf("%s | %s | %s", msg, res.Commit.Hash.String(), MessageToSummary(res.Commit.Message))
	} else {
		msg += " |"
	}
	if res.Message != "" {
		msg = fmt.Sprintf("%s\n    : %s", msg, res.Message)
	}

	return msg
}
