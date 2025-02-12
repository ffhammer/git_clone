package merge

import (
	"git_clone/gvc/refs"
	"git_clone/gvc/treediff"
)

func findMergeConflicts(mapA, mapB treediff.ChangeMap) []refs.MergeConflict {
	conflicts := make([]refs.MergeConflict, 0)

	for relPath, changeA := range mapA {

		changeB, ok := mapB[relPath]
		if !ok {
			continue
		}

		if changeA.Action == treediff.Delete && changeB.Action == treediff.Delete {
			continue
		}
		// if modified or added, its suffices to compare the new hashes
		if changeA.NewHash == changeB.NewHash {
			continue
		}

		conflicts = append(conflicts, refs.MergeConflict{RelPath: relPath,
			ActionA: changeA.Action, ActionB: changeB.Action, NewHashA: changeA.NewHash, NewHashB: changeB.NewHash})
	}

	return conflicts
}
