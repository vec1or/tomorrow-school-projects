package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

type ArtistInfo struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

type ArtistRelations struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

type ArtistPage struct {
	Artist    ArtistInfo
	Relations ArtistRelations
}

func Artists(w http.ResponseWriter, r *http.Request) {
	var artistInfo []ArtistInfo
	var artistRelations ArtistRelations
	var artistPage ArtistPage

	// need to fix to fetch only once not every time the handler is called!
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		msg := "Failed to connect to Artists API"
		fmt.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&artistInfo)
	if err != nil {
		msg := "Decoding error"
		fmt.Println(msg, err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	// fmt.Println("First Artist Name:", artistInfo[0].Name)

	tmpl, err := template.ParseFiles("templates/artists.html")
	if err != nil {
		msg := "HTML parsing error!"
		fmt.Println(msg, err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	tmpl2, err := template.ParseFiles("templates/individual_artist.html")
	if err != nil {
		msg := "HTML parsing error!"
		fmt.Println(msg, err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	idString := strings.TrimPrefix(r.URL.Path, "/artists/")
	if idString == "" {
		tmpl.Execute(w, artistInfo)
		return
	}

	idInt, err := strconv.Atoi(idString)
	if err != nil {
		msg := "String to Int conversion error!"
		http.Error(w, msg, http.StatusNotFound)
		fmt.Println(msg, err)
		return
	}

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

			tmpl2.Execute(w, artistPage)
			return
		}
	}
	http.Error(w, "No such ID found!", http.StatusNotFound)
}

func main() {
	http.HandleFunc("/artists/", Artists)

	fmt.Println("Starting server...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Failed to start HTTP server!")
		return
	}
}
