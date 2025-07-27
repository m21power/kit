package utils

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func WriteObject(filePath, fileType string) (string, error) {
	log.Println("Writing blob object for", filePath)
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	header := fmt.Sprintf("%s %d\x00", fileType, len(content))
	blob := append([]byte(header), content...)

	hash := sha1.Sum(blob)
	hashStr := fmt.Sprintf("%x", hash[:])

	objectDir := filepath.Join(".kit", "objects", hashStr[:2])
	objectPath := filepath.Join(objectDir, hashStr[2:])

	// If object already exists, skip writing
	if _, err := os.Stat(objectPath); err == nil {
		return hashStr, nil
	}

	err = os.MkdirAll(objectDir, 0755)
	if err != nil {
		return "", fmt.Errorf("failed to create object dir: %w", err)
	}

	var buf bytes.Buffer
	w := zlib.NewWriter(&buf)
	_, err = w.Write(blob)
	if err != nil {
		return "", fmt.Errorf("compression write error: %w", err)
	}
	w.Close()

	err = os.WriteFile(objectPath, buf.Bytes(), 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write object: %w", err)
	}

	return hashStr, nil
}
