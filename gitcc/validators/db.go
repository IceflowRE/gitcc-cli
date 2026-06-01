package validators

import (
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
	"fmt"
	"maps"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"slices"
	"strings"

	"golang.org/x/mod/semver"

	"github.com/IceflowRE/gitcc/v3/standalone/gitcc"
	"github.com/IceflowRE/gitcc/v3/standalone/gitcc/internal"
	"github.com/IceflowRE/gitcc/v3/standalone/gitcc/validators/regex"
	"github.com/IceflowRE/gitcc/v3/standalone/gitcc/validators/simpletag"
)

type DB struct {
	builtin          map[string]func() (gitcc.Validator, error)
	validatorDir     string
	customValidators []validatorMeta
}

func NewDB() (*DB, error) {
	valCacheDir, err := getValidatorCacheDir()
	if err != nil {
		return nil, err
	}
	db := &DB{
		builtin: map[string]func() (gitcc.Validator, error){
			regex.Name:     regex.NewValidator,
			simpletag.Name: simpletag.NewValidator,
		},
		validatorDir: valCacheDir,
	}
	err = db.refreshCustomValidators()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (db *DB) AvailableNames() []string {
	names := slices.Collect(maps.Keys(db.builtin))

	for _, meta := range db.customValidators {
		names = append(names, meta.Name)
	}

	return names
}

var ErrValidatorNotFound = fmt.Errorf("validator not found")

func (db *DB) GetBuiltin(name string) (gitcc.Validator, error) {
	validatorFn, ok := db.builtin[name]
	if ok {
		return validatorFn()
	}

	return nil, ErrValidatorNotFound
}

func (db *DB) GetCustom(path string) string {
	hash, err := getShortSha256(path)
	if err != nil {
		return ""
	}

	return db.getCustomByHash(hash)
}

func (db *DB) GetCustomByName(name string) string {
	idx := slices.IndexFunc(db.customValidators, func(elem validatorMeta) bool {
		return elem.Name == name
	})
	if idx == -1 {
		return ""
	}

	return filepath.Join(db.validatorDir, db.customValidators[idx].Filename())
}

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

// if hash is empty it will be calculated
func (db *DB) CompileCustom(path string, name string, hash string) (validatorPath string, err error) {
	if name == "" && path == "" {
		return "", fmt.Errorf("invalid name and path")
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
	defer os.RemoveAll(dir)

	// create main.go
	mainPath := filepath.Join(dir, "main.go")
	err = os.WriteFile(mainPath, mainFile, 0o600)
	if err != nil {
		return "", fmt.Errorf("write main.go: %w", err)
	}

	// copy custom validator
	err = internal.CopyFile(path, filepath.Join(dir, "validator.go"))
	if err != nil {
		return "", fmt.Errorf("copy validator: %w", err)
	}

	// create go.mod
	modData := []byte("module github.com/IceflowRE/gitcc/v3/standalone/custom")
	err = os.WriteFile(filepath.Join(dir, "go.mod"), modData, 0o600)
	if err != nil {
		return "", fmt.Errorf("write go.mod: %w", err)
	}

	// go mod tidy
	tidyCmd := exec.Command("go", "mod", "tidy")
	tidyCmd.Dir = dir
	out, err := tidyCmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("go mod tidy: %w\n%s", err, out)
	}

	// compile
	outPath := filepath.Join(db.validatorDir, executableName(name, hash))
	buildCmd := exec.Command("go", "build", "-o", outPath, ".")
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

func GetGitccCacheDir() (string, error) {
	userCacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	cacheDir := filepath.Join(userCacheDir, "gitcc")
	err = os.MkdirAll(cacheDir, 0o750)
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
	err = os.MkdirAll(dir, 0o750)
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
	data, err := os.ReadFile(path)
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
