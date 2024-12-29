package index

import "git_clone/gvc/config"

func GetUnstagedChanges() ([]ChangeEntry, error) {

	cwdTree, err := BuildTreeFromDir()
	if err != nil {
		return nil, err
	}

	stagedTree, err := BuildTreeFromIndex()
	if err != nil {
		return nil, err
	}

	changes := make([]ChangeEntry, 0)

	for key, val := range stagedTree {

		new_val, ok := cwdTree[key]
		if !ok {
			changes = append(changes, ChangeEntry{FileHash: config.DOES_NOT_EXIST_HASH, OldHash: val.FileHash, RelPath: val.RelPath, Action: Delete})
		} else if new_val.FileHash != val.FileHash {
			changes = append(changes, ChangeEntry{FileHash: new_val.FileHash, OldHash: val.FileHash, RelPath: val.RelPath, Action: Modify})
		}
	}

	for key, val := range cwdTree {
		_, ok := stagedTree[key]
		if !ok {
			changes = append(changes, ChangeEntry{FileHash: val.FileHash, OldHash: config.DOES_NOT_EXIST_HASH, RelPath: val.RelPath, Action: Add})
		}
	}

	return changes, nil
}
