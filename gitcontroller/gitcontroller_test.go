package gitcontroller

import (
	"os"
	"path/filepath"
	"testing"
	"time"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/stretchr/testify/assert"
)

func TestOpenRepository(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tmpRepoDir := t.TempDir()
		if _, err := git.PlainInit(tmpRepoDir, false); err != nil {
			t.Fatal(err)
		}
		repo, err := OpenRepository(tmpRepoDir)
		assert.NoError(t, err)
		assert.NotNil(t, repo)
	})
	t.Run("failure: not repository", func(t *testing.T) {
		tmpNotRepoDir := t.TempDir()
		repo, err := OpenRepository(tmpNotRepoDir)
		assert.Error(t, err)
		assert.Nil(t, repo)
	})
}

func TestLoadAutor(t *testing.T) {
	tmpRepoDir := t.TempDir()
	if _, err := git.PlainInit(tmpRepoDir, false); err != nil {
		t.Fatal(err)
	}
	repo, err := git.PlainOpenWithOptions(tmpRepoDir, &git.PlainOpenOptions{DetectDotGit: true})
	if err != nil {
		t.Fatal(err)
	}
	t.Run("success", func(t *testing.T) {
		t.Setenv("HOME", t.TempDir())
		t.Setenv("GIT_CONFIG_GLOBAL", "/dev/null")
		cfg := "[user]\n\tname = test\n\temail = test@example.com\n"
		if err := os.WriteFile(filepath.Join(tmpRepoDir, ".gitconfig"), []byte(cfg), 0666); err != nil {
			t.Fatal(err)
		}
		t.Setenv("HOME", tmpRepoDir)
		author, err := LoadAuthor(repo)
		assert.NoError(t, err)
		assert.NotNil(t, author)
	})
	t.Run("failure: malformed global config", func(t *testing.T) {
		t.Setenv("HOME", t.TempDir())
		t.Setenv("GIT_CONFIG_GLOBAL", "/dev/null")
		cfg := "[user\n\tname = test\n"
		if err := os.WriteFile(filepath.Join(tmpRepoDir, ".gitconfig"), []byte(cfg), 0666); err != nil {
			t.Fatal(err)
		}
		t.Setenv("HOME", tmpRepoDir)
		author, err := LoadAuthor(repo)
		assert.Error(t, err)
		assert.Nil(t, author)
	})
	t.Run("failure: undefined git config", func(t *testing.T) {
		t.Setenv("HOME", t.TempDir())
		t.Setenv("GIT_CONFIG_GLOBAL", "/dev/null")
		author, err := LoadAuthor(repo)
		assert.Error(t, err)
		assert.Nil(t, author)
	})
}

func TestCommit(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tmpRepoDir := t.TempDir()
		tmpFile := "sample.txt"
		f := filepath.Join(tmpRepoDir, tmpFile)
		err := os.WriteFile(f, []byte("hoge"), 0666)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := git.PlainInit(tmpRepoDir, false); err != nil {
			t.Fatal(err)
		}
		repo, err := git.PlainOpenWithOptions(tmpRepoDir, &git.PlainOpenOptions{DetectDotGit: true})
		if err != nil {
			t.Fatal(err)
		}
		w, err := repo.Worktree()
		if err != nil {
			t.Fatal(err)
		}
		if _, err := w.Add(tmpFile); err != nil {
			t.Fatal(err)
		}
		commitMsg := "feat: test"
		author := &object.Signature{
			Name:  "test",
			Email: "test@example.com",
			When:  time.Now(),
		}
		err = Commit(repo, commitMsg, author)
		assert.NoError(t, err)

		head, err := repo.Head()
		assert.NoError(t, err)
		assert.NotNil(t, head)
	})
	t.Run("failure: nothing staged", func(t *testing.T) {
		tmpRepoDir := t.TempDir()
		if _, err := git.PlainInit(tmpRepoDir, false); err != nil {
			t.Fatal(err)
		}
		repo, err := git.PlainOpenWithOptions(tmpRepoDir, &git.PlainOpenOptions{DetectDotGit: true})
		if err != nil {
			t.Fatal(err)
		}
		commitMsg := "feat: test"
		author := &object.Signature{
			Name:  "test",
			Email: "test@example.com",
			When:  time.Now(),
		}
		err = Commit(repo, commitMsg, author)
		assert.Error(t, err)
	})
}
