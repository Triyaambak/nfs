package routes

import (
	controller "github.com/Triyaambak/nfs/controllers"
	types "github.com/Triyaambak/nfs/types"

	"github.com/go-chi/chi"
)

func SetUpRoutes(router *chi.Mux, serverConfig *types.ServerConfig) {
	c := controller.Controller{}
	(*router).Mount("/", c.FileServer(serverConfig))

	(*router).Get("/ls/*", c.LS(serverConfig))
	(*router).Get("/cat/*", c.Cat(serverConfig))
	(*router).Get("/mv/*", c.MV(serverConfig))
	(*router).Get("/mkdir/*", c.Create(serverConfig, true))
	(*router).Get("/touch/*", c.Create(serverConfig, false))

}
