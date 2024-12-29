package commit

import (
	"fmt"
	"git_clone/gvc/diffalgos"
	"git_clone/gvc/index"
	"git_clone/gvc/objectio"
	"git_clone/gvc/utils"
)

// type CommitMetdata struct {
// 	ParentCommitHash string `json:"parent_commit_hash"`
// 	BranchName       string `json:"branch_name"`
// 	Author           string `json:"author"`
// 	CommitMessage    string `json:"commit_message"`
// 	Date             string `json:"date"`
// 	TreeHash         string `json:"tree_hash"`
// }

func CalculateNumberOfInsertionsAndDeletions() (int, int, error) {
	changes, err := index.LoadIndexChanges()
	if err != nil {
		return 0, 0, fmt.Errorf("can't load index: %w", err)
	}

	nInsertions := 0
	nDels := 0

	for _, val := range changes {
		switch val.Action {
		case index.Add:

			object, err := objectio.LoadObject(val.FileHash)
			if err != nil {
				return 0, 0, fmt.Errorf("can't load object for file '%s': %w", val.RelPath, err)
			}
			nInsertions += utils.CountLines(object)
		case index.Delete:
			object, err := objectio.LoadObject(val.OldHash)
			if err != nil {
				return 0, 0, fmt.Errorf("can't load object for file '%s': %w", val.RelPath, err)
			}
			nDels += utils.CountLines(object)
		case index.Modify:
			oldObject, err := objectio.LoadObject(val.OldHash)
			if err != nil {
				return 0, 0, fmt.Errorf("cant load object for file '%s': %w", val.RelPath, err)
			}

			newObject, err := objectio.LoadObject(val.FileHash)
			if err != nil {
				return 0, 0, fmt.Errorf("cant load object for file '%s': %w", val.RelPath, err)
			}

			diffs := diffalgos.MyersDiff(utils.SplitLines(string(oldObject)), utils.SplitLines(string(newObject)))

			for _, diff := range diffs {
				if diff.Action == diffalgos.Insert {
					nInsertions++
				} else if diff.Action == diffalgos.Delete {
					nDels++
				}
			}

		}

	}
	return nInsertions, nDels, nil
}
