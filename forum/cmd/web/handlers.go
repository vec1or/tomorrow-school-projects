package main

import (
	"errors"
	"fmt"
	"forum/internal/models"
	"forum/internal/validator"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	categories, err := app.categories.All()
	if err != nil {
		app.serverError(w, err)
		return
	}

	filters := models.PostFilters{}
	selectedCategoryID := 0
	if rawCategory := r.URL.Query().Get("category"); rawCategory != "" {
		selectedCategoryID, err = strconv.Atoi(rawCategory)
		if err != nil || selectedCategoryID < 1 {
			app.clientError(w, http.StatusBadRequest)
			return
		}
		exists, err := app.categories.Exists(selectedCategoryID)
		if err != nil {
			app.serverError(w, err)
			return
		}
		if !exists {
			app.notFound(w)
			return
		}
		filters.CategoryID = selectedCategoryID
	}

	filterName := r.URL.Query().Get("filter")
	userID := app.contextGetUserID(r)
	switch filterName {
	case "":
	case "created", "liked":
		if userID == 0 {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		if filterName == "created" {
			filters.CreatedByUserID = userID
		} else {
			filters.LikedByUserID = userID
		}
	default:
		app.clientError(w, http.StatusBadRequest)
		return
	}

	posts, err := app.posts.Filter(filters)
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Posts = posts
	data.Categories = categories
	data.ActiveFilter = filterName
	data.SelectedCategoryID = selectedCategoryID
	app.render(w, http.StatusOK, "home.tmpl", data)
}

type commentCreateForm struct {
	Content             string
	validator.Validator `form:"-"`
}

func (app *application) postView(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFound(w)
		return
	}
	app.renderPostView(w, r, id, commentCreateForm{}, http.StatusOK)
}

func (app *application) renderPostView(w http.ResponseWriter, r *http.Request, postID int, form commentCreateForm, status int) {
	post, err := app.posts.GetByID(postID)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
			return
		}
		app.serverError(w, err)
		return
	}

	comments, err := app.comments.ByPostID(postID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Post = post
	data.Comments = comments
	data.Form = form
	app.render(w, status, "view.tmpl", data)
}

func (app *application) commentCreate(w http.ResponseWriter, r *http.Request) {
	postID, err := app.readIDParam(r)
	if err != nil {
		app.notFound(w)
		return
	}
	if _, err = app.posts.GetByID(postID); err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	if err = r.ParseForm(); err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	form := commentCreateForm{Content: strings.TrimSpace(r.PostForm.Get("content"))}
	form.CheckField(validator.NotBlank(form.Content), "content", "Comment cannot be blank")
	form.CheckField(validator.MaxChars(form.Content, 2000), "content", "Comment cannot be more than 2000 characters long")
	if !form.Valid() {
		app.renderPostView(w, r, postID, form, http.StatusUnprocessableEntity)
		return
	}

	if _, err = app.comments.Insert(postID, app.contextGetUserID(r), form.Content); err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/post/view/%d#comments", postID), http.StatusSeeOther)
}

func (app *application) postReaction(w http.ResponseWriter, r *http.Request) {
	postID, err := app.readIDParam(r)
	if err != nil {
		app.notFound(w)
		return
	}
	if _, err = app.posts.GetByID(postID); err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	reactionType, ok := app.readReaction(r)
	if !ok {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	if err = app.reactions.ToggleForPost(app.contextGetUserID(r), postID, reactionType); err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/post/view/%d", postID), http.StatusSeeOther)
}

func (app *application) commentReaction(w http.ResponseWriter, r *http.Request) {
	commentID, err := app.readIDParam(r)
	if err != nil {
		app.notFound(w)
		return
	}
	comment, err := app.comments.GetByID(commentID)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	reactionType, ok := app.readReaction(r)
	if !ok {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	if err = app.reactions.ToggleForComment(app.contextGetUserID(r), commentID, reactionType); err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/post/view/%d#comment-%d", comment.PostID, commentID), http.StatusSeeOther)
}

func (app *application) readReaction(r *http.Request) (string, bool) {
	if err := r.ParseForm(); err != nil {
		return "", false
	}
	reactionType := r.PostForm.Get("reaction")
	return reactionType, reactionType == "like" || reactionType == "dislike"
}

type postCreateForm struct {
	Title               string
	Content             string
	CategoryIDs         []int
	validator.Validator `form:"-"`
}

func (app *application) postCreateForm(w http.ResponseWriter, r *http.Request) {
	app.renderPostCreate(w, r, postCreateForm{}, http.StatusOK)
}

func (app *application) renderPostCreate(w http.ResponseWriter, r *http.Request, form postCreateForm, status int) {
	categories, err := app.categories.All()
	if err != nil {
		app.serverError(w, err)
		return
	}
	data := app.newTemplateData(r)
	data.Form = form
	data.Categories = categories
	app.render(w, status, "create.tmpl", data)
}

func (app *application) postCreate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := postCreateForm{
		Title:   strings.TrimSpace(r.PostForm.Get("title")),
		Content: strings.TrimSpace(r.PostForm.Get("content")),
	}
	seen := make(map[int]bool)
	for _, raw := range r.PostForm["category"] {
		id, err := strconv.Atoi(raw)
		if err != nil || id < 1 {
			form.AddFieldError("category", "Select only valid categories")
			continue
		}
		if !seen[id] {
			seen[id] = true
			form.CategoryIDs = append(form.CategoryIDs, id)
		}
	}

	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "Title cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Content, 10000), "content", "Post cannot be more than 10000 characters long")
	form.CheckField(len(form.CategoryIDs) > 0, "category", "Select at least one category")

	allExist, err := app.categories.AllExist(form.CategoryIDs)
	if err != nil {
		app.serverError(w, err)
		return
	}
	form.CheckField(allExist, "category", "Select only valid categories")
	if !form.Valid() {
		app.renderPostCreate(w, r, form, http.StatusUnprocessableEntity)
		return
	}

	id, err := app.posts.Insert(app.contextGetUserID(r), form.Title, form.Content, form.CategoryIDs)
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/post/view/%d", id), http.StatusSeeOther)
}

type userSignupForm struct {
	Username            string `form:"username"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userSignupForm{}
	app.render(w, http.StatusOK, "signup.tmpl", data)
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	form := userSignupForm{
		Username: strings.TrimSpace(r.PostForm.Get("username")),
		Email:    strings.ToLower(strings.TrimSpace(r.PostForm.Get("email"))),
		Password: r.PostForm.Get("password"),
	}
	form.CheckField(validator.NotBlank(form.Username), "username", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Username, 40), "username", "Username cannot be more than 40 characters long")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "Enter a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "Password must be at least 8 characters long")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
		return
	}

	err := app.users.Insert(form.Username, form.Email, form.Password)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrDuplicateEmail):
			form.AddFieldError("email", "Email address is already in use")
		case errors.Is(err, models.ErrDuplicateUsername):
			form.AddFieldError("username", "Username is already in use")
		default:
			app.serverError(w, err)
			return
		}
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
		return
	}
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}
	app.render(w, http.StatusOK, "login.tmpl", data)
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	form := userLoginForm{
		Email:    strings.ToLower(strings.TrimSpace(r.PostForm.Get("email"))),
		Password: r.PostForm.Get("password"),
	}
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "Enter a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "login.tmpl", data)
		return
	}

	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddFieldError("email", "Email or password is incorrect")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "login.tmpl", data)
		} else {
			app.serverError(w, err)
		}
		return
	}
	token, err := app.sessions.Insert(id, 24*time.Hour)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie("session_token"); err == nil {
		_ = app.sessions.Delete(cookie.Value)
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
