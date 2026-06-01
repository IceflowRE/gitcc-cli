package gitcc

import (
	"strings"
)

func MessageToDescription(message string) string {
	_, description := SplitCommitMessage(message)
	return description
}

func MessageToSummary(message string) string {
	summary, _ := SplitCommitMessage(message)
	return summary
}

func SplitCommitMessage(message string) (summary string, description string) {
	summary, description, found := strings.Cut(message, "\n\n")

	if !found {
		return strings.TrimSpace(message), ""
	}

	return strings.TrimSpace(summary), strings.TrimSpace(description)
}
