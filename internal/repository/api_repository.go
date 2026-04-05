package repository

import (
	"encoding/json"
	"groupie-tracker/internal/models"
	"net/http"
)

type APIRepository struct {
	API string
}

func (a *APIRepository) FetchAll() ([]models.ArtistInfo, error) {
	var artists []models.ArtistInfo

	resp, err := http.Get(a.API)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&artists)
	if err != nil {
		return nil, err
	}
	return artists, nil
}

func (a *APIRepository) FetchRelation(RelationsURL string) (models.ArtistRelations, error) {
	var artistsRelations models.ArtistRelations

	respRelation, err := http.Get(RelationsURL)
	if err != nil {
		return models.ArtistRelations{}, err
	}
	defer respRelation.Body.Close()

	err = json.NewDecoder(respRelation.Body).Decode(&artistsRelations)
	if err != nil {
		return models.ArtistRelations{}, err
	}
	return artistsRelations, nil
}
