package treebuild

import (
	"fmt"
	"git_clone/gvc/treediff"
)

func GetUnstagedChangesList(ignoreAdditions bool) (treediff.ChangeList, error) {
	// get changes that are not in the index currently
	newTree, err := BuildTreeFromDir()
	if err != nil {
		return nil, fmt.Errorf("error computing uncommited tree changes:\n building tree from dir: \n%s", err)
	}
	oldTree, err := BuildTreeFromIndex()
	if err != nil {
		return nil, fmt.Errorf("error computing uncommited tree changes:\n building tree from index: \n%s", err)
	}

	var cl treediff.ChangeList = treediff.ChangeList{}
	treediff.TreeDiff[treediff.ChangeList](&cl, oldTree, newTree, ignoreAdditions)
	return cl, nil
}

func GetUnstagedChangesMap(ignoreAdditions bool) (treediff.ChangeMap, error) {
	// get changes that are not in the index currently
	newTree, err := BuildTreeFromDir()
	if err != nil {
		return nil, fmt.Errorf("error computing uncommited tree changes:\n building tree from dir: \n%s", err)
	}
	oldTree, err := BuildTreeFromIndex()
	if err != nil {
		return nil, fmt.Errorf("error computing uncommited tree changes:\n building tree from index: \n%s", err)
	}

	cm := treediff.ChangeMap{}
	treediff.TreeDiff[treediff.ChangeMap](cm, oldTree, newTree, ignoreAdditions)
	return cm, nil
}
