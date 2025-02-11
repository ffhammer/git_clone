package commands

import (
	"flag"
	"fmt"
	"git_clone/gvc/refs"
	"git_clone/gvc/switching"
)

func CheckoutCommand(inputArgs []string) string {

	flagset := flag.NewFlagSet("checkout", flag.ExitOnError)
	bFlag := flagset.Bool("b", false, "create new branch")

	if err := flagset.Parse(inputArgs); err != nil {
		return err.Error()
	}

	args := flagset.Args()

	if len(args) == 0 {
		return "need to specify branch"
	} else if len(args) > 1 {
		return "parse only one branch"
	}

	currentBranch, err := refs.LoadCurrentBranchName()
	if err != nil {
		return err.Error()
	}

	branchName := args[0]

	if branchName == currentBranch {
		return fmt.Sprintf("You are already on '%s'", branchName)
	}

	if *bFlag {
		if err := refs.CreateNewBranch(branchName); err != nil {
			return err.Error()
		}
	}

	if exist, err := refs.BranchExists(branchName); err != nil {
		return err.Error()
	} else if !exist {
		return fmt.Errorf("branch %s does not exist. to create it use -b", branchName).Error()
	}

	if err := switching.UpdateWorkingDirToBranch(branchName, "checkout"); err != nil {
		return err.Error()
	}

	if err := refs.UpdateHead(branchName); err != nil {
		return err.Error()
	}

	return ""
}
