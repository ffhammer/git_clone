package commands

import (
	"errors"
	"fmt"
	"git_clone/gvc/index"
	"git_clone/gvc/objectio"
	"git_clone/gvc/pointers"
	"git_clone/gvc/utils"
	"os"
	"path/filepath"

	"github.com/gobwas/glob"
)

func matchFileWithMapStringKey[T any](relPath string, m map[string]T) []string {
	files := make([]string, 0)

	if utils.IsGlob(relPath) {
		g := glob.MustCompile(relPath)

		for k, _ := range m {
			if g.Match(k) {
				files = append(files, k)
			}
		}

		return files
	}

	querySplittedParts := utils.SplitPath(relPath)

	for k, _ := range m {
		allMatched := true
		keySplittedParts := utils.SplitPath(k)

		if len(querySplittedParts) > len(keySplittedParts) {
			continue
		}

		for i := 0; i < len(querySplittedParts); i++ {

			allMatched = allMatched && (querySplittedParts[i] == keySplittedParts[i])
			if !allMatched {
				break
			}
		}

		if allMatched {
			files = append(files, k)
		}

	}

	return files
}

func Restore(absPath string, source string, staged, worktTree bool) error {
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

		matches := matchFileWithMapStringKey(relPath, changes)

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

	tree, err := pointers.GetLastCommitsTree()
	if err != nil {
		return err
	}

	matches := matchFileWithMapStringKey(relPath, tree)

	if len(matches) == 0 {
		return fmt.Errorf("pathspec '%s' did not match any file(s) known to git", absPath)
	}

	for _, matchedPath := range matches {

		entry, ok := tree[matchedPath]
		if !ok {
			return fmt.Errorf("this should not happen!")
		}

		oldVal, err := objectio.RetrieveFile(entry.FileHash)
		if err != nil {
			return fmt.Errorf("error retriving file '%s': %w", matchedPath, err)
		}

		if err := utils.MkdirIgnoreExists(filepath.Dir(matchedPath)); err != nil {
			return fmt.Errorf("error creating directories for '%s': %w", matchedPath, err)
		}

		if err := os.WriteFile(matchedPath, []byte(oldVal), os.ModePerm); err != nil {
			return fmt.Errorf("error writing file '%s': %w", matchedPath, err)
		}
	}

	return nil
}
