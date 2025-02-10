package index

import (
	"fmt"
)

func GetUncommitedChanges() ([]ChangeEntry, error) {
	cwdTree, err := BuildTreeFromDir()
	if err != nil {
		return nil, fmt.Errorf("error building tree from dir: \n%s", err)
	}
	stagedTree, err := BuildTreeFromIndex()
	if err != nil {
		return nil, fmt.Errorf("error building tree from index: \n%s", err)
	}

	return TreeDiff(stagedTree, cwdTree, true), nil

}
