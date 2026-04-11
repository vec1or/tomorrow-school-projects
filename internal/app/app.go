package app

import (
	"context"
	"flag"
	"groupie-tracker/internal/config"
	"groupie-tracker/internal/repository"
	"groupie-tracker/internal/server"
	"groupie-tracker/internal/service"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run() {
	addr := flag.String("addr", ":8080", "HTTP network address")

	flag.Parse()
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
		os.Exit(1)
	}
	// Parse the templates
	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error("failed to create template cache", "error", err)
		os.Exit(1)
	}

	// for name := range templateCache {
	// 	logger.Info("cached template", "name", name)
	// }

	application := &config.Application{
		Logger:        logger,
		Service:       svc,
		TemplateCache: templateCache,
	}

	// Graceful Shutdown
	server := &http.Server{
		Addr:    *addr,
		Handler: server.Routes(application),
	}

	go func() {
		logger.Info("starting server on", "addr", server.Addr)
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
