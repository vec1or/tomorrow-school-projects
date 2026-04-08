package main

import (
	"groupie-tracker/internal/handler"
	"groupie-tracker/internal/middleware"
	"groupie-tracker/internal/repository"
	"groupie-tracker/internal/service"
	"html/template"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	// Slog instance
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Create the repository
	repo := &repository.APIRepository{API: "https://groupietrackers.herokuapp.com/api/artists"}

	// Create the service
	svc := service.NewArtistService(logger, repo) // instance of ArtistService

	// Load the artist data
	err := svc.LoadArtists() // Method of ArtistService instance
	if err != nil {
		logger.Error("Failed to load Artists", "error", err)
		return
	}
	// Parse the templates
	AllArtistsTemplate, err := template.ParseFiles("templates/artists.html")
	if err != nil {
		logger.Error("Failed to Parse All Artists Template", "error", err)
		return
	}
	IndivArtistTemplate, err := template.ParseFiles("templates/individual_artist.html")
	if err != nil {
		logger.Error("Failed to Parse Individual Artist Template", "error", err)
		return
	}

	// Create the handler
	hndlr := handler.NewArtistHandler(logger, svc, AllArtistsTemplate, IndivArtistTemplate)

	// Register Routes
	mux := http.NewServeMux()
	mux.HandleFunc("/", hndlr.HandleRoot)
	mux.HandleFunc("/artists/", hndlr.HandleArtists)

	// Middleware
	var wrapped http.Handler = mux
	wrapped = middleware.AllowMethods([]string{http.MethodGet}, wrapped)
	wrapped = middleware.Logging(logger, wrapped)
	wrapped = middleware.Recovery(logger, wrapped)

	// Start the server
	if err := http.ListenAndServe(":8080", wrapped); err != nil {
		logger.Error("Failed to Start the Server!", "error", err)
		return
	}
}
