package git

import (
	"bytes"
	"compress/zlib"
	"encoding/hex"
	"fmt"
	"io"
	"kit/internals/utils"
	"kit/pkg"
	"log"
	"os"
	"path/filepath"
)

func StatusKit(hash string, fullPath, username string, visited map[string]bool) (map[string]pkg.IndexEntry, error) {
	if hash == "" {
		return nil, nil
	}
	workspaceDir := filepath.Join("workspaces", username)
	if visited[hash] {
		return nil, nil // Already processed
	}
	visited[hash] = true
	var result = make(map[string]pkg.IndexEntry)
	path := workspaceDir + "/.kit/objects/" + hash[:2]
	path = filepath.Join(path, hash[2:])
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("directory does not exist: %s", path)
	}
	content, err := os.ReadFile(path)

	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", path, err)
	}

	r, err := zlib.NewReader(bytes.NewReader(content))
	if err != nil {
		return nil, fmt.Errorf("failed to decompress: %w", err)
	}
	defer r.Close()

	decompressed, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("read decompressed error: %w", err)
	}

	entriesData, err := StripHeader(decompressed)
	if err != nil {
		log.Fatal(err)
	}

	entries, err := ParseTreeObject(entriesData) // your parser
	if err != nil {
		return nil, fmt.Errorf("failed to parse tree object: %w", err)
	}
	for _, entry := range entries {

		if entry.Type == "blob" {
			result[fullPath+"/"+entry.Name] = pkg.IndexEntry{
				Mode: entry.Mode,
				Hash: entry.Hash,
				Path: fullPath + "/" + entry.Name,
			}
		} else if entry.Type == "tree" {
			subEntries, err := StatusKit(entry.Hash, fullPath+"/"+entry.Name, username, visited)
			if err != nil {
				return nil, fmt.Errorf("failed to get sub entries: %w", err)
			}
			for k, v := range subEntries {
				result[k] = v
			}
		}

	}

	return result, nil
}
func ParseTreeObject(data []byte) ([]pkg.TreeEntry, error) {
	var entries []pkg.TreeEntry
	i := 0

	for i < len(data) {
		// 1. Read mode (e.g., 100644)
		modeEnd := bytes.IndexByte(data[i:], ' ')
		if modeEnd == -1 {
			return nil, fmt.Errorf("invalid mode format")
		}
		mode := string(data[i : i+modeEnd])
		i += modeEnd + 1

		// 2. Read filename
		nameEnd := bytes.IndexByte(data[i:], 0)
		if nameEnd == -1 {
			return nil, fmt.Errorf("invalid name format")
		}
		name := string(data[i : i+nameEnd])
		i += nameEnd + 1

		// 3. Read 20-byte SHA1 (raw binary)
		if i+20 > len(data) {
			return nil, fmt.Errorf("incomplete SHA1 hash")
		}
		rawHash := data[i : i+20]
		hash := hex.EncodeToString(rawHash)
		i += 20

		entries = append(entries, pkg.TreeEntry{
			Mode: mode,
			Name: name,
			Hash: hash,
			Type: getTypeFromMode(mode),
		})
	}

	return entries, nil
}

func getTypeFromMode(mode string) string {
	if mode == "040000" {
		return "tree"
	}
	return "blob"
}
func StripHeader(data []byte) ([]byte, error) {
	nullIndex := bytes.IndexByte(data, 0)
	if nullIndex == -1 {
		return nil, fmt.Errorf("invalid tree header (no null byte)")
	}
	return data[nullIndex+1:], nil
}

func IsChanged(oldTree map[string]pkg.IndexEntry, root string, username string) (map[string]pkg.Status, error) {

	changed := make(map[string]pkg.Status)

	normalizedTree, err := utils.GetNormalizedTree(oldTree, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get normalized tree: %w", err)
	}
	stagedMap, err := utils.GetStagedMap(username)
	if err != nil {
		return nil, fmt.Errorf("failed to get staged files: %w", err)
	}
	currentTree, err := utils.GetCurrentTree(username)
	if err != nil {
		return nil, fmt.Errorf("failed to get current tree: %w", err)
	}
	// 1- we have normalizedTree which is the commit tree
	// 2- we have stagedMap which is the staged entries
	// 3- we have currentTree which is the current working directory state
	// 4- we need to compare these three to determine the status of each file
	for path, entry := range currentTree {
		stagedEntry, inStaged := stagedMap[path]
		_, inCommit := normalizedTree[path]

		switch {
		case !inStaged && !inCommit:
			// New file, not staged
			changed[path] = pkg.Status{
				Staged:  false,
				Message: "created",
				Hash:    entry.Hash,
			}

		case inStaged && stagedEntry.Hash != entry.Hash:
			// Modified after staging
			changed[path] = pkg.Status{
				Staged:  true,
				Message: "modified",
				Hash:    entry.Hash,
			}

		case inStaged && stagedEntry.Hash == entry.Hash && !inCommit:
			// New file staged but not committed yet
			changed[path] = pkg.Status{
				Staged:  true,
				Message: "staged",
				Hash:    entry.Hash,
			}
		}
	}

	// Check for deleted files
	for path, commitEntry := range normalizedTree {
		if _, inCurrent := currentTree[path]; !inCurrent {
			if stagedEntry, inStaged := stagedMap[path]; inStaged {
				// File was staged (maybe for deletion)
				changed[path] = pkg.Status{
					Staged:  true,
					Message: "deleted (staged)",
					Hash:    stagedEntry.Hash,
				}
			} else {
				// File deleted in working directory but not staged
				changed[path] = pkg.Status{
					Staged:  false,
					Message: "deleted",
					Hash:    commitEntry.Hash,
				}
			}
		}
	}
	return changed, nil
}
