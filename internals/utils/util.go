package utils

import (
	"crypto/sha1"
	"fmt"
	"kit/pkg"
	"os"
	"path/filepath"
	"strings"
)

func GetCurrentTree(username string) (map[string]pkg.IndexEntry, error) {
	currentTree := make(map[string]pkg.IndexEntry)
	workspaceDir := filepath.Join("workspaces", username)
	absRoot := filepath.Join(workspaceDir, ".")

	// --- Walk working directory ---
	err := filepath.Walk(absRoot, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip .kit dir
		if info.IsDir() && filepath.Base(p) == ".kit" {
			return filepath.SkipDir
		}
		if info.IsDir() {
			return nil
		}

		hash, err := HashBlob(p)
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(workspaceDir, p)
		currentTree[relPath] = pkg.IndexEntry{
			Mode: "100644",
			Hash: hash,
			Path: relPath,
		}
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return currentTree, nil
}

func HashBlob(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	//"blob <size>\x00<content>"
	header := fmt.Sprintf("blob %d\x00", len(content))
	blob := append([]byte(header), content...)

	hash := sha1.Sum(blob)

	return fmt.Sprintf("%x", hash[:]), nil
}

func GetNormalizedTree(oldTree map[string]pkg.IndexEntry, username string) (map[string]pkg.IndexEntry, error) {
	normalizedTree := make(map[string]pkg.IndexEntry)
	workspaceDir := filepath.Join("workspaces", username)

	for path, entry := range oldTree {
		clean := strings.TrimPrefix(path, "/") // remove leading slash
		rel, _ := filepath.Rel(workspaceDir, clean)
		normalizedTree[rel] = pkg.IndexEntry{
			Mode: entry.Mode,
			Hash: entry.Hash,
			Path: rel,
		}
	}
	return normalizedTree, nil
}

func GetStagedMap(username string) (map[string]pkg.IndexEntry, error) {
	workspaceDir := filepath.Join("workspaces", username)
	stagedEntries, err := GetIndexEntry(username)
	if err != nil {
		return nil, fmt.Errorf("failed to get index entry: %w", err)
	}
	stagedMap := make(map[string]pkg.IndexEntry)
	for _, e := range stagedEntries {
		rel, _ := filepath.Rel(workspaceDir, e.Path)
		stagedMap[rel] = pkg.IndexEntry{
			Mode: e.Mode,
			Hash: e.Hash,
			Path: rel,
		}
	}
	return stagedMap, nil
}
func GetAllDir() ([]string, error) {
	workspaceDir := "workspaces"
	var result []string
	if _, err := os.Stat(workspaceDir); os.IsNotExist(err) {
		if err := os.Mkdir(workspaceDir, 0755); err != nil {
			return nil, err
		}
	}

	entries, err := os.ReadDir(workspaceDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read workspaces directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			result = append(result, entry.Name())
		}
	}
	return result, nil
}
