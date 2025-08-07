package utils

import (
	"bytes"
	"fmt"
	"kit/pkg"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func GetLogs(commitHash, username string) (pkg.CommitObject, error) {
	workspaceDir := filepath.Join("workspaces", username)

	objPath := filepath.Join(workspaceDir+"/.kit", "objects", commitHash[:2], commitHash[2:])
	objContent, err := os.ReadFile(objPath)
	if err != nil {
		return pkg.CommitObject{}, fmt.Errorf("failed to read commit object: %w", err)
	}

	nullIndex := bytes.IndexByte(objContent, 0)
	if nullIndex == -1 {
		return pkg.CommitObject{}, fmt.Errorf("invalid commit object format")
	}
	contentOnly := objContent[nullIndex+1:]
	lines := strings.Split(string(contentOnly), "\n")

	var commit pkg.CommitObject
	commit.Hash = commitHash

	messageStarted := false
	var messageLines []string

	for _, line := range lines {
		if messageStarted {
			messageLines = append(messageLines, line)
			continue
		}

		switch {
		case strings.HasPrefix(line, "tree "):
			// Skip tree hash
		case strings.HasPrefix(line, "parent "):
			commit.Parent = strings.TrimSpace(strings.TrimPrefix(line, "parent "))
		case strings.HasPrefix(line, "author "):
			// Format: author Name <email> 1234567890 +0300
			fields := strings.Fields(line)
			if len(fields) >= 6 {
				commit.Author = strings.Join(fields[1:len(fields)-4], " ")
				commit.Email = strings.Trim(fields[len(fields)-4], "<>")
				secs, err := strconv.ParseInt(fields[len(fields)-2], 10, 64)
				if err == nil {
					commit.Date = time.Unix(secs, 0)
				}
			}
		case line == "":
			messageStarted = true
		}
	}

	commit.Message = strings.Join(messageLines, "\n")
	return commit, nil
}
