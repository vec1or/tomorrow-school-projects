package handlers

import (
	"os"
	"html/template"
	"net/http"

	"ascii-art-web/ascii"
)

func AsciiArtHandler(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/ascii-art" {
        http.NotFound(w, r)
        return
    }

    if r.Method != http.MethodPost {
        http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        return
    }

    if err := r.ParseForm(); err != nil {
        http.Error(w, "Bad Request", http.StatusBadRequest)
        return
    }

	text := r.FormValue("text")
	banner := r.FormValue("banner")

	validBanners := map[string]bool{
    	"standard":   true,
    	"shadow":     true,
    	"thinkertoy": true,
	}	

	if text == "" || banner == "" {
    	http.Error(w, "Bad Request", http.StatusBadRequest)
    	return
	}

	if !validBanners[banner] {
    	http.Error(w, "Bad Request", http.StatusBadRequest)
    	return
	}

	result, err := ascii.Generate(text, banner)
	if err != nil {
    	if os.IsNotExist(err) {
        	http.Error(w, "Banner Not Found", http.StatusNotFound)
        	return
    	}
    	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    	return
	}

    tmpl, err := template.ParseFiles("templates/index.html")
    if err != nil {
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    data := struct {
        Input  string
        Banner string
        Result template.HTML
    }{
        Input:  text,
        Banner: banner,
        Result: template.HTML(result),
    }

    _ = tmpl.Execute(w, data)
}

