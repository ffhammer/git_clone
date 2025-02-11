package refs

import (
	"fmt"
	"git_clone/gvc/config"
	"git_clone/gvc/objectio"
	"git_clone/gvc/utils"
	"os"
	"path/filepath"
)

type CurrentBranchPointer struct {
	ParentCommitHash string
	BranchName       string
}

func SaveCurrentPointer(metadata CurrentBranchPointer) error {

	refsPath := filepath.Join(utils.RepoDir, config.RefsFolder, metadata.BranchName)

	err := os.WriteFile(refsPath, []byte(metadata.ParentCommitHash), 0644)
	if err != nil {
		return err
	}

	pathToCurrentPointer := filepath.Join(utils.RepoDir, config.CurrentBranchPointerFile)

	err = os.WriteFile(pathToCurrentPointer, []byte(metadata.BranchName), 0644)
	if err != nil {
		return fmt.Errorf("failed to write commit metadata to file %s: %w", pathToCurrentPointer, err)
	}

	return nil
}

func LoadCurrentPointer() (CurrentBranchPointer, error) {
	pathToCurrentPointer := filepath.Join(utils.RepoDir, config.CurrentBranchPointerFile)
	branchNameData, err := os.ReadFile(pathToCurrentPointer)
	if err != nil {
		return CurrentBranchPointer{}, fmt.Errorf("failed to read current branch pointer file %s: %w", pathToCurrentPointer, err)
	}

	branchName := string(branchNameData)
	refsPath := filepath.Join(utils.RepoDir, config.RefsFolder, branchName)
	parentCommitHashData, err := os.ReadFile(refsPath)
	if err != nil {
		return CurrentBranchPointer{}, fmt.Errorf("failed to read branch ref file %s: %w", refsPath, err)
	}

	return CurrentBranchPointer{
		ParentCommitHash: string(parentCommitHashData),
		BranchName:       branchName,
	}, nil
}

func GetLastCommit() (objectio.CommitMetdata, error) {

	branchPointer, err := LoadCurrentPointer()
	if err != nil {
		return objectio.CommitMetdata{}, fmt.Errorf("could not laod current branch pointer %s", err)
	}

	val, err := objectio.LoadCommit(branchPointer.ParentCommitHash)
	if err != nil {
		return objectio.CommitMetdata{}, fmt.Errorf("error while lading last commits from object storage: %w", err)
	}

	return val, nil

}

func GetLastCommitsTree() (objectio.TreeMap, error) {
	lastCommit, err := GetLastCommit()
	if err != nil {
		return objectio.TreeMap{}, err
	}

	val, err := objectio.LoadTree(lastCommit.TreeHash)
	if err != nil {
		return nil, fmt.Errorf("error while lading last commits tree from object storage: %w", err)
	}

	return val, nil
}

func LoadCommitTreeHeadAccpeted(commitHash string) (objectio.TreeMap, error) {
	if commitHash == config.HEAD {
		return GetLastCommitsTree()
	}

	return objectio.LoadTreeByCommitHash(commitHash)
}

func GetBranchCommitHash(branchName string) (string, error) {
	filePath := filepath.Join(utils.RepoDir, config.RefsFolder, branchName)

	file_bytes, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("cant load branch ref for '%s' error: '%w'", branchName, err)
	}

	return string(file_bytes), nil
}

func LoadBranchTree(branchName string) (objectio.TreeMap, error) {
	commitHash, err := GetBranchCommitHash(branchName)
	if err != nil {
		return nil, err
	}
	return objectio.LoadTreeByCommitHash(commitHash)
}
