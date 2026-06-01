package validators

import (
	"os"
	"path/filepath"
)

func PruneValidators() (deletedDirs []string, err error) {
	gitccDir, err := GetGitccCacheDir()
	if err != nil {
		return nil, err
	}
	curValidatorDir, err := getValidatorCacheDir()
	if err != nil {
		return nil, err
	}

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
			err := os.RemoveAll(path)
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
