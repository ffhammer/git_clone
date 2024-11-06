package gvc

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"

	"github.com/gobwas/glob"
)

var (
	ignorePatterns []string
	ignoreOnce     sync.Once
)

func getIgnorePatterns(repoPath string) []string {
	ignoreOnce.Do(func() {
		// Initialize ignorePatterns by parsing the ignore file
		patterns, err := parseIgnoreFile(repoPath)
		if err != nil {
			fmt.Printf("error parsing ignore file: %v. Ignoring patterns.\n", err)
			ignorePatterns = []string{}
		} else {
			ignorePatterns = patterns
		}
	})
	return ignorePatterns
}

func parseIgnoreFile(repoPath string) ([]string, error) {
	ignoreFilePath := filepath.Join(filepath.Dir(repoPath), IGNORE_PATH)

	// Check if the ignore file exists
	if _, err := os.Stat(ignoreFilePath); os.IsNotExist(err) {
		return nil, nil // No ignore file found, so return nil slice
	} else if err != nil {
		return nil, fmt.Errorf("error checking ignore file: %w", err)
	}

	// Open the ignore file
	file, err := os.Open(ignoreFilePath)
	if err != nil {
		return nil, fmt.Errorf("error opening ignore file: %w", err)
	}
	defer file.Close()

	var ignorePatterns []string
	scanner := bufio.NewScanner(file)

	// Read each line from the file
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())        // Trim whitespace
		if line != "" && !strings.HasPrefix(line, "#") { // Ignore empty lines and comments
			ignorePatterns = append(ignorePatterns, line)
		}
	}

	// Check for scanning errors
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading ignore file: %w", err)
	}

	return ignorePatterns, nil
}

func splitPath(path string) []string {
	subPath := path
	var result []string
	for {
		subPath = filepath.Clean(subPath) // Amongst others, removes trailing slashes (except for the root directory).

		dir, last := filepath.Split(subPath)
		if last == "" {
			if dir != "" { // Root directory.
				result = append(result, dir)
			}
			break
		}
		result = append(result, last)

		if dir == "" { // Nothing to split anymore.
			break
		}
		subPath = dir
	}

	slices.Reverse(result)
	return result
}

func isInIgnoreFile(relPath string, repoPath string) bool {
	ignorePatterns := getIgnorePatterns(repoPath)
	parts := splitPath(relPath)

	for _, pattern := range ignorePatterns {

		// Handle glob patterns (e.g., "*.log" or "dir/*")
		if strings.ContainsAny(pattern, "*?") {

			g := glob.MustCompile(pattern)
			if g.Match(relPath) {
				return true
			}
			continue
		}

		// Split pattern path to match components
		patternParts := splitPath(pattern)
		if len(parts) < len(patternParts) {
			// relPath is shorter than the pattern path, so it can't match
			continue
		}

		// Check if relPath starts with pattern path
		matches := true
		for i, part := range patternParts {
			if part != parts[i] {
				matches = false
				break
			}
		}
		if matches {
			return true
		}
	}

	return false
}
