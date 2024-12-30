package index

import (
	"encoding/json"
	"errors"
	"fmt"
	"git_clone/gvc/config"
	"git_clone/gvc/utils"
	"io"
	"os"
	"path"
)

func getIndexPath() string {
	return path.Join(utils.RepoDir, config.NEXT_COMMIT)
}

func getChangesPath() string {
	return path.Join(getIndexPath(), config.INDEX_CHANGES)
}

func saveIndexChanges(changes ChangeMap) error {

	data, err := json.Marshal(changes)
	if err != nil {
		return fmt.Errorf("error serializing tree map: %w", err)
	}
	changesPath := getChangesPath()

	file, err := os.OpenFile(changesPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("error opening changes file at %s: %w", changesPath, err)
	}
	defer file.Close()

	// Write the JSON data to the file
	if _, err := file.Write(data); err != nil {
		return fmt.Errorf("error writing to changes file at %s: %w", changesPath, err)
	}

	return nil

}

func LoadIndexChanges() (ChangeMap, error) {
	changesPath := getChangesPath()

	// Open the file for reading
	file, err := os.Open(changesPath)
	if errors.Is(err, os.ErrNotExist) {
		// If the file does not exist, return an empty ChangeMap
		return ChangeMap{}, nil
	}

	if err != nil {
		return nil, fmt.Errorf("error while loading index: could not load index changes: %w", err)
	}
	defer file.Close()

	// Read the file content into a byte slice
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("error while loading index: could not read changes file: %w", err)
	}

	// Deserialize the JSON data into a ChangeMap
	var changes ChangeMap
	err = json.Unmarshal(data, &changes)
	if err != nil {
		return nil, fmt.Errorf("error while loading index: error deserializing changes map: %w", err)
	}

	return changes, nil
}
