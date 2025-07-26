package utils

import (
	"fmt"
	"os"
	"strings"
)

func AddOrUpdateIndex(path, hash string) error {
	indexMap, err := ReadIndex()
	if err != nil {
		return err
	}

	indexMap[path] = hash

	var builder strings.Builder
	for p, h := range indexMap {
		builder.WriteString(fmt.Sprintf("100644 %s %s\n", h, p))
	}

	err = os.WriteFile(".kit/INDEX", []byte(builder.String()), 0644)
	if err != nil {
		return fmt.Errorf("failed to write index: %w", err)
	}

	return nil
}
