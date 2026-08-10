package main

import(
	"fmt"
	"net/http"
	"strconv"
	"log"
	"html/template"
)



func home(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	files := []string{
	"./ui/html/base.tmpl",
	"./ui/html/partials/nav.tmpl",
	"./ui/html/pages/home.tmpl",
	}

	// ts, err := template.ParseFiles("./ui/html/pages/home.tmpl")
	// if err != nil {
	// 	log.Print(err.Error())
	// 	http.Error(w, "Internal server error", 500)
	// 	return
	// }
	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal server error", 500)
		return
	}

	// err = ts.Execute(w, nil)
	// if err != nil {
	// 	log.Print(err.Error())
	// 	http.Error(w, "Internal server error", 500)
	// }

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal server error", 500)
	}


	//w.Write([]byte("Hello from Snippetbox"))
}

func snippetView(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	//w.Write([]byte("Display a specific snippet..."))

	fmt.Fprintf(w, "Display a specific snippet with ID: %d", id)
}

func snippetCreate(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		//w.WriteHeader(http.StatusMethodNotAllowed)
		//w.Write([]byte("Method not allowed"))
		return
	}

	w.Write([]byte("Create a new snippet..."))
}