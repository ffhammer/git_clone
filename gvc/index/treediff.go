package index

import (
	"git_clone/gvc/config"
	"git_clone/gvc/objectio"
)

func TreeDiff(oldTree, newTree objectio.TreeMap) ChangeList {

	changes := make(ChangeList, 0)

	for oldKey, oldVal := range oldTree {

		newVal, ok := newTree[oldKey]
		if !ok {
			changes = append(changes, ChangeEntry{RelPath: oldKey, NewHash: config.DOES_NOT_EXIST_HASH, OldHash: oldVal.FileHash, Action: Delete})
		} else if oldVal.FileHash != newVal.FileHash {
			changes = append(changes, ChangeEntry{RelPath: oldKey, NewHash: newVal.FileHash, OldHash: oldVal.FileHash, Action: Modify})
		}
	}

	for newKey, newVal := range newTree {
		_, ok := oldTree[newKey]
		if !ok {
			changes = append(changes, ChangeEntry{RelPath: newKey, NewHash: newVal.RelPath, OldHash: config.DOES_NOT_EXIST_HASH, Action: Add})
		}

	}

	return changes

}
