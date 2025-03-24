package commands

import (
	"flag"
	"fmt"
	"git_clone/gvc/log"
)

func LogCommand(args []string) string {
	flagSet := flag.NewFlagSet("log", flag.ExitOnError)
	help := flagSet.Bool("help", false, "Get help documentation")
	helpShort := flagSet.Bool("h", false, "Get help documentation")
	patch := flagSet.Bool("patch", false, "Show changes introduced by each commit (patch format)")
	since := flagSet.String("since", "", "Show commits after a given date (e.g. 2024-01-01)")
	until := flagSet.String("until", "", "Show commits before a given date")
	author := flagSet.String("author", "", "Filter commits by author name")
	grep := flagSet.String("grep", "", "Filter commits by message content")

	if err := flagSet.Parse(args); err != nil {
		return err.Error()
	}
	if *help || *helpShort {
		return "gvc log [options]\n" +
			"Display the commit history.\n\n" +
			"Options:\n" +
			"  --patch         Show full diff for each commit\n" +
			"  --since <date>  Only show commits after this date (YYYY-MM-DD)\n" +
			"  --until <date>  Only show commits before this date (YYYY-MM-DD)\n" +
			"  --author <name> Filter by author\n" +
			"  --grep <msg>    Filter by message substring\n"
	}

	output, err := log.StandardLog(*patch, *since, *until, *author, *grep)
	if err != nil {
		return fmt.Errorf("error while creating logs: %w", err).Error()
	}
	return output

}
