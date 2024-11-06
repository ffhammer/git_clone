package main

import (
	"flag"
	"fmt"
	"git_clone/gvc"
	"os"
	"strings"

	"github.com/fatih/color"
)

func main() {
	// Define the subcommand "init" with a description
	initCmd := flag.NewFlagSet("init", flag.ExitOnError)

	addCmd := flag.NewFlagSet("add", flag.ExitOnError)

	// Define the "-f" flag for force in addCmd
	force := addCmd.Bool("f", false, "force adding the file even if it is ignored")

	// Check if there are any arguments
	if len(os.Args) < 2 {
		fmt.Println("expected subcommand")
		os.Exit(1)
	}

	// Parse the subcommand
	if os.Args[1] == "init" {
		initCmd.Parse(os.Args[2:])
		err := gvc.InitGVC() // Call the InitGVC function from gvc package
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		fmt.Println("Repository initialized successfully.")
		return
	}

	// Find the repository directory
	repoDir, err := gvc.FindRepo() // Capitalize FindRepo to make it accessible
	if err != nil {
		fmt.Println("Can't find initialized repository")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "add":
		addCmd.Parse(os.Args[2:])
		// Check that at least one file path is provided after flags
		if len(addCmd.Args()) < 1 {
			fmt.Println("expected file paths to add")
			fmt.Printf("\n %s \n", os.Args)
			fmt.Printf("\n %s \n", addCmd.Args())
			os.Exit(1)
		}

		// Loop over each file path provided after "add"
		for _, filePath := range addCmd.Args() {
			messages := gvc.AddFiles(repoDir, filePath, *force)

			for _, message := range messages {
				if strings.HasPrefix(message, "added") {
					color.Green(message)
				} else {
					color.Red(message)
				}
			}
		}
	default:
		fmt.Printf("unrecognized command: %s\n", os.Args[1])
		os.Exit(1)
	}
}
