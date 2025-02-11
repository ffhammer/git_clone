package treebuild

import (
	"fmt"
	"git_clone/gvc/index"
	"git_clone/gvc/objectio"
	"git_clone/gvc/refs"
	"git_clone/gvc/treediff"
)

func BuildTreeFromIndex() (objectio.TreeMap, error) {
	lastTree, err := refs.GetLastCommitsTree()
	if err != nil {
		return objectio.TreeMap{}, fmt.Errorf("error while building tree from index: %w", err)
	}
	changes, err := index.LoadIndexChanges()
	if err != nil {
		return objectio.TreeMap{}, fmt.Errorf("error while building tree from index: %w", err)
	}

	for _, change := range changes {
		switch change.Action {
		case treediff.Delete:
			delete(lastTree, change.RelPath)
		case treediff.Add, treediff.Modify:
			lastTree[change.RelPath] = objectio.TreeEntry{RelPath: change.RelPath, FileHash: change.NewHash}
		}
	}
	return lastTree, nil
}
