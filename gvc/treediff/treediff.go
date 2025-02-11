package treediff

import (
	"git_clone/gvc/config"
	"git_clone/gvc/objectio"
)

func TreeDiff[T any](collector ChangeCollector[T], oldTree, newTree objectio.TreeMap, ignoreAdditions bool) {

	// Process deletions and modifications.
	for oldKey, oldVal := range oldTree {
		newVal, ok := newTree[oldKey]
		if !ok {
			collector.Add(ChangeEntry{
				RelPath: oldKey,
				NewHash: config.DOES_NOT_EXIST_HASH,
				OldHash: oldVal.FileHash,
				Action:  Delete,
			})
		} else if oldVal.FileHash != newVal.FileHash {
			collector.Add(ChangeEntry{
				RelPath: oldKey,
				NewHash: newVal.FileHash,
				OldHash: oldVal.FileHash,
				Action:  Modify,
			})
		}
	}

	// Process additions, if not ignored.
	if !ignoreAdditions {
		for newKey, newVal := range newTree {
			if _, ok := oldTree[newKey]; !ok {
				collector.Add(ChangeEntry{
					RelPath: newKey,
					NewHash: newVal.FileHash,
					OldHash: config.DOES_NOT_EXIST_HASH,
					Action:  Add,
				})
			}
		}
	}
}
