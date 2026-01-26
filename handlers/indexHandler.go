package handlers

import (
    "html/template"
    "log"
    "net/http"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" {
	http.Error(w, "404 Not Found", http.StatusNotFound)
	return
}

    tmpl, err := template.ParseFiles("templates/index.html")
    if err != nil {
        log.Println("TEMPLATE ERROR:", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    tmpl.Execute(w, map[string]string{
        "Result": "",
    })
}
