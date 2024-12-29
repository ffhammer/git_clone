package commands

import (
	"fmt"
	"git_clone/gvc/index"
	"git_clone/gvc/pointers"

	"github.com/fatih/color"
)

func Status() ([]string, error) {

	pointer, err := pointers.LoadCurrentPointer()
	if err != nil {
		return nil, fmt.Errorf("could not load pointer: %s", err)
	}

	messages := make([]string, 0)
	messages = append(messages, fmt.Sprintf("On branch %s\n\n", pointer.BranchName))

	changes, err := index.LoadIndexChanges()
	if err != nil {
		return nil, fmt.Errorf("could not load changes: %s", err)
	}

	if len(changes) > 0 {
		messages = append(messages, "Changes to be committed:\n\t(use 'gvc restore --staged <file>...' to unstage)")

		for _, change := range changes {
			messages = append(messages, color.GreenString(fmt.Sprintf("\t\t%9s\t%s", change.Action+":", change.RelPath)))
		}
		messages = append(messages, "\n")
	}

	unstashedChanges, err := index.GetUnstagedChanges()
	if err != nil {
		return nil, fmt.Errorf("could not load unstaged changes: %s", err)
	}

	addingIndex := 0

	if len(unstashedChanges) > 0 && unstashedChanges[0].Action != index.Add {
		messages = append(messages, "Changes not staged for commit:\n\t(use 'gvc add/rm <file>...' to update what will be committed)\n\t(use 'gvc restore <file>...' to discard changes in working directory)")

		for i := 0; i < len(unstashedChanges) && unstashedChanges[i].Action != index.Add; i++ {
			change := unstashedChanges[i]
			messages = append(messages, color.RedString(fmt.Sprintf("\t\t%9s\t%s", change.Action+":", change.RelPath)))
			addingIndex++
		}
		messages = append(messages, "\n")
	}

	if len(unstashedChanges) > addingIndex {
		messages = append(messages, "Untracked files:\n\t(use 'gvc add/rm <file>...' to update what will be committed)")

		for i := addingIndex; i < len(unstashedChanges) && unstashedChanges[i].Action != index.Add; i++ {
			change := unstashedChanges[i]
			messages = append(messages, color.RedString(fmt.Sprintf("\t\t%s", change.RelPath)))
		}
		messages = append(messages, "\n")
	}

	if len(changes) == 0 {
		messages = append(messages, "no changes added to commit (use 'gvc add')")
	}

	return messages, nil
}
