package objectio

import "git_clone/gvc/config"

type TreeEntry struct {
	RelPath  string `json:"relpath"`
	FileHash string `json:"filehash"`
}

type TreeMap map[string]TreeEntry

func LoadTree(fileHash string) (TreeMap, error) {
	if fileHash == config.DOES_NOT_EXIST_HASH {
		return TreeMap{}, nil
	}

	return LoadJsonObject[TreeMap](fileHash)
}
