package routes

import (
	controller "github.com/Triyaambak/nfs/controllers"
	types "github.com/Triyaambak/nfs/types"

	"github.com/go-chi/chi"
)

func SetUpRoutes(router *chi.Mux, serverConfig *types.ServerConfig) {
	c := controller.Controller{}
	(*router).Mount("/", c.FileServer(serverConfig))

	(*router).Get("/cat/*", c.Fetch(serverConfig))

	(*router).Post("/mkdir/*", c.Create(serverConfig, false))
	(*router).Post("/touch/*", c.Create(serverConfig, true))

}
