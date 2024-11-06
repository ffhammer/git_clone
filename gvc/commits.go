package gvc

type CommitMetdata struct {
	ParentCommitHash string `json:"parent_commit_hash"`
	BranchName       string `json:"branch_name"`
	Author           string `json:"author"`
	CommitMessage    string `json:"commit_message"`
	Date             string `json:"date"`
}

// func commit(repoPath string, message string) error {

// 	changesPath := filepath.Join(repoDir, NEXT_COMMIT)
// 	treePath := changesPath.Join(repoDir, FILE_TABLE)

// 	// if treePath not exists or len(treepath) lines < 2 {
// 	// 	return fmt.Errorf("no changes to commit")
// 	// }

// 	return nil

// }
