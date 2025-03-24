package commands

import (
	"flag"
	"fmt"
	"git_clone/gvc/logging"
	"git_clone/gvc/refs"
	"git_clone/gvc/restore"

	"os"
)

func Restore(args []string) string {
	restoreCmd := flag.NewFlagSet("restore", flag.ExitOnError)
	help := restoreCmd.Bool("help", false, "Get help documentation")
	helpShort := restoreCmd.Bool("h", false, "Get help documentation")
	restoreStaged := restoreCmd.Bool("staged", false, "Restore the index (staging area) content")
	restoreWorktree := restoreCmd.Bool("worktree", false, "Restore the working directory content")

	if err := restoreCmd.Parse(args); err != nil {
		return fmt.Errorf("error parsing arguments: %w", err).Error()
	}

	if *help || *helpShort {
		return "gvc restore [options] <path>\n" +
			"Restores files to their state in HEAD.\n\n" +
			"Options:\n" +
			"  --staged     Restore the index (remove from staging)\n" +
			"  --worktree   Restore the file content in the working directory\n\n" +
			"If neither --staged nor --worktree is given, neither are restored."
	}
	restoreCmd.Parse(args)

	if len(restoreCmd.Args()) < 1 {
		fmt.Println("Error: expected file paths to restore.")
		restoreCmd.Usage()
		os.Exit(1)
	}

	if refs.InMergeState && !(*restoreStaged) {
		return logging.NewError("in merge staged you can only restore with --staged").Error()
	}

	if !*restoreStaged && !*restoreWorktree {
		*restoreStaged = true
		*restoreWorktree = true
	}

	for _, filePath := range restoreCmd.Args() {

		var err error
		if refs.InMergeState {
			err = restore.InMergeRestore(filePath)

		} else {
			err = restore.StandardRestore(filePath, "HEAD", *restoreStaged, *restoreWorktree)
		}

		if err != nil {
			return logging.ErrorF("restore failed because: %w", err).Error()
		}

	}

	return ""
}
