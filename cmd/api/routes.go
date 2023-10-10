package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	// register middleware
	mux.Use(middleware.Recoverer)
	mux.Use(app.enableCORS)

	mux.Route("/user", func(mux chi.Router) {
		mux.Use(app.authRequired)
		mux.Put("/insert_article", app.InsertArticle)
		mux.Delete("/delete_articles/{id}", app.DeleteArticle)

		mux.Put("/insert_comment", app.InsertComment)
		mux.Delete("/delete_comment/{id}", app.DeleteComment)
		mux.Put("/insert_comment_good/{id}", app.InsertGoodComment)
		mux.Delete("/delete_comment_goos/{id}", app.DeleteGoodComment)

		mux.Put("/insert_article_good/{id}", app.InsertGoodArticle)
		mux.Delete("/delete_article_goos/{id}", app.DeleteGoodArticle)
	})

	mux.Get("/article/{id}", app.GetOneArticle)
	mux.Get("/user_articles/{userID}", app.GetUserArticles)
	mux.Get("/work_articles/{work}", app.GetWorkArticles)
	mux.Get("/articles", app.GetAllArticles)

	mux.Get("user_comments/{user_id}", app.GetUserComments)
	mux.Get("article_comments/{article_id}", app.GetArticleComments)

	mux.Post("/login", app.authenticate)
	mux.Post("/register", app.register)
	mux.Get("/refresh", app.refreshToken)
	mux.Get("/logout", app.logout)

	return mux
}
