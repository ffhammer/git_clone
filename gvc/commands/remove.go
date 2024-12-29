package commands

import (
	"errors"
	"fmt"
	"git_clone/gvc/index"
	"git_clone/gvc/utils"
	"os"
	"strings"
)

func RemoveFile(filePath string, cached, recursive, force bool) (string, error) {
	isdir, err := utils.IsDir(filePath)
	if err != nil {
		return "", err
	} else if isdir && !recursive {
		return "", fmt.Errorf("fatal: not removing '%s' recursively without -r", filePath)
	}

	files, err := utils.FindMatchingFiles(filePath)
	if err != nil {
		return "", fmt.Errorf("error while finding files: %s", err)
	}

	onlyFilesThatWereNotMatched := true
	allMatches := true

	messages := make([]string, 0)

	for _, file := range files {

		relpath, err := utils.MakePathRelativeToRepo(utils.RepoDir, file)
		if err != nil {
			return "", err
		}

		fileHash, err := utils.GetFileSHA1(file)
		if err != nil {
			return "", fmt.Errorf("error while hashing %s: %w", file, err)
		}

		// if there is something that is untracked but in the same folder and we delete a folder we can ignore/ if only -> return error
		err = index.RemoveFileFromIndex(relpath, fileHash, cached, force)
		var notExistError *index.FileNotPartOfIndexOrTreeError
		if errors.As(err, &notExistError) {
			allMatches = false
			continue
		}
		onlyFilesThatWereNotMatched = false

		if err != nil {
			return "", err
		}

		if !cached {
			if err := os.RemoveAll(file); err != nil {
				return "", fmt.Errorf("can't delete '%s' from disk because %w", file, err)
			}
		}

		messages = append(messages, fmt.Sprintf("rm '%s'", file))
	}

	if onlyFilesThatWereNotMatched {
		return "", fmt.Errorf("fatal: pathspec '%s' did not match any files", filePath)
	}

	if isdir && !cached && allMatches {
		if err := os.RemoveAll(filePath); err != nil {
			return "", fmt.Errorf("fatal: error while deleting dir '%s' at the end %w", filePath, err)
		}
	}

	return strings.Join(messages, "\n"), nil
}
