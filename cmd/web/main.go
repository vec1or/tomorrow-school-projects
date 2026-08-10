package main

import (
	//"fmt"
	"log"
	"net/http"
	//"strconv"
)

func main() {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	log.Print("Starting server on: 8080")
	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}
