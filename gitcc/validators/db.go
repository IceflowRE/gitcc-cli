// Package validators provides functionality for managing and compiling custom validators for commit messages.
package validators

import (
	"context"
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
	"errors"
	"fmt"
	"maps"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"slices"
	"strings"
	"time"

	"golang.org/x/mod/semver"

	"github.com/IceflowRE/gitcc-cli/v3/gitcc"
	"github.com/IceflowRE/gitcc-cli/v3/gitcc/internal"
	"github.com/IceflowRE/gitcc-cli/v3/gitcc/validators/regex"
	"github.com/IceflowRE/gitcc-cli/v3/gitcc/validators/simpletag"
)

// ErrValidatorNotFound is returned when a requested validator is not found in the database.
var ErrValidatorNotFound = errors.New("validator not found")

// DB represents a database of validators, including both built-in and custom validators.
type DB struct {
	builtin          map[string]func() (gitcc.Validator, error)
	validatorDir     string
	customValidators []validatorMeta
}

// NewDB initializes a new DB instance by loading built-in validators and refreshing custom validators from the cache directory.
func NewDB() (*DB, error) {
	valCacheDir, err := getValidatorCacheDir()
	if err != nil {
		return nil, err
	}
	db := &DB{
		builtin: map[string]func() (gitcc.Validator, error){
			regex.Name:     func() (gitcc.Validator, error) { return regex.NewValidator() },
			simpletag.Name: func() (gitcc.Validator, error) { return simpletag.NewValidator() },
		},
		validatorDir: valCacheDir,
	}
	err = db.refreshCustomValidators()
	if err != nil {
		return nil, err
	}

	return db, nil
}

// AvailableNames returns a list of all available validator names, including both built-in and custom validators.
func (db *DB) AvailableNames() []string {
	names := slices.Collect(maps.Keys(db.builtin))

	for _, meta := range db.customValidators {
		names = append(names, meta.Name)
	}

	return names
}

// GetBuiltin retrieves a built-in validator by its name.
// If the validator is not found, ErrValidatorNotFound is returned.
func (db *DB) GetBuiltin(name string) (gitcc.Validator, error) { //nolint:ireturn
	validatorFn, ok := db.builtin[name]
	if ok {
		return validatorFn()
	}

	return nil, ErrValidatorNotFound
}

// GetCustom retrieves a custom validator by the file path of its source code.
// Only absolute paths should be passed.
func (db *DB) GetCustom(path string) string {
	hash, err := getShortSha256(path)
	if err != nil {
		return ""
	}

	return db.getCustomByHash(hash)
}

// GetCustomByName retrieves a custom validator by its name.
func (db *DB) GetCustomByName(name string) string {
	idx := slices.IndexFunc(db.customValidators, func(elem validatorMeta) bool {
		return elem.Name == name
	})
	if idx == -1 {
		return ""
	}

	return filepath.Join(db.validatorDir, db.customValidators[idx].Filename())
}

// GetOrCompileCustom retrieves a custom validator by the file path or compiles it if it does not exist or is outdated.
func (db *DB) GetOrCompileCustom(path string, name string) (string, error) {
	hash, err := getShortSha256(path)
	if err != nil {
		return "", fmt.Errorf("failed to get validator hash: %w", err)
	}

	exePath := db.getCustomByHash(hash)
	if exePath != "" {
		return exePath, nil
	}

	return db.CompileCustom(path, name, hash)
}

// ErrInvalidNameAndPath is returned when both name and path are invalid.
var ErrInvalidNameAndPath = errors.New("invalid name and path")

// CompileCustom compiles a custom validator from the specified source file and stores it in the database.
// If a validator with the same name or hash already exists, it will be replaced.
// If the hash is not provided, it will be computed from the source file.
func (db *DB) CompileCustom(path string, name string, hash string) (validatorPath string, err error) {
	if name == "" && path == "" {
		return "", ErrInvalidNameAndPath
	}

	if hash == "" {
		hash, err = getShortSha256(path)
		if err != nil {
			return "", fmt.Errorf("failed to get validator hash: %w", err)
		}
	}

	// remove old validators
	for _, meta := range db.customValidators {
		if meta.Name == name || meta.Hash == hash {
			err := os.Remove(filepath.Join(db.validatorDir, meta.Filename()))
			if err != nil {
				return "", fmt.Errorf("failed to remove old validator: %w", err)
			}
		}
	}

	return db.compile(name, path, hash)
}

