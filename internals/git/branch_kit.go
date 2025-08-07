package git

import (
	"fmt"
	"kit/internals/utils"
	"kit/pkg"
	"os"
	"path/filepath"
)

func CreateBranch(name, username string) error {
	workspaceDir := filepath.Join("workspaces", username)

	path := workspaceDir + "/.kit/refs/heads/" + name
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return fmt.Errorf("branch already exists: %s", name)
	}
	branch, err := utils.GetHead(username)
	if err != nil {
		return fmt.Errorf("failed to get current branch: %w", err)
	}
	commit, err := utils.GetLastCommitHash(branch, username)
	if err != nil {
		return fmt.Errorf("failed to get last commit hash for branch %s: %w", name, err)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("failed to create directory for branch: %w", err)
	}
	if err := os.WriteFile(path, []byte(commit), 0644); err != nil {
		return fmt.Errorf("failed to create branch file: %w", err)
	}
	fmt.Printf("Branch '%s' created successfully\n", name)
	return nil
}

func CheckoutBranch(name, username string) (*pkg.FileSystemItem, error) {
	workspaceDir := filepath.Join("workspaces", username)
	path := workspaceDir + "/.kit/refs/heads/" + name
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("branch does not exist: %s", name)
	}
	oldBranch, err := utils.GetHead(username)
	if err != nil {
		return nil, fmt.Errorf("failed to get current branch: %w", err)
	}
	if oldBranch == name {
		fmt.Printf("Already on branch '%s'\n", name)
		return nil, fmt.Errorf("already on branch '%s'", name)
	}
	oldTreeCommit, err := utils.GetCommitTreeHash(oldBranch, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get last commit hash for branch %s: %w", oldBranch, err)
	}

	newTreeCommit, err := utils.GetCommitTreeHash(name, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get last commit hash for branch %s: %w", name, err)
	}
	if oldTreeCommit != newTreeCommit {
		newFiles, err := utils.GetFiles(newTreeCommit, "", username)
		if err != nil {
			return nil, fmt.Errorf("failed to get new files for branch %s: %w", name, err)
		}
		oldFiles, err := utils.GetFiles(oldTreeCommit, "", username)
		if err != nil {
			return nil, fmt.Errorf("failed to get old files for branch %s: %w", name, err)
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
				if err != nil && !os.IsNotExist(err) {
					return nil, fmt.Errorf("failed to delete removed file %s: %w", path, err)
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

	if err := os.WriteFile(workspaceDir+"/.kit/HEAD", []byte("ref: refs/heads/"+name), 0644); err != nil {
		return nil, fmt.Errorf("failed to checkout branch: %w", err)
	}
	rootPath := workspaceDir
	structure, err := utils.BuildFileTree(rootPath)
	if err != nil {
		return nil, fmt.Errorf("failed to build file structure: %w", err)
	}

	return &structure, nil

}
