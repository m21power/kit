package git

import (
	"fmt"
	"kit/internals/utils"
	"kit/pkg"
)

func KitRestore(username string, paths []string, oldTree map[string]pkg.IndexEntry) (map[string]bool, error) {
	normalizedTree, err := utils.GetNormalizedTree(oldTree, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get normalized tree: %w", err)
	}
	stagedMap, err := utils.GetStagedMap(username)
	if err != nil {
		return nil, fmt.Errorf("failed to get staged files: %w", err)
	}
	currentTree, err := utils.GetCurrentTree(username)
	if err != nil {
		return nil, fmt.Errorf("failed to get current tree: %w", err)
	}
	changed := make(map[string]pkg.Status)
	// Check for deleted files
	for path, commitEntry := range normalizedTree {
		if _, inCurrent := currentTree[path]; !inCurrent {
			if stagedEntry, inStaged := stagedMap[path]; inStaged {
				// File was staged (maybe for deletion)
				changed[path] = pkg.Status{
					Staged:  true,
					Message: "deleted (staged)",
					Hash:    stagedEntry.Hash,
				}
			} else {
				// File deleted in working directory but not staged
				changed[path] = pkg.Status{
					Staged:  false,
					Message: "deleted",
					Hash:    commitEntry.Hash,
				}
			}
		}
	}
	result := make(map[string]bool)
	if len(paths) == 0 {
		for path, entry := range changed {
			err := utils.RestoreFile(path, username, entry.Hash)
			if err != nil {
				result[path] = false
			} else {
				result[path] = true
			}
		}
	}
	for _, path := range paths {
		if entry, exists := changed[path]; exists {
			err := utils.RestoreFile(path, username, entry.Hash)
			if err != nil {
				result[path] = false
			} else {
				result[path] = true
			}
		} else {
			return nil, fmt.Errorf("file %s not found in changes", path)
		}

	}
	return result, nil

}
