package server

import (
	"groupie-tracker/internal/config"
	"groupie-tracker/internal/middleware"
	"net/http"
)

func Routes(app *config.Application) http.Handler {
	h := &Handler{App: app}
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", h.handleRoot)
	mux.HandleFunc("/artists/", h.handleArtists)

	return middleware.SecureHeaders(
		middleware.RecoverPanic(app,
			middleware.LogRequest(app,
				middleware.AllowMethods(mux))))
}
