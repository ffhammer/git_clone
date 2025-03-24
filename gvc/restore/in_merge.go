package restore

import (
	"fmt"
	"git_clone/gvc/index"
	"git_clone/gvc/logging"
	"git_clone/gvc/objectio"
	"git_clone/gvc/utils"
	"os"
	"path/filepath"
)

func InMergeRestore(absPath string) error {

	// for the moment souce is not implemented

	relPath, err := utils.MakePathRelativeToRepo(utils.RepoDir, absPath)
	if err != nil {
		return err
	}

	changes, err := index.LoadIndexChanges()
	if err != nil {
		return fmt.Errorf("cant load changes with index %w", err)
	}

	index_matches := utils.MatchFileWithMapStringKey(relPath, changes)

	if len(index_matches) == 0 {
		return logging.ErrorF("file '%s' is not in index. see 'gvc status' for more", relPath)
	}

	for _, matchedPath := range index_matches {

		entry, ok := changes[matchedPath]
		if !ok {
			return fmt.Errorf("this should not happens")
		}

		oldVal, err := objectio.RetrieveFile(entry.OldHash)
		if err != nil {
			return fmt.Errorf("error retriving file '%s': %w", matchedPath, err)
		}

		matchedAbsPath := utils.RelPathToAbs(matchedPath)
		if err := utils.MkdirIgnoreExists(filepath.Dir(matchedAbsPath)); err != nil {
			return fmt.Errorf("error creating directories for '%s': %w", matchedAbsPath, err)
		}

		if err := os.WriteFile(matchedAbsPath, []byte(oldVal), os.ModePerm); err != nil {
			return fmt.Errorf("error writing file '%s': %w", matchedAbsPath, err)
		}

		if err := index.RemoveFromIndex(matchedPath); err != nil {
			return fmt.Errorf("cant remove '%s' from index %w", matchedPath, err)
		}
	}

	return nil
}
