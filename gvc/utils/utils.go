package utils

import (
	"bufio"
	"crypto/sha1"
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

func IsDir(filePath string) (bool, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return false, fmt.Errorf("could not access %s: %w", filePath, err)
	}
	return fileInfo.IsDir(), nil
}
