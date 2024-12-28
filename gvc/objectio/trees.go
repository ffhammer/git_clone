package objectio

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

type treeEntry struct {
	RelPath  string `json:"relpath"`
	FileHash string `json:"filehash"`
}

type TreeMap map[string]treeEntry

func serializeTreeMap(tree TreeMap) (io.Reader, error) {
	// Serialize the tree map to JSON
	data, err := json.Marshal(tree)
	if err != nil {
		return nil, fmt.Errorf("error serializing tree map: %w", err)
	}

	// Wrap the JSON data in a bytes.Reader to create an io.Reader
	return bytes.NewReader(data), nil
}
func saveTree(tree TreeMap, fileHash string) error {
	// Serialize the TreeMap
	reader, err := serializeTreeMap(tree)
	if err != nil {
		return fmt.Errorf("error serializing tree: %w", err)
	}

	// Save the serialized tree as an object
	return SaveObject(fileHash, reader)
}
func deserializeTreeMap(data []byte) (TreeMap, error) {
	var tree TreeMap
	err := json.Unmarshal(data, &tree)
	if err != nil {
		return nil, fmt.Errorf("error deserializing tree map: %w", err)
	}
	return tree, nil
}

func loadTree(fileHash string) (TreeMap, error) {
	// Load the serialized tree object
	data, err := LoadObject(fileHash)
	if err != nil {
		return nil, fmt.Errorf("error loading tree object: %w", err)
	}

	// Deserialize the data into a TreeMap
	return deserializeTreeMap(data)
}
