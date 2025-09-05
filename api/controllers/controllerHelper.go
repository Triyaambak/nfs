package controllers

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"strings"

	types "github.com/Triyaambak/nfs/types"
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

func renameID(gid, uid int, name, group string) error {
	oldGroupName, err := user.LookupGroupId(strconv.Itoa(gid))
	if err != nil {
		return fmt.Errorf("group with gid %d cannot be found", gid)
	}

	cmd := exec.Command("sudo", "groupmod", "-n", group, oldGroupName.Name)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to rename group: %v, output: %s", err, string(output))
	}

	oldName, err := user.LookupId(strconv.Itoa(uid))
	if err != nil {
		return fmt.Errorf("user with uid %d cannot be found", uid)
	}

	cmd = exec.Command("sudo", "groupmod", "-n", name, oldName.Name)
	output, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to rename name: %v, output: %s", err, string(output))
	}

	return nil

}

func fetchContextData(ctxData *types.ContextDataType) (uid, gid int, name, group string) {

	uid = (*ctxData).Uid
	gid = (*ctxData).Gid
	name = (*ctxData).Name
	group = (*ctxData).Group

	return uid, gid, name, group
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
