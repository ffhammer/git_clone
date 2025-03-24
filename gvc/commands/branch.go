package commands

import (
	"flag"
	"fmt"
	"git_clone/gvc/refs"
)

func BranchCommand(inputArgs []string) string {
	flagset := flag.NewFlagSet("branch", flag.ExitOnError)
	help := flagset.Bool("help", false, "Get help documentation")
	helpShort := flagset.Bool("h", false, "Get help documentation")
	delFlag := flagset.Bool("d", false, "Delete a branch")

	if err := flagset.Parse(inputArgs); err != nil {
		return err.Error()
	}
	if *help || *helpShort {
		return "gvc branch [options] [<branch-name>...]\n" +
			"Create, list, or delete branches.\n\n" +
			"Options:\n" +
			"  -d        Delete the specified branch(es)\n" +
			"  --help    Show this message\n\n" +
			"Usage:\n" +
			"  gvc branch           List all branches\n" +
			"  gvc branch <name>    Create a new branch\n" +
			"  gvc branch -d <name> Delete branch"
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
