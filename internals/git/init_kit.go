package git

import (
	"os"
	"path/filepath"
)

func InitKit() error {
	permissions := 0755
	rootDir := ".kit"

	// Sub folder
	objectDir := filepath.Join(rootDir, "objects")
	refsHeadsDir := filepath.Join(rootDir, "refs", "heads")

	// Files
	headFile := filepath.Join(rootDir, "HEAD")
	indexFile := filepath.Join(rootDir, "INDEX") // staging are

	// Create root directory
	err := os.Mkdir(rootDir, os.FileMode(permissions))
	if err != nil && !os.IsExist(err) {
		return err
	}

	// Create objects/ and refs/heads/
	dirs := []string{objectDir, refsHeadsDir}
	for _, dir := range dirs {
		err := os.MkdirAll(dir, os.FileMode(permissions))
		if err != nil {
			return err
		}
	}

	//Create HEAD file
	headContent := []byte("ref: refs/heads/main\n")
	err = os.WriteFile(headFile, headContent, 0644)
	if err != nil {
		return err
	}

	// Create empty INDEX file
	_, err = os.Create(indexFile)
	if err != nil {
		return err
	}
	return nil
}
