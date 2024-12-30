package commands

import (
	"flag"
	"fmt"
	"git_clone/gvc/commit"
	"git_clone/gvc/index"
	"git_clone/gvc/objectio"
	"git_clone/gvc/pointers"
	"git_clone/gvc/utils"
	"os"
)

func Commit() string {
	commitCmd := flag.NewFlagSet("commit", flag.ExitOnError)
	commitMessage := commitCmd.String("m", "", "The commit message")
	commitUser := commitCmd.String("u", "", "The commit user")
	commitCmd.Parse(os.Args[2:])

	if *commitMessage == "" {
		return "Error: commit message (-m) is required."
	}
	if *commitUser == "" {
		return "Error: commit user (-u) is required."
	}

	output, err := commit_func(*commitMessage, *commitUser)

	if err != nil {
		return fmt.Errorf("commit failed: \n    %w", err).Error()
	}
	return output
}

func commit_func(message, author string) (string, error) {
	if message == "" {
		return "", fmt.Errorf("commit message cannot be empty")
	}
	if author == "" {
		return "", fmt.Errorf("author cannot be empty")
	}

	changes, err := index.LoadIndexChanges()
	if err != nil {
		return "", fmt.Errorf("cant load changes: %w", err)
	} else if len(changes) == 0 { // if no changes return status
		return status()
	}

	pointer, err := pointers.LoadCurrentPointer()
	if err != nil {
		return "", fmt.Errorf("cant load pointer %w", err)
	}

	tree, err := index.BuildTreeFromIndex()
	if err != nil {
		return "", fmt.Errorf("cant generate tree: %w", err)
	}

	nInsertions, nDeletions, err := commit.CalculateNumberOfInsertionsAndDeletions()
	if err != nil {
		return "", fmt.Errorf("cant calculate number of insertions and deletions: %w", err)
	}

	treeHash, err := objectio.SaveTree(tree)
	if err != nil {
		return "", fmt.Errorf("cant save tree: %w", err)
	}

	newCommit := objectio.CommitMetdata{
		ParentCommitHash: pointer.ParentCommitHash,
		BranchName:       pointer.BranchName,
		Author:           author,
		CommitMessage:    message,
		TreeHash:         treeHash,
		Date:             utils.GetCurrentTimeString(),
	}

	pointer.ParentCommitHash, err = objectio.SaveCommit(newCommit)
	if err != nil {
		return "", fmt.Errorf("cant save commit: %w", err)
	}

	if err := pointers.SaveCurrentPointer(pointer); err != nil {
		return "", fmt.Errorf("cant save current pointer: %w", err)
	}

	if err := index.ClearAllChanges(); err != nil {
		return "", fmt.Errorf("could not clear index %w", err)
	}

	return fmt.Sprintf("[%s %s] %s\n %d file(s) changed, %d insertions(+), %d deletions (-)", pointer.BranchName, pointer.ParentCommitHash, message, len(changes), nInsertions, nDeletions), nil

}
