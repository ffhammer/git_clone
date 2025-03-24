package commands

import (
	"flag"
	"fmt"
	"git_clone/gvc/refs"
	"git_clone/gvc/switching"
)

func CheckoutCommand(inputArgs []string) string {

	flagset := flag.NewFlagSet("checkout", flag.ExitOnError)
	help := flagset.Bool("help", false, "Get help documentation")
	helpShort := flagset.Bool("h", false, "Get help documentation")
	bFlag := flagset.Bool("b", false, "Create a new branch before switching to it")

	if err := flagset.Parse(inputArgs); err != nil {
		return err.Error()
	}
	if *help || *helpShort {
		return "gvc checkout [options] <branch>\n" +
			"Switch branches or restore working tree.\n\n" +
			"Options:\n" +
			"  -b        Create the branch before switching to it"
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

	if refs.InMergeState {
		return "Error: can't use checkout in open merge state"
	}

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
