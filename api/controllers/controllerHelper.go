package controllers

import (
	"fmt"
	"net/http"
)

func isEmptyDir(w http.ResponseWriter, dir string) bool {

	if dir == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Bad Request - No dir found in REQUEST body")
		return true
	}

	return false
}
