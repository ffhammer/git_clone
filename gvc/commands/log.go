package commands

import (
	"flag"
	"fmt"
	"git_clone/gvc/log"
)

func LogCommand(args []string) string {
	flagSet := flag.NewFlagSet("diff", flag.ExitOnError)
	patch := flagSet.Bool("patch", false, "Show changes introduced by each commit (patch format)")
	since := flagSet.String("since", "", "Show commits after a date")
	until := flagSet.String("until", "", "Show commits before a date")
	author := flagSet.String("author", "", "Filter commits by author")
	grep := flagSet.String("grep", "", "Filter commits by message pattern")

	if err := flagSet.Parse(args); err != nil {
		return err.Error()
	}

	output, err := log.StandardLog(*patch, *since, *until, *author, *grep)
	if err != nil {
		return fmt.Errorf("error while creating logs: %w", err).Error()
	}
	return output

}
