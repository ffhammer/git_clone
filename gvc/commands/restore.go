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
	// restoreSource := restoreCmd.String("source", "", "The branch or commit")
	restoreStaged := restoreCmd.Bool("staged", false, "")
	restoreWorktree := restoreCmd.Bool("worktree", false, "")

	restoreCmd.Parse(args)

	if len(restoreCmd.Args()) < 1 {
		fmt.Println("Error: expected file paths to restore.")
		restoreCmd.Usage()
		os.Exit(1)
	}

	if refs.InMergeState && !(*restoreStaged) {
		return logging.NewError("in merge staged you can only restore with --staged").Error()
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
