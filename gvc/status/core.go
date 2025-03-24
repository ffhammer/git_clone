package status

import (
	"fmt"
	"git_clone/gvc/index"
	"git_clone/gvc/logging"
	"git_clone/gvc/refs"
	"git_clone/gvc/treebuild"
	"git_clone/gvc/treediff"
	"strings"

	"github.com/fatih/color"
)

func Status() (string, error) {

	pointer, err := refs.LoadCurrentPointer()
	if err != nil {
		return "", fmt.Errorf("could not load pointer: %s", err)
	}

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("On branch %s\n\n", pointer.BranchName))

	if refs.InMergeState {

		openConflicts, err := index.GetOpenConflictFiles()
		if err != nil {
			return "", logging.Error(err)
		}

		if len(openConflicts) == 0 {
			builder.WriteString("You have handled all merge conflicts, commit now (make this better)")
		} else {
			builder.WriteString("You have unmerged paths.\n\t(fix conflicts and run 'git commit')\n\t(use 'git merge --abort' to abort the merge)\n")
			builder.WriteString("\t" + strings.Join(openConflicts, "\n\t"))
		}
		builder.WriteString("\n\n")
	}

	changes, err := index.LoadIndexChanges()
	if err != nil {
		return "", fmt.Errorf("could not load changes: %s", err)
	}

	if len(changes) > 0 {
		builder.WriteString("Changes to be committed:\n    (use 'gvc restore --staged <file>...' to unstage)\n")

		for _, change := range changes {
			builder.WriteString(color.GreenString(fmt.Sprintf("    %-9s    %s\n", change.Action+":", change.RelPath)))
		}
		builder.WriteString("\n")
	}

	unstashedChanges, err := treebuild.GetUnstagedChangesList(false)
	if err != nil {
		return "", fmt.Errorf("could not load unstaged changes: %s", err)
	}

	addingIndex := 0

	if len(unstashedChanges) > 0 && unstashedChanges[0].Action != treediff.Add {
		builder.WriteString("Changes not staged for commit:\n    (use 'gvc add/rm <file>...' to update what will be committed)\n    (use 'gvc restore <file>...' to discard changes in working directory)\n")

		for i := 0; i < len(unstashedChanges) && unstashedChanges[i].Action != treediff.Add; i++ {
			change := unstashedChanges[i]
			builder.WriteString(color.RedString(fmt.Sprintf("    %-9s    %s\n", change.Action+":", change.RelPath)))
			addingIndex++
		}
		builder.WriteString("\n")
	}

	if len(unstashedChanges) > addingIndex {
		builder.WriteString("Untracked files:\n    (use 'gvc add/rm <file>...' to update what will be committed)\n")

		for i := addingIndex; i < len(unstashedChanges); i++ {
			change := unstashedChanges[i]
			builder.WriteString(color.RedString(fmt.Sprintf("    %s\n", change.RelPath)))
		}
		builder.WriteString("\n")
	}

	if len(changes) == 0 {
		builder.WriteString("no changes added to commit (use 'gvc add')\n")
	}

	return builder.String(), nil
}
