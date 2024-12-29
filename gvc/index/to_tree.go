package index

import (
	"fmt"
	"git_clone/gvc/config"
	"git_clone/gvc/ignorefiles"
	"git_clone/gvc/objectio"
	"git_clone/gvc/pointers"
	"git_clone/gvc/utils"
	"io/fs"
	"path/filepath"
)

func BuildTreeFromIndex() (objectio.TreeMap, error) {
	lastTree, err := pointers.GetLastCommitsTree()
	if err != nil {
		return objectio.TreeMap{}, err
	}
	changes, err := LoadIndexChanges()
	if err != nil {
		return objectio.TreeMap{}, err
	}

	for _, change := range changes {
		switch change.Action {
		case Delete:
			delete(lastTree, change.RelPath)
		case Add, Modify:
			lastTree[change.RelPath] = objectio.TreeEntry{RelPath: change.RelPath, FileHash: change.FileHash}
		}
	}
	return lastTree, nil
}
func BuildTreeFromDir() (objectio.TreeMap, error) {
	directory := utils.RepoDIr // Use RepoDir directly for clarity

	// Initialize the tree
	tree := make(objectio.TreeMap)

	err := filepath.Walk(directory, func(absPath string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the special directory (e.g., ".git" equivalent)
		if info.IsDir() && info.Name() == config.OWN_FOLDER_NAME {
			fmt.Printf("skipping a dir without errors: %+v \n", info.Name())
			return filepath.SkipDir
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}
		// Convert absolute path to relative path
		relPath, err := filepath.Rel(directory, absPath)
		if err != nil {
			return fmt.Errorf("error getting relative path: %w", err)
		}

		// Check if the file should be ignored
		ignore := ignorefiles.IsIgnored(relPath)
		if ignore {
			return nil
		}

		// Hash the file
		fileHash, err := utils.GetFileSHA1(absPath)
		if err != nil {
			return fmt.Errorf("hashing the file failed: %w", err)
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
