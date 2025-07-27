package utils

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"kit/pkg"
	"os"
	"path/filepath"
	"strings"
)

func BuildTree(indexEntry []pkg.IndexEntry) (*pkg.TreeNode, error) {
	root := &pkg.TreeNode{
		Blobs: make(map[string]pkg.BlobEntry),
		Trees: make(map[string]*pkg.TreeNode),
	}

	for _, entry := range indexEntry {
		current := root
		parts := strings.Split(entry.Path, "/")
		for _, part := range parts[:len(parts)-1] {
			if current.Trees[part] == nil {
				current.Trees[part] = &pkg.TreeNode{
					Blobs: make(map[string]pkg.BlobEntry),
					Trees: make(map[string]*pkg.TreeNode),
				}
			}
			current = current.Trees[part]
		}
		fileName := parts[len(parts)-1]
		current.Blobs[fileName] = pkg.BlobEntry{
			Mode: entry.Mode,
			Hash: entry.Hash,
		}
	}
	return root, nil
}

func WriteTree(root *pkg.TreeNode) (string, error) {
	var entries []pkg.TreeEntry

	for name, blob := range root.Blobs {
		entries = append(entries, pkg.TreeEntry{
			Mode: blob.Mode,
			Type: "blob",
			Name: name,
			Hash: blob.Hash,
		})
	}

	for name, tree := range root.Trees {
		hash, err := WriteTree(tree)
		if err != nil {
			return "", err
		}
		entries = append(entries, pkg.TreeEntry{
			Mode: "040000",
			Type: "tree",
			Name: name,
			Hash: hash,
		})
	}

	var treeContent []byte
	for _, entry := range entries {
		line := entry.Mode + " " + entry.Name + "\x00"
		treeContent = append(treeContent, []byte(line)...)

		// Append raw binary SHA1 hash (20 bytes)
		rawHash, err := hex.DecodeString(entry.Hash)
		if err != nil {
			return "", err
		}
		treeContent = append(treeContent, rawHash...)
	}

	//tree object = "tree <len>\0" + content
	header := fmt.Sprintf("tree %d\x00", len(treeContent))
	full := append([]byte(header), treeContent...)

	//Hash and store
	hash := sha1.Sum(full)
	hashHex := hex.EncodeToString(hash[:])

	err := WriteTreeObject(hashHex, full)
	if err != nil {
		return "", err
	}

	return hashHex, nil
}
func WriteTreeObject(hash string, content []byte) error {
	if len(hash) < 2 {
		return fmt.Errorf("invalid hash: %s", hash)
	}

	dir := filepath.Join(".kit", "objects", hash[:2])
	file := filepath.Join(dir, hash[2:])

	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create object dir: %w", err)
	}

	if _, err := os.Stat(file); err == nil {
		return nil
	}

	err = os.WriteFile(file, content, 0644)
	if err != nil {
		return fmt.Errorf("failed to write object: %w", err)
	}

	return nil
}
