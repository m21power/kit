package git

import (
	"bytes"
	"fmt"
	"kit/internals/utils"
	"kit/pkg"
	"os"
	"path/filepath"
	"strings"
)

func ResetKit(username string, commitHash string) (pkg.FileSystemItem, error) {
	branch, err := utils.GetHead(username)
	if err != nil {
		return pkg.FileSystemItem{}, fmt.Errorf("failed to get current branch: %w", err)
	}
	file, err := ResetBranch(branch, username, commitHash)
	if err != nil {
		return pkg.FileSystemItem{}, fmt.Errorf("failed to reset branch %s: %w", branch, err)
	}

	return *file, nil
}
func ResetBranch(branch, username, commitHash string) (*pkg.FileSystemItem, error) {
	workspaceDir := filepath.Join("workspaces", username)
	path := workspaceDir + "/.kit/refs/heads/" + branch
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("branch does not exist: %s", branch)
	}
	oldTreeCommit, err := utils.GetCommitTreeHash(branch, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get last commit hash for branch %s: %w", branch, err)
	}

	newTreeCommit, err := GetTreeHash(branch, username, commitHash)
	if err != nil {
		return nil, fmt.Errorf("failed to get commit hash for the given hash %s: %w", branch, err)
	}
	if oldTreeCommit != newTreeCommit {
		newFiles, err := utils.GetFiles(newTreeCommit, "", username)
		if err != nil {
			return nil, fmt.Errorf("failed to get new files for branch %s: %w", branch, err)
		}
		oldFiles, err := utils.GetFiles(oldTreeCommit, "", username)
		if err != nil {
			return nil, fmt.Errorf("failed to get old files for branch %s: %w", branch, err)
		}

		// Handle added/modified files
		for path, newFile := range newFiles {
			oldFile, exists := oldFiles[path]

			if exists && oldFile.Hash == newFile.Hash {
				continue // same file, no need to update
			}
			fl := map[string]pkg.FileNode{
				newFile.Path: newFile,
			}
			// Either new file or modified
			err := utils.WriteNewFiles(fl, username)
			if err != nil {
				return nil, fmt.Errorf("failed to write file %s: %w", newFile.Path, err)
			}
		}

		// Handle deleted files (present in old, missing in new)
		for path := range oldFiles {
			if _, exists := newFiles[path]; !exists {
				err := os.Remove(path)
				if err != nil {
					if removeErr := os.RemoveAll(path); removeErr != nil {
						return nil, fmt.Errorf("failed to delete removed file or directory %s: %w", path, removeErr)
					}
				}
			}
		}

		// Rebuild the .kit/INDEX file

		f, err := os.OpenFile(workspaceDir+"/.kit/INDEX", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to open INDEX for writing: %w", err)
		}
		defer f.Close()

		for _, file := range newFiles {
			if file.Mode == "100644" {
				line := fmt.Sprintf("%s %s %s\n", file.Mode, file.Hash, file.Path)
				if _, err := f.WriteString(line); err != nil {
					return nil, fmt.Errorf("failed to write line to INDEX: %w", err)
				}
			}
		}
	}
	if err := os.WriteFile(path, []byte(commitHash), 0644); err != nil {
		return nil, fmt.Errorf("failed to write branch file: %w", err)
	}
	rootPath := workspaceDir
	structure, err := utils.BuildFileTree(rootPath)
	if err != nil {
		return nil, fmt.Errorf("failed to build file structure: %w", err)
	}

	return &structure, nil

}

func GetTreeHash(branch, username, commitHash string) (string, error) {
	workspaceDir := filepath.Join("workspaces", username)
	// Read the commit object
	objPath := filepath.Join(workspaceDir+"/.kit", "objects", commitHash[:2], commitHash[2:])
	objContent, err := os.ReadFile(objPath)
	if err != nil {
		return "", fmt.Errorf("failed to read commit object: %w", err)
	}

	// Remove the header up to the first null byte
	nullIndex := bytes.IndexByte(objContent, 0)
	if nullIndex == -1 {
		return "", fmt.Errorf("invalid commit object format")
	}
	contentOnly := objContent[nullIndex+1:] // skip the null byte

	lines := strings.Split(string(contentOnly), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "tree ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "tree ")), nil
		}
	}

	return "", fmt.Errorf("tree hash not found in commit object")
}
