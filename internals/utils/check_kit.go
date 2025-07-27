package utils

import (
	"fmt"
	"os"
)

func CheckKit() error {
	// Check if .kit directory exists
	if _, err := os.Stat(".kit"); os.IsNotExist(err) {
		return fmt.Errorf(".kit directory does not exist. Please initialize the repository using 'kit init'")
	}

	// Check if .kit/HEAD file exists
	headFile := ".kit/HEAD"
	if _, err := os.Stat(headFile); os.IsNotExist(err) {
		return fmt.Errorf("HEAD file does not exist. Please create a branch using 'kit branch <branch_name>'")
	}

	return nil
}
