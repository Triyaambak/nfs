package controllers

import (
	"fmt"
	"net/http"
	"os"
)

func isEmptyDir(w http.ResponseWriter, dir string) bool {

	if dir == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Bad Request - No dir found in REQUEST body")
		return true
	}

	return false
}

func createFolder(w http.ResponseWriter, path string) error {
	if err := os.MkdirAll(path, 0777); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error":"failed to create parent dirs: %v"}`, err)
		return err
	}
	return nil
}

func pathAlreadyTaken(w http.ResponseWriter, path string) (bool, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error accessing file: ", err)
		return false, err
	}

	return true, nil
}
