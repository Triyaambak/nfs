package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	routes "github.com/Triyaambak/nfs/routes"
	types "github.com/Triyaambak/nfs/types"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatal("Something went wrong while reading .env file :", err)
	}

	serverConfig := types.ServerConfig{}
	serverConfig.Port = os.Getenv("API_PORT")
	serverConfig.Dir = os.Getenv("NFS_DIR")
	serverConfig.Secret = []byte(os.Getenv("SECRET_KEY"))

	if serverConfig.Port == "" || serverConfig.Dir == "" {
		log.Fatalf("Both port and dir cannot be empty port:%s dir:%s", serverConfig.Port, serverConfig.Dir)
	}
	adr := fmt.Sprintf(":%s", serverConfig.Port)

	router := chi.NewMux()
	routes.SetUpRoutes(router, &serverConfig)

	server := http.Server{
		Addr:    adr,
		Handler: router,
	}

	fmt.Println("Starting server on port " + adr)
	err = server.ListenAndServe()
	if err != nil {
		fmt.Println("Server crashed")
		log.Fatal(err)
	}
}
