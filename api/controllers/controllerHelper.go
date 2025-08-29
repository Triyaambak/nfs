package controllers

import (
	"fmt"
	"os"
)

func isParamEmpty(dir string) error {
	if dir == "" {
		return fmt.Errorf("Bad Request - No dir found in REQUEST body")
	}

	return nil
}

func createFolder(path string) error {
	if err := os.MkdirAll(path, 0777); err != nil {
		return fmt.Errorf("Failed to create parent dirs: %v", err)
	}

	return nil
}

func isFile(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, fmt.Errorf("Something went wrong in isFile function %v", err)
	}
	if !info.Mode().IsRegular() {
		return false, nil
	}

	return true, nil
}

func pathExist(path string) (bool, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}
