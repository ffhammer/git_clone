package utils

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gobwas/glob"
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
	} else if IsGlob(filePath) {
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
		return "", fmt.Errorf("error getting absolute path of filePath '%s': %w", filePath, err)
	}

	// Check if absFilePath is within absRepoDir
	if !filepath.HasPrefix(absFilePath, absRepoDir) {
		return "", fmt.Errorf("file '%s' is outside of the repository directory", filePath)
	}

	// Return the relative path of filePath from RepoDIr
	relPath, err := filepath.Rel(absRepoDir, absFilePath)
	if err != nil {
		return "", fmt.Errorf("error computing relative path: %w", err)
	}

	return filepath.Clean(relPath), nil
}

func MatchFileWithMapStringKey[T any](relPath string, m map[string]T) []string {
	files := make([]string, 0)

	if IsGlob(relPath) {
		g := glob.MustCompile(relPath)

		for k, _ := range m {
			if g.Match(k) {
				files = append(files, k)
			}
		}

		return files
	}

	querySplittedParts := SplitPath(relPath)

	for k, _ := range m {
		allMatched := true
		keySplittedParts := SplitPath(k)

		if len(querySplittedParts) > len(keySplittedParts) {
			continue
		}

		for i := 0; i < len(querySplittedParts); i++ {

			allMatched = allMatched && (querySplittedParts[i] == keySplittedParts[i])
			if !allMatched {
				break
			}
		}

		if allMatched {
			files = append(files, k)
		}

	}

	return files
}

func FilterRelPathKeyMapWithAbsPaths[T any](m map[string]T, absPaths []string) (map[string]T, error) {

	if len(absPaths) == 0 {
		return nil, errors.New("this functios is intented not be called with emty paths")
	}

	filteredMap := make(map[string]T, 0)

	for _, absPath := range absPaths {
		relPath, err := MakePathRelativeToRepo(RepoDir, absPath)
		if err != nil {
			return nil, fmt.Errorf("could not generate relative path for '%s': %w", absPath, err)
		}

		for _, matchedPath := range MatchFileWithMapStringKey(relPath, m) {

			filteredMap[matchedPath] = m[matchedPath]
		}

	}
	return filteredMap, nil

}
