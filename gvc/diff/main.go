package diff

import (
	"fmt"
	"git_clone/gvc/index"
	"git_clone/gvc/objectio"
	"git_clone/gvc/utils"
	"strings"
)

func TreeToTree(oldTreeInput, newTreeInput objectio.TreeMap, ignoreAdditions bool, pathsToMatch []string) (string, error) {

	oldTree := oldTreeInput
	newTree := newTreeInput

	if len(pathsToMatch) > 0 {
		var err error
		oldTree, err = utils.FilterRelPathKeyMapWithAbsPaths(oldTree, pathsToMatch)
		if err != nil {
			return "", fmt.Errorf("error while matching tree with files: %w", err)
		}

		newTree, err = utils.FilterRelPathKeyMapWithAbsPaths(newTree, pathsToMatch)
		if err != nil {
			return "", fmt.Errorf("error while matching tree with files: %w", err)
		}
	}

	changes := index.TreeDiff(oldTree, newTree, ignoreAdditions)

	var builder strings.Builder

	for _, change := range changes {

		oldLines := []string{}
		if change.Action != index.Add {
			if file, err := objectio.RetrieveFile(change.OldHash); err != nil {
				return "", fmt.Errorf("cant retrieve file '%s': %w", utils.RelPathToAbs(change.RelPath), err)

			} else {
				oldLines = utils.SplitLines(file)
			}

		}

		newLines := []string{}
		if change.Action != index.Delete {

			if file, err := objectio.RetrieveFile(change.NewHash); err != nil {
				return "", fmt.Errorf("cant retrieve file '%s': %w", utils.RelPathToAbs(change.RelPath), err)

			} else {
				newLines = utils.SplitLines(file)
			}
		}

		if res, err := GenerateFileDiff(change.OldHash, change.RelPath, change.NewHash, change.RelPath, oldLines, newLines); err != nil {
			return "", fmt.Errorf("error generating diff file for '%s': %w", change.RelPath, err)

		} else {
			builder.WriteString(res)
		}
	}
	return builder.String(), nil

}
