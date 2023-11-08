package main

import (
	"errors"
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
	id := chi.URLParam(r, "userID")
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

	article, err := app.ADB.OneArticle(articleID, app.isLogin(w, r))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, article)

}

func (app *application) InsertArticle(w http.ResponseWriter, r *http.Request) {
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

	app.writeJSON(w, http.StatusOK, resp)
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

func (app *application) InsertGoodArticle(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	yourID := app.isLogin(w, r)
	if yourID == 0 {
		app.errorJSON(w, errors.New("you are not authenticated"), http.StatusUnauthorized)
		return
	}

	err = app.ADB.InsertGoodArticle(id, yourID)
	if err != nil {
		app.errorJSON(w, err)
	}

	resp := JSONResponse{
		Error:   false,
		Message: "article good",
	}

	app.writeJSON(w, http.StatusOK, resp)
}

func (app *application) DeleteGoodArticle(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	yourID := app.isLogin(w, r)
	if yourID == 0 {
		app.errorJSON(w, errors.New("you are not authenticated"), http.StatusUnauthorized)
		return
	}

	err = app.ADB.DeleteGoodArticle(id, yourID)
	if err != nil {
		app.errorJSON(w, err)
	}

	resp := JSONResponse{
		Error:   false,
		Message: "delete article good",
	}

	app.writeJSON(w, http.StatusOK, resp)
}
