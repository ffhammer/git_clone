package treebuild

import (
	"fmt"
	"git_clone/gvc/config"
	"git_clone/gvc/ignorefiles"
	"git_clone/gvc/objectio"
	"git_clone/gvc/refs"
	"git_clone/gvc/utils"
	"io/fs"
	"path/filepath"
)

func BuildTreeFromDir() (objectio.TreeMap, error) {
	directory := utils.GetBasePath() // Use RepoDir directly for clarity

	lastCommitTree, err := refs.GetLastCommitsTree()
	if err != nil {
		return nil, err
	}

	// Initialize the tree
	tree := make(objectio.TreeMap)

	err = filepath.Walk(directory, func(absPath string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the special directory (e.g., ".git" equivalent)
		if info.IsDir() && info.Name() == config.OWN_FOLDER_NAME {
			return filepath.SkipDir
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}
		// Convert absolute path to relative path
		relPath, err := filepath.Rel(directory, absPath)
		if err != nil {
			fmt.Printf("error getting relative path: %s", err)
			return fmt.Errorf("error getting relative path: %w", err)
		}

		// Check if the file should be ignored
		// Hash the file
		fileHash, err := utils.GetFileSHA1(absPath)
		if err != nil {
			fmt.Printf("error hashing the file: %s", err)

			return fmt.Errorf("hashing the file failed: %w", err)
		}

		// if ignored, check if the file is already tracked (like with add -f). in case of new file return nil -> skip
		if ignorefiles.IsIgnored(relPath) {

			if _, ok := lastCommitTree[relPath]; !ok {
				return nil
			}
		}

		// Add the file to the tree
		tree[relPath] = objectio.TreeEntry{RelPath: relPath, FileHash: fileHash}

		return nil
	})

	if err != nil {
		return tree, fmt.Errorf("error walking through repository directory %q: %v", directory, err)
	}

	return tree, nil
}
