package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/fatih/color"
)

func main() {
	checkArgs()
	switch os.Args[1] {
	case "add":
		addUnix()
	case "list":
		listUnix()
	case "rm-all":
		rmAllUnix()
	case "remove":
		id, err := strconv.Atoi(os.Args[2])
		if err != nil || id <= 0 {
			panic(err)
		}
		removeUnix(id)
	}
}

func printHelp(exit bool) {
	fmt.Println("Usage: notes <command> [arguments]")
	fmt.Println("Commands:")
	fmt.Println("  add - Add a note")
	fmt.Println("  list - List all notes")
	fmt.Println("  remove <id> - Remove a note")
	fmt.Println("  edit <id> - Edit a note")
	fmt.Println("  rm-all - Remove all notes")
	fmt.Println("  help - Print this help message")
	if exit {
		os.Exit(1)
	}
}

func checkArgs() {
	if len(os.Args) < 2 {
		printHelp(true)
	}
	switch os.Args[1] {
	case "remove":
		if len(os.Args) != 3 {
			color.Red("Error: Missing ID argument")
			os.Exit(1)
		}
	case "edit":
		if len(os.Args) < 3 {
			color.Red("Error: Missing ID and/or note argument")
			os.Exit(1)
		}
	}
}
