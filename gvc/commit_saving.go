package gvc

import (
	"fmt"
	"path/filepath"
)

// func getLast

type commitInfo struct {
	ParentCommitHash string `json:"parent_commit_hash"`
	BranchName       string `json:"branch_name"`
}

var emptyCommit = commitInfo{}

func loadCommitInfo(commitHash string) (commitInfo, error) {
	if commitHash == COMMIT_HASH_NO_COMMIT {
		return emptyCommit, nil
	}
}

func getLastCommit() (commitInfo, error) {
	branchPointer, err := LoadCommitInfo(filepath.Join(OWN_FOLDER_NAME, CurrentInfoFile))

	if err != nil {
		return commitInfo{}, fmt.Errorf("cant read current branch info: %s", err)
	}

}
