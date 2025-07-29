package git

import (
	"fmt"
	util "kit/internals/utils"
	"os"
	"path/filepath"
	"strings"
)

func AddKit(path string) error {
	isDir, err := util.IsDir(path)
	if err != nil {
		return err
	}

	if *isDir {
		return ProcessFolder(path)
	}
	if strings.HasSuffix(path, "main.go") {
		return nil
	}
	return handleFileSaving(path)
}

func handleFileSaving(path string) error {
	hash, err := HashBlob(path)
	if err != nil {
		return fmt.Errorf("failed to hash file %s: %w", path, err)
	}

	indexMap, err := util.ReadIndex()
	if err != nil {
		return fmt.Errorf("failed to read index: %w", err)
	}

	// Skip if already in index with same hash
	if existingHash, ok := indexMap[path]; ok && existingHash == hash {
		return nil
	}

	// Write the blob object
	_, err = util.WriteObject(path, "blob", hash)
	if err != nil {
		return fmt.Errorf("failed to write blob object for %s: %w", path, err)
	}

	// Update index
	err = util.AddOrUpdateIndex(path, hash)
	if err != nil {
		return fmt.Errorf("failed to update index: %w", err)
	}

	return nil
}

func ProcessFolder(root string) error {
	return filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip .kit directory
		if d.IsDir() && filepath.Base(path) == ".kit" {
			return filepath.SkipDir
		}

		// Skip anything inside .kit
		if strings.Contains(path, string(os.PathSeparator)+".kit"+string(os.PathSeparator)) {
			return nil
		}

		if !d.IsDir() {
			return AddKit(path)
		}

		return nil
	})
}
