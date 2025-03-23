package index

import (
	"git_clone/gvc/refs"
)

func GetOpenConflictFiles() ([]string, error) {

	metaData, err := refs.GetMergeMetaData()
	if err != nil {
		return nil, err
	}

	index, err := LoadIndexChanges()
	if err != nil {
		return nil, err
	}

	names := make([]string, 0)

	for _, conflict := range metaData.Conflicts {

		if _, ok := index[conflict.RelPath]; !ok {
			names = append(names, conflict.RelPath)
		}

	}

	return names, nil
}
