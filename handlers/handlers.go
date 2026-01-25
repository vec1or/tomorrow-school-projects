package handlers

import (
	"html/template"
	"net/http"

	"ascii-art-web/ascii"
)

type PageData struct {
	Result string
	Text   string
	Banner string
	Error  string
}

func Home(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusBadRequest)
		return
	}

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Template not found", http.StatusNotFound)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func AsciiArt(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}()

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	text := r.FormValue("text")
	banner := r.FormValue("banner")

	if text == "" {
		http.Error(w, "Empty text", http.StatusBadRequest)
		return
	}

	if banner == "" {
		http.Error(w, "Banner not selected", http.StatusBadRequest)
		return
	}
	result, err := ascii.Generate(text, banner)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	data := PageData{
		Result: result,
		Text:   text,
		Banner: banner,
	}

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, data)
}
