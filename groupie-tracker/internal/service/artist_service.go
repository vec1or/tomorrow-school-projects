package service

import (
	"fmt"
	"groupie-tracker/internal/models"
	"groupie-tracker/internal/repository"
	"log/slog"
)

type ArtistService struct {
	logger    *slog.Logger
	repo      repository.ArtistRepository
	artists   []models.ArtistInfo
	artistMap map[int]models.ArtistInfo
}

func NewArtistService(logger *slog.Logger, repo repository.ArtistRepository) *ArtistService {
	return &ArtistService{
		logger:    logger,
		repo:      repo,
		artistMap: make(map[int]models.ArtistInfo),
	}
}

func (a *ArtistService) LoadArtists() error {
	var err error
	a.artists, err = a.repo.FetchAll()
	if err != nil {
		a.logger.Error("Failed to Fetch All Artist Data", "error", err)
		return err
	}
	for _, val := range a.artists {
		a.artistMap[val.ID] = val
	}
	a.logger.Warn("All artists loaded", "count", len(a.artists))
	return nil
}

func (a *ArtistService) GetAll() []models.ArtistInfo {
	return a.artists
}

func (a *ArtistService) GetArtistPage(ID int) (models.ArtistPage, error) {
	artist, ok := a.artistMap[ID]
	if !ok {
		a.logger.Error("Not Found", "Artist ID", ID)
		return models.ArtistPage{}, fmt.Errorf("artist with ID %d not found", ID)
	}
	artistRelations, err := a.repo.FetchRelation(artist.Relations)
	if err != nil {
		a.logger.Error("Failed to fetch Relations Data!", "error", err)
		return models.ArtistPage{}, fmt.Errorf("failed to fetch relations: %w", err)
	}
	a.logger.Info("The artist loaded", "Artist ID", ID)
	return models.ArtistPage{
		Artist:    artist,
		Relations: artistRelations,
	}, nil
}
