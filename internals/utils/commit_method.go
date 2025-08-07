package utils

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func GetHead(username string) (string, error) {
	workspaceDir := filepath.Join("workspaces", username)
	headPath := filepath.Join(workspaceDir, ".kit/HEAD")
	content, err := os.ReadFile(headPath)
	if err != nil {
		return "", err
	}

	// ref: refs/heads/main
	lines := strings.Split(string(content), " ")
	if len(lines) < 2 {
		return "", fmt.Errorf("invalid HEAD format")
	}
	temp := strings.Split(lines[len(lines)-1], "/")
	return temp[len(temp)-1], nil
}

func GetCommitTreeHash(branch, username string) (string, error) {
	workspaceDir := filepath.Join("workspaces", username)
	rootDir := filepath.Join(workspaceDir, ".kit/refs/heads")
	commitPath := fmt.Sprintf("%s/%s", rootDir, branch)
	content, err := os.ReadFile(commitPath)
	if errors.Is(err, os.ErrNotExist) {
		return "", nil // first commit; no parent
	} else if err != nil {
		return "", fmt.Errorf("failed to read commit hash for branch %s: %w", branch, err)
	}

	commitHash := strings.TrimSpace(string(content))
	if commitHash == "" {
		return "", fmt.Errorf("no commit hash found for branch %s", branch)
	}

	// Read the commit object
	objPath := filepath.Join(workspaceDir+"/.kit", "objects", commitHash[:2], commitHash[2:])
	objContent, err := os.ReadFile(objPath)
	if err != nil {
		return "", fmt.Errorf("failed to read commit object: %w", err)
	}

	// Remove the header up to the first null byte
	nullIndex := bytes.IndexByte(objContent, 0)
	if nullIndex == -1 {
		return "", fmt.Errorf("invalid commit object format")
	}
	contentOnly := objContent[nullIndex+1:] // skip the null byte

	lines := strings.Split(string(contentOnly), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "tree ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "tree ")), nil
		}
	}

	return "", fmt.Errorf("tree hash not found in commit object")
}
func GetLastCommitHash(branch, username string) (string, error) {
	fmt.Println("Getting last commit hash for branch:", branch)
	workspaceDir := filepath.Join("workspaces", username)

	commitPath := fmt.Sprintf(workspaceDir+"/.kit/refs/heads/%s", branch)
	content, err := os.ReadFile(commitPath)
	if errors.Is(err, os.ErrNotExist) {
		return "", nil // first commit; no parent
	} else if err != nil {
		return "", fmt.Errorf("failed to read commit hash for branch %s: %w", branch, err)
	}

	commitHash := strings.TrimSpace(string(content))
	if commitHash == "" {
		return "", fmt.Errorf("no commit hash found for branch %s", branch)
	}

	return commitHash, nil
}
func WriteCommit(branch, newCommit, oldCommit, message, username string) error {
	author := username + " Author"
	email := "author@gmail.com"

	now := time.Now()
	timestamp := now.Unix()
	timezone := now.Format("-0700")

	commitContent := fmt.Sprintf("tree %s\n", newCommit)
	if oldCommit != "" {
		commitContent += fmt.Sprintf("parent %s\n", oldCommit)
	}
	commitContent += fmt.Sprintf("author %s <%s> %d %s\n", author, email, timestamp, timezone)
	commitContent += fmt.Sprintf("committer %s <%s> %d %s\n", author, email, timestamp, timezone)
	commitContent += "\n" + message + "\n"

	commitHash, err := WriteC(commitContent, "commit", username)
	fmt.Println("commit hash: ", commitHash)
	if err != nil {
		return fmt.Errorf("failed to write commit object: %w", err)
	}
	workspaceDir := filepath.Join("workspaces", username)

	branchPath := fmt.Sprintf(workspaceDir+"/.kit/refs/heads/%s", branch)
	err = os.WriteFile(branchPath, []byte(commitHash), 0644)
	if err != nil {
		return fmt.Errorf("failed to update branch reference: %w", err)
	}

	return nil
}

func WriteC(content string, objType string, username string) (string, error) {
	header := fmt.Sprintf("%s %d\u0000", objType, len(content))
	full := []byte(header + content)

	hash := sha1.Sum(full)
	hashStr := hex.EncodeToString(hash[:])
	workspaceDir := filepath.Join("workspaces", username)

	dir := filepath.Join(workspaceDir+"/.kit", "objects", hashStr[:2])
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create object dir: %w", err)
	}

	objPath := filepath.Join(dir, hashStr[2:])
	if err := os.WriteFile(objPath, full, 0644); err != nil {
		return "", fmt.Errorf("failed to write object: %w", err)
	}

	return hashStr, nil
}
