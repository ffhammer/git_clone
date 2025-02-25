package merge

import (
	"fmt"
	"git_clone/gvc/config"
	"git_clone/gvc/logging"
	"git_clone/gvc/objectio"
)

func getParentCommits(startCommitHash string) ([]string, error) {

	current, err := objectio.LoadCommit(startCommitHash)
	if err != nil {
		return nil, err
	}
	currentHash := startCommitHash

	hashes := make([]string, 0)
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

func findMergeBaseHash(commitHashA, commitHashB string) (string, error) {

	aList, err := getParentCommits(commitHashA)
	if err != nil {
		return "", fmt.Errorf("failed creating parent commit list %w", err)
	}
	logging.DebugF("hashA '%s' results:", commitHashA)
	for _, hash := range aList {
		logging.DebugF("- %s", hash)
	}

	bList, err := getParentCommits(commitHashB)
	if err != nil {
		return "", fmt.Errorf("failed creating parent commit list %w", err)
	}
	logging.DebugF("hashB '%s' results:", commitHashB)
	for _, hash := range bList {
		logging.DebugF("- %s", hash)
	}

	lastCommen := config.DOES_NOT_EXIST_HASH

	visited := make(map[string]bool)
	for _, hash := range aList {
		visited[hash] = true
	}

	for _, hash := range bList {
		if visited[hash] {
			lastCommen = hash
			break
		}
	}

	return lastCommen, nil

}
