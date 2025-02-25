package commands

import (
	"flag"
	"fmt"
	"git_clone/gvc/logging"
	"git_clone/gvc/merge"
	"git_clone/gvc/refs"
	"git_clone/gvc/settings"
)

func MergeCommand(inputArgs []string) string {
	flagset := flag.NewFlagSet("merge", flag.ExitOnError)

	if err := flagset.Parse(inputArgs); err != nil {
		logging.Log(fmt.Sprintf("Error parsing merge command arguments: %v", err), settings.ERROR)
		return err.Error()
	}

	args := flagset.Args()

	if len(args) == 0 {
		logging.Log("Merge command missing branch argument", settings.WARNING)
		return "Need to specify a branch to merge to."
	} else if len(args) > 1 {
		logging.Log("Merge command received multiple branch arguments", settings.WARNING)
		return "Pass only one branch."
	}

	currentBranch, err := refs.LoadCurrentBranchName()
	if err != nil {
		logging.Log(fmt.Sprintf("Failed to load current branch: %v", err), settings.ERROR)
		return err.Error()
	}

	branchName := args[0]

	if branchName == currentBranch {
		logging.Log(fmt.Sprintf("Attempted to merge branch '%s' into itself", branchName), settings.WARNING)
		return fmt.Sprintf("You are already on '%s'", branchName)
	}

	if exist, err := refs.BranchExists(branchName); err != nil {
		logging.Log(fmt.Sprintf("Error checking if branch '%s' exists: %v", branchName, err), settings.ERROR)
		return err.Error()
	} else if !exist {
		logging.Log(fmt.Sprintf("Merge failed: Branch '%s' does not exist", branchName), settings.WARNING)
		return fmt.Sprintf("Branch %s does not exist. To create it use -b", branchName)
	}

	logging.Log(fmt.Sprintf("Starting merge: Merging '%s' into '%s'", branchName, currentBranch), settings.INFO)
	if err := merge.Merge(currentBranch, branchName); err != nil {
		logging.Log(fmt.Sprintf("Merge failed: %v", err), settings.ERROR)
		return err.Error()
	}

	logging.Log(fmt.Sprintf("Merge of '%s' into '%s' completed successfully", branchName, currentBranch), settings.INFO)
	return "Merge completed successfully."
}
