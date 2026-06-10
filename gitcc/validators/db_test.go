package validators //nolint:testpackage

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCompile(t *testing.T) {
	t.Parallel()

	db := &DB{
		validatorDir: t.TempDir(),
	}

	buildDir := t.TempDir()
	outPath, err := db.compile("custom-test", filepath.Join("testdata", "validator.go"), "", buildDir)

	require.NoError(t, err)
	require.FileExists(t, outPath)
}
