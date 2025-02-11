package diff

import (
	"fmt"
	"git_clone/gvc/refs"
	"git_clone/gvc/treebuild"
)

func ToIndex(absInputPaths []string) (string, error) {

	dirTree, err := treebuild.BuildTreeFromDir()
	if err != nil {
		return "", err
	}

	indexTree, err := treebuild.BuildTreeFromIndex()
	if err != nil {
		return "", err
	}

	output, err := TreeToTree(indexTree, dirTree, true, absInputPaths, true) // include unstaged addition
	if err != nil {
		return "", fmt.Errorf("could not generate diff between trees: %w", err)
	}
	return output, nil
}

func CommitToWorkingDirectory(commitToCompareTo string, absInputPaths []string) (string, error) {
	dirTree, err := treebuild.BuildTreeFromDir()
	if err != nil {
		return "", err
	}

	commitTree, err := refs.LoadCommitTreeHeadAccpeted(commitToCompareTo)
	if err != nil {
		return "", fmt.Errorf("could not load tree for '%s': %w", commitToCompareTo, err)
	}

	output, err := TreeToTree(commitTree, dirTree, true, absInputPaths, true) // include unstaged addition
	if err != nil {
		return "", fmt.Errorf("could not generate diff between trees: %w", err)
	}
	return output, nil
}

func IndexToCommit(commitToCompareTo string, absInputPaths []string) (string, error) {
	// i guess i can do this if i have compare trees

	indexTree, err := treebuild.BuildTreeFromIndex()
	if err != nil {
		return "", err
	}

	commitTree, err := refs.LoadCommitTreeHeadAccpeted(commitToCompareTo)
	if err != nil {
		return "", fmt.Errorf("could not load tree for '%s': %w", commitToCompareTo, err)
	}

	output, err := TreeToTree(commitTree, indexTree, false, absInputPaths, false)
	if err != nil {
		return "", fmt.Errorf("could not generate diff between trees: %w", err)
	}
	return output, nil
}

func CommitToCommit(oldCommit, newCommit string, absInputPaths []string) (string, error) {
	// i guess i can do this if i have compare trees

	oldTree, err := refs.LoadCommitTreeHeadAccpeted(oldCommit)
	if err != nil {
		return "", fmt.Errorf("could not load tree for '%s': %w", oldCommit, err)
	}

	newTree, err := refs.LoadCommitTreeHeadAccpeted(newCommit)
	if err != nil {
		return "", fmt.Errorf("could not load tree for '%s': %w", newCommit, err)
	}

	output, err := TreeToTree(oldTree, newTree, false, absInputPaths, false)
	if err != nil {
		return "", fmt.Errorf("could not generate diff between trees: %w", err)
	}
	return output, nil
}
