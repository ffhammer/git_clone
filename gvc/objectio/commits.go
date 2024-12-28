package objectio

import "git_clone/gvc/config"

type CommitMetdata struct {
	ParentCommitHash string `json:"parent_commit_hash"`
	BranchName       string `json:"branch_name"`
	Author           string `json:"author"`
	CommitMessage    string `json:"commit_message"`
	Date             string `json:"date"`
	TreeHash         string `json:"tree_hash"`
}

func LoadCommit(fileHash string) (CommitMetdata, error) {
	if fileHash == config.DOES_NOT_EXIST_HASH {
		return CommitMetdata{TreeHash: config.DOES_NOT_EXIST_HASH}, nil
	}

	return LoadJsonObject[CommitMetdata](fileHash)
}
