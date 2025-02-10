package refs

import (
	"fmt"
	"git_clone/gvc/config"
	"git_clone/gvc/utils"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/fatih/color"
)

func UpdateHead(newBranchName string) error {
	pathToCurrentPointer := filepath.Join(utils.RepoDir, config.CurrentBranchPointerFile)

	err := os.WriteFile(pathToCurrentPointer, []byte(newBranchName), os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed updating head file %s: %w", pathToCurrentPointer, err)
	}
	return nil
}

func LoadCurrentBranchName() (string, error) {
	pointer, err := LoadCurrentPointer()
	if err != nil {
		return "", fmt.Errorf("error loading current branch name:\n\tcannot load current pointer: %w", err)
	}
	return pointer.BranchName, nil
}

func ListBranches() (string, error) {
	refsFolder := filepath.Join(utils.RepoDir, config.RefsFolder)

	entries, err := os.ReadDir(refsFolder)
	if err != nil {
		return "", fmt.Errorf("cannot list refs directory '%s': %w", refsFolder, err)
	}

	currentBranch, err := LoadCurrentBranchName()
	if err != nil {
		return "", err
	}

	branches := make([]string, 0)

	for _, e := range entries {
		if !e.IsDir() { // Ensure it's a file and not a directory
			branches = append(branches, e.Name())
		}
	}

	// Sort branches alphabetically (case-insensitive)
	sort.Slice(branches, func(i, j int) bool {
		return strings.ToLower(branches[i]) < strings.ToLower(branches[j])
	})

	var builder strings.Builder

	for index, branch := range branches {
		if branch == currentBranch {
			builder.WriteString("* " + color.GreenString("%s", branch))
		} else {
			builder.WriteString(fmt.Sprintf("  %s", branch))
		}
		if index+1 != len(branches) {
			builder.WriteString("\n")

		}
	}

	return builder.String(), nil
}

func BranchExists(name string) (bool, error) {
	filePath := filepath.Join(utils.RepoDir, config.RefsFolder, name)

	_, err := os.Stat(filePath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, fmt.Errorf("error checking branch existence: %w", err)
}

func CreateNewBranch(name string) error {
	// Validate branch name (basic example, can be extended)
	if strings.Contains(name, " ") || strings.ContainsAny(name, "\\/") {
		return fmt.Errorf("error: invalid branch name '%s'", name)
	}

	filePath := filepath.Join(utils.RepoDir, config.RefsFolder, name)

	if exists, err := BranchExists(name); err != nil {
		return err
	} else if exists {
		return fmt.Errorf("fatal: a branch named '%s' already exists", name)
	}

	pointer, err := LoadCurrentPointer()
	if err != nil {
		return fmt.Errorf("cannot load current pointer: %w", err)
	}

	err = os.WriteFile(filePath, []byte(pointer.ParentCommitHash), 0644)
	if err != nil {
		return fmt.Errorf("could not write ref '%s' for branch: %w", filePath, err)
	}

	return nil
}

func DeleteBranch(name string) error {
	filePath := filepath.Join(utils.RepoDir, config.RefsFolder, name)

	if exists, err := BranchExists(name); err != nil {
		return err
	} else if !exists {
		return fmt.Errorf("error: branch '%s' not found", name)
	}

	// Prevent deleting the currently checked-out branch
	pointer, err := LoadCurrentPointer()
	if err != nil {
		return fmt.Errorf("cannot load current pointer: %w", err)
	}
	if name == pointer.BranchName {
		return fmt.Errorf("error: cannot delete the branch currently checked out: '%s'", name)
	}

	return os.Remove(filePath)
}
