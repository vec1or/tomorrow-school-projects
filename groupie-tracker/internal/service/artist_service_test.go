package service

import (
	"fmt"
	"groupie-tracker/internal/models"
	"log/slog"
	"os"
	"testing"
)

type mockRepository struct {
	artists   []models.ArtistInfo
	relations models.ArtistRelations
	err       error
}

func (m *mockRepository) FetchAll() ([]models.ArtistInfo, error) {
	return m.artists, m.err
}
func (m *mockRepository) FetchRelation(RelationsURL string) (models.ArtistRelations, error) {
	return m.relations, m.err
}

func TestLoadArtists(t *testing.T) {
	mock := &mockRepository{
		artists: []models.ArtistInfo{
			{ID: 1, Name: "Queen"},
			{ID: 2, Name: "Beatles"},
		},
		err: nil,
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	svc := NewArtistService(logger, mock)

	err := svc.LoadArtists()

	if err != nil {
		t.Fatal(err)
	}

	artists := svc.GetAll()
	if len(artists) != 2 {
		t.Errorf("expected 2 artists, got %d", len(artists))
	}
}

func TestGetArtistPage(t *testing.T) {
	mock := &mockRepository{
		artists: []models.ArtistInfo{
			{ID: 1, Name: "Queen"},
		},
		relations: models.ArtistRelations{
			ID: 1,
			DatesLocations: map[string][]string{
				"dunedin-new_zealand": {"10-02-2020"},
				"georgia-usa":         {"22-08-2019"},
			},
		},
		err: nil,
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	svc := NewArtistService(logger, mock)

	err := svc.LoadArtists()
	if err != nil {
		t.Fatal(err)
	}

	page, err := svc.GetArtistPage(1)

	if err != nil {
		t.Fatal(err)
	}

	if page.Artist.Name != "Queen" {
		t.Errorf("Expected Queen, got %s", page.Artist.Name)
	}

	page, err = svc.GetArtistPage(999)

	if err == nil {
		t.Error("expected error for non-existent ID, got nil")
	}

	mock.err = fmt.Errorf("connection failed")

	page, err = svc.GetArtistPage(1)

	if err == nil {
		t.Error("expected error when fetching fails, got nil")
	}
}
