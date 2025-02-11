package commands

import (
	"flag"
	"fmt"
	"git_clone/gvc/merge"
	"git_clone/gvc/refs"
)

func MergeCommand(inputArgs []string) string {

	flagset := flag.NewFlagSet("checkout", flag.ExitOnError)

	if err := flagset.Parse(inputArgs); err != nil {
		return err.Error()
	}

	args := flagset.Args()

	if len(args) == 0 {
		return "need to specify branch to merge to"
	} else if len(args) > 1 {
		return "pass only one branch"
	}

	currentBranch, err := refs.LoadCurrentBranchName()
	if err != nil {
		return err.Error()
	}

	branchName := args[0]

	if branchName == currentBranch {
		return fmt.Sprintf("You are already on '%s'", branchName)
	}

	if exist, err := refs.BranchExists(branchName); err != nil {
		return err.Error()
	} else if !exist {
		return fmt.Errorf("branch %s does not exist. to create it use -b", branchName).Error()
	}

	if err := merge.Merge(currentBranch, branchName); err != nil {
		return err.Error()
	}

	return ""
}
