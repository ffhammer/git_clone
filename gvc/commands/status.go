package commands

import (
	"fmt"
	"git_clone/gvc/pointers"
)

func Status() ([]string, error) {

	pointer, err := pointers.LoadCurrentPointer()
	if err != nil {
		return nil, fmt.Errorf("could not load pointer: %s", err)
	}

	messages := make([]string, 0)
	messages = append(messages, fmt.Sprintf("On branch %s\n", pointer.BranchName))

	// 	On branch master
	// Your branch is up to date with 'origin/master'.

	// Changes to be committed:
	//   (use "git restore --staged <file>..." to unstage)
	//         modified:   gvc/index/index.go
	//         modified:   gvc/objectio/trees.go
	//         modified:   gvc/pointers/pointers.go

	// Changes not staged for commit:
	//   (use "git add/rm <file>..." to update what will be committed)
	//   (use "git restore <file>..." to discard changes in working directory)
	//         deleted:    gvc/changes.goo
	//         deleted:    gvc/commit_saving.goo
	//         modified:   gvc/index/index.go
	//         modified:   gvc/objectio/trees.go
	//         modified:   gvc/pointers/pointers.go

	// Untracked files:
	//   (use "git add <file>..." to include in what will be committed)
	//         gvc/index/to_tree.go
	//         gvc/status/

}
