package commands

import (
	"flag"
	"fmt"
	"git_clone/gvc/refs"
)

func Checkout(inputArgs []string) string {

	flagset := flag.NewFlagSet("checkout", flag.ExitOnError)
	bFlag := flag.Bool("b", false, "create new branch")

	if err := flagset.Parse(inputArgs); err != nil {
		return err.Error()
	}

	args := flagset.Args()

	if len(args) == 0 {
		return "need to specify branch"
	} else if len(args) > 1 {
		return "parse only one branch"
	}

	branchName := args[0]
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

	return ""

}
