package restore

import (
	"errors"
	"fmt"
	"git_clone/gvc/index"
	"git_clone/gvc/objectio"
	"git_clone/gvc/refs"
	"git_clone/gvc/utils"
	"os"
	"path/filepath"
)

func StandardRestore(absPath string, source string, staged, worktTree bool) error {

	// for the moment souce is not implemented
	if source != "HEAD" {
		return errors.New("source not implemtend yet")
	}

	relPath, err := utils.MakePathRelativeToRepo(utils.RepoDir, absPath)
	if err != nil {
		return err
	}
	if staged {
		changes, err := index.LoadIndexChanges()
		if err != nil {
			return fmt.Errorf("cant load changes with index %w", err)
		}

		matches := utils.MatchFileWithMapStringKey(relPath, changes)

		if len(matches) == 0 {
			return fmt.Errorf("pathspec '%s' did not match any file(s) known to git", absPath)
		}

		for _, matchedPath := range matches {
			if err := index.RemoveFromIndex(matchedPath); err != nil {
				return fmt.Errorf("cant remove '%s' from index %w", matchedPath, err)
			}
		}
	}

	if staged && !worktTree {
		return nil
	}

	tree, err := refs.GetLastCommitsTree()
	if err != nil {
		return err
	}

	matches := utils.MatchFileWithMapStringKey(relPath, tree)

	if len(matches) == 0 {
		return fmt.Errorf("pathspec '%s' did not match any file(s) known to git", absPath)
	}
	for _, matchedPath := range matches {

		entry, ok := tree[matchedPath]
		if !ok {
			return fmt.Errorf("this should not happens")
		}

		oldVal, err := objectio.RetrieveFile(entry.FileHash)
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
	}

	return nil
}
