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

func (c *Controller) Echo(serverConfig *types.ServerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dir := (*serverConfig).Dir

		serverConfig.MU.Lock()
		defer serverConfig.MU.Unlock()

		urlParam := chi.URLParam(r, "*")
		if err := isParamEmpty(urlParam); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		body, path, isAppend, err := getWriteMode(urlParam)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		fullPath := filepath.Join(dir, path)
		isTaken, err := pathExist(fullPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if !isTaken {
			http.Error(w, fmt.Sprintf("Soure path %s does not exist", path), http.StatusBadRequest)
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

		if body == "" {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "%s file modified", path)
			return
		}

		var f *os.File
		if !isAppend {
			f, err = os.OpenFile(fullPath, os.O_WRONLY|os.O_TRUNC, 0644)
		} else {
			f, err = os.OpenFile(fullPath, os.O_WRONLY|os.O_APPEND, 0644)
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer f.Close()

		_, err = f.Write([]byte(body))
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to write to file in path %s", path), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		if !isAppend {
			fmt.Fprintf(w, "%s file overwritten", path)
		} else {
			fmt.Fprintf(w, "%s file appended", path)
		}
	}
}

func (c *Controller) MV(serverConfig *types.ServerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dir := (*serverConfig).Dir

		serverConfig.MU.Lock()
		defer serverConfig.MU.Unlock()

		path := chi.URLParam(r, "*")
		if err := isParamEmpty(path); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		srcPath, destPath, err := splitPath(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		srcFullPath := filepath.Join(dir, srcPath)
		destFullPath := filepath.Join(dir, destPath)

		isTaken, err := pathExist(srcFullPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if !isTaken {
			http.Error(w, fmt.Sprintf("Soure path %s does not exist", srcPath), http.StatusBadRequest)
			return
		}

		isTaken, err = pathExist(destFullPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if isTaken {
			http.Error(w, fmt.Sprintf("Dest path %s already exist , cannot overwrite", destPath), http.StatusBadRequest)
			return
		}

		isTaken, err = pathExist(filepath.Dir(destFullPath))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if !isTaken {
			http.Error(w, fmt.Sprintf("Parent directory %s for dest path %s does not exist", filepath.Dir(destPath), destPath), http.StatusBadRequest)
			return
		}

		err = os.Rename(srcFullPath, destFullPath)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error moving src path %s to dest path %s: %v", srcPath, destPath, err), http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Moved src path %s to dst path %s", srcPath, destPath)

	}
}

func (c *Controller) LS(serverConfig *types.ServerConfig) http.HandlerFunc {
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
		isTaken, err := pathExist(fullPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
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

			w.WriteHeader(http.StatusOK)

			for _, e := range entries {
				if e.IsDir() {
					fmt.Fprintf(w, "%s/\n", e.Name())
				} else {
					fmt.Fprintf(w, "%s\n", e.Name())
				}
			}
		} else {
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
			return
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

		uid := r.Context().Value("uid").(int)
		gid := r.Context().Value("gid").(int)

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
			return
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

		if err := os.Chown(fullPath, uid, gid); err != nil {
			http.Error(w, fmt.Sprintf("failed to change ownership of %s to uid: %d and gid: %d", path, uid, gid), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		if isFolder {
			fmt.Fprintf(w, "Successfully created folder at path %s", path)
		} else {
			fmt.Fprintf(w, "Successfully created file at path %s", path)
		}
	}
}
