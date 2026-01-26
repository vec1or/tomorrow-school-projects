package handlers

import (
    "html/template"
    "log"
    "net/http"

    "ascii-art-web/ascii"
)
func AsciiArtHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("AsciiArtHandler called")

    if r.Method != http.MethodPost {
        http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        return
    }

    text := r.FormValue("text")
    banner := r.FormValue("banner")

    log.Println("TEXT:", text)
    log.Println("BANNER:", banner)

    if text == "" || banner == "" {
        http.Error(w, "Bad Request", http.StatusBadRequest)
        return
    }

    result, err := ascii.Generate(text, banner)
    if err != nil {
	http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
	return
}

    tmpl, err := template.ParseFiles("templates/index.html")
    if err != nil {
        log.Println("TEMPLATE ERROR:", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    tmpl.Execute(w, map[string]string{
        "Result": result,
    })
}
