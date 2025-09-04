package controllers

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

func splitPath(path string) (srcPath, destPath string, err error) {
	runes := []rune(path)
	idx := -1
	for i, c := range runes {
		if c == ' ' {
			if idx != -1 {
				return "", "", fmt.Errorf("%s contains multiple '-' , please keep only one whitespace", path)
			}
			idx = i
		}
	}

	if idx == -1 {
		return "", "", fmt.Errorf("%s does not contain even one whitespace to distinguish between src and dest path", path)
	}

	srcPath = path[:idx]
	destPath = path[idx+1:]

	return srcPath, destPath, nil
}

func getWriteMode(urlParam string) (body, path string, isAppend bool, err error) {
	startIndex := -1
	endIndex := -1

	for i, c := range urlParam {
		if c == '>' {
			if startIndex == -1 {
				startIndex = i
				endIndex = i
			} else {
				endIndex = i
			}
		}
	}

	if startIndex == -1 {
		return "", "", false, errors.New("url does not contain > to specify wether to write or append, please specify")
	}

	if startIndex == endIndex {
		isAppend = false
	} else {
		isAppend = true
	}

	body = strings.TrimSpace(urlParam[:startIndex])

	path = strings.TrimSpace(urlParam[endIndex+1:])

	return body, path, isAppend, nil
}

func isParamEmpty(dir string) error {
	if dir == "" {
		return fmt.Errorf("bad Request - No dir found in REQUEST body")
	}

	return nil
}

func createFolder(path string) error {
	if err := os.MkdirAll(path, 0777); err != nil {
		return fmt.Errorf("failed to create parent dirs: %v", err)
	}

	return nil
}

func isFile(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, fmt.Errorf("something went wrong in isFile function %v", err)
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
