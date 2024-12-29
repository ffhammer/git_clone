package utils

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func FindMatchingFiles(filePath string) ([]string, error) {
	var files []string

	// Check if the filePath is a directory
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not access %s: %w", filePath, err)
	}

	if fileInfo.IsDir() {
		// Add all files in the directory recursively
		err = filepath.Walk(filePath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			// Only add regular files (not subdirectories)
			if !info.IsDir() {
				files = append(files, path)
			}
			return nil
		})

		if err != nil {
			return nil, fmt.Errorf("error walking directory %s: %w", filePath, err)
		}
	} else if strings.ContainsAny(filePath, "*?") {
		// Handle glob pattern
		globMatches, err := filepath.Glob(filePath)
		if err != nil {
			return nil, fmt.Errorf("error processing glob pattern %s: %w", filePath, err)
		}
		files = append(files, globMatches...)
	} else {
		// Single file path, add directly
		files = append(files, filePath)

	}
	return files, nil
}

func GetBasePath() string {
	return filepath.Dir(RepoDir)
}

func RelPathToAbs(relPath string) string {
	return filepath.Join(GetBasePath(), relPath)
}

func MakePathRelativeToRepo(RepoDIr string, filePath string) (string, error) {

	absRepoDir, err := filepath.Abs(GetBasePath())
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

	// Return the relative path of filePath from RepoDIr
	relPath, err := filepath.Rel(absRepoDir, absFilePath)
	if err != nil {
		return "", fmt.Errorf("error computing relative path: %w", err)
	}

	return filepath.Clean(relPath), nil
}
