package diff

import (
	"fmt"
	"git_clone/gvc/index"
	"git_clone/gvc/objectio"
	"git_clone/gvc/utils"
	"os"
	"strings"
)

func TreeToTree(oldTreeInput, newTreeInput objectio.TreeMap, ignoreAdditions bool, pathsToMatch []string, loadNewFilesFromDisk bool) (string, error) {

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

			var file string

			if loadNewFilesFromDisk {
				if fileAsBytes, err := os.ReadFile(utils.RelPathToAbs(change.RelPath)); err != nil {
					fmt.Errorf("cant retrieve file '%s': %w", utils.RelPathToAbs(change.RelPath), err)
				} else {
					file = string(fileAsBytes)
				}

			} else {
				var err error
				file, err = objectio.RetrieveFile(change.NewHash)
				if err != nil {
					return "", fmt.Errorf("cant retrieve file '%s' %s: %w", utils.RelPathToAbs(change.RelPath), change.NewHash, err)

				}
			}
			newLines = utils.SplitLines(file)
		}

		if res, err := GenerateFileDiff(change.OldHash, change.RelPath, change.NewHash, change.RelPath, oldLines, newLines); err != nil {
			return "", fmt.Errorf("error generating diff file for '%s': %w", change.RelPath, err)

		} else {
			builder.WriteString(res)
		}
	}
	return builder.String(), nil

}
