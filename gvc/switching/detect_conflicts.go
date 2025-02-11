package switching

import (
	"errors"
	"fmt"
	"git_clone/gvc/index"
	"git_clone/gvc/treebuild"
	"strings"
)

func FindNotChangeableFiles(relPaths []string) error {

	change_set, err := index.LoadIndexChanges()
	if err != nil {
		return err
	}

	unstaged, err := treebuild.GetUnstagedChangesMap(false)
	if err != nil {
		return err
	}

	var builder strings.Builder

	for _, relPath := range relPaths {

		_, ok1 := change_set[relPath]
		_, ok2 := unstaged[relPath]
		if ok1 || ok2 {
			builder.WriteString(fmt.Sprintf("\t%s\n", relPath))
		}
	}

	output := builder.String()

	if output == "" {
		return nil
	}

	return errors.New(output)
}
