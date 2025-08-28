package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Controller struct{}

type Path struct {
	Dir string `json:"dir"`
}

func (c *Controller) FileServer(dir string) http.Handler {
	return http.StripPrefix("/", http.FileServer(http.Dir(dir)))
}

func (c *Controller) GetFile(dir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var path Path

		if err := json.NewDecoder(r.Body).Decode(&path); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "Something went wrong while reading request body , CreateFile function", err)
			return
		}

		if isEmptyDir(w, path.Dir) {
			return
		}

		fullPath := fmt.Sprintf("%s%s", dir, path.Dir)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintln(w, "File not found")
			return
		} else if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "Error accessing file: ", err)
			return
		}

		http.ServeFile(w, r, fullPath)
	}
}

func (c *Controller) CreateFile(dir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var path Path
		if err := json.NewDecoder(r.Body).Decode(&path); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "Something went wrong while reading request body , CreateFile function", err)
			return
		}

		if isEmptyDir(w, path.Dir) {
			return
		}

		fmt.Fprintln(w, path)
	}
}
