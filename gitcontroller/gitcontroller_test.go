package gitcontroller

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func TestLoadAuthor(t *testing.T) {
	// Creating a temporary directory
	tmpDir := t.TempDir()

	// Creating a Repository
	if _, err := git.PlainInit(tmpDir, false); err != nil {
		t.Fatal(err)
	}

	t.Run("success", func(t *testing.T) {
		author, err := LoadAuthor(tmpDir)
		if err != nil {
			t.Fatal(err)
		}
		if author.Name == "" || author.Email == "" {
			t.Errorf("expected non-empty author, got %+v", author)
		}
	})
	t.Run("undefined Config", func(t *testing.T) {
		t.Setenv("HOME", t.TempDir())
		t.Setenv("GIT_CONFIG_GLOBAL", "/dev/null")
		author, err := LoadAuthor(tmpDir)
		if err == nil {
			t.Fatalf("expected error when git config is undefined, got author=%+v", author)
		}
		if author != nil {
			t.Errorf("expected nil author when err != nil, got %+v", author)
		}
	})
}

func TestCommit(t *testing.T) {
	// Creating a temporary directory
	tmpDir := t.TempDir()
	tmpFile := "sample.txt"
	f := filepath.Join(tmpDir, tmpFile)
	err := os.WriteFile(f, []byte("hoge"), 0666)
	if err != nil {
		t.Fatal(err)
	}

	// Creating a Repository
	repo, err := git.PlainInit(tmpDir, false)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("success", func(t *testing.T) {
		// Stage the file
		w, err := repo.Worktree()
		if err != nil {
			t.Fatal(err)
		}

		if _, err := w.Add(tmpFile); err != nil {
			t.Fatal(err)
		}

		author := &object.Signature{
			Name:  "test",
			Email: "test@example.com",
			When:  time.Now(),
		}

		// Commit
		commitMsg := "feat: init"
		if err := Commit(tmpDir, commitMsg, author); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("not staging file", func(t *testing.T) {
		// Stage the file
		_, err := repo.Worktree()
		if err != nil {
			t.Fatal(err)
		}

		author := &object.Signature{
			Name:  "test",
			Email: "test@example.com",
			When:  time.Now(),
		}

		// Commit
		commitMsg := "feat: init"
		if err := Commit(tmpDir, commitMsg, author); err != nil {
			t.Log(err)
		}
	})
}
