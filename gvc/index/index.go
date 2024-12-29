package index

import (
	"fmt"
	"git_clone/gvc/config"
	"git_clone/gvc/pointers"
	"time"
)

type ChangeAction string

const (
	Add    ChangeAction = "added"
	Modify ChangeAction = "modified"
	Delete ChangeAction = "deleted"
	// Stash    ChangeAction = "stash"
	// Unmerged ChangeAction = "unmerged"
)

type fileStatus string

const (
	UNCHANGE_FILE fileStatus = "unchanged"
	NEW_FILE      fileStatus = "new"
	MODIFIED_FILE fileStatus = "delete"
)

type ChangeEntry struct {
	RelPath    string       `json:"relpath"`
	FileHash   string       `json:"filehash"`
	OldHash    string       `json:"oldHash"`
	EditedTime int64        `json:"editTime"`
	Action     ChangeAction `json:"actiion"`
}

type ChangeMap map[string]ChangeEntry

func partOfLastCommit(relPath, fileHash string) (fileStatus, string, error) {

	treeMap, err := pointers.GetLastCommitsTree()
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

	status, oldHash, err := partOfLastCommit(relPath, fileHash)
	if err != nil {
		return err
	}

	var newEntry ChangeEntry

	if status == MODIFIED_FILE {
		newEntry = ChangeEntry{RelPath: relPath, FileHash: fileHash, OldHash: oldHash, EditedTime: time.Now().Unix(), Action: Modify}
	} else if status == NEW_FILE {
		newEntry = ChangeEntry{RelPath: relPath, FileHash: fileHash, OldHash: oldHash, EditedTime: time.Now().Unix(), Action: Add}
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

func RemoveFileFromIndex(relPath, fileHash string, force, cached bool) error {
	// cached and force only  important if file is part of the index
	changes, err := LoadIndexChanges()
	if err != nil {
		return err
	}
	status, _, err := partOfLastCommit(relPath, fileHash)
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

	newEntry := ChangeEntry{RelPath: relPath, FileHash: fileHash, EditedTime: time.Now().Unix(), Action: Delete}
	changes[relPath] = newEntry

	err = saveIndexChanges(changes)
	if err != nil {
		return err
	}

	return nil
}

func ClearAllChanges() error {
	return saveIndexChanges(ChangeMap{})
}
