package utils

import (
	"errors"
	"fmt"
	"git_clone/gvc/config"
	"os"
	"path/filepath"
	"sync"
)

var (
	RepoDir string
	once    sync.Once
	initErr error
)

// FindRepo initializes the RepoDir variable and ensures it is set only once.
func FindRepo() error {
	once.Do(func() {
		// Get the current working directory
		cwd, err := os.Getwd()
		if err != nil {
			initErr = fmt.Errorf("error finding cwd: %w", err)
			return
		}

		for {
			repoPath := filepath.Join(cwd, config.OWN_FOLDER_NAME)
			if _, err := os.Stat(repoPath); err == nil {
				RepoDir = repoPath
				return
			}

			parentDir := filepath.Dir(cwd)
			if parentDir == cwd {
				initErr = errors.New("repository folder not found")
				return
			}

			cwd = parentDir
		}
	})

	return initErr
}
