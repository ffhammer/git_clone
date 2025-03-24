package merge

import (
	"fmt"
	"git_clone/gvc/commit"
	"git_clone/gvc/config"
	"git_clone/gvc/index"
	"git_clone/gvc/logging"
	"git_clone/gvc/objectio"
	"git_clone/gvc/refs"
	"git_clone/gvc/switching"
	"git_clone/gvc/treediff"
	"git_clone/gvc/utils"
	"os"
	"path/filepath"
)

// fastForwardMerge moves the branch pointer forward when no conflicts exist.
func fastForwardMerge(currentBranch, targetBranchName, userName string) (string, error) {
	logging.Info(fmt.Sprintf("Performing fast-forward merge to '%s'", targetBranchName))

	if err := switching.UpdateWorkingDirToBranch(targetBranchName, "merge"); err != nil {
		return "", logging.Error(err)
	}

	pathToCurrentPointer := filepath.Join(utils.RepoDir, config.CurrentBranchPointerFile)
	if err := os.WriteFile(pathToCurrentPointer, []byte(targetBranchName), os.ModePerm); err != nil {
		return "", logging.ErrorF("error updating HEAD pointer %w", err)
	}

	commitOutput, err := commit.Commit(fmt.Sprintf("Fastforward Merge branch '%s' into %s\n", targetBranchName, currentBranch), userName, false)
	logging.InfoF("commitOutput:\n%s", commitOutput)
	if err != nil {
		return "", err
	}

	logging.InfoF("Fast-forward merge to '%s' completed successfully", targetBranchName)
	return fmt.Sprintf("fast-forward merge to '%s' completed successfully\n%s", targetBranchName, commitOutput), nil
}

// Merge performs a three-way merge between the current branch and source branch.
func Merge(currentBranch, sourceBranchName, userName string) string {

	if refs.InMergeState {
		return "Error: can't use merge in open merge state"
	}

	logging.InfoF("Starting merge process from '%s' into '%s'", sourceBranchName, currentBranch)

	// Check for uncommitted changes
	if currentChanges, err := index.LoadIndexChanges(); err != nil {
		return logging.Error(err).Error()
	} else if len(currentChanges) > 0 {
		return logging.NewError("error: please commit your current changes before merging").Error()

	}
	currentHash, err := refs.GetBranchCommitHash(currentBranch)
	if err != nil {
		return logging.ErrorF("error: can't load commit hash for branch '%s'. error: %w", currentBranch, err).Error()
	}

	sourceBranchHash, err := refs.GetBranchCommitHash(sourceBranchName)
	if err != nil {
		return logging.ErrorF("error: can't load commit hash for branch '%s'. error: %w", sourceBranchName, err).Error()
	}

	mergeBase, err := findMergeBaseHash(currentHash, sourceBranchHash)
	if err != nil {
		return logging.Error(err).Error()
	}

	logging.DebugF("common merge base: '%s'", mergeBase)

	if mergeBase == config.DOES_NOT_EXIST_HASH {
		return logging.NewError("error: can't merge because no common parent commit was detected").Error()
	}

	if currentHash == mergeBase {
		logging.InfoF("Performing fast-forward merge from '%s' to '%s'", sourceBranchName, currentBranch)
		returnMessage, err := fastForwardMerge(currentBranch, sourceBranchName, userName)
		if err != nil {
			return fmt.Errorf("error: fast forward merge failed with: %w", err).Error()
		}
		return returnMessage
	}

	logging.Info("Performing three-way merge")

	baseTree, err := refs.LoadCommitTreeHeadAccpeted(mergeBase)
	if err != nil {
		return logging.ErrorF("can't retrieve merge base tree: %w", err).Error()
	}

	currentTree, err := objectio.LoadTreeByCommitHash(currentHash)
	if err != nil {
		return logging.ErrorF("error: can't retrieve last committed tree: %w", err).Error()
	}

	mergeFromTree, err := objectio.LoadTreeByCommitHash(sourceBranchHash)
	if err != nil {
		return logging.ErrorF("error: can't retrieve tree from which we will merge: %w", err).Error()
	}

	changesA := treediff.ChangeMap{}
	treediff.TreeDiff[treediff.ChangeMap](changesA, baseTree, currentTree, false)
	changesB := treediff.ChangeMap{}
	treediff.TreeDiff[treediff.ChangeMap](changesB, baseTree, mergeFromTree, false)

	mergeConflicts := findMergeConflicts(changesA, changesB)
	conflictRelPaths := make([]string, len(mergeConflicts))
	for index, conflict := range mergeConflicts {
		conflictRelPaths[index] = conflict.RelPath
	}

	if err := switching.FindNotChangeableFiles(conflictRelPaths); err != nil {
		return logging.ErrorF("error: Your local changes to the following files would be overwritten by merge:\n%w", err).Error()
	}

	if err := addFileAdditionsOfNewBranch(changesA, changesB); err != nil {
		return err.Error()
	}
	if err := prepareMergeState(mergeConflicts, currentBranch, sourceBranchHash, sourceBranchName); err != nil {
		return logging.ErrorF("error: preparing merge state :\n%w", err).Error()
	}

	return ""
}
