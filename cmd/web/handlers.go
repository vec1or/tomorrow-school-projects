package main

import(
	"fmt"
	"net/http"
	"strconv"
	//"log"
	"html/template"
)



func (app *application) home(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		app.notFound(w)
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
		app.serverError(w, err)
		return
	}

	// err = ts.Execute(w, nil)
	// if err != nil {
	// 	log.Print(err.Error())
	// 	http.Error(w, "Internal server error", 500)
	// }

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.serverError(w, err)
	}


	//w.Write([]byte("Hello from Snippetbox"))
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	//w.Write([]byte("Display a specific snippet..."))

	fmt.Fprintf(w, "Display a specific snippet with ID: %d", id)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		//w.WriteHeader(http.StatusMethodNotAllowed)
		//w.Write([]byte("Method not allowed"))
		return
	}

	w.Write([]byte("Create a new snippet..."))
}