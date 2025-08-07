package utils

import (
	"bufio"
	"fmt"
	"kit/pkg"
	"os"
	"path/filepath"
	"strings"
)

func GetIndexEntry(username string) ([]pkg.IndexEntry, error) {
	workspaceDir := filepath.Join("workspaces", username)

	indexPath := workspaceDir + "/.kit/INDEX"
	indexEntry := make([]pkg.IndexEntry, 0)

	file, err := os.Open(indexPath)
	if err != nil {
		if os.IsNotExist(err) {
			return indexEntry, nil
		}
		return nil, fmt.Errorf("failed to open index: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, " ", 3)
		if len(parts) != 3 {
			continue // skip invalid lines
		}
		mode := parts[0]
		hash := parts[1]
		path := parts[2]
		indexEntry = append(indexEntry, pkg.IndexEntry{
			Mode: mode,
			Hash: hash,
			Path: path,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading index: %w", err)
	}
	return indexEntry, nil
}
