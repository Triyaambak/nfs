package controllers

import (
	"net/http"
)

type Controller struct{}

func (c *Controller) FileServer(dir string) http.Handler {
	return http.FileServer(http.Dir(dir))
}
