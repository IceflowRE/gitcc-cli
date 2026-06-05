package internal_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/IceflowRE/gitcc/v3/standalone/gitcc"
	"github.com/IceflowRE/gitcc/v3/standalone/gitcc/internal"
	"github.com/IceflowRE/gitcc/v3/standalone/gitcc/validators/simpletag"
)

//go:generate ../../create_testdata.sh testdata/simpletag

func TestHeadSimpleTag(t *testing.T) {
	t.Parallel()

	repoOk, err := internal.LoadRepository("testdata/simpletag")
	require.NoError(t, err)

	repoFail, err := internal.LoadRepository("testdata/simpletag_fail")
	require.NoError(t, err)

	validator, err := simpletag.NewValidator()
	require.NoError(t, err)

	res, err := internal.ValidateHead(validator, repoOk)
	require.NoError(t, err)
	assert.Equal(t, gitcc.Valid, res.Status)

	res, err = internal.ValidateHead(validator, repoFail)
	require.NoError(t, err)
	assert.Equal(t, gitcc.Invalid, res.Status)
}

func TestHistorySimpleTag(t *testing.T) {
	t.Parallel()

	repoOk, err := internal.LoadRepository("testdata/simpletag")
	require.NoError(t, err)

	repoFail, err := internal.LoadRepository("testdata/simpletag_fail")
	require.NoError(t, err)

	validator, err := simpletag.NewValidator()
	require.NoError(t, err)

	res, err := internal.ValidateHistory(validator, repoOk, "")
	require.NoError(t, err)
	assert.ElementsMatch(t, []gitcc.Status{
		gitcc.Valid, gitcc.Valid, gitcc.Valid, gitcc.Valid, gitcc.Valid, gitcc.Valid, gitcc.Valid, gitcc.Valid, gitcc.Valid,
	}, collectStatus(res))

	res, err = internal.ValidateHistory(validator, repoFail, "")
	require.NoError(t, err)
	assert.ElementsMatch(t, []gitcc.Status{
		gitcc.Invalid, gitcc.Valid, gitcc.Invalid, gitcc.Invalid, gitcc.Valid, gitcc.Valid, gitcc.Invalid, gitcc.Invalid, gitcc.Valid,
	}, collectStatus(res))
}

func collectStatus(results []gitcc.Result) []gitcc.Status {
	statuses := make([]gitcc.Status, len(results))
	for idx, res := range results {
		statuses[idx] = res.Status
	}

	return statuses
}
