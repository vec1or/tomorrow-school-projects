package main

import (
	"errors"
	"fmt"
	"forum/internal/models"
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

	data := &templateData{
		Posts: posts,
	}

	app.render(w, http.StatusOK, "home.tmpl", data)
}

func (app *application) postView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	post, err := app.posts.GetByID(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
			return
		}
		app.serverError(w, err)
		return
	}

	data := &templateData{
		Post: post,
	}

	app.render(w, http.StatusOK, "view.tmpl", data)
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
