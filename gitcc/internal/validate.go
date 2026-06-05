package internal

import (
	"errors"
	"fmt"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/object"

	"github.com/IceflowRE/gitcc/v3/standalone/gitcc"
)

// ValidateCommit validates a specific commit by its SHA string.
func ValidateCommit(validator gitcc.Validator, repo *git.Repository, sha string) (gitcc.Result, error) {
	hash := plumbing.NewHash(sha)
	commit, err := repo.CommitObject(hash)
	if err != nil {
		return gitcc.Result{}, fmt.Errorf("failed to resolve commit %q: %w", sha, err)
	}

	res := validator.Validate(commit)
	res.Commit = commit

	return res, nil
}

// ValidateHead validates the current HEAD commit.
func ValidateHead(validator gitcc.Validator, repo *git.Repository) (gitcc.Result, error) {
	commit, err := getHeadCommit(repo)
	if err != nil {
		return gitcc.Result{}, fmt.Errorf("failed to get HEAD commit: %w", err)
	}
	res := validator.Validate(commit)
	res.Commit = commit

	return res, nil
}

// ValidateHistory validates the commit history starting from HEAD until the exitSha (exclusive). If exitSha is empty, it validates the entire history.
func ValidateHistory(validator gitcc.Validator, repo *git.Repository, exitSha string) ([]gitcc.Result, error) {
	ref, err := repo.Head()
	if err != nil {
		return nil, err
	}

	opts := &git.LogOptions{
		From: ref.Hash(),
	}
	if exitSha != "" {
		opts.To = plumbing.NewHash(exitSha)
	}
	cIter, err := repo.Log(opts)
	if err != nil {
		return nil, err
	}
	defer cIter.Close()

	var results []gitcc.Result

	err = cIter.ForEach(func(commit *object.Commit) error {
		res := validator.Validate(commit)
		res.Commit = commit
		results = append(results, res)

		return nil
	})

	return results, err
}

// ErrNoCommonAncestor is returned when no common ancestor is found between the source and target branches.
var ErrNoCommonAncestor = errors.New("no common ancestor found")

// GetMergeBase finds the common ancestor (merge base) between the current HEAD and the specified target branch.
// It returns ErrNoCommonAncestor if no common ancestor is found.
func GetMergeBase(repo *git.Repository, targetBranch string) (*object.Commit, error) {
	sourceCommit, err := getHeadCommit(repo)
	if err != nil {
		return nil, fmt.Errorf("failed to get HEAD commit: %w", err)
	}

	targetHash, err := repo.ResolveRevision(plumbing.Revision(targetBranch))
	if err != nil {
		return nil, fmt.Errorf("failed to resolve target branch: %w", err)
	}

	targetCommit, err := repo.CommitObject(*targetHash)
	if err != nil {
		return nil, err
	}

	res, err := sourceCommit.MergeBase(targetCommit)
	if err != nil || len(res) == 0 {
		return nil, ErrNoCommonAncestor
	}

	return res[0], nil
}

// LoadRepository opens a Git repository at the specified path. It detects the .git directory automatically.
func LoadRepository(path string) (*git.Repository, error) {
	repo, err := git.PlainOpenWithOptions(path, &git.PlainOpenOptions{
		DetectDotGit: true,
	})
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func getHeadCommit(repo *git.Repository) (*object.Commit, error) {
	ref, err := repo.Head()
	if err != nil {
		return nil, err
	}
	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		return nil, err
	}

	return commit, nil
}
