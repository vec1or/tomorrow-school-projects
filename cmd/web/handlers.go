package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	//"log"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	posts, err := app.posts.Latest()
	if err != nil {
		app.serverError(w, err)
	}

	// for _, post := range posts {
	// 	fmt.Fprintf(w, "%+v\n", post)
	// }

	data := &templateData{
		Posts: posts,
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

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, err)
	}

	//w.Write([]byte("Hello from Snippetbox"))
}

func (app *application) postView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	//w.Write([]byte("Display a specific snippet..."))
	post, err := app.posts.GetByID(id)
	if err != nil {
		app.notFound(w)
		return
	}

	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/partials/nav.tmpl",
		"./ui/html/pages/view.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := &templateData{
		Post: post,
	}

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, err)
	}

	//fmt.Fprintf(w, "ID: %d | Title: %s | Content: %s", post.ID, post.Title, post.Content)
}

func (app *application) postCreate(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		//w.WriteHeader(http.StatusMethodNotAllowed)
		//w.Write([]byte("Method not allowed"))
		return
	}

	userID := 1 // temp
	title := "Test title"
	content := "Test content"

	id, err := app.posts.Insert(userID, title, content)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post/view?id=%d", id), http.StatusSeeOther)
}
