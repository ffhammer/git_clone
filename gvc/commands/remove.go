package commands

import (
	"errors"
	"flag"
	"fmt"
	"git_clone/gvc/index"
	"git_clone/gvc/utils"
	"os"
	"strings"
)

func RemoveFile(filePath string, cached, recursive, force bool) (string, error) {
	// ich muss die pfade matchen mit dem tree wenn cached, also schon die gelöschten
	fileInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) && cached {
	} else if err != nil {
		return "", fmt.Errorf("could not access %s: %w", filePath, err)
	} else if fileInfo.IsDir() && !recursive {
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
		err = index.RemoveFile(relpath, fileHash, cached, force)
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

	if !cached && allMatches && fileInfo.IsDir() {
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
		fmt.Println("Error: expected file paths to rm.")
		rmCmd.Usage()
		os.Exit(1)
	}

	output := ""
	for _, filePath := range rmCmd.Args() {
		newOutput, err := RemoveFile(filePath, *rmChached, *rmRecursive, *rmForce)

		if err != nil {
			output += fmt.Errorf("error for file '%s': %w", filePath, err).Error()
			return output
		}

		output += newOutput + "\n"
	}
	return output

}
