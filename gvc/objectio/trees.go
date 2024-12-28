package objectio

type treeEntry struct {
	RelPath  string `json:"relpath"`
	FileHash string `json:"filehash"`
}

type TreeMap map[string]treeEntry
