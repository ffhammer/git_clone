package commands

import (
	"flag"
	"fmt"
	"git_clone/gvc/refs"
)

func BranchCommand(inputArgs []string) string {

	flagset := flag.NewFlagSet("branch", flag.ExitOnError)
	delFlag := flagset.Bool("d", false, "wether to delete a flag")

	if err := flagset.Parse(inputArgs); err != nil {
		return err.Error()
	}

	args := flagset.Args()

	if *delFlag && len(args) == 0 {
		return "fatal: branch name required"
	} else if *delFlag {
		for _, file := range args {
			if err := refs.DeleteBranch(file); err != nil {
				return err.Error()
			}

		}
		return ""
	}

	if len(args) == 0 {
		output, err := refs.ListBranches()
		if err != nil {
			return fmt.Errorf("error listing branches: %w", err).Error()
		}
		return output
	}

	for _, file := range args {
		if err := refs.CreateNewBranch(file); err != nil {
			return err.Error()
		}

	}

	return ""

}
