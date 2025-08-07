package utils

import (
	"fmt"
	"os"
	"strings"
)

func AddOrUpdateIndex(path, hash, username string) error {
	indexMap, err := ReadIndex(username)
	if err != nil {
		return err
	}

	indexMap[path] = hash

	var builder strings.Builder
	for p, h := range indexMap {
		builder.WriteString(fmt.Sprintf("100644 %s %s\n", h, p))
	}
	indexPath := "workspaces/" + username + "/.kit/INDEX"
	err = os.WriteFile(indexPath, []byte(builder.String()), 0644)
	if err != nil {
		return fmt.Errorf("failed to write index: %w", err)
	}

	return nil
}
