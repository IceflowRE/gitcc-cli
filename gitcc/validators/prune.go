package validators

import (
	"os"
	"path/filepath"
)

// PruneValidators removes all validator versions from the cache directory except for the current version and returns the list of deleted directories.
func PruneValidators() (deletedDirs []string, err error) {
	gitccDir, err := GetGitccCacheDir()
	if err != nil {
		return nil, err
	}
	curValidatorDir, err := GetValidatorCacheDir()
	if err != nil {
		return nil, err
	}

	return deleteValidators(gitccDir, curValidatorDir)
}

func deleteValidators(gitccDir string, curValidatorDir string) (deletedDirs []string, err error) {
	root, err := os.OpenRoot(gitccDir)
	if err != nil {
		return nil, err
	}
	defer root.Close() //nolint:errcheck

	deletedDirs = []string{}
	err = filepath.WalkDir(gitccDir, func(path string, dir os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == gitccDir {
			return nil
		}
		if path == curValidatorDir {
			return filepath.SkipDir
		}
		if dir.IsDir() {
			relPath, err := filepath.Rel(gitccDir, path)
			if err != nil {
				return err
			}
			err = root.RemoveAll(relPath)
			if err != nil {
				return err
			}
			deletedDirs = append(deletedDirs, filepath.Base(path))

			return filepath.SkipDir
		}

		return nil
	})

	return deletedDirs, err
}
