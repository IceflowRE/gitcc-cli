// Package internal contains internal utilities for gitcc.
package internal

import (
	"io"
	"os"
)

// CopyFile copies a file from src to dst. If dst does not exist, it will be created.
func CopyFile(src string, dst string) error {
	srcFile, err := os.Open(src) //nolint:gosec
	if err != nil {
		return err
	}

	defer func() {
		cerr := srcFile.Close()
		if err == nil {
			err = cerr
		}
	}()

	outFile, err := os.Create(dst) //nolint:gosec
	if err != nil {
		return err
	}

	defer func() {
		cerr := outFile.Close()
		if err == nil {
			err = cerr
		}
	}()

	_, err = io.Copy(outFile, srcFile)
	if err != nil {
		return err
	}

	err = outFile.Sync()

	return err
}
