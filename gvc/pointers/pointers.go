package pointers

import (
	"encoding/json"
	"fmt"
	"git_clone/gvc/config"
	"git_clone/gvc/objectio"
	"git_clone/gvc/utils"
	"os"
	"path/filepath"
)

type CurrentBranchPointer struct {
	ParentCommitHash string `json:"parent_commit_hash"`
	BranchName       string `json:"branch_name"`
}

func SaveCurrentPointer(metadata CurrentBranchPointer) error {

	pathToCurrentPointer := filepath.Join(utils.RepoDIr, config.CurrentBranchPointerFile)

	data, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize commit metadata: %w", err)
	}

	err = os.WriteFile(pathToCurrentPointer, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write commit metadata to file %s: %w", pathToCurrentPointer, err)
	}

	return nil
}

func LoadCurrentPointer() (CurrentBranchPointer, error) {
	filePath := filepath.Join(utils.RepoDIr, config.CurrentBranchPointerFile)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return CurrentBranchPointer{}, fmt.Errorf("current info file does not exist at %s", filePath)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return CurrentBranchPointer{}, fmt.Errorf("failed to read commit metadata file %s: %w", filePath, err)
	}

	var metadata CurrentBranchPointer
	err = json.Unmarshal(data, &metadata)
	if err != nil {
		return CurrentBranchPointer{}, fmt.Errorf("failed to deserialize commit metadata: %w", err)
	}

	return metadata, nil
}

func getLastCommit() (objectio.CommitMetdata, error) {

	branchPointer, err := LoadCurrentPointer()
	if err != nil {
		return objectio.CommitMetdata{}, fmt.Errorf("could not laod current branch pointer %s", err)
	}
	return objectio.LoadCommit(branchPointer.BranchName)

}
