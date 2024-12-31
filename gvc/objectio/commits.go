package objectio

import (
	"fmt"
	"git_clone/gvc/config"
	"git_clone/gvc/utils"
	"io"
	"strings"
)

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

func LoadTreeByCommitHash(fileHash string) (TreeMap, error) {
	commit, err := LoadCommit(fileHash)
	if err != nil {
		return nil, fmt.Errorf("error: could not retrieve commit '%s': %w", fileHash, err)
	}

	return LoadTree(commit.TreeHash)
}

func IsValidCommit(fileHash string) bool {

	if fileHash == config.HEAD {
		return true
	}

	if _, err := LoadCommit(fileHash); err != nil {
		return false
	}
	return true
}

func SaveCommit(commit CommitMetdata) (string, error) {
	reader, err := SerializeObject(commit)
	if err != nil {
		return "", err
	}

	// Convert reader to a buffer to allow re-reading
	buf := new(strings.Builder)
	_, err = io.Copy(buf, reader)
	if err != nil {
		return "", err
	}

	jsonString := buf.String()
	commitHash := utils.GetStringSHA1(jsonString)

	// Create a new reader from the buffer for SaveObject
	newReader := strings.NewReader(jsonString)

	return commitHash, SaveObject(commitHash, newReader)
}
