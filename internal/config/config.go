package config

import (
	"html/template"
	"log/slog"

	"groupie-tracker/internal/models"
)

type ArtistService interface {
	GetAll() []models.ArtistInfo
	GetArtistPage(id int) (models.ArtistPage, error)
}

type Application struct {
	Logger        *slog.Logger
	Service       ArtistService
	TemplateCache map[string]*template.Template
}
