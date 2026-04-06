package handler

import (
	"groupie-tracker/internal/service"
	"html/template"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
)

type ArtistHandler struct {
	logger              *slog.Logger
	Service             *service.ArtistService
	AllArtistsTemplate  *template.Template
	IndivArtistTemplate *template.Template
}

func NewArtistHandler(
	logger *slog.Logger,
	Service *service.ArtistService,
	AllArtistsTemplate *template.Template,
	IndivArtistTemplate *template.Template) *ArtistHandler {
	return &ArtistHandler{
		logger:              logger,
		Service:             Service,
		AllArtistsTemplate:  AllArtistsTemplate,
		IndivArtistTemplate: IndivArtistTemplate,
	}
}

func (a *ArtistHandler) HandleArtists(w http.ResponseWriter, r *http.Request) {
	RouteString := strings.TrimPrefix(r.URL.Path, "/artists/")
	if RouteString == "" {
		a.listArtist(w, r)
	} else if strings.Contains(RouteString, "/") {
		a.logger.Warn("Invalid route", "path", r.URL.Path)
		http.Error(w, "Not Found", http.StatusNotFound)
	} else {
		a.showArtists(w, r, RouteString)
	}
}

func (a *ArtistHandler) listArtist(w http.ResponseWriter, r *http.Request) {
	artists := a.Service.GetAll()
	if err := a.AllArtistsTemplate.Execute(w, artists); err != nil {
		a.logger.Error("All Artists Template Execution Error!", "error", err)
		http.Error(w, "All Artists Template Execution Error!", http.StatusInternalServerError)
	}
}

func (a *ArtistHandler) showArtists(w http.ResponseWriter, r *http.Request, idStr string) {
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		msg := "String to Integer conversion error!"
		a.logger.Warn("Invalid artist ID", "id", idStr, "error", err)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	indivArtist, err := a.Service.GetArtistPage(idInt)
	if err != nil {
		a.logger.Warn("Artist Not Found!", "Artist ID", idInt)
		http.Error(w, "Artist with the ID not found!", http.StatusNotFound)
		return
	}
	if err := a.IndivArtistTemplate.Execute(w, indivArtist); err != nil {
		a.logger.Error("Individual Artist Template Execution failed!", "error", err)
		http.Error(w, "Individual Artist Template Execution failed!", http.StatusInternalServerError)
		return
	}
}

func (a *ArtistHandler) HandleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.Redirect(w, r, "/artists/", http.StatusSeeOther)
	} else {
		a.logger.Warn("Not Found!", "error", http.StatusNotFound)
		http.Error(w, "Not Found!", http.StatusNotFound)
	}
}
