package git

import (
	"fmt"
	util "kit/internals/utils"
	"os"
	"path/filepath"
	"strings"
)

func AddKit(path, username string) ([]string, error) {
	isDir, err := util.IsDir(path)
	if err != nil {
		return nil, err
	}

	if *isDir {
		return ProcessFolder(path, username)
	}

	staged, err := handleFileSaving(path, username)
	if err != nil {
		return nil, err
	}

	if staged != "" {
		return []string{staged}, nil
	}
	return []string{}, nil
}

func handleFileSaving(path, username string) (string, error) {
	hash, err := util.HashBlob(path)
	if err != nil {
		return "", fmt.Errorf("failed to hash file %s: %w", path, err)
	}

	indexMap, err := util.ReadIndex(username)
	if err != nil {
		return "", fmt.Errorf("failed to read index: %w", err)
	}

	// Skip if already in index with same hash
	if existingHash, ok := indexMap[path]; ok && existingHash == hash {
		return "", nil
	}

	_, err = util.WriteObject(path, "blob", hash, username)
	if err != nil {
		return "", fmt.Errorf("failed to write blob object for %s: %w", path, err)
	}

	err = util.AddOrUpdateIndex(path, hash, username)
	if err != nil {
		return "", fmt.Errorf("failed to update index: %w", err)
	}
	//  "workspaces/mesay/mesay/src/main.go"
	p := strings.Split(path, "/")
	res := strings.Join(p[2:], "/")
	return res, nil
}

func ProcessFolder(root string, username string) ([]string, error) {
	var stagedFiles []string

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
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
			files, err := AddKit(path, username)
			if err != nil {
				return err
			}
			stagedFiles = append(stagedFiles, files...)
		}

		return nil
	})

	return stagedFiles, err
}
