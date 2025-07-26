package git

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
)

func HashObject(path string) (string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", err
	}

	if info.IsDir() {
		// 1. Recursively hash directory contents
		treeHash, err := HashTree(path)
		return treeHash, err
	} else {
		// 2. Hash file content as blob
		blobHash, err := HashBlob(path)
		return blobHash, err
	}
}

func HashBlob(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	//"blob <size>\x00<content>"
	header := fmt.Sprintf("blob %d\x00", len(content))
	blob := append([]byte(header), content...)

	hash := sha1.Sum(blob)

	return fmt.Sprintf("%x", hash[:]), nil
}

func HashTree(dirPath string) (string, error) {
	entries, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return "", err
	}

	var treeEntries []byte

	// Sort entries for consistency
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	for _, entry := range entries {
		entryPath := filepath.Join(dirPath, entry.Name())

		var hash string
		if entry.IsDir() {
			hash, err = HashTree(entryPath)
			if err != nil {
				return "", err
			}
			// mode for directory
			treeEntries = append(treeEntries, []byte("40000 "+entry.Name()+"\x00")...)
		} else {
			hash, err = HashBlob(entryPath)
			if err != nil {
				return "", err
			}
			// mode for normal file
			treeEntries = append(treeEntries, []byte("100644 "+entry.Name()+"\x00")...)
		}

		rawHash, err := hexStringToBytes(hash)
		if err != nil {
			return "", err
		}
		treeEntries = append(treeEntries, rawHash...)
	}

	header := fmt.Sprintf("tree %d\x00", len(treeEntries))
	h := sha1.New()
	h.Write([]byte(header))
	h.Write(treeEntries)

	treeHash := fmt.Sprintf("%x", h.Sum(nil))

	// Save tree object to .kit/objects/ (compressed) — implement separately

	return treeHash, nil
}

func hexStringToBytes(s string) ([]byte, error) {
	return hex.DecodeString(s)
}
