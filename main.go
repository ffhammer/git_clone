package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"git_clone/gvc/commands"
	"git_clone/gvc/utils"

	"github.com/fatih/color"
)

// Main usage information
func printGlobalHelp() {
	fmt.Println("Usage:")
	fmt.Println("  gvc <command> [options]")
	fmt.Println("\nCommands:")
	fmt.Println("  init       Initialize a new repository")
	fmt.Println("  add        Add files to the staging area")
	fmt.Println("\nUse 'gvc <command> -h' for more information about a command.")
}

func main() {
	// Define subcommands
	initCmd := flag.NewFlagSet("init", flag.ExitOnError)
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	// statusCmd := flag.NewFlagSet("status", flag.ExitOnError)

	// Flags for "add" subcommand
	force := addCmd.Bool("f", false, "Force adding the file even if it is ignored")

	// Check if a subcommand is provided
	if len(os.Args) < 2 {
		fmt.Println("Error: expected a subcommand.")
		printGlobalHelp()
		os.Exit(1)
	}

	if err := utils.FindRepo(); os.Args[1] != "init" && err != nil {
		fmt.Println("fatal: not a gvc repository (or any of the parent directories): .gvc")
		os.Exit(1)
	}

	// Handle subcommands
	switch os.Args[1] {
	case "init":
		initCmd.Parse(os.Args[2:])
		err := commands.InitGVC()
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		fmt.Println("Repository initialized successfully.")
	case "add":
		addCmd.Parse(os.Args[2:])
		if len(addCmd.Args()) < 1 {
			fmt.Println("Error: expected file paths to add.")
			addCmd.Usage()
			os.Exit(1)
		}
		for _, filePath := range addCmd.Args() {
			messages := commands.AddFiles(filePath, *force)
			for _, message := range messages {
				if strings.HasPrefix(message, "added") {
					color.Green(message)
				} else {
					color.Red(message)
				}
			}
		}
	case "status":

		lines, err := commands.Status()
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		for _, line := range lines {
			fmt.Println(line)
		}
	case "help", "-h", "--help":
		printGlobalHelp()
	default:
		fmt.Printf("Error: unrecognized command: %s\n", os.Args[1])
		printGlobalHelp()
		os.Exit(1)
	}
}
