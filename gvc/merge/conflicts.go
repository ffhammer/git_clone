package merge

import "git_clone/gvc/treediff"

type mergeConflict struct {
	RelPath  string
	ActionA  treediff.ChangeAction
	NewHashA string
	ActionB  treediff.ChangeAction
	NewHashB string
}

func findMergeConflicts(mapA, mapB treediff.ChangeMap) []mergeConflict {
	conflicts := make([]mergeConflict, 0)

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

		conflicts = append(conflicts, mergeConflict{RelPath: relPath,
			ActionA: changeA.Action, ActionB: changeB.Action, NewHashA: changeA.NewHash, NewHashB: changeB.NewHash})
	}

	return conflicts
}
