package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /", app.home)
	mux.HandleFunc("GET /post/view", app.postView)
	mux.HandleFunc("GET /post/create", app.postCreateForm)
	mux.HandleFunc("POST /post/create", app.postCreate)

	return app.recoverPanic(app.logRequests(secureHeaders(mux)))
}
