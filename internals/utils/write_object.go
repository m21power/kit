package utils

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"os"
	"path/filepath"
)

func WriteObject(filePath, fileType, hashStr string) (string, error) {
	objectDir := filepath.Join(".kit", "objects", hashStr[:2])
	objectPath := filepath.Join(objectDir, hashStr[2:])

	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("read error: %w", err)
	}

	if _, err := os.Stat(objectPath); err == nil {
		return hashStr, nil
	}

	header := fmt.Sprintf("%s %d\x00", fileType, len(content))
	full := append([]byte(header), content...)

	err = WriteZlibCompressedObject(hashStr, full)
	if err != nil {
		return "", err
	}

	return hashStr, nil
}

func WriteZlibCompressedObject(hash string, content []byte) error {
	dir := filepath.Join(".kit", "objects", hash[:2])
	file := filepath.Join(dir, hash[2:])

	if _, err := os.Stat(file); err == nil {
		return nil // already exists
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("mkdir error: %w", err)
	}

	var buf bytes.Buffer
	w := zlib.NewWriter(&buf)
	_, err := w.Write(content)
	if err != nil {
		return fmt.Errorf("zlib compress error: %w", err)
	}
	w.Close()

	return os.WriteFile(file, buf.Bytes(), 0644)
}
