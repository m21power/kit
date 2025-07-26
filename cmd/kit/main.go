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

	default:
		fmt.Println("Unknown command:", cmd)
	}
}
