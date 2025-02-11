package commands

import (
	"flag"
	"fmt"
	"git_clone/gvc/config"
	"git_clone/gvc/objectio"
	"git_clone/gvc/refs"
	"git_clone/gvc/treediff"
)

func getParentCommits(branchHash string) ([]string, error) {

	current, err := refs.GetLastCommit()
	if err != nil {
		return nil, err
	}
	currentHash := config.HEAD

	hashes := make([]string, 1)
	for currentHash != config.DOES_NOT_EXIST_HASH {

		hashes = append(hashes, currentHash)

		currentHash = current.ParentCommitHash
		current, err = objectio.LoadCommit(current.ParentCommitHash)
		if err != nil {
			return nil, err
		}
	}

	return hashes, nil
}

func findMergeBaseHash(hashA, hashB string) (string, error) {

	aList, err := getParentCommits(hashA)
	if err != nil {
		return "", fmt.Errorf("failed creating parent commit list %w", err)
	}
	bList, err := getParentCommits(hashB)
	if err != nil {
		return "", fmt.Errorf("failed creating parent commit list %w", err)
	}

	lastCommen := config.DOES_NOT_EXIST_HASH

	for i := 1; i < len(aList) && i < len(bList) && aList[len(aList)-i] == bList[len(bList)-i]; i++ {
		lastCommen = aList[len(aList)-i]
	}

	return lastCommen, nil

}

type mergeConflict struct {
	relPath  string
	actionA  treediff.ChangeAction
	newHashA string
	actionb  treediff.ChangeAction
	newHashb string
}

func findMergeConflicts(mapA, mapB treediff.ChangeMap) []mergeConflict {
	conflicts := make([]mergeConflict, 0)

	for relPath, changeA := range mapA {

		changeB, ok := mapB[relPath]
		if !ok {
			continue
		}

		if changeA.Action == treediff.Delete && changeB.Action == treediff.Delete {
			continue
		}
		// if modified or added, its suffices to compare the new hashes
		if changeA.NewHash == changeB.NewHash {
			continue
		}

		conflicts = append(conflicts, mergeConflict{relPath: relPath,
			actionA: changeA.Action, actionb: changeB.Action, newHashA: changeA.NewHash, newHashb: changeB.NewHash})
	}

	return conflicts
}

func merge(currentBranch, mergeFrom string) error {
	currentHash, err := refs.GetBranchCommitHash(currentBranch)
	if err != nil {
		return fmt.Errorf("cant load commit hash for branch '%s'. error: %w", currentBranch, err)
	}

	goalHash, err := refs.GetBranchCommitHash(mergeFrom)
	if err != nil {
		return fmt.Errorf("cant load commit hash for branch '%s'. error: %w", mergeFrom, err)
	}

	mergeBase, err := findMergeBaseHash(currentHash, goalHash)
	if err != nil {
		return err
	}

	if mergeBase == config.DOES_NOT_EXIST_HASH {
		fmt.Errorf("cant merge becase we did not detect common parent commit")
	}

	baseTree, err := objectio.LoadTreeByCommitHash(mergeBase)
	if err != nil {
		return fmt.Errorf("cant retrieve merge base tree: %w", err)
	}

	currentTree, err := objectio.LoadTreeByCommitHash(currentHash)
	if err != nil {
		return fmt.Errorf("cant retrieve last committed tree: %w", err)
	}

	mergeFromTree, err := objectio.LoadTreeByCommitHash(goalHash)
	if err != nil {
		return fmt.Errorf("cant retrieve tree from which we will merge: %w", err)
	}

	changesA := treediff.ChangeMap{}
	treediff.TreeDiff[treediff.ChangeMap](changesA, baseTree, currentTree, false)
	changesB := treediff.ChangeMap{}
	treediff.TreeDiff[treediff.ChangeMap](changesB, baseTree, mergeFromTree, false)

	mergeConflict := findMergeConflicts(changesA, changesB)
	mergeConflict = append(mergeConflict)
	return nil
}

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

	if err := merge(currentBranch, branchName); err != nil {
		return err.Error()
	}

	return ""
}
