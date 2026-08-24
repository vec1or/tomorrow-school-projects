package main

import "forum/internal/models"

type templateData struct {
	Post  *models.Post
	Posts []*models.Post
}
