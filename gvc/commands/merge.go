package commands

import (
	"flag"
	"fmt"
	"git_clone/gvc/config"
	"git_clone/gvc/logging"
	"git_clone/gvc/merge"
	"git_clone/gvc/refs"
	"git_clone/gvc/settings"
)

func MergeCommand(inputArgs []string) string {
	flagset := flag.NewFlagSet("merge", flag.ExitOnError)
	help := flagset.Bool("help", false, "Get help documentation")
	helpShort := flagset.Bool("h", false, "Get help documentation")
	commitUser := flagset.String("u", "", "The commit user")

	if err := flagset.Parse(inputArgs); err != nil {
		return logging.ErrorF("Error parsing merge command arguments: %v", err).Error()
	}
	if *help || *helpShort {
		return "gvc merge [options] <branch>\n" +
			"Merge the specified branch into the current branch.\n\n" +
			"Options:\n" +
			"  -u <user>   Specify commit user (overrides config setting)"
	}

	cfg, err := settings.LoadSettings()
	if err != nil {
		return err.Error()
	}
	if cfg.User == config.DOES_NOT_EXIST_HASH && *commitUser == "" {
		return "Error: user (-u) is required, since not set in settings\nUse gvc set --set User=username."
	}
	if *commitUser != "" {
		cfg.User = *commitUser
	}

	args := flagset.Args()
	if len(args) == 0 {
		logging.Warn("Merge command missing branch argument")
		return "Need to specify a branch to merge to."
	} else if len(args) > 1 {
		logging.Warn("Merge command received multiple branch arguments")
		return "Pass only one branch."
	}

	currentBranch, err := refs.LoadCurrentBranchName()
	if err != nil {
		return logging.ErrorF("Failed to load current branch: %v", err).Error()
	}

	branchName := args[0]

	if branchName == currentBranch {
		logging.WarnF("Attempted to merge branch '%s' into itself", branchName)
		return fmt.Sprintf("You are already on '%s'", branchName)
	}

	if exist, err := refs.BranchExists(branchName); err != nil {
		return logging.ErrorF("Error checking if branch '%s' exists: %v", branchName, err).Error()
	} else if !exist {
		logging.WarnF("Merge failed: Branch '%s' does not exist", branchName)
		return fmt.Sprintf("Branch %s does not exist. To create it use -b", branchName)
	}

	logging.InfoF("Starting merge: Merging '%s' into '%s'", branchName, currentBranch)

	return merge.Merge(currentBranch, branchName, cfg.User)
}
