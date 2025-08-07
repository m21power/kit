package git

import (
	"fmt"
	"os"
	"path/filepath"
)

func InitKit(username string) error {
	permissions := 0755
	workspaceDir := filepath.Join("workspaces", username)

	rootDir := filepath.Join(workspaceDir, ".kit")
	objectDir := filepath.Join(rootDir, "objects")
	refsHeadsDir := filepath.Join(rootDir, "refs", "heads")

	headFile := filepath.Join(rootDir, "HEAD")
	indexFile := filepath.Join(rootDir, "INDEX")

	// Check if workspace already exists
	if _, err := os.Stat(rootDir); err == nil {
		return fmt.Errorf("kit repository already initialized")
	} else if !os.IsNotExist(err) {
		return err
	}

	// Create workspace root
	err := os.MkdirAll(workspaceDir, os.FileMode(permissions))
	if err != nil {
		return err
	}

	// Create .kit subfolders
	dirs := []string{objectDir, refsHeadsDir}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, os.FileMode(permissions)); err != nil {
			return err
		}
	}

	// Create HEAD file
	headContent := []byte("ref: refs/heads/main")
	if err := os.WriteFile(headFile, headContent, 0644); err != nil {
		return err
	}

	// Create empty INDEX file
	if _, err := os.Create(indexFile); err != nil {
		return err
	}

	return nil
}
