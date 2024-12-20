package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

func main() {
	checkArgs()
}

func printHelp(exit bool) {
	fmt.Println("Usage: notes <command> [arguments]")
	fmt.Println("Commands:")
	fmt.Println("  add <note> - Add a note")
	fmt.Println("  list - List all notes")
	fmt.Println("  remove <id> - Remove a note")
	fmt.Println("  edit <id> <note> - Edit a note")
	if exit {
		os.Exit(1)
	}
}

func checkArgs() {
	if len(os.Args) < 2 {
		printHelp(true)
	}
	switch os.Args[1] {
	case "add":
		if len(os.Args) < 3 {
			color.Red("Error: Missing note argument")
			os.Exit(1)
		}
	case "remove":
		if len(os.Args) < 3 {
			color.Red("Error: Missing ID argument")
			os.Exit(1)
		}
	case "edit":
		if len(os.Args) < 4 {
			color.Red("Error: Missing ID and/or note argument")
			os.Exit(1)
		}
	}
}
