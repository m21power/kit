package utils

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"kit/pkg"
	"os"
	"path/filepath"
	"strings"
)

func WriteNewFiles(files map[string]pkg.FileNode, username string) error {
	for _, file := range files {

		if file.Mode == "100644" { // Regular file
			if err := os.MkdirAll(filepath.Dir(file.Path), 0755); err != nil {
				return fmt.Errorf("failed to create directory for file %s: %w", file.Path, err)
			}

			content, err := ReadObject(file.Hash, username)
			if err != nil {
				return fmt.Errorf("failed to read object for file %s: %w", file.Path, err)
			}

			if err := os.WriteFile(file.Path, content, 0644); err != nil {
				return fmt.Errorf("failed to write file %s: %w", file.Path, err)
			}
		}
	}
	return nil
}

func ReadObject(hash string, username string) ([]byte, error) {
	workspaceDir := filepath.Join("workspaces", username)
	path := workspaceDir + "/.kit/objects/" + hash[:2] + "/" + hash[2:]

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read object file %s: %w", path, err)
	}

	r, err := zlib.NewReader(bytes.NewReader(content))
	if err != nil {
		return nil, fmt.Errorf("zlib decompression failed: %w", err)
	}
	defer r.Close()

	decompressed, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read decompressed object: %w", err)
	}

	nullIndex := bytes.IndexByte(decompressed, 0)
	if nullIndex == -1 {
		return nil, fmt.Errorf("invalid object format: missing header separator")
	}
	objectType := string(decompressed[:nullIndex])
	contentOnly := decompressed[nullIndex+1:]

	if strings.HasPrefix(objectType, "blob") {
		return contentOnly, nil
	}

	return nil, fmt.Errorf("unsupported object type: %s", objectType)
}
