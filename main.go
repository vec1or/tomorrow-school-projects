package main

import (
	"encoding/json"
	"fmt"
	"groupie-tracker/internal/models"
	"html/template"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
)

// HTML template pointers declaration
var AllArtistsTemplate *template.Template
var IndivArtistTemplate *template.Template

// Artists data declaration
var artistInfo []models.ArtistInfo

func Artists(w http.ResponseWriter, r *http.Request) {
	var artistRelations models.ArtistRelations
	var artistPage models.ArtistPage

	// Route pattern should follow this structure /artists/{ID}
	// meaning if /artists/ is trimmed there should be no / symbol
	RouteString := strings.TrimPrefix(r.URL.Path, "/artists/")

	if RouteString == "" {
		if err := AllArtistsTemplate.Execute(w, artistInfo); err != nil {
			msg := "All Artists Templace Execution Error!"
			http.Error(w, msg, http.StatusInternalServerError)
			fmt.Println(msg, err)
			return
		}
		return
	} else if strings.Contains(RouteString, "/") {
		msg := "Invalid Route"
		http.Error(w, msg, http.StatusNotFound)
		fmt.Println(msg)
		return
	}

	idInt, err := strconv.Atoi(RouteString)
	if err != nil {
		msg := "Artist ID must be a number!"
		http.Error(w, msg, http.StatusBadRequest)
		fmt.Println(msg, err)
		return
	}

	// Fetching per detailed request for a specified artist is OK?
	// Loop and compare cause IDs may be shuffeled irl
	// Looping is fine for this project since the number of Artists is limited
	for i := 0; i < len(artistInfo); i++ {
		if artistInfo[i].ID == idInt {
			respRelation, err := http.Get(artistInfo[i].Relations)
			if err != nil {
				msg := "Failed to connect to Relations API!"
				fmt.Println(msg)
				http.Error(w, msg, http.StatusInternalServerError)
				return
			}
			defer respRelation.Body.Close()

			err = json.NewDecoder(respRelation.Body).Decode(&artistRelations)
			if err != nil {
				msg := "Decoding error"
				fmt.Println(msg, err)
				http.Error(w, msg, http.StatusInternalServerError)
				return
			}

			artistPage.Artist = artistInfo[i]
			artistPage.Relations = artistRelations

			if err := IndivArtistTemplate.Execute(w, artistPage); err != nil {
				msg := "Individual Artist Template Execution Error!"
				http.Error(w, msg, http.StatusInternalServerError)
				fmt.Println(msg, err)
				return
			}
			return
		}
	}
	http.Error(w, "Artist with such ID NOT found!", http.StatusNotFound)
}

func main() {
	// fetch API once here, decode and save, then use later in the handler
	logger := slog.Default()
	ArtistsAPI := "https://groupietrackers.herokuapp.com/api/artists"
	var err error

	resp, err := http.Get(ArtistsAPI)
	if err != nil {
		msg := "Failed to connect to Artists API"
		fmt.Println(msg)
		return
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&artistInfo)
	if err != nil {
		msg := "Decoding error"
		fmt.Println(msg, err)
		return
	}

	// initialize templates here, use later in the handler
	// declaration != initialization
	AllArtistsTemplate, err = template.ParseFiles("templates/artists.html")
	if err != nil {
		msg := "All Artists HTML parsing error!"
		fmt.Println(msg, err)
		return
	}

	IndivArtistTemplate, err = template.ParseFiles("templates/individual_artist.html")
	if err != nil {
		msg := "Individual Artist HTML parsing error!"
		fmt.Println(msg, err)
		return
	}

	// Call the /artists/ handler
	http.HandleFunc("/artists/", Artists)

	// Start the server
	logger.Info("Starting server...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		logger.Error("Failed to connect to Artists API", "error", err, "url", ArtistsAPI)
		return
	}
}
