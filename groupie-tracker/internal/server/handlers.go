package server

import (
	"net/http"
	"strconv"
	"strings"

	"groupie-tracker/internal/config"
	"groupie-tracker/internal/models"
)

// type ArtistHandler struct {
// 	logger              *slog.Logger
// 	Service             *service.ArtistService
// 	AllArtistsTemplate  *template.Template
// 	IndivArtistTemplate *template.Template
// }

// func NewArtistHandler(
// 	logger *slog.Logger,
// 	Service *service.ArtistService,
// 	AllArtistsTemplate *template.Template,
// 	IndivArtistTemplate *template.Template) *ArtistHandler {
// 	return &ArtistHandler{
// 		logger:              logger,
// 		Service:             Service,
// 		AllArtistsTemplate:  AllArtistsTemplate,
// 		IndivArtistTemplate: IndivArtistTemplate,
// 	}
// }

type Handler struct {
	App *config.Application
}

type mockService struct {
	artists []models.ArtistInfo
	page    models.ArtistPage
	err     error
}

func (m *mockService) GetAll() []models.ArtistInfo {
	return m.artists
}

func (m *mockService) GetArtistPage(id int) (models.ArtistPage, error) {
	return m.page, m.err
}

func (h *Handler) handleArtists(w http.ResponseWriter, r *http.Request) {
	RouteString := strings.TrimPrefix(r.URL.Path, "/artists/")
	if RouteString == "" {
		h.listArtist(w, r)
	} else if strings.Contains(RouteString, "/") {
		h.App.Logger.Warn("Invalid route", "path", r.URL.Path)
		http.Error(w, "Not Found", http.StatusNotFound)
	} else {
		h.showArtists(w, r, RouteString)
	}
}

func (h *Handler) listArtist(w http.ResponseWriter, r *http.Request) {
	artists := h.App.Service.GetAll()
	if err := h.App.TemplateCache["artists.html"].Execute(w, artists); err != nil {
		h.App.Logger.Error("All Artists Template Execution Error!", "error", err)
		http.Error(w, "All Artists Template Execution Error!", http.StatusInternalServerError)
	}
}

func (h *Handler) showArtists(w http.ResponseWriter, r *http.Request, idStr string) {
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		msg := "String to Integer conversion error!"
		h.App.Logger.Warn("Invalid artist ID", "id", idStr, "error", err)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	indivArtist, err := h.App.Service.GetArtistPage(idInt)
	if err != nil {
		h.App.Logger.Warn("Artist Not Found!", "Artist ID", idInt)
		http.Error(w, "Artist with the ID not found!", http.StatusNotFound)
		return
	}
	if err := h.App.TemplateCache["individual_artist.html"].Execute(w, indivArtist); err != nil {
		h.App.Logger.Error("Individual Artist Template Execution failed!", "error", err)
		http.Error(w, "Individual Artist Template Execution failed!", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.Redirect(w, r, "/artists/", http.StatusSeeOther)
	} else {
		h.App.Logger.Warn("Not Found!", "error", http.StatusNotFound)
		http.Error(w, "Not Found!", http.StatusNotFound)
	}
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
