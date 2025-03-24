package commands

import (
	"flag"
	"fmt"
	"git_clone/gvc/logging"
	"git_clone/gvc/status"
)

func Status(inputArgs []string) string {
	flagSet := flag.NewFlagSet("status", flag.ExitOnError)
	help := flagSet.Bool("help", false, "Get help documentation")
	helpShort := flagSet.Bool("h", false, "Get help documentation")
	if err := flagSet.Parse(inputArgs); err != nil {
		return fmt.Errorf("error parsing arguments: %w", err).Error()
	}
	if len(flagSet.Args()) > 0 {
		return logging.NewError("gvc status expects no arguments").Error()
	}
	if *help || *helpShort {
		return "gvc status\n" +
			"Show the current repository status, similar to 'git status'.\n" +
			"Displays the current branch, staged changes, modifications, untracked files, and merge state."
	}

	output, err := status.Status()
	if err != nil {
		return fmt.Errorf("status failed because: %w", err).Error()
	}
	return output
}
