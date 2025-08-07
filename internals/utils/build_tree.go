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

func WriteTree(root *pkg.TreeNode, username string) (string, error) {
	var entries []pkg.TreeEntry

	for name, blob := range root.Blobs {
		entries = append(entries, pkg.TreeEntry{
			Mode: "100644",
			Type: "blob",
			Name: name,
			Hash: blob.Hash,
		})
	}

	for name, tree := range root.Trees {
		hash, err := WriteTree(tree, username)
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

	err := WriteZlibCompressedObject(hashHex, username, full)
	if err != nil {
		return "", err
	}

	return hashHex, nil
}

func BuildFileTree(root string) (pkg.FileSystemItem, error) {
	info, err := os.Stat(root)
	if err != nil {
		return pkg.FileSystemItem{}, err
	}

	item := pkg.FileSystemItem{
		ID:   root, // optional
		Name: info.Name(),
		Path: root,
		Type: "folder",
	}

	entries, err := os.ReadDir(root)
	if err != nil {
		return item, err
	}

	for _, entry := range entries {
		// Skip internal .kit folder
		if entry.Name() == ".kit" {
			continue
		}

		childPath := filepath.Join(root, entry.Name())

		if entry.IsDir() {
			child, err := BuildFileTree(childPath)
			if err != nil {
				return item, err
			}
			item.Children = append(item.Children, child)
		} else {
			content, err := os.ReadFile(childPath)
			if err != nil {
				return item, err
			}
			child := pkg.FileSystemItem{
				ID:      childPath, // optional
				Name:    entry.Name(),
				Path:    childPath,
				Type:    "file",
				Content: string(content),
			}
			item.Children = append(item.Children, child)
		}
	}

	return item, nil
}
