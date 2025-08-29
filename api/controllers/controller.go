package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	types "github.com/Triyaambak/nfs/types"

	"github.com/go-chi/chi"
)

type Controller struct{}

func (c *Controller) FileServer(serverConfig *types.ServerConfig) http.Handler {
	dir := (*serverConfig).Dir
	return http.StripPrefix("/", http.FileServer(http.Dir(dir)))
}

func (c *Controller) LS(serverConig *types.ServerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dir := (*serverConig).Dir

		serverConig.MU.RLock()
		defer serverConig.MU.Unlock()

		path := chi.URLParam(r, "*")
		if err := isParamEmpty(path); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		fullPath := filepath.Join(dir, path)
		isTaken, err := pathExist(fullPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		if !isTaken {
			http.Error(w, fmt.Sprintf("Path %s does not exist", path), http.StatusBadRequest)
			return
		}

		info, _ := os.Stat(fullPath)

		if info.IsDir() {
			entries, err := os.ReadDir(fullPath)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to read directory: %v", err), http.StatusInternalServerError)
				return
			}

			for _, e := range entries {
				if e.IsDir() {
					fmt.Fprintf(w, "%s/\n", e.Name())
				} else {
					fmt.Fprintf(w, "%s\n", e.Name())
				}
			}
		} else {
			// It's a file → just print its name
			fmt.Fprintf(w, "%s\n", info.Name())
		}
	}
}

func (c *Controller) Cat(serverConfig *types.ServerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dir := (*serverConfig).Dir

		serverConfig.MU.RLock()
		defer serverConfig.MU.RUnlock()

		path := chi.URLParam(r, "*")

		if err := isParamEmpty(path); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		fullPath := filepath.Join(dir, path)
		if filepath.Ext(fullPath) == "" {
			http.Error(w, fmt.Sprintln("Bad Request , no file extension mentioned"), http.StatusBadRequest)
			return
		}
		isTaken, err := pathExist(fullPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		if !isTaken {
			http.Error(w, fmt.Sprintf("Path %s does not exist", path), http.StatusBadRequest)
			return
		}

		isF, err := isFile(fullPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if !isF {
			http.Error(w, fmt.Sprintf("Traget is not a file but a folder : %s", path), http.StatusBadRequest)
			return
		}

		http.ServeFile(w, r, fullPath)
	}
}

func (c *Controller) Create(serverConfig *types.ServerConfig, isFolder bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dir := (*serverConfig).Dir

		serverConfig.MU.Lock()
		defer serverConfig.MU.Unlock()

		w.Header().Set("Content-Type", "application/json")

		path := chi.URLParam(r, "*")

		if err := isParamEmpty(path); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		fullPath := filepath.Join(dir, path)
		isPathTaken, err := pathExist(fullPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		if isPathTaken {
			http.Error(w, fmt.Sprintf("Path %s already exists", path), http.StatusBadRequest)
			return
		}

		if isFolder {
			if err := createFolder(fullPath); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {

			if filepath.Ext(fullPath) == "" {
				http.Error(w, fmt.Sprintln("Bad Request , no file extension mentioned"), http.StatusBadRequest)
				return
			}

			doesFolderExist, err := pathExist(filepath.Dir(fullPath))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if !doesFolderExist {
				http.Error(w, fmt.Sprintf("Folder %s does not exist", filepath.Dir(fullPath)), http.StatusBadRequest)
				return
			}

			f, err := os.Create(fullPath)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, `{"error":"failed to create file: %v"}`, err)
				return
			}
			defer f.Close()
		}

		w.WriteHeader(http.StatusCreated)
		if isFolder {
			fmt.Fprintf(w, "Successfully created folder at path %s", path)
		} else {
			fmt.Fprintf(w, "Successfully created file at path %s", path)
		}
	}
}
