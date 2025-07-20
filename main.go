package main

import (
	"fmt"
	"os"
)

const (
	ColorRed   = "\033[31m"
	ColorReset = "\033[0m"
)

const helpMessage = `
Usage:
  bitgit <command> [options]

Available commands:
  add             Add file contents to the index
  cat-file        Provide content or type/size information for repository objects
  check-ignore    Debug what is being ignored and why
  checkout        Switch branches or restore working tree files
  commit          Record changes to the repository
  hash-object     Compute object ID and optionally creates a blob from a file
  init            Create an empty Git repository
  log             Show commit logs
  ls-files        Show information about files in the index
  ls-tree         List the contents of a tree object
  rev-parse       Pick out and massage parameters
  rm              Remove files from the working tree and from the index
  show-ref        List references in the repository
  status          Show the working tree status
  tag             Create, list, delete or verify a tag object
`

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: bitgit <commands> [args...]")
		fmt.Println("       bitgit help - for more information")
		return
	}
	commands := os.Args[1]
	switch commands {
	case "help":
		fmt.Print(helpMessage)
	case "init":
		invokeInit()
	case "cat-file":
		invokeCatFile()
	}
}

func invokeInit() {
	path := ""
	if len(os.Args) > 2 {
		path = os.Args[2]
	} else {
		path = "."
	}
	repo, err := repoCreate(path)
	if err != nil {
		fmt.Printf("%serror creating repository:%s %v\n", ColorRed, ColorReset, err)
		return
	}
	fmt.Printf("initialized an empty git repository in %s\n", repo.GitDir)
}

func invokeCatFile() {
	if len(os.Args) < 4 {
		fmt.Printf("%sInvalid format:%s <command> TYPE OBJECT>", ColorRed, ColorReset)
	}
	oType, object := os.Args[2], os.Args[3]
	err := CmdCatFile(object, oType)
	if err != nil {
		fmt.Printf("%sERROR:%s %w", ColorRed, ColorReset, err)
	}
}
