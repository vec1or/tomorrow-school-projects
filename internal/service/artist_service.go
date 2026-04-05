package service

import (
	"fmt"
	"groupie-tracker/internal/models"
	"groupie-tracker/internal/repository"
)

type ArtistService struct {
	repo      repository.ArtistRepository
	artists   []models.ArtistInfo
	artistMap map[int]models.ArtistInfo
}

func NewArtistService(repo repository.ArtistRepository) *ArtistService {
	return &ArtistService{
		repo:      repo,
		artistMap: make(map[int]models.ArtistInfo),
	}
}

func (a *ArtistService) LoadArtists() error {
	var err error
	a.artists, err = a.repo.FetchAll()
	if err != nil {
		return err
	}
	for _, val := range a.artists {
		a.artistMap[val.ID] = val
	}
	return nil
}

func (a *ArtistService) GetAll() []models.ArtistInfo {
	return a.artists
}

func (a *ArtistService) GetArtistPage(ID int) (models.ArtistPage, error) {
	artist, ok := a.artistMap[ID]
	if !ok {
		return models.ArtistPage{}, fmt.Errorf("artist with ID %d not found", ID)
	}
	artistRelations, err := a.repo.FetchRelation(artist.Relations)
	if err != nil {
		return models.ArtistPage{}, fmt.Errorf("failed to fetch relations: %w", err)
	}
	return models.ArtistPage{
		Artist:    artist,
		Relations: artistRelations,
	}, nil
}
