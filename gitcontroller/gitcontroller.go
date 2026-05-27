package gitcontroller

import (
	"errors"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// Commit ... A function that commits based on the received message
func Commit(message string) error {
	repo, err := git.PlainOpenWithOptions(".", &git.PlainOpenOptions{DetectDotGit: true})
	if err != nil {
		return err
	}

	w, err := repo.Worktree()
	if err != nil {
		return err
	}

	cfg, err := repo.ConfigScoped(config.GlobalScope)
	if err != nil {
		return err
	}

	if cfg.User.Name == "" || cfg.User.Email == "" {
		return errors.New("git user.name and user.email must be configured")
	}

	author := &object.Signature{
		Name:  cfg.User.Name,
		Email: cfg.User.Email,
		When:  time.Now(),
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

	_, err = w.Commit(message, &git.CommitOptions{Author: author})
	return err
}
