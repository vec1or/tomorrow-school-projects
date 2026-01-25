package main

import (
	"log"
	"net/http"

	"ascii-art-web/handlers"
)

func main() {
	
	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.Home)
	mux.HandleFunc("/ascii-art", handlers.AsciiArt)

	log.Println("Server started on http://localhost:8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}