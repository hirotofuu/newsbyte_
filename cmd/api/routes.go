package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	// register middleware
	mux.Use(middleware.Recoverer)
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "https://*"},
		AllowedMethods:   []string{"GET", "POST", "DELETE", "PUT", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-type", "X-CRSF-Token", "X-Requested-With"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	mux.Route("/user", func(mux chi.Router) {
		mux.Use(app.authRequired)
		mux.Put("/insert_article", app.InsertArticle)
		mux.Delete("/delete_articles/{id}", app.DeleteArticle)

		mux.Put("/insert_follow/{id}", app.InsertFollow)
		mux.Delete("/delete_follow/{id}", app.DeleteFollow)

		mux.Put("/insert_comment", app.InsertComment)
		mux.Delete("/delete_comment/{id}", app.DeleteComment)
		mux.Put("/insert_comment_good/{id}", app.InsertGoodComment)
		mux.Delete("/delete_comment_goos/{id}", app.DeleteGoodComment)

		mux.Put("/insert_article_good/{id}", app.InsertGoodArticle)
		mux.Delete("/delete_article_goos/{id}", app.DeleteGoodArticle)
	})

	mux.Get("/article/{id}", app.GetOneArticle)
	mux.Get("/edit_article/{id}", app.GetEditArticle)
	mux.Get("/good_article/{id}", app.StateGoodArticle)
	mux.Get("/user_articles/{userID}", app.GetUserArticles)
	mux.Get("/user_save_articles/{userID}", app.GetUserSaveArticles)
	mux.Get("/work_articles/{work}", app.GetWorkArticles)
	mux.Get("/articles", app.GetAllArticles)

	mux.Get("/user_comments/{user_id}", app.GetUserComments)
	mux.Get("/article_comments/{article_id}", app.GetArticleComments)

	mux.Get("/following_users/{id}", app.getFollowingUsers)
	mux.Get("/followed_users/{id}", app.getFollowedUsers)
	mux.Get("/all_users", app.getAllUsers)
	mux.Get("/users/{key_word}", app.getSearchUsers)
	mux.Get("/one_user/{id}", app.getOneUser)
	mux.Get("/one_id_name_user/{id_name}", app.getOneIdNameUser)
	mux.Post("/login", app.authenticate)
	mux.Post("/register", app.register)
	mux.Get("/refresh", app.refreshToken)
	mux.Get("/refresh_next", app.refreshToken_next)
	mux.Get("/logout", app.logout)

	return mux
}
