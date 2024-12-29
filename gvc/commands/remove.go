package commands

import (
	"errors"
	"fmt"
	"git_clone/gvc/index"
	"git_clone/gvc/utils"
	"os"
)

func RemoveFile(filePath string, cached, recursive, force bool) ([]string, error) {

	if isdir, err := utils.IsDir(filePath); err != nil {
		return nil, err
	} else if isdir && !recursive {
		return nil, fmt.Errorf("fatal: not removing '%s' recursively without -r", filePath)
	}

	files, err := utils.FindMatchingFiles(filePath)
	if err != nil {
		return nil, fmt.Errorf("error while finding files: %s", err)
	}

	onlyFilesThatWereNotMatched := true

	messages := make([]string, 0)

	for _, file := range files {

		relpath, err := utils.MakePathRelativeToRepo(utils.RepoDir, file)
		if err != nil {
			return nil, err
		}

		fileHash, err := utils.GetFileSHA1(file)
		if err != nil {
			return nil, fmt.Errorf("error while hashing %s: %w", file, err)
		}

		// if there is something that is untracked but in the same folder and we delete a folder we can ignore/ if only -> return error
		err = index.RemoveFileFromIndex(relpath, fileHash, cached, force)
		var notExistError *index.FileNotPartOfIndexOrTreeError
		if errors.As(err, &notExistError) {
			continue
		}
		onlyFilesThatWereNotMatched = false

		if err != nil {
			messages = append(messages, fmt.Sprintf("could not rm file %s: %v", file, err))
			continue
		}
		messages = append(messages, fmt.Sprintf("rm '%s'", file))
	}

	if onlyFilesThatWereNotMatched {
		return nil, fmt.Errorf("fatal: pathspec '%s' did not match any files", filePath)
	}

	if cached {
		return messages, nil
	}

	for _, file := range files {
		if err := os.RemoveAll(file); err != nil {
			messages = append(messages, fmt.Sprintf("can't delete '%s' from disk because %v", file, err))
		}
	}
	return nil, err
}
