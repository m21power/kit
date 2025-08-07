package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func RestoreFile(path string, username string, hash string) error {
	workspaceDir := filepath.Join("workspaces", username)
	fullPath := filepath.Join(workspaceDir, path)

	content, err := ReadObject(hash, username)
	if err != nil {
		return fmt.Errorf("failed to read object for file %s: %w", path, err)
	}
	fmt.Println("print content")
	fmt.Println(string(content))
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory for file %s: %w", path, err)
	}
	if err := os.WriteFile(fullPath, content, 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", path, err)
	}
	// Check if the file exists
	if _, err := os.Stat(fullPath); err != nil {
		return fmt.Errorf("file '%s' does not exist after restoration: %w", path, err)
	}

	return nil
}
