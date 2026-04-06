package main

import (
	"fmt"
	"groupie-tracker/internal/handler"
	"groupie-tracker/internal/repository"
	"groupie-tracker/internal/service"
	"html/template"
	"net/http"
)

func main() {
	// Create the repository
	repo := &repository.APIRepository{API: "https://groupietrackers.herokuapp.com/api/artists"}
	// Create the service
	svc := service.NewArtistService(repo) // instance of ArtistService
	// Load the artist data
	err := svc.LoadArtists() // Method of ArtistService instance
	if err != nil {
		fmt.Println(err)
		return
	}
	// Parse the templates
	AllArtistsTemplate, err := template.ParseFiles("templates/artists.html")
	if err != nil {
		fmt.Println(err)
		return
	}
	IndivArtistTemplate, err := template.ParseFiles("templates/individual_artist.html")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Create the handler
	hndlr := handler.NewArtistHandler(svc, AllArtistsTemplate, IndivArtistTemplate)

	// Register Routes
	http.HandleFunc("/", hndlr.HandleRoot)
	http.HandleFunc("/artists/", hndlr.HandleArtists)

	// Start the server
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println(err)
		return
	}
}
