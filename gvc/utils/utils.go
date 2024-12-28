package utils

import (
	"bufio"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func MkdirIgnoreExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.Mkdir(path, os.ModePerm)
	}
	return nil
}

func SplitPath(path string) []string {
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

func MakePathRelativeToRepo(RepoDIr string, filePath string) (string, error) {

	absRepoDir, err := filepath.Abs(filepath.Dir(RepoDIr))
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

func GetFileSHA1(filePath string) (string, error) {
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

func GetStringSHA1(s string) string {
	h := sha1.New()

	h.Write([]byte(s))

	return fmt.Sprintf("%x", h.Sum(nil))
}

func reader2String(reader io.Reader) (string, error) {
	buf := new(strings.Builder)
	_, err := io.Copy(buf, reader)
	return buf.String(), err
}

func SplitLines(s string) []string {
	var lines []string
	sc := bufio.NewScanner(strings.NewReader(s))
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	return lines
}
