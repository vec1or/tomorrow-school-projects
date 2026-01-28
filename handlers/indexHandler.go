package handlers

import (
	"html/template"
	"net/http"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" { // <-- ключевое
        http.NotFound(w, r)
        return
    }

    if r.Method != http.MethodGet {
        http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        return
    }

    tmpl, err := template.ParseFiles("templates/index.html")
    if err != nil {
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    _ = tmpl.Execute(w, nil)
}