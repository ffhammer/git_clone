package commit

import (
	"fmt"
	"git_clone/gvc/config"
	"git_clone/gvc/diffalgos"
	"git_clone/gvc/index"
	"git_clone/gvc/logging"
	"git_clone/gvc/objectio"
	"git_clone/gvc/refs"
	"git_clone/gvc/status"
	"git_clone/gvc/treebuild"
	"git_clone/gvc/treediff"
	"git_clone/gvc/utils"
)

func CalculateNumberOfInsertionsAndDeletions() (int, int, error) {
	changes, err := index.LoadIndexChanges()
	if err != nil {
		return 0, 0, fmt.Errorf("can't load index: %w", err)
	}

	nInsertions := 0
	nDels := 0

	for _, val := range changes {
		switch val.Action {
		case treediff.Add:

			object, err := objectio.LoadObject(val.NewHash)
			if err != nil {
				return 0, 0, fmt.Errorf("can't load object for file '%s': %w", val.RelPath, err)
			}
			nInsertions += utils.CountLines(object)
		case treediff.Delete:
			object, err := objectio.LoadObject(val.OldHash)
			if err != nil {
				return 0, 0, fmt.Errorf("can't load object for file '%s': %w", val.RelPath, err)
			}
			nDels += utils.CountLines(object)
		case treediff.Modify:
			oldObject, err := objectio.LoadObject(val.OldHash)
			if err != nil {
				return 0, 0, fmt.Errorf("cant load object for file '%s': %w", val.RelPath, err)
			}

			newObject, err := objectio.LoadObject(val.NewHash)
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

func Commit(message, author string, callStatusIfNoChanges bool) (string, error) {
	if message == "" {
		return "", fmt.Errorf("commit message cannot be empty")
	}
	if author == "" {
		return "", fmt.Errorf("author cannot be empty")
	}

	BParentCommitHash := config.DOES_NOT_EXIST_HASH
	if refs.InMergeState {
		openConflicts, err := index.GetOpenConflictFiles()
		if err != nil {
			return "", logging.ErrorF("can't commit because cant check openConflicts: %w", err)
		}
		mergeMetaData, err := refs.GetMergeMetaData()
		if err != nil {
			return "", logging.ErrorF("can't commit because cant check openConflicts: %w", err)
		}

		if len(openConflicts) > 0 {
			return "", logging.NewError("commit failed, because you still have unresolved conflicts.\nuse gvc status")
		}
		BParentCommitHash = mergeMetaData.MERGE_HEAD
	}

	changes, err := index.LoadIndexChanges()
	if err != nil {
		return "", fmt.Errorf("cant load changes: %w", err)
	} else if len(changes) == 0 && callStatusIfNoChanges { // if no changes return status
		logging.Info("found No changes, returning Status")
		return status.Status()
	}

	pointer, err := refs.LoadCurrentPointer()
	if err != nil {
		return "", fmt.Errorf("cant load pointer %w", err)
	}

	tree, err := treebuild.BuildTreeFromIndex()
	if err != nil {
		return "", fmt.Errorf("cant generate tree: %w", err)
	}

	nInsertions, nDeletions, err := CalculateNumberOfInsertionsAndDeletions()
	if err != nil {
		return "", fmt.Errorf("cant calculate number of insertions and deletions: %w", err)
	}

	treeHash, err := objectio.SaveTree(tree)
	if err != nil {
		return "", fmt.Errorf("cant save tree: %w", err)
	}

	newCommit := objectio.CommitMetdata{
		ParentCommitHash:  pointer.ParentCommitHash,
		BranchName:        pointer.BranchName,
		Author:            author,
		CommitMessage:     message,
		TreeHash:          treeHash,
		Date:              utils.GetCurrentTimeString(),
		BParentCommitHash: BParentCommitHash,
	}

	pointer.ParentCommitHash, err = objectio.SaveCommit(newCommit)
	if err != nil {
		return "", fmt.Errorf("cant save commit: %w", err)
	}

	logging.Info("Commit succesful")
	if err := refs.SaveCurrentPointer(pointer); err != nil {
		return "", fmt.Errorf("cant save current pointer: %w", err)
	}

	logging.Info("Commit succesful")
	if err := index.ClearAllChanges(); err != nil {
		return "", fmt.Errorf("could not clear index %w", err)
	}
	logging.Info("Commit succesful")
	return fmt.Sprintf("[%s %s] %s\n %d file(s) changed, %d insertions(+), %d deletions (-)", pointer.BranchName, pointer.ParentCommitHash, message, len(changes), nInsertions, nDeletions), nil
}
