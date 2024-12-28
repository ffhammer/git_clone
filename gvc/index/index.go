package index

import (
	"time"
)

type ChangeAction string

const (
	Add      ChangeAction = "add"
	Modify   ChangeAction = "modify"
	Delete   ChangeAction = "delete"
	Stash    ChangeAction = "stash"
	Unmerged ChangeAction = "unmerged"
)

type ChangeEntry struct {
	RelPath    string       `json:"relpath"`
	FileHash   string       `json:"filehash"`
	EditedTime int64        `json:"editTime"`
	Action     ChangeAction `json:"actiion"`
}

type ChangeMap map[string]ChangeEntry

func AddFile(relPath, fileHash string) error {
	// noch problematiscsh
	changes, err := loadIndexChanges()
	if err != nil {
		return err
	}

	newEntry := ChangeEntry{RelPath: relPath, FileHash: fileHash, EditedTime: time.Now().Unix(), Action: Add}
	changes[relPath] = newEntry

	return saveIndexChanges(changes)

}

func DeleteFile(relPath, fileHash string) error {

	changes, err := loadIndexChanges()
	if err != nil {
		return err
	}

	newEntry := ChangeEntry{RelPath: relPath, FileHash: fileHash, EditedTime: time.Now().Unix(), Action: Delete}
	changes[relPath] = newEntry

	return saveIndexChanges(changes)

}
