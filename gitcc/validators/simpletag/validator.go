// Package simpletag implements a simple tag validator.
package simpletag

import (
	"regexp"

	"github.com/go-git/go-git/v6/plumbing/object"

	"github.com/IceflowRE/gitcc-cli/v3/gitcc"
)

var (
	rxParser      = regexp.MustCompile(`^\[(.*)\] (.*)$`)
	rxCategory    = regexp.MustCompile(`^(?:\*|(?:[a-z0-9]{2,}|[a-z0-9][ -][a-z0-9]+)(?:[ -][a-z0-9]+)*(?:\|(?:[a-z0-9]{2,}|[a-z0-9][ -][a-z0-9]+)(?:[ -][a-z0-9]+)*)*)$`) //nolint:lll
	rxDescription = regexp.MustCompile(`^[A-Z0-9]\S*(?:\s\S*)+[^.!?,\s]$`)
)

// Name is the validators name.
const Name = "simpletag"

// Validator checks if the commit message summary starts with a tag in square brackets followed by a description.
// The tag should be either
// - a single '*' or
// - completely lowercase letters or numbers, at least 2 characters long, other allowed characters are: '|', '-' and spaces.
// The tag can also contain multiple categories separated by '|', for example: [feat|fix].
// The description should start with an uppercase letter or number, should be not to short and should not end with a punctuation.
type Validator struct{}

// NewValidator create a new simpletag validator.
func NewValidator(_options map[string]string) (*Validator, error) {
	return &Validator{}, nil
}

// Validate validates a commit.
func (v *Validator) Validate(commit *object.Commit) gitcc.Result {
	return v.validateSummary(gitcc.MessageToSummary(commit.Message))
}

func (*Validator) validateSummary(summary string) gitcc.Result {
	matches := rxParser.FindStringSubmatch(summary)

	if len(matches) != 3 { //nolint:mnd
		return gitcc.Result{
			Status:  gitcc.Invalid,
			Message: "Summary format is invalid. It must follow '[<tag>] <Good Description>'",
		}
	}

	if !rxCategory.MatchString(matches[1]) {
		return gitcc.Result{
			Status: gitcc.Invalid,
			Message: "Invalid category tag. It should be either a single '*' or completely lowercase " +
				"letters or numbers, at least 2 characters long, other allowed characters are: '|', '-' and spaces.",
		}
	}

	if !rxDescription.MatchString(matches[2]) {
		return gitcc.Result{
			Status: gitcc.Invalid,
			Message: "Invalid description. It should start with an uppercase letter or number, " +
				"should be not to short and should not end with a punctuation.",
		}
	}

	return gitcc.Result{
		Status:  gitcc.Valid,
		Message: "",
	}
}
