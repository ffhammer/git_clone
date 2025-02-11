package merge

import (
	"errors"
	"fmt"
	"git_clone/gvc/config"
	"git_clone/gvc/objectio"
	"git_clone/gvc/refs"
	"git_clone/gvc/switching"
	"git_clone/gvc/treediff"
	"git_clone/gvc/utils"
	"os"
	"path/filepath"
)

func fastForwardMerge(targetBranchName, targetBranchHash string) error {

	if err := switching.UpdateWorkingDirToBranch(targetBranchName, "merge"); err != nil {
		return err
	}

	pathToCurrentPointer := filepath.Join(utils.RepoDir, config.CurrentBranchPointerFile)
	if err := os.WriteFile(pathToCurrentPointer, []byte(targetBranchHash), os.ModePerm); err != nil {
		return fmt.Errorf("error updating HEAD pointer %w", err)
	}
	return nil

}

func Merge(currentBranch, sourceBranchName string) error {
	currentHash, err := refs.GetBranchCommitHash(currentBranch)
	if err != nil {
		return fmt.Errorf("cant load commit hash for branch '%s'. error: %w", currentBranch, err)
	}

	sourceBranchHash, err := refs.GetBranchCommitHash(sourceBranchName)
	if err != nil {
		return fmt.Errorf("cant load commit hash for branch '%s'. error: %w", sourceBranchName, err)
	}

	mergeBase, err := findMergeBaseHash(currentHash, sourceBranchHash)
	if err != nil {
		return err
	}

	if mergeBase == config.DOES_NOT_EXIST_HASH {
		return errors.New("cant merge becase we did not detect common parent commit")
	}

	if currentHash == mergeBase {
		return fastForwardMerge(sourceBranchName, sourceBranchHash)
	}

	baseTree, err := objectio.LoadTreeByCommitHash(mergeBase)
	if err != nil {
		return fmt.Errorf("cant retrieve merge base tree: %w", err)
	}

	currentTree, err := objectio.LoadTreeByCommitHash(currentHash)
	if err != nil {
		return fmt.Errorf("cant retrieve last committed tree: %w", err)
	}

	mergeFromTree, err := objectio.LoadTreeByCommitHash(sourceBranchHash)
	if err != nil {
		return fmt.Errorf("cant retrieve tree from which we will merge: %w", err)
	}

	changesA := treediff.ChangeMap{}
	treediff.TreeDiff[treediff.ChangeMap](changesA, baseTree, currentTree, false)
	changesB := treediff.ChangeMap{}
	treediff.TreeDiff[treediff.ChangeMap](changesB, baseTree, mergeFromTree, false)

	mergeConflict := findMergeConflicts(changesA, changesB)
	conflictRelPaths := make([]string, len(mergeConflict))
	for index, conflict := range mergeConflict {
		conflictRelPaths[index] = conflict.RelPath
	}

	if err := switching.FindNotChangeableFiles(conflictRelPaths); err != nil {
		return fmt.Errorf("error: Your local changes to the following files would be overwritten by merge:\n%w", err)
	}

	return nil
}
