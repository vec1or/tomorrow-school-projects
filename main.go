package main

import (
	"log"
	"net/http"

	"ascii-art-web/handlers"
)

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", logRequest(handlers.IndexHandler))
	mux.HandleFunc("/ascii-art", logRequest(handlers.AsciiArtHandler))

	log.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))

}
func logRequest(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next(w, r)
	}
}
