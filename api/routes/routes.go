package routes

import (
	middleware "github.com/Triyaambak/nfs/middleware"
	controller "github.com/Triyaambak/nfs/controllers"
	types "github.com/Triyaambak/nfs/types"

	"github.com/go-chi/chi"
)

func SetUpRoutes(router *chi.Mux, serverConfig *types.ServerConfig) {
	c := controller.Controller{}

	(*router).Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddle(serverConfig))
		
		r.Mount("/", c.FileServer(serverConfig))

		r.Get("/ls/*", c.LS(serverConfig))
		r.Get("/cat/*", c.Cat(serverConfig))
		r.Get("/mv/*", c.MV(serverConfig))
		r.Get("/mkdir/*", c.Create(serverConfig, true))
		r.Get("/touch/*", c.Create(serverConfig, false))
	})

}
