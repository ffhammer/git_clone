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
	help := rmCmd.Bool("help", false, "Get help documentation")
	helpShort := rmCmd.Bool("h", false, "Get help documentation")
	rmCached := rmCmd.Bool("cached", false, "Only remove from index, not filesystem")
	rmRecursive := rmCmd.Bool("r", false, "Recursively remove files in directories")
	rmForce := rmCmd.Bool("f", false, "Force removal even if file is staged")

	if err := rmCmd.Parse(args); err != nil {
		return fmt.Errorf("error parsing arguments: %w", err).Error()
	}
	if *help || *helpShort {
		return "gvc rm [options] <path>...\n" +
			"Removes file(s) from the working tree and/or the index.\n\n" +
			"Options:\n" +
			"  --cached   Only remove from index, keep file in working dir\n" +
			"  -f         Force removal, even if file is staged or ignored\n" +
			"  -r         Allow recursive removal of directories"
	}

	rmCmd.Parse(args)
	if len(rmCmd.Args()) < 1 {
		return "Error: expected file paths to rm."
	}

	if refs.InMergeState {
		return "Error: can't rm in open merge state"
	}

	var output strings.Builder
	for _, filePath := range rmCmd.Args() {
		newOutput, err := RemoveFile(filePath, *rmCached, *rmRecursive, *rmForce)
		if err != nil {
			return fmt.Sprintf("error for file '%s': %v", filePath, err)
		}
		output.WriteString(newOutput + "\n")
	}
	return output.String()
}
