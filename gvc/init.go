package gvc

import (
	"fmt"
	"os"
	"path/filepath"
)

func InitGVC() error {
	// Check if the repository already exists
	err := FindRepo()
	if err == nil {
		return fmt.Errorf("repository already exists at %s", repoDir)
	}

	if err := os.Mkdir(OWN_FOLDER_NAME, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", OWN_FOLDER_NAME, err)
	}

	commitsPath := filepath.Join(OWN_FOLDER_NAME, COMMITS_FOLDER)
	if err := os.Mkdir(commitsPath, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", commitsPath, err)
	}

	objetsPath := filepath.Join(OWN_FOLDER_NAME, OBJECT_FOLDER)
	if err := os.Mkdir(objetsPath, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", objetsPath, err)
	}

	inital_metdata := CurrentBranchPointer{ParentCommitHash: "none", BranchName: "main"}
	SaveCurrentBranchInfo(filepath.Join(OWN_FOLDER_NAME, CurrentInfoFile), inital_metdata)

	fmt.Println("Initialized a new repository at", OWN_FOLDER_NAME)
	return nil
}
