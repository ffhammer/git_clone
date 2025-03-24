package index

import (
	"fmt"
	"git_clone/gvc/config"
	"git_clone/gvc/refs"
	"git_clone/gvc/treediff"

	"time"
)

type fileStatus string

const (
	UNCHANGE_FILE fileStatus = "unchanged"
	NEW_FILE      fileStatus = "new"
	MODIFIED_FILE fileStatus = "delete"
)

func partOfLastCommit(relPath, fileHash string) (fileStatus, string, error) {

	treeMap, err := refs.GetLastCommitsTree()
	if err != nil {
		return NEW_FILE, "", err
	}

	val, ok := treeMap[relPath]

	if !ok {
		return NEW_FILE, config.DOES_NOT_EXIST_HASH, nil
	}

	if val.FileHash == fileHash {
		return UNCHANGE_FILE, val.FileHash, nil
	}

	return MODIFIED_FILE, val.FileHash, nil
}

func AddFile(relPath, fileHash string) error {
	changes, err := LoadIndexChanges()
	if err != nil {
		return err
	}

	var status fileStatus
	var oldHash string
	if !refs.InMergeState {
		status, oldHash, err = partOfLastCommit(relPath, fileHash)
		if err != nil {
			return err
		}
	} else {
		status = MODIFIED_FILE
		oldHash, err = refs.GetConflictFileHash(relPath)
		if err != nil {
			return err
		}
	}

	var newEntry treediff.ChangeEntry

	if status == MODIFIED_FILE {
		newEntry = treediff.ChangeEntry{RelPath: relPath, NewHash: fileHash, OldHash: oldHash, EditedTime: time.Now().Unix(), Action: treediff.Modify}
	} else if status == NEW_FILE {
		newEntry = treediff.ChangeEntry{RelPath: relPath, NewHash: fileHash, OldHash: oldHash, EditedTime: time.Now().Unix(), Action: treediff.Add}
	} else { // in case of neither added nor modifed -> do nothing
		return nil
	}

	changes[relPath] = newEntry

	return saveIndexChanges(changes)

}

type FileNotPartOfIndexOrTreeError struct{}

func (m *FileNotPartOfIndexOrTreeError) Error() string {
	return "the file was not wart of index or a tree when removing"
}

func RemoveFile(relPath, fileHash string, force, cached bool) error {
	// cached and force only  important if file is part of the index
	changes, err := LoadIndexChanges()
	if err != nil {
		return err
	}
	status, oldHash, err := partOfLastCommit(relPath, fileHash)
	if err != nil {
		return err
	}

	if _, ok := changes[relPath]; ok && (cached || force) {
		delete(changes, relPath)
		return saveIndexChanges(changes)
	} else if ok {
		return fmt.Errorf("error: the following file has changes staged in the index:\n    %s\n    (use --cached to keep the file, or -f to force removal)", relPath)
	} else if status == NEW_FILE {
		return &FileNotPartOfIndexOrTreeError{}
	}

	newEntry := treediff.ChangeEntry{RelPath: relPath, OldHash: oldHash, NewHash: config.DOES_NOT_EXIST_HASH, EditedTime: time.Now().Unix(), Action: treediff.Delete}
	changes[relPath] = newEntry

	err = saveIndexChanges(changes)
	if err != nil {
		return err
	}

	return nil
}

func ClearAllChanges() error {
	return saveIndexChanges(treediff.ChangeMap{})
}

func RemoveFromIndex(relPath string) error {
	changes, err := LoadIndexChanges()
	if err != nil {
		return fmt.Errorf("cant load changes %w", err)
	}

	delete(changes, relPath)
	return saveIndexChanges(changes)
}
