package merge

import (
	"git_clone/gvc/logging"
	"git_clone/gvc/refs"
	"git_clone/gvc/treediff"
)

func findMergeConflicts(mapA, mapB treediff.ChangeMap) []refs.MergeConflict {
	logging.Debug("Starting findMergeConflicts")

	// Log all changes in mapA
	logging.Debug("Changes in A:")
	for key, change := range mapA {
		logging.DebugF("Path: '%s', Action: %s, Hash: %s", key, change.Action, change.NewHash)
	}

	// Log all changes in mapB
	logging.Debug("Changes in B:")
	for key, change := range mapB {
		logging.DebugF("Path: '%s', Action: %s, Hash: %s", key, change.Action, change.NewHash)
	}

	conflicts := make([]refs.MergeConflict, 0)

	// Compare the two sets of changes
	for relPath, changeA := range mapA {
		changeB, ok := mapB[relPath]
		if !ok {
			logging.DebugF("File '%s' exists only in A, skipping.", relPath)
			continue
		}

		logging.DebugF("Checking conflict for '%s': A(%s, %s) vs B(%s, %s)",
			relPath, changeA.Action, changeA.NewHash, changeB.Action, changeB.NewHash)

		// If both branches deleted the file, no conflict
		if changeA.Action == treediff.Delete && changeB.Action == treediff.Delete {
			logging.DebugF("Both branches deleted '%s', skipping.", relPath)
			continue
		}

		// If modifications are the same, no conflict
		if changeA.NewHash == changeB.NewHash {
			logging.DebugF("Same change detected in both branches for '%s', skipping.", relPath)
			continue
		}

		// Otherwise, it's a conflict
		logging.WarnF("Conflict detected for '%s'", relPath)
		conflicts = append(conflicts, refs.MergeConflict{
			RelPath:  relPath,
			ActionA:  changeA.Action,
			ActionB:  changeB.Action,
			NewHashA: changeA.NewHash,
			NewHashB: changeB.NewHash,
		})
	}

	logging.InfoF("Merge conflict detection completed. Total conflicts: %d", len(conflicts))
	return conflicts
}
