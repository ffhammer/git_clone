package objectio

import (
	"fmt"
	"git_clone/gvc/config"
	"git_clone/gvc/utils"
	"io"
	"strings"
)

type TreeEntry struct {
	RelPath  string `json:"relpath"`
	FileHash string `json:"filehash"`
}

type TreeMap map[string]TreeEntry

func LoadTree(fileHash string) (TreeMap, error) {
	if fileHash == config.DOES_NOT_EXIST_HASH {
		return TreeMap{}, nil
	}

	val, err := LoadJsonObject[TreeMap](fileHash)
	if err != nil {
		return nil, fmt.Errorf("could not load tree '%s': %w", fileHash, err)
	}
	return val, nil
}

func SaveTree(tree TreeMap) (string, error) {
	reader, err := SerializeObject(tree)
	if err != nil {
		return "", err
	}

	// Convert reader to a buffer to allow re-reading
	buf := new(strings.Builder)
	_, err = io.Copy(buf, reader)
	if err != nil {
		return "", err
	}

	jsonString := buf.String()
	treeHash := utils.GetStringSHA1(jsonString)

	// Create a new reader from the buffer for SaveObject
	newReader := strings.NewReader(jsonString)

	return treeHash, SaveObject(treeHash, newReader)
}
