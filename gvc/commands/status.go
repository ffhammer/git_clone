package commands

import (
	"fmt"
	"git_clone/gvc/index"
	"git_clone/gvc/logging"
	"git_clone/gvc/merge"
	"git_clone/gvc/refs"
	"git_clone/gvc/treebuild"
	"git_clone/gvc/treediff"
	"strings"

	"github.com/fatih/color"
)

func status() (string, error) {

	pointer, err := refs.LoadCurrentPointer()
	if err != nil {
		return "", fmt.Errorf("could not load pointer: %s", err)
	}

	messages := make([]string, 0)
	messages = append(messages, fmt.Sprintf("On branch %s\n\n", pointer.BranchName))

	if refs.InMergeState {

		openConflicts, err := merge.GetOpenConflictFiles()
		if err != nil {
			return "", logging.Error(err)
		}

		if len(openConflicts) == 0 {
			messages = append(messages, "You have handled all merge conflicts, commit now (make this better)")
		} else {
			messages = append(messages, "You have unmerged paths.\n\t(fix conflicts and run 'git commit')\n\t(use 'git merge --abort' to abort the merge)")

			messages = append(messages, strings.Join(openConflicts, "\n\t"))
		}

	}

	changes, err := index.LoadIndexChanges()
	if err != nil {
		return "", fmt.Errorf("could not load changes: %s", err)
	}

	if len(changes) > 0 {
		messages = append(messages, "Changes to be committed:\n    (use 'gvc restore --staged <file>...' to unstage)")

		for _, change := range changes {
			messages = append(messages, color.GreenString(fmt.Sprintf("    %-9s    %s", change.Action+":", change.RelPath)))
		}
		messages = append(messages, "\n")
	}

	unstashedChanges, err := treebuild.GetUnstagedChangesList(false)
	if err != nil {
		return "", fmt.Errorf("could not load unstaged changes: %s", err)
	}

	addingIndex := 0

	if len(unstashedChanges) > 0 && unstashedChanges[0].Action != treediff.Add {
		messages = append(messages, "Changes not staged for commit:\n    (use 'gvc add/rm <file>...' to update what will be committed)\n    (use 'gvc restore <file>...' to discard changes in working directory)")

		for i := 0; i < len(unstashedChanges) && unstashedChanges[i].Action != treediff.Add; i++ {
			change := unstashedChanges[i]
			messages = append(messages, color.RedString(fmt.Sprintf("    %-9s    %s", change.Action+":", change.RelPath)))
			addingIndex++
		}
		messages = append(messages, "\n")
	}

	if len(unstashedChanges) > addingIndex {
		messages = append(messages, "Untracked files:\n    (use 'gvc add/rm <file>...' to update what will be committed)")

		for i := addingIndex; i < len(unstashedChanges); i++ {
			change := unstashedChanges[i]
			messages = append(messages, color.RedString(fmt.Sprintf("    %s", change.RelPath)))
		}
		messages = append(messages, "\n")
	}

	if len(changes) == 0 {
		messages = append(messages, "no changes added to commit (use 'gvc add')")
	}

	return strings.Join(messages, "\n"), nil
}

func Status() string {
	output, err := status()
	if err != nil {
		return fmt.Errorf("status failed because %w", err).Error()
	}
	return output
}
