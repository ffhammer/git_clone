package main

import (
	"flag"
	"fmt"
	"os"

	"git_clone/gvc/commands"
	"git_clone/gvc/utils"
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
	force := addCmd.Bool("f", false, "Force adding the file even if it is ignored")

	rmCmd := flag.NewFlagSet("rm", flag.ExitOnError)
	rmChached := rmCmd.Bool("cached", false, "Only deletes file from .gvc not the actual file")
	rmRecursive := rmCmd.Bool("r", false, "")
	rmForce := rmCmd.Bool("f", false, "")
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
	case "add":
		addCmd.Parse(os.Args[2:])
		if len(addCmd.Args()) < 1 {
			fmt.Println("Error: expected file paths to add.")
			addCmd.Usage()
			os.Exit(1)
		}
		for _, filePath := range addCmd.Args() {
			output := commands.AddFiles(filePath, *force)
			fmt.Println(output)
		}
	case "rm":
		rmCmd.Parse(os.Args[2:])
		if len(rmCmd.Args()) < 1 {
			fmt.Println("Error: expected file paths to remove.")
			rmCmd.Usage()
			os.Exit(1)
		}
		for _, filePath := range rmCmd.Args() {
			output, err := commands.RemoveFile(filePath, *rmChached, *rmRecursive, *rmForce)

			if err != nil {
				fmt.Print(err)
				os.Exit(1)
			}

			fmt.Println(output)
		}

	case "status":

		output, err := commands.Status()
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		fmt.Print(output)
	case "help", "-h", "--help":
		printGlobalHelp()
	default:
		fmt.Printf("Error: unrecognized command: %s\n", os.Args[1])
		printGlobalHelp()
		os.Exit(1)
	}
}
