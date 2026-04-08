package main

import (
	"context"
	"groupie-tracker/internal/handler"
	"groupie-tracker/internal/middleware"
	"groupie-tracker/internal/repository"
	"groupie-tracker/internal/service"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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

	// Graceful Shutdown
	server := &http.Server{
		Addr:    ":8080",
		Handler: wrapped,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server failed", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	logger.Info("Shutdown signal Ctrl+C received", "signal", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server Shutdown Error!", "Error", err)
	}
	logger.Info("Server stopped gracefully!")
}
