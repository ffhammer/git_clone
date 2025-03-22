package merge

import (
	"git_clone/gvc/logging"
	"git_clone/gvc/objectio"
	"git_clone/gvc/treediff"
	"git_clone/gvc/utils"
	"os"
)

func addFileAdditionsOfNewBranch(mapA, mapB treediff.ChangeMap) error {

	// Compare the two sets of changes
	for relPath, changeB := range mapB {
		_, ok := mapA[relPath]

		// the only case we want to handle
		if changeB.Action == treediff.Add && !ok {
			logging.DebugF("writing '%s'", relPath)
			file, err := objectio.RetrieveFile(changeB.NewHash)

			if err != nil {
				return logging.ErrorF("error in addFileAdditionsOfNewBranch: cant retrieve new hahs: %w", err)
			}

			err = os.WriteFile(utils.RelPathToAbs(relPath), []byte(file), os.ModePerm)
			if err != nil {
				return logging.ErrorF("error in addFileAdditionsOfNewBranch: cant write file '%s': %w", relPath, err)
			}
		}
	}
	return nil
}
