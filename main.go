package main

import (
	"log"
	"net/http"

	"ascii-art-web/handlers"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.IndexHandler)
	mux.HandleFunc("/ascii-art", handlers.AsciiArtHandler)

    log.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
