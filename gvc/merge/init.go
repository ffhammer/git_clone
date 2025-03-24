package merge

import (
	"errors"
	"fmt"
	"git_clone/gvc/config"
	"git_clone/gvc/diffalgos"
	"git_clone/gvc/objectio"
	"git_clone/gvc/refs"
	"git_clone/gvc/utils"
	"os"
	"strings"
)

const (
	rightArrows = ">>>>>>>>>>>>>"
	leftArrows  = "<<<<<<<<<<<<<"
	connection  = "============="
)

func retrieveFile(relPath, hash string) (string, error) {
	if hash == config.DOES_NOT_EXIST_HASH {
		return "", nil
	}
	f, err := objectio.RetrieveFile(hash)
	if err != nil {
		return "", fmt.Errorf("error creating conflict file: could not load '%s' with hash '%s': %w", relPath, hash, err)
	}
	return f, nil
}

func createConflictFile(relPath, hashA, hashB, currentBranchName, sourceBranchName string) (string, error) {
	if hashA == config.DOES_NOT_EXIST_HASH && hashB == config.DOES_NOT_EXIST_HASH {
		return "", errors.New("logic error: both of the files does not exist for conflict resolution")
	}
	fileA, err := retrieveFile(relPath, hashA)
	if err != nil {
		return "", err
	}

	fileB, err := retrieveFile(relPath, hashB)
	if err != nil {
		return "", err
	}

	linesA := utils.SplitLines(fileA)
	linesB := utils.SplitLines(fileB)
	diffs := diffalgos.MyersDiff(linesA, linesB)

	var builder strings.Builder

	for startIndex := 0; startIndex < len(diffs); {
		if diffs[startIndex].Action == diffalgos.Keep {
			builder.WriteString(linesA[diffs[startIndex].OldLineNumber] + "\n")
			startIndex++
			continue
		}

		currentChange := diffs[startIndex].Action
		endIndex := startIndex + 1
		for endIndex < len(diffs) && diffs[endIndex].Action == currentChange {
			endIndex++
		}

		if currentChange == diffalgos.Insert {

			modified := startIndex-1 >= 0 && diffs[startIndex-1].Action == diffalgos.Delete
			if !modified {
				builder.WriteString(fmt.Sprintf("%s %s\n", rightArrows, sourceBranchName))
			}

			for _, d := range diffs[startIndex:endIndex] {
				builder.WriteString(linesB[d.NewLineNumber] + "\n")
			}

			if modified {
				builder.WriteString(fmt.Sprintf("%s %s\n", leftArrows, sourceBranchName))
			} else {
				builder.WriteString(fmt.Sprintf("%s\n", leftArrows))
			}

		} else {
			builder.WriteString(fmt.Sprintf("%s %s\n", rightArrows, currentBranchName))
			for _, d := range diffs[startIndex:endIndex] {
				builder.WriteString(linesA[d.OldLineNumber] + "\n")
			}

			if endIndex < len(diffs) && currentChange == diffalgos.Delete && diffs[endIndex].Action == diffalgos.Insert {
				builder.WriteString(fmt.Sprintf("%s %s\n", connection, sourceBranchName))
			} else {
				builder.WriteString(fmt.Sprintf("%s\n", leftArrows))
			}
		}
		startIndex = endIndex
	}

	absPath := utils.RelPathToAbs(relPath)
	if err := os.WriteFile(absPath, []byte(builder.String()), os.ModePerm); err != nil {
		return "", fmt.Errorf("could not write merge file for %s: %w", relPath, err)
	}

	fileHash, err := utils.GetFileSHA1(absPath)
	if err != nil {
		return "", fmt.Errorf("can't add file %s to objects because %w", absPath, err)

	}

	err = objectio.AddFileToObjects(absPath, fileHash)
	if err != nil {
		return "", fmt.Errorf("can't add file %s to objects because %w", absPath, err)
	}

	return fileHash, nil
}

func createMergeMessage(mergeConflicts []refs.MergeConflict, currentBranchName, sourceBranchHash, sourceBranchName string) (string, error) {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("Merge branch '%s' into %s\n", sourceBranchName, currentBranchName))
	builder.WriteString(fmt.Sprintf("Merged commit: %s\n", sourceBranchHash))

	if len(mergeConflicts) > 0 {
		builder.WriteString("\nConflicts:\n")
		for _, conflict := range mergeConflicts {
			builder.WriteString(fmt.Sprintf(" - %s\n", conflict.RelPath))
		}
	}

	return builder.String(), nil
}

func prepareMergeState(mergeConflicts []refs.MergeConflict, currentBranchName, sourceBranchHash, sourceBranchName string) error {
	message, err := createMergeMessage(mergeConflicts, currentBranchName, sourceBranchHash, sourceBranchName)
	if err != nil {
		return fmt.Errorf("initializing merge state failed: error creating merge message: %w", err)
	}

	conflictFileHashes := make([]string, len(mergeConflicts))
	for idx, conflict := range mergeConflicts {
		if hash, err := createConflictFile(conflict.RelPath, conflict.NewHashA, conflict.NewHashB, currentBranchName, sourceBranchName); err != nil {
			return fmt.Errorf("initializing merge state failed: error creating conflict file for '%s': %w", conflict.RelPath, err)
		} else {
			conflictFileHashes[idx] = hash
		}
	}
	metaData := refs.MergeMetaData{
		MERGE_HEAD:     sourceBranchHash,
		CURRENT_HEAD:   sourceBranchHash,
		MERGE_MESSAGE:  message,
		Conflicts:      mergeConflicts,
		ConflictHashes: conflictFileHashes,
	}

	if err := refs.SaveMergeMetaData(metaData); err != nil {
		return err
	}

	return nil
}
