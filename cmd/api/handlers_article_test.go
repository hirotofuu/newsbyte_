package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

func Test_app_articleHandlers(t *testing.T) {
	var tests = []struct {
		name           string
		method         string
		json           string
		paramName      string
		paramID        string
		handler        http.HandlerFunc
		expectedStatus int
	}{
		{"allArticles", "GET", "", "", "", app.GetAllArticles, http.StatusOK},

		// fetch one test
		{"oneArticle", "GET", "", "id", "1", app.GetOneArticle, http.StatusOK},
		{"oneArticle invalid", "GET", "", "id", "100", app.GetOneArticle, http.StatusBadRequest},
		{"oneArticle bad URL param", "GET", "", "id", "Y", app.GetOneArticle, http.StatusBadRequest},

		// fetch user's test
		{"userArticle", "GET", "", "user_id", "1", app.GetUserArticles, http.StatusOK},

		// fetch work test
		{"workArticle", "GET", "", "work", "呪術廻戦", app.GetWorkArticles, http.StatusOK},

		// delete test
		{"deleteArticle", "DELETE", "", "id", "1", app.DeleteArticle, http.StatusOK},
		{"deleteArticle invalid", "DELETE", "", "id", "9", app.DeleteArticle, http.StatusBadRequest},
		{"deleteArticle  bad URL param", "DELETE", "", "id", "Y", app.DeleteArticle, http.StatusBadRequest},

		// insert article test
		{
			"insertArticle valid",
			"PUT",
			`{"title":"a","content":"it's great.","work":"jujukaisen","medium":1,"comment_ok":true,"main_img":"http:sss/sss","user_id":1}`,
			"",
			"",
			app.InsertArticle,
			http.StatusOK,
		},
		{
			"insertArticle invalid input",
			"PUT",
			`{"title":2,"content":"it's great.","work":"jujukaisen","medium":"1","comment_ok":true,"main_img":"http:sss/sss","user_id":1}`,
			"",
			"",
			app.InsertArticle,
			http.StatusBadRequest,
		},
	}

	for _, e := range tests {
		var req *http.Request
		if e.json == "" {
			req, _ = http.NewRequest(e.method, "/", nil)
		} else {
			req, _ = http.NewRequest(e.method, "/", strings.NewReader(e.json))
		}

		if e.paramName != "" {
			chiCtx := chi.NewRouteContext()
			chiCtx.URLParams.Add(e.paramName, e.paramID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(e.handler)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedStatus {
			t.Errorf("%s: wrong status returned; expected %d but got %d", e.name, e.expectedStatus, rr.Code)
		}
	}
}

func Test_app_commentHandlers(t *testing.T) {
	var tests = []struct {
		name           string
		method         string
		json           string
		paramName      string
		paramID        string
		handler        http.HandlerFunc
		expectedStatus int
	}{
		// fetch user's test
		{"userComment", "GET", "", "user_id", "1", app.GetUserComments, http.StatusOK},
		{"userComment invalid user_id", "GET", "", "user_id", "3", app.GetUserComments, http.StatusBadRequest},
		{"userComment invalid paramsName", "GET", "", "user_ide", "3", app.GetUserComments, http.StatusBadRequest},
		
		//fetch article's test
		{"userComment", "GET", "", "article_id", "1", app.GetArticleComments, http.StatusOK},
		{"articleComment invalid article_id", "GET", "", "article_id", "2", app.GetArticleComments, http.StatusBadRequest},
		{"userComment invalid paramsName", "GET", "", "article_ide", "1", app.GetArticleComments, http.StatusBadRequest},

		// delete comment
		{"userComment", "DELETE", "", "id", "1", app.DeleteComment, http.StatusOK},
		{"deleteComment invalid id", "DELETE", "", "id", "2", app.DeleteComment, http.StatusBadRequest},
		{"deleteComment invalid paramName", "DELETE", "", "ide", "2", app.DeleteComment, http.StatusBadRequest},
	}

	for _, e := range tests {
		var req *http.Request
		if e.json == "" {
			req, _ = http.NewRequest(e.method, "/", nil)
		} else {
			req, _ = http.NewRequest(e.method, "/", strings.NewReader(e.json))
		}

		if e.paramName != "" {
			chiCtx := chi.NewRouteContext()
			chiCtx.URLParams.Add(e.paramName, e.paramID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(e.handler)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedStatus {
			t.Errorf("%s: wrong status returned; expected %d but got %d", e.name, e.expectedStatus, rr.Code)
		}
	}
}
