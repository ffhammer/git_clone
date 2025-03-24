package commands

import (
	"errors"
	"flag"
	"fmt"
	"git_clone/gvc/index"
	"git_clone/gvc/logging"
	"git_clone/gvc/objectio"
	"git_clone/gvc/refs"

	"git_clone/gvc/utils"
	"os"
	"path/filepath"
)

func restore(absPath string, source string, staged, worktTree bool) error {
	if refs.InMergeState && !staged {
		return logging.NewError("in merge staged you can only restore with --staged")
	}

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

		oldVal := ""
		if !refs.InMergeState {
			oldVal, err = objectio.RetrieveFile(entry.FileHash)
			if err != nil {
				return fmt.Errorf("error retriving file '%s': %w", matchedPath, err)
			}
		} else {
			conflictMetadata, err := refs.GetMergeMetaData()
			if err != nil {
				return err
			}
			for idx, conflict := range conflictMetadata.Conflicts {
				if conflict.RelPath == matchedPath {
					oldVal = conflictMetadata.ConflictHashes[idx]
				}
			}
			if oldVal == "" {
				return logging.NewError(fmt.Sprintf("file '%s' is not part of conflicts. In merge state can only restore conflicts.", matchedPath))
			}
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

func Restore(args []string) string {
	restoreCmd := flag.NewFlagSet("restore", flag.ExitOnError)
	// restoreSource := restoreCmd.String("source", "", "The branch or commit")
	restoreStaged := restoreCmd.Bool("staged", false, "")
	restoreWorktree := restoreCmd.Bool("worktree", false, "")

	restoreCmd.Parse(args)

	if len(restoreCmd.Args()) < 1 {
		fmt.Println("Error: expected file paths to restore.")
		restoreCmd.Usage()
		os.Exit(1)
	}

	for _, filePath := range restoreCmd.Args() {
		if err := restore(filePath, "HEAD", *restoreStaged, *restoreWorktree); err != nil {
			return fmt.Errorf("restore failed because: %w", err).Error()
		}
	}
	return ""
}
