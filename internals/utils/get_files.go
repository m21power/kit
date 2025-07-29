package utils

import (
	"bytes"
	"compress/zlib"
	"encoding/hex"
	"fmt"
	"io"
	"kit/pkg"
	"os"
	"strings"
)

func GetFiles(hash string, prefix string) (map[string]pkg.FileNode, error) {
	files := make(map[string]pkg.FileNode)

	// Load and decompress the object
	data, objectType, err := loadObjectFromHash(hash)

	if err != nil {
		return nil, err
	}

	if objectType != "tree" {
		return nil, fmt.Errorf("object %s is not a tree", hash)
	}

	i := 0
	for i < len(data) {
		nullIdx := bytes.IndexByte(data[i:], 0)
		if nullIdx == -1 || i+nullIdx+20 > len(data) {
			return nil, fmt.Errorf("invalid tree format")
		}

		entry := data[i : i+nullIdx]
		parts := bytes.SplitN(entry, []byte(" "), 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid tree entry")
		}
		mode := string(parts[0])
		name := string(parts[1])
		fullPath := prefix + name

		hashStart := i + nullIdx + 1
		hash := data[hashStart : hashStart+20]
		hashStr := hex.EncodeToString(hash)

		i = hashStart + 20

		// Determine type (file/tree)
		_, objType, err := loadObjectFromHash(hashStr)

		if err != nil {
			return nil, err
		}

		node := pkg.FileNode{
			Path: fullPath,
			Hash: hashStr,
			Mode: mode,
			Type: objType,
		}
		files[fullPath] = node

		// Recursively get children if it's a tree (directory)
		if objType == "tree" {
			children, err := GetFiles(hashStr[:2]+"/"+hashStr[2:], fullPath+"/")
			if err != nil {
				return nil, err
			}
			for k, v := range children {
				files[k] = v
			}
		}
	}

	return files, nil
}
func loadObjectFromHash(fullHash string) ([]byte, string, error) {
	dir := fullHash[:2]
	file := fullHash[2:]
	path := ".kit/objects/" + dir + "/" + file

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read object at %s: %w", path, err)
	}

	r, err := zlib.NewReader(bytes.NewReader(content))
	if err != nil {
		return nil, "", err
	}
	defer r.Close()

	decompressed, err := io.ReadAll(r)
	if err != nil {
		return nil, "", err
	}

	nullIdx := bytes.IndexByte(decompressed, 0)
	if nullIdx == -1 {
		return nil, "", fmt.Errorf("invalid object header")
	}

	header := string(decompressed[:nullIdx])
	contentOnly := decompressed[nullIdx+1:]
	parts := strings.SplitN(header, " ", 2)

	return contentOnly, parts[0], nil
}
