package main

import (
	"context"
	"net/http"
)

type contextKey string

const isAuthenticatedContextKey = contextKey("isAuthenticated")

func (app *application) contextSetUser(r *http.Request, userID int) *http.Request {
	ctx := context.WithValue(r.Context(), isAuthenticatedContextKey, userID)
	return r.WithContext(ctx)
}

func (app *application) contextGetUserID(r *http.Request) int {
	userID, ok := r.Context().Value(isAuthenticatedContextKey).(int)
	if !ok {
		return 0
	}
	return userID
}