package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ReadIndex(username string) (map[string]string, error) {
	indexPath := fmt.Sprintf("workspaces/%s/.kit/INDEX", username)
	indexMap := make(map[string]string)

	file, err := os.Open(indexPath)
	if err != nil {
		if os.IsNotExist(err) {
			return indexMap, nil // no index yet
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
		hash := parts[1]
		path := parts[2]
		indexMap[path] = hash
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading index: %w", err)
	}

	return indexMap, nil
}
