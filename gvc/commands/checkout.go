package commands

import (
	"flag"
	"fmt"
	"git_clone/gvc/index"
	"git_clone/gvc/objectio"
	"git_clone/gvc/refs"
	"git_clone/gvc/utils"
	"os"
	"strings"
)

func checkForConflicts(changes index.ChangeList) error {

	uncomitted, err := index.GetUncommitedChanges()
	if err != nil {
		return err
	}

	unstaged, err := index.GetUnstagedChanges(true)
	if err != nil {
		return err
	}

	setForEffiency := make(map[string]bool, len(uncomitted))
	for _, i := range uncomitted {
		setForEffiency[i.RelPath] = false
	}
	for _, i := range unstaged {
		setForEffiency[i.RelPath] = false
	}

	var builder strings.Builder

	return_error := false
	for _, change := range changes {

		if _, ok := setForEffiency[change.RelPath]; ok {
			return_error = true
			builder.WriteString(fmt.Sprintf("\t\t-%s", change.RelPath))
		}
	}

	if return_error {
		return fmt.Errorf("error: Your local changes to the following files would be overwritten by checkout:\n%s", builder.String())
	} else {
		builder.Reset()
	}

	return nil
}

func checkout(branchName string) error {
	// main checkout command that does the heavy lifting
	origianalTree, err := refs.LoadBranchTree(branchName)
	if err != nil {
		return fmt.Errorf("could not load branch tree: %s", err)
	}

	currentTree, err := refs.GetLastCommitsTree()
	if err != nil {
		return err // already saved contained error message
	}

	// with this ordering a "deletion" would also entail a deletion
	changes := index.TreeDiff(currentTree, origianalTree, true)
	if err := checkForConflicts(changes); err != nil {
		return err
	}

	for _, change := range changes {

		absPath := utils.RelPathToAbs(change.RelPath)

		if change.Action == index.Delete {
			os.Remove(absPath)
			continue
		}

		// case modify/or add. either way -> write file
		newFile, err := objectio.RetrieveFile(change.NewHash)
		if err != nil {
			return fmt.Errorf("error retriving file '%s': %w", absPath, err)
		}

		if err := os.WriteFile(absPath, []byte(newFile), os.ModePerm); err != nil {
			return fmt.Errorf("error writing file '%s': %w", absPath, err)
		}

	}

	return nil
}

func CheckoutCommand(inputArgs []string) string {

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

	if err := checkout(branchName); err != nil {
		return err.Error()
	}

	if err := refs.UpdateHead(branchName); err != nil {
		return err.Error()
	}

	return ""
}
