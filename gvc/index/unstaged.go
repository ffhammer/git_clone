package index

import (
	"fmt"
)

func GetUnstagedChanges(ignoreAdditions bool) ([]ChangeEntry, error) {
	// get changes that are not in the index currently
	newTree, err := BuildTreeFromDir()
	if err != nil {
		return nil, fmt.Errorf("error computing uncommited tree changes:\n building tree from dir: \n%s", err)
	}
	oldTree, err := BuildTreeFromIndex()
	if err != nil {
		return nil, fmt.Errorf("error computing uncommited tree changes:\n building tree from index: \n%s", err)
	}

	return TreeDiff(oldTree, newTree, ignoreAdditions), nil
}
