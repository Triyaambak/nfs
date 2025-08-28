package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
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
		path := chi.URLParam(r, "*")

		if isEmptyDir(w, path) {
			return
		}

		fullPath := fmt.Sprintf("%s/%s", dir, path)
		isTaken, err := pathAlreadyTaken(w, fullPath)
		if err != nil {
			return
		}

		if !isTaken {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Path nfs/%s does not exist", path)
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

		fullPath := fmt.Sprintf("%s/%s", dir, path.Dir)
		isTaken, err := pathAlreadyTaken(w, fullPath)
		if err != nil {
			return
		}

		if isTaken {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Path nfs/%s already exists", path)
		}

		fmt.Fprintln(w, path)
	}
}
