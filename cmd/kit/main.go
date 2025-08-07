package main

import (
	"kit/routes"

	"github.com/gorilla/mux"
)

func main() {
	route := mux.NewRouter()
	r := routes.NewRouter(route)
	r.RegisterRoute()
	r.Run(":8080")
	// if len(os.Args) < 2 {
	// 	fmt.Println("Usage: kit <command> [args]")
	// 	os.Exit(1)
	// }

	// cmd := os.Args[1]
	// if cmd != "init" {
	// 	err := utils.CheckKit()
	// 	if err != nil {
	// 		fmt.Println("Error:", err)
	// 		os.Exit(1)
	// 	}
	// }
	// switch cmd {
	// // case "init":
	// // 	err := git.InitKit()
	// // 	if err != nil {
	// // 		fmt.Println("Error:", err)
	// // 	}
	// case "add":
	// 	if len(os.Args) < 3 {
	// 		fmt.Println("Usage: kit add <file|folder> [...]")
	// 		os.Exit(1)
	// 	}

	// 	// for _, path := range os.Args[2:] {
	// 	// 	err := git.AddKit(path)
	// 	// 	if err != nil {
	// 	// 		fmt.Println("Error adding", path, "→", err)
	// 	// 	}
	// 	// }
	// case "commit":
	// 	if len(os.Args) < 3 {
	// 		fmt.Println("Usage: kit commit <message>")
	// 		os.Exit(1)
	// 	}
	// 	message := os.Args[2]
	// 	result, err := git.CommitGit(message)
	// 	if err != nil {
	// 		fmt.Println("Error committing:", err)
	// 	}
	// 	fmt.Println(result)
	// case "log":
	// 	// count := 5
	// 	// if len(os.Args) > 2 {
	// 	// 	fmt.Sscanf(os.Args[2], "%d", &count)
	// 	// }
	// 	// logs, err := git.LogKit(count)
	// 	// if err != nil {
	// 	// 	fmt.Println("Error fetching logs:", err)
	// 	// } else {
	// 	// 	heaad, err := utils.GetHead()
	// 	// 	if err != nil {
	// 	// 		fmt.Println("Error getting head:", err)
	// 	// 	}
	// 	// 	fmt.Printf("\033[1;34mHEAD -> %s\033[0m\n", heaad) // blue bold
	// 	// 	for _, log := range logs {
	// 	// 		fmt.Printf("\033[1;33mcommit\033[0m %s\n", log.Hash)                                          // yellow bold
	// 	// 		fmt.Printf("\033[1;32mAuthor:\033[0m %s <%s>\n", log.Author, log.Email)                       // green bold
	// 	// 		fmt.Printf("\033[1;36mDate:\033[0m   %s\n", log.Date.Format("Mon Jan 2 15:04:05 2006 -0700")) // cyan bold
	// 	// 		fmt.Println()
	// 	// 		fmt.Printf("    %s\n", log.Message)
	// 	// 		fmt.Println()
	// 	// 	}

	// 	// }

	// case "status":
	// 	branch, err := utils.GetHead()
	// 	if err != nil {
	// 		fmt.Println("Error getting head: ", err)
	// 	}
	// 	treeComm, err := utils.GetCommitTreeHash(branch)
	// 	if err != nil {
	// 		fmt.Println("Error getting tree hash: ", err)
	// 	}
	// 	result, err := git.StatusKit(treeComm, "", make(map[string]bool))
	// 	if err != nil {
	// 		fmt.Println("Error showing content:", err)
	// 	} else {

	// 		res, err := git.IsChanged(result, ".")
	// 		if err != nil {
	// 			fmt.Println("Error checking changes:", err)
	// 		} else {
	// 			var stagedFiles []string
	// 			var unstagedFiles []string

	// 			for path, status := range res {
	// 				if status.Staged {
	// 					stagedFiles = append(stagedFiles, path)
	// 				} else {
	// 					unstagedFiles = append(unstagedFiles, path)
	// 				}
	// 			}

	// 			// Print staged files
	// 			if len(stagedFiles) > 0 {
	// 				fmt.Println("\033[1;36m Files ready to commit:\033[0m")
	// 				for _, file := range stagedFiles {
	// 					fmt.Printf("  \033[1;32m- %s\033[0m\n", file)
	// 				}
	// 			}
	// 			// Print unstaged files
	// 			if len(unstagedFiles) > 0 {
	// 				fmt.Println("\n\033[1;33m Not tracked yet. Run: kit add .\033[0m")
	// 				for _, file := range unstagedFiles {
	// 					fmt.Printf("  \033[1;31m- %s\033[0m\n", file)
	// 				}
	// 			}

	// 		}
	// 	}
	// case "branch":
	// 	if len(os.Args) < 3 {
	// 		branch, err := utils.GetHead()
	// 		if err != nil {
	// 			fmt.Println("Error getting current branch:", err)
	// 			os.Exit(1)
	// 		}
	// 		fmt.Println("Current branch:", branch)
	// 		os.Exit(0)
	// 	}
	// 	name := os.Args[2]
	// 	err := git.CreateBranch(name)
	// 	if err != nil {
	// 		fmt.Println("Error creating branch:", err)
	// 	}
	// case "checkout":
	// 	if len(os.Args) < 3 {
	// 		fmt.Println("Usage: kit checkout <branch_name>")
	// 		os.Exit(1)
	// 	}
	// 	name := os.Args[2]
	// 	err := git.CheckoutBranch(name)
	// 	if err != nil {
	// 		fmt.Println("Error checking out branch:", err)
	// 	}
	// default:
	// 	fmt.Println("Unknown command:", cmd)
	// }
}
