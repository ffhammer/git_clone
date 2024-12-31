package commands

import (
	"flag"
	"fmt"
	"git_clone/gvc/config"
	"git_clone/gvc/diff"
	"git_clone/gvc/index"
	"git_clone/gvc/objectio"
	"git_clone/gvc/pointers"
	"git_clone/gvc/utils"
	"os"
	"strings"
)

func diffToIndex() (string, error) {

	changes, err := index.GetUnstagedChanges(false)
	if err != nil {
		return "", fmt.Errorf("error while loading unstaged changes: %w", err)
	}

	var builder strings.Builder

	for _, change := range changes {

		file, err := objectio.RetrieveFile(change.OldHash)
		if err != nil {
			return "", fmt.Errorf("cant retrieve old version for file '%s': %w", change.RelPath, err)
		}
		oldLines := utils.SplitLines(file)

		newLines := []string{}
		if change.Action == index.Modify {

			if file, err := os.ReadFile(utils.RelPathToAbs(change.RelPath)); err != nil {
				return "", fmt.Errorf("cant read file '%s': %w", utils.RelPathToAbs(change.RelPath), err)

			} else {
				newLines = utils.SplitLines(string(file))
			}
		}

		if res, err := diff.GenerateFileDiff(change.OldHash, change.RelPath, change.NewHash, change.RelPath, oldLines, newLines); err != nil {
			return "", fmt.Errorf("error generating diff file for '%s': %w", change.RelPath, err)

		} else {
			builder.WriteString(res)
		}
	}
	return builder.String(), nil
}

func diffTreeToTree(oldTree, newTree objectio.TreeMap) (string, error) {

	changes := index.TreeDiff(oldTree, newTree)

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

		if res, err := diff.GenerateFileDiff(change.OldHash, change.RelPath, change.NewHash, change.RelPath, oldLines, newLines); err != nil {
			return "", fmt.Errorf("error generating diff file for '%s': %w", change.RelPath, err)

		} else {
			builder.WriteString(res)
		}
	}
	return builder.String(), nil

}

func diffIndexToCommit(commitToCompareTo string) (string, error) {
	// i guess i can do this if i have compare trees

	indexTree, err := index.BuildTreeFromIndex()
	if err != nil {
		return "", err
	}

	commitTree, err := pointers.LoadCommitTreeHeadAccpeted(commitToCompareTo)
	if err != nil {
		return "", fmt.Errorf("could not load tree for '%s': %w", commitToCompareTo, err)
	}

	output, err := diffTreeToTree(commitTree, indexTree)
	if err != nil {
		return "", fmt.Errorf("could not generate diff between trees: %w", err)
	}
	return output, nil
}

func DiffCommand(inputArgs []string) string {

	flag := flag.NewFlagSet("diff", flag.ExitOnError)
	noIndex := flag.Bool("no-index", false, "Compare files on hard drive")
	cached := flag.Bool("cached", false, "Compare staged changes to commit")

	if err := flag.Parse(inputArgs); err != nil {
		return fmt.Errorf("error while parsing the args: %w", err).Error()
	}

	args := flag.Args()

	if *noIndex && *cached {
		return "Invalid args. Cant have -no-index and -cached activated at same time"
	} else if *noIndex {

		if len(args) != 2 {
			return "in case of -no-index you need two valid files"
		}

		fileA, err := os.ReadFile(args[0])
		if err != nil {
			return fmt.Errorf("error while reading '%s': %w", args[0], err).Error()
		}
		hashA := utils.GetBytesSHA1(fileA)

		fileB, err := os.ReadFile(args[1])
		if err != nil {
			return fmt.Errorf("error while reading '%s': %w", args[0], err).Error()
		}
		hashB := utils.GetBytesSHA1(fileB)

		output, err := diff.GenerateFileDiff(hashA, args[0], hashB, args[1], utils.SplitLines(string(fileA)), utils.SplitLines(string(fileB)))
		if err != nil {
			return fmt.Errorf("error while diffing files: %w", err).Error()
		}
		return output
	} else if *cached {
		inputCommit := config.HEAD
		if len(args) == 1 {
			inputCommit = args[0]
		} else if len(args) > 1 {
			return "only a single commit arg accepted"
		}

		output, err := diffIndexToCommit(inputCommit)
		if err != nil {
			return err.Error()
		}
		return output

	}

	output, err := diffToIndex()
	if err != nil {
		return err.Error()
	}

	return output
}
