package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Triyaambak/nfs/types"
	"github.com/go-chi/chi"
)

type Controller struct{}

type Path struct {
	Dir string `json:"dir"`
}

func (c *Controller) FileServer(serverConfig *types.ServerConfig) http.Handler {
	dir := (*serverConfig).Dir
	return http.StripPrefix("/", http.FileServer(http.Dir(dir)))
}

func (c *Controller) GetFile(serverConfig *types.ServerConfig) http.HandlerFunc {
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

func (c *Controller) CreateFile(serverConfig *types.ServerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dir := (*serverConfig).Dir

		serverConfig.MU.Lock()
		defer serverConfig.MU.Unlock()

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
