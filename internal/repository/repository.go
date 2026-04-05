package repository

import "groupie-tracker/internal/models"

type ArtistRepository interface {
	FetchAll() ([]models.ArtistInfo, error)
	FetchRelation(RelationsURL string) (models.ArtistRelations, error)
}
