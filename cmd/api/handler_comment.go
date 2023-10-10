package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/hirotofuu/newsbyte/internal/models"
)

func (app *application) GetArticleComments(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "article_id")
	articleID, err := strconv.Atoi(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	comments, err := app.CDB.ArticleComments(articleID, 1)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, comments)
}

func (app *application) GetUserComments(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "user_id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	comments, err := app.CDB.UserComments(userID, 1)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, comments)
}

func (app *application) InsertComment(w http.ResponseWriter, r *http.Request) {
	var comment models.Comment

	err := app.readJSON(w, r, &comment)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()

	_, err = app.CDB.InsertComment(comment)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := JSONResponse{
		Error:   false,
		Message: "comment is inserted",
	}

	app.writeJSON(w, http.StatusOK, resp)
}

func (app *application) DeleteComment(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.CDB.DeleteComment(id)
	if err != nil {
		app.errorJSON(w, err)
	}

	resp := JSONResponse{
		Error:   false,
		Message: "comment deleted",
	}

	app.writeJSON(w, http.StatusOK, resp)

}

func (app *application) InsertGoodComment(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.CDB.InsertGoodComment(id, 1)
	if err != nil {
		app.errorJSON(w, err)
	}

	resp := JSONResponse{
		Error:   false,
		Message: "comment good",
	}

	app.writeJSON(w, http.StatusOK, resp)
}

func (app *application) DeleteGoodComment(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.CDB.DeleteGoodComment(id, 1)
	if err != nil {
		app.errorJSON(w, err)
	}

	resp := JSONResponse{
		Error:   false,
		Message: "delete comment good",
	}

	app.writeJSON(w, http.StatusOK, resp)
}
