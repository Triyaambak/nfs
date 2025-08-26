package main

import (
	"log"
	"net/http"
)

func main() {
	dir := ".././nfs"

	fs := http.FileServer(http.Dir(dir))

	http.Handle("/", fs)

	log.Println("Serving files from", dir, "on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
