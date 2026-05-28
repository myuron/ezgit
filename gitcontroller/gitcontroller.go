package gitcontroller

import (
	"errors"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// LoadAuthor ... A function that load the commit author
func LoadAuthor(path string) (*object.Signature, error) {
	repo, err := git.PlainOpenWithOptions(path, &git.PlainOpenOptions{DetectDotGit: true})
	if err != nil {
		return nil, err
	}

	cfg, err := repo.ConfigScoped(config.GlobalScope)
	if err != nil {
		return nil, err
	}

	if cfg.User.Name == "" || cfg.User.Email == "" {
		return nil, errors.New("git user.name and user.email must be configured")
	}

	author := &object.Signature{
		Name:  cfg.User.Name,
		Email: cfg.User.Email,
		When:  time.Now(),
	}

	return author, nil
}

// Commit ... A function that commits based on the received message
func Commit(path string, commitMsg string, author *object.Signature) error {
	repo, err := git.PlainOpenWithOptions(path, &git.PlainOpenOptions{DetectDotGit: true})
	if err != nil {
		return err
	}

	w, err := repo.Worktree()
	if err != nil {
		return err
	}

	status, err := w.Status()
	if err != nil {
		return err
	}

	hasStaged := false
	for _, s := range status {
		if s.Staging != git.Unmodified && s.Staging != git.Untracked {
			hasStaged = true
			break
		}
	}
	if !hasStaged {
		return errors.New("nothing staged to commit (use `git add` first)")
	}

	_, err = w.Commit(commitMsg, &git.CommitOptions{Author: author})
	return err
}
