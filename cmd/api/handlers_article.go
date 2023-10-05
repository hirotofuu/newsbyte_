package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/hirotofuu/newsbyte/internal/models"
)

func (app *application) GetAllArticles(w http.ResponseWriter, r *http.Request) {
	articles, err := app.ADB.AllArticles()
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, articles)
}

func (app *application) GetUserArticles(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "user_id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	articles, err := app.ADB.UserArticles(userID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, articles)
}

func (app *application) GetWorkArticles(w http.ResponseWriter, r *http.Request) {
	work := chi.URLParam(r, "work")

	articles, err := app.ADB.WorkArticles(work)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, articles)
}

func (app *application) GetOneArticle(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	articleID, err := strconv.Atoi(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	article, err := app.ADB.OneArticle(articleID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, article)

}

func (app *application) InsertMovie(w http.ResponseWriter, r *http.Request) {
	var article models.Article

	err := app.readJSON(w, r, &article)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	article.CreatedAt = time.Now()
	article.UpdatedAt = time.Now()

	_, err = app.ADB.InsertArticle(article)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := JSONResponse{
		Error:   false,
		Message: "article is inserted",
	}

	app.writeJSON(w, http.StatusAccepted, resp)
}

func (app *application) DeleteArticle(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.ADB.DeleteArticle(id)
	if err != nil {
		app.errorJSON(w, err)
	}

	resp := JSONResponse{
		Error:   false,
		Message: "article deleted",
	}

	app.writeJSON(w, http.StatusOK, resp)

}
