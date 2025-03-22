package commands

import (
	"flag"
	"fmt"
	"git_clone/gvc/logging"
	"git_clone/gvc/merge"
	"git_clone/gvc/refs"
)

func MergeCommand(inputArgs []string) string {
	flagset := flag.NewFlagSet("merge", flag.ExitOnError)

	if err := flagset.Parse(inputArgs); err != nil {
		return logging.ErrorF("Error parsing merge command arguments: %v", err).Error()
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

	return merge.Merge(currentBranch, branchName)
}
