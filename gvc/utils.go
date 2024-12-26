package gvc

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// mkdirIgnoreExists creates a directory if it doesn't already exist
func mkdirIgnoreExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.Mkdir(path, os.ModePerm)
	}
	return nil
}

func makePathRelativeToRepo(repoDir string, filePath string) (string, error) {

	absRepoDir, err := filepath.Abs(filepath.Dir(repoDir))
	if err != nil {
		return "", fmt.Errorf("error getting absolute path of repository: %w", err)
	}

	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		return "", fmt.Errorf("error getting absolute path of filePath: %w", err)
	}

	// Check if absFilePath is within absRepoDir
	if !filepath.HasPrefix(absFilePath, absRepoDir) {
		return "", errors.New("file is outside of the repository directory")
	}

	// Return the relative path of filePath from repoDir
	relPath, err := filepath.Rel(absRepoDir, absFilePath)
	if err != nil {
		return "", fmt.Errorf("error computing relative path: %w", err)
	}

	return filepath.Clean(relPath), nil
}

func getFileSHA1(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	hasher := sha1.New()

	if _, err := io.Copy(hasher, file); err != nil {
		return "", fmt.Errorf("error reading file content: %w", err)
	}

	hashSum := hasher.Sum(nil)

	return fmt.Sprintf("%x", hashSum), nil
}

func getStringSHA1(s string) string {
	h := sha1.New()

	h.Write([]byte(s))

	return fmt.Sprintf("%x", h.Sum(nil))
}
