package commands

import (
	"errors"
	"flag"
	"fmt"
	"git_clone/gvc/index"
	"git_clone/gvc/refs"
	"git_clone/gvc/utils"
	"os"

	"strings"
)

func RemoveFile(filePath string, cached, recursive, force bool) (string, error) {
	fileInfo, err := os.Stat(filePath)
	fileNotExists := os.IsNotExist(err)

	if err != nil && !fileNotExists {
		return "", fmt.Errorf("could not access %s: %w", filePath, err)
	}
	if !fileNotExists && fileInfo.IsDir() && !recursive {
		return "", fmt.Errorf("fatal: not removing '%s' recursively without -r", filePath)
	}

	tree, err := refs.GetLastCommitsTree()
	if err != nil {
		return "", err
	}

	relPath, err := utils.MakePathRelativeToRepo(utils.RepoDir, filePath)
	if err != nil {
		return "", fmt.Errorf("failed making path relative '%s': %w", filePath, err)
	}
	files := utils.MatchFileWithMapStringKey(relPath, tree)

	if len(files) == 0 {
		return "", fmt.Errorf("fatal: pathspec '%s' did not match any files", filePath)
	}

	anyMatched := false
	messages := make([]string, 0)

	for _, relpath := range files {
		err = index.RemoveFile(relpath, tree[relpath].FileHash, cached, force)
		var notExistError *index.FileNotPartOfIndexOrTreeError
		if errors.As(err, &notExistError) {
			continue
		}

		if err != nil {
			return "", err
		}
		anyMatched = true

		if !cached {
			if err := os.RemoveAll(utils.RelPathToAbs(relpath)); err != nil {
				return "", fmt.Errorf("can't delete '%s' from disk because %w", utils.RelPathToAbs(relpath), err)
			}
		}

		messages = append(messages, fmt.Sprintf("rm '%s'", utils.RelPathToAbs(relpath)))
	}

	if !anyMatched {
		return "", fmt.Errorf("fatal: pathspec '%s' did not match any files", filePath)
	}

	if !fileNotExists && !cached && fileInfo != nil && fileInfo.IsDir() {
		if err := os.RemoveAll(filePath); err != nil {
			return "", fmt.Errorf("fatal: error while deleting dir '%s' at the end %w", filePath, err)
		}
	}

	return strings.Join(messages, "\n"), nil
}

func RMCommand(args []string) string {
	rmCmd := flag.NewFlagSet("rm", flag.ExitOnError)
	rmChached := rmCmd.Bool("cached", false, "Only deletes file from .gvc not the actual file")
	rmRecursive := rmCmd.Bool("r", false, "")
	rmForce := rmCmd.Bool("f", false, "")

	rmCmd.Parse(os.Args[2:])
	if len(rmCmd.Args()) < 1 {
		return "Error: expected file paths to rm."
	}

	var output strings.Builder
	for _, filePath := range rmCmd.Args() {
		newOutput, err := RemoveFile(filePath, *rmChached, *rmRecursive, *rmForce)
		if err != nil {
			return fmt.Sprintf("error for file '%s': %v", filePath, err)
		}
		output.WriteString(newOutput + "\n")
	}
	return output.String()
}
