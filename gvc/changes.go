package gvc

// getPathToHashMap reads a CSV file and returns a map where the key is the relative path, and the value is the hash.

func getParentTree(repoDir string, parHash string) (map[string]string, error) {
	// returns a map of path to hash for the parent commit
	if parHash == "none" {
		return nil, nil
	}

	// treePath := filepath.Join(repoDir, COMMITS_FOLDER, parHash, FILE_TABLE)

	// then relpath hash
	// create hash map

	return nil, nil
}
