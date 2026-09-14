package server

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"groupie-tracker/internal/assert"
	"groupie-tracker/internal/config"
	"groupie-tracker/internal/models"
)

func TestHandleRoot(t *testing.T) {
	h := &Handler{
		App: &config.Application{
			Logger: slog.Default(),
		},
	}

	tests := []struct {
		name           string
		path           string
		expectedStatus int
	}{
		{"root redirects", "/", http.StatusSeeOther},
		{"unknown path", "/smth", http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			r, err := http.NewRequest(http.MethodGet, tt.path, nil)
			if err != nil {
				t.Fatal(err)
			}

			h.handleRoot(rr, r)

			rs := rr.Result()
			assert.Equal(t, rs.StatusCode, tt.expectedStatus)
		})
	}
}

func TestPing(t *testing.T) {
	rr := httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	ping(rr, r)

	rs := rr.Result()

	assert.Equal(t, rs.StatusCode, http.StatusOK)

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)

	assert.Equal(t, string(body), "OK")
}

func TestShowArtist(t *testing.T) {
	mock := &mockService{
		artists: []models.ArtistInfo{
			{ID: 1, Name: "Queen"},
			{ID: 2, Name: "Beatles"},
		},
		page: models.ArtistPage{},
		err:  nil,
	}

	tmpl := template.Must(template.New("individual_artist.html").Parse("{{.Artist.Name}}"))

	app := &config.Application{
		Logger:  slog.New(slog.NewTextHandler(os.Stdout, nil)),
		Service: mock,
		TemplateCache: map[string]*template.Template{
			"individual_artist.html": tmpl,
		},
	}
	tests := []struct {
		name           string
		idStr          string
		mock           *mockService
		expectedStatus int
	}{
		{
			"ID=abc",
			"abc",
			&mockService{},
			http.StatusBadRequest,
		},
		{
			"ID=999",
			"999",
			&mockService{err: fmt.Errorf("not found")},
			http.StatusNotFound,
		},
		{
			"ID=1",
			"1",
			&mockService{page: models.ArtistPage{}, err: nil},
			http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app.Service = tt.mock
			h := &Handler{App: app}

			rr := httptest.NewRecorder()
			r, _ := http.NewRequest(http.MethodGet, "/artists/"+tt.idStr, nil)

			h.showArtists(rr, r, tt.idStr)
			rs := rr.Result()

			assert.Equal(t, rs.StatusCode, tt.expectedStatus)
		})
	}
}

func TestListArtist(t *testing.T) {
	mock := &mockService{
		artists: []models.ArtistInfo{
			{ID: 1, Name: "Queen"},
			{ID: 2, Name: "Beatles"},
		},
		page: models.ArtistPage{},
		err:  nil,
	}

	tmpl := template.Must(template.New("artists.html").Parse("{{range .}}{{.Name}}{{end}}"))

	app := &config.Application{
		Logger:  slog.New(slog.NewTextHandler(os.Stdout, nil)),
		Service: mock,
		TemplateCache: map[string]*template.Template{
			"artists.html": tmpl,
		},
	}

	h := &Handler{App: app}

	rr := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/artists/", nil)

	h.listArtist(rr, r)
	rs := rr.Result()

	assert.Equal(t, rs.StatusCode, http.StatusOK)
}

func TestHandleArtists(t *testing.T) {
	h := &Handler{
		App: &config.Application{
			Logger: slog.Default(),
		},
	}

	rr := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "/artists/1/extra", nil)
	if err != nil {
		t.Fatal(err)
	}

	h.handleArtists(rr, r)

	rs := rr.Result()
	assert.Equal(t, rs.StatusCode, http.StatusNotFound)
}
