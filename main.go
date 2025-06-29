package main

import (
	"fmt"
	"os"
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
		fmt.Println(`Usage: git <commands> [args...]
                                      git help - for more information`)
		return
	}
	commands := os.Args[1]
	switch commands {
	case "help":
		fmt.Print(helpMessage)
	}
}
