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
	"strings"
)

func StatusKit(hash string, fullPath string, visited map[string]bool) (map[string]pkg.IndexEntry, error) {
	if visited[hash] {
		return nil, nil // Already processed
	}
	visited[hash] = true
	var result = make(map[string]pkg.IndexEntry)
	path := ".kit/objects/" + hash[:2]
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
			subEntries, err := StatusKit(entry.Hash, fullPath+"/"+entry.Name, visited)
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

func IsChanged(oldTree map[string]pkg.IndexEntry, root string) (map[string]pkg.Status, error) {
	changed := make(map[string]pkg.Status)
	staged, err := utils.GetIndexEntry()
	if err != nil {
		return nil, fmt.Errorf("failed to get index entry: %w", err)
	}
	err = filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(p, "main.go") {
			return nil
		}
		// Skip .kit directory entirely
		if info.IsDir() && filepath.Base(p) == ".kit" {
			return filepath.SkipDir
		}

		// Skip non-files (just in case)
		if info.IsDir() {
			return nil
		}

		// Compute blob hash of file
		hash, err := HashBlob(p)
		if err != nil {
			return err
		}

		// Get path relative to root
		relPath, err := filepath.Rel(root, p)
		if err != nil {
			return err
		}

		entry, exists := oldTree["/"+relPath]
		if !exists {
			changed[relPath] = pkg.Status{
				Staged:  false,
				Message: "created",
				Hash:    hash,
			}
		}
		if exists && entry.Hash != hash {
			changed[relPath] = pkg.Status{
				Staged:  false,
				Message: "modified",
				Hash:    hash,
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	for _, entry := range staged {
		change, exists := changed[entry.Path]
		if exists {
			if change.Hash == entry.Hash {
				changed[entry.Path] = pkg.Status{
					Staged:  true,
					Message: change.Message,
					Hash:    entry.Hash,
				}
			}

		} else {

		}
	}

	return changed, nil
}
