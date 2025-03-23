package commands

import (
	"flag"
	"fmt"
	"git_clone/gvc/commit"
	"git_clone/gvc/config"
	"git_clone/gvc/settings"
	"os"
)

func Commit() string {
	commitCmd := flag.NewFlagSet("commit", flag.ExitOnError)
	commitMessage := commitCmd.String("m", "", "The commit message")
	commitUser := commitCmd.String("u", "", "The commit user")
	commitCmd.Parse(os.Args[2:])

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
