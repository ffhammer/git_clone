package switching

import (
	"fmt"
	"git_clone/gvc/objectio"
	"git_clone/gvc/refs"
	"git_clone/gvc/treediff"
	"git_clone/gvc/utils"
	"os"
)

func UpdateWorkingDirToBranch(branchName, operationName string) error {
	origianalTree, err := refs.LoadBranchTree(branchName)
	if err != nil {
		return fmt.Errorf("could not load branch tree: %s", err)
	}

	currentTree, err := refs.GetLastCommitsTree()
	if err != nil {
		return err // already saved contained error message
	}

	// with this ordering a "deletion" would also entail a deletion
	var changes treediff.ChangeList = treediff.ChangeList{}
	treediff.TreeDiff[treediff.ChangeList](&changes, currentTree, origianalTree, false)

	relPathsOfChanges := make([]string, len(changes))
	for i, change := range changes {
		relPathsOfChanges[i] = change.RelPath
	}

	if err := FindNotChangeableFiles(relPathsOfChanges); err != nil {
		return fmt.Errorf("error: Your local changes to the following files would be overwritten by %s:\n%w", operationName, err)
	}

	for _, change := range changes {

		absPath := utils.RelPathToAbs(change.RelPath)

		if change.Action == treediff.Delete {
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
