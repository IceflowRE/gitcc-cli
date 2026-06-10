package simpletag //nolint:testpackage

import (
	"testing"
)

func Test(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		summary string
		want    string
	}{
		{
			name:    "valid summary",
			summary: "[feat] Add new feature",
			want:    "",
		},
		{
			name:    "invalid format",
			summary: "Add new feature",
			want:    "Summary has invalid format. Expected: '[<tag>] <Description>'",
		},
		{
			name:    "invalid category tag",
			summary: "[FeaT] Add new feature",
			want:    "Category tag must be a single '*' or lowercase letters/numbers (at least 2 characters), optionally separated by '|', '-' or spaces.",
		},
		{
			name:    "invalid description",
			summary: "[feat] add new feature.",
			want:    "Description must start with an uppercase letter or number, be sufficiently long and not end with punctuation.",
		},
		{
			name:    "valid category tag with spaces",
			summary: "[feat new] Add new feature",
			want:    "",
		},
		{
			name:    "valid category tag with special characters",
			summary: "[feat-new|cool] Add new feature",
			want:    "",
		},
	}

	for _, tt := range tests { //nolint:varnamelen
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			res := (&Validator{}).validateSummary(tt.summary)
			if res.Message != tt.want {
				t.Errorf("got '%s', want '%s'", res.Message, tt.want)
			}
		})
	}
}
