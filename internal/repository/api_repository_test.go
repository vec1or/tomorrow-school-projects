package repository

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchAll(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[{"id":1, "name":"Queen"}, {"id":2, "name":"Beatles"}]`))
	}))
	defer server.Close()

	repo := &APIRepository{API: server.URL}
	artists, err := repo.FetchAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(artists) != 2 {
		t.Errorf("expected 2 artist, got %d", len(artists))
	}
	if artists[0].Name != "Queen" {
		t.Errorf("expected Queen, got %s", artists[0].Name)
	}
}
func TestFetchRelation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"id":1,"datesLocations":{   "dunedin-new_zealand": ["10-02-2020"],"georgia-usa":["22-08-2019"]}}`))
	}))
	defer server.Close()

	repo := &APIRepository{API: server.URL}
	relations, err := repo.FetchRelation(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	if len(relations.DatesLocations) != 2 {
		t.Errorf("expected 2 artist, got %d", len(relations.DatesLocations))
	}
	if relations.DatesLocations["georgia-usa"][0] != "22-08-2019" {
		t.Errorf("expected Queen, got %s", relations.DatesLocations["georgia-usa"][0])
	}

}
