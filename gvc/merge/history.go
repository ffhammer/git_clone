package merge

import (
	"fmt"
	"git_clone/gvc/config"
	"git_clone/gvc/objectio"
	"git_clone/gvc/refs"
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
