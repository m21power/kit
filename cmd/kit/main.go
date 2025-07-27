package main

import (
	"fmt"
	git "kit/internals/git"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: kit <command> [args]")
		os.Exit(1)
	}

	cmd := os.Args[1]

	switch cmd {
	case "init":
		err := git.InitKit()
		if err != nil {
			fmt.Println("Error:", err)
		}
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Usage: kit add <file|folder> [...]")
			os.Exit(1)
		}

		for _, path := range os.Args[2:] {
			err := git.AddKit(path)
			if err != nil {
				fmt.Println("Error adding", path, "→", err)
			}
		}
	case "commit":
		if len(os.Args) < 3 {
			fmt.Println("Usage: kit commit <message>")
			os.Exit(1)
		}
		message := os.Args[2]
		result, err := git.CommitGit(message)
		if err != nil {
			fmt.Println("Error committing:", err)
		}
		fmt.Println(result)
	case "log":
		count := 5
		if len(os.Args) > 2 {
			fmt.Sscanf(os.Args[2], "%d", &count)
		}
		logs, err := git.LogKit(count)
		if err != nil {
			fmt.Println("Error fetching logs:", err)
		} else {
			for _, log := range logs {
				fmt.Printf("\033[1;33mcommit\033[0m %s\n", log.Hash)                                          // yellow bold
				fmt.Printf("\033[1;32mAuthor:\033[0m %s <%s>\n", log.Author, log.Email)                       // green bold
				fmt.Printf("\033[1;36mDate:\033[0m   %s\n", log.Date.Format("Mon Jan 2 15:04:05 2006 -0700")) // cyan bold
				fmt.Println()
				fmt.Printf("    %s\n", log.Message)
				fmt.Println()
			}

		}
	default:
		fmt.Println("Unknown command:", cmd)
	}
}
