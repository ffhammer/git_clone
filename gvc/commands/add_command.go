package commands

import (
	"flag"
	"fmt"
	"git_clone/gvc/config"
	"git_clone/gvc/ignorefiles"
	"git_clone/gvc/index"
	"git_clone/gvc/objectio"
	"git_clone/gvc/utils"
	"os"
	"path/filepath"
	"strings"
)

func addSingleFile(filePath string, force bool) error {
	next_commit := filepath.Join(utils.RepoDir, config.NEXT_COMMIT)

	utils.MkdirIgnoreExists(next_commit)

	if _, err := os.Stat(filePath); err != nil {
		return fmt.Errorf("can't find file %s: %w", filePath, err)
	}

	relPath, err := utils.MakePathRelativeToRepo(utils.RepoDir, filePath)
	if err != nil {
		return err
	}

	if ignorefiles.IsIgnored(relPath) && !force {
		return fmt.Errorf("file %s is ignored. Use add -f to force it", filePath)
	}

	fileHash, err := utils.GetFileSHA1(filePath)
	if err != nil {
		return fmt.Errorf("can't add file %s to objects because %w", filePath, err)

	}

	err = objectio.AddFileToObjects(filePath, fileHash)
	if err != nil {
		return fmt.Errorf("can't add file %s to objects because %w", filePath, err)
	}

	if err := index.AddFile(relPath, fileHash); err != nil {
		return err
	}

	return nil
}

func AddFiles(filePath string, force bool) string {

	files, err := utils.FindMatchingFiles(filePath)
	if err != nil {
		return fmt.Errorf("fatal: could not match files : %w", err).Error()
	}

	messages := make([]string, 0)
	// Add each file using addSingleFile
	for _, file := range files {
		err := addSingleFile(file, force)
		if err != nil {
			messages = append(messages, fmt.Sprintf("could not add file %s: %v", file, err))
		}
	}

	return strings.Join(messages, "\n")
}

func AddCommand(args []string) string {
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	force := addCmd.Bool("f", false, "Force adding the file even if it is ignored")
	addCmd.Parse(args)
	if len(addCmd.Args()) < 1 {
		fmt.Println("Error: expected file paths to add.")
		addCmd.Usage()
		os.Exit(1)
	}
	output := ""

	for _, filePath := range addCmd.Args() {

		if res := AddFiles(filePath, *force); res != "" {
			return res
		}
	}
	return output
}
