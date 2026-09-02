package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /post/view/{id}", app.postView)
	mux.HandleFunc("GET /user/signup", app.userSignup)
	mux.HandleFunc("POST /user/signup", app.userSignupPost)
	mux.HandleFunc("GET /user/login", app.userLogin)
	mux.HandleFunc("POST /user/login", app.userLoginPost)

	protected := http.NewServeMux()
	protected.HandleFunc("GET /post/create", app.postCreateForm)
	protected.HandleFunc("POST /post/create", app.postCreate)
	protected.HandleFunc("POST /user/logout", app.userLogoutPost)

	mux.Handle("/post/", app.requireAuthentication(protected))
	mux.Handle("POST /user/logout", app.requireAuthentication(protected))

	return app.recoverPanic(app.logRequests(secureHeaders(app.authenticate(mux))))
}
