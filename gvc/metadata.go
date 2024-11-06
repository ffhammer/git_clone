package gvc

import (
	"encoding/json"
	"fmt"
	"os"
)

type CurrentBranchPointer struct {
	ParentCommitHash string `json:"parent_commit_hash"`
	BranchName       string `json:"branch_name"`
}

const CurrentInfoFile = "current_info" // Define the file name for current info

func SaveCurrentBranchInfo(filePath string, metadata CurrentBranchPointer) error {

	data, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize commit metadata: %w", err)
	}

	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write commit metadata to file %s: %w", filePath, err)
	}

	return nil
}

func LoadCommitInfo(filePath string) (CurrentBranchPointer, error) {

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
