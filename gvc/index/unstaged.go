package index

import (
	"fmt"
	"git_clone/gvc/refs"
)

func GetUnstagedChanges(includeAdditions bool) ([]ChangeEntry, error) {

	cwdTree, err := BuildTreeFromDir()
	if err != nil {
		return nil, fmt.Errorf("error computing uncommited tree changes:\n building tree from dir: \n%s", err)
	}
	commitTree, err := refs.GetLastCommitsTree()
	if err != nil {
		return nil, fmt.Errorf("error computing uncommited tree changes:\n building tree from index: \n%s", err)
	}

	return TreeDiff(commitTree, cwdTree, includeAdditions), nil
}