func (db *DB) getCustomByHash(hash string) string {
	idx := slices.IndexFunc(db.customValidators, func(elem validatorMeta) bool {
		return elem.Hash == hash
	})
	if idx == -1 {
		return ""
	}

	return filepath.Join(db.validatorDir, db.customValidators[idx].Filename())
}

//go:embed main.tmpl
var mainFile []byte

func (db *DB) compile(name string, path string, hash string) (string, error) {
	dir, err := os.MkdirTemp("", "*")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(dir) //nolint:errcheck

	// create main.go
	mainPath := filepath.Join(dir, "main.go")
	err = os.WriteFile(mainPath, mainFile, 0o600) //nolint:mnd
	if err != nil {
		return "", fmt.Errorf("write main.go: %w", err)
	}

	// copy custom validator
	err = internal.CopyFile(path, filepath.Join(dir, "validator.go"))
	if err != nil {
		return "", fmt.Errorf("copy validator: %w", err)
	}

	// create go.mod
	modData := []byte("module github.com/IceflowRE/gitcc-cli/v3/custom")
	err = os.WriteFile(filepath.Join(dir, "go.mod"), modData, 0o600) //nolint:mnd
	if err != nil {
		return "", fmt.Errorf("write go.mod: %w", err)
	}

	// go mod tidy
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Minute) //nolint:mnd
	defer cancel()

	tidyCmd := exec.CommandContext(ctx, "go", "mod", "tidy")
	tidyCmd.Dir = dir
	out, err := tidyCmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("go mod tidy: %w\n%s", err, out)
	}

	// compile
	cCtx, cCancel := context.WithTimeout(context.Background(), 4*time.Minute) //nolint:mnd
	defer cCancel()

	outPath := filepath.Join(db.validatorDir, executableName(name, hash))
	buildCmd := exec.CommandContext(cCtx, "go", "build", "-o", outPath, ".") //nolint:gosec
	buildCmd.Dir = dir
	out, err = buildCmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("go build: %w\n%s", err, out)
	}

	return outPath, nil
}

func (db *DB) refreshCustomValidators() error {
	db.customValidators = []validatorMeta{}
	err := filepath.WalkDir(db.validatorDir, func(path string, dir os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if dir.IsDir() && path != db.validatorDir {
			return filepath.SkipDir
		}

		meta := getMetaFromName(dir.Name())
		if meta != nil {
			db.customValidators = append(db.customValidators, *meta)
		}

		return nil
	})
	if err != nil {
		db.customValidators = []validatorMeta{}

		return err
	}

	return nil
}

func executableName(name string, hash string) string {
	exeName := fmt.Sprintf("%s-%s", name, hash)
	if runtime.GOOS == "windows" {
		exeName += ".exe"
	}

	return exeName
}

// GetGitccCacheDir returns the path to the cache directory for gitcc, creating it if it does not exist.
func GetGitccCacheDir() (string, error) {
	userCacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	cacheDir := filepath.Join(userCacheDir, "gitcc")
	err = os.MkdirAll(cacheDir, 0o750) //nolint:mnd
	if err != nil {
		return "", err
	}

	return cacheDir, nil
}

func getValidatorCacheDir() (string, error) {
	cacheDir, err := GetGitccCacheDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(cacheDir, getCurrentVersion())
	err = os.MkdirAll(dir, 0o750) //nolint:mnd
	if err != nil {
		return "", err
	}

	return dir, nil
}

func shortSHA256(data []byte) string {
	sum := sha256.Sum256(data)

	return hex.EncodeToString(sum[:])[:10]
}

func getShortSha256(path string) (string, error) {
	data, err := os.ReadFile(path) //nolint:gosec
	if err != nil {
		return "", err
	}

	return shortSHA256(data), nil
}

type validatorMeta struct {
	Name string
	Hash string
}

func (meta *validatorMeta) Filename() string {
	return executableName(meta.Name, meta.Hash)
}

func getMetaFromName(filename string) *validatorMeta {
	filename = strings.TrimPrefix(filename, filepath.Ext(filename))

	idx := strings.LastIndex(filename, "-")
	if idx == -1 {
		return nil
	}

	return &validatorMeta{
		Name: filename[:idx],
		Hash: filename[idx+1:],
	}
}

func getCurrentVersion() string {
	info, ok := debug.ReadBuildInfo()
	if ok && semver.IsValid(info.Main.Version) {
		return strings.ReplaceAll(info.Main.Version, ".", "_")
	}

	return "v0"
}
