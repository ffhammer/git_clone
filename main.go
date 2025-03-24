package main

import (
	"fmt"
	"os"

	"git_clone/gvc/commands"
	"git_clone/gvc/refs"
	"git_clone/gvc/utils"
)

// Main usage information
func getGlobalHelp() string {
	output := "Usage:\n"
	output += "  gvc <command> [options]\n\n"
	output += "Commands:\n"
	output += "  init       Initialize a new repository              [gvc init -h]\n"
	output += "  add        Add files to the staging area            [gvc add -h]\n"
	output += "  rm         Remove files from the index/worktree     [gvc rm -h]\n"
	output += "  commit     Commit staged changes                    [gvc commit -h]\n"
	output += "  status     Show working tree and index state        [gvc status -h]\n"
	output += "  restore    Restore file content or index            [gvc restore -h]\n"
	output += "  diff       Show differences between states          [gvc diff -h]\n"
	output += "  log        Show commit logs                         [gvc log -h]\n"
	output += "  branch     Create/delete/list branches              [gvc branch -h]\n"
	output += "  checkout   Switch branches                          [gvc checkout -h]\n"
	output += "  merge      Merge another branch into current        [gvc merge -h]\n"
	output += "  set        View or change config settings           [gvc set -h]\n"
	output += "\nUse 'gvc <command> -h' for detailed help on a command.\n"
	return output
}

func main() {
	// Check if a subcommand is provided
	if len(os.Args) < 2 {
		fmt.Print(getGlobalHelp())
		os.Exit(0)
	}

	if err := utils.FindRepo(); os.Args[1] != "init" && err != nil {
		fmt.Println("fatal: not a gvc repository (or any of the parent directories): .gvc")
		os.Exit(1)
	} else if os.Args[1] != "init" {
		if err := refs.CheckForMergeState(); err != nil {
			fmt.Printf("fatal: failed to checked if in merge state: %s\n", err.Error())
			os.Exit(1)
		}

	}

	// Handle subcommands
	var output string
	switch os.Args[1] {
	case "init":
		output = commands.InitGVC(os.Args[2:])
	case "add":
		output = commands.AddCommand(os.Args[2:])
	case "rm":
		output = commands.RMCommand(os.Args[2:])
	case "commit":
		output = commands.Commit(os.Args[2:])
	case "status":
		output = commands.Status(os.Args[2:])
	case "restore":
		output = commands.Restore(os.Args[2:])
	case "diff":
		output = commands.DiffCommand(os.Args[2:])
	case "log":
		output = commands.LogCommand(os.Args[2:])
	case "branch":
		output = commands.BranchCommand(os.Args[2:])
	case "checkout":
		output = commands.CheckoutCommand(os.Args[2:])
	case "merge":
		output = commands.MergeCommand(os.Args[2:])
	case "set":
		output = commands.SettingsCommand(os.Args[2:])
	case "help", "-h", "--help":
		output = getGlobalHelp()
	default:
		output = getGlobalHelp()
	}

	if len(output) > 0 {

		fmt.Fprintln(os.Stderr, output)
	}

}
