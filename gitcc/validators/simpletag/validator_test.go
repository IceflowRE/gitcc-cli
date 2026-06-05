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
			want:    "Summary format is invalid. It must follow '[<tag>] <Good Description>'",
		},
		{
			name:    "invalid category tag",
			summary: "[FeaT] Add new feature",
			want:    "Invalid category tag. It should be either a single '*' or completely lowercase letters or numbers, at least 2 characters long, other allowed characters are: '|', '-' and spaces.", //nolint:lll
		},
		{
			name:    "invalid description",
			summary: "[feat] add new feature.",
			want:    "Invalid description. It should start with an uppercase letter or number, should be not to short and should not end with a punctuation.",
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
