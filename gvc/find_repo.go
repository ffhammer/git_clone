package gvc

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var (
	repoDir string
)

func FindRepo() error {
	// Get the current working directory

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error finding cwd: %w", err)
	}

	for {
		repoPath := filepath.Join(cwd, OWN_FOLDER_NAME)
		if _, err := os.Stat(repoPath); err == nil {
			repoDir = repoPath
			return nil
		}

		parentDir := filepath.Dir(cwd)
		if parentDir == cwd {
			return errors.New("repository folder not found")
		}

		cwd = parentDir
	}
}
