package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Triyaambak/nfs/types"
	"github.com/go-chi/chi"
)

type Controller struct{}

func (c *Controller) FileServer(serverConfig *types.ServerConfig) http.Handler {
	dir := (*serverConfig).Dir
	return http.StripPrefix("/", http.FileServer(http.Dir(dir)))
}

func (c *Controller) Fetch(serverConfig *types.ServerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dir := (*serverConfig).Dir

		serverConfig.MU.RLock()
		defer serverConfig.MU.RUnlock()

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
			return
		}

		http.ServeFile(w, r, fullPath)
	}
}

func (c *Controller) Create(serverConfig *types.ServerConfig, isFile bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dir := (*serverConfig).Dir

		serverConfig.MU.Lock()
		defer serverConfig.MU.Unlock()

		w.Header().Set("Content-Type", "application/json")

		path := chi.URLParam(r, "*")

		if isEmptyDir(w, path) {
			return
		}

		fullPath := fmt.Sprintf("%s/%s", dir, path)
		isTaken, err := pathAlreadyTaken(w, fullPath)
		if err != nil {
			return
		}

		if isTaken {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Path nfs/%s already exists", path)
			return
		}

		if !isFile {
			if err := createFolder(w, fullPath); err != nil {
				return
			}
		} else {
			if err := createFolder(w, filepath.Dir(path)); err != nil {
				return
			}
			_, err := os.Create(fullPath)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, `{"error":"failed to create file: %v"}`, err)
				return
			}
		}

		w.WriteHeader(http.StatusCreated)
		if !isTaken {
			fmt.Fprintf(w, "Successfully created folder at path nfs/%s", path)
		} else {
			fmt.Fprintf(w, "Successfully created file at path nfs/%s", path)
		}
	}
}
