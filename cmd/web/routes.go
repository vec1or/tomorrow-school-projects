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

	mux.Handle("GET /post/create", app.requireAuthentication(http.HandlerFunc(app.postCreateForm)))
	mux.Handle("POST /post/create", app.requireAuthentication(http.HandlerFunc(app.postCreate)))
	mux.Handle("POST /user/logout", app.requireAuthentication(http.HandlerFunc(app.userLogoutPost)))

	return app.recoverPanic(app.logRequests(secureHeaders(app.authenticate(mux))))
}
