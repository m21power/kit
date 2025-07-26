package utils

import "os"

func IsDir(path string) (*bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if info.IsDir() {
		result := true
		return &result, nil
	}
	result := false
	return &result, nil
}
