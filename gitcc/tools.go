package gitcc

import (
	"strings"
)

// MessageToDescription extracts the description part of a commit message.
func MessageToDescription(message string) string {
	_, description := SplitCommitMessage(message)

	return description
}

// MessageToSummary extracts the summary part of a commit message.
func MessageToSummary(message string) string {
	summary, _ := SplitCommitMessage(message)

	return summary
}

// SplitCommitMessage splits a commit message into summary and description parts.
func SplitCommitMessage(message string) (summary string, description string) {
	summary, description, found := strings.Cut(message, "\n\n")

	if !found {
		return strings.TrimSuffix(message, "\n"), ""
	}

	return summary, strings.TrimSuffix(description, "\n")
}
