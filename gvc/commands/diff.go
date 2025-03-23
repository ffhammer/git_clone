package commands

import (
	"flag"
	"fmt"
	"git_clone/gvc/config"
	"git_clone/gvc/diff"
	"git_clone/gvc/objectio"
	"git_clone/gvc/refs"
	"git_clone/gvc/utils"
	"os"
)

func DiffCommand(inputArgs []string) string {
	flagSet := flag.NewFlagSet("diff", flag.ExitOnError)
	noIndex := flagSet.Bool("no-index", false, "Compare files on hard drive")
	cached := flagSet.Bool("cached", false, "Compare staged changes to commit")

	if err := flagSet.Parse(inputArgs); err != nil {
		return fmt.Errorf("error parsing arguments: %w", err).Error()
	}

	args := flagSet.Args()

	// Ensure mutually exclusive flags
	if *noIndex && *cached {
		return "Invalid args: cannot use --no-index and --cached together"
	}

	if refs.InMergeState {
		return "Error: can't use diff in open merge state"
	}

	// Determine mode
	if *noIndex {
		return handleNoIndexMode(args)
	}
	if *cached {
		return handleCachedMode(args)
	}
	if len(args) >= 2 && objectio.IsValidCommit(args[0]) && objectio.IsValidCommit(args[1]) {
		return handleCommitToCommitMode(args)
	}
	if len(args) > 0 && objectio.IsValidCommit(args[0]) {
		return handleCommitToWorkingDirectoryMode(args)
	}

	return handleDiffToIndexMode(args)
}

func handleNoIndexMode(args []string) string {
	if len(args) != 2 {
		return "Invalid args: --no-index requires exactly two files"
	}

	fileA, err := os.ReadFile(args[0])
	if err != nil {
		return fmt.Errorf("error reading '%s': %w", args[0], err).Error()
	}
	hashA := utils.GetBytesSHA1(fileA)

	fileB, err := os.ReadFile(args[1])
	if err != nil {
		return fmt.Errorf("error reading '%s': %w", args[1], err).Error()
	}
	hashB := utils.GetBytesSHA1(fileB)

	output, err := diff.GenerateFileDiff(hashA, args[0], hashB, args[1], utils.SplitLines(string(fileA)), utils.SplitLines(string(fileB)))
	if err != nil {
		return fmt.Errorf("error generating diff: %w", err).Error()
	}
	return output
}

func handleCachedMode(args []string) string {
	commit := config.HEAD
	if len(args) > 0 && objectio.IsValidCommit(args[0]) {
		commit = args[0]
		args = args[1:]
	}

	output, err := diff.IndexToCommit(commit, args)
	if err != nil {
		return fmt.Errorf("error generating diff for cached mode: %w", err).Error()
	}
	return output
}

func handleCommitToCommitMode(args []string) string {
	output, err := diff.CommitToCommit(args[0], args[1], args[2:])
	if err != nil {
		return fmt.Errorf("error generating commit-to-commit diff: %w", err).Error()
	}
	return output
}

func handleCommitToWorkingDirectoryMode(args []string) string {
	commit := args[0]
	paths := args[1:]

	output, err := diff.CommitToWorkingDirectory(commit, paths)
	if err != nil {
		return fmt.Errorf("error generating diff against working directory: %w", err).Error()
	}
	return output
}

func handleDiffToIndexMode(args []string) string {
	output, err := diff.ToIndex(args)
	if err != nil {
		return fmt.Errorf("error generating diff to index: %w", err).Error()
	}
	return output
}
