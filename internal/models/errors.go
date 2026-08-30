package models

import "errors"

var(
	ErrNoRecord = errors.New("models: no matching record found")
	ErrInvalidCredentials = errors.New("model: invalid credentials")
	ErrDublicateEmail = errors.New("models: dublicate email")
	ErrDublicateUsername = errors.New("models: dublicate username")
)
