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

func Artists(w http.ResponseWriter, r *http.Request) {
	var artistInfo []ArtistInfo

	// need to fix to fetch only once not every time the handler is called!
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		http.Error(w, "Failed to connect to Artists API", http.StatusInternalServerError)
		fmt.Println("Failed to connect to Artists API")
		return
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&artistInfo)
	if err != nil {
		http.Error(w, "Decode error", http.StatusInternalServerError)
		fmt.Println("decode error:", err)
		return
	}

	fmt.Println("First Artist Name:", artistInfo[0].Name)

	tmpl, err := template.ParseFiles("templates/artists.html")
	if err != nil {
		fmt.Println("HTML parsing error!", err)
		http.Error(w, "HTML parsing error!", http.StatusInternalServerError)
		return
	}

	tmpl2, err := template.ParseFiles("templates/individual_artist.html")
	if err != nil {
		fmt.Println("HTML parsing error!", err)
		http.Error(w, "HTML parsing error!", http.StatusInternalServerError)
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
			tmpl2.Execute(w, artistInfo[i])
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
