package commands

import (
	"flag"
	"fmt"
	"git_clone/gvc/commit"
	"git_clone/gvc/config"
	"git_clone/gvc/logging"
	"git_clone/gvc/settings"
)

func Commit(inputArgs []string) string {
	commitCmd := flag.NewFlagSet("commit", flag.ExitOnError)
	help := commitCmd.Bool("help", false, "Get Help Documentation")
	helpShort := commitCmd.Bool("h", false, "Get Help Documentation")
	commitMessage := commitCmd.String("m", "", "The commit message")
	commitUser := commitCmd.String("u", "", "The commit user")

	if err := commitCmd.Parse(inputArgs); err != nil {
		return fmt.Errorf("error parsing arguments: %w", err).Error()
	}
	args := commitCmd.Args()
	if len(args) > 0 {
		return logging.NewError("gvc commit expects no positional arguments").Error()
	}
	if *help || *helpShort {
		return "gvc commit -m <message> -u <user>\nCommits the staged changes with the provided commit message and user.\nIf no user is provided the user set by 'gvc set' is used\nFails if any positional arguments are supplied."
	}

	if *commitMessage == "" {
		return "Error: commit message (-m) is required."
	}

	cfg, err := settings.LoadSettings()
	if err != nil {
		return err.Error()
	}

	if cfg.User == config.DOES_NOT_EXIST_HASH && *commitUser == "" {
		return "Error: commit user (-u) is required, since not set in settings\nUse gvc set --set User=username."
	}

	if *commitUser != "" {
		cfg.User = *commitUser
	}

	output, err := commit.Commit(*commitMessage, cfg.User, true)

	if err != nil {
		return fmt.Errorf("commit failed: \n    %w", err).Error()
	}
	return output
}
