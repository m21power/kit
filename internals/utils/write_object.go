package utils

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"kit/pkg"
	"os"
	"path/filepath"
)

func WriteObject(filePath, fileType, hashStr, username string) (string, error) {
	workspaceDir := filepath.Join("workspaces", username)
	rootDir := filepath.Join(workspaceDir, ".kit")
	objectDir := filepath.Join(rootDir, "objects", hashStr[:2])
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

	err = WriteZlibCompressedObject(hashStr, username, full)
	if err != nil {
		return "", err
	}

	return hashStr, nil
}

func WriteZlibCompressedObject(hash, username string, content []byte) error {
	workspaceDir := filepath.Join("workspaces", username)
	rootDir := filepath.Join(workspaceDir, ".kit")

	dir := filepath.Join(rootDir, "objects", hash[:2])
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

func WriteStructure(basePath string, item pkg.FileSystemItem) error {
	targetPath := filepath.Join(basePath, item.Path)

	// If the path exists, remove it first
	if _, err := os.Stat(targetPath); err == nil {
		if err := os.RemoveAll(targetPath); err != nil {
			return fmt.Errorf("failed to remove existing path %s: %w", targetPath, err)
		}
	}

	if item.Type == "folder" {
		if err := os.MkdirAll(targetPath, 0755); err != nil {
			return err
		}
		for _, child := range item.Children {
			if err := WriteStructure(basePath, child); err != nil {
				return err
			}
		}
	} else if item.Type == "file" {
		// Ensure parent folder exists
		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return err
		}
		if err := os.WriteFile(targetPath, []byte(item.Content), 0644); err != nil {
			return err
		}
	}

	return nil
}
